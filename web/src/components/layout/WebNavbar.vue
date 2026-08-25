<template>
  <div :class="{'web-navbar': true,show: isShow}">
    <div class="container">
      <logo/>
      <div class="web-menu">
        <el-menu mode="horizontal" :ellipsis="false" :router="true" :default-active="$route.path">
          <template v-for="item in menuList">
            <el-menu-item :index="item.name"><span>{{ item.title }}</span></el-menu-item>
          </template>
        </el-menu>
      </div>
      <auth-popover/>
    </div>
  </div>
</template>



<style scoped lang="scss">
.web-navbar {
  display: flex;
  justify-content: center;
  width: 100%;
  height: 68px;
  position: fixed;
  top: 0;
  left: 0;
  z-index: 100;
  color: #aaa;
  background: #1c1c1c;
  border-bottom: 1px solid #2d2d2d;
  box-shadow: 0 4px 18px rgba(0, 0, 0, 0.08);
  transition: all 0.35s ease;

  .container {
    display: flex;
    align-items: center;
    max-width: 1400px;
    width: 100%;
    height: 68px;

    .logo {
      width: 200px;
      height: 60px;
    }

    .web-menu {
      margin-left: 20px;

      .el-menu {
        height: 68px;
        background: transparent !important;
        border: none !important;

        .el-menu-item {
          position: relative;
          height: 68px;
          line-height: 68px;
          padding: 0 18px;
          border: none !important;
          background: transparent !important;

          color: #aaa !important;
          font-size: 16px;
          font-weight: 400;
          letter-spacing: 0.5px;

          transition: color 0.25s ease, transform 0.25s ease;

          &:hover {
            color: #fff !important;
            background: transparent !important;
            transform: translateY(-1px);
          }

          &::after {
            content: "";
            position: absolute;
            left: 50%;
            bottom: 12px;
            width: 0;
            height: 2px;
            background: #fff;
            border-radius: 2px;
            transform: translateX(-50%);
            transition: width 0.3s ease;
          }

          &:hover::after {
            width: 28px;
          }

          &.is-active {
            color: #fff !important;
            background: transparent !important;

            &::after {
              width: 28px;
              background: #fff;
            }
          }
        }
      }
    }

    .auth-popover {
      margin-left: auto;
      margin-top: auto;
      margin-bottom: auto;
      padding-right: 20px;
    }
  }

  // ============================
  // 滚动后的导航栏
  // ============================

  &.show {
    background: rgba(255, 255, 255, 0.96);
    border-bottom: 1px solid rgba(0, 0, 0, 0.08);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);

    backdrop-filter: blur(18px);
    -webkit-backdrop-filter: blur(18px);

    .container {
      margin-top: 0;
    }

    .web-menu {
      .el-menu {
        .el-menu-item {
          color: #555 !important;
          background: transparent !important;

          &:hover {
            color: #111 !important;
            background: transparent !important;
          }

          &.is-active {
            color: #111 !important;
            background: transparent !important;

            &::after {
              background: #111;
            }
          }
        }
      }
    }
  }
}
</style>
<script setup lang="ts">
import AuthPopover from "@/components/common/AuthPopover.vue";
import Logo from "@/components/widgets/Logo.vue";
import {ref} from "vue";
import {onUnmounted} from "vue";

const isShow = ref(true)

const props = defineProps<{
  noScroll?: boolean
}>()

if (!props.noScroll) {
  isShow.value = false
  window.addEventListener("scroll", scroll)
  scroll()
}

function scroll() {
  let top = document.documentElement.scrollTop
  isShow.value = top >= 100;
}

onUnmounted(() => {
  if (!props.noScroll) {
    window.removeEventListener("scroll", scroll)
  }
})

interface MenuItem {
  title: string;
  name: string;
}

const menuList: MenuItem[] = [
  {
    title: "首页",
    name: "/",
  },
  {
    title: "搜索",
    name: "/search",
  },
  {
    title: "新闻",
    name: "/news",
  },
  {
    title: "友链",
    name: "/friend-link",
  },
  {
    title: "关于",
    name: "/about",
  }
]

</script>
