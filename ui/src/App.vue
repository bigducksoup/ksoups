<script setup>
import {NMessageProvider,NConfigProvider} from 'naive-ui'
import { onMounted } from 'vue';
import { baseUrl }  from './state'
import { darkTheme } from 'naive-ui'
import {useRouter} from 'vue-router'

onMounted(()=>{
        if(import.meta.env.DEV === false){
                baseUrl.value = window.location.href
        }

        let sid = window.localStorage.getItem('sid')
        
        const router = useRouter()

        if(!sid){
                router.push('/login')
                return
        }

        console.log(sid)


        fetch(baseUrl.value + 'api/auth/check_login', {
                method: 'POST',
                headers: {
                        'Content-Type': 'application/json',
                        'sid' : sid
                },
        }).then((res)=>{
               return res.json()
        }).then(json=>{
                if(json.code !== 200){
                        router.push('/login')
                }
        })




})





</script>

<template>
        <n-message-provider>
                <n-config-provider :theme="darkTheme" class="w-full h-full">
                        <router-view>
                        </router-view>
                </n-config-provider>
        </n-message-provider>
</template>

<style scoped></style>
