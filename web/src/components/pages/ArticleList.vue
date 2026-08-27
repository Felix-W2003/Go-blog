<template>
  <el-card class="article-list">
    <el-row class="title">文章列表</el-row>
    <div class="search">
      <el-input v-model="articleSearchRequest.query" placeholder="请输入搜索内容" prefix-icon="Search" maxlength="50"
                @change="changeArticleSearchItem"/>
      <el-button @click="changeArticleSearchItem">搜索</el-button>
    </div>

    <div class="category">
      <el-row size="large">类别</el-row>
      <el-radio-group v-model="articleSearchRequest.category" @change="changeArticleSearchItem">
        <el-radio-button label="全部" value=""/>
        <template v-for="item in categoryArr">
          <el-radio-button :label="item" :value="item"/>
        </template>
      </el-radio-group>
    </div>

    <div class="tag">
      <el-row size="large">标签</el-row>
      <el-radio-group v-model="articleSearchRequest.tag" @change="changeArticleSearchItem">
        <el-radio-button label="全部" value=""/>
        <template v-for="item in tagArr">
          <el-radio-button :label="item" :value="item"/>
        </template>
      </el-radio-group>
    </div>

    <div class="sort">
      <el-row size="large">排序</el-row>
      <el-button @click="handleSortClick();changeArticleSearchItem()">
        <el-icon :color="downColor">
          <component is="SortDown"></component>
        </el-icon>
        <el-icon :color="upColor">
          <component is="SortUp"></component>
        </el-icon>
      </el-button>
      <el-radio-group v-model="articleSearchRequest.sort" v-for="item in sortArr"
                      @change="changeArticleSearchItem">
        <el-radio-button :label="item.label" :value="item.value"/>
      </el-radio-group>
    </div>

    <el-table :data="articleTableData" :show-header="false" :row-style="{height: '150px'}">
      <el-table-column label="cover" width="200">
        <template #default="scope:{ row: any, column: any, $index: number }">
          <el-image style="width: 160px; height: 100px" :src="scope.row._source.cover" alt=""/>
        </template>
      </el-table-column>
      <el-table-column label="description">
        <template #default="scope:{ row: Hit<Article>, column: any, $index: number }">
          <div class="description" @click="handleArticleJumps(scope.row._id)">
            <el-row class="title">{{ scope.row._source.title }}</el-row>
            <el-text class="abstract" size="large">{{ scope.row._source.abstract }}</el-text>
            <el-text class="footer">
            <div class="tags">
  <el-tag
      v-for="(item, index) in scope.row._source.tags"
      :key="item"
      :class="'tag-color-' + (index % 5)"
  >
    {{ item }}
  </el-tag>
</div>
              <div class="status">
                发布时间：{{ scope.row._source.created_at }}
                <el-icon>
                  <component is="View"/>
                </el-icon>
                {{ scope.row._source.views }}
                <el-icon>
                  <component is="ChatDotRound"/>
                </el-icon>
                {{ scope.row._source.comments }}
                <el-icon>
                  <component is="Star"/>
                </el-icon>
                {{ scope.row._source.likes }}
              </div>
            </el-text>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
        :current-page="page"
        :page-size="page_size"
        :page-sizes="[10, 30, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
    />
  </el-card>
</template>

<script setup lang="ts">
import type {Hit} from "@/api/common";
import {type Article, articleCategory, articleSearch, type ArticleSearchRequest, articleTags} from "@/api/article";
import {computed, reactive, ref} from "vue";

const articleSearchRequest = reactive<ArticleSearchRequest>({
  query: "",
  category: "",
  tag: "",
  sort: "",
  order: "desc",
  page: 1,
  page_size: 10,
})

const categoryArr = ref<string[]>([])
const tagArr = ref<string[]>([])
const sortArr = [
  {label: "默认", value: ""},
  {label: "时间", value: "time"},
  {label: "评论", value: "comment"},
  {label: "浏览", value: "view"},
  {label: "点赞", value: "like"},
]

const downColor = computed(() => {
  return articleSearchRequest.order === "desc" ? "#333333" : "#C4C4C4"
})

const upColor = computed(() => {
  return articleSearchRequest.order === "desc" ? "#C4C4C4" : "#333333"
})
const handleSortClick = () => {
  articleSearchRequest.order = articleSearchRequest.order === "desc" ? "asc" : "desc"
}

const getArticleCategory = async () => {
  const res = await articleCategory()
  if (res.code === 0) {
    res.data.forEach((item) => {
      categoryArr.value.push(item.category)
    })
  }
}

getArticleCategory()

const getArticleTags = async () => {
  const res = await articleTags()
  if (res.code === 0) {
    res.data.forEach((item) => {
      tagArr.value.push(item.tag)
    })
  }
}

getArticleTags()

const page = ref(1)
const page_size = ref(10)
const total = ref(0)
const articleTableData = ref<Hit<Article>[]>()

const getArticleSearchTableData = async () => {
  articleSearchRequest.page = page.value;
  articleSearchRequest.page_size = page_size.value;

  const table = await articleSearch(articleSearchRequest)

  if (table.code === 0) {
    articleTableData.value = table.data.list;
    total.value = table.data.total;
  }
}

getArticleSearchTableData()

const changeArticleSearchItem = () => {
  getArticleSearchTableData()
}

const handleArticleJumps = (id: string) => {
  window.open("/article/" + id)
}

const handleSizeChange = (val: number) => {
  page_size.value = val
  getArticleSearchTableData()
}

const handleCurrentChange = (val: number) => {
  page.value = val
  getArticleSearchTableData()
}
</script>



<style scoped lang="scss">

// =====================================================
// 文章列表
// =====================================================

.article-list {
  position: relative;

  overflow: hidden;

  margin-bottom: 20px;

  color: #333333;

  background: #ffffff;

  border: 1px solid #eeeeee;

  border-radius: 12px;

  box-shadow:
    0 4px 20px rgba(0, 0, 0, 0.045);

  transition:
    box-shadow 0.3s ease,
    transform 0.3s ease;

  &:hover {
    box-shadow:
      0 8px 28px rgba(0, 0, 0, 0.07);
  }


  // ===================================================
  // 主标题
  // ===================================================

  > .title {
    display: flex;

    align-items: center;

    height: 56px;

    margin: 0;

    padding: 0 22px;

    color: #222222;

    font-size: 21px;

    font-weight: 600;

    letter-spacing: 0.5px;

    border-bottom: 1px solid #f0f0f0;

    background: #ffffff;

    &::before {
      content: "";

      width: 4px;

      height: 20px;

      margin-right: 10px;

      border-radius: 4px;

      background: linear-gradient(
        to bottom,
        #222222,
        #aaaaaa
      );
    }
  }


  // ===================================================
  // 搜索
  // ===================================================

  .search {
    display: flex;

    align-items: center;

    padding: 18px 22px 10px;

    .el-input {
      width: 320px;

      margin-left: auto;
    }

    .el-button {
      height: 34px;

      margin-left: 8px;

      padding: 0 17px;

      color: #ffffff;

      font-size: 13px;

      background: #333333;

      border: 1px solid #333333;

      border-radius: 7px;

      transition:
        background 0.2s ease,
        transform 0.2s ease,
        box-shadow 0.2s ease;

      &:hover {
        background: #222222;

        transform: translateY(-1px);

        box-shadow:
          0 4px 10px rgba(0, 0, 0, 0.12);
      }

      &:active {
        transform: translateY(0);
      }
    }
  }


  // ===================================================
  // 分类 / 标签 / 排序
  // ===================================================

  .category,
  .tag,
  .sort {
    display: flex;

    align-items: flex-start;

    margin: 5px 22px;

    padding: 8px 0;

    .el-row {
      flex-shrink: 0;

      width: 42px;

      margin-right: 20px;

      padding-top: 5px;

      color: #777777;

      font-size: 13px;

      font-weight: 500;
    }

    .el-radio-group {
      display: flex;

      flex-wrap: wrap;

      gap: 6px;

      max-width: calc(100% - 65px);
    }
  }


  // ===================================================
  // Radio Button
  // ===================================================

  :deep(.el-radio-button) {

    .el-radio-button__inner {
      height: 30px;

      padding: 0 11px;

      display: flex;

      align-items: center;

      justify-content: center;

      color: #777777;

      font-size: 12px;

      background: #fafafa;

      border: 1px solid #e8e8e8;

      border-radius: 6px !important;

      box-shadow: none;

      transition:
        color 0.2s ease,
        background 0.2s ease,
        border-color 0.2s ease,
        transform 0.2s ease;

      &:hover {
        color: #333333;

        background: #f5f5f5;

        border-color: #d5d5d5;

        transform: translateY(-1px);
      }
    }

    &.is-active {

      .el-radio-button__inner {
        color: #ffffff;

        background: #333333;

        border-color: #333333;

        box-shadow: none;
      }
    }
  }


  // ===================================================
  // 排序
  // ===================================================

  .sort {

    > .el-button {
      display: flex;

      align-items: center;

      justify-content: center;

      width: 34px;

      height: 30px;

      padding: 0;

      margin-right: 7px;

      color: #555555;

      background: #fafafa;

      border: 1px solid #e8e8e8;

      border-radius: 6px;

      transition:
        background 0.2s ease,
        border-color 0.2s ease,
        transform 0.2s ease;

      &:hover {
        background: #f2f2f2;

        border-color: #d8d8d8;

        transform: translateY(-1px);
      }

      .el-icon {
        font-size: 14px;
      }
    }
  }


  // ===================================================
  // 文章表格
  // ===================================================

  .el-table {

    margin-top: 14px;

    color: #444444;

    background: transparent;

    border-top: 1px solid #f0f0f0;

    &::before {
      display: none;
    }

    :deep(.el-table__inner-wrapper) {
      background: transparent;

      &::before {
        display: none;
      }
    }

    :deep(.el-table__body-wrapper) {
      background: transparent;
    }


    // -------------------------------------------------
    // 行
    // -------------------------------------------------

    :deep(.el-table__row) {

      background: #ffffff;

      transition:
        background 0.25s ease;

      &:hover {
        background: #fafafa !important;
      }

      td {
        background: transparent !important;

        border-bottom: 1px solid #f0f0f0;
      }
    }


    // -------------------------------------------------
    // 单元格
    // -------------------------------------------------

    :deep(.el-table__cell) {

      padding: 14px 10px;

      color: #444444;

      background: transparent !important;

      border-bottom-color: #f0f0f0;
    }


    // =================================================
    // 文章封面
    // =================================================

    :deep(.el-image) {

      width: 160px;

      height: 100px;

      overflow: hidden;

      border-radius: 8px;

      border: 1px solid #eeeeee;

      box-shadow:
        0 3px 12px rgba(0, 0, 0, 0.07);

      transition:
        transform 0.35s ease,
        box-shadow 0.35s ease,
        border-color 0.35s ease;
    }


    :deep(.el-table__row:hover .el-image) {

      transform: scale(1.035);

      border-color: #dddddd;

      box-shadow:
        0 6px 18px rgba(0, 0, 0, 0.11);
    }


    // =================================================
    // 文章描述
    // =================================================

    .description {

      height: 120px;

      display: flex;

      flex-direction: column;

      padding: 2px 6px;

      cursor: pointer;


      // -----------------------------------------------
      // 标题
      // -----------------------------------------------

      .title {

        display: block;

        overflow: hidden;

        margin-bottom: 7px;

        color: #222222;

        font-size: 20px;

        font-weight: 600;

        line-height: 1.4;

        letter-spacing: 0.2px;

        white-space: nowrap;

        text-overflow: ellipsis;

        transition:
          color 0.25s ease,
          transform 0.25s ease;
      }


      &:hover {

        .title {

          color: #000000;

          transform: translateX(3px);
        }
      }


      // -----------------------------------------------
      // 摘要
      // -----------------------------------------------

      .abstract {

        display: -webkit-box;

        overflow: hidden;

        margin-right: auto;

        color: #777777;

        font-size: 13px;

        line-height: 1.75;

        -webkit-box-orient: vertical;

        -webkit-line-clamp: 2;

        transition: color 0.25s ease;
      }


      &:hover {

        .abstract {

          color: #666666;
        }
      }


      // -----------------------------------------------
      // 底部
      // -----------------------------------------------

      .footer {

        display: flex;

        align-items: center;

        width: 100%;

        margin-top: auto;

        color: #999999;

        font-size: 11px;


        // ---------------------------------------------
        // 标签
        // ---------------------------------------------

      .tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;

  margin-right: auto;
  max-width: 55%;

  .el-tag {
    height: 24px;
    padding: 0 9px;

    font-size: 11px;
    font-weight: 400;

    border-radius: 5px;

    transition:
      transform 0.2s ease,
      box-shadow 0.2s ease,
      filter 0.2s ease;

    &:hover {
      transform: translateY(-1px);

      filter: brightness(0.97);

      box-shadow:
        0 3px 8px rgba(0, 0, 0, 0.08);
    }
  }

  // 蓝灰
  .tag-color-0 {
    color: #52708a;

    background: #edf4f8;

    border-color: #d5e4ed;
  }

  // 绿色
  .tag-color-1 {
    color: #587866;

    background: #eef6f0;

    border-color: #d7e8dc;
  }

  // 紫灰
  .tag-color-2 {
    color: #706584;

    background: #f2eff7;

    border-color: #e1dceb;
  }

  // 橙灰
  .tag-color-3 {
    color: #967252;

    background: #faf3eb;

    border-color: #eadccc;
  }

  // 灰蓝
  .tag-color-4 {
    color: #667080;

    background: #f0f2f5;

    border-color: #dfe3e8;
  }
}

        // ---------------------------------------------
        // 文章状态
        // ---------------------------------------------

        .status {

          display: flex;

          align-items: center;

          justify-content: flex-end;

          gap: 5px;

          margin-left: auto;

          color: #999999;

          white-space: nowrap;

          .el-icon {

            margin-left: 3px;

            color: #aaaaaa;

            font-size: 13px;

            transition: color 0.2s ease;
          }
        }
      }
    }
  }


  // ===================================================
  // 分页
  // ===================================================

  .el-pagination {

    display: flex;

    justify-content: center;

    align-items: center;

    margin-top: 14px;

    padding: 8px 0 18px;

    color: #888888;


    :deep(.el-pagination__total) {

      color: #888888;
    }


    :deep(.el-pagination__sizes) {

      .el-select {

        .el-select__wrapper {

          color: #777777;

          background: #fafafa;

          border: 1px solid #e8e8e8;

          box-shadow: none;
        }
      }
    }


    :deep(.btn-prev),
    :deep(.btn-next) {

      color: #777777;

      background: #fafafa;

      border: 1px solid #eeeeee;

      border-radius: 6px;

      transition:
        color 0.2s ease,
        background 0.2s ease,
        border-color 0.2s ease;

      &:hover {

        color: #333333;

        background: #f3f3f3;

        border-color: #dddddd;
      }
    }


    :deep(.el-pager) {

      li {

        color: #888888;

        background: transparent;

        border-radius: 6px;

        transition:
          color 0.2s ease,
          background 0.2s ease;

        &:hover {

          color: #333333;

          background: #f2f2f2;
        }


        &.is-active {

          color: #ffffff;

          background: #333333;
        }
      }
    }
  }
}


// =====================================================
// 搜索框
// =====================================================

:deep(.search .el-input__wrapper) {

  min-height: 34px;

  padding: 0 12px;

  background: #fafafa;

  border: 1px solid #e5e5e5;

  border-radius: 7px;

  box-shadow: none;

  transition:
    background 0.25s ease,
    border-color 0.25s ease,
    box-shadow 0.25s ease;

  &:hover {

    border-color: #d2d2d2;
  }

  &.is-focus {

    background: #ffffff;

    border-color: #bdbdbd;

    box-shadow:
      0 0 0 3px rgba(0, 0, 0, 0.035);
  }
}


:deep(.search .el-input__inner) {

  color: #333333;

  font-size: 13px;

  &::placeholder {

    color: #aaaaaa;
  }
}


:deep(.search .el-input__prefix-inner) {

  color: #888888;
}


// =====================================================
// 标签
// =====================================================

:deep(.el-tag) {

  --el-tag-bg-color: #f7f7f7;

  --el-tag-border-color: #e8e8e8;

  --el-tag-text-color: #777777;
}


// =====================================================
// 空数据
// =====================================================

:deep(.el-table__empty-block) {

  background: #ffffff;
}


:deep(.el-table__empty-text) {

  color: #aaaaaa;
}


// =====================================================
// 文章加载动画
// =====================================================

@keyframes articleFadeIn {

  from {

    opacity: 0;

    transform: translateY(8px);
  }

  to {

    opacity: 1;

    transform: translateY(0);
  }
}

</style>


