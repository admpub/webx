package top

import (
	"html/template"
	"testing"

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
