<script setup>
import {NButton, NCode, NModal} from 'naive-ui'
import {formatTime} from "../../../services/time.js";
import {onMounted, ref, watch} from "vue";
import {useFSOperation} from "../../../hooks/file.js";

const props = defineProps(['data'])

const code = ref('')

const { GetFileContent } = useFSOperation()

const loadCode = async () => {


  if(props.data.shortcutType === 0){
    code.value = props.data.payload
    return
  }

  if (props.data.shortcutType === 1){
    code.value = await GetFileContent(props.data.probeId,props.data.payload)
  }
}


watch(()=>props.data,()=>{
  loadCode()
})

onMounted(()=>{
  loadCode()
})




const showContent = ref('')
const showModal = ref(false)


</script>

<template>
  <div class="w-full h-full p-1">
      <div v-if="props.data == null" class="h-full w-full flex items-center justify-center">
        点击节点以查看执行状况
      </div>
      <div v-else>
          <div>
            名称 : {{ props.data.nodeName }}
          </div>
        <div>
          执行时间 : {{ formatTime(props.data.createTime) }}
        </div>
        <div>
          状态 : <span :class="props.data.ok ? 'text-green-500' : 'text-red-500'"> {{ props.data.ok ? '执行成功' : '执行失败' }} </span>
        </div>
        <div>
          代码: <n-button size="small" @click="() => {showContent = code;showModal = true}"> 查看代码 </n-button>
        </div>
        <div class="w-full line-clamp-5 text-ellipsis whitespace-wrap overflow-hidden text-blue-500">
          <pre>{{ code }}</pre>
        </div>
        <div class="mt-5">
          标准输出: <n-button size="small" @click="() => {showContent = props.data.stdOut;showModal = true}"> 查看标准输出 </n-button>
        </div>
        <div class="w-full line-clamp-5 text-ellipsis whitespace-wrap overflow-hidden text-green-400">
          {{ props.data.stdOut }}
        </div>
        <div class="mt-5">
          错误输出: <n-button size="small" @click="() => {showContent = props.data.stdErr;showModal = true}"> 查看错误输出 </n-button>
        </div>
        <div class="mt-5 w-full line-clamp-5 text-ellipsis whitespace-wrap overflow-hidden">
          {{ props.data.stdErr }}
        </div>
      </div>


    <n-modal v-model:show="showModal">
      <div class=" w-[95%] h-[800px] sm:w-[1000px] rounded p-2 bg-gray-800 overflow-scroll">
        <n-code :code="showContent"  class="text-green-500 w-full h-full">
        </n-code>
      </div>
    </n-modal>

  </div>
</template>

<style scoped>

</style>