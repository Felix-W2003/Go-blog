# Go-blog 项目代码解读系列

这是一套围绕我个人博客系统源码的技术文章，**后端与前端各成一条线**，可以顺着读，也可以按需挑着看。

## 项目是什么

一个前后端分离的个人博客系统：包含**博客前台**（首页、文章、搜索、新闻热榜、友链、留言）与**管理后台**（内容管理、用户管理、图片管理、系统配置、数据看板）两套子系统。

| 层 | 技术选型 |
|---|---|
| 后端 | Go · Gin · GORM · RESTful API |
| 存储与检索 | MySQL 8（业务数据）· Redis（会话/计数/缓存）· Elasticsearch 8（文章正文与全文检索） |
| 认证 | JWT 双 Token（Access 2h + Refresh 7d）· RBAC 三级权限 · Redis 黑名单 |
| 中间件与组件 | Zap + Lumberjack 日志 · robfig/cron 定时任务 · base64Captcha 验证码 · urfave/cli 命令行工具 · endless 平滑重启 |
| 前端 | Vue 3 · TypeScript · Vite · Element Plus · Pinia · Vue Router · Axios · ECharts · md-editor-v3 |
| 第三方集成 | QQ 互联登录 · QQ SMTP 邮件 · 七牛云对象存储 · 高德地图 IP 定位 |

## 一次请求的流转

```
                浏览器  Vue 3 + TS + Vite
                          │  HTTP / JSON
                          ▼
        ┌─────────────────────────────────────┐
        │  Gin 路由 + 中间件                   │
        │  public │ private │ admin            │
        │  JWT 鉴权 · 结构化访问日志 · panic 恢复│
        └─────────────────┬───────────────────┘
                          ▼
              api 层：参数绑定 + 响应封装
                          ▼
              service 层：业务逻辑
                 │        │        │
        ┌────────▼──┐ ┌───▼────┐ ┌─▼──────────┐
        │  MySQL    │ │ Redis  │ │ Elasticsearch│
        │ 用户/评论 │ │ 会话   │ │ 文章正文     │
        │ 友链/图片 │ │ 计数   │ │ 全文检索     │
        │ 登录日志  │ │ 缓存   │ │ 浏览量/收藏数│
        └───────────┘ └────────┘ └─────────────┘
                          ▲
                          │ 定时任务（cron）
                浏览量批量落库 · 热门榜缓存刷新
```

## 后端系列

| # | 文章 | 一句话 |
|---|---|---|
| 1 | [架构总览与分层设计](backend-01-architecture.md) | 目录为什么这样分、一次请求经过哪些层、MySQL/Redis/ES 各自负责什么 |
| 2 | [配置管理、日志与启动流程](backend-02-config-and-logging.md) | YAML 配置驱动的取舍、Zap + Lumberjack 日志体系、结构化访问日志怎么写 |
| 3 | [JWT 双 Token 认证与 RBAC](backend-03-jwt-auth.md) | 双 Token 为什么必要、共享 jti、Redis 黑名单的 TTL 对齐，以及三个真实的坑 |
| 4 | [Elasticsearch 存储与全文检索](backend-04-elasticsearch.md) | 文章为什么只存 ES、Mapping 选型、Bool Query 组合、Painless 原子更新 |
| 5 | [浏览量削峰、定时任务与命令行工具](backend-05-views-cron-cli.md) | Redis 计数 + cron 批量落库、榜单缓存设计、urfave/cli 运维工具 |

## 前端系列

| # | 文章 | 一句话 |
|---|---|---|
| 6 | [工程组织与请求层设计](frontend-01-setup-and-request.md) | 目录划分、自动按需引入、Axios 封装与 Token 自动续签、路由守卫 |
| 7 | [Markdown 编辑器与文章导入](frontend-02-editor-and-import.md) | md-editor-v3 接入、图片上传钩子、手写 front matter 解析器与四个踩坑 |
| 8 | [组件主题化与 URL 驱动数据流](frontend-03-theme-and-dataflow.md) | 一个组件适配深色/浅色两套主题，以及搜索页「改了参数不刷新」的完整排查 |

## 推荐阅读路径

**想快速了解整体设计** → 第 1 篇（地图）→ 第 3 篇（认证）→ 第 4 篇（检索）

**只关心后端** → 1 → 2 → 3 → 4 → 5

**只关心前端** → 6 → 7 → 8

**想看清「设计与取舍」** → 第 3 篇与第 5 篇的技术决策密度最高；第 8 篇是两个完整的问题排查过程

## 关于这套代码

- 仓库规模：后端约 150 个 Go 文件，前端 37 个页面视图与 35 个组件
- 已配置 GitHub Actions：每次提交自动跑 `go vet` + `go build` + `go test -race`
- 文章中出现的所有代码片段都来自仓库的真实实现

---

> 写这套文章的初衷：这个项目从设计到上线踩过的坑，比最终写出来的代码更有价值。把「为什么这么做」和「当时错在哪」记下来，既是复盘，也希望能帮到遇到同类问题的人。
