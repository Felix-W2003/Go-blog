# Markdown 编辑器接入与 .md 文件导入：一个博客后台编辑器的实现与取舍

做博客后台最费心思的地方不是接口，而是写作体验。文章写不顺，功能再多也没人用。

这篇讲我在这个 Go + Vue 3 的博客里把 Markdown 编辑器接进后台的完整过程：编辑器怎么选、图片上传为什么要自己接管、发布弹窗和页面之间怎么分工传数据，以及最后加的「导入 .md 文件」功能——包括 front matter 解析器的设计取舍，和四个在真实文件上踩到的坑（BOM、CRLF、代码块、标题重复）。读完能拿到一套可以直接照搬的接入写法、一个不依赖任何 yaml 库的 front matter 解析器，以及一个判断标准：**什么时候值得手写解析，什么时候该老老实实引依赖。**

## 一、编辑器选型：为什么是 md-editor-v3

候选不少，我最终选 `md-editor-v3`，主要看三点。它是 Vue 3 原生的，不是 React 组件的包装，props 和事件都是 Vue 习惯，不会出现「数据变了视图不动」这种跨框架边界问题。它的编辑器、预览、目录是三个独立组件：写作场景要「能改」，阅读场景只要「能看」，让全功能编辑器去承担阅读渲染是浪费；拆开之后就能按场景只引需要的那个——后台发布页用 `MdEditor`，前台详情页用 `MdPreview` 渲染正文，详情页侧边栏用 `MdCatalog` 生成目录导航。`MdPreview` 也不只服务文章正文——评论区同样用它渲染每条评论的 Markdown（`web/src/components/common/CommentItem.vue`），那个组件里因此有一大段 `:deep(.md-editor-preview)` 覆盖，专门调配正文、标题、链接、引用、代码块的配色，让它能同时适配深色卡片和浅色页面。第三点是它的上传钩子暴露得足够底层：不强制你用内置逻辑，而是给你一个回调，把 `File[]` 交给你、你返回 URL 数组，这一点直接决定了我能不能带上自己的鉴权头。

## 二、把编辑器装进发布页

发布页在 `web/src/views/dashboard/articles/article-publish.vue`，核心就一行，`text` 是正文，一个普通的 `ref`：

```vue
<!-- web/src/views/dashboard/articles/article-publish.vue -->
<MdEditor v-model="text" @onUploadImg="onUploadImg" @onSave="onSave" @onChange="onChange"/>
```

### 2.1 样式为什么要单独引

`md-editor-v3` 不把样式打进组件，得手动引。后台只要编辑器样式 `import 'md-editor-v3/lib/style.css'`；前台详情页（`web/src/views/web/article/index.vue`）还要再引一个预览样式：

```ts
import 'md-editor-v3/lib/style.css';
import 'md-editor-v3/lib/preview.css';
```

一开始我觉得这是设计缺陷，后来想明白了：编辑器和预览的样式表体积都不小，而两者的使用场景往往不重叠。如果库强制打包，前台就要白白背上一整套工具栏、下拉菜单、弹窗的 CSS。让使用方按需引入是把选择权交出来了，代价是容易漏引——漏了 `preview.css` 的表现是预览区排版全乱，但代码不报错，只能靠肉眼发现。

### 2.2 图片上传：为什么不能吃默认行为

这是接入时第一个真正要思考的点。编辑器的默认上传带不上我的鉴权头，而后端上传接口要求 `x-access-token`，所以 `onUploadImg` 必须自己实现：

```ts
// web/src/views/dashboard/articles/article-publish.vue
const onUploadImg = async (files: File[], callback: (urls: string[]) => void) => {
  const res = await Promise.all(
      files.map((file) => {
        const form = new FormData();
        form.append('image', file);
        return axios.post('/api/image/upload', form, {
          headers: {'Content-Type': 'multipart/form-data'},
          withCredentials: true,
        });
      })
  );
  callback(res.map((item) => item.data.data.url));
};
```

三个细节。**并发上传用 `Promise.all` 而不是串行**：一次拖进去五张图，串行要等五个往返，并发只需要最慢的那一张，代价是后端要能承受瞬时并发，个人博客场景完全没问题。**`callback` 是编辑器的回填出口**：拿到 URL 数组后必须调用它，编辑器才会把 `![](url)` 插到光标处，忘了调的表现是「图片传上去了但正文里什么都没出现」，很容易误判成上传失败。**表单字段名写死成 `image`**，这必须和后端 `c.FormFile("image")` 对齐，属于跨端契约，改名要两边一起改。

### 2.3 自动保存

编辑器给了 `onChange` 和 `onSave` 两个钩子。我没有后端草稿接口，就用 `localStorage` 顶上：

```ts
// web/src/views/dashboard/articles/article-publish.vue
const onSave = (v: string, _: Promise<string>) => {
  localStorage.setItem('article', v)
};
const onChange = (v: string) => {
  if (isAutoSaveEnabled.value) { onSave(v, Promise.resolve('')) }
}
```

恢复在初始化时做：`const text = savedArticle ? ref(savedArticle) : ref('')`。开关状态本身也持久化（键名 `isAutoSaveEnabled`，默认开），否则每次刷新都要重新打开一次，等于没有开关。

这里有两个坑。**`localStorage` 只认字符串**：写的时候必须 `String(bool)`，恢复时必须 `=== 'true'` 显式比较，直接拿 `getItem` 的返回值当布尔用会得到 `"false"` 这个**真值**——我踩过，表现是关掉自动保存后刷新又自己开了。**`onSave` 的第二个参数我没用上**：这个钩子原本是给「保存到远端」设计的，期望你返回 Promise 表示保存结果，我用本地存储就传个 `Promise.resolve('')` 占位，这是接口设计与实际用法不匹配时的常见妥协。

## 三、发布流程：页面持正文，弹窗持元数据

这是整个功能里最该讲清楚的一处设计，因为它决定了后面「导入」怎么落地。页面 `article-publish.vue` 手里有**标题**和**正文**（正文属于编辑器）；点「发布文章」弹出对话框，里面是 `ArticleCreateForm`，负责**封面、类别、标签、简介**这些元数据：

```vue
<!-- web/src/views/dashboard/articles/article-publish.vue -->
<el-dialog v-model="articleCreateVisible" width="500" align-center
           destroy-on-close :before-close="articleCreateVisibleSynchronization">
  <article-create-form
      :title="title" :content="text"
      :category="importedMeta.category" :tags="importedMeta.tags"
      :abstract="importedMeta.abstract" :cover="importedMeta.cover"
  />
</el-dialog>
```

表单那边用 `props` 初始化：

```ts
// web/src/components/forms/ArticleCreateForm.vue
const articleCreateFormData = reactive<ArticleCreateRequest>({
  cover: props.cover ?? '',
  title: props.title,
  category: props.category ?? '',
  // 必须复制一份：否则在弹窗里增删标签会连带改掉父组件传进来的数组
  tags: props.tags ? [...props.tags] : [],
  abstract: props.abstract ?? '',
  content: props.content,
})
```

### 3.1 `destroy-on-close` 不是可选项

注意 `el-dialog` 上的 `destroy-on-close`，它必须开着，原因在 `reactive` 的语义上：`reactive({...})` 只在组件实例创建时执行一次，`props` 后续再变也不会重新读。所以弹窗不销毁的话，第二次打开复用的是上一次的实例，表单里还是上一次的数据。我在模板里留了句注释提醒自己：`<!--  这里必须销毁，不然不会重新加载props-->`。同样的道理，`ArticleUpdateForm` 也从 `props.article._source` 一次性初始化所有字段，靠的是列表页每次打开都新建弹窗。

### 3.2 这个分工的取舍

**好处**是正文数据只有一份，就在编辑器里，弹窗只负责补充元数据。**代价**是数据流变成两段：正文靠 `props` 单向流入弹窗，但弹窗里如果要编辑正文，改动就没法回流到页面——所以两个表单对正文的处理完全不同，发布表单只读接收，更新表单则把编辑器搬进 `el-drawer` 里直接改。如果重来一次，我可能会把整个「文章草稿」抽成一个 Pinia store，页面和弹窗都读写同一个 store，就不存在「谁是唯一数据源」的问题。现在的写法胜在直观，代价是每加一个字段就要在 props 上多接一条线——「导入文章」往弹窗里送四个新字段时，这个代价就很明显了。

## 四、导入 .md：交互层

需求是：发布时多一个按钮，选一个 `.md` 文件，自动把内容填到对应位置。

### 4.1 隐藏的原生 input

没有用 `el-upload`，用的是最朴素的原生控件：

```vue
<!-- web/src/views/dashboard/articles/article-publish.vue -->
<el-button icon="Upload" @click="triggerImport">导入文章</el-button>
<input ref="fileInputRef" class="import-file-input" type="file"
       accept=".md,.markdown,text/markdown" @change="handleImportFile" />
```

```ts
const triggerImport = () => { fileInputRef.value?.click() }
```

样式里把它藏掉：`.import-file-input { display: none; }`。

**为什么不直接用 `el-upload`？** 因为我不需要它提供的上传能力——文件根本不往服务器传，是在浏览器里读出来解析的，而 `el-upload` 的整套机制（action、headers、上传状态、文件列表）在这里全是负担。原生 input 只负责「弹出选文件对话框」这一件事，剩下的自己控制，这也是 `accept` 能写得很精确的原因。至于为什么要隐藏原生 input：它的样式几乎不可控，不同浏览器渲染出来的按钮长得完全不一样，内部文字也改不了，用一个自己的按钮去 `click()` 它是标准做法。

### 4.2 `input.value = ''` 为什么必须在读文件之前

这一行不写会有个很迷惑的现象：**第二次选同一个文件，什么都不会发生。**

```ts
const handleImportFile = async (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''            // 立刻清空，保证连续导入同一个文件也能再次触发 change
  if (!file) { return }
  // ...
}
```

原因是 `change` 事件的触发条件包含「值发生了变化」。选完文件后 input 的 `value` 就停在那份文件的路径上，再选同一个文件，值没变，事件就不触发。注意我是**先取出 `file` 再清空 `value`**——文件对象已经拿到了，清空 input 不影响它，反过来写就丢了。另外这一行放在 `if (!file) return` **之前**，是因为用户点开对话框又点取消时 `files` 是空的，但 value 也可能被浏览器置空，先清一次没有副作用，逻辑上更稳。

### 4.3 读文本：`File.text()`

`const parsed = parseMarkdownArticle(await file.text(), file.name)`。不需要 `FileReader`——`File` 继承自 `Blob`，而 `Blob.text()` 返回 Promise，直接 `await` 就行。这是这些年 Web API 最实在的一次简化，原来那套 `new FileReader()` + `onload` + `result` 的写法现在一行顶掉。

## 五、front matter：为什么不引 yaml 依赖

Markdown 博客的元数据惯例是 front matter，用 `---` 夹一段声明：

```markdown
---
title: 从零实现 JWT 双 Token 认证
category: 技术
tags: [Go, JWT, Redis]
abstract: 用 Redis 黑名单实现可主动下线的双 Token 方案
cover: /uploads/image/xxx.png
---

正文从这里开始……
```

### 5.1 手写还是引依赖

标准答案是引 `yaml` 或 `js-yaml`。我没引，理由很实际：**这个解析器只需要覆盖「我自己写出来的格式」**。我不是在做通用 Markdown 工具，只有一个人写文章，front matter 的写法可控；为了覆盖 YAML 规范里的锚点、多行标量、嵌套映射、类型推断去引一个几十 KB 的库，收益为零。

但这不是无条件成立的，判断标准是：**如果输入来自不受控的第三方，就必须引库。** 这里是「我自己的文章」，手写没问题；哪天要支持从任意平台导入，手写解析器的边界会立刻变成 bug 来源。我在文件头把支持的子集写清楚了，就是为了让这个边界显式：

```ts
// web/src/utils/markdown.ts
/**
 * 只支持以下写法：
 *        key: value          key: "带引号的值"
 *        key: [a, b, c]      key: a, b, c
 *        key:
 *          - a
 *          - b
 */
```

### 5.2 块状列表的处理

行内的 `[a, b, c]` 好办，`split(',')` 就完了。麻烦的是跨行写法，这里得做**前瞻**：遇到 `key:` 后面是空值时，往下看连续多少行是 `- xxx` 的形式，把它们收成数组，然后把外层循环的游标跳过去。

```ts
// web/src/utils/markdown.ts
if (value === "") {
  const items: string[] = [];
  let j = i + 1;
  while (j < lines.length && /^\s*-\s+/.test(lines[j])) {
    items.push(stripQuotes(lines[j].replace(/^\s*-\s+/, "").trim()));
    j++;
  }
  if (items.length) {
    result[key] = items;
    i = j - 1;   // 跳过已经消费掉的行
  }
  continue;
}
```

`i = j - 1` 这行容易写错。外层是 `for (let i = 0; ...)`，每轮结束会 `i++`，所以这里要减一才能正好停在被消费的最后一行；写成 `i = j` 会漏掉一行，而且只在「块状列表后面紧跟另一个字段」时才暴露——这种 bug 特别难查。

### 5.3 字段别名：为了少填几次表

```ts
// web/src/utils/markdown.ts
const FIELD_ALIASES: Record<string, keyof ParsedArticle> = {
  title: "title",  "标题": "title",
  category: "category", categories: "category", "类别": "category", "分类": "category",
  tags: "tags", tag: "tags", "标签": "tags",
  abstract: "abstract", description: "abstract", summary: "abstract",
  excerpt: "abstract", "简介": "abstract", "摘要": "abstract",
  cover: "cover", image: "cover", thumbnail: "cover", "封面": "cover",
};
```

做这个是因为**不同导出源的字段名根本不统一**：Hexo 用 `description`，Hugo 偏爱 `summary`，有的工具写 `excerpt`；中文用户手写时又倾向于直接写「标签」「分类」。别名表让这些写法都能落到位，多花二十行，换掉的是「导入完还得手动改一遍字段名」——而这恰恰是导入功能想省掉的事。查表时做了大小写兜底 `FIELD_ALIASES[rawKey] ?? FIELD_ALIASES[rawKey.toLowerCase()]`，`Title` 和 `title` 都认；中文字符 `toLowerCase()` 原样返回，所以这条兜底对中文键没有副作用。

## 六、没有 front matter 时的启发式

纯文本文件丢进来也应该能填出个大概。标题的优先级是三级，文件名兜底很粗糙但很实用——把文件命名成「Redis 缓存穿透笔记.md」再导入，标题直接就有了：

```ts
// web/src/utils/markdown.ts
// 标题优先级：前置声明 > 正文首个一级标题 > 文件名
let title = (frontMatter.title ?? "").trim();
if (!title) {
  const heading = trimmedBody.match(/^#[ \t]+(.+)$/m);
  if (heading) {
    title = heading[1].trim();
    content = trimmedBody.replace(/^#[ \t]+.+\n?/, "").trim();
  } else {
    title = fileName.replace(/\.(md|markdown)$/i, "").trim();
  }
}
```

简介的提取核心是「找第一段**正常**文字」，判断很挑剔：

```ts
// web/src/utils/markdown.ts
if (!text) continue;
if (/^#{1,6}\s/.test(text)) continue;         // 标题
if (/^!\[/.test(text)) continue;              // 图片
if (/^>/.test(text)) continue;                // 引用
if (/^([-*+]|\d+\.)\s/.test(text)) continue;  // 列表
if (/^\|/.test(text)) continue;               // 表格
if (/^</.test(text)) continue;                // HTML
```

为什么？因为「第一段」经常不是正文。文章开头放封面图的写法太常见了，直接取第一段的话简介就变成一串图片语法；列表和表格同理，它们作为开篇时抽出来当简介读起来很怪。抽到之后还要洗掉行内标记：

```ts
// web/src/utils/markdown.ts
return text
  .replace(/!\[[^\]]*\]\([^)]*\)/g, "")     // 图片整体去掉
  .replace(/\[([^\]]*)\]\([^)]*\)/g, "$1")  // 链接只保留文字
  .replace(/[*_~`]+/g, "")                  // 强调、行内代码的标记符
  .replace(/\s+/g, " ").trim();
```

链接那条用 `$1` 保留文字是刻意的：`[Redis 官网](https://example.com)` 洗成 `Redis 官网`，比整段删掉信息量更大。最后截断到 150 字（`ABSTRACT_MAX_LENGTH = 150`）。这个数字是拍脑袋定的，但有个约束：文章列表卡片上简介是两行省略，中文一行约 25 到 30 字，两行就是 60 字左右；150 字保证能填满，多余的交给 CSS 的 `-webkit-line-clamp` 裁掉。**宁可多给一点让 CSS 去裁，也不要提前截得太短**——简介在详情页顶部也会展示，那里空间更大。

## 七、四个坑

前面都是设计，这一段是实打实踩出来的。

### 7.1 BOM：文件开头那个看不见的字符

Windows 上的编辑器（尤其记事本和一部分 Markdown 工具）保存文件时，会在开头写一个 `\uFEFF`——字节序标记。它在编辑器里完全看不见，但**会让「文件必须以 `---` 开头」的判断失效**。症状很隐蔽：front matter 明明写对了，导入之后一个字段都没填上，标题还退化成了文件名。修复在解析入口，一行：

```ts
// web/src/utils/markdown.ts
const normalized = raw.replace(/^\uFEFF/, "").replace(/\r\n/g, "\n");
```

`^\uFEFF` 必须在最前面处理，因为后面所有正则的 `^` 锚点都假设文件是真的从那三个短横线开始的。

### 7.2 CRLF：不统一换行符，所有按行匹配的正则都会失配

这是同一个坑的另一半，而且更隐蔽。front matter 的判定正则写的是 `\n`：

```ts
const matched = raw.match(/^---\n([\s\S]*?)\n---[ \t]*(?:\n|$)/);
```

但 Windows 文件的行尾是 `\r\n`，所以开头的 `---\r\n` 匹配不上 `^---\n`，**整个 front matter 直接被跳过**。为什么不把正则改成 `\r?\n` 兼容两种？可以，但要改的地方太多了：解析 front matter 要改、按行切分要改、提取简介的正则也全部要改。**与其在每个正则上打补丁，不如在入口处做一次归一化**——就是 7.1 里那行 `replace` 的后半段。之后所有代码都按 `\n` 假设写，一处都不用操心。把平台差异挡在边界上，是成本最低的做法。

### 7.3 代码块内部被当成简介

这个是写测试用例时才发现的。我原本的 `extractAbstract` 是逐行判断，跳过 ` ``` ` 开头的那一行：`if (/^```/.test(text)) continue;`。看起来没问题，但**跳过围栏本身，不等于跳过围栏里面的内容**。如果文章是这样开头的：

````markdown
# 标题

```go
func main() { fmt.Println("hello") }
```

真正的第一段文字在这里。
````

那么 `func main()` 那一行既不是空行、也不是标题、不是图片、不是列表、不是表格——它会被当成「第一段正常文字」，于是**简介变成了一段 Go 代码**。正确做法是用状态变量跟踪围栏的开关：

```ts
// web/src/utils/markdown.ts
let inCodeBlock = false;
for (const line of body.split("\n")) {
  const text = line.trim();
  // 代码围栏：切换状态；围栏本身与它内部的所有行都要跳过
  if (/^(```|~~~)/.test(text)) {
    inCodeBlock = !inCodeBlock;
    continue;
  }
  if (inCodeBlock) continue;
  // ...后面的启发式判断
}
```

顺带把 `~~~` 也认了，Markdown 里它是 ` ``` ` 的等价写法；同时我把原来那行单独的 `/^```/` 判断删掉了，因为围栏在更早的地方就被 `continue` 掉了，留着是永远走不到的死代码。**这个坑的教训是：写「逐行扫描」的逻辑时，行与行之间往往有状态。** 只要碰到成对出现的标记（围栏、注释块、HTML 标签），就该立刻想到用状态机，而不是对单行做无状态判断。

### 7.4 标题重复

最后一个坑在输出侧。纯 `.md` 文件通常第一行就是 `# 一级标题`，我把这个 H1 提取成了「文章标题」，但正文里那一行**还在**，于是前台详情页会把标题显示两遍：一次来自页面顶部的 `articleInfo.title`，一次来自 Markdown 渲染出来的 H1。

```vue
<!-- web/src/views/web/article/index.vue -->
<el-row class="title">{{ articleInfo.title }}</el-row>
<MdPreview :id="mdID" :modelValue="articleInfo.content"/>
```

所以提取标题的同时要把它从正文里摘掉，而这里用了**两个不同的正则**，分工是明确的：

```ts
// web/src/utils/markdown.ts
title = heading[1].trim();                                 // 找标题：/^#[ \t]+(.+)$/m，带 m
content = trimmedBody.replace(/^#[ \t]+.+\n?/, "").trim(); // 删标题：不带 m，只删最开头那一个
```

找标题带 `m` 标志，能匹配正文中任意位置的一级标题；删标题**不带 `m`**，只删正文最开头的那一个。为什么不统一用带 `m` 的？因为语义上我要删的是「被我当作标题的那一行」，它的位置就是正文第一行，不带 `m` 的锚点正好表达这个意思，也顺手保证了只删一处。还有一个刻意的选择：**只有当标题来自 H1 时才删。** 如果 front matter 里已经给了 `title`，正文里的 H1 我原样保留——这时作者是显式写了元数据的，正文里的 `#` 有更大概率是章节标题而不是文章标题，替他做主反而危险。

## 八、解析结果怎么分发

`parseMarkdownArticle` 返回一个规整的对象（`title` / `content` / `category` / `tags` / `abstract` / `cover`）。分发时要分两路，因为标题正文和元数据住在不同的地方：

```ts
// web/src/views/dashboard/articles/article-publish.vue
const parsed = parseMarkdownArticle(await file.text(), file.name)

title.value = parsed.title          // → 页面顶部的标题输入框
text.value = parsed.content         // → 页面上的编辑器
importedMeta.category = parsed.category    // ┐
importedMeta.tags = parsed.tags            // │ → 通过 props 传给发布弹窗
importedMeta.abstract = parsed.abstract    // │
importedMeta.cover = parsed.cover          // ┘
```

`importedMeta` 是个 `reactive`，四项都挂在上面。而表单侧有一处必须注意的细节——就是前面那段代码块里的 `tags: props.tags ? [...props.tags] : []`：**数组 props 一定要复制。** JS 里数组是引用传递，如果直接写 `tags: props.tags`，那么弹窗里点标签的 `×` 时执行的是 `articleCreateFormData.tags.splice(...)`，`splice` 是原地修改，它会**同时改掉父组件那个 `importedMeta.tags`**。症状是：在弹窗里删掉一个标签、点取消、重新打开弹窗，那个标签还是不见的，因为父组件的数据已经被悄悄改了，而 `props` 重新传入的还是那个被改过的数组。用扩展运算符复制一层就切断了引用。

这里还有个现存的小隐患值得提一下：`ArticleUpdateForm` 里有一模一样的写法，但它没有复制——`tags: props.article._source.tags`。更新表单从列表页的 ES 命中结果里初始化，所以在这里删标签会改动父组件持有的那条文章数据。目前没出问题，只是因为列表页每次操作都会重新拉取、把对象整个换掉，但它和创建表单的不一致本身就是隐患。

## 九、一个交互上的收尾：覆盖确认

导入会整段替换编辑器内容，如果正在写一篇文章、手滑点了「导入文章」再选中文件，前面的东西就没了，所以覆盖前加了确认：

```ts
// web/src/views/dashboard/articles/article-publish.vue
if (text.value.trim()) {
  try {
    await ElMessageBox.confirm('导入会覆盖当前编辑器里的标题与正文，是否继续？', '提示',
        {confirmButtonText: '覆盖导入', cancelButtonText: '取消', type: 'warning'})
  } catch {
    return
  }
}
```

两个细节。**用 `text.value.trim()` 判空而不是长度**：编辑器里留一个空格、一个换行，`length` 就不是零，会弹出多余的确认框。**`ElMessageBox.confirm` 取消时是 reject 不是 resolve**，所以必须包 `try/catch` 并在 catch 里 `return`；不写 catch 的话取消会抛一个未处理的 Promise 拒绝，控制台报错但功能看着还是对的——这种「看起来能用但控制台一片红」的问题最容易被漏掉。

导入完成后，正文也主动存一次本地缓存，并且尊重自动保存开关：

```ts
// web/src/views/dashboard/articles/article-publish.vue
if (isAutoSaveEnabled.value) {
  localStorage.setItem('article', text.value)
}
```

**为什么要手动存？** 因为 `localStorage` 的写入挂在 `onChange` 上，而 `onChange` 是给「用户在编辑器里打字」设计的，程序化地改 `v-model` 绑定的变量，编辑器不一定会回调它。赌它触发是不可靠的，不如显式写一次；开关判断也要带上，不然关掉自动保存的用户刷新一下还是能找回内容，行为就不一致了。

## 小结

把这一段做下来，几个判断标准比具体代码更值得留下来。

**样式按需引入是好事，但要接受「漏引不报错」的代价**——库把选择权交给你，就得自己记住有几份样式。**能用原生控件的地方就用原生控件**：文件选择这种一次性动作，`el-upload` 的整套状态机是负担，原生 `input` 加一个自己样式的按钮反而更可控，代价是自己处理 `value` 重置这类浏览器怪癖。

**依赖该不该引，看输入是否受控。** front matter 是「我自己写的内容」，手写子集解析器两百行搞定，边界写在注释里；哪天要支持从任意平台导入，就必须换成正经的 yaml 库。这不是「手写更优雅」，而是覆盖率和体积之间的具体权衡。

**平台差异要在边界一次性抹平。** BOM 和 CRLF 这两个坑，如果在每个正则上都写 `\r?\n`、每处都判 BOM，代码会脏得没法看；放在解析入口做一次归一化，后面所有逻辑都能干净地假设输入是规范的——这是整个解析器里最划算的两行代码。**成对标记的处理一定要用状态**：代码围栏那个 bug 说明「跳过标记行」和「跳过标记区间」是两件事。

最后一条是设计层面的：**`reactive` 只在初始化时读一次 props，这个语义决定了弹窗必须销毁重建。** 它不是配置项，而是必须理解的行为——想通这一点，后面「导入的四个字段怎么送进弹窗」才有答案。而如果字段继续变多，把草稿抽成独立 store、让页面和弹窗共享同一份数据，会比继续加 props 更值得做。
