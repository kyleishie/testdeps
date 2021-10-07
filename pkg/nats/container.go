package nats

import (
	"context"
	"testing"
	"time"

	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testdeps/pkg/options"
)

const (
	image               = "nats"
	mappedPort          = "4222/tcp"
	readyLog            = "Server is ready"
	proto               = "nats"
	cmdJetStreamEnabled = "-js"
	testDuration        = time.Minute * 2
)

type container struct {
	tc.Container
	req              tc.ContainerRequest
	ConnectionString string
}

// Run creates and starts a docker container with the `nats/nats` image.
// Defaults to `nats/nats:latest` if no option sets image tag.
// A default context is used with a timeout of two minutes. To customize use RunWithContext.
func Run(options ...options.Option) (*container, error) {
	ctx, cancel := context.WithTimeout(context.Background(), testDuration)
	defer cancel()
	return RunWithContext(ctx, options...)
}

// RunWithContext creates and starts a docker container with the `nats/nats` image.
// Defaults to `nats/nats:latest` if no option sets image tag.
// A context can be provided to configure things such as timeout.
func RunWithContext(ctx context.Context, options ...options.Option) (con *container, err error) {
	cReq := tc.ContainerRequest{
		Image: image,
		ExposedPorts: []string{
			mappedPort,
		},
		WaitingFor: wait.ForLog(readyLog),
		AutoRemove: true,
	}

	/// Apply options
	for _, opt := range options {
		err = opt(&cReq)
		if err != nil {
			return
		}
	}

	c, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: cReq,
		Started:          true,
	})
	if err != nil {
		return
	}

	connStr, err := c.PortEndpoint(ctx, mappedPort, proto)
	if err != nil {
		return
	}

	con = &container{
		Container:        c,
		req:              cReq,
		ConnectionString: connStr,
	}

	return
}

// RunTest creates and starts a docker container with the `nats/nats` image.
// Defaults to `nats/nats:latest` if no option sets image tag.
// The container is automatically terminated after the test is finished.
// A default context is used with a timeout of two minutes. To customize use RunTestWithContext.
func RunTest(t *testing.T, opts ...options.Option) (con *container, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), testDuration)
	defer cancel()
	return RunTestWithContext(t, ctx, opts...)
}

// RunTestWithContext creates and starts a docker container with the `nats/nats` image.
// Defaults to `nats/nats:latest` if no option sets image tag.
// A context can be provided to configure things such as timeout.
// The container is automatically terminated after the test is finished.
func RunTestWithContext(t *testing.T, ctx context.Context, opts ...options.Option) (con *container, err error) {
	c, err := RunWithContext(ctx, opts...)
	if err != nil {
		return nil, err
	}

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), testDuration)
		defer cancel()
		if err := c.Terminate(ctx); err != nil {
			t.Error(err)
		}
	})

	return c, nil
}
