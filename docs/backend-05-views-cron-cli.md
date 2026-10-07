# 浏览量削峰、定时任务与命令行工具

这个博客有两类会被反复触发的写入：每一次文章浏览、每一次首页刷新。前者如果直接落到 Elasticsearch，写代价远高于读取；后者虽然只是读，却会被每个访客重复触发。

这篇讲我怎么把这两类压力处理掉：**用 Redis Hash 聚合浏览量增量、用定时任务批量落库、用 Redis 缓存承载读多写少的榜单**，以及配套的定时任务体系和命令行运维工具。读完之后你应该能拿到几个具体的东西：

- 计数为什么必须用 `HINCRBY`，以及「读出来加一再写回去」在并发下会怎么丢数据；
- 批量同步任务里最容易写错的地方（我确实写错了，而且错了很久没被发现）；
- 缓存榜单的几个决策点：TTL 设多长、查询失败要不要清缓存、空结果要不要缓存；
- 为什么有些运维动作应该做成命令行工具而不是 HTTP 接口。

---

## 一、浏览量为什么要绕一圈 Redis

### 直接写 ES 的问题

文章正文只存在 Elasticsearch 里（MySQL 里没有 article 表），所以「浏览量 +1」最直觉的写法就是更新 ES 文档。但一次文档更新要经过解析脚本、更新内部索引结构、写 translog、等待 refresh 才能被搜到——代价远高于一次读取。首页那篇文章被点十次就是十次这样的操作。

写放大在量小的时候看不出来，但它是**随访问量线性增长、且系数远大于 1** 的东西。所以改成两段：

```
浏览时：  HINCRBY article_views <文章id> 1      （内存操作）
每小时：  读出全部增量 → 批量更新 ES → 清理已同步条目
```

把 N 次写合并成 1 次。代价是**浏览量滞后最多一小时**，对个人博客完全可以接受——但必须自己清楚，因为后面做热门榜会直接受它影响。

### 计数器类型

`server/service/article_stat.go`：

```go
type CountDB struct {
	Index string
}

// Set 在原有基础上加一
func (c CountDB) Set(id string) error {
	return global.Redis.HIncrBy(c.Index, id, 1).Err()
}

// ClearOne 只清除指定文章的计数
func (c CountDB) ClearOne(id string) error {
	return global.Redis.HDel(c.Index, id).Err()
}
```

**为什么用 Hash 而不是每篇文章一个 String。** Hash 让「所有文章的待同步增量」天然聚在一个 key 下：`GetInfo()`（内部是一次 `HGETALL`）能一次拿到全量快照，清理可以精确到 field，不需要先 `SCAN` 再拼 key 前缀。

**为什么 `Set` 用 `HIncrBy`。** 我最早的版本是读改写：`HGet` 取值、加一、`HSet` 写回。同一篇文章在两个请求之间被读了同样的值，两边各自加一写回，结果只增加了一次——**丢更新**。`HIncrBy` 是单条命令、Redis 单线程执行，天然原子，把「浏览量偶尔偏少」这类极难复现的问题从根上消掉。

调用它的地方是异步的，`server/service/article.go`：

```go
func (articleService *ArticleService) ArticleInfoByID(id string) (elasticsearch.Article, error) {
	// 异步更新浏览量
	go func() {
		articleView := articleService.NewArticleView()
		_ = articleView.Set(id)
	}()
	return articleService.Get(id)
}
```

把计数丢进 goroutine 是为了不让一次 Redis 往返拖慢文章详情。这也意味着**并发调用 `Set` 是常态**，更说明读改写方案不可行——靠 goroutine 只会让竞态窗口撞得更频繁。

**为什么清理要按条删除。** `CountDB` 里还留着一个 `Clear()` 做整表 `DEL`，但同步任务用的是 `ClearOne()`。原因是：批量同步过程中只要有一篇失败，整表清空就会把那篇**还没同步出去的增量一起抹掉**，等于永久少算。按条删除让失败条目留到下一轮重试，做到不丢不重。代价是清理变成 N 次 `HDEL`，相对 ES 的写入开销可以忽略。

（顺带一提，改成 `ClearOne` 之后 `Clear()` 已经没有任何调用点了，属于该删的历史代码。）

---

## 二、批量同步任务：循环里那行 `return`

同步任务在 `server/task/article_views.go`。ES 侧用 Painless 脚本做原子自增，脚本内容就是 `"ctx._source.views += " + strconv.Itoa(num)`，配合 `types.Script{..., Lang: &scriptlanguage.Painless}` 下发。用脚本而不是「读出来加一写回去」，理由和 Redis 那边一样：避免读改写互相覆盖。

### 我写错的那一版

```go
for id, num := range viewsInfo {
	if num == 0 {
		continue
	}
	_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), id).
		Script(&script).Do(context.TODO())
	return err          // ← 错在这里
}
articleView.Clear()
```

`return err` 写在了循环体内，带来三个连锁后果，其中第二个我完全没预料到：

1. **每小时只同步一篇文章。** 循环第一次执行就返回，其余文章的增量一直躺在 Redis 里，永远写不出去。
2. **已同步的那篇会被重复累加，浏览量虚高。** 因为 `return` 跳过了后面的 `Clear()`，Redis 里的计数根本没被清掉。下一轮 Go 的 map 遍历顺序变了，可能换一篇文章，把**同一批旧计数再加一遍**。所以这个 bug 不只是漏同步，还是多算。
3. **`Clear()` 变成死代码。** 它只在「所有计数都为 0」时才可能执行到——而那正是没有东西需要同步的情况。

第二个后果最阴险：页面上的数字看起来在涨，你不会怀疑它坏，只会怀疑「怎么这篇这么火」。

### 正确的写法

同样是 `server/task/article_views.go`：

```go
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
// 循环结束后：errs 为空就返回 nil；否则记一条 error 日志，并用 errors.Join 统一返回
```

要点四条：单篇失败不中断循环，累积后用 `errors.Join` 一次返回；`ClearOne` 放在 ES 更新**成功之后**；成功时不写错误日志，避免每小时刷噪音；**失败条目保留在 Redis 里等下一轮重试**。

这里的选择是：宁可这一小时少算几个浏览量，也不能因为一次 ES 抖动就丢掉浏览记录。计数这种东西，少算比多算安全，丢数据比少算更严重。

---

## 三、热门文章缓存：TTL 必须大于刷新周期

首页要展示浏览量最高的 10 篇。这是典型的读多写少：每个访客开首页都会触发一次查询，而榜单每小时才变。如果每次访问都去 ES 跑「按 views 降序取前 10」，等于把第一节省下的写压力从写侧搬回读侧。所以加一层 Redis 缓存。

实现都在 `server/service/article.go`，`server/task/hot_articles.go` 只是一个很薄的定时包装。

```go
const hotArticleLimit = 10
const hotArticleCacheKey = "article:hot"
const hotArticleCacheTTL = 2 * time.Hour

type hotArticleCache struct {
	List  []types.Hit `json:"list"`
	Total int64       `json:"total"`
}

func (articleService *ArticleService) RefreshHotArticleCache() error {
	list, total, err := articleService.queryHotArticles()
	if err != nil {
		return err // 查询失败，保留旧缓存
	}
	return articleService.setHotArticleCache(list, total)
}
```

### 决策一：TTL 设成刷新周期的两倍

刷新每小时一次，TTL 却设两小时，这不是笔误。如果 TTL 等于刷新周期，就会出现「缓存刚过期、下一轮定时任务还没执行」的窗口：这期间所有首页请求都未命中、回源 ES——恰好把你想要消除的读压力又放了回来，而且集中在同一个时间点。

TTL 设成两倍后，key 总会在过期之前被下一次刷新覆盖，**永远不会出现空窗**。此时 TTL 的角色就从「控制数据新鲜度」退化成「定时任务失效时的兜底」：万一 cron 挂了，脏数据最多存活两小时。

一句话：**当有外部机制主动刷新时，TTL 不该用来表达业务上的新鲜度，只该表达「出故障时愿意忍受多久的陈旧」。**

### 决策二：缓存 ES 命中的原样结构

`hotArticleCache` 装的是 `[]types.Hit`，即 ES 返回的原始命中结构，不是自定义精简结构体。好处是**接口契约不变**，`ArticleHot()` 返回的东西和以前直查 ES 一模一样，前端一行都不用改；代价是多存了 `_index`、`sort` 这类用不上的字段，10 条大约 3KB。榜单要放大到几百条时，就该换成精简结构了。

这个方案成立的前提是 `types.Hit` 能安全 JSON 往返。我确认过：它在 typed client 里**没有自定义 `MarshalJSON`**，走默认结构体标签序列化，`_source` 字段是 `json.RawMessage`，所以 `json.Marshal` 存进 Redis 再解回来，`_id` 和 `_source` 都不会变形。如果它有自定义序列化逻辑，这个方案就得重新设计。

### 决策三：查询失败时保留旧缓存

`RefreshHotArticleCache` 失败时什么都不做，把旧 key 留在那里。常见的错误写法是「先删缓存、再查库、再写回」——一旦查询失败缓存就是空的，紧接着的请求全部回源，可能直接把 ES 压垮，形成故障放大。语义必须明确：**宁可返回一小时前的榜单，也不要返回错误，更不要把缓存清空。**

### 决策四：空结果也要缓存

`setHotArticleCache` 在列表为空时**不跳过写入**，只额外记一条 warn 日志（`"Hot article query returned an empty list, caching it anyway"`）。空列表也是一种有效结果。如果因为「查出来是空的」就不写缓存，那么当 ES 里确实没有数据（或索引有问题持续返回空）时，每个请求都会去问一次 ES——这就是缓存穿透。缓存照写，但不能让这件事悄无声息地过去。

### 决策五：未命中时自愈

对外暴露的入口 `ArticleHot()` 在同一个文件里：

```go
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
```

「读缓存 → 未命中回源 ES → 顺手写回」这条回源路径不是多余的：**服务刚启动、Redis 被清空、定时任务失效**这三种情况下，第一个请求会把缓存重建起来。读缓存时也区分了错误：`redis.Nil`（key 不存在）是正常未命中，不记日志；其它错误才记 error，避免噪音。

### 一个遗留缺口

严格说这里还差一层**缓存击穿**防护。如果 `article:hot` 恰好过期，同时一批并发请求进来，它们会同时回源 ES。标准做法是用 singleflight 让同一时刻的相同 key 只查一次。个人博客的并发量可以忽略这个问题，所以暂时没做——但如果以后要补，这是第一个该补的地方。

---

## 四、定时任务体系

### 注册方式

所有任务集中在 `server/task/enter.go` 注册：

```go
func RegisterScheduledTasks(c *cron.Cron) error {
	if _, err := c.AddFunc("@hourly", func() {
		if err := UpdateArticleViewsSyncTask(); err != nil {
			global.Log.Error("Failed to update article views:", zap.Error(err))
		}
	}); err != nil {
		return err
	}
	// 必须排在浏览量同步之后：先落库最新浏览量，再据此刷新热门榜
	if _, err := c.AddFunc("@hourly", func() {
		if err := UpdateHotArticlesCacheTask(); err != nil {
			global.Log.Error("Failed to refresh hot article cache:", zap.Error(err))
		}
	}); err != nil {
		return err
	}
	// 热搜抓取 @hourly、日历数据 @daily 同理
	return nil
}
```

用 `robfig/cron/v3` 的预定义表达式（三个 `@hourly`、一个 `@daily`），而不是 `0 0 * * * *` 这类六位表达式——这几个任务的语义本来就是「每小时」「每天」，写成 cron 表达式反而多一层需要解读的东西。

### 注册顺序就是执行顺序

这是本节最需要注意的地方。对**同一个时间表达式**的任务，cron 按注册顺序执行。所以：**浏览量同步必须排在热门榜刷新之前**，否则榜单读到的永远是上一小时落库的 `views`，稳定地慢一拍。这个 bug 不报错、不崩溃，只会让你觉得「数据怎么总是对不上」。我把这条约束直接写成注册处的一行注释——代码本身表达不出「顺序有语义」，只能靠注释拦住后来的人。

### cron 日志接入 Zap

`server/initialize/cron.go` 做了一件容易被忽略的事：

```go
func InitCron() {
	c := cron.New(cron.WithLogger(NewZapLogger()))
	if err := task.RegisterScheduledTasks(c); err != nil {
		global.Log.Error("Error scheduling cron job:", zap.Error(err))
		os.Exit(1)
	}
	c.Start()
}
```

其中 `ZapLogger` 实现了 cron 的 `Logger` 接口（`Info` / `Error` 两个方法），把 cron 自己的日志转发给 Zap。意义在于：默认情况下 cron 会用自己的格式往标准输出打日志，定时任务的记录就散在两个地方，出问题要在两种格式之间来回对照；接进 Zap 之后，调度信息、任务自身的错误、Gin 的访问日志全在同一份结构化日志里，可以按 level 和字段一起检索。

另外，注册失败直接 `os.Exit(1)`：调度都注册不上，服务就没必要继续跑——**启动期失败远好过运行到半夜才发现某个任务从来没执行过**。

### 初始化顺序与两个任务的取舍

`server/main.go` 里的顺序有讲究：`flag.InitFlag()` → `initialize.InitCron()` → `core.RunServer()`。把 `InitFlag()` 放在 `InitCron()` 之前，意味着**命令行模式下根本不会启动 cron**——`InitFlag` 识别到参数、执行完命令后直接 `os.Exit(0)`（下一节细说）。所以「跑一次建表命令」不会顺带把常驻服务和调度器拉起来。

**热搜抓取**（`server/task/hot_search.go`）遍历 `baidu`、`zhihu`、`kuaishou`、`toutiao` 四个源，轮流调用 `hotSearch.NewSource(sourceStr).GetHotSearchData(30)` 并写入 Redis（TTL 一小时）。问题在于它是**任一源失败就直接 `return err`**——这和第二节犯的是同一类错误：**一个源失败会打断后面所有源**，而失败一次的代价是「接下来几个源这一小时都不更新」。更好的写法和浏览量同步一致——累积错误、继续处理、最后统一返回。我把它列为已知待改进项。

**日历数据**（`server/task/calendar.go`）把日期拼进 key、TTL 一天：`"calendar-" + time.Now().Format("2006/0102")`。日期进 key 意味着不同日期的数据天然隔离，不需要在代码里判断「今天的数据是不是旧的」。

---

## 五、命令行工具

### 为什么要有 CLI

有几个操作天然不适合做成 HTTP 接口：建表/迁移结构（跑一次就该结束）、ES 索引初始化（会**删除已有索引**）、数据导入导出（大文件、长耗时、产出是本地文件）、创建管理员（需要交互式读密码并关闭终端回显）。放成接口等于给服务开了一堆高权限、能被远程触发的后门；做成命令行之后，它们只能在能登录服务器的人手里执行——**把权限收敛到「能上机器」这个更硬的门槛上**。

### 命令分发

`server/flag/enter.go` 用 `urfave/cli` 定义一组 flag，再在 `switch` 里分发。入口刻意做得很克制：`InitFlag()` 只在 `len(os.Args) > 1` 时才构造 `NewApp()` 并执行，执行完直接 `os.Exit(0)`；没有参数就完全跳过。这才让同一个二进制既能当常驻服务，又能当运维工具。

分发时还有一条限制：入口先检查 `c.NumFlags() > 1`，一旦发现同时给了多个 flag 就报 `"Only one command can be specified"` 并退出。这些操作里有几个是破坏性的（重建索引、覆盖导入），允许组合执行意味着「一次手抖同时干两件危险的事」。多写几行判断，换掉一整类误操作。

### 建表与 SQL 导出

建表在 `server/flag/sql.go`：

```go
func SQL() error {
	return global.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&database.Advertisement{},
		&database.Comment{},
		&database.User{},
		// ...共 11 张表
	)
}
```

表结构由模型推导，不维护单独的建表 SQL——**模型和表结构只有一处真相**，改字段不会忘记改 SQL。代价是复杂变更（改列类型、加复合索引）的掌控力不如手写 SQL，那种情况我会单独处理。

导出走 `exec.Command("docker", "exec", "mysql", "mysqldump", "-u"+..., "-p"+..., dbName)`——直接调容器里的 `mysqldump`。好处是拿到标准 SQL dump，能灌进任何 MySQL；坏处是**硬依赖「MySQL 跑在一个叫 mysql 的容器里」这个事实**，换部署方式就失效。对单机个人项目，这个假设成立。

### ES 索引初始化：不可逆操作要二次确认

`server/flag/es.go` 在索引已存在时**停下来问一句**，因为删掉重建意味着所有文章正文都要重新导入。输入非法时递归重新询问，而不是默认继续——对破坏性操作，默认值必须是「不做」。

### ES 数据迁移：导出用 Scroll，导入用 Bulk

导出不能一次 `Search` 拉完（默认有返回条数上限），所以用 Scroll 游标分页，每批 1000 条、游标保留 1 分钟，循环到没有更多数据，最后清理游标释放资源。落盘的中间结构定义在 `server/model/other/es_index.go`：`Data` 只有 `ID *string` 和 `Doc json.RawMessage` 两个字段，文档本体用 `json.RawMessage` 原样透传、不做结构体映射——这样导出导入不会因为模型改动而丢字段。

导入时先删旧索引、按 mapping 重建，再用 Bulk 批量写入，并在请求上指定 `Refresh(refresh.True)`。正常写入不该强制 refresh（每次都有代价），但导入是**一次性、离线**的场景，等它慢慢刷不如直接刷新让文档立即可搜。

### 迁移的边界：附件不在 JSON 里

有个容易被忽略的点：导出的只是**文档数据**，文章引用的图片不在里面。图片走另一套东西，`server/utils/upload/upload.go`：

```go
type OSS interface {
	UploadImage(file *multipart.FileHeader) (string, string, error)
	DeleteImage(key string) error
}

func NewOss() OSS {
	switch global.Config.System.OssType {
	case "local":
		return &Local{}
	case "qiniu":
		return &Qiniu{}
	default:
		return &Local{}
	}
}
```

接口抽象掉了本地磁盘和七牛云两种存储，切换只改配置。但这也意味着**数据迁移要分两步**：文档用 CLI 导入，附件要另外拷贝（本地目录，或换 bucket）。上传的文件名是 `utils.MD5V([]byte(name)) + "-" + time.Now().Format("20060102150405") + ext`——带时间戳，同一张图重复上传不会互相覆盖，但也不会去重。上传前会校验扩展名白名单和大小上限，白名单是显式列出的 `map[string]struct{}`，比黑名单可靠：漏掉的类型是被拒绝，而不是被放行。

### 创建管理员：密码不能出现在终端里

`Admin()` 先取标准输入的文件描述符，用 `term.MakeRaw(fd)` 关掉终端回显，读完再 `defer term.Restore(fd, oldState)` 恢复——两次输入不一致就拒绝。管理员密码是这套系统里权限最高的凭据，**不能出现在屏幕、shell history 或日志里**，所以这个初始化动作必须走交互式 CLI，而不是「接口传个 JSON 就建管理员」。用户名和地址直接取配置文件里的站点信息，省一次输入；密码限制 8~20 位，入库前用 bcrypt 哈希。

---

## 六、定时任务的可观测性

上面这套机制目前的可观测性只有两样：任务失败时的 error 日志，和热门榜刷新成功时的一行 Info。这对「任务崩了」够用，但对**第二节那个 bug 完全无效**——它每小时都「成功」返回，日志干干净净，只有对比 Redis 和 ES 的数字才能发现问题。这类「静默地少做了事」的故障靠日志抓不住，需要指标。

如果只能埋两个，我会选：

| 指标 | 类型 | 为什么是它 |
|---|---|---|
| `article_views_pending` | Gauge | Redis Hash 里待同步的条目数与增量总和 |
| `cron_task_last_success_timestamp{task}` | Gauge | 每个任务上一次成功执行的时间戳 |

第一个是**积压量**，能一眼看出「只同步一条」这类问题：任务报成功，但积压量不降反升，趋势图会立刻变得可疑。如果当初有这个指标，那个 bug 上线当天就会暴露。

第二个是**心跳**，解决「任务根本没跑」这一类问题——cron 注册失败、进程重启后调度丢失、时间表达式写错。判断方法很直接：当前时间减去这个时间戳超过两个刷新周期就该告警。它比「任务失败计数」更有用，因为任务压根没执行时失败计数也是零。

再往下排值得考虑的是榜单缓存的年龄（key 的剩余 TTL，用于确认刷新链路还活着）和 ES 调用耗时直方图，但这两个属于「有了更好」，不像上面两个能直接定位一类具体故障。

现实问题是这类指标要有 Prometheus 之类的采集端才有意义，对单机个人博客，专门部署一整套监控的成本高于收益。所以更现实的做法是**先用日志把同样的信息打出来**——每小时固定输出一行「待同步条目数 / 增量总和 / 本次成功同步条数 / 耗时」，再用最简单的日志告警盯住「待同步条目数连续三次上涨」这一条规则。等它真的长大到需要，再换成正式指标。

---

## 小结

回顾一遍，几个决策其实是同一套思路在不同地方的重复：

**把高频写换成低频批量写。** 浏览量先落 Redis 聚合，再由定时任务合并落库；代价是最多滞后一小时，换来写放大被摊平。

**需要原子性的地方，不要自己拼「读—改—写」。** Redis 侧用 `HINCRBY`，ES 侧用 Painless 脚本，都是把并发正确性交给存储引擎，而不是交给自己的代码。

**批量任务里单点失败不该拖垮整体。** 循环内 `return` 是我付出过真实代价的地方；正确做法是累积错误、继续处理，并且只把**真正成功**的部分从队列里摘掉，失败的下轮重试。

**缓存的新鲜度靠主动刷新保证，TTL 只用来兜底。** 所以 TTL 必须大于刷新周期；查询失败时保留旧缓存而不是清空；空结果也要缓存，避免穿透。

**危险操作收敛到命令行。** 重建索引、导入数据、创建管理员这些事只该在能登录服务器的人手里发生，而且要一次只做一件、不可逆的还要二次确认。

最后一点感受：这套东西里最难发现的，从来不是「报错的任务」，而是「每小时都成功、但只做了一部分事」的任务。错误日志能覆盖前者，覆盖不了后者——这也是为什么我觉得，给定任务埋一个「积压量」和「上次成功时间」，比多写几个 try-catch 更值得。
