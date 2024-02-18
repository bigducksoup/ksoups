<script setup>
import {inject, onMounted, ref} from "vue";
import {Cog, RefreshCircle,CloseCircle} from '@vicons/ionicons5'
import {NIcon} from 'naive-ui'

const state = ref(0)

const getNode = inject('getNode',)
const getGraph = inject('getGraph')

const node = ref(null)
const graph = ref(null)

const data = ref({
  proto: {
    name: 'loading'
  }
})

const onClickClose = () => {
    graph.value.trigger('node:onClickDeleteBtn',node.value)
}

const onClickNode = () => {
  graph.value.trigger('node:onClickNode',node.value)
}

onMounted(() => {
  node.value = getNode()
  graph.value = getGraph()
  data.value = node.value.getData()
});
</script>

<template>
  <div @click="onClickNode" class="w-36 h-14 bg-blue-500 rounded flex items-center justify-center overflow-hidden">
    <div :class="`state-${state}`">

    </div>
    <n-icon size="20" class="ml-2 ">
      <Cog/>
    </n-icon>
    <span class="ml-1"> {{ data.proto.name }} </span>
    <span v-if="state === 0" class="relative flex h-3 w-3 ml-auto mr-2">
      <span class="absolute inline-flex h-full w-full rounded-full bg-white opacity-75"></span>
      <span class="relative inline-flex rounded-full h-3 w-3 bg-white"></span>
    </span>

    <span v-if="state === 1" class="relative flex h-3 w-3 ml-auto mr-2">
      <span class="absolute inline-flex h-full w-full rounded-full bg-green-500 opacity-75"></span>
      <span class="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
    </span>

    <span v-if="state === 2" class="relative flex h-3 w-3 ml-auto mr-2">
      <span class="absolute animate-pulse inline-flex h-full w-full rounded-full bg-red-500 opacity-75"></span>
      <span class="relative animate-pulse inline-flex rounded-full h-3 w-3 bg-red-500"></span>
    </span>

    <n-icon  v-if="state === 3" size="20" class="relative flex ml-auto mr-2 text-yellow-500">
      <RefreshCircle class="animate-spin"/>
    </n-icon>


    <n-icon @click="onClickClose" size="20" class="absolute right-0 top-0 translate-x-1/2 -translate-y-1/2 text-red-500 hover:cursor-pointer">
      <CloseCircle/>
    </n-icon>

  </div>
</template>


<style scoped>

.state-0 {
  @apply bg-white mark
}

.state-1 {
  @apply bg-green-500 mark
}

.state-2 {
  @apply bg-red-500 mark
}

.state-3 {
  @apply bg-yellow-500 mark
}


.mark {
  @apply w-1 h-full absolute left-0
}




</style>