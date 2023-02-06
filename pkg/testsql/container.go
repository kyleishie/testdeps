package testsql

import (
	tc "github.com/testcontainers/testcontainers-go"
)

// / Container represents a docker container
type Container struct {
	tc.Container
	ConnectionString string
	driver           string
}

func New(c tc.Container, driver, connectionString string) *Container {
	return &Container{
		Container:        c,
		ConnectionString: connectionString,
		driver:           driver,
	}
}
