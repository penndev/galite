package system

import (
	"errors"
	"log"
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
		config.Logger.Error("Captcha", zap.Error(err))
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "获取验证码出错"})
		return
	}
	c.JSON(http.StatusOK, bindCaptcha{
		CaptchaID:  vd.ID,
		CaptchaURL: vd.PngBase64,
	})
}

func Login(c *gin.Context) {
	var request bindLoginInput
	if err := c.ShouldBindJSON(&request); err != nil {
		config.Logger.Warn("登录失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}

	// 创建验证码
	if !captcha.Verify(request.CaptchaId, request.Captcha) {
		c.JSON(http.StatusForbidden, bind.ErrorMessage{Message: "验证码错误"})
		return
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
		"exp": jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour))}).SignedString([]byte(config.JWTSecret))
	if err != nil {
		config.Logger.Error("用户登录失败", zap.Error(err))
		c.JSON(http.StatusForbidden, bind.ErrorMessage{Message: "用户登录失败(jwt签名错误)"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token":  token,
		"routes": "*", //前端菜单解决方案
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

func SystemUserList(c *gin.Context) {
	param := &bindSystemUserParam{}
	if err := c.BindQuery(&param); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	var total int64
	var list []system.SysUser

	m := param.Param() //处理筛选
	m.List(&total, &list)
	c.JSON(http.StatusOK, bind.DataList{Total: total, List: list})
}
