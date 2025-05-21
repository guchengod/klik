package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

// 数据库连接信息
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "klik"
)

// 视频结构体
type Video struct {
	AwemeID      string      `json:"aweme_id"`
	Desc         string      `json:"desc"`
	CreateTime   int64       `json:"create_time"`
	AuthorUserID interface{} `json:"author_user_id"`
	Type         string      `json:"type"`
}

func main() {
	// 连接数据库
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	defer db.Close()

	// 测试连接
	err = db.Ping()
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	fmt.Println("成功连接到数据库!")

	// 创建表（如果不存在）
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS videos (
			id SERIAL PRIMARY KEY,
			aweme_id VARCHAR(50) UNIQUE NOT NULL,
			desc TEXT,
			create_time BIGINT,
			author_user_id VARCHAR(50),
			type VARCHAR(50) DEFAULT 'recommend-video'
		)
	`)
	if err != nil {
		log.Fatalf("创建表失败: %v", err)
	}

	// 获取项目根目录
	rootDir, err := filepath.Abs("../..")
	if err != nil {
		log.Fatalf("无法获取项目根目录: %v", err)
	}

	// 读取JSON文件
	jsonFile := filepath.Join(rootDir, "src", "assets", "data", "posts6.json")
	jsonData, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		log.Fatalf("无法读取JSON文件: %v", err)
	}

	// 解析JSON数据
	var videos []Video
	err = json.Unmarshal(jsonData, &videos)
	if err != nil {
		log.Fatalf("无法解析JSON数据: %v", err)
	}

	fmt.Printf("成功解析 %d 条视频数据\n", len(videos))

	// 导入数据
	for _, video := range videos {
		// 将作者ID转换为字符串
		var authorID string
		switch v := video.AuthorUserID.(type) {
		case string:
			authorID = v
		case float64:
			authorID = fmt.Sprintf("%d", int64(v))
		default:
			authorID = ""
		}

		// 插入数据
		_, err = db.Exec(`
			INSERT INTO videos (aweme_id, desc, create_time, author_user_id, type)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (aweme_id) DO NOTHING
		`, video.AwemeID, video.Desc, video.CreateTime, authorID, "recommend-video")
		if err != nil {
			log.Printf("插入数据失败: %v", err)
			continue
		}
	}

	// 查询数据
	rows, err := db.Query(`SELECT COUNT(*) FROM videos`)
	if err != nil {
		log.Fatalf("查询数据失败: %v", err)
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Fatalf("读取数据失败: %v", err)
		}
	}

	fmt.Printf("数据库中共有 %d 条视频数据\n", count)
}
