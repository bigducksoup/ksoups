<script setup>
import {Add, ArrowDownCircle, CaretDown, Play, Remove, ReturnUpForward, SettingsOutline} from '@vicons/ionicons5'
import {NButton, NIcon} from 'naive-ui'
import {onMounted, ref} from "vue";
import {useChainExecInfo} from "../../hooks/chain.js";
import DispatchRecord from "../../components/DispatchRecord.vue";
import NodeExecResult from "../../components/NodeExecResult.vue";
import {useAnimation} from "../../hooks/animation.js";

const props = defineProps(['showNodes', 'chainId'])
const showHistory = ref(false)

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

const {SameTimeAnimate} = useAnimation([test])


const {SameTimeAnimate: childAnim} = useAnimation([his, cur])

const selectDispatch = (index) => {
  if (index === selectedIndex.value) return
  SameTimeAnimate('animate__zoomIn', 'animate__fast')
  childAnim('animate__zoomIn', 'animate__fast')
  selectedIndex.value = index
  GetDispatchInfo(execHistory.value[index].id)
}


onMounted(() => {

  SameTimeAnimate('animate__zoomIn', 'animate__fast')
  childAnim('animate__zoomIn', 'animate__fast')


  setInterval(() => {
    if (dispatchInfo.value.dispatchLog.id) {
      if (dispatchInfo.value.dispatchLog.done === true) return
      GetDispatchInfo(dispatchInfo.value.dispatchLog.id)
    }
  }, 5000)

})


</script>

<template>

  <div class="w-full h-full overflow-scroll  z-0">
    <div ref="test" class="console  h-full p-2 flex">

      <div v-show="execHistory.length !== 0" class="w-11/12 h-full p-5">

        <div ref="his"
             class="pre w-full h-72 flex  mb-5  border-2 border-gray-700   overflow-x-scroll overflow-y-hidden p-2 gap-2">
          <div class="w-96 h-full shrink-0" v-for="item in dispatchInfo.preNodes">
            <node-exec-result
                class="w-full h-full"
                :item="item"
                :key="item.id"
            >

            </node-exec-result>
          </div>
        </div>

        <div ref="cur" class="cur flex items-center justify-center">

          <div v-if="dispatchInfo.curNode" class="w-full h-full flex p-2">
            <div v-if="dispatchInfo.curNode.shortcut" class="w-4/6 h-full">
              <div class="w-full text-xl">
                将在 {{ dispatchInfo.curNode.shortcut.probeId }} 上执行:
              </div>
              <div class="font-bold text-base text-blue-500">
                {{ dispatchInfo.curNode.shortcut.payload }}
              </div>
            </div>
            <div class="w-2/6 h-full text-xl flex flex-col items-end">
               <span class="mt-auto"> 超时时间: {{ dispatchInfo.curNode.shortcut.timeout }}ms</span>
            </div>



          </div>
          <span v-else class=" text-3xl">
            没有要执行的Node
          </span>

        </div>

        <div class=" arrow w-full h-20 mb-1 flex items-center justify-around">
          <n-icon size="80" class="text-green-500">
            <ArrowDownCircle/>
          </n-icon>


          <n-icon size="80" class="text-red-500">
            <ArrowDownCircle/>
          </n-icon>
        </div>

        <div class="next w-full h-36 flex mb-5 items-center justify-between">
          <div ref="success" class="success">

            <div v-if="dispatchInfo.successThen" class="h-full w-full p-2 flex flex-row">
              <div class="flex flex-col w-2/6 ">
                <span>名称:{{ dispatchInfo.successThen.node.name }}</span>
                <span>描述:{{ dispatchInfo.successThen.node.description }}</span>
              </div>
              <div class="w-4/6 flex flex-col" v-if="dispatchInfo.successThen.shortcut">
                <span>在 {{ dispatchInfo.successThen.shortcut.probeId }}上执行</span>
                <span>payload:
                  <br>
                  <span class="text-blue-500">{{ dispatchInfo.successThen.shortcut.payload }}</span>
                </span>
              </div>

            </div>

            <span v-else class="text-2xl">
              没有可能执行的NODE
            </span>

          </div>
          <div ref="fail" class="fail">
            <div v-if="dispatchInfo.failThen" class="h-full w-full p-2 flex flex-row">
              <div class="flex flex-col w-2/6 ">
                <span>名称:{{ dispatchInfo.failThen.node.name }}</span>
                <span>描述:{{ dispatchInfo.failThen.node.description }}</span>
              </div>
              <div class="w-4/6 flex flex-col" v-if="dispatchInfo.successThen.shortcut">
                <span>在 {{ dispatchInfo.failThen.shortcut.probeId }}上执行</span>
                <span>payload:
                  <br>
                  <span class="text-blue-500">{{ dispatchInfo.failThen.shortcut.payload }}</span>
                </span>
              </div>
            </div>

            <span v-else class="text-2xl">
              没有可能执行的NODE
            </span>
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
          SameTimeAnimate('animate__zoomIn','animate__fast');
          childAnim('animate__zoomIn','animate__fast');
        }" circle size="large" class="mb-5 bg-green-500">
          <n-icon size="30">
            <Add/>
          </n-icon>
        </n-button>
        <n-button @click="()=>{
          AllStepDispatch();
          SameTimeAnimate('animate__bounce','animate__fast');
        }" circle size="large" class="pl-1 mb-3">
          <n-icon size="30">
            <Play/>
          </n-icon>
        </n-button>
        <n-button @click="()=>{
          SingleStepDispatch();
          childAnim('animate__bounce','animate__fast');
        }" circle size="large" class="">
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
  @apply w-[48%] h-4/5  transition flex-center border-[2px] border-green-500 rounded;
}

.flex-center {
  @apply flex items-center justify-center
}

.fail {
  @apply w-[48%] h-4/5  flex-center border-[2px] border-red-500 rounded;
}

.cur {
  @apply w-full h-40 border-white border-2 rounded mb-1 transition;
}


</style>