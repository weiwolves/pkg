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

package scope_test

import (
	"context"
	"testing"

	"github.com/weiwolves/pkg/store/scope"
	"github.com/weiwolves/pkg/util/assert"
)

func TestFromContext(t *testing.T) {
	tests := []struct {
		stID uint32
		wID  uint32
		want bool
	}{
		{0, 0, true},
		{1, 1, true},
		{11, 11, true},
		{0, 1, true},
		{1, 0, true},
	}
	for i, test := range tests {
		ctx := scope.WithContext(context.TODO(), test.wID, test.stID)
		haveWebsiteID, haveStoreID, haveOK := scope.FromContext(ctx)
		if have, want := haveOK, test.want; have != want {
			t.Errorf("(%d) Have: %v Want: %v", i, have, want)
		}
		if have, want := haveStoreID, test.stID; have != want {
			t.Errorf("Current Have: %v Want: %v", have, want)
		}
		if have, want := haveWebsiteID, test.wID; have != want {
			t.Errorf("Parent Have: %v Want: %v", have, want)
		}
	}
	w, st, ok := scope.FromContext(context.Background())
	assert.Exactly(t, uint32(0), st)
	assert.Exactly(t, uint32(0), w)
	assert.False(t, ok)
}
