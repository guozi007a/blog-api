package plugins

import (
	"blog-api/global"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 生成token
func CreateToken(userId int, nickname string) string {
	claims := global.CustomClaims{
		UserId:   userId,
		NickName: nickname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(global.ActiveToken)), // 过期时间，必须在nbf之后
			IssuedAt:  jwt.NewNumericDate(time.Now()),                         // 表示创建token的时间
			NotBefore: jwt.NewNumericDate(time.Now()),                         // 如果生成时间或者有效期时间在该字段时间之前，则token无效
			Issuer:    "login token",                                          // 表示谁创建的这个token，一般使用项目名或者服务器域名等
			Subject:   "activity",                                             // 表示该token对应的主题。可以包含任何内容，但应该能够唯一地标识一个实体。在某些情况下可能与ID字段相同。
			ID:        "activity",                                             // 该token的唯一标识符
			Audience:  []string{"a", "b", "c"},                                // 给定一个接收者验证token是否为其本人使用
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	sign, err := token.SignedString(global.SecreteKey)
	if err != nil {
		panic(err)
	}
	return sign
}

// 解析token
func ParseToken(tokenStr string) *global.CustomClaims {
	claims := &global.CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return global.SecreteKey, nil
	})
	if err != nil {
		panic(err)
	} else if claims, ok := token.Claims.(*global.CustomClaims); ok {
		// fmt.Printf("foo: %s\npt: %s\n", claims.Foo, claims.RegisteredClaims.ID)
		return claims
	} else {
		fmt.Println("unknown claims type, cannot proceed")
	}
	return claims
}
