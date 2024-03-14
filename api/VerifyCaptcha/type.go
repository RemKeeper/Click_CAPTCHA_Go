package VerifyCaptcha

type CaptchaReq struct {
	ImgId  string      `json:"imgId"`
	Answer []AnswerRaw `json:"answer"`
}

type AnswerRaw struct {
	X int
	Y int
}
