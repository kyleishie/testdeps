package mongo

import (
	"context"
	"fmt"
	"time"

	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testdeps/internal/options"
)

const (
	image        = "mongo"
	mappedPort   = "27017/tcp"
	readyLog     = "Waiting for connections"
	proto        = "mongodb"
	testDuration = time.Minute * 2
)

type container struct {
	tc.Container
	ConnectionString string
}

func Run(opts ...options.Option) (*container, error) {
	return RunWithContext(context.Background(), opts...)
}

func RunWithContext(ctx context.Context, opts ...options.Option) (con *container, err error) {
	cReq, err := makeContainerRequest(opts)
	if err != nil {
		return
	}

	c, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: cReq,
		Started:          true,
	})
	if err != nil {
		return
	}

	host, err := c.Host(ctx)
	port, err := c.MappedPort(ctx, mappedPort)
	if err != nil {
		return
	}

	con = &container{
		Container:        c,
		ConnectionString: fmt.Sprintf("%s://%s%s:%d", proto, makeRootUserPrefix(cReq), host, port.Int()),
	}

	return
}

func makeContainerRequest(opts []options.Option) (cReq tc.ContainerRequest, err error) {
	defer func() {
		if err != nil {
			cReq = tc.ContainerRequest{}
		}
	}()

	cReq = tc.ContainerRequest{
		Image: image,
		ExposedPorts: []string{
			mappedPort,
		},
		WaitingFor: wait.ForLog(readyLog),
		AutoRemove: true,
	}

	/// Apply opts
	for _, opt := range opts {
		err = opt(&cReq)
		if err != nil {
			return
		}
	}

	return
}

func makeRootUserPrefix(cReq tc.ContainerRequest) (auth string) {
	username, uExists := cReq.Env[env_MONGO_INITDB_ROOT_USERNAME]
	password, pExists := cReq.Env[env_MONGO_INITDB_ROOT_PASSWORD]
	if uExists && pExists {
		auth = fmt.Sprintf("%s:%s@", username, password)
	}
	return
}
