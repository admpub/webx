package top

import (
	"html/template"
	"testing"
	"time"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/stretchr/testify/assert"
)

func TestOutputContent(t *testing.T) {
	v := OutputContent(`1 1 2 3`, `text`)
	assert.Equal(t, template.HTML(`1 1 2 3`), v)

	v = OutputContent(`1 &nbsp; 1 2 3`, `text`)
	assert.Equal(t, template.HTML(`1 &amp;nbsp; 1 2 3`), v)

	v = OutputContent("1 \r\n 1 2 3", `text`)
	assert.Equal(t, template.HTML(`1 <br /> 1 2 3`), v)
}

func TestMakeEncodedURL(t *testing.T) {
	cfg := config.NewConfig().SetDefaults(``)
	cfg.AsDefault()
	expiry := time.Now().Add(time.Hour * 24).Unix()
	rawURL := `https://coscms.com/download?a=测试&b=1#c`
	v, err := MakeEncodedURL(rawURL, expiry)
	assert.NoError(t, err)
	assert.True(t, len(v) > 0)
	_rawURL, _expiry, err := ParseEncodedURL(v)
	assert.NoError(t, err)
	assert.Equal(t, _rawURL, rawURL)
	assert.Equal(t, _expiry, expiry)
}
