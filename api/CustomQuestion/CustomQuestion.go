package CustomQuestion

import (
	"Click_CAPTCHA_Go/Env"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// AddCustomQuestion Todo 添加Oss支持
func AddCustomQuestion(c *gin.Context) {
	var ReqData CustomData
	if err := c.BindJSON(&ReqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	questionMd5, err := ReqData.GenerateMD5()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"系统错误": err.Error()})
		return
	}

	var QuestionDir = Env.CustomQuestionPath + questionMd5
	err = os.Mkdir(QuestionDir, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"系统错误": err.Error()})
		return
	}
	err = os.WriteFile(QuestionDir+"/question.txt", []byte(ReqData.QuestionText), 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"系统错误": err.Error()})
		return
	}
	decodeByte, err := base64.StdEncoding.DecodeString(ReqData.Image[23:])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"图像解码错误": err.Error()})
		return
	}
	err = os.WriteFile(QuestionDir+"/image.png", decodeByte, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"系统错误": err.Error()})
		return
	}
	err = os.WriteFile(QuestionDir+"/answer.json", []byte(ReqData.Answer), 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"系统错误": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "创建成功"})

}
