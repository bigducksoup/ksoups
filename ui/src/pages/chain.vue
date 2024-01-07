<script setup>
import {useRouter} from 'vue-router'
import {onMounted, ref} from "vue";
import { getChainList } from '../services/chain.js'


const router = useRouter()

const headerStyle = ref({
  animation : 'full-to-normal 0.70s ease-in-out',
})


const chainList = ref([])


const backToMain = () => {
  headerStyle.value.animation = 'none'
  requestAnimationFrame(() => {
    headerStyle.value.animation = 'full-to-normal 0.70s ease-in-out reverse'
  })
  setTimeout(() => {
    router.push('/')
  }, 650)
}


const init = async () => {
  let res = await getChainList()

  chainList.value = res.data

}

onMounted(() => {
  init()
})

</script>

<template>

  <div class="w-screen min-h-screen bg-black flex flex-col">
    <div class="header" :style="headerStyle">
      <div class="w-1/3 h-full items-center flex pl-5">
        <span @click="backToMain" class="hover:cursor-pointer hover:text-white hover:underline transition">
          返回首页
        </span>
      </div>
      <div class="w-1/3 h-full flex items-center justify-center">
        当前正在浏览
      </div>
      <div class="w-1/3 h-full"></div>
    </div>
    <div class="body text-white">
        <div v-for="item in chainList">
          {{item}}
        </div>
    </div>
  </div>

</template>

<style>



.header {
  @apply w-full h-10 bg-green-500
}

.anim {
  animation: full-to-normal 1s ease-in-out;
}

.anim-reverse {
  animation: full-to-normal 1s ease-in-out reverse;
}

@keyframes full-to-normal {

  0% {
    height: 100vh;
    background-color: #272727;
    color: transparent;
  }

  30%{
    background-color: #272727;
    color: transparent;
  }

  80% {
    height: 2.2rem;
    //--tw-bg-opacity: 1;
    //background-color: rgb(239 68 68 / var(--tw-bg-opacity));
    --tw-bg-opacity: 1;
    background-color: rgb(34 197 94 / var(--tw-bg-opacity));
    color: transparent;

  }
  100% {
    height: 2.5rem;
    //--tw-bg-opacity: 1;
    //background-color: rgb(239 68 68 / var(--tw-bg-opacity));
    --tw-bg-opacity: 1;
    background-color: rgb(34 197 94 / var(--tw-bg-opacity));
    //background-color: #272727;
  }

}


</style>