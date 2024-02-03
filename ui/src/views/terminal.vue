<script setup>
import MyTerm from '../components/ssh/MyTerm.vue'
import { useRoute } from 'vue-router';

const route = useRoute();




</script>

<template>
    <div class="w-full h-full">
        <MyTerm :ssh-id="route.query.sshId"></MyTerm>
    </div>
</template>