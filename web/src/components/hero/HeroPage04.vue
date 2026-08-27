<template>
  <section class="hero-page">

    <!-- ================================
         背景装饰
    ================================= -->

    <div class="background-number">
      04
    </div>

    <div class="background-line line-01"></div>
    <div class="background-line line-02"></div>
    <div class="background-line line-03"></div>


    <!-- ================================
         顶部
    ================================= -->

    <div class="top-bar">

      <div class="page-number">
        04 / ABOUT
      </div>

      <div class="top-right">
        PERSONAL BLOG · 2026
      </div>

    </div>


    <!-- ================================
         主体
    ================================= -->

    <div class="hero-container">

      <!-- ==========================
           左侧
      =========================== -->

      <div class="hero-left">

        <div class="small-title">
          FELIX BLOG
        </div>

        <div class="vertical-line"></div>

        <div class="left-description">
          PERSONAL
          <br>
          TECHNICAL
          <br>
          JOURNAL
        </div>

      </div>


      <!-- ==========================
           中间
      =========================== -->

      <div class="hero-main">

        <div class="eyebrow">

          <span></span>

          ABOUT THIS BLOG

        </div>


        <h1>

          记录技术，
          <br>

          <strong>分享思考。</strong>

        </h1>


        <p class="description">

          这里记录我在技术学习、项目实践以及
          日常探索中的一些思考。

          <br>

          希望把走过的路留下来，
          也希望这些内容能够帮助正在前进的人。

        </p>


        <!-- ======================
             Quote
        ======================= -->

        <div class="quote">

          <span class="quote-mark">
            “
          </span>

          <div>

            <p>
              Stay curious.
            </p>

            <p>
              Keep building.
            </p>

            <p>
              Never stop exploring.
            </p>

          </div>

        </div>


        <!-- ======================
             Button
        ======================= -->

        <button
          class="about-button"
          @click="$emit('about')"
        >

          <span>
            了解更多
          </span>

          <span class="arrow">
            →
          </span>

        </button>

      </div>


      <!-- ==========================
           右侧信息
      =========================== -->

      <div class="hero-right">

        <div class="profile-card">

          <!-- 顶部 -->

          <div class="card-header">

            <span>
              PROFILE
            </span>

            <span>
              001
            </span>

          </div>


          <div class="card-line"></div>


          <!-- 名字 -->

          <div class="profile-name">

            FELIX

            <span>
              / DEVELOPER
            </span>

          </div>


          <!-- 描述 -->

          <p class="profile-description">

            Backend Developer
            <br>
            Technology Explorer
            <br>
            Lifelong Learner

          </p>


          <!-- 技术 -->

          <div class="skills">

            <div class="skill">
              <span>GO</span>
              <span>01</span>
            </div>

            <div class="skill">
              <span>GIN</span>
              <span>02</span>
            </div>

            <div class="skill">
              <span>MYSQL</span>
              <span>03</span>
            </div>

            <div class="skill">
              <span>REDIS</span>
              <span>04</span>
            </div>

            <div class="skill">
              <span>DOCKER</span>
              <span>05</span>
            </div>

          </div>


          <!-- 底部 -->

          <div class="card-footer">

            <span>
              SHANGHAI
            </span>

            <span>
              2026
            </span>

          </div>

        </div>


        <!-- 圆形装饰 -->

        <div class="circle-decoration">

          <div class="circle-inner"></div>

          <span>
            KEEP
          </span>

        </div>

      </div>

    </div>


    <!-- ================================
         底部
    ================================= -->

    <div class="bottom-bar">

      <div>
        FELIX BLOG
      </div>

      <div class="bottom-center">
        THANKS FOR READING
      </div>

      <div>
        SCROLL TO EXPLORE ↓
      </div>

    </div>

  </section>
</template>


<script setup lang="ts">

defineEmits<{
  about: []
}>()

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

  background: #f7f7f5;

  color: #181818;

}


/* =====================================================
   Background
===================================================== */

.background-number {

  position: absolute;

  right: -20px;

  bottom: -90px;

  color: rgba(0, 0, 0, 0.025);

  font-size: 260px;

  line-height: 1;

  font-weight: 800;

  letter-spacing: -20px;

  user-select: none;

}


.background-line {

  position: absolute;

  height: 1px;

  background: rgba(0, 0, 0, 0.06);

}


.line-01 {

  width: 30%;

  top: 120px;

  left: 0;

}


.line-02 {

  width: 20%;

  bottom: 95px;

  right: 0;

}


.line-03 {

  width: 1px;

  height: 100%;

  left: calc(50% + 320px);

  top: 0;

  background: rgba(0, 0, 0, 0.045);

}


/* =====================================================
   Top
===================================================== */

.top-bar {

  position: absolute;

  z-index: 5;

  top: 22px;

  left: 40px;

  right: 40px;

  display: flex;

  justify-content: space-between;

  align-items: center;

  color: #999;

  font-size: 8px;

  letter-spacing: 2px;

}


.page-number {

  color: #222;

  font-weight: 600;

}


.top-right {

  color: #aaa;

}


/* =====================================================
   Container
===================================================== */

.hero-container {

  position: relative;

  z-index: 3;

  width: 1180px;

  height: 100%;

  margin: 0 auto;

  display: flex;

  align-items: center;

}


/* =====================================================
   Left
===================================================== */

.hero-left {

  width: 150px;

  height: 245px;

  padding-top: 20px;

  border-right: 1px solid #dededb;

}


.small-title {

  color: #222;

  font-size: 12px;

  font-weight: 700;

  letter-spacing: 2px;

}


.vertical-line {

  width: 20px;

  height: 1px;

  margin-top: 25px;

  margin-bottom: 22px;

  background: #222;

}


.left-description {

  color: #aaa;

  font-size: 8px;

  line-height: 2;

  letter-spacing: 2px;

}


/* =====================================================
   Main
===================================================== */

.hero-main {

  width: 560px;

  padding-left: 55px;

}


/* =====================================================
   Eyebrow
===================================================== */

.eyebrow {

  display: flex;

  align-items: center;

  gap: 10px;

  margin-bottom: 15px;

  color: #999;

  font-size: 8px;

  letter-spacing: 3px;

}


.eyebrow span {

  display: block;

  width: 28px;

  height: 1px;

  background: #555;

}


/* =====================================================
   Title
===================================================== */

.hero-main h1 {

  margin: 0;

  color: #1b1b1b;

  font-size: 47px;

  line-height: 1.15;

  font-weight: 300;

  letter-spacing: -3px;

}


.hero-main h1 strong {

  font-weight: 700;

}


/* =====================================================
   Description
===================================================== */

.description {

  width: 430px;

  margin-top: 17px;

  color: #888;

  font-size: 11px;

  line-height: 1.9;

}


/* =====================================================
   Quote
===================================================== */

.quote {

  display: flex;

  align-items: flex-start;

  gap: 12px;

  margin-top: 15px;

}


.quote-mark {

  color: #bbb;

  font-family: Georgia, serif;

  font-size: 30px;

  line-height: 20px;

}


.quote p {

  margin: 0;

  color: #444;

  font-family: Georgia, serif;

  font-size: 10px;

  font-style: italic;

  line-height: 1.6;

}


/* =====================================================
   Button
===================================================== */

.about-button {

  display: flex;

  align-items: center;

  justify-content: space-between;

  width: 105px;

  height: 31px;

  margin-top: 17px;

  padding: 0 11px;

  border: 1px solid #222;

  border-radius: 2px;

  background: transparent;

  color: #222;

  font-size: 9px;

  cursor: pointer;

  transition: all 0.25s ease;

}


.about-button .arrow {

  font-size: 13px;

}


.about-button:hover {

  background: #1c1c1c;

  color: white;

  transform: translateY(-2px);

}


/* =====================================================
   Right
===================================================== */

.hero-right {

  position: relative;

  width: 390px;

  height: 330px;

  margin-left: auto;

}


/* =====================================================
   Profile Card
===================================================== */

.profile-card {

  position: absolute;

  z-index: 4;

  top: 42px;

  right: 20px;

  width: 245px;

  height: 270px;

  padding: 18px;

  box-sizing: border-box;

  border: 1px solid #dcdcd9;

  background: rgba(255, 255, 255, 0.75);

  box-shadow:
    0 20px 45px rgba(0, 0, 0, 0.05);

}


/* =====================================================
   Card Header
===================================================== */

.card-header {

  display: flex;

  justify-content: space-between;

  color: #999;

  font-size: 7px;

  letter-spacing: 2px;

}


.card-header span:last-child {

  color: #bbb;

}


.card-line {

  width: 100%;

  height: 1px;

  margin-top: 13px;

  background: #e5e5e2;

}


/* =====================================================
   Name
===================================================== */

.profile-name {

  margin-top: 20px;

  color: #1c1c1c;

  font-size: 23px;

  font-weight: 700;

  letter-spacing: 1px;

}


.profile-name span {

  margin-left: 5px;

  color: #aaa;

  font-size: 7px;

  font-weight: 400;

  letter-spacing: 1px;

}


/* =====================================================
   Description
===================================================== */

.profile-description {

  margin-top: 10px;

  color: #999;

  font-size: 8px;

  line-height: 1.8;

}


/* =====================================================
   Skills
===================================================== */

.skills {

  margin-top: 17px;

}


.skill {

  display: flex;

  align-items: center;

  justify-content: space-between;

  height: 21px;

  border-bottom: 1px solid #eeeeeb;

  color: #555;

  font-size: 7px;

  letter-spacing: 1px;

}


.skill span:last-child {

  color: #bbb;

}


/* =====================================================
   Card Footer
===================================================== */

.card-footer {

  position: absolute;

  bottom: 14px;

  left: 18px;

  right: 18px;

  display: flex;

  justify-content: space-between;

  color: #aaa;

  font-size: 7px;

  letter-spacing: 1px;

}


/* =====================================================
   Circle
===================================================== */

.circle-decoration {

  position: absolute;

  right: -5px;

  top: 12px;

  width: 185px;

  height: 185px;

  border: 1px solid rgba(0, 0, 0, 0.08);

  border-radius: 50%;

}


.circle-decoration::before {

  content: "";

  position: absolute;

  inset: 18px;

  border: 1px solid rgba(0, 0, 0, 0.06);

  border-radius: 50%;

}


.circle-inner {

  position: absolute;

  left: 50%;

  top: 50%;

  width: 7px;

  height: 7px;

  border-radius: 50%;

  background: #222;

  transform: translate(-50%, -50%);

}


.circle-decoration span {

  position: absolute;

  right: -8px;

  top: 82px;

  padding: 4px 7px;

  background: #222;

  color: white;

  font-size: 7px;

  letter-spacing: 2px;

}


/* =====================================================
   Bottom
===================================================== */

.bottom-bar {

  position: absolute;

  z-index: 5;

  bottom: 16px;

  left: 40px;

  right: 40px;

  display: flex;

  align-items: center;

  color: #aaa;

  font-size: 7px;

  letter-spacing: 1.5px;

}


.bottom-bar > div:first-child {

  color: #444;

}


.bottom-center {

  position: absolute;

  left: 50%;

  transform: translateX(-50%);

  color: #bbb;

}


.bottom-bar > div:last-child {

  margin-left: auto;

}


/* =====================================================
   Responsive
===================================================== */

@media screen and (max-width: 1250px) {

  .hero-container {

    width: calc(100% - 100px);

  }

}


@media screen and (max-width: 1050px) {

  .hero-left {

    width: 120px;

  }

  .hero-main {

    padding-left: 35px;

    width: 500px;

  }

  .hero-main h1 {

    font-size: 42px;

  }

  .hero-right {

    transform: scale(0.85);

    transform-origin: center right;

  }

}


@media screen and (max-width: 900px) {

  .hero-container {

    width: calc(100% - 60px);

  }

  .hero-left {

    display: none;

  }

  .hero-main {

    padding-left: 0;

    width: 55%;

  }

  .hero-right {

    width: 40%;

    transform: scale(0.8);

  }

}


@media screen and (max-width: 760px) {

  .hero-page {

    height: 420px;

  }

  .hero-container {

    width: calc(100% - 40px);

  }

  .hero-main {

    width: 100%;

  }

  .hero-main h1 {

    font-size: 38px;

  }

  .description {

    width: 100%;

  }

  .hero-right {

    display: none;

  }

}


@media screen and (max-width: 500px) {

  .hero-container {

    width: calc(100% - 30px);

  }

  .hero-main h1 {

    font-size: 34px;

    letter-spacing: -2px;

  }

  .description {

    font-size: 10px;

  }

  .top-bar {

    left: 20px;

    right: 20px;

  }

  .top-right {

    display: none;

  }

  .bottom-bar {

    left: 20px;

    right: 20px;

  }

  .bottom-center {

    display: none;

  }

}

</style>