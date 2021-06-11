package model

import (
	logger2 "goblog/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//db gorm.db 对象
var DB *gorm.DB

//connectDB 初始化模型
func ConnectDB() *gorm.DB  {

	var err error

	config := mysql.New(mysql.Config{
		DSN:"rdbuser:shop2db123#@tcp(101.37.150.149:3306)/blog?charset=utf-8&parseTime=True&loc=Local",
	})

	//准备数据库连接池
	DB, err = gorm.Open(config, &gorm.Config{})

	logger2.LogError(err)

	return DB
}
