<template>
  <el-card class="profile-card" shadow="never">
    <!-- 背景装饰 -->
    <div class="background-grid"></div>
    <div class="background-glow"></div>
    <!-- 顶部状态 -->
    <div class="profile-status">
      <span class="status-dot"></span>
      <span>AVAILABLE FOR PROJECTS</span>
    </div>
    <!-- 用户信息 -->
    <div class="profile-main">
      <!-- 头像 -->
      <div class="avatar-wrapper">
        <div class="avatar-ring"></div>
        <el-avatar
            class="avatar"
            :size="82"
            :src="userInfo.avatar"
        />

      </div>


      <!-- 名字 -->
      <div class="profile-name">
        {{ userInfo.name }}
      </div>


      <!-- 职业 -->
      <div class="profile-role">
        Backend Developer
      </div>


      <!-- 简介 -->
      <div class="profile-description">
        专注后端开发与系统设计，
        热爱技术、记录与分享。
      </div>

    </div>


    <!-- 技术栈 -->
    <div class="tech-stack">

      <div class="section-label">
        TECH STACK
      </div>

      <div class="tech-list">

        <span
            v-for="item in techList"
            :key="item.name"
            class="tech-item"
        >
          <el-icon>
            <component :is="item.icon"/>
          </el-icon>

          {{ item.name }}
        </span>

      </div>

    </div>


    <!-- 分割线 -->
    <div class="divider"></div>


    <!-- 社交链接 -->
    <div class="social-links">

      <a
          v-for="item in socialList"
          :key="item.name"
          :href="item.url"
          target="_blank"
          class="social-item"
      >

        <el-icon :size="18">
          <component :is="item.icon"/>
        </el-icon>

        <span>
          {{ item.name }}
        </span>

      </a>

    </div>


    <!-- 底部 -->
    <div class="profile-footer">

      <span>
        BUILD · LEARN · SHARE
      </span>

      <span class="footer-arrow">
        ↗
      </span>

    </div>

  </el-card>
</template>


<script setup lang="ts">

import {
  Monitor,
  Coin,
  DataAnalysis,
  Connection,
  Link
} from "@element-plus/icons-vue";
interface TechItem {
  name: string;
  icon: any;
}
interface SocialItem {
  name: string;
  icon: any;
  url: string;
}
/* ================================
   用户信息
================================ */
const userInfo = {
  name: "Felix",
  avatar: "/image/avatar.jpg"
};
/* ================================
   技术栈
================================ */
const techList: TechItem[] = [
  {
    name: "Go",
    icon: Monitor
  },
  {
    name: "Gin",
    icon: Connection
  },
  {
    name: "MySQL",
    icon: Coin
  },
  {
    name: "Redis",
    icon: DataAnalysis
  },
  {
    name: "Docker",
    icon: Monitor
  }
];
/* ================================
   社交链接
================================ */
const socialList: SocialItem[] = [
  {
    name: "GitHub",
    icon: Link,
    url: "https://github.com/"
  },
  {
    name: "Email",
    icon: Link,
    url: "mailto:example@email.com"
  }
];

</script>
<style scoped lang="scss">

/* =====================================================
   Profile Card
===================================================== */

.profile-card {
  position: relative;
  width: 100%;
  min-height: 420px;
  padding: 0;
  box-sizing: border-box;
  overflow: hidden;
  border: 1px solid #2c2c2c;
  border-radius: 6px;
  background:
      linear-gradient(
          145deg,
          #515151,
          #0d0d0d
      );

  color: #fff;

  transition:
      transform 0.4s ease,
      box-shadow 0.4s ease,
      border-color 0.4s ease;


  &:hover {

    transform:
        translateY(-4px);

    border-color: #444;

    box-shadow:
        0 20px 50px rgba(0, 0, 0, 0.25);

  }

}


/* =====================================================
   Element Plus Card
===================================================== */

:deep(.el-card__body) {

  position: relative;

  min-height: 420px;

  padding: 28px;

  box-sizing: border-box;

  overflow: hidden;

}


/* =====================================================
   Background Grid
===================================================== */

.background-grid {

  position: absolute;

  inset: 0;

  opacity: 0.08;

  background-image:
      linear-gradient(
          rgba(255,255,255,0.08) 1px,
          transparent 1px
      ),
      linear-gradient(
          90deg,
          rgba(255,255,255,0.08) 1px,
          transparent 1px
      );

  background-size:
      30px 30px;

  mask-image:
      linear-gradient(
          to bottom right,
          #000,
          transparent 75%
      );

}


/* =====================================================
   Glow
===================================================== */

.background-glow {

  position: absolute;

  width: 220px;

  height: 220px;

  top: -100px;

  right: -80px;

  border-radius: 50%;

  background:
      radial-gradient(
          circle,
          rgba(255,255,255,0.12),
          transparent 70%
      );

  filter:
      blur(10px);

  animation:
      glowMove 6s ease-in-out infinite alternate;

}


@keyframes glowMove {

  from {

    transform:
        translate(0, 0);

  }

  to {

    transform:
        translate(-35px, 30px);

  }

}


/* =====================================================
   Status
===================================================== */

.profile-status {

  position: relative;

  display: flex;

  align-items: center;

  gap: 8px;

  font-size: 9px;

  letter-spacing: 1.5px;

  color: #777;

}


/* 状态圆点 */

.status-dot {

  width: 6px;

  height: 6px;

  border-radius: 50%;

  background: #fff;

  box-shadow:
      0 0 0 4px rgba(255,255,255,0.06);

  animation:
      statusPulse 2s infinite;

}


@keyframes statusPulse {

  0%,
  100% {

    box-shadow:
        0 0 0 3px
        rgba(255,255,255,0.05);

  }

  50% {

    box-shadow:
        0 0 0 7px
        rgba(255,255,255,0.02);

  }

}


/* =====================================================
   Main
===================================================== */

.profile-main {

  position: relative;

  margin-top: 30px;

}


/* =====================================================
   Avatar
===================================================== */

.avatar-wrapper {

  position: relative;

  width: 82px;

  height: 82px;

}


/* 外圈 */

.avatar-ring {

  position: absolute;

  inset: -5px;

  border: 1px solid rgba(255,255,255,0.18);

  border-radius: 50%;

  animation:
      avatarRotate 8s linear infinite;

}


.avatar-ring::after {

  content: "";

  position: absolute;

  width: 6px;

  height: 6px;

  top: 5px;

  left: 50%;

  border-radius: 50%;

  background: #fff;

}


@keyframes avatarRotate {

  from {

    transform:
        rotate(0deg);

  }

  to {

    transform:
        rotate(360deg);

  }

}


/* Avatar */

.avatar {

  position: relative;

  display: block;

  border: 2px solid #222;

  background: #222;

  transition:
      transform 0.4s ease,
      filter 0.4s ease;

}


.profile-card:hover {

  .avatar {

    transform:
        scale(1.04);

    filter:
        brightness(1.08);

  }

}


/* =====================================================
   Name
===================================================== */

.profile-name {

  margin-top: 18px;

  font-size: 27px;

  font-weight: 600;

  letter-spacing: 1px;

}


/* =====================================================
   Role
===================================================== */

.profile-role {

  margin-top: 5px;

  font-family:
      monospace;

  font-size: 11px;

  letter-spacing: 1.5px;

  color: #777;

}


/* =====================================================
   Description
===================================================== */

.profile-description {

  max-width: 260px;

  margin-top: 16px;

  font-size: 13px;

  line-height: 1.9;

  color: #999;

}


/* =====================================================
   Tech Stack
===================================================== */

.tech-stack {

  position: relative;

  margin-top: 24px;

}


.section-label {

  margin-bottom: 10px;

  font-size: 9px;

  letter-spacing: 2px;

  color: #555;

}


.tech-list {

  display: flex;

  flex-wrap: wrap;

  gap: 7px;

}


.tech-item {

  display: inline-flex;

  align-items: center;

  gap: 5px;

  padding: 5px 9px;

  border: 1px solid #292929;

  border-radius: 3px;

  background: rgba(255,255,255,0.025);

  font-size: 11px;

  color: #aaa;

  cursor: default;

  transition:
      color 0.25s ease,
      border-color 0.25s ease,
      background 0.25s ease,
      transform 0.25s ease;


  &:hover {

    color: #fff;

    border-color: #555;

    background:
        rgba(255,255,255,0.07);

    transform:
        translateY(-2px);

  }

}


/* =====================================================
   Divider
===================================================== */

.divider {
  height: 1px;
  margin:
      22px 0 16px;
  background:
      linear-gradient(
          90deg,
          #292929,
          transparent
      );

}


/* =====================================================
   Social
===================================================== */

.social-links {

  display: flex;

  gap: 18px;

}


.social-item {

  display: inline-flex;

  align-items: center;

  gap: 6px;

  color: #888;

  font-size: 12px;

  text-decoration: none;

  transition:
      color 0.25s ease,
      transform 0.25s ease;


  &:hover {

    color: #fff;

    transform:
        translateY(-2px);

  }
}
/* =====================================================
   Footer
===================================================== */

.profile-footer {

  position: absolute;

  left: 28px;

  right: 28px;

  bottom: 20px;

  display: flex;

  align-items: center;

  justify-content: space-between;

  font-family:
      monospace;

  font-size: 8px;

  letter-spacing: 1.5px;

  color: #444;

}
.footer-arrow {
  font-size: 16px;
  color: #666;
  transition:
      transform 0.3s ease,
      color 0.3s ease;
}

.profile-card:hover {

  .footer-arrow {

    color: #fff;

    transform:
        translate(3px, -3px);

  }

}


/* =====================================================
   Mobile
===================================================== */

@media screen and (max-width: 700px) {

  .profile-card {

    min-height: 400px;

  }


  :deep(.el-card__body) {

    min-height: 400px;

    padding: 24px;

  }


  .profile-footer {

    left: 24px;

    right: 24px;

  }

}

</style>
