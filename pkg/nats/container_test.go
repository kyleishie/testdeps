package nats

import (
	"context"
	"errors"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go"
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
			conn, err := nats.Connect(con.ConnectionString)
			assert.NoError(t, err)
			assert.Equal(t, nats.CONNECTED, conn.Status())
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
			conn, err := nats.Connect(con.ConnectionString)
			assert.NoError(t, err)
			assert.Equal(t, nats.CONNECTED, conn.Status())
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
		ctx, cancel := context.WithTimeout(context.Background(), testDuration)
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
			conn, err := nats.Connect(con.ConnectionString)
			assert.NoError(t, err)
			assert.Equal(t, nats.CONNECTED, conn.Status())
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
			conn, err := nats.Connect(con.ConnectionString)
			assert.NoError(t, err)
			assert.Equal(t, nats.CONNECTED, conn.Status())
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
		ctx, cancel := context.WithTimeout(context.Background(), testDuration)
		cancel()
		con, newErr := RunTestWithContext(t, ctx)
		assert.Error(t, newErr)
		assert.Nil(t, con)
	})
}
