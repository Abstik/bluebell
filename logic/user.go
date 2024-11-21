package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
)

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

func Login(p *models.LoginParam) error {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	err := mysql.Login(user)
	return err
}
