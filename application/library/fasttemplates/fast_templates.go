package fasttemplates

import (
	"sync"

	"github.com/admpub/fasttemplate"
)

func New() *FastTemplates {
	return &FastTemplates{
		m: map[string]*fasttemplate.Template{},
		l: sync.RWMutex{},
	}
}

type FastTemplates struct {
	m map[string]*fasttemplate.Template
	l sync.RWMutex
}

func (f *FastTemplates) Set(key string, t *fasttemplate.Template) {
	f.l.Lock()
	f.m[key] = t
	f.l.Unlock()
}

func (f *FastTemplates) SetBy(content, startTag, endTag string) {
	k := content + `,` + startTag + endTag
	t := fasttemplate.New(content, startTag, endTag)
	f.Set(k, t)
}

func (f *FastTemplates) Get(key string) (t *fasttemplate.Template) {
	f.l.RLock()
	t = f.m[key]
	f.l.RUnlock()
	return
}

func (f *FastTemplates) GetBy(content, startTag, endTag string) *fasttemplate.Template {
	k := content + `,` + startTag + endTag
	return f.Get(k)
}

func (f *FastTemplates) GetOrSetBy(content, startTag, endTag string) *fasttemplate.Template {
	k := content + `,` + startTag + endTag
	t := f.Get(k)
	if t == nil {
		t = fasttemplate.New(content, startTag, endTag)
		f.Set(k, t)
	}
	return t
}
