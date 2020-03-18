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

package directory

import (
	"github.com/corestoreio/errors"
	"github.com/weiwolves/pkg/config"
	"github.com/weiwolves/pkg/config/cfgmodel"
	"github.com/weiwolves/pkg/store/scope"
	"golang.org/x/text/currency"
)

// ConfigCurrency currency type for the configuration based on text/currency pkg.
type ConfigCurrency struct {
	cfgmodel.Str
}

// NewConfigCurrency creates a new currency configuration type.
func NewConfigCurrency(path string, opts ...cfgmodel.Option) ConfigCurrency {
	return ConfigCurrency{
		Str: cfgmodel.NewStr(path, opts...),
	}
}

// GetDefault returns the default currency without considering the scope.
func (cc ConfigCurrency) GetDefault(sg config.Getter) (cur Currency, err error) {
	p, err := cc.ToPath(scope.Default, 0)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	raw, err := sg.String(p)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	cur.Unit, err = currency.ParseISO(raw)
	return
}

// Get tries to retrieve a currency considering the scope
func (cc ConfigCurrency) Get(sg config.Scoped) (cur Currency, err error) {
	raw, err := cc.Str.Get(sg)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	if raw == "" {
		scp, scpID := sg.ScopeID()
		err = errors.Errorf("Empty currency for path: %q, scope: %q, scopeID: %d", cc.String(), scp, scpID)
		return
	}
	cur.Unit, err = currency.ParseISO(raw)
	return
}

// Writes a currency to the configuration storage.
func (cc ConfigCurrency) Write(w config.Setter, v Currency, s scope.Type, id int64) error {
	cur := v.String()

	if err := cc.ValidateString(cur); err != nil {
		return errors.WithStack(err)
	}
	return cc.Str.Write(w, cur, s, id)
}
