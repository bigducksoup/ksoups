<template>
  <div class="layout flex flex-col w-screen ">

    <div
        class="w-full h-10 flex"
        :style="topBarStyle"
    >
      <div class="w-1/3 h-full items-center flex pl-5">
        <span @click="backToMain" class="hover:cursor-pointer hover:text-white hover:underline transition">
          返回首页
        </span>
      </div>
      <div class="w-1/3 h-full flex items-center justify-center">
        当前正在浏览 <span class="text-purple-900">{{ selectedItem['name'] }}</span>
      </div>
      <div class="w-1/3 h-full"></div>
    </div>

    <div v-if="nodeMap.size !== 0" class="layout-mid flex flex-row h-screen overflow-scroll bg-[#121417]">

      <div class="layout-left w-60 flex flex-col items-center bg-[#13151D]">
        <div class="w-full text-green-500 flex items-center justify-center h-10 mt-5">
          选择节点
          <n-popselect
              v-model:value="selectedItem['id']"
              :options="nodeOptions"
              size="medium"
              scrollable
              trigger="click"
              :on-update-value="chooseNode"
          >
            <n-button size="small" class="ml-2">{{ selectedItem['id'] || '选择一个node' }}</n-button>
          </n-popselect>
        </div>

        <n-menu
            v-model:value="activeKey"
            :indent="12"
            :options="menuOptions"
            class="w-full"
        />
      </div>

      <div class="layout-right flex w-full">

        <router-view v-if="selectedItem['online']===true" v-slot="{ Component }">
          <keep-alive>
            <component :key="$route.params.addr + $route.path" :is="Component"/>
          </keep-alive>
        </router-view>

        <div v-else class="text-white w-full flex items-center justify-center text-3xl animate-pulse">
          当前Probe不可用...
        </div>

      </div>
    </div>

    <div v-else class="grow bg-[#121417] flex items-center justify-center">
      <span class="text-red-500">暂无在线节点</span>
    </div>

  </div>
</template>


<script setup>
import {NButton, NMenu, NPopselect, useMessage} from 'naive-ui'
import {h, onMounted, ref} from 'vue'
import {baseUrl} from '../state'
import {RouterLink, useRouter} from 'vue-router'

//当前选中的node地址
const selectedItem = ref({
  'id': "none"
})
const activeKey = ref('file')


//在线的node
const nodeMap = ref({})
const nodeOptions = ref([])
const router = useRouter()


const message = useMessage()


const topBarStyle = ref({
  backgroundColor: '#262626',
  transition: '0.9s',
  animation : 'screen-up 0.6s',
  animationIterationCount: '1'
})


const chooseNode = (id) => {

  selectedItem.value = nodeMap.value[id]


  topBarStyle.value.backgroundColor = selectedItem.value['online']===true? '#22C55D':'#EF4444'

  if (selectedItem.value['online'] === false) {
    console.log(selectedItem.value)
    message.error(selectedItem.value['name'] + " is offline")
    return
  }

  pushToApp(selectedItem, activeKey.value)

}


// 选择node
const pushToApp = (item, app) => {
  // if (selectedItem.value === item) return
  // selectedItem.value = item

  if (app === 'file') {
    router.push('/probe/file/' + selectedItem.value['id'])
    return;
  }

  if (app === 'shortcut') {
    router.push('/probe/shortcut/' + selectedItem.value['id'])
  }
}
const backToMain = () => {

  // Reset the animation
  topBarStyle.value.animation = 'none';
  topBarStyle.value.backgroundColor = '#262626';

  // Use requestAnimationFrame to ensure the style changes are applied
  requestAnimationFrame(() => {

    // Re-run the animation
    topBarStyle.value.animation = 'screen-up 0.68s reverse';
  });

  setTimeout(() => {
    // Show the content
    router.push('/');
  }, 650);
};


onMounted(() => {

  let sid = window.localStorage.getItem('sid')

  fetch(baseUrl.value + "api/info/nodes", {
    method: "GET",
    headers: {
      "sid": sid
    }
  }).then(res => {
    return res.json()
  }).then(json => {

    for (let item of json) {
      nodeMap.value[item['id']] = item
      nodeOptions.value.push({
        'label': () => renderNode(item),
        'value': item['id']
      })
    }

    let defaultNode = json[0]
    chooseNode(defaultNode['id'])
  })
})


const renderNode = (item) => {

  let publicStyle = 'w-2 h-2 rounded-full mr-2 ml-2'

  return h(
      'div',
      {
        'class': 'flex flex-row items-center'
      },
      {
        default: () => [
          h(
              'div',
              {
                'class': item['online'] === true ? publicStyle + ' bg-green-500 animate-pulse' : publicStyle + ' bg-red-500 animate-ping'
              }
          ),
          h(
              'span',
              null,
              {
                default: () => item['name']
              }
          )
        ]
      }
  )
}


const menuOptions = [
  {
    label: () => h(
        RouterLink,
        {
          to: '/probe/file/' + selectedItem.value['id'],
          'class': 'ml-2'
        },
        {
          default: () => '文件'
        }
    ),
    key: 'file',
  },
  {
    label: () => h(
        RouterLink,
        {
          to: '/probe/shortcut/' + selectedItem.value['id'],
          'class': 'ml-2'
        },
        {
          default: () => '捷径'
        }
    ),
    key: 'shortcut'
  }
]
</script>


<style>


@keyframes screen-up {
  0% {
    height: 100vh;
    width: 100vw;
    color: transparent;
  }

  60% {
    height: 2.2rem;
    width: 100vw;
  }

  100% {
    height: 2.5rem;
  }
}


</style>