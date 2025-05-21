package controller

import (
	"github.com/gin-gonic/gin"
	"klik/server/model"
	"net/http"
	"strconv"
)

// GetRecommendedVideos 获取推荐视频
func GetRecommendedVideos(c *gin.Context) {
	// 获取参数
	start, _ := strconv.Atoi(c.Query("start"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))

	// 检查参数
	if pageSize <= 0 {
		pageSize = 10
	}

	var videos []model.Video
	var err error

	// 从数据库加载视频数据
	videos, err = model.GetRecommendVideosFromDB(start, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "加载视频数据失败: " + err.Error(),
			Data: nil,
		})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "",
		Data: model.PageResponse{
			Total: 844,
			List:  videos,
		},
	})
}

// GetLongRecommendedVideos 获取长视频推荐
func GetLongRecommendedVideos(c *gin.Context) {
	// 获取参数
	var params model.PageParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "参数错误",
			Data: nil,
		})
		return
	}

	// 检查参数
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	// 计算分页参数
	start := (params.PageNo - 1) * params.PageSize
	pageSize := params.PageSize

	// 从数据库加载视频数据
	videos, err := model.GetLongRecommendVideosFromDB(start, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "加载视频数据失败: " + err.Error(),
			Data: nil,
		})
		return
	}

	// 获取长视频总数
	total, err := model.GetVideoCountFromDB("long-video")
	if err != nil {
		// 如果获取总数失败，使用默认值
		total = 844
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "",
		Data: model.PageResponse{
			PageNo: params.PageNo,
			Total:  total,
			List:   videos,
		},
	})
}

// GetVideoComments 获取视频评论
func GetVideoComments(c *gin.Context) {
	// 获取参数
	videoID := c.Query("id")

	// 加载评论数据
	comments, err := model.GetVideoCommentsFromPostgres(videoID)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "加载评论数据失败: " + err.Error(),
			Data: nil,
		})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "",
		Data: comments,
	})
}

// GetPrivateVideos 获取私有视频
func GetPrivateVideos(c *gin.Context) {
	// 获取参数
	var params model.PageParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "参数错误",
			Data: nil,
		})
		return
	}

	// 检查参数
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	// 计算分页参数
	start := (params.PageNo - 1) * params.PageSize
	pageSize := params.PageSize

	// 从数据库加载视频数据
	videos, err := model.GetPrivateVideosFromDB(start, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "加载视频数据失败: " + err.Error(),
			Data: nil,
		})
		return
	}

	// 获取私有视频总数
	total, err := model.GetVideoCountFromDB("private-video")
	if err != nil {
		// 如果获取总数失败，使用默认值
		total = 10
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "",
		Data: model.PageResponse{
			PageNo: params.PageNo,
			Total:  total,
			List:   videos,
		},
	})
}

// GetLikedVideos 获取喜欢的视频
func GetLikedVideos(c *gin.Context) {
	// 获取参数
	var params model.PageParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "参数错误",
			Data: nil,
		})
		return
	}

	// 检查参数
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	// 计算分页参数
	start := (params.PageNo - 1) * params.PageSize
	pageSize := params.PageSize

	// 从数据库加载视频数据
	videos, err := model.GetLikedVideosFromDB(start, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "加载视频数据失败: " + err.Error(),
			Data: nil,
		})
		return
	}

	// 获取喜欢的视频总数
	// 注意：这里需要一个特殊的函数来获取当前用户喜欢的视频总数
	// 暂时使用固定值或者视频列表长度
	total := len(videos)
	if total == 0 {
		total = 150 // 默认值
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "",
		Data: model.PageResponse{
			PageNo: params.PageNo,
			Total:  total,
			List:   videos,
		},
	})
}

// GetMyVideos 获取我的视频
func GetMyVideos(c *gin.Context) {
	// 获取参数
	var params model.PageParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "参数错误",
			Data: nil,
		})
		return
	}

	// 检查参数
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	// 计算分页参数
	start := (params.PageNo - 1) * params.PageSize
	pageSize := params.PageSize

	// 从数据库加载视频数据
	videos, err := model.GetMyVideosFromDB(start, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "加载视频数据失败: " + err.Error(),
			Data: nil,
		})
		return
	}

	// 获取我的视频总数
	// 注意：这里需要一个特殊的函数来获取当前用户的视频总数
	// 暂时使用视频列表长度
	total := len(videos)

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "",
		Data: model.PageResponse{
			PageNo: params.PageNo,
			Total:  total,
			List:   videos,
		},
	})
}

// GetHistoryVideos 获取历史视频
func GetHistoryVideos(c *gin.Context) {
	// 获取参数
	var params model.PageParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "参数错误",
			Data: nil,
		})
		return
	}

	// 检查参数
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	// 计算分页参数
	start := (params.PageNo - 1) * params.PageSize
	pageSize := params.PageSize

	// 从数据库加载视频数据
	videos, err := model.GetHistoryVideosFromDB(start, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "加载视频数据失败: " + err.Error(),
			Data: nil,
		})
		return
	}

	// 获取历史视频总数
	// 注意：这里需要一个特殊的函数来获取当前用户的历史视频总数
	// 暂时使用视频列表长度或默认值
	total := len(videos)
	if total == 0 {
		total = 150 // 默认值
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "",
		Data: model.PageResponse{
			PageNo: params.PageNo,
			Total:  total,
			List:   videos,
		},
	})
}
