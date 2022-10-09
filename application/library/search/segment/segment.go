/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package segment

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/once"
	"github.com/webx-top/echo"
)

type Filter func(string) bool

var (
	DefaultEngine  = `sego` // gojieba / sego / jiebago
	stopWords      []string
	stopWordsMap   map[string]bool
	Filters        []Filter
	defaultSegment Segment
	onceSegment    once.Once
	onceStopword   once.Once
)

func initDefaultSegment() {
	log.Debug("[segment]Default engine:", DefaultEngine)
	defaultSegment = Get(DefaultEngine)
}

func IsInitialized() bool {
	return defaultSegment != nil
}

func Default() Segment {
	onceSegment.Do(initDefaultSegment)
	return defaultSegment
}

func ResetSegment() {
	onceSegment.Reset()
}

func ResetStopwords() {
	onceStopword.Reset()
}

func StopWords() []string {
	onceStopword.Do(initLoadStopWordsDict)
	return stopWords
}

func LoadStopWordsDict(stopWordsFile string, args ...bool) {
	var rebuild bool
	if len(args) > 0 {
		rebuild = args[0]
	}
	b, err := os.ReadFile(stopWordsFile)
	if err != nil {
		log.Debug(stopWordsFile+`:`, err)
		return
	}
	words := strings.Split(strings.TrimSpace(string(b)), "\n")
	if rebuild {
		stopWords = []string{}
		stopWordsMap = nil
	}
	for _, word := range words {
		word = strings.TrimSpace(word)
		if len(word) > 0 {
			stopWords = append(stopWords, word)
		}
	}
}

func initLoadStopWordsDict() {
	var stopWordsFiles []string
	stopWordsFiles = append(stopWordsFiles, filepath.Join(echo.Wd(), `data`, `sego/stopwords.txt`))
	goPath := os.Getenv(`GOPATH`)
	if len(goPath) > 0 {
		stopWordsFiles = append(stopWordsFiles, filepath.Join(goPath, `src`, `github.com/admpub/webx/application/library/search/segment/stopwords.txt`))
	}
	for _, stopWordsFile := range stopWordsFiles {
		_, err := os.Stat(stopWordsFile)
		if err == nil {
			LoadStopWordsDict(stopWordsFile)
			return
		}
		log.Debug(stopWordsFile+`:`, err)
	}
}

func CleanStopWords(v string) string {
	for _, word := range StopWords() {
		v = strings.Replace(v, word, ` `, -1)
	}
	return v
}

func CleanStopWordsFromSlice(v []string) (r []string) {
	if stopWordsMap == nil {
		stopWordsMap = make(map[string]bool)
		for _, word := range StopWords() {
			stopWordsMap[word] = true
		}
	}
	r = make([]string, 0)
	for _, w := range v {
		if _, ok := stopWordsMap[w]; !ok {
			r = append(r, w)
		}
	}
	return r
}

func DoFilter(v string) bool {
	for _, f := range Filters {
		if !f(v) {
			return false
		}
	}
	return true
}

func AddFilter(filter Filter) {
	Filters = append(Filters, filter)
}

// Segment interface
type Segment interface {
	//载入词典（词典路径，词典类型）
	LoadDict(string, ...string) error

	//分词（文本，词性）
	Segment(string, ...string) []string

	//分词（文本，分词模式，词性）
	SegmentBy(string, string, ...string) []string

	//关闭或释放资源
	Close() error
}

type nopSegment struct {
}

func (s *nopSegment) LoadDict(dictFile string, dictType ...string) error {
	return nil
}

func (s *nopSegment) Segment(text string, args ...string) []string {
	return []string{}
}

func (s *nopSegment) SegmentBy(text string, mode string, args ...string) []string {
	return []string{}
}

func (s *nopSegment) Close() error {
	return nil
}
