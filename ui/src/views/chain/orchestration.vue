<script setup>
import {NButton, NSwitch, NTabPane, NTabs, useMessage} from "naive-ui";
import ControlBar from "../chain/orchestration/control-bar.vue";
import Graph from '../chain/orchestration/graph.vue'
import ShortcutTree from '../chain/orchestration/shortcut-tree.vue'
import {useRoute} from "vue-router";
import {onMounted, ref, watch} from "vue";
import ShortcutDetail from "./orchestration/shortcut-detail.vue";
import {processGraphDataToChainData} from "../../services/graph.js";
import {useChainInfo, useMachine} from "../../hooks/chain.js";
import {getExecDetail, getExecList, loadData} from "../../services/chain.js";
import ExecHostory from "./orchestration/exec-hostory.vue";
import NodeExecDetail from "./orchestration/node-exec-detail.vue";


const route = useRoute();
const message = useMessage()

// 编辑模式是否开启
const editing = ref(false)

const graphRef = ref(null)

// 图数据
const data = ref([])

const { chainDetail } = useChainInfo(route.params.chainId)

watch(chainDetail,() => {
  data.value =  chainDetail.value.chain.originData === "" ? [] : JSON.parse(chainDetail.value.chain.originData)
})

// 选中的快捷指令（用于展示详情）
const selectedNode = ref(null)

const setNode = (node) => {
  selectedNode.value = node
}

// 用于拖动来创建节点 详情见AntV/x6 官网插件dnd
const onStartDrag = (e,option) => {
    graphRef.value.startDrag(e,option)
}

// 编辑模式切换时触发
const onSwitch = (val) => {
  if (val === false){
  }
}

const onSave = async () => {
  let graphData = graphRef.value.exportData()
  let cells = graphData['cells']
  let rootCount = 0

  console.log(cells)

  cells.filter(item => item.shape === 'custom-node').forEach(item => item.data.root ?? false ? rootCount++ : null)

  if (rootCount !== 1) {
    message.error('需要设定一个根节点为入口')
    return
  }

  const processedData =  processGraphDataToChainData(graphData,route.params.chainId)



  let res = await loadData(processedData)

  if (!res.data) {
    message.error(res.msg)
    return
  }

  data.value = graphData
  message.success('保存成功')

}

///////////////////////////////////////////
// 图引用
const dispatchGraphRef = ref(null)

// 详细信息
const graphData = ref([])
const nodeExecResults = ref([])
const execInfo = ref({})

const execList = ref([])
onMounted(() => {
  getExecList(route.params.chainId).then(res => {
    if (res['code'] !== 200){
      return
    }
    execList.value = res['data']
  })
})

const onClickExecItem = async (dispatch) => {
  const res = await getExecDetail(dispatch.id);
  if (res['code'] !== 200 ){
    return
  }

  if (res['data']['info']['chainData'] !== ''){
    graphData.value = JSON.parse(res['data']['info']['chainData'])
  }

  execInfo.value = res['data']['info']
  nodeExecResults.value = res['data']['execResults']
}



const selectedNodeResult = ref(null)
const onClickNodeExecResult = (node) => {
  selectedNodeResult.value = node.getData()['nodeExecResult']
}


const { NewMachine,MachineExecOne,MachineExecAll } = useMachine(route.params.chainId)

const onClickNew = async () => {
  const execInfo = await NewMachine(route.params.chainId)
  execList.value.unshift(execInfo)
  await onClickExecItem(execInfo)
}


const onClickNextOne = async () => {
  const res = await MachineExecOne(execInfo.value.id)
  if (res === null) return
  nodeExecResults.value = res
}

const onClickNextAll = async () => {
  const res = await MachineExecAll(execInfo.value.id)
  if (res === null) return
  nodeExecResults.value = res
}


const onClickBtnInControlBar = (type) => {
  switch (type){
    case 'add': onClickNew();break;
    case 'next': onClickNextOne();break;
    case 'nextAll': onClickNextAll();break;
    case 'remove': message.warning('remove'); break;
    default: message.warning('unknown behavior')
  }
}

</script>

<template>

  <div class="h-full w-full">

    <div class="flex gap-2 absolute top-5 left-1/2 z-20 items-center -translate-x-1/2">
      编辑模式
      <n-switch v-model:value="editing" :on-update-value="onSwitch">
      </n-switch>
      <n-button v-show="editing" @click="onSave" type="success" class="h-6 bg-green-500">保存</n-button>
      <div class="w-72 h-full">
      </div>
    </div>

    <div v-if="editing" class="flex h-full w-full relative">
      <div class="grow relative">
        <graph key="1" ref="graphRef" :editable="true" :data="data" @onClickNode="setNode"></graph>
      </div>
      <div class="info w-72 h-full border-[1px] p-1 border-gray-600 rounded">
        <n-tabs class="h-full" type="line" animated justify-content="space-evenly">
          <n-tab-pane name="指令详情" tab="指令详情"> <shortcut-detail :node="selectedNode"></shortcut-detail> </n-tab-pane>
          <n-tab-pane name="指令树" tab="指令树"> <shortcut-tree @startDrag="onStartDrag"></shortcut-tree> </n-tab-pane>
        </n-tabs>
      </div>
    </div>


    <div v-if="!editing" class="flex h-full w-full">
      <div class="grow relative">
        <control-bar @onClickBtn="onClickBtnInControlBar" class="absolute left-5 top-1/2 -translate-y-1/2 z-30"></control-bar>
        <graph key="2" :editable="false" ref="dispatchGraphRef" @onClickNodeWithResult="onClickNodeExecResult" :data="graphData" :node-exec-results="nodeExecResults"></graph>
      </div>
      <div class="info w-72 min-w-[18rem] h-full border-[1px] right-1 border-gray-600 rounded">
        <n-tabs class="h-full" type="line" animated justify-content="space-evenly" :pane-style="{height:'100%',width:'100%',padding:'0'}" :pane-wrapper-style="{height:'100%',width:'100%'}">
          <n-tab-pane  name="调度记录" tab="调度记录">
            <exec-hostory :data="execList" @onClickItem="onClickExecItem"></exec-hostory>
          </n-tab-pane>
          <n-tab-pane name="节点详情" tab="节点详情">
            <node-exec-detail :data="selectedNodeResult"></node-exec-detail>
          </n-tab-pane>
        </n-tabs>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
