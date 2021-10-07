package mongo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestContainer_NewDatabase(t *testing.T) {
	t.Parallel()
	con, _ := Run()
	t.Run("database named correctly", func(subT *testing.T) {
		testName := "testDb"
		db, err := con.NewDatabase(testName)
		assert.NoError(t, err)
		assert.Equal(subT, testName, db.Name())
	})
	t.Run("can use database", func(t *testing.T) {
		db, _ := con.NewDatabase("testName")
		collection := db.Collection("testCollection")
		ctx := context.Background()
		_, insertErr := collection.InsertOne(ctx, bson.D{{"test", true}})
		assert.NoError(t, insertErr)
	})
	t.Run("forwards errors", func(t *testing.T) {
		t.Parallel()
		con, _ := Run()
		con.ConnectionString = ""
		db, err := con.NewDatabase("testName")
		assert.Nil(t, db)
		assert.Error(t, err)
	})
}

func TestContainer_NewDatabaseWithContext(t *testing.T) {
	t.Parallel()
	con, _ := Run()
	t.Run("database named correctly", func(subT *testing.T) {
		testName := "testDb"
		ctx := context.Background()
		db, err := con.NewDatabaseWithContext(ctx, testName)
		assert.NoError(t, err)
		assert.Equal(subT, testName, db.Name())
	})
	t.Run("can use database", func(t *testing.T) {
		ctx := context.Background()
		db, _ := con.NewDatabaseWithContext(ctx, "testName")
		collection := db.Collection("testCollection")
		_, insertErr := collection.InsertOne(ctx, bson.D{{"test", true}})
		assert.NoError(t, insertErr)
	})
	t.Run("forwards errors", func(t *testing.T) {
		t.Parallel()
		con, _ := Run()
		con.ConnectionString = ""
		ctx := context.Background()
		db, err := con.NewDatabaseWithContext(ctx, "testName")
		assert.Nil(t, db)
		assert.Error(t, err)
	})
}

func TestContainer_NewTestDatabase(t *testing.T) {
	t.Parallel()
	con, _ := Run()

	t.Run("can use database", func(t *testing.T) {
		ctx := context.Background()
		db, _ := con.NewTestDatabase(t)
		collection := db.Collection("testCollection")
		_, insertErr := collection.InsertOne(ctx, bson.D{{"test", true}})
		assert.NoError(t, insertErr)
	})

	t.Run("forwards errors", func(t *testing.T) {
		t.Parallel()
		con, _ := Run()
		con.ConnectionString = ""
		db, err := con.NewTestDatabase(t)
		assert.Nil(t, db)
		assert.Error(t, err)
	})
}

func TestContainer_NewTestDatabaseWithContext(t *testing.T) {
	t.Parallel()
	con, _ := Run()

	t.Run("can use database", func(t *testing.T) {
		ctx := context.Background()
		db, _ := con.NewTestDatabaseWithContext(t, ctx)
		collection := db.Collection("testCollection")
		_, insertErr := collection.InsertOne(ctx, bson.D{{"test", true}})
		assert.NoError(t, insertErr)
	})
	t.Run("forwards errors", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		con, _ := Run()
		con.ConnectionString = ""
		db, err := con.NewTestDatabaseWithContext(t, ctx)
		assert.Nil(t, db)
		assert.Error(t, err)
	})
}
