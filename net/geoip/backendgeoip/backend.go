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

package backendgeoip

import (
	"github.com/weiwolves/pkg/config/cfgmodel"
	"github.com/weiwolves/pkg/config/cfgsource"
	"github.com/weiwolves/pkg/config/element"
	"github.com/weiwolves/pkg/net/geoip"
)

// Configuration just exported for the sake of documentation. See fields for more
// information.
type Configuration struct {
	*geoip.OptionFactories

	// AllowedCountries list of countries which are currently allowed.
	// Separated via comma, e.g.: DE,CH,AT,AU,NZ,
	//
	// Path: net/geoip/allowed_countries
	AllowedCountries cfgmodel.StringCSV

	// AlternativeRedirect redirects the client to this URL if their
	// country hasn't been granted access to the next middleware handler.
	//
	// Path: net/geoip/alternative_redirect
	AlternativeRedirect cfgmodel.URL

	// AlternativeRedirectCode HTTP redirect code.
	//
	// Path: net/geoip/alternative_redirect_code
	AlternativeRedirectCode cfgmodel.Int

	// DataSource defines to either load the Geo location data from a MaxMind
	// "file" or from the MaxMind "webservice".
	//
	// Path: net/geoip_maxmind/data_source
	DataSource cfgmodel.Str

	// MaxmindLocalFile path to a file name stored on the server.
	//
	// Path: net/geoip_maxmind/local_file
	MaxmindLocalFile cfgmodel.Str

	// MaxmindWebserviceUserID user id
	//
	// Path: net/geoip_maxmind/webservice_userid
	MaxmindWebserviceUserID cfgmodel.Str

	// MaxmindWebserviceLicense license name
	//
	// Path: net/geoip_maxmind/webservice_license
	MaxmindWebserviceLicense cfgmodel.Str

	// MaxmindWebserviceTimeout HTTP request time out
	//
	// Path: net/geoip_maxmind/webservice_timeout
	MaxmindWebserviceTimeout cfgmodel.Duration

	// MaxmindWebserviceRedisURL an URL to the Redis server
	//
	// Path: net/geoip_maxmind/webservice_redisurl
	MaxmindWebserviceRedisURL cfgmodel.URL
}

// New initializes the backend configuration models containing the cfgpath.Route
// variable to the appropriate entries. The function Load() will be executed to
// apply the Sections to all models. See Load() for more details.
func New(cfgStruct element.Sections, opts ...cfgmodel.Option) *Configuration {
	be := &Configuration{
		OptionFactories: geoip.NewOptionFactories(),
	}

	opts = append(opts, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	optsRedir := append([]cfgmodel.Option{}, opts...)
	optsRedir = append(optsRedir, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(redirects))

	be.AllowedCountries = cfgmodel.NewStringCSV(`net/geoip/allowed_countries`, opts...)
	be.AlternativeRedirect = cfgmodel.NewURL(`net/geoip/alternative_redirect`, opts...)
	be.AlternativeRedirectCode = cfgmodel.NewInt(`net/geoip/alternative_redirect_code`, optsRedir...)

	be.DataSource = cfgmodel.NewStr(`net/geoip_maxmind/data_source`, append(opts, cfgmodel.WithSourceByString(
		"file", "File on this server",
		"webservice", "Maxmind web service",
	))...)
	be.MaxmindLocalFile = cfgmodel.NewStr(`net/geoip_maxmind/local_file`, opts...)
	be.MaxmindWebserviceUserID = cfgmodel.NewStr(`net/geoip_maxmind/webservice_userid`, opts...)
	be.MaxmindWebserviceLicense = cfgmodel.NewStr(`net/geoip_maxmind/webservice_license`, opts...)
	be.MaxmindWebserviceTimeout = cfgmodel.NewDuration(`net/geoip_maxmind/webservice_timeout`, opts...)
	be.MaxmindWebserviceRedisURL = cfgmodel.NewURL(`net/geoip_maxmind/webservice_redisurl`, opts...)

	return be
}

// Load creates the configuration models for each PkgBackend field. Internal
// mutex will protect the fields during loading. The argument Sections will
// be applied to all models.

var redirects = cfgsource.NewByInt(
	cfgsource.Ints{
		{301, "301 moved permanently"},
		{302, "302 found"},
		{303, "303 see other"},
		{308, "308 permanent redirect "},
	},
)
