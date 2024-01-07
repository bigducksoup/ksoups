<script setup>

import {useRouter} from 'vue-router'

const router = useRouter()


</script>

<!-- 在你的 Vue 文件中 -->
<template>
  <div class="w-screen min-h-screen bg-[#13151D] flex backGround">
    <div class="w-0 h-screen md:w-60  p-5 flex flex-col">
      <div class="nav-text" @click="router.push('/node')">
        节点管理
      </div>
      <div class="nav-text" @click="router.push('/chain')">
        命令调度
      </div>
    </div>
  </div>
</template>

<style scoped>

.backGround {
  background-image: url("/bg.jpg");
  background-repeat: no-repeat;
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