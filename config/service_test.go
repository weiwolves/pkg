// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config_test

import (
	"bytes"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log/logw"
	"github.com/weiwolves/pkg/config"
	"github.com/weiwolves/pkg/config/storage"
	"github.com/weiwolves/pkg/store/scope"
	"github.com/weiwolves/pkg/sync/bgwork"
	"github.com/weiwolves/pkg/util/assert"
	"github.com/weiwolves/pkg/util/conv"
	"github.com/fortytw2/leaktest"
)

var (
	_ config.Setter             = (*config.Service)(nil)
	_ config.Subscriber         = (*config.Service)(nil)
	_ config.ObserverRegisterer = (*config.Service)(nil)
)

func TestMustNewService_ShouldPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				assert.True(t, errors.VerificationFailed.Match(err), "%+v", err)
			} else {
				t.Errorf("Panic should contain an error but got:\n%+v", r)
			}
		} else {
			t.Error("Expecting a panic but got nothing")
		}
	}()

	srv := config.MustNewService(storage.NewMap(), config.Options{
		EnablePubSub: false,
	}, config.MakeLoadDataOption(func(_ *config.Service) error {
		return errors.VerificationFailed.Newf("Ups")
	}))

	assert.Nil(t, srv)
}

func TestNotKeyNotFoundError(t *testing.T) {
	srv := config.MustNewService(storage.NewMap(), config.Options{})
	assert.NotNil(t, srv)

	scopedSrv := srv.Scoped(1, 1)

	flat := scopedSrv.Get(scope.Default, "catalog/product/enable_flat")
	assert.False(t, flat.IsValid(), "Should not find the key")
	assert.True(t, flat.IsEmpty(), "should be empty")

	val := scopedSrv.Get(scope.Store, "catalog")
	assert.False(t, val.IsValid())
	assert.True(t, val.IsEmpty(), "should be empty")
}

func TestService_Put(t *testing.T) {
	srv := config.MustNewService(storage.NewMap(), config.Options{})
	assert.NotNil(t, srv)

	var p1 config.Path
	err := srv.Set(p1, []byte{})
	assert.True(t, errors.Empty.Match(err), "Error: %s", err)
}

func TestService_Write_Get_Value_Success(t *testing.T) {
	runner := func(p config.Path, value []byte) func(*testing.T) {
		return func(t *testing.T) {
			srv := config.MustNewService(storage.NewMap(), config.Options{})
			assert.NoError(t, srv.Set(p, value), "Writing Value in Test %q should not fail", t.Name())

			haveStr, ok, haveErr := srv.Get(p).Str()
			assert.NoError(t, haveErr, "No error should occur when retrieving a value")
			assert.True(t, ok)
			assert.Exactly(t, string(value), haveStr)
		}
	}

	basePath := config.MustMakePath("aa/bb/cc")

	t.Run("stringDefault", runner(basePath, []byte("Gopher")))
	t.Run("stringWebsite", runner(basePath.BindWebsite(10), []byte("Gopher")))
	t.Run("stringStore", runner(basePath.BindStore(22), []byte("Gopher")))
}

func TestScoped_IsValid(t *testing.T) {
	t.Parallel()
	cfg := config.NewFakeService(storage.NewMap())
	tests := []struct {
		websiteID uint32
		storeID   uint32
		want      bool
	}{
		{0, 0, true},
		{1, 0, true},
		{1, 1, true},
		{0, 1, false},
		{1, 0, true},
		{1, 1, true},
		{0, 1, false},
	}
	for i, test := range tests {
		s := cfg.Scoped(test.websiteID, test.storeID)
		if have, want := s.IsValid(), test.want; have != want {
			t.Errorf("Idx %d => Have: %v Want: %v", i, have, want)
		}
	}
}

func TestScopedServiceScope(t *testing.T) {
	t.Parallel()
	tests := []struct {
		websiteID, storeID uint32
		wantScope          scope.Type
		wantID             uint32
	}{
		{0, 0, scope.Default, 0},
		{1, 0, scope.Website, 1},
		{1, 3, scope.Store, 3},
		{0, 3, scope.Store, 3},
	}
	for i, test := range tests {
		sg := config.NewFakeService(storage.NewMap()).Scoped(test.websiteID, test.storeID)
		haveScope, haveID := sg.ScopeID().Unpack()
		assert.Exactly(t, test.wantScope, haveScope, "Index %d", i)
		assert.Exactly(t, test.wantID, haveID, "Index %d", i)
	}
}

func TestScopedServicePath(t *testing.T) {
	t.Parallel()
	basePath := config.MustMakePath("aa/bb/cc")
	tests := []struct {
		desc               string
		fqpath             string
		route              string
		perm               scope.Type
		websiteID, storeID uint32
		wantErrKind        errors.Kind
		wantTypeIDs        scope.TypeIDs
	}{
		{
			"Default ScopedGetter should return default scope",
			basePath.String(), "aa/bb/cc", scope.Absent, 0, 0, errors.NoKind, scope.TypeIDs{scope.DefaultTypeID},
		},
		{
			"Website ID 1 ScopedGetter should fall back to default scope",
			basePath.String(), "aa/bb/cc", scope.Website, 1, 0, errors.NoKind, scope.TypeIDs{scope.DefaultTypeID, scope.Website.WithID(1)},
		},
		{
			"Website ID 10 ScopedGetter should fall back to website 10 scope",
			basePath.BindWebsite(10).String(), "aa/bb/cc", scope.Website, 10, 0, errors.NoKind, scope.TypeIDs{scope.Website.WithID(10)},
		},
		{
			"Website ID 10 + Store 22 ScopedGetter should fall back to website 10 scope",
			basePath.BindWebsite(10).String(), "aa/bb/cc", scope.Store, 10, 22, errors.NoKind, scope.TypeIDs{scope.Website.WithID(10), scope.Store.WithID(22)},
		},
		{
			"Website ID 10 + Store 22 ScopedGetter should return Store 22 scope",
			basePath.BindStore(22).String(), "aa/bb/cc", scope.Store, 10, 22, errors.NoKind, scope.TypeIDs{scope.Store.WithID(22)},
		},
		{
			"Website ID 10 + Store 42 ScopedGetter should return nothing",
			basePath.BindStore(22).String(), "aa/bb/cc", scope.Store, 10, 42, errors.NotFound, scope.TypeIDs{scope.DefaultTypeID},
		},
		{
			"Path consists of only two elements which is incorrect",
			basePath.String(), "aa/bb", scope.Store, 0, 0, errors.NotValid, nil,
		},
	}

	// vals stores all possible types for which we have functions in config.ScopedGetter
	vals := []interface{}{"Gopher", true, float64(3.14159), int(2016), time.Now(), []byte(`Hellö Dear Goph€rs`), time.Hour}
	const customTestTimeFormat = "2006-01-02 15:04:05"

	for vi, wantVal := range vals {
		for xi, test := range tests {

			fqVal := conv.ToString(wantVal)
			if wv, ok := wantVal.(time.Time); ok {
				fqVal = wv.Format(customTestTimeFormat)
			}
			cg := config.NewFakeService(storage.NewMap(test.fqpath, fqVal))

			sg := cg.Scoped(test.websiteID, test.storeID)
			haveVal := sg.Get(test.perm, test.route)

			if !test.wantErrKind.Empty() {
				assert.False(t, haveVal.IsValid(), "Index %d/%d scoped path value must be found ", vi, xi)
				continue
			}

			switch wv := wantVal.(type) {
			case []byte:
				var buf strings.Builder
				_, err := haveVal.WriteTo(&buf)
				assert.NoError(t, err, "Error: %+v\n\n%s", err, test.desc)
				assert.Exactly(t, string(wv), buf.String(), test.desc)

			case string:
				hs, _, err := haveVal.Str()
				assert.NoError(t, err)
				assert.Exactly(t, wv, hs)
			case bool:
				hs, _, err := haveVal.Bool()
				assert.NoError(t, err)
				assert.Exactly(t, wv, hs)
			case float64:
				hs, _, err := haveVal.Float64()
				assert.NoError(t, err)
				assert.Exactly(t, wv, hs)
			case int:
				hs, _, err := haveVal.Int()
				assert.NoError(t, err)
				assert.Exactly(t, wv, hs)
			case time.Time:
				hs, _, err := haveVal.Time()
				assert.NoError(t, err)
				assert.Exactly(t, wv.Format(customTestTimeFormat), hs.Format(customTestTimeFormat))
			case time.Duration:
				hs, _, err := haveVal.Duration()
				assert.NoError(t, err)
				assert.Exactly(t, wv, hs)

			default:
				t.Fatalf("Unsupported type: %#v in vals index %d", wantVal, vi)
			}

			assert.Exactly(t, test.wantTypeIDs, cg.AllInvocations().ScopeIDs(), "Index %d/%d", vi, xi)
		}
	}
}

func TestScopedServicePermission_All(t *testing.T) {
	t.Parallel()
	basePath := config.MustMakePath("aa/bb/cc")
	mapStorage := storage.NewMap(
		basePath.Bind(scope.DefaultTypeID).String(), "a",
		basePath.BindWebsite(1).String(), "b",
		basePath.BindStore(1).String(), "c",
	)

	tests := []struct {
		s       scope.Type
		want    string
		wantIDs scope.TypeIDs
	}{
		{scope.Default, "a", scope.TypeIDs{scope.DefaultTypeID}},
		{scope.Website, "b", scope.TypeIDs{scope.Website.WithID(1)}},
		{scope.Group, "a", scope.TypeIDs{scope.DefaultTypeID}},
		{scope.Store, "c", scope.TypeIDs{scope.Store.WithID(1)}},
		{scope.Absent, "c", scope.TypeIDs{scope.Store.WithID(1)}}, // because ScopedGetter bound to store scope
	}
	for i, test := range tests {
		srv := config.NewFakeService(mapStorage)

		haveVal := srv.Scoped(1, 1).Get(test.s, "aa/bb/cc")
		s, ok, err := haveVal.Str()
		if err != nil {
			t.Fatal("Index", i, "Error", err)
		}
		assert.True(t, ok, "Scoped path value must be found")
		assert.Exactly(t, test.want, s, "Index %d", i)
		assert.Exactly(t, test.wantIDs, srv.Invokes().ScopeIDs(), "Index %d", i)
	}
	// assert.Exactly(t, []string{"default/0/aa/bb/cc", "stores/1/aa/bb/cc", "websites/1/aa/bb/cc"}, sm.Invokes().Paths())
}

func TestScopedServicePermission_One(t *testing.T) {
	t.Parallel()
	basePath1 := config.MustMakePath("aa/bb/cc")
	basePath2 := config.MustMakePath("dd/ee/ff")
	basePath3 := config.MustMakePath("dd/ee/gg")

	const WebsiteID = 3
	const StoreID = 5

	sm := config.NewFakeService(storage.NewMap(
		basePath1.BindDefault().String(), "a",
		basePath1.BindWebsite(WebsiteID).String(), "b",
		basePath1.BindStore(StoreID).String(), "c",
		basePath2.BindWebsite(WebsiteID).String(), "bb2",
		basePath3.String(), "cc3",
	))

	t.Run("query1 by scope.Default, matches default", func(t *testing.T) {
		s, ok, err := sm.Scoped(WebsiteID, StoreID).Get(scope.Default, "aa/bb/cc").Str()
		assert.True(t, ok, "scoped path value must be found")
		assert.NoError(t, err)
		assert.Exactly(t, "a", s) // because ScopedGetter bound to store scope
	})

	t.Run("query1 by scope.Website, matches website", func(t *testing.T) {
		s, ok, err := sm.Scoped(WebsiteID, StoreID).Get(scope.Website, "aa/bb/cc").Str()
		assert.True(t, ok, "scoped path value must be found")
		assert.NoError(t, err)
		assert.Exactly(t, "b", s) // because ScopedGetter bound to store scope
	})

	t.Run("query1 by scope.Store, matches store", func(t *testing.T) {
		s, ok, err := sm.Scoped(WebsiteID, StoreID).Get(scope.Store, "aa/bb/cc").Str()
		assert.True(t, ok, "scoped path value must be found")
		assert.NoError(t, err)
		assert.Exactly(t, "c", s) // because ScopedGetter bound to store scope
	})
	t.Run("query1 by scope.Absent, matches store", func(t *testing.T) {
		s, ok, err := sm.Scoped(WebsiteID, StoreID).Get(scope.Absent, "aa/bb/cc").Str()
		assert.True(t, ok, "scoped path value must be found")
		assert.NoError(t, err)
		assert.Exactly(t, "c", s) // because ScopedGetter bound to store scope
	})

	t.Run("query2 by scope.Store, fallback to website", func(t *testing.T) {
		s, ok, err := sm.Scoped(WebsiteID, StoreID).Get(scope.Store, "dd/ee/ff").Str()
		assert.True(t, ok, "scoped path value must be found")
		assert.NoError(t, err)
		assert.Exactly(t, "bb2", s) // because ScopedGetter bound to store scope
	})

	t.Run("query3 by scope.Store, fallback to default", func(t *testing.T) {
		s, ok, err := sm.Scoped(WebsiteID, StoreID).Get(scope.Store, "dd/ee/gg").Str()
		assert.True(t, ok, "scoped path value must be found")
		assert.NoError(t, err)
		assert.Exactly(t, "cc3", s) // because ScopedGetter bound to store scope
	})
	t.Run("query3 by scope.Website, fallback to default", func(t *testing.T) {
		s, ok, err := sm.Scoped(WebsiteID, StoreID).Get(scope.Website, "dd/ee/gg").Str()
		assert.True(t, ok, "scoped path value must be found")
		assert.NoError(t, err)
		assert.Exactly(t, "cc3", s) // because ScopedGetter bound to store scope
	})
}

func TestWithLRU(t *testing.T) {
	t.Parallel()
	mustFloat := func(val float64, ok bool, err error) float64 {
		if err != nil {
			t.Fatal(err)
		}
		return val
	}

	var buf bytes.Buffer
	l := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithWriter(&buf),
	)

	srv := config.MustNewService(storage.NewMap(), config.Options{
		Level1: storage.NewLRU(5),
		Log:    l,
	})

	p1 := config.MustMakePath("carrier/dhl/enabled")
	p2 := p1.BindWebsite(2)
	p3 := p1.BindStore(3)

	val := srv.Get(p1)
	assert.False(t, val.IsValid(), "value should NOT be valid and not found")

	assert.NoError(t, srv.Set(p1, []byte(`1.001`)))
	assert.NoError(t, srv.Set(p2, []byte(`2.002`)))
	assert.NoError(t, srv.Set(p3, []byte(`3.003`)))

	// NOT from LRU
	assert.Exactly(t, 1.001, mustFloat(srv.Get(p1).Float64()), "Path1: %s", p1.String())
	assert.Exactly(t, 2.002, mustFloat(srv.Get(p2).Float64()), "Path2: %s", p2.String())
	assert.Exactly(t, 3.003, mustFloat(srv.Get(p3).Float64()), "Path3: %s", p3.String())

	// Now from LRU
	assert.Exactly(t, 1.001, mustFloat(srv.Get(p1).Float64()), "Path1: %s", p1.String())
	assert.Exactly(t, 2.002, mustFloat(srv.Get(p2).Float64()), "Path2: %s", p2.String())
	assert.Exactly(t, 3.003, mustFloat(srv.Get(p3).Float64()), "Path3: %s", p3.String())

	lStr := buf.String()

	tests := []struct {
		contains string
		strCount int
	}{
		{`"default/0/carrier/dhl/enabled" found: "NO"`, 1},
		{`"default/0/carrier/dhl/enabled" data_length: 5`, 1},
		{`"websites/2/carrier/dhl/enabled" data_length: 5`, 1},
		{`"stores/3/carrier/dhl/enabled" data_length: 5`, 1},
		{`found: "Level2"`, 3},
		{`found: "Level1"`, 3},
		{`"websites/2/carrier/dhl/enabled" found: "Level1"`, 1},
	}
	for _, test := range tests {
		assert.Contains(t, lStr, test.contains)
		assert.Exactly(t, test.strCount, strings.Count(lStr, test.contains), "%s", test.contains)
	}
}

func TestService_Scoped_LRU_Parallel(t *testing.T) {
	srv := config.MustNewService(storage.NewMap(), config.Options{
		Level1: storage.NewLRU(5),
	})

	const route1 = "carrier/dhl/enabled"
	const route2 = "payment/paypal/active"
	p1 := config.MustMakePath(route1)
	p2 := config.MustMakePath(route2)
	paths := config.PathSlice{
		p1,
		p1.BindWebsite(2),
		p1.BindStore(3),
		p2,
		p2.BindWebsite(2),
		p2.BindStore(3),
	}

	scpd := srv.Scoped(2, 3)

	bgwork.Wait(len(paths), func(idx int) {
		true := []byte(`1`)
		p := paths[idx]
		if p.RouteHasPrefix(route2) {
			true = []byte(`0`)
		}
		if err := srv.Set(p, true); err != nil {
			panic(err)
		}

		time.Sleep(time.Millisecond)

		v1, ok, err := scpd.Get(scope.Website, route1).Bool()
		if !ok {
			panic("route1 Value must be found")
		}
		if err != nil {
			panic(err)
		}
		if !v1 {
			panic("route1 Value must be true")
		}

		v2, ok, err := scpd.Get(scope.Website, route2).Bool()
		if !ok {
			panic("route2 Value must be found")
		}
		if err != nil {
			panic(err)
		}
		if v2 {
			panic("route2 Value must be false")
		}
	})
}

func TestService_EnvName_From_OSEnv(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer leaktest.Check(t)()

		assert.NoError(t, os.Setenv(config.DefaultOSEnvVariableName, "STAGING"))
		defer func() { assert.NoError(t, os.Unsetenv(config.DefaultOSEnvVariableName)) }()

		srv := config.MustNewService(storage.NewMap(), config.Options{})
		defer func() { assert.NoError(t, srv.Close()) }()

		assert.Exactly(t, "STAGING", srv.EnvName())
		assert.Exactly(t, "__STAGING__", srv.ReplaceEnvName("__"+config.EnvNamePlaceHolder+"__"))
	})

	t.Run("env does contain non-letters", func(t *testing.T) {
		assert.NoError(t, os.Setenv(config.DefaultOSEnvVariableName, "STAGING"))
		defer func() { assert.NoError(t, os.Unsetenv(config.DefaultOSEnvVariableName)) }()

		srv, err := config.NewService(storage.NewMap(), config.Options{})
		assert.Nil(t, srv)
		assert.True(t, errors.NotValid.Match(err), "%+v", err)
		assert.EqualError(t, err, "[config] Environment key \"CS_ENV\" contains invalid non-letter characters: \"STAG\\uf8ffING\"")
	})
}

func TestHotReload(t *testing.T) {
	defer leaktest.Check(t)()

	p := config.MustMakePath("ww/ee/rr")
	reloadCounter := new(int64)
	srv := config.MustNewService(storage.NewMap(), config.Options{
		EnableHotReload:  true,
		HotReloadSignals: []os.Signal{syscall.SIGUSR1},
	},
		config.MakeLoadDataOption(func(s *config.Service) error {
			return s.Set(p, strconv.AppendInt(nil, atomic.AddInt64(reloadCounter, 1), 10))
		}),
		config.WithFieldMeta(
			&config.FieldMeta{
				Route:          "carrier/dhl/timeout",
				WriteScopePerm: scope.PermDefault,
				ScopeID:        scope.DefaultTypeID, // gets ignored
				Default:        "3601s",
			},
		).WithSortOrder(2),
	)
	defer func() { assert.NoError(t, srv.Close()) }()

	assert.Exactly(t, `"1"`, srv.Get(p).String())

	syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)

	time.Sleep(time.Millisecond * 2)
	assert.Exactly(t, `"2"`, srv.Get(p).String())

	pTimeout := config.MustMakePath("carrier/dhl/timeout")
	assert.Exactly(t, `"3601s"`, srv.Get(pTimeout).String())
}

type keyer interface {
	Keys(ret ...string) []string
}

func TestService_PathWithEnvName(t *testing.T) {
	p := config.MustMakePath("payment/datatrans/username").BindWebsite(4).WithEnvSuffix()

	sMap := storage.NewMap()
	srv := config.MustNewService(sMap, config.Options{
		EnvName: "MY_LOCAL_MACBOOK",
	})
	defer func() { assert.NoError(t, srv.Close()) }()

	assert.NoError(t, srv.Set(p, []byte(`TESCHT1`)))
	p.UseEnvSuffix = false
	assert.NoError(t, srv.Set(p, []byte(`TESCHT2`)))

	assert.Exactly(t, `"TESCHT1"`, srv.Get(p.WithEnvSuffix()).String())
	p.UseEnvSuffix = false
	assert.Exactly(t, `"TESCHT2"`, srv.Get(p).String())

	keys := sMap.(keyer).Keys()
	sort.Strings(keys)
	assert.Exactly(t, []string{
		"Type(Website) ID(4)/payment/datatrans/username",
		"Type(Website) ID(4)/payment/datatrans/username/MY_LOCAL_MACBOOK",
	}, keys)
}

func TestService_DifferentStorageLevels(t *testing.T) {
	sMap := storage.NewMap()
	srv := config.MustNewService(sMap, config.Options{})
	defer func() { assert.NoError(t, srv.Close()) }()
}

type testObserver struct {
	err     error
	rawData []byte
	observe func(p config.Path, rawData []byte, found bool) (rawData2 []byte, err error)
}

func (to testObserver) Observe(p config.Path, rawData []byte, found bool) (rawData2 []byte, err error) {
	if to.observe != nil {
		return to.observe(p, rawData, found)
	}
	if to.rawData != nil {
		return to.rawData, nil
	}
	return rawData, to.err
}

func TestService_Observer(t *testing.T) {
	srv := config.MustNewService(storage.NewMap(), config.Options{})
	defer func() { assert.NoError(t, srv.Close()) }()

	pUser := config.MustMakePath("carrier/dhl/username")
	dUserName := []byte(`96703400169141436bc769418a3577e5`)

	assert.NoError(t, srv.Set(pUser, dUserName))

	t.Run("Register with out of range event ID", func(t *testing.T) {
		err := srv.RegisterObserver(81, "aa/bb/cc", nil)
		assert.True(t, errors.OutOfRange.Match(err), "%+v", err)
	})
	t.Run("DeregisterObserver with out of range event ID", func(t *testing.T) {
		err := srv.DeregisterObserver(82, "aa/bb/cc")
		assert.True(t, errors.OutOfRange.Match(err), "%+v", err)
	})

	t.Run("Before GET returns error", func(t *testing.T) {
		assert.NoError(t, srv.RegisterObserver(config.EventOnBeforeGet, "carrier/dhl/username", testObserver{
			err: errors.AlreadyInUse.Newf("Ups"),
		}))

		str, ok, err := srv.Get(pUser).Str()
		assert.Empty(t, str, "Get should return empty string")
		assert.False(t, ok, "Get should confirm string cannot be found due to an error")
		assert.True(t, errors.AlreadyInUse.Match(err), "%+v", err)

		assert.NoError(t, srv.DeregisterObserver(config.EventOnBeforeGet, "carrier/dhl/username"))
	})

	t.Run("After GET returns error", func(t *testing.T) {
		assert.NoError(t, srv.RegisterObserver(config.EventOnAfterGet, "carrier/dhl/username", testObserver{
			err: errors.AlreadyCaptured.Newf("Ups"),
		}))

		str, ok, err := srv.Get(pUser).Str()
		assert.Empty(t, str, "Get should return empty string")
		assert.False(t, ok, "Get should confirm string cannot be found due to an error")
		assert.True(t, errors.AlreadyCaptured.Match(err), "%+v", err)

		assert.NoError(t, srv.DeregisterObserver(config.EventOnAfterGet, "carrier/dhl/username"))
	})

	t.Run("Before SET returns error", func(t *testing.T) {
		assert.NoError(t, srv.RegisterObserver(config.EventOnBeforeSet, "aa/bb/cc", testObserver{
			err: errors.AlreadyInUse.Newf("Ups"),
		}))

		p := config.MustMakePathWithScope(scope.Website.WithID(2), "aa/bb/cc")
		data := []byte(`4711`)
		err := srv.Set(p, data)
		assert.True(t, errors.AlreadyInUse.Match(err), "%+v", err)

		assert.NoError(t, srv.DeregisterObserver(config.EventOnBeforeSet, "aa/bb/cc"))
	})

	t.Run("After SET returns error", func(t *testing.T) {
		assert.NoError(t, srv.RegisterObserver(config.EventOnAfterSet, "aa/bb/cc", testObserver{
			err: errors.AlreadyCaptured.Newf("Ups"),
		}))

		p := config.MustMakePathWithScope(scope.Website.WithID(2), "aa/bb/cc")
		data := []byte(`4711`)
		err := srv.Set(p, data)
		assert.True(t, errors.AlreadyCaptured.Match(err), "%+v", err)

		assert.NoError(t, srv.DeregisterObserver(config.EventOnAfterSet, "aa/bb/cc"))
	})

	t.Run("dispatch", func(t *testing.T) {
		assert.NoError(t, srv.RegisterObserver(config.EventOnBeforeSet, "aa/bb/dd", testObserver{
			rawData: []byte(`0816`),
		}))

		getsCalledGet := false
		assert.NoError(t, srv.RegisterObserver(config.EventOnBeforeGet, "aa/bb", testObserver{
			observe: func(p config.Path, rawData []byte, found bool) (rawData2 []byte, err error) {
				assert.False(t, found)
				assert.Nil(t, rawData)
				assert.Exactly(t, `websites/2/aa/bb/dd`, p.String())
				getsCalledGet = true
				return rawData, nil
			},
		}))

		getsCalledSet := false
		assert.NoError(t, srv.RegisterObserver(config.EventOnAfterSet, "aa/bb/dd", testObserver{
			observe: func(p config.Path, rawData []byte, found bool) (rawData2 []byte, err error) {
				assert.Exactly(t, `0816`, string(rawData))
				getsCalledSet = true
				return rawData, nil
			},
		}))

		p := config.MustMakePathWithScope(scope.Website.WithID(2), "aa/bb/dd")
		data := []byte(`4711`)
		assert.NoError(t, srv.Set(p, data))

		assert.Exactly(t, `"0816"`, srv.Get(p).String())
		assert.True(t, getsCalledSet, "Event after set should get called")
		assert.True(t, getsCalledGet, "Event before get should get called")
		assert.NoError(t, srv.DeregisterObserver(config.EventOnBeforeSet, "/aa/bb/dd"))
	})
}

func genFieldMetaChan(fms ...*config.FieldMeta) func(*config.Service) (<-chan *config.FieldMeta, <-chan error) {
	return func(s *config.Service) (<-chan *config.FieldMeta, <-chan error) {
		fmc := make(chan *config.FieldMeta)
		errc := make(chan error)

		go func() {
			defer close(fmc)
			defer close(errc)

			for _, fm := range fms {
				fmc <- fm
			}
		}()

		return fmc, errc
	}
}

func TestService_FieldMetaData_CheckScope_Sync(t *testing.T) {
	t.Parallel()
	srv, err := config.NewService(storage.NewMap(), config.Options{},
		config.WithFieldMetaGenerator(
			genFieldMetaChan(&config.FieldMeta{
				Route:          "carrier/dhl/timeout",
				WriteScopePerm: scope.PermWebsite,
				ScopeID:        scope.Store.WithID(4),
				Default:        "3600s",
			}),
		),
	)
	assert.Nil(t, srv)
	assert.True(t, errors.NotAcceptable.Match(err), "%+v", err)
	assert.EqualError(t, err, "[config] WriteScopePerm \"websites\" and ScopeID \"Type(Store) ID(4)\" cannot be set at once for path \"carrier/dhl/timeout\"")
}

func TestService_FieldMetaData_CheckScope_Async(t *testing.T) {
	t.Parallel()
	srv, err := config.NewService(storage.NewMap(), config.Options{},
		config.WithFieldMeta(
			&config.FieldMeta{
				Route:          "carrier/dhl/timeout",
				WriteScopePerm: scope.PermWebsite,
				ScopeID:        scope.Store.WithID(4),
				Default:        "3600s",
			},
		),
	)
	assert.Nil(t, srv)
	assert.True(t, errors.NotAcceptable.Match(err), "%+v", err)
	assert.EqualError(t, err, "[config] WriteScopePerm \"websites\" and ScopeID \"Type(Store) ID(4)\" cannot be set at once for path \"carrier/dhl/timeout\"")
}

func TestService_FieldMetaData_Fallbacks(t *testing.T) {
	t.Parallel()

	srv := config.MustNewService(storage.NewMap(), config.Options{},
		config.WithFieldMeta(
			&config.FieldMeta{
				Route:          "catalog/layered_navigation/display_product_count",
				WriteScopePerm: 0, // writing allowed from all scopes
				Default:        "false",
			},
			&config.FieldMeta{
				Route:          "carrier/dhl/timeout",
				WriteScopePerm: scope.PermDefault,
				ScopeID:        scope.DefaultTypeID, // gets ignored
				Default:        "3600s",
			},
			&config.FieldMeta{
				Route:          "carrier/dhl/username", // for all other websites/stores this is the default value
				WriteScopePerm: scope.PermWebsite,      // Perm can only be set when ScopeID is <= DefaultID.
				Default:        "prdUser0",
			},
			&config.FieldMeta{
				Route:   "carrier/dhl/username",
				ScopeID: scope.Website.WithID(1), // for Website 1 this is the default value
				Default: "prdUser1",
			},
			&config.FieldMeta{
				Route:   "carrier/dhl/username",
				ScopeID: scope.Website.WithID(2), // for Website 2 this is the default value
				Default: "prdUser2",
			},
		).WithSortOrder(3),

		config.WithApplySections(
			&config.Section{
				ID: `customer`,
				Groups: config.MakeGroups(
					&config.Group{
						ID: `address`,
						Fields: config.MakeFields(
							&config.Field{
								ID:        `prefix_options`,
								Default:   `Bitte wählen;Frau;Herr`,
								Type:      config.TypeSelect,
								Label:     `Name Title`,
								SortOrder: 100,
								Visible:   true,
							},
						),
					},
				),
			},
		).WithSortOrder(1),
	)
	defer func() { assert.NoError(t, srv.Close()) }()

	dhlTimeout := config.MustMakePath("carrier/dhl/timeout")
	dhlUsername := config.MustMakePath("carrier/dhl/username")

	t.Run("timeout: store scope cannot write default scope", func(t *testing.T) {
		dhlTimeout := dhlTimeout.BindStore(3)
		err := srv.Set(dhlTimeout, []byte(`1m`))
		assert.True(t, errors.NotAllowed.Match(err), "%+v", err)

		str, ok, err := srv.Get(dhlTimeout).Str()
		assert.True(t, ok)
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, `3600s`, str)
	})
	t.Run("timeout: website scope cannot access default scope", func(t *testing.T) {
		dhlTimeout := dhlTimeout.BindWebsite(3)
		err := srv.Set(dhlTimeout, []byte(`1m`))
		assert.True(t, errors.NotAllowed.Match(err), "%+v", err)

		str, ok, err := srv.Get(dhlTimeout).Str()
		assert.True(t, ok)
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, `3600s`, str)
	})
	t.Run("timeout: default scope can access default scope", func(t *testing.T) {
		err := srv.Set(dhlTimeout, []byte(`1m`))
		assert.NoError(t, err, "%+v", err)

		d, ok, err := srv.Get(dhlTimeout).Duration()
		assert.True(t, ok)
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, time.Duration(60000000000), d)
	})

	t.Run("username: store scope cannot access website scope", func(t *testing.T) {
		dhlUsername := dhlUsername.BindStore(3)
		err := srv.Set(dhlUsername, []byte(`guestUser`))
		assert.True(t, errors.NotAllowed.Match(err), "%+v", err)

		str, ok, err := srv.Get(dhlUsername).Str()
		assert.False(t, ok)
		assert.NoError(t, err, "%+v", err)
		assert.Empty(t, str) // direct hit must return empty value because no default value for default scope set.
	})
	t.Run("username: website scope 3 can access website scope and gets default scope default value", func(t *testing.T) {
		scpd := srv.Scoped(3, 4)
		str := scpd.Get(scope.Website, "carrier/dhl/username").String()
		assert.Exactly(t, `"prdUser0"`, str)
	})
	t.Run("username: website scope 1 can access website scope and gets its scope default value", func(t *testing.T) {
		scpd := srv.Scoped(1, 4)
		str := scpd.Get(scope.Store, "carrier/dhl/username").String()
		assert.Exactly(t, `"prdUser1"`, str)
	})
	t.Run("username: website scope 2 can access website scope and gets its scope default value", func(t *testing.T) {
		scpd := srv.Scoped(2, 4)
		str := scpd.Get(scope.Website, "carrier/dhl/username").String()
		assert.Exactly(t, `"prdUser2"`, str)
	})

	t.Run("prefix_option: scope 5,6 falls back to default", func(t *testing.T) {
		scpd := srv.Scoped(5, 6)
		str := scpd.Get(scope.Website, "customer/address/prefix_options").String()
		assert.Exactly(t, "\"Bitte wählen;Frau;Herr\"", str)
	})
}
