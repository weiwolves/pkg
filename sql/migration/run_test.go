package migration_test

import (
	"testing"

	"github.com/weiwolves/pkg/sql/migration"
)

func TestRun(t *testing.T) {
	err := migration.Run()
	t.Error(err)
}
