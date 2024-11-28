package redis

import (
	"bluebell/models"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

// 根据键名和索引查询指定范围的id列表
func getIDsFormKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1

	//ZREVRANGE 按分数从大到小查询指定数量的元素
	result, err := client.ZRevRange(key, start, end).Result()
	return result, err
}

// 根据排序方式和索引范围，查询id列表
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

// 按社区查询该社区下的id列表
func GetCommunityPostIDsInOrder(p *models.CommunityPostListParam) (ids []string, err error) {
	//根据指定的排序方式，确定要操作的redis中的key
	// orderKey指定排序方式的键名，按时间排序则是KeyPostTimeZSet，按分数排序则是KeyPostScoreZSet
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}

	//从KeyCommunitySetPF中查询该社区下的帖子id列表，根据id列表去KeyPostTimeZSet或KeyPostScoreZSet中去查询时间或分数
	//也就是查询交集，将查询到的内容（帖子postID和对应的时间或分数）保存到新的自定义的key中

	//社区的key
	communityKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))
	key := orderKey + strconv.Itoa(int(p.CommunityID)) //自定义新key，用来存储两表交集的，值为postID和对应的时间或分数，表示此社区分类下的帖子和时间/分数
	if client.Exists(key).Val() < 1 {                  //判断key是否存在，如果不存在返回值为0
		//key不存在，需要计算
		pipeline := client.Pipeline()
		//通过 ZInterStore 对有序集合communityKey和orderKey进行交集运算，结果存储到key中
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX", //表示交集的分数取较大的值，如果将一个普通的set（无序集合）与一个zset（有序集合）一起参与ZINTERSTORE操作，Redis会自动将set视为一个所有成员分数为1的特殊zset
		}, communityKey, orderKey)
		pipeline.Expire(key, 60*time.Second) // 设置超时时间
		_, err = pipeline.Exec()
		if err != nil {
			return
		}
	}

	//查询指定索引范围的id列表
	return getIDsFormKey(key, p.Page, p.Size)
}
