<script setup>
import { AtCircleOutline, Terminal, Cloud, EyeOffOutline } from "@vicons/ionicons5";
import {NButton, NIcon} from "naive-ui";

/**
 * info <info>
 * { "id": "f51296cf-bcfc-4b7b-accb-aeb0ee0d7f0c", "addrPort": "127.0.0.1:22", "username": "ducksoup", "password": "******", "groupId": "root" }
 * @type {Prettify<Readonly<{[key in string]?: any}>>}
 */
const props = defineProps(["info"]);
</script>

<template>
  <div
    class="hover:cursor-pointer text-[#0078D7] w-full min-w-max h-10 p-1 flex items-center border-b-[1px] border-b-gray-600 hover:bg-gray-800"
  >
    <n-icon size="20" class="mr-2 text-2xl">
      <Terminal />
    </n-icon>

    <span class="single">
      <n-icon size="20" class="mr-2 text-2xl">
        <Cloud />
      </n-icon>
      {{ info.addrPort }}
    </span>

    <span class="single">
      <n-icon size="20" class="mr-2 text-2xl">
        <AtCircleOutline />
      </n-icon>
      {{ info.username }}
    </span>

    <span class="single">
      <n-icon size="20" class="mr-2 text-2xl">
        <EyeOffOutline />
      </n-icon>
      {{ info.password }}
    </span>


    <n-button size="small" class="ml-auto"  > 详情 </n-button>

  </div>
</template>

<style scoped>
.single {
  @apply ml-20
  w-36
  flex items-center;
}
</style>
