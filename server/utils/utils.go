package utils

import (
	"encoding/json"
	"io/ioutil"
	"klik/server/config"
	"klik/server/model"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

var (
	// 推荐视频列表
	recommendVideos []model.Video
	// 用户视频列表
	userVideos []model.Video
	// 推荐帖子列表
	recommendPosts []model.Post
	// 资源数据
	resourceData map[string]interface{}
)

func init() {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())
}

// LoadResource 加载资源数据
func LoadResource() (map[string]interface{}, error) {
	if resourceData != nil {
		return resourceData, nil
	}

	// 加载资源文件
	resourcePath := filepath.Join(config.DataPath, "resource.js")
	data, err := ioutil.ReadFile(resourcePath)
	if err != nil {
		return nil, err
	}

	// 解析JSON数据
	err = json.Unmarshal(data, &resourceData)
	if err != nil {
		return nil, err
	}

	return resourceData, nil
}

// LoadRecommendVideos 加载推荐视频
func LoadRecommendVideos() ([]model.Video, error) {
	if len(recommendVideos) > 0 {
		return recommendVideos, nil
	}

	// 加载视频数据
	postsPath := filepath.Join(config.DataPath, "posts6.json")
	data, err := ioutil.ReadFile(postsPath)
	if err != nil {
		return nil, err
	}

	var videos []model.Video
	err = json.Unmarshal(data, &videos)
	if err != nil {
		return nil, err
	}

	// 设置视频类型
	for i := range videos {
		videos[i].Type = "recommend-video"
	}

	recommendVideos = videos
	return recommendVideos, nil
}

// LoadUserVideos 加载用户视频
func LoadUserVideos() ([]model.Video, error) {
	if len(userVideos) > 0 {
		return userVideos, nil
	}

	// 加载用户视频数据
	userVideoPath := filepath.Join(config.DataPath, "user_video_list", "user-12345xiaolaohu.md")
	data, err := ioutil.ReadFile(userVideoPath)
	if err != nil {
		return nil, err
	}

	var videos []model.Video
	err = json.Unmarshal(data, &videos)
	if err != nil {
		return nil, err
	}

	userVideos = videos
	return userVideos, nil
}

// LoadRecommendPosts 加载推荐帖子
func LoadRecommendPosts() ([]model.Post, error) {
	if len(recommendPosts) > 0 {
		return recommendPosts, nil
	}

	// 加载帖子数据
	postsPath := filepath.Join(config.DataPath, "posts.md")
	data, err := ioutil.ReadFile(postsPath)
	if err != nil {
		return nil, err
	}

	var posts []model.Post
	err = json.Unmarshal(data, &posts)
	if err != nil {
		return nil, err
	}

	recommendPosts = posts
	return recommendPosts, nil
}

// LoadUsers 加载用户数据
func LoadUsers() ([]model.User, error) {
	// 加载用户数据
	usersPath := filepath.Join(config.DataPath, "users.md")
	data, err := ioutil.ReadFile(usersPath)
	if err != nil {
		return nil, err
	}

	var users []model.User
	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// LoadVideoComments 加载视频评论
func LoadVideoComments(videoID string) ([]model.Comment, error) {
	// 视频ID列表
	videoIDs := []string{
		"7260749400622894336",
		"7128686458763889956",
		"7293100687989148943",
		"6923214072347512068",
		"7005490661592026405",
		"7161000281575148800",
		"7267478481213181238",
		"6686589698707590411",
		"7321200290739326262",
		"7194815099381484860",
		"6826943630775831812",
		"7110263965858549003",
		"7295697246132227343",
		"7270431418822446370",
		"6882368275695586568",
		"7000587983069957383",
	}

	// 如果视频ID不在列表中，随机选择一个
	found := false
	for _, id := range videoIDs {
		if id == videoID {
			found = true
			break
		}
	}
	if !found {
		videoID = videoIDs[rand.Intn(len(videoIDs))]
	}

	// 加载评论数据
	commentsPath := filepath.Join(config.DataPath, "comments", "video_id_"+videoID+".md")
	data, err := ioutil.ReadFile(commentsPath)
	if err != nil {
		return nil, err
	}

	var comments []model.Comment
	err = json.Unmarshal(data, &comments)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// LoadGoods 加载商品数据
func LoadGoods() ([]model.Good, error) {
	// 加载商品数据
	goodsPath := filepath.Join(config.DataPath, "goods.md")
	data, err := ioutil.ReadFile(goodsPath)
	if err != nil {
		return nil, err
	}

	var goods []model.Good
	err = json.Unmarshal(data, &goods)
	if err != nil {
		return nil, err
	}

	return goods, nil
}

// LoadUserVideoList 加载用户视频列表
func LoadUserVideoList(userID string) ([]model.Video, error) {
	// 加载用户视频列表
	userVideoPath := filepath.Join(config.DataPath, "user_video_list", "user-"+userID+".md")
	data, err := ioutil.ReadFile(userVideoPath)
	if err != nil {
		return nil, err
	}

	var videos []model.Video
	err = json.Unmarshal(data, &videos)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

// Random 生成随机数
func Random(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// CloneDeep 深度克隆
func CloneDeep(src interface{}, dst interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}
