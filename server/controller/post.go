package controller

import (
	"github.com/gin-gonic/gin"
	"klik/server/model"
	"net/http"
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

	// 计算分页
	offset, limit := model.GetPageRange(params)

	// 从数据库加载帖子数据
	posts, err := model.GetRecommendPostsFromDB(offset, limit)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "加载帖子数据失败: " + err.Error(),
			Data: nil,
		})
		return
	}

	// 获取帖子总数
	total, err := model.GetPostCountFromDB()
	if err != nil {
		total = len(posts) // 如果获取总数失败，使用当前列表长度
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "",
		Data: model.PageResponse{
			PageNo: params.PageNo,
			Total:  total,
			List:   posts,
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

	// 计算分页
	offset, limit := model.GetPageRange(params)

	// 从数据库加载商品数据
	goods, err := model.GetGoodsFromDB(offset, limit)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code: 500,
			Msg:  "加载商品数据失败: " + err.Error(),
			Data: nil,
		})
		return
	}

	// 获取商品总数
	total, err := model.GetGoodCountFromDB()
	if err != nil {
		total = len(goods) // 如果获取总数失败，使用当前列表长度
	}

	// 返回数据
	c.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "",
		Data: model.PageResponse{
			PageNo: params.PageNo,
			Total:  total,
			List:   goods,
		},
	})
}
