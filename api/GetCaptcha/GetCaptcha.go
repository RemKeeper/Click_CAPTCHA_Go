package GetCaptcha

import (
	"Click_CAPTCHA_Go/DrawImg"
	utils "Click_CAPTCHA_Go/Utils"
	Config "Click_CAPTCHA_Go/config"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetCaptcha(c *gin.Context) {
	DrawType := c.Query("type")
	DrawTypeInt, err := strconv.Atoi(DrawType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type参数错误"})
		return
	}
	var CaptchaRespRaw CaptchaResp
	switch Config.GlobalConfig.DrawServerMode {
	case 0:
		bytes, AnswerRaw, AnswerIcon := DrawImg.DrawCaptchaImg(DrawTypeInt)
		imageMd5, err := utils.GenerateImageMd5(bytes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"系统错误": err.Error()})
			return
		}
		marshal, err := json.Marshal(AnswerRaw)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"系统错误": err.Error()})
			return
		}
		CreateTimeInt64 := utils.GetCurrentTime()
		ImgBase64 := utils.EncodeBase64(bytes)
		var AnswerIconBase64Slice []string
		if DrawTypeInt == 0 {
			AnswerIconBase64Slice = make([]string, Config.GlobalConfig.DrawIconsCount)
			for i := 0; i < Config.GlobalConfig.DrawIconsCount; i++ {
				AnswerIconBase64Slice[i] = utils.EncodeBase64(AnswerIcon[i].ImgByte)
			}
		}

		CaptchaRespRaw = CaptchaResp{
			ImgId:        imageMd5,
			Img:          ImgBase64,
			CreateTime:   CreateTimeInt64,
			AnswerPrompt: AnswerIconBase64Slice,
			ExpireTime:   Config.GlobalConfig.CaptchaExpireTime,
		}
		utils.SetCaptchaToRedis(imageMd5, ImgBase64, string(marshal), Config.GlobalConfig.CaptchaExpireTime)
		c.JSON(http.StatusOK, CaptchaRespRaw)
	case 1:

	}
}
