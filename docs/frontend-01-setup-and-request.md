# 个人博客前端（一）：工程组织与请求层设计

这个博客的前端是 Vue 3 + TypeScript + Vite 的前后端分离项目，同时承载两套界面：面向访客的博客前台，和面向站长的管理后台。两套界面共享同一套组件库、状态管理和请求层，但页面形态、权限要求和交互密度差别很大。这一篇不谈具体页面，只谈两件更底层的事：**项目怎么组织**，以及**HTTP 请求怎么收口**。前者决定三个月后你还能不能快速找到文件，后者决定 token 续签、错误提示、强制重登这些横切关注点会不会散落到几十个组件里。读完可以带走三样东西：一套「页面放哪、组件放哪」的判断标准而不是背目录名；一个把 HTTP 细节收口到单个 axios 实例的完整写法，含 Access Token 续签在前端的落点；以及几处真实的取舍与坑。

## 一、目录结构

### 1.1 顶层划分

```
web/src
├─ api/         每个后端模块一个文件，只放"请求函数 + 类型"
├─ assets/      全局样式（目前只有 base.css）
├─ components/  可复用组件，按"用途"分五类
├─ router/      路由表 + 全局守卫（单文件）
├─ stores/      Pinia，按业务域拆分
├─ utils/       跨域基建（request.ts 等）
├─ views/       页面，分 web/ 与 dashboard/ 两套
├─ App.vue / main.ts
```

判断一个 `.vue` 放哪，两个问题就够了：**它会被路由直接访问吗**（会 → `views/`）；**只有某一个页面在用、且耦合很紧吗**（是 → 留在那个页面目录；会被两处以上复用或本身就是独立业务块 → `components/`）。顺着这条标准，首页里的 `ArticleList`、`TagCloud` 该进 `components/pages/`——它们既被首页用，也能被别的页面用；而只服务于搜索页的内部结构，留在 `views/web/search/` 更省事。

### 1.2 views 分成两套

```
views/
├─ web/          博客前台
│  ├─ index.vue          前台布局（只提供页脚）
│  ├─ index/index.vue    首页
│  ├─ article/index.vue  文章详情
│  ├─ search/  news/  friend-link/  about/
│  └─ test/              早期 demo，路由已注释
├─ dashboard/    管理后台
│  ├─ index.vue          后台布局（侧边栏 + 顶栏 + 内容区）
│  └─ home/  articles/  images/  users/  system/  user-center/
├─ login/  error/（404 页）
```

两套布局职责完全不同。`views/web/index.vue` 只包了一层 `el-container` 和 `<web-footer/>`，然后就是 `<router-view/>`——**导航栏不在这里**，而是每个前台页面自己 `<WebNavbar />`（`web/index/index.vue`、`web/search/index.vue`、`web/article/index.vue` 等 7 处都是这么写的）。好处是页面能自由决定整页结构，代价是新增前台页面时容易漏掉导航栏。`views/dashboard/index.vue` 则是完整骨架：`el-aside` 放 `Logo` 和 `DashboardMenu`，`el-header` 放 `Breadcrumb` 和 `AuthPopover`，再用 `DashboardTag` 做多标签页，`<router-view/>` 落在 `el-main` 里——后台页面形态高度一致，骨架一次给足更划算。另外 `views/web/article/index.vue` 对应的是**顶层路由** `/article/:id`，不挂在前台布局下，这样它不被布局的容器和页脚约束。

### 1.3 components 的三个判断维度

```
components/
├─ common/    跨页面通用件：AuthPopover / UserCardPopover / CommentItem
├─ forms/     弹窗表单：LoginForm / ArticleCreateForm / ...
├─ layout/    布局与导航：WebNavbar / WebFooter / Carousel / Breadcrumb / DashboardMenu
├─ pages/     内容区块：ArticleList / TagCloud / ProfileCard / DailyNews / HotArticles / ...
├─ widgets/   纯展示小单元：Logo / UserCard / UserActivityChart
└─ hero/      首页 Hero 备选方案（HeroPage01~04，路由已注释）
```

这套分类按**用途**切，不按技术类型切：

- **`forms/` 单独拎出来**，因为这类组件遵循同一约定：自身持有一份表单数据、自己提交、自己关闭弹窗；由父级的可见性开关控制挂载，不对外暴露复杂 props。集中放置，将来复用表单结构时不用满仓库找。
- **`pages/`** 是"首页的一块砖"（`TagCloud`、`ProfileCard`、`RecentComments`），各自带取数逻辑、互不依赖，因此可以自由增删排序。
- **`common/` 与 `widgets/` 的区别在粒度**：`common/` 是"带行为的小组件"（`CommentItem` 内建了回复框、删除、表情选择），`widgets/` 是纯展示单元（`Logo`、`UserCard`）。

历史遗留：`components/hero/` 的四个 Hero 组件和 `views/web/test/` 的几个 demo，路由里都是注释状态，不及时清理会持续干扰"哪些是活代码"的判断。

### 1.4 其余四个目录

- **`api/`**：11 个文件，一个后端模块一个。除 `common.ts` 只放类型外，其余 10 个都从 `@/utils/request` 引入同一个 axios 实例。这一层的产出物只有**请求函数**和**入参/出参类型**，不含业务逻辑和状态。
- **`router/`**：单文件同时容纳路由表和守卫。规模不大时比拆成 `routes.ts` + `guard.ts` 更好读——守卫和它依赖的 `meta` 约定挨在一起，改一处不会忘另一处。
- **`stores/`**：见第四节。
- **`assets/`**：`base.css` 只做了 `box-sizing` 重置，样式都写在各组件的 `<style scoped lang="scss">` 里，没有设计令牌层，这是目前最大的技术债之一。

## 二、工程配置

### 2.1 vite.config.ts：别名、代理、监听地址

```ts
// web/vite.config.ts
export default defineConfig({
  plugins: [vue(), AutoImport({resolvers: [ElementPlusResolver()]}), Components({resolvers: [ElementPlusResolver()]})],
  resolve: { alias: {'@': fileURLToPath(new URL('./src', import.meta.url))} },
  server: {
    host: "0.0.0.0", port: 80,
    proxy: {
      "/api":     {target: env.VITE_SERVER_URL, changeOrigin: true},
      "/uploads": {target: env.VITE_SERVER_URL, changeOrigin: true},
    }
  }
})
```

- **`@` 别名**指向 `src`，同时要在 `tsconfig.app.json` 的 `paths` 里配一份。只配一边的话，要么构建能过但 TS 报错，要么反过来。
- **代理不做 rewrite**，`/api` 原样转发——后端路由前缀本身就配成了 `api`。`/uploads` 也要代理，因为封面和插图存的是 `/uploads/image/xxx.png` 这样的相对路径，由后端作为静态资源提供；不代理的话前端 dev server 会把它当成自己的资源去找，直接 404。
- **`host: "0.0.0.0"`** 让 dev server 监听所有网卡，手机平板和同局域网的机器都能打开；`port: 80` 则让本地地址不带端口号。

### 2.2 .env 里两个变量各管一段

```
# web/.env
VITE_SERVER_URL = http://127.0.0.1:8080
VITE_BASE_API = /api
```

| 变量 | 谁在用 | 用途 |
|---|---|---|
| `VITE_SERVER_URL` | `vite.config.ts` | 仅作为 dev server 的**代理目标**（后端地址） |
| `VITE_BASE_API` | `src/utils/request.ts` | axios 实例的 `baseURL`，即浏览器实际发出的路径前缀 |

分工是刻意的：**浏览器永远只知道 `/api`，不知道后端在哪台机器**。跨域那一跳由 dev server 代劳，所以开发期后端不用开 CORS；换后端地址只改 `.env`，业务代码零改动。加载用 `loadEnv("", process.cwd())`，第一个参数是前缀过滤，传空串表示不过滤，因此配置里能拿到全部变量。而 `env.d.ts` 只显式声明了一个变量：

```ts
// web/env.d.ts
export interface ImportMetaEnv { VITE_SERVER_URL: string }
```

那 `import.meta.env.VITE_BASE_API` 为什么也能通过类型检查？因为 Vite 自带的 `ImportMetaEnv` 有一个 `Record<string, any>` 的兜底索引签名——**没声明的变量不报错，但类型是 `any`**。所以只声明一个不是遗漏，而是"只给需要强类型的那一个上锁"。

### 2.3 自动按需引入：收益与代价

```ts
AutoImport({resolvers: [ElementPlusResolver()]}),   // 自动引入 API
Components({resolvers: [ElementPlusResolver()]}),   // 自动注册组件
```

**收益**是两件事同时达成：模板里写 `<el-card>`、`<el-tabs>` 不用 import；脚本里写 `ElMessage.error(...)`、`ElMessageBox.confirm(...)` 也不用 import。打开 `src/utils/request.ts` 会发现它从头到尾没引入 `ElMessage` 却直接在用——这是被自动注入的。同时按需引入让只有真正用到的组件代码进包，而不是整个 Element Plus。**代价**是依赖变"看不见"了，有三个表现：

1. **静态检索失效**——想找"哪些地方弹过提示"，`grep "import.*ElMessage"` 搜不到，只能直接搜 `ElMessage`。
2. **报错信息变间接**——哪次构建没生成 `auto-imports.d.ts`，编辑器就满屏 `Cannot find name 'ElMessage'`，而代码本身没问题，问题在生成物。
3. **多出两份必须入库的生成文件**，分别是 `unplugin-auto-import` 与 `unplugin-vue-components` 的产物。`web/auto-imports.d.ts` 只声明 `ElMessage` 和 `ElMessageBox` 两个全局常量（Element Plus 里"是函数而非组件"的 API 主要就这两位）；`web/components.d.ts` 则同时登记了 Element Plus 组件（`ElAnchor`、`ElAside`……）和项目自己的组件（`ArticleCreateForm`、`ArticleList`……），因为 `unplugin-vue-components` 默认也会扫描 `src/components`——不过项目里实际写法仍是**显式 import 自己的组件**，所以这份声明更像"能力登记表"。两个文件顶部都带 `// @ts-nocheck`，是生成物，**不要手改**。

### 2.4 图标全局注册，所以能写 `is="View"`

```ts
// web/src/main.ts
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
```

这段循环把所有图标注册成全局组件，于是模板里可以直接用字符串解析（`web/src/components/pages/ArticleList.vue`）：

```vue
<el-icon><component is="View"/></el-icon>{{ scope.row._source.views }}
```

`is` 传的是字符串，Vue 会在**全局注册表**里查 `View`，所以这种写法必须搭配全局注册，光靠 `import` 无法用字符串写组件名。**取舍**是：全局注册换来渲染自由（尤其在循环里按数据决定图标时很方便），代价是**所有图标都会进包**，摇不掉。图标库不大，这个代价目前可接受；真成问题时改成按需引入 + 组件内显式 import，代价是失去 `is="字符串"` 的写法。

### 2.5 index.html 里的小技巧与一个真实隐患

```html
<!-- web/index.html -->
<link rel="icon" href="http://127.0.0.1:80/api/website/logo">
<title></title>
<script>
  fetch('/api/website/title').then(r => r.json()).then(d => { document.title = d.title })
</script>
```

这个 `fetch` 解决一个实际问题：**站点标题存在后端配置里**，而 `index.html` 是静态文件，构建时并不知道标题。所以在首屏 HTML 里直接发一个请求写进 `document.title`——不用等 Vue 挂载，标签页标题会尽快变成正确值。它能这么写，是因为后端的 `/website/title` 是个**刻意的例外**：不走项目统一的 `{code, msg, data}` 包装，而是直接返回 `{"title": "..."}`（见 `server/api/website.go` 的 `WebsiteTitle`）；否则这里的 `data.title` 要写成一长串 `data.data.title`。用一个"不统一"的接口换一段极简的内联脚本是划算的，但**它必须被记录下来**——否则以后有人按统一规范"顺手修正"这个接口，首页标题会静默变成 `undefined`。

隐患在 favicon 那一行：`http://127.0.0.1:80/...` 是**写死的本地地址**，上线后每个访客的浏览器都会去请求他自己机器的 80 端口，favicon 必然拿不到。应改成相对路径（`/api/website/logo`），或在构建时用环境变量替换主机名。另外 `<style>` 块被放在了 `</html>` 之后，浏览器通常会挪进 `body` 从而"看起来正常"，但这不是合法结构。

## 三、请求层：一个 axios 实例收口所有 HTTP 细节

`src/utils/request.ts` 是整个前端唯一创建 axios 实例的地方，`api/` 下所有函数都走它。这一节是全文重点。

### 3.1 实例与类型约定

```ts
// web/src/utils/request.ts
const service = axios.create({ baseURL: import.meta.env.VITE_BASE_API, timeout: 10000 })

export interface ApiResponse<T> { code: number; msg: string; data: T }
```

后端响应统一是 `{code, msg, data}`：`code === 0` 为业务成功，非 0 为业务失败（HTTP 状态码仍是 200）。把这条约定固化成 `ApiResponse<T>` 后，api 函数的返回类型就能写成 `Promise<ApiResponse<具体类型>>`，调用处统一 `if (res.code === 0)`。分页与检索相关的类型放在 `src/api/common.ts`：

```ts
// web/src/api/common.ts
export interface PageInfo { page: number; page_size: number }
export interface PageResult<T> { list: T[]; total: number }
export interface Hit<T> { _id: string; _source: T }
```

`PageInfo` / `PageResult<T>` 对应关系型分页，`Hit<T>` 对应 Elasticsearch 的命中结构——ES 的文档 id 在 `_id`、字段在 `_source` 里，形状和普通实体不同，所以类型上也分开表达。

### 3.2 请求拦截器：注入 token 与展开顺序

```ts
service.interceptors.request.use(
    (config: AxiosRequestConfig) => {
        const userStore = useUserStore();
        config.headers = {
            'Content-Type': 'application/json',
            'x-access-token': userStore.state.accessToken,
            ...config.headers,
        }
        return config as InternalAxiosRequestConfig
    },
    (error: AxiosError) => { /* 弹错误提示并 reject */ }
)
```

这是全项目 Access Token 的**唯一注入点**，api 函数只需声明 url 和 method，鉴权头由这一处补齐。关键在于 `...config.headers` 放在**最后**：对象字面量里后写的键覆盖先写的，所以这个顺序等于"默认值在前、调用方显式传入的头优先"。反过来写则会埋一个运行时才暴露的坑——将来任何需要发 `multipart/form-data` 的接口都会被这里强行改回 `application/json`，而且不报错。诚实的现状是：目前 `api/` 下 10 个文件**没有任何一处传自定义 headers**，所以这个顺序眼下是防御性的，价值在于为将来留出覆盖能力，成本为零。顺带一提，图片上传在 `views/dashboard/articles/article-publish.vue` 里用的是**裸 `axios.post`**，绕过了这个实例——拦截器提供的 token 注入、响应头回写、统一错误提示都不覆盖它。这类"逃逸"是请求层收口时最容易出现的漏洞。

### 3.3 响应拦截器：业务码、续签落地点、强制重登

```ts
service.interceptors.response.use(
    (response: AxiosResponse) => {
        const userStore = useUserStore()
        if (response.headers['new-access-token']) {
            userStore.state.accessToken = (response.headers['new-access-token'])
        }
        if (response.data.code !== 0) {
            ElMessage.error(response.data.msg)
            if (response.data.data && response.data.data.reload) {
                userStore.reset()
                localStorage.clear()
                const layoutStore = useLayoutStore()
                router.push({name: 'index', replace: true}).then(() => {
                    layoutStore.state.popoverVisible = true
                    layoutStore.state.loginVisible = true
                })
            }
        }
        return response.data
    },
    /* 错误分支见 3.4 */
);
```

**① Access Token 的自动续签就落在第一行。** 后端在 Access 过期时用 Refresh Token 换发新的，通过响应头 `new-access-token` 下发；前端不发起任何"续签请求"，只在**每次响应**时检查这个头，有就写回 store，下一次请求的请求拦截器自然会带上新 token，整个过程对业务代码透明。这个设计省掉一次往返：与其让前端发现 401 再发续签请求、然后重放原请求（典型的双请求模式），不如让后端在原本就要返回的响应里捎带新 token，前端只当搬运工。代价是必须相信"每个响应都可能携带新 token"，所以这段检查必须在响应拦截器最前面。

**② 业务失败统一提示。** 只要 `code !== 0` 就弹 `ElMessage.error(response.data.msg)`，各页面不用再写一遍提示逻辑。注意这里**没有 return 也没有 reject**，提示完仍走到 `return response.data`，所以约定是"先看 `res.code`"——它把"提示"和"流程控制"分开了：调用方可以按需覆盖，比如表单校验失败时自己展示更详细的错误。

**③ `reload` 标记触发强制重登。** 后端在账号被冻结、会话被踢等场景返回带 `reload` 标志的数据，前端识别后重置用户 store、清空 localStorage、跳回首页并弹出登录框。注意这里调用的是自定义的 `userStore.reset()`——这是正确写法，原因见 3.5。

**④ 返回的是 `response.data` 而不是 `response`**，所以调用处拿到的是 `{code, msg, data}` 而非 `AxiosResponse`。这也解释了为什么 `ApiResponse` 定义在 `request.ts` 里：它和拦截器是一体的，类型描述的就是拦截器"剥壳"之后的形状。

### 3.4 错误分层：为什么网络错误和 500/404/403 要分开

```ts
(error: AxiosError) => {
    if (!error.response) {
        ElMessageBox.confirm(`<p>检测到请求错误</p><p>${error.message}</p>`, '请求报错', {
            dangerouslyUseHTMLString: true,
            confirmButtonText: '稍后重试', cancelButtonText: '取消',
        }).then()
        return Promise.reject(error)
    }

    switch (error.response.status) {
        case 500: return handleSpecificError(500, error)
        case 404: return handleSpecificError(404, error)
        case 403: return handleSpecificError(403, error)
    }
    return Promise.reject(error)
}
```

**没有 `error.response`** 意味着请求根本没走到服务端：断网、后端没启动、DNS 失败、超时（实例设了 10 秒）、或 CORS 被拦。这类问题与具体接口无关，所以提示只说"请求错误 + 原始 message"，并给「稍后重试」——因为重试确实可能好。

**有 `error.response`** 说明后端给了明确状态码，不同状态码指向完全不同的排查方向，所以每种都给了定制文案（`handleSpecificError` 里以 HTML 字符串维护）：

| 状态码 | 文案里的判断 | 引导动作 |
|---|---|---|
| 500 | "通常由后台服务器发生不可预料的错误（如 panic）引起，请先查看后台日志" | 清理缓存 / 重新登录 |
| 404 | "通常表示接口未注册（或服务未重启），或请求路径（方法）与 API 不符" | 检查 URL 与方法 |
| 403 | "您没有权限访问此路由（admin）" | 确认当前用户角色 |

价值不在措辞，而在**把排查方向直接写进提示**。个人项目没有监控和值班，出问题时盯着屏幕的就是自己，"看到提示就知道下一步去哪看"能省掉大量回想时间——比如 404 那句"或服务未重启"，对应的就是改完后端路由忘重启的常见情形。

### 3.5 这一层的一个真实坑：`$reset` 只对 option store 有效

接着 3.4 看「清理缓存」按钮的回调：

```ts
// web/src/utils/request.ts
ElMessageBox.confirm(errorMessages[status], '接口报错', { /* ... */ }).then(() => {
    const userStore = useUserStore()
    userStore.$reset()          // ← 问题在这一行
    localStorage.clear()
    router.push({name: 'index', replace: true})
});
```

而 `user` store 是 setup 语法定义的，导出的是自定义的 **`reset`**，不是 `$reset`：

```ts
// web/src/stores/user.ts
export const useUserStore = defineStore('user', () => {
    const state = ref(initState())
    const reset = () => { /* 手写的重置逻辑 */ }
    return {state, reset, loginIn, /* ... */};
});
```

Pinia 对 `$reset` 的实现是：

```js
// web/node_modules/pinia/dist/pinia.js
const $reset = isOptionsStore ? function $reset() { /* 用 state() 重建初始状态 */ }
    : process.env.NODE_ENV !== "production"
        ? () => { throw new Error(`🍍: Store "${$id}" is built using the setup syntax and does not implement $reset().`) }
        : noop;
```

也就是说，对 setup store 调用 `$reset`：**开发环境**直接抛错，`.then()` 里后面的清缓存和跳转都不执行——用户点了什么都没发生，只有控制台多一条红字；**生产环境** `$reset` 是 `noop`，静默什么都不做，但后续清缓存与跳转照常执行。同一个问题在两种环境下表现不同，是这类 bug 最难查的地方。修复只需一个字符：`userStore.$reset()` → `userStore.reset()`，和 3.3 里 `reload` 分支用的是同一个方法。**教训**是用 setup 语法写 store 时 Pinia 不会自动提供重置能力，必须自己实现，并且全程统一用自己实现的那个名字；想彻底避免混用，最省事的办法是把自己的重置函数直接命名为 `$reset`，覆盖掉 Pinia 的默认实现。

## 四、状态管理：按域拆分

### 4.1 五个 store 的分工

```
stores/
├─ index.ts   创建 pinia 实例 + 注册插件
├─ user.ts    用户信息、accessToken、登录/注册/登出
├─ layout.ts  纯 UI 状态：弹窗可见性 + 表格刷新标志 + 侧边栏折叠
├─ website.ts 站点配置（标题、logo、备案号等）
├─ tag.ts     后台多标签页
└─ article.ts （目前是空壳，内容全部注释）
```

拆分依据是**"谁会改它"**：`user` 的变更都来自认证流程，集中放便于一次清干净。`layout` 装界面状态，这里有个值得注意的设计：**把"弹窗可见性"和"表格刷新标志"统一放在一个 store 里**——`articleCreateVisible`、`shouldRefreshArticleTable`、`shouldRefreshCommentList` 这些标志让"表单提交成功 → 关闭弹窗 → 通知列表刷新"这条链路不用层层 emit：表单改开关，列表 `watch` 开关重新取数。代价是它会不断膨胀成"全局开关抽屉"，字段多了以后很难看出哪些还有用。`website` 的初始化方式比较特别，通过 pinia 插件在 store 创建时自动取数：

```ts
// web/src/stores/index.ts
pinia.use(({ store }) => {
    if (store.$id === 'website') { store.initializeWebsite() }   // 确保在 store 创建时调用
})
```

配合内部的 `websiteInfoInitialized` 标志做幂等，效果是"任何组件第一次 `useWebsiteStore()`，站点配置就已经在路上了"，省掉了在每个页面 `onMounted` 里重复初始化。

### 4.2 accessToken 只放内存的取舍

```ts
// web/src/stores/user.ts
function initState() {
    const userInfo = ref<User>({ /* ... */ })
    const savedIsUserLoggedInBefore = localStorage.getItem('isUserLoggedInBefore');
    return {
        userInfo,
        accessToken: '',                                            // 不持久化
        userInfoInitialized: false,
        isUserLoggedInBefore: savedIsUserLoggedInBefore === 'true'   // 持久化
    }
}
watch(() => state.value.isUserLoggedInBefore, (newIsUserLoggedInBefore) => {
    localStorage.setItem('isUserLoggedInBefore', String(newIsUserLoggedInBefore));
})
```

这里做了两个相反的决定：**`accessToken` 只存在内存里**，刷新页面就丢，任何脚本（包括 XSS）也读不到 localStorage 里的令牌；**`isUserLoggedInBefore` 写进 localStorage**，它只是一个布尔标记"这台机器上曾经登录过"。配合 `initializeUserInfo()` 的惰性取数（用 `userInfoInitialized` 保证只跑一次），刷新后的流程是：读到标记为 `true` → 主动请求一次 `/user/info` → 用返回数据恢复 `userInfo`。Access Token 不需要恢复——Refresh Token 存在 HttpOnly Cookie 里，后端中间件会在 Access 过期时自动换发并用 `new-access-token` 响应头送回来（就是 3.3 那段）。

**这个组合的关键在于"登录态"有两个不同判据**：

```ts
// web/src/stores/user.ts
const isLoggedIn = computed(() => state.value.userInfo.role_id !== 0);
const isAdmin    = computed(() => state.value.userInfo.role_id === 2);
```

判断是否登录看的是 **`role_id !== 0`**（来自 `/user/info` 的真实数据），而不是"有没有 accessToken"。这样刷新的那一瞬间不会因为 token 是空串而误判成未登录——token 空不空无所谓，只要 `/user/info` 能拿到数据就是登录态。很小的设计，但把"令牌"和"身份"解耦了，避免大量"刷新后闪一下未登录"的问题。代价是 `role_id` 用了魔法数字（`0` = 游客、`2` = 管理员），散落在 computed 和守卫里，定义成常量更好。

## 五、路由与守卫

### 5.1 路由结构

```ts
// web/src/router/index.ts
const routes = [
  { path: '/', name: 'web', component: () => import('@/views/web/index.vue'),
    children: [
      { path: '/',      name: 'index',  component: () => import('@/views/web/index/index.vue'),  meta: {title: '首页'} },
      { path: 'search', name: 'search', component: () => import('@/views/web/search/index.vue'), meta: {title: '搜索'} },
      { path: 'news',   name: 'news',   /* ... */ },
      { path: 'about',  name: 'about',  /* ... */ },
    ]},
  { path: '/login',       name: 'login',   component: () => import('@/views/login/index.vue') },
  { path: '/article/:id', name: 'article', component: () => import('@/views/web/article/index.vue') },
  { path: '/dashboard',   name: 'dashboard', component: () => import('@/views/dashboard/index.vue'),
    meta: {title: '控制面板', requiresAuth: true},
    children: [ /* home / user-center / users / articles / images / system */ ]},
  { path: '/404',           name: '404', component: () => import('@/views/error/index.vue') },
  { path: '/:catchAll(.*)',              component: () => import('@/views/error/index.vue') },
]
```

四个要点：**全部 `() => import(...)` 懒加载**，首屏只下载首页那块代码，后台几十个页面在点进去之前不加载；**`meta.title` 只写不用**——路由里几乎每个都标了中文标题，但项目里没有读取它去设 `document.title`（标题靠 2.5 那个 fetch），这是现成的扩展点；**`/article/:id` 是顶层路由**，理由见 1.2；**两级兜底**，`/404` 给一个显式地址、`/:catchAll(.*)` 兜住所有未匹配路径，两者指向同一组件。

### 5.2 meta 与守卫

为了让 `meta` 有类型，`web/env.d.ts` 对 vue-router 的 `RouteMeta` 做了声明合并：

```ts
// web/env.d.ts
declare module 'vue-router' {
    interface RouteMeta { requiresAuth?: boolean; requiresAdmin?: boolean; title?: string }
}
```

没有这段的话 `route.meta.requiresAuth` 会被推断成 `unknown`，每次用都要断言。这是 TS 项目里很值得做的一件小事：**把"约定"写成类型**。守卫的逻辑是"先恢复身份，再按 meta 判断"：

```ts
// web/src/router/index.ts
router.beforeEach((to, from, next) => {
    const userStore = useUserStore()
    const layoutStore = useLayoutStore()
    userStore.initializeUserInfo().then(() => {
        const isAuthenticated = userStore.isLoggedIn
        const isAdmin = userStore.isAdmin

        if (to.matched.some(record => record.meta.requiresAuth)) {
            if (!isAuthenticated) {
                ElMessageBox.confirm('登录已过期，需要重新登录，是否跳转到登录页？', 'Warning', { /* ... */ })
                    .then(() => { router.push({name: 'index', replace: true}) /* 并弹出登录框 */ })
                    .catch(() => { router.push({name: from.name as string}) })   // ← 见 5.3
            } else if (to.matched.some(record => record.meta.requiresAdmin) && !isAdmin) {
                ElMessageBox.confirm('权限不足，请确认您的用户角色是否具备访问该页面的权限。', 'Warning', { /* ... */ })
                    .then(() => { router.push({name: from.name as string}) })   // ← 见 5.3
            } else { next() }
        } else { next() }
    })
})
```

- **用 `to.matched.some(...)` 而不是 `to.meta.requiresAuth`**。`matched` 是匹配到的整条路由链，所以只要父级（如 `/dashboard`）标了 `requiresAuth`，它下面所有子路由都自动受保护，不必给几十个子路由逐个加 `meta`。这是把权限声明写在**层级**而非**叶子**上的好处。
- **权限分两级**：`requiresAuth` 只要求登录，`requiresAdmin` 再要求 `role_id === 2`。后台的"个人中心"只挂前者，管理功能挂后者。判断顺序也重要——先看登不登录，再看是不是管理员，未登录时给的是"去登录"而非"权限不足"。
- **先 `initializeUserInfo()` 再判断**。刷新后 `userInfo` 是空的（只在内存里），必须先等 `/user/info` 回来 `isLoggedIn` 才有意义。放在守卫里而不是 `main.ts` 里，是为了让判断和取数在时序上绑定，避免"页面已开始渲染但身份还没恢复"。

### 5.3 一个真实隐患：`from.name` 在首次导航时是 `undefined`

上面两处 `.catch(() => { router.push({name: from.name as string}) })` 想表达"用户点了取消，就退回原来所在的页面"。用 `.catch()` 是对的——`ElMessageBox.confirm` 在取消或关闭时会 reject。问题在 **`from`**：Vue Router 在第一次导航（用户直接打开某个 URL）时，`from` 是内部的 `START_LOCATION`，它的 `name` 是 **`undefined`**，于是这一行变成：

```ts
router.push({name: undefined})   // 无法解析出目标路由
```

这不是一次能解析出目标的导航——vue-router 找不到可匹配的路由名，会给出告警，用户也不会被送回任何地方。**触发场景**：直接输入 `/dashboard`（或在后台页面按 F5），此时未登录或权限不足 → 弹出提示 → 点「取消」→ 命中这个分支。**修法**是给一个兜底目标：`const backTo = from.name ? {name: from.name} : {name: 'index'}; router.push(backTo)`——`from.name` 存在就回原页面，不存在（首次导航）就回首页，两条路径都有明确落点。同一段代码还有一个相关问题：**这些分支里始终没有调用 `next()` 或 `next(false)`**。vue-router 的守卫约定是"必须表态"——要么放行要么拦截。这里靠"随后再发起一次 `router.push`"来结束，被取消的那次导航会一直保持挂起。更好的写法是把决定权交回守卫：确认跳登录页时 `next({name: 'index'})`，取消时 `next(false)` 直接放弃这次导航。

## 小结

**分层靠"谁会用它"，不靠技术类型。** `views` 按前台/后台两套界面分，`components` 按用途分，`api` 一个后端模块一个文件，`stores` 按业务域拆。判断落点时问"会被路由直接访问吗"和"只有这一个地方用吗"，比记目录名可靠得多。

**横切关注点必须收口到一个地方。** 请求层是最典型的例子：token 注入、续签响应头回写、业务错误码提示、HTTP 错误分层，全部集中在 `src/utils/request.ts`，业务代码只需写 `if (res.code === 0)`。最大收益不是"代码少"，而是**改一处、全局生效**——调超时时间或改 403 文案都只有一个文件要动。

**约定要写进类型。** `ApiResponse<T>`、`PageResult<T>`、`Hit<T>` 描述了后端响应形状；`RouteMeta` 的声明合并把 `requiresAuth`/`requiresAdmin` 从口头约定变成编译期约束；`env.d.ts` 给需要强类型的环境变量上了锁。这些类型没有运行时代价，但能让"写错了"和"忘了写"更早暴露。

**顺手记录"刻意的例外"。** `/website/title` 故意不遵守统一响应包装，这是为了首屏标题那段内联脚本足够短——好设计，但必须被记住，否则会被后来的人"修正"掉。

最后，这一路也留下了几笔明确的技术债：favicon 里写死的 `127.0.0.1`、`index.html` 里位置不合法的 `<style>`、`userStore.$reset()` 与 `userStore.reset()` 的混用、守卫里 `from.name` 的 `undefined`、以及 `layout` store 正在长成"全局开关抽屉"。每一个都很小，但都属于"不改会持续消耗排查时间"的那一类——比起新增功能，这类修补的性价比往往更高。
