<template>
  <div class="search">
   <WebNavbar />
    <el-container class="main-content">
      <div class="container">
        <el-main>
          <div class="search">
            <el-input v-model="articleSearchRequest.query" placeholder="请输入搜索内容" prefix-icon="Search"
                      maxlength="50"
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
        </el-main>
      </div>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import WebNavbar from "@/components/layout/WebNavbar.vue";
import type {Hit} from "@/api/common";
import {type Article, articleCategory, articleSearch, type ArticleSearchRequest, articleTags} from "@/api/article";
import {computed, nextTick, onMounted, reactive, ref, watch} from "vue";
import {useRoute, useRouter} from "vue-router";

const articleSearchRequest = reactive<ArticleSearchRequest>({
  query: "",
  category: "",
  tag: "",
  sort: "",
  order: "desc",
  page: 1,
  page_size: 10,
})

const route = useRoute()
const router = useRouter()

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
  return articleSearchRequest.order === "desc" ? "#333333" : "#C8C8C8"
})

const upColor = computed(() => {
  return articleSearchRequest.order === "desc" ? "#C8C8C8" : "#333333"
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

onMounted(() => {
  articleSearchRequest.query = route.query.query as string || ""
  articleSearchRequest.category = route.query.category as string || ""
  articleSearchRequest.tag = route.query.tag as string || ""
  articleSearchRequest.sort = route.query.sort as string || ""
  articleSearchRequest.order = route.query.order as string || "desc"
  page.value = Number(route.query.page) || 1
  page_size.value = Number(route.query.page_size) || 10
})


const getArticleSearchTableData = async () => {
  articleSearchRequest.page = page.value;
  articleSearchRequest.page_size = page_size.value;

  const table = await articleSearch(articleSearchRequest)

  if (table.code === 0) {
    articleTableData.value = table.data.list;
    total.value = table.data.total;
  }

  await router.push({
    path: router.currentRoute.value.path,
    query: {
      query: articleSearchRequest.query,
      category: articleSearchRequest.category,
      tag: articleSearchRequest.tag,
      sort: articleSearchRequest.sort,
      order: articleSearchRequest.order,
      page: articleSearchRequest.page,
      page_size: articleSearchRequest.page_size,
    }
  })
}

watch(() => route.query, (newQuery) => {
  articleSearchRequest.query = newQuery.query as string || ""
  articleSearchRequest.category = newQuery.category as string || ""
  articleSearchRequest.tag = newQuery.tag as string || ""
  articleSearchRequest.sort = newQuery.sort as string || ""
  articleSearchRequest.order = newQuery.order as string || "desc"
  articleSearchRequest.page = Number(newQuery.page) || 1
  articleSearchRequest.page_size = Number(newQuery.page_size) || 10
}, {immediate: true})

nextTick(() => {
  getArticleSearchTableData()
})

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

// ======================================================
// 页面整体
// ======================================================

.search {
  min-height: 100vh;

  padding-top: 1px;

  background: #f5f5f5;

  color: #333;

  // 页面进入动画
  animation: pageFadeIn 0.45s ease;
}


// ======================================================
// 主内容
// ======================================================

.search {

  .main-content {

    margin-top: 70px;

    padding: 20px 20px 40px;

    display: flex;

    justify-content: center;


    .container {

      width: 100%;

      max-width: 1400px;


      .el-main {

        padding: 0;

        background: #fff;

        border: 1px solid #ededed;

        border-radius: 12px;

        overflow: hidden;

        box-shadow:
          0 4px 20px rgba(0, 0, 0, 0.04);
      }
    }
  }
}


// ======================================================
// 搜索框区域
//
// 注意：这里使用 .main-content .container .el-main > .search
// 不再直接给所有 .search 设置布局
// ======================================================

.search {

  .main-content {

    .container {

      .el-main {

        > .search {

          min-height: unset;

          padding: 20px 28px 14px;

          display: flex;

          align-items: center;

          background: #fff;


          .el-input {

            width: 360px;

            margin-left: auto;
          }


          .el-button {

            height: 36px;

            margin-left: 8px;

            padding: 0 18px;

            color: #fff;

            font-size: 13px;

            background: #333;

            border: 1px solid #333;

            border-radius: 7px;

            transition: all 0.2s ease;


            &:hover {

              background: #222;

              border-color: #222;

              transform: translateY(-1px);

              box-shadow:
                0 4px 10px rgba(0, 0, 0, 0.12);
            }


            &:active {

              transform: translateY(0);

              box-shadow: none;
            }
          }
        }
      }
    }
  }
}


// ======================================================
// 搜索输入框
// ======================================================

.search {

  :deep(.el-input__wrapper) {

    min-height: 36px;

    padding: 0 12px;

    background: #fafafa;

    border: 1px solid #e3e3e3;

    border-radius: 7px;

    box-shadow: none;

    transition: all 0.25s ease;


    &:hover {

      border-color: #ccc;
    }


    &.is-focus {

      background: #fff;

      border-color: #bbb;

      box-shadow:
        0 0 0 3px rgba(0, 0, 0, 0.035);
    }
  }


  :deep(.el-input__inner) {

    color: #333;

    font-size: 13px;


    &::placeholder {

      color: #aaa;
    }
  }


  :deep(.el-input__prefix-inner) {

    color: #888;
  }
}


// ======================================================
// 分类 / 标签 / 排序
// ======================================================

.search {

  .main-content {

    .container {

      .el-main {

        > .category,
        > .tag,
        > .sort {

          display: flex;

          align-items: flex-start;

          margin: 0;

          padding: 9px 28px;

          border-top: 1px solid #f3f3f3;


          > .el-row {

            flex-shrink: 0;

            width: 42px;

            margin-right: 22px;

            padding-top: 6px;

            color: #777;

            font-size: 13px;

            font-weight: 500;
          }


          > .el-radio-group {

            display: flex;

            flex-wrap: wrap;

            gap: 6px;

            max-width: calc(100% - 70px);
          }
        }
      }
    }
  }
}


// ======================================================
// Radio 按钮
// ======================================================

.search {

  :deep(.el-radio-button) {

    .el-radio-button__inner {

      height: 30px;

      padding: 0 12px;

      display: flex;

      align-items: center;

      justify-content: center;

      color: #777;

      font-size: 12px;

      background: #fafafa;

      border: 1px solid #e6e6e6;

      border-radius: 6px !important;

      box-shadow: none;

      transition: all 0.2s ease;


      &:hover {

        color: #333;

        background: #f3f3f3;

        border-color: #d2d2d2;

        transform: translateY(-1px);
      }
    }


    &.is-active {

      .el-radio-button__inner {

        color: #fff;

        background: #333;

        border-color: #333;

        box-shadow: none;
      }
    }
  }
}


// ======================================================
// 排序按钮
// ======================================================

.search {

  .sort {

    > .el-button {

      width: 34px;

      height: 30px;

      margin-right: 7px;

      padding: 0;

      display: flex;

      align-items: center;

      justify-content: center;

      color: #555;

      background: #fafafa;

      border: 1px solid #e6e6e6;

      border-radius: 6px;

      transition: all 0.2s ease;


      &:hover {

        background: #f2f2f2;

        border-color: #d2d2d2;

        transform: translateY(-1px);
      }


      .el-icon {

        font-size: 14px;
      }
    }
  }
}


// ======================================================
// 文章表格
// ======================================================

.search {

  .el-table {

    margin-top: 12px;

    background: #fff;

    color: #444;

    border-top: 1px solid #eee;


    &::before {

      display: none;
    }


    :deep(.el-table__inner-wrapper)::before {

      display: none;
    }


    :deep(.el-table__row) {

      background: #fff;

      transition: background 0.25s ease;


      &:hover {

        background: #fafafa !important;
      }


      td {

        background: transparent !important;

        border-bottom: 1px solid #eee;
      }
    }


    :deep(.el-table__cell) {

      padding: 14px 10px;

      background: transparent !important;

      border-bottom-color: #eee;
    }


    // ==================================================
    // 文章封面
    // ==================================================

    :deep(.el-image) {

      width: 160px;

      height: 100px;

      overflow: hidden;

      border: 1px solid #eee;

      border-radius: 9px;

      box-shadow:
        0 3px 12px rgba(0, 0, 0, 0.07);

      transition: all 0.35s ease;
    }


    :deep(.el-table__row:hover .el-image) {

      transform: scale(1.035);

      border-color: #ddd;

      box-shadow:
        0 7px 18px rgba(0, 0, 0, 0.1);
    }


    // ==================================================
    // 文章描述
    // ==================================================

    .description {

      height: 120px;

      display: flex;

      flex-direction: column;

      padding: 2px 6px;

      cursor: pointer;


      .title {

        display: block;

        overflow: hidden;

        margin-bottom: 7px;

        color: #222;

        font-size: 20px;

        font-weight: 600;

        line-height: 1.4;

        white-space: nowrap;

        text-overflow: ellipsis;

        transition: all 0.25s ease;
      }


      &:hover {

        .title {

          color: #000;

          transform: translateX(3px);
        }
      }


      .abstract {

        display: -webkit-box;

        overflow: hidden;

        margin-right: auto;

        color: #777;

        font-size: 13px;

        line-height: 1.7;

        -webkit-box-orient: vertical;

        -webkit-line-clamp: 2;
      }


      // =================================================
      // 底部
      // =================================================

      .footer {

        width: 100%;

        margin-top: auto;

        display: flex;

        align-items: center;

        color: #999;

        font-size: 11px;


        // =================================================
        // 标签
        // =================================================

        .tags {

          display: flex;

          flex-wrap: wrap;

          gap: 6px;

          max-width: 55%;

          margin-right: auto;


          .el-tag {

            height: 24px;

            padding: 0 9px;

            font-size: 11px;

            font-weight: 400;

            border-radius: 5px;

            transition: all 0.2s ease;


            &:hover {

              transform: translateY(-1px);

              box-shadow:
                0 3px 8px rgba(0, 0, 0, 0.08);
            }
          }


          // 蓝色
          .tag-color-0 {

            color: #356d9a;

            background: #eaf3fb;

            border-color: #cfe3f3;
          }


          // 绿色
          .tag-color-1 {

            color: #397052;

            background: #eaf6ef;

            border-color: #cce8d7;
          }


          // 紫色
          .tag-color-2 {

            color: #65538a;

            background: #f1ecf9;

            border-color: #ddd2ef;
          }


          // 橙色
          .tag-color-3 {

            color: #9a6534;

            background: #fbf0e4;

            border-color: #eedbc5;
          }


          // 灰蓝
          .tag-color-4 {

            color: #596675;

            background: #eef1f4;

            border-color: #dce2e8;
          }
        }


        // =================================================
        // 发布时间 / 浏览 / 评论 / 点赞
        // =================================================

        .status {

          margin-left: auto;

          display: flex;

          align-items: center;

          justify-content: flex-end;

          gap: 5px;

          color: #999;

          white-space: nowrap;


          .el-icon {

            margin-left: 3px;

            color: #aaa;

            font-size: 13px;
          }
        }
      }
    }
  }
}


// ======================================================
// 分页
// ======================================================

.search {

  .el-pagination {

    margin-top: 14px;

    padding: 10px 0 22px;

    display: flex;

    justify-content: center;

    align-items: center;

    color: #888;


    :deep(.el-pagination__total) {

      color: #888;
    }


    :deep(.btn-prev),
    :deep(.btn-next) {

      color: #777;

      background: #fafafa;

      border: 1px solid #eee;

      border-radius: 6px;

      transition: all 0.2s ease;


      &:hover {

        color: #333;

        background: #f3f3f3;

        border-color: #ddd;
      }
    }


    :deep(.el-pager) {

      li {

        color: #888;

        background: transparent;

        border-radius: 6px;

        transition: all 0.2s ease;


        &:hover {

          color: #333;

          background: #f2f2f2;
        }


        &.is-active {

          color: #fff;

          background: #333;
        }
      }
    }
  }
}


// ======================================================
// 空数据
// ======================================================

.search {

  :deep(.el-table__empty-block) {

    background: #fff;
  }


  :deep(.el-table__empty-text) {

    color: #aaa;
  }
}


// ======================================================
// 页面动画
// ======================================================

@keyframes pageFadeIn {

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