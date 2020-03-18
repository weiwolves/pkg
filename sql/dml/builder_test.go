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
	"bytes"
	"fmt"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/weiwolves/pkg/util/assert"
)

// check if the types implement the interfaces

var (
	_ fmt.Stringer = (*Delete)(nil)
	_ fmt.Stringer = (*Insert)(nil)
	_ fmt.Stringer = (*Update)(nil)
	_ fmt.Stringer = (*Select)(nil)
	_ fmt.Stringer = (*Union)(nil)
	_ fmt.Stringer = (*With)(nil)
	_ fmt.Stringer = (*Show)(nil)
)

var (
	_ QueryBuilder = (*Select)(nil)
	_ QueryBuilder = (*Delete)(nil)
	_ QueryBuilder = (*Update)(nil)
	_ QueryBuilder = (*Insert)(nil)
	_ QueryBuilder = (*Show)(nil)
	_ QueryBuilder = (*Union)(nil)
	_ QueryBuilder = (*With)(nil)
)

func TestSqlObjToString(t *testing.T) {
	t.Parallel()
	t.Run("error", func(t *testing.T) {
		s := sqlObjToString("", errors.Aborted.Newf("Query aborted"))
		assert.Contains(t, s, "[dml] String Error: Query aborted\n")
	})
	t.Run("DELETE", func(t *testing.T) {
		b := NewDelete("tableX").Where(Column("columnA").Greater().Int64(2))
		assert.Exactly(t, "DELETE FROM `tableX` WHERE (`columnA` > 2)", b.String())
	})
	t.Run("INSERT", func(t *testing.T) {
		b := NewInsert("tableX").AddColumns("columnA", "columnB").BuildValues()
		// keep the place holder for columnA,columnB because we're not using interpolation
		assert.Exactly(t, "INSERT INTO `tableX` (`columnA`,`columnB`) VALUES (?,?)", b.String())
	})
	t.Run("SELECT", func(t *testing.T) {
		b := NewSelect("columnA").FromAlias("tableX", "X").Where(Column("columnA").LessOrEqual().Float64(2.4))
		assert.Exactly(t, "SELECT `columnA` FROM `tableX` AS `X` WHERE (`columnA` <= 2.4)", b.String())
	})
	t.Run("UPDATE", func(t *testing.T) {
		b := NewUpdate("tableX").AddClauses(
			Column("columnA").Int64(4),
		).Where(Column("columnB").Between().Ints(5, 7))
		// keep the place holder for columnA because we're not using interpolation
		assert.Exactly(t, "UPDATE `tableX` SET `columnA`=4 WHERE (`columnB` BETWEEN 5 AND 7)", b.String())
	})
	t.Run("WITH", func(t *testing.T) {
		b := NewWith(WithCTE{Name: "sel", Select: NewSelect().Unsafe().AddColumns("1")}).
			Select(NewSelect().Star().From("sel"))
		assert.Exactly(t, "WITH `sel` AS (SELECT 1)\nSELECT * FROM `sel`", b.String())
	})
	t.Run("SHOW", func(t *testing.T) {
		b := NewShow().Variable().Like().String()
		assert.Exactly(t, "SHOW VARIABLES LIKE ?", b)
	})
}

func TestWriteInsertPlaceholders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		rowCount    uint
		columnCount uint
		want        string
	}{
		{1, 0, ""},
		{1, 1, "(?)"},
		{2, 1, "(?),(?)"},
		{1, 2, "(?,?)"},
		{1, 3, "(?,?,?)"},
		{1, 10, "(?,?,?,?,?,?,?,?,?,?)"},
		{1, 11, "(?,?,?,?,?,?,?,?,?,?,?)"},
		{1, 12, "(?,?,?,?,?,?,?,?,?,?,?,?)"},
		{1, 13, "(?,?,?,?,?,?,?,?,?,?,?,?,?)"},
		{1, 14, "(?,?,?,?,?,?,?,?,?,?,?,?,?,?)"},
		{1, 15, "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"},
		{1, 60, "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"},
		{1, 11, "(?,?,?,?,?,?,?,?,?,?,?)"},
		{2, 3, "(?,?,?),(?,?,?)"},
		{4, 3, "(?,?,?),(?,?,?),(?,?,?),(?,?,?)"},
	}

	for i, test := range tests {
		var buf bytes.Buffer
		writeTuplePlaceholders(&buf, test.rowCount, test.columnCount)
		assert.Exactly(t, test.want, buf.String(), "Index %d", i)
	}
}
