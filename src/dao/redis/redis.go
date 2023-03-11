package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"robTickets/src/configs"
)

var rdb *redis.Client

// InitRedis ,初始化redis
func InitRedis(initConfig *configs.AppConfig) (err error) {
	conf := initConfig.Redis
	addr := fmt.Sprintf("%s:%s", conf.Addr, conf.Port)
	//addr := fmt.Sprintf("%s:%s", conf.Addr, "3300")
	rdb = redis.NewClient(&redis.Options{
		Addr: addr,
		//Password: conf.Password,
		DB: conf.DB,
	})

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}
