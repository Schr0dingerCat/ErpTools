# ErpTools
backend use go + gin, frontend use vue3 + element-plus + axios

# Q&A
Q1、后端需要定义与前端对应的结构体，前端setup中使用ref()定义const常量对象a，需要使用a.value进行赋值
Q2、使用 el-select option 获取的数据量大时，会有卡顿
A2、之后使用 el-select-v2，默认的v2列表框显示文字位置不对
Q3、前端使用setup导出 form，form.xxx 无法取到select 选择的值