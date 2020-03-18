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

package store_test

import (
	"testing"

	"github.com/corestoreio/errors"
	"github.com/weiwolves/pkg/store"
	"github.com/weiwolves/pkg/util/assert"
)

func TestValidateStoreCode(t *testing.T) {
	tests := []struct {
		have       string
		wantErrBhf errors.Kind
	}{
		{"@de", errors.NotValid},
		{" de", errors.NotValid},
		{"de", errors.NoKind},
		{"DE", errors.NoKind},
		{"deCH09_", errors.NoKind},
		{"_de", errors.NotValid},
		{"", errors.NotValid},
		{"\U0001f41c", errors.NotValid},
		{"au_en", errors.NoKind},
		{"au-fr", errors.NotValid},
		{"Hello GoLang", errors.NotValid},
		{"Hello€GoLang", errors.NotValid},
		{"HelloGoLdhashdfkjahdjfhaskjdfhuiwehfiawehfuahweldsnjkasfkjkwejqwehqang", errors.NotValid},
	}
	for i, test := range tests {
		haveErr := store.CodeIsValid(test.have)
		if !test.wantErrBhf.Empty() {
			assert.True(t, test.wantErrBhf.Match(haveErr), "Index %d => %s", i, haveErr)
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
		}
	}
}
