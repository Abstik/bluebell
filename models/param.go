package models

// 定义请求的参数结构体（DTO）

// SignUpParam 用户注册请求参数
type SignUpParam struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"` // 该字段必须有并且要和Password字段相等
}

// LoginParam 用户登录请求参数
type LoginParam struct {
	Username string `json:"username" bind:"required"`
	Password string `json:"password" bind:"required"`
}
