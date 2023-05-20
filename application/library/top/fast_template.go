package top

import (
	"github.com/admpub/fasttemplate"
	"github.com/admpub/webx/application/library/fasttemplates"
	"github.com/webx-top/echo/param"
)

var fastTemplates = fasttemplates.New()

func ReplacePlaceholder(content string, values map[string]interface{}, startAndEndTag ...string) string {
	startTag := "{"
	endTag := "}"
	if len(startAndEndTag) > 0 {
		if len(startAndEndTag[0]) > 0 {
			startTag = startAndEndTag[0]
		}
		if len(startAndEndTag) > 1 && len(startAndEndTag[1]) > 0 {
			startTag = startAndEndTag[1]
		}
	}
	for k, v := range values {
		values[k] = param.AsString(v)
	}
	t := fastTemplates.GetOrSetBy(content, startTag, endTag)
	return t.ExecuteString(values)
}

func ReplacePlaceholderx(t *fasttemplate.Template, values map[string]interface{}, args ...interface{}) string {
	for k, v := range values {
		values[k] = param.AsString(v)
	}
	var k string
	for i, j := 0, len(args); i < j; i++ {
		if i%2 == 0 {
			k = param.AsString(args[i])
			continue
		}
		values[k] = param.AsString(args[i])
		k = ``
	}
	if len(k) > 0 {
		values[k] = ``
	}
	return t.ExecuteString(values)
}
