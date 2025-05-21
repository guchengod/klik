package router

import (
	"klik/server/controller"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()

	// 配置跨域
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 静态文件服务
	r.StaticFS("/api/file", http.Dir("public/data"))

	// API路由
	api := r.Group("/api")
	{
		// 视频相关接口
		video := api.Group("/video")
		{
			video.GET("/recommended", controller.GetRecommendedVideos)
			video.GET("/long/recommended", controller.GetLongRecommendedVideos)
			video.GET("/comments", controller.GetVideoComments)
			video.GET("/private", controller.GetPrivateVideos)
			video.GET("/like", controller.GetLikedVideos)
			video.GET("/my", controller.GetMyVideos)
			video.GET("/history", controller.GetHistoryVideos)
		}

		// 用户相关接口
		user := api.Group("/user")
		{
			user.GET("/collect", controller.GetUserCollect)
			user.GET("/video_list", controller.GetUserVideoList)
			user.GET("/panel", controller.GetUserPanel)
			user.GET("/friends", controller.GetUserFriends)
		}

		// 帖子相关接口
		post := api.Group("/post")
		{
			post.GET("/recommended", controller.GetRecommendedPosts)
		}

		// 商店相关接口
		shop := api.Group("/shop")
		{
			shop.GET("/recommended", controller.GetRecommendedShop)
		}
	}

	return r
}
