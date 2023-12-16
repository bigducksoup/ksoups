<script setup>
import {NMessageProvider,NConfigProvider} from 'naive-ui'
import { onMounted } from 'vue';
import { baseUrl }  from './state'
import { darkTheme } from 'naive-ui'


onMounted(()=>{
        if(import.meta.env.DEV == false){
                baseUrl.value = window.location.href
        }
})

</script>

<template>
        <n-message-provider>
                <n-config-provider :theme="darkTheme">
                        <router-view>
                        </router-view>
                </n-config-provider>
        </n-message-provider>
</template>

<style scoped></style>
