package wafcdn

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/model/wafcdn"
)

// 对nginx提供接口 获取域名配置信息
// @url=/@wafcdn/domain?host=@host
// @return 配置信息
func handleDomain(c *gin.Context) {
	// 读取文件内容
	data, err := os.ReadFile("../wafcdn/docs/domain.json")
	if err != nil {
		fmt.Println("读取文件错误:", err)
		return
	}

	c.String(200, string(data))
}

func handleGetCache(c *gin.Context) {
	param := &wafcdn.Cache{}
	if err := c.BindQuery(param); err != nil {
		log.Println("参数错误", err.Error())
		c.JSON(400, gin.H{"message": "参数错误" + err.Error()})
		return
	}

	if err := param.Bind(param).First(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "查询失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, param)
	}
}

func handlePutCache(c *gin.Context) {
	param := &wafcdn.Cache{}
	if err := c.BindJSON(param); err != nil {
		log.Println("参数错误", err.Error())
		c.JSON(400, gin.H{"message": "参数错误" + err.Error()})
		return
	}

	if err := param.Bind(param).Assign(*param).FirstOrCreate(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "创建失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "完成"})
	}
}
