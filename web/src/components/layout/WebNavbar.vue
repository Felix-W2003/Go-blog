<template>
  <header class="web-navbar" :class="{ show: isScrolled }">
    <div class="container">
      <!-- Logo -->
      <div class="logo-wrapper">
        <div class="blog-logo">
          <span class="logo-main">
            FELIX
          </span>
          <span class="logo-sub">
            BLOG
          </span>
        </div>
      </div>

      <!-- 导航 -->
      <div class="web-menu">
        <el-menu mode="horizontal" :ellipsis="false" :router="true" :default-active="$route.path">
          <el-menu-item v-for="item in menuList" :key="item.name" :index="item.name">
            <span>
              {{ item.title }}
            </span>
          </el-menu-item>
        </el-menu>
      </div>


      <!-- 右侧 -->
      <div class="navbar-right">
        <!-- 搜索 -->
        <div class="search-box">
          <el-input v-model="searchKeyword" placeholder="搜索文章..." clearable @keyup.enter="handleSearch">
            <template #prefix>
              <el-icon>
                <Search />
              </el-icon>
            </template>
          </el-input>
        </div>



        <!-- 登录 -->
        <AuthPopover />


      </div>


    </div>


  </header>
</template>

<script setup lang="ts">
import {
  ref,
  onMounted,
  onUnmounted
} from "vue";
import {
  Search
} from "@element-plus/icons-vue";
import {
  useRouter
} from "vue-router";
import AuthPopover
  from "@/components/common/AuthPopover.vue";
/**
 * router
 */
const router = useRouter();
/**
 * props
 */
const props = defineProps<{
  noScroll?: boolean

}>();
/**
 * 是否滚动
 */
const isScrolled = ref(false);
/**
 * 搜索关键词
 */
const searchKeyword = ref("");
/**
 * 菜单
 */
interface MenuItem {
  title: string;
  name: string;
}
const menuList: MenuItem[] = [
  {
    title: "首页",
    name: "/"
  },
  {
    title: "搜索",
    name: "/search"
  },
  {
    title: "新闻",
    name: "/news"
  },
  {
    title: "友链",
    name: "/friend-link"
  },
  {
    title: "关于",
    name: "/about"
  }
];
/**
 * 搜索
 */
const handleSearch = () => {
  const keyword =
    searchKeyword.value.trim();
    console.log(keyword)
  // 空搜索不处理
  if (!keyword) {
      console.log("this is blanket")
    return;
  }
  console.log("push")
  router.push({
    path: "/search",
    query: {
      query: keyword
    }
  });
};
/**
 * 滚动监听
 */
function handleScroll() {
  const top =
    window.scrollY ||
    document.documentElement.scrollTop ||
    document.body.scrollTop ||
    0;
  isScrolled.value =
    top >= 100;
}
onMounted(() => {
  if (props.noScroll) {
    isScrolled.value = true;
    return;
  }
  handleScroll();
  window.addEventListener(
    "scroll",
    handleScroll,
    {
      passive: true
    }
  );
});
onUnmounted(() => {
  window.removeEventListener(
    "scroll",
    handleScroll
  );
});

</script>
<style scoped lang="scss">
/* =====================================================
   Navbar
===================================================== */
.web-navbar {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 68px;
  z-index: 1000;
  display: flex;
  justify-content: center;
  /**
   * 顶部状态：
   * 深黑色
   */
  background: rgba(18, 18, 18, 0.94);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  box-shadow:
    0 4px 18px rgba(0, 0, 0, 0.08);
  transition:
    background 0.35s ease,
    border-color 0.35s ease,
    box-shadow 0.35s ease,
    backdrop-filter 0.35s ease;
}

/* =====================================================
   Container
===================================================== */
.container {
  width: 100%;
  max-width: 1400px;
  height: 68px;
  padding: 0 40px;
  display: flex;
  align-items: center;
}

/* =====================================================
   Logo
===================================================== */
.logo-wrapper {
  width: 200px;
  height: 68px;
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

/* =========================================
   FELIX BLOG Logo
========================================= */
.blog-logo {
  display: flex;
  align-items: baseline;
  gap: 7px;
  color: #ffffff;
  cursor: pointer;
  user-select: none;
  transition:
    color 0.35s ease;
}

.logo-main {
  font-size: 20px;
  font-weight: 800;
  letter-spacing: 1px;
}

.logo-sub {
  font-size: 12px;
  font-weight: 400;
  letter-spacing: 2.5px;
  opacity: 0.65;
}

/* =====================================================
   Menu
===================================================== */
.web-menu {
  margin-left: 20px;
  flex: 1;
}

/* Element Plus Menu */
.web-menu {
  :deep(.el-menu) {
    height: 68px;
    display: flex;
    align-items: center;
    background: transparent !important;
    border-bottom: none !important;
  }

  :deep(.el-menu-item) {
    position: relative;
    height: 68px;
    padding: 0 18px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #a8a8a8 !important;
    background: transparent !important;
    border-bottom: none !important;
    font-size: 15px;
    font-weight: 400;
    letter-spacing: 0.5px;
    transition:
      color 0.25s ease,
      transform 0.25s ease;
  }

  /* hover */
  :deep(.el-menu-item:hover) {
    color: #ffffff !important;
    background: transparent !important;
    transform: translateY(-1px);
  }

  /* =================================================
     自定义激活下划线
  ================================================= */
  :deep(.el-menu-item::after) {
    content: "";
    position: absolute;
    left: 50%;
    bottom: 11px;
    width: 0;
    height: 2px;
    transform: translateX(-50%);
    background: #ffffff;
    border-radius: 10px;
    transition:
      width 0.3s ease;
  }

  /* hover 下划线 */

  :deep(.el-menu-item:hover::after) {

    width: 26px;

  }


  /* active */

  :deep(.el-menu-item.is-active) {

    color: #ffffff !important;

    background: transparent !important;

  }


  :deep(.el-menu-item.is-active::after) {

    width: 28px;

    background: #ffffff;

  }

}


/* =====================================================
   Navbar Right
===================================================== */

.navbar-right {

  margin-left: auto;

  display: flex;

  align-items: center;

  gap: 20px;

  flex-shrink: 0;

}


/* =====================================================
   Search
===================================================== */

.search-box {

  width: 165px;

}


.search-box {

  :deep(.el-input__wrapper) {

    height: 34px;

    padding: 0 12px;

    background:
      rgba(255, 255, 255, 0.08);

    border: 1px solid rgba(255, 255, 255, 0.08);

    box-shadow: none !important;

    border-radius: 20px;

    transition:
      background 0.3s ease,
      border-color 0.3s ease;

  }


  :deep(.el-input__inner) {

    color: #ffffff;

    font-size: 12px;

  }


  :deep(.el-input__inner::placeholder) {

    color:
      rgba(255, 255, 255, 0.45);

  }


  :deep(.el-input__prefix) {

    color: #888;

  }

}


/* =====================================================
   AuthPopover
===================================================== */

.navbar-right {

  :deep(.auth-popover) {

    margin: 0;

    padding: 0;

  }

}


/* =====================================================
   滚动后的 Navbar
===================================================== */

.web-navbar.show {

  /**
   * 白色半透明
   */

  background:
    rgba(255, 255, 255, 0.88);

  /**
   * 毛玻璃
   */

  backdrop-filter:
    blur(18px);

  -webkit-backdrop-filter:
    blur(18px);

  /**
   * 边框
   */

  border-bottom:
    1px solid rgba(0, 0, 0, 0.07);

  /**
   * 阴影
   */

  box-shadow:
    0 5px 25px rgba(0, 0, 0, 0.08);

  .blog-logo {
    color: #111;
  }

  /* =================================================
     Menu
  ================================================= */

  .web-menu {

    :deep(.el-menu-item) {

      color: #666 !important;

    }


    :deep(.el-menu-item:hover) {

      color: #111 !important;

      background:
        transparent !important;

    }


    :deep(.el-menu-item.is-active) {

      color: #111 !important;

      background:
        transparent !important;

    }


    :deep(.el-menu-item::after) {

      background: #111;

    }

  }


  /* =================================================
     Search
  ================================================= */

  .search-box {

    :deep(.el-input__wrapper) {

      background:
        rgba(0, 0, 0, 0.045);

      border-color:
        rgba(0, 0, 0, 0.06);

    }


    :deep(.el-input__inner) {

      color: #222;

    }


    :deep(.el-input__inner::placeholder) {

      color: #999;

    }


    :deep(.el-input__prefix) {

      color: #777;

    }

  }

}


/* =====================================================
   Responsive
===================================================== */

@media screen and (max-width: 1200px) {

  .container {

    padding: 0 25px;

  }


  .logo-wrapper {

    width: 170px;

  }


  .web-menu {

    margin-left: 10px;

  }


  :deep(.el-menu-item) {

    padding: 0 14px !important;

  }

}


@media screen and (max-width: 900px) {

  .logo-wrapper {

    width: 150px;

  }


  .web-menu {

    margin-left: 0;

  }


  :deep(.el-menu-item) {

    padding: 0 10px !important;

    font-size: 14px !important;

  }


  .search-box {

    display: none;

  }

}


@media screen and (max-width: 650px) {

  .web-navbar {

    height: 60px;

  }


  .container {

    height: 60px;

    padding: 0 15px;

  }


  .logo-wrapper {

    width: auto;

  }


  .web-menu {

    display: none;

  }

}
</style>