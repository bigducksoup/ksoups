<script setup>
import {NList,NListItem} from  'naive-ui'
import {formatTime} from "../../../services/time.js";


const props = defineProps(['data'])
const emits = defineEmits(['onClickItem'])



</script>

<template>
  <div class="w-full h-full  flex flex-col overflow-scroll">
    <n-list clickable hoverable>
      <n-list-item  @click="() => {
        emits('onClickItem',item)
      }" v-for="item in data">
        执行时间： {{ formatTime(item.createTime) }}
        执行状态： {{ item.status }}
      </n-list-item>
    </n-list>
  </div>
</template>

<style scoped>

</style>