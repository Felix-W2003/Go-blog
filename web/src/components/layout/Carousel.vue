<template>
  <div class="carousel">
    <!-- <el-carousel trigger="click" height="700px">
      <el-carousel-item v-for="item in imgList" :key="item">
        <el-image fit="cover" :src=item alt=""></el-image>
      </el-carousel-item>
    </el-carousel>
    -->
    <el-carousel trigger="click" height="700px">
     <el-carousel-item>
        <demo1/>
      </el-carousel-item>
      <el-carousel-item>
        <demo2/>
      </el-carousel-item>
      <el-carousel-item>
        <demo3/>
      </el-carousel-item>
    </el-carousel>
  </div>
</template>

<script setup lang="ts">
import {ref} from "vue";
import {websiteCarousel} from "@/api/website";
import demo1 from "@/views/web/test/demo1.vue"
import demo2 from "@/views/web/test/demo2.vue"
import demo3 from "@/views/web/test/demo3.vue"

const imgList = ref<string[]>([
  '/image/carousel_1.jpg',
  '/image/carousel_2.jpg',
  '/image/carousel_3.jpg',
  '/image/carousel_4.jpg',
])

const getWebsiteCarousel = async () => {
  const res = await websiteCarousel()
  if (res.code === 0 && res.data.length !== 0) {
    imgList.value = res.data
  }
}

getWebsiteCarousel()
</script>

<style scoped lang="scss">
.carousel {
  width: 100%;
  max-width: 1400px;

  position: relative;
  margin:0 auto;
  margin-top: 0px;
  .el-carousel {
    .el-image {
      width: 100%;
      height: 100%;
    }
  }
}
</style>
