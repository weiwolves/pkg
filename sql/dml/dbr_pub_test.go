package dml_test

import (
	"context"
	"testing"

	"github.com/weiwolves/pkg/sql/dml"
	"github.com/weiwolves/pkg/sql/dmltest"
	"github.com/weiwolves/pkg/util/assert"
)

func TestDBR_Prepare(t *testing.T) {
	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)
	ctx := context.Background()

	prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("SELECT * FROM `core_config_data` WHERE ((`entity_id`, `attribute_id`, `store_id`, `source_id`) IN ((?,?,?,?),(?,?,?,?)))"))
	prep.ExpectQuery().WithArgs(1, 2, 3, 4, 6, 7, 8, 9).
		WillReturnRows(dmltest.MustMockRows(dmltest.WithFile("testdata/core_config_data.csv")))

	dbrStmt, err := dbc.SelectFrom("core_config_data").Star().Where(
		dml.Columns("entity_id", "attribute_id", "store_id", "source_id").In().Tuples(),
	).WithDBR().TupleCount(4, 2).ExpandPlaceHolders().Prepare(ctx)

	assert.NoError(t, err)

	ccd := &TableCoreConfigDataSlice{}
	rc, err := dbrStmt.Load(ctx, ccd, 1, 2, 3, 4, 6, 7, 8, 9)
	assert.NoError(t, err)
	assert.Exactly(t, 7, int(rc))
}
