<template>
    <div class="flex flex-row grow w-full">
        <div class="flex flex-row items-center overflow-hidden">



            <n-ellipsis class="w-96">
                <n-icon v-if="props.item['isDir']" size="16" class="mr-1 text-green-200">
                    <folder-sharp></folder-sharp>
                </n-icon>

                <n-icon v-else size="16" class="mr-1 text-yellow-200">
                    <document-text-sharp></document-text-sharp>
                </n-icon>


                <span :class="props.item['isDir'] ? dirStyle : otherStyle">
                    <span @click="onClick" class="hover:underline hover:cursor-pointer">{{ props.item['name'] }}</span>
                </span>

            </n-ellipsis>

        </div>

        <div class="w-1/3 text-[#add2d2]">
            <span v-if="!props.item['isDir']">{{ showSize }}</span>
        </div>

    </div>
</template>


<script setup>
import {  onMounted, ref } from 'vue';
import { FolderSharp, DocumentTextSharp } from '@vicons/ionicons5'
import { NIcon,NEllipsis } from 'naive-ui'

// props
const props = defineProps(['item', 'parent'])

// vars




//emits
const emits = defineEmits(['click'])

//methods
const onClick = () => {
    emits('click', props.item, props.parent)
}

// styles
const dirStyle = "text-white"
const otherStyle = "text-green-600"


const showSize = ref('')


onMounted(() => {


    let size = props.item['size']
    //if size > 1MB
    if (size > 1024 * 1024) {
        size = (size / 1024 / 1024).toFixed(2) + 'MB'
        showSize.value = size
        return
    }

    size = (size / 1024).toFixed(2) + 'KB'

    showSize.value = size


})


</script>