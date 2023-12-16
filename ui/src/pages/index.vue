<template>
  <div class="layout flex flex-col w-screen h-screen overflow-hidden">
    <div class="layout-top">
      <div class="bar w-full h-[50px] flex justify-center items-center bg-[#23272E] border-b-2 border-[#1E2228]">
        <span class="text-green-500"> 配置中心 </span>
      </div>
    </div>

    <div v-if="nodeArr.length != 0" class="layout-mid flex flex-row overflow-scroll grow bg-[#23272E]">

      <div class="layout-left  w-60 flex flex-col items-center">

        <div v-for="(item, index) in nodeArr" :class="item['addr'] == selectAddress ? seletedStyle : unselectStyle"
          @click="selectNode(item['addr'])">
          {{ item['addr'] }}
        </div>



      </div>

      <div  class="layout-right flex grow">


        <router-view   v-slot="{ Component }">
          <keep-alive>
            <component :key="$route.params.addr" :is="Component" />
          </keep-alive>
        </router-view>




      </div>
    </div>

    <div v-else class="grow bg-[#23272E] flex items-center justify-center">
      <span class="text-red-500">暂无在线节点</span>
    </div>

    <div class="layout-footer flex justify-center items-center text-white bg-[#23272E] border-t-2 border-[#1E2228]">
      <!-- powered by native-ui, Vue3, TailWindCss -->
    </div>
  </div>
</template>




<script setup>
import { ref } from 'vue'
import { onMounted } from 'vue';
import { baseUrl } from '../state'

import { useRouter } from 'vue-router'

//当前选中的node地址
const selectAddress = ref('')

//在线的node
const nodeArr = ref([])





const unselectStyle = "w-[95%] m-2 rounded-md p-2 border-2 border-transparent hover:border-2 hover:border-black bg-red-500"

const seletedStyle = "w-[95%] m-2 rounded-md p-2 border-2 border-green-500 bg-green-500"

const router = useRouter()

// 选择node
const selectNode = (addr) => {
  if (selectAddress.value == addr) return
  selectAddress.value = addr

  router.push("/" + addr)

}


onMounted(() => {


  fetch(baseUrl.value + "api/info/nodes").then(res => {
    return res.json()
  }).then(json => {
    nodeArr.value = json
    let defaultAddr = nodeArr.value[0]['addr']
    selectNode(defaultAddr)
  })
})


</script>