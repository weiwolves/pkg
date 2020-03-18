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

package dml_test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log/logw"
	"github.com/weiwolves/pkg/sql/dml"
	"github.com/weiwolves/pkg/sql/dmltest"
	"github.com/weiwolves/pkg/util/assert"
)

func TestWithLogger_Insert(t *testing.T) {
	uniID := new(int32)
	rConn := createRealSession(t)
	defer dmltest.Close(t, rConn)

	uniqueIDFunc := func() string {
		return fmt.Sprintf("UNIQ%02d", atomic.AddInt32(uniID, 4))
	}

	buf := new(bytes.Buffer)
	lg := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithWriter(buf),
		logw.WithFlag(0), // no flags at all
	)
	assert.NoError(t, rConn.Options(dml.WithLogger(lg, uniqueIDFunc)))

	t.Run("Conn1Pool", func(t *testing.T) {
		d := rConn.InsertInto("dml_people").Replace().AddColumns("email", "name")

		t.Run("Prepare", func(t *testing.T) {
			defer buf.Reset()
			stmt, err := d.BuildValues().Prepare(context.TODO())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ04\" insert_id: \"UNIQ08\" table: \"dml_people\" duration: 0 error: \"<nil>\" sql: \"REPLACE /*ID$UNIQ08*/ INTO `dml_people` (`email`,`name`) VALUES (?,?)\"\n",
				buf.String())
		})

		t.Run("Exec", func(t *testing.T) {
			defer buf.Reset()
			_, err := d.WithDBR().Interpolate().ExecContext(context.TODO(), "a@b.c", "John")
			assert.NoError(t, err)

			assert.Exactly(t, "DEBUG Exec conn_pool_id: \"UNIQ04\" insert_id: \"UNIQ08\" table: \"dml_people\" duration: 0 sql: \"REPLACE /*ID$UNIQ08*/ INTO `dml_people` (`email`,`name`) VALUES ('a@b.c','John')\" length_args: 0 length_raw_args: 2 source: \"i\" error: \"<nil>\"\n",
				buf.String())
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()
			err := rConn.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
				_, err := tx.InsertInto("dml_people").Replace().AddColumns("email", "name").WithDBR().ExecContext(context.TODO(), "a@b.c", "John")
				return err
			})
			assert.NoError(t, err)
			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ04\" tx_id: \"UNIQ12\"\nDEBUG Exec conn_pool_id: \"UNIQ04\" tx_id: \"UNIQ12\" insert_id: \"UNIQ16\" table: \"dml_people\" duration: 0 sql: \"REPLACE /*ID$UNIQ16*/ INTO `dml_people` (`email`,`name`) VALUES (?,?)\" length_args: 2 length_raw_args: 2 source: \"i\" error: \"<nil>\"\nDEBUG Commit conn_pool_id: \"UNIQ04\" tx_id: \"UNIQ12\" duration: 0\n",
				buf.String())
		})
	})

	t.Run("Conn2", func(t *testing.T) {
		conn, err := rConn.Conn(context.TODO())
		assert.NoError(t, err)

		oIns := conn.InsertInto("dml_people").Replace().AddColumns("email", "name").BuildValues()

		t.Run("Exec", func(t *testing.T) {
			defer buf.Reset()
			_, err := oIns.WithDBR().Interpolate().ExecContext(context.TODO(), "a@b.zeh", "J0hn")
			assert.NoError(t, err)

			assert.Exactly(t, "DEBUG Exec conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" insert_id: \"UNIQ24\" table: \"dml_people\" duration: 0 sql: \"REPLACE /*ID$UNIQ24*/ INTO `dml_people` (`email`,`name`) VALUES ('a@b.zeh','J0hn')\" length_args: 0 length_raw_args: 2 source: \"i\" error: \"<nil>\"\n",
				buf.String())
		})

		t.Run("Prepare", func(t *testing.T) {
			defer buf.Reset()
			stmt, err := oIns.BuildValues().Prepare(context.TODO())
			// oIns.IsBuildValues = false
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" insert_id: \"UNIQ24\" table: \"dml_people\" duration: 0 error: \"<nil>\" sql: \"REPLACE /*ID$UNIQ24*/ INTO `dml_people` (`email`,`name`) VALUES (?,?)\"\n",
				buf.String())
		})

		t.Run("Prepare Exec", func(t *testing.T) {
			defer buf.Reset()
			stmt, err := oIns.BuildValues().Prepare(context.TODO())
			// oIns.IsBuildValues = false
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)

			_, err = stmt.WithDBR().ExecContext(context.TODO(), "mail@e.de", "Hans")
			assert.NoError(t, err)

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" insert_id: \"UNIQ24\" table: \"dml_people\" duration: 0 error: \"<nil>\" sql: \"REPLACE /*ID$UNIQ24*/ INTO `dml_people` (`email`,`name`) VALUES (?,?)\"\nDEBUG Exec conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" insert_id: \"UNIQ24\" table: \"dml_people\" duration: 0 sql: \"\" length_args: 2 length_raw_args: 2 source: \"i\" error: \"<nil>\"\n",
				buf.String())
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()

			err := conn.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
				_, err := tx.InsertInto("dml_people").Replace().AddColumns("email", "name").WithDBR().ExecContext(context.TODO(), "a@b.c", "John")
				return err
			})
			assert.NoError(t, err)

			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" tx_id: \"UNIQ28\"\nDEBUG Exec conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" tx_id: \"UNIQ28\" insert_id: \"UNIQ32\" table: \"dml_people\" duration: 0 sql: \"REPLACE /*ID$UNIQ32*/ INTO `dml_people` (`email`,`name`) VALUES (?,?)\" length_args: 2 length_raw_args: 2 source: \"i\" error: \"<nil>\"\nDEBUG Commit conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" tx_id: \"UNIQ28\" duration: 0\n",
				buf.String())
		})

		t.Run("Tx Rollback", func(t *testing.T) {
			defer buf.Reset()

			assert.Error(t, conn.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
				_, err := tx.InsertInto("dml_people").Replace().AddColumns("email", "name").
					WithDBR().Interpolate().ExecContext(context.TODO(), "only one arg provided")
				return err
			}))

			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" tx_id: \"UNIQ36\"\nDEBUG Exec conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" tx_id: \"UNIQ36\" insert_id: \"UNIQ40\" table: \"dml_people\" duration: 0 sql: \"\" length_args: 0 length_raw_args: 1 source: \"i\" error: \"[dml] Interpolation failed: \\\"REPLACE /*ID$UNIQ40*/ INTO `dml_people` (`email`,`name`) VALUES (?,?)\\\": [dml] Number of place holders (2) vs number of arguments (1) do not match.\"\nDEBUG Rollback conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" tx_id: \"UNIQ36\" duration: 0\n",
				buf.String())
		})

		t.Run("Tx WithDBR", func(t *testing.T) {
			defer buf.Reset()

			// This INSERT statement does not have an ID as the others above.
			insA := dml.NewInsert("dml_people").Replace().AddColumns("email", "name").WithDBR()

			err := conn.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
				_, err := insA.WithTx(tx).ExecContext(context.TODO(), "a@b.c", "John")
				return err
			})
			assert.NoError(t, err)

			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" tx_id: \"UNIQ44\"\nDEBUG Exec conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" tx_id: \"UNIQ44\" duration: 0 sql: \"REPLACE INTO `dml_people` (`email`,`name`) VALUES (?,?)\" length_args: 2 length_raw_args: 2 source: \"i\" error: \"<nil>\"\nDEBUG Commit conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" tx_id: \"UNIQ44\" duration: 0\n",
				buf.String())
		})
	})
}

func TestWithLogger_Delete(t *testing.T) {
	uniID := new(int32)
	rConn := createRealSession(t)
	defer dmltest.Close(t, rConn)

	uniqueIDFunc := func() string {
		return fmt.Sprintf("UNIQUEID%02d", atomic.AddInt32(uniID, 1))
	}

	buf := new(bytes.Buffer)
	lg := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithWriter(buf),
		logw.WithFlag(0), // no flags at all
	)
	assert.NoError(t, rConn.Options(dml.WithLogger(lg, uniqueIDFunc)))

	t.Run("ConnPool", func(t *testing.T) {
		d := rConn.DeleteFrom("dml_people").Where(dml.Column("id").GreaterOrEqual().Float64(34.56))

		t.Run("Exec", func(t *testing.T) {
			defer func() {
				buf.Reset()
			}()
			_, err := d.WithDBR().Interpolate().ExecContext(context.TODO())
			assert.NoError(t, err)

			assert.Exactly(t, "DEBUG Exec conn_pool_id: \"UNIQUEID01\" delete_id: \"UNIQUEID02\" table: \"dml_people\" duration: 0 sql: \"DELETE /*ID$UNIQUEID02*/ FROM `dml_people` WHERE (`id` >= 34.56)\" length_args: 0 length_raw_args: 0 source: \"d\" error: \"<nil>\"\n",
				buf.String())
		})

		t.Run("Prepare", func(t *testing.T) {
			defer buf.Reset()
			stmt, err := d.Prepare(context.TODO())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)
			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQUEID01\" delete_id: \"UNIQUEID02\" table: \"dml_people\" duration: 0 error: \"<nil>\" sql: \"DELETE /*ID$UNIQUEID02*/ FROM `dml_people` WHERE (`id` >= 34.56)\"\n",
				buf.String())
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()
			assert.NoError(t, rConn.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
				_, err := tx.DeleteFrom("dml_people").Where(dml.Column("id").GreaterOrEqual().Float64(36.56)).WithDBR().Interpolate().ExecContext(context.TODO())
				return err
			}))
			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQUEID01\" tx_id: \"UNIQUEID03\"\nDEBUG Exec conn_pool_id: \"UNIQUEID01\" tx_id: \"UNIQUEID03\" delete_id: \"UNIQUEID04\" table: \"dml_people\" duration: 0 sql: \"DELETE /*ID$UNIQUEID04*/ FROM `dml_people` WHERE (`id` >= 36.56)\" length_args: 0 length_raw_args: 0 source: \"d\" error: \"<nil>\"\nDEBUG Commit conn_pool_id: \"UNIQUEID01\" tx_id: \"UNIQUEID03\" duration: 0\n",
				buf.String())
		})
	})

	t.Run("Conn", func(t *testing.T) {
		conn, err := rConn.Conn(context.TODO())
		assert.NoError(t, err)

		d := conn.DeleteFrom("dml_people").Where(dml.Column("id").GreaterOrEqual().PlaceHolder())

		t.Run("Exec", func(t *testing.T) {
			defer func() {
				buf.Reset()
			}()

			_, err := d.WithDBR().Interpolate().ExecContext(context.TODO(), 39.56)
			assert.NoError(t, err)

			assert.Exactly(t, "DEBUG Exec conn_pool_id: \"UNIQUEID01\" conn_id: \"UNIQUEID05\" delete_id: \"UNIQUEID06\" table: \"dml_people\" duration: 0 sql: \"DELETE /*ID$UNIQUEID06*/ FROM `dml_people` WHERE (`id` >= 39.56)\" length_args: 0 length_raw_args: 1 source: \"d\" error: \"<nil>\"\n",
				buf.String())
		})

		t.Run("Prepare", func(t *testing.T) {
			defer buf.Reset()

			stmt, err := d.Prepare(context.TODO())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQUEID01\" conn_id: \"UNIQUEID05\" delete_id: \"UNIQUEID06\" table: \"dml_people\" duration: 0 error: \"<nil>\" sql: \"DELETE /*ID$UNIQUEID06*/ FROM `dml_people` WHERE (`id` >= ?)\"\n",
				buf.String())
		})

		t.Run("Prepare Exec", func(t *testing.T) {
			defer buf.Reset()

			stmt, err := d.Prepare(context.TODO())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)

			_, err = stmt.WithDBR().ExecContext(context.TODO(), 41.57)
			assert.NoError(t, err)

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQUEID01\" conn_id: \"UNIQUEID05\" delete_id: \"UNIQUEID06\" table: \"dml_people\" duration: 0 error: \"<nil>\" sql: \"DELETE /*ID$UNIQUEID06*/ FROM `dml_people` WHERE (`id` >= ?)\"\nDEBUG Exec conn_pool_id: \"UNIQUEID01\" conn_id: \"UNIQUEID05\" delete_id: \"UNIQUEID06\" table: \"dml_people\" duration: 0 sql: \"\" length_args: 1 length_raw_args: 1 source: \"d\" error: \"<nil>\"\n",
				buf.String())
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()

			assert.NoError(t, conn.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
				_, err := tx.DeleteFrom("dml_people").Where(dml.Column("id").GreaterOrEqual().Float64(37.56)).
					WithDBR().Interpolate().ExecContext(context.TODO())
				return err
			}))

			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQUEID01\" conn_id: \"UNIQUEID05\" tx_id: \"UNIQUEID07\"\nDEBUG Exec conn_pool_id: \"UNIQUEID01\" conn_id: \"UNIQUEID05\" tx_id: \"UNIQUEID07\" delete_id: \"UNIQUEID08\" table: \"dml_people\" duration: 0 sql: \"DELETE /*ID$UNIQUEID08*/ FROM `dml_people` WHERE (`id` >= 37.56)\" length_args: 0 length_raw_args: 0 source: \"d\" error: \"<nil>\"\nDEBUG Commit conn_pool_id: \"UNIQUEID01\" conn_id: \"UNIQUEID05\" tx_id: \"UNIQUEID07\" duration: 0\n",
				buf.String())
		})

		t.Run("Tx Rollback", func(t *testing.T) {
			defer buf.Reset()

			assert.Error(t, conn.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
				_, err := tx.DeleteFrom("dml_people").Where(dml.Column("id").GreaterOrEqual().PlaceHolder()).WithDBR().Interpolate().ExecContext(context.TODO())
				return err
			}))

			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQUEID01\" conn_id: \"UNIQUEID05\" tx_id: \"UNIQUEID09\"\nDEBUG Exec conn_pool_id: \"UNIQUEID01\" conn_id: \"UNIQUEID05\" tx_id: \"UNIQUEID09\" delete_id: \"UNIQUEID10\" table: \"dml_people\" duration: 0 sql: \"DELETE /*ID$UNIQUEID10*/ FROM `dml_people` WHERE (`id` >= ?)\" length_args: 0 length_raw_args: 0 source: \"d\" error: \"<nil>\"\nDEBUG Rollback conn_pool_id: \"UNIQUEID01\" conn_id: \"UNIQUEID05\" tx_id: \"UNIQUEID09\" duration: 0\n",
				buf.String())
		})
	})
}

func TestWithLogger_Select(t *testing.T) {
	uniID := new(int32)
	rConn := createRealSession(t)
	defer dmltest.Close(t, rConn)

	uniqueIDFunc := func() string {
		return fmt.Sprintf("UNIQ%02d", atomic.AddInt32(uniID, 1))
	}

	buf := new(bytes.Buffer)
	lg := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithWriter(buf),
		logw.WithFlag(0), // no flags at all
	)
	assert.NoError(t, rConn.Options(dml.WithLogger(lg, uniqueIDFunc)))

	t.Run("ConnPool", func(t *testing.T) {
		pplSel := rConn.SelectFrom("dml_people").AddColumns("email").Where(dml.Column("id").Greater().PlaceHolder())

		t.Run("Query Error interpolation with iFace slice", func(t *testing.T) {
			defer buf.Reset()
			rows, err := pplSel.WithDBR().Interpolate().QueryContext(context.TODO(), 67896543123)
			assert.NotNil(t, rows)
			assert.NoError(t, err)
		})
		t.Run("Query Correct", func(t *testing.T) {
			defer buf.Reset()
			rows, err := pplSel.WithDBR().QueryContext(context.TODO(), 67896543123)
			assert.NoError(t, err)
			assert.NoError(t, rows.Close())

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dml_people\" duration: 0 sql: \"SELECT /*ID$UNIQ02*/ `email` FROM `dml_people` WHERE (`id` > ?)\" length_args: 1 source: \"s\" error: \"<nil>\"\n",
				buf.String())
		})

		t.Run("Load", func(t *testing.T) {
			defer buf.Reset()
			p := &dmlPerson{}
			_, err := pplSel.WithDBR().Load(context.TODO(), p, 67896543113)
			assert.NoError(t, err)

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dml_people\" duration: 0 sql: \"SELECT /*ID$UNIQ02*/ `email` FROM `dml_people` WHERE (`id` > ?)\" length_args: 1 source: \"s\" error: \"<nil>\"\nDEBUG Load conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dml_people\" duration: 0 id: \"UNIQ02\" error: \"<nil>\" ColumnMapper: \"*dml_test.dmlPerson\" row_count: 0\n",
				buf.String())
		})

		t.Run("LoadInt64", func(t *testing.T) {
			defer buf.Reset()
			_, _, err := pplSel.WithDBR().LoadNullInt64(context.TODO(), 67896543124)
			if !errors.NotFound.Match(err) {
				assert.NoError(t, err)
			}

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dml_people\" duration: 0 sql: \"SELECT /*ID$UNIQ02*/ `email` FROM `dml_people` WHERE (`id` > ?)\" length_args: 1 source: \"s\" error: \"<nil>\"\nDEBUG LoadPrimitive conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dml_people\" duration: 0 id: \"UNIQ02\" error: \"<nil>\" ptr_type: \"*null.Int64\"\n",
				buf.String())
		})

		t.Run("LoadInt64s", func(t *testing.T) {
			defer buf.Reset()
			_, err := pplSel.WithDBR().LoadInt64s(context.TODO(), nil, 67896543125)
			assert.NoError(t, err)

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dml_people\" duration: 0 sql: \"SELECT /*ID$UNIQ02*/ `email` FROM `dml_people` WHERE (`id` > ?)\" length_args: 1 source: \"s\" error: \"<nil>\"\nDEBUG LoadInt64s conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dml_people\" duration: 0 row_count: 0 error: \"<nil>\"\n",
				buf.String())
		})

		t.Run("LoadUint64", func(t *testing.T) {
			defer buf.Reset()
			_, _, err := pplSel.WithDBR().LoadNullUint64(context.TODO(), 67896543126)
			if !errors.NotFound.Match(err) {
				assert.NoError(t, err)
			}
			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dml_people\" duration: 0 sql: \"SELECT /*ID$UNIQ02*/ `email` FROM `dml_people` WHERE (`id` > ?)\" length_args: 1 source: \"s\" error: \"<nil>\"\nDEBUG LoadPrimitive conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dml_people\" duration: 0 id: \"UNIQ02\" error: \"<nil>\" ptr_type: \"*null.Uint64\"\n",
				buf.String())
		})

		t.Run("LoadUint64s", func(t *testing.T) {
			defer buf.Reset()
			_, err := pplSel.WithDBR().LoadUint64s(context.TODO(), nil, 67896543127)
			assert.NoError(t, err)

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dml_people\" duration: 0 sql: \"SELECT /*ID$UNIQ02*/ `email` FROM `dml_people` WHERE (`id` > ?)\" length_args: 1 source: \"s\" error: \"<nil>\"\nDEBUG LoadUint64s conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dml_people\" duration: 0 row_count: 0 id: \"UNIQ02\" error: \"<nil>\"\n",
				buf.String())
		})

		t.Run("LoadFloat64", func(t *testing.T) {
			defer buf.Reset()
			_, _, err := pplSel.WithDBR().LoadNullFloat64(context.TODO(), 678965.43125)
			if !errors.NotFound.Match(err) {
				assert.NoError(t, err)
			}
			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dml_people\" duration: 0 sql: \"SELECT /*ID$UNIQ02*/ `email` FROM `dml_people` WHERE (`id` > ?)\" length_args: 1 source: \"s\" error: \"<nil>\"\nDEBUG LoadPrimitive conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dml_people\" duration: 0 id: \"UNIQ02\" error: \"<nil>\" ptr_type: \"*null.Float64\"\n",
				buf.String())
		})

		t.Run("LoadFloat64s", func(t *testing.T) {
			defer buf.Reset()
			_, err := pplSel.WithDBR().LoadFloat64s(context.TODO(), nil, 6789654.3125)
			assert.NoError(t, err)

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dml_people\" duration: 0 sql: \"SELECT /*ID$UNIQ02*/ `email` FROM `dml_people` WHERE (`id` > ?)\" length_args: 1 source: \"s\" error: \"<nil>\"\nDEBUG LoadFloat64s conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dml_people\" duration: 0 id: \"UNIQ02\" error: \"<nil>\"\n",
				buf.String())
		})

		t.Run("LoadString", func(t *testing.T) {
			defer buf.Reset()
			_, _, err := pplSel.WithDBR().LoadNullString(context.TODO(), "hello")
			if !errors.NotFound.Match(err) {
				assert.NoError(t, err)
			}

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dml_people\" duration: 0 sql: \"SELECT /*ID$UNIQ02*/ `email` FROM `dml_people` WHERE (`id` > ?)\" length_args: 1 source: \"s\" error: \"<nil>\"\nDEBUG LoadPrimitive conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dml_people\" duration: 0 id: \"UNIQ02\" error: \"<nil>\" ptr_type: \"*null.String\"\n",
				buf.String())
		})

		t.Run("LoadStrings", func(t *testing.T) {
			defer buf.Reset()

			_, err := pplSel.WithDBR().LoadStrings(context.TODO(), nil, 99987)
			assert.NoError(t, err)

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dml_people\" duration: 0 sql: \"SELECT /*ID$UNIQ02*/ `email` FROM `dml_people` WHERE (`id` > ?)\" length_args: 1 source: \"s\" error: \"<nil>\"\nDEBUG LoadStrings conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dml_people\" duration: 0 row_count: 0 id: \"UNIQ02\" error: \"<nil>\"\n",
				buf.String())
		})

		t.Run("Prepare", func(t *testing.T) {
			defer buf.Reset()
			stmt, err := pplSel.Prepare(context.TODO())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ01\" select_id: \"UNIQ02\" table: \"dml_people\" duration: 0 error: \"<nil>\" sql: \"SELECT /*ID$UNIQ02*/ `email` FROM `dml_people` WHERE (`id` > ?)\"\n",
				buf.String())
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()
			assert.NoError(t, rConn.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
				rows, err := tx.SelectFrom("dml_people").
					AddColumns("name", "email").Where(dml.Column("id").In().Int64s(7, 9)).
					WithDBR().QueryContext(context.TODO())
				assert.NoError(t, err)
				return rows.Close()
			}))
			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ01\" tx_id: \"UNIQ03\"\nDEBUG Query conn_pool_id: \"UNIQ01\" tx_id: \"UNIQ03\" select_id: \"UNIQ04\" table: \"dml_people\" duration: 0 sql: \"SELECT /*ID$UNIQ04*/ `name`, `email` FROM `dml_people` WHERE (`id` IN (7,9))\" length_args: 0 source: \"s\" error: \"<nil>\"\nDEBUG Commit conn_pool_id: \"UNIQ01\" tx_id: \"UNIQ03\" duration: 0\n",
				buf.String())
		})
	})

	t.Run("ConnSingle", func(t *testing.T) {
		conn, err := rConn.Conn(context.TODO())
		defer dmltest.Close(t, conn)
		assert.NoError(t, err)

		pplSel := conn.SelectFrom("dml_people", "dp2").AddColumns("name", "email").Where(dml.Column("id").Less().PlaceHolder())

		t.Run("Query", func(t *testing.T) {
			defer buf.Reset()

			rows, err := pplSel.WithDBR().QueryContext(context.TODO(), -3)
			assert.NoError(t, err)
			dmltest.Close(t, rows)

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" select_id: \"UNIQ06\" table: \"dml_people\" duration: 0 sql: \"SELECT /*ID$UNIQ06*/ `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` < ?)\" length_args: 1 source: \"s\" error: \"<nil>\"\n",
				buf.String())
		})

		t.Run("Load", func(t *testing.T) {
			defer buf.Reset()
			p := &dmlPerson{}
			_, err := pplSel.WithDBR().Load(context.TODO(), p, -2)
			assert.NoError(t, err)

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" select_id: \"UNIQ06\" table: \"dml_people\" duration: 0 sql: \"SELECT /*ID$UNIQ06*/ `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` < ?)\" length_args: 1 source: \"s\" error: \"<nil>\"\nDEBUG Load conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" select_id: \"UNIQ06\" table: \"dml_people\" duration: 0 id: \"UNIQ06\" error: \"<nil>\" ColumnMapper: \"*dml_test.dmlPerson\" row_count: 0\n",
				buf.String())
		})

		t.Run("Prepare", func(t *testing.T) {
			stmt, err := pplSel.Prepare(context.TODO())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)
			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" select_id: \"UNIQ06\" table: \"dml_people\" duration: 0 error: \"<nil>\" sql: \"SELECT /*ID$UNIQ06*/ `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` < ?)\"\n",
				buf.String())
			buf.Reset()

			t.Run("QueryRow", func(t *testing.T) {
				defer buf.Reset()
				rows := stmt.WithDBR().QueryRowContext(context.TODO(), -8)
				var x string
				err := rows.Scan(&x)
				assert.True(t, errors.Cause(err) == sql.ErrNoRows, "but got this error: %#v", err)
				_ = x

				assert.Exactly(t, "DEBUG QueryRowContext conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" select_id: \"UNIQ06\" table: \"dml_people\" duration: 0 sql: \"\" source: \"s\" error: \"<nil>\"\n",
					buf.String())
			})

			t.Run("Query", func(t *testing.T) {
				defer buf.Reset()
				rows, err := stmt.WithDBR().QueryContext(context.TODO(), -4)
				assert.NoError(t, err)
				dmltest.Close(t, rows)
				assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" select_id: \"UNIQ06\" table: \"dml_people\" duration: 0 sql: \"\" length_args: 1 source: \"s\" error: \"<nil>\"\n",
					buf.String())
			})

			t.Run("Load", func(t *testing.T) {
				defer buf.Reset()
				p := &dmlPerson{}
				_, err := stmt.WithDBR().Load(context.TODO(), p, -6)
				assert.NoError(t, err)
				assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" select_id: \"UNIQ06\" table: \"dml_people\" duration: 0 sql: \"\" length_args: 1 source: \"s\" error: \"<nil>\"\nDEBUG Load conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" select_id: \"UNIQ06\" table: \"dml_people\" duration: 0 id: \"UNIQ06\" error: \"<nil>\" ColumnMapper: \"*dml_test.dmlPerson\" row_count: 0\n",
					buf.String())
			})

			t.Run("LoadInt64", func(t *testing.T) {
				defer buf.Reset()
				_, _, err := stmt.WithDBR().LoadNullInt64(context.TODO(), -7)
				if !errors.NotFound.Match(err) {
					assert.NoError(t, err)
				}
				assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" select_id: \"UNIQ06\" table: \"dml_people\" duration: 0 sql: \"\" length_args: 1 source: \"s\" error: \"<nil>\"\nDEBUG LoadPrimitive conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" select_id: \"UNIQ06\" table: \"dml_people\" duration: 0 id: \"UNIQ06\" error: \"<nil>\" ptr_type: \"*null.Int64\"\n",
					buf.String())
			})

			t.Run("LoadInt64s", func(t *testing.T) {
				defer buf.Reset()
				iSl, err := stmt.WithDBR().LoadInt64s(context.TODO(), nil, -7)
				assert.NoError(t, err)
				assert.Nil(t, iSl)
				assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" select_id: \"UNIQ06\" table: \"dml_people\" duration: 0 sql: \"\" length_args: 1 source: \"s\" error: \"<nil>\"\nDEBUG LoadInt64s conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" select_id: \"UNIQ06\" table: \"dml_people\" duration: 0 row_count: 0 error: \"<nil>\"\n",
					buf.String())
			})
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()
			assert.NoError(t, conn.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
				rows, err := tx.SelectFrom("dml_people").AddColumns("name", "email").Where(dml.Column("id").In().Int64s(71, 91)).
					WithDBR().QueryContext(context.TODO())
				if err != nil {
					return err
				}
				return rows.Close()
			}))
			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ07\"\nDEBUG Query conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ07\" select_id: \"UNIQ08\" table: \"dml_people\" duration: 0 sql: \"SELECT /*ID$UNIQ08*/ `name`, `email` FROM `dml_people` WHERE (`id` IN (71,91))\" length_args: 0 source: \"s\" error: \"<nil>\"\nDEBUG Commit conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ07\" duration: 0\n",
				buf.String())
		})

		t.Run("Tx Rollback", func(t *testing.T) {
			defer buf.Reset()
			assert.Error(t, conn.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
				rows, err := tx.SelectFrom("dml_people").AddColumns("name", "email").Where(dml.Column("id").In().PlaceHolder()).
					WithDBR().QueryContext(context.TODO())
				if err != nil {
					return err
				}
				return rows.Close()
			}))

			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ09\"\nDEBUG Query conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ09\" select_id: \"UNIQ10\" table: \"dml_people\" duration: 0 sql: \"SELECT /*ID$UNIQ10*/ `name`, `email` FROM `dml_people` WHERE (`id` IN ?)\" length_args: 0 source: \"s\" error: \"<nil>\"\nDEBUG Rollback conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ09\" duration: 0\n",
				buf.String())
		})
	})
}

func TestWithLogger_Union(t *testing.T) {
	uniID := new(int32)
	rConn := createRealSession(t)
	defer dmltest.Close(t, rConn)

	uniqueIDFunc := func() string {
		return fmt.Sprintf("UNIQ%02d", atomic.AddInt32(uniID, 1))
	}

	buf := new(bytes.Buffer)
	lg := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithWriter(buf),
		logw.WithFlag(0), // no flags at all
	)
	assert.NoError(t, rConn.Options(dml.WithLogger(lg, uniqueIDFunc)))

	t.Run("ConnPool", func(t *testing.T) {
		u := rConn.Union(
			dml.NewSelect("name").AddColumnsAliases("email", "email").From("dml_people"),
			dml.NewSelect("name", "email").FromAlias("dml_people", "dp2").Where(dml.Column("id").In().Int64s(6, 8)),
		)

		t.Run("Query", func(t *testing.T) {
			defer buf.Reset()
			rows, err := u.WithDBR().QueryContext(context.TODO())
			assert.NoError(t, err)
			assert.NoError(t, rows.Close())

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" union_id: \"UNIQ02\" tables: \"dml_people, dml_people\" duration: 0 sql: \"(SELECT /*ID$UNIQ02*/ `name`, `email` AS `email` FROM `dml_people`)\\nUNION\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (6,8)))\" length_args: 0 source: \"n\" error: \"<nil>\"\n",
				buf.String())
		})

		t.Run("Load", func(t *testing.T) {
			defer buf.Reset()
			p := &dmlPerson{}
			_, err := u.WithDBR().Interpolate().Load(context.TODO(), p)
			assert.NoError(t, err)

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" union_id: \"UNIQ02\" tables: \"dml_people, dml_people\" duration: 0 sql: \"(SELECT /*ID$UNIQ02*/ `name`, `email` AS `email` FROM `dml_people`)\\nUNION\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (6,8)))\" length_args: 0 source: \"n\" error: \"<nil>\"\nDEBUG Load conn_pool_id: \"UNIQ01\" union_id: \"UNIQ02\" tables: \"dml_people, dml_people\" duration: 0 id: \"UNIQ02\" error: \"<nil>\" ColumnMapper: \"*dml_test.dmlPerson\" row_count: 0\n",
				buf.String())
		})

		t.Run("Prepare", func(t *testing.T) {
			defer buf.Reset()
			stmt, err := u.Prepare(context.TODO())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ01\" union_id: \"UNIQ02\" tables: \"dml_people, dml_people\" duration: 0 error: \"<nil>\" sql: \"(SELECT /*ID$UNIQ02*/ `name`, `email` AS `email` FROM `dml_people`)\\nUNION\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (6,8)))\"\n",
				buf.String())
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()
			assert.NoError(t, rConn.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
				rows, err := tx.Union(
					dml.NewSelect("name").AddColumnsAliases("email", "email").From("dml_people"),
					dml.NewSelect("name", "email").FromAlias("dml_people", "dp2").Where(dml.Column("id").In().Int64s(7, 9)),
				).WithDBR().Interpolate().QueryContext(context.TODO())

				assert.NoError(t, rows.Close())
				return err
			}))
			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ01\" tx_id: \"UNIQ03\"\nDEBUG Query conn_pool_id: \"UNIQ01\" tx_id: \"UNIQ03\" union_id: \"UNIQ04\" tables: \"dml_people, dml_people\" duration: 0 sql: \"(SELECT /*ID$UNIQ04*/ `name`, `email` AS `email` FROM `dml_people`)\\nUNION\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (7,9)))\" length_args: 0 source: \"n\" error: \"<nil>\"\nDEBUG Commit conn_pool_id: \"UNIQ01\" tx_id: \"UNIQ03\" duration: 0\n",
				buf.String())
		})
	})

	t.Run("Conn", func(t *testing.T) {
		conn, err := rConn.Conn(context.TODO())
		assert.NoError(t, err)

		u := conn.Union(
			dml.NewSelect("name").AddColumnsAliases("email", "email").From("dml_people"),
			dml.NewSelect("name", "email").FromAlias("dml_people", "dp2").Where(dml.Column("id").In().Int64s(61, 81)),
		)
		t.Run("Query", func(t *testing.T) {
			defer buf.Reset()

			rows, err := u.WithDBR().Interpolate().QueryContext(context.TODO())
			assert.NoError(t, err)
			assert.NoError(t, rows.Close())

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" union_id: \"UNIQ06\" tables: \"dml_people, dml_people\" duration: 0 sql: \"(SELECT /*ID$UNIQ06*/ `name`, `email` AS `email` FROM `dml_people`)\\nUNION\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (61,81)))\" length_args: 0 source: \"n\" error: \"<nil>\"\n",
				buf.String())
		})

		t.Run("Load", func(t *testing.T) {
			defer buf.Reset()
			p := &dmlPerson{}
			_, err := u.WithDBR().Load(context.TODO(), p)
			assert.NoError(t, err)

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" union_id: \"UNIQ06\" tables: \"dml_people, dml_people\" duration: 0 sql: \"(SELECT /*ID$UNIQ06*/ `name`, `email` AS `email` FROM `dml_people`)\\nUNION\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (61,81)))\" length_args: 0 source: \"n\" error: \"<nil>\"\nDEBUG Load conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" union_id: \"UNIQ06\" tables: \"dml_people, dml_people\" duration: 0 id: \"UNIQ06\" error: \"<nil>\" ColumnMapper: \"*dml_test.dmlPerson\" row_count: 0\n",
				buf.String())
		})

		t.Run("Prepare", func(t *testing.T) {
			defer buf.Reset()

			stmt, err := u.Prepare(context.TODO())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" union_id: \"UNIQ06\" tables: \"dml_people, dml_people\" duration: 0 error: \"<nil>\" sql: \"(SELECT /*ID$UNIQ06*/ `name`, `email` AS `email` FROM `dml_people`)\\nUNION\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (61,81)))\"\n",
				buf.String())
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()
			assert.NoError(t, conn.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
				rows, err := tx.Union(
					dml.NewSelect("name").AddColumnsAliases("email", "email").From("dml_people"),
					dml.NewSelect("name", "email").FromAlias("dml_people", "dp2").Where(dml.Column("id").In().Int64s(71, 91)),
				).WithDBR().Interpolate().QueryContext(context.TODO())
				if err != nil {
					return err
				}
				return rows.Close()
			}))
			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ07\"\nDEBUG Query conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ07\" union_id: \"UNIQ08\" tables: \"dml_people, dml_people\" duration: 0 sql: \"(SELECT /*ID$UNIQ08*/ `name`, `email` AS `email` FROM `dml_people`)\\nUNION\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (71,91)))\" length_args: 0 source: \"n\" error: \"<nil>\"\nDEBUG Commit conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ07\" duration: 0\n",
				buf.String())
		})

		t.Run("Tx Rollback", func(t *testing.T) {
			defer buf.Reset()
			assert.Error(t, conn.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
				rows, err := tx.Union(
					dml.NewSelect("name").AddColumnsAliases("email", "email").From("dml_people"),
					dml.NewSelect("name", "email").FromAlias("dml_people", "dp2").Where(dml.Column("id").In().PlaceHolder()),
				).WithDBR().Interpolate().QueryContext(context.TODO())
				if err != nil {
					return err
				}
				return rows.Close()
			}))

			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ09\"\nDEBUG Query conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ09\" union_id: \"UNIQ10\" tables: \"dml_people, dml_people\" duration: 0 sql: \"(SELECT /*ID$UNIQ10*/ `name`, `email` AS `email` FROM `dml_people`)\\nUNION\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN ?))\" length_args: 0 source: \"n\" error: \"<nil>\"\nDEBUG Rollback conn_pool_id: \"UNIQ01\" conn_id: \"UNIQ05\" tx_id: \"UNIQ09\" duration: 0\n",
				buf.String())
		})
	})
}

func TestWithLogger_Update(t *testing.T) {
	uniID := new(int32)
	rConn := createRealSession(t)
	defer dmltest.Close(t, rConn)

	uniqueIDFunc := func() string {
		return fmt.Sprintf("UNIQ%02d", atomic.AddInt32(uniID, 3))
	}

	buf := new(bytes.Buffer)
	lg := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithWriter(buf),
		logw.WithFlag(0), // no flags at all
	)
	assert.NoError(t, rConn.Options(dml.WithLogger(lg, uniqueIDFunc)))

	t.Run("ConnPool", func(t *testing.T) {
		d := rConn.Update("dml_people").AddClauses(
			dml.Column("email").Str("new@email.com"),
		).Where(dml.Column("id").GreaterOrEqual().Float64(78.31))

		t.Run("Exec", func(t *testing.T) {
			defer buf.Reset()
			_, err := d.WithDBR().ExecContext(context.TODO())
			assert.NoError(t, err)

			assert.Exactly(t, "DEBUG Exec conn_pool_id: \"UNIQ03\" update_id: \"UNIQ06\" table: \"dml_people\" duration: 0 sql: \"UPDATE /*ID$UNIQ06*/ `dml_people` SET `email`='new@email.com' WHERE (`id` >= 78.31)\" length_args: 0 length_raw_args: 0 source: \"u\" error: \"<nil>\"\n",
				buf.String())
		})

		t.Run("Prepare", func(t *testing.T) {
			defer buf.Reset()
			stmt, err := d.Prepare(context.TODO())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ03\" update_id: \"UNIQ06\" table: \"dml_people\" duration: 0 error: \"<nil>\" sql: \"UPDATE /*ID$UNIQ06*/ `dml_people` SET `email`='new@email.com' WHERE (`id` >= 78.31)\"\n",
				buf.String())
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()
			assert.NoError(t, rConn.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
				_, err := tx.Update("dml_people").AddClauses(
					dml.Column("email").Str("new@email.com"),
				).Where(dml.Column("id").GreaterOrEqual().Float64(36.56)).WithDBR().ExecContext(context.TODO())
				return err
			}))
			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ03\" tx_id: \"UNIQ09\"\nDEBUG Exec conn_pool_id: \"UNIQ03\" tx_id: \"UNIQ09\" update_id: \"UNIQ12\" table: \"dml_people\" duration: 0 sql: \"UPDATE /*ID$UNIQ12*/ `dml_people` SET `email`='new@email.com' WHERE (`id` >= 36.56)\" length_args: 0 length_raw_args: 0 source: \"u\" error: \"<nil>\"\nDEBUG Commit conn_pool_id: \"UNIQ03\" tx_id: \"UNIQ09\" duration: 0\n",
				buf.String())
		})
	})

	t.Run("Conn", func(t *testing.T) {
		conn, err := rConn.Conn(context.TODO())
		assert.NoError(t, err)

		d := conn.Update("dml_people").AddClauses(
			dml.Column("email").Str("new@email.com"),
		).Where(dml.Column("id").GreaterOrEqual().Float64(21.56))

		t.Run("Exec", func(t *testing.T) {
			defer buf.Reset()

			_, err := d.WithDBR().ExecContext(context.TODO())
			assert.NoError(t, err)

			assert.Exactly(t, "DEBUG Exec conn_pool_id: \"UNIQ03\" conn_id: \"UNIQ15\" update_id: \"UNIQ18\" table: \"dml_people\" duration: 0 sql: \"UPDATE /*ID$UNIQ18*/ `dml_people` SET `email`='new@email.com' WHERE (`id` >= 21.56)\" length_args: 0 length_raw_args: 0 source: \"u\" error: \"<nil>\"\n",
				buf.String())
		})

		t.Run("Prepare", func(t *testing.T) {
			defer buf.Reset()

			stmt, err := d.Prepare(context.TODO())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ03\" conn_id: \"UNIQ15\" update_id: \"UNIQ18\" table: \"dml_people\" duration: 0 error: \"<nil>\" sql: \"UPDATE /*ID$UNIQ18*/ `dml_people` SET `email`='new@email.com' WHERE (`id` >= 21.56)\"\n",
				buf.String())
		})

		t.Run("Prepare Exec", func(t *testing.T) {
			defer buf.Reset()

			stmt, err := d.Prepare(context.TODO())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)

			_, err = stmt.WithDBR().ExecContext(context.TODO())
			assert.NoError(t, err)

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ03\" conn_id: \"UNIQ15\" update_id: \"UNIQ18\" table: \"dml_people\" duration: 0 error: \"<nil>\" sql: \"UPDATE /*ID$UNIQ18*/ `dml_people` SET `email`='new@email.com' WHERE (`id` >= 21.56)\"\nDEBUG Exec conn_pool_id: \"UNIQ03\" conn_id: \"UNIQ15\" update_id: \"UNIQ18\" table: \"dml_people\" duration: 0 sql: \"\" length_args: 0 length_raw_args: 0 source: \"u\" error: \"<nil>\"\n",
				buf.String())
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()
			assert.NoError(t, conn.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
				_, err := tx.Update("dml_people").AddClauses(
					dml.Column("email").Str("new@email.com"),
				).Where(dml.Column("id").GreaterOrEqual().Float64(39.56)).WithDBR().ExecContext(context.TODO())
				return err
			}))

			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ03\" conn_id: \"UNIQ15\" tx_id: \"UNIQ21\"\nDEBUG Exec conn_pool_id: \"UNIQ03\" conn_id: \"UNIQ15\" tx_id: \"UNIQ21\" update_id: \"UNIQ24\" table: \"dml_people\" duration: 0 sql: \"UPDATE /*ID$UNIQ24*/ `dml_people` SET `email`='new@email.com' WHERE (`id` >= 39.56)\" length_args: 0 length_raw_args: 0 source: \"u\" error: \"<nil>\"\nDEBUG Commit conn_pool_id: \"UNIQ03\" conn_id: \"UNIQ15\" tx_id: \"UNIQ21\" duration: 0\n",
				buf.String())
		})

		t.Run("Tx Rollback", func(t *testing.T) {
			defer buf.Reset()
			assert.Error(t, conn.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
				_, err := tx.Update("dml_people").AddClauses(
					dml.Column("email").Str("new@email.com"),
				).Where(dml.Column("id").GreaterOrEqual().PlaceHolder()).WithDBR().ExecContext(context.TODO())
				return err
			}))

			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ03\" conn_id: \"UNIQ15\" tx_id: \"UNIQ27\"\nDEBUG Exec conn_pool_id: \"UNIQ03\" conn_id: \"UNIQ15\" tx_id: \"UNIQ27\" update_id: \"UNIQ30\" table: \"dml_people\" duration: 0 sql: \"UPDATE /*ID$UNIQ30*/ `dml_people` SET `email`='new@email.com' WHERE (`id` >= ?)\" length_args: 0 length_raw_args: 0 source: \"u\" error: \"<nil>\"\nDEBUG Rollback conn_pool_id: \"UNIQ03\" conn_id: \"UNIQ15\" tx_id: \"UNIQ27\" duration: 0\n",
				buf.String())
		})
	})
}

func TestWithLogger_WithCTE(t *testing.T) {
	uniID := new(int32)
	rConn := createRealSession(t)
	defer dmltest.Close(t, rConn)

	uniqueIDFunc := func() string {
		return fmt.Sprintf("UNIQ%02d", atomic.AddInt32(uniID, 2))
	}

	buf := new(bytes.Buffer)
	lg := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithWriter(buf),
		logw.WithFlag(0), // no flags at all
	)
	assert.NoError(t, rConn.Options(dml.WithLogger(lg, uniqueIDFunc)))

	cte := dml.WithCTE{
		Name:    "zehTeEh",
		Columns: []string{"name2", "email2"},
		Union: dml.NewUnion(
			dml.NewSelect("name").AddColumnsAliases("email", "email").From("dml_people"),
			dml.NewSelect("name", "email").FromAlias("dml_people", "dp2").Where(dml.Column("id").In().Int64s(6, 8)),
		).All(),
	}
	cteSel := dml.NewSelect().Star().From("zehTeEh")

	t.Run("ConnPool", func(t *testing.T) {
		wth := rConn.With(cte).Select(cteSel)

		t.Run("Query", func(t *testing.T) {
			defer buf.Reset()
			rows, err := wth.WithDBR().Interpolate().QueryContext(context.TODO())
			assert.NoError(t, err)
			assert.NoError(t, rows.Close())

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ02\" with_cte_id: \"UNIQ04\" tables: \"zehTeEh\" duration: 0 sql: \"WITH /*ID$UNIQ04*/ `zehTeEh` (`name2`,`email2`) AS ((SELECT `name`, `email` AS `email` FROM `dml_people`)\\nUNION ALL\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (6,8))))\\nSELECT * FROM `zehTeEh`\" length_args: 0 source: \"w\" error: \"<nil>\"\n",
				buf.String())
		})

		t.Run("Load", func(t *testing.T) {
			defer buf.Reset()
			p := &dmlPerson{}
			_, err := wth.WithDBR().Interpolate().Load(context.TODO(), p)
			assert.NoError(t, err)

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ02\" with_cte_id: \"UNIQ04\" tables: \"zehTeEh\" duration: 0 sql: \"WITH /*ID$UNIQ04*/ `zehTeEh` (`name2`,`email2`) AS ((SELECT `name`, `email` AS `email` FROM `dml_people`)\\nUNION ALL\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (6,8))))\\nSELECT * FROM `zehTeEh`\" length_args: 0 source: \"w\" error: \"<nil>\"\nDEBUG Load conn_pool_id: \"UNIQ02\" with_cte_id: \"UNIQ04\" tables: \"zehTeEh\" duration: 0 id: \"UNIQ04\" error: \"<nil>\" ColumnMapper: \"*dml_test.dmlPerson\" row_count: 0\n",
				buf.String())
		})

		t.Run("Prepare", func(t *testing.T) {
			defer buf.Reset()
			stmt, err := wth.Prepare(context.TODO())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ02\" with_cte_id: \"UNIQ04\" tables: \"zehTeEh\" duration: 0 error: \"<nil>\" sql: \"WITH /*ID$UNIQ04*/ `zehTeEh` (`name2`,`email2`) AS ((SELECT `name`, `email` AS `email` FROM `dml_people`)\\nUNION ALL\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (6,8))))\\nSELECT * FROM `zehTeEh`\"\n",
				buf.String())
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()
			assert.NoError(t, rConn.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
				rows, err := tx.With(
					dml.WithCTE{
						Name:    "zehTeEh",
						Columns: []string{"name2", "email2"},
						Union: dml.NewUnion(
							dml.NewSelect("name").AddColumnsAliases("email", "email").From("dml_people"),
							dml.NewSelect("name", "email").FromAlias("dml_people", "dp2").Where(dml.Column("id").In().Int64s(6, 8)),
						).All(),
					},
				).Recursive().
					Select(dml.NewSelect().Star().From("zehTeEh")).WithDBR().Interpolate().QueryContext(context.TODO())

				assert.NoError(t, err)
				return rows.Close()
			}))
			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ02\" tx_id: \"UNIQ06\"\nDEBUG Query conn_pool_id: \"UNIQ02\" tx_id: \"UNIQ06\" with_cte_id: \"UNIQ08\" tables: \"zehTeEh\" duration: 0 sql: \"WITH /*ID$UNIQ08*/ RECURSIVE `zehTeEh` (`name2`,`email2`) AS ((SELECT `name`, `email` AS `email` FROM `dml_people`)\\nUNION ALL\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (6,8))))\\nSELECT * FROM `zehTeEh`\" length_args: 0 source: \"w\" error: \"<nil>\"\nDEBUG Commit conn_pool_id: \"UNIQ02\" tx_id: \"UNIQ06\" duration: 0\n",
				buf.String())
		})
	})

	t.Run("Conn", func(t *testing.T) {
		conn, err := rConn.Conn(context.TODO())
		assert.NoError(t, err)

		u := conn.With(cte).Select(cteSel)

		t.Run("Query", func(t *testing.T) {
			defer buf.Reset()

			rows, err := u.WithDBR().Interpolate().QueryContext(context.TODO())
			assert.NoError(t, err)
			assert.NoError(t, rows.Close())

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ02\" conn_id: \"UNIQ10\" with_cte_id: \"UNIQ12\" tables: \"zehTeEh\" duration: 0 sql: \"WITH /*ID$UNIQ12*/ `zehTeEh` (`name2`,`email2`) AS ((SELECT `name`, `email` AS `email` FROM `dml_people`)\\nUNION ALL\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (6,8))))\\nSELECT * FROM `zehTeEh`\" length_args: 0 source: \"w\" error: \"<nil>\"\n",
				buf.String())
		})

		t.Run("Load", func(t *testing.T) {
			defer buf.Reset()
			p := &dmlPerson{}
			_, err := u.WithDBR().Load(context.TODO(), p)
			assert.NoError(t, err)

			assert.Exactly(t, "DEBUG Query conn_pool_id: \"UNIQ02\" conn_id: \"UNIQ10\" with_cte_id: \"UNIQ12\" tables: \"zehTeEh\" duration: 0 sql: \"WITH /*ID$UNIQ12*/ `zehTeEh` (`name2`,`email2`) AS ((SELECT `name`, `email` AS `email` FROM `dml_people`)\\nUNION ALL\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (6,8))))\\nSELECT * FROM `zehTeEh`\" length_args: 0 source: \"w\" error: \"<nil>\"\nDEBUG Load conn_pool_id: \"UNIQ02\" conn_id: \"UNIQ10\" with_cte_id: \"UNIQ12\" tables: \"zehTeEh\" duration: 0 id: \"UNIQ12\" error: \"<nil>\" ColumnMapper: \"*dml_test.dmlPerson\" row_count: 0\n",
				buf.String())
		})

		t.Run("Prepare", func(t *testing.T) {
			defer buf.Reset()

			stmt, err := u.Prepare(context.TODO())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ02\" conn_id: \"UNIQ10\" with_cte_id: \"UNIQ12\" tables: \"zehTeEh\" duration: 0 error: \"<nil>\" sql: \"WITH /*ID$UNIQ12*/ `zehTeEh` (`name2`,`email2`) AS ((SELECT `name`, `email` AS `email` FROM `dml_people`)\\nUNION ALL\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (6,8))))\\nSELECT * FROM `zehTeEh`\"\n",
				buf.String())
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()
			assert.NoError(t, conn.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
				rows, err := tx.With(cte).Select(cteSel).WithDBR().QueryContext(context.TODO())
				if err != nil {
					return err
				}
				return rows.Close()
			}))
			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ02\" conn_id: \"UNIQ10\" tx_id: \"UNIQ14\"\nDEBUG Query conn_pool_id: \"UNIQ02\" conn_id: \"UNIQ10\" tx_id: \"UNIQ14\" with_cte_id: \"UNIQ16\" tables: \"zehTeEh\" duration: 0 sql: \"WITH /*ID$UNIQ16*/ `zehTeEh` (`name2`,`email2`) AS ((SELECT `name`, `email` AS `email` FROM `dml_people`)\\nUNION ALL\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (6,8))))\\nSELECT * FROM `zehTeEh`\" length_args: 0 source: \"w\" error: \"<nil>\"\nDEBUG Commit conn_pool_id: \"UNIQ02\" conn_id: \"UNIQ10\" tx_id: \"UNIQ14\" duration: 0\n",
				buf.String())
		})

		t.Run("Tx Rollback", func(t *testing.T) {
			defer buf.Reset()
			assert.Error(t, conn.Transaction(context.TODO(), nil, func(tx *dml.Tx) error {
				rows, err := tx.With(cte).Select(cteSel.Where(dml.Column("email").In().PlaceHolder())).WithDBR().QueryContext(context.TODO())
				if err != nil {
					return err
				}
				return rows.Close()
			}))

			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ02\" conn_id: \"UNIQ10\" tx_id: \"UNIQ18\"\nDEBUG Query conn_pool_id: \"UNIQ02\" conn_id: \"UNIQ10\" tx_id: \"UNIQ18\" with_cte_id: \"UNIQ20\" tables: \"zehTeEh\" duration: 0 sql: \"WITH /*ID$UNIQ20*/ `zehTeEh` (`name2`,`email2`) AS ((SELECT `name`, `email` AS `email` FROM `dml_people`)\\nUNION ALL\\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (6,8))))\\nSELECT * FROM `zehTeEh` WHERE (`email` IN ?)\" length_args: 0 source: \"w\" error: \"<nil>\"\nDEBUG Rollback conn_pool_id: \"UNIQ02\" conn_id: \"UNIQ10\" tx_id: \"UNIQ18\" duration: 0\n",
				buf.String())
		})
	})
}
