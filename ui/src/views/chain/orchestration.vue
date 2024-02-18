<script setup>
import { NTabs, NTabPane, NSwitch } from "naive-ui";
import ControlBar from "../chain/orchestration/control-bar.vue";
import Graph from '../chain/orchestration/graph.vue'
import ShortcutTree from '../chain/orchestration/shortcut-tree.vue'
import { useRoute } from "vue-router";
import {onMounted, ref, watch} from "vue";
import ShortcutDetail from "./orchestration/shortcut-detail.vue";
import {processGraphDataToChainData} from "../../services/graph.js";
import {useChainInfo} from "../../hooks/chain.js";
const route = useRoute();

// 编辑模式是否开启
const editing = ref(false)

const graphRef = ref(null)

// 图数据
const data = ref([])

const { chainDetail } = useChainInfo(route.params.chainId)

watch(chainDetail,() => {
  data.value = JSON.parse(chainDetail.value.chain.originData)
})

// 选中的快捷指令（用于展示详情）
const selectedShortcut = ref(null)

const setShortcut = (shortcut) => {
  selectedShortcut.value = shortcut
}

// 用于拖动来创建节点 详情见AntV/x6 官网插件dnd
const onStartDrag = (e,option) => {
    graphRef.value.startDrag(e,option)
}

// 编辑模式切换时触发
const onSwitch = (val) => {
  if (val === false){
    data.value =  graphRef.value.exportData()
    console.log(data.value)
    const processedData =  processGraphDataToChainData(data.value,route.params.chainId)
    console.log(processedData)
  }
}




</script>

<template>

  <div class="h-full w-full">

    <div class="flex gap-2 absolute top-5 left-1/2 z-20 -translate-x-1/2">
      编辑模式
      <n-switch v-model:value="editing" :on-update-value="onSwitch">

      </n-switch>
      <div class="w-72 h-full">

      </div>
    </div>

    <div v-if="editing" class="flex h-full w-full relative">
      <div class="graph grow relative">
        <graph key="1" ref="graphRef" :editable="true" :data="data" @onClickShortcut="setShortcut"></graph>
      </div>
      <div class="info w-72 h-full border-[1px] p-1 border-gray-600 rounded">
        <n-tabs type="line" animated justify-content="space-evenly">
          <n-tab-pane name="指令详情" tab="指令详情"> <shortcut-detail :shortcut="selectedShortcut"></shortcut-detail> </n-tab-pane>
          <n-tab-pane name="指令树" tab="指令树"> <shortcut-tree @startDrag="onStartDrag"></shortcut-tree> </n-tab-pane>
        </n-tabs>
      </div>
    </div>


    <div v-if="!editing" class="flex h-full w-full">
      <div class="graph grow relative">
        <control-bar class="absolute left-5 top-1/2 -translate-y-1/2 z-30"></control-bar>
        <graph key="2" :editable="false" :data="data"></graph>
      </div>
      <div class="info w-72 h-full border-[1px] border-gray-600 rounded">
        <n-tabs type="line" animated justify-content="space-evenly">
          <n-tab-pane name="调度记录" tab="调度记录"> 调度记录 </n-tab-pane>
          <n-tab-pane name="节点详情" tab="节点详情"> 节点详情 </n-tab-pane>
          <n-tab-pane name="指令树" tab="指令树"> <shortcut-tree></shortcut-tree> </n-tab-pane>
        </n-tabs>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
