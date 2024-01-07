<template>
  <div class="w-full h-full flex flex-col p-2">
    <div class="w-full h-16 flex flex-row bg-[#0C0C0C]  items-center  p-1">
      <n-button circle class="bg-[rgb(18,20,23)]" :on-click="onRefresh">
        <template #icon>
          <n-icon>
            <Refresh/>
          </n-icon>
        </template>
      </n-button>

      <div
          class="h-9 flex flex-row items-center rounded-md grow bg-transparent p-1 ml-2 border-[1px] border-solid border-[#303438]">
        <Path :path="parent" @click="onClickPath"></Path>
      </div>

    </div>
    <div class="w-full h-16 flex flex-row items-center p-1">

      <n-button @click="showCreateForm = true">
          创建
      </n-button>


      <n-button-group class="ml-2">
        <n-button>复制</n-button>
        <n-button>移动</n-button>
      </n-button-group>

      <div class="grow w-1/3"></div>

      <div class="w-1/4">
        <n-input placeholder="搜索" v-model:value="search" :on-update-value="nameSearch" :on-change="nameSearch">
          <template #prefix>
            <n-icon :component="FlashOutline"/>
          </template>
        </n-input>
      </div>

    </div>
    <div class="w-full  p-1">
      <n-data-table :columns="tableCols" :data="data" :pagination="pagination" default-expand-all/>
    </div>
  </div>


  <n-modal
      v-model:show="codeModal"
      preset="card" :title="'编辑 ' + opendFilePath + ' 中'"
      class="w-[900px] h-[750px]"
  >
    <n-tabs type="line" default-value="edit" animated>
      <n-tab-pane name="edit" tab="编辑">
        <n-scrollbar class="w-full max-h-[550px]">
          <Code v-model="file.content" height="550px"></Code>
        </n-scrollbar>

      </n-tab-pane>
      <n-tab-pane name="diff" tab="对比">
        <n-scrollbar class="w-full max-h-[550px]">
          <Diff :current="file.content" :origin="codeBackUp" @get-diff-res="getDiffRes"></Diff>
        </n-scrollbar>
      </n-tab-pane>
    </n-tabs>
    <template #footer>
      <div class="w-full flex flex-row-reverse">
        <n-button class="ml-2 bg-blue-500" color="#60A5FA" @click="modifyContent">保存</n-button>
        <n-button @click="codeModal = false">关闭</n-button>
      </div>
    </template>
  </n-modal>



  <Transition
      name="bounce"
  >
    <create-form
        @close="showCreateForm = false"
        :parent="parent"
        :probe-id="probeId"
        v-if="showCreateForm"
    >
    </create-form>
  </Transition>


</template>


<script setup>
import {
  NButton,
  NButtonGroup,
  NDataTable,
  NDropdown,
  NIcon,
  NInput,
  NModal,
  NScrollbar,
  NTabPane,
  NTabs,
  useMessage,
} from 'naive-ui'
import {ChevronDown, DocumentTextOutline, FlashOutline, FolderOutline, Link, Refresh} from '@vicons/ionicons5'
import Path from '../../components/path.vue'
import Code from '../../components/code.vue'
import DragDown from '../../components/DropDown.vue'
import {h, onMounted, ref} from 'vue';
import {scanDir} from '../../services/dir.js'
import {diffLines} from '../../services/diff.js'
import {getFileContent, modifyFile} from '../../services/file.js'
import {useRoute} from 'vue-router';
import Diff from '../../components/Diff.vue'
import {formatTime} from '../../services/time.js'
import CreateForm from "./file/create-form.vue";

const message = useMessage()
const route = useRoute()

// 当前ip地址
const probeId = ref('')

// 父路径
const parent = ref("/")

// 搜索的内容
const search = ref("")

// data-table 的（目录，文件）数据
const data = ref([])

// data的备份
const backup = ref([])
// 打开的文件路径
const opendFilePath = ref('')

// 编辑modal是否打开
const codeModal = ref(false)
// 打开的文件信息
const file = ref({})
// 打开的文件内容备份
const codeBackUp = ref('')

//diff 后的 结果
const diffRes = ref([])


const showCreateForm = ref(false)

// 显示文件内容编辑modal
const showCodeModal = async (path) => {

  let response = await getFileContent(probeId.value, path)

  if (response['code'] !== 200) {
    message.error(response['msg'])
    return
  }

  diffRes.value = []

  file.value = response['data']

  codeBackUp.value = file.value['content']

  codeModal.value = true
}

// 初始化
const init = () => {
  goToPath(parent.value)
}

// 进入路径
const goToPath = async (path) => {

  let res = await scanDir(probeId.value, path)
  if (res['code'] !== 200) {
    message.error(res['msg'])
    return
  }

  parent.value = res['data']['parent']

  res['data']['items'].forEach((item) => {
    item.key = item['name']
  })

  data.value = res['data']['items']
  backup.value = data.value
  search.value = ''
}

// 点击文件或目录名
const onClickName = (item) => {

  let p = parent.value
  if (item['isDir']) {
    goToPath(p === '/' ? p + item['name'] : p + '/' + item['name'])
    return
  }

  let path = p + '/' + item['name']
  opendFilePath.value = path
  showCodeModal(path)
}

// 点击路径
const onClickPath = (path) => {
  goToPath(path)
}
// 刷新
const onRefresh = () => {
  goToPath(parent.value)
}

// 搜索
const nameSearch = () => {

  if (search.value === "") {
    data.value = backup.value
    return
  }

  let s = search.value
  let b = backup.value
  b = b.filter((item) => {
    return item['name'].toLowerCase().indexOf(s.toLowerCase()) !== -1
  })
  data.value = b

}

const getDiffRes = (res) => {
  diffRes.value = res
  console.log(diffRes.value)
}

const modifyContent = async () => {

  if (diffRes.value.length === 0) {
    diffRes.value = diffLines(codeBackUp.value, file.value.content)
  }

  let res = await modifyFile(probeId.value, opendFilePath.value, diffRes.value)

  if (res['code'] !== 200) {
    message.error(res['msg'])
  }
  codeModal.value = false
  message.info(res['msg'])
}







onMounted(() => {
  probeId.value = route.params.probeId
  init()
})




const pagination = ref({
  pageSize: 100
})

const tableCols = [
  {
    type: "selection"
  },
  {
    title: "名称",
    key: "name",
    width: 300,
    algin: "left",
    render: (row) => renderName(row, () => onClickName(row))
  },
  {
    title: "权限",
    key: "permission",
    render: (row) => renderTextButton(row, row.permission, () => console.log(row))
  },
  {
    title: "用户",
    key: "user",
    render: (row) => renderTextButton(row, row.user, () => console.log(1))
  },
  {
    title: "用户组",
    key: "usergroup",
    render: (row) => renderTextButton(row, row.usergroup, () => console.log(1))
  },
  {
    title: "大小",
    key: "size",
    render: (row) => h(
        'span',
        {
          class: 'text-blue-400'
        },
        {
          default: () => calcSize(row)
        }
    )
  },
  {
    title: "修改时间",
    key: "modTime",
    render: (row) => renderTextButton(row, formatTime(row.modTime), () => console.log(1))
  },
  {
    title: "操作",
    key: "operation",
    render: (row) => [
      h(
          'div',
          {
            class: "text-blue-500"
          },
          {
            default: () => (row.isDir || row.isLink) ?  '':[
              h('button',
                  {
                    class: 'mr-2 hover:underline'
                  },
                  {
                    default: () => "下载"
                  }),
              renderMoreOption(row)
            ]
          }
      ),
    ]
  },

]


const moreOptions = [
  {
    label: "链接到触发器",
    key: "marina bay sands",
  },
  {
    label: () => h(
        'span',
        {
          'class' : 'text-red-500'
        },
        {
          default : () => '删除'
        }
    ),
    key: "brown's hotel, london"
  }
]


const calcSize = (row) => {

  if (row.isDir) {
    return ''
  }

  let size = row.size


  //if size > 1MB
  if (size > 1024 * 1024) {
    size = (size / 1024 / 1024).toFixed(2) + 'MB'
    return size
  }

  size = (size / 1024).toFixed(2) + 'KB'

  return size

}


const renderMoreOption = (row) => h(
    NDropdown,
    {
      options: moreOptions,
      trigger: 'click'
    },
    {
      default: () => h(
          'button',
          {
            class: 'hover:underline'
          },
          {
            default: () => "更多"
          })
    }
)

const renderTextButton = (row, label, onClick) => h(
    'button',
    {
      class: "text-blue-400 hover:text-blue-500 w-full flex flex-row text-left",
      onClick: () => onClick(row)
    },
    {
      default: () => label
    })


const renderName = (row, onClick) => h(
    'div',
    {
      class: 'inline-flex  w-full items-center text-blue-400 hover:text-blue-600 hover:cursor-pointer',
      onClick: () => onClick(row)
    },
    {
      default: () => [
        h(NIcon, {
          size: '18',
          class: 'mr-2'
        }, {
          default: () => {
            if (row.isDir) return h(FolderOutline, {class: 'text-yellow-500'})
            if (row.isLink) return h(Link)
            return h(DocumentTextOutline, {class: 'text-white'})
          }
        }),
        h('div', {
          class: 'h-fit'
        }, {
          default: () => row.isLink ? row.name + ' ---> ' + row.linkTo : row.name
        })
      ]
    }
)

</script>