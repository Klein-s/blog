package database

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	logger2 "goblog/pkg/logger"
	"time"
)

var DB *sql.DB

func Initialize()  {
	initDB()
	createTables()
}
//数据库驱动初始化
func initDB()  {

	var err error
	config := mysql.Config{
		User: "rdbuser",
		Passwd: "shop2db123#",
		Addr: "101.37.150.149:3306",
		Net: "tcp",
		DBName: "blog",
		AllowNativePasswords: true,
	}

	//准备数据库连接池
	DB, err = sql.Open("mysql", config.FormatDSN())
	logger2.LogError(err)

	//设置最大连接数
	DB.SetMaxOpenConns(25)
	//设置最大空闲连接数
	DB.SetMaxIdleConns(25)
	//设置每个连接的过期时间
	DB.SetConnMaxLifetime(5 * time.Minute)

	//尝试连接 失败报错
	err = DB.Ping()
	logger2.LogError(err)
}


func createTables() {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
    id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    body longtext COLLATE utf8mb4_unicode_ci
); `

	_, err := DB.Exec(createArticlesSQL)
	logger2.LogError(err)
}
