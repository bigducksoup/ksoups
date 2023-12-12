<template>
    <div class="flex flex-row items-center">
        <n-icon v-if="props.item['isDir']" size="20" class="mr-1">
            <folder-sharp></folder-sharp>
        </n-icon>

        <n-icon v-else size="20" class="mr-1">
            <document-text-sharp></document-text-sharp>
        </n-icon>

        <span :class="props.item['isDir'] ? dirStyle : otherStyle">
            <span @click="onClick" class="hover:underline hover:cursor-pointer">{{ props.item['name'] }}</span>
        </span>
    </div>
</template>


<script setup>
import { onMounted } from 'vue';
import { FolderSharp, DocumentTextSharp } from '@vicons/ionicons5'
import { NIcon } from 'naive-ui'

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
const dirStyle = "text-black"
const otherStyle = "text-green-600"



onMounted(() => {

})


</script>