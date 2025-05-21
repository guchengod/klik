package config

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	// DB 全局数据库连接
	DB *sqlx.DB
)

// InitDB 初始化数据库连接
func InitDB() {
	// 构建连接字符串
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		AppConfig.Database.Host, 
		AppConfig.Database.Port, 
		AppConfig.Database.User, 
		AppConfig.Database.Password, 
		AppConfig.Database.DBName, 
		AppConfig.Database.SSLMode)

	// 连接数据库
	var err error
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}

	// 设置连接池参数
	DB.SetMaxOpenConns(AppConfig.Database.MaxOpenConns)
	DB.SetMaxIdleConns(AppConfig.Database.MaxIdleConns)

	log.Println("数据库连接成功")
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("数据库连接已关闭")
	}
}
