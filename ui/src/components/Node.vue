<script setup>
import {h, onMounted, reactive, ref, watch} from "vue";
import {NButton, NDropdown, NPopconfirm, NTreeSelect, useMessage} from "naive-ui";
import {setRoot} from "../services/chain.js";
import {useNodeOperation} from "../hooks/chain.js";

const message = useMessage();

const {DelNode, BindShortcut, UnbindShortcut, LkNode, UnlkNode, stRoot} = useNodeOperation()

const emit = defineEmits(["onDelete", "onSetRoot"]);

const props = defineProps({
  // {
  //     "probe": {
  //         "id": "c0fc4aa9-d454-4038-b30a-09d9787e4c89",
  //         "name": "Node001",
  //         "chainId": "e2e0c7fe-52c2-4171-887b-1da1ad4fa1bb",
  //         "description": "the 1st probe description in the world",
  //         "root": false
  //     },
  //     "shortcut": {
  //         "id": "14af9713-4337-4e5e-8b20-841398a6fadd",
  //         "name": "test",
  //         "description": "001",
  //         "type": 0,
  //         "createTime": "2024-01-10T21:05:23.208881+08:00",
  //         "timeout": 1000,
  //         "justRun": false,
  //         "payload": "ls",
  //         "probeId": "mac-os"
  //     },
  //     "successThenId": "b0ab7ef7-2588-494e-b05a-6365dfdb5305",
  //     "successThenName": "Node002",
  //     "failThenId": "959231cb-c157-4526-8429-ea7f99cd9a51",
  //     "failThenName": "Node003"
  // },
  item: {
    type: Object,
    required: true,
  },
  // {
  //     "mac-os": [
  //       {
  //         "id": "cb457915-053f-4218-b567-e5503b2086e1",
  //         "name": "8081端口占用",
  //         "description": "查看8081端口占用",
  //         "type": 0,
  //         "createTime": "2024-01-12T00:11:54.429639+08:00",
  //         "timeout": 1000,
  //         "justRun": false,
  //         "payload": "lsof -i :8081",
  //         "probeId": "mac-os"
  //       }
  //     ]
  //   }
  shortcutGroup: {
    type: Object,
    default: () => {
    },
  },
  allNodes: {
    type: Array,
    default: () => [],
  },
});

// 捷径树形选择
const shortcutTreeSelect = ref([]);

// 选中的捷径
const selectedShortcut = reactive({id: "none", name: "暂无"});

const dropDownOptions = ref([
  {
    render: () => h(
        'div',
        {
          'class': 'm-1 h-8 pt-1 flex justify-center item-center  border-[1px] border-transparent  hover:bg-[#59595D] text-green-500 transition rounded hover:cursor-pointer ',
          'onClick': () => setSelfRoot()
        },
        {
          default: () => '设置为ROOT'
        }
    ),
    key: "setRoot",
    type: 'render'
  },
  {
    render: () => h(
        NPopconfirm,
        {
          'class': 'text-hover',
          'onPositiveClick': () => deleteSelf(),
          'positiveText': '确定',
          'negativeText': '取消',
        },
        {
          trigger: () => h(
              'div',
              {
                'class': 'm-1 h-8 pt-1 flex justify-center item-center  border-[1px] border-transparent  hover:bg-[#59595D] text-red-500 transition rounded hover:cursor-pointer ',
              },
              {
                default: () => '删除'
              }
          ),
          default: () => '你确定要删除吗？'
        }
    ),
    key: "delete",
    type: 'render'
  },
]);

// 选中捷径触发的函数
const onShortcutSelect = async (value, meta) => {
  // 如果选中的捷径为null，且当前捷径不为none，则解绑捷径
  if (value === null && selectedShortcut.id !== "none") {
    let success = await UnbindShortcut(props.item.node.id, selectedShortcut.id);
    if (success) {
      selectedShortcut.id = "none";
      selectedShortcut.name = "暂无";
    }
  }

  // 绑定捷径
  let success = await BindShortcut(props.item.node.id, value);

  if (success) {
    selectedShortcut.id = value;
    selectedShortcut.name = meta.label;
  }
};

const otherNodesTreeSelect = ref([]);
const successThen = reactive({id: "none", name: "暂无"});
const failThen = reactive({id: "none", name: "暂无"});

const handleNodeSelection = async (value, meta, node, type) => {
  if (value === null && node.id !== "none") {
    let success = await UnlkNode(props.item.node.id, node.id, props.item.node.chainId, type);
    if (success) {
      node.id = "none";
      node.name = "暂无";
    }
    return;
  }

  let success = await LkNode(props.item.node.id, value, props.item.node.chainId, type);
  if (success) {
    node.id = value;
    node.name = meta.label;
  }
};

const onSuccessThenSelect = (value, meta) => handleNodeSelection(value, meta, successThen, 0);
const onFailThenSelect = (value, meta) => handleNodeSelection(value, meta, failThen, 1);
/**
 * 将快捷方式组对象转换为树形选择数据结构。
 *
 * @param {Object} shortcutGroup - 要转换的快捷方式组对象。
 * @returns {Array} - 转换后的树形选择数据结构。
 */
const toTreeSelect = (shortcutGroup) => {
  let res = [];
  for (let key in shortcutGroup) {
    let children = [];
    for (let item of shortcutGroup[key]) {
      children.push({
        label: item.name,
        key: item.id,
      });
    }
    res.push({
      label: key,
      key: key,
      disabled: true,
      children: children,
    });
  }
  return res;
};


watch(() => props.shortcutGroup,
    () => shortcutTreeSelect.value = toTreeSelect(props.shortcutGroup)
)


watch(
    () => props.allNodes,
    () => filterNodes()
);

const filterNodes = () => {
  otherNodesTreeSelect.value = props.allNodes
      .filter((item) => item.id !== props.item.node.id)
      .map((item) => {
        return {
          label: item.name,
          key: item.id,
        };
      });
};

const deleteSelf = async () => {
  let success = await DelNode(props.item.node.id);
  if (success) {
    emit("onDelete", props.item.node.id);
  }
};


const setSelfRoot = async () => {
  let success = await stRoot(props.item.node.id)
  success && emit("onSetRoot", props.item.node.id);
};


onMounted(() => {
  shortcutTreeSelect.value = toTreeSelect(props.shortcutGroup);
  filterNodes();

  if (props.item.successThenId !== null) {
    successThen.id = props.item.successThenId;
    successThen.name = props.item.successThenName;
  }

  if (props.item.failThenId !== null) {
    failThen.id = props.item.failThenId;
    failThen.name = props.item.failThenName;
  }

  if (props.item.shortcut !== null) {
    selectedShortcut.id = props.item.shortcut.id;
    selectedShortcut.name = props.item.shortcut.name;
  }
});
</script>

<template>
  <div class="relative">
    <div
        class="w-60 h-72 overflow-hidden shadow-black shadow-md relative text-white flex items-center justify-center"
    >
      <div
          v-if="item.node.root"
          class="w-full h-full absolute left-0 top-0  z-0 text-3xl flex items-center justify-center text-green-500 opacity-20">
        ROOT
      </div>
      <!-- <div class="small-ball bg-green-500 right-0 top-1/4 translate-x-1/2">
            </div>
            <div class="small-ball bg-red-500 border-2 border-black right-0 bottom-1/4 translate-x-1/2">
            </div>
            <div class="small-ball bg-white left-0 top-[50%] -translate-y-1/2 -translate-x-1/2">
            </div> -->
      <div class="w-full h-full bg-transparent p-4 flex flex-col z-10">
        <div class="flex items-center w-full">
          <span>📄:{{ props.item.node.name }}</span>

          <n-dropdown
              :options="dropDownOptions"
              trigger="click"
          >
            <n-button
                size="small"
                class="ml-auto text-white"
            >...
            </n-button>
          </n-dropdown>
        </div>
        <div
            class="h-fit w-full whitespace-nowrap text-ellipsis overflow-hidden"
        >
          {{ props.item.node.description }}
        </div>
        <div class="flex flex-col mb-2 mt-2">
          <span class="text-green-500">成功后：</span>
          <n-tree-select
              v-model:value="successThen.id"
              :default-value="successThen.id"
              :options="otherNodesTreeSelect"
              placeholder="成功后"
              :on-update-value="onSuccessThenSelect"
              clearable
              size="medium"
          />
        </div>
        <div class="mb-auto flex flex-col">
          <span class="text-red-500">失败后：</span>
          <n-tree-select
              v-model:value="failThen.id"
              :default-value="failThen.id"
              :options="otherNodesTreeSelect"
              placeholder="失败后"
              :on-update-value="onFailThenSelect"
              clearable
              size="medium"
              style="width: 100%"
          />
        </div>
        <div>
          <n-tree-select
              v-model:value="selectedShortcut.id"
              :default-value="selectedShortcut.id"
              :options="shortcutTreeSelect"
              placeholder="请选择快捷方式"
              :on-update-value="onShortcutSelect"
              clearable
              size="medium"
              style="width: 100%"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.small-ball {
  @apply absolute w-7 h-7 rounded-full z-0 absolute;
}

.text-hover {
  @apply hover:text-blue-700 hover:cursor-pointer hover:underline;
}
</style>
