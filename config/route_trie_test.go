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

// The MIT License (MIT)
//
// Copyright (c) 2014 Dalton Hubble
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package config

import (
	"testing"

	"github.com/weiwolves/pkg/store/scope"
	"github.com/weiwolves/pkg/util/assert"
)

func TestBuildTrieKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		key  string
		scp  scope.TypeID
		want string
	}{
		{"aa/bb/cc", 0, "/aa/bb/cc"},
		{"/aa/bb/cc", 0, "/aa/bb/cc"},
		{"a1/bb/cc", scope.DefaultTypeID, "/a1/bb/cc"},
		{"/a1a/bb/cc", scope.DefaultTypeID, "/a1a/bb/cc"},
		{"a2/bb/cc", scope.Website.WithID(3), "/a2/bb/cc/websites/3"},
		{"/a2a/bb/cc", scope.Website.WithID(3), "/a2a/bb/cc/websites/3"},
		{"a3/bb/cc", scope.Group.WithID(3), "/a3/bb/cc"},
		{"/a3a/bb/cc", scope.Group.WithID(3), "/a3a/bb/cc"},
		{"a4/bb/cc", scope.Store.WithID(3), "/a4/bb/cc/stores/3"},
		{"/a5/bb/cc", scope.Store.WithID(31), "/a5/bb/cc/stores/31"},
	}
	for _, test := range tests {
		assert.Exactly(t, test.want, buildTrieKey(test.key, test.scp), "Path: %q", test.key)
	}
}

var benchmarkBuildTrieKey string

func BenchmarkBuildTrieKey(b *testing.B) {
	const route = "customer/create_account/generate_human_friendly_id"
	const route2 = "/" + route
	b.Run("scope=0 with prefix", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			benchmarkBuildTrieKey = buildTrieKey(route2, 0)
		}
		if want := "/customer/create_account/generate_human_friendly_id"; benchmarkBuildTrieKey != want {
			b.Fatalf("want %q\nhave %q", want, benchmarkBuildTrieKey)
		}
	})
	b.Run("scope=0 without", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			benchmarkBuildTrieKey = buildTrieKey(route, 0)
		}
		if want := "/customer/create_account/generate_human_friendly_id"; benchmarkBuildTrieKey != want {
			b.Fatalf("want %q\nhave %q", want, benchmarkBuildTrieKey)
		}
	})
	b.Run("Store=10", func(b *testing.B) {
		scp := scope.Store.WithID(10)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchmarkBuildTrieKey = buildTrieKey(route, scp)
		}
		if want := "/customer/create_account/generate_human_friendly_id/stores/10"; benchmarkBuildTrieKey != want {
			b.Fatalf("want %q\nhave %q", want, benchmarkBuildTrieKey)
		}
	})
	b.Run("Website=20", func(b *testing.B) {
		scp := scope.Website.WithID(20)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchmarkBuildTrieKey = buildTrieKey(route, scp)
		}
		if want := "/customer/create_account/generate_human_friendly_id/websites/20"; benchmarkBuildTrieKey != want {
			b.Fatalf("want %q\nhave %q", want, benchmarkBuildTrieKey)
		}
	})
}

type noopObserver int

func (noopObserver) Observe(p Path, rawData []byte, found bool) (rawData2 []byte, err error) {
	return rawData, nil
}

var (
	noopCB0 = new(noopObserver)
	noopCB1 = new(noopObserver)
	noopCB2 = new(noopObserver)
)

func TestPathTrieNormal(t *testing.T) {
	trie := newTrieRoute()

	cases := []struct {
		key   string
		value Observer
		meta  FieldMeta
	}{
		{"fish", noopCB1, FieldMeta{}},
		{"/cat", noopCB2, FieldMeta{valid: true, Default: "cat", WriteScopePerm: scope.PermStore}},
		{"/dog", noopCB1, FieldMeta{valid: true, Default: "dog", WriteScopePerm: scope.PermStore}},
		{"/cats", noopCB2, FieldMeta{valid: true, Default: "cats", WriteScopePerm: scope.PermStore}},
		{"/caterpillar", noopCB2, FieldMeta{valid: true, Default: "caterpillar", WriteScopePerm: scope.PermWebsite}},
		{"/cat/gideon", noopCB2, FieldMeta{valid: true, Default: "cat/gideon", WriteScopePerm: scope.PermStore}},
		{"/cat/giddy", noopCB1, FieldMeta{valid: true, Default: "cat/giddy", WriteScopePerm: scope.PermDefault}},
	}

	// get missing keys
	for _, c := range cases {
		if value := trie.Get(c.key); value.valid {
			t.Errorf("expected key %s to be missing, found value %#v", c.key, value)
		}
	}

	// initial put
	for _, c := range cases {
		if isNew := trie.PutEvent(EventOnAfterGet, c.key, noopCB0); !isNew {
			t.Errorf("expected key %s to be missing", c.key)
		}
		fm := &FieldMeta{
			WriteScopePerm: c.meta.WriteScopePerm,
			DefaultValid:   true,
			Default:        c.meta.Default,
		}
		if isNew := trie.PutMeta(c.key, fm); isNew {
			t.Errorf("expected key %s to be missing", c.key)
		}
	}

	// subsequent put
	for _, c := range cases {
		if isNew := trie.PutEvent(EventOnAfterGet, c.key, c.value); isNew {
			t.Errorf("expected key %s to have a value already", c.key)
		}
	}

	// get
	for _, c := range cases {
		value := trie.Get(c.key)
		assert.Same(t, value.Events[EventOnAfterGet][1], c.value)
		assert.Exactly(t, c.meta.Default, value.Default)
		assert.Exactly(t, c.meta.WriteScopePerm, value.WriteScopePerm)
	}

	// delete, expect Delete to return true indicating a node was nil'd
	for _, c := range cases {
		if deleted := trie.Delete(c.key); !deleted {
			t.Errorf("expected key %s to be deleted", c.key)
		}
	}

	// delete cleaned all the way to the first character
	// expect Delete to return false bc no node existed to nil
	for _, c := range cases {
		if deleted := trie.Delete(string(c.key[0])); deleted {
			t.Errorf("expected key %s to be cleaned by delete", string(c.key[0]))
		}
	}

	// get deleted keys
	for _, c := range cases {
		if value := trie.Get(c.key); value.valid {
			t.Errorf("expected key %s to be deleted, got value %#v", c.key, value)
		}
	}
}

func TestPathTrieNilBehavior(t *testing.T) {
	trie := newTrieRoute()
	cases := []struct {
		key   string
		value Observer
	}{
		{"/cat", noopCB1},
		{"/catamaran", noopCB2},
		{"/caterpillar", nil},
	}
	expectNilValues := []string{"/", "/c", "/ca", "/caterpillar", "/other"}

	// initial put
	for _, c := range cases {
		if isNew := trie.PutEvent(EventOnAfterGet, c.key, c.value); !isNew {
			t.Errorf("expected key %s to be missing", c.key)
		}
	}

	// get nil
	for _, key := range expectNilValues {
		if value := trie.Get(key); len(value.Events[EventOnAfterGet]) != 0 {
			t.Errorf("expected key %s to have value nil, got %#v", key, value)
		}
	}
}

func TestPathTrieRoot(t *testing.T) {
	trie := newTrieRoute()

	if value := trie.Get(""); value.valid {
		t.Errorf("expected key '' to be missing, found value %#v", value)
	}
	if !trie.PutEvent(EventOnAfterGet, "", noopCB0) {
		t.Error("expected key '' to be missing")
	}
	if trie.PutEvent(EventOnAfterGet, "", noopCB1) {
		t.Error("expected key '' to have a value already")
	}
	value := trie.Get("")
	assert.Same(t, value.Events[EventOnAfterGet][1], noopCB1)

	if !trie.Delete("") {
		t.Error("expected key '' to be deleted")
	}
	if value := trie.Get(""); value.valid {
		t.Errorf("expected key '' to be deleted, got value %#v", value)
	}
}

func TestPathTrieWalk(t *testing.T) {
	trie := newTrieRoute()
	table := map[string]Observer{
		"/fish":        noopCB0,
		"/cat":         noopCB1,
		"/dog":         noopCB1,
		"/cats":        noopCB2,
		"/caterpillar": noopCB2,
		"/notes":       noopCB1,
		"/notes/new":   noopCB2,
		"/notes/:id":   noopCB2,
	}
	// key -> times walked
	walked := make(map[string]int)
	for key := range table {
		walked[key] = 0
	}

	for key, value := range table {
		if isNew := trie.PutEvent(EventOnAfterGet, key, value); !isNew {
			t.Errorf("expected key %s to be missing", key)
		}
	}

	walker := func(key string, value FieldMeta) error {
		// value for each walked key is correct

		if value.Events[EventOnAfterGet] != nil {
			assert.Same(t, value.Events[EventOnAfterGet][0], table[key])
		}
		walked[key]++
		return nil
	}
	trie.Walk(walker)

	// each key/value walked exactly once
	for key, walkedCount := range walked {
		if walkedCount != 1 {
			t.Errorf("expected key %s to be walked exactly once, got %v", key, walkedCount)
		}
	}
}

// test splitting /path/keys/ into parts (e.g. /path, /keys, /)
func TestPathSegmenter(t *testing.T) {
	cases := []struct {
		key     string
		parts   []string
		indices []int // indexes to use as next start, in order
	}{
		{"", []string{""}, []int{-1}},
		{"/", []string{"/"}, []int{-1}},
		{"static_file", []string{"static_file"}, []int{-1}},
		{"/users/scott", []string{"/users", "/scott"}, []int{6, -1}},
		{"users/scott", []string{"users", "/scott"}, []int{5, -1}},
		{"/users/ramona/", []string{"/users", "/ramona", "/"}, []int{6, 13, -1}},
		{"users/ramona/", []string{"users", "/ramona", "/"}, []int{5, 12, -1}},
		{"//", []string{"/", "/"}, []int{1, -1}},
		{"/a/b/c", []string{"/a", "/b", "/c"}, []int{2, 4, -1}},
	}

	for _, c := range cases {
		partNum := 0
		for prefix, i := segmentRoute(c.key, 0); ; prefix, i = segmentRoute(c.key, i) {
			if prefix != c.parts[partNum] {
				t.Errorf("expected part %d of key '%s' to be '%s', got '%s'", partNum, c.key, c.parts[partNum], prefix)
			}
			if i != c.indices[partNum] {
				t.Errorf("in iteration %d, expected next index of key '%s' to be '%d', got '%d'", partNum, c.key, c.indices[partNum], i)
			}
			partNum++
			if i == -1 {
				break
			}
		}
		if partNum != len(c.parts) {
			t.Errorf("expected '%s' to have %d parts, got %d", c.key, len(c.parts), partNum)
		}
	}
}

func TestPathSegmenterEdgeCases(t *testing.T) {
	cases := []struct {
		path      string
		start     int
		segment   string
		nextIndex int
	}{
		{"", 0, "", -1},
		{"", 10, "", -1},
		{"/", 0, "/", -1},
		{"/", 10, "", -1},
		{"/", -10, "", -1},
		{"/", 1, "", -1},
		{"//", 0, "/", 1},
		{"//", 1, "/", -1},
		{"//", 2, "", -1},
		{" /", 0, " ", 1},
		{" /", 1, "/", -1},
	}

	for _, c := range cases {
		segment, nextIndex := segmentRoute(c.path, c.start)
		if segment != c.segment {
			t.Errorf("expected segment %s starting at %d in path %s, got %s", c.segment, c.start, c.path, segment)
		}
		if nextIndex != c.nextIndex {
			t.Errorf("expected nextIndex %d starting at %d in path %s, got %d", c.nextIndex, c.start, c.path, nextIndex)
		}
	}
}
