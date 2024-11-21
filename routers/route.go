package routers

import (
	"bluebell/controller"
	"bluebell/logger"
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
	//用户注册
	r.POST("/signup", controller.SignUpHandler)
	//用户登录
	r.POST("/login", controller.LoginHandler)

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
