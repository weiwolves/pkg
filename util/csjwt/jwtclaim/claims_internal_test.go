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
	"testing"

	"github.com/weiwolves/pkg/util/assert"
	"github.com/weiwolves/pkg/util/csjwt"
)

func TestKeys(t *testing.T) {
	tests := []struct {
		claim    csjwt.Claimer
		wantKeys []string
	}{
		{&Standard{}, allKeys[:7]},
		{NewStore(), allKeys[:]},
		{Map{
			"k1": "v1", // only one key because of random map output >8-)
		}, []string{"k1"}},
	}
	for i, test := range tests {
		assert.Exactly(t, test.wantKeys, test.claim.Keys(), "Index %d", i)
	}
}
