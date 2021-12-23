# ErpTools
backend use go + gin, frontend use vue3 + element-plus + axios

# Q&A
Q1、后端需要定义与前端对应的结构体，前端setup中使用ref()定义const常量对象a，需要使用a.value进行赋值
Q2、使用 el-select option 获取的数据量大时，会有卡顿
A2、之后使用 el-select-v2，默认的v2列表框显示文字位置不对,通过设置 .el-select-v2 style
Q3、前端使用setup导出 form，<template>中使用 ref="form"，form.xxx 无法取到select 选择的值，具体需要查看 template Refs的文档
Q4、前端取到的时间为UTC时间非本地时间
A4、前端根据dayjs文档，引入，但并未设置正确，程序在后端进行UTC转本地
Q5、bug：设置stripe后，在改变行颜色，偶数行无法改变