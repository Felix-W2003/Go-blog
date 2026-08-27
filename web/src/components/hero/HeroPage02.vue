<template>
  <section class="hero-page">

    <!-- =========================================
         背景
    ========================================== -->

    <div class="grid-bg"></div>

    <div class="glow glow-top"></div>
    <div class="glow glow-bottom"></div>


    <!-- =========================================
         背景装饰
    ========================================== -->

    <div class="background-label">
      REDIS / DATA STRUCTURE
    </div>

    <div class="vertical-line vertical-line-left"></div>
    <div class="vertical-line vertical-line-right"></div>


    <!-- =========================================
         主内容
    ========================================== -->

    <div class="hero-container">


      <!-- =======================================
           左侧文字
      ======================================== -->

      <div class="hero-content">

        <div class="hero-label">

          <span class="label-number">
            02
          </span>

          <span class="label-line"></span>

          DATA

        </div>


        <h1>
          深入理解
          <br>

          <strong>Redis</strong>
        </h1>


        <p class="description">
          从 String 到 Hash、List、Set、ZSet，
          理解 Redis 的核心数据结构与底层实现，
          探索缓存、计数、排行榜等真实应用场景。
        </p>


        <!-- =================================
             Button
        ================================== -->

        <div class="actions">

          <button class="primary-btn">

            探索 Redis

            <span>→</span>

          </button>


          <button class="text-btn">
            查看文章
          </button>

        </div>


        <!-- =================================
             数据类型
        ================================== -->

        <div class="data-types">

          <div class="type-item active">
            <span class="type-dot"></span>
            STRING
          </div>

          <div class="type-item">
            <span class="type-dot"></span>
            HASH
          </div>

          <div class="type-item">
            <span class="type-dot"></span>
            LIST
          </div>

          <div class="type-item">
            <span class="type-dot"></span>
            SET
          </div>

          <div class="type-item">
            <span class="type-dot"></span>
            ZSET
          </div>

        </div>

      </div>


      <!-- =======================================
           右侧 Redis 可视化
      ======================================== -->

      <div class="data-visual">


        <!-- =====================================
             Redis 主节点
        ====================================== -->

        <div class="redis-core">

          <div class="core-ring ring-one"></div>
          <div class="core-ring ring-two"></div>
          <div class="core-ring ring-three"></div>


          <div class="core-content">

            <span class="core-label">
              DATA STORE
            </span>

            <span class="core-title">
              REDIS
            </span>

            <span class="core-status">
              ● ONLINE
            </span>

          </div>

        </div>


        <!-- =====================================
             数据节点
        ====================================== -->

        <div class="data-node node-string">

          <div class="node-header">

            <span class="node-icon">
              S
            </span>

            <span>
              STRING
            </span>

          </div>

          <div class="node-value">

            <span class="key">
              user:name
            </span>

            <span class="value">
              Felix
            </span>

          </div>

        </div>


        <div class="data-node node-hash">

          <div class="node-header">

            <span class="node-icon">
              H
            </span>

            <span>
              HASH
            </span>

          </div>

          <div class="node-value">

            <span class="key">
              user:1001
            </span>

            <span class="value">
              5 fields
            </span>

          </div>

        </div>


        <div class="data-node node-list">

          <div class="node-header">

            <span class="node-icon">
              L
            </span>

            <span>
              LIST
            </span>

          </div>

          <div class="node-value">

            <span class="key">
              task:queue
            </span>

            <span class="value">
              128 items
            </span>

          </div>

        </div>


        <div class="data-node node-set">

          <div class="node-header">

            <span class="node-icon">
              S
            </span>

            <span>
              SET
            </span>

          </div>

          <div class="node-value">

            <span class="key">
              tags
            </span>

            <span class="value">
              24 members
            </span>

          </div>

        </div>


        <div class="data-node node-zset">

          <div class="node-header">

            <span class="node-icon">
              Z
            </span>

            <span>
              ZSET
            </span>

          </div>

          <div class="node-value">

            <span class="key">
              ranking
            </span>

            <span class="value">
              1024 scores
            </span>

          </div>

        </div>


        <!-- =====================================
             中央数据流
        ====================================== -->

        <div class="data-flow flow-one"></div>
        <div class="data-flow flow-two"></div>
        <div class="data-flow flow-three"></div>
        <div class="data-flow flow-four"></div>
        <div class="data-flow flow-five"></div>


        <!-- 数据包 -->

        <div class="data-packet packet-one"></div>
        <div class="data-packet packet-two"></div>
        <div class="data-packet packet-three"></div>

      </div>

    </div>


    <!-- =========================================
         底部信息
    ========================================== -->

    <div class="bottom-info">

      <span>
        KEY → VALUE
      </span>

      <span>
        IN-MEMORY
      </span>

      <span>
        HIGH PERFORMANCE
      </span>

      <span>
        02 / 04
      </span>

    </div>


    <!-- =========================================
         右上角编号
    ========================================== -->

    <div class="page-index">

      <span class="current">
        02
      </span>

      <i></i>

      <span>
        04
      </span>

    </div>

  </section>
</template>


<script setup lang="ts">

</script>


<style scoped lang="scss">

/* =====================================================
   Hero Page
===================================================== */

.hero-page {

  position: relative;

  width: 100%;

  height: 430px;

  overflow: hidden;

  color: #ffffff;

  background:
    linear-gradient(
      115deg,
      #0b0b0b 0%,
      #111111 48%,
      #181818 100%
    );

}


/* =====================================================
   Grid Background
===================================================== */

.grid-bg {

  position: absolute;

  inset: 0;

  opacity: 0.14;

  background-image:

    linear-gradient(
      rgba(255,255,255,0.06) 1px,
      transparent 1px
    ),

    linear-gradient(
      90deg,
      rgba(255,255,255,0.06) 1px,
      transparent 1px
    );

  background-size:
    42px 42px;

  mask-image:
    radial-gradient(
      circle at 70% 50%,
      black 0%,
      transparent 70%
    );

}


/* =====================================================
   Glow
===================================================== */

.glow {

  position: absolute;

  width: 420px;

  height: 420px;

  border-radius: 50%;

  filter: blur(130px);

  opacity: 0.08;

  pointer-events: none;

}


.glow-top {

  top: -300px;

  right: 100px;

  background: #ffffff;

}


.glow-bottom {

  bottom: -300px;

  right: 350px;

  background: #888888;

}


/* =====================================================
   Background Label
===================================================== */

.background-label {

  position: absolute;

  left: 50%;

  top: 50%;

  transform:
    translate(-10%, -50%)
    rotate(-90deg);

  color: rgba(255,255,255,0.025);

  font-size: 72px;

  font-weight: 700;

  letter-spacing: 15px;

  white-space: nowrap;

  pointer-events: none;

}


/* =====================================================
   Vertical Lines
===================================================== */

.vertical-line {

  position: absolute;

  top: 0;

  width: 1px;

  height: 100%;

  background:
    linear-gradient(
      to bottom,
      transparent,
      rgba(255,255,255,0.08),
      transparent
    );

}


.vertical-line-left {

  left: 7%;

}


.vertical-line-right {

  right: 7%;

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

  justify-content: space-between;

}


/* =====================================================
   Left Content
===================================================== */

.hero-content {

  width: 46%;

}


/* =====================================================
   Label
===================================================== */

.hero-label {

  display: flex;

  align-items: center;

  gap: 10px;

  margin-bottom: 18px;

  color: #777777;

  font-size: 10px;

  letter-spacing: 3px;

}


.label-number {

  color: #ffffff;

  font-size: 11px;

}


.label-line {

  width: 28px;

  height: 1px;

  background: #666666;

}


/* =====================================================
   Title
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
   Description
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
   Actions
===================================================== */

.actions {

  display: flex;

  align-items: center;

  gap: 22px;

  margin-top: 25px;

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

  transition:
    transform 0.25s ease,
    background 0.25s ease;

}


.primary-btn span {

  margin-left: 10px;

}


.primary-btn:hover {

  transform:
    translateY(-2px);

  background: #eeeeee;

}


.text-btn {

  padding: 0;

  border: none;

  background: transparent;

  color: #777777;

  font-size: 12px;

  cursor: pointer;

  transition:
    color 0.2s ease;

}


.text-btn:hover {

  color: #ffffff;

}


/* =====================================================
   Data Types
===================================================== */

.data-types {

  display: flex;

  align-items: center;

  gap: 18px;

  margin-top: 35px;

}


.type-item {

  display: flex;

  align-items: center;

  gap: 5px;

  color: #555555;

  font-size: 8px;

  letter-spacing: 1.5px;

  transition:
    color 0.25s ease;

}


.type-item.active {

  color: #eeeeee;

}


.type-dot {

  width: 4px;

  height: 4px;

  border-radius: 50%;

  background: #444444;

}


.type-item.active .type-dot {

  background: #ffffff;

}


/* =====================================================
   Data Visual
===================================================== */

.data-visual {

  position: relative;

  width: 530px;

  height: 360px;

}


/* =====================================================
   Redis Core
===================================================== */

.redis-core {

  position: absolute;

  left: 50%;

  top: 50%;

  width: 120px;

  height: 120px;

  transform:
    translate(-50%, -50%);

  border:
    1px solid rgba(255,255,255,0.25);

  border-radius: 50%;

  background:
    rgba(20,20,20,0.95);

  box-shadow:
    0 0 50px rgba(255,255,255,0.04);

}


/* =====================================================
   Core Rings
===================================================== */

.core-ring {

  position: absolute;

  border:
    1px solid rgba(255,255,255,0.08);

  border-radius: 50%;

  pointer-events: none;

}


.ring-one {

  inset: -12px;

  animation:
    rotate-ring 16s linear infinite;

}


.ring-two {

  inset: -25px;

  border-style: dashed;

  animation:
    rotate-ring-reverse 22s linear infinite;

}


.ring-three {

  inset: -42px;

  border-color:
    rgba(255,255,255,0.04);

}


/* =====================================================
   Core Content
===================================================== */

.core-content {

  position: relative;

  width: 100%;

  height: 100%;

  display: flex;

  flex-direction: column;

  align-items: center;

  justify-content: center;

}


.core-label {

  color: #555555;

  font-size: 7px;

  letter-spacing: 2px;

}


.core-title {

  margin-top: 5px;

  color: #ffffff;

  font-size: 24px;

  font-weight: 700;

  letter-spacing: 3px;

}


.core-status {

  margin-top: 8px;

  color: #777777;

  font-size: 7px;

  letter-spacing: 1px;

}


/* =====================================================
   Data Nodes
===================================================== */

.data-node {

  position: absolute;

  z-index: 5;

  width: 145px;

  padding: 10px 12px;

  border:
    1px solid rgba(255,255,255,0.09);

  border-radius: 4px;

  background:
    rgba(20,20,20,0.88);

  box-shadow:
    0 12px 35px rgba(0,0,0,0.3);

  backdrop-filter:
    blur(12px);

  transition:
    transform 0.3s ease,
    border-color 0.3s ease;

}


.data-node:hover {

  transform:
    translateY(-4px);

  border-color:
    rgba(255,255,255,0.25);

}


/* =====================================================
   Node Position
===================================================== */

.node-string {

  left: 5px;

  top: 35px;

}


.node-hash {

  right: 5px;

  top: 35px;

}


.node-list {

  left: 5px;

  bottom: 32px;

}


.node-set {

  right: 5px;

  bottom: 32px;

}


.node-zset {

  left: 50%;

  bottom: -18px;

  transform:
    translateX(-50%);

}


.node-zset:hover {

  transform:
    translateX(-50%)
    translateY(-4px);

}


/* =====================================================
   Node Header
===================================================== */

.node-header {

  display: flex;

  align-items: center;

  gap: 7px;

  color: #999999;

  font-size: 8px;

  letter-spacing: 1.5px;

}


.node-icon {

  width: 18px;

  height: 18px;

  display: flex;

  align-items: center;

  justify-content: center;

  border:
    1px solid #555555;

  border-radius: 3px;

  color: #ffffff;

  font-size: 8px;

  font-weight: 600;

}


/* =====================================================
   Node Value
===================================================== */

.node-value {

  display: flex;

  align-items: center;

  justify-content: space-between;

  margin-top: 8px;

  padding-top: 7px;

  border-top:
    1px solid rgba(255,255,255,0.06);

}


.node-value .key {

  color: #555555;

  font-size: 7px;

}


.node-value .value {

  color: #bbbbbb;

  font-size: 7px;

}


/* =====================================================
   Data Flow
===================================================== */

.data-flow {

  position: absolute;

  z-index: 1;

  height: 1px;

  background:
    linear-gradient(
      to right,
      transparent,
      rgba(255,255,255,0.18),
      transparent
    );

  transform-origin: left center;

}


.flow-one {

  width: 150px;

  left: 82px;

  top: 105px;

  transform:
    rotate(23deg);

}


.flow-two {

  width: 150px;

  right: 82px;

  top: 105px;

  transform:
    rotate(157deg);

}


.flow-three {

  width: 150px;

  left: 82px;

  bottom: 102px;

  transform:
    rotate(-23deg);

}


.flow-four {

  width: 150px;

  right: 82px;

  bottom: 102px;

  transform:
    rotate(-157deg);

}


.flow-five {

  width: 85px;

  left: 50%;

  bottom: 48px;

  transform:
    translateX(-50%)
    rotate(90deg);

}


/* =====================================================
   Data Packets
===================================================== */

.data-packet {

  position: absolute;

  z-index: 2;

  width: 4px;

  height: 4px;

  border-radius: 50%;

  background: #ffffff;

  box-shadow:
    0 0 8px rgba(255,255,255,0.5);

}


.packet-one {

  left: 130px;

  top: 116px;

  animation:
    packet-move-one 3s ease-in-out infinite;

}


.packet-two {

  right: 130px;

  top: 116px;

  animation:
    packet-move-two 3.5s ease-in-out infinite;

}


.packet-three {

  left: 50%;

  bottom: 70px;

  animation:
    packet-move-three 3s ease-in-out infinite;

}


/* =====================================================
   Page Index
===================================================== */

.page-index {

  position: absolute;

  z-index: 6;

  top: 28px;

  right: 35px;

  display: flex;

  align-items: center;

  gap: 8px;

  color: #555555;

  font-size: 9px;

  letter-spacing: 2px;

}


.page-index .current {

  color: #ffffff;

}


.page-index i {

  width: 24px;

  height: 1px;

  background: #555555;

}


/* =====================================================
   Bottom Info
===================================================== */

.bottom-info {

  position: absolute;

  z-index: 6;

  left: 35px;

  right: 35px;

  bottom: 20px;

  display: flex;

  align-items: center;

  gap: 28px;

  color: #444444;

  font-size: 7px;

  letter-spacing: 1.5px;

}


.bottom-info span:last-child {

  margin-left: auto;

  color: #555555;

}


/* =====================================================
   Animations
===================================================== */

@keyframes rotate-ring {

  from {
    transform:
      rotate(0deg);
  }

  to {
    transform:
      rotate(360deg);
  }

}


@keyframes rotate-ring-reverse {

  from {
    transform:
      rotate(360deg);
  }

  to {
    transform:
      rotate(0deg);
  }

}


@keyframes packet-move-one {

  0% {
    opacity: 0;
    transform:
      translate(0, 0);
  }

  30% {
    opacity: 1;
  }

  70% {
    opacity: 1;
  }

  100% {
    opacity: 0;
    transform:
      translate(90px, 35px);
  }

}


@keyframes packet-move-two {

  0% {
    opacity: 0;
    transform:
      translate(0, 0);
  }

  30% {
    opacity: 1;
  }

  70% {
    opacity: 1;
  }

  100% {
    opacity: 0;
    transform:
      translate(-90px, 35px);
  }

}


@keyframes packet-move-three {

  0% {
    opacity: 0;
    transform:
      translateY(30px);
  }

  30% {
    opacity: 1;
  }

  70% {
    opacity: 1;
  }

  100% {
    opacity: 0;
    transform:
      translateY(-20px);
  }

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

  .hero-container {

    width:
      calc(100% - 70px);

  }

  .data-visual {

    width: 460px;

    transform:
      scale(0.9);

  }

}


@media screen and (max-width: 900px) {

  .hero-content {

    width: 55%;

  }

  .data-visual {

    width: 400px;

    transform:
      scale(0.8);

  }

  .description {

    width: 100%;

  }

}


@media screen and (max-width: 760px) {

  .hero-page {

    height: 400px;

  }

  .hero-container {

    width:
      calc(100% - 50px);

  }

  .hero-content {

    width: 100%;

  }

  .hero-content h1 {

    font-size: 40px;

  }

  .data-visual {

    display: none;

  }

  .background-label {

    font-size: 60px;

  }

}


@media screen and (max-width: 500px) {

  .hero-container {

    width:
      calc(100% - 35px);

  }

  .hero-content h1 {

    font-size: 34px;

  }

  .description {

    font-size: 12px;

  }

  .data-types {

    gap: 10px;

  }

  .type-item {

    font-size: 7px;

  }

  .bottom-info {

    gap: 12px;

  }

  .bottom-info span:nth-child(2),
  .bottom-info span:nth-child(3) {

    display: none;

  }

}

</style>