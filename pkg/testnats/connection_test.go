package testnats

import (
	"context"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestContainer_NewConnectionWithContext(t *testing.T) {
	t.Run("can connect", func(t *testing.T) {
		c, _ := RunTest(t)
		conn, err := c.NewConnectionWithContext(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		assert.Equal(t, nats.CONNECTED, conn.Status())
		conn.Close()
	})
	t.Run("respects ctx cancel", func(t *testing.T) {
		c, _ := RunTest(t)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		conn, err := c.NewConnectionWithContext(ctx)
		assert.Error(t, err)
		assert.Nil(t, conn)
	})
	t.Run("forwards connection error", func(t *testing.T) {
		c, _ := RunTest(t)
		c.ConnectionString = ""
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		conn, err := c.NewConnectionWithContext(ctx)
		assert.Error(t, err)
		assert.Nil(t, conn)
	})
}
