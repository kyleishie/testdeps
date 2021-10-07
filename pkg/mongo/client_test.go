package mongo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainer_NewClient(t *testing.T) {
	t.Parallel()

	con, _ := Run()
	t.Run("should connect", func(t *testing.T) {
		ctx := context.Background()
		client, err := con.NewClient()
		assert.NoError(t, err)
		assert.NoError(t, client.Ping(ctx, nil))
	})
	t.Run("should error", func(t *testing.T) {
		con.ConnectionString = ""
		client, err := con.NewClient()
		assert.Error(t, err)
		assert.Nil(t, client)
	})
}

func TestContainer_NewClientWithContext(t *testing.T) {
	t.Parallel()
	con, _ := Run()
	t.Run("should connect", func(t *testing.T) {
		ctx := context.Background()
		client, err := con.NewClientWithContext(ctx)
		assert.NoError(t, err)
		assert.NoError(t, client.Ping(ctx, nil))
	})
	t.Run("should error", func(t *testing.T) {
		ctx := context.Background()
		con, _ := Run()
		con.ConnectionString = ""
		client, err := con.NewClientWithContext(ctx)
		assert.Error(t, err)
		assert.Nil(t, client)
	})
	t.Run("good ctx", func(t *testing.T) {
		ctx := context.Background()
		client, err := con.NewClientWithContext(ctx)
		assert.NoError(t, err)
		assert.NoError(t, client.Ping(ctx, nil))
	})
	t.Run("cancelled ctx", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), testDuration)
		cancel()
		client, err := con.NewClientWithContext(ctx)
		assert.Error(t, err)
		assert.Nil(t, client)
	})
}

func TestContainer_NewTestClient(t *testing.T) {
	t.Parallel()

	con, _ := Run()
	t.Run("should connect", func(t *testing.T) {
		ctx := context.Background()
		client, err := con.NewTestClient(t)
		assert.NoError(t, err)
		assert.NoError(t, client.Ping(ctx, nil))
	})
	t.Run("should error", func(t *testing.T) {
		con.ConnectionString = ""
		client, err := con.NewClient()
		assert.Error(t, err)
		assert.Nil(t, client)
	})
}

func TestContainer_NewTestClientWithContext(t *testing.T) {
	t.Parallel()
	con, _ := Run()
	t.Run("should connect", func(t *testing.T) {
		ctx := context.Background()
		client, err := con.NewTestClientWithContext(t, ctx)
		assert.NoError(t, err)
		assert.NoError(t, client.Ping(ctx, nil))
	})
	t.Run("should error", func(t *testing.T) {
		ctx := context.Background()
		con, _ := Run()
		con.ConnectionString = ""
		client, err := con.NewTestClientWithContext(t, ctx)
		assert.Error(t, err)
		assert.Nil(t, client)
	})
	t.Run("good ctx", func(t *testing.T) {
		ctx := context.Background()
		client, err := con.NewTestClientWithContext(t, ctx)
		assert.NoError(t, err)
		assert.NoError(t, client.Ping(ctx, nil))
	})
	t.Run("cancelled ctx", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), testDuration)
		cancel()
		client, err := con.NewTestClientWithContext(t, ctx)
		assert.Error(t, err)
		assert.Nil(t, client)
	})

}
