package plugins

import (
	"blog-api/db_server/tables"
	"blog-api/global"
	"net/http"
	"time"

	"slices"

	"strconv"

	"github.com/gin-gonic/gin"
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
		return ""
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
		return claims
	}
	if _claims, ok := token.Claims.(*global.CustomClaims); ok {
		// fmt.Printf("foo: %s\npt: %s\n", claims.Foo, claims.RegisteredClaims.ID)
		claims = _claims
	} else {
		return claims
	}
	return claims
}

// 验证token中间件 (只验证传递的token是否失效)
func VerifyTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 有些请求是不需要进行token验证的，如登录请求
		excludePath := []string{"/v3/login"}
		path := c.Request.URL.Path
		if slices.Contains[[]string, string](excludePath, path) {
			c.Next()
			return
		}

		userId := c.Request.Header.Get("userId")
		token := c.Request.Header.Get("token")
		if token == "" {
			c.Next()
			return
		}

		claims := ParseToken(token)
		t := claims.RegisteredClaims.ExpiresAt.UnixMilli()
		if t > time.Now().UnixMilli() {
			c.Next()
			return
		}

		uid, err := strconv.Atoi(userId)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeParamError,
				"message": "参数有误",
				"data":    nil,
			})
			return
		}

		db := global.GlobalDB
		result := db.Model(&tables.IdInfo{}).Select("isLogin, token").Where("userId = ?", uid).Updates(tables.IdInfo{IsLogin: false, Token: ""})
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeUpdateFailed,
				"message": "更新失败",
				"data":    nil,
			})
			return
		}

		// 给前端发送401码
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    global.CodeTokenInvalid,
			"message": "Token失效",
			"data":    nil,
		})
	}
}
