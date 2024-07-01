package system

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/penndev/wga/admin/bind"
	"github.com/penndev/wga/config"
	"github.com/penndev/wga/model/system"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"

	"github.com/penndev/gopkg/captcha"
	"golang.org/x/crypto/bcrypt"
)

func Captcha(c *gin.Context) {
	vd, err := captcha.NewImg()
	if err != nil {
		config.Logger.Warn("Captcha", zap.Error(err))
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "获取验证码出错"})
		return
	}
	c.JSON(http.StatusOK, bindCaptcha{
		ID:  vd.ID,
		Img: vd.PngBase64,
	})
}

func Login(c *gin.Context) {
	var request bindLoginInput
	c.ShouldBindJSON(&request)
	// 创建验证码
	if !captcha.Verify(request.CaptchaId, request.Captcha) {
		c.JSON(http.StatusForbidden, bind.ErrorMessage{Message: "验证码错误"})
		// return
	}
	res, err := system.GetSysUserByName(request.Username)
	if err != nil {
		var msg = "获取用户失败"
		if request.Username == "admin" && errors.Is(err, logger.ErrRecordNotFound) {
			msg = "已初始化管理员，请再次点击登录"
			res.Name = request.Username
			str, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.MinCost)
			if err != nil {
				config.Logger.Error("初始化管理员失败", zap.Error(err))
				msg = "初始化管理员失败，请查看错误日志"
			} else {
				res.Password = string(str)
				system.CreateSysUser(&res)
			}
		}
		c.JSON(http.StatusForbidden, bind.ErrorMessage{Message: msg})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(request.Password)) != nil {
		c.JSON(http.StatusForbidden, bind.ErrorMessage{Message: "用户密码错误"})
		return
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": request.Username,
		"exp": jwt.NewNumericDate(time.Now())}).SignedString([]byte(config.JWTSecret))
	if err != nil {
		config.Logger.Error("用户登录失败", zap.Error(err))
		c.JSON(http.StatusForbidden, bind.ErrorMessage{Message: "用户登录失败(jwt签名错误)"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// code, err := config.Redis.Get(context.Background(), "captcha:admin:"+request.CaptchaId).Result()
// if err != nil {
// 	c.JSON(http.StatusBadRequest, gin.H{
// 		"message": "验证码错误(0)",
// 	})
// 	return
// }
// config.Redis.Del(context.Background(), "captcha:admin:"+request.CaptchaId)
// if code != request.Captcha {
// 	c.JSON(http.StatusBadRequest, gin.H{
// 		"message": "验证码错误(1)",
// 	})
// 	return
// }
