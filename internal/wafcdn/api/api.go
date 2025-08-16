package api

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/admin/bind"
	"github.com/penndev/galite/internal/logger"
	"github.com/penndev/galite/internal/wafcdn/model"
	"github.com/penndev/galite/pkg/util"
	"github.com/penndev/gopkg/ip2region"
	"go.uber.org/zap"
)

var (
	logBuffer     = make([]*model.Log, 0, 1000) // 初始容量1000的缓冲区
	bufferMutex   = &sync.Mutex{}
	maxBufferSize = 1000            // 缓冲区最大大小
	flushInterval = 5 * time.Second // 定期刷新间隔
)

// 实际执行批量写入操作
func flushBuffer() {
	if len(logBuffer) == 0 {
		return
	}

	// 复制当前缓冲区内容并清空缓冲区
	logsToInsert := make([]*model.Log, len(logBuffer))
	copy(logsToInsert, logBuffer)
	logBuffer = logBuffer[:0]

	// 批量写入数据库
	if err := (&model.Log{}).DB().Create(logsToInsert).Error; err != nil {
		logger.L.Error("批量写入日志失败", zap.Error(err))
		return
	}
	log.Printf("成功批量写入 %d 条日志", len(logsToInsert))
}

// 启动定期刷新协程
func init() {
	go func() {
		ticker := time.NewTicker(flushInterval)
		defer ticker.Stop()

		for range ticker.C {
			bufferMutex.Lock()
			if len(logBuffer) > 0 {
				flushBuffer()
			}
			bufferMutex.Unlock()
		}
	}()
}

func HandlePutLog(c *gin.Context) {
	param := &model.Log{}
	if err := c.BindJSON(param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误" + err.Error()})
		return
	}

	// 将日志添加到缓冲区
	bufferMutex.Lock()
	logBuffer = append(logBuffer, param)
	if len(logBuffer) >= maxBufferSize {
		// 如果缓冲区满了，立即刷新
		flushBuffer()
	}
	bufferMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"Message": "完成"})
}

// func HandlePutLog(c *gin.Context) {
// 	param := &model.Log{}
// 	if err := c.BindJSON(param); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误" + err.Error()})
// 		return
// 	}
// 	param.DB().Create(param)
// 	c.JSON(http.StatusOK, gin.H{"Message": "完成"})
// }

func HandleIpCheck(c *gin.Context) {
	site := c.Query("site")
	siteID, err := util.StrConv[uint](site)
	if err != nil {
		c.JSON(http.StatusNotFound, bind.Message{})
		return
	}
	siteModel, err := model.GetSiteByID(siteID)
	if err != nil {
		c.JSON(http.StatusNotFound, bind.Message{})
		return
	}
	ipStruct := siteModel.Security.Ip
	if !ipStruct.Status {
		c.JSON(http.StatusNotFound, bind.Message{})
		return
	}
	if len(ipStruct.Region) > 0 {
		reqRegion := ip2region.Find(c.Query("ip"))
		for _, region := range ipStruct.Region {
			if region.Country != "" && region.Country != reqRegion.Country { // 国家
				continue
			}
			if region.Province != "" && region.Province != reqRegion.Province { // 省份
				continue
			}
			if region.City != "" && region.City != reqRegion.City { // 城市
				continue
			}
			if region.County != "" && region.County != reqRegion.County { // 县级
				continue
			}
			if ipStruct.Allowed { // 如果能匹配到任意一条规则，则根据 Allowed 返回
				c.JSON(http.StatusOK, bind.Message{Message: "ok"})
			} else {
				c.JSON(http.StatusNotFound, bind.Message{Message: "fail"})
			}
			return
		}
	}

	if !ipStruct.Allowed { // 最终效果需要取反
		c.JSON(http.StatusOK, bind.Message{Message: "ok"})
	} else {
		c.JSON(http.StatusNotFound, bind.Message{Message: "fail"})
	}
}
