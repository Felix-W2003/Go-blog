<template>
  <el-card class="daily-news" shadow="never">

    <!-- ==============================
         Header
    =============================== -->

    <div class="news-header">

      <div class="title-wrapper">
        <div class="title">
          每日新闻
        </div>

        <div class="title-line"></div>
      </div>

      <span class="title-en">
        DAILY NEWS
      </span>

    </div>


    <!-- ==============================
         News Tabs
    =============================== -->

    <el-tabs
        v-model="activeTab"
        class="news-tabs"
        @tab-click="handleNewsTabClick"
    >

      <el-tab-pane
          v-for="item in newsTypeList"
          :key="item.name"
          :name="item.name"
      >

        <!-- Tab Label -->

        <template #label>

          <div class="tab-label">

            <div class="tab-icon-wrapper">

              <el-image
                  class="tab-icon"
                  :src="item.src"
                  alt=""
              />

            </div>

            <span>
              {{ item.label }}
            </span>

          </div>

        </template>


        <!-- ==========================
             News Table
        =========================== -->

        <div class="table-wrapper">

          <el-table
              :data="newsTableData"
              class="news-table"
              :row-style="{ height: '64px' }"
              :show-header="true"
          >

            <!-- 序号 -->

            <el-table-column
                prop="index"
                label="序号"
                width="70"
            >

              <template #default="scope">

                <span
                    class="news-index"
                    :class="{
                      top: scope.row.index <= 3
                    }"
                >
                  {{ scope.row.index }}
                </span>

              </template>

            </el-table-column>


            <!-- 标题 -->

            <el-table-column label="标题">

              <template #default="scope">

                <div
                    class="news-title"
                    @click="handleNewsTableClick(scope.row)"
                >

                  <span class="title-text">
                    {{ scope.row.title }}
                  </span>

                  <span class="arrow">
                    →
                  </span>

                </div>

              </template>

            </el-table-column>


            <!-- 热度 -->

            <el-table-column
                prop="popularity"
                label="热度"
                width="130"
            >

              <template #default="scope">

                <span class="popularity">
                  {{ scope.row.popularity }}
                </span>

              </template>

            </el-table-column>

          </el-table>

        </div>

      </el-tab-pane>

    </el-tabs>


    <!-- ==============================
         Footer
    =============================== -->

    <div class="news-footer">

      <span class="tip-text">
        仅显示前 7 条热榜数据
      </span>

      <el-link
          class="more-link"
          href="/news"
          :underline="false"
      >

        查看完整热榜

        <span class="more-arrow">
          →
        </span>

      </el-link>

    </div>

  </el-card>
</template>

<script setup lang="ts">

import { ref } from "vue";

import {
  type HotItem,
  websiteNews,
  type WebsiteNewsRequest
} from "@/api/website";

import type { TabsPaneContext } from "element-plus";


/* =====================================================
   当前 Tab
===================================================== */

const activeTab = ref("baidu");


/* =====================================================
   新闻数据
===================================================== */

const newsTableData = ref<HotItem[]>([]);


/* =====================================================
   新闻类型
===================================================== */

interface NewsTypeItem {

  name: string;

  label: string;

  src: string;

}


const newsTypeList: NewsTypeItem[] = [

  {
    name: "baidu",
    label: "百度热搜",
    src: "/image/baidu.png"
  },

  {
    name: "zhihu",
    label: "知乎热榜",
    src: "/image/zhihu.png"
  },

  {
    name: "kuaishou",
    label: "快手热榜",
    src: "/image/kuaishou.png"
  },

  {
    name: "toutiao",
    label: "头条热榜",
    src: "/image/toutiao.png"
  }

];


/* =====================================================
   新闻缓存
===================================================== */

const newsMap = new Map<string, HotItem[]>();


/* =====================================================
   Tab 切换
===================================================== */

const handleNewsTabClick = (
    tab: TabsPaneContext,
    _: Event
) => {

  const source = tab.paneName as string;

  getNewsTableData(source);

};


/* =====================================================
   获取新闻
===================================================== */

const getNewsTableData = async (
    source: string
) => {

  /*
   * 如果缓存中已经有数据
   * 直接使用
   */

  if (newsMap.has(source)) {

    newsTableData.value =
        newsMap.get(source) || [];

    return;

  }


  const newsRequest: WebsiteNewsRequest = {

    source

  };


  const res = await websiteNews(
      newsRequest
  );


  if (res.code === 0) {

    const data =
        res.data.hot_list
            .slice(0, 7)
            .map((item, index) => ({
              ...item,
              index: index + 1
            }));


    newsMap.set(
        source,
        data
    );


    newsTableData.value = data;

  }

};


/* =====================================================
   点击新闻
===================================================== */

const handleNewsTableClick = (
    item: HotItem
) => {

  window.open(
      item.url,
      "_blank"
  );

};


/* =====================================================
   初始化
===================================================== */

getNewsTableData("baidu");

</script>

<style scoped lang="scss">

/* =====================================================
   Daily News
===================================================== */

.daily-news {

  position: relative;
  width: 100%;
  margin-bottom: 24px;
  padding: 28px 30px 22px;
  box-sizing: border-box;
  border: 1px solid #e8e8e8;
  border-radius: 20px;
  background: #ffffff;
  overflow: hidden;
  transition:
      transform 0.35s ease,
      box-shadow 0.35s ease,
      border-color 0.35s ease;


  /* subtle hover */

  &:hover {

    transform: translateY(-2px);

    border-color: #dedede;

    box-shadow:
        0 12px 35px rgba(0, 0, 0, 0.06);

  }

}


/* =====================================================
   Header
===================================================== */

.news-header {
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
  font-size:40px;
  font-weight: 400;
  line-height: 1.3;
  letter-spacing: 7px;
  color: #464646;

}


/* 左侧小黑线 */

.title-line {
  width: 34px;
  height: 3px;
  margin-top: 9px;
  background:  rgb(33, 33, 33),;
  transition:
      width 0.75s ease;

}


/*
 * 鼠标进入整个卡片
 * 标题线展开
 */

.daily-news:hover {

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
   Tabs
===================================================== */

.news-tabs {

  width: 100%;

}


/* Element Plus */

:deep(.el-tabs__header) {

  margin: 0;

  border-bottom: 1px solid #eeeeee;

}


:deep(.el-tabs__nav-wrap::after) {

  height: 1px;

  background-color: #eeeeee;

}


/* Tab */

:deep(.el-tabs__item) {

  position: relative;

  height: 52px;

  padding: 0 20px;

  color: #777;

  font-size: 14px;

  transition:
      color 0.3s ease;

}


:deep(.el-tabs__item:hover) {

  color: #111;

}


:deep(.el-tabs__item.is-active) {

  color: #111;

  font-weight: 500;

}


/*
 * active 下划线
 */

:deep(.el-tabs__active-bar) {

  height: 2px;

  border-radius: 2px;

  background: #111;

  transition:
      all 0.35s ease;

}


/* =====================================================
   Tab Label
===================================================== */

.tab-label {

  display: flex;

  align-items: center;

  gap: 8px;

  height: 100%;

  white-space: nowrap;

}


.tab-icon-wrapper {

  display: flex;

  align-items: center;

  justify-content: center;

  width: 24px;

  height: 24px;

}


.tab-icon {

  width: 22px;

  height: 22px;

  object-fit: contain;

  transition:
      transform 0.35s ease;

}


/*
 * Hover 图标轻微放大
 */

:deep(.el-tabs__item:hover) {

  .tab-icon {

    transform:
        scale(1.12)
        rotate(-3deg);

  }

}


/*
 * Active 图标
 */

:deep(.el-tabs__item.is-active) {

  .tab-icon {

    transform:
        scale(1.08);

  }

}


/* =====================================================
   Table Wrapper
===================================================== */

.table-wrapper {

  width: 100%;
  margin-top: 5px;

}


/* =====================================================
   Table
===================================================== */

.news-table {

  width: 100%;

  background: transparent;

}


/*
 * Header
 */

:deep(.el-table__header-wrapper) {

  background: #fafafa;

}


:deep(.el-table th.el-table__cell) {

  height: 44px;

  padding: 0;

  border-bottom: 1px solid #eeeeee;

  background: #fafafa;

  color: #0b0b0b;

  font-size: 12px;

  font-weight: 400;

}


/*
 * Body
 */

:deep(.el-table td.el-table__cell) {

  padding: 0;

  border-bottom: 1px solid #f0f0f0;

  background: transparent;

}


/*
 * 去掉 Element 默认边框
 */

:deep(.el-table::before) {

  display: none;

}


:deep(.el-table__inner-wrapper::before) {

  display: none;

}


/* =====================================================
   Row
===================================================== */

:deep(.el-table__body tr) {

  transition:
      background-color 0.25s ease;

}


:deep(.el-table__body tr:hover > td.el-table__cell) {

  background: #fafafadf !important;

}


/* =====================================================
   Index
===================================================== */

.news-index {

  display: inline-flex;

  align-items: center;

  justify-content: center;

  width: 25px;

  height: 25px;

  font-size: 13px;

  color: #999;

  transition:
      all 0.3s ease;

}


.news-index.top {

  color: #222;

  font-weight: 600;

}


/*
 * Row hover
 */

:deep(.el-table__body tr:hover) {

  .news-index {

    transform:
        translateX(3px);

    color: #111;

  }

}


/* =====================================================
   News Title
===================================================== */

.news-title {

  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-right: 20px;
  cursor: pointer;
  color: #444;
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


/*
 * Arrow
 */

.arrow {

  flex-shrink: 0;

  margin-left: 15px;

  opacity: 0;

  transform:
      translateX(-8px);

  color: #111;

  transition:
      all 0.3s ease;

}


/*
 * Title hover
 */

.news-title:hover {

  .title-text {

    color: #111;

    transform:
        translateX(4px);

  }

  .arrow {

    opacity: 1;

    transform:
        translateX(0);

  }

}


/* =====================================================
   Popularity
===================================================== */

.popularity {

  font-size: 13px;

  font-variant-numeric: tabular-nums;

  color: #777;

  transition:
      color 0.25s ease,
      transform 0.25s ease;

}


:deep(.el-table__body tr:hover) {

  .popularity {

    color: #111;

    transform:
        translateX(-3px);

  }

}


/* =====================================================
   Footer
===================================================== */

.news-footer {

  display: flex;

  align-items: center;

  justify-content: space-between;

  margin-top: 18px;

  padding-top: 16px;

  border-top: 1px solid #eeeeee;

}


.tip-text {

  font-size: 12px;

  color: #aaa;

}


/* =====================================================
   More Link
===================================================== */

.more-link {

  display: inline-flex;

  align-items: center;

  gap: 8px;

  font-size: 12px;

  color: #666;

  transition:
      color 0.25s ease;

}


.more-arrow {

  transition:
      transform 0.3s ease;

}


.more-link:hover {

  color: #111;

}


.more-link:hover {

  .more-arrow {

    transform:
        translateX(5px);

  }

}


/* =====================================================
   Animation
===================================================== */

/*
 * 表格内容进入动画
 */

:deep(.el-table__body tr) {

  animation:
      newsRowIn 0.4s ease both;

}


@keyframes newsRowIn {

  from {

    opacity: 0;

    transform:
        translateY(8px);

  }

  to {

    opacity: 1;

    transform:
        translateY(0);

  }

}


/*
 * 每一行稍微错开
 */

:deep(.el-table__body tr:nth-child(1)) {

  animation-delay: 0.03s;

}


:deep(.el-table__body tr:nth-child(2)) {

  animation-delay: 0.06s;

}


:deep(.el-table__body tr:nth-child(3)) {

  animation-delay: 0.09s;

}


:deep(.el-table__body tr:nth-child(4)) {

  animation-delay: 0.12s;

}


:deep(.el-table__body tr:nth-child(5)) {

  animation-delay: 0.15s;

}


:deep(.el-table__body tr:nth-child(6)) {

  animation-delay: 0.18s;

}


:deep(.el-table__body tr:nth-child(7)) {

  animation-delay: 0.21s;

}


/* =====================================================
   Mobile
===================================================== */
@media screen and (max-width: 700px) {
  .daily-news {

    padding:
        22px 18px 18px;

  }

  .title {

    font-size: 20px;

  }
  .title-en {

    display: none;

  }


  :deep(.el-tabs__item) {

    padding:
        0 10px;

    font-size: 13px;

  }


  .tab-label {

    gap: 5px;

  }


  .tab-icon {

    width: 19px;

    height: 19px;

  }


  :deep(.el-table th:nth-child(1)),
  :deep(.el-table td:nth-child(1)) {

    width: 45px !important;

  }

  :deep(.el-table th:nth-child(3)),
  :deep(.el-table td:nth-child(3)) {

    width: 80px !important;

  }

  .news-title {

    padding-right: 5px;

  }

  .arrow {

    display: none;

  }
  .news-footer {

    align-items: flex-start;

    flex-direction: column;

    gap: 10px;

  }

}

</style>