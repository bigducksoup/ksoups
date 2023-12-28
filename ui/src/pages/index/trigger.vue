<template>

<div class="text-white">123{{route.params.probeId}}</div>

</template>


<script setup>
import {useRoute} from 'vue-router'

const route = useRoute()

</script>

<style scoped>

</style>