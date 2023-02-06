package testmongo

import (
	"github.com/kyleishie/testdeps/pkg/options"
	"github.com/testcontainers/testcontainers-go"
)

const (
	env_MONGO_INITDB_ROOT_USERNAME = "MONGO_INITDB_ROOT_USERNAME"
	env_MONGO_INITDB_ROOT_PASSWORD = "MONGO_INITDB_ROOT_PASSWORD"
)

// WithRootUser sets MONGO_INITDB_ROOT_USERNAME & MONGO_INITDB_ROOT_PASSWORD to the given username & password.
func WithRootUser(username, password string) options.Option {
	return func(cr *testcontainers.ContainerRequest) error {
		if cr.Env == nil {
			cr.Env = make(map[string]string)
		}
		cr.Env[env_MONGO_INITDB_ROOT_USERNAME] = username
		cr.Env[env_MONGO_INITDB_ROOT_PASSWORD] = password
		return nil
	}
}
