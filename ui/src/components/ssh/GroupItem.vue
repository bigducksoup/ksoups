<script setup>
import {FolderOutline} from '@vicons/ionicons5'
import {NButton, NDrawer, NDrawerContent, NForm, NFormItem, NIcon, NInput, NPopconfirm, NTreeSelect} from 'naive-ui'
import {onMounted, reactive, ref, toRaw, watch, watchEffect} from "vue";
import {useSSHGroupOperation} from "../../hooks/ssh.js";

/**
 * GItem <GItem>
 * { "id": "f39ff0df-b0e6-4a9b-8845-d52e0bb8cdfe", "name": "testg1", "parent": "root" }
 * @type {Prettify<Readonly<{[key in string]?: any}>>}
 */
const props = defineProps(['GItem', 'groupTree'])

const emit = defineEmits(['onDelete','onUpdate'])


const {DeleteSSHGroup, UpdateSSHGroup} = useSSHGroupOperation()


const tree = ref([])

const filterSelfAndChild = (tree, id) => {
  for (let i = 0; i < tree.length; i++) {
    if (tree[i].id === id) {
      tree.splice(i, 1)
      return
    }
    if (tree[i].children) {
      filterSelfAndChild(tree[i].children, id)
    }
  }
}


const item = reactive({
  id: '',
  name: '',
  parent: ''
})

// 保存更改
const onSubmit = async () => {

  let success = await UpdateSSHGroup(item)

  if (success) {
    showEditGroup.value = false
    emit('onUpdate', item)
  }
}


// 是否显示更改drawer
const showEditGroup = ref(false)
// 显示更改drawer
const onClickDetail = (e) => {
  e.stopPropagation()
  showEditGroup.value = true
}


// 点击删除
const onClickDelete = async () => {
  let success = await DeleteSSHGroup(item.id)
  if (success) {
    showEditGroup.value = false
    emit('onDelete', item.id)
  }
}


watchEffect(() => {
  item.id = props.GItem.id
  item.name = props.GItem.name
  item.parent = props.GItem.parent

  tree.value = JSON.parse(JSON.stringify(toRaw(props.groupTree)))
  filterSelfAndChild(tree.value, item.id)
})




</script>

<template>
  <div
      class="hover:cursor-pointer text-[#76D1FB] w-full min-w-max h-10 p-1 flex items-center border-b-[1px] border-b-gray-600 hover:bg-gray-800">
    <n-icon size="20" class="mr-2">
      <FolderOutline/>
    </n-icon>
    <span class="text-xl">
      {{ props.GItem.name }}
    </span>

    <n-button size="small" class="ml-auto" @click="onClickDetail"> 详情</n-button>
  </div>

  <n-drawer v-model:show="showEditGroup" width="500">
    <n-drawer-content title="详情">
      <n-form label-width="100px" ref="formRef">
        <n-form-item label="名称">
          <n-input v-model:value="item.name" placeholder="请输入名称"/>
        </n-form-item>

        <n-form-item label="父组">
          <n-tree-select
              default-expand-all
              :on-update-value="(value) => item.parent = value"
              :options="tree"
              label-field="name"
              key-field="id"
              children-field="children"
              placeholder="如果想要移动组，请选择目标组"
          />
        </n-form-item>
      </n-form>


      <template #footer>
        <n-popconfirm
            :on-positive-click="onClickDelete"
            positive-text="确定"
            negative-text="取消"
        >
          <template #trigger>
            <n-button type="error" class="mr-auto bg-red-500">
              删除此组
            </n-button>
          </template>
          此操作会<span class="text-red-500">删除当前组内所有内容</span>，你确定要继续吗！

        </n-popconfirm>
        <n-button type="primary" class="bg-green-500" @click="onSubmit"> 保存</n-button>
      </template>


    </n-drawer-content>
  </n-drawer>


</template>

<style scoped>

</style>