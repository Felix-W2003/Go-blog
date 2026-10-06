2026年9月29日15:16:03  
1、优化了jwt黑名单的存储结构，原本是存在系统内存，现在改为存在redis，为之后学习微服务分布式系统做好准备；

2、优化了jwt在redis里的存储结构，现在Redis存储用户登录状态使用uuid生成固定字长的jti，即结构为Reds[用户的uuid]=jti

3、优化了存储黑名单的逻辑，Access / Refresh 共享 jti，黑名单只存一个 key，中间件先验 Access、查黑名单，过期才走 Refresh 换发，登录踢旧会话、登出、改密、冻结 / 解冻全部接入黑名单

4、冻结 / 解冻时清理 Redis[uuid]，状态干净

5、TokenNext 去掉了变量遮蔽和重复代码


2026年10月5日19点19分
1、修复定时任务只同步一条 —— task/article_views.go:30 的 return err 写在了 for 循环体内：第一次迭代就返回，① 每小时只落一篇的浏览量，其余增量永久积压；② 循环后的 articleView.Clear() 永不执行。

2、修复登出接口 panic —— utils/claims.go 的 GetJti 在 Refresh Token 解析失败时会解引用 nil claims，
   导致「连续登出 / 手动清 Cookie / 换设备」等场景下接口 500；service/user.go 的 Logout 存在同一处隐患。
   现 GetJti 改为返回 (string, error)；Logout 改用 Access Token 的 jti 作为兜底（Access/Refresh 共享 jti），
   并保证 Cookie 与服务端会话无论黑名单是否写入成功都会被清理。
3、修复黑名单 key 可能永不过期 —— service/jwt.go 的 JoinInBlacklist 在解析
   RefreshTokenExpiryTime 失败时只记日志不返回，exp 留在零值 0；而 Redis.Set 的过期时间
   为 0 并不表示「立刻过期」，而是不带 EX/PX 的永久 key。另外 ParseDuration("0s")
   会返回 (0, nil)，只判断 err 拦不住，因此加了「解析失败即返回错误 + 值 <= 0 拒绝写入」
   两道防线；并把硬编码 3 处的 "jwt:BlackList:" 抽成 blacklistKey()，避免前缀写歪导致
   黑名单静默失效。触发路径：后台 JWT 配置页会热更新 global.Config 与 config.yaml，
   而运行中的进程不会重走 OtherInit 的启动校验 —— 现 UpdateJwt 已补上格式与密钥校验。

4、优化冻结用户的黑名单有效期 —— 原先固定用配置里的 7 天兜底，现改为读取会话 key 的
   Redis TTL，使黑名单生命周期与 Refresh Token 的真实剩余寿命对齐；写黑名单失败不再
   静默忽略，改为记录日志。

5、修复异地登录踢旧会话时黑名单时长取错对象 —— api/user.go 的 TokenNext 原先用
   utils.ParseRefreshExp(refreshClaims.ExpiresAt) 计算旧会话的黑名单 TTL，但 refreshClaims 是
   【新】令牌的 claims，等于拿新令牌的寿命给旧令牌上黑名单。常态下只是过度拉黑（多留几个
   到点自动过期的无用 key），但一旦管理员把 refresh_token_expiry_time 调短（例如 7d → 2h），
   旧 jti 只被拉黑 2 小时，而旧 Refresh Token 仍有 6 天以上有效期 —— 2 小时后旧会话会「复活」，
   异地登录强制下线形同虚设。现改为读取 Redis 中 uuid 这个 key 的 TTL（它就是旧 Refresh Token
   的剩余寿命，且必须在 SetRedisJWT 覆盖它之前读取），读不到时交由 JoinInBlacklist 用配置值兜底。

6、同源加固：service/jwt.go 的 SetRedisJWT 也补上「值 <= 0 拒绝写入」的校验 —— 否则配置写成 "0s"
   时，会话 key（uuid → jti）同样会变成永不过期；UpdateJwt 的校验同步升级为「err 或 值 <= 0」
   双重判断。另外移除 config/conf_mysql.go 与 api/article.go 里残留的调试输出（前者会打印含
   明文密码的 MySQL DSN）。

7、新增首页「热门文章」—— 文章正文只存在 Elasticsearch（MySQL 无 article 表），故新增 public 接口
   GET /api/article/hot：match_all + 按 views 降序取前 10（条数由 hotArticleLimit 常量集中定义），
   并用 SourceIncludes 只回传 title / views / cover / created_at 以减小响应体。
   前端新增 HotArticles.vue（置于主内容区首位，每日新闻下移），列表样式与每日新闻共用一套浅色
   语言；顺带修复 web/src/api/article.ts 中 Article 接口缺失 created_at / updated_at 导致的
   5 处类型错误。注意：views 由定时任务每小时从 Redis 同步，榜单最多滞后一小时。

8、工程化：补充 server/utils 单元测试（ParseDuration / ParseRefreshExp / JtiSHA256 / Bcrypt /
   MD5V，共 21 个断言），并添加 GitHub Actions 工作流（go vet + go build + go test -race）；
   同时修复 utils/hotSearch/zhihu.go 把抓取到的 URL 当格式串传给 fmt.Sprintf 的问题
   （URL 含 % 时会输出 %!E(MISSING) 之类乱码，且被 go vet 判定为非恒定格式串）。