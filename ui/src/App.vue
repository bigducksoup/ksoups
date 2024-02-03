<script setup>
import {NMessageProvider,NConfigProvider,NNotificationProvider} from 'naive-ui'
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
                  return
                }
                // router.push('/')
        })




})





</script>

<template>
        <n-message-provider>
          <n-notification-provider>
            <n-config-provider :theme="darkTheme" class="w-full h-full">
                <router-view v-slot="{ Component,route }">

                  <Transition :name="route.meta.transition" >
                    <component  :is="Component"/>
                  </Transition>

                </router-view>
            </n-config-provider>
          </n-notification-provider>
        </n-message-provider>
</template>

<style scoped></style>
