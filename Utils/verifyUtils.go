package utils

import (
	RedisUtils "Click_CAPTCHA_Go/Db/redis"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"time"
)

func GenerateImageMd5(imageRaw []byte) (string, error) {
	harsher := md5.New()
	harsher.Write(imageRaw)
	return hex.EncodeToString(harsher.Sum(nil)), nil
}

func EncodeBase64(Raw []byte) string {
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(Raw)
}

func GetCurrentTime() int64 {
	return time.Now().Unix()
}

func SetCaptchaToRedis(CaptchaId, CaptchaRaw, Answer string, ExpireTime int64) {
	RedisUtils.SetValueWithExpiration("Captcha:"+CaptchaId+":raw", CaptchaRaw, time.Duration(ExpireTime)*time.Second)
	RedisUtils.SetValueWithExpiration("Captcha:"+CaptchaId+":answer", Answer, time.Duration(ExpireTime)*time.Second)
}

func GetCaptchaFromRedis(CaptchaId string) string {
	return RedisUtils.GetValue("Captcha:" + CaptchaId + ":answer")
}
