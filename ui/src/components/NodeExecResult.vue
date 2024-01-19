<script setup>
import {NPopover} from "naive-ui";
import {formatTime} from "../services/time.js";
import {toHtml} from "../services/utils.js";
// {
//   "createTime": "2024-01-16T15:20:41.399683+08:00",
//     "ok": true,
//     "out": "CO",
//     "oodeName": "root",
//     "shortcutName": "8081端口占用",
//     "payload": "lsof -i :8081",
//     "shortcutType": 0
// }
const props = defineProps(['item'])




</script>

<template>
   <div class="public bg-gray-800">
     <div class="text-xl font-bold mb-4 flex">
        <span class="mr-auto">节点名称:{{ item.nodeName }}</span>
       <span>
         {{ formatTime(item.createTime)}}
       </span>
     </div>

     <div class="text-base mb-1">
       快捷指令: <span class="mr-5">{{ item.shortcutName }}</span>
     </div>

     <div class="mb-1 flex-col h-15 w-full">
       命令 :
       <br>
       <span class="text-blue-500">{{ item.payload }}</span>
     </div>



     <n-popover trigger="hover" style="width: fit-content;max-width: 1400px;height: fit-content;max-height: 700px" scrollable placement="bottom">
       <template #trigger>
         <div :class="item.ok ? 'green':'red'">
           <span class="text-white">
             output:
           </span>
           <br>
           {{ item.out }}
         </div>
       </template>
         <div :class="item.ok ? 'text-green-500':'text-red-500'" v-html="toHtml(item.out)"></div>
      </n-popover>

   </div>
</template>

<style scoped>

.red {
  @apply text-red-500 w-full h-20 text-ellipsis overflow-auto;
}

.public {
  @apply flex flex-col p-1 rounded;
}

.green {
  @apply text-green-500 w-full h-20 text-ellipsis overflow-auto;
}

</style>