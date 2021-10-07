package nats

import (
	tc "github.com/testcontainers/testcontainers-go"
	"testdeps/pkg/options"
)

// WithJetStream enables JetStream in the NATS container.
// This is the equivalent of running nats -js.
func WithJetStream() options.Option {
	return func(request *tc.ContainerRequest) error {
		request.Cmd = append(request.Cmd, cmdJetStreamEnabled)
		return nil
	}
}
