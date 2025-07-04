package middle

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/penndev/galite/internal/admin/bind"
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
	// auth
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("x-token")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, bind.Message{Message: "登录验证失败01"})
			c.Abort()
			return
		}
		token, err := jwt.Parse(tokenStr, keyFunc)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, bind.Message{Message: "登录验证失败02"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, bind.Message{Message: "登录验证失败03"})
			c.Abort()
			return
		}
		sub, err := claims.GetSubject()
		if err != nil {
			c.JSON(http.StatusUnauthorized, bind.Message{Message: "登录验证失败04"})
			c.Abort()
			return
		}
		adminID, err := strconv.Atoi(sub) // 验证 sub 是否为数字
		if err != nil {
			c.JSON(http.StatusUnauthorized, bind.Message{Message: "登录验证失败05"})
			c.Abort()
			return
		}
		c.Set("adminID", adminID)
	}
}
