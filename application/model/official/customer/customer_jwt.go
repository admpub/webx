package customer

import (
	"errors"
	"fmt"
	"time"

	"github.com/admpub/nging/v5/application/cmd/bootconfig"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/webx/application/dbschema"
	"github.com/golang-jwt/jwt/v4"
	"github.com/webx-top/db"
	mwJWT "github.com/webx-top/echo/middleware/jwt"
	"github.com/webx-top/echo/param"
)

// JWTMaxLifeTime JWT寿命(单位:秒)
var JWTMaxLifeTime int64 = 90 * 86400

func (f *Customer) JWTClaims(customers ...*dbschema.OfficialCustomer) *jwt.StandardClaims {
	var customer *dbschema.OfficialCustomer
	if len(customers) > 0 {
		customer = customers[0]
	} else {
		customer = f.OfficialCustomer
	}
	nowTS := time.Now().Unix()
	lifetime := config.Setting(`base`).Int64(`jwtMaxLifetime`)
	if lifetime <= 0 {
		lifetime = JWTMaxLifeTime
	}
	endTS := nowTS + lifetime
	return &jwt.StandardClaims{
		Audience:  f.Context().Session().MustID(),
		ExpiresAt: endTS,
		Id:        param.AsString(customer.Id),
		IssuedAt:  nowTS,
		Issuer:    bootconfig.SoftwareName,
		NotBefore: nowTS,
		Subject:   customer.Name,
	}
}

func (f *Customer) JWTSignedString(key interface{}, customers ...*dbschema.OfficialCustomer) (string, error) {
	claims := f.JWTClaims(customers...)
	if key == nil {
		key = []byte(config.FromFile().Cookie.HashKey)
	}
	return mwJWT.BuildStandardSignedString(claims, key)
}

var ErrInvalidSession = errors.New(`invalid session`)

func (f *Customer) GetByJWT() (*dbschema.OfficialCustomer, error) {
	//echo.Dump(f.Context().Internal().Get(`jwtUser`))
	token, ok := f.Context().Internal().Get(`jwtUser`).(*jwt.Token)
	if !ok {
		return nil, nil
	}
	//echo.Dump(token.Claims)
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, nil
	}
	customerID := param.AsUint64(claims.Id)
	if customerID == 0 {
		return nil, nil
	}
	err := f.Get(nil, `id`, customerID)
	if err != nil {
		if db.ErrNoMoreRows == err {
			return nil, nil
		}
		return nil, fmt.Errorf(`GetByJWT: %w`, err)
	}
	if len(claims.Audience) > 0 {
		if f.SessionId != claims.Audience {
			return nil, ErrInvalidSession
		}
		if err := f.Context().Session().SetID(claims.Audience); err != nil {
			return nil, err
		}
		sid := f.Context().GetCookie(f.Context().SessionOptions().Name)
		if len(sid) == 0 {
			f.Context().SetCookie(f.Context().SessionOptions().Name, claims.Audience)
		}
	}
	customer := f.ClearPasswordData()
	return &customer, err
}
