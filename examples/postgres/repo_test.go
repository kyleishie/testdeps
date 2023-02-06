package postgres

import (
	"github.com/kyleishie/testdeps/pkg/testsql/testpostgres"
	"github.com/stretchr/testify/assert"
	"testing"
)

var container = testpostgres.MustRun()

func TestCreateUser(t *testing.T) {

	user := User{
		Name:  "example",
		Email: "exmaple@example.com",
	}

	t.Run("can create user", func(t *testing.T) {
		db := container.NewTestDatabase(t, "migrations")
		err := CreateUser(db, user)
		assert.NoError(t, err)

		rows, err := db.Query("select * from users where email = $1", user.Email)
		assert.NoError(t, err)
		defer rows.Close()

		/// 1 rows should be available
		var u User
		assert.True(t, rows.Next())
		assert.NoError(t, rows.Scan(&u.Name, &u.Email))

		/// No other rows should be available
		assert.False(t, rows.Next())
	})

}
