package mongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go"
)

func TestWithRootUser(t *testing.T) {
	testUser := "user"
	testPass := "pass"
	fn := WithRootUser(testUser, testPass)
	cReq := tc.ContainerRequest{}
	err := fn(&cReq)
	assert.NoError(t, err)
	assert.Equal(t, testUser, cReq.Env["MONGO_INITDB_ROOT_USERNAME"])
	assert.Equal(t, testPass, cReq.Env["MONGO_INITDB_ROOT_PASSWORD"])
}
