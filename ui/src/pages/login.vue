<template>

<div class="w-full bg-[#23272E]  flex flex-col items-center justify-center h-screen">


  <div class="bg-white shadow-md rounded-md p-8 max-w-xs w-full">
    <h2 class="text-2xl font-bold mb-4">登录账号</h2>
    <form>
      <div class="mb-4">
        <label for="username" class="block mb-2 text-sm font-medium text-gray-600">Account</label>
        <input v-model="account" type="username" id="username" name="username" class="border border-gray-300 rounded-md px-3 py-2 w-full focus:outline-none focus:border-blue-500" required>
      </div>
      <div class="mb-4">
        <label for="password" class="block mb-2 text-sm font-medium text-gray-600">Password</label>
        <input v-model="password" type="password" id="password" name="password" class="border border-gray-300 rounded-md px-3 py-2 w-full focus:outline-none focus:border-blue-500" required>
      </div>
      <div class="flex items-center justify-between">
        <button @click="submit" type="button" class="bg-blue-500 hover:bg-blue-600 text-white font-medium py-2 px-4 rounded-md focus:outline-none">登录</button>
        <a href="#" class="text-sm text-blue-500 hover:text-blue-600">Forgot Password?</a>
      </div>
    </form>
  </div>

</div>


</template>


<script setup>
import { ref } from 'vue';
import { baseUrl } from '../state';
import { useRouter } from 'vue-router';
import CryptoJS from 'crypto-js'

const account = ref('')
const password = ref('')


const router = useRouter()


const submit = ()=>{

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
        }else{
            alert(json['msg'])
        }

    })
}


const md5 = str => CryptoJS.MD5(str).toString().toLowerCase()

</script>