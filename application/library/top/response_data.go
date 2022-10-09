package top

import "github.com/webx-top/echo"

type ResponseData struct {
	Code  int         `json:"Code"`
	State string      `json:"State"`
	Info  interface{} `json:"Info,omitempty" xml:"Info,omitempty" swaggertype:"string"`
	Zone  interface{} `json:"Zone,omitempty" xml:"Zone,omitempty" swaggertype:"string"`
	URL   string      `json:"URL,omitempty" xml:"URL,omitempty"`
	Data  interface{} `json:"Data" swaggertype:"object"`
}

func Response(ctx echo.Context) error {
	data := ctx.Data()
	return ctx.JSON(data.SetData(ctx.Stored(), data.GetCode().Int()))
}
