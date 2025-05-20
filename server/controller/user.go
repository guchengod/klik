package controller

import (
	"klik/server/model"
	"klik/server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserCollect 获取用户收藏
func GetUserCollect(c *gin.Context) {
	// 加载视频数据
	videos, err := utils.LoadRecommendVideos()
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "加载视频数据失败",
			Data: nil,
		})
		return
	}

	// 加载资源数据
	resource, err := utils.LoadResource()
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "加载资源数据失败",
			Data: nil,
		})
		return
	}

	// 获取音乐数据
	var music []model.Music
	if musicData, ok := resource["music"].([]interface{}); ok {
		for _, item := range musicData {
			if m, ok := item.(map[string]interface{}); ok {
				music = append(music, model.Music{
					ID:     m["id"].(string),
					Title:  m["title"].(string),
					Artist: m["artist"].(string),
					Cover:  m["cover"].(string),
				})
			}
		}
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "",
		Data: model.CollectResponse{
			Video: model.PageResponse{
				Total: 50,
				List:  videos[350:400],
			},
			Music: model.PageResponse{
				Total: len(music),
				List:  music,
			},
		},
	})
}

// GetUserVideoList 获取用户视频列表
func GetUserVideoList(c *gin.Context) {
	// 获取参数
	userID := c.Query("id")

	// 加载用户视频列表
	videos, err := utils.LoadUserVideoList(userID)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
		})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Data: videos,
	})
}

// GetUserPanel 获取用户面板信息
func GetUserPanel(c *gin.Context) {
	// 加载用户数据
	users, err := utils.LoadUsers()
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
		})
		return
	}

	// 查找指定用户
	var user model.User
	for _, u := range users {
		if u.UID == "2739632844317827" {
			user = u
			break
		}
	}

	// 返回数据
	if user.UID != "" {
		c.JSON(http.StatusOK, model.Response{
			Code: 200,
			Data: user,
		})
	} else {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
		})
	}
}

// GetUserFriends 获取用户好友
func GetUserFriends(c *gin.Context) {
	// 加载用户数据
	users, err := utils.LoadUsers()
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
		})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Data: users,
	})
}
