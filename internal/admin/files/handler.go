package files

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/admin/bind"
)

const maxReadSize = 2 << 20 // 2MB

var root string

func init() {
	abs, err := filepath.Abs("/home")
	if err != nil {
		panic(err)
	}
	root = filepath.Clean(abs)
}

func resolvePath(p string) (string, error) {
	if strings.TrimSpace(p) == "" {
		return root, nil
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	abs = filepath.Clean(abs)
	if abs != root {
		rel, err := filepath.Rel(root, abs)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return "", fmt.Errorf("路径超出允许范围")
		}
	}
	return abs, nil
}

func handleList(c *gin.Context) {
	param := &bindList{}
	if err := c.BindQuery(param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	abs, err := resolvePath(param.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	items, err := os.ReadDir(abs)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	entries := make([]gin.H, 0, len(items))
	for _, item := range items {
		name := item.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		fi, err := item.Info()
		if err != nil {
			continue
		}
		entries = append(entries, gin.H{
			"name":    name,
			"path":    filepath.Join(abs, name),
			"isDir":   fi.IsDir(),
			"size":    fi.Size(),
			"mode":    fi.Mode(),
			"modTime": fi.ModTime(),
		})
	}
	c.JSON(http.StatusOK, bind.DataList{Data: gin.H{"path": abs, "entries": entries}})
}

func handleGetContent(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "缺少 path 参数"})
		return
	}
	abs, err := resolvePath(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	info, err := os.Stat(abs)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	if info.IsDir() {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "是目录"})
		return
	}
	if info.Size() > maxReadSize {
		c.JSON(http.StatusBadRequest, bind.Message{Message: fmt.Sprintf("文件过大（%d 字节），最大 %d 字节", info.Size(), maxReadSize)})
		return
	}
	data, err := os.ReadFile(abs)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, bind.DataList{Data: map[string]string{
		"path":    path,
		"content": string(data),
	}})
}

func handlePutContent(c *gin.Context) {
	param := &bindContent{}
	if err := c.BindJSON(param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	abs, err := resolvePath(param.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	if err := os.WriteFile(abs, []byte(param.Content), 0o644); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, bind.Message{Message: "完成"})
}

func handleDownload(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "缺少 path 参数"})
		return
	}
	abs, err := resolvePath(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	info, err := os.Stat(abs)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	if info.IsDir() {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "是目录"})
		return
	}
	f, err := os.Open(abs)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	defer f.Close()
	c.Header("Content-Disposition", "attachment; filename="+info.Name())
	c.Header("Content-Type", "application/octet-stream")
	c.DataFromReader(http.StatusOK, info.Size(), "application/octet-stream", f, nil)
}

func handleUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	defer f.Close()
	dirAbs, err := resolvePath(c.PostForm("path"))
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	info, err := os.Stat(dirAbs)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	if !info.IsDir() {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "不是目录"})
		return
	}
	target := filepath.Join(dirAbs, filepath.Base(file.Filename))
	out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	defer out.Close()
	if _, err := io.Copy(out, f); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, bind.DataList{Data: map[string]string{
		"path": target,
		"name": file.Filename,
	}})
}

func handleMkdir(c *gin.Context) {
	param := &bindPath{}
	if err := c.BindJSON(param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	abs, err := resolvePath(param.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, bind.Message{Message: "完成"})
}

func handleDelete(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "缺少 path 参数"})
		return
	}
	abs, err := resolvePath(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	if err := os.RemoveAll(abs); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, bind.Message{Message: "完成"})
}

func handleRename(c *gin.Context) {
	param := &bindRename{}
	if err := c.BindJSON(param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	oldAbs, err := resolvePath(param.OldPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	newAbs, err := resolvePath(param.NewPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	if err := os.Rename(oldAbs, newAbs); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, bind.Message{Message: "完成"})
}

func handleChmod(c *gin.Context) {
	param := &bindChmod{}
	if err := c.BindJSON(param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	mode, err := strconv.ParseUint(param.Mode, 8, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "mode 须为八进制，如 0755"})
		return
	}
	abs, err := resolvePath(param.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	if err := os.Chmod(abs, os.FileMode(mode)); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, bind.Message{Message: "完成"})
}
