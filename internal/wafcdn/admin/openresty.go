package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/admin/bind"
	"github.com/penndev/galite/internal/lib"
)

func OpenrestyStatus(c *gin.Context) {
	version, _ := lib.Nginx.Version(false)
	info, _ := lib.Nginx.Version(true)
	c.JSON(200, gin.H{
		"status":  lib.Nginx.Status(),
		"version": version,
		"info":    info,
	})
}

func OpenrestyStart(c *gin.Context) {
	if err := lib.Nginx.Start(); err != nil {
		c.JSON(400, bind.Message{
			Message: err.Error(),
		})
		return
	}
	c.JSON(200, bind.Message{
		Message: "操作完成",
	})
}

func OpenrestyStop(c *gin.Context) {
	if err := lib.Nginx.Stop(); err != nil {
		c.JSON(400, bind.Message{
			Message: err.Error(),
		})
		return
	}
	c.JSON(200, bind.Message{
		Message: "操作完成",
	})
}

func OpenrestyReload(c *gin.Context) {
	if err := lib.Nginx.Reload(); err != nil {
		c.JSON(400, bind.Message{
			Message: err.Error(),
		})
		return
	}
	c.JSON(200, bind.Message{
		Message: "操作完成",
	})
}
