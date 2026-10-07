# JWT 双 Token 认证 + RBAC 三级权限 + Redis 黑名单

登录态管理是每个后端项目都绕不开的一块，也是最容易"看起来能跑、实际上经不起推敲"的一块。我在自己的博客系统里前后重构过三轮：第一版单 Token，第二版改成双 Token 把黑名单放在本地内存，第三版把黑名单整体迁到 Redis，并顺手修掉了一批隐藏很深的问题。这篇讲第三版的完整设计。读完大概能收获：单 Token 为什么一定会在"安全"和"体验"之间二选一；一个会话为什么只该对应一个黑名单 key；Redis 黑名单的三个实现细节（为什么存摘要、TTL 取谁的值、前缀为什么要抽常量）；以及三个真实踩过的坑——其中两个能让接口直接 500，一个能让"强制下线"在几小时后悄悄失效。技术栈 Go + Gin + GORM + MySQL + Redis，前端 Vue 3 + TypeScript，代码都来自仓库里的真实实现。

---

## 一、把权限写进路由树，而不是写进 handler

常见写法是在每个 handler 开头判断角色。问题不在于能不能跑，而在于**判断散落到几十个 handler 之后，你就没法回答"哪些接口需要管理员"**——只能一个个文件翻。我用的是三级路由分组：

```go
// server/initialize/router.go
publicGroup := Router.Group(global.Config.System.RouterPrefix)
privateGroup := Router.Group(global.Config.System.RouterPrefix)
privateGroup.Use(middleware.JWTAuth())
adminGroup := Router.Group(global.Config.System.RouterPrefix)
adminGroup.Use(middleware.JWTAuth()).Use(middleware.AdminAuth())
```

三组共用同一个前缀（`api`，来自配置项 `system.router_prefix`），区别只在挂了几层中间件：

| 分组 | 中间件链 | 语义 |
|---|---|---|
| `publicGroup` | 无 | 游客可访问：登录、注册、找回密码、验证码 |
| `privateGroup` | `JWTAuth` | 必须登录：登出、改资料、收藏、评论 |
| `adminGroup` | `JWTAuth` → `AdminAuth` | 管理员：用户列表、冻结、删文章、系统配置 |

各模块的路由注册函数统一接收三个 group，往对应分组挂：`userRouter.POST("logout", …)` 挂 private，`userLoginRouter.POST("login", …)` 挂 public，`userAdminRouter.PUT("freeze", …)` 挂 admin（见 `server/router/user.go`）。这样做的三个收益：**权限可枚举**——想知道哪些接口要管理员，看哪个 group 注册的就行，不用读实现；**职责可组合**——`AdminAuth` 不需要自己解析 Token、查黑名单、判断用户是否存在，那是 `JWTAuth` 的活，它只做一件事：

```go
// server/middleware/admin.go
roleID := utils.GetRoleID(c)
if roleID != appTypes.Admin {
    response.Forbidden("Access denied. Admin privileges are required", c)
    c.Abort()
    return
}
c.Next()
```

`GetRoleID` 直接从 `gin.Context` 取 `JWTAuth` 放好的 claims，不重复鉴权。**第三，`adminGroup` 天然带上了 `JWTAuth`**，管理员接口同时也是登录接口，不需要额外记得"admin 也要先验 Token"。角色定义在 `server/model/appTypes/user_role.go`，用 `iota`：`Guest`(0) / `User`(1) / `Admin`(2)——游客是 0，前端判断"未登录"就是判断 `role_id === 0`。

---

## 二、为什么一个 Token 不够用

单 Token 的死结在于**有效期只有一个旋钮**，两边拧都有问题。**设短**（15 分钟）安全，但用户正写评论就被踢出去，体验不可接受；**设长**（7 天）体验好，但泄露就是 7 天的敞口，更糟的是 **JWT 无状态，服务端没法主动让它失效**——改了密码、发现被盗，也只能等它自然过期。有人说"加黑名单不就行了"。加了黑名单 JWT 就不无状态了，但那是必要代价。问题在于：如果只有一个长 Token，黑名单会把"每次请求都查一次"的成本拉满，而它本来是为了解决"少数情况下要提前作废"这个罕见需求。**双 Token 的思路是把旋钮拆成两个，各管一件事：**

| | Access Token | Refresh Token |
|---|---|---|
| 有效期 | 2 小时 | 7 天 |
| 职责 | 每个业务请求的鉴权凭证 | 只用来换新的 Access Token |
| 使用频率 | 极高 | 极低（Access 过期时才用） |
| 泄露后果 | 最多 2 小时 | 7 天，但服务端可主动作废 |

配置在 `config.yaml` 的 `jwt` 段，两个 Token 用**不同的密钥**签名（`server/utils/jwt.go` 的 `JWT` 结构体分别持有 `AccessTokenSecret` 和 `RefreshTokenSecret`）。密钥隔离的意义是：即使 Access 的签名密钥以某种方式泄露，攻击者也没法伪造 Refresh 换取长期访问——两类令牌信任级别不同，密钥就该分开。Claims 也分开定义：Access 用 `JwtCustomClaims`，内嵌 `BaseClaims`（`UserID` / `UUID` / `RoleID`）加标准的 `jwt.RegisteredClaims`；Refresh 用 `JwtCustomRefreshClaims`，只有 `UserID` 和标准声明（见 `server/model/request/jwt.go`）。Access 里塞了 `RoleID`，这样 `AdminAuth` 不用查库就能判角色；Refresh 刻意不带角色——它只负责证明"这个会话还有效"，权限由换发出来的新 Access 承载。

---

## 三、两个 Token 分别放哪

两个 Token 的安全性要求完全不同，存放位置也不同。**Access 走请求头 `x-access-token`**（`utils.GetAccessToken` 就是读这个头），前端每次请求显式带上；**Refresh 走 HttpOnly Cookie `x-refresh-token`**。因为 Refresh 生命周期长、价值高，**绝不能放在能被 JavaScript 读取的地方**——前端一旦用 `localStorage` 存它，任何 XSS 都能直接偷走并换出一串新 Access，等于永久接管账号。HttpOnly Cookie 由浏览器自动携带、JS 读不到，把 XSS 的影响面压到"最多偷 2 小时的 Access"。

`setCookie` 里有个容易忽略的分支——**domain 要区分 IP 和域名**：

```go
// server/utils/claims.go
func setCookie(c *gin.Context, name, value string, maxAge int, host string) {
    if net.ParseIP(host) != nil {
        // IP 场景：不能把 domain 设成 IP，否则部分浏览器会直接丢弃这个 Cookie
        c.SetCookie(name, value, maxAge, "/", "", false, true)
    } else {
        // 域名场景：显式设置 domain，保证子域之间共享
        c.SetCookie(name, value, maxAge, "/", host, false, true)
    }
}
```

这个分支是被本地开发逼出来的。开发环境访问 `127.0.0.1:8080`，如果照样把 domain 设成 `127.0.0.1`，浏览器行为并不一致——有的接受，有的当非法域丢弃，表现就是"登录接口返回成功但没有 Refresh Cookie"，然后 Access 一过期就莫名掉线。判断一下 IP 并留空 domain，Cookie 就绑定到当前 host，问题消失。最后两个参数是 `secure` 和 `httpOnly`：`httpOnly` 必须为真，`secure=false` 是本地 HTTP 的妥协，部署到 HTTPS 时应该改成 `true`。CSRF 方面：所有需鉴权的接口都要求 `x-access-token` 请求头，而跨站请求带不上自定义头（会被 CORS 预检拦掉），实际攻击面很小。若以后引入纯 Cookie 鉴权的接口（比如文件下载），就必须补 CSRF Token 或 SameSite。

---

## 四、共享 jti：一个会话 = 一个黑名单 key

这是整套设计里我最满意的一个决定。JWT 标准里的 `jti`（JWT ID）用来唯一标识一个令牌，绝大多数实现会给 Access 和 Refresh 各生成一个——毕竟它们确实是两个不同的令牌。但我是**只生成一次，两个令牌复用**：

```go
// server/api/user.go
jti := j.CreateJti()                          // 只生成一次
accessClaims := j.CreateAccessClaims(baseClaims, jti)
refreshClaims := j.CreateRefreshClaims(baseClaims, jti)
```

`CreateJti` 就是一个 UUID v4 字符串（`server/utils/jwt.go`），两个 `Create*Claims` 都把它写进 `RegisteredClaims.ID`。

**为什么？** 因为 jti 在我的设计里承担的不是"标识一个令牌"，而是**"标识一个会话"**。一个用户在一个设备登录一次产生一个会话，会话里有两个令牌，我要作废的是**会话**而不是某个令牌。如果两个令牌各有各的 jti，踢旧会话就要往黑名单写两条记录，中间件每次要判断"这是 Access 还是 Refresh"再决定查哪个 key，登出同样写两条，黑名单 key 数量直接翻倍。共享之后，**一个会话只需要一个黑名单 key**，中间件里因此能这么写：

```go
// server/middleware/jwt.go
// 因为 Access/Refresh 共享 jti，这里查的是同一个 key
res, bErr := jwtService.IsInBlacklist(refreshClaims.ID)
```

用 Access 的 jti 查、用 Refresh 的 jti 查，查的是同一个 key，没有任何分支。这也是后面能顺畅实现"登出时用 Access 的 jti 兜底"的前提。Redis 里还有一份记录支撑多地点登录拦截：`<user_uuid>` → `<当前活跃会话的 jti>`。写入逻辑在 `SetRedisJWT`，注意它的 TTL 用的是 **Refresh Token 的有效期**：

```go
// server/service/jwt.go
dr, err := utils.ParseDuration(global.Config.Jwt.RefreshTokenExpiryTime)
if err != nil {
    return err
}
if dr <= 0 {
    return errors.New("RefreshTokenExpiryTime 非法，拒绝写入以免会话 key 永不过期")
}
return global.Redis.Set(uuid.String(), jti, dr).Err()
```

这个细节在第五节判断"旧会话还剩多久"时会变成关键。

---

## 五、登录流程：三条分支

登录集中在 `server/api/user.go` 的 `TokenNext`，三种登录方式（邮箱密码、注册后自动登录、QQ 互联）最后都汇到这里。结构是一个 `switch`：

```go
// server/api/user.go
if !global.Config.System.UseMultipoint {
    userApi.respondWithLogin(c, ...)
    return
}
oldJti, err := jwtService.GetRedisJWT(user.UUID)
switch {
case errors.Is(err, redis.Nil):
    // 分支一：首次登录
case err != nil:
    // 分支二：Redis 查询出错
default:
    // 分支三：已存在活跃会话，踢掉旧的
}
```

**分支零：多地点登录开关关闭。** 配置项 `system.use_multipoint` 为 `false` 时直接发令牌、不记录会话、不踢任何人，也就是允许一个账号多处同时登录。关掉它，`uuid → jti` 就不会被写入，"踢旧会话"逻辑也不会触发。

**分支一：首次登录。** `GetRedisJWT` 返回 `redis.Nil`，说明该 uuid 在 Redis 里没有记录，直接调 `SetRedisJWT` 写入新会话。**分支二：Redis 查询出错**——这里和分支一必须分开：`redis.Nil` 是"键不存在"这个**业务语义**，其他 error 是"Redis 挂了"这个**故障语义**，混在一起会导致 Redis 故障时被当成首次登录，把已有会话悄悄覆盖掉。分支二只记日志并返回失败，宁可登录失败也不破坏会话状态。**分支三：踢掉旧会话**，最复杂的一条：

```go
// 旧会话的剩余寿命 = Redis 中 uuid 这个 key 的 TTL
// 必须在下面 SetRedisJWT 覆盖它之前读取
oldTTL, ttlErr := global.Redis.TTL(user.UUID.String()).Result()
if ttlErr != nil {
    global.Log.Warn("读取旧会话剩余有效期失败，改用黑名单默认有效期", zap.Error(ttlErr))
    oldTTL = 0
}
if err := jwtService.JoinInBlacklist(oldJti, oldTTL); err != nil {
    response.FailWithMessage("Failed to invalidate jwt", c)
    return
}
if err := jwtService.SetRedisJWT(refreshClaims.ID, user.UUID); err != nil {
    response.FailWithMessage("Failed to set login status", c)
    return
}
```

这里有严格的顺序，三个动作不能换位置：**读 TTL → 写黑名单 → 覆盖会话记录**。第 3 步会把这个 key 的 TTL 重置成新令牌的 7 天，如果顺序颠倒，第 1 步读到的就是新令牌的寿命——这正是第八节第三个坑要展开的问题。

三条分支最后由 `respondWithLogin` 收尾：现算 `ttl = refreshClaims.ExpiresAt - now`（而不是直接用配置里的 7 天，这样 Cookie 的 max-age 和令牌真实寿命严格对齐），调 `utils.SetRefreshToken` 下发 Cookie，并把 Access Token 放进 JSON body 返回。Access 走 body、Refresh 只走 `Set-Cookie`，**Refresh 不出现在响应体里**——这点很重要，如果它同时出现在 JSON 里，前端就有了把它写进 localStorage 的诱惑，HttpOnly 的保护也就白费了。

---

## 六、请求进来之后：中间件的自动续签

`server/middleware/jwt.go` 的 `JWTAuth` 是整套机制里分支最多的地方：

```
① 解析 Access Token → 成功 → 查黑名单 → 不在 → 放行
   └─ 失败 ↓
② 只有「Access 为空」或「Access 过期」才继续，否则直接 401
   ↓
③ 解析 Refresh → 查黑名单 → 查用户 → 换发新 Access → 响应头下发
```

### 为什么"过期"和"格式错误"要区别对待

这是中间件里最关键的一行判断：

```go
// server/middleware/jwt.go
// 只有「Access 为空」或「Access 过期」才尝试用 Refresh 换新
if accessToken != "" && !errors.Is(err, utils.TokenExpired) {
    utils.ClearRefreshToken(c)
    response.NoAuth("Invalid access token", c)
    c.Abort()
    return
}
```

为了能做这个判断，`server/utils/jwt.go` 把解析错误分成四类哨兵错误：`TokenExpired` / `TokenNotValidYet` / `TokenMalformed` / `TokenInvalid`，分类逻辑在 `parseToken` 里解包 `jwt.ValidationError` 的位标志来决定返回哪一个。为什么必须区分？**过期**是**正常**的生命周期事件，每个用户每天都会遇到，正是为了处理它才设计了 Refresh Token，所以必须放行到续签流程；而**签名错误 / 格式错误**是**异常**信号，意味着令牌被篡改、伪造或密钥用错，这种请求不但不该续签，还必须立刻清掉 Cookie——**因为如果攻击者拿一个伪造的 Access Token 就能触发"用 Refresh 换新 Access"，等于给了他一个免费的探测接口。**`accessToken != ""` 这个条件也值得说：Access 完全为空（前端刚启动、内存里还没令牌）是合法的首次续签场景，应该允许走 Refresh 流程，所以判空和判错误类型是"或"的关系。

### 换发新 Access Token

走到第三步时 Refresh 已通过签名校验和黑名单检查，`refreshClaims.UserID` 可以信任了，但**角色信息不在 Refresh Claims 里**，必须回库查一次。这里用 `Select("uuid", "role_id")` 只取两个字段而不是捞整行——续签路径上每个用户的 Access 过期时都会执行这个查询，能省的列就省：

```go
// server/middleware/jwt.go
var user database.User
if err := global.DB.Select("uuid", "role_id").Take(&user, refreshClaims.UserID).Error; err != nil {
    utils.ClearRefreshToken(c)
    response.NoAuth("The user does not exist", c)
    c.Abort()
    return
}
```

然后换发新令牌，**复用 Refresh 的 jti**：

```go
// 关键：复用 Refresh 的 jti，保证会话 ID 不变
newAccessClaims := j.CreateAccessClaims(request.BaseClaims{
    UserID: refreshClaims.UserID, UUID: user.UUID, RoleID: user.RoleID,
}, refreshClaims.ID)
newAccessToken, err := j.CreateAccessToken(newAccessClaims)
c.Header("new-access-token", newAccessToken)
c.Header("new-access-expires-at", strconv.FormatInt(newAccessClaims.ExpiresAt.Unix(), 10))
c.Set("claims", &newAccessClaims)
c.Next()
```

**续签时 jti 必须保持不变。** 如果这里生成新 jti，会话就分裂成"旧 jti 的 Refresh"和"新 jti 的 Access"两个身份，共享 jti 的设计立刻崩塌。保持 jti 不变还有第二个好处：**只要用户活跃，就能无限续签下去**——每次 Access 过期就用同一个 jti 换一张新的，会话记录和黑名单状态都不用动，用户感知不到 Token 过期，除非 7 天没来。新 Access 通过响应头下发，因为**中间件没机会改响应体**（业务 handler 马上就要往里写 JSON 了），响应头是唯一不需要侵入业务代码的位置。

---

## 七、Redis 黑名单的三个实现细节

黑名单本体在 `server/service/jwt.go`，只有几十行，但每行都有取舍。

### 存 SHA256 摘要，不存明文 jti

```go
// server/service/jwt.go
const blacklistPrefix = "jwt:BlackList:"
func blacklistKey(jti string) string {
    return blacklistPrefix + utils.JtiSHA256(jti)
}
```

jti 本身是 UUID 字符串，直接当 key 完全可行。选摘要主要考虑两点：**一是避免明文凭据落在缓存里**——jti 虽然单独拿出来不能反推出可用令牌，但它是一个会话标识符，而 Redis 很容易被运维工具直接看到（`MONITOR`、慢查询日志、`KEYS` 输出、备份文件），把会话标识符以明文散落在这些地方不是好习惯；**二是查询只需要存在性判断**，黑名单只做 `EXISTS`，不需要还原原文，单向哈希完全够用（`utils.JtiSHA256` 就是 `sha256` 后转十六进制）。

### TTL 必须和令牌剩余寿命对齐

黑名单的本质是"这个会话提前作废"，但作废是有期限的——**令牌自然过期后，黑名单记录就是垃圾，必须自动消失**，否则 Redis 会被历史记录撑爆。所以每次写入都要带精确 TTL：`Redis.Set(blacklistKey(jti), 1, expTime)`。值写 `1` 是因为判断只用到 `Exists`（`IsInBlacklist` 就是 `global.Redis.Exists(blacklistKey(jti)).Result()`），中间件拿到 `res > 0` 就认为在黑名单里。用 `EXISTS` 比 `GET` 语义更直接，也不用管值是什么。**TTL 到底取多少**是这个设计里最容易出错的地方，第八节第三个坑专门展开。

### 前缀抽成常量

```go
const blacklistPrefix = "jwt:BlackList:"
```

这一行看着无关紧要，但它是我踩坑之后补上的。原来这个字符串硬编码在三处：写入、查询、`Exists`。三处必须完全一致，**只要有一处写歪，会发生什么？** 答案是：**什么都不会发生**。写入用一个 key，查询用另一个 key，`EXISTS` 永远返回 0，黑名单永远查不到——**登出、冻结、踢人全部静默失效，日志里一条错误都没有**，用户表现为"点了登出，但旧 Token 还能用"。这种"不发脾气只做错事"的 bug 最难排查。抽成常量加构造函数后，写入和查询走同一个函数，结构上就不可能不一致。

---

## 八、三个真实的坑

三个问题的共同点是：**代码能编译、手测能通过、但在特定时序下会做出完全错误的事**。

### 坑一：Redis 的"永不过期"，以及只判断 err 是不够的

**现象。** Redis 迁移后验收时我用 `redis-cli` 查黑名单 key，发现 `TTL` 返回 `-1`——含义是"key 存在但没有设置过期时间"。**这些本该几小时后自动消失的记录会永远留在 Redis 里。** 原来的 `JoinInBlacklist` 长这样：

```go
// 修复前
exp, err := utils.ParseDuration(global.Config.Jwt.RefreshTokenExpiryTime)
if err != nil {
    global.Log.Error("时间解析失败", zap.Error(err))  // 只记日志，没有 return
}
jtiHash := utils.JtiSHA256(jti)
if expTime > 0 {
    return global.Redis.Set("jwt:BlackList:"+jtiHash, 1, expTime).Err()
} else {
    return global.Redis.Set("jwt:BlackList:"+jtiHash, 1, exp).Err()
}
```

两个问题叠在一起。**第一，解析失败时 `exp` 保持零值 `0`，而 `err` 只被记了日志。** 这里必须知道 go-redis 的一个行为——`Set` 只有在过期时间 `> 0` 时才会在命令里带上 `EX`/`PX`。所以 **`Set(key, value, 0)` 不是"立刻过期"，而是"完全不带过期参数"，等价于 `SET key value`，也就是永久 key**，负数同理。这是语义陷阱：直觉上 `0` 表示"没有剩余寿命"，实际表示"永生"。**第二，光把 `err` 返回还不够。** 我在写回归测试时才发现，`ParseDuration("0s")` 会返回 **`(0, nil)`**——解析成功、没有错误、但值是 0。看实现就明白了：它按 `d`/`h`/`m`/`s` 四个单位逐个扫字符串，`"0s"` 匹配到单位 `s`，前面部分是 `"0"`，`strconv.Atoi("0")` 得到 0 且 `err` 为 nil，累加进去后最终返回 0 也不报错（见 `server/utils/parse.go`）。同理 `"0d"` 也是 `(0, nil)`，所以**配置写成 `0s` 时，`if err != nil` 这道防线根本不触发**，照样写出永久 key。**这个 bug 怎么被触发？** 我原本以为"配置写错服务就起不来"——启动时的 `OtherInit` 确实会校验并 `os.Exit(1)`。但 JWT 配置**可以在后台管理页热更新**，而热更新路径只做了 JSON 绑定、没有校验。管理员把它填成 `0s` 或 `7天`，运行中的进程就带着坏配置继续跑，直到下次重启才起不来。**修法**是两道防线，一道判错误、一道判值：

```go
// server/service/jwt.go
func (jwtService *JwtService) JoinInBlacklist(jti string, expTime time.Duration) error {
    if expTime <= 0 {
        exp, err := utils.ParseDuration(global.Config.Jwt.RefreshTokenExpiryTime)
        if err != nil {
            return fmt.Errorf("解析 RefreshTokenExpiryTime 失败，无法确定黑名单有效期: %w", err)
        }
        expTime = exp
    }
    // 兜底：ParseDuration("0s") 返回 (0, nil)，解析成功但值为 0
    if expTime <= 0 {
        return errors.New("黑名单有效期非法，拒绝写入以避免 key 永不过期")
    }
    return global.Redis.Set(blacklistKey(jti), 1, expTime).Err()
}
```

判值用 `<= 0` 而不是 `== 0`，因为负数行为和 0 完全一样。同一轮里我给 `SetRedisJWT` 也补了同样的校验——会话 key `uuid → jti` 有完全相同的风险。这类"值不合法"的校验更该在**入口**做，所以我在配置热更新的服务方法里也加了校验，把非法值挡在写进 `config.yaml` 之前。

### 坑二：登出接口的 panic

**现象。** 连续点两次登出，第二次接口返回 500，日志里是 `invalid memory address or nil pointer dereference`。**触发条件比想象的宽**，任何"带着有效 Access Token、但没有有效 Refresh Cookie"的情况都会触发：连续登出（第一次已清 Cookie）、用户手动清了 Cookie 但页面内存里还有 Access、换设备只拿到 Access 没有 Refresh。**为什么中间件没拦住？** 因为这些请求携带的是**有效的 Access Token**，`JWTAuth` 走第一条分支——校验 Access、查黑名单、放行，**整个流程根本不碰 Refresh Cookie**，请求就这样进到了 `Logout` 里面。**根因在两处，是同一类错误**，第一处在工具函数 `GetJti`：

```go
// 修复前
func GetJti(c *gin.Context) string {
    refreshToken, _ := c.Cookie("x-refresh-token")
    refreshClaims, err := j.ParseRefreshToken(refreshToken)
    if err != nil {
        global.Log.Error("refresh token 解析失败:", zap.Error(err))  // 只记日志
    }
    return refreshClaims.ID   // ← 失败时 refreshClaims 是 nil，直接 panic
}
```

`ParseRefreshToken` 失败时返回 `(nil, err)`，代码只记了日志就解引用。第二处在 `Logout` 内部，同样写法（取到 `refreshClaims` 后直接访问 `ExpiresAt`）。**这里有个必须一起修的坑**：只修 `Logout` 那处，`GetJti` 会先炸；只修 `GetJti`，panic 只是从一行挪到下一行。同一个病，必须同一轮处理。**修法分两步**，第一步让工具函数返回错误，而不是内部消化：

```go
// server/utils/claims.go
func GetJti(c *gin.Context) (string, error) {
    refreshToken, _ := c.Cookie("x-refresh-token")
    if refreshToken == "" {
        return "", errors.New("x-refresh-token 不存在")
    }
    refreshClaims, err := NewJWT().ParseRefreshToken(refreshToken)
    if err != nil {
        return "", err
    }
    return refreshClaims.ID, nil
}
```

第二步更有意思：**Refresh Cookie 不可用时，用 Access Token 的 jti 兜底**。

```go
// server/service/user.go
// 兜底来源：Access Token。JWTAuth 已经校验过它，claims 一定在 context 里。
// Access / Refresh 共享同一个 jti，所以即使 Refresh Cookie 丢了也能拉黑当前令牌
if claims := utils.GetUserInfo(c); claims != nil && claims.ExpiresAt != nil {
    jti = claims.ID
    blacklistTTL = time.Until(claims.ExpiresAt.Time)
}
// 首选来源：Refresh Token，它的剩余有效期更长，覆盖更彻底
if refreshClaims, err := utils.GetRefreshClaims(c); err != nil {
    global.Log.Warn("登出时 x-refresh-token 不可用，改用 Access Token 的 jti 兜底", zap.Error(err))
} else {
    jti = refreshClaims.ID
    blacklistTTL = utils.ParseRefreshExp(refreshClaims.ExpiresAt)
}
```

**为什么这个兜底成立？** 因为第四节讲的共享 jti——两者 jti 相同，拿 Access 的 jti 写黑名单等价于作废整个会话，包括那个没解析出来的 Refresh Token。如果这里只是 `return` 不管，用户登出后**那个还没过期的 Access Token 依然有效，最多能再用 2 小时**。这不是体验问题，是安全问题。还有一处顺序讲究——**清理动作绝不能因为写黑名单失败而被短路**：

```go
// server/service/user.go
// 无论黑名单能不能写，都保证 Cookie 清掉、Redis 会话删掉，用户不会"退不出去"
utils.ClearRefreshToken(c)
global.Redis.Del(uuid.String())
if jti == "" || blacklistTTL <= 0 {
    global.Log.Warn("登出时没有可用的 jti 或剩余有效期，跳过黑名单写入")
    return
}
if err := ServiceGroupApp.JwtService.JoinInBlacklist(jti, blacklistTTL); err != nil {
    global.Log.Error("登出写入黑名单失败", zap.Error(err))
}
```

如果写成"先写黑名单，失败就 return"，Redis 抖一下用户的 Cookie 就清不掉，会卡在"点了登出但还是登录状态"的诡异状态。**先做不可回退的清理，再做尽力而为的拉黑**，这个顺序是对的。顺便说一句，`Logout` 保持"尽力而为、不返回错误"是刻意的——即使黑名单写失败，接口也返回登出成功，因为本地 Cookie 已经清掉了，对用户来说登出确实完成了。

### 坑三：拿别人的寿命给自己的令牌上黑名单

**这是三个坑里最隐蔽的，因为它平时完全正常。** `TokenNext` 里踢旧会话时，黑名单 TTL 这么算：

```go
// 修复前
remainExp := utils.ParseRefreshExp(refreshClaims.ExpiresAt)
jwtService.JoinInBlacklist(oldJti, remainExp)
```

而 `refreshClaims` 是**刚刚创建的新令牌**的 claims，`ExpiresAt` 是"现在 + 7 天"。也就是说，**我用新令牌的寿命去给旧令牌上黑名单。** 平时为什么没问题？因为算出来的是新令牌的**完整**寿命（7 天），而旧令牌剩余寿命一定**小于** 7 天，所以常态下这只是"过度拉黑"——旧令牌明明 1 小时后自然失效，却在黑名单里躺 7 天，代价只是 Redis 里多留几个到点自动过期的无用 key，看不出异常。

**真正的危险场景需要配置被调短。** 想象这条时间线：① 用户在电脑上登录，当时 `refresh_token_expiry_time: 7d`，Redis 里 `uuid → jti_old` 的 TTL ≈ 7 天；② 管理员在后台把有效期收紧成 `2h`（很合理的操作，热更新立刻生效，不用重启）；③ 用户在手机上再登录 → 走"踢旧会话"分支 → `remainExp` 算的是**新**令牌的寿命 = **2 小时** → 旧 jti 只被拉黑 2 小时；④ **2 小时后黑名单 key 自动过期，而电脑上那个旧 Refresh Token 还剩 6 天多的有效期。** 结果是**旧会话"复活"了**——被"强制下线"的设备可以重新换发 Access Token 继续使用，异地登录拦截形同虚设。这个问题在测试环境极难复现，它需要"签发旧令牌时的有效期"和"踢人时的有效期"不一致。**修法：TTL 必须取被拉黑那张令牌自己**的剩余寿命。而旧令牌的剩余寿命就藏在 Redis 里：**`uuid` 这个会话 key 的 TTL**。因为 `SetRedisJWT` 当初就是用"当时配置的 Refresh 有效期"把 key 写进去的，TTL 一直在倒数，任意时刻读出来就是那张 Refresh Token 的真实剩余寿命。

```go
// server/api/user.go
// 旧会话的剩余寿命 = Redis 中 uuid 这个 key 的 TTL，必须在 SetRedisJWT 覆盖它之前读取
oldTTL, ttlErr := global.Redis.TTL(user.UUID.String()).Result()
if ttlErr != nil {
    global.Log.Warn("读取旧会话剩余有效期失败，改用黑名单默认有效期", zap.Error(ttlErr))
    oldTTL = 0 // 交给 JoinInBlacklist 用配置里的 Refresh 有效期兜底
}
if err := jwtService.JoinInBlacklist(oldJti, oldTTL); err != nil {
    // ...
}
```

**最要紧的是执行顺序**——读 TTL → 写黑名单 → 覆盖会话记录，第 3 步会把 key 的 TTL 重置成新令牌的 7 天，顺序颠倒就等于没修。**关于失败方向**：如果 `uuid` 已不存在，`TTL` 返回 `-2`；如果 key 存在但没设过期时间，返回 `-1`，这两种都会落进 `JoinInBlacklist` 的 `expTime <= 0` 分支，用配置里的 Refresh 有效期兜底，**失败方向是"多拉黑"而不是"少拉黑"**，这个方向是安全的——过度拉黑只浪费几个 key，拉黑不足会让旧会话复活。

修完后我顺手把冻结用户那条路径统一成同一个手法：

```go
// server/service/user.go
jti, _ := ServiceGroupApp.JwtService.GetRedisJWT(user.UUID)
exp, _ := global.Redis.TTL(user.UUID.String()).Result()
if jti != "" {
    if err := ServiceGroupApp.JwtService.JoinInBlacklist(jti, exp); err != nil {
        global.Log.Error("冻结用户时写入黑名单失败", zap.Uint("user_id", user.ID), zap.Error(err))
    }
}
global.Redis.Del(user.UUID.String())
```

原来是 `JoinInBlacklist(jti, 0)`，无脑用配置里的 7 天兜底；现在读会话 key 的真实 TTL，黑名单有效期和令牌真实寿命严格对齐。同时错误从 `_ =` 丢弃改成记日志——冻结时如果黑名单写失败，**管理员以为冻结生效了，实际那个用户的令牌还能继续用最多 2 小时**，这种事必须留痕。

---

## 九、三条下线路径，一个出口

需要"提前作废某个会话"的场景一共三个，最终都收敛到同一个函数：

| 场景 | 入口 | jti 来源 | TTL 来源 |
|---|---|---|---|
| 用户主动登出 | `UserService.Logout` | Refresh 的 jti，失败时退化到 Access 的 jti | 令牌自身的剩余寿命 |
| 异地登录踢旧会话 | `UserApi.TokenNext` | Redis 里存的旧 jti | 会话 key 的 TTL（必须在覆盖前读） |
| 管理员冻结用户 | `UserService.UserFreeze` | Redis 里存的 jti | 会话 key 的 TTL |

三条路径的差异只在**怎么拿到 jti 和 TTL**，一旦拿到动作完全一样：`JoinInBlacklist(jti, ttl)`。这种"多入口、单出口"的收敛比代码复用本身更有价值——它意味着黑名单的写入规则（前缀、摘要、值、TTL 校验）**只有一处实现**，坑一里那个"零值写出永久 key"的问题只需要在一个地方修，而不是三个地方各修一遍。解冻则不需要写黑名单，只清会话记录（`UserUnfreeze` 里就是 `Update("freeze", false)` 加一次 `Redis.Del`）。因为冻结时已经把旧会话拉黑了，解冻后用户重新登录会拿到全新 jti——**冻结-解冻不试图"恢复"原会话，而是让用户重新登录**，状态更干净。

---

## 十、前端要做的事很少

后端设计得好的一个标志是：前端的配合代码很短。请求拦截器只做一件事——把内存里的 Access Token 塞进请求头：

```ts
// web/src/utils/request.ts
config.headers = {
    'Content-Type': 'application/json',
    'x-access-token': userStore.state.accessToken,
    ...config.headers,
}
```

注意 `x-access-token` 被放在展开 `...config.headers` **之前**，这样个别接口传特殊头时不会覆盖掉认证头——顺序反了会导致那些接口莫名 401。响应拦截器负责接收续签结果：

```ts
service.interceptors.response.use((response: AxiosResponse) => {
    const userStore = useUserStore()
    if (response.headers['new-access-token']) {
        userStore.state.accessToken = (response.headers['new-access-token'])
    }
    return response.data
})
```

就这么几行。中间件续签成功后通过 `new-access-token` 响应头下发新令牌，前端拦截器替换状态里的旧值，**用户和业务代码都不需要知道刚才发生过一次续签**。这里有个我特意保留的设计：**Access Token 只放在 Pinia 状态里，不落 localStorage**，代价是刷新页面后内存清空、Access 丢失。但这不会导致掉线——因为 Refresh Cookie 还在：页面重载时第一个需要鉴权的请求带着空 `x-access-token` 和有效的 Refresh Cookie 发出去，中间件走"Access 为空"的分支完成续签，用户感觉不到。也就是说，**"刷新页面仍保持登录"这个体验是由 Refresh Cookie 提供的，而不是靠持久化 Access Token**，这个分工让长期凭据只存在于 HttpOnly Cookie 一个地方，XSS 偷不到。

---

## 小结

这套方案的核心是**把"会话"从令牌里抽出来，变成 Redis 里的一等公民**：一个会话 = 一个 jti = 一个黑名单 key = 一条 `uuid → jti` 记录；Access 是 2 小时一换的短期凭证；Refresh 是只在 HttpOnly Cookie 里的长期凭证，可被服务端主动作废；中间件把续签做得完全透明，前端只把响应头里的新令牌存回内存。最终链路：

```
请求 → x-access-token 解析成功？
        ├─ 是 → EXISTS jwt:BlackList:<sha256(jti)>
        │       ├─ 不在 → 放行，claims 写入 Context
        │       └─ 在   → 401 + 清 Cookie
        └─ 否 → 错误是"过期"或"Access 为空"？
                ├─ 是 → 用 Refresh 换发新 Access（复用 jti）→ 响应头下发 → 放行
                └─ 否 → 401 + 清 Cookie
```

还有哪些没解决、我知道但暂时接受的。**单点 Redis**：黑名单和会话记录都在一个实例上，它挂了"主动下线"能力就没了；目前中间件对 Redis 错误是**快速失败**（返回 `Service unavailable` 而不是放行），这是有意选择，鉴权环节宁可拒绝服务，也不能把"黑名单查不到"当成"没有黑名单"，要真正解决得上哨兵或集群。**没有刷新令牌轮换**：现在同一个 Refresh Token 在 7 天内可反复使用（只要不触发踢旧会话），更严格的做法是每次续签都换发新 Refresh 并作废旧的，这样旧的一旦被重放就能立刻发现，代价是每次续签都要写 Redis，还要考虑换发失败时的事务性。**改密码后旧令牌不失效**：改密码走的是校验原密码再更新，但没有把现有会话写进黑名单，要做到"改密码即下线所有设备"，需要一个按用户维度作废会话的机制——比如给用户加 `token_version` 字段，签发时写进 claims、校验时比对，改密码就自增。这三个坑如果只能记一个，我希望是坑一里那句：**Redis 的过期时间传 0 是"永不过期"，不是"立刻过期"**。这个语义陷阱在任何用 Redis 做过期控制的系统里都存在，后果往往是"资源缓慢泄漏"，等发现时已经积累了几十万个永久 key。而坑三更值得警惕的地方在于——它平时完全正常，只在配置变更后才暴露。**这类"条件性正确"的代码，比直接报错的代码危险得多。**
