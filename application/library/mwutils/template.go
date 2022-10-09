package mwutils

import "github.com/webx-top/echo"

type TmplFuncGenerator func(ctx echo.Context) interface{}
type TmplFuncGenerators map[string]TmplFuncGenerator

func (t *TmplFuncGenerators) Add(name string, gen TmplFuncGenerator) {
	(*t)[name] = gen
}

func (t *TmplFuncGenerators) Delete(name string) {
	delete(*t, name)
}

func (t *TmplFuncGenerators) Get(name string) TmplFuncGenerator {
	return (*t)[name]
}

func (t *TmplFuncGenerators) Apply(ctx echo.Context) {
	for name, funcGenerator := range *t {
		ctx.SetFunc(name, funcGenerator(ctx))
	}
}
