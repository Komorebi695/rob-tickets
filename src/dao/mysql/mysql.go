package mysql

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"robTickets/src/configs"
	"time"
)

var db *gorm.DB

// InitMysql 初始化MySQL
func InitMysql(initConfig *configs.AppConfig) (err error) {
	// 获取配置信息
	database := initConfig.MySQL

	// 日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Warn, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 禁用彩色打印
		},
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", database.User, database.Pwd, database.Host, database.Port, database.Database)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return
	}

	// 迁移
	err = db.AutoMigrate()
	if err != nil {
		return err
	}

	return err
}
