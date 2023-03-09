package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"robTickets/src/configs"
	"robTickets/src/controller"
	"robTickets/src/dao/mysql"
)

func main() {
	// 创建路由
	r := gin.Default()
	//r := gin.New()

	// 加载配置文件
	config := configs.InitConfig()

	// 创建数据库连接
	err := mysql.InitMysql(config)
	if err != nil {
		panic("MySQL init error：" + err.Error())
	}

	gin.SetMode(gin.ReleaseMode)
	// 路由
	//r.POST("/auth", controller2.AuthHandler)
	r.POST("/buy", controller.JWTAuthMiddleware(), controller.BuyTicket)

	// 监听端口
	addr := fmt.Sprintf(":%s", config.Port)
	if err := r.Run(addr); err != nil {
		panic("Project startup error!")
	}
}
