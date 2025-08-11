package lib

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var Nginx NginxManager

type NginxManager struct {
	Binary     string // nginx 可执行文件路径
	Prefix     string // -p 指定的路径
	OutputFile string // 用于日志输出的文件（stdout + stderr）
}

func (n *NginxManager) Start() error {
	if err := n.TestConfig(); err != nil {
		return err
	}

	if n.isRunning() {
		return errors.New("nginx is already running")
	}
	cmd := exec.Command(n.Binary, "-p", n.Prefix)
	return n.runCommandWithOutput(cmd)
}

func (n *NginxManager) Stop() error {
	if !n.isRunning() {
		return errors.New("nginx is not running")
	}
	cmd := exec.Command(n.Binary, "-p", n.Prefix, "-s", "stop")
	return n.runCommandWithOutput(cmd)
}

func (n *NginxManager) Reload() error {
	if err := n.TestConfig(); err != nil {
		return err
	}

	if !n.isRunning() {
		return errors.New("nginx is not running")
	}
	cmd := exec.Command(n.Binary, "-p", n.Prefix, "-s", "reload")
	return n.runCommandWithOutput(cmd)
}

func (n *NginxManager) Restart() error {
	if err := n.TestConfig(); err != nil {
		return fmt.Errorf("config test failed: %w", err)
	}
	_ = n.Stop()
	return n.Start()
}

func (n *NginxManager) TestConfig() error {
	cmd := exec.Command(n.Binary, "-p", n.Prefix, "-t")

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()

	// 无论 err 是否为 nil，都要检查 stderr 内容中是否包含 "syntax is ok" 等
	output := outBuf.String() + errBuf.String()
	if err != nil || !strings.Contains(output, "syntax is ok") {
		return fmt.Errorf("nginx config test failed: %s", strings.TrimSpace(output))
	}
	return err
}

func (n *NginxManager) Status() bool {
	return n.isRunning()
}

func (n *NginxManager) Version(verbose bool) (string, error) {
	var arg string
	if verbose {
		arg = "-V"
	} else {
		arg = "-v"
	}

	cmd := exec.Command(n.Binary, arg)

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()

	// 合并 stdout 和 stderr
	output := outBuf.String()
	if errBuf.Len() > 0 {
		output += errBuf.String()
	}

	output = strings.TrimSpace(output)

	if !verbose {
		// 提取版本号，如 openresty/1.27.1.1 或 nginx/1.24.0
		re := regexp.MustCompile(`(?:nginx|openresty)/[\d\.]+`)
		match := re.FindString(output)
		if match != "" {
			return match, nil
		}
		return "", fmt.Errorf("unexpected version format: %s", output)
	}

	return output, err
}

func (n *NginxManager) isRunning() bool {
	pidFile := filepath.Join(n.Prefix, "logs", "nginx.pid")

	// 读取 PID 文件
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return false
	}
	pidStr := strings.TrimSpace(string(data))
	if pidStr == "" {
		return false
	}

	// Windows 系统：使用 tasklist 检查
	if runtime.GOOS == "windows" {
		cmd := exec.Command("tasklist", "/FI", "PID eq "+pidStr)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return false
		}
		return strings.Contains(string(out), pidStr)
	}

	// 非 Windows（Linux / macOS）：检查 /proc/<pid> 是否存在
	procPath := filepath.Join("/proc", pidStr)
	_, err = os.Stat(procPath)
	if err != nil {
		return false
	}

	// 可选增强：确保是 nginx 进程
	cmdlinePath := filepath.Join(procPath, "cmdline")
	cmdline, err := os.ReadFile(cmdlinePath)
	if err != nil {
		return false
	}
	if !bytes.Contains(cmdline, []byte("nginx")) {
		return false
	}

	return true
}

// 支持输出写入文件
func (n *NginxManager) runCommandWithOutput(cmd *exec.Cmd) error {
	if n.Prefix != "" {
		cmd.Dir = n.Prefix
	}
	if n.OutputFile != "" {
		f, err := os.OpenFile(n.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		cmd.Stdout = f
		cmd.Stderr = f
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Start()
}
