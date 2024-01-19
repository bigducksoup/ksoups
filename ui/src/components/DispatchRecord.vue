<script setup>
import {formatTime} from '../services/time.js'

const props = defineProps({
  item: {
    type: Object,
    required: true
  },
  selected: {
    type: Boolean,
    default: false
  }
})


const statusMap = {
  0: '执行中',
  1: '执行完成',
  2: '执行失败',
  3: '放弃执行',
}

</script>

<template>

  <div :class="selected ? 'selected':'normal'">
    <div>
      执行日期 : {{ formatTime(item.createTime) }}
    </div>
    <div class="grow flex">
      <div class="w-1/2">
        是否完成：{{ item.done ? '✅' : '❌' }}
      </div>
      <div class="w-1/2">
        状态 : {{ statusMap[item.status] }}
      </div>
    </div>
  </div>

</template>

<style scoped>

.selected {
  @apply public rounded border-green-500 hover:cursor-pointer
}

.public {
  @apply w-full h-20 flex-col flex p-1 border-[1px]
}

.normal {
  @apply public hover:bg-gray-600 border-transparent hover:border-green-500 hover:cursor-pointer;
}

</style>