package model

import (
	"fmt"
	"goblog/pkg/config"
	logger2 "goblog/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//db gorm.db 对象
var DB *gorm.DB

//connectDB 初始化模型
func ConnectDB() *gorm.DB  {

	var err error

	//初始化 mysql 连接信息
	var (
		host = config.GetString("database.mysql.host")
		port = config.GetString("database.mysql.port")
		database = config.GetString("database.mysql.database")
		username = config.GetString("database.mysql.username")
		password = config.GetString("database.mysql.password")
		charset = config.GetString("database.mysql.charset")
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
		username, password, host, port, database, charset, true, "Local")
	gormConfig := mysql.New(mysql.Config{
		DSN: dsn,
	})

	var level logger.LogLevel
	if config.GetBool("app.debug") {
		level = logger.Warn
	} else {
		level = logger.Error
	}

	//准备数据库连接池
	DB, err = gorm.Open(gormConfig, &gorm.Config{
		Logger: logger.Default.LogMode(level),
	})

	logger2.LogError(err)

	return DB
}
