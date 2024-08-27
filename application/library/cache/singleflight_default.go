package cache

var defaultSingleflight = Singleflight()

func SingleflightDo(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	return defaultSingleflight.Do(key, fn)
}

func SingleflightDoChan(key string, fn func() (interface{}, error)) <-chan SinglefightResult {
	return defaultSingleflight.DoChan(key, fn)
}

func SingleflightForget(key string) {
	defaultSingleflight.Forget(key)
}
