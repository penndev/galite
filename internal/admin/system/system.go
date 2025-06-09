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
	"github.com/penndev/galite/internal/admin/bind"
	"github.com/penndev/galite/internal/admin/model/system"
	"github.com/penndev/galite/internal/cache"
	"github.com/penndev/galite/internal/config"
	"github.com/penndev/galite/internal/logger"
	"github.com/penndev/gopkg/captcha"
	"github.com/penndev/gopkg/otp"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	loggerGorm "gorm.io/gorm/logger"
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
		logger.ZapLogger.Error("Captcha", zap.Error(err))
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "获取验证码出错"})
		return
	}
	data := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
	id := uuid.New().String()
	cmd := cache.Redis.Set(context.TODO(), "captcha:"+id, randText, 5*time.Minute)
	if err := cmd.Err(); err != nil {
		logger.ZapLogger.Error("Redis错误", zap.Error(err))
	}
	c.JSON(http.StatusOK, bindCaptcha{
		CaptchaID:  id,
		CaptchaURL: data,
	})
}

// 登录成功根据 sysadmin 返回 jwt token 数据
func loginInfo(res *system.SysAdmin) (map[string]any, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(int(res.ID)),
		"exp": jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		"iat": jwt.NewNumericDate(time.Now()),
	}).SignedString([]byte(config.Secret()))
	if err != nil {
		return gin.H{}, err
	}
	result := gin.H{
		"token":     token,
		"routes":    res.SysRole.Menu, //前端菜单解决方案
		"nickname":  res.Nickname,
		"otpStatus": res.OtpStatus,
	}
	if res.SysRoleID == nil || *res.SysRoleID < 1 {
		result["routes"] = "*" // 超级管理员 则替换为通配符
	}
	return result, nil
}

func Login(c *gin.Context) {
	var request bindLoginInput
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.ZapLogger.Warn("登录失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误" + err.Error()})
		return
	}

	// 创建验证码
	captcha, err := cache.Redis.Get(context.Background(), "captcha:"+request.CaptchaId).Result()
	if err != nil {
		logger.ZapLogger.Warn("Redis错误", zap.Error(err))
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
		if request.Username == "admin" && errors.Is(err, loggerGorm.ErrRecordNotFound) {
			msg = "已初始化管理员，请再次点击登录"
			res.Email = request.Username
			bcryptPasswd, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
			if err != nil {
				logger.ZapLogger.Error("初始化管理员失败", zap.Error(err))
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

	// 登录需要二次验证时，先将用户信息存入 Redis，等待 OTP 验证
	if res.OtpStatus == 1 {
		key := "otp:login:" + strconv.Itoa(int(res.ID))
		data, err := cache.Encode(res)
		if err != nil {
			logger.ZapLogger.Error("cache.Encode", zap.Error(err))
			c.JSON(http.StatusInternalServerError, bind.ErrorMessage{Message: "编码失败"})
			return
		}
		cmd := cache.Redis.Set(context.TODO(), key, data, 5*time.Minute)
		if err := cmd.Err(); err != nil {
			logger.ZapLogger.Error("Redis错误", zap.Error(err))
			c.JSON(http.StatusInternalServerError, bind.ErrorMessage{Message: "Redis错误"})
			return
		}
		// 返回
		c.JSON(http.StatusOK, gin.H{
			"otpStatus": res.OtpStatus,
			"id":        res.ID,
			"otpTitle":  res.OtpTitle,
		})
		return
	}
	result, err := loginInfo(res)
	if err != nil {
		logger.ZapLogger.Error("用户登录失败", zap.Error(err))
		c.JSON(http.StatusForbidden, bind.ErrorMessage{Message: "用户登录失败(jwt签名错误)"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// 两步验证登录
func LoginOTP(c *gin.Context) {
	var request struct {
		ID   int    `form:"id" binding:"required"`         // 用户ID
		Code string `form:"code" binding:"required,len=6"` // 验证码
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.ZapLogger.Warn("二步验证登录失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	key := "otp:login:" + strconv.Itoa(request.ID)
	data, err := cache.Redis.Get(context.TODO(), key).Bytes()
	if err != nil {
		logger.ZapLogger.Warn("Redis错误", zap.Error(err))
		c.JSON(http.StatusForbidden, bind.ErrorMessage{Message: "登录信息已过期，请重新登录"})
		return
	}
	var res system.SysAdmin
	if err := cache.Decode(string(data), &res); err != nil {
		logger.ZapLogger.Error("cache.Decode", zap.Error(err))
		c.JSON(http.StatusInternalServerError, bind.ErrorMessage{Message: "解码失败"})
		return
	}
	// 校验 OTP
	code, err := otp.GenerateOTPWithTime(res.OtpSecret, time.Now())
	if err != nil {
		c.JSON(http.StatusForbidden, bind.ErrorMessage{Message: "二步验证失败"})
		return
	}
	if code != request.Code {
		c.JSON(http.StatusForbidden, bind.ErrorMessage{Message: "二步验证码错误"})
		return
	}

	result, err := loginInfo(&res)
	if err != nil {
		logger.ZapLogger.Error("用户登录失败", zap.Error(err))
		c.JSON(http.StatusForbidden, bind.ErrorMessage{Message: "用户登录失败(jwt签名错误)"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// 用户修改密码
func ChangePasswd(c *gin.Context) {
	var request bindChangePasswdInput
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.ZapLogger.Warn("修改失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	res, err := system.SysAdminGetByID(c.GetInt("adminID"))
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
		logger.ZapLogger.Error("创建管理员密码失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "创建密码失败"})
		return
	}
	res.Passwd = string(pwd)
	res.Bind(res).Updates(res)
	c.JSON(http.StatusOK, bind.ErrorMessage{Message: "修改完成"})
}

// 用户重置自己的OTP验证器
func ChangeOTP(c *gin.Context) {
	var request struct {
		OtpTitle  string `form:"otpTitle"`  // 密码
		OtpSecret string `form:"otpSecret"` // 密码
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.ZapLogger.Warn("修改失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	res, err := system.SysAdminGetByID(c.GetInt("adminID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "修改失败" + err.Error()})
		return
	}
	res.OtpStatus = 1
	res.OtpTitle = request.OtpTitle
	res.OtpSecret = request.OtpSecret
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
		logger.ZapLogger.Error("生成二次验证器失败", zap.Error(err))
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
	// 用户改密请求体
	var request struct {
		Code   string `form:"code" binding:"required,len=6"` // 验证码
		Secret string `form:"secret" binding:"required"`     // 密钥
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.ZapLogger.Warn("二次验证失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误" + err.Error()})
		return
	}
	code, err := otp.GenerateOTPWithTime(request.Secret, time.Now())
	if code != request.Code {
		logger.ZapLogger.Warn("二次验证失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "二次验证失败"})
		return
	}

	c.JSON(http.StatusOK, bind.ErrorMessage{Message: "二次验证成功"})
}
