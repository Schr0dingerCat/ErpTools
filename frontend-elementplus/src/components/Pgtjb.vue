<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import axios from "axios";

const innerHeight = ref(window.innerHeight);
const tableMaxHeight = ref(innerHeight.value - 50 * 3 - 40 - 34);

const fullscreenLoading = ref(true);

const form = reactive({
  prdno: "",
  zcno: "",
  datepg: [""],
});
// 货品列表
const options = ref([
  {
    value: "",
    label: "",
  },
]);
// tab初始化
const activeTab = ref("tab1");

const date1 = ref("");
const date2 = ref("");

const tableData1 = reactive([
  {
    prdno: "",
    zcno: "",
    zcname: "",
    qtypg: "",
    qtywr: "",
    qtywrlost: "",
  },
]);
const tableData2 = reactive([
  {
    pgdd: "",
    pgno: "",
    prdno: "",
    zcname: "",
    sbno: "",
    tzno: "",
    batno: "",
    qtypg: "",
    qtywr: "",
    qtywrlost: "",
  },
]);

tableData1.length = 0;
tableData2.length = 0;

onMounted(() => {
  // 加载 /erptools界面是自动获取数据
  axios
    .post("/erptools", {
      cmd: "getprdno",
    })
    .then(function (response) {
      if (response.data.options != null) {
        options.value.length = 0;
        options.value = response.data.options;
      }
    })
    .catch(function (error) {
      console.log(error);
    })
    .finally(function () {
      fullscreenLoading.value = false;
    });
});

// 获取派工单统计表
const onSubmit = () => {
  form.datepg = [date1.value, date2.value];
  form.zcno = "";
  fullscreenLoading.value = true;
  axios
    .post("/erptools", {
      cmd: "getpgtjlist",
      args: JSON.stringify(form),
    })
    .then((response) => {
      tableData1.length = 0;
      if (response.data.pgtjdatas != null) {
        tableData1.push(...response.data.pgtjdatas);
      }
      tableData2.length = 0;
    })
    .catch((error) => {
      console.log(error);
    })
    .finally(() => {
      fullscreenLoading.value = false;
      activeTab.value = "tab1";
      total1.value = tableData1.length;
    });
};
//
const onRowClick = (row: any, column: any, event: any) => {
  if (row.prdno != "") {
    // loadling
    fullscreenLoading.value = true;
    axios
      .post("/erptools", {
        cmd: "getpgmxlist",
        args: JSON.stringify({
          prdno: row["prdno"],
          zcno: row["zcno"],
          datepg: form.datepg,
        }),
      })
      .then((response) => {
        tableData2.length = 0;
        if (response.data.pgmxdatas != null) {
          tableData2.push(...response.data.pgmxdatas);
        }
        // 切换tab
        activeTab.value = "tab2";
      })
      .catch((error) => {
        console.log(error);
      })
      .finally(() => {
        fullscreenLoading.value = false;
        total2.value = tableData2.length;
      });
  }
};

// table1分页
const currentPage1 = ref(1);
const pageSize1 = ref((tableMaxHeight.value - 48) / 48);
const total1 = ref(0);
const handleSizeChange1 = (val: any) => {};
const handleCurrentChange1 = (val: any) => {};
// table2分页
const currentPage2 = ref(1);
const pageSize2 = ref((tableMaxHeight.value - 48) / 48);
const total2 = ref(0);
</script>

<template>
  <div
    v-loading.fullscreen.lock="fullscreenLoading"
    element-loading-text="Loading..."
    element-loading-background="rgba(0, 0, 0, 0.8)"
  ></div>
  <el-row justify="center">
    <el-form :model="form">
      <el-space wrap>
        <el-form-item>
          <el-date-picker
            v-model="date1"
            type="date"
            placeholder="派工起始日期"
          ></el-date-picker>
        </el-form-item>
        <el-form-item>
          <el-date-picker
            v-model="date2"
            type="date"
            placeholder="派工截止日期"
          ></el-date-picker>
        </el-form-item>
        <el-form-item>
          <el-select-v2
            v-model="form.prdno"
            :options="options"
            clearable
            filterable
            placeholder="货品代号"
          ></el-select-v2>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSubmit">查询</el-button>
        </el-form-item>
      </el-space>
    </el-form>
  </el-row>
  <el-tabs v-model="activeTab" stretch>
    <el-tab-pane label="派工单统计" name="tab1">
      <div>
        <el-table
          :data="
            tableData1.slice((currentPage1 - 1) * pageSize1, currentPage1 * pageSize1)
          "
          stripe
          :height="tableMaxHeight"
          :max-height="tableMaxHeight"
          highlight-current-row
          @row-click="onRowClick"
        >
          <el-table-column
            prop="prdno"
            label="货品代号"
            width="130"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="zcno"
            label="工序代号"
            width="80"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="zcname"
            label="工序名称"
            width="100"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="qtypg"
            label="派工数量"
            width="80"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="qtywr"
            label="上帐数量"
            width="80"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="qtywrlost"
            label="上帐废品"
            width="80"
            :show-overflow-tooltip="true"
          />
        </el-table>
      </div>
      <el-row justify="start">
        <div>
          <el-pagination
            v-model:currentPage="currentPage1"
            :page-size="pageSize1"
            :pager-count="5"
            :hide-on-single-page="true"
            layout="prev, pager, next, jumper"
            :total="total1"
            @size-change="handleSizeChange1"
            @current-change="handleCurrentChange1"
          >
          </el-pagination>
        </div>
      </el-row>
    </el-tab-pane>
    <el-tab-pane label="派工单明细" name="tab2">
      <div>
        <el-table
          :data="
            tableData2.slice((currentPage2 - 1) * pageSize2, currentPage2 * pageSize2)
          "
          stripe
          highlight-current-row
          :height="tableMaxHeight"
          :max-height="tableMaxHeight"
        >
          <el-table-column
            prop="pgdd"
            label="派工日期"
            width="100"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="pgno"
            label="派工单号"
            width="130"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="prdno"
            label="货品代号"
            width="130"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="zcname"
            label="工序名称"
            width="80"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="sbno"
            label="设备代号"
            width="100"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="tzno"
            label="通知单号"
            width="80"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="batno"
            label="批次号"
            width="80"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="qtypg"
            label="派工数量"
            width="80"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="qtywr"
            label="上帐数量"
            width="80"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="qtywrlost"
            label="上帐废品"
            width="80"
            :show-overflow-tooltip="true"
          />
        </el-table>
      </div>
      <el-row justify="start">
        <div>
          <el-pagination
            v-model:currentPage="currentPage2"
            :page-size="pageSize2"
            :pager-count="5"
            :hide-on-single-page="true"
            layout="prev, pager, next, jumper"
            :total="total2"
          >
          </el-pagination>
        </div>
      </el-row>
    </el-tab-pane>
  </el-tabs>
</template>

<style scoped>
.el-form-item {
  margin-bottom: 0px;
}
.el-select-v2 {
  width: 220px;
  text-align: left;
}
.el-table {
  width: 100%;
}
</style>
