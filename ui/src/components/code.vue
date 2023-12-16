<template>
    <!-- :indent-with-tab="true" 是否自动获取焦点-->
    <codemirror v-model="code"  :style="{ height : props.height ,width: '100%',backgroundColor: '#23272E' }" :autofocus="true"
        :tabSize="2" :extensions="extensions"  />
</template>


<script setup>
import { Codemirror } from "vue-codemirror";
import { ref } from "vue";
import { oneDark } from "@codemirror/theme-one-dark";



const code = ref(``);

const props =  defineProps(['height'])


const extensions = [oneDark];

</script>
