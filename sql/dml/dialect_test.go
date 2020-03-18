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

package dml

import (
	"context"
	"testing"

	"github.com/weiwolves/pkg/storage/null"
	"github.com/weiwolves/pkg/util/naughtystrings"
)

// They both must be kept in sync
var (
	_ null.Dialecter = (*mysqlDialect)(nil)
	_ dialecter      = (*mysqlDialect)(nil)
)

func TestEscapeWith_NaughtyStrings(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)
	defer testCloser(t, s)

	sel := s.SelectFrom("dml_people").AddColumns("id", "name", "email").OrderBy("id")

	for _, nstr := range naughtystrings.Unencoded() {
		var people dmlPersons
		sel.Where(Column("name").Str(nstr))
		count, err := sel.WithDBR().Load(context.TODO(), &people)
		if err != nil {
			t.Fatalf("DB Error: %+v\n\nWith string: %q", err, nstr)
		}
		if count > 0 {
			t.Fatalf("Should not find any rows, but got %d for string: %q", count, nstr)
		}

		sel.Wheres = sel.Wheres[:0]
	}
}
