package segment

import (
	"regexp"
	"strings"
)

var (
	reHanExists = regexp.MustCompile(`\p{Han}+`)
	reNoWords   = regexp.MustCompile(`^([\d]+[\w]*|[\w]+[\d]+)$`)
	// DictFile 分词词典文件
	DictFile           string
	lastLoadedDictFile string
)

func init() {
	AddFilter(func(v string) bool {
		switch {
		case reHanExists.MatchString(v):
			if len(v) < 6 { // 忽略掉单字
				return false
			}
		case reNoWords.MatchString(v):
			return false
		default:
			return true
		}
		return true
	})
}

// SplitWords 分词
func SplitWords(b []byte, args ...string) []string {
	if len(DictFile) > 0 && lastLoadedDictFile != DictFile {
		lastLoadedDictFile = DictFile
		ReloadDict()
	}
	return Default().Segment(string(b), args...)
}

// SplitWordsBy 按模式分词
func SplitWordsBy(b []byte, mode string, args ...string) []string {
	if len(DictFile) > 0 && lastLoadedDictFile != DictFile {
		lastLoadedDictFile = DictFile
		ReloadDict()
	}
	return Default().SegmentBy(string(b), mode, args...)
}

// ReloadDict 重新加载词典
func ReloadDict(dictFiles ...string) error {
	dictFile := DictFile
	if len(dictFiles) > 0 {
		dictFile = dictFiles[0]
	}
	return Default().LoadDict(dictFile)
}

// SplitWordsAsString 将分词结果作为字串返回
func SplitWordsAsString(b []byte, args ...string) string {
	words := SplitWords(b, args...)
	content := strings.Join(words, ` `)
	content = strings.TrimSpace(content)
	return content
}
