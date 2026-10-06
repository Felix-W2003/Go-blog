<template>
  <el-card class="recent-comments">
    <el-row class="title">
      <span class="title-line"></span>
      <span>最新评论</span>
      <span class="title-sub">RECENT</span>
    </el-row>

    <div class="comments-content">
      <comment-item :comments="comments" theme="dark"/>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import CommentItem from "@/components/common/CommentItem.vue";
import { ref, watch } from "vue";
import { type Comment, commentNew } from "@/api/comment";
import { useLayoutStore } from "@/stores/layout";

const comments = ref<Comment[]>([]);

const getRecentCommentInfo = async () => {
  const res = await commentNew();

  if (res.code === 0) {
    comments.value = res.data;
  }
};

getRecentCommentInfo();

const layoutStore = useLayoutStore();

watch(
  () => layoutStore.state.shouldRefreshCommentList,
  (newVal) => {
    if (newVal) {
      getRecentCommentInfo();
    }
  }
);
</script>
<style scoped lang="scss">
.recent-comments {
  position: relative;

  margin-bottom: 20px;

  overflow: hidden;

  // =========================
  // 深灰高级背景
  // =========================
  background: linear-gradient(
    145deg,
    #242424 0%,
    #1d1d1d 55%,
    #191919 100%
  );

  border: 1px solid rgba(255, 255, 255, 0.08);

  border-radius: 12px;

  box-shadow:
    0 8px 30px rgba(0, 0, 0, 0.25),
    inset 0 1px 0 rgba(255, 255, 255, 0.03);

  transition:
    transform 0.35s ease,
    box-shadow 0.35s ease,
    border-color 0.35s ease;

  // =========================
  // 背景装饰光
  // =========================
  &::before {
    content: "";

    position: absolute;

    width: 200px;
    height: 200px;

    top: -120px;
    right: -80px;

    background: radial-gradient(
      circle,
      rgba(255, 255, 255, 0.07),
      transparent 70%
    );

    pointer-events: none;
  }

  // 左下角微弱光晕
  &::after {
    content: "";

    position: absolute;

    width: 160px;
    height: 160px;

    left: -100px;
    bottom: -100px;

    background: radial-gradient(
      circle,
      rgba(255, 255, 255, 0.035),
      transparent 70%
    );

    pointer-events: none;
  }

  // =========================
  // 卡片 Hover
  // =========================
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
    z-index: 2;

    display: flex;
    align-items: center;

    margin-bottom: 20px;

    font-size: 20px;
    font-weight: 600;

    letter-spacing: 1px;

    color: #eeeeee;

    // 左侧装饰线
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

      box-shadow:
        0 0 10px rgba(255, 255, 255, 0.15);

      transition:
        height 0.3s ease,
        box-shadow 0.3s ease;
    }

    // 副标题
    .title-sub {
      margin-left: 10px;

      font-size: 10px;
      font-weight: 400;

      letter-spacing: 2px;

      color: #666666;

      transition: color 0.3s ease;
    }
  }

  // Hover 标题动画
  &:hover {
    .title-line {
      height: 24px;

      box-shadow:
        0 0 14px rgba(255, 255, 255, 0.25);
    }

    .title-sub {
      color: #888888;
    }
  }

  // =========================
  // 评论区域
  // =========================
  .comments-content {
    position: relative;
    z-index: 2;

    animation: commentsFadeIn 0.5s ease both;
  }
}

// =========================
// 评论淡入
// =========================
@keyframes commentsFadeIn {
  from {
    opacity: 0;

    transform: translateY(8px);
  }

  to {
    opacity: 1;

    transform: translateY(0);
  }
}

// =========================
// Element Plus Card
// =========================
:deep(.el-card__body) {
  padding: 22px;
}

// =========================
// CommentItem 深色适配
// =========================
//
// 如果 CommentItem.vue 内部存在白色背景，
// 这里可以覆盖它的外层容器。
//
:deep(.comment-item) {
  transition:
    background 0.25s ease,
    transform 0.25s ease;
}
</style>