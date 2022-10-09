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

package jiebago

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/admpub/log"
	. "github.com/admpub/webx/application/library/search/segment"
	"github.com/wangbin/jiebago"
	"github.com/webx-top/echo"
)

func init() {
	Register(`jiebago`, New)
}

func New() Segment {
	return &Jieba{
		segmenter:   &jiebago.Segmenter{},
		defaultDict: filepath.Join(echo.Wd(), `data`, `sego/dict.txt`),
	}
}

type Jieba struct {
	segmenter   *jiebago.Segmenter
	defaultDict string
	dictLoaded  bool
	once        sync.Once
}

func (s *Jieba) LoadDict(dictFile string, dictType ...string) error {
	for index, file := range strings.Split(dictFile, `,`) { //多个字典文件用半角“,”逗号分隔
		var err error
		if index == 0 {
			log.Debug(`[jiebago]Load dictionary:`, file)
			err = s.segmenter.LoadDictionary(file)
		} else {
			log.Debug(`[jiebago]Load user dictionary:`, file)
			err = s.segmenter.LoadUserDictionary(file)
		}
		if err != nil {
			log.Error(`[jiebago]LoadDict:`, err)
		}
	}
	s.dictLoaded = true
	return nil
}

func (s *Jieba) Segment(text string, args ...string) []string {
	if !s.dictLoaded {
		s.once.Do(func() {
			s.LoadDict(s.defaultDict, `default`)
		})
	}
	var (
		words = []string{}
		ch    <-chan string //精确模式
	)
	ch = s.segmenter.Cut(text, false)

	for word := range ch {
		word = strings.TrimSpace(word)
		if len(word) > 0 && DoFilter(word) {
			words = append(words, word)
		}
	}
	words = CleanStopWordsFromSlice(words)
	return words
}

func (s *Jieba) SegmentBy(text string, mode string, args ...string) []string {
	if !s.dictLoaded {
		s.once.Do(func() {
			s.LoadDict(s.defaultDict, `default`)
		})
	}
	var (
		words = []string{}
		ch    <-chan string //精确模式
	)
	switch mode {
	case `all`:
		//log.Println(`all mode:`, text)
		ch = s.segmenter.CutAll(text)
	case `new`: //新词识别
		//log.Println(`new mode:`, text)
		ch = s.segmenter.Cut(text, true)
	case `search`: //搜索引擎模式
		//log.Println(`search mode:`, text)
		ch = s.segmenter.CutForSearch(text, true)

	//TODO
	//case `tag`: //词性标注
	//case `keywords`: //关键词提取

	default: //精确模式
		ch = s.segmenter.Cut(text, false)
	}
	for word := range ch {
		word = strings.TrimSpace(word)
		if len(word) > 0 && DoFilter(word) {
			words = append(words, word)
		}
	}
	words = CleanStopWordsFromSlice(words)
	return words
}

func (s *Jieba) Close() error {
	return nil
}
