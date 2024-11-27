package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"time"
)

/*
投票的情况：
direction=1时：
    1、之前没有投过票，现在投赞成票 --> 更新分数和投票记录
    2、之前投反对票，现在改投赞成票 --> 更新分数和投票记录
direction=-1时：
    1、之前没有投过票，现在投反对票 --> 更新分数和投票记录
    2、之前投赞成票，现在改投反对票 --> 更新分数和投票记录
direction=0时：
    1、之前投过赞成票，现在要取消投票 --> 更新分数和投票记录
    2、之前投反对票，现在要取消投票 --> 更新分数和投票记录
无论哪种情况，用本次投票数减去以前投票数即为此时的实际投票数

投票的限制：
	每个帖子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许投票
    1、到期之后将redis中保存的赞成票及反对票存储到mysql表中
	2、到期之后删除 KeyPostVotedZSetPF
*/

const (
	oneWeekInSeconds = 3600 * 24 * 7 //一周的秒数
	scorePerVote     = 432           //每一票的票数
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已结束")
)

// 为帖子投票
func VoteForPost(userID, postID string, value float64) error {
	//1.判断投票限制
	//利用redis获取帖子发布时间
	//ZScore函数的两个参数：键名和成员名，获取该成员的分数score
	//Val将结果转换为float64类型
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()

	//利用redis获取帖子发布时间
	//如果帖子发布时间超过一周，则不能投票
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	//2.更新帖子分数
	//查询当前用户(userID)给当前帖子(postID)的投票记录
	ov := client.ZScore(getRedisKey(KeyPostVotedZSetPF+postID), userID).Val() // 上次投票类型：1 or 0 or -1
	diff := value - ov                                                        //计算两次投票类型的差值

	//开启事务
	pipeline := client.TxPipeline()

	//给指定的键和成员名增加分数
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), diff*scorePerVote, postID)

	//3.更新用户为该帖子投票的数据
	//if value == 0 { //如果是取消投票，根据userID移除成员名(userID)和分数对
	//	pipeline.ZRem(getRedisKey(KeyPostVotedZSetPF+postID), userID)
	//} else { //如果投了赞成票或反对票，在该贴子的投票记录中增加此次投票的数据
	pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
		Score:  value, //投票类型
		Member: userID,
	})
	//}

	//执行事务
	_, err := pipeline.Exec()
	return err
}

// 新建帖子
func CreatePost(postId int64) error {
	//开启事务
	pipeline := client.TxPipeline()

	//在redis中更新帖子创建时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postId,
	})

	//在redis中更新帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postId,
	})

	_, err := pipeline.Exec()
	return err
}
