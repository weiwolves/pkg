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

// +build csall json yaml

package storage_test

import (
	"io"
	"io/ioutil"
	"testing"

	"github.com/weiwolves/pkg/config"
	"github.com/weiwolves/pkg/config/storage"
	"github.com/weiwolves/pkg/util/assert"
)

func TestWithGlob_Placeholder(t *testing.T) {
	cfgSrv, err := config.NewService(
		storage.NewMap(), config.Options{
			EnvName: "path",
		},
	)
	assert.NoError(t, err)

	shouldGetCalled := false
	storage.WithGlob("testdata/malformed_"+config.EnvNamePlaceHolder+".yaml")(cfgSrv, func(_ config.Setter, r io.Reader) error {
		data, err := ioutil.ReadAll(r)
		assert.NoError(t, err)

		assert.Exactly(t, `vendorbarenvironment:
  default:
    0: 192.168.33.1
`, string(data))

		shouldGetCalled = true
		return nil
	})

	assert.True(t, shouldGetCalled, "cb of WithGlob should have been called")
}
