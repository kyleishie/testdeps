package nats

import (
	"context"

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
)

type container struct {
	tc.Container
	req              tc.ContainerRequest
	ConnectionString string
}

func Run(options ...options.Option) (*container, error) {
	return RunWithContext(context.Background(), options...)
}

func RunWithContext(ctx context.Context, options ...options.Option) (con *container, err error) {
	if ctx == nil {
		ctx = context.Background()
	}

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

	connStr, err := c.Endpoint(ctx, proto)
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

func (c *container) Start() error {
	ctx := context.Background()
	return c.Container.Start(ctx)
}

func (c *container) Terminate() error {
	ctx := context.Background()
	return c.Container.Terminate(ctx)
}
