<template>
  <el-form :model="form" label-width="120px">
    <el-form-item label="订单预交日">
      <el-date-picker
        v-model="form.date1"
        type="daterange"
        start-placeholder="起始日期"
        end-placeholder="截止日期"
      ></el-date-picker>
    </el-form-item>
    <el-form-item label="货品代号">
      <el-select v-model="form.prdno" clearable filterable placeholder="请选择">
        <el-option
          v-for="item in options"
          :key="item.value"
          :label="item.label"
          :value="item.value"
        ></el-option>
      </el-select>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="onSubmit">查询</el-button>
    </el-form-item>
  </el-form>
</template>

<script lang="ts">
import { ref, reactive, defineComponent } from "vue";
import axios from "axios";
export default defineComponent({
  setup() {
    // 定义常量对象，使用时需要用 xx.value
    const options = ref([
      {
        value: "",
        label: "",
      },
    ]);

    // 加载 /erptools界面是自动获取数据
    axios
      .post("/erptools", {
        cmd: "getprdno",
      })
      .then(function (response) {
        options.value = response.data.options;
      })
      .catch(function (error) {
        console.log(error);
      });

    // TODO: 发送数据，后台查询
    const onSubmit = () => {
      console.log("submit!");
      axios
        .post("/erptools", {
          cmd: "getsolist",
        })
        .then((response) => {})
        .catch((error) => {});
    };

    return {
      form: ref({
        prdno: "",
        date1: "",
      }),
      options,
      onSubmit,
    };
  },
});
</script>
