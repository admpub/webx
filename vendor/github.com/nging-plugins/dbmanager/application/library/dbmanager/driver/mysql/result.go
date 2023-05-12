/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

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

package mysql

import (
	"database/sql"
	"encoding/gob"
	"html"
	"html/template"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/nging-plugins/dbmanager/application/library/dbmanager/result"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo/param"
)

func init() {
	gob.Register(&Result{})
}

var _ result.Resulter = &Result{}
var resultPool = &sync.Pool{
	New: func() interface{} {
		return &Result{}
	},
}

func AcquireResult() *Result {
	return resultPool.Get().(*Result)
}

func ReleaseResult(r *Result) {
	r.Reset()
	resultPool.Put(r)
}

type Result struct {
	SQL          string
	SQLs         []string
	RowsAffected int64
	TimeStart    time.Time
	TimeEnd      time.Time
	Started      string
	Elapsed      string
	err          error
	ErrorString  string
}

func (r *Result) Reset() {
	r.SQL = ``
	r.SQLs = []string{}
	r.RowsAffected = 0
	r.TimeStart = param.EmptyTime
	r.TimeEnd = param.EmptyTime
	r.Started = ``
	r.Elapsed = ``
	r.err = nil
	r.ErrorString = ``
}

func (r *Result) GetSQL() string {
	return r.SQL
}

func (r *Result) GetSQLs() []string {
	if r.SQLs == nil || len(r.SQLs) == 0 {
		return []string{r.SQL}
	}
	return r.SQLs
}

func (r *Result) GetBeginTime() string {
	return r.Started
}

var tableNameRegex = regexp.MustCompile("(?i)((?:[\\s]+INTO|[\\s]+FROM|[\\s]*UPDATE)[\\s]+`)([^\\s]+)(`[\\s]+)")

func (r *Result) ToHTML(urlPrefix string, sql string) template.HTML {
	sql = html.EscapeString(sql)
	sql = tableNameRegex.ReplaceAllString(sql, `$1<a href="`+urlPrefix+`&table=$2">$2</a>$3`)
	return template.HTML(sql)
}

func (r *Result) GetElapsedTime() string {
	return r.Elapsed
}
func (r *Result) GetAffected() int64 {
	return r.RowsAffected
}

func (r *Result) GetError() string {
	return r.ErrorString
}

func (r *Result) Error() error {
	return r.err
}

func (r *Result) SetError(err error) {
	r.err = err
}

func (r *Result) elapsed() time.Duration {
	return r.TimeEnd.Sub(r.TimeStart)
}

func (r *Result) start() *Result {
	r.TimeStart = time.Now()
	return r
}

func (r *Result) end() *Result {
	r.TimeEnd = time.Now()
	r.Started = r.TimeStart.Format(`2006-01-02 15:04:05`)
	r.Elapsed = r.elapsed().String()
	if r.err != nil {
		r.ErrorString = r.err.Error()
	}
	return r
}

func (r *Result) Exec(p *factory.Param) *Result {
	r.start()
	defer r.end()
	result, err := p.SetCollection(r.SQL).Exec()
	r.err = err
	if err != nil {
		return r
	}
	r.RowsAffected, r.err = result.RowsAffected()
	return r
}

func (r *Result) Query(p *factory.Param, readRows func(*sql.Rows) error) *Result {
	r.start()
	defer r.end()
	rows, err := p.SetCollection(r.SQL).Query()
	r.err = err
	if err != nil {
		return r
	}
	r.err = readRows(rows)
	rows.Close()
	return r
}

func (r *Result) QueryRow(p *factory.Param, recvs ...interface{}) *Result {
	r.start()
	defer r.end()
	row, err := p.SetCollection(r.SQL).QueryRow()
	if err != nil {
		r.err = err
		return r
	}
	r.err = row.Scan(recvs...)
	return r
}

func (r *Result) Execs(p *factory.Param) *Result {
	if r.TimeStart.IsZero() {
		r.start()
	}
	defer r.end()
	if r.SQL == `` {
		return r
	}
	if r.SQLs == nil {
		r.SQLs = []string{}
	}
	if strings.HasSuffix(r.SQL, `;`) {
		r.SQLs = append(r.SQLs, "DELIMITER ;;\n"+r.SQL+";\nDELIMITER ;")
	} else {
		r.SQLs = append(r.SQLs, r.SQL+";")
	}
	result, err := p.DB().Exec(r.SQL)
	r.err = err
	if err != nil {
		return r
	}
	r.RowsAffected, r.err = result.RowsAffected()
	return r
}

func (r *Result) Queries(p *factory.Param, readRows func(*sql.Rows) error) *Result {
	if r.TimeStart.IsZero() {
		r.start()
	}
	defer r.end()
	if r.SQL == `` {
		return r
	}
	if r.SQLs == nil {
		r.SQLs = []string{}
	}
	if strings.HasSuffix(r.SQL, `;`) {
		r.SQLs = append(r.SQLs, "DELIMITER ;;\n"+r.SQL+";\nDELIMITER ;")
	} else {
		r.SQLs = append(r.SQLs, r.SQL+";")
	}
	rows, err := p.SetCollection(r.SQL).Query()
	r.err = err
	if err != nil {
		return r
	}
	r.err = readRows(rows)
	rows.Close()
	return r
}
