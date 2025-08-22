package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/lib"
)

func HandleAcme(c *gin.Context) {
	// 处理ACME证书申请逻辑
	var token string
	err := lib.Cache.GetAny("acme:"+c.Query("token"), &token)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid token"+err.Error())
		return
	}
	c.JSON(http.StatusOK, token)
}
