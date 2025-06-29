package ginhelper

import (
	"crypto/hmac"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type SignatureConfig struct {
	Key         []byte           // 签名密钥
	Hash        func() hash.Hash // HMAC算法，如 "sha256"
	SignName    string           // 签名参数名称
	ExpiredName string           // 过期时间参数名称
	OnError     func(c *gin.Context, status int, err error)
}

// 验证url签名签名步骤
// js客户端示例 https://gist.github.com/penndev/96ef7ddaba72e1eb42b09ce24f8ff734
// example:
//
//	var signMiddle = middle.Signature(middle.SignatureConfig{
//		Key:         []byte(os.Getenv("APP_SECRET")), // 签名密钥
//		Hash:        sha256.New,                      // HMAC算法
//		SignName:    "sign",                          // 签名参数名称
//		ExpiredName: "expired",                       // 过期时间参数名称
//		OnError:     onError,                         // 自定义异常处理结果
//	})
func Signature(cfg SignatureConfig) gin.HandlerFunc {
	// 验证参数是否合规
	if len(cfg.Key) == 0 {
		panic("SignatureConfig: Key must be set")
	}
	if cfg.Hash == nil {
		panic("SignatureConfig: Hash must be set")
	}
	if cfg.SignName == "" {
		panic("SignatureConfig: SignName must be set")
	}
	if cfg.ExpiredName == "" {
		panic("SignatureConfig: ExpiredName must be set")
	}
	if cfg.OnError == nil {
		panic("SignatureConfig: OnError must be set")
	}
	return func(c *gin.Context) {

		// 验证过期时间
		expired, err := strconv.Atoi(c.Query(cfg.ExpiredName))
		if err != nil {
			cfg.OnError(c, http.StatusForbidden, errors.New("expired time bad"))
			c.Abort()
			return
		}
		if expired <= 0 || int64(expired) < time.Now().Unix() {
			cfg.OnError(c, http.StatusForbidden, errors.New("expired time denied"))
			c.Abort()
			return
		}
		// c.Request.URL = /ping?param=value&expired=100&sign=XItiH_n47ZZJa7oKgWUvg-qMv_qFvoLCb9n8d63DHR4
		queryList := strings.Split(fmt.Sprint(c.Request.URL), "&"+cfg.SignName+"=")
		if len(queryList) != 2 {
			cfg.OnError(c, http.StatusForbidden, errors.New("url is bad"))
			c.Abort()
			return
		}
		// 验证url签名
		signMsg := queryList[0]
		signConn := queryList[1]
		hmacMethod := hmac.New(cfg.Hash, cfg.Key)
		_, err = hmacMethod.Write([]byte(signMsg))
		if err != nil {
			cfg.OnError(c, http.StatusForbidden, errors.New("hmac Write fail"))
			c.Abort()
			return
		}
		signGen := base64.RawURLEncoding.EncodeToString(hmacMethod.Sum(nil))
		if signGen != strings.TrimSpace(signConn) {
			cfg.OnError(c, http.StatusForbidden, errors.New("url sign denied"))
			c.Abort()
			return
		}
	}
}
