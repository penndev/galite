package lib

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	return n.runCommandWithOutput(cmd)
}

func (n *NginxManager) Status() bool {
	return n.isRunning()
}

func (n *NginxManager) isRunning() bool {
	pidFile := filepath.Join(n.Prefix, "logs", "nginx.pid")
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return false
	}
	pidStr := strings.TrimSpace(string(data))
	if pidStr == "" {
		return false
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tasklist", "/FI", "PID eq "+pidStr)
	} else {
		cmd = exec.Command("ps", "-p", pidStr)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), pidStr)
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
