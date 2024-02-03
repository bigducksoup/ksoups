<template>

  <div class="w-full bg-[#05253A]  flex items-center justify-end p-10 h-screen">


    <div class=" w-full sm:w-[400px] h-full bg-[#0C3451] rounded-2xl shadow flex flex-col p-8 text-white">
      <div class="w-full h-10 text-base mb-5">
        新概念运维中心
      </div>
      <div class="w-full h-12 text-2xl">
        更轻松，更高效
        <div>
          <span class="text-sm text-gray-500">当前只支持Linux, MacOS</span>
        </div>
      </div>

      <div class="w-full grow mt-10 flex flex-col gap-8">

        <n-input v-model:value="account" class="rounded-xl h-12 flex items-center border-gray-400 border-[1px]"
                 placeholder="输入您配置的账号">
          <template #prefix>
            <n-icon size="20">
              <PersonOutline class="w-6 h-6 text-gray-400"/>
            </n-icon>
          </template>
        </n-input>

        <n-input type="password" v-model:value="password"
                 class="rounded-xl h-12 flex items-center border-gray-400 border-[1px]" placeholder="输入您配置的密码">
          <template #prefix>
            <n-icon size="20">
              <LockClosedOutline class="w-6 h-6 text-gray-400"/>
            </n-icon>
          </template>
        </n-input>


        <div class="w-full flex justify-center mt-auto">
          <NButton type="primary" class="w-full rounded-xl h-12 bg-blue-400" :onClick="submit">即刻开始</NButton>
        </div>
      </div>

      <div class="w-full h-10  mt-10 flex justify-center items-center gap-2">
          <n-icon size="30">
            <LogoGithub/>
          </n-icon>
        <a href="https://github.com/bigducksoup">https://github.com/bigducksoup</a>
      </div>


    </div>

  </div>


</template>


<script setup>
import {NButton, NIcon, NInput,useMessage} from 'naive-ui'
import {LockClosedOutline, PersonOutline,LogoGithub} from '@vicons/ionicons5'
import {ref} from 'vue';
import {baseUrl} from '../state';
import {useRouter} from 'vue-router';
import CryptoJS from 'crypto-js'

const message = useMessage()



const account = ref('')
const password = ref('')


const router = useRouter()


const submit = () => {

  if (account.value === '') {
    message.error('账号不能为空')
    return
  }

  if (password.value === '') {
    message.error('密码不能为空')
    return
  }

  let p = md5(password.value)

  fetch(baseUrl.value + 'api/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      account: account.value,
      password: p
    })
  }).then(res => {
    return res.json()
  }).then(json => {
    console.log(json)
    if (json['code'] === 200) {
      window.localStorage.setItem('sid', json['data'])
      router.push('/')
    } else {
      message.error(json['msg'])
    }

  })
}


const md5 = str => CryptoJS.MD5(str).toString().toLowerCase()

</script>