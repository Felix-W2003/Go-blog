<template>
  <el-card class="hot-articles" shadow="never">

    <!-- ==============================
         Header
    =============================== -->

    <div class="hot-header">

      <div class="title-wrapper">
        <div class="title">
          热门文章
        </div>

        <div class="title-line"></div>
      </div>

      <span class="title-en">
        HOT
      </span>

    </div>


    <!-- ==============================
         Hot List
    =============================== -->

    <div class="hot-list">

      <div
          v-if="hotArticles.length === 0"
          class="hot-empty"
      >
        暂无数据
      </div>

      <div
          v-for="(item, index) in hotArticles"
          :key="item._id"
          class="hot-item"
          :style="{ '--delay': `${index * 0.03}s` }"
          @click="handleArticleJump(item._id)"
      >

        <!-- 排名 -->

        <span
            class="hot-index"
            :class="{ top: index < 3 }"
        >
          {{ index + 1 }}
        </span>


        <!-- 标题 -->

        <div class="hot-title">

          <span class="title-text">
            {{ item._source.title }}
          </span>

          <span class="arrow">
            →
          </span>

        </div>


        <!-- 浏览量 -->

        <span class="hot-views">

          <el-icon>
            <component is="View"/>
          </el-icon>

          {{ item._source.views }}

        </span>

      </div>

    </div>

  </el-card>
</template>

<script setup lang="ts">
import { ref } from "vue";
import type { Hit } from "@/api/common";
import { type Article, articleHot } from "@/api/article";

const hotArticles = ref<Hit<Article>[]>([]);

const getHotArticles = async () => {
  const res = await articleHot();

  if (res.code === 0) {
    hotArticles.value = res.data.list;
  }
};

getHotArticles();

const handleArticleJump = (id: string) => {
  window.open("/article/" + id);
};
</script>

<style scoped lang="scss">

/* =====================================================
   Hot Articles
   （卡片外壳与「每日新闻」保持同一套浅色语言）
===================================================== */

.el-card {
  border: 0;
}

.hot-articles {
  position: relative;
  width: 100%;
  margin-bottom: 24px;
  padding: 28px 30px 22px;
  box-sizing: border-box;
  border-radius: 10px;
  background: #fcfcfc;
  overflow: hidden;

  box-shadow:
      2px 1px rgba(0, 0, 0, 0.077);

  transition:
      transform 0.35s ease,
      box-shadow 0.35s ease,
      border-color 0.35s ease;

  &:hover {
    transform: translateY(-2px);

    box-shadow:
        5px 5px rgba(0, 0, 0, 0.173);
  }
}


/* =====================================================
   Header
===================================================== */

.hot-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  margin-bottom: 22px;
}


.title-wrapper {
  display: flex;
  flex-direction: column;
}


.title {
  position: relative;

  font-size: 25px;
  font-weight: 600;
  line-height: 1.3;

  color: #464646;
}


/* 左侧小黑线 */

.title-line {
  width: 34px;
  height: 3px;

  margin-top: 9px;

  background: rgb(33, 33, 33);

  transition:
      width 0.75s ease;
}


/*
 * 鼠标进入整个卡片
 * 标题线展开
 */

.hot-articles:hover {
  .title-line {
    width: 940px;
    height: 3px;

    background:
        linear-gradient(
            to right,
            rgb(255, 255, 255),
            rgb(159, 158, 158),
            rgb(0, 0, 0)
        );
  }
}


/* 英文 */

.title-en {
  position: absolute;
  right: 20px;
  top: 20px;

  font-size: 20px;
  font-weight: 500;
  letter-spacing: 3px;

  color: #b0b0b0;
}


/* =====================================================
   List
===================================================== */

.hot-list {
  width: 100%;
}


.hot-empty {
  padding: 26px 0;

  font-size: 13px;
  text-align: center;

  color: #aaaaaa;
}


/* =====================================================
   Item
===================================================== */

.hot-item {
  display: flex;
  align-items: center;

  padding: 0 10px;
  height: 52px;

  cursor: pointer;

  border-bottom: 1px solid #f0f0f0;

  animation: hotRowIn 0.4s ease both;
  animation-delay: var(--delay);

  transition:
      background-color 0.25s ease;

  &:last-child {
    border-bottom: none;
  }

  &:hover {
    background: #fafafadf;
  }
}


/* =====================================================
   Index（排名）
===================================================== */

.hot-index {
  display: inline-flex;
  align-items: center;
  justify-content: center;

  flex-shrink: 0;

  width: 25px;
  height: 25px;

  font-size: 13px;

  color: #999999;

  transition:
      all 0.3s ease;
}


/* 前三名加粗加深 */

.hot-index.top {
  color: #222222;
  font-weight: 600;
}


.hot-item:hover {
  .hot-index {
    transform: translateX(3px);

    color: #111111;
  }
}


/* =====================================================
   Title
===================================================== */

.hot-title {
  display: flex;
  align-items: center;
  justify-content: space-between;

  flex: 1;
  min-width: 0;

  padding-right: 20px;

  color: #444444;

  transition:
      color 0.25s ease;
}


.title-text {
  overflow: hidden;

  text-overflow: ellipsis;
  white-space: nowrap;

  transition:
      transform 0.3s ease,
      color 0.3s ease;
}


/* 箭头 */

.arrow {
  flex-shrink: 0;

  margin-left: 15px;

  opacity: 0;
  transform: translateX(-8px);

  color: #111111;

  transition:
      all 0.3s ease;
}


.hot-item:hover {
  .title-text {
    color: #111111;

    transform: translateX(4px);
  }

  .arrow {
    opacity: 1;

    transform: translateX(0);
  }
}


/* =====================================================
   Views（浏览量）
===================================================== */

.hot-views {
  display: inline-flex;
  align-items: center;
  gap: 5px;

  flex-shrink: 0;

  font-size: 13px;
  font-variant-numeric: tabular-nums;

  color: #777777;

  transition:
      color 0.25s ease,
      transform 0.25s ease;

  .el-icon {
    font-size: 14px;

    color: #aaaaaa;

    transition: color 0.25s ease;
  }
}


.hot-item:hover {
  .hot-views {
    color: #111111;

    transform: translateX(-3px);

    .el-icon {
      color: #111111;
    }
  }
}


/* =====================================================
   Animation
===================================================== */

@keyframes hotRowIn {

  from {
    opacity: 0;

    transform: translateY(8px);
  }

  to {
    opacity: 1;

    transform: translateY(0);
  }

}


/* =====================================================
   Mobile
===================================================== */

@media screen and (max-width: 700px) {

  .hot-articles {
    padding:
        22px 18px 18px;
  }

  .title {
    font-size: 20px;
  }

  .title-en {
    display: none;
  }

  .hot-item {
    height: 46px;

    padding: 0 4px;
  }

  .hot-title {
    padding-right: 8px;
  }

  .arrow {
    display: none;
  }

}

</style>
