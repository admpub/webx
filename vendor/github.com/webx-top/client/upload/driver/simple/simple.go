/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package simple

import (
	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/echo"
)

func init() {
	uploadClient.Register(`simple`, func() uploadClient.Client {
		return New()
	})
}

var FormField = `filedata`

func New() uploadClient.Client {
	client := &Simple{}
	client.BaseClient = uploadClient.New(client, FormField)
	return client
}

type Simple struct {
	*uploadClient.BaseClient
}

func (a *Simple) BuildResult() uploadClient.Client {
	status := 1
	var errMsg string
	if a.GetError() != nil {
		status = 0
		errMsg = a.ErrorString()
	}
	a.RespData = echo.H{
		`Status`:  status,
		`Message`: errMsg,
		`Data`: echo.H{
			`Url`: a.Data.FileURL,
			`Id`:  a.Data.FileIDString(),
		},
	}
	return a
}
