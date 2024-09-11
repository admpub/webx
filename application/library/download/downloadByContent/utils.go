package downloadByContent

import (
	"regexp"
	"sync"

	uploadLibrary "github.com/coscms/webcore/library/upload"
)

var (
	htmlImageRegexp     *regexp.Regexp
	markdownImageRegexp *regexp.Regexp
	once                sync.Once
)

func initRegexp(forces ...bool) {
	var force bool
	if len(forces) > 0 {
		force = forces[0]
	}
	if !force && htmlImageRegexp != nil {
		return
	}
	ruleEnd := uploadLibrary.ExtensionRegexpEnd(true)
	htmlImageRegexp = regexp.MustCompile(`(?s)(?i)<img[\s](?:[^>]+\s)?src=["'](https?://[^'"]+` + ruleEnd + `)["'][^>]*>`)
	markdownImageRegexp = regexp.MustCompile(`(?i)!\[[^]]*\]\((https?://[^\s\)]+` + ruleEnd + `)[^\)]*\)`)
	//println(`OutsideImage: initialize regexp`)
}

func onceInitRegexp() {
	initRegexp()
}

func OutsideImage(content string, contentType string) map[string]string {
	once.Do(onceInitRegexp)
	result := map[string]string{}
	switch contentType {
	case `markdown`:
		matches := markdownImageRegexp.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			result[match[0]] = match[1]
		}
		fallthrough
	case `html`:
		matches := htmlImageRegexp.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			result[match[0]] = match[1]
		}
	}
	return result
}
