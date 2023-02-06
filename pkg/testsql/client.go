package testsql

import (
	"context"
	"database/sql"
	"github.com/kyleishie/testdeps/pkg/common"
	"testing"
)

// NewClient creates a new *sql.DB using the connection string and driver of the Container.
// The connection is tested once before returning the new client.
func (c *Container) NewClient() (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), common.DefaultConnTimeout)
	defer cancel()
	return c.NewClientWithContext(ctx)
}

// NewClientWithContext creates a new *sql.DB using the connection string and driver of the Container.
// The connection is tested once before returning the new client.
// NewClientWithContext exists to allow you to customize the connection process, e.g., apply timeout.
func (c *Container) NewClientWithContext(ctx context.Context) (*sql.DB, error) {
	client, err := sql.Open(c.driver, c.ConnectionString)
	if err != nil {
		return nil, err
	}

	if err := client.PingContext(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

// NewTestClient creates a *sql.DB for testing purposes.
// The *sql.DB will be disconnected automatically after the test finishes.
// Note: A default context is used with a timeout of two minutes.
func (c *Container) NewTestClient(t *testing.T) (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), common.DefaultConnTimeout)
	defer cancel()
	return c.NewTestClientWithContext(t, ctx)
}

// NewTestClientWithContext creates a *sql.DB for testing purposes.
// The *sql.DB will be closed automatically after the test finishes.
func (c *Container) NewTestClientWithContext(t *testing.T, ctx context.Context) (*sql.DB, error) {
	client, err := c.NewClientWithContext(ctx)
	if err != nil {
		return nil, err
	}

	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Error(err)
		}
	})

	return client, nil
}
