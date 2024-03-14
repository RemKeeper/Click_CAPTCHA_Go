package VerifyCaptcha

import (
	utils "Click_CAPTCHA_Go/Utils"
	Config "Click_CAPTCHA_Go/config"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"math"
)

func VerifyCaptcha(c *gin.Context) {
	var CaptchaReqRaw CaptchaReq
	err := c.BindJSON(&CaptchaReqRaw)
	if err != nil {
		c.JSON(400, gin.H{"Bad Request": err.Error()})
		return
	}
	answerStr := utils.GetCaptchaFromRedis(CaptchaReqRaw.ImgId)

	if len(answerStr) == 0 {
		c.JSON(400, gin.H{"Bad Request": "CAPTCHA expired"})
		return
	}

	var RealAnswer []AnswerRaw
	json.Unmarshal([]byte(answerStr), &RealAnswer)
	if len(CaptchaReqRaw.Answer) != len(RealAnswer) {
		c.JSON(400, gin.H{"Bad Request": "incorrect answer"})
		return
	}
	if !VerifyAnswer(CaptchaReqRaw, RealAnswer) {
		c.JSON(400, gin.H{"Bad Request": "incorrect answer"})
		return
	}
	c.JSON(200, gin.H{"success": "true"})
}

func VerifyAnswer(GetAnswer CaptchaReq, RealAnswer []AnswerRaw) bool {
	for i := 0; i < len(GetAnswer.Answer); i++ {
		if math.Abs(float64(GetAnswer.Answer[i].X-RealAnswer[i].X)) > float64(Config.GlobalConfig.ClickToleranceValue) || math.Abs(float64(GetAnswer.Answer[i].Y-RealAnswer[i].Y)) > float64(Config.GlobalConfig.ClickToleranceValue) {
			return false
		}
	}
	return true
}
