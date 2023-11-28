// token key
package global

import (
	"github.com/golang-jwt/jwt/v5"
)

var SecreteKey = []byte("ZF23%Diz46-7u_wW/Nbk20YjK")

type CustomClaims struct {
	UserId   int    `json:"userId"`
	NickName string `json:"nickname"`
	jwt.RegisteredClaims
}
