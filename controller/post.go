package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// 创建帖子
func CreatePostHandler(c *gin.Context) {
	//1.获取参数及参数的校验
	p := new(models.Post)
	//将参数绑定到p中
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("请求参数错误", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//在请求上下文中获取userID
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID

	//2.创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, CodeSuccess)
}

// 获取帖子详情
func GetPostDetailHandler(c *gin.Context) {
	//1.获取参数（从URL中获取帖子的id）
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("获取帖子详情的参数不正确", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	//2.根据id取出帖子数据（查数据库）
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("获取帖子详情失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, data)
}

// 获取帖子列表
func GetPostListHandler(c *gin.Context) {
	//获取分页参数
	//page表示第几页，size表示每页几条数据
	pageNumStr := c.Query("page")
	pageSizeStr := c.Query("size")

	var (
		pageNum  int64
		pageSize int64
		err      error
	)

	pageNum, err = strconv.ParseInt(pageNumStr, 10, 64)
	if err != nil {
		pageNum = 0
	}
	pageSize, err = strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		pageSize = 0
	}

	//获取数据
	data, err := logic.GetPostList(pageNum, pageSize)
	if err != nil {
		zap.L().Error("获取帖子列表失败", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
