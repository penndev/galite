package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/logger"
	"github.com/penndev/galite/internal/wafcdn/model"
	"go.uber.org/zap"
)

// 处理日志的批量写入
func HandlePutLog(c *gin.Context) {
	logsToInsert := []model.Log{}
	if err := c.BindJSON(&logsToInsert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误" + err.Error()})
		return
	}

	if err := (&model.Log{}).DB().Create(logsToInsert).Error; err != nil {
		logger.L.Error("批量写入日志失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"Message": "批量写入日志失败" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "完成"})
}
