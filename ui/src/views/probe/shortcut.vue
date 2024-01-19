<template>


<div ref="root" class="w-full  p-2">
  <div class="ml-2 text-3xl text-[#BFCFE7]">
    快捷指令
  </div>
  <div class="w-full flex flex-row flex-wrap text-white">
    <ShortCutCard
        v-for="item in shortcuts"
        @onClickPlay="runOLShortcut(item)"
        @onClickMore="shortcutDetail(item)"
        :title="item['name']"
        :content="item['payload']"
        class="w-56  mt-2 ml-2 aspect-[6/4]"
    >
    </ShortCutCard>
    <div
        @click="createSCShow = true"
        class="w-56 rounded-sm mt-2 ml-2 aspect-[6/4] hover:cursor-pointer  border-[#303438] border-[1px] flex justify-center items-center"
    >
      <n-icon size="80">
        <Add/>
      </n-icon>
    </div>
  </div>




  <Transition name="bounce">
    <OneLineForm
        v-if="createSCShow"
        @on-click-close="createSCShow=false"
        @on-success="init"
        :probe-id="route.params.probeId"
        class="absolute z-10 left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2"
    >
    </OneLineForm>
  </Transition>

  <Transition name="bounce">
    <ShortCutDetail
      v-show="shortcutDetailShow"
      @onClickClose="shortcutDetailShow=false"
      @onDelete="init"
      class="absolute z-10 left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2"
      :shortcut="detailShortcut"
    >

    </ShortCutDetail>
  </Transition>



</div>


</template>


<script setup>
import {useRoute} from 'vue-router'
import {onBeforeUpdate, onMounted, onUpdated, ref, toRef} from "vue";

import {NIcon,useNotification} from 'naive-ui'
import {Add} from '@vicons/ionicons5'

import OneLineForm from "./shortcut/one-line-form.vue";
import ShortCutCard from "../../components/ShortcutCard.vue";

import {listShortcut,runShortcut} from '../../services/shortcut.js'
import ShortCutDetail from "./shortcut/short-cut-detail.vue";
import CreateForm from "./file/create-form.vue";
import {useAnimation} from "../../hooks/animation.js";




const route = useRoute()
const notification = useNotification()

const shortcuts = ref([])

const createSCShow = ref(false)
const shortcutDetailShow = ref(false)

const detailShortcut = ref({})


const root = ref(null)

const { SameTimeAnimate:RootAnim } = useAnimation([root])


const init = async() => {
  let res  = await listShortcut(route.params.probeId)
  shortcuts.value = res['data']
}



const runOLShortcut = async(shortcut)=>{
  // let date = new Date()

  let res = await runShortcut(shortcut['id'])
  let type = res['ok'] ? 'success' : 'error'

  notification[type]({
    title:shortcut.Name,
    description: res['ok'] ? "执行成功" : "执行失败",
    content: res['out'],
    meta: new Date().toLocaleString()
  })
}


const shortcutDetail = (shortcut)=>{
  shortcutDetailShow.value = true
  detailShortcut.value = shortcut
}

onMounted(()=>{
  RootAnim('animate__slideInUp','animate__fast')
  init()
})



</script>

<style scoped>
.bounce-enter-active {
  animation: bounce-in 0.5s;
}
.bounce-leave-active {
  animation: bounce-in 0.5s reverse;
}

</style>