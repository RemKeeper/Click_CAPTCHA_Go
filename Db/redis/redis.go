package RedisUtils

import (
	Config "Click_CAPTCHA_Go/config"
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

var Rdb *redis.Client

var RedisCtx = context.Background()

func ConnectRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     Config.GlobalConfig.RedisEndpoint,
		Password: Config.GlobalConfig.RedisPassword, // 没有密码，默认值
		DB:       Config.GlobalConfig.RedisDbIndex,  // 默认DB 0
	})
	//Rdb.Set(RedisCtx, "key", "value", 0)
}

func SetValue(key, value string) {
	Rdb.Set(RedisCtx, key, value, 0)
}

func GetValue(key string) string {
	val, _ := Rdb.Get(RedisCtx, key).Result()
	return val
}

func SetValueWithExpiration(key, value string, expiration time.Duration) {
	Rdb.Set(RedisCtx, key, value, expiration)
}
