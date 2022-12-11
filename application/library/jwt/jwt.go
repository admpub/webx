package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/webx-top/com"
)

// MakeSignedString 生成签名字符串
func MakeSignedString(customerID uint64, tokenPassword string) (tokenString string, err error) {
	claims := jwt.MapClaims{`customerId`: customerID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	tokenString, err = token.SignedString(com.Str2bytes(tokenPassword))
	return
}
