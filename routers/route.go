package routers

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

func maxArea(height []int) int {
	left, right, answer := 0, len(height)-1, 0

	for left < right {
		if height[left] <= height[right] {
			answer = max(height[left], answer)
			left++
		} else {
			answer = max(height[right], answer)
			right--
		}
	}

	return answer
}

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) //gin设置成发布模式
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	//注册业务路由
	v1 := r.Group("/api/v1")
	//用户注册
	v1.POST("/signup", controller.SignUpHandler)
	//用户登录
	v1.POST("/login", controller.LoginHandler)

	//应用JWT认证中间件
	v1.Use(middlewares.JWTAuthMiddleware())

	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		v1.GET("/posts", controller.GetPostListHandler)

		v1.POST("/vote", controller.PostVoteController)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
