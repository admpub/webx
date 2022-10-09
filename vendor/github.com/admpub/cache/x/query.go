package x

// QueryFunc 查询过程签名
type QueryFunc func() error

// Query 查询过程实现Querier接口
func (q QueryFunc) Query() error {
	return q()
}

// Querier 查询接口
type Querier interface {
	Query() error
}
