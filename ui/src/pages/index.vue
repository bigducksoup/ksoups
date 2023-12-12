<template>
    <div class="layout flex flex-col w-full h-full">
      <div class="layout-top">
        <div class="bar w-full h-[30px] p-2 flex justify-center items-center border-2 border-[#EFEFF4]">
          <span> VVADMIN </span>
        </div>
      </div>

      <div class="layout-mid flex flex-row grow">

        <div class="layout-left  w-60 bg-white flex flex-col items-center">
          <div v-for="(item, index) in nodeArr" :class="item['addr'] == selectAddress ? seletedStyle : unselectStyle"
            @click="selectNode(item['addr'])">
            {{ item['addr'] }}
          </div>
        </div>

        <div class="layout-right flex flex-col p-2 w-full">

          <Path :path="parent" @click="onClickPath"></Path>

          <div class="flex flex-row w-full">

            <div class="w-1/2">
              <div v-for="item in items">
                <FdItem :item="item" :parent="parent" @click="onClickItem"></FdItem>
              </div>
            </div>


            <div class="w-1/2">

              <textarea class="w-full h-96" v-model="fileContent"></textarea>

            </div>

          </div>



        </div>

      </div>


      <div class="layout-footer flex justify-center items-center text-white">
        powered by native-ui, Vue3, TailWindCss
      </div>


    </div>
</template>




<script setup>
import { ref } from 'vue'
import { onMounted } from 'vue';
import FdItem from '../components/fdItem.vue';
import Path from '../components/path.vue'
import { useMessage } from 'naive-ui'

//当前选中的node地址
const selectAddress = ref('')

//在线的node
const nodeArr = ref([])

//当前文件夹的内容
const items = ref([])

//文件内容
const fileContent = ref("")
//父路径
const parent = ref("")



const unselectStyle = "w-[95%] m-2 rounded-md p-2 border-2 border-transparent hover:border-2 hover:border-black bg-red-500"

const seletedStyle = "w-[95%] m-2 rounded-md p-2 border-2 border-green-500 bg-green-500"


const message = useMessage()



const onClickItem = (item, parent) => {
  if (item['isDir']) {
    console.log(item['name'])
    scanDir(selectAddress.value, parent == '/' ? parent + item['name'] : parent + '/' + item['name'])
    return
  }

  if (!item['isDir']) {
    scanFile(selectAddress.value, parent + '/' + item['name'])
  }

}


const onClickPath = (path) => {
  scanDir(selectAddress.value, path)
}



// 扫描文件夹
const scanDir = (address, path) => {

  fetch("http://127.0.0.1:8080/dir/read?path=" + path + "&address=" + address + "&fileOnly=false").then(res => {
    return res.json()
  }).then(json => {
    items.value = json['data']['items']
    parent.value = json['data']['parent']
  })

}

// 查看文件内容
const scanFile = (address, path) => {


  fetch("http://127.0.0.1:8080/file/read?path=" + path + "&address=" + address).then(res => {
    return res.json()
  }).then(json => {
    if (json['code'] != 200) {
      message.error(json['message'])
      return
    }
    fileContent.value = json['data']['content']
  })

}

// 选择node
const selectNode = (addr) => {
  if (selectAddress.value == addr) return
  selectAddress.value = addr
  scanDir(addr, '/')
}


onMounted(() => {
  fetch("http://127.0.0.1:8080/info/nodes").then(res => {
    return res.json()
  }).then(json => {
    nodeArr.value = json
    let defaultAddr = nodeArr.value[0]['addr']
    selectNode(defaultAddr)

  })
})


</script>