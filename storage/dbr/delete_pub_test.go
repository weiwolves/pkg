// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dbr_test

import (
	"bytes"
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log/logw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDelete_Prepare(t *testing.T) {
	t.Parallel()

	t.Run("ToSQL Error", func(t *testing.T) {
		compareToSQL(t, dbr.NewDelete("").Where(dbr.Column("a").Int64(1)), errors.IsEmpty, "", "")
	})

	t.Run("Prepare Error", func(t *testing.T) {
		d := &dbr.Delete{
			BuilderBase: dbr.BuilderBase{
				Table: dbr.MakeIdentifier("table"),
			},
		}
		d.DB = dbMock{
			error: errors.NewAlreadyClosedf("Who closed myself?"),
		}
		d.Where(dbr.Column("a").Int(1))
		stmt, err := d.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
	})

	t.Run("ExecArgs One Row", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(cstesting.SQLMockQuoteMeta("DELETE FROM `customer_entity` WHERE (`email` = ?) AND (`group_id` = ?)"))
		prep.ExpectExec().WithArgs("a@b.c", 33).WillReturnResult(sqlmock.NewResult(0, 1))
		prep.ExpectExec().WithArgs("x@y.z", 44).WillReturnResult(sqlmock.NewResult(0, 2))

		stmt, err := dbr.NewDelete("customer_entity").
			Where(dbr.Column("email").PlaceHolder(), dbr.Column("group_id").PlaceHolder()).
			WithDB(dbc.DB).
			Prepare(context.TODO())
		require.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			require.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		tests := []struct {
			email   string
			groupID int
			affRows int64
		}{
			{"a@b.c", 33, 1},
			{"x@y.z", 44, 2},
		}

		args := dbr.MakeArgs(3)
		for i, test := range tests {
			args = args[:0]

			res, err := stmt.WithArguments(args.Str(test.email).Int(test.groupID)).Exec(context.TODO())
			if err != nil {
				t.Fatalf("Index %d => %+v", i, err)
			}
			ra, err := res.RowsAffected()
			if err != nil {
				t.Fatalf("Result index %d with error: %s", i, err)
			}
			assert.Exactly(t, test.affRows, ra, "Index %d has different RowsAffected", i)
		}
	})

	t.Run("ExecRecord One Row", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(cstesting.SQLMockQuoteMeta("DELETE FROM `dbr_person` WHERE (`name` = ?) AND (`email` = ?)"))
		prep.ExpectExec().WithArgs("Peter Gopher", "peter@gopher.go").WillReturnResult(sqlmock.NewResult(0, 4))
		prep.ExpectExec().WithArgs("John Doe", "john@doe.go").WillReturnResult(sqlmock.NewResult(0, 5))

		stmt, err := dbr.NewDelete("dbr_person").
			Where(dbr.Column("name").PlaceHolder(), dbr.Column("email").PlaceHolder()).
			WithDB(dbc.DB).
			Prepare(context.TODO())
		require.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			require.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		tests := []struct {
			name     string
			email    string
			insertID int64
		}{
			{"Peter Gopher", "peter@gopher.go", 4},
			{"John Doe", "john@doe.go", 5},
		}

		for i, test := range tests {

			p := &dbrPerson{
				Name:  test.name,
				Email: dbr.MakeNullString(test.email),
			}

			res, err := stmt.WithRecords(dbr.Qualify("", p)).Exec(context.TODO())
			if err != nil {
				t.Fatalf("Index %d => %+v", i, err)
			}
			lid, err := res.RowsAffected()
			if err != nil {
				t.Fatalf("Result index %d with error: %s", i, err)
			}
			assert.Exactly(t, test.insertID, lid, "Index %d has different RowsAffected", i)
		}
	})

	t.Run("ExecContext", func(t *testing.T) {
		dbc, dbMock := cstesting.MockDB(t)
		defer cstesting.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(cstesting.SQLMockQuoteMeta("DELETE FROM `dbr_person` WHERE (`name` = ?) AND (`email` = ?)"))
		prep.ExpectExec().WithArgs("Peter Gopher", "peter@gopher.go").WillReturnResult(sqlmock.NewResult(0, 4))

		stmt, err := dbr.NewDelete("dbr_person").
			Where(dbr.Column("name").PlaceHolder(), dbr.Column("email").PlaceHolder()).
			WithDB(dbc.DB).
			Prepare(context.TODO())
		require.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			require.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		res, err := stmt.Exec(context.TODO(), "Peter Gopher", "peter@gopher.go")
		require.NoError(t, err, "failed to execute ExecContext")

		lid, err := res.RowsAffected()
		if err != nil {
			t.Fatal(err)
		}
		assert.Exactly(t, int64(4), lid, "Different RowsAffected")
	})
}

func TestDelete_WithLogger(t *testing.T) {
	uniID := new(int32)
	rConn := createRealSession(t)
	defer cstesting.Close(t, rConn)

	var uniqueIDFunc = func() string {
		return fmt.Sprintf("UNIQUEID%02d", atomic.AddInt32(uniID, 1))
	}

	buf := new(bytes.Buffer)
	lg := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithWriter(buf),
		logw.WithFlag(0), // no flags at all
	)
	require.NoError(t, rConn.Options(dbr.WithLogger(lg, uniqueIDFunc)))

	t.Run("ConnPool", func(t *testing.T) {
		d := rConn.DeleteFrom("dbr_people").Where(dbr.Column("id").GreaterOrEqual().Float64(34.56))

		t.Run("Exec", func(t *testing.T) {
			defer func() {
				buf.Reset()
				d.IsInterpolate = false
			}()
			_, err := d.Interpolate().Exec(context.TODO())
			require.NoError(t, err)

			assert.Exactly(t, "DEBUG Exec conn_pool_id: \"UNIQUEID01\" delete_id: \"UNIQUEID02\" table: \"dbr_people\" duration: 0 sql: \"DELETE /*ID:UNIQUEID02*/ FROM `dbr_people` WHERE (`id` >= 34.56)\"\n",
				buf.String())
		})

		t.Run("Prepare", func(t *testing.T) {
			defer buf.Reset()
			stmt, err := d.Prepare(context.TODO())
			require.NoError(t, err)
			defer stmt.Close()

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQUEID01\" delete_id: \"UNIQUEID02\" table: \"dbr_people\" duration: 0 sql: \"DELETE /*ID:UNIQUEID02*/ FROM `dbr_people` WHERE (`id` >= ?)\"\n",
				buf.String())
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()
			tx, err := rConn.BeginTx(context.TODO(), nil)
			require.NoError(t, err)
			require.NoError(t, tx.Wrap(func() error {
				_, err := tx.DeleteFrom("dbr_people").Where(dbr.Column("id").GreaterOrEqual().Float64(36.56)).Interpolate().Exec(context.TODO())
				if err != nil {
					return err
				}

				assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQUEID01\" tx_id: \"UNIQUEID03\"\nDEBUG Exec conn_pool_id: \"UNIQUEID01\" tx_id: \"UNIQUEID03\" delete_id: \"UNIQUEID04\" table: \"dbr_people\" duration: 0 sql: \"DELETE /*ID:UNIQUEID04*/ FROM `dbr_people` WHERE (`id` >= 36.56)\"\n",
					buf.String())

				return nil
			}))
		})
	})

	t.Run("Conn", func(t *testing.T) {
		conn, err := rConn.Conn(context.TODO())
		require.NoError(t, err)

		d := conn.DeleteFrom("dbr_people").Where(dbr.Column("id").GreaterOrEqual().Float64(39.56))
		t.Run("Exec", func(t *testing.T) {
			defer func() {
				buf.Reset()
				d.IsInterpolate = false
			}()

			_, err := d.Interpolate().Exec(context.TODO())
			require.NoError(t, err)

			assert.Exactly(t, "DEBUG Exec conn_pool_id: \"UNIQUEID01\" conn_id: \"UNIQUEID05\" delete_id: \"UNIQUEID06\" table: \"dbr_people\" duration: 0 sql: \"DELETE /*ID:UNIQUEID06*/ FROM `dbr_people` WHERE (`id` >= 39.56)\"\n",
				buf.String())
		})

		t.Run("Prepare", func(t *testing.T) {
			defer buf.Reset()

			stmt, err := d.Prepare(context.TODO())
			require.NoError(t, err)
			defer stmt.Close()

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQUEID01\" conn_id: \"UNIQUEID05\" delete_id: \"UNIQUEID06\" table: \"dbr_people\" duration: 0 sql: \"DELETE /*ID:UNIQUEID06*/ FROM `dbr_people` WHERE (`id` >= ?)\"\n",
				buf.String())
		})

		t.Run("Prepare Exec", func(t *testing.T) {
			defer buf.Reset()

			stmt, err := d.Prepare(context.TODO())
			require.NoError(t, err)
			defer stmt.Close()

			_, err = stmt.Exec(context.TODO(), 41.57)
			require.NoError(t, err)

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQUEID01\" conn_id: \"UNIQUEID05\" delete_id: \"UNIQUEID06\" table: \"dbr_people\" duration: 0 sql: \"DELETE /*ID:UNIQUEID06*/ FROM `dbr_people` WHERE (`id` >= ?)\"\nDEBUG Exec conn_pool_id: \"UNIQUEID01\" conn_id: \"UNIQUEID05\" delete_id: \"UNIQUEID06\" table: \"dbr_people\" duration: 0 arg_len: 1\n",
				buf.String())
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()
			tx, err := conn.BeginTx(context.TODO(), nil)
			require.NoError(t, err)
			require.NoError(t, tx.Wrap(func() error {
				_, err := tx.DeleteFrom("dbr_people").Where(dbr.Column("id").GreaterOrEqual().Float64(37.56)).Interpolate().Exec(context.TODO())
				return err
			}))

			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQUEID01\" conn_id: \"UNIQUEID05\" tx_id: \"UNIQUEID07\"\nDEBUG Exec conn_pool_id: \"UNIQUEID01\" conn_id: \"UNIQUEID05\" tx_id: \"UNIQUEID07\" delete_id: \"UNIQUEID08\" table: \"dbr_people\" duration: 0 sql: \"DELETE /*ID:UNIQUEID08*/ FROM `dbr_people` WHERE (`id` >= 37.56)\"\nDEBUG Commit conn_pool_id: \"UNIQUEID01\" conn_id: \"UNIQUEID05\" tx_id: \"UNIQUEID07\" duration: 0\n",
				buf.String())
		})

		t.Run("Tx Rollback", func(t *testing.T) {
			defer buf.Reset()
			tx, err := conn.BeginTx(context.TODO(), nil)
			require.NoError(t, err)
			require.Error(t, tx.Wrap(func() error {
				_, err := tx.DeleteFrom("dbr_people").Where(dbr.Column("id").GreaterOrEqual().PlaceHolder()).Interpolate().Exec(context.TODO())
				return err
			}))

			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQUEID01\" conn_id: \"UNIQUEID05\" tx_id: \"UNIQUEID09\"\nDEBUG Exec conn_pool_id: \"UNIQUEID01\" conn_id: \"UNIQUEID05\" tx_id: \"UNIQUEID09\" delete_id: \"UNIQUEID10\" table: \"dbr_people\" duration: 0 sql: \"DELETE /*ID:UNIQUEID10*/ FROM `dbr_people` WHERE (`id` >= ?)\"\nDEBUG Rollback conn_pool_id: \"UNIQUEID01\" conn_id: \"UNIQUEID05\" tx_id: \"UNIQUEID09\" duration: 0\n",
				buf.String())
		})
	})

}
