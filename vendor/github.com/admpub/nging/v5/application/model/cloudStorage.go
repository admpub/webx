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

package model

import (
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/common"
)

func NewCloudStorage(ctx echo.Context) *CloudStorage {
	m := &CloudStorage{
		NgingCloudStorage: dbschema.NewNgingCloudStorage(ctx),
	}
	return m
}

type CloudStorage struct {
	*dbschema.NgingCloudStorage
}

func (s *CloudStorage) check() error {
	s.Bucket = strings.TrimSpace(s.Bucket)
	s.Endpoint = strings.TrimSpace(s.Endpoint)
	if len(s.Baseurl) > 0 {
		s.Baseurl = strings.TrimSuffix(s.Baseurl, `/`)
	}
	s.Secret = strings.TrimSpace(s.Secret)
	s.Secret = common.Crypto().Encode(s.Secret)
	return nil
}

func (s *CloudStorage) RawSecret() string {
	return common.Crypto().Decode(s.Secret)
}

func (s *CloudStorage) BaseURL() string {
	if len(s.Baseurl) > 0 {
		return s.Baseurl
	}
	return `https://` + s.Bucket + `.` + s.Endpoint
}

func (s *CloudStorage) CachedList() map[string]*dbschema.NgingCloudStorage {
	items, ok := s.Context().Internal().Get(`NgingCloudStorages`).(map[string]*dbschema.NgingCloudStorage)
	if !ok {
		s.ListByOffset(nil, nil, 0, -1)
		list := s.Objects()
		items = s.KeyBy(`Id`, list)
		s.Context().Internal().Set(`NgingCloudStorages`, items)
	}
	return items
}

func (s *CloudStorage) Get(mw func(db.Result) db.Result, args ...interface{}) error {
	err := s.NgingCloudStorage.Get(mw, args...)
	if err != nil {
		return err
	}
	s.Secret = s.RawSecret()
	return nil
}

func (s *CloudStorage) Add() (pk interface{}, err error) {
	if err = s.check(); err != nil {
		return nil, err
	}
	return s.NgingCloudStorage.Insert()
}

func (s *CloudStorage) Edit(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	if err = s.check(); err != nil {
		return err
	}
	return s.NgingCloudStorage.Update(mw, args...)
}
