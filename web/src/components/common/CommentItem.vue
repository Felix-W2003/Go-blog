<template>
  <div class="comment-item">
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
        <comment-item :comments="item.children" />
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

defineProps<{
  comments: Comment[];
}>();

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

  color: #dedede;

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

    border: 1px solid rgba(255, 255, 255, 0.08);

    border-radius: 9px;

    background: rgba(255, 255, 255, 0.025);

    box-shadow:
      inset 0 1px 0 rgba(255, 255, 255, 0.025);

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

      background: rgba(255, 255, 255, 0.045);

      border-color: rgba(255, 255, 255, 0.15);

      box-shadow:
        0 8px 24px rgba(0, 0, 0, 0.2),
        inset 0 1px 0 rgba(255, 255, 255, 0.035);
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

        border: 1px solid rgba(255, 255, 255, 0.12);

        box-shadow:
          0 3px 10px rgba(0, 0, 0, 0.25);

        transition:
          transform 0.3s ease,
          border-color 0.3s ease;
      }

      &:hover {
        .user-avatar {
          transform: scale(1.05);

          border-color: rgba(255, 255, 255, 0.28);
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

        color: #f1f1f1;

        transition: color 0.25s ease;
      }

      // ------------------------------------
      // 时间
      // ------------------------------------

      .time {
        margin-left: auto;

        font-size: 11px;

        color: #8e8e8e;

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
        color: #dcdcdc;

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

      border-top: 1px solid rgba(255, 255, 255, 0.055);

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

      border-top: 1px solid rgba(255, 255, 255, 0.06);

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

          border: 1px solid rgba(255, 255, 255, 0.08);

          transition:
            transform 0.25s ease,
            border-color 0.25s ease;

          &:hover {
            transform: scale(1.08);

            border-color: rgba(255, 255, 255, 0.2);
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

      background: linear-gradient(
        to bottom,
        rgba(255, 255, 255, 0.12),
        rgba(255, 255, 255, 0.02)
      );
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
  color: #eeeeee;

  background: rgba(255, 255, 255, 0.045);

  border: 1px solid rgba(255, 255, 255, 0.10);

  border-radius: 7px;

  box-shadow: none;

  resize: vertical;

  transition:
    border-color 0.25s ease,
    background 0.25s ease,
    box-shadow 0.25s ease;

  &::placeholder {
    color: #777777;
  }

  &:hover {
    border-color: rgba(255, 255, 255, 0.18);
  }

  &:focus {
    border-color: rgba(255, 255, 255, 0.28);

    background: rgba(255, 255, 255, 0.055);

    box-shadow:
      0 0 0 2px rgba(255, 255, 255, 0.035);
  }
}

// ========================================
// Markdown Preview 深色适配
// ========================================

:deep(.md-editor-preview) {
  color: #dcdcdc;

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
    color: #f1f1f1;

    border-bottom-color: rgba(255, 255, 255, 0.08);
  }

  // 普通文字
  p {
    color: #dcdcdc;
  }

  // 加粗
  strong {
    color: #f0f0f0;
  }

  // 链接
  a {
    color: #c8c8c8;

    transition: color 0.2s ease;

    &:hover {
      color: #ffffff;
    }
  }

  // 引用
  blockquote {
    color: #aaaaaa;

    background: rgba(255, 255, 255, 0.035);

    border-left-color: #777777;
  }

  // 行内代码
  code {
    color: #d6d6d6;

    background: rgba(255, 255, 255, 0.07);

    border-radius: 4px;
  }

  // 代码块
  pre {
    background: #151515;

    border: 1px solid rgba(255, 255, 255, 0.07);

    border-radius: 7px;
  }

  // 分割线
  hr {
    border-color: rgba(255, 255, 255, 0.08);
  }

  // 列表
  li {
    color: #dcdcdc;
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