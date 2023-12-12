<script setup>
import {NMessageProvider,NConfigProvider} from 'naive-ui'

</script>

<template>
        <n-message-provider>
                <n-config-provider>
                        <router-view>
                        </router-view>
                </n-config-provider>
        </n-message-provider>
</template>

<style scoped></style>
