<script setup>
import {Add, ArrowDownCircle, CaretDown, Play, Remove, ReturnUpForward, SettingsOutline} from '@vicons/ionicons5'
import {NButton, NDataTable, NDrawer, NDrawerContent, NIcon, NLog,NEmpty,NThing} from 'naive-ui'
import {h, onMounted, onUnmounted, ref, watch} from "vue";
import {useChainExecInfo} from "../../hooks/chain.js";
import DispatchRecord from "../../components/DispatchRecord.vue";
import {useAnimation} from "../../hooks/animation.js";
import {useDataTable} from "../../hooks/common.js";
import {formatTime} from "../../services/time.js";

const props = defineProps(['showNodes', 'chainId'])
const showHistory = ref(false)

const showNodeExecDetail = ref(false)
const curNodeExecDetail = ref(null)
const onClickShowNodeExecDetail = (row) => {
  curNodeExecDetail.value = row
  showNodeExecDetail.value = true
}


const showCurNodeDetail = ref(false)
const curNodeDetail = ref(null)
const onClickShowCurNodeDetail = (row) => {
  curNodeDetail.value = row
  showCurNodeDetail.value = true
}


const test = ref(null)
const his = ref(null)
const cur = ref(null)
const success = ref(null)
const fail = ref(null)


const {
  execHistory,
  dispatchInfo,
  GetDispatchInfo,
  NewDispatch,
  SingleStepDispatch,
  AllStepDispatch
} = useChainExecInfo(props.chainId)
const selectedIndex = ref(0)

const {SameTimeAnimate: fatherAnim, SetAnimationClass: SetFatherAnim} = useAnimation([test])


const {SameTimeAnimate: childAnim, SetAnimationClass: SetChildAnim} = useAnimation([his, cur])


const {tableData, columns, SetData} = useDataTable([
  {
    title: '节点名称',
    key: 'nodeName'
  },
  {
    title: '输出',
    key: 'out',
    width: 200,
    ellipsis: {
      tooltip: true
    }
  },
  {
    title: '状况',
    key: 'ok',
    render: (row) => h('span', {
      class: row.ok ? 'text-green-500' : 'text-red-500'
    }, {default: () => row.ok ? '成功' : '失败'})
  },
  {
    title: '指令名称',
    key: 'shortcutName'
  },
  {
    title: '指令',
    key: 'payload',
    ellipsis: {
      tooltip: true
    }
  },
  {
    title: '执行时间',
    key: 'createTime',
    width: 160,
    render: (row) => h('span', null, {
      default: () => formatTime(row.createTime)
    })
  },
  {
    title: '操作',
    key: 'action',
    width: 100,
    render: (row) => h(NButton, {
      size: 'small',
      onClick: () => onClickShowNodeExecDetail(row)
    }, ()=>'详情')
  }
], dispatchInfo.value.preNodes)

watch(() => dispatchInfo.value.preNodes, (newVal) => {
  SetData(newVal)
})


const selectDispatch = (index) => {
  if (index === selectedIndex.value) return
  fatherAnim('animate__zoomIn', 'animate__fast')
  childAnim('animate__zoomIn', 'animate__fast')
  selectedIndex.value = index
  GetDispatchInfo(execHistory.value[index].id)
}


const onClickSingleStep = () => {
  SingleStepDispatch()
}

const onClickAllStep = () => {
  AllStepDispatch()
}

let timer = null

onMounted(() => {

  fatherAnim('animate__slideInUp', 'animate__fast')

  timer = setInterval(() => {
    if (dispatchInfo.value.dispatchLog.id) {
      if (dispatchInfo.value.dispatchLog.done === true) return
      GetDispatchInfo(dispatchInfo.value.dispatchLog.id)
    }
  }, 5000)
})


onUnmounted(() => {
  clearInterval(timer)
})


defineExpose({
  fatherAnim,
  childAnim,
})


</script>

<template>

  <n-drawer v-model:show="showCurNodeDetail" placement="right" width="600">
    <n-drawer-content title="将要执行的NODE">
        <span class="block">NODE名称:
          <span class="block text-xl mb-3">{{ curNodeDetail.node.name }}</span>
        </span>
      <span class="block">描述:
        <span class="block text-xl mb-3">{{ curNodeDetail.node.description }}</span>
      </span>
      <span class="block">指令名称:
        <span class="block text-xl mb-3">{{ curNodeDetail.shortcut.name }}</span>
      </span>
      <span class="block">指令描述:
        <span class="block text-xl mb-3">{{ curNodeDetail.shortcut.description }}</span>
      </span>
      <span class="block">超时时间:
        <span class="block text-xl mb-3">{{ curNodeDetail.shortcut.timeout }} ms</span>
      </span>
      <span class="block">仅运行:
        <span class="block text-xl mb-3">{{ curNodeDetail.shortcut.justRun ? '是' : '否' }}</span>
      </span>
      <span class="block">目标机器:
        <span class="block text-xl mb-3">{{ curNodeDetail.shortcut.probeId }}</span>
      </span>
      <span class="block">指令:
        <span class="block text-xl text-blue-500 mb-3">{{ curNodeDetail.shortcut.payload }}</span>
      </span>
    </n-drawer-content>
  </n-drawer>

  <n-drawer v-model:show="showNodeExecDetail" placement="right" width="100vw">
    <n-drawer-content title="执行详情" closable>
      <div>
        <span class="text-2xl text-blue-500 block mb-1">输出：</span>
        <n-log :log="curNodeExecDetail.out">
        </n-log>
      </div>
      <div>
        <span class="text-xl  block mb-1">
          状况：<span :class="curNodeExecDetail.ok ? 'text-green-500' : 'text-red-500'">{{
            curNodeExecDetail.ok ? '成功' : '失败'
          }}</span>
        </span>
      </div>
      <div>
        <span class="text-xl block mb-1">名称：{{ curNodeExecDetail.nodeName }}</span>
      </div>
      <div>
        <span class="text-xl  block mb-1">执行时间：{{ formatTime(curNodeExecDetail.createTime) }}</span>
      </div>
      <div>
        <span class="text-xl  block mb-1">指令名称：{{ curNodeExecDetail.shortcutName }} </span>
      </div>
      <div>
        <span class="text-xl  block mb-1">目标机器：{{ curNodeExecDetail.probeId }}</span>
      </div>
      <div>
        <span class="text-xl  block mb-1">指令：{{ curNodeExecDetail.payload }}</span>
      </div>
    </n-drawer-content>
  </n-drawer>

  <div class="w-full h-full overflow-scroll  z-0">
    <div ref="test" class="console  h-full p-2 flex">

      <div v-show="execHistory.length !== 0" class="w-11/12 h-full p-5">

        <div ref="his"
             class="pre w-full h-72 flex  mb-5">
          <n-data-table :columns="columns" :data="tableData" max-height="18rem" ></n-data-table>
        </div>

        <div ref="cur" class="cur flex items-center justify-center">

          <div v-if="dispatchInfo.curNode" @click="onClickShowCurNodeDetail(dispatchInfo.curNode)"
               class="w-full h-full flex flex-col justify-between">
            <div class="h-10 bg-[#222226] w-full flex">
              <div class="table-header">
                名称
              </div>
              <div class="table-header">
                描述
              </div>
              <div class="table-header">
                指令名称
              </div>
              <div class="table-header">
                指令内容
              </div>
              <div class="table-header">
                仅运行
              </div>
              <div class="table-header">
                超时时间
              </div>
              <div class="table-header">
                目标机器
              </div>
            </div>
            <div class="grow  w-full flex flex-row hover:cursor-pointer hover:bg-gray-600 transition duration-300">
              <div class="w-full h-full flex items-center justify-start pl-1">
                {{ dispatchInfo.curNode.node.name }}
              </div>
              <div class="w-full h-full flex items-center justify-start pl-1 whitespace-nowrap text-ellipsis">
                {{ dispatchInfo.curNode.node.description }}
              </div>
              <div class="w-full h-full flex items-center justify-start pl-1">
                {{ dispatchInfo.curNode.shortcut.name }}
              </div>
              <div class="w-full h-full flex items-center justify-start pl-1 whitespace-nowrap text-ellipsis">
                {{ dispatchInfo.curNode.shortcut.payload }}
              </div>
              <div class="w-full h-full flex items-center justify-start pl-1">
                {{ dispatchInfo.curNode.shortcut.justRun ? '是' : '否' }}
              </div>
              <div class="w-full h-full flex items-center justify-start pl-1">
                {{ dispatchInfo.curNode.shortcut.timeout }} ms
              </div>
              <div class="w-full h-full flex items-center justify-start pl-1">
                {{ dispatchInfo.curNode.shortcut.probeId }}
              </div>
            </div>
          </div>
          <n-empty v-else description="没有要执行的指令">

          </n-empty>

        </div>

        <div class=" arrow w-full h-14 mb-1 flex items-center justify-around">
          <n-icon size="60" class="text-green-500">
            <ArrowDownCircle/>
          </n-icon>


          <n-icon size="60" class="text-red-500">
            <ArrowDownCircle/>
          </n-icon>
        </div>

        <div class="next w-full h-24 flex mb-5 items-center justify-between">
          <div ref="success" class="success overflow-hidden">
            <div v-if="dispatchInfo.successThen" @click="onClickShowCurNodeDetail(dispatchInfo.successThen)"
                 class="w-full h-full flex flex-col justify-between">
              <div class="h-10 bg-[#222226] w-full flex">
                <div class="table-header">
                  名称
                </div>
                <div class="table-header">
                  描述
                </div>
                <div class="table-header">
                  指令名称
                </div>
                <div class="table-header">
                  指令内容
                </div>
                <div class="table-header">
                  仅运行
                </div>
                <div class="table-header">
                  超时时间
                </div>
                <div class="table-header">
                  目标机器
                </div>
              </div>
              <div class="grow  w-full flex flex-row hover:cursor-pointer hover:bg-gray-600 transition duration-300">
                <div class="w-full h-full flex items-center justify-start pl-1">
                  {{ dispatchInfo.successThen.node.name }}
                </div>
                <div class="w-full h-full flex items-center justify-start pl-1 whitespace-nowrap text-ellipsis">
                  {{ dispatchInfo.successThen.node.description }}
                </div>
                <div class="w-full h-full flex items-center justify-start pl-1">
                  {{ dispatchInfo.successThen.shortcut.name }}
                </div>
                <div class="w-full h-full flex items-center justify-start pl-1 whitespace-nowrap text-ellipsis">
                  {{ dispatchInfo.successThen.shortcut.payload }}
                </div>
                <div class="w-full h-full flex items-center justify-start pl-1">
                  {{ dispatchInfo.successThen.shortcut.justRun ? '是' : '否' }}
                </div>
                <div class="w-full h-full flex items-center justify-start pl-1">
                  {{ dispatchInfo.successThen.shortcut.timeout }} ms
                </div>
                <div class="w-full h-full flex items-center justify-start pl-1">
                  {{ dispatchInfo.successThen.shortcut.probeId }}
                </div>
              </div>
            </div>

            <n-empty v-else description="上一节点成功后无操作">

            </n-empty>

          </div>
          <div ref="fail" class="fail">
            <div v-if="dispatchInfo.failThen" @click="onClickShowCurNodeDetail(dispatchInfo.failThen)"
                 class="w-full h-full flex flex-col justify-between">
              <div class="h-10 bg-[#222226] w-full flex">
                <div class="table-header">
                  名称
                </div>
                <div class="table-header">
                  描述
                </div>
                <div class="table-header">
                  指令名称
                </div>
                <div class="table-header">
                  指令内容
                </div>
                <div class="table-header">
                  仅运行
                </div>
                <div class="table-header">
                  超时时间
                </div>
                <div class="table-header">
                  目标机器
                </div>
              </div>
              <div class="grow  w-full flex flex-row hover:cursor-pointer hover:bg-gray-600 transition duration-300">
                <div class="w-full h-full flex items-center justify-start pl-1">
                  {{ dispatchInfo.failThen.node.name }}
                </div>
                <div class="w-full h-full flex items-center justify-start pl-1 whitespace-nowrap text-ellipsis">
                  {{ dispatchInfo.failThen.node.description }}
                </div>
                <div class="w-full h-full flex items-center justify-start pl-1">
                  {{ dispatchInfo.failThen.shortcut.name }}
                </div>
                <div class="w-full h-full flex items-center justify-start pl-1 whitespace-nowrap text-ellipsis">
                  {{ dispatchInfo.failThen.shortcut.payload }}
                </div>
                <div class="w-full h-full flex items-center justify-start pl-1">
                  {{ dispatchInfo.failThen.shortcut.justRun ? '是' : '否' }}
                </div>
                <div class="w-full h-full flex items-center justify-start pl-1">
                  {{ dispatchInfo.failThen.shortcut.timeout }} ms
                </div>
                <div class="w-full h-full flex items-center justify-start pl-1">
                  {{ dispatchInfo.failThen.shortcut.probeId }}
                </div>
              </div>
            </div>

            <n-empty v-else description="上一节点失败后无操作">

            </n-empty>
          </div>
        </div>

      </div>
      <div v-show="execHistory.length === 0"
           class="w-11/12 h-full p-5 flex items-center justify-center -translate-y-10 text-3xl animate-pulse">
        没有调度记录，点击新建按钮，创建一个调度器
      </div>
      <div class="w-1/12 h-full  flex flex-col items-center">
        <n-button class="mb-10 mt-5" size="large" circle @click="props.showNodes()">
          <n-icon size="30">
            <SettingsOutline/>
          </n-icon>
        </n-button>
        <n-button circle size="large" color="red" class="bg-red-400 mb-5">
          <n-icon size="30">
            <Remove/>
          </n-icon>
        </n-button>
        <n-button @click="()=>{
          NewDispatch();
          childAnim('animate__zoomIn','animate__fast');
        }" circle size="large" class="mb-5 bg-green-500">
          <n-icon size="30">
            <Add/>
          </n-icon>
        </n-button>
        <n-button @click="onClickAllStep" circle size="large" class="pl-1 mb-3">
          <n-icon size="30">
            <Play/>
          </n-icon>
        </n-button>
        <n-button @click="onClickSingleStep" circle size="large" class="">
          <n-icon size="30">
            <ReturnUpForward/>
          </n-icon>
        </n-button>
      </div>

    </div>


    <div
        class="history bg-gray-700 h-[600px] transition duration-300 shadow-2xl rounded w-72 absolute z-10 bottom-0 right-1 p-2 flex flex-col "
        :style="{
          transform: showHistory ? 'translateY(0)' : 'translateY(88%)'
        }"
    >
      <div @click="showHistory = !showHistory"
           class="w-full h-6 flex items-center justify-center group hover:cursor-pointer"
      >
        <n-icon size="30" class="group-hover:text-blue-500 transition" :style="{
          transform: showHistory ? 'rotate(0deg)' : 'rotate(180deg)'
        }">
          <CaretDown/>
        </n-icon>

      </div>

      <div class="grow  p-2 overflow-scroll">
        <dispatch-record
            :item="item"
            :selected="selectedIndex === index"
            @click="selectDispatch(index)"
            v-for="(item,index) in execHistory"
            :key="index"
            class="mb-2"
        >
        </dispatch-record>
      </div>
    </div>
  </div>


</template>

<style scoped>


.success {
  @apply w-[48%] h-full  transition flex-center border-gray-700 border-[1px] rounded;
}

.flex-center {
  @apply flex items-center justify-center
}

.fail {
  @apply w-[48%] h-full flex-center border-gray-700 border-[1px] rounded;
}

.cur {
  @apply w-full h-24 border-gray-700 border-[1px] rounded mb-1 transition;
}


.table-header {
  @apply w-full flex items-center justify-start pl-1;
}


</style>