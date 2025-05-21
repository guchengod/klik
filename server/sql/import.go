package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
)

// 视频结构体
type Video struct {
	AwemeID      string      `json:"aweme_id"`
	Desc         string      `json:"desc"`
	CreateTime   int64       `json:"create_time"`
	Music        Music       `json:"music"`
	VideoInfo    VideoInfo   `json:"video"`
	Statistics   Statistics  `json:"statistics"`
	Status       Status      `json:"status"`
	AuthorUserID interface{} `json:"author_user_id"`
	Duration     int         `json:"duration"`
	Author       Author      `json:"author"`
}

// 音乐结构体
type Music struct {
	ID            interface{} `json:"id"`
	Title         string      `json:"title"`
	Author        string      `json:"author"`
	CoverMedium   Cover       `json:"cover_medium"`
	PlayURL       Cover       `json:"play_url"`
	Duration      int         `json:"duration"`
	OwnerID       string      `json:"owner_id"`
	OwnerNickname string      `json:"owner_nickname"`
	IsOriginal    bool        `json:"is_original"`
}

// 视频信息结构体
type VideoInfo struct {
	PlayAddr       Cover  `json:"play_addr"`
	Cover          Cover  `json:"cover"`
	Height         int    `json:"height"`
	Width          int    `json:"width"`
	Ratio          string `json:"ratio"`
	UseStaticCover bool   `json:"use_static_cover"`
	Duration       int    `json:"duration"`
}

// 封面结构体
type Cover struct {
	URI     string   `json:"uri"`
	URLList []string `json:"url_list"`
	Width   int      `json:"width"`
	Height  int      `json:"height"`
}

// 统计结构体
type Statistics struct {
	CommentCount int `json:"comment_count"`
	DiggCount    int `json:"digg_count"`
	CollectCount int `json:"collect_count"`
	PlayCount    int `json:"play_count"`
	ShareCount   int `json:"share_count"`
}

// 状态结构体
type Status struct {
	IsDelete      bool `json:"is_delete"`
	AllowShare    bool `json:"allow_share"`
	IsProhibited  bool `json:"is_prohibited"`
	InReviewing   bool `json:"in_reviewing"`
	PrivateStatus int  `json:"private_status"`
}

// 作者结构体
type Author struct {
	Avatar168x168  Cover       `json:"avatar_168x168"`
	Avatar300x300  Cover       `json:"avatar_300x300"`
	UID            string      `json:"uid"`
	Nickname       string      `json:"nickname"`
	Signature      string      `json:"signature"`
	Gender         int         `json:"gender"`
	FollowerCount  int         `json:"follower_count"`
	FollowingCount int         `json:"following_count"`
	TotalFavorited int         `json:"total_favorited"`
	UniqueID       string      `json:"unique_id"`
	ShortID        string      `json:"short_id"`
	CoverURL       []Cover     `json:"cover_url"`
	WhiteCoverURL  []Cover     `json:"white_cover_url"`
	IPLocation     string      `json:"ip_location"`
	Province       string      `json:"province"`
	City           string      `json:"city"`
	Country        string      `json:"country"`
	AwemeCount     int         `json:"aweme_count"`
}

// 数据库连接信息
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "klik"
)

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

	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("无法开始事务: %v", err)
	}

	// 导入数据
	for _, video := range videos {
		// 导入用户数据
		userID, err := importUser(tx, video.Author)
		if err != nil {
			tx.Rollback()
			log.Fatalf("导入用户数据失败: %v", err)
		}

		// 导入音乐数据
		musicID, err := importMusic(tx, video.Music)
		if err != nil {
			tx.Rollback()
			log.Fatalf("导入音乐数据失败: %v", err)
		}

		// 导入视频数据
		videoID, err := importVideo(tx, video, musicID)
		if err != nil {
			tx.Rollback()
			log.Fatalf("导入视频数据失败: %v", err)
		}

		// 导入视频播放地址
		err = importVideoPlayAddress(tx, videoID, video.VideoInfo.PlayAddr)
		if err != nil {
			tx.Rollback()
			log.Fatalf("导入视频播放地址失败: %v", err)
		}

		// 导入视频封面
		err = importVideoCover(tx, videoID, video.VideoInfo.Cover)
		if err != nil {
			tx.Rollback()
			log.Fatalf("导入视频封面失败: %v", err)
		}

		// 导入视频统计
		err = importVideoStatistics(tx, videoID, video.Statistics)
		if err != nil {
			tx.Rollback()
			log.Fatalf("导入视频统计失败: %v", err)
		}

		// 导入视频状态
		err = importVideoStatus(tx, videoID, video.Status)
		if err != nil {
			tx.Rollback()
			log.Fatalf("导入视频状态失败: %v", err)
		}

		fmt.Printf("成功导入视频: %s\n", video.AwemeID)
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		log.Fatalf("无法提交事务: %v", err)
	}

	fmt.Println("数据导入完成!")
}

// 导入用户数据
func importUser(tx *sql.Tx, author Author) (int, error) {
	var userID int

	// 检查用户是否已存在
	err := tx.QueryRow("SELECT id FROM users WHERE uid = $1", author.UID).Scan(&userID)
	if err == nil {
		// 用户已存在，返回用户ID
		return userID, nil
	} else if err != sql.ErrNoRows {
		// 发生其他错误
		return 0, err
	}

	// 插入用户数据
	err = tx.QueryRow(`
		INSERT INTO users (
			uid, nickname, gender, signature, ip_location, province, city, country,
			follower_count, following_count, total_favorited, aweme_count, unique_id, short_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`, author.UID, author.Nickname, author.Gender, author.Signature, author.IPLocation,
		author.Province, author.City, author.Country, author.FollowerCount, author.FollowingCount,
		author.TotalFavorited, author.AwemeCount, author.UniqueID, author.ShortID).Scan(&userID)
	if err != nil {
		return 0, err
	}

	// 插入头像
	if len(author.Avatar168x168.URI) > 0 {
		var avatarURL string
		if len(author.Avatar168x168.URLList) > 0 {
			avatarURL = author.Avatar168x168.URLList[0]
		}
		_, err = tx.Exec(`
			UPDATE users SET avatar_168x168_uri = $1, avatar_168x168_url = $2
			WHERE id = $3
		`, author.Avatar168x168.URI, avatarURL, userID)
		if err != nil {
			return 0, err
		}
	}

	if len(author.Avatar300x300.URI) > 0 {
		var avatarURL string
		if len(author.Avatar300x300.URLList) > 0 {
			avatarURL = author.Avatar300x300.URLList[0]
		}
		_, err = tx.Exec(`
			UPDATE users SET avatar_300x300_uri = $1, avatar_300x300_url = $2
			WHERE id = $3
		`, author.Avatar300x300.URI, avatarURL, userID)
		if err != nil {
			return 0, err
		}
	}

	// 插入封面URL
	for _, cover := range author.CoverURL {
		var coverURL string
		if len(cover.URLList) > 0 {
			coverURL = cover.URLList[0]
		}
		_, err = tx.Exec(`
			INSERT INTO cover_urls (user_id, uri, url, type)
			VALUES ($1, $2, $3, $4)
		`, userID, cover.URI, coverURL, "cover")
		if err != nil {
			return 0, err
		}
	}

	// 插入白色封面URL
	for _, cover := range author.WhiteCoverURL {
		var coverURL string
		if len(cover.URLList) > 0 {
			coverURL = cover.URLList[0]
		}
		_, err = tx.Exec(`
			INSERT INTO cover_urls (user_id, uri, url, type)
			VALUES ($1, $2, $3, $4)
		`, userID, cover.URI, coverURL, "white_cover")
		if err != nil {
			return 0, err
		}
	}

	return userID, nil
}

// 导入音乐数据
func importMusic(tx *sql.Tx, music Music) (int64, error) {
	var musicID int64
	var musicIDStr string

	// 将音乐ID转换为字符串
	switch v := music.ID.(type) {
	case float64:
		musicID = int64(v)
		musicIDStr = fmt.Sprintf("%d", musicID)
	case string:
		musicIDStr = v
		// 尝试将字符串转换为int64
		fmt.Sscanf(musicIDStr, "%d", &musicID)
	default:
		// 生成一个随机ID
		musicID = 0
		musicIDStr = "0"
	}

	// 检查音乐是否已存在
	var existingID int64
	err := tx.QueryRow("SELECT id FROM music WHERE id = $1", musicID).Scan(&existingID)
	if err == nil {
		// 音乐已存在，返回音乐ID
		return existingID, nil
	} else if err != sql.ErrNoRows {
		// 发生其他错误
		return 0, err
	}

	// 获取封面URL
	var coverURI, coverURL string
	if len(music.CoverMedium.URI) > 0 {
		coverURI = music.CoverMedium.URI
		if len(music.CoverMedium.URLList) > 0 {
			coverURL = music.CoverMedium.URLList[0]
		}
	}

	// 获取播放URL
	var playURL string
	if len(music.PlayURL.URLList) > 0 {
		playURL = music.PlayURL.URLList[0]
	}

	// 插入音乐数据
	_, err = tx.Exec(`
		INSERT INTO music (
			id, title, author, cover_uri, cover_url, play_url,
			duration, owner_id, owner_nickname, is_original
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, musicID, music.Title, music.Author, coverURI, coverURL, playURL,
		music.Duration, music.OwnerID, music.OwnerNickname, music.IsOriginal)
	if err != nil {
		return 0, err
	}

	return musicID, nil
}

// 导入视频数据
func importVideo(tx *sql.Tx, video Video, musicID int64) (int, error) {
	var videoID int

	// 检查视频是否已存在
	err := tx.QueryRow("SELECT id FROM videos WHERE aweme_id = $1", video.AwemeID).Scan(&videoID)
	if err == nil {
		// 视频已存在，返回视频ID
		return videoID, nil
	} else if err != sql.ErrNoRows {
		// 发生其他错误
		return 0, err
	}

	// 获取作者用户ID
	var authorUserID string
	switch v := video.AuthorUserID.(type) {
	case string:
		authorUserID = v
	case float64:
		authorUserID = fmt.Sprintf("%d", int64(v))
	default:
		authorUserID = video.Author.UID
	}

	// 插入视频数据
	err = tx.QueryRow(`
		INSERT INTO videos (
			aweme_id, desc, create_time, music_id, author_user_id,
			duration, type, is_top, prevent_download
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`, video.AwemeID, video.Desc, video.CreateTime, musicID, authorUserID,
		video.Duration, "recommend-video", false, false).Scan(&videoID)
	if err != nil {
		return 0, err
	}

	return videoID, nil
}

// 导入视频播放地址
func importVideoPlayAddress(tx *sql.Tx, videoID int, playAddr Cover) error {
	var url string
	if len(playAddr.URLList) > 0 {
		url = playAddr.URLList[0]
	}

	_, err := tx.Exec(`
		INSERT INTO video_play_addresses (
			video_id, uri, url, width, height
		) VALUES ($1, $2, $3, $4, $5)
	`, videoID, playAddr.URI, url, playAddr.Width, playAddr.Height)
	return err
}

// 导入视频封面
func importVideoCover(tx *sql.Tx, videoID int, cover Cover) error {
	var url string
	if len(cover.URLList) > 0 {
		url = cover.URLList[0]
	}

	_, err := tx.Exec(`
		INSERT INTO video_covers (
			video_id, uri, url, width, height
		) VALUES ($1, $2, $3, $4, $5)
	`, videoID, cover.URI, url, cover.Width, cover.Height)
	return err
}

// 导入视频统计
func importVideoStatistics(tx *sql.Tx, videoID int, statistics Statistics) error {
	_, err := tx.Exec(`
		INSERT INTO video_statistics (
			video_id, comment_count, digg_count, collect_count, play_count, share_count
		) VALUES ($1, $2, $3, $4, $5, $6)
	`, videoID, statistics.CommentCount, statistics.DiggCount,
		statistics.CollectCount, statistics.PlayCount, statistics.ShareCount)
	return err
}

// 导入视频状态
func importVideoStatus(tx *sql.Tx, videoID int, status Status) error {
	_, err := tx.Exec(`
		INSERT INTO video_status (
			video_id, is_delete, allow_share, is_prohibited, in_reviewing, private_status
		) VALUES ($1, $2, $3, $4, $5, $6)
	`, videoID, status.IsDelete, status.AllowShare,
		status.IsProhibited, status.InReviewing, status.PrivateStatus)
	return err
}
