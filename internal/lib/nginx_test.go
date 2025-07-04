package lib_test

import (
	"testing"

	"github.com/penndev/galite/internal/lib"
)

func TestNginxManager(t *testing.T) {
	nm := lib.NginxManager{
		Binary:     "openresty",
		Prefix:     "C:/Users/Penn/Dev/github/wafcdn",
		OutputFile: "./nginx.log",
	}
	t.Log("配置文件")
	t.Log(nm.TestConfig())
	t.Log(nm.Version(false))
	t.Log(nm.Version(true))
	t.Log("准备启动")
	t.Log(nm.Start())
	t.Log("查看状态")
	t.Log(nm.Status())
	t.Log("重新启动")
	t.Log(nm.Reload())
	t.Log("关闭启动")
	t.Log(nm.Stop())
	t.Log("结束")
}
