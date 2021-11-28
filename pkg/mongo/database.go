package mongo

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewDatabase creates a new mongo client then a new database with then given name and options.
func (c *Container) NewDatabase(name string, opts ...*options.DatabaseOptions) (*mongo.Database, error) {
	return c.NewDatabaseWithContext(context.Background(), name, opts...)
}

// NewDatabaseWithContext creates a new mongo.Database with then given name and options.
// NewDatabaseWithContext exists to allow you to customize the connection process, e.g., apply timeout.
func (c *Container) NewDatabaseWithContext(ctx context.Context, name string, opts ...*options.DatabaseOptions) (*mongo.Database, error) {
	client, err := c.NewClientWithContext(ctx)
	if err != nil {
		return nil, err
	}
	return client.Database(name, opts...), nil
}

// NewTestDatabase creates a new Database with a random name within the Container.
// The database is automatically dropped after to test it finished.
// Note: A default context is used with a timeout of two minutes.
func (c *Container) NewTestDatabase(t *testing.T, opts ...*options.DatabaseOptions) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), testDuration)
	defer cancel()
	return c.NewTestDatabaseWithContext(t, ctx, opts...)
}

// NewTestDatabaseWithContext creates a new mongo.Database with a random name within the Container.
// The database will be named randomly.
// The database is automatically dropped and the underlying mongo.Client will be disconnected after to test it finished.
func (c *Container) NewTestDatabaseWithContext(t *testing.T, ctx context.Context, opts ...*options.DatabaseOptions) (*mongo.Database, error) {
	dbId := uuid.New().String()
	dbId = strings.ReplaceAll(dbId, "-", "")

	db, err := c.NewDatabaseWithContext(ctx, dbId, opts...)
	if err != nil {
		return nil, err
	}

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), testDuration)
		defer cancel()
		if err := db.Drop(ctx); err != nil {
			t.Error(err)
		}
		if err := db.Client().Disconnect(ctx); err != nil {
			t.Error(err)
		}
	})

	return db, err
}
