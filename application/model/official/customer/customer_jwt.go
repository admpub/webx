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

func (f *Customer) JWTClaims(customers ...*dbschema.OfficialCustomer) *jwt.RegisteredClaims {
	var customer *dbschema.OfficialCustomer
	if len(customers) > 0 {
		customer = customers[0]
	} else {
		customer = f.OfficialCustomer
	}
	now := time.Now()
	lifetime := config.Setting(`base`).Int64(`jwtMaxLifetime`)
	if lifetime <= 0 {
		lifetime = JWTMaxLifeTime
	}
	expires := now.Add(time.Second * time.Duration(lifetime))
	return &jwt.RegisteredClaims{
		Audience:  jwt.ClaimStrings{f.Context().Session().MustID()},
		ExpiresAt: &jwt.NumericDate{Time: expires},
		ID:        param.AsString(customer.Id),
		IssuedAt:  &jwt.NumericDate{Time: now},
		Issuer:    bootconfig.SoftwareName,
		NotBefore: &jwt.NumericDate{Time: now},
		Subject:   customer.Name,
	}
}

func (f *Customer) JWTSignedString(key interface{}, customers ...*dbschema.OfficialCustomer) (string, error) {
	claims := f.JWTClaims(customers...)
	if key == nil {
		key = []byte(config.FromFile().Cookie.HashKey)
	}
	return mwJWT.BuildRegisteredSignedString(claims, key)
}

var ErrInvalidSession = errors.New(`invalid session`)

func (f *Customer) GetByJWT() (*dbschema.OfficialCustomer, error) {
	//echo.Dump(f.Context().Internal().Get(`jwtUser`))
	token, ok := f.Context().Internal().Get(`jwtUser`).(*jwt.Token)
	if !ok {
		return nil, nil
	}
	//echo.Dump(token.Claims)
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, nil
	}
	customerID := param.AsUint64(claims.ID)
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
		var found bool
		for _, v := range claims.Audience {
			if f.SessionId == v {
				found = true
				break
			}
		}
		if !found {
			return nil, ErrInvalidSession
		}
		if err := f.Context().Session().SetID(f.SessionId); err != nil {
			return nil, err
		}
		sid := f.Context().GetCookie(f.Context().SessionOptions().Name)
		if len(sid) == 0 {
			f.Context().SetCookie(f.Context().SessionOptions().Name, f.SessionId)
		}
	}
	customer := f.ClearPasswordData()
	return &customer, err
}
