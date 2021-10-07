package nats

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go"
)

func TestWithJetStream(t *testing.T) {
	cReq := tc.ContainerRequest{}
	fn := WithJetStream()
	err := fn(&cReq)
	assert.NoError(t, err)
	assert.Contains(t, cReq.Cmd, "-js")
}
