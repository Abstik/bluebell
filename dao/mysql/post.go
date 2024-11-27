package mysql

import (
	"bluebell/models"
	"github.com/jmoiron/sqlx"
	"strings"
)

// 创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post (post_id, title, content, author_id, community_id) VALUES (?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// 根据帖子id查询帖子详情
func GetPostById(pid int64) (data *models.Post, err error) {
	post := new(models.Post)
	sqlStr := `select post_id, title, content, author_id, community_id, create_time from post where post_id = ?`
	err = db.Get(post, sqlStr, pid)
	return post, err
}

// 查询帖子列表
func GetPostList(pageNum, pageSize int64) (data []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time 
				from post
			    order by create_time desc 
                limit ?,?`
	data = make([]*models.Post, 0, 2)
	err = db.Select(&data, sqlStr, (pageNum-1)*pageSize, pageSize)
	return
}

// 根据给定的id列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time 
               from post
               where post.post_id in (?)
               order by FIND_IN_SET(post_id, ?)`
	//order by FIND_IN_SET(post_id, ?) 表示根据 post_id 在另一个给定字符串列表中的位置进行排序。
	//? 是另一个占位符，将被替换为一个包含多个ID的字符串，例如 "1,3,2"。

	//将传入的 sqlStr 和 ids 转换为适合 SQL 查询的格式，并生成查询语句 query 和参数 args
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return
	}

	//使用 db.Rebind 方法重新绑定查询语句，以适应不同的数据库驱动
	query = db.Rebind(query)
	//执行查询并将结果存储在 postList 变量中
	err = db.Select(&postList, query, args...)
	return
}
