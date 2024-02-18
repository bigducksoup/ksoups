<script setup>
import { NButton, NIcon } from "naive-ui";
import { Remove,Add,Play,ReturnUpForward } from "@vicons/ionicons5";
</script>

<template>
  <div class="option-bar">
    <n-button circle size="large" color="red" class="bg-red-400 hover:after:content-['终止调度'] after:text-red-500 tip">
      <n-icon size="30">
        <Remove />
      </n-icon>
    </n-button>
    <n-button circle size="large" class="hover:after:content-['新建调度'] tip">
      <n-icon size="30">
        <Add />
      </n-icon>
    </n-button>
    <n-button circle size="large" class="pl-1 hover:after:content-['执行全部'] tip">
      <n-icon size="30">
        <Play />
      </n-icon>
    </n-button>
    <n-button  circle size="large" class="hover:after:content-['下一步'] tip">
      <n-icon size="30">
        <ReturnUpForward />
      </n-icon>
    </n-button>
  </div>
</template>

<style scoped>
.option-bar {
  @apply min-w-fit h-fit p-1  flex flex-col items-center border-[1px] border-gray-500 rounded gap-1 justify-end text-white;
}



.tip {
 @apply after:absolute after:z-20 after:left-full;
}

</style>
