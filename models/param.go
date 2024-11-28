package models

// 定义请求的参数结构体（DTO）

// 指定排序方式的常量
const (
	OrderTime  = "time"
	OrderScore = "score"
)

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

// 投票数据
type VoteData struct {
	//UserID 从请求中获取当前的用户， 不用在结构体中添加
	PostID    string `json:"post_id" bind:"required"`              //帖子id
	Direction int8   `json:"direction,string" bind:"oneof=-1,0,1"` //赞成票(1)or反对票(-1)or取消投票(0) TODO这里oneof不生效
	//oneof=-1,0,1表示该字段的值要求只可能是-1或0或1
}

// 获取帖子列表的参数，get请求参数拼接在url中不是json
type PostListParam struct {
	Page  int64  `json:"page" form:"page"`   //查询第几页的数据
	Size  int64  `json:"size" form:"size"`   //每页数据条数
	Order string `json:"order" form:"order"` //排序方式
}

// 根据社区查询帖子详情
type CommunityPostListParam struct {
	PostListParam       //嵌入获取帖子列表的参数的结构体
	CommunityID   int64 `json:"community_id" form:"community_id"` //社区id
}
