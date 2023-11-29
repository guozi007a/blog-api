// token key
package global

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SecreteKey = []byte("ZF23%Diz46-7u_wW/Nbk20YjK")
var ActiveCookie int = 7 * 24 * 60 * 60 // 有效时间，默认7天
var ActiveToken time.Duration = 7 * 24 * time.Hour

type CustomClaims struct {
	UserId   int    `json:"userId"`
	NickName string `json:"nickname"`
	jwt.RegisteredClaims
}
