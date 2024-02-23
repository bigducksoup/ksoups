<script setup>
import {NCode,NButton,NModal} from 'naive-ui'

import {formatTime} from "../../../services/time.js";
import {onMounted, ref, watch} from "vue";
import {useFSOperation} from "../../../hooks/file.js";
const props = defineProps(['node'])
const emits = defineEmits(['onSetRoot'])

const code = ref('')
const shortcut = ref({})

const { GetFileContent } = useFSOperation()

const loadCode = async () => {

  if(!props.node)return

  shortcut.value = props.node.getData().proto

  if(shortcut.value.type === 0){
    code.value = shortcut.value.payload
    return
  }

  if (shortcut.value.type === 1){
    code.value = await GetFileContent(shortcut.value.probeId,shortcut.value.payload)
  }
}


const showContentModal = ref(false)

const onClickSetRoot = () => {
  emits('onSetRoot',props.node)
  props.node.setData({
    root:true
  })
}

watch(()=>props.node,()=>{
  loadCode()
})

onMounted(()=>{
  loadCode()
})

</script>

<template>
  <div v-if="props.node === null" class="flex w-full h-full items-center justify-center">
    点击节点以查看详情
  </div>
  <div v-else class="flex flex-col w-full ml-1 mr-1 ">
    <div>{{ shortcut.name }}</div>
    <div class="text-xs text-gray-600 mb-2 h-6 overflow-ellipsis w-full">{{ shortcut.description }}</div>
    <div class="w-full mt-auto">
      <n-button @click="onClickSetRoot">
        设为根节点
      </n-button>
    </div>
    <div class="mt-2 text-blue-500 mb-1"> <span class="text-white">创建时间:</span> {{ formatTime(shortcut.createTime) }}</div>
    <div class="mb-1"> 超时时间: <span class="text-blue-500">{{ shortcut.timeout }}ms</span> </div>
    <div class="mb-1"> 仅运行: <span class="text-blue-500">{{ shortcut.justRun ? '是' : '否' }} </span></div>
    <div><span>Code:</span></div>

    <n-modal v-model:show="showContentModal">
      <div class=" w-[95%] h-[700px] sm:w-[800px] rounded p-2 bg-gray-800 overflow-scroll">
        <n-code :code="code"  class="text-green-500 w-full h-full">
        </n-code>
      </div>
    </n-modal>

    <div class="w-full h-[100px]  overflow-hidden">
      <n-code :code="code"  class="text-green-500 w-full">
      </n-code>
    </div>

    <n-button class="w-[95%] mt-1" @click="showContentModal = true" size="small"> 查看完整内容 </n-button>


  </div>


</template>

<style scoped>

</style>