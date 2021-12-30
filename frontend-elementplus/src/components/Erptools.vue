<script setup lang="ts">
import { ref } from "vue";
import { Setting } from "@element-plus/icons-vue";
import Pgtjb from "./Pgtjb.vue";
import Scdtb from "./Scdtb.vue";

const innerHeight = ref(window.innerHeight);
const outerHeight = ref(window.outerHeight);

// el-header command 事件
const mytitle = ref("生产动态表");
const showscdtb = ref(true);
const showpgtjb = ref(false);
const handleCommand = (command: any) => {
  switch (command) {
    case "scdtb":
      mytitle.value = "生产动态表";
      showscdtb.value = true;
      showpgtjb.value = false;
      break;
    case "pgtjb":
      mytitle.value = "派工统计表";
      showscdtb.value = false;
      showpgtjb.value = true;
      break;
    default:
      break;
  }
};
</script>

<template>
  <el-container direction="vertical">
    <el-affix position="top">
      <div>
        <el-header style="text-align: right; font-size: 25px">
          <el-row>
            <el-col :span="23" style="text-align: center; font-size: 20px">
              {{ mytitle }}
            </el-col>
            <el-col :span="1">
              <el-dropdown @command="handleCommand">
                <el-icon><setting /></el-icon>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="scdtb">生产动态表</el-dropdown-item>
                    <el-dropdown-item command="pgtjb">派工统计表</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </el-col>
          </el-row>
        </el-header>
      </div>
    </el-affix>
    <el-main>
      <div v-show="showscdtb">
        <scdtb></scdtb>
      </div>
      <div v-show="showpgtjb">
        <pgtjb></pgtjb>
      </div>
    </el-main>
  </el-container>
</template>

<style scoped>
.el-header {
  background-color: #b3c0d1;
  color: var(--el-text-color-primary);
  max-height: 30px;
}
</style>
