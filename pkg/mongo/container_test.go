package mongo

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	opts "testdeps/pkg/options"
)

func TestRun(t *testing.T) {
	t.Parallel()

	var con *container

	t.Run("should start container", func(t *testing.T) {
		var newErr error
		con, newErr = Run()
		assert.NoError(t, newErr)
		assert.NotNil(t, con)
	})

	testContainer(t, con)

	t.Run("should error", func(t *testing.T) {
		t.Run("forwards opts errors", func(t *testing.T) {
			testErr := errors.New("test error")
			var newErr error
			con, newErr = Run(func(request *tc.ContainerRequest) error {
				return testErr
			})
			assert.ErrorIs(t, newErr, testErr)
			assert.Nil(t, con)
		})
		t.Run("forwards container errors", func(t *testing.T) {
			t.Run("container creation error", func(t *testing.T) {
				var newErr error
				con, newErr = Run(func(request *tc.ContainerRequest) error {
					request.Image = ""
					return nil
				})
				assert.Error(t, newErr)
				assert.Nil(t, con)
			})
			t.Run("running container error", func(t *testing.T) {
				var newErr error
				con, newErr = Run(func(request *tc.ContainerRequest) error {
					request.ExposedPorts = nil
					return nil
				})
				assert.Error(t, newErr)
				assert.Nil(t, con)
			})
		})
	})

}

func TestRunWithContext(t *testing.T) {
	t.Parallel()

	var con *container

	t.Run("should start container", func(t *testing.T) {
		ctx := context.Background()
		var newErr error
		con, newErr = RunWithContext(ctx)
		assert.NoError(t, newErr)
		assert.NotNil(t, con)
	})

	testContainer(t, con)

	t.Run("should error", func(t *testing.T) {
		t.Run("forwards opts errors", func(t *testing.T) {
			testErr := errors.New("test error")
			ctx := context.Background()
			var newErr error
			con, newErr = RunWithContext(ctx, func(request *tc.ContainerRequest) error {
				return testErr
			})
			assert.ErrorIs(t, newErr, testErr)
			assert.Nil(t, con)
		})
		t.Run("forwards container errors", func(t *testing.T) {
			t.Run("container creation error", func(t *testing.T) {
				ctx := context.Background()
				var newErr error
				con, newErr = RunWithContext(ctx, func(request *tc.ContainerRequest) error {
					request.Image = ""
					return nil
				})
				assert.Error(t, newErr)
				assert.Nil(t, con)
			})
			t.Run("running container error", func(t *testing.T) {
				ctx := context.Background()
				var newErr error
				con, newErr = RunWithContext(ctx, func(request *tc.ContainerRequest) error {
					request.ExposedPorts = nil
					return nil
				})
				assert.Error(t, newErr)
				assert.Nil(t, con)
			})
		})

	})
}

func testContainer(t *testing.T, con *container) {
	t.Run("connection string", func(t *testing.T) {
		assert.NotEmpty(t, con.ConnectionString)
	})

	t.Run("can connect to container", func(t *testing.T) {
		/// Test that container is alive by connecting to it
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
}

func TestMakeContainerRequest(t *testing.T) {

	t.Run("should succeed", func(t *testing.T) {
		t.Run("default", func(t *testing.T) {
			cReq, err := makeContainerRequest([]opts.Option{})
			assert.NoError(t, err)
			assert.NotEmpty(t, cReq)
			assert.Equal(t, image, cReq.Image)
			assert.Equal(t, []string{mappedPort}, cReq.ExposedPorts)
			assert.Equal(t, wait.ForLog(readyLog), cReq.WaitingFor)
			assert.True(t, cReq.AutoRemove)
		})
		t.Run("applies opts", func(t *testing.T) {
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
	})

	t.Run("should error", func(t *testing.T) {
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
