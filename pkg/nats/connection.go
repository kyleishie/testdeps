package nats

import (
	"context"
	"errors"
	"testing"

	"github.com/nats-io/nats.go"
)

// NewConnection creates a NATS connection to the underlying docker container.
// A default context is used with a timeout of two minutes. To customize use NewConnectionWithContext.
func (c *container) NewConnection(options ...nats.Option) (*nats.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), testDuration)
	defer cancel()
	return c.NewConnectionWithContext(ctx, options...)
}

// NewConnectionWithContext creates a NATS connection to the underlying docker container.
func (c *container) NewConnectionWithContext(ctx context.Context, options ...nats.Option) (*nats.Conn, error) {
	connChan := make(chan *nats.Conn, 1)
	errChan := make(chan error, 1)

	go func(connectionStr string, connChan chan *nats.Conn, errChan chan error, options ...nats.Option) {
		conn, err := nats.Connect(connectionStr, options...)
		if err != nil {
			errChan <- err
			return
		}
		connChan <- conn
	}(c.ConnectionString, connChan, errChan, options...)

	select {
	case conn := <-connChan:
		return conn, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, errors.New("ctx timeout")
	}
}

// NewTestConnection creates a NATS connection to the underlying docker container.
// The connection is automatically closed after t is finished.
//
// A default context is used with a timeout of two minutes. To customize use NewTestConnectionWithContext.
func (c *container) NewTestConnection(t *testing.T, options ...nats.Option) (*nats.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), testDuration)
	defer cancel()
	return c.NewTestConnectionWithContext(t, ctx, options...)
}

// NewTestConnectionWithContext creates a NATS connection to the underlying docker container.
// The connection is automatically closed after t is finished.
func (c *container) NewTestConnectionWithContext(t *testing.T, ctx context.Context, options ...nats.Option) (*nats.Conn, error) {
	conn, err := c.NewConnectionWithContext(ctx, options...)
	if err != nil {
		return nil, err
	}

	t.Cleanup(func() {
		conn.Close()
	})

	return conn, err
}
