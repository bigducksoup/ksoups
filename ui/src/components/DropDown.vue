<template>
    <n-dropdown trigger="hover" :options="options" @select="handleSelect">
        <n-button class="bg-[#025DEA]" color="#025DEA">
            <template #default>
                <slot name="default"></slot>
            </template>
        </n-button>
    </n-dropdown>
</template>


<script setup>

import { NButton, NDropdown } from 'naive-ui'
import { useMessage } from "naive-ui";
import { onMounted, ref } from 'vue';



const props = defineProps(['options'])
const emits = defineEmits(['select'])

const options = ref([])


const handleSelect = (key) => {
    emits('select', key)
}


onMounted(()=>{

    options.value = props.options

})


</script>