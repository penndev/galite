package system

// 用户登录请求体
type bindLoginInput struct {
	Username  string // 用户名
	Password  string // 密码
	Captcha   string // 验证码
	CaptchaId string // 验证码ID
}
