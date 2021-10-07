package nats

import (
	"context"
	"testing"

	"github.com/nats-io/nats.go"
)

// NewJetStream creates a JetStreamContext and its NATS connection to the underlying docker container.
// A default context is used with a timeout of two minutes. To customize use NewJetStreamWithContext.
//
// If you need to use nats.Option's to customize the underlying connection use NewConnection and create the
// JetStreamContext yourself.
//
// To use JetStream the container must be created with the WithJetStream option.
func (c *container) NewJetStream(opts ...nats.JSOpt) (nats.JetStreamContext, error) {
	ctx, cancel := context.WithTimeout(context.Background(), testDuration)
	defer cancel()
	return c.NewJetStreamWithContext(ctx, opts...)
}

// NewJetStreamWithContext creates a JetStreamContext and its NATS connection to the underlying docker container.
//
// If you need to use nats.Option's to customize the underlying connection use NewConnectionWithContext and create the
// JetStreamContext yourself.
//
// To use JetStream the container must be created with the WithJetStream option.
func (c *container) NewJetStreamWithContext(ctx context.Context, opts ...nats.JSOpt) (nats.JetStreamContext, error) {
	conn, err := c.NewConnectionWithContext(ctx)
	if err != nil {
		return nil, err
	}
	return conn.JetStream(opts...)
}

// NewTestJetStream creates a JetStreamContext and its NATS connection to the underlying docker container.
// The connection is automatically closed after t is finished.
// A default context is used with a timeout of two minutes. To customize use NewTestJetStreamWithContext.
//
// If you need to use nats.Option's to customize the underlying connection use NewTestConnection and create the
// JetStreamContext yourself.
//
// To use JetStream the container must be created with the WithJetStream option.
func (c *container) NewTestJetStream(t *testing.T, opts ...nats.JSOpt) (nats.JetStreamContext, error) {
	ctx, cancel := context.WithTimeout(context.Background(), testDuration)
	defer cancel()
	return c.NewTestJetStreamWithContext(t, ctx, opts...)
}

// NewTestJetStreamWithContext creates a JetStreamContext and its NATS connection to the underlying docker container.
// The connection is automatically closed after t is finished.
//
// If you need to use nats.Option's to customize the underlying connection use NewTestConnectionWithContext and create the
// JetStreamContext yourself.
//
// To use JetStream the container must be created with the WithJetStream option.
func (c *container) NewTestJetStreamWithContext(t *testing.T, ctx context.Context, opts ...nats.JSOpt) (nats.JetStreamContext, error) {
	conn, err := c.NewTestConnectionWithContext(t, ctx)
	if err != nil {
		return nil, err
	}

	t.Cleanup(func() {
		conn.Close()
	})

	return conn.JetStream(opts...)
}
