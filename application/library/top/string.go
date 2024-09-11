package top

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"
	"github.com/webx-top/echo/param"
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
	expiryTagRegex   = regexp.MustCompile(`(?s)(?i)\[expiry(\:[^]]+)?\](.*?)\[/expiry\]`)
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
			return v
		}
		buf := bytes.NewBuffer(nil)
		err = t.Execute(buf, nil)
		if err != nil {
			v = err.Error()
		} else {
			v = buf.String()
		}
		return v
	})
}

func parseExpiryContent(content string) string {
	return expiryTagRegex.ReplaceAllStringFunc(content, func(v string) string {
		if len(v) <= 17 { // [expiry][/expiry]
			return ``
		}
		index := strings.Index(v, `]`)
		tagStart := v[0:index]
		v = v[index+1:]
		v = v[0 : len(v)-9]
		splited := strings.Split(tagStart, `:`) // expiry:<duration>:<linkTitle>
		var duration, linkTitle string
		if len(splited) > 1 {
			duration = splited[1]
		}
		if len(splited) > 2 {
			linkTitle = splited[2]
		}
		result := MakeEncodedURLOrLink(v, duration, linkTitle)
		switch r := result.(type) {
		case string:
			return r
		case template.HTML:
			return string(r)
		default:
			return fmt.Sprint(r)
		}
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
		v = parseExpiryContent(v)
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

const PrefixEncoded = `encoded:`

func MakeEncodedURL(urlStr string, expiry int64) (string, error) {
	d := echo.H{`url`: urlStr, `expiry`: expiry}
	b, err := json.Marshal(d)
	if err != nil {
		return urlStr, err
	}
	urlStr = `/article/redirect?url=` +
		url.QueryEscape(PrefixEncoded+config.FromFile().Encode256(com.Bytes2str(b))) +
		`&expiry=` + param.AsString(expiry)
	return urlStr, err
}

func MakeEncodedURLOrLink(url string, expiry interface{}, linkTitle ...string) interface{} {
	var err error
	var _expiry int64
	switch exp := expiry.(type) {
	case string:
		if len(exp) > 0 {
			dur, err := ParseDuration(exp)
			if err != nil {
				return err.Error()
			}
			_expiry = int64(time.Now().Add(dur).Unix())
		} else {
			_expiry = int64(time.Now().AddDate(0, 0, 1).Unix())
		}
	case int64:
		_expiry = exp
	case uint64:
		_expiry = int64(exp)
	case int:
		_expiry = int64(exp)
	case uint:
		_expiry = int64(exp)
	case time.Duration:
		_expiry = int64(time.Now().Add(exp).Unix())
	default:
		_expiry = int64(time.Now().AddDate(0, 0, 1).Unix())
	}
	url = strings.ReplaceAll(url, `&amp;`, `&`)
	url, err = MakeEncodedURL(url, _expiry)
	if err != nil {
		url = err.Error()
		return url
	}
	if len(linkTitle) > 0 && len(linkTitle[0]) > 0 {
		return template.HTML(`<a href="` + url + `" target="_blank">` + linkTitle[0] + `</a>`)
	}
	return url
}

func ParseEncodedURL(encodedURL string) (string, int64, error) {
	rawURL := encodedURL
	var expiry int64
	if strings.HasPrefix(rawURL, PrefixEncoded) {
		rawURL = strings.TrimPrefix(rawURL, PrefixEncoded)
		rawURL = config.FromFile().Decode256(rawURL)
		if len(rawURL) == 0 {
			return rawURL, expiry, nil
		}
		data := echo.H{}
		jsonBytes := com.Str2bytes(rawURL)
		err := json.Unmarshal(jsonBytes, &data)
		if err != nil {
			return rawURL, expiry, common.JSONBytesParseError(err, jsonBytes)
		}
		rawURL = data.String(`url`)
		if len(rawURL) == 0 {
			return rawURL, expiry, nil
		}
		expiry = data.Int64(`expiry`)
	}
	return rawURL, expiry, nil
}
