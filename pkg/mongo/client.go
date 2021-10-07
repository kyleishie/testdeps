package mongo

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewClient creates a new mongo client using the connection string of the Container.
// The connection is tested once before returning the new client.
func (c *container) NewClient() (*mongo.Client, error) {
	return c.NewClientWithContext(context.Background())
}

// NewClientWithContext creates a new mongo client using the connection string of the Container.
// The connection is tested once before returning the new client.
// NewClientWithContext exists to allow you to customize the connection process, e.g., apply timeout.
func (c *container) NewClientWithContext(ctx context.Context) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(c.ConnectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// NewTestClient creates a mongo.Client for testing purposes.
// The mongo.Client will be disconnected automatically after the test finishes.
// Note: A default context is used with a timeout of two minutes.
func (c *container) NewTestClient(t *testing.T) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), testDuration)
	defer cancel()
	return c.NewTestClientWithContext(t, ctx)
}

// NewTestClientWithContext creates a mongo.Client for testing purposes.
// The mongo.Client will be disconnected automatically after the test finishes.
func (c *container) NewTestClientWithContext(t *testing.T, ctx context.Context) (*mongo.Client, error) {
	client, err := c.NewClientWithContext(ctx)
	if err != nil {
		return nil, err
	}

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), testDuration)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			t.Error(err)
		}
	})

	return client, nil
}
