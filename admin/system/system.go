package system

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/penndev/galite/admin/bind"
	"github.com/penndev/galite/config"
	"github.com/penndev/galite/model/system"
	"github.com/penndev/gopkg/captcha"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/logger"
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
	res, err := system.SysAdminGetByEmail(request.Username)
	if err != nil {
		var msg = "获取用户失败"
		if request.Username == "admin" && errors.Is(err, logger.ErrRecordNotFound) {
			msg = "已初始化管理员，请再次点击登录"
			res.Email = request.Username
			str, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
			if err != nil {
				config.Logger.Error("初始化管理员失败", zap.Error(err))
				msg = "初始化管理员失败，请查看错误日志"
			} else {
				res.Passwd = string(str)
				if err = res.Bind(res).Create(&res).Error; err != nil {
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
		"sub": res.ID,
		"exp": jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour))}).SignedString([]byte(config.JWTSecret))
	if err != nil {
		config.Logger.Error("用户登录失败", zap.Error(err))
		c.JSON(http.StatusForbidden, bind.ErrorMessage{Message: "用户登录失败(jwt签名错误)"})
		return
	}
	// 超级管理员
	if res.SysRoleID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"token":  token,
			"routes": "*", //超级管理员
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"token":  token,
			"routes": res.SysRole.Menu, //前端菜单解决方案
		})
	}
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

// 用户菜单鉴权
func Role(isLog bool) gin.HandlerFunc {

	return func(c *gin.Context) {
		admin, err := system.SysAdminGetByID(c.GetString("jwtAuth"))
		if err != nil {
			c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "用户鉴权失败"})
			c.Abort()
			return
		}
		// log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", admin.SysRoleID)
		if admin.SysRoleID != 0 {
			routes := admin.SysRole.Route
			pass := false
			for _, route := range routes {
				if route.Method == c.Request.Method && route.Path == c.Request.URL.Path {
					pass = true
					break
				}
			}

			if !pass {
				c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "用户鉴权失败"})
				c.Abort()
				return
			}
		}

		if isLog { // 记录日志
			access := &system.SysAccessLog{
				SysAdminID: admin.ID,
				Method:     c.Request.Method,
				Path:       fmt.Sprint(c.Request.URL),
				IP:         c.ClientIP(),
			}
			if err := access.Bind(access).Create(access).Error; err != nil {
				c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "日志记录失败:" + err.Error()})
				c.Abort()
				return
			}
			c.Next()
			httpRequest, _ := httputil.DumpRequest(c.Request, false)
			access.Payload = string(httpRequest)
			access.Status = c.Writer.Status()
			access.Bind(access).Updates(access)
		}
	}
}

func AccessLog(c *gin.Context) {
	param := &bindSysAccessParam{}
	if err := c.BindQuery(&param); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	var total int64
	var list []system.SysAccessLog

	m := param.Param() //处理筛选
	m.List(&total, &list)
	c.JSON(http.StatusOK, bind.DataList{Total: total, Data: list})
}
