package search

var engines = map[string]Searcher{}

func Register(name string, engine Searcher) {
	engines[name] = engine
}

func Unregister(name string) {
	if _, ok := engines[name]; ok {
		delete(engines, name)
	}
}

func Get(name string) Searcher {
	if search, ok := engines[name]; ok {
		return search
	}
	return DefaultSearch
}
