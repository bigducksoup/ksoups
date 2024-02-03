<script setup>
import 'xterm/css/xterm.css'
import {NIcon, NMenu} from 'naive-ui'
import {h, onMounted, ref} from "vue";
import {useRouter} from "vue-router";
import {useMenu} from "../hooks/common.js";
import {List,TerminalOutline} from "@vicons/ionicons5";

const router = useRouter()

const { MenuOption,Key ,onUpdate, } = useMenu([
  {
    label:'主机',
    key: 'host',
    icon: () => h(NIcon,null,{
      default : () => h(List)
    })
  },
  // {
  //   label: '终端',
  //   key: 'term',
  //   icon: () => h(NIcon,null,{
  //     default : () => h(TerminalOutline)
  //   })
  // }
],(key, item)=>{
    Key.value = key
    router.push('/ssh/' + key)
})

const header = ref(null)


onMounted(() => {
  header.value.classList.add('anim')
  onUpdate('host',null)
})


const leave = () => {
  header.value.classList.remove('anim')
  requestAnimationFrame(() => {
    header.value.classList.add('anim-rev')
  })
  setTimeout(() => {
    router.push('/')
  }, 650)
}


</script>

<template>

  <div ref="root" id="root-term" class="w-full h-full">
    <div ref="header" class="flex items-center pl-2">
      <button @click="leave" class="hover:underline hover:text-white">
        返回首页
      </button>
    </div>
    <div class="h-fit w-full bg-sky-700 flex">
      <div class="md:w-0 w-0 bg-[#13151D]   h-screen overflow-hidden">
        <n-menu
            :options="MenuOption"
            class="w-full"
            :on-update-value="onUpdate"
            v-model:value="Key"
        >
        </n-menu>
      </div>
      <div class="grow min-h-screen overflow-x-scroll">
        <RouterView v-slot="{ Component }">
          <keep-alive>
            <component ref="detail" :key="$route.name" :is="Component"/>
          </keep-alive>
        </RouterView>
      </div>
    </div>
  </div>

</template>

<style scoped>


.anim {
  animation: full-to-normal 0.7s ease-in-out forwards;
}

.anim-rev {
  animation: full-to-normal 0.7s ease-in-out reverse;
}

@keyframes full-to-normal {
  0% {
    height: 100vh;
    background-color: #05253A;
    color: transparent;
  }

  30% {
    background-color: #05253A;
    color: transparent;
  }

  80% {
    height: 2.2rem;
    --tw-bg-opacity: 1;
    background-color: rgb(34 197 94 / var(--tw-bg-opacity));
    color: transparent;
  }
  100% {
    height: 2.5rem;
    --tw-bg-opacity: 1;
    background-color: rgb(34 197 94 / var(--tw-bg-opacity));
  }
}

</style>