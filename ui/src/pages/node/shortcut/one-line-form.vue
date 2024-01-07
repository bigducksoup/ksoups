<script setup>

import {ref} from "vue";
import {baseUrl} from "../../../state/index.js";
import {useMessage} from 'naive-ui'
const props = defineProps(['probeId'])
const emits = defineEmits(['onClickClose','onSuccess'])

const message = useMessage()

const form = ref({
  "probeId":props.probeId,
  "name":"",
  "description":"",
  "timeout":0,
  "justRun":false,
  "payload":"",
  "type":0
})

const clear = ()=>{
  form.value = {
    "probeId":props.probeId,
    "Name":"",
    "Description":"",
    "Timeout":0,
    "JustRun":false,
    "Command":""
  }
}

const submit = async ()=>{

  let sid = window.localStorage.getItem('sid')

  let res =  await fetch(baseUrl.value + "api/shortcut/create" , {
    method:'POST',
    headers: {
      'sid': sid
    },
    body:JSON.stringify(form.value)
  })

  let json = await res.json()

  if (json['code']!==200){
    message.error(json['msg'])
    return
  }
  message.success("创建成功")
  cancel()
  emits('onSuccess')
}

const cancel = ()=>{
  clear()
  emits('onClickClose')
}

</script>

<template>
  <div class="h-screen w-screen flex items-center justify-center">
    <div class=" w-11/12 sm:w-[30rem]">
      <form  class="bg-[#4B5563] shadow-md rounded-2xl px-8 pt-6 pb-8 mb-4">
        <div class="mb-4">
          <label class="block text-white text-sm font-bold mb-2" for="name">
            快捷指令名称
          </label>
          <input
              v-model="form.name"
              class="appearance-none  bg-[#5D6673] rounded w-full py-2 px-3 text-gray-300 leading-tight focus:outline-none focus:shadow-outline"
              id="name"
              type="text"
              placeholder="输入指令名称"
          />
        </div>
        <div class="mb-4">
          <label class="block text-white text-sm font-bold mb-2" for="description">
            描述
          </label>
          <input
              v-model="form.description"
              class="appearance-none bg-[#5D6673] rounded w-full py-2 px-3 text-gray-300 leading-tight focus:outline-none focus:shadow-outline"
              id="description"
              type="text"
              placeholder="输入描述"
          />
        </div>
        <div class="mb-4">
          <label class="block text-white text-sm font-bold mb-2" for="timeout">
            超时时间（ms）
          </label>
          <input
              v-model="form.timeout"
              class="appearance-none bg-[#5D6673] rounded w-full py-2 px-3 text-gray-300 leading-tight focus:outline-none focus:shadow-outline"
              id="timeout"
              type="number"
              placeholder="ms"
          />
        </div>
        <div class="mb-4">
          <label class="block text-white text-sm font-bold mb-2" for="justRun">
            仅运行
          </label>
          <input
              v-model="form.justRun"
              class="mr-2 leading-tight text-gray-300"
              id="justRun"
              type="checkbox"
          />
          <span class="text-sm text-gray-300">忽略运行结果</span>
        </div>
        <div class="mb-4">
          <label class="block text-white text-sm font-bold mb-2" for="command">
            Command
          </label>
          <textarea
              v-model="form.payload"
              class="appearance-none bg-[#5D6673] rounded w-full py-2 px-3 text-gray-300 leading-tight focus:outline-none focus:shadow-outline"
              id="command"
              rows="4"
              placeholder="Enter the command"
          ></textarea>
        </div>
        <div class="flex items-center flex-row-reverse">
          <button
              @click="cancel"
              class="bg-gray-800 ml-2 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
              type="button"
          >
            取消
          </button>
          <button
              @click="submit"
              class="bg-green-800 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
              type="button"
          >
            提交
          </button>
        </div>
      </form>
    </div>
  </div>

</template>

<style scoped>



</style>