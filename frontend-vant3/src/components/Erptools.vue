<script setup lang="ts">
import { ref } from "vue";
const showloading = ref(false);

const onSubmit = (values) => {
  console.log("submit", values);
};

// calendar
const currentDate = ref(new Date());
const minDate = ref(new Date(2010, 0, 1));
const maxDate = ref(new Date(2030, 0, 1));
const date = ref("");
const showCalendar = ref(false);
const formatDate = (date) => `${date.getMonth() + 1}/${date.getDate()}`;
const onCalendarClick = () => {
  showCalendar.value = true;
};
const onCalendarConfirm = (values) => {
  const [start, end] = values;
  showCalendar.value = false;
  date.value = `${formatDate(start)} - ${formatDate(end)}`;
};
</script>

<template>
  <van-overlay :show="showloading">
    <div class="wrapper">
      <van-loading color="#1989fa">加载中...</van-loading>
    </div>
  </van-overlay>

  <van-form @submit="onSubmit">
    <van-cell-group inset>
      <van-cell title="选择日期区间" :value="date" @click="onCalendarClick" />
      <van-calendar
        v-model:show="showCalendar"
        type="range"
        :min-date="minDate"
        :max-date="maxDate"
        :show-confirm="false"
        @confirm="onCalendarConfirm"
      />
    </van-cell-group>
    <div style="margin: 16px">
      <van-button round block type="primary" native-type="submit">提交</van-button>
    </div>
  </van-form>
</template>

<style scoped>
.wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}
</style>
