package admin

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/admin/bind"
	"github.com/penndev/galite/pkg/util"
)

type filesItem struct {
	Name    string    `json:"name"`
	Dir     bool      `json:"dir"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
}

type filesMkdirParam struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

func filesRoot() (string, error) {
	root := util.GetEnv("WAFCDN_DATA_DIR", "/data")
	abs, err := filepath.Abs(filepath.Join(root, "static"))
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(abs, 0755); err != nil {
		return "", err
	}
	return abs, nil
}

func FilesResolve(rel string) (string, error) {
	root, err := filesRoot()
	if err != nil {
		return "", err
	}
	rel = strings.TrimSpace(rel)
	rel = strings.ReplaceAll(rel, "\\", "/")
	rel = strings.TrimPrefix(rel, "/")
	if rel == "" || rel == "." {
		return root, nil
	}
	for _, seg := range strings.Split(rel, "/") {
		if seg == "" || seg == "." || seg == ".." {
			return "", errors.New("非法路径")
		}
	}
	abs, err := filepath.Abs(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		return "", err
	}
	sep := string(os.PathSeparator)
	if abs != root && !strings.HasPrefix(abs, root+sep) {
		return "", errors.New("非法路径")
	}
	return abs, nil
}

func filesJoin(dir, name string) (string, error) {
	name = strings.TrimSpace(name)
	name = filepath.Base(strings.ReplaceAll(name, "\\", "/"))
	if name == "" || name == "." || name == ".." {
		return "", errors.New("非法名称")
	}
	return FilesResolve(strings.Trim(dir+"/"+name, "/"))
}

func FilesList(c *gin.Context) {
	abs, err := FilesResolve(c.Query("path"))
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	info, err := os.Stat(abs)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "目录不存在"})
		return
	}
	if !info.IsDir() {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "不是目录"})
		return
	}
	entries, err := os.ReadDir(abs)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "读取失败(" + err.Error() + ")"})
		return
	}
	dirs := make([]filesItem, 0)
	files := make([]filesItem, 0)
	for _, entry := range entries {
		item := filesItem{Name: entry.Name(), Dir: entry.IsDir()}
		if fi, err := entry.Info(); err == nil {
			item.Size = fi.Size()
			item.ModTime = fi.ModTime()
		}
		if item.Dir {
			dirs = append(dirs, item)
		} else {
			files = append(files, item)
		}
	}
	list := append(dirs, files...)
	c.JSON(http.StatusOK, bind.DataList{Total: int64(len(list)), Data: list})
}

func FilesMkdir(c *gin.Context) {
	param := &filesMkdirParam{}
	if err := c.BindJSON(param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
		return
	}
	abs, err := filesJoin(param.Path, param.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	if err := os.Mkdir(abs, 0755); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "创建失败(" + err.Error() + ")"})
		return
	}
	c.JSON(http.StatusOK, bind.Message{Message: "完成"})
}

func FilesRename(c *gin.Context) {
	param := &filesMkdirParam{}
	if err := c.BindJSON(param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
		return
	}
	oldAbs, err := FilesResolve(param.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	root, err := filesRoot()
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	if oldAbs == root {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "不能重命名根目录"})
		return
	}
	if _, err := os.Stat(oldAbs); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "目标不存在"})
		return
	}
	rel := strings.Trim(strings.ReplaceAll(param.Path, "\\", "/"), "/")
	dir := ""
	if idx := strings.LastIndex(rel, "/"); idx >= 0 {
		dir = rel[:idx]
	}
	newAbs, err := filesJoin(dir, param.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	if newAbs == oldAbs {
		c.JSON(http.StatusOK, bind.Message{Message: "完成"})
		return
	}
	if _, err := os.Stat(newAbs); err == nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "已存在同名文件或目录"})
		return
	}
	if err := os.Rename(oldAbs, newAbs); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "重命名失败(" + err.Error() + ")"})
		return
	}
	c.JSON(http.StatusOK, bind.Message{Message: "完成"})
}

func FilesUpload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "请选择文件"})
		return
	}
	defer file.Close()
	abs, err := filesJoin(c.PostForm("path"), header.Filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	info, err := os.Stat(abs)
	if err == nil {
		if c.PostForm("overwrite") != "1" || info.IsDir() {
			c.JSON(http.StatusBadRequest, bind.Message{Message: "已存在同名文件或目录"})
			return
		}
	} else if !os.IsNotExist(err) {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "保存失败(" + err.Error() + ")"})
		return
	}
	dst, err := os.OpenFile(abs, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "保存失败(" + err.Error() + ")"})
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "保存失败(" + err.Error() + ")"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "完成", "name": filepath.Base(abs)})
}

func FilesDownload(c *gin.Context) {
	abs, err := FilesResolve(c.Query("path"))
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	info, err := os.Stat(abs)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "文件不存在"})
		return
	}
	if info.IsDir() {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "不能下载目录"})
		return
	}
	c.FileAttachment(abs, info.Name())
}

func FilesDelete(c *gin.Context) {
	path := c.Query("path")
	if strings.TrimSpace(path) == "" {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "不能删除根目录"})
		return
	}
	root, err := filesRoot()
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	abs, err := FilesResolve(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	if abs == root {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "不能删除根目录"})
		return
	}
	if _, err := os.Stat(abs); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "目标不存在"})
		return
	}
	if err := os.RemoveAll(abs); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "删除失败(" + err.Error() + ")"})
		return
	}
	c.JSON(http.StatusOK, bind.Message{Message: "完成"})
}
