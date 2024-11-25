package mysql

import (
	"bluebell/models"
	"database/sql"
	"errors"
	"go.uber.org/zap"
)

// 查询社区列表
func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "select community_id, community_name from community"
	err = db.Select(&communityList, sqlStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			zap.L().Warn("数据库中没有数据")
			err = nil
		}
	}
	return
}

// 根据id查询社区分类详情
func GetCommunityDetailById(id int64) (*models.CommunityDetail, error) {
	communityDetail := new(models.CommunityDetail)
	sqlStr := `select community_id, community_name, introduction, create_time 
			   from community where community_id = ?`

	err := db.Get(communityDetail, sqlStr, id)
	//如果查询失败，返回"无效的id"错误
	if err != nil {
		err = ErrorInvalidID
	}

	return communityDetail, err
}
