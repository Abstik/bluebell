package controller

import (
	"bluebell/logic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// 社区模块

// 查询到所有的的社区，以列表形式返回
func CommunityHandler(c *gin.Context) {
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("获取社区列表失败", zap.Error(err))
		return
	}
	ResponseSuccess(c, data)
}

// 社区分类查询
func CommunityDetailHandler(c *gin.Context) {
	//1.获取社区id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	//如果获取请求参数失败
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	//查询到所有的社区，以列表形式返回
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("获取社区列表失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
