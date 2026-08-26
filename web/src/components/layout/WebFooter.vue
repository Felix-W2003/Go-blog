<template>
  <footer class="web-footer">

    <!-- 顶部装饰线 -->
    <div class="footer-line"></div>

    <!-- 背景装饰 -->
    <div class="footer-grid"></div>
    <div class="footer-glow"></div>

    <div class="container">

      <!-- ================= 左侧 ================= -->
      <div class="footer-left">

        <div class="brand">
          <span class="brand-main">FELIX</span>
          <span class="brand-sub">BLOG</span>
        </div>

        <div class="description">
          {{ websiteStore.state.websiteInfo.description }}
        </div>

        <div class="slogan">
          RECORD · SHARE · EXPLORE
        </div>

      </div>


      <!-- ================= 中间 ================= -->
      <div class="footer-center">

        <!-- 导航 -->
        <div class="navigation">

          <div class="section-title">
            NAVIGATION
          </div>

          <div class="footer-links">

            <el-link
                v-for="item in footerLinkList"
                :key="item.title"
                :href="item.link"
                :underline="false"
            >
            {{ item.title }}
            </el-link>

          </div>

        </div>


        <!-- 网站运行时间 -->
        <div class="website-status">

          <div class="section-title">
            WEBSITE STATUS
          </div>

          <div class="status-text">

            <span>
              建站日期：
            </span>

            <strong>
              {{ websiteStore.state.websiteInfo.created_at }}
            </strong>

          </div>

          <div class="running-time">

            <span class="status-dot"></span>

            <span>
              网站已运行
            </span>

            <strong>
              {{ elapsedTime }}
            </strong>

          </div>

        </div>


        <!-- 备案 -->
        <div class="filing">

          <el-image
              src="/image/filing.png"
              alt=""
          />

          <el-link
              href="https://beian.miit.gov.cn/#/Integrated/index"
              :underline="false"
          >
            {{ websiteStore.state.websiteInfo.icp_filing }}
          </el-link>

          <el-link
              :href="publicSecurityFilingLink"
              :underline="false"
          >
            {{ websiteStore.state.websiteInfo.public_security_filing }}
          </el-link>

        </div>


        <!-- 底部 -->
        <div class="footer-bottom">

          <div class="version">

            <span class="version-label">
              VERSION
            </span>

            <span class="version-number">
              {{ websiteStore.state.websiteInfo.version }}
            </span>

          </div>


          <div class="social-link">

            <el-link
                v-for="socialLink in socialLinks"
                :key="socialLink.url"
                :href="socialLink.url"
                :underline="false"
                target="_blank"
            >

              <el-image
                  :src="socialLink.src"
                  :alt="socialLink.alt"
              />

            </el-link>

          </div>

        </div>

      </div>


      <!-- ================= 右侧人物 ================= -->
      <div class="footer-right">

    

        <el-image
            src="/image/Felix.png"
            alt=""
        />

      </div>

    </div>


    <!-- ================= 最底部版权 ================= -->

    <div class="copyright">

      <span>
        © {{ new Date().getFullYear() }} FELIX BLOG
      </span>

      <span class="copyright-divider">
        /
      </span>

      <span>
        ALL RIGHTS RESERVED.
      </span>

    </div>

  </footer>
</template>


<script setup lang="ts">

import { useWebsiteStore } from "@/stores/website";
import { computed, ref, onUnmounted } from "vue";
import {
  type FooterLink,
  websiteFooterLink
} from "@/api/website";


/* ===============================
   Website
================================ */

const websiteStore = useWebsiteStore();


/* ===============================
   Footer Link
================================ */

const footerLinkList = ref<FooterLink[]>([]);


const getFooterLinkList = async () => {

  const res = await websiteFooterLink();

  if (res.code === 0) {

    footerLinkList.value = res.data;

  }

};


getFooterLinkList();


/* ===============================
   Website Running Time
================================ */

let timerId: number | null = null;

const elapsedTime = ref("");


function updateElapsedTime() {

  const creationDate =
      websiteStore.state.websiteInfo.created_at;

  if (!creationDate) {

    return;

  }


  const creationTimestamp =
      new Date(creationDate).getTime();

  const currentTimestamp =
      Date.now();

  const diff =
      Math.max(
          0,
          currentTimestamp - creationTimestamp
      );


  const totalSeconds =
      Math.floor(diff / 1000);

  const days =
      Math.floor(
          totalSeconds / 86400
      );

  const hours =
      Math.floor(
          (totalSeconds % 86400) / 3600
      );

  const minutes =
      Math.floor(
          (totalSeconds % 3600) / 60
      );

  const seconds =
      totalSeconds % 60;


  elapsedTime.value =
      `${days}天${hours}时${minutes}分${seconds}秒`;

}


function initializeTimer() {

  updateElapsedTime();

  timerId =
      window.setInterval(
          updateElapsedTime,
          1000
      );

}


initializeTimer();


onUnmounted(() => {

  if (timerId !== null) {

    clearInterval(timerId);

  }

});


/* ===============================
   Filing
================================ */

const publicSecurityFilingLink = computed(() => {

  const filing =
      websiteStore.state.websiteInfo
          .public_security_filing;

  const match =
      filing?.match(/\d+/);

  return match
      ? `http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=${match[0]}`
      : "";

});


/* ===============================
   Social
================================ */

const socialLinks = computed(() => [

  {
    src: "/image/bilibili.png",
    alt: "Bilibili",
    url:
        websiteStore.state.websiteInfo
            .bilibili_url
  },

  {
    src: "/image/gitee.png",
    alt: "Gitee",
    url:
        websiteStore.state.websiteInfo
            .gitee_url
  },

  {
    src: "/image/github.png",
    alt: "GitHub",
    url:
        websiteStore.state.websiteInfo
            .github_url
  }

]);

</script>


<style scoped lang="scss">

/* =====================================================
   Footer
===================================================== */
.web-footer {
  position: relative;
  width: 100%;
  overflow: hidden;
  background: #0d0d0d;
  color: #fff;
}


/* =====================================================
   顶部装饰线
===================================================== */

.footer-line {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 1px;
  background:
      linear-gradient(
          90deg,
          transparent,
          #555,
          #fff,
          #555,
          transparent
      );

  opacity: 0.7;

}


/* =====================================================
   背景网格
===================================================== */

.footer-grid {

  position: absolute;

  inset: 0;

  pointer-events: none;

  opacity: 0.095;

  background-image:

      linear-gradient(
          rgb(255, 255, 255) 1px,
          transparent 1px
      ),

      linear-gradient(
          90deg,
          rgba(255,255,255,0.8) 1px,
          transparent 1px
      );

  background-size: 40px 40px;

}


/* =====================================================
   背景光晕
===================================================== */

.footer-glow {

  position: absolute;

  width: 500px;
  height: 500px;

  right: 500px;
  bottom: 300px;

  border-radius: 50%;

  background:
      radial-gradient(
          circle,
          rgba(255, 255, 255, 0.479),
          transparent 70%
      );

  filter: blur(30px);

  pointer-events: none;

  animation:
      footerGlow 8s ease-in-out infinite
      alternate;

}


@keyframes footerGlow {

  from {
    transform: translate(0, 0);
  }

  to {
    transform: translate(-80px, -40px);
  }

}


/* =====================================================
   Container
===================================================== */

.container {

  position: relative;

  display: flex;

  max-width: 1400px;

  width: 100%;

  min-height: 330px;

  margin: 0 auto;

  padding: 65px 48px 45px;

  box-sizing: border-box;

}


/* =====================================================
   左侧
===================================================== */

.footer-left {
  width: 30%;
  padding-top: 8px;

}


/* =====================================================
   Logo
===================================================== */

.brand {

  display: flex;

  align-items: baseline;

  gap: 12px;

  margin-bottom: 22px;

}


.brand-main {

  font-size: 30px;

  font-weight: 700;

  letter-spacing: 2px;

  color: #ffffff;

}


.brand-sub {

  font-size: 14px;

  letter-spacing: 4px;

  color: #aaa;

}


/* =====================================================
   描述
===================================================== */

.description {

  max-width: 310px;

  color: #c0c0c0;

  font-size: 15px;

  line-height: 1.9;

}


/* =====================================================
   Slogan
===================================================== */

.slogan {

  margin-top: 28px;

  font-family: monospace;

  font-size: 11px;

  letter-spacing: 3px;

  color: #888;

}


/* =====================================================
   中间
===================================================== */

.footer-center {

  width: 50%;

}


/* =====================================================
   小标题
===================================================== */

.section-title {

  margin-bottom: 15px;

  font-family: monospace;

  font-size: 11px;

  font-weight: 600;

  letter-spacing: 2px;

  color: #999;

}


/* =====================================================
   导航
===================================================== */

.navigation {

  margin-bottom: 30px;

}


.footer-links {

  display: flex;

  flex-wrap: wrap;

  gap: 10px 30px;

}


.footer-links :deep(.el-link) {

  color: #d0d0d0;

  font-size: 15px;

  transition:
      color 0.25s ease,
      transform 0.25s ease;

}


.footer-links :deep(.el-link:hover) {

  color: #ffffff;

  transform:
      translateY(-2px);

}


/* =====================================================
   网站状态
===================================================== */

.website-status {

  margin-bottom: 22px;

}


.status-text {

  margin-bottom: 10px;

  font-size: 14px;

  color: #aaa;

}


.status-text strong {

  margin-left: 5px;

  color: #ddd;

  font-weight: 400;

}


.running-time {

  display: flex;

  align-items: center;

  gap: 9px;

  font-size: 14px;

  color: #aaa;

}


.running-time strong {

  color: #e0e0e0;

  font-weight: 400;

  font-family: monospace;

  letter-spacing: 0.5px;

}


/* =====================================================
   状态灯
===================================================== */

.status-dot {

  width: 7px;

  height: 7px;

  flex-shrink: 0;

  border-radius: 50%;

  background: #d8d8d8;

  box-shadow:
      0 0 10px
      rgba(255,255,255,0.65);

  animation:
      statusPulse 2s infinite;

}


@keyframes statusPulse {

  0%,
  100% {

    opacity: 0.55;

    transform: scale(0.85);

  }

  50% {

    opacity: 1;

    transform: scale(1.15);

  }

}


/* =====================================================
   备案
===================================================== */

.filing {

  display: flex;

  align-items: center;

  gap: 10px;

  padding-top: 17px;

  border-top:
      1px solid #303030;

}


.filing :deep(.el-image) {

  width: 16px;

  height: 16px;

  opacity: 0.8;

}


.filing :deep(.el-link) {

  color: #aaa;

  font-size: 13px;

  transition:
      color 0.25s ease;

}


.filing :deep(.el-link:hover) {

  color: #ffffff;

}


/* =====================================================
   Bottom
===================================================== */

.footer-bottom {

  display: flex;

  align-items: center;

  justify-content: space-between;

  margin-top: 22px;

}


/* =====================================================
   Version
===================================================== */

.version {

  display: flex;

  align-items: center;

  overflow: hidden;

  border:
      1px solid #383838;

  border-radius: 4px;

}


.version-label {

  padding: 6px 10px;

  background: #222;

  color: #aaa;

  font-family: monospace;

  font-size: 10px;

  font-weight: 600;

  letter-spacing: 1px;

}


.version-number {

  padding: 6px 11px;

  color: #ddd;

  font-family: monospace;

  font-size: 11px;

  background: #151515;

}


/* =====================================================
   社交图标
===================================================== */

.social-link {

  display: flex;

  align-items: center;

  gap: 14px;

}


.social-link :deep(.el-link) {

  display: flex;

  align-items: center;

  justify-content: center;

  width: 34px;

  height: 34px;

  border:
      1px solid #383838;

  border-radius: 50%;

  background: #181818;

  transition:
      transform 0.3s ease,
      border-color 0.3s ease,
      background 0.3s ease,
      box-shadow 0.3s ease;

}


.social-link :deep(.el-link:hover) {

  transform:
      translateY(-4px);

  border-color: #777;

  background: #242424;

  box-shadow:
      0 8px 20px
      rgba(0,0,0,0.5);

}


.social-link :deep(.el-image) {

  width: 17px;

  height: 17px;

  object-fit: contain;

  opacity: 0.8;

  transition:
      opacity 0.25s ease;

}


.social-link :deep(.el-link:hover .el-image) {

  opacity: 1;

}


/* =====================================================
   右侧人物
===================================================== */

.footer-right {

  position: relative;

  display: flex;

  align-items: flex-start;

  justify-content: flex-end;

  width: 20%;

}


/* 人物 */

.footer-right :deep(.el-image) {

  position: relative;

  z-index: 2;

  width: auto;

  height: 150px;

  margin-bottom: -45px;

  filter:
      grayscale(10%);

  transition:
      transform 0.5s ease,
      filter 0.5s ease;

}


.footer-right:hover :deep(.el-image) {

  transform:
      translateY(-10px);

  filter:
      grayscale(0%);

}


/* =====================================================
   Copyright
===================================================== */

.copyright {

  position: relative;

  display: flex;

  justify-content: center;

  align-items: center;

  gap: 12px;

  padding:20px;

  border-top:
      1px solid #282828;

  color: #777;

  font-family: monospace;

  font-size: 11px;

  letter-spacing: 1px;

}


.copyright-divider {

  color: #555;

}


/* =====================================================
   Responsive
===================================================== */
@media screen and (max-width: 900px) {
  .container {

    padding:
        50px 25px 40px;

  }

  .footer-left {

    width: 30%;

  }

  .footer-center {

    width: 55%;

  }

  .footer-right {

    width: 15%;

  }

  .footer-right :deep(.el-image) {

    height: 190px;

  }

}


@media screen and (max-width: 768px) {

  .container {

    flex-direction: column;

    gap: 40px;

  }

  .footer-left,
  .footer-center,
  .footer-right {

    width: 100%;

  }

  .footer-right {

    display: none;

  }

  .footer-bottom {

    flex-direction: column;

    align-items: flex-start;

    gap: 20px;

  }

  .copyright {

    font-size: 9px;

  }

}

</style>