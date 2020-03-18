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

package maxmindwebservice_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/weiwolves/pkg/config/cfgmock"
	"github.com/weiwolves/pkg/net/geoip"
	"github.com/weiwolves/pkg/net/geoip/backendgeoip"
	"github.com/weiwolves/pkg/net/geoip/maxmindwebservice"
	"github.com/weiwolves/pkg/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/weiwolves/pkg/util/assert"
)

var backend *backendgeoip.Configuration

func init() {

	cfgStruct, err := backendgeoip.NewConfigStructure()
	if err != nil {
		panic(err)
	}
	backend = backendgeoip.New(cfgStruct)
}

func TestConfiguration_WithGeoIP2Webservice_Redis(t *testing.T) {

	t.Run("Error_API", testBackend_WithGeoIP2Webservice_Redis(
		func() *http.Client {
			// http://dev.maxmind.com/geoip/geoip2/web-services/#Errors
			return &http.Client{
				Transport: cstesting.NewHTTPTrip(402, `{"error":"The license key you have provided is out of queries.","code":"OUT_OF_QUERIES"}`, nil),
			}
		},
		func(t *testing.T) http.Handler {
			return http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
				panic("Should not get called")
			})
		},
		http.StatusServiceUnavailable,
	))

	t.Run("Error_JSON", testBackend_WithGeoIP2Webservice_Redis(
		func() *http.Client {
			// http://dev.maxmind.com/geoip/geoip2/web-services/#Errors
			return &http.Client{
				Transport: cstesting.NewHTTPTrip(200, `{"error":"The license ... wow this JSON isn't valid.`, nil),
			}
		},
		func(t *testing.T) http.Handler {
			return http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		},
		http.StatusServiceUnavailable,
	))

	var calledSuccessHandler int32
	t.Run("Success", testBackend_WithGeoIP2Webservice_Redis(
		func() *http.Client {
			return &http.Client{
				Transport: cstesting.NewHTTPTrip(200, `{ "continent": { "code": "EU", "geoname_id": 6255148, "names": { "de": "Europa", "en": "Europe", "es": "Europa", "fr": "Europe", "ja": "ヨーロッパ", "pt-BR": "Europa", "ru": "Европа", "zh-CN": "欧洲" } }, "country": { "geoname_id": 2921044, "iso_code": "DE", "names": { "de": "Deutschland", "en": "Germany", "es": "Alemania", "fr": "Allemagne", "ja": "ドイツ連邦共和国", "pt-BR": "Alemanha", "ru": "Германия", "zh-CN": "德国" } }, "registered_country": { "geoname_id": 2921044, "iso_code": "DE", "names": { "de": "Deutschland", "en": "Germany", "es": "Alemania", "fr": "Allemagne", "ja": "ドイツ連邦共和国", "pt-BR": "Alemanha", "ru": "Германия", "zh-CN": "德国" } }, "traits": { "autonomous_system_number": 1239, "autonomous_system_organization": "Linkem IR WiMax Network", "domain": "example.com", "is_anonymous_proxy": true, "is_satellite_provider": true, "isp": "Linkem spa", "ip_address": "1.2.3.4", "organization": "Linkem IR WiMax Network", "user_type": "traveler" }, "maxmind": { "queries_remaining": 54321 } }`, nil),
			}
		},
		func(t *testing.T) http.Handler {
			return http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				cty, ok := geoip.FromContextCountry(r.Context())
				assert.True(t, ok)
				assert.Exactly(t, "DE", cty.Country.IsoCode)
				atomic.AddInt32(&calledSuccessHandler, 1)
			})
		},
		http.StatusOK,
	))
	assert.Exactly(t, int32(80), atomic.LoadInt32(&calledSuccessHandler), "calledSuccessHandler")
}

func testBackend_WithGeoIP2Webservice_Redis(
	hcf func() *http.Client,
	finalHandler func(t *testing.T) http.Handler,
	wantCode int,
) func(*testing.T) {

	return func(t *testing.T) {
		rd := miniredis.NewMiniRedis()
		if err := rd.Start(); err != nil {
			t.Fatal(err)
		}
		defer rd.Close()
		redConURL := fmt.Sprintf("redis://%s/3", rd.Addr())

		// test if we get the correct country and if the country has
		// been successfully stored in redis and can be retrieved.

		//be.WebServiceClient = hcf()
		backend.Register(maxmindwebservice.NewOptionFactory(
			hcf(),
			backend.MaxmindWebserviceUserID,
			backend.MaxmindWebserviceLicense,
			backend.MaxmindWebserviceTimeout,
			backend.MaxmindWebserviceRedisURL,
		))

		cfgSrv := cfgmock.NewService(cfgmock.PathValue{
			// @see structure.go for the limitation to scope.Default
			backend.DataSource.MustFQ():                maxmindwebservice.OptionName,
			backend.MaxmindWebserviceUserID.MustFQ():   `TestUserID`,
			backend.MaxmindWebserviceLicense.MustFQ():  `TestLicense`,
			backend.MaxmindWebserviceTimeout.MustFQ():  `21s`,
			backend.MaxmindWebserviceRedisURL.MustFQ(): redConURL,
		})
		cfgScp := cfgSrv.NewScoped(1, 2) // Website ID 2 == euro / Store ID == 2 Austria ==> here doesn't matter

		geoSrv := geoip.MustNew()

		req := func() *http.Request {
			req, _ := http.NewRequest("GET", "http://corestore.io", nil)
			req.Header.Set("X-Cluster-Client-Ip", "2a02:d180::") // Germany
			return req
		}()

		scpFnc := backend.PrepareOptionFactory()
		if err := geoSrv.Options(scpFnc(cfgScp)...); err != nil {
			t.Fatalf("%+v", err)
		}
		// food for the race detector
		hpu := cstesting.NewHTTPParallelUsers(8, 10, 500, time.Millisecond) // 8,10
		hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
			assert.Exactly(t, wantCode, rec.Code)
		}
		hpu.ServeHTTP(req, geoSrv.WithCountryByIP(finalHandler(t)))
	}
}

func TestConfiguration_Path_Errors(t *testing.T) {

	backend.Register(maxmindwebservice.NewOptionFactory(
		&http.Client{},
		backend.MaxmindWebserviceUserID,
		backend.MaxmindWebserviceLicense,
		backend.MaxmindWebserviceTimeout,
		backend.MaxmindWebserviceRedisURL,
	))

	tests := []struct {
		toPath string
		val    interface{}
		errBhf errors.BehaviourFunc
	}{
		0: {backend.MaxmindWebserviceUserID.MustFQ(), struct{}{}, errors.IsNotValid},
		1: {backend.MaxmindWebserviceLicense.MustFQ(), struct{}{}, errors.IsNotValid},
		2: {backend.MaxmindWebserviceTimeout.MustFQ(), struct{}{}, errors.IsNotValid},
		3: {backend.MaxmindWebserviceRedisURL.MustFQ(), struct{}{}, errors.IsNotValid},
	}
	for i, test := range tests {

		cfgSrv := cfgmock.NewService(cfgmock.PathValue{
			backend.DataSource.MustFQ(): maxmindwebservice.OptionName,
			test.toPath:                 test.val,
		})

		gs := geoip.MustNew(
			geoip.WithRootConfig(cfgSrv),
			geoip.WithOptionFactory(backend.PrepareOptionFactory()),
		)
		assert.NoError(t, gs.ClearCache())
		_, err := gs.ConfigByScope(0, 0)
		assert.True(t, test.errBhf(err), "Index %d Error: %+v", i, err)
	}
}

func TestNewOptionFactory_Invalid_ConfigValue(t *testing.T) {

	backend.Register(maxmindwebservice.NewOptionFactory(
		&http.Client{},
		backend.MaxmindWebserviceUserID,
		backend.MaxmindWebserviceLicense,
		backend.MaxmindWebserviceTimeout,
		backend.MaxmindWebserviceRedisURL,
	))

	cfgSrv := cfgmock.NewService(cfgmock.PathValue{
		backend.DataSource.MustFQ(): maxmindwebservice.OptionName,
	})

	gs := geoip.MustNew(
		geoip.WithRootConfig(cfgSrv),
		geoip.WithOptionFactory(backend.PrepareOptionFactory()),
	)
	//assert.NoError(t, gs.ClearCache())
	_, err := gs.ConfigByScope(1, 0)
	assert.True(t, errors.IsNotValid(err), " Error: %+v", err)
	assert.Contains(t, err.Error(), `Incomplete WebService configuration: User: "" License "" Ti`)
}
