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

package backendauth

import (
	"github.com/corestoreio/errors"
	"github.com/weiwolves/pkg/config"
	"github.com/weiwolves/pkg/config/cfgmodel"
	"github.com/weiwolves/pkg/net"
)

// ConfigIPRange defines how IP ranges are stored and handled.
// A valid IP range string looks like for example:
// 		IPv4: 74.50.153.0-74.50.153.4
// 		IPv6: ::ffff:192.0.2.128-::ffff:192.0.2.250
// 		IPv6: 2001:0db8:85a3:0000:0000:8a2e:0370:7334-2001:0db8:85a3:0000:0000:8a2e:0370:8334
// No white spaces! Multiple entries supported via \r and/or \n.
type ConfigIPRange struct {
	cfgmodel.CSV
}

// NewConfigIPRange ....
func NewConfigIPRange(path string, opts ...cfgmodel.Option) ConfigIPRange {
	cip := ConfigIPRange{
		CSV: cfgmodel.NewCSV(path, opts...),
	}
	cip.Comma = '-'
	return cip
}

// Get ...
func (cc ConfigIPRange) Get(sg config.Scoped) (net.IPRanges, error) {
	data, err := cc.CSV.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[backendauth] Str.Get")
	}

	var rngs net.IPRanges
	for _, row := range data {
		if len(row) != 2 {
			return nil, errors.NotValid.Newf("[backendauth] IP Range %q not in expected format: IP.From-IP.To", row)
		}
		if row[0] != "" && row[1] != "" {
			rngs = append(rngs, net.MakeIPRange(row[0], row[1]))
		}
	}
	return rngs, nil
}
