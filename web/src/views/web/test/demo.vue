<template>
  <div class="blog-home">

    <!-- ==================== 导航栏 ==================== -->
    <header
      class="navbar"
      :class="{ 'navbar-scrolled': isScrolled }"
    >
      <div class="navbar-container">

        <!-- Logo -->
        <div class="logo">
          <span class="logo-main">FELIX</span>
          <span class="logo-sub">BLOG</span>
        </div>

        <!-- PC菜单 -->
        <nav class="nav-menu">
          <div
            v-for="item in menuList"
            :key="item.name"
            class="nav-item"
            :class="{ active: activeMenu === item.name }"
            @click="activeMenu = item.name"
          >
            {{ item.name }}
          </div>
        </nav>

        <!-- 右侧操作 -->
        <div class="nav-actions">

          <div class="search-box">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索文章..."
              :prefix-icon="Search"
              clearable
            />
          </div>

          <el-icon class="github-icon">
            <svg
              viewBox="0 0 24 24"
              width="22"
              height="22"
              fill="currentColor"
            >
              <path
                d="M12 .5C5.65.5.5 5.65.5 12c0 5.08 3.29 9.39 7.86 10.91.57.1.78-.25.78-.55v-2.14c-3.2.69-3.87-1.35-3.87-1.35-.52-1.33-1.28-1.69-1.28-1.69-1.04-.71.08-.7.08-.7 1.15.08 1.76 1.18 1.76 1.18 1.02 1.75 2.68 1.25 3.34.96.1-.75.4-1.25.73-1.54-2.55-.29-5.23-1.28-5.23-5.7 0-1.26.45-2.29 1.18-3.1-.12-.29-.51-1.46.11-3.05 0 0 .96-.31 3.15 1.18a10.93 10.93 0 0 1 5.74 0c2.19-1.49 3.15-1.18 3.15-1.18.62 1.59.23 2.76.11 3.05.73.81 1.18 1.84 1.18 3.1 0 4.43-2.69 5.41-5.25 5.69.41.36.78 1.08.78 2.18v3.23c0 .3.2.65.79.54A11.5 11.5 0 0 0 23.5 12C23.5 5.65 18.35.5 12 .5Z"
              />
            </svg>
          </el-icon>

        </div>

      </div>
    </header>


    <!-- ==================== Hero轮播 ==================== -->
    <section class="hero">

      <el-carousel
        height="520px"
        indicator-position="outside"
        :interval="5000"
        arrow="always"
      >

        <el-carousel-item
          v-for="item in banners"
          :key="item.title"
        >

          <div
            class="hero-slide"
            :style="{ backgroundImage: `url(${item.image})` }"
          >

            <div class="hero-mask"></div>

            <div class="hero-content">

              <div class="hero-label">
                FEATURED
              </div>

              <h1>
                {{ item.title }}
              </h1>

              <p>
                {{ item.description }}
              </p>

              <el-button
                class="hero-button"
                @click="handleRead(item)"
              >
                阅读全文
                <el-icon>
                  <ArrowRight />
                </el-icon>
              </el-button>

            </div>

          </div>

        </el-carousel-item>

      </el-carousel>

    </section>


    <!-- ==================== 主体 ==================== -->
    <main class="main-container">

      <!-- 左侧 -->
      <section class="article-section">

        <div class="section-header">
          <div>
            <h2>最新文章</h2>
            <span></span>
          </div>

          <el-button text>
            查看全部
            <el-icon>
              <ArrowRight />
            </el-icon>
          </el-button>
        </div>


        <!-- 文章列表 -->
        <div class="article-list">

          <article
            v-for="article in articles"
            :key="article.id"
            class="article-card"
          >

            <div class="article-image">
              <img
                :src="article.image"
                :alt="article.title"
              />
            </div>

            <div class="article-content">

              <div class="article-meta">

                <el-tag
                  size="small"
                  effect="plain"
                >
                  {{ article.category }}
                </el-tag>

                <span>
                  <el-icon>
                    <Calendar />
                  </el-icon>
                  {{ article.date }}
                </span>

                <span>
                  <el-icon>
                    <View />
                  </el-icon>
                  {{ article.views }}
                </span>

                <span>
                  <el-icon>
                    <ChatDotRound />
                  </el-icon>
                  {{ article.comments }}
                </span>

              </div>

              <h3>
                {{ article.title }}
              </h3>

              <p>
                {{ article.description }}
              </p>

              <div class="read-more">
                阅读文章
                <el-icon>
                  <ArrowRight />
                </el-icon>
              </div>

            </div>

          </article>

        </div>


        <!-- 加载更多 -->
        <div class="load-more">

          <el-button
            plain
            class="load-button"
          >
            加载更多文章
            <el-icon>
              <ArrowDown />
            </el-icon>
          </el-button>

        </div>

      </section>


      <!-- ==================== 右侧Sidebar ==================== -->
      <aside class="sidebar">


        <!-- 关于我 -->
        <div class="sidebar-card">

          <div class="sidebar-title">
            <h3>关于我</h3>
          </div>

          <div class="profile">

            <div class="avatar">
              F
            </div>

            <div>
              <h4>Felix</h4>
              <p>后端开发者 · 技术爱好者</p>
            </div>

          </div>

          <p class="profile-description">
            热爱编程与技术，喜欢记录与分享，
            相信分享的力量可以让我们走得更远。
          </p>

          <div class="social-list">

            <el-icon>
              <Link />
            </el-icon>

            <el-icon>
              <Message />
            </el-icon>

            <span class="zhihu">
              知
            </span>

          </div>

        </div>


        <!-- 热门标签 -->
        <div class="sidebar-card">

          <div class="sidebar-title">
            <h3>热门标签</h3>
          </div>

          <div class="tags">

            <el-tag
              v-for="tag in tags"
              :key="tag"
              effect="plain"
            >
              {{ tag }}
            </el-tag>

          </div>

        </div>


        <!-- 归档 -->
        <div class="sidebar-card">

          <div class="sidebar-title">
            <h3>文章归档</h3>
          </div>

          <div class="archives">

            <div
              v-for="item in archives"
              :key="item.month"
              class="archive-item"
            >
              <span>
                {{ item.month }}
              </span>

              <span>
                {{ item.count }}
              </span>
            </div>

          </div>

          <div class="more-archive">
            查看更多归档 →
          </div>

        </div>


        <!-- 热门文章 -->
        <div class="sidebar-card">

          <div class="sidebar-title">
            <h3>热门文章</h3>
          </div>

          <div class="popular-list">

            <div
              v-for="article in popularArticles"
              :key="article.id"
              class="popular-item"
            >

              <img
                :src="article.image"
                :alt="article.title"
              />

              <div>

                <h4>
                  {{ article.title }}
                </h4>

                <span>
                  <el-icon>
                    <View />
                  </el-icon>

                  {{ article.views }}
                </span>

              </div>

            </div>

          </div>

        </div>

      </aside>

    </main>


    <!-- ==================== Footer ==================== -->
    <footer class="footer">

      <div class="footer-container">

        <div class="footer-brand">

          <h2>
            FELIX BLOG
          </h2>

          <p>
            记录技术，分享思考，
            <br />
            持续成长，探索无限可能。
          </p>

          <div class="footer-social">

            <el-icon>
              <Link />
            </el-icon>

            <el-icon>
              <Message />
            </el-icon>

            <span>
              知
            </span>

          </div>

        </div>


        <div class="footer-column">

          <h3>
            导航
          </h3>

          <a
            v-for="item in menuList"
            :key="item.name"
          >
            {{ item.name }}
          </a>

        </div>


        <div class="footer-column">

          <h3>
            资源
          </h3>

          <a>GitHub</a>
          <a>掘金</a>
          <a>知乎</a>
          <a>个人简历</a>

        </div>


        <div class="footer-column subscribe">

          <h3>
            订阅更新
          </h3>

          <p>
            订阅我的博客，接收最新文章更新通知。
          </p>

          <div class="subscribe-box">

            <el-input
              placeholder="输入你的邮箱"
            />

            <el-button>
              订阅
            </el-button>

          </div>

        </div>

      </div>


      <div class="copyright">
        © 2026 Felix Blog. All rights reserved.
      </div>

    </footer>


    <!-- 返回顶部 -->
    <el-backtop
      :right="30"
      :bottom="30"
    />

  </div>
</template>


<script setup lang="ts">

import {
  ref,
  onMounted,
  onBeforeUnmount
} from 'vue'

import {
  Search,
  ArrowRight,
  ArrowDown,
  Calendar,
  View,
  ChatDotRound,
  Link,
  Message
} from '@element-plus/icons-vue'


/* ==================== Navigation ==================== */

const menuList = [
  { name: '首页' },
  { name: '博客' },
  { name: '分类' },
  { name: '标签' },
  { name: '归档' },
  { name: '关于我' }
]

const activeMenu = ref('首页')

const searchKeyword = ref('')


/* ==================== Navbar Scroll ==================== */

const isScrolled = ref(false)

const handleScroll = () => {
  isScrolled.value = window.scrollY > 50
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', handleScroll)
})


/* ==================== Banner ==================== */

const banners = [
  {
    title: '不断探索，永不停步',
    description:
      '保持好奇，细拾未知的领域，记录成长，分享思考，在技术的世界里不断前行。',
    image:
      'https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?auto=format&fit=crop&w=1800&q=90'
  },
  {
    title: '代码改变世界',
    description:
      '用代码解决问题，用技术创造价值，让每一次学习都成为成长的一部分。',
    image:
      'https://images.unsplash.com/photo-1518770660439-4636190af475?auto=format&fit=crop&w=1800&q=90'
  },
  {
    title: '探索技术的边界',
    description:
      'Go、Redis、MySQL、Docker、Elasticsearch，持续深入后端技术。',
    image:
      'https://images.unsplash.com/photo-1558494949-ef010cbdcc31?auto=format&fit=crop&w=1800&q=90'
  }
]


/* ==================== Articles ==================== */

const articles = [
  {
    id: 1,
    title: '使用 Go + Gin 搭建高性能 Web 服务',
    category: 'Go',
    date: '2026-05-20',
    views: 1287,
    comments: 36,
    description:
      '记录使用 Go 语言与 Gin 框架搭建 RESTful API 的过程，包含项目结构设计、路由管理、中间件封装与性能优化实践。',
    image:
      'https://images.unsplash.com/photo-1515879218367-8466d910aaa4?auto=format&fit=crop&w=800&q=85'
  },
  {
    id: 2,
    title: 'Redis 核心数据结构与应用场景解析',
    category: 'Redis',
    date: '2026-05-18',
    views: 942,
    comments: 24,
    description:
      '深入理解 Redis 的五大数据结构以及其底层实现原理，并结合实际场景探讨如何在项目中合理选择与使用。',
    image:
      'https://images.unsplash.com/photo-1448375240586-882707db888b?auto=format&fit=crop&w=800&q=85'
  },
  {
    id: 3,
    title: 'Docker 入门到实践：容器化你的应用',
    category: 'Docker',
    date: '2026-05-15',
    views: 1563,
    comments: 42,
    description:
      '从基础概念到实际部署，带你掌握 Docker 的核心用法，构建可移植、可扩展的应用环境。',
    image:
      'https://images.unsplash.com/photo-1497366754035-f200968a6e72?auto=format&fit=crop&w=800&q=85'
  },
  {
    id: 4,
    title: 'MySQL 索引原理与慢查询优化',
    category: 'MySQL',
    date: '2026-05-12',
    views: 1124,
    comments: 31,
    description:
      '剖析 MySQL 索引的底层原理，常见慢查询场景分析与优化方案，提升数据库性能。',
    image:
      'https://images.unsplash.com/photo-1500534623283-312aade485b7?auto=format&fit=crop&w=800&q=85'
  }
]


/* ==================== Tags ==================== */

const tags = [
  'Go',
  'Gin',
  'Redis',
  'MySQL',
  'Docker',
  'Kubernetes',
  'Linux',
  '分布式',
  '后端',
  '算法'
]


/* ==================== Archives ==================== */

const archives = [
  { month: '2026年5月', count: 8 },
  { month: '2026年4月', count: 12 },
  { month: '2026年3月', count: 9 },
  { month: '2026年2月', count: 7 },
  { month: '2026年1月', count: 10 }
]


/* ==================== Popular ==================== */

const popularArticles = articles.slice(0, 3)


/* ==================== Events ==================== */

const handleRead = (article: any) => {
  console.log('阅读文章：', article.title)
}

</script>


<style lang="scss" scoped>

/* =========================================================
   Global
========================================================= */

.blog-home {
  min-height: 100vh;
  background: #f7f7f7;
  color: #222;
}


/* =========================================================
   Navbar
========================================================= */

.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;

  height: 72px;

  z-index: 1000;

  background: rgba(0, 0, 0, 0.88);

  transition: all 0.3s ease;

  border-bottom: 1px solid rgba(255,255,255,0.06);

  &.navbar-scrolled {
    background: rgba(255,255,255,0.88);

    backdrop-filter: blur(18px);
    -webkit-backdrop-filter: blur(18px);

    box-shadow: 0 5px 25px rgba(0,0,0,0.08);

    .logo {
      color: #111;
    }

    .nav-item {
      color: #555;

      &:hover,
      &.active {
        color: #111;
      }

      &.active::after {
        background: #111;
      }
    }

    .github-icon {
      color: #111;
    }
  }
}

.navbar-container {
  width: 1180px;
  height: 100%;
  margin: 0 auto;

  display: flex;
  align-items: center;
  justify-content: space-between;
}


/* Logo */

.logo {
  display: flex;
  align-items: center;
  gap: 6px;

  color: white;

  cursor: pointer;

  .logo-main {
    font-size: 20px;
    font-weight: 800;
    letter-spacing: 1px;
  }

  .logo-sub {
    font-size: 13px;
    font-weight: 400;
    letter-spacing: 2px;
    opacity: 0.7;
  }
}


/* Menu */

.nav-menu {
  display: flex;
  height: 100%;
  align-items: center;
  gap: 40px;
}

.nav-item {
  position: relative;

  height: 100%;

  display: flex;
  align-items: center;

  color: rgba(255,255,255,0.85);

  font-size: 14px;

  cursor: pointer;

  transition: 0.25s;

  &:hover {
    color: white;
  }

  &.active {
    color: white;
  }

  &.active::after {
    content: '';

    position: absolute;

    left: 50%;
    bottom: 0;

    width: 28px;
    height: 3px;

    transform: translateX(-50%);

    background: white;

    border-radius: 5px;
  }
}


/* Right */

.nav-actions {
  display: flex;
  align-items: center;
  gap: 20px;
}

.search-box {
  width: 145px;

  :deep(.el-input__wrapper) {
    background: rgba(255,255,255,0.08);

    box-shadow: none;

    border-radius: 20px;
  }

  :deep(.el-input__inner) {
    color: white;
    font-size: 12px;
  }

  :deep(.el-input__inner::placeholder) {
    color: rgba(255,255,255,0.5);
  }
}

.github-icon {
  color: white;

  cursor: pointer;

  transition: 0.2s;

  &:hover {
    transform: scale(1.1);
  }
}


/* =========================================================
   Hero
========================================================= */

.hero {
  margin-top: 72px;

  :deep(.el-carousel__indicators) {
    bottom: 15px;
  }

  :deep(.el-carousel__button) {
    width: 8px;
    height: 8px;
    border-radius: 50%;
  }
}

.hero-slide {
  position: relative;

  width: 100%;
  height: 520px;

  background-size: cover;
  background-position: center;
}

.hero-mask {
  position: absolute;
  inset: 0;

  background:
    linear-gradient(
      90deg,
      rgba(0,0,0,0.78),
      rgba(0,0,0,0.25),
      rgba(0,0,0,0.25)
    );
}

.hero-content {
  position: relative;

  width: 1180px;
  height: 100%;

  margin: 0 auto;

  display: flex;
  flex-direction: column;

  justify-content: center;

  color: white;

  z-index: 2;
}

.hero-label {
  margin-bottom: 20px;

  font-size: 13px;

  letter-spacing: 4px;

  opacity: 0.7;
}

.hero-content h1 {
  margin: 0 0 20px;

  font-size: 48px;

  font-weight: 700;

  letter-spacing: 3px;
}

.hero-content p {
  width: 500px;

  line-height: 1.9;

  color: rgba(255,255,255,0.8);

  font-size: 15px;

  margin-bottom: 30px;
}

.hero-button {
  width: 120px;
  height: 42px;

  border: none;

  color: #111;

  background: white;

  font-weight: 600;

  border-radius: 3px;

  transition: 0.3s;

  &:hover {
    background: #eeeeee;

    transform: translateY(-2px);
  }
}


/* =========================================================
   Main
========================================================= */

.main-container {
  width: 1180px;

  margin: 0 auto;

  padding: 65px 0 80px;

  display: grid;

  grid-template-columns: 1fr 310px;

  gap: 35px;
}


/* =========================================================
   Section Header
========================================================= */

.section-header {
  display: flex;

  justify-content: space-between;
  align-items: center;

  margin-bottom: 25px;

  h2 {
    margin: 0;

    font-size: 25px;

    font-weight: 700;
  }

  span {
    display: block;

    width: 32px;
    height: 3px;

    margin-top: 10px;

    background: #111;
  }

  .el-button {
    color: #666;
  }
}


/* =========================================================
   Article
========================================================= */

.article-list {
  display: flex;
  flex-direction: column;

  gap: 18px;
}

.article-card {
  display: flex;

  min-height: 190px;

  padding: 16px;

  background: white;

  border-radius: 4px;

  transition: all 0.3s;

  border: 1px solid #eeeeee;

  &:hover {
    transform: translateY(-3px);

    box-shadow:
      0 12px 35px rgba(0,0,0,0.08);
  }
}

.article-image {
  width: 225px;

  flex-shrink: 0;

  overflow: hidden;

  border-radius: 3px;

  img {
    width: 100%;
    height: 100%;

    object-fit: cover;

    transition: 0.5s;
  }
}

.article-card:hover {
  .article-image img {
    transform: scale(1.05);
  }
}

.article-content {
  padding: 10px 15px;

  display: flex;

  flex-direction: column;

  flex: 1;

  h3 {
    margin: 12px 0 10px;

    font-size: 19px;

    font-weight: 650;

    color: #222;

    cursor: pointer;

    transition: 0.2s;

    &:hover {
      color: #666;
    }
  }

  p {
    margin: 0;

    color: #777;

    font-size: 13px;

    line-height: 1.8;

    display: -webkit-box;

    -webkit-line-clamp: 2;

    -webkit-box-orient: vertical;

    overflow: hidden;
  }
}

.article-meta {
  display: flex;

  align-items: center;

  gap: 15px;

  color: #999;

  font-size: 12px;

  span {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  :deep(.el-tag) {
    border-color: #e5e5e5;
    color: #555;
    background: #fafafa;
  }
}

.read-more {
  margin-top: auto;

  display: flex;

  align-items: center;

  gap: 5px;

  color: #555;

  font-size: 12px;

  cursor: pointer;

  transition: 0.2s;

  &:hover {
    color: #111;

    gap: 9px;
  }
}


/* =========================================================
   Load More
========================================================= */

.load-more {
  display: flex;

  justify-content: center;

  margin-top: 35px;
}

.load-button {
  height: 42px;

  padding: 0 28px;

  color: #555;

  border-color: #ddd;

  background: white;

  &:hover {
    color: #111;
    border-color: #999;
    background: white;
  }
}


/* =========================================================
   Sidebar
========================================================= */

.sidebar {
  display: flex;

  flex-direction: column;

  gap: 18px;
}

.sidebar-card {
  padding: 22px;

  background: white;

  border: 1px solid #eeeeee;

  border-radius: 4px;
}

.sidebar-title {
  padding-bottom: 14px;

  margin-bottom: 18px;

  border-bottom: 1px solid #eeeeee;

  h3 {
    margin: 0;

    font-size: 15px;

    font-weight: 650;
  }
}


/* Profile */

.profile {
  display: flex;

  align-items: center;

  gap: 14px;
}

.avatar {
  width: 55px;
  height: 55px;

  display: flex;

  align-items: center;
  justify-content: center;

  border-radius: 50%;

  background: #111;

  color: white;

  font-size: 22px;

  font-weight: 600;
}

.profile h4 {
  margin: 0 0 5px;

  font-size: 16px;
}

.profile p {
  margin: 0;

  color: #999;

  font-size: 11px;
}

.profile-description {
  margin: 18px 0;

  color: #777;

  line-height: 1.8;

  font-size: 12px;
}

.social-list {
  display: flex;

  align-items: center;

  gap: 20px;

  color: #555;

  .el-icon {
    cursor: pointer;

    &:hover {
      color: #111;
    }
  }
}

.zhihu {
  font-weight: bold;

  cursor: pointer;
}


/* Tags */

.tags {
  display: flex;

  flex-wrap: wrap;

  gap: 8px;

  :deep(.el-tag) {
    color: #555;

    border-color: #e2e2e2;

    background: #fafafa;

    cursor: pointer;

    transition: 0.2s;

    &:hover {
      color: white;

      background: #111;

      border-color: #111;
    }
  }
}


/* Archives */

.archives {
  display: flex;

  flex-direction: column;

  gap: 15px;
}

.archive-item {
  display: flex;

  justify-content: space-between;

  color: #666;

  font-size: 13px;

  cursor: pointer;

  &:hover {
    color: #111;
  }

  span:last-child {
    color: #aaa;
  }
}

.more-archive {
  margin-top: 18px;

  color: #777;

  font-size: 12px;

  cursor: pointer;

  &:hover {
    color: #111;
  }
}


/* Popular */

.popular-list {
  display: flex;

  flex-direction: column;

  gap: 18px;
}

.popular-item {
  display: flex;

  gap: 12px;

  img {
    width: 65px;
    height: 55px;

    flex-shrink: 0;

    object-fit: cover;

    border-radius: 3px;
  }

  h4 {
    margin: 0 0 8px;

    font-size: 12px;

    line-height: 1.5;

    font-weight: 600;

    display: -webkit-box;

    -webkit-line-clamp: 2;

    -webkit-box-orient: vertical;

    overflow: hidden;

    cursor: pointer;

    &:hover {
      color: #666;
    }
  }

  span {
    display: flex;

    align-items: center;

    gap: 4px;

    color: #aaa;

    font-size: 11px;
  }
}


/* =========================================================
   Footer
========================================================= */

.footer {
  padding: 60px 0 25px;

  color: rgba(255,255,255,0.75);

  background: #111;
}

.footer-container {
  width: 1180px;

  margin: 0 auto;

  display: grid;

  grid-template-columns: 2fr 1fr 1fr 2fr;

  gap: 60px;
}

.footer-brand h2 {
  margin: 0 0 18px;

  color: white;

  font-size: 20px;

  letter-spacing: 2px;
}

.footer-brand p {
  margin: 0;

  color: #888;

  line-height: 1.8;

  font-size: 13px;
}

.footer-social {
  display: flex;

  gap: 20px;

  margin-top: 25px;

  color: #aaa;

  .el-icon,
  span {
    cursor: pointer;

    &:hover {
      color: white;
    }
  }
}

.footer-column {
  display: flex;

  flex-direction: column;

  gap: 13px;

  h3 {
    margin: 0 0 10px;

    color: white;

    font-size: 14px;
  }

  a {
    color: #888;

    font-size: 13px;

    cursor: pointer;

    transition: 0.2s;

    &:hover {
      color: white;
    }
  }
}

.subscribe {
  p {
    margin: 0;

    color: #888;

    font-size: 13px;

    line-height: 1.8;
  }
}

.subscribe-box {
  display: flex;

  margin-top: 18px;

  :deep(.el-input__wrapper) {
    background: #1d1d1d;

    box-shadow: none;

    border: 1px solid #333;

    border-radius: 3px 0 0 3px;
  }

  :deep(.el-input__inner) {
    color: white;
  }

  .el-button {
    height: 40px;

    border-radius: 0 3px 3px 0;

    color: #111;

    background: white;

    border: none;
  }
}

.copyright {
  width: 1180px;

  margin: 50px auto 0;

  padding-top: 25px;

  border-top: 1px solid #292929;

  text-align: center;

  color: #666;

  font-size: 12px;
}


/* =========================================================
   Responsive
========================================================= */

@media screen and (max-width: 1200px) {

  .navbar-container,
  .hero-content,
  .main-container,
  .footer-container,
  .copyright {
    width: calc(100% - 40px);
  }

}

@media screen and (max-width: 900px) {

  .nav-menu {
    gap: 15px;
  }

  .search-box {
    display: none;
  }

  .main-container {
    grid-template-columns: 1fr;
  }

  .sidebar {
    display: grid;

    grid-template-columns: repeat(2, 1fr);
  }

  .footer-container {
    grid-template-columns: repeat(2, 1fr);
  }

}

@media screen and (max-width: 650px) {

  .navbar {
    height: 60px;
  }

  .navbar-container {
    width: calc(100% - 25px);
  }

  .nav-menu {
    display: none;
  }

  .hero {
    margin-top: 60px;
  }

  .hero-slide {
    height: 450px;
  }

  .hero-content {
    width: calc(100% - 40px);
  }

  .hero-content h1 {
    font-size: 32px;
  }

  .hero-content p {
    width: 100%;

    font-size: 13px;
  }

  .main-container {
    width: calc(100% - 25px);

    padding-top: 40px;
  }

  .article-card {
    flex-direction: column;
  }

  .article-image {
    width: 100%;
    height: 200px;
  }

  .article-content {
    padding: 15px 5px 5px;
  }

  .sidebar {
    display: flex;
  }

  .footer-container {
    width: calc(100% - 30px);

    grid-template-columns: 1fr;

    gap: 35px;
  }

  .copyright {
    width: calc(100% - 30px);
  }

}

</style>