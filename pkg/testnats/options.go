package testnats

import (
	"github.com/kyleishie/testdeps/pkg/options"
	tc "github.com/testcontainers/testcontainers-go"
)

// WithJetStream enables JetStream in the NATS Container.
// This is the equivalent of running nats -js.
func WithJetStream() options.Option {
	return func(request *tc.ContainerRequest) error {
		request.Cmd = append(request.Cmd, cmdJetStreamEnabled)
		return nil
	}
}
