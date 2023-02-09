package testpostgres

import (
	"context"
	"fmt"
	"github.com/kyleishie/testdeps/pkg/common"
	"github.com/kyleishie/testdeps/pkg/options"
	"github.com/kyleishie/testdeps/pkg/testsql"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

const (
	image           = "postgres"
	mappedPort      = "5432/tcp"
	readyLog        = "database system is ready to accept connections"
	driver          = "postgres"
	defaultUser     = "postgres"
	defaultPassword = "postgres"
)

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
		// Note: This is printed twice during container init and the server is not ready until the second log.
		WaitingFor: wait.ForLog(readyLog).WithOccurrence(2),
		AutoRemove: true,
		Cmd:        []string{"postgres", "-c", "fsync=off"},
	}

	/// Apply opts
	for _, opt := range opts {
		err = opt(&cReq)
		if err != nil {
			return
		}
	}

	if pass, ok := cReq.Env[env_POSTGRES_PASSWORD]; !ok || pass != "" {
		/// The caller did not specify a password so let's prevent the container error.
		WithPassword(defaultPassword)(&cReq)
	}

	return
}

// Run creates and starts a docker Container with the `postgres` image.
// Defaults to `postgres:latest` if no option sets image tag.
// A default context is used with a timeout of two minutes. To customize use RunWithContext.
func Run(opts ...options.Option) (*testsql.Container, error) {
	ctx, cancel := context.WithTimeout(context.Background(), common.DefaultConnTimeout)
	defer cancel()
	return RunWithContext(ctx, opts...)
}

// MustRun is the same as Run except it panics on error.
func MustRun(opts ...options.Option) *testsql.Container {
	con, err := Run(opts...)
	if err != nil {
		panic(fmt.Sprintf("error creating container: %s", err.Error()))
	}
	return con
}

// RunWithContext creates and starts a docker Container with the `postgres` image.
// Defaults to `postgres:latest` if no option sets image tag.
// A context can be provided to configure things such as timeout.
func RunWithContext(ctx context.Context, opts ...options.Option) (con *testsql.Container, err error) {
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

	user, userExists := cReq.Env[env_POSTGRES_USER]
	if !userExists {
		user = defaultUser
	}

	password, passwordExists := cReq.Env[env_POSTGRES_PASSWORD]
	if !passwordExists {
		password = defaultPassword
	}

	con = testsql.New(c, driver, fmt.Sprintf("postgres://%s:%s@%s:%s?sslmode=disable", user, password, host, port.Port()))

	return
}

// RunForTest creates and starts a docker Container with the `postgres` image.
// Defaults to `postgres:latest` if no option sets image tag.
// The Container is automatically terminated after the test has finished.
// A default context is used with a timeout of two minutes. To customize use RunTestWithContext.
func RunForTest(t *testing.T, opts ...options.Option) *testsql.Container {
	ctx, cancel := context.WithTimeout(context.Background(), common.DefaultConnTimeout)
	defer cancel()
	return RunForTestWithContext(t, ctx, opts...)
}

// RunForTestWithContext creates and starts a docker Container with the `postgres` image.
// Defaults to `latest` if no option sets image tag.
// A context can be provided to configure things such as timeout.
// The Container is automatically terminated after the test has finished.
func RunForTestWithContext(t *testing.T, ctx context.Context, opts ...options.Option) *testsql.Container {
	c, err := RunWithContext(ctx, opts...)
	if err != nil {
		t.Fatalf("error starting container: %s", err.Error())
	}

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), common.DefaultConnTimeout)
		defer cancel()
		if err := c.Terminate(ctx); err != nil {
			t.Error(err)
		}
	})

	return c
}
