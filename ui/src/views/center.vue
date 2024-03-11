<script setup>
import {NIcon, NMenu} from 'naive-ui'
import {h, onMounted, ref} from "vue";
import {useRouter} from "vue-router";
import {useMenu} from "../hooks/common.js";
import {Key,FitnessOutline} from "@vicons/ionicons5";

const router = useRouter()

const { MenuOption,Key:menuKey ,onUpdate, } = useMenu([
  {
    label:'监控',
    key: 'monitor',
    icon: () => h(NIcon,null,{
      default : () => h(FitnessOutline)
    })
  },
  {
    label:'密钥',
    key: 'keys',
    icon: () => h(NIcon,null,{
      default : () => h(Key)
    })
  },
],(key, item)=>{
  console.log(key, item)
  menuKey.value = key
  console.log(menuKey.value, key)
  router.push(`/center/${key}`)
})

const header = ref(null)


onMounted(() => {
  header.value.classList.add('anim')
  onUpdate('keys',null)
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
  <div ref="root" id="root-term" class="w-full h-full bg-black">
    <div ref="header" class="flex items-center pl-2">
      <button @click="leave" class="hover:underline hover:text-white">
        返回首页
      </button>
    </div>
    <div class="h-full w-full overflow-scroll flex">
      <div class="md:w-[180px] w-40 bg-[#13151D]">
        <n-menu
            :options="MenuOption"
            class="w-full"
            :on-update-value="onUpdate"
            v-model:value="menuKey"
        >
        </n-menu>
      </div>
      <div class="grow w-0 h-fit">
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