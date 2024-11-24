package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
)

const ContextUserIDKey string = "userID"

var ErrorUserNotLogin = errors.New("用户未登录")

func GetCurrentUser(c *gin.Context) (userID int64, err error) {
	//如果用户以登录则可以在请求上下文中获取到userID, 如果获取不到则用户未登录
	uid, ok := c.Get(ContextUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}

	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}

	return
}
