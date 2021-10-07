package mongo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kyleishie/testdeps/pkg/options"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
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

// Run creates and starts a docker container with the `mongo` image.
// Defaults to `mongo:latest` if no option sets image tag.
// A default context is used with a timeout of two minutes. To customize use RunWithContext.
func Run(opts ...options.Option) (*container, error) {
	ctx, cancel := context.WithTimeout(context.Background(), testDuration)
	defer cancel()
	return RunWithContext(ctx, opts...)
}

// RunWithContext creates and starts a docker container with the `mongo` image.
// Defaults to `mongo:latest` if no option sets image tag.
// A context can be provided to configure things such as timeout.
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

// RunTest creates and starts a docker container with the `mongo` image.
// Defaults to `mongo:latest` if no option sets image tag.
// The container is automatically terminated after the test is finished.
// A default context is used with a timeout of two minutes. To customize use RunTestWithContext.
func RunTest(t *testing.T, opts ...options.Option) (con *container, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), testDuration)
	defer cancel()
	return RunTestWithContext(t, ctx, opts...)
}

// RunTestWithContext creates and starts a docker container with the `mongo` image.
// Defaults to `mongo:latest` if no option sets image tag.
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
