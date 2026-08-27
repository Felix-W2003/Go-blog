<template>
  <section class="hero-page">

    <!-- =========================================
         背景装饰
    ========================================== -->

    <div class="background-grid"></div>

    <div class="background-circle"></div>

    <div class="background-text">
      JOURNEY
    </div>


    <!-- =========================================
         顶部信息
    ========================================== -->

    <div class="top-info">

      <span class="page-label">
        JOURNEY / 03
      </span>

      <span class="page-description">
        KEEP LEARNING · KEEP BUILDING
      </span>

    </div>


    <!-- =========================================
         主体
    ========================================== -->

    <div class="hero-container">


      <!-- =======================================
           左侧内容
      ======================================== -->

      <div class="hero-content">

        <div class="eyebrow">

          <span class="line"></span>

          MY JOURNEY

        </div>


        <h1>

          一直向前
          <br>

          <span>不断成长</span>

        </h1>


        <p class="description">
          技术不是终点，而是一场持续探索的旅程。
          从写下第一行代码，到构建真正的 Web 服务，
          每一次实践，都让我离想成为的人更近一步。
        </p>


        <button class="journey-button">

          查看我的技术栈

          <span>↗</span>

        </button>


        <!-- =====================================
             时间线
        ====================================== -->

        <div class="timeline">

          <!-- 2025 -->

          <div class="timeline-item">

            <div class="timeline-year">
              2025
            </div>

            <div class="timeline-point"></div>

            <div class="timeline-content">

              <h3>
                Backend · Start
              </h3>

              <p>
                从基础开始，探索后端开发。
              </p>

              <div class="tags">

                <span>Python</span>

                <span>MySQL</span>

                <span>Linux</span>

              </div>

            </div>

          </div>


          <!-- 2026 -->

          <div class="timeline-item active">

            <div class="timeline-year">
              2026
            </div>

            <div class="timeline-point"></div>

            <div class="timeline-content">

              <h3>
                Go · Build
              </h3>

              <p>
                深入 Go 后端，构建真实项目。
              </p>

              <div class="tags">

                <span>Go</span>

                <span>Gin</span>

                <span>Redis</span>

              </div>

            </div>

          </div>


          <!-- NEXT -->

          <div class="timeline-item next">

            <div class="timeline-year">
              NEXT
            </div>

            <div class="timeline-point"></div>

            <div class="timeline-content">

              <h3>
                Distributed Systems
              </h3>

              <p>
                探索更大的系统与更复杂的问题。
              </p>

              <div class="tags">

                <span>Docker</span>

                <span>ES</span>

                <span>Cloud</span>

              </div>

            </div>

          </div>

        </div>

      </div>


      <!-- =======================================
           右侧视觉
      ======================================== -->

      <div class="hero-visual">


        <!-- 大年份 -->

        <div class="year">

          <span class="year-main">
            26
          </span>

          <span class="year-small">
            20
          </span>

        </div>


        <!-- 圆形轨迹 -->

        <div class="orbit orbit-one"></div>

        <div class="orbit orbit-two"></div>

        <div class="orbit orbit-three"></div>


        <!-- 中心 -->

        <div class="visual-center">

          <span class="center-number">
            03
          </span>

          <span class="center-title">
            KEEP
          </span>

          <span class="center-title">
            BUILDING
          </span>

        </div>


        <!-- 浮动标签 -->

        <div class="float-card card-go">

          <span class="card-number">
            01
          </span>

          <span>
            GO
          </span>

        </div>


        <div class="float-card card-docker">

          <span class="card-number">
            02
          </span>

          <span>
            DOCKER
          </span>

        </div>


        <div class="float-card card-es">

          <span class="card-number">
            03
          </span>

          <span>
            ELASTICSEARCH
          </span>

        </div>


      </div>

    </div>


    <!-- =========================================
         底部
    ========================================== -->

    <div class="bottom-info">

      <span>
        2025 — 2026
      </span>

      <span>
        LEARNING IN PROGRESS
      </span>

      <span>
        FELIX BLOG
      </span>

    </div>

  </section>
</template>


<script setup lang="ts">

</script>


<style scoped lang="scss">


/* =====================================================
   Hero
===================================================== */

.hero-page {

  position: relative;

  width: 100%;

  height: 430px;

  overflow: hidden;

  color: #181818;

  background:
    linear-gradient(
      135deg,
      #fafafa 0%,
      #f5f5f3 55%,
      #eeeeeb 100%
    );

}


/* =====================================================
   Background Grid
===================================================== */

.background-grid {

  position: absolute;

  inset: 0;

  opacity: 0.35;

  background-image:

    linear-gradient(
      rgba(0,0,0,0.035) 1px,
      transparent 1px
    ),

    linear-gradient(
      90deg,
      rgba(0,0,0,0.035) 1px,
      transparent 1px
    );

  background-size: 45px 45px;

  mask-image:
    linear-gradient(
      to right,
      black,
      transparent 75%
    );

}


/* =====================================================
   Background Circle
===================================================== */

.background-circle {

  position: absolute;

  right: -180px;

  top: -200px;

  width: 650px;

  height: 650px;

  border-radius: 50%;

  border:
    1px solid rgba(0,0,0,0.05);

}


.background-circle::before {

  content: "";

  position: absolute;

  inset: 45px;

  border:
    1px solid rgba(0,0,0,0.04);

  border-radius: 50%;

}


.background-circle::after {

  content: "";

  position: absolute;

  inset: 100px;

  border:
    1px solid rgba(0,0,0,0.035);

  border-radius: 50%;

}


/* =====================================================
   Background Text
===================================================== */

.background-text {

  position: absolute;

  right: -35px;

  bottom: -40px;

  color:
    rgba(0,0,0,0.025);

  font-size: 150px;

  font-weight: 800;

  letter-spacing: 12px;

  pointer-events: none;

}


/* =====================================================
   Top Info
===================================================== */

.top-info {

  position: absolute;

  z-index: 10;

  left: 40px;

  right: 40px;

  top: 25px;

  display: flex;

  align-items: center;

  justify-content: space-between;

  color: #888;

  font-size: 8px;

  letter-spacing: 2px;

}


.page-label {

  color: #222;

  font-weight: 600;

}


.page-description {

  color: #aaa;

}


/* =====================================================
   Container
===================================================== */

.hero-container {

  position: relative;

  z-index: 5;

  width: 1180px;

  height: 100%;

  margin: 0 auto;

  display: flex;

  align-items: center;

}


/* =====================================================
   Hero Content
===================================================== */

.hero-content {

  width: 58%;

  padding-top: 20px;

}


/* =====================================================
   Eyebrow
===================================================== */

.eyebrow {

  display: flex;

  align-items: center;

  gap: 10px;

  margin-bottom: 14px;

  color: #999;

  font-size: 9px;

  font-weight: 500;

  letter-spacing: 3px;

}


.eyebrow .line {

  width: 30px;

  height: 1px;

  background: #777;

}


/* =====================================================
   Title
===================================================== */

.hero-content h1 {

  margin: 0;

  font-size: 48px;

  line-height: 1.1;

  font-weight: 300;

  letter-spacing: -2px;

}


.hero-content h1 span {

  font-weight: 700;

}


/* =====================================================
   Description
===================================================== */

.description {

  width: 490px;

  margin-top: 17px;

  color: #777;

  font-size: 12px;

  line-height: 1.9;

}


/* =====================================================
   Button
===================================================== */

.journey-button {

  display: flex;

  align-items: center;

  gap: 25px;

  margin-top: 17px;

  height: 34px;

  padding:
    0 14px;

  border:
    1px solid #222;

  border-radius: 2px;

  background: transparent;

  color: #222;

  font-size: 10px;

  cursor: pointer;

  transition:
    all 0.25s ease;

}


.journey-button span {

  font-size: 14px;

}


.journey-button:hover {

  background: #1c1c1c;

  color: #fff;

  transform:
    translateY(-2px);

}


/* =====================================================
   Timeline
===================================================== */

.timeline {

  display: flex;

  gap: 0;

  margin-top: 25px;

  width: 590px;

}


/* =====================================================
   Timeline Item
===================================================== */

.timeline-item {

  position: relative;

  width: 33.333%;

  padding-right: 20px;

}


.timeline-item::before {

  content: "";

  position: absolute;

  left: 0;

  right: 15px;

  top: 27px;

  height: 1px;

  background: #d5d5d2;

}


.timeline-item:last-child::before {

  display: none;

}


/* =====================================================
   Year
===================================================== */

.timeline-year {

  margin-bottom: 12px;

  color: #aaa;

  font-size: 9px;

  font-weight: 600;

  letter-spacing: 1px;

}


.timeline-item.active .timeline-year {

  color: #222;

}


/* =====================================================
   Point
===================================================== */

.timeline-point {

  position: relative;

  z-index: 2;

  width: 7px;

  height: 7px;

  margin-bottom: 13px;

  border-radius: 50%;

  border:
    2px solid #aaa;

  background: #f5f5f3;

}


.timeline-item.active .timeline-point {

  width: 9px;

  height: 9px;

  margin-top: -1px;

  border-color: #222;

  background: #222;

  box-shadow:
    0 0 0 4px rgba(0,0,0,0.06);

}


.timeline-item.next .timeline-point {

  border-style: dashed;

}


/* =====================================================
   Timeline Content
===================================================== */

.timeline-content h3 {

  margin: 0;

  color: #333;

  font-size: 10px;

  font-weight: 600;

}


.timeline-content p {

  margin:
    5px 0 7px;

  color: #999;

  font-size: 8px;

}


/* =====================================================
   Tags
===================================================== */

.tags {

  display: flex;

  gap: 5px;

}


.tags span {

  padding:
    3px 6px;

  border:
    1px solid #dededb;

  color: #999;

  font-size: 7px;

  background: rgba(255,255,255,0.5);

}


.timeline-item.active .tags span {

  border-color: #c7c7c4;

  color: #555;

}


/* =====================================================
   Hero Visual
===================================================== */

.hero-visual {

  position: relative;

  width: 440px;

  height: 360px;

  margin-left: auto;

}


/* =====================================================
   Big Year
===================================================== */

.year {

  position: absolute;

  right: 25px;

  top: 45px;

}


.year-main {

  position: relative;

  z-index: 2;

  display: block;

  color: #1a1a1a;

  font-size: 155px;

  font-weight: 800;

  line-height: 0.8;

  letter-spacing: -14px;

}


.year-small {

  position: absolute;

  left: -18px;

  top: -18px;

  color: #c9c9c5;

  font-size: 38px;

  font-weight: 300;

  letter-spacing: -3px;

}


/* =====================================================
   Orbit
===================================================== */

.orbit {

  position: absolute;

  left: 50%;

  top: 50%;

  border:
    1px solid rgba(0,0,0,0.09);

  border-radius: 50%;

  transform:
    translate(-50%, -50%)
    rotate(-20deg);

}


.orbit-one {

  width: 270px;

  height: 120px;

}


.orbit-two {

  width: 320px;

  height: 170px;

  transform:
    translate(-50%, -50%)
    rotate(30deg);

}


.orbit-three {

  width: 380px;

  height: 230px;

  transform:
    translate(-50%, -50%)
    rotate(-40deg);

  border-color:
    rgba(0,0,0,0.045);

}


/* =====================================================
   Center
===================================================== */

.visual-center {

  position: absolute;

  left: 50%;

  top: 50%;

  transform:
    translate(-50%, -50%);

  display: flex;

  flex-direction: column;

  align-items: center;

}


.center-number {

  color: #aaa;

  font-size: 8px;

  letter-spacing: 3px;

}


.center-title {

  color: #222;

  font-size: 11px;

  font-weight: 700;

  line-height: 1.3;

  letter-spacing: 3px;

}


/* =====================================================
   Floating Cards
===================================================== */

.float-card {

  position: absolute;

  z-index: 5;

  display: flex;

  align-items: center;

  gap: 9px;

  padding:
    8px 11px;

  border:
    1px solid rgba(0,0,0,0.08);

  background:
    rgba(250,250,248,0.85);

  box-shadow:
    0 10px 25px rgba(0,0,0,0.05);

  backdrop-filter:
    blur(8px);

  color: #444;

  font-size: 8px;

  letter-spacing: 1px;

}


.card-number {

  color: #aaa;

  font-size: 7px;

}


.card-go {

  left: 5px;

  top: 105px;

}


.card-docker {

  right: 10px;

  bottom: 105px;

}


.card-es {

  left: 45px;

  bottom: 35px;

}


/* =====================================================
   Bottom
===================================================== */

.bottom-info {

  position: absolute;

  z-index: 10;

  left: 40px;

  right: 40px;

  bottom: 17px;

  display: flex;

  align-items: center;

  gap: 25px;

  color: #aaa;

  font-size: 7px;

  letter-spacing: 1.5px;

}


.bottom-info span:first-child {

  color: #444;

}


.bottom-info span:last-child {

  margin-left: auto;

  color: #777;

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


@media screen and (max-width: 1050px) {

  .hero-content {

    width: 60%;

  }

  .hero-visual {

    width: 380px;

    transform:
      scale(0.9);

  }

}


@media screen and (max-width: 900px) {

  .hero-container {

    width:
      calc(100% - 70px);

  }

  .hero-content {

    width: 100%;

  }

  .hero-visual {

    width: 350px;

    transform:
      scale(0.8);

  }

  .description {

    width: 100%;

  }

}


@media screen and (max-width: 760px) {

  .hero-page {

    height: 420px;

  }

  .hero-container {

    width:
      calc(100% - 45px);

  }

  .hero-content {

    width: 100%;

  }

  .hero-content h1 {

    font-size: 40px;

  }

  .description {

    width: 100%;

  }

  .hero-visual {

    display: none;

  }

  .timeline {

    width: 100%;

  }

}


@media screen and (max-width: 500px) {

  .hero-container {

    width:
      calc(100% - 30px);

  }

  .hero-content h1 {

    font-size: 34px;

  }

  .timeline {

    display: none;

  }

  .description {

    font-size: 11px;

  }

  .top-info {

    left: 20px;

    right: 20px;

  }

  .page-description {

    display: none;

  }

  .bottom-info {

    left: 20px;

    right: 20px;

  }

}

</style>