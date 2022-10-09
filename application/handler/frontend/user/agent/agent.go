package agent

import (
	"github.com/webx-top/echo"
)

// 代理

// AgentIndex 代理中心首页
var AgentIndex = func(c echo.Context) error {
	var err error
	return c.Render(`/user/agent/index`, err)
}

// AgentEdit 修改代理资料
var AgentEdit = func(c echo.Context) error {
	var err error
	return c.Render(`/user/agent/edit`, err)
}
