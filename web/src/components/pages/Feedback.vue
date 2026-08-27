<template>
  <el-card class="feedback">
    <el-row class="title">意见反馈</el-row>
    <el-input type="textarea" :rows="4" v-model="feedbackCreateFormData.content" maxlength="100" show-word-limit
              placeholder="请输入反馈建议"></el-input>
    <div class="content">
      <el-text>tip:请登录后再进行反馈!</el-text>
      <div class="button-group">
        <el-button @click="submitForm" type="primary">确定</el-button>
        <el-button @click="feedbackCreateFormData.content=''">取消</el-button>
      </div>
    </div>
    <el-row class="title-sub">反馈列表</el-row>
    <div class="footer">
      <div class="feedback-new" v-for="item in feedbackInfoList">
        <el-row>{{ item.content }}</el-row>
        <el-row class="container">
          <div class="time">{{ item.time }}</div>
        </el-row>
        <div class="reply">
          <el-text v-if="item.reply!==''">回复：{{ item.reply }}</el-text>
        </div>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import {reactive, ref, watch} from "vue";
import {feedbackCreate, type FeedbackCreateRequest, feedbackNew} from "@/api/feedback";

const feedbackCreateFormData = reactive<FeedbackCreateRequest>({
  content: '',
})

interface FeedbackNew {
  content: string;
  reply: string;
  time: string;
}

const feedbackInfoList = ref<FeedbackNew[]>([])

const shouldRefreshFeedbackInfoTable = ref(false)
watch(() => shouldRefreshFeedbackInfoTable.value, (newVal) => {
  if (newVal) {
    getFeedbackNew()
    shouldRefreshFeedbackInfoTable.value = false;
  }
})

const submitForm = async () => {
  const res = await feedbackCreate(feedbackCreateFormData)
  if (res.code === 0) {
    ElMessage.success(res.msg)
    feedbackCreateFormData.content = ''
    shouldRefreshFeedbackInfoTable.value = true
  }
}

const getFeedbackNew = async () => {
  feedbackInfoList.value = []
  const res = await feedbackNew()
  if (res.code === 0) {
    res.data.forEach(value => {
      const date = new Date(value.created_at);
      const info: FeedbackNew = {
        content: value.content,
        reply: value.reply,
        time: date.toLocaleString(),
      }
      feedbackInfoList.value.push(info)
    })
  }
}

getFeedbackNew()
</script>
<style scoped lang="scss">

.feedback {
  margin-bottom: 20px;

  padding: 4px;

  color: #e8e8e8;

  background: #242424;

  border: 1px solid #303030;

  border-radius: 10px;

  box-shadow:
    0 8px 25px rgba(0, 0, 0, 0.18);

  animation: feedbackFadeIn 0.45s ease;


  // ==================================================
  // 标题
  // ==================================================

  .title {
    margin-bottom: 20px;

    color: #f5f5f5;

    font-size: 23px;

    font-weight: 600;

    letter-spacing: 0.5px;

    position: relative;

    padding-left: 12px;


    // 左侧装饰线

    &::before {
      content: "";

      position: absolute;

      left: 0;

      top: 50%;

      width: 3px;

      height: 22px;

      transform: translateY(-50%);

      background: #8c9aa8;

      border-radius: 3px;
    }
  }


  // ==================================================
  // 输入框区域
  // ==================================================

  :deep(.el-textarea__inner) {

    min-height: 110px;

    padding: 14px 16px;

    color: #eeeeee;

    font-size: 14px;

    line-height: 1.7;

    letter-spacing: 0.2px;

    background: #1c1c1c;

    border: 1px solid #383838;

    border-radius: 8px;

    resize: vertical;

    box-shadow: none;

    transition:
      border-color 0.25s ease,
      background-color 0.25s ease,
      box-shadow 0.25s ease;


    &::placeholder {

      color: #777;

    }


    &:hover {

      border-color: #4a4a4a;

    }


    &:focus {

      background: #1e1e1e;

      border-color: #687786;

      box-shadow:
        0 0 0 3px rgba(120, 135, 150, 0.08);

    }
  }


  // 字数统计

  :deep(.el-input__count) {

    color: #666;

    background: transparent;

  }


  // ==================================================
  // 输入框下面的提示 + 按钮
  // ==================================================

  .content {

    display: flex;

    align-items: center;

    margin-top: 16px;

    margin-bottom: 24px;

    padding: 0 2px;


    :deep(.el-text) {

      color: #8f8f8f;

      font-size: 12px;

      letter-spacing: 0.2px;

    }


    .button-group {

      margin-left: auto;

      display: flex;

      gap: 8px;


      .el-button {

        height: 34px;

        padding: 0 18px;

        font-size: 13px;

        border-radius: 6px;

        transition:
          all 0.2s ease;


        // 确定

        &:first-child {

          color: #fff;

          background: #586978;

          border-color: #586978;


          &:hover {

            background: #687b8b;

            border-color: #687b8b;

            transform: translateY(-1px);

            box-shadow:
              0 5px 14px rgba(80, 100, 120, 0.2);

          }


          &:active {

            transform: translateY(0);

          }

        }


        // 取消

        &:last-child {

          color: #aaa;

          background: #2b2b2b;

          border-color: #444;


          &:hover {

            color: #ddd;

            background: #333;

            border-color: #555;

          }

        }
      }
    }
  }


  // ==================================================
  // 「反馈列表」
  // ==================================================

  .title-sub {

    margin-top: 8px;

    margin-bottom: 16px;

    color: #d8d8d8;

    font-size: 16px;

    font-weight: 500;

    letter-spacing: 0.3px;

    padding-bottom: 10px;

    border-bottom: 1px solid #333;
  }


  // ==================================================
  // 反馈列表
  // ==================================================

  .footer {

    .feedback-new {

      position: relative;

      margin-bottom: 12px;

      padding: 16px 18px;

      color: #dedede;

      background: #1d1d1d;

      border: 1px solid #303030;

      border-radius: 8px;

      overflow: hidden;

      transition:
        background-color 0.25s ease,
        border-color 0.25s ease,
        transform 0.25s ease,
        box-shadow 0.25s ease;


      // 左侧细线

      &::before {

        content: "";

        position: absolute;

        left: 0;

        top: 14px;

        bottom: 14px;

        width: 2px;

        background: #596775;

        border-radius: 2px;

        opacity: 0;

        transition:
          opacity 0.25s ease;

      }


      &:hover {

        background: #202020;

        border-color: #3d3d3d;

        transform: translateY(-2px);

        box-shadow:
          0 7px 18px rgba(0, 0, 0, 0.16);


        &::before {

          opacity: 1;

        }

      }


      // =================================================
      // 用户反馈内容
      // =================================================

      > .el-row {

        color: #e2e2e2;

        font-size: 14px;

        line-height: 1.7;

        word-break: break-word;

      }


      // =================================================
      // 时间
      // =================================================

      .container {

        display: flex;

        margin-top: 10px;

        padding-bottom: 9px;

        border-bottom: 1px solid #303030;


        .time {

          margin-left: auto;

          color: #777;

          font-size: 11px;

          letter-spacing: 0.2px;

        }

      }


      // =================================================
      // 回复
      // =================================================

      .reply {

        padding-top: 10px;


        :deep(.el-text) {

          color: #9baab8;

          font-size: 12px;

          line-height: 1.7;

        }


        // 回复前的小装饰

        &::before {

          content: "↳";

          margin-right: 7px;

          color: #657482;

          font-size: 13px;

        }

      }
    }


    // 最后一条不需要额外间距

    .feedback-new:last-child {

      margin-bottom: 0;

    }
  }
}


// ======================================================
// 动画
// ======================================================

@keyframes feedbackFadeIn {

  from {

    opacity: 0;

    transform: translateY(8px);

  }

  to {

    opacity: 1;

    transform: translateY(0);

  }
}


// ======================================================
// 深色主题下的 Element Plus Card
// ======================================================

:deep(.el-card__body) {

  padding: 20px 22px;

  background: #242424;

}


// ======================================================
// 响应式
// ======================================================

@media screen and (max-width: 768px) {

  .feedback {

    .content {

      align-items: flex-start;

      flex-direction: column;

      gap: 12px;


      .button-group {

        width: 100%;

        margin-left: 0;

        justify-content: flex-end;

      }
    }


    .footer {

      .feedback-new {

        padding: 14px;

      }
    }
  }
}

</style>
