# 组件主题化与 URL 驱动数据流：两个前端的真实排查

这个博客项目的前端有两套并存的视觉语言：**主内容区是浅色**（白底、浅灰边框），**首页侧边栏是深色卡片**（深灰渐变）。这种「一个页面里两种背景」的布局，加上 Vue Router 的组件复用机制，先后给我制造了两个必须认真排查的问题。

第一个是**可读性事故**：一个原本按深色背景写的评论组件被复用到了浅色页面上，文字几乎隐形——肉眼只觉得「有点看不清」，实际对比度已经低到 1.1:1。第二个是**交互失效**：在搜索页里用顶部导航栏搜索，URL 变了、搜索框里的字也变了，结果列表却一动不动。

这篇把两次排查完整写下来：怎么量化「感觉看不清」这种模糊问题、为什么同一个组件不能只有一套配色、为什么「在 watch 里补一句查询」是错的、以及 **URL 作为唯一数据源**带来的结构性收益。读完能拿到三样东西：一套可以直接抄的组件主题化写法、一个判断「状态该放哪」的思路，以及一个关于开发期缓存假象的排错经验。

---

## 一、一个组件，三个调用点，两套背景

### 问题现场

`CommentItem.vue` 负责渲染一条评论及其子评论，它出现在三个地方：

| 调用点 | 所在容器 | 容器背景 |
|---|---|---|
| `web/src/views/web/article/index.vue` | 文章详情页评论区 | 浅色（白底，边框 `#DCDFE6`） |
| `web/src/views/dashboard/user-center/user-comment.vue` | 后台「我的评论」 | 浅色（边框 `#DCDFE6`） |
| `web/src/components/pages/RecentComments.vue` | 首页侧边栏「最新评论」 | 深色渐变卡片 |

而它最初的配色是**纯深色**的：

```scss
/* 改动前：web/src/components/common/CommentItem.vue */
.comment-item {
  color: #dedede;
  .item-card {
    border: 1px solid rgba(255, 255, 255, 0.08);
    background: rgba(255, 255, 255, 0.025);
  }
  .name    { color: #f1f1f1; }
  .time    { color: #8e8e8e; }
  .content { color: #dcdcdc; }
}
```

另外还有一大段给 markdown 预览做的深色适配，把标题设成 `#f1f1f1`、正文设成 `#dcdcdc`、代码块背景设成 `#151515`。这套值放在首页那张深灰卡片里完全没问题——问题是它在文章页和后台页面上也生效了，**近白色的文字压在白色背景上**。

### 把「有点看不清」量化出来

「看不清」是个主观描述，容易让人低估严重程度。按 WCAG 的对比度公式算一遍（sRGB 通道先线性化，再按 `0.2126R + 0.7152G + 0.0722B` 求相对亮度，最后 `(L亮 + 0.05) / (L暗 + 0.05)`）：

| 元素 | 改动前 | 在白底上的对比度 | WCAG 要求 |
|---|---|---|---|
| 评论正文 | `#dcdcdc` | **≈ 1.37 : 1** | 正文 ≥ 4.5 : 1 |
| 用户名 | `#f1f1f1` | **≈ 1.13 : 1** | 正文 ≥ 4.5 : 1 |
| 时间 | `#8e8e8e` | ≈ 3.17 : 1 | 次要文字 ≥ 4.5 : 1 |

对比度的下限是 1:1，也就是**完全同色、什么都看不见**。所以「有点看不清」的真实含义是：正文只剩 1.37:1，用户名只剩 1.13:1——**已经接近隐形了**。量化还有个额外好处：它把「我觉得不好看」变成了「不满足无障碍标准」，需求从主观偏好变成客观缺陷，该不该改、优先改哪个就都清楚了。

### 为什么不能简单改成浅色

最直觉的修法是把配色改成浅色，但这样会**把问题推给第三个调用点**：首页侧边栏的「最新评论」卡片本身是深灰渐变背景（`linear-gradient(145deg, #242424 0%, #1d1d1d 55%, #191919 100%)`），浅色文字放进去照样看不见。而且这不是个别现象，`TagCloud.vue`、`Feedback.vue` 用的都是同一套深色语言。

而主内容区是另一套：`DailyNews.vue` 和 `HotArticles.vue` 的卡片都是 `background: #fcfcfc`，配浅色边框。**同一个首页，左侧一列浅色卡片 + 右侧一列深色卡片**——这个项目的设计语言本来就是双轨的：

```vue
<!-- web/src/views/web/index/index.vue -->
<section class="main-column">
  <HotArticles />
  <DailyNews />
  <ArticleList />
  <!-- …其余卡片略 -->
</section>
<aside class="sidebar">
  <TagCloud />
  <RecentComments />
  <!-- …其余卡片略 -->
</aside>
```

> 关于这个双轨制我还有个具体教训：热门文章卡片最早是按**侧边栏深色风格**写的，后来把它移到主内容区首位时，深色卡片夹在一列浅色卡片中间非常突兀。最后照着 `DailyNews.vue` 的头部参数（25px 标题、`#464646` 字色、34×3px 标题线、右上角绝对定位的英文副标题）整个重写了一遍，两张卡片才看起来是一套。**跨区域复用组件时，配色往往比结构更需要重做。**

结论：不能全局改成浅色，也不能全局保持深色——**必须让组件知道自己处在哪种背景里**。

### 方案对比

| 方案 | 问题 |
|---|---|
| `v-if` 渲染两套模板 | DOM 结构和逻辑重复两遍，改一处要记得改另一处 |
| 传两套 class、样式写两遍 | 每条规则都要重复，主题越多越糟 |
| 拆成两个组件 | 逻辑重复，且调用方需要知道「该用哪个」 |
| **`theme` prop + CSS 自定义属性** | 结构只写一遍，只有调色板需要写两份 |

第四种是 Vue 生态里的常规做法，也是我最终采用的。

### 实现：两套调色板，样式只引用变量

组件根节点根据 prop 挂上互斥的 class，也就是 `:class="'theme-' + theme"`；prop 本身用字面量联合类型约束，并给出浅色默认值：

```ts
/* web/src/components/common/CommentItem.vue */
withDefaults(
  defineProps<{
    comments: Comment[];
    theme?: "light" | "dark";
  }>(),
  { theme: "light" }
);
```

**默认值选浅色是个关键决策**：三个调用点里有两个是浅色页面，把浅色设为默认，这两个地方一行都不用改就修好了；只有深色卡片那一处需要显式声明。

然后在根节点的两套 class 里定义调色板（只摘录有代表性的部分）：

```scss
/* web/src/components/common/CommentItem.vue */
.comment-item {
  &.theme-light {              /* 浅色（默认）：文章页 / 后台「我的评论」 */
    --ci-text: #333333;
    --ci-text-strong: #111111;
    --ci-card-bg: #fbfbfb;
    --ci-card-border: #ebebeb;
    --ci-input-focus-ring: rgba(0, 0, 0, 0.045);
    --ci-pre-bg: #f7f7f7;
    /* …两个主题各有 31 个变量 */
  }

  &.theme-dark {               /* 深色：首页侧边栏「最新评论」深色卡片 */
    --ci-text: #dcdcdc;
    --ci-text-strong: #f1f1f1;
    --ci-card-bg: rgba(255, 255, 255, 0.025);
    --ci-card-border: rgba(255, 255, 255, 0.08);
    --ci-input-focus-ring: rgba(255, 255, 255, 0.035);
    --ci-pre-bg: #151515;
  }

  color: var(--ci-text);
}
```

**深色那一套的数值是逐个照抄原有实现的**，所以首页侧边栏的观感完全没有变化——重构不该顺手改设计。下面的样式规则只引用变量，一个写死的颜色都不留：

```scss
/* web/src/components/common/CommentItem.vue */
.item-card {
  border: 1px solid var(--ci-card-border);
  background: var(--ci-card-bg);
  box-shadow: var(--ci-card-shadow);
  &:hover {
    background: var(--ci-card-bg-hover);
    border-color: var(--ci-card-border-hover);
  }
}
```

调用方也变得非常直白——**声明自己在什么背景上**：文章页 `<comment-item :comments="comments"/>` 和后台 `<comment-item :comments="[item]"/>` 都不传，走浅色默认值；首页侧边栏那处写 `<comment-item :comments="comments" theme="dark"/>`。

改完之后的关键对比度：

| 元素 | 改动后 | 对比度 |
|---|---|---|
| 评论正文 | `#333333` | **≈ 12.2 : 1** |
| 用户名 | `#111111` | **≈ 18.3 : 1** |
| 时间 | `#767676` | ≈ 4.4 : 1 |

时间那一项要特意说明：`#767676` 在白底上算出来是 **4.39:1，刚好卡在 AA 门槛（4.5:1）下面一点点**。它是 11px 的次要信息，观感没问题，但要严格过 AA 还得再深一档。把数字算清楚的价值就在这里——你会知道哪些是真达标了，哪些只是「看起来还行」。

### 递归组件必须把主题透传下去

评论是树形结构，子评论靠组件递归渲染，递归调用处是 `:comments="item.children"` 加上 **`:theme="theme"`**。

那行 `:theme` **不能省**。只传 `comments` 的话，子组件会拿到 prop 的默认值 `light`——结果是深色卡片里的一级评论是深色配色，**展开的子评论却是浅色**，一层一层颜色不一样。这类问题的特点是「主路径看不出问题，一有嵌套数据就露馅」，所以递归调用处要专门检查一遍。

### 一个容易踩的坑：变量定义在根元素上，父级覆盖不了

我一开始想的是另一种写法：组件里只写一套中性默认值，由父容器注入变量来切主题，比如在侧边栏卡片上写 `--ci-text: #dcdcdc`。

**这个思路行不通。** CSS 自定义属性的取值规则是：**元素自身的声明优先于从祖先继承来的值**。如果 `--ci-text` 定义在 `.comment-item` 自己身上，那么父级 `.recent-comments` 就算定义同名变量，也只是被继承——继承的优先级低于自身声明，最终还是用自己的那套。

所以「父级注入」只适用于组件自身不声明、只消费的写法，而那样每个变量都得带兜底值（`var(--ci-card-bg, #fbfbfb)`），规则会变得很啰嗦。我最后选了「**两套值都定义在根节点、用互斥 class 切换**」：谁生效完全由 `theme` 这一个 prop 决定，不受外部环境影响，推理成本最低。代价是调色板写在组件内部，但它本来就该由组件自己拥有。

### 变量完备性：两个主题必须严格对齐

这是主题化最容易出事、也最难排查的地方。所有样式引用的是同一个变量名（比如 `var(--ci-pre-bg)`）。如果某个变量**只在 `theme-light` 里定义了、`theme-dark` 里漏了**，深色主题下这个变量就没有值——`var()` 会回退成「继承值或初始值」，表现是**某个元素的颜色莫名其妙变成环境色**（通常是黑色或透明），而且只在一种主题下出现，极难定位。

所以两个主题的变量集合必须完全相同。我没有靠眼睛核对，而是从编译后的 CSS 里分别抽出两个 class 的变量名集合求差集，结果是：**light 定义 31 个、dark 定义 31 个、样式引用 31 个，差集为空，无未定义引用**。这三个数字相等就从静态上排除了「某个主题缺变量」这类问题，以后加变量时也会自然成对添加。

### 顺带修正的两处深色假设

主题化过程中暴露出两个「只在深色背景下成立」的写死值：

```scss
/* web/src/components/common/CommentItem.vue */

// 其一：代码块背景。原来写死 #151515，而文章页引的是 md-editor-v3 的浅色预览样式，
// 两者叠在一起是「浅色语法配色 + 深黑色底」，在浅色页面上就是一块突兀的黑斑
pre {
  background: var(--ci-pre-bg);      /* 浅色 #f7f7f7 / 深色 #151515 */
  border: 1px solid var(--ci-pre-border);
}

// 其二：输入框聚焦光晕。原来用 rgba(255,255,255,0.035) 做外发光，
// 白光晕只有落在深色背景上才看得见，在白底上等于没写
&:focus {
  box-shadow:
    0 0 0 2px var(--ci-input-focus-ring);
  /* 浅色 rgba(0,0,0,0.045) / 深色 rgba(255,255,255,0.035) */
}
```

这类「白色半透明用于高光、黑色半透明用于阴影」的写法在深浅两套之间是**必须互换**的。凡是在主题化时看到 `rgba(255, 255, 255, x)` 用作背景或阴影，都值得停下来想想它在浅色下是否成立。

---

## 二、搜索页「URL 变了，结果不刷新」

### 现象

首页顶部导航栏有个搜索框，回车会跳转到 `/search?query=关键词`：

```ts
/* web/src/components/layout/WebNavbar.vue */
const handleSearch = () => {
  const keyword = searchKeyword.value.trim();
  if (!keyword) {
    return;
  }
  router.push({
    path: "/search",
    query: {
      query: keyword
    }
  });
};
```

从首页搜没问题。但**已经在搜索页时**再搜一次，就出现了一个很奇怪的状态：地址栏确实从 `/search?query=A` 变成了 `/search?query=B`，页面上搜索框里的词也变成了 B，**但下面的结果列表还是 A 的结果**。三个地方各自为政，只有列表不动。

### 根因一：`watch` 只同步状态，从不触发查询

先看搜索页原来的写法：

```ts
/* 改动前：web/src/views/web/search/index.vue */
onMounted(() => {
  articleSearchRequest.query = route.query.query as string || ""
  articleSearchRequest.category = route.query.category as string || ""
  /* …tag / sort / order 同理 */
  page.value = Number(route.query.page) || 1
  page_size.value = Number(route.query.page_size) || 10
})

watch(() => route.query, (newQuery) => {
  articleSearchRequest.query = newQuery.query as string || ""
  articleSearchRequest.category = newQuery.category as string || ""
  /* …其余字段同样只是赋值，全程没有调用查询 */
  articleSearchRequest.page = Number(newQuery.page) || 1
  articleSearchRequest.page_size = Number(newQuery.page_size) || 10
}, {immediate: true})

nextTick(() => { getArticleSearchTableData() })
```

问题就在这：**`watch` 的函数体只是把新参数赋值给本地状态，从头到尾没有调用过查询。** 真正的查询只在初始化时通过那个 `nextTick` 跑了一次。

那为什么从首页跳进来是好的？因为那时组件是**首次挂载**，`nextTick` 会执行。而从 `/search?query=A` 跳到 `/search?query=B` 时：两者属于**同一个路由记录**（路由表里都是 `name: "search"` 那一条），Vue Router 会**复用组件实例**、不销毁重建，所以 `onMounted` 不再执行、`nextTick` 也不会再跑；`watch` 倒是触发了（`route.query` 是个新对象），可它只改状态、不查询。于是状态更新了、UI 的一部分跟着变了（搜索框绑的就是 `articleSearchRequest.query`），**唯独列表没有重新拉取**。

### 根因二：分页状态有两份，会互相打架

如果只是「没触发查询」，那在 `watch` 里补一句就好了。但这里还有个更隐蔽的问题。

分页组件绑的是 **ref**（`:current-page="page"`、`:page-size="page_size"`），而真正发给后端的却是 **reactive 对象**，并且在查询函数里用 ref 去覆盖它——函数第一行就是 `articleSearchRequest.page = page.value;` 和 `articleSearchRequest.page_size = page_size.value;`。

而前面那个 `watch` 改的是**另一份**（`articleSearchRequest.page`），没碰 `page` 这个 ref。后果就很具体了——**在第 3 页时用导航栏搜索**：

1. 导航栏 push 的是 `/search?query=B`，URL 里**没有 page 参数**
2. `watch` 把 `articleSearchRequest.page` 重置为 1
3. 但 `page` 这个 ref 仍然是 3
4. 查询函数第一步就把刚同步好的 1 **覆盖回 3**
5. 于是带着 `page=3` 去查一个新关键词 → 结果大概率是**空列表**

同一个状态存在两份、且两处各自更新一半，这类 bug 的特点就是「正常路径看不出问题，一旦有多来源输入就错得莫名其妙」。

### 为什么不能只在 `watch` 里补一行查询

最省事的修法是在 `watch` 里加一句 `getArticleSearchTableData()`。但查询函数结尾会把本地状态**写回 URL**（`router.push` 到当前 path，query 里放上全部查询条件）。于是补上那句之后就形成闭环：

**用户操作 → 查询 → 改 URL → watch 触发 → 再查询 → 再改 URL → …**

实际表现是点一次翻页或筛选**会发两次后端请求**（第一次是操作直接调的，第二次是 watch 被自己写的 URL 触发的），而且这个循环能不能停，取决于 Vue Router 对「重复导航」的中止行为——太脆弱了。**根子在数据结构上：状态同时被 URL 和本地变量两边写，还各写一半。** 只补一行是打补丁，不是修结构。

### 重构：让 URL 成为唯一数据源

明确一条规则：**用户操作只写 URL，URL 变化统一触发查询。** 本地状态永远只是 URL 的镜像，不反过来驱动数据。按这个规则把逻辑拆成三个职责单一的函数：

```ts
/* web/src/views/web/search/index.vue */

// 只读 URL：同时同步 reactive 字段与两个 ref，消除双份状态
const syncFromRoute = () => {
  articleSearchRequest.query = route.query.query as string || ""
  articleSearchRequest.category = route.query.category as string || ""
  page.value = Number(route.query.page) || 1
  page_size.value = Number(route.query.page_size) || 10
  articleSearchRequest.page = page.value          // ← 两份状态写同一个值
  articleSearchRequest.page_size = page_size.value
}

// 只查询：不碰 URL
const getArticleSearchTableData = async () => {
  const table = await articleSearch(articleSearchRequest)

  if (table.code === 0) {
    articleTableData.value = table.data.list;
    total.value = table.data.total;
  }
}

// 用户操作：只把当前条件写进 URL
const applySearch = () => {
  router.push({
    path: route.path,
    query: {
      query: articleSearchRequest.query,
      category: articleSearchRequest.category,
      /* …tag / sort / order 同理 */
      page: String(page.value),
      page_size: String(page_size.value),
    }
  })
}

// URL 是唯一入口：无论变化来自导航栏、筛选、排序还是翻页，都在这里触发
watch(() => route.query, () => {
  syncFromRoute()
  getArticleSearchTableData()
}, {immediate: true})
```

注意 `syncFromRoute` 里那两行 `articleSearchRequest.page = page.value`——**它同时写了 reactive 字段和 ref**，这就把「双份状态」从根上消掉了：后面不管哪个函数读哪一份，值都一样。

调用点也全部改成「只写 URL」：筛选/排序相关的 `changeArticleSearchItem` 变成 `applySearch()`；`handleSizeChange` 先写 `page_size.value` 再调 `applySearch()`；`handleCurrentChange` 先写 `page.value` 再调 `applySearch()`。原本那两个 `onMounted` / `nextTick` 代码块可以直接删掉——`watch` 带了 `immediate: true`，首次渲染时它就会完成「同步 + 查询」，初始加载和后续变更走的是**完全同一条路径**。

### 收益

| | 改动前 | 改动后 |
|---|---|---|
| 导航栏搜索 | 只在首次进入页面生效 | 任何时刻都生效 |
| 在第 N 页搜索 | 带着 page=N 查，返回空 | `syncFromRoute` 把 page 重置为 1 |
| 状态份数 | 分页状态两份，易不一致 | 一份，由 URL 派生 |
| 请求次数 | 补一行 fetch 的写法会发两次 | 一次操作一次请求 |
| 循环风险 | 操作 → 查询 → 改 URL → 再查询 | 结构上不可能发生 |
| 初始加载 | `onMounted` + `nextTick` 单独一条路径 | 与变更共用同一条路径 |

最后一行是我最看重的：**初始加载和后续变更不再有两条代码路径。** 这个 bug 之所以出现，本质就是「初始化时同步并查询」和「参数变化时只同步」是两段独立逻辑；只要有两段逻辑处理同一件事，早晚会漂移。

### 一个必须讲清楚的行为变化

重构后有个可观测的变化：**条件完全没变时，点搜索不会重新请求。** 原因是 `router.push` 到完全相同的 query 时，Vue Router 会判定为重复导航并中止它，`route.query` 不产生变化，watch 自然不触发。

这算不算 bug？我认为是**可接受的取舍**：正常使用不受影响（导航栏搜索会丢掉其它参数，URL 必然变化，一定重新查），只有「在搜索页里点搜索按钮、且关键词和筛选条件一字未改」这一种情况命中，而这时界面上显示的结果**本来就是对的**。如果确实希望「点搜索永远重新查一次」，可以在 `applySearch` 里先比较新条件与当前 URL，相同就直接调查询、不同才 push——但那会重新引入「一个操作可能有两条执行路径」的分支，我倾向于保持现在的单向结构。

---

## 三、尾声：别把开发期的「陈旧 UI」当成代码 bug

排查上面两个问题的过程中，我还撞见一个很容易误判的现象。有一轮我把「加进主栏」和「从侧栏删掉」两条改动放在了**同一次操作里**，两次写文件相隔只有几毫秒，结果浏览器里出现了「新旧两份同时存在」的画面——主内容区有那个卡片，侧边栏也有。

代码查了半天没找出问题，最后比对两份东西才定位：**磁盘上的源文件是对的**（只有一处使用），而 **dev server 实际返回给浏览器的编译产物是错的**（渲染函数里两处都在）。也就是说，文件监听把两次写入**合并**成了一次事件，读到的却是**中间状态**；之后文件不再变化，这个中间态就被缓存下来持续提供给浏览器。

判断方法很直接：**去看 dev server 返回的编译结果，而不是盯着源码。** 以 Vite 为例，直接请求那个模块就能拿到转换后的代码，编译错误会返回 500 并带上信息：

```bash
# 让 dev server 重新编译指定模块（本项目 dev server 在 80 端口，见 web/vite.config.ts）
curl "http://127.0.0.1:80/src/components/pages/HotArticles.vue"
```

处置顺序也就两步：**先硬刷新**（`Ctrl + F5`，多数情况 HMR 已更新模块，但浏览器里还留着旧的组件树），**还不行就重启 dev server**（文件监听确实漏了事件时，只有重启才能重建缓存）。

两条经验：一是**修改同一个文件的多处内容时不要挤在一次操作里**，给文件监听留出处理每一次写入的机会；二是**遇到「改了但没生效」，先怀疑工具链的缓存，再去怀疑自己的代码**——尤其当源码看起来完全正确的时候。这不代表代码一定没问题，但排查顺序反了会很浪费时间。

---

## 小结

这两个问题表面上一个是 CSS、一个是状态管理，但底层的教训是同一个：**同一份事实不要存在两个地方。**

第一个问题里，「这个组件现在处于什么背景」这个事实，原本隐含在「组件内部写死了一套深色」里；一旦组件被放到别处，这个隐含假设就失效了。显式声明成 `theme` prop、把颜色收敛成一组变量之后，事实只剩一处，三种背景下的表现都可推导。

第二个问题里，「当前要查什么」这个事实，原本同时存在于 URL、reactive 对象、两个 ref 里，四处各更新一部分，于是必然漂移。把 URL 定为唯一来源、其余全部由它派生之后，多来源输入（导航栏、筛选、排序、翻页、直接刷新页面）就都收敛到了同一条路径上——bug 不是被修掉的，而是变得**不可能发生**了。

要从整篇里带走一条最实用的做法，我推荐这个：**遇到「感觉不对但说不清」的问题时，先想办法把它变成一个数字。** 对比度是 1.37:1 还是 12.2:1、请求是一次还是两次、状态有一份还是两份——能算出数字的问题，通常已经解决一半了。
