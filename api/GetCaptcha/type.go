package GetCaptcha

type CaptchaResp struct {
	ImgId        string   `json:"imgId"`
	Img          string   `json:"img"`
	AnswerPrompt []string `json:"answerPrompt"`
	CreateTime   int64    `json:"createTime"`
	ExpireTime   int64    `json:"expireTime"`
}
