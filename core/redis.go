package core

import (
	"context"
	"gin_chat/global"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"), // 没有密码，默认值
		DB:           viper.GetInt("redis.db"),          // 默认DB 0
		PoolSize:     viper.GetInt("redis.pool_size"),
		MinIdleConns: viper.GetInt("redis.min_idle_conns"),
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		global.Log.Errorf("初始化redis失败：%v", err.Error())
		return nil
	}
	global.Log.Info("初始化redis成功")
	return rdb
}
