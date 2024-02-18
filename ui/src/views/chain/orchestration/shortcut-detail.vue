<script setup>
import {NCode,NButton,NModal} from 'naive-ui'

import {formatTime} from "../../../services/time.js";
import {onMounted, ref, watch} from "vue";
import {useFSOperation} from "../../../hooks/file.js";
const props = defineProps(['shortcut'])

const code = ref('')

const { GetFileContent } = useFSOperation()

const loadCode = async () => {

  if(!props.shortcut)return

  if(props.shortcut.type === 0){
    code.value = props.shortcut.payload
    return
  }

  if (props.shortcut.type === 1){
    code.value = await GetFileContent(props.shortcut.probeId,props.shortcut.payload)
  }
}


const showContentModal = ref(false)

watch(()=>props.shortcut,()=>{
  loadCode()
})

onMounted(()=>{
  loadCode()
})

</script>

<template>
  <div v-if="props.shortcut === null" class="flex w-full h-full items-center justify-center">
    点击节点以查看详情
  </div>
  <div v-else class="flex flex-col w-full ml-1 mr-1 overflow-scroll">
    <div>{{ shortcut.name }}</div>
    <div class="text-xs text-gray-600 mb-5">{{ shortcut.description }}</div>
    <div class="mt-2 text-blue-500 mb-1"> <span class="text-white">创建时间:</span> {{ formatTime(shortcut.createTime) }}</div>
    <div class="mb-1"> 超时时间: <span class="text-blue-500">{{ shortcut.timeout }}ms</span> </div>
    <div class="mb-1"> 仅运行: <span class="text-blue-500">{{ shortcut.justRun ? '是' : '否' }} </span></div>
    <div><span>Code:</span> <span> <n-button v-if="shortcut.type === 1" @click="showContentModal = true" size="small"> 查看 </n-button> </span></div>

    <n-modal v-model:show="showContentModal">
      <div class=" w-[95%] h-[700px] sm:w-[800px] rounded p-2 bg-gray-800 overflow-scroll">
        <n-code :code="code"  class="text-green-500 w-full h-full">
        </n-code>
      </div>
    </n-modal>

    <div class="w-full h-[550px] overflow-x-scroll">
      <n-code :code="code"  class="text-green-500 w-full">
      </n-code>
    </div>
  </div>


</template>

<style scoped>

</style>