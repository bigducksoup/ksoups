<script setup>
import {NDataTable,NButton,useMessage,NPopconfirm,NModal,NInput,NIcon} from 'naive-ui'
import {GolfSharp,Close,Checkmark} from '@vicons/ionicons5'
import {useDataTable, useInput, useModal} from "../../hooks/common.js";
import {formatTime} from "../../services/time.js";
import {h, onMounted} from "vue";
import {deleteKeyPair, generateKeyPair, keyList} from "../../services/info.js";

const message = useMessage()

const { tableData,columns } = useDataTable([
  {
    title: 'ID',
    key: 'id',
    width: 50,
  },
  {
    title: '名称',
    key: 'name',
    width: 50,
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 50,
  },
  {
    title: '公钥',
    key: 'public_key',
    width: 200
  },
  {
    title: '公钥MD5',
    key: 'public_key_md5',
    width: 100,
    ellipsis: true
  },
  {
    title: '操作',
    key: 'action',
    width: 50,
    render: (row) => h('div', {
      class: 'flex text-white space-x-2'
    },[
      h(NPopconfirm,{
        onPositiveClick: () => deleteKey(row.id),
        positiveText: '删除',
        negativeText: '取消',
      },{
        default: () => '确定要删除吗',
        trigger: () => h(NButton, {type: 'error',size:'small',class:'hover:underline',text: true}, {default: () => '删除'})
      })
    ])
  }
],[

])


const deleteKey = async (id) => {
  let response = await deleteKeyPair(id)

  if (response.code !== 200) {
    message.error(response.msg)
    return
  }

  message.success('删除成功')
  tableData.value = tableData.value.filter((item) => item.id !== id)
}


const loadKeyTableData = async () => {
  const res = await keyList()

  if (res.code !== 200) {
    message.error(res.msg)
  }

  tableData.value = res.data.map((item) => {
    return {
      id: item.id,
      created_at: formatTime(item.created_at),
      public_key: item.public_key,
      public_key_md5: item.public_key_md5,
      name: item.name
    }
  })

}

const {input,clear} = useInput('')

const generateKey = async () => {

  if (input.value === '') {
    message.error('请输入密钥别名')
    return
  }

  console.log(input.value)

  const response = await generateKeyPair(input.value)

  if (response.code !== 200) {
    message.error(response.msg)
    return
  }
  close()
  clear()

  message.success('密钥生成成功')
  await loadKeyTableData()
}

const { visible,open,close } = useModal()


onMounted(()=>{
  loadKeyTableData()
})




</script>

<template>
  <div class="grow p-2 overflow-auto">

    <div class="w-full h-12 bg-[#26262A] flex items-center pl-2 pr-2">
      <span class="text-3xl text-white">
        密钥列表
      </span>

      <n-button type="success" class="text-white ml-auto" @click="visible = true">
        创建密钥
      </n-button>
    </div>

    <n-data-table
        :columns="columns"
        :data="tableData"
        :bordered="false"
        :pagination="{pageSize: 10}"
    >
    </n-data-table>


    <n-modal v-model:show="visible">
      <div class="w-80 h-80 rounded-lg bg-gray-800 p-5 flex flex-col ">
        <div class="w-full flex h-10">
          <n-button circle class="ml-auto" @click="close">
            <template #icon>
              <n-icon>
                <Close/>
              </n-icon>
            </template>
          </n-button>
        </div>
        <div class="w-full h-60 flex flex-col items-center justify-center">
          <n-icon size="80" class="text-purple-500">
            <GolfSharp/>
          </n-icon>
          <span class="text-2xl">
            输入密钥别名
          </span>
        </div>
        <div class="w-full flex mb-2">
          <n-input v-model:value="input" placeholder="" type="text" round size="large" class="mt-auto mr-3"></n-input>
          <n-button @click="generateKey" type="success" circle size="large" class="ml-auto text-white">
            <template #icon>
              <n-icon>
                <Checkmark/>
              </n-icon>
            </template>
          </n-button>
        </div>
      </div>
    </n-modal>

  </div>
</template>

<style scoped>

</style>