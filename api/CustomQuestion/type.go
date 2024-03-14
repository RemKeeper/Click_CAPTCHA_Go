package CustomQuestion

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

type CustomData struct {
	QuestionText string `json:"questionText"`
	Image        string `json:"image"`
	Answer       string `json:"answer"`
}

func (cq *CustomData) GenerateMD5() (string, error) {
	bytes, err := json.Marshal(cq)
	if err != nil {
		return "", err
	}
	harsher := md5.New()
	harsher.Write(bytes)
	return hex.EncodeToString(harsher.Sum(nil)), nil
}
