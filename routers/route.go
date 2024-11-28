package routers

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
		//查询社区列表
		v1.GET("/community", controller.CommunityHandler)
		//查询社区分类
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		//创建帖子
		v1.POST("/post", controller.CreatePostHandler)
		//查询帖子详情
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		//查询帖子列表(分页)
		v1.GET("/posts", controller.GetPostListHandler)
		//指定顺序查询帖子列表（分页）
		v1.GET("post2", controller.GetPostListHandler2)
		//指定社区并按顺序查询帖子详情（分页）
		v1.GET("post3", controller.GetCommunityPostListHandler)

		//投票
		v1.POST("/vote", controller.PostVoteController)

	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
