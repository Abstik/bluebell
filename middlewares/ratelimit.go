package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

// 利用令牌桶进行限流
func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	//根据放令牌速度和令牌同容量创建令牌桶，2秒放一个令牌，总容量最多为1
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(c *gin.Context) {
		// 如果取不到令牌就中断本次请求返回 rate limit...
		if bucket.TakeAvailable(1) == 0 {
			c.String(http.StatusOK, "rate limit...")
			c.Abort()
			return
		}
		//取到令牌就放行
		c.Next()
	}
}
