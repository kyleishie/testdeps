package nats

import (
	"github.com/testcontainers/testcontainers-go"
	"testdeps/internal/options"
)

// WithJetStream enables JetStream in the NATS container.
// This is the equivalent of running nats -js.
func WithJetStream() options.Option {
	return func(request *testcontainers.ContainerRequest) error {
		request.Cmd = []string{
			cmdJetStreamEnabled,
		}
		return nil
	}
}
