package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

// 用户注册
func SingUp(p *models.SignUpParam) error {
	//1.判断用户存不存在
	err := mysql.CheckUserExist(p.Username)
	if err != nil {
		//数据库查询出错
		return err
	}

	//2.生成UID
	userID := snowflake.GenID()
	//构造一个user实例
	user := models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	//3.保存进数据库
	err = mysql.InsertUser(&user)
	return err
}

// 用户登录
func Login(p *models.LoginParam) (string, error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//可以从user中拿到UserID
	err := mysql.Login(user)
	if err != nil {
		return "", err
	}
	//生成JWT
	token, err := jwt.GenToken(user.UserID, user.Username)
	return token, err
}
