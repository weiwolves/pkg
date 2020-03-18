package dmltestgenerated4

import (
	"context"
	"sort"
	"testing"

	"github.com/weiwolves/pkg/sql/ddl"
	"github.com/weiwolves/pkg/sql/dml"
	"github.com/weiwolves/pkg/sql/dmltest"
	"github.com/weiwolves/pkg/util/assert"
	"github.com/weiwolves/pkg/util/pseudo"
)

func TestNewDBManager_Manual(t *testing.T) {
	// var logbuf bytes.Buffer
	// defer func() { println("\n", logbuf.String(), "\n") }()
	// l := logw.NewLog(logw.WithLevel(logw.LevelDebug), logw.WithWriter(&logbuf))
	// db := dmltest.MustConnectDB(t, dml.WithLogger(l, shortid.MustGenerate))

	db := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, db)
	defer dmltest.SQLDumpLoad(t, "../testdata/test_*_tables.sql", &dmltest.SQLDumpOptions{
		SkipDBCleanup: true,
	}).Deferred()

	availableEvents := []dml.EventFlag{
		dml.EventFlagBeforeInsert, dml.EventFlagAfterInsert,
		dml.EventFlagBeforeUpsert, dml.EventFlagAfterUpsert,
		dml.EventFlagBeforeSelect, dml.EventFlagAfterSelect,
	}
	const (
		eventIdxEntity = iota
		eventIdxCollection
		eventIdxMax
	)
	calledEvents := [eventIdxMax][dml.EventFlagMax]int{}
	ctx := context.Background()

	opts := &DBMOption{
		TableOptions: []ddl.TableOption{ddl.WithConnPool(db)},
		InitSelectFn: func(d *dml.Select) *dml.Select {
			d.Limit(0, 1000) // adds to every SELECT the LIMIT clause, for testing purposes
			return d
		},
	}

	for _, eventID := range availableEvents {
		eventID := eventID
		opts = opts.AddEventCoreConfiguration(eventID, func(_ context.Context, cc *CoreConfigurations, c *CoreConfiguration) error {
			if cc != nil {
				calledEvents[eventIdxCollection][eventID]++ // set to 2 to verify that it has been called
			} else if c != nil {
				calledEvents[eventIdxEntity][eventID]++ // set to 2 to verify that it has been called
			}
			return nil
		})
	}

	// used for debugging or different query styles
	shouldInterpolateFn := func(dbr *dml.DBR) {
		// dbr.Interpolate()
	}

	dbm, err := NewDBManager(ctx, opts)
	assert.NoError(t, err)

	ps := pseudo.MustNewService(0, &pseudo.Options{Lang: "de", FloatMaxDecimals: 6, MaxLenStringLimit: 41})
	var entityInsertCount int64
	t.Run("Entity", func(t *testing.T) {
		var eFake CoreConfiguration // e=entity => entityFake or entityLoaded
		assert.NoError(t, ps.FakeData(&eFake))

		t.Run("Insert", func(t *testing.T) {
			res, err := eFake.Insert(ctx, dbm, shouldInterpolateFn)
			assert.NoError(t, err)

			assert.NoError(t, dml.ExecValidateOneAffectedRow(res, err))

			lid, _ := res.LastInsertId()
			ra, _ := res.RowsAffected()
			assert.True(t, lid > 0, "LastInsertID should be greater than 0")
			assert.True(t, ra > 0, "RowsAffected should be greater than 0")
			t.Logf("LastInsertId(%d) RowsAffected(%d)", lid, ra)
			entityInsertCount += ra
		})
		t.Run("Upsert", func(t *testing.T) {
			// this test, runs the ON DUPLICATE KEY clause as the table core_config_data has a unique key.
			res, err := eFake.Upsert(ctx, dbm, shouldInterpolateFn)
			assert.NoError(t, err)
			lid, _ := res.LastInsertId()
			ra, _ := res.RowsAffected()
			assert.True(t, lid == 0, "LastInsertID should be zero")
			assert.True(t, ra == 0, "RowsAffected should be zero")
			t.Logf("LastInsertId(%d) RowsAffected(%d)", lid, ra)
		})
		t.Run("Load", func(t *testing.T) {
			eLoaded := &CoreConfiguration{}
			err = eLoaded.Load(ctx, dbm, eFake.ConfigID)
			assert.NoError(t, err)
			assert.NotEmpty(t, eLoaded.ConfigID)
			assert.NotEmpty(t, eLoaded.Scope)
			assert.NotEmpty(t, eLoaded.Path)
		})
	})

	t.Run("Collection", func(t *testing.T) {
		var ec CoreConfigurations
		assert.NoError(t, ps.FakeData(&ec))
		ec.Each(func(c *CoreConfiguration) {
			c.ConfigID = 0
		}) // reset configIDs
		t.Run("DBInsert", func(t *testing.T) {
			res, err := ec.DBInsert(ctx, dbm, shouldInterpolateFn)
			assert.NoError(t, err)
			lid, _ := res.LastInsertId()
			ra, _ := res.RowsAffected()
			assert.True(t, lid > 0, "LastInsertID should be greater than 0")
			assert.True(t, ra > 0, "RowsAffected should be greater than 0")
			t.Logf("LastInsertId(%d) RowsAffected(%d) RowsIn:%d Len:%d", lid, ra, lid+ra, len(ec.Data))
		})
		t.Run("DBUpsert", func(t *testing.T) {
			// this test, runs the ON DUPLICATE KEY clause as the table core_config_data has a unique key.
			res, err := ec.DBUpsert(ctx, dbm, shouldInterpolateFn)
			assert.NoError(t, err)
			lid, _ := res.LastInsertId()
			ra, _ := res.RowsAffected()
			assert.True(t, lid == 0, "LastInsertID should be zero")
			assert.True(t, ra == 0, "RowsAffected should be zero")
			t.Logf("LastInsertId(%d) RowsAffected(%d)", lid, ra)
		})

		t.Run("validate auto increment", func(t *testing.T) {
			calls := 0
			ec.Each(func(c *CoreConfiguration) {
				assert.True(t, c.ConfigID > 0, "c.ConfigID must be greater than zero")
				calls++
			})
			assert.Exactly(t, calls, len(ec.Data), "Length of ec must be equal")
			t.Logf("calls %d == %d len(ec.Data)", calls, len(ec.Data))
		})
		t.Run("DBLoad All", func(t *testing.T) {
			var eca CoreConfigurations
			assert.NoError(t, eca.DBLoad(ctx, dbm, []uint32{}))
			assert.Exactly(t, len(ec.Data)+int(entityInsertCount), len(eca.Data), "former collection must have the same length as the loaded one")
		})
		t.Run("DBLoad partial IDs", func(t *testing.T) {
			var eca CoreConfigurations
			assert.NoError(t, eca.DBLoad(ctx, dbm, ec.ConfigIDs()))
			assert.Exactly(t, len(ec.Data), len(eca.Data), "former collection must have the same length as the loaded one")
		})
		t.Run("DBDelete", func(t *testing.T) {
			t.Skip("asdadasds")
			res, err := ec.DBDelete(ctx, dbm)
			assert.NoError(t, err)
			lid, _ := res.LastInsertId()
			ra, _ := res.RowsAffected()
			assert.True(t, lid == 0, "LastInsertID should be zero")
			assert.Exactly(t, int64(len(ec.Data)), ra, "RowsAffected should be same as ec.Data length")
			t.Logf("LastInsertId(%d) RowsAffected(%d)", lid, ra)
		})
	})

	t.Run("Events and cached queries", func(t *testing.T) {
		cq := dbm.CachedQueries()
		queries := make([]string, 0, len(cq)*2)
		for k, v := range cq {
			queries = append(queries, k+"::"+v)
		}
		sort.Strings(queries)

		assert.Exactly(t, []string{
			"CoreConfigurationDeleteByPK::DELETE FROM `core_configuration` WHERE (`config_id` IN ?)",
			"CoreConfigurationInsert::INSERT INTO `core_configuration` (`scope`,`scope_id`,`expires`,`path`,`value`) VALUES (?,?,?,?,?)",
			"CoreConfigurationSelectByPK::SELECT `config_id`, `scope`, `scope_id`, `expires`, `path`, `value` FROM `core_configuration` AS `main_table` WHERE (`config_id` = ?) LIMIT 0,1000",
			"CoreConfigurationUpdateByPK::UPDATE `core_configuration` SET `scope`=?, `scope_id`=?, `expires`=?, `path`=?, `value`=? WHERE (`config_id` IN ?)",
			"CoreConfigurationUpsertByPK::INSERT INTO `core_configuration` (`scope`,`scope_id`,`expires`,`path`,`value`) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE `scope`=VALUES(`scope`), `scope_id`=VALUES(`scope_id`), `expires`=VALUES(`expires`), `path`=VALUES(`path`), `value`=VALUES(`value`)",
			"CoreConfigurationsSelectAll::SELECT `config_id`, `scope`, `scope_id`, `expires`, `path`, `value` FROM `core_configuration` AS `main_table` LIMIT 0,1000",
			"CoreConfigurationsSelectByPK::SELECT `config_id`, `scope`, `scope_id`, `expires`, `path`, `value` FROM `core_configuration` AS `main_table` WHERE (`config_id` IN ?) LIMIT 0,1000",
			"SalesOrderStatusStateSelectByPK::SELECT `status`, `state`, `is_default`, `visible_on_front` FROM `sales_order_status_state` AS `main_table` WHERE ((`status`, `state`) = /*TUPLES=002*/) LIMIT 0,1000",
			"SalesOrderStatusStatesSelectAll::SELECT `status`, `state`, `is_default`, `visible_on_front` FROM `sales_order_status_state` AS `main_table` LIMIT 0,1000",
			"SalesOrderStatusStatesSelectByPK::SELECT `status`, `state`, `is_default`, `visible_on_front` FROM `sales_order_status_state` AS `main_table` WHERE ((`status`, `state`) IN /*TUPLES=002*/) LIMIT 0,1000",
			"ViewCustomerAutoIncrementSelectByPK::SELECT `ce_entity_id`, `email`, `firstname`, `lastname`, `city` FROM `view_customer_auto_increment` AS `main_table` WHERE (`ce_entity_id` = ?) LIMIT 0,1000",
			"ViewCustomerAutoIncrementsSelectAll::SELECT `ce_entity_id`, `email`, `firstname`, `lastname`, `city` FROM `view_customer_auto_increment` AS `main_table` LIMIT 0,1000",
			"ViewCustomerAutoIncrementsSelectByPK::SELECT `ce_entity_id`, `email`, `firstname`, `lastname`, `city` FROM `view_customer_auto_increment` AS `main_table` WHERE (`ce_entity_id` IN ?) LIMIT 0,1000",
		}, queries)

		for _, eventID := range availableEvents {
			assert.Exactly(t, 1, calledEvents[eventIdxEntity][eventID])
		}
	})
}
