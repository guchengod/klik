package controller

import (
	"klik/server/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserCollect 获取用户收藏
func GetUserCollect(c *gin.Context) {
	// 获取用户ID，实际应该从认证信息中获取
	userID := "2739632844317827" // 模拟当前登录用户ID

	// 计算分页参数
	start := 0
	pageSize := 50

	// 从数据库获取用户收藏的视频
	videos, err := model.GetUserCollectVideosFromDB(userID, start, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "加载视频数据失败: " + err.Error(),
			Data: nil,
		})
		return
	}

	// 从数据库获取用户收藏的音乐
	music, err := model.GetUserCollectMusicFromDB(userID, start, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "加载音乐数据失败: " + err.Error(),
			Data: nil,
		})
		return
	}

	// 获取收藏视频和音乐的总数
	videoCount, err := model.GetUserCollectVideosCountFromDB(userID)
	if err != nil {
		videoCount = len(videos) // 如果获取总数失败，使用当前列表长度
	}

	musicCount, err := model.GetUserCollectMusicCountFromDB(userID)
	if err != nil {
		musicCount = len(music) // 如果获取总数失败，使用当前列表长度
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "",
		Data: model.CollectResponse{
			Video: model.PageResponse{
				Total: videoCount,
				List:  videos,
			},
			Music: model.PageResponse{
				Total: musicCount,
				List:  music,
			},
		},
	})
}

// GetUserVideoList 获取用户视频列表
func GetUserVideoList(c *gin.Context) {
	// 获取参数
	userID := c.Query("id")

	// 从数据库加载用户视频列表
	videos, err := model.GetUserVideoListFromDB(userID)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "加载用户视频列表失败: " + err.Error(),
			Data: nil,
		})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "",
		Data: videos,
	})
}

// GetUserPanel 获取用户面板信息
func GetUserPanel(c *gin.Context) {
	// 获取用户ID，实际应该从认证信息中获取
	userID := "2739632844317827" // 模拟当前登录用户ID

	// 从数据库加载用户信息
	user, err := model.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "获取用户信息失败: " + err.Error(),
			Data: nil,
		})
		return
	}

	// 返回数据
	if user.UID != "" {
		c.JSON(http.StatusOK, model.Response{
			Code: 200,
			Msg:  "",
			Data: user,
		})
	} else {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "用户不存在",
			Data: nil,
		})
	}
}

// GetUserFriends 获取用户好友
func GetUserFriends(c *gin.Context) {
	// 获取用户ID，实际应该从认证信息中获取
	userID := "2739632844317827" // 模拟当前登录用户ID

	// 从数据库加载用户好友列表
	friends, err := model.GetUserFriendsFromDB(userID)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "获取用户好友列表失败: " + err.Error(),
			Data: nil,
		})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "",
		Data: friends,
	})
}
