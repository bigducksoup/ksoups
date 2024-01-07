<script setup>
import {NButton, NCheckbox, NCheckboxGroup, NForm, NFormItem, NInput, NRadio, useMessage} from 'naive-ui'
import {onMounted, ref} from "vue";
import {createDir} from "../../../services/dir.js";
import {createFile} from "../../../services/file.js";

const props = defineProps(['parent', 'probeId'])
const emit = defineEmits(['create-success', 'close'])
const message = useMessage()


const type = ref('file')
const name = ref('')
const userPermission = ref([])
const groupPermission = ref([])
const otherPermission = ref([])

const form = ref({
  path: '',
  probeId: '',
  permission: ''
})

const fillForm = () => {

  form.value.path = props.parent + (props.parent === '/' ? '' : '/') + name.value
  let userPermissionNum = 0
  let groupPermissionNum = 0
  let otherPermissionNum = 0

  userPermission.value.forEach(item => {
    userPermissionNum += item
  })
  groupPermission.value.forEach(item => {
    groupPermissionNum += item
  })
  otherPermission.value.forEach(item => {
    otherPermissionNum += item
  })

  form.value.permission = userPermissionNum.toString() + groupPermissionNum.toString() + otherPermissionNum.toString()
}


const create = () => {

  fillForm()

  if (type.value === 'file') {
    createFile(form.value.probeId, form.value.path, form.value.permission).then(res => {
      if (res.code === 200) {
        message.success('创建成功')
        emit('create-success')
        emit('close')
      } else {
        message.error('创建失败')
      }
    })
  } else {
    createDir(form.value.probeId, form.value.path, form.value.permission).then(res => {
      if (res.code === 200) {
        message.success('创建成功')
        emit('create-success')
        emit('close')
      } else {
        message.error('创建失败')
      }
    })
  }

}

onMounted(() => {
  form.value.probeId = props.probeId
})


</script>

<template>
  <div class="w-screen h-screen z-10 absolute left-0 top-0 bg-transparent flex items-center justify-center">
    <div class="w-11/12 h-5/6 rounded-2xl shadow-2xl lg:w-[30rem] p-5 bg-gray-600 flex-col flex">
      <div class="mb-4 w-full flex justify-between">
        <span class="text-white"> 在 <span class="text-pink-400"> {{ props.parent }}</span> 内创建 </span>

        <span @click="emit('close')" class="hover:cursor-pointer">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 50 50">
          <line x1="10" y1="10" x2="40" y2="40" stroke="black" stroke-width="5"/>
          <line x1="40" y1="10" x2="10" y2="40" stroke="black" stroke-width="5"/>
          </svg>
        </span>

      </div>
      <n-form>
        <n-form-item label="类型">
          <n-radio value="file" @change="type = 'file'" :checked="type === 'file'">
            文件
          </n-radio>
          <n-radio value="dir" @change="type = 'dir'" :checked="type === 'dir'">
            文件夹
          </n-radio>
        </n-form-item>
        <n-form-item label="名称">
          <n-input v-model:value="name" placeholder="名称"></n-input>
        </n-form-item>
        <n-form-item label="权限">

          <div class="w-full flex flex-col">
            用户：
            <div class="w-full flex pb-3">
              <n-checkbox-group v-model:value="userPermission">
                <n-checkbox :value="4">
                  读
                </n-checkbox>
                <n-checkbox :value="2">
                  写
                </n-checkbox>
                <n-checkbox :value="1">
                  执行
                </n-checkbox>
              </n-checkbox-group>
            </div>
            用户组：
            <div class="w-full flex">
              <n-checkbox-group v-model:value="groupPermission">
                <n-checkbox :value="4">
                  读
                </n-checkbox>
                <n-checkbox :value="2">
                  写
                </n-checkbox>
                <n-checkbox :value="1">
                  执行
                </n-checkbox>
              </n-checkbox-group>
            </div>

            所有人：
            <div class="w-full flex pb-3">
              <n-checkbox-group v-model:value="otherPermission">
                <n-checkbox :value="4">
                  读
                </n-checkbox>
                <n-checkbox :value="2">
                  写
                </n-checkbox>
                <n-checkbox :value="1">
                  执行
                </n-checkbox>
              </n-checkbox-group>
            </div>

          </div>

        </n-form-item>
      </n-form>
      <div class="w-full grow"></div>
      <div class="w-full">
        <n-button class="w-full" type="primary" @click="create">创建</n-button>
      </div>
    </div>
  </div>
</template>

<style scoped>


</style>