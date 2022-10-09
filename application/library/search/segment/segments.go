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

import "github.com/admpub/log"

var (
	segments   = make(map[string]func() Segment)
	defaultNop = &nopSegment{}
)

func Register(name string, c func() Segment) {
	segments[name] = c
}

func IsNop(segment Segment) bool {
	return defaultNop == segment
}

func Get(name string) Segment {
	fn, ok := segments[name]
	if !ok || fn == nil {
		log.Error(`[segment]Not found engine:`, name)
		fn = func() Segment {
			return defaultNop
		}
	}
	return fn()
}

func Has(name string) bool {
	_, ok := segments[name]
	return ok
}

func Unregister(name string) {
	delete(segments, name)
}
