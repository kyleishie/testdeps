package options

import (
	"errors"
	"strings"

	"github.com/testcontainers/testcontainers-go"
)

// Option is a function that can customize the given testcontainers.ContainerRequest.
// The intended use of Option is to provide extension points for implementers of Container.
type Option func(*testcontainers.ContainerRequest) error

const imageTagDelimiter = ":"

// WithCustomTag sets the tag for the underlying docker image.
// If the tag is already set WithCustomTag replaces it.
func WithCustomTag(tag string) Option {
	return func(request *testcontainers.ContainerRequest) error {
		image := request.Image
		parts := strings.Split(image, imageTagDelimiter)
		if len(parts) == 0 || parts[0] == "" {
			return errors.New("image not set")
		}

		if tag == "" {
			request.Image = parts[0]
		} else {
			request.Image = parts[0] + imageTagDelimiter + strings.ReplaceAll(tag, imageTagDelimiter, "")
		}

		return nil
	}
}
