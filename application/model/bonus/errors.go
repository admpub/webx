package bonus

import "errors"

var (
	ErrDeadLoop            = errors.New(`推荐人关系异常，同一推荐人多次出现，已经构成死循环`)
	ErrAgentLevelNotExists = errors.New(`代理等级不存在`)
)
