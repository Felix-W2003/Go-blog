<template>
  <el-card class="tag-cloud">
    <el-row class="title">
      <span class="title-line"></span>
      <span>标签云</span>
      <span class="title-sub">TAGS</span>
    </el-row>

    <div class="tags-wrapper">
      <el-tag
        v-for="(item, index) in tagCloudArray"
        :key="item.tag"
        :type="item.type"
        size="large"
        effect="plain"
        :style="{ '--delay': `${index * 0.05}s` }"
        @click="handleSearchJumps(item.tag)"
      >
        <span class="tag-name">{{ item.tag }}</span>
        <span class="tag-number">{{ item.number }}</span>
      </el-tag>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { type ArticleTag, articleTags } from "@/api/article";

// 深灰系配色
const tagTypes = ["primary", "info", "success", "warning", "danger"];

interface TagCloudItem {
  tag: string;
  number: number;
  type: string;
}

const tagCloudArray = ref<TagCloudItem[]>([]);

const getTagCloudArray = async () => {
  let tagsArray: ArticleTag[];
  const res = await articleTags();

  if (res.code === 0) {
    tagsArray = res.data;

    for (let i = 0; i < tagsArray.length; i++) {
      const item = tagsArray[i];

      const tagCloud: TagCloudItem = {
        tag: item.tag,
        number: item.number,
        type: tagTypes[i % tagTypes.length]
      };

      tagCloudArray.value.push(tagCloud);
    }
  }
};

getTagCloudArray();

const handleSearchJumps = (tag: string) => {
  window.open("/search?tag=" + tag);
};
</script>

<style scoped lang="scss">
.tag-cloud {
  position: relative;
  margin-bottom: 20px;
  overflow: hidden;

  // 深灰背景
  background: linear-gradient(
    145deg,
    #242424 0%,
    #1d1d1d 55%,
    #191919 100%
  );

  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;

  // 去掉 Element Plus 默认阴影
  box-shadow:
    0 8px 30px rgba(0, 0, 0, 0.25),
    inset 0 1px 0 rgba(255, 255, 255, 0.03);

  transition:
    transform 0.35s ease,
    box-shadow 0.35s ease,
    border-color 0.35s ease;

  // 背景装饰光
  &::before {
    content: "";
    position: absolute;
    width: 180px;
    height: 180px;
    top: -100px;
    right: -70px;

    background: radial-gradient(
      circle,
      rgba(255, 255, 255, 0.07),
      transparent 70%
    );

    pointer-events: none;
  }

  &:hover {
    transform: translateY(-3px);

    border-color: rgba(255, 255, 255, 0.13);

    box-shadow:
      0 14px 40px rgba(0, 0, 0, 0.35),
      inset 0 1px 0 rgba(255, 255, 255, 0.04);
  }

  // =========================
  // 标题
  // =========================
  .title {
    position: relative;
    display: flex;
    align-items: center;

    margin-bottom: 22px;

    font-size: 20px;
    font-weight: 600;
    letter-spacing: 1px;

    color: #eeeeee;

    .title-line {
      width: 4px;
      height: 20px;

      margin-right: 10px;

      border-radius: 4px;

      background: linear-gradient(
        to bottom,
        #ffffff,
        #777777
      );

      box-shadow: 0 0 10px rgba(255, 255, 255, 0.15);
    }

    .title-sub {
      margin-left: 10px;

      font-size: 10px;
      font-weight: 400;
      letter-spacing: 2px;

      color: #666666;
    }
  }

  // =========================
  // 标签区域
  // =========================
  .tags-wrapper {
    display: flex;
    flex-wrap: wrap;
    gap: 11px;

    // 标签进入动画
    .el-tag {
      opacity: 0;

      animation: tagFadeIn 0.45s ease forwards;
      animation-delay: var(--delay);

      cursor: pointer;

      height: 36px;
      padding: 0 13px;

      border-radius: 7px;

      font-size: 13px;

      background: rgba(255, 255, 255, 0.035);

      border: 1px solid rgba(255, 255, 255, 0.10);

      color: #bdbdbd;

      backdrop-filter: blur(6px);

      transition:
        transform 0.25s ease,
        color 0.25s ease,
        background 0.25s ease,
        border-color 0.25s ease,
        box-shadow 0.25s ease;

      // 不同标签稍微有一点层次
      &:nth-child(3n + 1) {
        background: rgba(255, 255, 255, 0.045);
      }

      &:nth-child(3n + 2) {
        background: rgba(255, 255, 255, 0.03);
      }

      &:nth-child(3n) {
        background: rgba(255, 255, 255, 0.06);
      }

      &:hover {
        transform: translateY(-4px) scale(1.03);

        color: #ffffff;

        background: rgba(255, 255, 255, 0.09);

        border-color: rgba(255, 255, 255, 0.28);

        box-shadow:
          0 6px 18px rgba(0, 0, 0, 0.3),
          0 0 12px rgba(255, 255, 255, 0.05);
      }

      &:active {
        transform: translateY(-1px) scale(0.98);
      }

      // 标签文字
      .tag-name {
        transition: color 0.25s ease;
      }

      // 数量
      .tag-number {
        display: inline-flex;
        align-items: center;
        justify-content: center;

        min-width: 20px;
        height: 20px;

        margin-left: 7px;
        padding: 0 5px;

        border-radius: 10px;

        font-size: 10px;

        color: #777777;

        background: rgba(255, 255, 255, 0.06);

        transition:
          color 0.25s ease,
          background 0.25s ease;
      }

      &:hover .tag-number {
        color: #eeeeee;

        background: rgba(255, 255, 255, 0.12);
      }
    }
  }
}

// =========================
// 标签进入动画
// =========================
@keyframes tagFadeIn {
  0% {
    opacity: 0;
    transform: translateY(8px);
  }

  100% {
    opacity: 1;
    transform: translateY(0);
  }
}

// =========================
// Element Plus 深色适配
// =========================
:deep(.el-card__body) {
  padding: 22px;
}

:deep(.el-tag--primary),
:deep(.el-tag--success),
:deep(.el-tag--info),
:deep(.el-tag--warning),
:deep(.el-tag--danger) {
  --el-tag-bg-color: transparent;
  --el-tag-border-color: transparent;
  --el-tag-text-color: inherit;
}
</style>