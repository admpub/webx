package top

import "github.com/webx-top/echo"

func HeaderP3P(ctx echo.Context) {
	ctx.Response().Header().Set(`P3P`, `CP="CURa ADMa DEVa PSAo PSDo OUR BUS UNI PUR INT DEM STA PRE COM NAV OTC NOI DSP COR"`)
}
