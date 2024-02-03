<template>

    <div class="w-full h-6 bg-[#242424] relative flex items-center border-b-[1px] border-black">

        <div class="w-2"></div>
        <div @click="onClose" class="w-4 h-4 rounded-full flex items-center justify-center bg-red-500">
            <NIcon>
                <Close />
            </NIcon>
        </div>
        <div class="w-2"></div>
        <div class="w-4 h-4 rounded-full bg-yellow-500"></div>
        <div class="w-2"></div>
        <div class="w-4 h-4 rounded-full bg-green-500"></div>
        

        <div class="w-full h-6 absolute bg-transparent text-white flex items-center justify-center pointer-events-none">
            {{ props.title }}
        </div>

    </div>



</template>



<script setup>
import {Close} from '@vicons/ionicons5'
import { NIcon } from 'naive-ui'


const emits = defineEmits(['onClose'])

const props = defineProps(['title'])

const onClose = () => {
    emits('onClose')
}




</script>