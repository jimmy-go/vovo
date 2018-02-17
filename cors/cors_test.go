package cors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	origin := "abc"
	h := Handler(origin)
	assert.NotNil(t, h)
}
