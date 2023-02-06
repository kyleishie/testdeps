package testsql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/kyleishie/testdeps/pkg/common"
	"strings"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (c *Container) NewDatabase(name, migrationsFilepath string) (*sql.DB, error) {
	return c.NewDatabaseWithContext(context.Background(), name, migrationsFilepath)
}

// NewDatabaseWithContext creates a new mongo.Database with then given name and options.
// NewDatabaseWithContext exists to allow you to customize the connection process, e.g., apply timeout.
func (c *Container) NewDatabaseWithContext(ctx context.Context, name, migrationsFilepath string) (*sql.DB, error) {

	client, err := c.NewClientWithContext(ctx)
	if err != nil {
		return nil, err
	}

	_, err = client.Exec(fmt.Sprintf("create database %s", name))
	if err != nil {
		return nil, err
	}

	if err = client.Close(); err != nil {
		return nil, err
	}

	//TODO: Implement migration level here for schema backwards compatability tests
	if migrationsFilepath != "" {
		if !strings.HasPrefix(migrationsFilepath, "file://") {
			migrationsFilepath = fmt.Sprintf("file://%s", migrationsFilepath)
		}
		m, err := migrate.New(migrationsFilepath, c.ConnectionString)
		if err = m.Up(); err != nil {
			return nil, err
		}
	}

	return c.NewClientWithContext(ctx)
}

// NewTestDatabase creates a new Database with a random name within the Container.
// The database is automatically dropped after to test it finished.
// Note: A default context is used with a timeout of two minutes.
func (c *Container) NewTestDatabase(t testing.TB, migrationsFilepath string) *sql.DB {
	ctx, cancel := context.WithTimeout(context.Background(), common.DefaultConnTimeout)
	defer cancel()
	return c.NewTestDatabaseWithContext(t, ctx, migrationsFilepath)
}

// NewTestDatabaseWithContext creates a new mongo.Database with a random name within the Container.
// The database will be named randomly.
// The database is automatically dropped and the underlying mongo.Client will be disconnected after to test it finished.
// Any error that occurs will result in a t.Fatal
func (c *Container) NewTestDatabaseWithContext(t testing.TB, ctx context.Context, migrationsFilepath string) *sql.DB {
	name := common.GenerateId()

	db, err := c.NewDatabaseWithContext(ctx, name, migrationsFilepath)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if _, err := db.Exec(fmt.Sprintf("drop database %s", name)); err != nil {
			t.Error(err)
		}
		if err = db.Close(); err != nil {
			t.Error(err)
		}
	})

	return db
}
