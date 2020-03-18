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

package geoip

import (
	"net/http"

	"github.com/weiwolves/pkg/net/mw"
)

// DefaultAlternativeHandler gets called when detected Country cannot be found
// within the list of allowed countries. This handler can be overridden to provide
// a fallback for all scopes. To set a alternative handler for a website or store
// use the With*() options. This function gets called in WithIsCountryAllowedByIP.
//
// Status is StatusServiceUnavailable
var DefaultAlternativeHandler mw.ErrorHandler = mw.ErrorWithStatusCode(http.StatusServiceUnavailable)
