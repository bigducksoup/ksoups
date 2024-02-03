<script setup>
import {useRouter} from "vue-router";
import {onMounted, ref} from "vue";
import {createChain, getChainList} from "../services/chain.js";
import {Add} from "@vicons/ionicons5";
import {NButton, NIcon, NInput, NModal, useMessage} from "naive-ui";

const message = useMessage();

const router = useRouter();
const headerStyle = ref({
  animation: "full-to-normal 0.50s ease-in-out",
});

const detail = ref(null);

const selectedChain = ref({
  id: "none",
});

const showAddChainModal = ref(false);

const chainCreateForm = ref({
  name: "",
  description: "",
});

const chainList = ref([]);

const chooseChain = (chain) => {
  selectedChain.value = chain;
  router.push("/chain/" + selectedChain.value["id"]);
  if (detail.value.dispatch){
    detail.value.dispatch.fatherAnim("animate__slideInUp", "animate__fast")
  }
};

const CreateChain = async () => {
  let res = await createChain(
      chainCreateForm.value.name,
      chainCreateForm.value.description
  );
  if (res.code !== 200) {
    message.error(res.msg);
    return;
  }
  message.success("创建成功");
  chainList.value.push(res.data);
  showAddChainModal.value = false;
};

const backToMain = () => {
  headerStyle.value.animation = "none";
  requestAnimationFrame(() => {
    headerStyle.value.animation = "full-to-normal 0.70s ease-in-out reverse";
  });
  setTimeout(() => {
    router.push("/");
  }, 650);
};

const init = async () => {
  let res = await getChainList();
  chainList.value = res.data;
};

onMounted(() => {
  init();
});
</script>

<template>
  <div class="w-screen min-h-screen bg-black flex flex-col">
    <div class="header" :style="headerStyle">
      <div class="w-1/3 h-full items-center flex pl-5">
        <span
            @click="backToMain"
            class="hover:cursor-pointer hover:text-white hover:underline transition"
        >
          返回首页
        </span>
      </div>
      <div class="w-1/3 h-full flex items-center justify-center text-white">
        {{selectedChain.id === 'none' ? '暂无' : selectedChain.name}}
      </div>
      <div class="w-1/3 h-full"></div>
    </div>

    <div class="body text-white w-full grow flex bg-[#13151D]">
      <!-- 左边chain列表 -->
      <div class="chains w-2/12 min-h-screen p-3">
        <div
            @click="showAddChainModal = true"
            class="rounded group hover:border-green-400 hover:border-[1px] border-[1px] border-gray-700 transition hover:cursor-pointer h-14 mb-2 flex items-center justify-center"
        >
          <div
              class="w-10 h-10 border-[1px] border-gray-700 rounded-full flex items-center justify-center transition group-hover:bg-green-500 group-hover:text-black"
          >
            <n-icon size="25">
              <add/>
            </n-icon>
          </div>
        </div>
        <div
            @click="chooseChain(item)"
            v-for="item in chainList"
            :class="selectedChain.id === item.id ? 'choosen' : 'normal'"
        >
          {{ item.name }}
          <br/>
          {{ item.description }}
        </div>
      </div>
      <!-- 右侧展示区 -->
      <div class="w-10/12 h-screen overflow-scroll">
        <RouterView v-slot="{ Component }">
          <keep-alive>
            <component ref="detail" :key="$route.params.chainId" :is="Component"/>
          </keep-alive>
        </RouterView>
      </div>
    </div>
    <n-modal v-model:show="showAddChainModal" title="添加Chain">
      <div class="w-[400px] h-[300px] rounded p-3 bg-gray-700 flex flex-col">
        名称:
        <n-input
            placeholder="输入Chain名称"
            v-model:value="chainCreateForm.name"
        ></n-input>
        <div class="h-3"></div>
        描述:
        <n-input
            placeholder="输入Chain描述"
            v-model:value="chainCreateForm.description"
        ></n-input>
        <div class="mb-auto"></div>
        <n-button @click="CreateChain">创建</n-button>
      </div>
    </n-modal>
  </div>
</template>

<style>
.choosen {
  @apply bg-[#203434] rounded border-[1px] border-green-500 p-1 text-green-400 transition hover:cursor-pointer h-14 mb-2;
}

.normal {
  @apply bg-[#13151D]  p-1 hover:bg-gray-800 border-[1px] border-transparent rounded hover:border-green-400 transition hover:cursor-pointer h-14 mb-2;
}

.header {
  @apply w-full h-10 bg-green-500 z-10 flex;
}

.anim {
  animation: full-to-normal 0.6s ease-in-out;
}

.anim-reverse {
  animation: full-to-normal 0.6s ease-in-out reverse;
}

@keyframes full-to-normal {
  0% {
    height: 100vh;
    background-color: #05253A;
    color: transparent;
  }

  30% {
    background-color: #05253A;
    color: transparent;
  }

  80% {
    height: 2.2rem;
    --tw-bg-opacity: 1;
    background-color: rgb(34 197 94 / var(--tw-bg-opacity));
    color: transparent;
  }
  100% {
    height: 2.5rem;
    --tw-bg-opacity: 1;
    background-color: rgb(34 197 94 / var(--tw-bg-opacity));
  }
}
</style>
