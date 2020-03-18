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

package backendjwt_test

import (
	"testing"

	"github.com/weiwolves/pkg/config/cfgmock"
	"github.com/weiwolves/pkg/net/jwt/backendjwt"
	"github.com/weiwolves/pkg/store/scope"
	"github.com/corestoreio/errors"
	"github.com/weiwolves/pkg/util/assert"
)

func TestNewConfigSigningMethodGetDefaultPathError(t *testing.T) {
	ccModel := backendjwt.NewConfigSigningMethod("a/x/c")
	cr := cfgmock.NewService()
	sm, err := ccModel.Get(cr.NewScoped(1, 1))
	assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
	assert.Nil(t, sm)
	assert.Exactly(t, scope.TypeIDs{}, cr.AllInvocations().ScopeIDs())
}

func TestNewConfigSigningMethodGetPathError(t *testing.T) {
	ccModel := backendjwt.NewConfigSigningMethod("a//c")
	cr := cfgmock.NewService()
	sm, err := ccModel.Get(cr.NewScoped(0, 0))
	assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
	assert.Nil(t, sm)
	assert.Exactly(t, scope.TypeIDs{}, cr.AllInvocations().ScopeIDs())
}
