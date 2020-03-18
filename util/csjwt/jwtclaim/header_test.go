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

package jwtclaim_test

import (
	"encoding/json"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/weiwolves/pkg/util/assert"
	"github.com/weiwolves/pkg/util/csjwt"
	"github.com/weiwolves/pkg/util/csjwt/jwtclaim"
)

var _ csjwt.Header = (*jwtclaim.HeadSegments)(nil)

func TestHeadSegmentsParseJSON(t *testing.T) {
	var sc csjwt.Header
	sc = &jwtclaim.HeadSegments{
		Algorithm: `ES999`,
		Type:      jwtclaim.ContentTypeJWT,
	}
	rawJSON, err := json.Marshal(sc)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, rawJSON, 27, "JSON: %s", rawJSON)

	scNew := &jwtclaim.HeadSegments{}
	if err := json.Unmarshal(rawJSON, &scNew); err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, sc, scNew)
}

func TestHeadSegmentsAlgTyp(t *testing.T) {
	var sc csjwt.Header
	sc = jwtclaim.NewHeadSegments(`ES999`)
	assert.Exactly(t, "ES999", sc.Alg())
	assert.Exactly(t, jwtclaim.ContentTypeJWT, sc.Typ())
}

func TestHeadSegmentsGetSet(t *testing.T) {
	tests := []struct {
		sc             csjwt.Header
		key            string
		val            string
		wantSetErrKind errors.Kind
		wantGetErrKind errors.Kind
	}{
		{&jwtclaim.HeadSegments{}, jwtclaim.HeaderAlg, "", errors.NoKind, errors.NoKind},
		{&jwtclaim.HeadSegments{}, jwtclaim.HeaderTyp, "Go", errors.NoKind, errors.NoKind},
		{&jwtclaim.HeadSegments{}, "ext", "Test", errors.NotSupported, errors.NotSupported},
	}
	for i, test := range tests {

		haveSetErr := test.sc.Set(test.key, test.val)
		if !test.wantSetErrKind.Empty() {
			assert.True(t, test.wantSetErrKind.Match(haveSetErr), "Index %d => %s", i, haveSetErr)
		} else {
			assert.NoError(t, haveSetErr, "Index %d", i)
		}

		haveVal, haveGetErr := test.sc.Get(test.key)
		if !test.wantGetErrKind.Empty() {
			assert.True(t, test.wantGetErrKind.Match(haveGetErr), "Index %d => %s", i, haveGetErr)
			continue
		} else {
			assert.NoError(t, haveGetErr, "Index %d", i)
		}

		if test.wantSetErrKind.Empty() {
			assert.Exactly(t, test.val, haveVal, "Index %d", i)
		}
	}
}
