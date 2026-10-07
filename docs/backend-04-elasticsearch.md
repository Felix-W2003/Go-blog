# 用 Elasticsearch 做文章存储与全文检索：一个没有 article 表的博客后端

常见做法是 MySQL 存正文、另起一套搜索引擎做检索，两边靠双写或 binlog 保持一致。这个博客没走这条路——**文章正文只存在于 Elasticsearch 里，MySQL 中根本没有 `article` 表**。

这个选择消掉了一批问题，也换来另一批必须自己扛的责任。下面从索引映射、检索实现、写入链路、并发计数一路讲到命令行的数据搬迁工具，重点说明每处为什么这么设计、代价落在哪里。

## 一、一个反直觉的决定：文章不进 MySQL

`server/model/database/` 里有 `article_category.go`、`article_tag.go`、`article_like.go`，却没有 `article.go`；`server/flag/sql.go` 的表清单同样印证——共 11 张表，没有文章表。文章实体定义在 ES 模型里：

```go
// server/model/elasticsearch/article.go（节选，字段注释为对应的 Mapping 类型）
type Article struct {
    CreatedAt string   `json:"created_at"`  // DateProperty
    Title     string   `json:"title"`       // TextProperty：参与相关性检索
    Keyword   string   `json:"keyword"`     // KeywordProperty：整值精确匹配
    Category  string   `json:"category"`    // KeywordProperty：精确筛选
    Tags      []string `json:"tags"`        // []KeywordProperty
    Abstract  string   `json:"abstract"`    // TextProperty
    Content   string   `json:"content"`     // TextProperty
    Views     int      `json:"views"`       // IntegerNumberProperty
    Comments  int      `json:"comments"`
    Likes     int      `json:"likes"`
}
```

**动机是把自己从"双份存储"里解放出来。** 按常规做，MySQL 存一份、ES 存一份，立刻要面对：同步选双写还是 CDC、怎么做一致性对账、索引重建时数据从哪来、加字段要同时改 DDL 和 mapping。而检索需求天生需要分词与相关性排序，翻成 SQL 就是走不了索引的 `LIKE '%关键词%'`——既然检索侧必须是搜索引擎，MySQL 那一份就退化成"只为备份而存在、却要付出全部同步成本"的数据。于是只保留一份。

代价同样明确：

1. **没有跨实体的原子性**——MySQL 事务管不到 ES 写入，第八节详细讨论。
2. **没有 JOIN，计数只能反规范化**——`Views` / `Comments` / `Likes` 直接存进 ES 文档。MySQL 的 `article_like` 记"谁收藏了哪篇"，ES 的 `likes` 记"共被收藏多少次"，**同一份事实存两处，一致性由应用代码维护**。
3. **强依赖 ES 可用性**——`server/initialize/es.go` 里 `NewTypedClient` 失败会直接 `os.Exit(1)`，而所有文章读写都走 ES，没有降级路径。
4. **备份迁移得自己造工具**——没有 `mysqldump` 等价物，导出只能走 Scroll API 写自定义工具；mapping 也不能原地改，只能重建索引再灌数据。
5. **文档 ID 由 ES 生成**（形如 `bR5uH6ABNUHRlLEfuB-A` 的不透明字符串）。`comment` 和 `article_like` 里的 `article_id` 存的就是它，而 **MySQL 没有外键能校验**——删文章必须由应用代码手动级联，这正是 `ArticleDelete` 里那段循环在做的事。

## 二、索引 Mapping：把检索需求翻译成字段类型

选字段类型本质是在回答：**这个字段将来会被"搜"，还是被"精确匹配/排序/聚合"？**

- `text`：会被分词器切成词条。搜 "JWT 双 Token" 能命中含 "Token" 的文档，但字段本身不能排序或精确匹配。
- `keyword`：整体作为一个不可分割的词条。适合精确匹配、排序、聚合，但搜部分关键词命中不了。

```go
// server/model/elasticsearch/article.go（ArticleMapping，节选）
"title":    types.TextProperty{},
"keyword":  types.KeywordProperty{},
"category": types.KeywordProperty{},
"tags":     []types.KeywordProperty{},
"abstract": types.TextProperty{},
"content":  types.TextProperty{},
"views":    types.IntegerNumberProperty{},
"comments": types.IntegerNumberProperty{},
"likes":    types.IntegerNumberProperty{},
```

**title 与 keyword：同一个值，两种类型。** 注意 `title` 是 text，而名为 `keyword` 的字段是 keyword 类型——写入时两者取的是同一个值（`Title: req.Title, Keyword: req.Title`），等于同一个标题存两份。

冗余的原因是不同查询需要不同物理表示：`title` 负责相关性检索，`keyword` 负责整值精确匹配，用来判断标题是否已存在：

```go
// server/service/article_helpers.go
func (articleService *ArticleService) Exits(title string) (bool, error) {
    req := &search.Request{Query: &types.Query{
        Match: map[string]types.MatchQuery{"keyword": {Query: title}},
    }}
    res, err := global.ESClient.Search().Index(elasticsearch.ArticleIndex()).Request(req).Size(1).Do(context.TODO())
    return res.Hits.Total.Value > 0, nil
}
```

拿 text 型 `title` 判重会误判：搜 "Token" 会把标题为 "JWT 双 Token 认证" 的文章也算成"已存在"。keyword 字段不分词，`Match` 等价整值比较，得到严格相等语义（语义上 `Term` 更贴切，结果一样）。代价是每篇多存一份标题——这个数据量下用空间换语义清晰很划算。

**category / tags 必须是 keyword。** 它们只参与"筛"和"聚合"：按类别/标签过滤、给标签云算数量。建成 text 的话 "Go" 会被分词，`Term` 查不到，聚合出的桶名也是拆开的。

**日期用 DateProperty 并显式声明格式**（`yyyy-MM-dd HH:mm:ss`），与写入时 `time.Now().Format("2006-01-02 15:04:05")` 严格对应。这里有个易忽略的耦合：**Go 的布局字符串和 ES 的 format 必须一致**，否则 ES 解析失败、字段无法排序。建成字符串的话，这个格式的字典序恰好等于时间序，排序也能工作，但会失去范围查询和日期运算能力。

**计数用 integer**，因为它们要参与排序（按浏览量排热门）和数值自增——Painless 的 `ctx._source.views += n` 要求字段是数值，字符串上做加法会报错。

顺带指出一处瑕疵：`cover` 被建成 text，但它实际是个 URL，既不参与全文检索、将来也可能要精确匹配，标准分析器会把它按 `/`、`.`、`-` 切碎并小写化。目前它只作为字符串原样返回（分词只影响索引，`_source` 存原值），**功能没问题但类型不够准确**。

## 三、typed client：把 JSON 拼装换成结构体

项目用的是 ES 8 官方 typed 客户端（`server/go.mod` 里的 `github.com/go-elasticsearch/v8 v8.19.6`）。老式客户端要自己把查询序列化成 `map[string]interface{}` 再塞进 `strings.NewReader`，字段名写错只能等运行时；typed 版是链式配置加结构体请求响应：

```go
// server/service/article_helpers.go
res, err := global.ESClient.Search().
    Index(elasticsearch.ArticleIndex()).
    Request(req).      // req 是 *search.Request，字段写错编译期就报错
    Size(1).
    Do(context.TODO())
```

**收益是类型安全与可发现性**：IDE 能补全 `types.BoolQuery` 的成员，响应 `res.Hits.Hits` 直接是 `[]types.Hit`，不用逐层断言。

**代价是啰嗦**。typed API 由规范自动生成，排序要写成 `[]types.SortCombinations` 套 `types.SortOptions` 再套 `map[string]types.FieldSort`，`Order` 还要枚举指针。更明显的是**部分更新仍要手写 JSON**——`update.Request.Doc` 的类型是 `json.RawMessage`：

```go
// server/service/article_helpers.go
bytes, err := json.Marshal(v)
_, err = global.ESClient.Update(elasticsearch.ArticleIndex(), articleID).
    Request(&update.Request{Doc: bytes}).
    Refresh(refresh.True).
    Do(context.TODO())
```

typed 覆盖了"查询构建"这条主路径，但在"部分文档更新"这种天然是任意 JSON 的地方，类型系统帮不上忙——这是 union / 自由结构在强类型语言里的固有摩擦。

## 四、检索：Bool Query 的三段式组合

`ArticleSearch` 把用户输入翻译成 Bool Query，三个子句各司其职：

```go
// server/service/article.go（ArticleSearch 片段）
// 关键词：多字段匹配，放 should —— 命中越多分越高
if info.Query != "" {
    boolQuery.Should = []types.Query{
        {Match: map[string]types.MatchQuery{"title": {Query: info.Query}}},
        {Match: map[string]types.MatchQuery{"keyword": {Query: info.Query}}},
        {Match: map[string]types.MatchQuery{"abstract": {Query: info.Query}}},
        {Match: map[string]types.MatchQuery{"content": {Query: info.Query}}},
    }
}

// 标签：放 must —— 必须满足且参与打分
if info.Tag != "" {
    boolQuery.Must = []types.Query{{Match: map[string]types.MatchQuery{"tags": {Query: info.Tag}}}}
}

// 类别：放 filter —— 必须满足但不打分
if info.Category != "" {
    boolQuery.Filter = []types.Query{{Term: map[string]types.TermQuery{"category": {Value: info.Category}}}}
}
```

**为什么过滤条件用 `filter` 而不是 `must`？** 两点：类别是精确筛选，用户选"技术"就是只要技术类，不存在"技术类更相关"的程度问题，参与打分纯属浪费 CPU；二是 filter 子句结果会被 ES 缓存成 bitset，重复的过滤条件（大量用户在同一类别翻页）能直接复用，而 must 的结果与具体查询词绑定，缓存收益小得多。经验法则是：**能回答"是/否"的条件放 filter，回答"有多相关"的放 should/must。**

**一个容易忽略的语义细节**：`should` 的含义会随上下文变化。当 bool 里**只有 should** 时（只填关键词），ES 默认 `minimum_should_match = 1`，四个字段至少命中一个；当**同时存在 must 或 filter** 时（关键词 + 标签），should 默认退化成 0——**它变成纯粹的加分项，不再过滤**。这恰好是想要的效果：选了标签后，关键词成为排序依据而非过滤条件。

多维度排序用一个字段映射收敛了四种维度：

```go
// server/service/article.go（ArticleSearch 片段）
switch info.Sort {
case "time":    sortField = "created_at"
case "view":    sortField = "views"
case "comment": sortField = "comments"
case "like":    sortField = "likes"
default:        sortField = "created_at"
}
```

两个考虑：**默认降序**（`info.Order != "asc"` 把"没传"和"传 desc"归为一类，缺省即最新/最热优先）；**排序字段走白名单**（default 兜底到 `created_at`，用户传 `sort=content` 也拼不出对正文排序的请求）。排序字段要拼进查询体，不能直接采信外部输入。另外可以留意：`ArticleList` 只在 MatchAll 分支设默认排序，一旦按标题搜索就没有 sort 子句、结果按 `_score` 排——对后台检索是合理的。

**SourceIncludes 裁剪返回字段。** 列表接口不需要正文：

```go
// server/service/article.go（ArticleSearch 片段）
SourceIncludes: []string{"created_at", "cover", "title", "abstract", "category", "tags", "views", "comments", "likes"},
```

注意**没有 `content` 和 `keyword`**：前者是列表不渲染的正文，后者是给精确匹配用的冗余字段。热门榜单裁得更彻底，只取 `[]string{"title", "views", "created_at"}`。单看省得不多，配上 Redis 缓存意义就大了——**缓存体积正比于返回字段数**，字段越少，缓存条目越小、序列化越便宜。

## 五、分页工具：把参数收敛成一个对象

`server/model/other/pagination.go` 为两种存储各定义了一个参数对象：

```go
// server/model/other/pagination.go
type MySQLOption struct {
    request.PageInfo
    Order   string
    Where   *gorm.DB
    Preload []string
}

type EsOption struct {
    request.PageInfo
    Index          string
    Request        *search.Request
    SourceIncludes []string
}
```

`EsOption` 把一次 ES 分页需要的东西收在一处：**页码与每页条数**（内嵌 `request.PageInfo`）、**目标索引**、**查询请求体**、**返回字段**。调用方只描述"要查什么"，"怎么翻页"交给工具：

```go
// server/utils/pagination.go（省略错误处理）
func EsPagination(ctx context.Context, option other.EsOption) (list []types.Hit, total int64, err error) {
    if option.Page < 1 { option.Page = 1 }              // 默认值兜底做在工具里
    if option.PageSize < 1 { option.PageSize = 10 }

    from := (option.Page - 1) * option.PageSize
    option.Request.Size = &option.PageSize
    option.Request.From = &from

    res, err := global.ESClient.Search().
        Index(option.Index).
        Request(option.Request).
        SourceIncludes_(option.SourceIncludes...).
        Do(ctx)
    if err != nil {
        return nil, 0, err
    }
    return res.Hits.Hits, res.Hits.Total.Value, nil
}
```

**默认值兜底放在工具里**是个好决定：每个调用方都不必重复这段防御，前端传 `page=0&page_size=0` 也拼不出非法请求。

**泛型的边界在这里体现得很清楚。** 同一个文件里 MySQL 版本用了泛型：

```go
// server/utils/pagination.go
func MySQLPagination[T any](model *T, option other.MySQLOption) (list []T, total int64, err error)
```

GORM 能把任意 `*T` 映射成表并填充 `[]T`，泛型很自然。但 `EsPagination` 返回的是原始命中信封 `[]types.Hit`，不是 `[]elasticsearch.Article`——**这不是偷懒，而是由 ES 的响应形态决定的**：ES 永远返回信封（`_index` / `_id` / `_score` / `_source`），业务文档只是 `_source` 里的一段 JSON。要变成 `[]T` 就必须逐个 `json.Unmarshal(hit.Source_, &item)`，而"反序列化成什么类型"是调用方的知识。

泛型解决的是"返回类型统一、具体类型未知"，而这里的问题是**同一份响应在不同场景下要投影成不同形状**。项目里两种做法都有：列表接口让前端直接消费信封（按 `_id` 取文档 ID、按 `_source.xxx` 取字段）；收藏列表则在调用方手动投影，现定义一个带 `_id` / `_source` 标签的匿名结构体来承接（见 `ArticleLikesList`）。两者都说明：**ES 的强类型只覆盖到请求构建，响应侧的类型信息在 `_source` 之前就断了。**

还有个副作用值得一提：`EsPagination` 收的是值传递的 `EsOption`，但里面的 `Request` 是**指针**，所以填 `Size`/`From` 改的是**调用方那个请求对象**。因为每次都现场构造新请求，实际不会出问题；但若要复用同一请求分多次查询，后一次会覆盖前一次。

## 六、写入链路：索引、批量、部分更新

ES 的写操作集中在 `server/service/article_helpers.go`。**创建时必须立即刷新**：

```go
// server/service/article_helpers.go
func (articleService *ArticleService) Create(a *elasticsearch.Article) error {
    _, err := global.ESClient.Index(elasticsearch.ArticleIndex()).
        Request(a).
        Refresh(refresh.True).
        Do(context.TODO())
    return err
}
```

ES 是准实时的：文档写入先进内存缓冲，默认**最多 1 秒**才被 refresh 成可搜索的段。不加这个参数，"发布成功"后立刻去列表页可能搜不到，用户会以为发布失败。代价是每次 refresh 都产生新 Lucene 段，段多了触发合并、写吞吐下降。

在博客场景这个取舍划算——发文是低频操作而"所见即所得"要求很高。反过来，**浏览量和收藏这类高频且允许延迟的更新就不该带 refresh**：

```go
// server/service/article.go（ArticleLike 片段）
source := "ctx._source.likes += " + strconv.Itoa(num)
script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), req.ArticleID).
    Script(&script).
    Do(context.TODO())   // 没有 Refresh(refresh.True)
```

同一条写入路径上，两种操作对实时性要求不同，于是用了不同策略——**按操作特征分配代价，而不是无脑统一**。

**批量删除用 Bulk 而非循环单删**，把 N 次网络往返压成 1 次：

```go
// server/service/article_helpers.go（片段）
var request bulk.Request
for _, id := range ids {
    request = append(request, types.OperationContainer{Delete: &types.DeleteOperation{Id_: &id}})
}
_, err := global.ESClient.Bulk().Request(&request).
    Index(elasticsearch.ArticleIndex()).Refresh(refresh.True).Do(context.TODO())
```

`bulk.Request` 是 `[]types.OperationContainer`，每个元素包一个动作。这里 `Id_: &id` 取循环变量地址——Go 1.22 之前这是经典陷阱（变量复用导致所有元素指向同一值），1.22 起每次迭代新建变量，而 `go.mod` 声明的是 `go 1.26.2`，所以安全。

**部分更新时刻意用窄结构体。** `ArticleUpdate` 构造的不是完整 `Article`，而是**不含计数字段**的匿名结构：

```go
// server/service/article.go（ArticleUpdate 片段）
articleToUpdate := struct {
    UpdatedAt string   `json:"updated_at"`
    Cover     string   `json:"cover"`
    Title     string   `json:"title"`
    Keyword   string   `json:"keyword"`
    Category  string   `json:"category"`
    Tags      []string `json:"tags"`
    Abstract  string   `json:"abstract"`
    Content   string   `json:"content"`
    // 故意不含 views / comments / likes / created_at
}{}
```

因为 ES 的部分更新是 `doc` 合并语义：只提交要改的字段，其余原样保留。如果图省事传完整实体，计数字段就是零值，**每次编辑文章都会把浏览量、评论数、收藏数清零**——而这三个数字来自 Redis 计数器和 MySQL 关联表，ES 里清掉就再也恢复不回来。用窄结构体更新，是个看起来不起眼、实际上保命的约定。

## 七、Painless 脚本：并发安全的计数

收藏计数没有"读出来 +1 再写回"，而是发一段 Painless 脚本：

```go
// server/service/article.go（ArticleLike 片段）
source := "ctx._source.likes += " + strconv.Itoa(num)   // num 为 +1 或 -1
script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), req.ArticleID).Script(&script).Do(context.TODO())
```

**读改写为什么会丢更新**：假设 `likes = 10`，两个用户几乎同时点赞——

```
请求 A: GET  -> 读到 10
请求 B: GET  -> 读到 10
请求 A: PUT likes = 11
请求 B: PUT likes = 11      <- A 的那次 +1 被覆盖，最终少了 1
```

这就是经典的**丢失更新**，低并发下也会偶发，很难复现——表现是"收藏数偶尔对不上"。

把 `+= 1` 送到 ES 端执行就不一样：更新在分片内**针对文档当前版本**应用，脚本执行与文档写入是一体的，不存在"我读到的值已被别人改过"的窗口。并发请求各自作用在当前最新值上，两个 `+1` 都不会丢。ES 的 update API 内部用乐观并发控制（文档版本号）处理冲突并在版本变化时重试，上层不需要自己写重试循环。同一模式也用在浏览量同步（`server/task/article_views.go`）里，脚本同样是 `ctx._source.views += <增量>`。

**结论：只要字段是"计数"语义且会被并发修改，更新方式就该是"发送自增意图"而不是"提交最终值"。** 这条规则与具体用 Redis 的 `HINCRBY` 还是 ES 的 Painless 无关。

> `views` 并非每次访问都直接写 ES：文章被读取时先在一个 Redis Hash 里累加（`server/service/article_stat.go` 的 `CountDB.Set`），再由定时任务按小时批量同步进 ES，把 N 次随机写合并成 1 次顺序写。这条削峰链路单独成篇，这里只交代它在写入链路中的位置。

## 八、双写一致性：坦率说清楚做不到什么

既然"多处写入 + 无法原子"是既定事实，能做的只有两件事：**选择写入顺序**、**明确失败时的行为**。

`ArticleCreate` 的实际顺序是：开启 MySQL 事务 → 更新类别计数 → 更新标签计数 → 更新图片归类 → 写入 ES（**不在事务内**）→ 提交事务。由此产生几种失败时序：

- **ES 写入失败**：事务函数返回错误、MySQL 回滚。结果一致（都没写），**这是安全的**。把 ES 写入放在事务函数最后一步，正是为了让这个失败发生在 MySQL 还没提交时，能被干净回滚。
- **ES 成功但 MySQL 提交失败**：比如提交瞬间连接断开。MySQL 回滚，但 **ES 里已经躺着一篇文章**，而类别/标签计数没增加。
- **MySQL 提交成功但进程在返回响应前崩溃**：数据完整，用户看到失败提示可能重试——重试会被 `Exits` 拦下（标题重复），不会产生重复文章。这个时序反而最无害。

目前的实际风险点，直说：

- **第二种时序没有补偿机制。** 没有对账任务，没有"ES 成功但 MySQL 失败"的重试队列或补偿日志。一旦发生，文章会出现在列表里，但类别计数少 1、标签计数少 1，**没有任何自动机制会发现它**。
- **没有幂等键。** ES 写入用自动生成的文档 ID（`Index(...)` 不带 `Id_`），同一个逻辑操作重放一次就会产生第二个文档。标题查重只能防人工重试，防不住系统级重放。
- **删除路径窗口更大。** `ArticleDelete` 在事务里循环删除每篇文章的评论，最后统一调 ES 批量删除。ES 删成功而事务提交失败时，MySQL 里的评论和收藏关系还在，ES 里的文章却没了。

**继续加固的方向是明确的**：用业务唯一键作为 ES 文档 ID 引入幂等、把"ES 写入意图"落一张本地消息表由后台任务重试并对账。这会让架构复杂不少——对个人博客来说，当前"能干净回滚的部分占多数、孤儿文档概率极低"是可以接受的；但**知道风险在哪、知道怎么补**，比假装它不存在重要。

## 九、CLI：索引初始化与数据搬迁

因为文章数据没有 MySQL 那份"权威副本"，ES 侧运维必须自己提供。项目用 `urfave/cli` 做了三个开关（入口 `server/flag/enter.go`）：`-es`、`-es-export`、`-es-import`，它们**不启动 HTTP 服务**，跑完即退出。

**`-es` 索引初始化**：先 `IndexExists` 判断，已存在时打印提示并要求从 stdin 输入 `y` 或 `n`（`y` 删除重建、`n` 退出、其它输入递归重试），最后调用 `IndexCreate(elasticsearch.ArticleIndex(), elasticsearch.ArticleMapping())`。

默认不破坏数据——索引已存在时必须人工输 `y`。因为"重建索引"和"保留现有数据"是两个完全不同的意图，让机器猜不如让操作者确认。同时 `ArticleMapping()` 是索引定义的唯一来源，CLI 和导入流程都调它，不会出现机器上 mapping 与代码不一致。

**`-es-export` 用 Scroll 分批导出**：先发一次 `Search`，带上 `Scroll("1m")`（游标保留 1 分钟）与 `Size(1000)`（每批 1000 条），收集这一批后循环调用 `Scroll` 直到返回空，最后 `ClearScroll` 释放游标。

**为什么不能用 `from + size` 翻页导出？** 深分页时协调节点要把前 `from + size` 条命中的排序结果全部取回归并，`from` 越大越贵，ES 还会用 `max_result_window`（默认 10000）直接拒绝。Scroll 是在某个时间点给结果集建立快照游标，之后每批从游标继续取，代价与深度无关。

导出文件的结构也刻意做窄：`server/model/other/es_index.go` 里每条记录只有 `ID` 和 `Doc`，而 `Doc` 是 `json.RawMessage`——**导出时完全不反序列化文档内容**，原样搬运，省 CPU 且天然无损，将来加字段旧数据也能原样带回。最后落盘为按日期命名的文件（`es_%s.json`）。可留意两处：循环里 `res, err := ...` 用 `:=` **在循环作用域重新声明了 `res`**，所以循环结束后的 `ClearScroll` 用的是第一次的 scroll id（ES 一般仍能定位上下文，但更稳妥是保存最后一次返回的 id）；另外 `os.Create` 写的是当前工作目录。

**`-es-import` 重建索引加 Bulk 灌数据：**

```go
// server/flag/es_import.go（片段，省略索引删除与重建）
var request bulk.Request
for _, data := range response.Data {
    request = append(request, types.OperationContainer{Index: &types.IndexOperation{Id_: data.ID}})
    request = append(request, data.Doc)
}
_, err = global.ESClient.Bulk().Request(&request).
    Index(elasticsearch.ArticleIndex()).Refresh(refresh.True).Do(context.TODO())
```

几点说明：**它是破坏性的**（先删索引再建，等于全量覆盖），导入到一半失败会留下半空索引；**`Id_: data.ID` 保留了原始文档 ID**，这点很关键——MySQL 里的评论和收藏存的就是这个 ID，**如果导入时让 ES 重新生成 ID，所有评论和收藏都会指向不存在的文档**；**Bulk 请求全量组装在内存里**，文章上千时占用可观，生产级做法是按批提交或使用 `esutil.BulkIndexer`。

这三件工具覆盖了 ES-only 下最关键的两个场景：**迁移重建**（导出再导入，因为 mapping 来自代码，新环境一定得到与代码一致的结构）和**危险操作前的兜底**（重建索引会清空数据，先导出是唯一保险）。反过来说，它也暴露了代价：**"备份数据库"从一条 `mysqldump` 变成了一套必须自己维护、自己测试的代码**。

## 十、小结

**要不要把主数据放进搜索引擎。** 当数据天然是文档型（无跨表事务需求）、检索是核心功能时，单一存储能省掉整套同步与对账成本。一旦涉及强关系、跨实体事务或复杂报表，就该把权威副本放回关系库。这个项目能用这条路，很大程度是因为博客文章恰好是自包含的文档。

**索引映射要从查询倒推。** `text` 与 `keyword` 的分工应该问"这个字段要支持什么查询"，而不是凭字段名猜。同一个标题同时存 `title`(text) 和 `keyword`(keyword) 就是这个思路的典型：同一个值、两种查询、两种物理表示。类型选错往往不立刻报错，而是等到某天排序失灵、聚合结果诡异时才暴露。

**写入策略要跟着操作特征走。** 发布文章带 `Refresh(refresh.True)`（用户马上要看到），计数更新不带（允许延迟、保护写吞吐）；计数一律用 Painless 脚本表达自增意图，绝不读改写；更新文档用窄结构体，避免把不该动的计数字段清零。

**跨两个存储就要接受"不存在原子边界"。** 然后回答三个问题：谁先写、失败时谁能回滚、孤儿数据谁来发现。这里的答案是"先 MySQL 后 ES、ES 写在事务最后一步以便回滚、孤儿数据目前无人发现"——最后一句是已知的空白，也是下一步最该补的地方。

**当数据库不再自带备份，备份就成了你自己的代码。** Scroll 导出加 Bulk 导入本身不难，难的是**保留文档 ID**、**让 mapping 有单一来源**、**给破坏性操作加确认**——这三件事做全了，一条自建的备份通道才算可靠。

这不是"用 ES 更高级"，而是在明确知道代价的前提下做的一次交换：用"自己写运维工具、自己扛一致性风险"，换掉"双份存储 + 无限期的一致性维护"。判断这类交换是否划算，标准只有一个——**你为它额外写的代码，是不是少于你省下的代码。**
