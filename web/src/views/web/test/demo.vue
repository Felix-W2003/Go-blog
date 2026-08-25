<template>
  <div class="archive-page">
    <!-- 顶部导航 -->
    <header class="page-header">
      <div class="logo">ARCHIVE</div>
      <nav class="nav-menu">
        <span class="nav-item">文章</span>
        <span class="nav-item">摄影</span>
        <span class="nav-item">设计</span>
        <span class="nav-item">关于</span>
      </nav>
    </header>
    <div class="divider-line"></div>

    <!-- 主内容区 -->
    <main class="archive-container">
      <!-- 首张大卡片 -->
      <div class="card-large">
        <div class="card-img-wrap">
          <img src="https://picsum.photos/id/1048/420/280" alt="cover" />
        </div>
        <div class="card-content">
          <h2 class="card-title">缓慢写作的力量</h2>
          <p class="card-subtitle">在碎片化时代重新建立长期表达</p>
          <div class="card-footer">
            <span class="date">2026.06</span>
            <span class="tag">● 随笔</span>
          </div>
        </div>
      </div>

      <!-- 下方卡片网格 -->
      <div class="card-grid">
        <div class="card-item">
          <img src="https://picsum.photos/id/1048/220/140" alt="" />
          <div class="item-text">
            <h3>缓慢写作的力量</h3>
            <p>在碎片化时代重新建立长期表达，调合倦怠与系统</p>
          </div>
        </div>
        <div class="card-item">
          <img src="https://picsum.photos/id/1048/220/140" alt="" />
          <div class="item-text">
            <h3>缓慢写作的力量</h3>
            <p>在碎片化时代重新建立长期表达，留净清晰目光对于感受</p>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
// 无业务逻辑，纯UI展示
</script>

<style scoped lang="scss">
.archive-page {
  background-color: #f7f7f4;
  min-height: 100vh;
  padding: 40px 60px;
  font-family: "PingFang SC", "Microsoft YaHei", sans-serif;
  color: #111111;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  .logo {
    font-size: 26px;
    font-weight: 700;
    letter-spacing: 2px;
  }
  .nav-menu {
    display: flex;
    gap: 36px;
    .nav-item {
      font-size: 18px;
      cursor: pointer;
      &:hover {
        opacity: 0.6;
      }
    }
  }
}

.divider-line {
  height: 1px;
  background: #222;
  margin: 16px 0 80px;
}

.archive-container {
  max-width: 1100px;
  margin: 0 auto;
}

.card-large {
  background: #ffffff;
  border-radius: 12px;
  display: flex;
  width: 720px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.04);
  margin-bottom: 64px;
  .card-img-wrap {
    flex-shrink: 0;
    img {
      width: 240px;
      height: 180px;
      object-fit: cover;
      border-radius: 6px;
    }
  }
  .card-content {
    padding-left:28px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    .card-title {
      font-size:42px;
      font-weight: bold;
      margin:0 0 8px;
    }
    .card-subtitle {
      font-size:18px;
      color:#333;
      margin:0 0 32px;
    }
    .card-footer {
      display:flex;
      justify-content: space-between;
      align-items: center;
      font-size: 20px;
      .tag {
        display:flex;
        align-items:center;
        gap:6px;
      }
    }
  }
}

.card-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap:32px;
  .card-item {
    background:#fff;
    border-radius:10px;
    display:flex;
    padding:16px;
    box-shadow: 0 1px 8px rgba(0,0,0,0.03);
    img {
      width:140px;
      height:100px;
      object-fit:cover;
      border-radius:6px;
    }
    .item-text {
      padding-left:16px;
      h3 {
        font-size:26px;
        margin:4px 0 8px;
      }
      p {
        font-size:14px;
        color:#555;
        margin:0;
      }
    }
  }
}
</style>
