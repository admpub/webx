package top

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplacePlaceholder(t *testing.T) {
	result := ReplacePlaceholder(`00{a}1{b}2{c}99`, map[string]interface{}{
		`a`: `x`,
		`b`: `y`,
		`c`: `z`,
	})
	assert.Equal(t, `00x1y2z99`, result)
	result = ReplacePlaceholder(`00{a}1{b}2{c}99`, map[string]interface{}{
		`a`: 6,
		`b`: 7,
		`c`: 8,
	})
	assert.Equal(t, `006172899`, result)
}
