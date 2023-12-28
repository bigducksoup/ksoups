<template>
  <div class="layout flex flex-col w-screen h-screen overflow-hidden">

    <div v-if="nodeMap.size !== 0" class="layout-mid flex flex-row overflow-scroll grow bg-[#121417]">

      <div class="layout-left w-60 flex flex-col items-center bg-[#13151D]">

        <div class="w-full text-green-500 flex items-center justify-center h-10 mt-5">
          <n-popselect
            v-model:value="selectedItem['id']"
            :options="nodeOptions"
            size="medium"
            scrollable
            trigger="click"
            :on-update-value="chooseNode"
          >
            <n-button style="margin-right: 8px">
              {{ selectedItem['id'] || '选择一个node' }}
            </n-button>

          </n-popselect>

        </div>

        <n-menu
            v-model:value="activeKey"
            :indent="12"
            :options="menuOptions"
            class="w-full"
        />

      </div>

      <div class="layout-right flex h-fit w-full">

        <router-view v-if="selectedItem['online']===true" v-slot="{ Component }">
          <keep-alive>
            <component :key="$route.params.addr + $route.path" :is="Component"/>
          </keep-alive>
        </router-view>

        <div v-else class="text-white h-[99vh] w-full flex items-center justify-center text-3xl animate-pulse">
          当前Probe不可用...
        </div>

      </div>
    </div>

    <div v-else class="grow bg-[#121417] flex items-center justify-center">
      <span class="text-red-500">暂无在线节点</span>
    </div>

    <div class="layout-footer flex justify-center items-center text-white bg-[#121417] border-t-2 border-[#1E2228]">
      <!-- powered by native-ui, Vue3, TailWindCss -->
    </div>
  </div>
</template>


<script setup>
import {NMenu,NPopselect,NButton,useMessage} from 'naive-ui'
import {h, onMounted, ref} from 'vue'
import {baseUrl} from '../state'
import {RouterLink, useRouter} from 'vue-router' //当前选中的node地址

//当前选中的node地址
const selectedItem = ref({
  'id':"none"
})
const activeKey = ref('file')


//在线的node
const nodeMap = ref({})
const nodeOptions = ref([])
const router = useRouter()



const message = useMessage()


const chooseNode = (id)=>{

  selectedItem.value = nodeMap.value[id]


  if (selectedItem.value['online']===false){
    console.log(selectedItem.value)
    message.error(selectedItem.value['name']+" is offline")
    return
  }

  pushToApp(selectedItem,activeKey.value)

}




// 选择node
const pushToApp = (item, app) => {
  // if (selectedItem.value === item) return
  // selectedItem.value = item

  if (app === 'file'){
    router.push('/file/' + selectedItem.value['id'])
    return;
  }

  if (app === 'trigger'){
    router.push('/trigger/' + selectedItem.value['id'])
  }

}


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
        'label' : ()=> renderNode(item),
        'value' : item['id']
      })
    }

    let defaultNode= json[0]

    chooseNode(defaultNode['id'])
  })
})


const renderNode = (item) => {

  let publicStyle = 'w-2 h-2 rounded-full mr-2 ml-2'

  return h(
      'div',
      {
        'class':'flex flex-row items-center'
      },
      {
        default :() => [
          h(
              'div',
              {
                'class':item['online']===true ? publicStyle + ' bg-green-500 animate-pulse': publicStyle + ' bg-red-500 animate-ping'
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
          to : '/file/'+selectedItem.value['id']
        },
        {
          default: ()=>'文件'
        }
    ),
    key: 'file',
  },
  {
    label: () => h(
        RouterLink,
        {
          to : '/trigger/'+selectedItem.value['id']
        },
        {
          default: ()=>'触发器'
        }
    ),
    key: 'trigger'
  }
]




</script>