package cache

import (
	"gin_mall/conf"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

var RedisClient *redis.Client

func InitCache() {
	Redis()
}

func Redis() {
	db, _ := strconv.ParseUint(conf.RedisDbName, 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr:     conf.RedisAddr,
		Password: conf.RedisPw,
		DB:       int(db),
	})
	_, err := client.Ping().Result()

	if err != nil {
		logrus.Info(err)
		panic(err)
	}
	RedisClient = client

}
