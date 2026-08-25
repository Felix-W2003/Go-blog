<template>
  <div class="blog-container">
    <!-- Hero 头部区域，带网格背景 -->
    <section class="hero">
      <div class="hero-content">
        <h1 class="hero-title">记录思考与生活2</h1>
        <p class="hero-sub">Design · Code · Life</p>
      </div>
      <div class="grid-bg"></div>
    </section>

    <!-- 文章卡片区域 -->
    <section class="article-list">
      <div class="article-card card‑light">
        <h3>如何建立个人博客风格</h3>
        <div class="card‑bg‑light"></div>
      </div>
      <div class="article-card card‑dark">
        <h3>黑白灰配色实践笔记</h3>
        <div class="card‑bg‑dark"></div>
      </div>
    </section>

  </div>
</template>

<script setup lang="ts">


</script>

<style lang="scss" scoped>
// 使用全局黑白灰scss变量（vite已配置additionalData自动导入theme.scss）
.blog-container {
  min-height: 100vh;
  background: $white;
  color: $black;
  margin: 0 auto;
}

// 导航栏
.blog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 48px;
  border-bottom: 1px solid $light-gray;

  .logo {
    font-size: 32px;
    font-weight: 700;
  }
  .nav-menu {
    display: flex;
    gap: 60px;
    .nav-item {
      font-size: 20px;
      color: $black;
      text-decoration: none;
      position: relative;
      &.active::after {
        content: "";
        position: absolute;
        width: 100%;
        height: 2px;
        bottom: -10px;
        left:0;
        background: $black;
      }
    }
  }
}

// Hero区域 + CSS网格背景
.hero {
  position: relative;
  padding: 80px 48px;
  border-bottom:1px solid $light-gray;
  overflow:hidden;

  .hero-content {
    position: relative;
    z-index: 2;
    h1.hero-title {
      font-size: 80px;
      font-weight:900;
      line-height:1.1;
      margin-bottom:12px;
    }
    .hero-sub {
      font-size:36px;
    }
  }
  // 网格背景，还原截图的点阵网格
  .grid-bg {
    position:absolute;
    top:0;
    left:0;
    width:100%;
    height:100%;
    background-image: 
      linear-gradient(to right, #e2e2e2 1px, transparent 1px),
      linear-gradient(to bottom, #e2e2e2 1px, transparent 1px);
    background-size:40px 40px;
    opacity:0.35;
    z-index:1;
  }
}

// 文章卡片区域
.article-list {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap:32px;
  padding:48px;

  .article-card {
    border:1px solid $light-gray;
    height:220px;
    position:relative;
    padding:32px;
    display:flex;
    flex-direction:column;
    justify-content:space-between;
    h3 {
      font-size:26px;
      font-weight:bold;
      z-index:2;
      position:relative;
    }
  }
  .card‑light {
    .card‑bg‑light {
      position:absolute;
      bottom:0;
      left:0;
      width:100%;
      height:80px;
      background:linear-gradient(90deg,#e6e6e6 65%,#f1f1f1);
    }
  }
  .card‑dark {
    background:#1a1a1a;
    color:$white;
    .card‑bg‑dark {
      position:absolute;
      bottom:0;
      left:0;
      width:100%;
      height:80px;
      display:grid;
      grid-template-columns: repeat(4,1fr);
      >div {
        height:100%;
        &:nth-child(1){background:#222222;}
        &:nth-child(2){background:#444444;}
        &:nth-child(3){background:#666666;}
        &:nth-child(4){background:#888888;}
      }
    }
  }
}

// 页脚
.blog-footer {
  padding:32px 48px;
  border-top:1px solid $light-gray;
  font-size:16px;
  color:$mid-gray;
}
</style>
