<script setup>
import {useRouter} from "vue-router";
import {onMounted, ref} from "vue";
import {createChain, getChainList} from "../services/chain.js";
import {useMenu, useSiderControl} from "../hooks/common.js";
import {NButton, NInput, NLayout, NLayoutContent, NLayoutSider, NMenu, NModal, useMessage} from "naive-ui";

const message = useMessage();

const {collapsed, expand, collapse} = useSiderControl();

const {} = useMenu([]);

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
  if (detail.value.dispatch) {
    detail.value.dispatch.fatherAnim("animate__slideInUp", "animate__fast");
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
  chainCreateForm.value.name = '';
  chainCreateForm.value.description = '';
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
  <div class="w-screen min-h-0 h-screen overflow-visible bg-black flex flex-col">
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
        {{ selectedChain.id === "none" ? "暂无" : selectedChain.name }}
      </div>
      <div class="w-1/3 h-full flex items-center">
        <n-button type="primary" size="small" @click="showAddChainModal = true" class=" ml-auto mr-2 bg-blue-500">
          创建调度
        </n-button>
      </div>
    </div>

    <div class="grow body overflow-hidden min-h-0 h-0 text-white flex bg-[#13151D]">
      <n-layout has-sider class="min-h-0 w-full h-[830px] overflow-hidden">
        <n-layout-sider
            class="h-full min-h-0"
            :native-scrollbar="false"
            bordered
            collapse-mode="width"
            :collapsed-width="1"
            :width="200"
            :collapsed="collapsed"
            show-trigger
            @collapse="collapse"
            @expand="expand"
        >
          <n-menu
              :collapsed="collapsed"
              :collapsed-width="1"
              :collapsed-icon-size="22"
              key-field="id"
              label-field="name"
              :options="chainList"
              :on-update-value="(key,item)=>chooseChain(item)"
          >
          </n-menu>
        </n-layout-sider>
        <n-layout-content class="overflow-auto h-full" content-style="padding: 14px;height:100%" :native-scrollbar="false">
          <RouterView v-slot="{ Component }">
              <component ref="detail" :key="$route.params.chainId" :is="Component"/>
          </RouterView>
        </n-layout-content>
      </n-layout>

      <!-- <div class="chains w-2/12 min-h-screen p-3 relative" ref="chainListRef">

        <button @click="" class="absolute right-0 top-1/2 -translate-y-1/2 h-10 w-10 bg-green-500 translate-x-1/2 z-10 rounded-full">
            收起
        </button>

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
        </div>
      </div>
      <div class="w-10/12 h-screen overflow-scroll">
        <RouterView v-slot="{ Component }">
          <keep-alive>
            <component ref="detail" :key="$route.params.chainId" :is="Component"/>
          </keep-alive>
        </RouterView>
      </div> -->
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
  @apply mid bg-[#203434] rounded border-[1px] border-green-500 p-1 text-green-400 transition hover:cursor-pointer h-14 mb-2;
}

.normal {
  @apply mid bg-[#13151D]   p-1 hover:bg-gray-800 border-[1px] border-transparent rounded hover:border-green-400 transition hover:cursor-pointer h-14 mb-2 border-gray-700;
}

.mid {
  @apply flex items-center justify-center;
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
    background-color: #05253a;
    color: transparent;
  }

  30% {
    background-color: #05253a;
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
