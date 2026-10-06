<template>
  <div
    class="comment-item"
    :class="'theme-' + theme"
  >
    <div
      v-for="item in comments"
      :key="item.id"
      class="comment-wrapper"
    >
      <div class="item-card">

        <!-- 用户信息 -->
        <div class="title">
          <el-popover width="280">
            <template #reference>
              <el-avatar
                class="user-avatar"
                :src="item.user.avatar"
              />
            </template>

            <template #default>
              <user-card
                :uuid="''"
                :user-card-info="{
                  uuid: item.user.uuid,
                  username: item.user.username,
                  avatar: item.user.avatar,
                  address: item.user.address,
                  signature: item.user.signature
                }"
              />
            </template>
          </el-popover>

          <div class="name">
            {{ item.user.username }}
          </div>

          <div class="time">
            {{ getTime(item.created_at) }}
          </div>
        </div>

        <!-- 评论内容 -->
        <div class="content-wrapper">
          <MdPreview
            class="content"
            :modelValue="item.content"
          />
        </div>

        <!-- 操作按钮 -->
        <div class="footer">
          <div class="button-group">

            <el-button
              v-if="replyFlag === item.id"
              type="primary"
              @click="submitReply(item); content = ''"
            >
              确定
            </el-button>

            <el-button
              v-if="replyFlag === item.id"
              @click="content = ''; replyFlag = 0"
            >
              取消
            </el-button>

            <el-button
              v-if="replyFlag !== item.id"
              type="primary"
              @click="replyFlag = item.id"
            >
              回复
            </el-button>

            <el-button
              v-if="
                (item.user_uuid === userStore.state.userInfo.uuid ||
                userStore.isAdmin) &&
                replyFlag !== item.id
              "
              type="danger"
              @click="handleDelete(item.id)"
            >
              删除
            </el-button>

          </div>
        </div>

        <!-- 回复 -->
        <div
          v-if="replyFlag === item.id"
          class="reply"
        >
          <el-input
            class="comment-input"
            v-model="content"
            :autosize="{ minRows: 4, maxRows: 8 }"
            type="textarea"
            placeholder="在这里输入您的回复..."
            maxlength="320"
          />

          <div class="comment-tool">
            <el-popover
              width="502"
              trigger="click"
            >
              <template #reference>
                <el-avatar
                  class="emoji-avatar"
                  src="/emoji/xiaochun_emoji_01.png"
                />
              </template>

              <template #default>
                <el-image
                  v-for="number in numbers"
                  :key="number"
                  :src="'/emoji/xiaochun_emoji_' + number + '.png'"
                  @click="
                    content =
                      content +
                      '![](/emoji/xiaochun_emoji_' +
                      number +
                      '.png)'
                  "
                />
              </template>
            </el-popover>
          </div>
        </div>

      </div>

      <!-- 子评论 -->
      <div
        v-if="item.children && item.children.length"
        class="item-children"
      >
        <comment-item
          :comments="item.children"
          :theme="theme"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  type Comment,
  commentCreate,
  type CommentCreateRequest,
  commentDelete,
  type CommentDeleteRequest
} from "@/api/comment";

import UserCard from "@/components/widgets/UserCard.vue";
import { useUserStore } from "@/stores/user";
import { MdPreview } from "md-editor-v3";
import { ref } from "vue";
import { useLayoutStore } from "@/stores/layout";

/**
 * theme:
 *   light（默认）—— 文章页、后台「我的评论」等浅色容器
 *   dark         —— 首页侧边栏「最新评论」深色卡片
 */
withDefaults(
  defineProps<{
    comments: Comment[];
    theme?: "light" | "dark";
  }>(),
  {
    theme: "light",
  }
);

const userStore = useUserStore();

const getTime = (date: Date): string => {
  const time = new Date(date);
  return time.toLocaleString();
};

const replyFlag = ref(0);
const content = ref("");

const numbers = ref(
  Array.from(
    { length: 50 },
    (_, i) => (i + 1).toString().padStart(2, "0")
  )
);

const layoutStore = useLayoutStore();

const submitReply = async (item: Comment) => {
  const commentCreateRequest: CommentCreateRequest = {
    article_id: item.article_id,
    p_id: item.id,
    content: content.value,
  };

  const res = await commentCreate(commentCreateRequest);

  if (res.code === 0) {
    ElMessage.success(res.msg);

    layoutStore.state.shouldRefreshCommentList = true;

    replyFlag.value = 0;
  }
};

const handleDelete = async (id: number) => {
  const ids: number[] = [];
  ids.push(id);

  const commentDeleteRequest: CommentDeleteRequest = {
    ids: ids
  };

  const res = await commentDelete(commentDeleteRequest);

  if (res.code === 0) {
    ElMessage.success(res.msg);

    layoutStore.state.shouldRefreshCommentList = true;
  }
};
</script>

<style scoped lang="scss">

// ========================================
// 评论整体
// ========================================

.comment-item {
  width: 100%;

  // ======================================
  // 主题变量
  // 所有配色统一走变量，下面的样式只引用变量，
  // 这样同一套结构可以同时适配浅色页面与深色卡片。
  // ======================================

  // 浅色（默认）：文章页 / 后台「我的评论」
  &.theme-light {
    --ci-text: #333333;
    --ci-text-strong: #111111;
    --ci-text-muted: #767676;

    --ci-card-bg: #fbfbfb;
    --ci-card-bg-hover: #f5f5f5;
    --ci-card-border: #ebebeb;
    --ci-card-border-hover: #d8d8d8;
    --ci-card-shadow: 0 1px 2px rgba(0, 0, 0, 0.03);
    --ci-card-shadow-hover: 0 6px 18px rgba(0, 0, 0, 0.07);

    --ci-divider: #f0f0f0;

    --ci-avatar-border: #e8e8e8;
    --ci-avatar-border-hover: #d0d0d0;
    --ci-avatar-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);

    --ci-thread-line: linear-gradient(to bottom, #e4e4e4, #f8f8f8);

    --ci-input-text: #333333;
    --ci-input-bg: #ffffff;
    --ci-input-bg-focus: #ffffff;
    --ci-input-border: #e2e2e2;
    --ci-input-border-hover: #cfcfcf;
    --ci-input-border-focus: #b0b0b0;
    --ci-input-focus-ring: rgba(0, 0, 0, 0.045);
    --ci-placeholder: #b5b5b5;

    --ci-link: #4a6fa5;
    --ci-link-hover: #2f5286;

    --ci-quote-text: #666666;
    --ci-quote-bg: #f7f7f7;
    --ci-quote-border: #d4d4d4;

    --ci-code-text: #444444;
    --ci-code-bg: #f2f2f2;

    --ci-pre-bg: #f7f7f7;
    --ci-pre-border: #ececec;
  }

  // 深色：首页侧边栏「最新评论」深色卡片
  &.theme-dark {
    --ci-text: #dcdcdc;
    --ci-text-strong: #f1f1f1;
    --ci-text-muted: #8e8e8e;

    --ci-card-bg: rgba(255, 255, 255, 0.025);
    --ci-card-bg-hover: rgba(255, 255, 255, 0.045);
    --ci-card-border: rgba(255, 255, 255, 0.08);
    --ci-card-border-hover: rgba(255, 255, 255, 0.15);
    --ci-card-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.025);
    --ci-card-shadow-hover:
      0 8px 24px rgba(0, 0, 0, 0.2),
      inset 0 1px 0 rgba(255, 255, 255, 0.035);

    --ci-divider: rgba(255, 255, 255, 0.06);

    --ci-avatar-border: rgba(255, 255, 255, 0.12);
    --ci-avatar-border-hover: rgba(255, 255, 255, 0.28);
    --ci-avatar-shadow: 0 3px 10px rgba(0, 0, 0, 0.25);

    --ci-thread-line: linear-gradient(
      to bottom,
      rgba(255, 255, 255, 0.12),
      rgba(255, 255, 255, 0.02)
    );

    --ci-input-text: #eeeeee;
    --ci-input-bg: rgba(255, 255, 255, 0.045);
    --ci-input-bg-focus: rgba(255, 255, 255, 0.055);
    --ci-input-border: rgba(255, 255, 255, 0.1);
    --ci-input-border-hover: rgba(255, 255, 255, 0.18);
    --ci-input-border-focus: rgba(255, 255, 255, 0.28);
    --ci-input-focus-ring: rgba(255, 255, 255, 0.035);
    --ci-placeholder: #777777;

    --ci-link: #c8c8c8;
    --ci-link-hover: #ffffff;

    --ci-quote-text: #aaaaaa;
    --ci-quote-bg: rgba(255, 255, 255, 0.035);
    --ci-quote-border: #777777;

    --ci-code-text: #d6d6d6;
    --ci-code-bg: rgba(255, 255, 255, 0.07);

    --ci-pre-bg: #151515;
    --ci-pre-border: rgba(255, 255, 255, 0.07);
  }

  color: var(--ci-text);

  // ======================================
  // 单条评论
  // ======================================

  .comment-wrapper {
    animation: commentFadeIn 0.45s ease both;
  }

  .item-card {
    position: relative;

    padding: 14px;

    margin-bottom: 10px;

    border: 1px solid var(--ci-card-border);

    border-radius: 9px;

    background: var(--ci-card-bg);

    box-shadow: var(--ci-card-shadow);

    transition:
      transform 0.3s ease,
      background 0.3s ease,
      border-color 0.3s ease,
      box-shadow 0.3s ease;

    // --------------------------------------
    // 评论 Hover
    // --------------------------------------

    &:hover {
      transform: translateY(-2px);

      background: var(--ci-card-bg-hover);

      border-color: var(--ci-card-border-hover);

      box-shadow: var(--ci-card-shadow-hover);
    }

    // ======================================
    // 用户信息
    // ======================================

    .title {
      display: flex;
      align-items: center;

      min-height: 40px;

      .user-avatar {
        width: 38px;
        height: 38px;

        margin-right: 10px;

        border: 1px solid var(--ci-avatar-border);

        box-shadow: var(--ci-avatar-shadow);

        transition:
          transform 0.3s ease,
          border-color 0.3s ease;
      }

      &:hover {
        .user-avatar {
          transform: scale(1.05);

          border-color: var(--ci-avatar-border-hover);
        }
      }

      // ------------------------------------
      // 用户名
      // ------------------------------------

      .name {
        display: flex;
        align-items: center;

        height: 38px;

        font-size: 14px;

        font-weight: 500;

        letter-spacing: 0.3px;

        color: var(--ci-text-strong);

        transition: color 0.25s ease;
      }

      // ------------------------------------
      // 时间
      // ------------------------------------

      .time {
        margin-left: auto;

        font-size: 11px;

        color: var(--ci-text-muted);

        white-space: nowrap;

        transition: color 0.25s ease;
      }
    }

    // ======================================
    // 评论正文
    // ======================================

    .content-wrapper {
      margin-top: 10px;

      padding: 0 3px;

      .content {
        color: var(--ci-text);

        font-size: 14px;

        line-height: 1.8;
      }
    }

    // ======================================
    // Footer
    // ======================================

    .footer {
      display: flex;

      margin-top: 8px;

      padding-top: 8px;

      border-top: 1px solid var(--ci-divider);

      .button-group {
        display: flex;

        gap: 5px;

        margin-left: auto;
      }
    }

    // ======================================
    // 回复区域
    // ======================================

    .reply {
      margin-top: 12px;

      padding-top: 12px;

      border-top: 1px solid var(--ci-divider);

      .comment-input {
        margin-top: 2px;
      }

      .comment-tool {
        display: flex;

        align-items: center;

        margin-top: 8px;

        .emoji-avatar {
          width: 32px;
          height: 32px;

          cursor: pointer;

          background-color: transparent;

          border: 1px solid var(--ci-card-border);

          transition:
            transform 0.25s ease,
            border-color 0.25s ease;

          &:hover {
            transform: scale(1.08);

            border-color: var(--ci-card-border-hover);
          }
        }
      }
    }
  }

  // ========================================
  // 子评论
  // ========================================

  .item-children {
    position: relative;

    margin-left: 18px;

    padding-left: 18px;

    // 左侧层级线
    &::before {
      content: "";

      position: absolute;

      left: 0;
      top: 0;
      bottom: 18px;

      width: 1px;

      background: var(--ci-thread-line);
    }
  }
}

// ========================================
// Element Plus 按钮
// ========================================

:deep(.el-button) {
  height: 28px;

  padding: 0 11px;

  border-radius: 6px;

  font-size: 12px;

  transition:
    transform 0.2s ease,
    background 0.2s ease,
    border-color 0.2s ease;

  &:hover {
    transform: translateY(-1px);
  }
}

// ========================================
// 回复输入框
// ========================================

:deep(.comment-input .el-textarea__inner) {
  color: var(--ci-input-text);

  background: var(--ci-input-bg);

  border: 1px solid var(--ci-input-border);

  border-radius: 7px;

  box-shadow: none;

  resize: vertical;

  transition:
    border-color 0.25s ease,
    background 0.25s ease,
    box-shadow 0.25s ease;

  &::placeholder {
    color: var(--ci-placeholder);
  }

  &:hover {
    border-color: var(--ci-input-border-hover);
  }

  &:focus {
    border-color: var(--ci-input-border-focus);

    background: var(--ci-input-bg-focus);

    box-shadow:
      0 0 0 2px var(--ci-input-focus-ring);
  }
}

// ========================================
// Markdown Preview 主题适配
// ========================================

:deep(.md-editor-preview) {
  color: var(--ci-text);

  background: transparent;

  font-size: 14px;

  line-height: 1.8;

  // 标题
  h1,
  h2,
  h3,
  h4,
  h5,
  h6 {
    color: var(--ci-text-strong);

    border-bottom-color: var(--ci-card-border);
  }

  // 普通文字
  p {
    color: var(--ci-text);
  }

  // 加粗
  strong {
    color: var(--ci-text-strong);
  }

  // 链接
  a {
    color: var(--ci-link);

    transition: color 0.2s ease;

    &:hover {
      color: var(--ci-link-hover);
    }
  }

  // 引用
  blockquote {
    color: var(--ci-quote-text);

    background: var(--ci-quote-bg);

    border-left-color: var(--ci-quote-border);
  }

  // 行内代码
  code {
    color: var(--ci-code-text);

    background: var(--ci-code-bg);

    border-radius: 4px;
  }

  // 代码块
  pre {
    background: var(--ci-pre-bg);

    border: 1px solid var(--ci-pre-border);

    border-radius: 7px;
  }

  // 分割线
  hr {
    border-color: var(--ci-card-border);
  }

  // 列表
  li {
    color: var(--ci-text);
  }
}

// ========================================
// Emoji
// ========================================

:deep(.el-popover.el-popper) {
  .el-image {
    width: 50px;

    height: 50px;

    margin: 2px;

    cursor: pointer;

    border-radius: 6px;

    transition:
      transform 0.2s ease;

    &:hover {
      transform: scale(1.08);
    }
  }
}

// ========================================
// 评论进入动画
// ========================================

@keyframes commentFadeIn {
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
