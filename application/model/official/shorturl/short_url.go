package shorturl

import (
	"net/url"
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/errors"
	shortid "github.com/admpub/go-shortid"
	"github.com/admpub/nging/v5/application/library/ip2region"
	ua "github.com/admpub/useragent"
	"github.com/admpub/webx/application/dbschema"
)

var (
	ShortURLGenerator   *shortid.Gen
	ErrLongURLExists    = errors.New(`网址已经存在`)
	ErrNeedURLPassword  = errors.New(`请输入密码`)
	ErrWrongURLPassword = errors.New(`密码错误`)
	ErrExpiredURL       = errors.New(`短网址已过期`)
	ErrNotExistsURL     = errors.New(`短网址不存在`)
)

func init() {
	ShortURLGenerator = shortid.Generator()
}

func NewShortURL(ctx echo.Context) *ShortURL {
	return &ShortURL{
		OfficialShortUrl: dbschema.NewOfficialShortUrl(ctx),
		Domain:           dbschema.NewOfficialShortUrlDomain(ctx),
		Visit:            dbschema.NewOfficialShortUrlVisit(ctx),
	}
}

type ShortURL struct {
	*dbschema.OfficialShortUrl
	urlInfo *url.URL
	Domain  *dbschema.OfficialShortUrlDomain
	Visit   *dbschema.OfficialShortUrlVisit
}

func (f *ShortURL) check() error {
	f.LongUrl = strings.TrimSpace(f.LongUrl)
	if len(f.LongUrl) == 0 {
		return f.Context().E(`请输入网址`)
	}
	pos := strings.Index(f.LongUrl, `://`)
	if pos <= 0 {
		return f.Context().E(`网址无效: %s`, f.LongUrl)
	}
	switch strings.ToLower(f.LongUrl[0:pos]) {
	case `http`, `https`:
	default:
		return errors.New(f.Context().T(`网址无效: %s`, f.LongUrl) + `, ` + f.Context().T(`必须以“http://”或“https://”开头`))
	}
	f.LongHash = com.Md5(f.LongUrl)
	return nil
}

func (f *ShortURL) ParseURL(urls ...string) (urlInfo *url.URL, err error) {
	rurl := f.LongUrl
	if len(urls) > 0 {
		rurl = urls[0]
	} else {
		if f.urlInfo != nil {
			urlInfo = f.urlInfo
			return
		}
	}
	urlInfo, err = url.Parse(rurl)
	if err != nil {
		return
	}
	if len(urlInfo.Host) == 0 {
		err = f.Context().E(`无效网址: %v`, rurl)
		return
	}
	urlInfo.Host = urlInfo.Hostname() //不含端口
	if len(urls) == 0 {
		f.urlInfo = urlInfo
	}
	return
}

func (f *ShortURL) Add() (pk interface{}, err error) {
	err = f.Context().Begin()
	if err != nil {
		return
	}
	defer func() {
		f.Context().End(err == nil)
	}()
	err = f.check()
	if err != nil {
		return
	}
	err = f.ExistsHash(f.LongHash)
	if err != nil {
		if err != ErrLongURLExists {
			return
		}
		err = f.Context().E(`网址已存在: %s`, f.LongUrl)
		return
	}
	var urlInfo *url.URL
	urlInfo, err = f.ParseURL()
	if err != nil {
		return
	}
	err = f.Domain.Get(nil, db.Cond{`domain`: urlInfo.Host})
	if err != nil {
		if err != db.ErrNoMoreRows {
			return
		}
		f.Domain.Reset()
		f.Domain.Domain = urlInfo.Host
		f.Domain.OwnerId = f.OwnerId
		f.Domain.OwnerType = f.OwnerType
		f.Domain.UrlCount = 1
		_, err = f.Domain.Insert()
	} else {
		if f.Domain.Disabled == `Y` {
			err = f.Context().E(`操作失败！域名“%v”已经被加入黑名单。如有疑问，请联系我们`, f.Domain.Domain)
			return
		}
		err = f.Domain.UpdateField(nil, `url_count`, db.Raw(`url_count+1`), `id`, f.Domain.Id)
	}
	if err != nil {
		return
	}
	f.DomainId = f.Domain.Id
	f.Visited = 0
	f.Visits = 0
	if len(f.ShortUrl) == 0 {
		f.ShortUrl = ShortURLGenerator.Generate()
	}
	err = f.Exists(f.ShortUrl)
	if err != nil {
		return
	}
	if len(f.Password) > 0 {
		f.Password = com.Md5(f.Password)
	}
	pk, err = f.Insert()
	return
}

func (f *ShortURL) Edit(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	err = f.Context().Begin()
	if err != nil {
		return
	}
	defer func() {
		f.Context().End(err == nil)
	}()
	err = f.check()
	if err != nil {
		return
	}
	var urlInfo *url.URL
	urlInfo, err = f.ParseURL()
	if err != nil {
		return
	}
	err = f.ExistsHashOther(f.LongHash, f.Id)
	if err != nil {
		return
	}
	oldURL := dbschema.NewOfficialShortUrl(nil)
	oldURL.CPAFrom(f.OfficialShortUrl)
	err = oldURL.Get(nil, args...)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return f.Context().E(`不存在ID为“%d”的数据`, f.Id)
		}
		return
	}
	oldDomain := dbschema.NewOfficialShortUrlDomain(nil)
	oldDomain.CPAFrom(f.OfficialShortUrl)
	err = oldDomain.Get(nil, db.Cond{`id`: oldURL.DomainId})
	var isNewDomain bool
	if err != nil {
		if err != db.ErrNoMoreRows {
			return
		}
		err = nil
		isNewDomain = true
	} else {
		if oldDomain.Domain != urlInfo.Host {
			if oldDomain.UrlCount < 2 {
				err = oldDomain.Delete(nil, db.Cond{`id`: oldURL.DomainId})
			} else {
				err = oldDomain.UpdateField(nil, `url_count`, db.Raw(`url_count-1`), db.Cond{`id`: oldURL.DomainId})
			}
			if err != nil {
				return
			}
			isNewDomain = true
		}
	}
	if isNewDomain {
		err = f.Domain.Get(nil, db.Cond{`domain`: urlInfo.Host})
		if err != nil {
			if err != db.ErrNoMoreRows {
				return
			}
			f.Domain.Reset()
			f.Domain.Domain = urlInfo.Host
			f.Domain.OwnerId = f.OwnerId
			f.Domain.OwnerType = f.OwnerType
			f.Domain.UrlCount = 1
			_, err = f.Domain.Insert()
		} else {
			if f.Domain.Disabled == `Y` {
				err = f.Context().E(`操作失败！域名“%v”已经被加入黑名单。如有疑问，请联系我们`, f.Domain.Domain)
				return
			}
			err = f.Domain.UpdateField(nil, `url_count`, db.Raw(`url_count+1`), `id`, f.Domain.Id)
		}
		if err != nil {
			return
		}
	}
	err = f.ExistsOther(f.ShortUrl, f.Id)
	if err != nil {
		return
	}
	modifyPassword := f.Context().Formx(`modifyPassword`).Bool()
	if modifyPassword {
		if len(f.Password) > 0 {
			f.Password = com.Md5(f.Password)
		} else {
			f.Password = ``
		}
	} else {
		f.Password = oldURL.Password
	}
	err = f.Update(mw, args...)
	return
}

func (f *ShortURL) Exists(shortURL string) error {
	exists, err := f.OfficialShortUrl.Exists(nil, db.Cond{`short_url`: shortURL})
	if err != nil {
		return err
	}
	if exists {
		err = f.Context().E(`短网址名称“%s”已经存在`, shortURL)
	}
	return err
}

func (f *ShortURL) ExistsHash(longHash string) error {
	exists, err := f.OfficialShortUrl.Exists(nil, db.Cond{`long_hash`: longHash})
	if err != nil {
		return err
	}
	if exists {
		return ErrLongURLExists
	}
	return err
}

func (f *ShortURL) ExistsHashOther(longHash string, id uint64) error {
	exists, err := f.OfficialShortUrl.Exists(nil, db.Cond{`long_hash`: longHash}, db.Cond{`id`: db.NotEq(id)})
	if err != nil {
		return err
	}
	if exists {
		return ErrLongURLExists
	}
	return err
}

func (f *ShortURL) ExistsOther(shortURL string, id uint64) error {
	exists, err := f.OfficialShortUrl.Exists(nil, db.And(
		db.Cond{`short_url`: shortURL},
		db.Cond{`id`: db.NotEq(id)},
	))
	if err != nil {
		return err
	}
	if exists {
		err = f.Context().E(`短网址名称“%s”已经使用过了`, shortURL)
	}
	return err
}

func (f *ShortURL) Find(shortURL string) (longURL string, err error) {
	err = f.Get(nil, db.And(
		db.Cond{`short_url`: shortURL},
	))
	if err != nil {
		if err != db.ErrNoMoreRows {
			return
		}
		err = ErrNotExistsURL
		return
	}
	t := time.Now()
	ts := uint(t.Unix())
	if f.Expired > 0 && f.Expired < ts {
		err = errors.WithMessage(ErrExpiredURL, f.Context().T(`网址已经于“%s”过期`, time.Unix(int64(f.Expired), 0).Format(`2006-01-02 15:04:05`)))
		return
	}
	if len(f.Password) > 0 {
		inputPwd := f.Context().Formx(`password`).String()
		if len(inputPwd) == 0 {
			err = ErrNeedURLPassword
			return
		}
		if f.Password != com.Md5(inputPwd) {
			err = ErrWrongURLPassword
			return
		}
	}
	longURL = f.LongUrl
	err = f.VisitAdd(t)
	if err != nil {
		return
	}
	err = f.UpdateFields(nil, echo.H{
		`visited`: ts,
		`visits`:  db.Raw(`visits+1`),
	}, `id`, f.Id)
	return
}

func (f *ShortURL) GenVisitData(t time.Time) error {
	f.Visit.UrlId = f.Id
	f.Visit.DomainId = f.DomainId
	f.Visit.OwnerId = f.OwnerId
	f.Visit.OwnerType = f.OwnerType
	f.Visit.Year = uint(t.Year())
	f.Visit.Month = uint(t.Month())
	f.Visit.Day = uint(t.Day())
	f.Visit.Hour = uint(t.Hour())
	f.Visit.Ip = f.Context().RealIP()
	acceptLanguage := f.Context().Header(`Accept-Language`)
	pos := strings.Index(acceptLanguage, `,`)
	if pos > 0 {
		acceptLanguage = strings.TrimSpace(acceptLanguage[0:pos])
	}
	acceptLanguage = strings.TrimSpace(strings.SplitN(acceptLanguage, `;`, 2)[0])
	f.Visit.Language = acceptLanguage
	f.Visit.Referer = f.Context().Referer()
	info, err := ip2region.IPInfo(f.Visit.Ip)
	if err != nil {
		if !ip2region.ErrIsInvalidIP(err) {
			return err
		}
	}
	if info.ISP == `0` {
		info.ISP = ``
	}
	if info.City == `0` {
		info.City = ``
	}
	if info.Province == `0` {
		info.Province = ``
	}
	if info.Region == `0` {
		info.Region = ``
	}
	if info.Country == `0` {
		info.Country = ``
	}
	f.Visit.Country = info.Country
	f.Visit.Region = info.Region
	f.Visit.Province = info.Province
	f.Visit.City = info.City
	f.Visit.Isp = info.ISP
	userAgent := f.Context().Request().UserAgent()
	infoUA := ua.Parse(userAgent)
	f.Visit.Os = infoUA.OS
	f.Visit.OsVersion = infoUA.OSVersion
	f.Visit.Browser = infoUA.Name
	f.Visit.BrowserVersion = infoUA.Version
	return nil
}

func (f *ShortURL) VisitAdd(t time.Time) error {
	if err := f.GenVisitData(t); err != nil {
		return err
	}
	_, err := f.Visit.Insert()
	return err
}

func (f *ShortURL) VisitListWithURL(rows []*dbschema.OfficialShortUrlVisit) ([]*ShortURLVisitWithURL, error) {
	lastList := make([]*ShortURLVisitWithURL, len(rows))
	urlIDs := []uint64{}
	for index, row := range rows {
		lastList[index] = &ShortURLVisitWithURL{
			OfficialShortUrlVisit: row,
			URL:                   dbschema.NewOfficialShortUrl(f.Context()),
		}
		if !com.InUint64Slice(row.UrlId, urlIDs) {
			urlIDs = append(urlIDs, row.UrlId)
		}
	}
	if len(urlIDs) > 0 {
		_, err := f.ListByOffset(nil, nil, 0, -1, db.Cond{`id`: db.In(urlIDs)})
		if err != nil {
			return lastList, err
		}
		for _, urlRow := range f.Objects() {
			for index, row := range lastList {
				if row.UrlId == urlRow.Id {
					lastList[index].URL = urlRow
				}
			}
		}
	}
	return lastList, nil
}

func (f *ShortURL) VisitListFillData(rows []*ShortURLVisitWithURL) ([]*ShortURLVisitWithURL, error) {
	urlIDs := []uint64{}
	for _, row := range rows {
		if !com.InUint64Slice(row.UrlId, urlIDs) {
			urlIDs = append(urlIDs, row.UrlId)
		}
	}
	if len(urlIDs) > 0 {
		_, err := f.ListByOffset(nil, nil, 0, -1, db.Cond{`id`: db.In(urlIDs)})
		if err != nil {
			return rows, err
		}
		for _, urlRow := range f.Objects() {
			for index, row := range rows {
				if row.UrlId == urlRow.Id {
					rows[index].URL = urlRow
				}
			}
		}
	}
	return rows, nil
}
