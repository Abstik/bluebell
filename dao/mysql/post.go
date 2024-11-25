package mysql

import "bluebell/models"

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
	sqlStr := `select post_id, title, content, author_id, community_id, create_time from post 
                                                                     limit ?,?`
	data = make([]*models.Post, 0, 2)
	err = db.Select(&data, sqlStr, (pageNum-1)*pageSize, pageSize)
	return
}
