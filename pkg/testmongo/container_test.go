package testmongo

import (
	"context"
	"errors"
	"github.com/kyleishie/testdeps/pkg/common"
	"testing"

	opts "github.com/kyleishie/testdeps/pkg/options"
	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestRun(t *testing.T) {
	t.Parallel()

	t.Run("should start Container", func(t *testing.T) {
		t.Parallel()
		con, newErr := Run()
		assert.NoError(t, newErr)
		assert.NotNil(t, con)

		t.Run("connection string", func(t *testing.T) {
			assert.NotEmpty(t, con.ConnectionString)
		})

		t.Run("can connect to Container", func(t *testing.T) {
			/// Test that Container is alive by connecting to it
			ctx := context.Background()
			clientOptions := options.Client().ApplyURI(con.ConnectionString)
			client, connErr := mongo.Connect(ctx, clientOptions)
			assert.NoError(t, connErr)
			assert.NotNil(t, client)
			assert.NoError(t, client.Ping(ctx, nil))
			assert.NoError(t, client.Disconnect(ctx))
		})

		t.Run("can terminate", func(t *testing.T) {
			ctx := context.Background()
			termErr := con.Terminate(ctx)
			assert.NoError(t, termErr)
		})
	})
	t.Run("forwards opts errors", func(t *testing.T) {
		t.Parallel()
		testErr := errors.New("test error")
		con, newErr := Run(func(request *tc.ContainerRequest) error {
			return testErr
		})
		assert.ErrorIs(t, newErr, testErr)
		assert.Nil(t, con)
	})
	t.Run("forwards Container creation error", func(t *testing.T) {
		t.Parallel()
		con, newErr := Run(func(request *tc.ContainerRequest) error {
			request.Image = ""
			return nil
		})
		assert.Error(t, newErr)
		assert.Nil(t, con)
	})
	t.Run("forwards running Container error", func(t *testing.T) {
		t.Parallel()
		con, newErr := Run(func(request *tc.ContainerRequest) error {
			request.ExposedPorts = nil
			return nil
		})
		assert.Error(t, newErr)
		assert.Nil(t, con)
	})
}

func TestRunWithContext(t *testing.T) {
	t.Parallel()

	t.Run("should start Container", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		con, newErr := RunWithContext(ctx)
		assert.NoError(t, newErr)
		assert.NotNil(t, con)

		t.Run("connection string", func(t *testing.T) {
			assert.NotEmpty(t, con.ConnectionString)
		})

		t.Run("can connect to Container", func(t *testing.T) {
			/// Test that Container is alive by connecting to it
			ctx := context.Background()
			clientOptions := options.Client().ApplyURI(con.ConnectionString)
			client, connErr := mongo.Connect(ctx, clientOptions)
			assert.NoError(t, connErr)
			assert.NotNil(t, client)
			assert.NoError(t, client.Ping(ctx, nil))
			assert.NoError(t, client.Disconnect(ctx))
		})

		t.Run("can terminate", func(t *testing.T) {
			ctx := context.Background()
			termErr := con.Terminate(ctx)
			assert.NoError(t, termErr)
		})
	})
	t.Run("forwards opts errors", func(t *testing.T) {
		t.Parallel()
		testErr := errors.New("test error")
		ctx := context.Background()
		con, newErr := RunWithContext(ctx, func(request *tc.ContainerRequest) error {
			return testErr
		})
		assert.ErrorIs(t, newErr, testErr)
		assert.Nil(t, con)
	})
	t.Run("forwards Container creation error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		con, newErr := RunWithContext(ctx, func(request *tc.ContainerRequest) error {
			request.Image = ""
			return nil
		})
		assert.Error(t, newErr)
		assert.Nil(t, con)
	})
	t.Run("forwards running Container error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		con, newErr := RunWithContext(ctx, func(request *tc.ContainerRequest) error {
			request.ExposedPorts = nil
			return nil
		})
		assert.Error(t, newErr)
		assert.Nil(t, con)
	})
	t.Run("respects ctx", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithTimeout(context.Background(), common.DefaultConnTimeout)
		cancel()
		con, newErr := RunWithContext(ctx)
		assert.Error(t, newErr)
		assert.Nil(t, con)
	})
}

func TestRunTest(t *testing.T) {
	t.Parallel()

	t.Run("terminates Container after test", func(t *testing.T) {
		t.Parallel()
		var con *Container
		/// Using sub test as testing scoep for termination
		t.Run("sub test scope", func(t *testing.T) {
			con, _ = RunTest(t)
		})

		ip, err := con.Container.ContainerIP(context.Background())
		assert.Error(t, err)
		assert.Empty(t, ip)
	})

	t.Run("should start Container", func(t *testing.T) {
		t.Parallel()
		con, newErr := RunTest(t)
		assert.NoError(t, newErr)
		assert.NotNil(t, con)

		t.Run("connection string", func(t *testing.T) {
			assert.NotEmpty(t, con.ConnectionString)
		})

		t.Run("can connect to Container", func(t *testing.T) {
			/// Test that Container is alive by connecting to it
			ctx := context.Background()
			clientOptions := options.Client().ApplyURI(con.ConnectionString)
			client, connErr := mongo.Connect(ctx, clientOptions)
			assert.NoError(t, connErr)
			assert.NotNil(t, client)
			assert.NoError(t, client.Ping(ctx, nil))
			assert.NoError(t, client.Disconnect(ctx))
		})
	})
	t.Run("forwards opts errors", func(t *testing.T) {
		t.Parallel()
		testErr := errors.New("test error")
		con, newErr := RunTest(t, func(request *tc.ContainerRequest) error {
			return testErr
		})
		assert.ErrorIs(t, newErr, testErr)
		assert.Nil(t, con)
	})
	t.Run("forwards Container creation error", func(t *testing.T) {
		t.Parallel()
		con, newErr := RunTest(t, func(request *tc.ContainerRequest) error {
			request.Image = ""
			return nil
		})
		assert.Error(t, newErr)
		assert.Nil(t, con)
	})
	t.Run("forwards running Container error", func(t *testing.T) {
		t.Parallel()
		con, newErr := RunTest(t, func(request *tc.ContainerRequest) error {
			request.ExposedPorts = nil
			return nil
		})
		assert.Error(t, newErr)
		assert.Nil(t, con)
	})
}

func TestRunTestWithContext(t *testing.T) {
	t.Parallel()

	t.Run("terminates Container after test", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		var con *Container
		/// Using sub test as testing scoep for termination
		t.Run("sub test scope", func(t *testing.T) {
			con, _ = RunTestWithContext(t, ctx)
		})

		ip, err := con.Container.ContainerIP(context.Background())
		assert.Error(t, err)
		assert.Empty(t, ip)
	})

	t.Run("should start Container", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		con, newErr := RunTestWithContext(t, ctx)
		assert.NoError(t, newErr)
		assert.NotNil(t, con)

		t.Run("connection string", func(t *testing.T) {
			assert.NotEmpty(t, con.ConnectionString)
		})

		t.Run("can connect to Container", func(t *testing.T) {
			/// Test that Container is alive by connecting to it
			ctx := context.Background()
			clientOptions := options.Client().ApplyURI(con.ConnectionString)
			client, connErr := mongo.Connect(ctx, clientOptions)
			assert.NoError(t, connErr)
			assert.NotNil(t, client)
			assert.NoError(t, client.Ping(ctx, nil))
			assert.NoError(t, client.Disconnect(ctx))
		})
	})

	t.Run("forwards opts errors", func(t *testing.T) {
		t.Parallel()
		testErr := errors.New("test error")
		ctx := context.Background()
		con, newErr := RunTestWithContext(t, ctx, func(request *tc.ContainerRequest) error {
			return testErr
		})
		assert.ErrorIs(t, newErr, testErr)
		assert.Nil(t, con)
	})
	t.Run("forwards Container creation error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		con, newErr := RunTestWithContext(t, ctx, func(request *tc.ContainerRequest) error {
			request.Image = ""
			return nil
		})
		assert.Error(t, newErr)
		assert.Nil(t, con)
	})
	t.Run("forwards running Container error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		con, newErr := RunTestWithContext(t, ctx, func(request *tc.ContainerRequest) error {
			request.ExposedPorts = nil
			return nil
		})
		assert.Error(t, newErr)
		assert.Nil(t, con)
	})
	t.Run("respects ctx", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithTimeout(context.Background(), common.DefaultConnTimeout)
		cancel()
		con, newErr := RunTestWithContext(t, ctx)
		assert.Error(t, newErr)
		assert.Nil(t, con)
	})
}

func TestMakeContainerRequest(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		cReq, err := makeContainerRequest([]opts.Option{})
		assert.NoError(t, err)
		assert.NotEmpty(t, cReq)
		assert.Equal(t, image, cReq.Image)
		assert.Equal(t, []string{mappedPort}, cReq.ExposedPorts)
		assert.Equal(t, wait.ForLog(readyLog), cReq.WaitingFor)
		assert.True(t, cReq.AutoRemove)
	})
	t.Run("applies options", func(t *testing.T) {
		image := "test"
		cReq, err := makeContainerRequest([]opts.Option{
			func(request *tc.ContainerRequest) error {
				request.Image = image
				return nil
			},
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, cReq)
		assert.Equal(t, image, cReq.Image)
	})
	t.Run("error in opt", func(t *testing.T) {
		testErr := errors.New("test error")
		cReq, err := makeContainerRequest([]opts.Option{
			func(request *tc.ContainerRequest) error {
				return testErr
			},
		})
		assert.ErrorIs(t, err, testErr)
		assert.Empty(t, cReq)
	})
}

func TestMakeRootUserPrefix(t *testing.T) {
	testUser := "user"
	testPass := "pass"
	expectation := testUser + ":" + testPass + "@"
	cReq := tc.ContainerRequest{
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": testUser,
			"MONGO_INITDB_ROOT_PASSWORD": testPass,
		},
	}

	authStr := makeRootUserPrefix(cReq)
	assert.Equal(t, expectation, authStr)
}
