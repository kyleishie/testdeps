package options

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go"
)

func TestWithCustomTag(t *testing.T) {
	t.Run("version number", func(t *testing.T) {
		testTag := "1.0.1"
		fn := WithCustomTag(testTag)
		cReq := tc.ContainerRequest{
			Image: "banana",
		}
		err := fn(&cReq)
		assert.NoError(t, err)
		assert.True(t, strings.HasSuffix(cReq.Image, testTag))
		assert.Len(t, strings.Split(cReq.Image, ":"), 2)
	})
	t.Run("latest", func(t *testing.T) {
		testTag := "latest"
		fn := WithCustomTag(testTag)
		cReq := tc.ContainerRequest{
			Image: "banana",
		}
		err := fn(&cReq)
		assert.NoError(t, err)
		assert.True(t, strings.HasSuffix(cReq.Image, testTag))
		assert.Len(t, strings.Split(cReq.Image, ":"), 2)
	})
	t.Run("empty", func(t *testing.T) {
		testTag := ""
		fn := WithCustomTag(testTag)
		cReq := tc.ContainerRequest{
			Image: "banana",
		}
		err := fn(&cReq)
		assert.NoError(t, err)
		assert.NotContains(t, cReq.Image, ":")
	})
	t.Run("has : prefix", func(t *testing.T) {
		testTag := ":1.0.1"
		fn := WithCustomTag(testTag)
		cReq := tc.ContainerRequest{
			Image: "banana",
		}
		err := fn(&cReq)
		assert.NoError(t, err)
		assert.True(t, strings.HasSuffix(cReq.Image, testTag))
		assert.Len(t, strings.Split(cReq.Image, ":"), 2)
	})
	t.Run("image empty", func(t *testing.T) {
		testTag := ":1.0.1"
		fn := WithCustomTag(testTag)
		cReq := tc.ContainerRequest{}
		err := fn(&cReq)
		assert.Error(t, err)
		assert.Empty(t, cReq.Image)
	})
}
