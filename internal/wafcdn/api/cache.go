package api

import (
	"bytes"
	"log"
	"net/http"

	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/lib"
	"github.com/penndev/galite/internal/logger"
	"github.com/penndev/galite/internal/wafcdn/model"
	"go.uber.org/zap"
)

func HandleGetCache(c *gin.Context) {
	param := &model.Cache{}
	if err := c.BindQuery(param); err != nil {
		log.Println("参数错误", err.Error())
		c.JSON(400, gin.H{"message": "参数错误" + err.Error()})
		return
	}
	// 直接从sql中取导致io瓶颈用key vals 做持久化加速
	// err := param.Bind(param).Where("site_id = ? and method = ? and uri = ?", param.SiteID, param.Method, param.Uri).First(param).Error;
	// 赋值直接给param
	err := lib.Badger.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(param.CacheKey()))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			return lib.Decode(bytes.NewBuffer(val), param)
		})
		return err
	})
	// 解决容错问题判断是否设置一致性。
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "查询失败(" + err.Error() + ")"})
		return
	} else {
		c.JSON(http.StatusOK, param)
	}
}

func HandlePutCache(c *gin.Context) {
	param := &model.Cache{}
	if err := c.BindJSON(param); err != nil {
		log.Println("参数错误", err.Error())
		c.JSON(400, gin.H{"message": "参数错误" + err.Error()})
		return
	}
	if err := lib.Badger.Update(func(txn *badger.Txn) error {
		data, err := lib.Encode(&model.Cache{
			Header: param.Header,
			Path:   param.Path,
			Time:   param.Time,
		})
		if err != nil {
			return err
		}
		err = txn.Set([]byte(param.CacheKey()), data.Bytes())
		return err
	}); err != nil {
		logger.L.Error("wafcdn/cache", zap.Error(err))
	}

	if err := param.Bind(param).Where("site_id = ? and method = ? and uri = ?", param.SiteID, param.Method, param.Uri).Assign(*param).FirstOrCreate(&model.Cache{
		SiteID: param.SiteID,
		Method: param.Method,
		Uri:    param.Uri,
		Path:   param.Path,
	}).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "创建失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "完成"})
	}
}
