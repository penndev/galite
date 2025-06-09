package analyze

import "github.com/gin-gonic/gin"

type configApp struct {
	Version string `json:"version"`
	Url     string `json:"url"`
	Remark  string `json:"remark"`
}

type resultConfig struct {
	Android configApp `json:"android"`
	Ios     configApp `json:"ios"`
}

type requestStart struct {
	App      string `json:"app"`       // 渠道码
	Channel  string `json:"channel"`   // 渠道码
	DeviceID string `json:"device_id"` // 设备ID
}

func Start(c *gin.Context) {
	c.ShouldBindJSON(&requestStart{})
	c.JSON(200, resultConfig{
		Android: configApp{
			Version: "1.0.0",
			Url:     "https://example.com/android.apk",
			Remark:  "Android应用最新版本",
		},
		Ios: configApp{
			Version: "1.0.0",
			Url:     "https://example.com/ios.ipa",
			Remark:  "iOS应用最新版本",
		},
	})
}
