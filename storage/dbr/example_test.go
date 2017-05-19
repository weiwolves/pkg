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
	"fmt"
	"os"
	"time"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/wordwrap"
)

func writeToSqlAndPreprocess(qb dbr.QueryBuilder) {
	sqlStr, args, err := qb.ToSQL()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Println("Prepared Statement:")
	wordwrap.Fstring(os.Stdout, sqlStr, 80)
	if len(args) > 0 {
		fmt.Printf("\nArguments: %v\n\n", args.Interfaces())
	} else {
		fmt.Print("\n")
	}

	sqlStr, err = dbr.Interpolate(sqlStr, args...)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Println("Preprocessed Statement:")
	wordwrap.Fstring(os.Stdout, sqlStr, 80)
}

func ExampleNewInsert() {
	i := dbr.NewInsert("tableA").
		AddColumns("b", "c", "d", "e").
		AddValues(dbr.ArgInt(1), dbr.ArgInt64(2), dbr.ArgString("Three"), dbr.ArgNull()).
		AddValues(dbr.ArgInt(5), dbr.ArgInt64(6), dbr.ArgString("Seven"), dbr.ArgFloat64(3.14156))
	writeToSqlAndPreprocess(i)

	// Output:
	//Prepared Statement:
	//INSERT INTO `tableA` (`b`,`c`,`d`,`e`) VALUES (?,?,?,?),(?,?,?,?)
	//Arguments: [1 2 Three <nil> 5 6 Seven 3.14156]
	//
	//Preprocessed Statement:
	//INSERT INTO `tableA` (`b`,`c`,`d`,`e`) VALUES
	//(1,2,'Three',NULL),(5,6,'Seven',3.14156)
}

func ExampleNewInsert_withoutColumns() {
	i := dbr.NewInsert("catalog_product_link").
		AddValues(dbr.ArgInt64(2046), dbr.ArgInt64(33), dbr.ArgInt64(3)).
		AddValues(dbr.ArgInt64(2046), dbr.ArgInt64(34), dbr.ArgInt64(3)).
		AddValues(dbr.ArgInt64(2046), dbr.ArgInt64(35), dbr.ArgInt64(3))
	writeToSqlAndPreprocess(i)

	// Output:
	//Prepared Statement:
	//INSERT INTO `catalog_product_link` VALUES (?,?,?),(?,?,?),(?,?,?)
	//Arguments: [2046 33 3 2046 34 3 2046 35 3]
	//
	//Preprocessed Statement:
	//INSERT INTO `catalog_product_link` VALUES (2046,33,3),(2046,34,3),(2046,35,3)
}

func ExampleInsert_AddValues() {
	// Without any columns you must for each row call AddValues. Here we insert
	// three rows at once.
	i := dbr.NewInsert("catalog_product_link").
		AddValues(dbr.ArgInt64(2046), dbr.ArgInt64(33), dbr.ArgInt64(3)).
		AddValues(dbr.ArgInt64(2046), dbr.ArgInt64(34), dbr.ArgInt64(3)).
		AddValues(dbr.ArgInt64(2046), dbr.ArgInt64(35), dbr.ArgInt64(3))
	writeToSqlAndPreprocess(i)
	fmt.Print("\n\n")

	// Specifying columns allows to call only one time AddValues but inserting
	// three rows at once. Of course you can also insert only one row ;-)
	i = dbr.NewInsert("catalog_product_link").
		AddColumns("product_id", "linked_product_id", "link_type_id").
		AddValues(
			dbr.ArgInt64(2046), dbr.ArgInt64(33), dbr.ArgInt64(3),
			dbr.ArgInt64(2046), dbr.ArgInt64(34), dbr.ArgInt64(3),
			dbr.ArgInt64(2046), dbr.ArgInt64(35), dbr.ArgInt64(3),
		)
	writeToSqlAndPreprocess(i)

	// Output:
	//Prepared Statement:
	//INSERT INTO `catalog_product_link` VALUES (?,?,?),(?,?,?),(?,?,?)
	//Arguments: [2046 33 3 2046 34 3 2046 35 3]
	//
	//Preprocessed Statement:
	//INSERT INTO `catalog_product_link` VALUES (2046,33,3),(2046,34,3),(2046,35,3)
	//
	//Prepared Statement:
	//INSERT INTO `catalog_product_link`
	//(`product_id`,`linked_product_id`,`link_type_id`) VALUES (?,?,?),(?,?,?),(?,?,?)
	//Arguments: [2046 33 3 2046 34 3 2046 35 3]
	//
	//Preprocessed Statement:
	//INSERT INTO `catalog_product_link`
	//(`product_id`,`linked_product_id`,`link_type_id`) VALUES
	//(2046,33,3),(2046,34,3),(2046,35,3)
}

func ExampleInsert_AddOnDuplicateKey() {
	i := dbr.NewInsert("dbr_people").
		AddColumns("id", "name", "email").
		AddValues(dbr.ArgInt64(1), dbr.ArgString("Pik'e"), dbr.ArgString("pikes@peak.com")).
		AddOnDuplicateKey("name", dbr.ArgString("Pik3")).
		AddOnDuplicateKey("email", nil)
	writeToSqlAndPreprocess(i)

	// Output:
	//Prepared Statement:
	//INSERT INTO `dbr_people` (`id`,`name`,`email`) VALUES (?,?,?) ON DUPLICATE KEY
	//UPDATE `name`=?, `email`=VALUES(`email`)
	//Arguments: [1 Pik'e pikes@peak.com Pik3]
	//
	//Preprocessed Statement:
	//INSERT INTO `dbr_people` (`id`,`name`,`email`) VALUES
	//(1,'Pik\'e','pikes@peak.com') ON DUPLICATE KEY UPDATE `name`='Pik3',
	//`email`=VALUES(`email`)
}

func ExampleInsert_FromSelect() {
	ins := dbr.NewInsert("tableA")

	argEq := dbr.Eq{"int64B": dbr.In.Int64(1, 2, 3)}

	ins.FromSelect(
		dbr.NewSelect().AddColumns("something_id", "user_id").
			AddColumns("other").
			From("some_table").
			Where(
				dbr.ParenthesisOpen(),
				dbr.Column("int64A", dbr.GreaterOrEqual.Int64(1)),
				dbr.Column("string", dbr.ArgString("wat")).Or(),
				dbr.ParenthesisClose(),
			).
			Where(argEq).
			OrderByDesc("id").
			Paginate(1, 20),
	)
	writeToSqlAndPreprocess(ins)
	// Output:
	//Prepared Statement:
	//INSERT INTO `tableA` SELECT `something_id`, `user_id`, `other` FROM `some_table`
	//WHERE ((`int64A` >= ?) OR (`string` = ?)) AND (`int64B` IN ?) ORDER BY `id` DESC
	//LIMIT 20 OFFSET 0
	//Arguments: [1 wat 1 2 3]
	//
	//Preprocessed Statement:
	//INSERT INTO `tableA` SELECT `something_id`, `user_id`, `other` FROM `some_table`
	//WHERE ((`int64A` >= 1) OR (`string` = 'wat')) AND (`int64B` IN (1,2,3)) ORDER BY
	//`id` DESC LIMIT 20 OFFSET 0
}

func ExampleNewDelete() {
	d := dbr.NewDelete("tableA").Where(
		dbr.Column("a", dbr.Like.Str("b'%")),
		dbr.Column("b", dbr.In.Int(3, 4, 5, 6)),
	).
		Limit(1).OrderBy("id")
	writeToSqlAndPreprocess(d)
	// Output:
	//Prepared Statement:
	//DELETE FROM `tableA` WHERE (`a` LIKE ?) AND (`b` IN ?) ORDER BY `id` LIMIT 1
	//Arguments: [b'% 3 4 5 6]
	//
	//Preprocessed Statement:
	//DELETE FROM `tableA` WHERE (`a` LIKE 'b\'%') AND (`b` IN (3,4,5,6)) ORDER BY
	//`id` LIMIT 1
}

// ExampleNewUnion constructs a UNION with three SELECTs. It preserves the
// results sets of each SELECT by simply adding an internal index to the columns
// list and sort ascending with the internal index.
func ExampleNewUnion() {

	u := dbr.NewUnion(
		dbr.NewSelect().AddColumnsAlias("a1", "A", "a2", "B").From("tableA").Where(dbr.Column("a1", dbr.ArgInt64(3))),
		dbr.NewSelect().AddColumnsAlias("b1", "A", "b2", "B").From("tableB").Where(dbr.Column("b1", dbr.ArgInt64(4))),
	)
	// Maybe more of your code ...
	u.Append(
		dbr.NewSelect().AddColumnsExprAlias("concat(c1,?,c2)", "A").
			AddArguments(dbr.ArgString("-")).
			AddColumnsAlias("c2", "B").
			From("tableC").Where(dbr.Column("c2", dbr.Equal.Str("ArgForC2"))),
	).
		OrderBy("A").       // Ascending by A
		OrderByDesc("B").   // Descending by B
		All().              // Enables UNION ALL syntax
		PreserveResultSet() // Maintains the correct order of the result set for all SELECTs.
	// Note that the final ORDER BY statement of a UNION creates a temporary
	// table in MySQL.
	writeToSqlAndPreprocess(u)
	// Output:
	//Prepared Statement:
	//(SELECT `a1` AS `A`, `a2` AS `B`, 0 AS `_preserve_result_set` FROM `tableA`
	//WHERE (`a1` = ?))
	//UNION ALL
	//(SELECT `b1` AS `A`, `b2` AS `B`, 1 AS `_preserve_result_set` FROM `tableB`
	//WHERE (`b1` = ?))
	//UNION ALL
	//(SELECT concat(c1,?,c2) AS `A`, `c2` AS `B`, 2 AS `_preserve_result_set` FROM
	//`tableC` WHERE (`c2` = ?))
	//ORDER BY `_preserve_result_set`, `A` ASC, `B` DESC
	//Arguments: [3 4 - ArgForC2]
	//
	//Preprocessed Statement:
	//(SELECT `a1` AS `A`, `a2` AS `B`, 0 AS `_preserve_result_set` FROM `tableA`
	//WHERE (`a1` = 3))
	//UNION ALL
	//(SELECT `b1` AS `A`, `b2` AS `B`, 1 AS `_preserve_result_set` FROM `tableB`
	//WHERE (`b1` = 4))
	//UNION ALL
	//(SELECT concat(c1,'-',c2) AS `A`, `c2` AS `B`, 2 AS `_preserve_result_set` FROM
	//`tableC` WHERE (`c2` = 'ArgForC2'))
	//ORDER BY `_preserve_result_set`, `A` ASC, `B` DESC
}

func ExampleNewUnionTemplate() {

	u := dbr.NewUnionTemplate(
		dbr.NewSelect().AddColumns("t.value", "t.attribute_id", "t.store_id").From("catalog_product_entity_{type}", "t").
			Where(dbr.Column("entity_id", dbr.ArgInt64(1561)), dbr.Column("store_id", dbr.In.Int64(1, 0))),
	).
		StringReplace("{type}", "varchar", "int", "decimal", "datetime", "text").
		PreserveResultSet().
		All().OrderBy("attribute_id", "store_id")
	writeToSqlAndPreprocess(u)
	// Output:
	//Prepared Statement:
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS
	//`_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE
	//(`entity_id` = ?) AND (`store_id` IN ?))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS
	//`_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE
	//(`entity_id` = ?) AND (`store_id` IN ?))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS
	//`_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE
	//(`entity_id` = ?) AND (`store_id` IN ?))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS
	//`_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE
	//(`entity_id` = ?) AND (`store_id` IN ?))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS
	//`_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE
	//(`entity_id` = ?) AND (`store_id` IN ?))
	//ORDER BY `_preserve_result_set`, `attribute_id` ASC, `store_id` ASC
	//Arguments: [1561 1 0 1561 1 0 1561 1 0 1561 1 0 1561 1 0]
	//
	//Preprocessed Statement:
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS
	//`_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE
	//(`entity_id` = 1561) AND (`store_id` IN (1,0)))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS
	//`_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE
	//(`entity_id` = 1561) AND (`store_id` IN (1,0)))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS
	//`_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE
	//(`entity_id` = 1561) AND (`store_id` IN (1,0)))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS
	//`_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE
	//(`entity_id` = 1561) AND (`store_id` IN (1,0)))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS
	//`_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE
	//(`entity_id` = 1561) AND (`store_id` IN (1,0)))
	//ORDER BY `_preserve_result_set`, `attribute_id` ASC, `store_id` ASC
}

// ExampleUnionTemplate_Interpolate interpolates the SQL string with its
// placeholders and puts for each placeholder the correct encoded and escaped
// value into it. Eliminates the need for prepared statements by sending in one
// round trip the query and its arguments directly to the database server. If
// you execute a query multiple times within a short time you should use
// prepared statements.
func ExampleUnionTemplate_Interpolate() {

	u := dbr.NewUnionTemplate(
		dbr.NewSelect().AddColumns("t.value", "t.attribute_id", "t.store_id").From("catalog_product_entity_{type}", "t").
			Where(dbr.Column("entity_id", dbr.ArgInt64(1561)), dbr.Column("store_id", dbr.In.Int64(1, 0))),
	).
		StringReplace("{type}", "varchar", "int", "decimal", "datetime", "text").
		PreserveResultSet().
		All().OrderBy("attribute_id", "store_id")
	writeToSqlAndPreprocess(u)
	// Output:
	//Prepared Statement:
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS
	//`_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE
	//(`entity_id` = ?) AND (`store_id` IN ?))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS
	//`_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE
	//(`entity_id` = ?) AND (`store_id` IN ?))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS
	//`_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE
	//(`entity_id` = ?) AND (`store_id` IN ?))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS
	//`_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE
	//(`entity_id` = ?) AND (`store_id` IN ?))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS
	//`_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE
	//(`entity_id` = ?) AND (`store_id` IN ?))
	//ORDER BY `_preserve_result_set`, `attribute_id` ASC, `store_id` ASC
	//Arguments: [1561 1 0 1561 1 0 1561 1 0 1561 1 0 1561 1 0]
	//
	//Preprocessed Statement:
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 0 AS
	//`_preserve_result_set` FROM `catalog_product_entity_varchar` AS `t` WHERE
	//(`entity_id` = 1561) AND (`store_id` IN (1,0)))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 1 AS
	//`_preserve_result_set` FROM `catalog_product_entity_int` AS `t` WHERE
	//(`entity_id` = 1561) AND (`store_id` IN (1,0)))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 2 AS
	//`_preserve_result_set` FROM `catalog_product_entity_decimal` AS `t` WHERE
	//(`entity_id` = 1561) AND (`store_id` IN (1,0)))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 3 AS
	//`_preserve_result_set` FROM `catalog_product_entity_datetime` AS `t` WHERE
	//(`entity_id` = 1561) AND (`store_id` IN (1,0)))
	//UNION ALL
	//(SELECT `t`.`value`, `t`.`attribute_id`, `t`.`store_id`, 4 AS
	//`_preserve_result_set` FROM `catalog_product_entity_text` AS `t` WHERE
	//(`entity_id` = 1561) AND (`store_id` IN (1,0)))
	//ORDER BY `_preserve_result_set`, `attribute_id` ASC, `store_id` ASC
}

func ExampleInterpolate() {
	sqlStr, err := dbr.Interpolate("SELECT * FROM x WHERE a IN ? AND b IN ? AND c NOT IN ? AND d BETWEEN ? AND ?",
		dbr.In.Int(1),
		dbr.In.Int(1, 2, 3),
		dbr.In.Int64(5, 6, 7),
		dbr.Between.Str("wat", "ok"),
	)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Printf("%s\n", sqlStr)
	// Output:
	// SELECT * FROM x WHERE a IN (1) AND b IN (1,2,3) AND c NOT IN (5,6,7) AND d BETWEEN 'wat' AND 'ok'
}

func ExampleRepeat() {
	sl := []string{"a", "b", "c", "d", "e"}

	sqlStr, args, err := dbr.Repeat("SELECT * FROM `table` WHERE id IN (?) AND name IN (?)",
		dbr.ArgInt(5, 7, 9), dbr.ArgString(sl...))

	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	fmt.Printf("%s\nArguments: %v\n", sqlStr, args)
	// Output:
	// SELECT * FROM `table` WHERE id IN (?,?,?) AND name IN (?,?,?,?,?)
	// Arguments: [5 7 9 a b c d e]
}

// ExampleArgument is duplicate of ExampleColumn
func ExampleArgument() {

	argPrinter := func(arg ...dbr.Argument) {
		sqlStr, args, err := dbr.NewSelect().AddColumns("a", "b").
			From("c").Where(dbr.Column("d", arg...)).ToSQL()
		if err != nil {
			fmt.Printf("%+v\n", err)
		} else {
			fmt.Printf("%q", sqlStr)
			if len(args) > 0 {
				fmt.Printf(" Arguments: %v", args.Interfaces())
			}
			fmt.Print("\n")
		}
	}

	argPrinter(dbr.ArgNull())
	argPrinter(dbr.NotNull.Null())
	argPrinter(dbr.ArgInt(2))
	argPrinter(dbr.Null.Int(3))
	argPrinter(dbr.NotNull.Int(4))
	argPrinter(dbr.In.Int(7, 8, 9))
	argPrinter(dbr.NotIn.Int(10, 11, 12))
	argPrinter(dbr.Between.Int(13, 14))
	argPrinter(dbr.NotBetween.Int(15, 16))
	argPrinter(dbr.Greatest.Int(17, 18, 19))
	argPrinter(dbr.Least.Int(20, 21, 22))
	argPrinter(dbr.Equal.Int(30))
	argPrinter(dbr.NotEqual.Int(31))

	argPrinter(dbr.Less.Int(32))
	argPrinter(dbr.Greater.Int(33))
	argPrinter(dbr.LessOrEqual.Int(34))
	argPrinter(dbr.GreaterOrEqual.Int(35))

	argPrinter(dbr.Like.Str("Goph%"))
	argPrinter(dbr.NotLike.Str("Cat%"))

	//Output:
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NOT NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)" Arguments: [2]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NOT NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IN ?)" Arguments: [7 8 9]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT IN ?)" Arguments: [10 11 12]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` BETWEEN ? AND ?)" Arguments: [13 14]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT BETWEEN ? AND ?)" Arguments: [15 16]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` GREATEST (?))" Arguments: [17 18 19]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` LEAST (?))" Arguments: [20 21 22]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)" Arguments: [30]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` != ?)" Arguments: [31]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` < ?)" Arguments: [32]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` > ?)" Arguments: [33]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` <= ?)" Arguments: [34]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` >= ?)" Arguments: [35]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` LIKE ?)" Arguments: [Goph%]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT LIKE ?)" Arguments: [Cat%]
}

// ExampleColumn is a duplicate of ExampleArgument
func ExampleColumn() {
	argPrinter := func(arg ...dbr.Argument) {
		sqlStr, args, err := dbr.NewSelect().AddColumns("a", "b").
			From("c").Where(dbr.Column("d", arg...)).ToSQL()
		if err != nil {
			fmt.Printf("%+v\n", err)
		} else {
			fmt.Printf("%q", sqlStr)
			if len(args) > 0 {
				fmt.Printf(" Arguments: %v", args.Interfaces())
			}
			fmt.Print("\n")
		}
	}

	argPrinter(dbr.ArgNull())
	argPrinter(dbr.NotNull.Null())
	argPrinter(dbr.ArgInt(2))
	argPrinter(dbr.Null.Int(3))
	argPrinter(dbr.NotNull.Int(4))
	argPrinter(dbr.In.Int(7, 8, 9))
	argPrinter(dbr.NotIn.Int(10, 11, 12))
	argPrinter(dbr.Between.Int(13, 14))
	argPrinter(dbr.NotBetween.Int(15, 16))
	argPrinter(dbr.Greatest.Int(17, 18, 19))
	argPrinter(dbr.Least.Int(20, 21, 22))
	argPrinter(dbr.Equal.Int(30))
	argPrinter(dbr.NotEqual.Int(31))

	argPrinter(dbr.Less.Int(32))
	argPrinter(dbr.Greater.Int(33))
	argPrinter(dbr.LessOrEqual.Int(34))
	argPrinter(dbr.GreaterOrEqual.Int(35))

	argPrinter(dbr.Like.Str("Goph%"))
	argPrinter(dbr.NotLike.Str("Cat%"))

	//Output:
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NOT NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)" Arguments: [2]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IS NOT NULL)"
	//"SELECT `a`, `b` FROM `c` WHERE (`d` IN ?)" Arguments: [7 8 9]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT IN ?)" Arguments: [10 11 12]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` BETWEEN ? AND ?)" Arguments: [13 14]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT BETWEEN ? AND ?)" Arguments: [15 16]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` GREATEST (?))" Arguments: [17 18 19]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` LEAST (?))" Arguments: [20 21 22]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)" Arguments: [30]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` != ?)" Arguments: [31]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` < ?)" Arguments: [32]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` > ?)" Arguments: [33]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` <= ?)" Arguments: [34]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` >= ?)" Arguments: [35]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` LIKE ?)" Arguments: [Goph%]
	//"SELECT `a`, `b` FROM `c` WHERE (`d` NOT LIKE ?)" Arguments: [Cat%]
}

func ExampleSubSelect() {
	s := dbr.NewSelect("sku", "type_id").
		From("catalog_product_entity").
		Where(dbr.SubSelect(
			"entity_id", dbr.In,
			dbr.NewSelect().From("catalog_category_product").
				AddColumns("entity_id").Where(dbr.Column("category_id", dbr.ArgInt64(234))),
		))
	writeToSqlAndPreprocess(s)
	// Output:
	//Prepared Statement:
	//SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` IN
	//(SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` = ?)))
	//Arguments: [234]
	//
	//Preprocessed Statement:
	//SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` IN
	//(SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` =
	//234)))
}

func ExampleNewSelectFromSub() {
	sel3 := dbr.NewSelect().From("sales_bestsellers_aggregated_daily", "t3").
		AddColumnsExprAlias("DATE_FORMAT(t3.period, '%Y-%m-01')", "period").
		AddColumns("t3.store_id", "t3.product_id", "t3.product_name").
		AddColumnsExprAlias("AVG(`t3`.`product_price`)", "avg_price", "SUM(t3.qty_ordered)", "total_qty").
		Where(dbr.Column("product_name", dbr.ArgString("Canon%"))).
		GroupBy("t3.store_id").
		GroupByExpr("DATE_FORMAT(t3.period, '%Y-%m-01')").
		GroupBy("t3.product_id", "t3.product_name").
		OrderBy("t3.store_id").
		OrderByExpr("DATE_FORMAT(t3.period, '%Y-%m-01')").
		OrderByDesc("total_qty")

	sel1 := dbr.NewSelectFromSub(sel3, "t1").
		AddColumns("t1.period", "t1.store_id", "t1.product_id", "t1.product_name", "t1.avg_price", "t1.qty_ordered").
		Where(dbr.Column("product_name", dbr.ArgString("Sony%"))).
		OrderBy("t1.period", "t1.product_id")
	writeToSqlAndPreprocess(sel1)
	// Output:
	//Prepared Statement:
	//SELECT `t1`.`period`, `t1`.`store_id`, `t1`.`product_id`, `t1`.`product_name`,
	//`t1`.`avg_price`, `t1`.`qty_ordered` FROM (SELECT DATE_FORMAT(t3.period,
	//'%Y-%m-01') AS `period`, `t3`.`store_id`, `t3`.`product_id`,
	//`t3`.`product_name`, AVG(`t3`.`product_price`) AS `avg_price`,
	//SUM(t3.qty_ordered) AS `total_qty` FROM `sales_bestsellers_aggregated_daily` AS
	//`t3` WHERE (`product_name` = ?) GROUP BY `t3`.`store_id`, DATE_FORMAT(t3.period,
	//'%Y-%m-01'), `t3`.`product_id`, `t3`.`product_name` ORDER BY `t3`.`store_id`,
	//DATE_FORMAT(t3.period, '%Y-%m-01'), `total_qty` DESC) AS `t1` WHERE
	//(`product_name` = ?) ORDER BY `t1`.`period`, `t1`.`product_id`
	//Arguments: [Canon% Sony%]
	//
	//Preprocessed Statement:
	//SELECT `t1`.`period`, `t1`.`store_id`, `t1`.`product_id`, `t1`.`product_name`,
	//`t1`.`avg_price`, `t1`.`qty_ordered` FROM (SELECT DATE_FORMAT(t3.period,
	//'%Y-%m-01') AS `period`, `t3`.`store_id`, `t3`.`product_id`,
	//`t3`.`product_name`, AVG(`t3`.`product_price`) AS `avg_price`,
	//SUM(t3.qty_ordered) AS `total_qty` FROM `sales_bestsellers_aggregated_daily` AS
	//`t3` WHERE (`product_name` = 'Canon%') GROUP BY `t3`.`store_id`,
	//DATE_FORMAT(t3.period, '%Y-%m-01'), `t3`.`product_id`, `t3`.`product_name` ORDER
	//BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `total_qty` DESC) AS
	//`t1` WHERE (`product_name` = 'Sony%') ORDER BY `t1`.`period`, `t1`.`product_id`
}

func ExampleSQLIfNull() {
	fmt.Println(dbr.SQLIfNull("column1"))
	fmt.Println(dbr.SQLIfNull("table1.column1"))
	fmt.Println(dbr.SQLIfNull("column1", "column2"))
	fmt.Println(dbr.SQLIfNull("table1.column1", "table2.column2"))
	fmt.Println(dbr.SQLIfNull("column2", "1/0", "alias"))
	fmt.Println(dbr.SQLIfNull("SELECT * FROM x", "8", "alias"))
	fmt.Println(dbr.SQLIfNull("SELECT * FROM x", "9 ", "alias"))
	fmt.Println(dbr.SQLIfNull("column1", "column2", "alias"))
	fmt.Println(dbr.SQLIfNull("table1.column1", "table2.column2", "alias"))
	fmt.Println(dbr.SQLIfNull("table1", "column1", "table2", "column2"))
	fmt.Println(dbr.SQLIfNull("table1", "column1", "table2", "column2", "alias"))
	fmt.Println(dbr.SQLIfNull("table1", "column1", "table2", "column2", "alias", "x"))
	fmt.Println(dbr.SQLIfNull("table1", "column1", "table2", "column2", "alias", "x", "y"))
	//Output:
	//IFNULL(`column1`,(NULL ))
	//IFNULL(`table1`.`column1`,(NULL ))
	//IFNULL(`column1`,`column2`)
	//IFNULL(`table1`.`column1`,`table2`.`column2`)
	//IFNULL(`column2`,(1/0)) AS `alias`
	//IFNULL((SELECT * FROM x),`8`) AS `alias`
	//IFNULL((SELECT * FROM x),(9 )) AS `alias`
	//IFNULL(`column1`,`column2`) AS `alias`
	//IFNULL(`table1`.`column1`,`table2`.`column2`) AS `alias`
	//IFNULL(`table1`.`column1`,`table2`.`column2`)
	//IFNULL(`table1`.`column1`,`table2`.`column2`) AS `alias`
	//IFNULL(`table1`.`column1`,`table2`.`column2`) AS `alias_x`
	//IFNULL(`table1`.`column1`,`table2`.`column2`) AS `alias_x_y`
}

func ExampleSQLIf() {
	s := dbr.NewSelect().AddColumns("a", "b", "c").
		From("table1").Where(
		dbr.Expression(
			dbr.SQLIf("a > 0", "b", "c"),
			dbr.Greater.Int(4711),
		))
	writeToSqlAndPreprocess(s)

	// Output:
	//Prepared Statement:
	//SELECT `a`, `b`, `c` FROM `table1` WHERE (IF((a > 0), b, c) > ?)
	//Arguments: [4711]
	//
	//Preprocessed Statement:
	//SELECT `a`, `b`, `c` FROM `table1` WHERE (IF((a > 0), b, c) > 4711)
}

func ExampleSQLCase_update() {
	u := dbr.NewUpdate("cataloginventory_stock_item").
		Set("qty", dbr.ArgExpr(dbr.SQLCase("`product_id`", "qty",
			"3456", "qty+?",
			"3457", "qty+?",
			"3458", "qty+?",
		), dbr.ArgInt(3, 4, 5))).
		Where(
			dbr.Column("product_id", dbr.In.Int64(345, 567, 897)),
			dbr.Column("website_id", dbr.ArgInt64(6)),
		)
	writeToSqlAndPreprocess(u)

	// Output:
	//Prepared Statement:
	//UPDATE `cataloginventory_stock_item` SET `qty`=CASE `product_id` WHEN 3456 THEN
	//qty+? WHEN 3457 THEN qty+? WHEN 3458 THEN qty+? ELSE qty END WHERE (`product_id`
	//IN ?) AND (`website_id` = ?)
	//Arguments: [3 4 5 345 567 897 6]
	//
	//Preprocessed Statement:
	//UPDATE `cataloginventory_stock_item` SET `qty`=CASE `product_id` WHEN 3456 THEN
	//qty+3 WHEN 3457 THEN qty+4 WHEN 3458 THEN qty+5 ELSE qty END WHERE (`product_id`
	//IN (345,567,897)) AND (`website_id` = 6)
}

// ExampleSQLCase_select is a duplicate of ExampleSelect_AddArguments
func ExampleSQLCase_select() {
	// time stamp has no special meaning ;-)
	start := dbr.ArgTime(time.Unix(1257894000, 0))
	end := dbr.ArgTime(time.Unix(1257980400, 0))
	s := dbr.NewSelect().AddColumns("price", "sku", "name", "title", "description").
		AddColumnsExprAlias(
			dbr.SQLCase("", "`closed`",
				"date_start <= ? AND date_end >= ?", "`open`",
				"date_start > ? AND date_end > ?", "`upcoming`",
			),
			"is_on_sale",
		).
		AddArguments(start, end, start, end).
		From("catalog_promotions").Where(
		dbr.Column("promotion_id", dbr.NotIn.Int(4711, 815, 42)))
	writeToSqlAndPreprocess(s)

	// Output:
	//Prepared Statement:
	//SELECT `price`, `sku`, `name`, `title`, `description`, CASE  WHEN date_start <=
	//? AND date_end >= ? THEN `open` WHEN date_start > ? AND date_end > ? THEN
	//`upcoming` ELSE `closed` END AS `is_on_sale` FROM `catalog_promotions` WHERE
	//(`promotion_id` NOT IN ?)
	//Arguments: [2009-11-11 00:00:00 +0100 CET 2009-11-12 00:00:00 +0100 CET 2009-11-11 00:00:00 +0100 CET 2009-11-12 00:00:00 +0100 CET 4711 815 42]
	//
	//Preprocessed Statement:
	//SELECT `price`, `sku`, `name`, `title`, `description`, CASE  WHEN date_start <=
	//'2009-11-11 00:00:00' AND date_end >= '2009-11-12 00:00:00' THEN `open` WHEN
	//date_start > '2009-11-11 00:00:00' AND date_end > '2009-11-12 00:00:00' THEN
	//`upcoming` ELSE `closed` END AS `is_on_sale` FROM `catalog_promotions` WHERE
	//(`promotion_id` NOT IN (4711,815,42))
}

// ExampleSelect_AddArguments is duplicate of ExampleSQLCase_select
func ExampleSelect_AddArguments() {
	// time stamp has no special meaning ;-)
	start := dbr.ArgTime(time.Unix(1257894000, 0))
	end := dbr.ArgTime(time.Unix(1257980400, 0))
	s := dbr.NewSelect().AddColumns("price", "sku", "name", "title", "description").
		AddColumnsExprAlias(
			dbr.SQLCase("", "`closed`",
				"date_start <= ? AND date_end >= ?", "`open`",
				"date_start > ? AND date_end > ?", "`upcoming`",
			),
			"is_on_sale",
		).
		AddArguments(start, end, start, end).
		From("catalog_promotions").Where(
		dbr.Column("promotion_id", dbr.NotIn.Int(4711, 815, 42)))
	writeToSqlAndPreprocess(s)

	// Output:
	//Prepared Statement:
	//SELECT `price`, `sku`, `name`, `title`, `description`, CASE  WHEN date_start <=
	//? AND date_end >= ? THEN `open` WHEN date_start > ? AND date_end > ? THEN
	//`upcoming` ELSE `closed` END AS `is_on_sale` FROM `catalog_promotions` WHERE
	//(`promotion_id` NOT IN ?)
	//Arguments: [2009-11-11 00:00:00 +0100 CET 2009-11-12 00:00:00 +0100 CET 2009-11-11 00:00:00 +0100 CET 2009-11-12 00:00:00 +0100 CET 4711 815 42]
	//
	//Preprocessed Statement:
	//SELECT `price`, `sku`, `name`, `title`, `description`, CASE  WHEN date_start <=
	//'2009-11-11 00:00:00' AND date_end >= '2009-11-12 00:00:00' THEN `open` WHEN
	//date_start > '2009-11-11 00:00:00' AND date_end > '2009-11-12 00:00:00' THEN
	//`upcoming` ELSE `closed` END AS `is_on_sale` FROM `catalog_promotions` WHERE
	//(`promotion_id` NOT IN (4711,815,42))
}

func ExampleParenthesisOpen() {
	s := dbr.NewSelect("columnA", "columnB").
		Distinct().
		From("tableC", "ccc").
		Where(
			dbr.ParenthesisOpen(),
			dbr.Column("d", dbr.ArgInt(1)),
			dbr.Column("e", dbr.ArgString("wat")).Or(),
			dbr.ParenthesisClose(),
			dbr.Eq{"f": dbr.ArgInt(2)},
		).
		GroupBy("ab").
		Having(
			dbr.Expression("j = k"),
			dbr.ParenthesisOpen(),
			dbr.Column("m", dbr.ArgInt(33)),
			dbr.Column("n", dbr.ArgString("wh3r3")).Or(),
			dbr.ParenthesisClose(),
		).
		OrderBy("l").
		Limit(7).
		Offset(8)
	writeToSqlAndPreprocess(s)

	// Output:
	//Prepared Statement:
	//SELECT DISTINCT `columnA`, `columnB` FROM `tableC` AS `ccc` WHERE ((`d` = ?) OR
	//(`e` = ?)) AND (`f` = ?) GROUP BY `ab` HAVING (j = k) AND ((`m` = ?) OR (`n` =
	//?)) ORDER BY `l` LIMIT 7 OFFSET 8
	//Arguments: [1 wat 2 33 wh3r3]
	//
	//Preprocessed Statement:
	//SELECT DISTINCT `columnA`, `columnB` FROM `tableC` AS `ccc` WHERE ((`d` = 1) OR
	//(`e` = 'wat')) AND (`f` = 2) GROUP BY `ab` HAVING (j = k) AND ((`m` = 33) OR
	//(`n` = 'wh3r3')) ORDER BY `l` LIMIT 7 OFFSET 8
}