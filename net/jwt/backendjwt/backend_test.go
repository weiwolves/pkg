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
	"github.com/weiwolves/pkg/config/cfgmodel"
	"github.com/weiwolves/pkg/net/jwt/backendjwt"
)

// backend overall backend models for all tests
var backend *backendjwt.Configuration

var _ cfgmodel.Encrypter = (*noopCrypt)(nil)
var _ cfgmodel.Decrypter = (*noopCrypt)(nil)

type noopCrypt struct{}

func (noopCrypt) Encrypt(s []byte) ([]byte, error) {
	return s, nil
}

func (noopCrypt) Decrypt(s []byte) ([]byte, error) {
	return s, nil
}

// this would belong into the test suit setup
func init() {

	cfgStruct, err := backendjwt.NewConfigStructure()
	if err != nil {
		panic(err)
	}

	backend = backendjwt.New(cfgStruct)
	backend.HmacPassword.Encrypter = noopCrypt{}
	backend.HmacPassword.Decrypter = noopCrypt{}
}
