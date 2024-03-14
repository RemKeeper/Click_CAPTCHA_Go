package main

import (
	"Click_CAPTCHA_Go/Db/redis"
	"Click_CAPTCHA_Go/Utils"
	"Click_CAPTCHA_Go/api/CustomQuestion"
	"Click_CAPTCHA_Go/api/GetCaptcha"
	"Click_CAPTCHA_Go/api/VerifyCaptcha"
	"Click_CAPTCHA_Go/config"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	Config.GlobalConfig = Config.GetConfig()
	utils.ReadBackgroundsToByteSlice()
	utils.ReadIconsToByteSlice()
	utils.ReadWordTextToSlice()
	utils.ReadWordFontToByteSlice()

	RedisUtils.ConnectRedis()

	rand.NewSource(time.Now().Unix())
	r := gin.Default()

	r.GET("/CAPTCHA", GetCaptcha.GetCaptcha)

	r.GET("/web/addCustomQuestion", func(context *gin.Context) {
		file, err := os.ReadFile("./web/CustomQuestionWeb/customQuestion.html")
		if err != nil {
			return
		}
		context.Data(http.StatusOK, "text/html", file)
	})

	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/verify", VerifyCaptcha.VerifyCaptcha)
		apiGroup.POST("/addCustomQuestion", CustomQuestion.AddCustomQuestion)
	}

	r.Run(":8080")
}
