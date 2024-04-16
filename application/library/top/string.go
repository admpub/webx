package top

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"
)

func TrimOverflowText(text string, maxLength int, seperators ...string) string {
	var seperator string
	if len(seperators) > 0 {
		seperator = seperators[0]
	}
	if len(text) <= maxLength {
		return text
	}
	text = text[:maxLength]
	if len(seperator) > 0 {
		if p := strings.LastIndex(text, seperator); p > 0 {
			text = text[0:p]
		}
	}
	return text
}

func OutputContent(content string, contypes ...string) interface{} {
	var contype string
	if len(contypes) > 0 {
		contype = contypes[0]
	}
	switch contype {
	case `text`:
		return template.HTML(tplfunc.Nl2br(content))
	case `html`:
		return template.HTML(content)
	case `markdown`:
		return template.HTML(`<textarea class="markdown-preview-text hidden">` + com.HTMLEncode(content) + `</textarea>`)
	default:
		return content
	}
}

var (
	hideTagRegex     = regexp.MustCompile(`(?s)(?i)\[hide(\:[^]]+)?\](.*?)\[/hide\]`)
	parseTagRegex    = regexp.MustCompile(`(?s)(?i)\[parse\](.*?)\[/parse\]`)
	brTagRegex       = regexp.MustCompile(`(?s)(?i)(^<br[ ]*[/]?>|<br[ ]*[/]?>$)`)
	defaultMsgOnHide = `此处内容需要评论回复后方可阅读`
)

type HideDetector func(hideType string, hideContent string, args ...string) (hide bool, msgOnHide string)

func parseGoTemplateContent(content string, funcMap map[string]interface{}) string {
	return parseTagRegex.ReplaceAllStringFunc(content, func(v string) string {
		if len(v) <= 15 { // [parse][/parse]
			return ``
		}
		index := strings.Index(v, `]`)
		v = v[index+1:]
		v = v[0 : len(v)-8]
		name := com.Md5(v)
		t := template.New(name)
		if funcMap != nil {
			t.Funcs(funcMap)
		}
		defer func() {
			if e := recover(); e != nil {
				v = fmt.Sprintf(`%v`, e)
			}
		}()
		_, err := t.Parse(v)
		if err != nil {
			err = echo.ParseTemplateError(err, v)
			v = err.Error()
		}
		return v
	})
}

func HideContent(content string, contype string, hide HideDetector, funcMap map[string]interface{}) (result string) {

	if hide == nil {
		hide = func(hideType string, hideContent string, args ...string) (hide bool, msgOnHide string) {
			return true, defaultMsgOnHide
		}
	}
	var hideStart, hideEnd string
	var showStart, showEnd string
	filter := func(v string) string {
		return v
	}
	switch contype {
	case `text`:
		hideStart, hideEnd = `[ `, ` ]`
	case `html`:
		hideStart, hideEnd = `<div class="hide-block-content show-after-comment mission-uncompleted">`, `</div>`
		showStart, showEnd = `<div class="hide-block-content show-after-comment mission-completed">`, `</div>`
		filter = func(v string) string {
			v = strings.TrimSpace(v)
			return brTagRegex.ReplaceAllString(v, "")
		}
	case `markdown`:
		//hideStart, hideEnd = "\n> **", "**\n"
		hideStart, hideEnd = `<div class="hide-block-content show-after-comment mission-uncompleted">`, `</div>`
	default:
		hideStart, hideEnd = `[ `, ` ]`
	}
	result = hideTagRegex.ReplaceAllStringFunc(content, func(v string) string {
		if len(v) <= 13 { // [hide][/hide]
			return ``
		}
		index := strings.Index(v, `]`)
		tagStart := v[0:index]
		v = v[index+1:]
		v = v[0 : len(v)-7]
		splited := strings.Split(tagStart, `:`) // hide:<hideType>:<arg1>:<...arg2>
		hideType := `comment`
		var args []string
		if len(splited) > 1 {
			hideType = splited[1]
		}
		if len(splited) > 2 {
			args = splited[2:]
		}
		v = filter(v)
		hidden, msgOnHide := hide(hideType, v, args...)
		if hidden {
			if len(msgOnHide) == 0 {
				msgOnHide = defaultMsgOnHide
			}
			return hideStart + msgOnHide + hideEnd
		}
		return showStart + parseGoTemplateContent(v, funcMap) + showEnd
	})
	return
}

func MakeKVCallback(cb func(k interface{}, v interface{}) error, args ...interface{}) (err error) {
	var k interface{}
	for i, j := 0, len(args); i < j; i++ {
		if i%2 == 0 {
			k = args[i]
			continue
		}
		if err = cb(k, args[i]); err != nil {
			return
		}
		k = nil
	}
	if k != nil {
		err = cb(k, nil)
	}
	return
}
