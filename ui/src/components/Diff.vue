<template>
    <n-scrollbar class="w-full h-full">
        <div v-for="item in diffRes" class="text-base">

            <div v-if="item.added" class="text-green-500">
                <div v-html="toHtml(item.value)" ></div>
            </div>

            <div v-else-if="item.removed" class="text-red-500">
                <div v-html="toHtml(item.value)"></div>
            </div>

            <div v-else class="text-white">
                <div v-html="toHtml(item.value)"></div>
            </div>
        </div>
    </n-scrollbar>
</template>


<script setup>
import {onMounted, ref} from 'vue';
import {diffLines} from '../services/diff.js'
import {NScrollbar} from 'naive-ui'


const props = defineProps(['origin', 'current'])
const emits = defineEmits(['GetDiffRes'])


const diffRes = ref([])



const toHtml = (text) => {

  let one = text.replace(/\n/g, "<br>");

  return one.replace(/ /g, '&nbsp&nbsp&nbsp')
}

onMounted(() => {
    diffRes.value = diffLines(props.origin, props.current)
    emits('GetDiffRes',diffRes.value)
})


</script>