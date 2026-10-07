# Go 博客后端：架构总览与分层设计

这篇文章拆解一个已经上线的个人博客后端——Go + Gin，配合 GORM、MySQL、Redis 和 Elasticsearch 8，共 152 个 Go 文件、79 个 REST 接口。重点不是"用了什么框架"，而是**为什么这样分层、边界划在哪里、每个存储为什么承担它现在的角色**。

读完能收获三件事：一套可复用的 Go Web 分层判断标准；一个把"全局单例 + 聚合注入"讲透的实践样本；以及用"代价"而不是"功能"做存储选型的思路。

## 一、这套后端是做什么的，为什么是 Go + Gin

它是博客的内容侧后端：对外提供文章、评论、友链、留言、广告、站点配置等接口，对内提供管理后台所需的管理接口，外加一套命令行工具和若干定时任务。

选 Go + Gin 主要是三条务实理由。**部署成本**：Go 编译成单个静态二进制，丢到服务器就能跑，个人项目没有专职运维，这比语言特性更重要。**客户端生态**：GORM、go-redis、go-elasticsearch 都是官方或事实标准，尤其 ES 8 的 typed client 让查询能拿到编译期检查。**Gin 的分组与中间件模型够用且直白**，`Router.Group()` 加 `Use()` 就能表达"公开 / 登录 / 管理员"三档权限，不需要额外抽象层。

相比之下，Spring Boot 对一个人维护的博客抽象太重；Node 在处理 ES 查询和并发写计数时，类型与并发模型都不如 Go 直接。

## 二、目录结构：先有一张地图

```
server/
├── main.go                 程序入口，唯一负责"顺序"的文件
├── api/                    11   HTTP 处理器：绑参 + 调 service + 写响应
├── service/                20   业务逻辑：事务、多数据源编排、缓存策略
├── router/                 11   把 handler 挂到三个权限组上
├── middleware/              4   日志/Recovery、JWT 鉴权、管理员校验、登录记录
├── model/
│   ├── request/ response/       入参与出参契约（10 + 8）
│   ├── database/               11 张 MySQL 表模型
│   ├── elasticsearch/          文章文档模型 + 索引名 + Mapping
│   ├── appTypes/               业务枚举：角色、注册来源、图片类别、存储类型
│   └── other/                  跨层传递的中间结构（分页选项、热搜、日历）
├── initialize/              6   GORM、Redis、ES、路由、cron 的初始化
├── global/                  2   全局句柄声明（Config / Log / DB / Redis / ES）
├── config/                 14   配置结构体，与 config.yaml 一一对应
├── core/                    5   框架级组装：读配置、建 Logger、起 HTTP 服务
├── flag/                    8   命令行工具：建表、导入导出、创建管理员
├── task/                    5   定时任务：浏览量落库、热门榜、热搜、日历
└── utils/                  16   无状态工具 + hotSearch/ + upload/
```

目录是按**职责**切，不是按技术层堆。`model` 拆成六个子包，是因为它们的"变化原因"不同：`request` / `response` 跟着接口契约走，`database` 跟着表结构走，`elasticsearch` 跟着索引 Mapping 走，`appTypes` 是一组被各层共享的纯枚举，`other` 放的是既不落库也不返回给前端、只在层间传递的结构。全塞进一个包，改一个字段就要在几十个文件里做影响分析。

`utils` 下再分 `hotSearch` 和 `upload`，是因为这两块各自有"多种实现"——热搜有百度、知乎、快手、头条四个来源，对象存储有本地和七牛云两条通道，各自用一个构造函数把差异收在内部。`appTypes` 这类枚举包看似可省，但它让 `appTypes.Admin` 比裸 `2` 可读得多，也把"角色只有三个取值"固定在一处。

## 三、启动顺序：依赖是怎么被拉起来的

整个程序的骨架只有 25 行：

```go
// server/main.go
func main() {
	global.Config = core.InitConf()
	global.Log = core.InitLogger()
	initialize.OtherInit()
	global.DB = initialize.InitGorm()
	global.Redis = initialize.ConnectRedis()
	global.ESClient = initialize.ConnectEs()

	defer global.Redis.Close()

	flag.InitFlag()
	initialize.InitCron()
	core.RunServer()
}
```

这个顺序每一步都有依赖理由。**配置必须最先**，后面每个初始化都要从 `global.Config` 取自己那一段；它失败时用标准库 `log.Fatalf`，因为此时 logger 还不存在。**日志第二**，之后所有初始化失败都能记结构化日志再退出，而不是黑盒退出。**`OtherInit` 第三**，它做启动期校验：解析 JWT 的两个有效期字符串，任一个格式非法就 `os.Exit(1)`——把"配置写错"放在连数据库之前暴露，能在 0.1 秒内失败，而不是等 MySQL 连完再说。

**三个存储连接排在 CLI 之前**，这点容易忽略：`flag.InitFlag()` 里的建表、导入导出命令同样需要 DB 和 ES 已就绪。三个连接函数策略一致——**连不上就退出，不做降级**。缺了主存储的博客后端没有任何有意义的工作，带着残缺状态启动只会把错误推迟到第一个请求。

`flag.InitFlag()` 是个"截胡"设计：`os.Args` 里有参数就执行对应命令然后 `os.Exit(0)`，压根走不到后面的 HTTP 服务。同一个二进制既是服务端也是运维工具，不必额外维护一套脚本。

## 四、一次请求的完整流转

以 `GET /api/article/hot` 为例，请求要穿过五层。

### 4.1 路由初始化：三个权限组

```go
// server/initialize/router.go
publicGroup := Router.Group(global.Config.System.RouterPrefix)
privateGroup := Router.Group(global.Config.System.RouterPrefix)
privateGroup.Use(middleware.JWTAuth())
adminGroup := Router.Group(global.Config.System.RouterPrefix)
adminGroup.Use(middleware.JWTAuth()).Use(middleware.AdminAuth())
```

三个组共享同一前缀（配置里的 `router_prefix`），差别只在挂了几层中间件。随后每个 router 同时拿到这三个组，**由 router 自己决定某个接口该放哪一档**：

```go
// server/router/article.go（节选）
func (a *ArticleRouter) InitArticleRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup, AdminRouter *gin.RouterGroup) {
	articleRouter := Router.Group("article")
	articlePublicRouter := PublicRouter.Group("article")
	articleAdminRouter := AdminRouter.Group("article")
	articleApi := api.ApiGroupApp.ArticleApi

	articlePublicRouter.GET("hot", articleApi.ArticleHot)
	articlePublicRouter.GET(":id", articleApi.ArticleInfoByID)
	articleAdminRouter.POST("create", articleApi.ArticleCreate)
}
```

"把三个组一起传进去"让**权限归属和接口定义挨在一起**。读 `router/article.go` 一个文件就知道哪些接口公开、哪些要登录、哪些必须是管理员，不用跳去别处对照。

### 4.2 中间件链

全局中间件只有两个，在 `InitRouter` 里挂上：`GinLogger` 把 path、status、耗时、IP、UA 记进 zap；`GinRecovery` 捕获 panic，并**特判 broken pipe**——客户端断开连接不是程序缺陷，不该刷一屏堆栈。

`public` 组的 `hot` 接口到这里就结束；`admin` 组还要穿过 `JWTAuth` 和 `AdminAuth`。后者被压到极小，只做一件事——从上下文取角色并比对（`server/middleware/admin.go`）。取角色的细节（先看 `gin.Context` 里有没有解析好的 claims，没有才回头解析 token）被 `utils.GetRoleID` 收走了。中间件保持"薄"，才能被自由组合。

### 4.3 handler：只做三件事

```go
// server/api/article.go
func (articleApi *ArticleApi) ArticleHot(c *gin.Context) {
	list, total, err := articleService.ArticleHot()
	if err != nil {
		global.Log.Error("Failed to get hot articles:", zap.Error(err))
		response.FailWithMessage("Failed to get hot articles", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}
```

`api` 包 11 个文件、79 个接口，处理器全部是这个形状：**绑参 → 调 service → 写响应**，中间加一次错误日志。

这一点可以量化验证：在整个 `api/` 目录里搜 `global.DB` 和 `global.ESClient`，**结果是零次**；`global.Redis` 只出现一次。除此之外 api 层用到的全局对象只有 `global.Log`（记日志）和 `global.Config`（读验证码尺寸、OSS 类型这类展示性配置）。

那唯一一处 `global.Redis` 是**边界上的泄漏**，值得单独指出：它为了在踢掉旧会话时算出黑名单的存活时长，直接读了 Redis 的 TTL（`api/user.go`）。更干净的做法是把"取旧会话剩余有效期"包成 service 方法。我留着它是因为只有一行、语义也清晰——但这正说明分层边界需要主动维护，它不会自己保持干净。

### 4.4 service：业务逻辑与数据编排

同一个接口在 service 层要复杂得多，因为"热门文章"背后有三件事：查 ES、维护 Redis 缓存、以及决定每个失败分支的后果。

```go
// server/service/article.go
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

这段是 service 层价值的缩影：**它组合多个数据源，并决定每个失败分支的后果**。缓存写失败只记日志、不影响响应；ES 查询失败则把错误抛给上层。这类判断放 api 层会让 handler 迅速膨胀，放 model 层又缺少调用上下文。

再看一个跨存储的例子。评论写进 MySQL，但文章评论数存在 ES 里，两者必须一致，于是同步被放进 GORM 钩子：

```go
// server/model/database/comment.go（节选）
func (c *Comment) AfterCreate(_ *gorm.DB) error {
	source := "ctx._source.comments += 1"
	script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
	_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), c.ArticleID).Script(&script).Do(context.TODO())
	return err
}
```

好处是任何创建评论的路径都自动带上这一步，不会有人漏写；代价是**写入路径变得隐式**——只看 service 层，你意识不到"插一条评论"还会改 ES。这是"一致性 vs 可预测性"的取舍，我选了前者，因为这个计数必须准。

## 五、各层职责与"什么该放哪一层"

| 层 | 只应该做的事 | 判断标准 |
| --- | --- | --- |
| `api` | 绑参校验、调 service、封装响应 | 这行代码是否只与 HTTP 有关？ |
| `service` | 业务规则、事务、多数据源编排、缓存策略 | 换一个入口（CLI / cron）还需不需要它？ |
| `model/database` | 表结构、关联、少量一致性钩子 | 它是否只描述"数据长什么样"？ |
| `model/elasticsearch` | 文档结构、索引名、Mapping | 同上，只是目标存储不同 |
| `utils` | 无状态、无业务语义的纯函数 | 它需要知道"现在是谁在调用"吗？ |

三条实践中总结出的红线。**api 层不出现 `global.DB` / `global.ESClient`**：一旦出现，业务逻辑就漏到了 HTTP 层，这段逻辑再也无法被 CLI 或定时任务复用。**service 层尽量不出现 `*gin.Context`**：这套代码里有个例外，`UserService.Logout(c *gin.Context)` 直接收了上下文，因为它要读写 Cookie——为了少一层参数搬运而妥协，代价是该方法的可测试性变差。可以接受妥协，但要知道自己在妥协什么。**`utils` 只放纯函数**：某个工具函数一旦开始读 `global.DB`，它就该改名搬去 service。

判断标准里最有用的是第二条："**换一个入口还需不需要它？**"这套后端有三个入口——HTTP、命令行（`flag`）、定时任务（`task`），它们共用同一个 service 层。定时任务和 HTTP 接口调用的是同一个方法（`service.ServiceGroupApp.ArticleService.NewArticleView()`）。这就是分层的实际回报：业务逻辑写一次，三个入口都能用。

## 六、ServiceGroupApp：聚合式依赖注入

service 包有 20 个文件，却没有一个构造函数，它们靠一个聚合结构体组织起来：

```go
// server/service/enter.go
type ServiceGroup struct {
	EsService
	JwtService
	UserService
	ArticleService
	// ...其余 12 个
}

var ServiceGroupApp = new(ServiceGroup)
```

`api` 层同样处理一遍，并多做一步——把常用 service 提升成包级变量：

```go
// server/api/enter.go（节选）
type ApiGroup struct {
	BaseApi
	ArticleApi
	// ...
}

var ApiGroupApp = new(ApiGroup)
var articleService = service.ServiceGroupApp.ArticleService
```

于是 handler 里写 `articleService.ArticleHot()`，而不是一长串链式访问；中间件也用了同一招（`var jwtService = service.ServiceGroupApp.JwtService`）。

这样组织带来三个收益。**依赖只有一个构造入口**：要知道项目用了哪些 service，读 `service/enter.go` 一个文件就够，不必在代码里追 `new`。**依赖方向被强制成单向**：`api → service → model`，`middleware` 和 `task` 可以引用 service，service 不反过来引用它们——从结构上消除了循环依赖的可能，比靠约定和 review 去防更可靠。**替换实现只需改一处**：换掉某个 service，改 `ServiceGroup` 的字段类型即可。

但必须诚实地说，**它不是真正的依赖注入**。服务是无状态空结构体、方法挂在指针接收者上，所以 `new(ServiceGroup)` 的零值就能用；依赖通过包级变量取得，而不是构造函数传入。更准确的叫法是"**聚合 + 服务定位器**"。代价也很明确：**依赖不体现在函数签名里**——看 `ArticleHot()` 的签名，你看不出它用了 Redis 和 ES；想给它写单元测试，得先构造 `global.ESClient` 和 `global.Redis`。如果哪天需要"同一个 service 有 MySQL 和 ES 两种实现"，或要给 service 层补大量测试，就得改成真正的字段注入。

## 七、数据存储的分工

三个存储不是"都用上更完整"，而是各解决一类问题。

### 7.1 为什么文章正文只存在 Elasticsearch

这是最容易招质疑的决定：**MySQL 里根本没有 article 表**。证据有两处。

```go
// server/model/elasticsearch/article.go（节选）
type Article struct {
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	Cover    string   `json:"cover"`
	Title    string   `json:"title"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	Abstract string   `json:"abstract"`
	Content  string   `json:"content"`

	Views int `json:"views"`
}

func ArticleIndex() string {
	return "article_index"
}
```

```go
// server/flag/sql.go（节选）
return global.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
	&database.Advertisement{},
	&database.ArticleCategory{},
	&database.Comment{},
	&database.Image{},
	&database.Login{},
	&database.User{},
	// ...共 11 张表
)
```

**为什么这么做**：博客的核心场景是"按关键词搜正文"。若正文存 MySQL、再同步一份到 ES，就必须处理双写一致性——文章创建成功但索引失败怎么办？正文更新了但索引没刷新怎么办？这类问题在个人项目里没有精力保证。更彻底的做法是**取消双写**：文章只作为 ES 文档存在，MySQL 只保存真正需要事务和关联的数据。

代价很明确：**没有事务**，文章与类别/标签计数的更新只能靠应用层编排（`service/article_helpers.go` 的 `UpdateCategoryCount` / `UpdateTagsCount`）；**没有外键约束**，评论表的 `article_id` 只是字符串，删文章要靠应用层处理关联数据；**ES 挂了整个内容侧不可用**。

这个决定本质上是"用一种一致性风险换掉另一种"——从"两个存储之间可能不一致"换成"一个存储内部没有跨表事务"。对只有文章一种核心实体的博客，后者的问题小得多。

### 7.2 MySQL 存什么

11 张表，共同特征是**需要事务、需要关联、或需要可靠落地**：

| 表 | 为什么需要 MySQL |
| --- | --- |
| `user` / `login` | 账号体系需要唯一约束；登录审计需要按用户查询 |
| `comment` | 自关联（`PID` 指向父评论），需要事务 |
| `feedback` / `advertisement` | 留言与广告位配置，需要可靠持久化 |
| `image` | 图片元数据，`url` 有唯一约束 |
| `friend_link` / `footer_link` | 友链与页脚链接，需要可靠持久化 |
| `article_category` / `article_tag` | 文章分类与标签的**计数**，需要事务保证加减一致 |
| `article_like` | 收藏关系表，需要唯一约束防重复收藏 |

注意最后两张表：它们不存文章本身，只存文章的**衍生数据**。分工很清晰——ES 存内容，MySQL 存关系和计数。

所有表模型都嵌同一个基类：`ID` + `CreatedAt` + `UpdatedAt` + `DeletedAt`（`server/global/model.go`）。用 `gorm.DeletedAt` 而不是物理删除，是为了让"删错了"还有回旋余地；代价是所有查询都自动带 `deleted_at IS NULL`，索引也要为它单独留一位。

### 7.3 Redis 承担什么

Redis 在这里不是"缓存层"，而是四种职责的集合。

**会话状态**：`uuid → jti` 的映射，用来判断某个用户当前活跃的会话是哪一个，多端登录时新登录踢掉旧会话。**JWT 黑名单**：以 jti 的 SHA256 摘要为 key，TTL 与令牌剩余有效期对齐——存摘要而非明文令牌，既省空间，也避免明文凭据落到缓存里。

**浏览量计数**：浏览量先写 Redis Hash 而不是直接写 ES：

```go
// server/service/article_stat.go
func (c CountDB) Set(id string) error {
	err := global.Redis.HIncrBy(c.Index, id, 1).Err()
	return err
}
```

用 `HIncrBy` 而不是"读出来加一再写回"，是因为后者在并发下会丢更新。随后定时任务每小时把增量批量刷进 ES，用 Painless 脚本让 ES 侧的自增也是原子的（`server/task/article_views.go`）。把 N 次写合并成 1 次，**代价是榜单和浏览量最多滞后一小时**——对博客完全可以接受。

**外部数据缓存**：热搜和日历都走同一个读穿模式——先 `global.Redis.Get`，未命中就回源，序列化后 `Set` 写回并设过期时间。热搜缓存一小时（`server/service/hot_search.go`），日历缓存一天（`server/service/calendar.go`）。好处是缓存可以随时清空、服务可以随时重启，第一个请求会自己把它重建起来。

### 7.4 分页：把重复的取数逻辑收进 utils

MySQL 和 ES 的分页写法完全不同，但暴露给 service 的接口被统一成同一形状：

```go
// server/utils/pagination.go（仅列函数签名）
func MySQLPagination[T any](model *T, option other.MySQLOption) (list []T, total int64, err error)
func EsPagination(ctx context.Context, option other.EsOption) (list []types.Hit, total int64, err error)
```

两者都在内部兜默认值（页码 < 1 就当 1，每页 < 1 就当 10），于是 service 里写分页只需要三行（`server/service/friend_link.go`）。`MySQLPagination` 用了泛型，`EsPagination` 没有——因为 ES 的返回永远是 `[]types.Hit`，类型固定，加泛型只增加噪音。**泛型用在哪，取决于类型参数是否真的会变。**

## 八、几个刻意的取舍

### 8.1 为什么不用 wire / fx

依赖图是一棵深度为 3 的树：`main → initialize → (service → api → router)`，没有循环依赖、没有多种实现、没有需要按生命周期启停的组件。在这种形状下，代码生成（wire）或反射容器（fx）的收益接近于零，代价却是：多一个构建步骤、出错时报错信息变得间接、新人要先学框架约定才能读懂启动流程。

现在的 25 行 `main.go` 就是完整的启动说明。**当"手写"的复杂度还没超过"框架"的认知成本时，手写是更优解。** 触发重构的信号很明确：当依赖出现多种实现、或启动顺序需要按环境分支时，再引入容器不迟。

### 8.2 为什么不做微服务拆分

按业务域拆，这套系统能分成用户、内容、评论、配置四个服务。但拆完会发生什么：一次"取文章详情 + 作者信息"的请求，从进程内两次函数调用变成两次网络调用，外加序列化、超时、重试、熔断；评论创建时要更新文章评论数，原本靠事务和钩子能处理的动作，变成跨服务的分布式一致性问题；原本不存在的服务发现、配置中心、链路追踪、日志聚合，全都变成必需品。

而微服务真正的收益——独立扩缩容、技术栈异构、团队解耦——在这个场景下**一个都不成立**：单人维护、单机部署、日访问量很低，也没有第二个团队需要解耦。

所以现在的形态是**单体 + 清晰分层**。分层保留了微服务最核心的那部分好处（边界明确、依赖单向、改一个功能知道改哪几个文件），却不用付分布式的代价。如果哪天真要拆分，现有的 service 层就是天然的拆分线——业务逻辑已经全部集中在那里，没有散落在 handler 里。

## 九、小结

这套后端的架构可以压缩成四句话：

- **入口有三个**（HTTP / CLI / 定时任务），但业务逻辑只有一份，全部落在 `service` 层。
- **边界靠职责划，不靠技术划**：`api` 只碰 HTTP，`service` 只碰业务与数据编排，`utils` 只放纯函数。一条可量化的红线是——`api/` 目录里不出现 `global.DB` 和 `global.ESClient`。
- **存储按问题分工**：ES 存内容并承担检索，MySQL 存关系与计数，Redis 存会话、黑名单、高频计数与热点缓存。每个存储进来都是为了解决一个具体问题，而不是为了让技术栈看起来完整。
- **依赖用最朴素的方式组织**：`ServiceGroupApp` 聚合 + 包级单例，没有 DI 框架。这是"聚合式服务定位器"而非真正的注入——换来一目了然的依赖清单，代价是可测试性与签名的自解释性。

如果只记一件事：**架构决定的价值不在于选了什么，而在于能不能说清每个选择的代价，以及什么条件下会推翻它。** 上面每一处取舍都留下了触发重构的信号——分层出现隐式写入、`api` 层开始摸数据库、service 需要多实现、单体遇到扩缩容瓶颈——这些信号出现的时刻，就是这套设计该被修改的时刻。
