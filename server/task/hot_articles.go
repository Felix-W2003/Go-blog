package task

import (
	"server/global"
	"server/service"
)

// UpdateHotArticlesCacheTask 重新查询浏览量最高的文章并覆盖 Redis 缓存
//
// 目的：首页「热门文章」是访问最频繁的接口之一，而榜单数据每小时才变化一次，
// 因此把 Elasticsearch 的查询结果缓存进 Redis，让首页的读请求不再直接压到 ES。
//
// 注册顺序上排在 UpdateArticleViewsSyncTask 之后，这样读到的是本轮刚同步完的最新浏览量。
func UpdateHotArticlesCacheTask() error {
	if err := service.ServiceGroupApp.ArticleService.RefreshHotArticleCache(); err != nil {
		return err
	}

	global.Log.Info("Hot article cache refreshed")
	return nil
}
