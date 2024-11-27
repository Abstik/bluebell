package redis

import (
	"bluebell/models"
	"github.com/go-redis/redis"
)

func GetPostIDsInOrder(p *models.PostListParam) ([]string, error) {
	//从redis中获取id
	//1.根据用户请求中携带的order参数（排序方式）确定要查询的redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}

	//2.确定查询的索引起始点
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1

	//3.ZREVRANGE 按分数从大到小查询指定数量的元素
	result, err := client.ZRevRange(key, start, end).Result()

	return result, err
}

// 根据ids列表查询每篇帖子的投赞成票的数据
func GetPostVoteDataByIDs(ids []string) (data []int64, err error) {
	pipeline := client.Pipeline()

	for _, id := range ids {
		//在KeyPostVotedZSetPF后拼接帖子的id，构成完整的key
		key := getRedisKey(KeyPostVotedZSetPF + id)
		//计算该帖子的赞成票总数
		pipeline.ZCount(key, "1", "1").Val()
	}

	cmders, err := pipeline.Exec()
	if err != nil {
		return
	}

	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// 根据id查询此帖子的赞成票数
func GetPostVoteDataByID(id string) int64 {
	key := getRedisKey(KeyPostVotedZSetPF + id)
	data := client.ZCount(key, "1", "1").Val()
	return data
}
