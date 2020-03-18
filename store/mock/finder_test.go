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

package mock_test

import (
	"testing"

	"github.com/corestoreio/errors"
	"github.com/weiwolves/pkg/store/mock"
	"github.com/weiwolves/pkg/util/assert"
)

func TestFind(t *testing.T) {
	tests := []struct {
		f *mock.Find
	}{
		{
			mock.NewDefaultStoreID(2, 1, errors.Fatal.Newf("Whooops2"),
				mock.NewStoreIDbyCode(4, 3, errors.Fatal.Newf("Whooops1")),
			),
		},
		{
			mock.NewStoreIDbyCode(4, 3, errors.Fatal.Newf("Whooops1"),
				mock.NewDefaultStoreID(2, 1, errors.Fatal.Newf("Whooops2")),
			),
		},
	}
	for _, test := range tests {
		wID, sID, err := test.f.DefaultStoreID(0)
		assert.Exactly(t, uint32(1), sID)
		assert.Exactly(t, uint32(2), wID)
		assert.True(t, errors.Fatal.Match(err))
		assert.Exactly(t, 1, test.f.DefaultStoreIDInvoked())

		wID, sID, err = test.f.StoreIDbyCode(0, "x")
		assert.Exactly(t, uint32(3), sID)
		assert.Exactly(t, uint32(4), wID)
		assert.True(t, errors.Fatal.Match(err))
		assert.Exactly(t, 1, test.f.StoreIDbyCodeInvoked())
	}
}
