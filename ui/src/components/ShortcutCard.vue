<script setup>
import {NIcon,NPopconfirm} from 'naive-ui'
import {EllipsisVertical,Play} from '@vicons/ionicons5'
import {ref} from "vue";
const props = defineProps({
  'title':{
    type: String,
    default : 'This place should be a title'
  },
  'content':{
    type:String,
    default:'This place should be a command like "docker ps -a | grep redis"'
  }
})
const emits = defineEmits(['onClickPlay','onClickMore'])

const root = ref(null)

defineExpose({
  root
})

</script>

<template>
  <div ref="root" class="shadow-sm rounded-sm overflow-hidden flex flex-col border-[#303438] border-[1px] hover:border-blue-500">
    <div class="pl-2 pt-2 mb-2 overflow-ellipsis text-xl text-white w-full truncate">{{props.title}}</div>
    <code class="grow pl-2 text-xs text-blue-500 overflow-scroll">
      {{props.content}}
    </code>
    <div class="h-9 w-full bg-[#26262A] flex flex-row-reverse items-center pr-2">
      <button @click="emits('onClickMore')" class="w-7 h-7 ml-1 rounded-full flex hover:text-green-500 items-center justify-center">
        <n-icon size='20' >
          <EllipsisVertical/>
        </n-icon>
      </button>
      <n-popconfirm
          @positive-click="emits('onClickPlay')"
          positive-text="确认"
          negative-text="取消"
      >
        <template #trigger>
          <button  class="w-7 h-7 bg-gray-600 hover:bg-blue-600 rounded-full pl-0.5 flex items-center justify-center">
            <n-icon size='20' >
              <Play/>
            </n-icon>
          </button>
        </template>
        你可想好了嗷！
      </n-popconfirm>
    </div>
  </div>

</template>

<style scoped>

</style>