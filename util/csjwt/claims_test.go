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

package csjwt_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/weiwolves/pkg/util/assert"
	"github.com/weiwolves/pkg/util/csjwt"
	"github.com/weiwolves/pkg/util/csjwt/jwtclaim"
)

var _ csjwt.Header = (*csjwt.Head)(nil)

type claimMock struct {
	validErr error
	setErr   error
	getErr   error
}

func (c claimMock) Valid() error                            { return c.validErr }
func (c claimMock) Expires() time.Duration                  { return 0 }
func (c claimMock) Set(key string, value interface{}) error { return c.setErr }
func (c claimMock) Get(key string) (interface{}, error) {
	return nil, c.getErr
}
func (c claimMock) Keys() []string { return []string{"k1"} }

func TestNewHead(t *testing.T) {
	var h csjwt.Header
	h = csjwt.NewHead("X")
	assert.Exactly(t, "X", h.Alg())
	assert.Exactly(t, csjwt.ContentTypeJWT, h.Typ())
}

func TestHeadSetGet(t *testing.T) {
	var h csjwt.Header
	h = csjwt.NewHead("X")

	assert.NoError(t, h.Set(jwtclaim.HeaderAlg, "Y"))
	g, err := h.Get(jwtclaim.HeaderAlg)
	assert.NoError(t, err)
	assert.Exactly(t, "Y", g)

	assert.NoError(t, h.Set(jwtclaim.HeaderTyp, "JWE"))
	g, err = h.Get(jwtclaim.HeaderTyp)
	assert.NoError(t, err)
	assert.Exactly(t, "JWE", g)

	assert.True(t, errors.NotSupported.Match(h.Set("x", "y")))
	g, err = h.Get("x")
	assert.True(t, errors.NotSupported.Match(err))
	assert.Empty(t, g)
}

func TestMergeClaims(t *testing.T) {
	tests := []struct {
		dst               *csjwt.Token
		srcs              csjwt.Claimer
		wantSigningString string
		wantErrKind       errors.Kind
	}{
		{csjwt.NewToken(nil), nil, `eyJ0eXAiOiJKV1QifQ.bnVsbA`, errors.NoKind},
		{csjwt.NewToken(jwtclaim.Map{}), claimMock{getErr: errors.Fatal.Newf("claimMerge get error")}, ``, errors.Fatal},
		{csjwt.NewToken(jwtclaim.Map{"k1": "v1"}), jwtclaim.Map{"k2": 2}, `eyJ0eXAiOiJKV1QifQ.eyJrMSI6InYxIiwiazIiOjJ9`, errors.NoKind},
		{csjwt.NewToken(jwtclaim.NewStore()), jwtclaim.Map{"k2": 2}, ``, errors.NotSupported},
		{csjwt.NewToken(&jwtclaim.Standard{}), &jwtclaim.Store{
			UserID: "Gopher",
		}, ``, errors.NotSupported},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("Index_%02d", i), func(t *testing.T) {
			haveErr := csjwt.MergeClaims(test.dst.Claims, test.srcs)
			if !test.wantErrKind.Empty() {
				assert.ErrorIsKind(t, test.wantErrKind, haveErr)
				return
			}
			assert.NoError(t, haveErr)
			buf, err := test.dst.SigningString(nil)
			assert.NoError(t, err)
			assert.Exactly(t, test.wantSigningString, string(buf), "Index %d", i)
		})
	}
}

func TestClaimExpiresSkew(t *testing.T) {
	st := jwtclaim.NewStore()
	st.ExpiresAt = time.Now().Unix() - 2
	st.Store = "HelloWorld"
	tk := csjwt.NewToken(st)

	pwKey := csjwt.WithPasswordRandom()
	hs256 := csjwt.NewSigningMethodHS256()
	token, err := tk.SignedString(hs256, pwKey)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	vrf := csjwt.NewVerification(hs256)

	parsedTK := csjwt.NewToken(&jwtclaim.Store{
		Standard: jwtclaim.Standard{
			TimeSkew: 0,
		},
	})
	parsedErr := vrf.Parse(parsedTK, token, csjwt.NewKeyFunc(hs256, pwKey))
	assert.ErrorIsKind(t, errors.NotValid, parsedErr)
	assert.False(t, parsedTK.Valid, "Token must be not valid")

	// now adjust skew
	parsedTK = csjwt.NewToken(&jwtclaim.Store{
		Standard: jwtclaim.Standard{
			TimeSkew: time.Second * 3,
		},
	})
	parsedErr = vrf.Parse(parsedTK, token, csjwt.NewKeyFunc(hs256, pwKey))
	assert.NoError(t, parsedErr, "Error: %s", parsedErr)
	assert.True(t, parsedTK.Valid, "Token must be valid")

	haveSt, err := parsedTK.Claims.Get(jwtclaim.KeyStore)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, "HelloWorld", haveSt)
}
