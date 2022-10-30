/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package license

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/admpub/errors"
	"github.com/admpub/log"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/restclient"
)

var (
	ErrConnectionFailed       = errors.New(`连接授权服务器失败`)
	ErrOfficialDataUnexcepted = errors.New(`官方数据返回异常`)
	ErrLicenseDownloadFailed  = errors.New(`下载证书失败：官方数据返回异常`)
)

type OfficialData struct {
	License   string
	Timestamp int64
}

type OfficialResponse struct {
	Code int
	Info string
	Zone string        `json:",omitempty" xml:",omitempty"`
	Data *OfficialData `json:",omitempty" xml:",omitempty"`
}

type ValidateResponse struct {
	Code int
	Info string
	Zone string    `json:",omitempty" xml:",omitempty"`
	Data Validator `json:",omitempty" xml:",omitempty"`
}

type Validator interface {
	Validate() error
}

var NewValidateResponse = func() *ValidateResponse {
	return &ValidateResponse{
		Data: ValidateResultInitor(),
	}
}

var ValidateResultInitor = func() Validator {
	return &ValidateResult{}
}

type ValidateResult struct {
}

func (v *ValidateResult) Validate() error {
	return nil
}

func validateFromOfficial(ctx echo.Context) error {
	client := restclient.Resty()
	client.SetHeader("Accept", "application/json")
	result := NewValidateResponse()
	client.SetResult(result)
	fullURL := FullLicenseURL(ctx)
	response, err := client.Get(fullURL)
	if err != nil {
		if strings.Contains(err.Error(), `connection refused`) {
			return ErrConnectionFailed
		}
		return errors.Wrap(err, `Connection to the license server failed`)
	}
	if response == nil {
		return ErrConnectionFailed
	}
	switch response.StatusCode() {
	case http.StatusOK:
		if result.Code != 1 {
			return errors.New(result.Info)
		}
		if result.Data == nil {
			return ErrOfficialDataUnexcepted
		}
		return result.Data.Validate()
	case http.StatusNotFound:
		return ErrConnectionFailed
	default:
		return errors.New(response.Status())
	}
}

type VersionResponse struct {
	Code int
	Info string
	Zone string          `json:",omitempty" xml:",omitempty"`
	Data *ProductVersion `json:",omitempty" xml:",omitempty"`
}

func LatestVersion(ctx echo.Context, forceDownload bool) (*ProductVersion, error) {
	client := restclient.Resty()
	client.SetHeader("Accept", "application/json")
	result := &VersionResponse{}
	client.SetResult(result)
	response, err := client.Get(versionURL + `?` + URLValues(ctx).Encode())
	if err != nil {
		return nil, errors.Wrap(err, `Check for the latest version failed`)
	}
	if response == nil {
		return nil, ErrConnectionFailed
	}
	switch response.StatusCode() {
	case http.StatusOK:
		if result.Code != 1 {
			return nil, errors.New(result.Info)
		}
		if result.Data == nil {
			return nil, ErrOfficialDataUnexcepted
		}
		hasNew := config.Version.IsNew(result.Data.Version, result.Data.Type)
		if hasNew {
			log.Okay(`New version found: v`, result.Data.Version)
			if forceDownload || result.Data.ForceUpgrade == `Y` {
				if len(result.Data.DownloadUrl) > 0 {
					log.Okay(`Automatically download the new version v`, result.Data.Version)
					saveTo := filepath.Join(echo.Wd(), `data/cache/nging-new-version`)
					err = com.MkdirAll(saveTo, os.ModePerm)
					if err != nil {
						return result.Data, err
					}
					saveTo += echo.FilePathSeparator + result.Data.Version + `_` + path.Base(result.Data.DownloadUrl)
					result.Data.DownloadedPath = saveTo
					if com.FileExists(saveTo) {
						log.Okay(`The file already exists: `, saveTo)
						return result.Data, nil
					}
					log.Okayf(`Downloading %s => %s`, result.Data.DownloadUrl, saveTo)
					err = com.RangeDownload(result.Data.DownloadUrl, saveTo)
					if err != nil {
						if len(result.Data.DownloadUrlOther) > 0 {
							log.Okay(`Try to download from the mirror URL `, result.Data.DownloadUrlOther)
							err = com.RangeDownload(result.Data.DownloadUrlOther, saveTo)
						}
					}
					if err != nil {
						return result.Data, err
					}
					var signList []string
					if len(result.Data.Sign) > 0 {
						signList = strings.Split(result.Data.Sign, `,`)
					}
					if len(signList) > 0 {
						fileMd5 := com.Md5file(saveTo)
						var matched bool
						for _, sign := range signList {
							if sign == fileMd5 {
								matched = true
								break
							}
						}
						if !matched {
							return result.Data, com.ErrMd5Unmatched
						}
					}
					//OK
				}
			}
			return result.Data, nil
		}

		log.Okay(`No new version`)
		return result.Data, nil
	case http.StatusNotFound:
		return nil, ErrConnectionFailed
	default:
		return nil, errors.New(response.Status())
	}
}
