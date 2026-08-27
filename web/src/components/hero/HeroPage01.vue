<template>
  <section class="hero-page">

    <!-- 背景网格 -->
    <div class="grid-bg"></div>

    <!-- 背景光晕 -->
    <div class="glow glow-left"></div>
    <div class="glow glow-right"></div>

    <!-- 装饰线 -->
    <div class="decor-line decor-line-top"></div>
    <div class="decor-line decor-line-bottom"></div>


    <div class="hero-container">

      <!-- =========================
           左侧内容
      ========================== -->
      <div class="hero-content">

        <div class="hero-label">
          <span class="label-line"></span>
          ENGINEERING / 01
        </div>


        <h1>
          构建高性能
          <br>
          <strong>Web 服务</strong>
        </h1>


        <p class="description">
          用 Go 构建简单、可靠、可扩展的后端系统，
          从 Gin 到 Redis，再到 MySQL 与 Docker，
          探索现代 Web 服务的工程实践。
        </p>


        <div class="actions">

          <button class="primary-btn">
            阅读文章
            <span>→</span>
          </button>

          <button class="text-btn">
            查看技术栈
          </button>

        </div>


        <!-- 技术栈 -->
        <div class="tech-stack">

          <span>GO</span>
          <span>GIN</span>
          <span>REDIS</span>
          <span>MYSQL</span>
          <span>DOCKER</span>

        </div>

      </div>


      <!-- =========================
           右侧视觉
      ========================== -->
      <div class="hero-visual">

        <!-- 代码窗口 -->
        <div class="code-window">

          <div class="window-header">

            <div class="window-dots">
              <span></span>
              <span></span>
              <span></span>
            </div>

            <span class="window-title">
              server.go
            </span>

            <span class="window-status">
              ●
            </span>

          </div>


          <div class="code-content">

            <div>
              <span class="line-number">01</span>
              <span class="keyword">package</span>
              main
            </div>

            <div>
              <span class="line-number">02</span>
            </div>

            <div>
              <span class="line-number">03</span>
              <span class="keyword">import</span>
              "github.com/gin-gonic/gin"
            </div>

            <div>
              <span class="line-number">04</span>
            </div>

            <div>
              <span class="line-number">05</span>
              <span class="keyword">func</span>
              main() {
            </div>

            <div class="indent">
              <span class="line-number">06</span>
              router := gin.Default()
            </div>

            <div class="indent">
              <span class="line-number">07</span>
              router.GET(
            </div>

            <div class="indent-2">
              <span class="line-number">08</span>
              "/api/hello",
            </div>

            <div class="indent-2">
              <span class="line-number">09</span>
              func(c *gin.Context) {
            </div>

            <div class="indent-2">
              <span class="line-number">10</span>
              c.JSON(200, gin.H{
            </div>

            <div class="indent-3">
              <span class="line-number">11</span>
              "message": "hello",
            </div>

            <div class="indent-2">
              <span class="line-number">12</span>
              })
            </div>

            <div class="indent-2">
              <span class="line-number">13</span>
              })
            </div>

            <div class="indent">
              <span class="line-number">14</span>
              router.Run(":8080")
            </div>

            <div>
              <span class="line-number">15</span>
              }
            </div>

          </div>


          <div class="window-footer">

            <span>GO 1.24</span>

            <span>UTF-8</span>

            <span class="online">
              ● SERVER ONLINE
            </span>

          </div>

        </div>


        <!-- 数据卡片 -->
        <div class="stats-card">

          <div class="stat">
            <span class="stat-value">99.9%</span>
            <span class="stat-label">UPTIME</span>
          </div>

          <div class="stat">
            <span class="stat-value">12ms</span>
            <span class="stat-label">LATENCY</span>
          </div>

          <div class="stat">
            <span class="stat-value">24K</span>
            <span class="stat-label">REQUESTS</span>
          </div>

        </div>

      </div>

    </div>


    <!-- 页面编号 -->
    <div class="page-index">
      <span>01</span>
      <i></i>
      <span>04</span>
    </div>


    <!-- 右下角装饰 -->
    <div class="corner-text">
      BUILD / CREATE / EXPLORE
    </div>

  </section>
</template>


<script setup lang="ts">

</script>


<style scoped lang="scss">

/* =====================================================
   页面
===================================================== */

.hero-page {
  position: relative;

  width: 100%;
  height: 430px;

  overflow: hidden;

  color: #ffffff;

  background:
    linear-gradient(
      110deg,
      #080808 0%,
      #111111 55%,
      #181818 100%
    );
}


/* =====================================================
   背景网格
===================================================== */

.grid-bg {
  position: absolute;

  inset: 0;

  opacity: 0.16;

  background-image:
    linear-gradient(
      rgba(255, 255, 255, 0.06) 1px,
      transparent 1px
    ),
    linear-gradient(
      90deg,
      rgba(255, 255, 255, 0.06) 1px,
      transparent 1px
    );

  background-size: 42px 42px;

  mask-image:
    linear-gradient(
      to right,
      black,
      transparent 85%
    );
}


/* =====================================================
   光晕
===================================================== */

.glow {
  position: absolute;

  width: 420px;
  height: 420px;

  border-radius: 50%;

  filter: blur(120px);

  opacity: 0.12;

  pointer-events: none;
}


.glow-left {
  left: -200px;
  top: -220px;

  background: #ffffff;
}


.glow-right {
  right: -200px;
  bottom: -250px;

  background: #999999;
}


/* =====================================================
   装饰线
===================================================== */

.decor-line {
  position: absolute;

  width: 1px;

  height: 100px;

  background:
    linear-gradient(
      to bottom,
      transparent,
      rgba(255,255,255,0.2),
      transparent
    );
}


.decor-line-top {
  left: 8%;

  top: 0;
}


.decor-line-bottom {
  right: 8%;

  bottom: 0;
}


/* =====================================================
   主容器
===================================================== */

.hero-container {
  position: relative;

  z-index: 2;

  width: 1180px;
  height: 100%;

  margin: 0 auto;

  display: flex;

  align-items: center;

  justify-content: space-between;
}


/* =====================================================
   左侧
===================================================== */

.hero-content {
  width: 48%;
}


.hero-label {
  display: flex;

  align-items: center;

  gap: 10px;

  margin-bottom: 20px;

  color: #888888;

  font-size: 10px;

  letter-spacing: 3px;
}


.label-line {
  width: 30px;

  height: 1px;

  background: #777777;
}


/* =====================================================
   标题
===================================================== */

.hero-content h1 {
  margin: 0;

  font-size: 48px;

  line-height: 1.15;

  font-weight: 300;

  letter-spacing: -1px;
}


.hero-content h1 strong {
  font-weight: 700;
}


/* =====================================================
   描述
===================================================== */

.description {
  width: 500px;

  margin:
    22px 0 0;

  color: #999999;

  font-size: 13px;

  line-height: 1.9;
}


/* =====================================================
   操作
===================================================== */

.actions {
  display: flex;

  align-items: center;

  gap: 22px;

  margin-top: 26px;
}


.primary-btn {
  height: 38px;

  padding: 0 20px;

  border: none;

  border-radius: 3px;

  background: #ffffff;

  color: #111111;

  font-size: 12px;

  cursor: pointer;

  transition: all 0.25s ease;
}


.primary-btn span {
  margin-left: 10px;
}


.primary-btn:hover {
  transform: translateY(-2px);

  background: #eeeeee;
}


.text-btn {
  padding: 0;

  border: none;

  background: transparent;

  color: #888888;

  font-size: 12px;

  cursor: pointer;

  transition: color 0.2s ease;
}


.text-btn:hover {
  color: #ffffff;
}


/* =====================================================
   技术栈
===================================================== */

.tech-stack {
  display: flex;

  gap: 24px;

  margin-top: 34px;
}


.tech-stack span {
  color: #555555;

  font-size: 9px;

  letter-spacing: 2px;
}


/* =====================================================
   右侧
===================================================== */

.hero-visual {
  position: relative;

  width: 500px;
}


/* =====================================================
   代码窗口
===================================================== */

.code-window {
  width: 500px;

  overflow: hidden;

  border:
    1px solid rgba(255,255,255,0.1);

  border-radius: 6px;

  background:
    rgba(15,15,15,0.92);

  box-shadow:
    0 30px 80px rgba(0,0,0,0.45);

  backdrop-filter: blur(15px);
}


/* =====================================================
   Window Header
===================================================== */

.window-header {
  height: 42px;

  padding: 0 16px;

  display: flex;

  align-items: center;

  border-bottom:
    1px solid rgba(255,255,255,0.07);
}


.window-dots {
  display: flex;

  gap: 6px;
}


.window-dots span {
  width: 7px;
  height: 7px;

  border-radius: 50%;

  background: #555555;
}


.window-title {
  margin-left: 20px;

  color: #777777;

  font-family:
    "JetBrains Mono",
    monospace;

  font-size: 10px;
}


.window-status {
  margin-left: auto;

  color: #999999;

  font-size: 9px;
}


/* =====================================================
   Code
===================================================== */

.code-content {
  height: 255px;

  padding: 15px 18px;

  overflow: hidden;

  color: #bbbbbb;

  font-family:
    "JetBrains Mono",
    "Fira Code",
    monospace;

  font-size: 10px;

  line-height: 1.7;
}


.code-content > div {
  white-space: nowrap;
}


.line-number {
  display: inline-block;

  width: 28px;

  color: #444444;

  user-select: none;
}


.keyword {
  color: #ffffff;

  font-weight: 600;
}


.indent {
  padding-left: 20px;
}


.indent-2 {
  padding-left: 40px;
}


.indent-3 {
  padding-left: 60px;
}


/* =====================================================
   Window Footer
===================================================== */

.window-footer {
  height: 30px;

  padding: 0 16px;

  display: flex;

  align-items: center;

  gap: 18px;

  border-top:
    1px solid rgba(255,255,255,0.07);

  color: #555555;

  font-size: 8px;
}


.online {
  margin-left: auto;

  color: #888888;
}


/* =====================================================
   数据卡片
===================================================== */

.stats-card {
  position: absolute;

  right: -28px;

  bottom: -34px;

  width: 310px;

  padding: 14px 18px;

  display: flex;

  justify-content: space-between;

  border:
    1px solid rgba(255,255,255,0.1);

  border-radius: 5px;

  background:
    rgba(25,25,25,0.88);

  box-shadow:
    0 15px 40px rgba(0,0,0,0.35);

  backdrop-filter: blur(15px);
}


.stat {
  display: flex;

  flex-direction: column;

  gap: 3px;
}


.stat-value {
  color: #ffffff;

  font-size: 13px;

  font-weight: 600;
}


.stat-label {
  color: #555555;

  font-size: 8px;

  letter-spacing: 1px;
}


/* =====================================================
   页面编号
===================================================== */

.page-index {
  position: absolute;

  right: 35px;
  top: 30px;

  z-index: 3;

  display: flex;

  align-items: center;

  gap: 8px;

  color: #555555;

  font-size: 9px;

  letter-spacing: 2px;
}


.page-index span:first-child {
  color: #ffffff;
}


.page-index i {
  width: 25px;

  height: 1px;

  background: #555555;
}


/* =====================================================
   右下角文字
===================================================== */

.corner-text {
  position: absolute;

  right: 35px;
  bottom: 25px;

  color: #444444;

  font-size: 8px;

  letter-spacing: 2px;
}


/* =====================================================
   Responsive
===================================================== */

@media screen and (max-width: 1250px) {

  .hero-container {
    width:
      calc(100% - 100px);
  }

}


@media screen and (max-width: 1000px) {

  .hero-container {
    width:
      calc(100% - 70px);
  }

  .hero-visual {
    width: 420px;
  }

  .code-window {
    width: 420px;
  }

  .description {
    width: 100%;
  }

}


@media screen and (max-width: 850px) {

  .hero-content {
    width: 100%;
  }

  .hero-visual {
    display: none;
  }

  .hero-content h1 {
    font-size: 42px;
  }

}


@media screen and (max-width: 600px) {

  .hero-page {
    height: 400px;
  }

  .hero-container {
    width:
      calc(100% - 40px);
  }

  .hero-content h1 {
    font-size: 34px;
  }

  .tech-stack {
    gap: 14px;
  }

  .corner-text {
    display: none;
  }

}

</style>