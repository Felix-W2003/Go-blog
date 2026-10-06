package service

import (
	"context"
	"encoding/json"
	"errors"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/elasticsearch"
	"server/model/other"
	"server/model/request"
	"server/utils"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ArticleService struct {
}

func (articleService *ArticleService) ArticleInfoByID(id string) (elasticsearch.Article, error) {
	// 异步更新浏览量
	go func() {
		articleView := articleService.NewArticleView()
		_ = articleView.Set(id)
	}()
	return articleService.Get(id)
}

func (articleService *ArticleService) ArticleSearch(info request.ArticleSearch) (interface{}, int64, error) {
	req := &search.Request{
		Query: &types.Query{},
	}

	boolQuery := &types.BoolQuery{}

	// 根据查询字段查询
	if info.Query != "" {
		boolQuery.Should = []types.Query{
			{Match: map[string]types.MatchQuery{"title": {Query: info.Query}}},
			{Match: map[string]types.MatchQuery{"keyword": {Query: info.Query}}},
			{Match: map[string]types.MatchQuery{"abstract": {Query: info.Query}}},
			{Match: map[string]types.MatchQuery{"content": {Query: info.Query}}},
		}
	}

	// 根据标签筛选
	if info.Tag != "" {
		boolQuery.Must = []types.Query{
			{Match: map[string]types.MatchQuery{"tags": {Query: info.Tag}}},
		}
	}

	// 根据类别筛选
	if info.Category != "" {
		boolQuery.Filter = []types.Query{
			{Term: map[string]types.TermQuery{"category": {Value: info.Category}}},
		}
	}

	// 如果有查询条件，则使用 Bool 查询，否则使用 MatchAll 查询
	if boolQuery.Should != nil || boolQuery.Must != nil || boolQuery.Filter != nil {
		req.Query.Bool = boolQuery
	} else {
		req.Query.MatchAll = &types.MatchAllQuery{}
	}

	// 设置排序字段
	if info.Sort != "" {
		var sortField string
		switch info.Sort {
		case "time":
			sortField = "created_at"
		case "view":
			sortField = "views"
		case "comment":
			sortField = "comments"
		case "like":
			sortField = "likes"
		default:
			sortField = "created_at"
		}

		var order sortorder.SortOrder
		if info.Order != "asc" {
			order = sortorder.Desc
		} else {
			order = sortorder.Asc
		}

		req.Sort = []types.SortCombinations{
			types.SortOptions{
				SortOptions: map[string]types.FieldSort{
					sortField: {Order: &order},
				},
			},
		}
	}

	option := other.EsOption{
		PageInfo:       info.PageInfo,
		Index:          elasticsearch.ArticleIndex(),
		Request:        req,
		SourceIncludes: []string{"created_at", "cover", "title", "abstract", "category", "tags", "views", "comments", "likes"},
	}
	return utils.EsPagination(context.TODO(), option)
}

// hotArticleLimit 首页「热门文章」展示的条数
const hotArticleLimit = 10

// hotArticleCacheKey 热门文章在 Redis 中的缓存 key
const hotArticleCacheKey = "article:hot"

// hotArticleCacheTTL 热门文章缓存的有效期
//
// 定时任务每小时刷新一次，TTL 特意设成刷新周期的两倍：
// 正常情况下 key 总会在过期之前被下一次刷新覆盖，首页不会出现「刚好过期、
// 只能回源 ES」的空窗；TTL 只作为定时任务失效时的兜底，防止脏数据永久驻留。
const hotArticleCacheTTL = 2 * time.Hour

// hotArticleCache 热门文章缓存的载体
// 直接缓存 ES 命中的原样结构，可以原封不动地回给前端，无需改动接口契约
type hotArticleCache struct {
	List  []types.Hit `json:"list"`
	Total int64       `json:"total"`
}

// queryHotArticles 从 Elasticsearch 查询浏览量最高的前 hotArticleLimit 篇文章
func (articleService *ArticleService) queryHotArticles() ([]types.Hit, int64, error) {
	req := &search.Request{
		Query: &types.Query{
			MatchAll: &types.MatchAllQuery{},
		},
		// 按浏览量降序排列
		Sort: []types.SortCombinations{
			types.SortOptions{
				SortOptions: map[string]types.FieldSort{
					"views": {Order: &sortorder.Desc},
				},
			},
		},
	}

	option := other.EsOption{
		PageInfo: request.PageInfo{
			Page:     1,
			PageSize: hotArticleLimit,
		},
		Index:   elasticsearch.ArticleIndex(),
		Request: req,
		// 首页小卡片只用到标题和浏览量，不回传摘要与正文，减小响应体
		SourceIncludes: []string{"title", "views", "created_at"},
	}
	return utils.EsPagination(context.TODO(), option)
}

// getHotArticleFromCache 读缓存，第二个返回值表示是否命中
// Redis 报错或数据损坏都只记日志并返回未命中，交给上层回源 ES，不阻断首页
func (articleService *ArticleService) getHotArticleFromCache() ([]types.Hit, int64, bool) {
	cached, err := global.Redis.Get(hotArticleCacheKey).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			global.Log.Error("Failed to read hot article cache:", zap.Error(err))
		}
		return nil, 0, false
	}

	var cache hotArticleCache
	if err := json.Unmarshal([]byte(cached), &cache); err != nil {
		global.Log.Error("Failed to unmarshal hot article cache:", zap.Error(err))
		return nil, 0, false
	}

	return cache.List, cache.Total, true
}

// setHotArticleCache 覆盖写缓存
func (articleService *ArticleService) setHotArticleCache(list []types.Hit, total int64) error {
	if len(list) == 0 {
		// 空榜单也是有效结果，照样缓存，避免每次请求都打到 ES
		global.Log.Warn("Hot article query returned an empty list, caching it anyway")
	}

	payload, err := json.Marshal(hotArticleCache{List: list, Total: total})
	if err != nil {
		return err
	}

	return global.Redis.Set(hotArticleCacheKey, payload, hotArticleCacheTTL).Err()
}

// RefreshHotArticleCache 重新查询 ES 并覆盖缓存，供定时任务调用
// 查询失败时直接返回错误、保留旧缓存 —— 宁可返回一小时前的榜单，也不要把它清空
func (articleService *ArticleService) RefreshHotArticleCache() error {
	list, total, err := articleService.queryHotArticles()
	if err != nil {
		return err
	}

	return articleService.setHotArticleCache(list, total)
}

// ArticleHot 获取首页热门文章：优先读 Redis 缓存，未命中再回源 ES 并顺手写回
//
// 定期刷新由 task.UpdateHotArticlesCacheTask 负责，这里的回源是为了自愈：
// 服务刚启动、Redis 被清空、或定时任务失效时，第一个请求会重建缓存。
func (articleService *ArticleService) ArticleHot() (interface{}, int64, error) {
	if list, total, ok := articleService.getHotArticleFromCache(); ok {
		return list, total, nil
	}

	list, total, err := articleService.queryHotArticles()
	if err != nil {
		return nil, 0, err
	}

	if err := articleService.setHotArticleCache(list, total); err != nil {
		// 写缓存失败不影响本次返回
		global.Log.Error("Failed to warm hot article cache:", zap.Error(err))
	}

	return list, total, nil
}

func (articleService *ArticleService) ArticleCategory() ([]database.ArticleCategory, error) {
	var category []database.ArticleCategory
	if err := global.DB.Find(&category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (articleService *ArticleService) ArticleTags() ([]database.ArticleTag, error) {
	var tags []database.ArticleTag
	if err := global.DB.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (articleService *ArticleService) ArticleLike(req request.ArticleLike) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		var al database.ArticleLike
		var num int

		// 如果用户未收藏，则创建收藏记录
		if errors.Is(tx.Where("user_id = ? AND article_id = ?", req.UserID, req.ArticleID).First(&al).Error, gorm.ErrRecordNotFound) {
			if err := tx.Create(&database.ArticleLike{UserID: req.UserID, ArticleID: req.ArticleID}).Error; err != nil {
				return err
			}
			num = 1
		} else { // 如果用户已经收藏，则取消收藏
			if err := tx.Delete(&al).Error; err != nil {
				return err
			}
			num = -1
		}

		// 更新文章收藏数
		source := "ctx._source.likes += " + strconv.Itoa(num)
		script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
		_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), req.ArticleID).Script(&script).Do(context.TODO())
		return err
	})
}

func (articleService *ArticleService) ArticleIsLike(req request.ArticleLike) (bool, error) {
	return !errors.Is(global.DB.Where("user_id = ? AND article_id = ?", req.UserID, req.ArticleID).First(&database.ArticleLike{}).Error, gorm.ErrRecordNotFound), nil
}

func (articleService *ArticleService) ArticleLikesList(info request.ArticleLikesList) (interface{}, int64, error) {
	db := global.DB.Where("user_id = ?", info.UserID)
	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}

	l, total, err := utils.MySQLPagination(&database.ArticleLike{}, option)
	if err != nil {
		return nil, 0, err
	}
	var list []struct {
		Id_     string                `json:"_id"`
		Source_ elasticsearch.Article `json:"_source"`
	}

	for _, articleLike := range l {
		article, err := articleService.Get(articleLike.ArticleID)
		if err != nil {
			return nil, 0, err
		}
		article.UpdatedAt = ""
		article.Keyword = ""
		article.Content = ""
		list = append(list, struct {
			Id_     string                `json:"_id"`
			Source_ elasticsearch.Article `json:"_source"`
		}{
			Id_:     articleLike.ArticleID,
			Source_: article,
		})
	}
	return list, total, nil
}

func (articleService *ArticleService) ArticleCreate(req request.ArticleCreate) error {
	b, err := articleService.Exits(req.Title)
	if err != nil {
		return err
	}
	if b {
		return errors.New("the article already exists")
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	articleToCreate := elasticsearch.Article{
		CreatedAt: now,
		UpdatedAt: now,
		Cover:     req.Cover,
		Title:     req.Title,
		Keyword:   req.Title,
		Category:  req.Category,
		Tags:      req.Tags,
		Abstract:  req.Abstract,
		Content:   req.Content,
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 同时更新文章类别表中的数据
		if err := articleService.UpdateCategoryCount(tx, "", articleToCreate.Category); err != nil {
			return err
		}

		// 同时更新文章标签表中的数据
		if err := articleService.UpdateTagsCount(tx, []string{}, articleToCreate.Tags); err != nil {
			return err
		}

		// 同时更新图片表中的图片类别
		if err := utils.ChangeImagesCategory(tx, []string{articleToCreate.Cover}, appTypes.Cover); err != nil {
			return err
		}
		illustrations, err := utils.FindIllustrations(articleToCreate.Content)
		if err != nil {
			return err
		}
		if err := utils.ChangeImagesCategory(tx, illustrations, appTypes.Illustration); err != nil {
			return err
		}

		return articleService.Create(&articleToCreate)
	})
}

func (articleService *ArticleService) ArticleDelete(req request.ArticleDelete) error {
	if len(req.IDs) == 0 {
		return nil
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		commentService := new(CommentService)
		for _, id := range req.IDs {
			articleToDelete, err := articleService.Get(id)
			if err != nil {
				return err
			}
			// 同时更新文章类别表中的数据
			if err := articleService.UpdateCategoryCount(tx, articleToDelete.Category, ""); err != nil {
				return err
			}
			// 同时更新文章标签表中的数据
			if err := articleService.UpdateTagsCount(tx, articleToDelete.Tags, []string{}); err != nil {
				return err
			}
			// 同时更新图片表中的图片类别
			if err := utils.InitImagesCategory(tx, []string{articleToDelete.Cover}); err != nil {
				return err
			}
			illustrations, err := utils.FindIllustrations(articleToDelete.Content)
			if err != nil {
				return err
			}
			if err := utils.InitImagesCategory(tx, illustrations); err != nil {
				return err
			}
			// 同时删除该文章下的所有评论
			comments, err := commentService.CommentInfoByArticleID(request.CommentInfoByArticleID{ArticleID: id})
			// if err != nil {
			// 	return err
			// }
			// // 同时删除该文章下的所有评论
			// comments, err := ServiceGroupApp.CommentService.CommentInfoByArticleID(request.CommentInfoByArticleID{ArticleID: id})
			if err != nil {
				return err
			}
			for _, comment := range comments {
				if err := ServiceGroupApp.CommentService.DeleteCommentAndChildren(tx, comment.ID); err != nil {
					return err
				}
			}
		}

		return articleService.Delete(req.IDs)
	})
}

func (articleService *ArticleService) ArticleUpdate(req request.ArticleUpdate) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	articleToUpdate := struct {
		UpdatedAt string   `json:"updated_at"`
		Cover     string   `json:"cover"`
		Title     string   `json:"title"`
		Keyword   string   `json:"keyword"`
		Category  string   `json:"category"`
		Tags      []string `json:"tags"`
		Abstract  string   `json:"abstract"`
		Content   string   `json:"content"`
	}{
		UpdatedAt: now,
		Cover:     req.Cover,
		Title:     req.Title,
		Keyword:   req.Title,
		Category:  req.Category,
		Tags:      req.Tags,
		Abstract:  req.Abstract,
		Content:   req.Content,
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		oldArticle, err := articleService.Get(req.ID)
		if err != nil {
			return err
		}

		// 同时更新文章类别表中的数据
		if err := articleService.UpdateCategoryCount(tx, oldArticle.Category, articleToUpdate.Category); err != nil {
			return err
		}

		// 同时更新文章标签表中的数据
		if err := articleService.UpdateTagsCount(tx, oldArticle.Tags, articleToUpdate.Tags); err != nil {
			return err
		}

		// 同时更新图片表中的图片类别
		if articleToUpdate.Cover != oldArticle.Cover {
			if err := utils.InitImagesCategory(tx, []string{oldArticle.Cover}); err != nil {
				return err
			}
			if err := utils.ChangeImagesCategory(tx, []string{articleToUpdate.Cover}, appTypes.Cover); err != nil {
				return err
			}
		}
		oldIllustrations, err := utils.FindIllustrations(oldArticle.Content)
		if err != nil {
			return err
		}
		newIllustrations, err := utils.FindIllustrations(articleToUpdate.Content)
		if err != nil {
			return err
		}
		addedIllustrations, removedIllustrations := utils.DiffArrays(oldIllustrations, newIllustrations)
		if err := utils.InitImagesCategory(tx, removedIllustrations); err != nil {
			return err
		}
		if err := utils.ChangeImagesCategory(tx, addedIllustrations, appTypes.Illustration); err != nil {
			return err
		}

		return articleService.Update(req.ID, articleToUpdate)
	})
}

func (articleService *ArticleService) ArticleList(info request.ArticleList) (list interface{}, total int64, err error) {
	req := &search.Request{
		Query: &types.Query{},
	}

	boolQuery := &types.BoolQuery{}

	// 根据标题查询
	if info.Title != nil {
		boolQuery.Must = append(boolQuery.Must, types.Query{Match: map[string]types.MatchQuery{"title": {Query: *info.Title}}})
	}

	// 根据简介查询
	if info.Abstract != nil {
		boolQuery.Must = append(boolQuery.Must, types.Query{Match: map[string]types.MatchQuery{"abstract": {Query: *info.Abstract}}})
	}

	// 根据类别筛选
	if info.Category != nil {
		boolQuery.Filter = []types.Query{
			{
				Term: map[string]types.TermQuery{
					"category": {Value: info.Category},
				},
			},
		}
	}

	// 根据条件执行查询
	if boolQuery.Must != nil || boolQuery.Filter != nil {
		req.Query.Bool = boolQuery
	} else {
		req.Query.MatchAll = &types.MatchAllQuery{}
		req.Sort = []types.SortCombinations{
			types.SortOptions{
				SortOptions: map[string]types.FieldSort{
					"created_at": {Order: &sortorder.Desc},
				},
			},
		}
	}

	option := other.EsOption{
		PageInfo: info.PageInfo,
		Index:    elasticsearch.ArticleIndex(),
		Request:  req,
	}
	return utils.EsPagination(context.TODO(), option)
}
