package controller

import (
	"klik/server/model"
	"klik/server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetRecommendedPosts 获取推荐帖子
func GetRecommendedPosts(c *gin.Context) {
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

	// 加载帖子数据
	posts, err := utils.LoadRecommendPosts()
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "加载帖子数据失败",
			Data: nil,
		})
		return
	}

	// 计算分页
	offset, limit := model.GetPageRange(params)
	if offset > len(posts) {
		offset = 0
	}
	if limit > len(posts) {
		limit = len(posts)
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "",
		Data: model.PageResponse{
			PageNo: params.PageNo,
			Total:  len(posts),
			List:   posts[offset:limit],
		},
	})
}

// GetRecommendedShop 获取推荐商品
func GetRecommendedShop(c *gin.Context) {
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

	// 加载商品数据
	goods, err := utils.LoadGoods()
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "加载商品数据失败",
			Data: nil,
		})
		return
	}

	// 计算分页
	offset, limit := model.GetPageRange(params)
	if offset > len(goods) {
		offset = 0
	}
	if limit > len(goods) {
		limit = len(goods)
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Data: model.PageResponse{
			Total: len(goods),
			List:  goods[offset:limit],
		},
	})
}
