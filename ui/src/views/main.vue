<script setup>

import {useRouter} from 'vue-router'
import {initMessagePush} from "../hooks/websocket.js";
const router = useRouter()
initMessagePush()


</script>

<!-- 在你的 Vue 文件中 -->
<template>
  <div class="w-screen min-h-screen bg-[#05253A] flex backGround relative">
    <div
        class="w-full h-screen -translate-y-14  p-5 flex flex-col transition-all items-center justify-center md:justify-normal md:translate-y-0 md:w-60 md:items-start">
      <div class="nav-text" @click="router.push('/probe')">
        PROBE管理
      </div>
      <div class="nav-text" @click="router.push('/center')">
        CENTER管理
      </div>
      <div class="nav-text" @click="router.push('/chain')">
        命令调度
      </div>
      <div class="nav-text" @click="router.push('/ssh')">
        终端
      </div>
    </div>

    <div class="absolute left-1/2 top-1/2 md:-translate-y-1/2 transition-all translate-y-36 -translate-x-1/2 text-blue-500 text-[100px] ">
      KS<span class="text-green-500">O</span>U<span class="text-green-500">PS</span>
    </div>



  </div>
</template>

<style scoped>

.backGround {
  background-repeat: no-repeat;
  background-position: center;
  background-size: cover;
}

.nav-text {
  @apply
  text-3xl
  text-white
  italic w-fit
  underline-animation
  mb-4
  hover:cursor-pointer
  hover:text-green-500
  ani
}

.ani {
  animation: TextShowUp 0.3s ease-out;
}

/* 在你的 CSS 文件中 */
.underline-animation {
  position: relative;
  display: inline-block;
}

.underline-animation::after {
  content: '';
  position: absolute;
  width: 0;
  height: 2px;
  bottom: 0;
  left: 0;
  background-color: white;
  transition: width 0.3s ease-in-out;
}

.underline-animation:hover::after {
  width: 100%;
}


@keyframes TextShowUp {


  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(0);
  }

}

</style>