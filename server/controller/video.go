package controller

import (
	"klik/server/model"
	"klik/server/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetRecommendedVideos 获取推荐视频
func GetRecommendedVideos(c *gin.Context) {
	// 获取参数
	start, _ := strconv.Atoi(c.Query("start"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))

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

	// 计算分页
	end := start + pageSize
	if end > len(videos) {
		end = len(videos)
	}
	if start > len(videos) {
		start = 0
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "",
		Data: model.PageResponse{
			Total: 844,
			List:  videos[start:end],
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

	// 计算分页
	offset, limit := model.GetPageRange(params)
	if offset > len(videos) {
		offset = 0
	}
	if limit > len(videos) {
		limit = len(videos)
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "",
		Data: model.PageResponse{
			Total: 844,
			List:  videos[offset:limit],
		},
	})
}

// GetVideoComments 获取视频评论
func GetVideoComments(c *gin.Context) {
	// 获取参数
	videoID := c.Query("id")

	// 加载评论数据
	comments, err := utils.LoadVideoComments(videoID)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Data: nil,
		})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
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

	// 计算分页
	offset, limit := model.GetPageRange(params)
	start := 100
	end := 110
	if start+offset > len(videos) {
		offset = 0
	}
	if start+limit > len(videos) {
		limit = len(videos) - start
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "",
		Data: model.PageResponse{
			Total: 10,
			List:  videos[start:end][offset:limit],
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

	// 计算分页
	offset, limit := model.GetPageRange(params)
	start := 200
	end := 350
	if start+offset > len(videos) {
		offset = 0
	}
	if start+limit > len(videos) {
		limit = len(videos) - start
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "",
		Data: model.PageResponse{
			Total: 150,
			List:  videos[start:end][offset:limit],
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

	// 加载视频数据
	videos, err := utils.LoadUserVideos()
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "加载视频数据失败",
			Data: nil,
		})
		return
	}

	// 计算分页
	offset, limit := model.GetPageRange(params)
	if offset > len(videos) {
		offset = 0
	}
	if limit > len(videos) {
		limit = len(videos)
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "",
		Data: model.PageResponse{
			PageNo: params.PageNo,
			Total:  len(videos),
			List:   videos[offset:limit],
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

	// 计算分页
	offset, limit := model.GetPageRange(params)
	start := 200
	end := 350
	if start+offset > len(videos) {
		offset = 0
	}
	if start+limit > len(videos) {
		limit = len(videos) - start
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "",
		Data: model.PageResponse{
			Total: 150,
			List:  videos[start:end][offset:limit],
		},
	})
}

// GetHistoryOther 获取其他历史记录
func GetHistoryOther(c *gin.Context) {
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

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "",
		Data: model.PageResponse{
			PageNo: params.PageNo,
			Total:  0,
			List:   []interface{}{},
		},
	})
}
