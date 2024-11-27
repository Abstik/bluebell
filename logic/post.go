package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"go.uber.org/zap"
	"strconv"
)

// 创建帖子
func CreatePost(p *models.Post) (err error) {
	//1.生成post id
	p.ID = snowflake.GenID()
	//2.保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return
	}
	//3.保存到redis
	err = redis.CreatePost(p.ID)
	return
}

// 查询帖子详情
func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {
	//查询帖子详情
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById failed", zap.Error(err))
		return
	}

	//根据作者id查询作者信息
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById failed", zap.Error(err))
		return
	}

	//根据社区id查询社区详情
	community, err := mysql.GetCommunityDetailById(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailById failed", zap.Error(err))
		return
	}

	//封装查询到的信息
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
		VoteNum:         redis.GetPostVoteDataByID(strconv.Itoa(int((post.ID)))),
	}

	return
}

// 查询帖子列表
func GetPostList(pageNum, pageSize int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(pageNum, pageSize)
	if err != nil {
		return
	}

	data = make([]*models.ApiPostDetail, 0, len(posts))

	for _, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed", zap.Error(err))
			continue
		}

		//根据社区id查询社区详情
		community, err := mysql.GetCommunityDetailById(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailById failed", zap.Error(err))
			continue
		}

		//封装查询到的信息
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
			VoteNum:         redis.GetPostVoteDataByID(strconv.Itoa(int((post.ID)))),
		}

		data = append(data, postDetail)
	}
	return
}

func GetPostList2(p *models.PostListParam) (data []*models.ApiPostDetail, err error) {
	//1.从redis中，根据指定的排序方式和查询数量，查询符合条件的id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		return
	}

	//2.根据id列表去数据库查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	//根据帖子id列表查询所有帖子的赞成票数
	voteData, err := redis.GetPostVoteDataByIDs(ids)
	if err != nil {
		return
	}

	//将帖子作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed", zap.Error(err))
			continue
		}

		//根据社区id查询社区详情
		community, err := mysql.GetCommunityDetailById(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailById failed", zap.Error(err))
			continue
		}

		//封装查询到的信息
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
			VoteNum:         voteData[idx],
		}

		data = append(data, postDetail)
	}
	return
}
