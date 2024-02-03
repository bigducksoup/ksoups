<script setup>
import {AtCircleOutline, Cloud, EyeOffOutline, Terminal} from "@vicons/ionicons5";
import {NButton, NDrawer, NDrawerContent, NForm, NFormItem, NIcon, NInput, NPopconfirm, NTreeSelect} from "naive-ui";
import {ref, watchEffect} from "vue";
import {useSSHInfoOperation} from "../../hooks/ssh.js";

/**
 * info <info>
 * { "id": "f51296cf-bcfc-4b7b-accb-aeb0ee0d7f0c", "addrPort": "127.0.0.1:22", "username": "ducksoup", "password": "******", "groupId": "root" }
 * @type {Prettify<Readonly<{[key in string]?: any}>>}
 */
const props = defineProps(["info", 'groupTree']);

const emit = defineEmits(['onDelete', 'onUpdate'])


const tree = ref([])

const info = ref({
  id: "",
  addrPort: "",
  username: "",
  password: "",
  groupId: "",
})


const rule = {
  addrPort: [
    {required: true, message: "地址端口不能为空"},
  ],
  username: [
    {required: true, message: "用户名不能为空"},
  ],
  password: [
    {required: false},
  ],
  groupId: [
    {required: true, message: "组ID不能为空"},
  ],
}
const showEditDrawer = ref(false)
const {UpdateSSHInfo,DeleteSSHInfo} = useSSHInfoOperation()

const submit = async () => {
  let success = await UpdateSSHInfo(info.value)

  if (success) {
    showEditDrawer.value = false
    emit('onUpdate', info.value)
  }
}

const deleteInfo = async () => {
  let success = await DeleteSSHInfo (info.value.id)

  if (success) {
    showEditDrawer.value = false
    emit('onDelete', info.value.id)
  }
}


watchEffect(() => {
  info.value = props.info
  info.value.password = ''
  tree.value = props.groupTree
})


</script>

<template>
  <div
      class="hover:cursor-pointer text-[#0078D7] w-full min-w-max h-10 p-1 flex items-center border-b-[1px] border-b-gray-600 hover:bg-gray-800"
  >
    <n-icon size="20" class="mr-2 text-2xl">
      <Terminal/>
    </n-icon>

    <span class="single">
      <n-icon size="20" class="mr-2 text-2xl">
        <Cloud/>
      </n-icon>
      {{ info.addrPort }}
    </span>

    <span class="single">
      <n-icon size="20" class="mr-2 text-2xl">
        <AtCircleOutline/>
      </n-icon>
      {{ info.username }}
    </span>

    <span class="single">
      <n-icon size="20" class="mr-2 text-2xl">
        <EyeOffOutline/>
      </n-icon>
      {{ info.password }}
    </span>


    <n-button size="small" class="ml-auto" @click="(e)=>{
      e.stopPropagation()
      showEditDrawer = true
    }"> 详情
    </n-button>

  </div>


  <n-drawer v-model:show="showEditDrawer" width="500">
    <n-drawer-content title="详情">
      <n-form :model="info" :rules="rule">

        <n-form-item label="地址端口" path="addrPort">
          <n-input v-model:value="info.addrPort"/>
        </n-form-item>

        <n-form-item label="用户名" path="username">
          <n-input v-model:value="info.username"/>
        </n-form-item>

        <n-form-item label="密码" path="password">
          <n-input type="password" v-model:value="info.password" placeholder="如需更改密码，请输入新密码"/>
        </n-form-item>


        <n-form-item label="父组">
          <n-tree-select
              default-expand-all
              :on-update-value="(value) => info.groupId = value"
              :options="tree"
              label-field="name"
              key-field="id"
              children-field="children"
              placeholder="如果想要移动位置，请选择目标组"
          />
        </n-form-item>

      </n-form>

      <template #footer>
        <n-popconfirm
            positive-text='确认'
            negative-text='取消'
            :on-positive-click="deleteInfo"
        >
          <template #trigger>
            <n-button type="error" class="bg-red-500 mr-auto">删除</n-button>
          </template>
          此操作会<span class="text-red-500">删除当前信息</span>，你确定要继续吗！
        </n-popconfirm>

        <n-button @click="submit" type="primary" class="bg-green-500">保存更改</n-button>
      </template>

    </n-drawer-content>
  </n-drawer>

</template>

<style scoped>
.single {
  @apply ml-20
  w-36
  flex items-center;
}
</style>
