package service

import (
	"server/global"
	"strconv"
)

func (articleService *ArticleService) NewArticleView() CountDB {
	return CountDB{
		Index: "article_views",
	}
}

type CountDB struct {
	Index string
}

// Set 在原有基础上加一
func (c CountDB) Set(id string) error {
	err := global.Redis.HIncrBy(c.Index, id, 1).Err()
	return err
}

// GetInfo 取出数据
func (c CountDB) GetInfo() map[string]int {
	var Info = map[string]int{}
	maps := global.Redis.HGetAll(c.Index).Val()
	for id, val := range maps {
		num, _ := strconv.Atoi(val)
		Info[id] = num
	}
	return Info
}

// Clear 清除数据
func (c CountDB) Clear() {
	global.Redis.Del(c.Index)
}

// ClearOne 只清楚指定文章的计数
// 同步成功一篇就摘掉一篇，避免整表清空导致数据丢失
func (c CountDB) ClearOne(id string) error {
	return global.Redis.HDel(c.Index, id).Err()
}
