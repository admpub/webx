package top

import (
	"strings"

	"github.com/admpub/pinyin-golang/pinyin"
)

var PinyinDict = pinyin.NewDict()

func FirstLetter(s string) string {
	firstLetters := PinyinDict.Abbr(s, ``)
	if len(firstLetters) > 1 {
		return strings.ToUpper(firstLetters[0:1])
	}
	if len(firstLetters) == 1 {
		return strings.ToUpper(firstLetters)
	}
	return ``
}

func Pinyin(s string, seps ...string) string {
	var sep string
	if len(seps) > 0 {
		sep = seps[0]
	}
	result := PinyinDict.Convert(s, sep)
	return result.None()
}
