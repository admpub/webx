package xtemplate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRGB2Hex(t *testing.T) {
	v := RGB2Hex(`rgb(190, 53, 220);`)
	expected := `#BE35DC`
	assert.Equal(t, expected, v)

	v = RGB2Hex(`rgb(190, 53, 220)`)
	assert.Equal(t, expected, v)

	v = RGB2Hex(`rgba(190, 53, 220, 0.1)`)
	assert.Equal(t, expected, v)
}
