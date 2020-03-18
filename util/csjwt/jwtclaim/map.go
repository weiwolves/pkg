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

package jwtclaim

import (
	"encoding/json"
	"time"

	"github.com/corestoreio/errors"
	"github.com/weiwolves/pkg/util/conv"
)

// Map default type for the Claim field in a token. Slowest but most flexible
// type. For speed, use a custom struct type with embedding StandardClaims and
// generated en-/decoder.
//easyjson:json
type Map map[string]interface{}

// VerifyAudience compares the aud claim against cmp. If required is false, this
// method will return true if the value matches or is unset
func (m Map) VerifyAudience(cmp string, req bool) bool {
	aud := conv.ToByte(m["aud"])
	return verifyConstantTime(aud, []byte(cmp), req)
}

func (m Map) exp() int64 {
	return conv.ToInt64(m["exp"])
}

func (m Map) iat() int64 {
	return conv.ToInt64(m["iat"])
}

func (m Map) nbf() int64 {
	return conv.ToInt64(m["nbf"])
}

func (m Map) skew() time.Duration {
	if d, ok := m[KeyTimeSkew]; ok {
		return conv.ToDuration(d)
	}
	return 0
}

// Compares the exp claim against cmp. If required is false, this method will
// return true if the value matches or is unset
func (m Map) VerifyExpiresAt(cmp int64, req bool) bool {
	return verifyExp(m.skew(), m.exp(), cmp, req)
}

// Compares the iat claim against cmp. If required is false, this method will
// return true if the value matches or is unset.
func (m Map) VerifyIssuedAt(cmp int64, req bool) bool {
	return verifyIat(m.iat(), cmp, req)
}

// Compares the iss claim against cmp. If required is false, this method will
// return true if the value matches or is unset.
func (m Map) VerifyIssuer(cmp string, req bool) bool {
	iss := conv.ToByte(m["iss"])
	return verifyConstantTime(iss, []byte(cmp), req)
}

// Compares the nbf claim against cmp. If required is false, this method will
// return true if the value matches or is unset.
func (m Map) VerifyNotBefore(cmp int64, req bool) bool {
	return verifyNbf(m.skew(), m.nbf(), cmp, req)
}

// Validates time based claims "exp, iat, nbf". There is no accounting for clock
// skew. As well, if any of the above claims are not in the token, it will still
// be considered a valid claim.
func (m Map) Valid() error {
	now := TimeFunc().Unix()

	switch {
	case len(m) == 0:
		return errors.NotValid.Newf(`[jwtclaim] token claims validation failed1`)

	// case m.exp() == 0 && m.iat() == 0 && m.nbf() == 0:
	//	return errors.NotValid.Newf(`[jwtclaim] token claims validation failed2`)

	case !m.VerifyExpiresAt(now, false):
		return errors.NotValid.Newf(`[jwtclaim] token is expired %s ago`, TimeFunc().Sub(time.Unix(m.exp(), 0)))

	case !m.VerifyIssuedAt(now, false):
		return errors.NotValid.Newf(`[jwtclaim] token used before issued, clock skew issue? Diff %s`, time.Unix(m.iat(), 0).Sub(TimeFunc()))

	case !m.VerifyNotBefore(now, false):
		return errors.NotValid.Newf(`[jwtclaim] token is not valid yet. Diff %s`, time.Unix(m.nbf(), 0).Sub(TimeFunc()))
	}

	return nil
}

func (m Map) Set(key string, value interface{}) error {
	m[key] = value
	return nil
}

// Get can return nil,nil
func (m Map) Get(key string) (interface{}, error) {
	return m[key], nil
}

func (m Map) Keys() []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

// Expires duration when a token expires.
func (m Map) Expires() (exp time.Duration) {
	if cexp, ok := m["exp"]; ok {
		fexp := conv.ToFloat64(cexp)
		if fexp > 0.001 {
			tm := time.Unix(int64(fexp), 0)
			if remainer := tm.Sub(time.Now()); remainer > 0 {
				exp = remainer
			}
		}
	}
	return
}

// String human readable output via JSON, slow.
func (m Map) String() string {
	b, err := json.Marshal(m)
	if err != nil {
		return errors.Fatal.Newf("[jwtclaim] Map.String(): json.Marshal Error: %s", err).Error()
	}
	return string(b)
}
