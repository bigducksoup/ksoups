<script setup>
import {darkTheme, NConfigProvider, NMessageProvider, NNotificationProvider} from 'naive-ui'
import {onMounted} from 'vue';
import {baseUrl,baseHost} from './state'
import {useRoute, useRouter} from 'vue-router'

const router = useRouter()

onMounted(() => {

  if (import.meta.env.DEV === false) {

    const url = new URL(window.location.href)
    baseUrl.value = url.origin + '/'
    baseHost.value = url.host
  }

  let sid = window.localStorage.getItem('sid')

  const router = useRouter()

  if (!sid) {
    router.push('/login')
    return
  }


  fetch(baseUrl.value + 'api/auth/check_login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'sid': sid
    },
  }).then((res) => {
    return res.json()
  }).then(json => {
    if (json.code !== 200) {
      router.push('/login')
    }
  }).catch((e) => {
    router.push('/login')
  })


})


</script>

<template>
  <n-message-provider>
    <n-notification-provider>
      <n-config-provider :theme="darkTheme" class="w-full h-full">
        <router-view v-slot="{ Component,route }">

          <Transition :name="route.meta.transition">
            <component :is="Component"/>
          </Transition>

        </router-view>
      </n-config-provider>
    </n-notification-provider>
  </n-message-provider>
</template>

<style scoped></style>
