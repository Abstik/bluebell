package models

import "time"

type Post struct {
	ID          int64     `json:"id,string" db:"post_id"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	AuthorID    int64     `json:"author_id" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

// 帖子详情接口的结构体
type ApiPostDetail struct {
	AuthorName       string             `json:"author_name"`
	*Post            `json:"post"`      // 嵌入帖子结构体
	*CommunityDetail `json:"community"` // 嵌入社区信息
	VoteNum          int64              `json:"vote_num"`
}
