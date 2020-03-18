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

// +build csall bigcache

package storage_test

import (
	"testing"

	"github.com/allegro/bigcache"
	"github.com/corestoreio/errors"
	"github.com/weiwolves/pkg/config"
	"github.com/weiwolves/pkg/config/storage"
	"github.com/weiwolves/pkg/store/scope"
	"github.com/weiwolves/pkg/util/assert"
)

func TestCacheGet(t *testing.T) {
	bgc, err := storage.NewBigCache(bigcache.Config{
		Shards: 64,
	})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		scp        scope.TypeID
		path       string
		val        []byte
		wantSetErr error
		wantGetErr error
	}{
		{scope.DefaultTypeID, "aa/bb/cc", []byte(`DataXYZ`), nil, nil},
		{scope.Store.WithID(3), "aa/bb/cc", []byte(`DataXYA`), nil, nil},
	}
	for idx, test := range tests {
		p := config.MustMakePathWithScope(test.scp, test.path)
		haveSetErr := bgc.Set(p, test.val)
		if test.wantSetErr != nil {
			assert.EqualError(t, haveSetErr, test.wantSetErr.Error(), "Index %d", idx)
		} else {
			assert.NoError(t, haveSetErr, "Index %d", idx)
		}

		haveData, haveOK, haveGetErr := bgc.Get(p)
		if test.wantGetErr != nil {
			assert.EqualError(t, haveGetErr, test.wantGetErr.Error(), "Index %d", idx)
			assert.False(t, haveOK)
		} else {
			assert.NoError(t, haveGetErr, "Index %d", idx)
		}
		// don't do this 2x conv casting in production code
		assert.Exactly(t, test.val, haveData, "Index %d", idx)
	}
}

func TestCacheGetNotFound(t *testing.T) {
	sc, err := storage.NewBigCache(bigcache.Config{
		Shards: 64,
	})
	if err != nil {
		t.Fatal(err)
	}
	haveVal, haveFound, haveGetErr := sc.Get(config.MustMakePath("aa/bb/cc"))
	assert.False(t, haveFound)
	assert.NoError(t, haveGetErr)
	assert.Empty(t, haveVal)
}

func TestCacheError(t *testing.T) {
	sc, err := storage.NewBigCache(bigcache.Config{
		Shards: 63,
	})
	assert.True(t, errors.Fatal.Match(err), "Error: %s", err)
	assert.Empty(t, sc)
}
