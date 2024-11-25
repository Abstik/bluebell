package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

// 把每一步数据库操作封装成函数
// 待logic层根据业务需求调用

// 对密码进行加密的盐
const secret = "Ji Bowen"

// 根据用户名查询用户是否存在
func CheckUserExist(username string) error {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int64
	err := db.Get(&count, sqlStr, username)
	if err != nil {
		return err
	}
	//如果用户已存在，返回错误
	if count > 0 {
		return ErrorUserExist
	}
	return nil
}

// 新增用户
func InsertUser(user *models.User) error {
	//对密码进行加密
	password := encryptPassword(user.Password)

	//执行SQL语句入库
	sqlStr := `insert into user (user_id, username, password) values (?,?,?)`
	_, err := db.Exec(sqlStr, user.UserID, user.Username, password)
	return err
}

// 对密码进行加密
func encryptPassword(oPasswoed string) string {
	h := md5.New()          // 创建一个 MD5 哈希对象
	h.Write([]byte(secret)) // 向哈希对象中写入 `secret` 的字节数据
	//把 secret 的字节数据写入到 MD5 哈希的内部状态，开始计算哈希值。 相当于让 secret 成为一个固定的输入。
	return hex.EncodeToString(h.Sum([]byte(oPasswoed)))
	// h.Sum([]byte(oPasswoed))：将 oPassword 的字节作为已有哈希值的“附加值”，生成最终的哈希
	//hex.EncodeToString：将计算出的 MD5 哈希值（16 字节）转换成一个可读的十六进制字符串，便于存储或显示。
}

// 用户登录
func Login(user *models.User) (err error) {
	oPassWord := user.Password // 用户登录密码
	//从数据库中查询用户
	sqlStr := `select user_id, username, password from user where user.username = ?`
	err = db.Get(user, sqlStr, user.Username)
	//如果没查询到用户，返回用户不存在错误
	if errors.Is(err, sql.ErrNoRows) {
		return ErrorUserNotExist
	}

	//判断密码是否正确
	password := encryptPassword(oPassWord)
	//如果密码不正确，返回密码不正确错误
	if password != user.Password {
		return ErrorInvalidPassword
	}

	return err
}

// 根据用户id查询用户
func GetUserById(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username from user where user_id = ?`
	err = db.Get(user, sqlStr, uid)
	return
}
