<template>
  <div class="carousel">

    <!-- ==========================
         Carousel
    =========================== -->

  <el-carousel
    ref="carouselRef" 
    trigger="click"
    height="515px"
    :interval="6000"
    indicator-position="none"
    arrow="never"
    @change="handleChange"
>

      <!-- Page 01 -->
      <el-carousel-item>
        <HeroPage01/>
      </el-carousel-item>

      <!-- Page 02 -->
      <el-carousel-item>
        <HeroPage02/>
      </el-carousel-item>

      <!-- Page 03 -->
      <el-carousel-item>
        <HeroPage03/>
      </el-carousel-item>

      <!-- Page 04 -->
      <el-carousel-item>
        <HeroPage04/>
      </el-carousel-item>

    </el-carousel>


    <!-- ==========================
         左右切换按钮
    =========================== -->

    <button
        class="carousel-arrow carousel-prev"
        @click="prev"
        aria-label="上一页"
    >
      <span>←</span>
    </button>


    <button
        class="carousel-arrow carousel-next"
        @click="next"
        aria-label="下一页"
    >
      <span>→</span>
    </button>


    <!-- ==========================
         底部控制区域
    =========================== -->

    <div class="carousel-control">

      <!-- 页码 -->

      <div class="page-number">

        <span class="current">
          {{ String(currentIndex + 1).padStart(2, '0') }}
        </span>

        <span class="divider">
          /
        </span>

        <span class="total">
          04
        </span>

      </div>


      <!-- 指示器 -->

      <div class="indicators">

        <button
            v-for="(_, index) in 4"
            :key="index"
            :class="[
              'indicator',
              {
                active: currentIndex === index
              }
            ]"
            @click="goTo(index)"
        >

          <span class="indicator-number">
            {{ String(index + 1).padStart(2, '0') }}
          </span>

          <span class="indicator-line"></span>

        </button>

      </div>


      <!-- 右侧文字 -->

      <div class="carousel-label">
        FELIX BLOG
      </div>

    </div>


    <!-- ==========================
         底部装饰线
    =========================== -->

    <div class="bottom-line"></div>

  </div>
</template>


<script setup lang="ts">

import { ref } from "vue";

import HeroPage01 from "@/views/web/test/HeroPage01.vue";
import HeroPage02 from "@/views/web/test/HeroPage02.vue";
import HeroPage03 from "@/views/web/test/HeroPage03.vue";
import HeroPage04 from "@/views/web/test/HeroPage04.vue";


/**
 * 当前轮播页
 */
const currentIndex = ref(0);


/**
 * Carousel 实例
 */
const carouselRef = ref();


/**
 * 切换事件
 */
const handleChange = (index: number) => {

  currentIndex.value = index;

};


/**
 * 上一页
 */
const prev = () => {

  carouselRef.value?.prev();

};


/**
 * 下一页
 */
const next = () => {

  carouselRef.value?.next();

};


/**
 * 跳转指定页面
 */
const goTo = (index: number) => {

  carouselRef.value?.setActiveItem(index);

};

</script>


<style scoped lang="scss">

/* =====================================================
   Carousel
===================================================== */

.carousel {

  position: relative;

  width: 100%;

  max-width: 1400px;

  margin: 0 auto;

  padding: 0 20px;

  box-sizing: border-box;

}


/* =====================================================
   Element Plus Carousel
===================================================== */

:deep(.el-carousel) {

  width: 100%;

  overflow: hidden;

  border-radius: 0 0 4px 4px;

}


/*
 * 防止 Element Plus 默认背景影响 Hero
 */

:deep(.el-carousel__container) {
  height: 515px;
}

:deep(.el-carousel__item) {
  height: 515px;
}


/*
 * 页面切换动画
 */

:deep(.el-carousel__item) {

  transition:
      opacity 0.8s ease,
      transform 0.8s ease;

}


/* =====================================================
   左右按钮
===================================================== */
.carousel-arrow {
  position: absolute;
  z-index: 20;

  top: 45%;
  width: 42px;
  height: 42px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  border: 1px solid rgba(255, 255, 255, 0.35);
  border-radius: 50%;
  background: rgba(20, 20, 20, 0.18);

  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  color: #fff;

  font-size: 16px;

  cursor: pointer;

  transform: translateY(-50%);

  transition: all 0.3s ease;
}

.carousel-arrow span {

  display: block;

  transition:
      transform 0.3s ease;

}


.carousel-arrow:hover {

  background: rgba(20,20,20,0.75);

  border-color: rgba(255,255,255,0.7);

  box-shadow:
      0 8px 25px rgba(0,0,0,0.15);

}


.carousel-arrow:hover span {

  transform: translateX(2px);

}


.carousel-prev {

  left: 38px;

}


.carousel-next {

  right: 38px;

}


.carousel-prev:hover span {

  transform: translateX(-2px);

}


/* =====================================================
   底部控制栏
===================================================== */

.carousel-control {

  position: absolute;

  z-index: 30;

  left: 40px;

  right: 40px;

  bottom: 20px;

  height: 38px;

  display: flex;

  align-items: center;

  pointer-events: none;

}


/* =====================================================
   页码
===================================================== */

.page-number {

  display: flex;

  align-items: center;

  gap: 7px;

  min-width: 70px;

  color: rgb(196, 24, 24);

  font-size: 9px;

  letter-spacing: 2px;

}


.page-number .current {

  color: rgb(78, 78, 78);

  font-size: 11px;

  font-weight: 600;

}


.page-number .divider {

  color: rgb(78, 78, 78);

}


.page-number .total {
  color: rgb(78, 78, 78);
}


/* =====================================================
   Indicators
===================================================== */

.indicators {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-left: 30px;
}


.indicator {
  display: flex;
  align-items: center;
  gap: 7px;
  padding: 0;
  border: none;
  outline: none;
  background: transparent;
  color: rgb(78, 78, 78);
  cursor: pointer;
  pointer-events: auto;
}
.indicator-number {
  font-size: 10px;
  letter-spacing: 1px;
  transition:
      color 0.25s ease;

}
.indicator-line {
  display: block;
  width: 18px;
  height: 1px;
  background: rgba(71, 71, 71, 0.25);
  transition:
      width 0.35s ease,
      background 0.35s ease;
}


.indicator:hover {

  color: rgb(47, 48, 48);

}


.indicator:hover .indicator-line {

  background: rgb(47, 48, 48);

}


/* Active */

.indicator.active {
  font-weight: bold;
  color: rgb(47, 48, 48);

}


.indicator.active .indicator-line {

  width: 42px;

  background: rgb(47, 48, 48);

}


/* =====================================================
   右侧文字
===================================================== */

.carousel-label {
  margin-left: auto;
  color: rgba(19, 18, 18, 0.45);
  font-weight: bold;
  font-size:10px;
  letter-spacing: 3px;

}


/* =====================================================
   底部装饰线
===================================================== */

.bottom-line {

  position: absolute;
  z-index: 30;
  left: 40px;
  right: 40px;
  bottom: 0;
  height: 1px;
  background:
      linear-gradient(
          to right,
          rgb(0, 0, 0),
          rgba(56, 52, 52, 0.05),
          rgba(86, 86, 86, 0.35)
      );
}


/* =====================================================
   Responsive
===================================================== */

@media screen and (max-width: 900px) {

  .carousel {

    padding: 0;

  }


  .carousel-arrow {

    width: 36px;

    height: 36px;

  }


  .carousel-prev {

    left: 20px;

  }


  .carousel-next {

    right: 20px;

  }


  .carousel-control {

    left: 25px;

    right: 25px;

  }


  .indicators {

    margin-left: 15px;

    gap: 8px;

  }


  .carousel-label {

    display: none;

  }

}


@media screen and (max-width: 600px) {

  .carousel {

    padding: 0;

  }


  :deep(.el-carousel__container) {

    height: 500px;

  }


  :deep(.el-carousel__item) {

    height: 500px;

  }


  .carousel-arrow {

    display: none;

  }


  .carousel-control {

    left: 18px;

    right: 18px;

    bottom: 15px;

  }


  .indicators {

    margin-left: 10px;

  }


  .indicator-number {

    display: none;

  }


  .indicator-line {

    width: 15px;

  }


  .indicator.active .indicator-line {

    width: 30px;

  }

}

</style>