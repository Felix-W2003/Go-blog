package task

import (
	"context"
	"errors"
	"fmt"
	"server/global"
	"server/model/elasticsearch"
	"server/service"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"go.uber.org/zap"
)

// UpdateArticleViewsSyncTask 将 Redis 中的文章浏览量（增量）同步到 Elasticsearch
func UpdateArticleViewsSyncTask() error {
	articleView := service.ServiceGroupApp.ArticleService.NewArticleView()

	viewsInfo := articleView.GetInfo()

	var errs []error
	for id, num := range viewsInfo {
		if num == 0 {
			continue
		}

		source := "ctx._source.views += " + strconv.Itoa(num)
		script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}

		if _, err := global.ESClient.Update(elasticsearch.ArticleIndex(), id).
			Script(&script).Do(context.TODO()); err != nil {
			// 单篇失败不影响其它文章，累积下来最后统一返回
			errs = append(errs, fmt.Errorf("同步文章 %s 的浏览量(+%d)失败: %w", id, num, err))
			continue
		}

		// 只有真正写进 ES 的条目才从 Redis 摘除，失败的留到下一轮重试
		if err := articleView.ClearOne(id); err != nil {
			errs = append(errs, fmt.Errorf("清理文章 %s 的浏览量缓存失败: %w", id, err))
		}
	}

	if len(errs) == 0 {
		return nil
	}
	global.Log.Error("文章浏览量同步存在失败项", zap.Int("failed", len(errs)))
	return errors.Join(errs...)
}
