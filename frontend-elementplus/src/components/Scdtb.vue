<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import axios from "axios";

const innerHeight = ref(window.innerHeight);
const outerHeight = ref(window.outerHeight);
const tableMaxHeight = ref(innerHeight.value - 50 * 3 - 40 - 34);

// 定义常量对象，使用时需要用 xx.value
// 全屏loading
const fullscreenLoading = ref(true);
// tab初始化
const activeTab = ref("tab1");
const prdnocp = ref();
const qtycp = ref();
// 货品列表
const options = ref([
  {
    value: "",
    label: "",
  },
]);
const form = reactive({
  prdno: "",
  dateso: [""],
});
// el-date-picker 为了适配手机端，不使用daterange区间模式
const dateso1 = ref("");
const dateso2 = ref("");

const tableData1 = reactive([
  {
    sono: "",
    estitmso: "",
    prdno: "",
    prdname: "",
    qtyso: "",
    cusname: "",
    estdd: "",
    clsmpid: "",
    mono: "",
    qtysolj: "",
    biltype: "",
    status: "",
  },
]);
const tableData2 = reactive([
  {
    tzno: "",
    depname: "",
    zcname: "",
    qty: "",
    qtyfin: "",
    qtylost: "",
    qtybf: "",
    qtysy: "",
    qtypgs: "",
    mydinge: "",
  },
]);
const tableData3 = reactive([
  {
    prdno: "",
    zcname: "",
    qty: "",
  },
]);
const tableData4 = reactive([
  {
    prdno: "",
    zcname: "",
    batno: "",
    qty: "",
  },
]);
// 初始化
tableData1.length = 0;
tableData2.length = 0;
tableData3.length = 0;
tableData4.length = 0;

// 界面加载后调用
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

// TODO: 发送数据，后台查询
const onSubmit = () => {
  form.dateso = [dateso1.value, dateso2.value];
  fullscreenLoading.value = true;
  axios
    .post("/erptools", {
      cmd: "getsolist",
      args: JSON.stringify(form),
    })
    .then((response) => {
      tableData1.length = 0;
      if (response.data.sodatas != null) {
        tableData1.push(...response.data.sodatas);
      }
      tableData2.length = 0;
      tableData3.length = 0;
      tableData4.length = 0;
      prdnocp.value = "";
      qtycp.value = "";
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

// 根据单元格值设置行颜色
// 设置stripe后，在改变行颜色，偶数行无法改变
const setRowStyle = ({ row, rowIndex }) => {
  let stylejson = {};
  if (rowIndex % 2 === 0) {
    stylejson = { "background-color": "#fafafa" };
  } else {
    stylejson = { "background-color": "#ffffff" };
  }
  if (row.status === "未完成") {
    stylejson = { "background-color": "yellow" };
  }
  return stylejson;
};
// 受订单列表点击某行
const onSoRowClick = (row: any, column: any, event: any) => {
  if (row["mono"] != "") {
    // loadling
    fullscreenLoading.value = true;
    axios
      .post("/erptools", {
        cmd: "gettzlist",
        mono: row["mono"],
        prdno: row["prdno"],
      })
      .then((response) => {
        tableData2.length = 0;
        if (response.data.tzdatas != null) {
          tableData2.push(...response.data.tzdatas);
        }
        prdnocp.value = row["prdno"];
        qtycp.value = response.data.qtycp;
        tableData3.length = 0;
        if (response.data.bcpsdatas != null) {
          tableData3.push(...response.data.bcpsdatas);
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
        total3.value = tableData3.length;
      });
  }
};

// 半成品列表点击某行
const onBcpRowClick = (row: any, column: any, event: any) => {
  // loading
  fullscreenLoading.value = true;
  axios
    .post("/erptools", {
      cmd: "getbcplist",
      prdno: row["prdno"],
    })
    .then((response) => {
      tableData4.length = 0;
      if (response.data.bcpdatas != null) {
        tableData4.push(...response.data.bcpdatas);
      }
      // 切换tab
      activeTab.value = "tab4";
    })
    .catch((error) => {
      console.log(error);
    })
    .finally(() => {
      fullscreenLoading.value = false;
      total4.value = tableData4.length;
    });
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
// table3分页
const currentPage3 = ref(1);
const pageSize3 = ref((tableMaxHeight.value - 48) / 71);
const total3 = ref(0);
// table4分页
const currentPage4 = ref(1);
const pageSize4 = ref((tableMaxHeight.value - 48) / 71);
const total4 = ref(0);
</script>

<template>
  <div
    v-loading.fullscreen.lock="fullscreenLoading"
    element-loading-text="Loading..."
    element-loading-background="rgba(0, 0, 0, 0.8)"
  ></div>
  <el-row justify="center">
    <el-form :model="form">
      <!-- <el-form-item label="订单预交日"> -->
      <el-space wrap>
        <el-form-item>
          <el-date-picker
            v-model="dateso1"
            type="date"
            placeholder="订单预交日起始日期"
          ></el-date-picker>
        </el-form-item>
        <el-form-item>
          <el-date-picker
            v-model="dateso2"
            type="date"
            placeholder="订单预交日截止日期"
          ></el-date-picker>
        </el-form-item>
        <!-- <el-form-item label="货品代号"> -->
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
    <el-tab-pane label="受订单" name="tab1">
      <div>
        <el-table
          :data="
            tableData1.slice((currentPage1 - 1) * pageSize1, currentPage1 * pageSize1)
          "
          :height="tableMaxHeight"
          :max-height="tableMaxHeight"
          highlight-current-row
          :row-style="setRowStyle"
          @row-click="onSoRowClick"
        >
          <el-table-column
            prop="sono"
            label="受订单号"
            width="120"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="prdno"
            label="货品代号"
            width="130"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="qtyso"
            label="受订数量"
            width="80"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="estdd"
            label="预交日期"
            width="100"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="mono"
            label="制令单号"
            width="80"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="status"
            label="是否完成"
            width="80"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="cusname"
            label="客户名称"
            width="150"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="prdname"
            label="货品名称"
            width="80"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="biltype"
            label="订单类型"
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
    <el-tab-pane label="制令单" name="tab2">
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
            prop="zcname"
            label="制程名称"
            width="80"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="qty"
            label="数量"
            width="70"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="qtyfin"
            label="已完数量"
            width="80"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="qtysy"
            label="剩余数量"
            width="80"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="qtypgs"
            label="已派工量"
            width="80"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="mydinge"
            label="当前定额"
            width="80"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="qtylost"
            label="不合格数"
            width="80"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="qtybf"
            label="报废数量"
            width="80"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="tzno"
            label="通知单号"
            width="80"
            :show-overflow-tooltip="true"
          />
          <el-table-column
            prop="depname"
            label="部门名称"
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
    <el-tab-pane label="库存" name="tab3">
      <el-row justify="start">
        <div>成品库存: {{ prdnocp }}: {{ qtycp }}</div>
      </el-row>
      <div>
        <el-table
          :data="
            tableData3.slice((currentPage3 - 1) * pageSize3, currentPage3 * pageSize3)
          "
          stripe
          :height="tableMaxHeight - 21"
          :max-height="tableMaxHeight - 21"
          highlight-current-row
          @row-click="onBcpRowClick"
        >
          <el-table-column prop="prdno" label="货品代号" width="80" />
          <el-table-column prop="zcname" label="制程名称" width="80" />
          <el-table-column prop="qty" label="数量" width="80" />
        </el-table>
      </div>
      <el-row justify="start">
        <div>
          <el-pagination
            v-model:currentPage="currentPage3"
            :page-size="pageSize3"
            :pager-count="5"
            :hide-on-single-page="true"
            layout="prev, pager, next, jumper"
            :total="total3"
          >
          </el-pagination>
        </div>
      </el-row>
    </el-tab-pane>
    <el-tab-pane label="半成品区位" name="tab4">
      <div>
        <el-table
          :data="
            tableData4.slice((currentPage4 - 1) * pageSize4, currentPage4 * pageSize4)
          "
          stripe
          :height="tableMaxHeight"
          :max-height="tableMaxHeight"
          highlight-current-row
        >
          <el-table-column prop="prdno" label="货品代号" width="80" />
          <el-table-column prop="zcname" label="制程名称" width="80" />
          <el-table-column prop="batno" label="批次号" width="80" />
          <el-table-column prop="qty" label="数量" width="80" />
        </el-table>
      </div>
      <el-row justify="start">
        <div>
          <el-pagination
            v-model:currentPage="currentPage4"
            :page-size="pageSize4"
            :pager-count="5"
            :hide-on-single-page="true"
            layout="prev, pager, next, jumper"
            :total="total4"
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
