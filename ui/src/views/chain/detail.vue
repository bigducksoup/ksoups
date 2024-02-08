<script setup>
import { useRoute } from "vue-router";
import Node from "../../components/Node.vue";
import { onMounted, reactive, ref } from "vue";
import { Add, Close } from "@vicons/ionicons5";
import { NButton, NIcon, NInput, NModal, useMessage } from "naive-ui";
import { useChainInfo, useNodeOperation } from "../../hooks/chain.js";
import { useShortcutGroup } from "../../hooks/shortcut.js";
import DispatchConsole from "./dispatch-console.vue";
import NodeManage from "./node-manage.vue";

const message = useMessage();
const route = useRoute();

const ctl = ref(null);
const manage = ref(null);

const showNodes = ref(false);

const showMode = ref(true);

const showAddNodeModal = ref(false);

const dispatch = ref(null);

defineExpose({
  dispatch,
});


// TODO 适当时机更新chainDetail
const { chainDetail, nodes, getChainData } = useChainInfo(route.params.chainId);

const { shortcutGroup } = useShortcutGroup();

const { NewNode } = useNodeOperation();

const nodeCreateForm = reactive({
  chainId: "",
  name: "",
  description: "",
});

/**
 * 创建节点
 *
 * @param {string} chainId - 链ID
 * @param {string} name - 节点名称
 * @param {string} description - 节点描述
 * @returns {Promise} - 返回创建节点的结果
 */
const CreateNode = async () => {
  if (nodeCreateForm.name === "") {
    message.warning("Node名称不能为空");
    return;
  }

  let chainId = chainDetail.value.chain.id;
  let name = nodeCreateForm.name;
  let description = nodeCreateForm.description;

  let success = await NewNode(chainId, name, description);

  if (success) {
    showAddNodeModal.value = false;
    nodeCreateForm.name = "";
    nodeCreateForm.description = "";
    await getChainData();
  }
};

const onDeleteNode = (nodeId) => {
  chainDetail.value.nodes = chainDetail.value.nodes.filter(
    (item) => item.node.id !== nodeId
  );
};

const onSetRoot = (nodeId) => {
  getChainData();
};

const showNodeManage = () => {
  showNodes.value = true;
};
</script>

<template>
  <div class="p-2 pb-16 bg-black h-screen overflow-hidden relative">
    <div ref="ctl" class="h-screen w-full p-3 z-10">
      <dispatch-console
        ref="dispatch"
        :show-nodes="showNodeManage"
        :chain-id="route.params.chainId"
      ></dispatch-console>
    </div>
  </div>

  <Transition name="bounce">

      <div
      v-if="showNodes"
        ref="manage"
        class="h-screen w-screen left-0 flex p-1 absolute top-0 z-20 bg-[#29292A]"
      >
        <div class="flex flex-col mt-1 w-full overflow-scroll">
          <div class="w-full p-2 flex items-center text-3xl border-b-[1px] border-gray-600">
            NODES:
            <n-button
              type="info"
              class=" text-white ml-auto"
              @click="showMode = !showMode"
              
            >
              {{ !showMode ? "切换到编辑模式" : "切换到展示模式"  }}
            </n-button>
            <n-button
              @click="showAddNodeModal = true"
              class="bg-green-500 ml-3"
            >
              <n-icon size="20">
                <add />
              </n-icon>
            </n-button>

            <n-button
              @click="showNodes = false"
              circle
              class="ml-3 bg-gray-500"
            >
              <n-icon size="20">
                <close />
              </n-icon>
            </n-button>

          </div>
          <div
            class="w-full rounded shadow flex flex-row gap-2 items-center flex-wrap pl-2"
          >
            <node
              v-show="showMode"
              :key="item.node.id"
              v-for="item in chainDetail.nodes"
              :item="item"
              :shortcut-group="shortcutGroup"
              :all-nodes="nodes"
              @on-delete="onDeleteNode"
              @on-set-root="onSetRoot"
            >
            </node>

            <node-manage v-show="!showMode" :chain-info="chainDetail"></node-manage>
          </div>
        </div>
        <n-modal v-model:show="showAddNodeModal">
          <div
            class="w-[400px] h-[300px] rounded p-3 bg-gray-700 flex flex-col"
          >
            名称:
            <n-input
              placeholder="输入Node名称"
              v-model:value="nodeCreateForm.name"
            ></n-input>
            <div class="h-3"></div>
            描述:
            <n-input
              placeholder="输入Node描述"
              v-model:value="nodeCreateForm.description"
            ></n-input>
            <div class="mb-auto"></div>
            <n-button @click="CreateNode">创建</n-button>
          </div>
        </n-modal>
      </div>
  </Transition>
</template>

<style scoped>
h1 {
  color: blue;
}
</style>
