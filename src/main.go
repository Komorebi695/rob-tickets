package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"robTickets/src/configs"
	"robTickets/src/controller"
	"robTickets/src/dao/mysql"
	"robTickets/src/dao/redis"
	"robTickets/src/logic"
)

const InitialCount = 1900

func main() {
	// 创建路由
	r := gin.Default()
	//r := gin.New()

	// 关闭Kafka
	defer logic.ProducerClose()
	// 加载配置 & 初始化
	config := initConfig()
	gin.SetMode(gin.ReleaseMode)
	// 路由
	r.POST("/buy", controller.JWTAuthMiddleware(), controller.BuyTicket)

	// 监听端口
	addr := fmt.Sprintf(":%s", config.Port)
	if err := r.Run(addr); err != nil {
		panic("Project startup error!")
	}
}

func initConfig() *configs.AppConfig {
	// 加载配置文件
	config := configs.InitConfig()
	// 创建mysql连接
	if err := mysql.InitMysql(config); err != nil {
		panic("MySQL init error：" + err.Error())
	}

	// 创建redis连接
	if err := redis.InitRedis(config); err != nil {
		panic("Redis init error：" + err.Error())
	}

	// 初始化redis缓存中的的初始票数量
	if err := logic.SetRedisTicketNumber("h1k4J7Dyt0", InitialCount); err != nil {
		panic("Redis set init ticket error：" + err.Error())
	}

	// 初始化Kafka
	if err := logic.InitKafka(); err != nil {
		panic("Kafka init error: " + err.Error())
	}
	return config
}
