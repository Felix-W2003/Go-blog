<template>
  <div class="index">
    <!-- ==============================
         顶部导航
    =============================== -->
    <WebNavbar />
    <!-- ==============================
         Hero
    =============================== -->
    <section class="hero-section">
      <Carousel />
    </section>
    <!-- ==============================
         首页主体
    =============================== -->
    <main class="main-content">
      <div class="container">
        <!-- ==========================
             左侧主要内容
        =========================== -->
        <section class="main-column">
          <DailyNews />
          <Advertisement />
          <ArticleList />
        </section>
        <!-- ==========================
             右侧 Sidebar
        =========================== -->
        <aside class="sidebar">

          <ProfileCard />
                    <Calendar />
          <TagCloud />
          <RecentComments />
          <Feedback />
        </aside>
      </div>
    </main>
  </div>
</template>
<script setup lang="ts">
import Carousel from "@/components/layout/Carousel.vue";
import WebNavbar from "@/components/layout/WebNavbar.vue";
import Advertisement from "@/components/pages/Advertisement.vue";
import DailyNews from "@/components/pages/DailyNews.vue";
import ArticleList from "@/components/pages/ArticleList.vue";
import Calendar from "@/components/pages/Calendar.vue";
import ProfileCard from "@/components/pages/ProfileCard.vue";
import TagCloud from "@/components/pages/TagCloud.vue";
import RecentComments from "@/components/pages/RecentComments.vue";
import Feedback from "@/components/pages/Feedback.vue";
</script>
<style scoped lang="scss">

/* =====================================================
   首页
===================================================== */
.index {
  min-height: 100vh;
  background: #fcfcfc;
  color: #222;
}
/* =====================================================
   Hero
===================================================== */
.hero-section {
  width: 100%;
  padding-top: 100px;
 }
/* =====================================================
   Main
===================================================== */
.main-content {
  width: 100%;
  padding:
    45px 0 80px;
  box-sizing: border-box;
}
/* =====================================================
   Main Container
===================================================== */
.container {
  display: flex;
  align-items: flex-start;
  gap: 30px;
  width: calc(100% - 80px);
  max-width: 1400px;
  margin: 0 auto;
}
/* =====================================================
   Main Column
===================================================== */
.main-column {
  flex: 1;
  min-width: 0;
}
/* =====================================================
   Sidebar
===================================================== */
.sidebar {
  flex: 0 0 320px;
  width: 320px;
}

/* =====================================================
   主内容组件间距
===================================================== */

.main-column > * {
  margin-bottom: 24px;
}


.main-column > *:last-child {

  margin-bottom: 0;

}


/* =====================================================
   Sidebar 组件间距
===================================================== */

.sidebar > * {
  margin-bottom: 20px;
}
.sidebar > *:last-child {
  margin-bottom: 0;
}
/* =====================================================
   Tablet
===================================================== */

@media screen and (max-width: 1100px) {
  .container {
    width: calc(100% - 50px);
    gap: 20px;
  }
  .sidebar {
    flex-basis: 280px;
    width: 280px;
  }

}

/* =====================================================
   Mobile
===================================================== */
@media screen and (max-width: 850px) {

  .main-content {

    padding-top: 30px;

  }

  .container {

    display: block;

    width: calc(100% - 30px);

  }

  .main-column {

    width: 100%;

  }

  .sidebar {

    width: 100%;

    margin-top: 30px;

  }

}


/* =====================================================
   Small Mobile
===================================================== */

@media screen and (max-width: 600px) {

  .hero-section {

    padding-top: 60px;

  }

  .main-content {

    padding:
      25px 0 50px;

  }

}

</style>