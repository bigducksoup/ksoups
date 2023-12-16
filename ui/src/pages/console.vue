
<template>
  <div class="flex flex-row grow">

    <div class="grow">
      <n-scrollbar>
        <Path :path="parent" @click="onClickPath"></Path>
        <div v-for="item in items">
          <FdItem :item="item" :parent="parent" @click="onClickItem"></FdItem>
        </div>
      </n-scrollbar>
    </div>


    <div class="w-[800px] flex flex-col pr-3">
      <WindowBar :title="editingFilePath" @on-close="closeFile"></WindowBar>
      <Code v-model="fileContent" height="650px"></Code>
      <div class="w-full grow flex flex-row-reverse p-3">
        <n-button class="text-white" type="primary" @click="showDiffContent">
          提交
        </n-button>
      </div>
    </div>


    <n-modal v-model:show="showDiff" preset="card" size="huge" :bordered="false" title="对比结果"
      :style="{ width: '800px', height: '700px', backgroundColor: '#23272E' }">

      <n-scrollbar class="h-[500px]">
        <div v-for="item in diffRes">

          <div v-if="item.added" class="text-green-500">
            <div v-html="toHtml(item.value)"></div>
          </div>

          <div v-else-if="item.removed" class="text-red-500">
            <div v-html="toHtml(item.value)"></div>
          </div>

          <div v-else class="text-white">
            <div v-html="toHtml(item.value)"></div>
          </div>
        </div>
      </n-scrollbar>

      <template #footer>
        <div class="w-full h-full flex flex-row-reverse">
          <n-button class="text-white" type="primary" @click="submitModify">确认修改</n-button>
        </div>
      </template>
    </n-modal>


  </div>
</template>


<script setup>
import FdItem from '../components/fdItem.vue';
import Path from '../components/path.vue'
import Code from '../components/code.vue'
import WindowBar from '../components/window-bar.vue'
import { useMessage, NScrollbar, NButton, NModal } from 'naive-ui'
import { onMounted } from 'vue';
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { baseUrl } from '../state'
import * as Diff from 'diff'


onMounted(() => {
  selectAddress.value = route.params.addr
  scanDir(route.params.addr, parent.value)
})


//选择的addr
const selectAddress = ref("")

const editingFilePath = ref(null)


//当前文件夹的内容
const items = ref([])

//文件内容
const fileContent = ref("")
//文件原始内容
const originContent = ref("")
//父路径
const parent = ref("/")

//显示diff
const showDiff = ref(false)
//diff 后的结果
const diffRes = ref([])




const message = useMessage()


const route = useRoute()


const toHtml = (text) => {
  return text.replace(/\n/g, "<br>")
}


const onClickItem = (item, parent) => {
  if (item['isDir']) {
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

  fetch(baseUrl.value + "api/dir/read?path=" + path + "&address=" + address + "&fileOnly=false").then(res => {
    return res.json()
  }).then(json => {
    items.value = json['data']['items']
    parent.value = json['data']['parent']
  })

}

// 查看文件内容
const scanFile = (address, path) => {

  fetch(baseUrl.value + "api/file/read?path=" + path + "&address=" + address).then(res => {
    return res.json()
  }).then(json => {
    if (json['code'] != 200) {
      message.error(json['message'])
      return
    }
    fileContent.value = json['data']['content']
    originContent.value = json['data']['content']
    editingFilePath.value = path
  })

}


const closeFile = () => {

  editingFilePath.value = null
  fileContent.value = ""
  originContent.value = ""
  diffRes.value = []

}


const showDiffContent = () => {

  showDiff.value = true

  diffRes.value = Diff.diffLines(originContent.value, fileContent.value)

}




const submitModify = () => {

  console.log(diffRes.value)

  let changes = []


  for (let diff of diffRes.value) {

    let operation = 2

    if (diff.added) {
      operation = 0
    }

    if (diff.removed) {
      operation = 1
    }

    let change = {
      "count": diff.count,
      "operation": operation
    }

    changes.push(change)
  }


  let modifyParams = {
    "path": editingFilePath.value,
    "addr": selectAddress.value,
    "changes": changes
  }

  reqModify(modifyParams)

}


const reqModify = (params) => {

  fetch(baseUrl.value + "api/file/modify", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(params)
  }).then(res => {
    return res.json()
  }).then(json => {
    if (json['code'] != 200) {
      message.error(json['message'])
      return
    }
    message.success(json['message'])
  })

}





</script>

