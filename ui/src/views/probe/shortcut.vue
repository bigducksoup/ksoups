<template>


  <div ref="root" class="w-full  p-2">
    <div class="ml-2 text-3xl flex text-[#BFCFE7]">
      <div>
        快捷指令
      </div>
      <div class="ml-auto mr-2">
        <n-button
            type="primary"
            class="bg-green-500"
            @click="createSCShow = true"
        >
          <n-icon>
            <Add/>
          </n-icon>
          创建快捷指令
        </n-button>
      </div>
    </div>
    <div class="w-full flex flex-row flex-wrap text-white p-2">
      <n-data-table :columns="columns" :data="tableData" :pagination="{pageSize:50}"></n-data-table>
    </div>


    <n-drawer v-model:show="createSCShow" :width="500">
      <n-drawer-content title="创建快捷指令" :native-scrollbar="false">
        <n-form :model="shortcutForm" :label-width="160">

          <n-form-item label="类型">
            <n-radio-group v-model:value="shortcutForm.type" name="typeChoice">
              <n-radio-button  :value="0">
                单行指令
              </n-radio-button>
              <n-radio-button :value="1">
                脚本指令
              </n-radio-button>
            </n-radio-group>
          </n-form-item>

          <n-form-item label="指令名称">
            <n-input v-model:value="shortcutForm.name" placeholder="输入指令名称"/>
          </n-form-item>
          <n-form-item label="指令描述">
            <n-input v-model:value="shortcutForm.description" placeholder="输入指令描述"/>
          </n-form-item>
          <n-form-item label="超时时间(ms) （超过超时时间后，会kill指令）">
            <n-input-number v-model:value="shortcutForm.timeout"/>
          </n-form-item>
          <n-form-item label="是否仅运行">
            <n-switch
                v-model:value="shortcutForm.justRun"
                @update-value="(val)=>{shortcutForm.justRun = val}"
            >
            </n-switch>
          </n-form-item>
          <n-form-item label="指令内容">
            <Code v-if="shortcutForm.type === 1" v-model="shortcutForm.payload" height="250px">
            </Code>
            <n-input v-else v-model:value="shortcutForm.payload" placeholder="请输入指令"></n-input>
          </n-form-item>
          <n-form-item>
            <div class="flex w-full flex-row justify-end">
              <n-button
                  type="primary"
                  class="ml-2 bg-green-500"
                  @click="submitShortcutForm"
              >
                创建
              </n-button>
            </div>
          </n-form-item>
        </n-form>
      </n-drawer-content>
    </n-drawer>


    <n-drawer v-model:show="shortcutDetailShow" :width="500">
      <n-drawer-content title="快捷指令详情" :native-scrollbar="false">
        <n-form
            :model="detailShortcut"
            :label-width="160"
        >
          <n-form-item label="指令名称">
            <n-input v-model:value="detailShortcut.name" placeholder="Input"/>
          </n-form-item>
          <n-form-item label="指令描述">
            <n-input v-model:value="detailShortcut.description" placeholder="Input"/>
          </n-form-item>
          <n-form-item label="超时时间(ms)">
            <n-input-number v-model:value="detailShortcut.timeout" placeholder="Input"/>
          </n-form-item>
          <n-form-item label="是否仅运行">
            <n-switch
                v-model:value="detailShortcut.justRun"
                @update-value="(val)=>{detailShortcut.justRun = val}"
            >
            </n-switch>
          </n-form-item>
          <n-form-item label="指令内容">
            <Code v-model="detailShortcut.payload" height="300px">
            </Code>
          </n-form-item>
          <n-form-item>
            <div class="flex w-full flex-row justify-end">
              <n-button
                  type="error"
                  class="bg-red-500"
                  @click="onClickDeleteShortcut"
              >
                删除
              </n-button>
              <n-button
                  type="primary"
                  class="ml-2 bg-green-500"
                  @click="onClickUpdateShortcut"
              >
                更新
              </n-button>
            </div>
          </n-form-item>

        </n-form>
      </n-drawer-content>
    </n-drawer>


  </div>


</template>


<script setup>
import {useRoute} from 'vue-router'
import {h, onMounted, ref} from "vue";

import {
  NButton,
  NDataTable,
  NDrawer,
  NDrawerContent,
  NForm,
  NFormItem,
  NIcon,
  NInput,
  NInputNumber,
  NPopconfirm, NRadioButton, NRadioGroup,
  NSwitch,
  useNotification
} from 'naive-ui'
import {Add, Play} from '@vicons/ionicons5'
import Code from '../../components/code.vue'
import {listShortcut, runShortcut} from '../../services/shortcut.js'
import {useAnimation} from "../../hooks/animation.js";
import {useShortcutOperation} from "../../hooks/shortcut.js";
import {useDataTable, useForm} from "../../hooks/common.js";


const route = useRoute()
const notification = useNotification()

const shortcuts = ref([])

const createSCShow = ref(false)
const shortcutDetailShow = ref(false)

const detailShortcut = ref({})

const {form: shortcutForm, Submit: submitShortcutForm} = useForm({
  'probeId': route.params.probeId,
  'name': '',
  'description': '',
  'timeout': 0,
  'justRun': false,
  'payload': '',
  'type': 0
}, async () => {
  let success = await CreateShortcut(shortcutForm)
  if (success) {
    await init()
    createSCShow.value = false
  }
})

const root = ref(null)

const {SameTimeAnimate: RootAnim} = useAnimation([root])


const {DeleteShortcutById, UpdateShortcut, CreateShortcut} = useShortcutOperation()


const {tableData, columns} = useDataTable([
  {
    title: '名称',
    key: 'name',
    render: (row) => h('div', {
      class: 'text-blue-500 hover:cursor-pointer',
      onClick: () => shortcutDetail(row)
    }, row.name)
  },
  {
    title: '描述',
    key: 'description',
    width: 200,
    ellipsis: {
      tooltip: true
    }
  },
  {
    title: '超时时间(ms)',
    key: 'timeout',
    render: (row) => {
      return row.timeout
    },
  },
  {
    title: '是否仅运行',
    key: 'justRun',
    render: (row) => {
      return row.justRun ? 'Yes' : 'No'
    },
    width: '100px'
  },
  {
    title: '操作',
    key: 'action',
    width: 200,
    render: (row) => h('div', {
      class: 'flex flex-row items-center'
    }, [
      h(NPopconfirm,
          {
            onPositiveClick: () => runOLShortcut(row),
            positiveText: '确认',
            negativeText: '取消'
          },
          {
            trigger: () => h(NButton, {
              class: 'w-7 h-7 bg-gray-600 hover:bg-blue-600 pl-0.5 flex items-center justify-center',
              circle: true
            }, () => h(NIcon, {size: 20}, () => h(Play))),
            default: () => '你可想好了嗷！'
          },
      ),
      h(NButton, {
        class: 'ml-2 flex hover:text-green-500 items-center justify-center',
        onClick: () => {
          shortcutDetail(row)
        }
      }, () => '详情'),
        h(NButton,{
          class: 'ml-2 flex hover:text-red-500 items-center justify-center',
        },() => '执行记录')
    ])
  }
], shortcuts)


const init = async () => {
  let res = await listShortcut(route.params.probeId)
  shortcuts.value = res.data
}


const runOLShortcut = async (shortcut) => {
  // let date = new Date()

  let res = await runShortcut(shortcut['id'])
  let type = res['ok'] ? 'success' : 'error'

  notification[type]({
    title: shortcut.Name,
    description: res['ok'] ? "执行成功" : "执行失败",
    content: res['out'],
    meta: new Date().toLocaleString()
  })
}


const shortcutDetail = (shortcut) => {
  shortcutDetailShow.value = true
  detailShortcut.value = shortcut
}


const onClickDeleteShortcut = async () => {
  let success = await DeleteShortcutById(detailShortcut.value['id'])
  if (success) {
    await init()
    shortcutDetailShow.value = false
  }
}

const onClickUpdateShortcut = async () => {
  let success = UpdateShortcut(detailShortcut.value)
  if (success) {
    shortcutDetailShow.value = false
  }
}

onMounted(() => {
  RootAnim('animate__slideInUp', 'animate__fast')
  init()
})


</script>

<style scoped>
.bounce-enter-active {
  animation: bounce-in 0.5s;
}

.bounce-leave-active {
  animation: bounce-in 0.5s reverse;
}

</style>