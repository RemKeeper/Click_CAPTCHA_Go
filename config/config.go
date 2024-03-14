package Config

import (
	"encoding/json"
	"fmt"
	"os"
)

var GlobalConfig Config

type Config struct {
	//绘制模式 0:即时绘制 1:预先绘制
	//
	//即时绘制模式:每次请求都会重新绘制一张图片
	//
	//预先绘制模式:预先绘制一批图片,存入缓存,每次请求从缓存中取一张图片
	DrawServerMode int
	//预先绘制更新时间(秒)
	BeforehandDrawUpdateTimes int
	//点击允许偏移量(像素)
	ClickToleranceValue int
	//图标绘制数量
	DrawIconsCount int
	//是否开启字符旋转绘制

	//绘制文字大小
	DrawWordSize float64

	//DrawWordRotation bool

	//是否开启字符变形绘制
	DrawWordDeformation bool
	//是否开启字符颜色绘制
	DrawColorfulWord bool

	//验证码过期时间(秒)
	CaptchaExpireTime int64

	//	运行配置
	RedisEndpoint string

	RedisPassword string

	RedisDbIndex int
}

func GetConfig() Config {
	file, err := os.ReadFile("config.json")
	if err != nil {
		CreateConfig()
	}
	var config Config
	_ = json.Unmarshal(file, &config)
	return config
}

func CreateConfig() {
	config := Config{
		DrawServerMode:            0,
		BeforehandDrawUpdateTimes: 0,
		ClickToleranceValue:       50,
		DrawIconsCount:            3,
		DrawWordSize:              30,
		//DrawWordRotation:          true,
		DrawWordDeformation: true,
		DrawColorfulWord:    true,
	}
	marshal, _ := json.Marshal(config)
	_ = os.WriteFile("config.json", marshal, 0644)
	fmt.Println("配置文件创建成功,请修改配置文件后重启程序")
	os.Exit(0)
}
