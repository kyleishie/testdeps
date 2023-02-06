package tests

import (
	"github.com/kyleishie/testdeps/pkg/testsql/testpostgres"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContainer_NewDatabaseWithContext(t *testing.T) {
	con := testpostgres.RunForTest(t)

	t.Run("migrations work", func(t *testing.T) {
		db := con.NewTestDatabase(t, "testdata/migrations")

		rows, err := db.Query("SELECT * FROM users")
		assert.NoError(t, err)
		columns, err := rows.ColumnTypes()
		assert.NoError(t, err)

		assert.Len(t, columns, 2)

		// Email Column
		assert.Equal(t, "email", columns[0].Name())
		assert.Equal(t, "VARCHAR", columns[0].DatabaseTypeName())

		// Name Column
		assert.Equal(t, "name", columns[1].Name())
		assert.Equal(t, "VARCHAR", columns[1].DatabaseTypeName())
	})

}
