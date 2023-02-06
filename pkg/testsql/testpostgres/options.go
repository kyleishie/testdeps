package testpostgres

import (
	"github.com/kyleishie/testdeps/pkg/options"
	"github.com/testcontainers/testcontainers-go"
)

const (
	env_POSTGRES_HOST_AUTH_METHOD = "POSTGRES_HOST_AUTH_METHOD"
	env_POSTGRES_USER             = "POSTGRES_USER"
	env_POSTGRES_PASSWORD         = "POSTGRES_PASSWORD"
	env_POSTGRES_DB               = "POSTGRES_DB"
)

func WithTrust() options.Option {
	return func(cr *testcontainers.ContainerRequest) error {
		if cr.Env == nil {
			cr.Env = make(map[string]string)
		}
		cr.Env[env_POSTGRES_HOST_AUTH_METHOD] = "trust"
		return nil
	}
}

// WithPassword sets POSTGRES_PASSWORD to the given password.
func WithPassword(password string) options.Option {
	return func(cr *testcontainers.ContainerRequest) error {
		if cr.Env == nil {
			cr.Env = make(map[string]string)
		}
		cr.Env[env_POSTGRES_PASSWORD] = password
		return nil
	}
}

// WithUser sets POSTGRES_USER to the given username.
func WithUser(user string) options.Option {
	return func(cr *testcontainers.ContainerRequest) error {
		if cr.Env == nil {
			cr.Env = make(map[string]string)
		}
		cr.Env[env_POSTGRES_USER] = user
		return nil
	}
}

// WithInitialDatabase sets POSTGRES_DB
func WithInitialDatabase(dbName string) options.Option {
	return func(cr *testcontainers.ContainerRequest) error {
		if cr.Env == nil {
			cr.Env = make(map[string]string)
		}
		cr.Env[env_POSTGRES_DB] = dbName
		return nil
	}
}
