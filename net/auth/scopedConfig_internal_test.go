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

package auth

import (
	"testing"

	"github.com/corestoreio/errors"
	"github.com/weiwolves/pkg/util/assert"
)

func TestAuthProviders_Do_Len0(t *testing.T) {
	ap := providers{}
	err := ap.do(0, nil)
	assert.True(t, errors.NotImplemented.Match(err), "%+v", err)
}

func TestAuthTriggers_Do_Len0(t *testing.T) {
	assert.True(t, triggers{}.do(true, nil))
}
