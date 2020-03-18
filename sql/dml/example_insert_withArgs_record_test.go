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
	"fmt"

	"github.com/corestoreio/errors"
	"github.com/weiwolves/pkg/sql/dml"
	"github.com/weiwolves/pkg/storage/null"
)

// Make sure that type productEntity implements interface.
var _ dml.ColumnMapper = (*productEntity)(nil)

// productEntity represents just a demo record.
type productEntity struct {
	EntityID       int64 // Auto Increment
	AttributeSetID int64
	TypeID         string
	SKU            null.String
	HasOptions     bool
}

func (pe productEntity) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		// This case gets executed when an INSERT statement doesn't contain any
		// columns.
		return cm.Int64(&pe.EntityID).Int64(&pe.AttributeSetID).String(&pe.TypeID).NullString(&pe.SKU).Bool(&pe.HasOptions).Err()
	}
	// This case gets executed when an INSERT statement requests specific columns.
	for cm.Next() {
		switch c := cm.Column(); c {
		case "attribute_set_id":
			cm.Int64(&pe.AttributeSetID)
		case "type_id":
			cm.String(&pe.TypeID)
		case "sku":
			cm.NullString(&pe.SKU)
		case "has_options":
			cm.Bool(&pe.HasOptions)
		default:
			return errors.NotFound.Newf("[dml_test] Column %q not found", c)
		}
	}
	return errors.WithStack(cm.Err())
}

// ExampleInsert_AddRecords inserts new data into table
// `catalog_product_entity`. First statement by specifying the exact column
// names. In the second example all columns values are getting inserted and you
// specify the number of place holders per record.
func ExampleInsert_WithArgs_record() {
	objs := []productEntity{
		{1, 5, "simple", null.MakeString("SOA9"), false},
		{2, 5, "virtual", null.String{}, true},
	}

	i := dml.NewInsert("catalog_product_entity").AddColumns("attribute_set_id", "type_id", "sku", "has_options").
		WithDBR().TestWithArgs(dml.Qualify("", objs[0]), dml.Qualify("", objs[1]))
	writeToSQLAndInterpolate(i)

	fmt.Print("\n\n")
	i = dml.NewInsert("catalog_product_entity").SetRecordPlaceHolderCount(5).
		WithDBR().TestWithArgs(dml.Qualify("", objs[0]), dml.Qualify("", objs[1]))
	writeToSQLAndInterpolate(i)

	// Output:
	// Prepared Statement:
	// INSERT INTO `catalog_product_entity`
	//(`attribute_set_id`,`type_id`,`sku`,`has_options`) VALUES (?,?,?,?),(?,?,?,?)
	// Arguments: [5 simple SOA9 false 5 virtual <nil> true]
	//
	// Interpolated Statement:
	// INSERT INTO `catalog_product_entity`
	//(`attribute_set_id`,`type_id`,`sku`,`has_options`) VALUES
	//(5,'simple','SOA9',0),(5,'virtual',NULL,1)
	//
	// Prepared Statement:
	// INSERT INTO `catalog_product_entity` VALUES (?,?,?,?,?),(?,?,?,?,?)
	// Arguments: [1 5 simple SOA9 false 2 5 virtual <nil> true]
	//
	// Interpolated Statement:
	// INSERT INTO `catalog_product_entity` VALUES
	//(1,5,'simple','SOA9',0),(2,5,'virtual',NULL,1)
}
