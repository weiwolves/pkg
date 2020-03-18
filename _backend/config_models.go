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

package backend

import (
	"net/http"

	"github.com/weiwolves/pkg/config"
	"github.com/weiwolves/pkg/config/cfgmodel"
	"github.com/weiwolves/pkg/config/source"
	"github.com/weiwolves/pkg/store/scope"
)

// ConfigRedirectToBase enables if a redirect to the base URL should
// happen and with which status code.
type ConfigRedirectToBase struct {
	cfgmodel.Int
}

// NewConfigRedirectToBase creates a new type.
func NewConfigRedirectToBase(path string, opts ...cfgmodel.Option) ConfigRedirectToBase {
	return ConfigRedirectToBase{
		Int: cfgmodel.NewInt(
			path,
			append(opts, cfgmodel.WithSourceByInt(source.Ints{
				0: {0, "No"},
				1: {1, "Yes (302 Found)"},                // old from Magento :-(
				2: {http.StatusFound, "Yes (302 Found)"}, // new correct
				3: {http.StatusMovedPermanently, "Yes (301 Moved Permanently)"},
			}))...,
		),
	}
}

// Write writes an int value and checks if the int value is within the allowed Options.
func (p ConfigRedirectToBase) Write(w config.Setter, v int, s scope.Type, id int64) error {

	if err := p.ValidateInt(v); err != nil {
		return err
	}

	return p.Int.Write(w, v, s, id)
}
