package common

import "time"

const (
	// DefaultConnTimeout is the standard timeout that all containers should use when attempting a connection.
	DefaultConnTimeout = time.Minute * 2
)
