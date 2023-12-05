// token key
package global

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SecreteKey = []byte("ZF23%Diz46-7u_wW/Nbk20YjK")
var ActiveToken time.Duration = 7 * 24 * time.Hour

type CustomClaims struct {
	UserId   int    `json:"userId"`
	NickName string `json:"nickname"`
	jwt.RegisteredClaims
}
