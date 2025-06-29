package ginhelper

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"github.com/gin-gonic/gin"
)

// 对gin的write进行封装
type encryptResponseWriter struct {
	OnWrite func([]byte) []byte
	gin.ResponseWriter
}

func (encrypt *encryptResponseWriter) Write(body []byte) (int, error) {
	cipher := encrypt.OnWrite(body)
	return encrypt.ResponseWriter.Write(cipher)
}

type EncryptionConfig struct {
	Secret          string
	IvName          string
	ContentTypeName string
	OnError         func(c *gin.Context, status int, err error)
}

// 加解密中间件
// js客户端示例 https://gist.github.com/penndev/96ef7ddaba72e1eb42b09ce24f8ff734
// example:
//
//	var encryptMiddle = middle.Encryption(middle.EncryptionConfig{
//		Secret:          os.Getenv("APP_SECRET"),
//		IvName:          "X-Iv",
//		ContentTypeName: "application/x-buffer",
//		OnError:         onError, // 自定义异常处理结果
//	})
func Encryption(cfg EncryptionConfig) gin.HandlerFunc {
	// 加解密需要固定的key长度，所以必须处理key
	shaHash := sha256.New()
	_, err := shaHash.Write(([]byte(cfg.Secret)))
	if err != nil {
		panic(err)
	}
	key := shaHash.Sum(nil)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		// 对请求体进行解密处理
		if c.GetHeader("Content-Type") == cfg.ContentTypeName {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				cfg.OnError(c, 500, err)
				c.Abort()
			}
			nonce, err := base64.RawURLEncoding.DecodeString(c.GetHeader(cfg.IvName))
			if err != nil {
				cfg.OnError(c, 403, err)
				c.Abort()
			}
			if len(nonce) != aesGCM.NonceSize() {
				cfg.OnError(c, 403, errors.New("nonce length error"))
				c.Abort()
			}
			plaintext, err := aesGCM.Open(nil, nonce, bodyBytes, nil)
			if err != nil {
				cfg.OnError(c, 403, err)
				c.Abort()
			}
			c.Request.Body = io.NopCloser(bytes.NewReader(plaintext))
		}
		// 对响应体进行加密处理
		c.Writer = &encryptResponseWriter{
			OnWrite: func(body []byte) []byte {
				nonce := make([]byte, aesGCM.NonceSize())
				if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
					// 处理报错
				}
				c.Header("Content-Type", cfg.ContentTypeName)
				c.Header(cfg.IvName, base64.RawURLEncoding.EncodeToString(nonce))
				return aesGCM.Seal(nil, nonce, body, nil)
			},
			ResponseWriter: c.Writer,
		}
	}
}
