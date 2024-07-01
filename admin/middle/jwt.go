package middle

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/penndev/wga/admin/bind"
)

// jwt验证用户登录
func JWTAuth(jwtSecret []byte) gin.HandlerFunc {
	// key func
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	}
	// auth func
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("x-token")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, bind.ErrorMessage{Message: "登录验证失败01"})
			c.Abort()
			return
		}
		// token 验证包含了 expired
		token, err := jwt.Parse(tokenStr, keyFunc)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, bind.ErrorMessage{Message: "登录验证失败02"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, bind.ErrorMessage{Message: "登录验证失败03"})
			c.Abort()
			return
		}

		sub, ok := claims["sub"]
		if !ok {
			c.JSON(http.StatusUnauthorized, bind.ErrorMessage{Message: "登录验证失败03"})
			c.Abort()
			return
		}

		c.Set("jwtAuth", sub)
		c.Next()
	}
}
