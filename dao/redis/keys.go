package redis

// redis key 注意使用命名空间
const (
	Prefix             = "bluebell"    //项目key前缀
	KeyPostTimeZSet    = "post:time"   // zset;成员名:帖子(postID),分数:发帖时间
	KeyPostScoreZSet   = "post:score"  // zset;成员名:帖子(postID),分数:投票分数
	KeyPostVotedZSetPF = "post:voted:" //zset;记录某一个帖子投票的用户及投票类型，其后加postId构成完整键名（指定哪个帖子）
	//成员名:用户id(userID),分数:投票类型(1,-1)

	KeyCommunitySetPF = "community:" //set;保存每个分区下所有帖子的id，键名后加分区communityid构成完整键名
)

func getRedisKey(key string) string {
	return Prefix + key
}
