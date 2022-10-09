package top

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsMobile(t *testing.T) {
	r := IsMobile(`Mozilla/5.0 (iPhone; CPU OS 10_14 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/11.1.1 Mobile/14E304 Safari/605.1.15`)
	assert.True(t, r)
	r = IsMobile(`iPhone`)
	assert.True(t, r)
	r = IsMobile(`Mozilla/5.0 (Linux; Android 9.0; SAMSUNG-SM-T377A Build/NMF26X) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Mobile Safari/537.36`)
	assert.True(t, r)
	r = IsMobile(`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/11.1.1 Safari/605.1.15`)
	assert.False(t, r)
}
