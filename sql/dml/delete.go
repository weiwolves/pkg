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
	"context"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// Delete contains the clauses for a DELETE statement.
//
// InnoDB Tables: If you are deleting many rows from a large table, you may
// exceed the lock table size for an InnoDB table. To avoid this problem, or
// simply to minimize the time that the table remains locked, the following
// strategy (which does not use DELETE at all) might be helpful:
//
// Select the rows not to be deleted into an empty table that has the same
// structure as the original table:
//	INSERT INTO t_copy SELECT * FROM t WHERE ... ;
// Use RENAME TABLE to atomically move the original table out of the way and
// rename the copy to the original name:
//	RENAME TABLE t TO t_old, t_copy TO t;
// Drop the original table:
//	DROP TABLE t_old;
// No other sessions can access the tables involved while RENAME TABLE executes,
// so the rename operation is not subject to concurrency problems.
// TODO(CyS) add DELETE ... JOIN ... statement SQLStmtDeleteJoin
type Delete struct {
	BuilderBase
	BuilderConditional
	// Listeners allows to dispatch certain functions in different
	// situations.
	Listeners ListenersDelete
	// Select used in case a DELETE statement should be build from a SELECT.
	Select *Select
}

// NewDelete creates a new Delete object.
func NewDelete(from string) *Delete {
	return &Delete{
		BuilderBase: BuilderBase{
			Table: MakeIdentifier(from),
		},
		BuilderConditional: BuilderConditional{
			Wheres: make(Conditions, 0, 2),
		},
	}
}

func newDeleteFrom(db QueryExecPreparer, idFn uniqueIDFn, l log.Logger, from string) *Delete {
	id := idFn()
	if l != nil {
		l = l.With(log.String("delete_id", id), log.String("table", from))
	}
	return &Delete{
		BuilderBase: BuilderBase{
			builderCommon: builderCommon{
				id:  id,
				Log: l,
				DB:  db,
			},
			Table: MakeIdentifier(from),
		},
		BuilderConditional: BuilderConditional{
			Wheres: make(Conditions, 0, 2),
		},
	}
}

// DeleteFrom creates a new Delete for the given table
func (c *ConnPool) DeleteFrom(from string) *Delete {
	return newDeleteFrom(c.DB, c.makeUniqueID, c.Log, from)
}

// DeleteFrom creates a new Delete for the given table
// in the context for a single database connection.
func (c *Conn) DeleteFrom(from string) *Delete {
	return newDeleteFrom(c.DB, c.makeUniqueID, c.Log, from)
}

// DeleteFrom creates a new Delete for the given table
// in the context for a transaction
func (tx *Tx) DeleteFrom(from string) *Delete {
	return newDeleteFrom(tx.DB, tx.makeUniqueID, tx.Log, from)
}

// FromSelect derives a DELETE with a SELECT statement. It supports the multi
// table syntax. Tables argument can be the actual table names or their aliases.
func (b *Delete) FromSelect(s *Select, tables ...string) *Delete {
	panic("TODO implement")
	//DELETE [LOW_PRIORITY] [QUICK] [IGNORE]
	//tbl_name[.*] [, tbl_name[.*]] ...
	//FROM table_references
	//[WHERE where_condition]
	//
	//DELETE [LOW_PRIORITY] [QUICK] [IGNORE]
	//FROM tbl_name[.*] [, tbl_name[.*]] ...
	//USING table_references
	//[WHERE where_condition]
	return b
}

// Alias sets an alias for the table name.
func (b *Delete) Alias(alias string) *Delete {
	b.Table.Aliased = alias
	return b
}

// WithDB sets the database query object. DB can be either a *sql.DB (connection
// pool), a *sql.Conn (a single dedicated database session) or a *sql.Tx (an
// in-progress database transaction).
func (b *Delete) WithDB(db QueryExecPreparer) *Delete {
	b.DB = db
	return b
}

// Unsafe see BuilderBase.IsUnsafe which weakens security when building the SQL
// string. This function must be called before calling any other function.
func (b *Delete) Unsafe() *Delete {
	b.IsUnsafe = true
	return b
}

// Where appends a WHERE clause to the statement whereSQLOrMap can be a string
// or map. If it'ab a string, args wil replaces any places holders.
func (b *Delete) Where(wf ...*Condition) *Delete {
	b.Wheres = append(b.Wheres, wf...)
	return b
}

// OrderBy appends columns to the ORDER BY statement for ascending sorting. A
// column gets always quoted if it is a valid identifier otherwise it will be
// treated as an expression. When you use ORDER BY or GROUP BY to sort a column
// in a DELETE, the server sorts arguments using only the initial number of
// bytes indicated by the max_sort_length system variable.
func (b *Delete) OrderBy(columns ...string) *Delete {
	b.OrderBys = b.OrderBys.AppendColumns(b.IsUnsafe, columns...)
	return b
}

// OrderByDesc appends columns to the ORDER BY statement for descending sorting.
// A column gets always quoted if it is a valid identifier otherwise it will be
// treated as an expression. When you use ORDER BY or GROUP BY to sort a column
// in a DELETE, the server sorts arguments using only the initial number of
// bytes indicated by the max_sort_length system variable.
func (b *Delete) OrderByDesc(columns ...string) *Delete {
	b.OrderBys = b.OrderBys.AppendColumns(b.IsUnsafe, columns...).applySort(len(columns), sortDescending)
	return b
}

// Limit sets a LIMIT clause for the statement; overrides any existing LIMIT
func (b *Delete) Limit(limit uint64) *Delete {
	b.LimitCount = limit
	b.LimitValid = true
	return b
}

// WithArgs returns a new type to support multiple executions of the underlying
// SQL statement and reuse of memory allocations for the arguments. WithArgs
// builds the SQL string and sets the optional raw interfaced arguments for the
// later execution. It copies the underlying connection and settings from the
// current DML type (Delete, Insert, Select, Update, Union, With, etc.). The
// query executor can still be overwritten. Interpolation does not support the
// raw interfaces.
func (b *Delete) WithArgs(args ...interface{}) *Arguments {
	b.source = dmlSourceDelete
	return b.withArgs(b, args...)
}

// ToSQL generates the SQL string and might caches it internally, if not
// disabled. The returned interface slice is always nil.
func (b *Delete) ToSQL() (string, []interface{}, error) {
	b.source = dmlSourceDelete
	rawSQL, err := b.buildToSQL(b)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}
	return string(rawSQL), nil, nil
}

func (b *Delete) writeBuildCache(sql []byte, qualifiedColumns []string) {
	b.rwmu.Lock()
	b.qualifiedColumns = qualifiedColumns
	if !b.IsBuildCacheDisabled {
		b.BuilderConditional = BuilderConditional{}
		b.cachedSQL = sql
	}
	b.rwmu.Unlock()
}

// DisableBuildCache if enabled it does not cache the SQL string as a final
// rendered byte slice. Allows you to rebuild the query with different
// statements.
func (b *Delete) DisableBuildCache() *Delete {
	b.IsBuildCacheDisabled = true
	return b
}

// ToSQL serialized the Delete to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Delete) toSQL(w *bytes.Buffer, placeHolders []string) ([]string, error) {
	b.defaultQualifier = b.Table.qualifier()

	if err := b.Listeners.dispatch(OnBeforeToSQL, b); err != nil {
		return nil, errors.WithStack(err)
	}

	if b.RawFullSQL != "" {
		_, err := w.WriteString(b.RawFullSQL)
		return nil, err
	}

	if b.Table.Name == "" {
		return nil, errors.Empty.Newf("[dml] Delete: Table is missing")
	}

	w.WriteString("DELETE ")
	writeStmtID(w, b.id)
	w.WriteString("FROM ")
	placeHolders, err := b.Table.writeQuoted(w, placeHolders)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// TODO(CyS) add SQLStmtDeleteJoin
	placeHolders, err = b.Wheres.write(w, 'w', placeHolders)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sqlWriteOrderBy(w, b.OrderBys, false)
	sqlWriteLimitOffset(w, b.LimitValid, b.LimitCount, false, 0)

	return placeHolders, nil
}

// Prepare executes the statement represented by the Delete to create a prepared
// statement. It returns a custom statement type or an error if there was one.
// Provided arguments or records in the Delete are getting ignored. The provided
// context is used for the preparation of the statement, not for the execution
// of the statement. If debug mode for logging has been enabled it logs the
// duration taken and the SQL string. The returned Stmter is not safe for
// concurrent use, despite the underlying *sql.Stmt is.
func (b *Delete) Prepare(ctx context.Context) (*Stmt, error) {
	return b.prepare(ctx, b.DB, b, dmlSourceDelete)
}
