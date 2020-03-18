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
	"context"
	"fmt"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/weiwolves/pkg/sql/dml"
	"github.com/weiwolves/pkg/sql/dmltest"
	"github.com/weiwolves/pkg/storage/null"
)

// Make sure that type salesOrder implements interface.
var _ dml.ColumnMapper = (*salesOrder)(nil)

// salesOrder represents just a demo record.
type salesOrder struct {
	EntityID   int64  // Auto Increment
	State      string // processing, pending, shipped,
	StoreID    int64
	CustomerID int64
	GrandTotal null.Float64
}

func (so *salesOrder) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		return cm.Int64(&so.EntityID).String(&so.State).Int64(&so.StoreID).Err()
	}
	for cm.Next() {
		switch c := cm.Column(); c {
		case "entity_id":
			cm.Int64(&so.EntityID)
		case "state":
			cm.String(&so.State)
		case "store_id":
			cm.Int64(&so.StoreID)
		case "customer_id":
			cm.Int64(&so.CustomerID)
		case "grand_total":
			cm.NullFloat64(&so.GrandTotal)
		default:
			return errors.NotFound.Newf("[dml_test] Column %q not found", c)
		}
		if cm.Err() != nil {
			return cm.Err()
		}
	}
	return nil
}

// ExampleUpdate_Prepare can create a prepared statement or interpolated statements
// to run updates on table  `sales_order` with different objects. The SQL UPDATE
// statement acts as a template.
func ExampleUpdate_Prepare() {
	// <ignore_this>
	dbc, dbMock := dmltest.MockDB(nil)
	defer dmltest.MockClose(nil, dbc, dbMock)

	prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta(
		"UPDATE `sales_order` SET `state`=?, `customer_id`=?, `grand_total`=? WHERE (`shipping_method` IN ('DHL','UPS')) AND (`entity_id` = ?)",
	))

	prep.ExpectExec().WithArgs(
		"pending", int64(5678), 31.41459, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	prep.ExpectExec().WithArgs(
		"processing", int64(8912), nil, 2).
		WillReturnResult(sqlmock.NewResult(0, 1))
	// </ignore_this>

	// Create the prepared update statement
	stmt, err := dml.NewUpdate("sales_order").
		AddColumns("state", "customer_id", "grand_total").
		Where(
			dml.Column("shipping_method").In().Strs("DHL", "UPS"),
			dml.Column("entity_id").PlaceHolder(),
		).
		WithDB(dbc.DB).
		Prepare(context.TODO())
	if err != nil {
		fmt.Printf("Exec Error: %+v\n", err)
		return
	}
	defer stmt.Close()

	// Our objects which should update the columns in the database table
	// `sales_order`.
	collection := []*salesOrder{
		{1, "pending", 5, 5678, null.MakeFloat64(31.41459)},
		{2, "processing", 7, 8912, null.Float64{}},
	}
	for _, record := range collection {
		// We're not using an alias in the query so Qualify can have an empty
		// qualifier, which falls back to the default table name "sales_order".
		result, err := stmt.WithDBR().ExecContext(context.Background(), dml.Qualify("", record))
		if err != nil {
			fmt.Printf("Exec Error: %+v\n", err)
			return
		}

		ra, err := result.RowsAffected()
		if err != nil {
			fmt.Printf("RowsAffected Error: %+v\n", err)
			return
		}
		fmt.Printf("RowsAffected %d\n", ra)
	}

	// Output:
	// RowsAffected 1
	// RowsAffected 1
}
