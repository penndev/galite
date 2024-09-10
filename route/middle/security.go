package middle

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type message struct {
	Code    int    `json:"code"`
	Data    string `json:"data"`
	Message string `json:"message"`
}

type ginWrite struct {
	key [16]byte
	iv  [16]byte
	gin.ResponseWriter
}

func (w *ginWrite) Write(body []byte) (int, error) {
	plaintext := body
	key := w.key[:]
	iv := w.iv[:]
	block, err := aes.NewCipher(key)
	if err != nil {
		// config.Logger.Error("ginWrite", zap.Error(err))
		return 0, err
	}

	var paddedPlaintext []byte
	padding := block.BlockSize() - len(plaintext)%block.BlockSize()
	if padding == 0 {
		paddedPlaintext = []byte(plaintext)
	} else { //PKCS7填充
		padtext := bytes.Repeat([]byte{byte(padding)}, padding)
		paddedPlaintext = append([]byte(plaintext), padtext...)
	}
	blockMode := cipher.NewCBCEncrypter(block, iv)
	aesPlainText := make([]byte, len(paddedPlaintext))
	blockMode.CryptBlocks(aesPlainText, paddedPlaintext)

	m := message{
		Data: base64.StdEncoding.EncodeToString(aesPlainText),
	}

	jsonMsg, err := json.Marshal(m)
	if err != nil {
		// config.Logger.Error("ginWrite", zap.Error(err))
		return 0, err
	}
	return w.ResponseWriter.Write(jsonMsg)
}

// 验证url签名
// 验证body的签名
// 验证app key
// 接口加密
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {

		queryList := strings.Split(fmt.Sprintln(c.Request.URL), "&sign=")
		if len(queryList) != 2 {
			c.JSON(http.StatusForbidden, message{Message: "URL is Bad"})
			c.Abort()
			return
		}
		queryBase64 := make([]byte, base64.StdEncoding.EncodedLen(len(queryList[0])))
		base64.StdEncoding.Encode(queryBase64, []byte(queryList[0]))
		log.Print(string(queryBase64))
		sign := sha256.Sum256(queryBase64)
		signHex := hex.EncodeToString(sign[:])
		if signHex != strings.TrimSpace(queryList[1]) {
			c.JSON(http.StatusForbidden, message{Message: "Query Sign Denied"})
			c.Abort()
			return
		}

		// 验证app key
		appKey := c.GetHeader("X-App-Key")
		if appKey == "" {
			c.JSON(http.StatusForbidden, message{Message: "AppKey Bad"})
			c.Abort()
			return
		}
		// 真实验证 - 读取数据
		c.Set("appKey", appKey)

		// 验证http的query的t参数，判断是否是实时的请求

		// 验证http的query的s参数 http body的md5结果
		switch c.Request.Method {
		case http.MethodPut, http.MethodPost:
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, message{Message: "Body Sign Fail"})
				c.Abort()
				return
			}
			bodyMd5Sum := md5.Sum(bodyBytes)
			bodyMd5Hex := hex.EncodeToString(bodyMd5Sum[:])
			bodySign := c.Request.URL.Query().Get("s")

			if bodySign != bodyMd5Hex {
				c.JSON(http.StatusForbidden, message{Message: "Body Sign Denied"})
				c.Abort()
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 返回结果加密
		writer := &ginWrite{
			key:            md5.Sum([]byte(signHex)),
			iv:             md5.Sum([]byte(appKey)),
			ResponseWriter: c.Writer,
		}
		c.Writer = writer
	}
}
