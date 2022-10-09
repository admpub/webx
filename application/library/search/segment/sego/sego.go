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

package sego

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/admpub/log"
	. "github.com/admpub/webx/application/library/search/segment"
	"github.com/huichen/sego"
	"github.com/webx-top/echo"
)

func init() {
	Register(`sego`, New)
}

func New() Segment {
	return &Sego{
		segmenter:   &sego.Segmenter{},
		defaultDict: filepath.Join(echo.Wd(), `data`, `sego/dict.txt`),
	}
}

type Sego struct {
	segmenter   *sego.Segmenter
	defaultDict string
	dictLoaded  bool
	once        sync.Once
}

func (s *Sego) LoadDict(dictFile string, dictType ...string) error {
	log.Debug(`[sego]LoadDict:`, dictFile)
	s.segmenter.LoadDictionary(dictFile) //多个字典文件用半角“,”逗号分隔
	s.dictLoaded = true
	return nil
}

func (s *Sego) Segment(text string, args ...string) []string {
	if !s.dictLoaded {
		s.once.Do(func() {
			s.LoadDict(s.defaultDict, `default`)
		})
	}
	segments := s.segmenter.Segment([]byte(text))
	var (
		words     []string
		wordTypes []string //获取指定类型的词语,如仅仅获取名词则为n
	)
	if len(args) > 0 {
		wordTypes = args
	}
	typeLength := len(wordTypes)
	for _, seg := range segments {
		//排除指定词性的词
		if typeLength > 0 {
			var ok bool
			for _, wt := range wordTypes {
				if seg.Token().Pos() == wt {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}
		content := seg.Token().Text()
		content = strings.Replace(content, `　`, ``, -1)
		content = strings.TrimSpace(content)
		if len(content) > 0 && DoFilter(content) {
			words = append(words, content)
		}
	}
	words = CleanStopWordsFromSlice(words)
	return words
}

func (s *Sego) SegmentBy(text string, mode string, args ...string) []string {
	return s.Segment(text, args...)
}

func (s *Sego) Close() error {
	return nil
}
