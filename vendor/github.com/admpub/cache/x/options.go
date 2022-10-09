package x

// getOptions Get方法的可选参数项
type getOptions struct {
	querier Querier
	ttl     int64 // seconds

	// DisableCacheUsage disables the cache.
	// It can be useful during debugging.
	disableCacheUsage bool

	// UseFreshData will ignore content in the cache and always pull fresh data.
	// The pulled data will subsequently be saved in the cache.
	useFreshData bool
}

// GetOption Get方法的可选参数项结构，不需要直接调用。
type GetOption struct {
	apply func(options *getOptions)
}

// Query 为Get操作定制查询过程
func Query(querier Querier) GetOption {
	return GetOption{
		apply: func(options *getOptions) {
			options.querier = querier
		},
	}
}

// TTL 为Get操作定制TTL
func TTL(ttl int64) GetOption {
	return GetOption{
		apply: func(options *getOptions) {
			options.ttl = ttl
		},
	}
}

func DisableCacheUsage(disableCacheUsage bool) GetOption {
	return GetOption{
		apply: func(options *getOptions) {
			options.disableCacheUsage = disableCacheUsage
		},
	}
}

func UseFreshData(useFreshData bool) GetOption {
	return GetOption{
		apply: func(options *getOptions) {
			options.useFreshData = useFreshData
		},
	}
}
