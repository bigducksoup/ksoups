<script setup>
import Code from  '../../../components/code.vue'
import {NButton,NPopconfirm} from 'naive-ui'
import {ref, toRaw, watch, watchEffect} from "vue";
import {deleteShortcut} from "../../../services/shortcut.js";


const props = defineProps(['shortcut'])
const emit = defineEmits(['onClickClose','onDelete'])
const code = ref('')

const sc = ref({})

watch(props,(newProps)=>{
  sc.value = newProps.shortcut
  setCode()
})

const cancel = ()=>{
  emit('onClickClose')
}

const deleteSC = async ()=>{
  let res = await deleteShortcut(sc.value['id'])
  if (res['code'] !== 200){
    return
  }
  emit('onDelete')
  cancel()
}


const setCode = ()=>{

  if (sc.value['type'] === 0){
    code.value = sc.value['payload']
    return
  }

  if (sc.value['type'] === 1){
    code.value = 'should be script content'
  }

}

</script>

<template>
  <div class="w-screen h-screen bg-transparent flex items-center  justify-center" @click="cancel">
    <div @click.stop class="md:w-[32rem] md:h-[40rem] w-10/12 h-5/6 flex flex-col overflow-hidden bg-[#282828] text-white rounded-3xl">
      <div class="w-full bg-green-500 h-12">
          <h4 class="text-center text-xl font-medium pt-2">快捷指令详情</h4>
      </div>
      <div class="w-full h-full p-4 overflow-auto flex-col flex">
        <input v-model="sc.name" class="mb-2 bg-[#282828] text-xl focus:outline-none " />
        <input v-model="sc.description" class="bg-[#282828] mb-1 focus:outline-none" />
        <div v-if="sc.type === 1" class="text-xs mb-2 font-light text-gray-500">
          文件路径:{{sc.payload}}
        </div>
        <Code v-model="code" height="400px"></Code>
        <div class="grow bg-[#282C34]"></div>
        <div class="h-14 flex items-center flex-row justify-center">
          <n-popconfirm
            @positive-click="deleteSC"
            positive-text="确认"
            negative-text="取消"
          >
            <template #trigger>
              <n-button color="red" class="bg-red-500 mr-2" >删除指令</n-button>
            </template>
            <span>是否确认删除该指令?</span>
          </n-popconfirm>
          <n-button color="yellow" class="bg-yellow-500 mr-2">保存变更</n-button>
        </div>
      </div>
    </div>

  </div>

</template>

<style scoped>

</style>