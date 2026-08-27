<template>
  <el-card class="advertisement">
    <el-row class="title">独家推广</el-row>
    <el-carousel :interval="5000" type="card" height="320px">
      <el-carousel-item v-for="advertisement in advertisementList" :key="advertisement">
        <el-image :src=advertisement.ad_image alt="" @click=handleAdverisementClick(advertisement)></el-image>
        <!-- <el-row class="title">{{ advertisement.title }}</el-row> -->
        <el-text class="text">{{ advertisement.content }}</el-text>
      </el-carousel-item>
    </el-carousel>
  </el-card>
</template>

<script setup lang="ts">
import {ref} from "vue";
import {type Advertisement, advertisementInfo} from "@/api/advertisement";

const advertisementList = ref<Advertisement[]>()
const getAdvertisementList = async () => {
  const res = await advertisementInfo()
  if (res.code == 0) {
    advertisementList.value = res.data.list
  }
}
getAdvertisementList()

const handleAdverisementClick = (advertisement: Advertisement) => {
  window.open(advertisement.link)
}
</script>

<style scoped lang="scss">
.advertisement {
  margin-bottom: 20px;
  border-radius: 5px;
    background: linear-gradient(
    145deg,
    #000000 0%,
    #1d1d1d 55%,
    #2a2a2a 100%
  );
  .title {
    font-size: 25px;
    font-weight: 600;
    color: #fff;
    margin-bottom: 20px;
  }

  .el-carousel__item {
    background-color: #fcfcfc;
     border-radius: 1px;
    .el-image {
      height: 260px;
      width: 100%;
    }
    .text{
      color:#000000;
      font-size: 16px;
    }
    .title{
      color: #000000;
      font-size: 20px;
    }
  }
}
</style>
