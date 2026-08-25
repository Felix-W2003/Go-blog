<template>
  <el-card class="calendar">
    <div class="title">
      <div class="title-icon">
        <CalendarIcon />
      </div>
      <span>今日日历</span>
      <span class="live-dot"></span>
    </div>

    <div class="calendar-content">
      <div class="calendar-item">
        <Clock />
        <span class="label">时间</span>
        <span class="value">{{ calendarInfo.date }} {{ currentTime }}</span>
      </div>

      <div class="calendar-item">
        <Collection />
        <span class="label">农历</span>
        <span class="value">{{ calendarInfo.lunar_date }}</span>
      </div>

      <div class="calendar-item">
        <DataAnalysis />
        <span class="label">干支</span>
        <span class="value">{{ calendarInfo.ganzhi }}</span>
      </div>

      <div class="calendar-item">
        <Star />
        <span class="label">星座</span>
        <span class="value">{{ calendarInfo.zodiac }}</span>
      </div>

      <div class="calendar-item">
        <Timer />
        <span class="label">天次</span>
        <span class="value">{{ calendarInfo.day_of_year }}</span>
      </div>

      <div class="calendar-item">
        <Sunny />
        <span class="label">节气</span>
        <span class="value">{{ calendarInfo.solar_term }}</span>
      </div>

      <div class="calendar-item">
        <CircleCheck />
        <span class="label">宜项</span>
        <span class="value">{{ calendarInfo.auspicious }}</span>
      </div>

      <div class="calendar-item">
        <CircleClose />
        <span class="label">禁忌</span>
        <span class="value">{{ calendarInfo.inauspicious }}</span>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { onUnmounted, ref } from "vue";
import {
  Calendar as CalendarIcon,
  Clock,
  Collection,
  DataAnalysis,
  Star,
  Timer,
  Sunny,
  CircleCheck,
  CircleClose
} from "@element-plus/icons-vue";
import { websiteCalendar, type WebsiteCalendarResponse } from "@/api/website";

const calendarInfo = ref<WebsiteCalendarResponse>({
  date: '',
  lunar_date: '',
  ganzhi: '',
  zodiac: '',
  day_of_year: '',
  solar_term: '',
  auspicious: '',
  inauspicious: '',
});

let timerId: number | null = null;
const currentTime = ref('');

function updateCurrentTime() {
  currentTime.value = new Date().toLocaleTimeString();
}

function initializeTimer() {
  updateCurrentTime();
  timerId = setInterval(updateCurrentTime, 1000);
}

onUnmounted(() => {
  if (timerId) {
    clearInterval(timerId);
  }
});

initializeTimer();

const getCalendarInfo = async () => {
  const res = await websiteCalendar();
  if (res.code == 0) {
    calendarInfo.value = res.data;
  }
};

getCalendarInfo();
</script>

<style scoped lang="scss">
.calendar {
  position: relative;
  margin-bottom: 20px;
  overflow: hidden;
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 2px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.04);
  transition: all 0.3s ease;

  :deep(.el-card__body) {
    padding: 22px 24px;
  }

  &:hover {
    border-color: #d5d5d5;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.07);
    transform: translateY(-2px);
  }

  .title {
    display: flex;
    align-items: center;
    height: 36px;
    padding-bottom: 14px;
    margin-bottom: 8px;
    color: #222;
    font-size: 20px;
    font-weight: 600;
    letter-spacing: 1px;
    border-bottom: 1px solid #eeeeee;

    .title-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 30px;
      height: 30px;
      margin-right: 10px;
      color: #444;
      background: #f2f2f2;
      border-radius: 50%;
      transition: all 0.3s ease;

      :deep(svg) {
        width: 16px;
        height: 16px;
      }
    }

    .live-dot {
      width: 5px;
      height: 5px;
      margin-left: 9px;
      background: #888;
      border-radius: 50%;
      animation: breathe 2s ease-in-out infinite;
    }
  }

  .calendar-content {
    display: flex;
    flex-direction: column;
  }

  .calendar-item {
    display: flex;
    align-items: center;
    min-height: 32px;
    color: #666;
    font-size: 13px;
    line-height: 1.6;
    transition: all 0.25s ease;

    :deep(svg) {
      width: 15px;
      height: 15px;
      margin-right: 10px;
      color: #888;
      transition: all 0.25s ease;
    }

    .label {
      width: 42px;
      color: #999;
      flex-shrink: 0;
    }

    .value {
      color: #555;
      word-break: break-all;
    }

    &:hover {
      padding-left: 4px;

      :deep(svg) {
        color: #222;
        transform: scale(1.1);
      }

      .label {
        color: #666;
      }

      .value {
        color: #222;
      }
    }
  }
}

@keyframes breathe {
  0%,
  100% {
    opacity: 0.35;
    transform: scale(0.8);
  }

  50% {
    opacity: 1;
    transform: scale(1.2);
  }
}
</style>