package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 用户注册
func SignUpHandler(c *gin.Context) {
	//1.获取参数和参数校验
	var p models.SignUpParam
	err := c.ShouldBindJSON(&p)
	if err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("参数校验失败", zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		/*字段级别的验证错误（例如某个字段没有满足required），
		或结构体级别的验证错误（例如 Password 和 RePassword 不一致），都会返回 validator.ValidationErrors 类型的错误。*/
		errs, ok := err.(validator.ValidationErrors)
		//如果不是，则可能是其他错误，例如：json格式不正确
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		//如果是，说明字段校验失败，调用自定义的函数进行错误信息的翻译
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	//2.业务处理
	err = logic.SingUp(&p)
	//如果出现错误
	if err != nil {
		zap.L().Error("注册失败", zap.Error(err))
		//如果是用户已存在的错误
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		//如果是其他错误，返回服务端繁忙错误信息
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回成功响应
	ResponseSuccess(c, nil)
	return
}

// 用户登录
func LoginHandler(c *gin.Context) {
	//1.获取请求参数以及参数校验
	p := new(models.LoginParam)
	err := c.ShouldBindJSON(p)
	if err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("参数校验失败", zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		//如果不是
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		//如果是
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	//2.业务逻辑处理
	token, err := logic.Login(p)
	if err != nil {
		zap.L().Error("登录失败", zap.String("uername", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) { //如果是用户不存在错误
			ResponseError(c, CodeUserNotExist)
			return
		} else if errors.Is(err, mysql.ErrorInvalidPassword) { //如果是密码不正确错误
			ResponseError(c, CodeInvalidPassword)
			return
		} else { //否则返回服务端繁忙错误
			ResponseError(c, CodeServerBusy)
			return
		}
	}

	//3.登陆成功，直接将token返回给用户
	ResponseSuccess(c, token)
	return
}
