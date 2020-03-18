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

package bufferpool_test

import (
	"sync"
	"testing"

	"github.com/weiwolves/pkg/util/assert"
	"github.com/weiwolves/pkg/util/bufferpool"
)

var singleBuf = bufferpool.New(4096)

func TestBufferPoolSize(t *testing.T) {
	t.Parallel()

	const iterations = 10
	var wg sync.WaitGroup
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func(wg *sync.WaitGroup) {
			b := singleBuf.Get()
			defer func() { singleBuf.Put(b); wg.Done() }()

			assert.Exactly(t, 4096, b.Cap())
			assert.Exactly(t, 0, b.Len())
		}(&wg)
	}
	wg.Wait()
}
