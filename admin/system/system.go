package system

import (
	"context"
	"encoding/base64"
	"errors"
	"image/color"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/penndev/galite/admin/bind"
	"github.com/penndev/galite/cache"
	"github.com/penndev/galite/config"
	"github.com/penndev/galite/model/system"
	"github.com/penndev/gopkg/captcha"
	"github.com/penndev/gopkg/otp"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/logger"
)

func Captcha(c *gin.Context) {
	randText := captcha.RandText(4)
	buf, err := captcha.NewPngImg(captcha.Option{
		Width:     120,
		Height:    30,
		DPI:       90,
		Text:      randText,
		FontSize:  20,
		TextColor: color.RGBA{0, 0, 0, 255},
	})
	if err != nil {
		config.Logger.Error("Captcha", zap.Error(err))
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "获取验证码出错"})
		return
	}
	data := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
	id := uuid.New().String()
	cmd := cache.Redis.Set(context.TODO(), "captcha:"+id, randText, 5*time.Minute)
	if err := cmd.Err(); err != nil {
		config.Logger.Error("Redis错误", zap.Error(err))
	}
	c.JSON(http.StatusOK, bindCaptcha{
		CaptchaID:  id,
		CaptchaURL: data,
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
	captcha, err := cache.Redis.Get(context.Background(), "captcha:"+request.CaptchaId).Result()
	if err != nil {
		config.Logger.Warn("Redis错误", zap.Error(err))
		c.JSON(http.StatusForbidden, bind.ErrorMessage{Message: "验证码错误"})
		return
	}

	if !strings.EqualFold(captcha, request.Captcha) {
		c.JSON(http.StatusForbidden, bind.ErrorMessage{Message: "验证码错误"})
		return
	}
	res, err := system.SysAdminGetByEmail(request.Username)
	if err != nil {
		var msg = "获取用户失败"
		if request.Username == "admin" && errors.Is(err, logger.ErrRecordNotFound) {
			msg = "已初始化管理员，请再次点击登录"
			res.Email = request.Username
			bcryptPasswd, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
			if err != nil {
				config.Logger.Error("初始化管理员失败", zap.Error(err))
				msg = "初始化管理员失败，请查看错误日志"
			} else {
				res.Passwd = string(bcryptPasswd)
				res.Status = 1
				res.Nickname = "超级管理员"
				if err = res.Bind(res).Create(res).Error; err != nil {
					msg = "初始化管理员失败，请查看错误日志(1)"
				}
			}
		}
		c.JSON(http.StatusForbidden, bind.ErrorMessage{Message: msg})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(res.Passwd), []byte(request.Password)) != nil {
		c.JSON(http.StatusForbidden, bind.ErrorMessage{Message: "用户密码错误"})
		return
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(int(res.ID)),
		"exp": jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		"iat": jwt.NewNumericDate(time.Now()),
	}).SignedString([]byte(config.JWTSecret))
	if err != nil {
		config.Logger.Error("用户登录失败", zap.Error(err))
		c.JSON(http.StatusForbidden, bind.ErrorMessage{Message: "用户登录失败(jwt签名错误)"})
		return
	}

	result := gin.H{
		"token":    token,
		"routes":   res.SysRole.Menu, //前端菜单解决方案
		"nickname": res.Nickname,
	}
	if res.SysRoleID == nil || *res.SysRoleID < 1 {
		result["routes"] = "*" // 超级管理员 则替换为通配符
	}
	c.JSON(http.StatusOK, result)

}

func ChangePasswd(c *gin.Context) {
	var request bindChangePasswdInput
	if err := c.ShouldBindJSON(&request); err != nil {
		config.Logger.Warn("修改失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	res, err := system.SysAdminGetByID(c.GetString("jwtAuth"))
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "修改失败" + err.Error()})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(res.Passwd), []byte(request.Passwd)) != nil {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "旧的用户密码错误"})
		return
	}
	pwd, err := bcrypt.GenerateFromPassword([]byte(request.NewPasswd), bcrypt.DefaultCost)
	if err != nil {
		config.Logger.Error("创建管理员密码失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "创建密码失败"})
		return
	}
	res.Passwd = string(pwd)
	res.Bind(res).Updates(res)
	c.JSON(http.StatusOK, bind.ErrorMessage{Message: "修改完成"})
}

// 获取二次验证器密钥
func GetOTPSecret(c *gin.Context) {
	topic := c.Query("topic")
	if topic == "" {
		topic = "未知主题"
	}
	title := c.Query("title")
	if topic == "" {
		topic = "未知名称"
	}
	// 生成一个新的二次验证器
	secret, err := otp.GenerateSecret()
	if err != nil {
		config.Logger.Error("生成二次验证器失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, bind.ErrorMessage{Message: "生成二次验证器失败"})
		return
	}
	uri := otp.GenerateOTPURI("totp", topic, title, secret)
	c.JSON(http.StatusOK, gin.H{
		"secret": secret,
		"uri":    uri,
	})
}

// 验证二次验证器客户端是否正常
func VerifyOTPSecret(c *gin.Context) {
	var request bindVerifyOPTInput
	if err := c.ShouldBindJSON(&request); err != nil {
		config.Logger.Warn("二次验证失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误" + err.Error()})
		return
	}
	code, err := otp.GenerateOTPWithTime(request.Secret, time.Now())
	if code != request.Code {
		config.Logger.Warn("二次验证失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "二次验证失败"})
		return
	}

	c.JSON(http.StatusOK, bind.ErrorMessage{Message: "二次验证成功"})

}
