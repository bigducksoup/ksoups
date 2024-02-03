<script setup>
import {
  NButton,
  NDrawer,
  NDrawerContent,
  NForm,
  NFormItem,
  NIcon,
  NInput,
  NRadioButton,
  NRadioGroup,
  useMessage,
} from "naive-ui";
import {Add, ChevronBack, Refresh} from "@vicons/ionicons5";
import {GetSSHGroupTree, useSSH} from "../../hooks/ssh.js";
import {onMounted, ref} from "vue";
import GroupItem from "../../components/ssh/GroupItem.vue";
import SSHInfoItem from "../../components/ssh/SSHInfoItem.vue";
import {useRouter} from "vue-router";
import {useForm} from "../../hooks/common";

const router = useRouter();
const message = useMessage();

const lastGroupIds = ref([]);
const curGroupId = ref("root");
const {GetSSHGroupContent, groupInfo, SaveSSHGroup, SaveSSHInfo} = useSSH();

const onClickItem = (item) => {
  if (item.type === 0) {
    lastGroupIds.value.push(item.payload.parent);
    GetSSHGroupContent(item.payload.id);
    curGroupId.value = item.payload.id;
    return;
  }

  if (item.type === 1) {
    const {href} = router.resolve({
      name: 'index',
      query: {
        sshId: item.payload.id,
      },
    });
    // open and set title
    let W = window.open(href, "_blank");
    W.onload = () => {
      W.document.title = "SSH-" + item.payload.addrPort;
    };
  }
};

const back = () => {
  let id = lastGroupIds.value.pop();
  if (!id) return;
  GetSSHGroupContent(id);
  curGroupId.value = id;
};

// 是否显示抽屉
const showCreateDrawer = ref(false);

// 创建类型 0:创建组 1:创建SSH信息
const createType = ref(0);

// 组名
const groupName = ref("");

/**
 * 保存组
 * 检查组名是否为空
 */
const handleSaveGroup = () => {
  if (!groupName.value) {
    message.error("组名不能为空");
    return;
  }

  SaveSSHGroup(groupName.value, curGroupId.value).then((success) => {
    if (success) {
      GetSSHGroupContent(curGroupId.value);
      showCreateDrawer.value = false;
      groupName.value = "";
    }
  });
};

/**
 * 保存SSH信息
 * 检查表单是否通过
 * 保存
 */
const {form: SaveSSHInfoForm, Submit, Clear} = useForm(
    {
      address: "",
      port: "",
      username: "",
      password: "",
    },
    (form) => {
      let info = {
        addrPort: form.address + ":" + form.port,
        username: form.username,
        password: form.password,
        groupId: curGroupId.value,
      };
      SaveSSHInfo(info).then((success) => {
        if (success) {
          GetSSHGroupContent(curGroupId.value)
          showCreateDrawer.value = false;
          Clear()
        }
      });
    }
);

const formRef = ref(null);

const rules = {
  address: {
    required: true,
    message: "地址不能为空",
  },
  port: {
    required: true,
    message: "端口不能为空",
  },
  username: {
    required: true,
    message: "用户名不能为空",
  },
  password: {
    required: true,
    message: "密码不能为空",
  },
};

const handleValidateClick = (e) => {
  e.preventDefault();
  formRef.value?.validate((errors) => {
    if (!errors) {
      Submit();
    }
  });
};


const groupTree = ref([])

const GetGroupTree = () => {
  GetSSHGroupTree().then((tree) => {
    groupTree.value = tree;
  });
}

onMounted(() => {
  GetGroupTree()
})


</script>

<template>
  <div class="bg-[#121417] w-full h-full text-white p-2 overflow-x-scroll">
    <div class="op-bar w-full h-10 flex items-center pl-2 mb-2 gap-2">
      <n-button circle @click="back">
        <n-icon size="20">
          <ChevronBack/>
        </n-icon>
      </n-button>
      <n-button circle @click="()=>{
        GetSSHGroupContent(curGroupId)
      }">
        <n-icon size="20">
          <Refresh/>
        </n-icon>
      </n-button>

      <n-button circle @click="showCreateDrawer = true">
        <n-icon size="20">
          <Add/>
        </n-icon>
      </n-button>
    </div>

    <div class="w-full h-fit border-[1px] border-gray-700">
      <div v-for="item in groupInfo" @click="onClickItem(item)">
        <group-item v-if="item.type === 0" :GItem="item.payload" :group-tree="groupTree"
                    @onDelete="()=>{
                      GetSSHGroupContent(curGroupId)
                      GetGroupTree()
                    }"
                    @onUpdate="()=>{
                      GetSSHGroupContent(curGroupId)
                      GetGroupTree()
                    }"></group-item>
        <SSHInfoItem v-if="item.type === 1" :info="item.payload" :group-tree="groupTree"
                     @onDelete="()=>{
                       GetSSHGroupContent(curGroupId)
                     }"
                     @onUpdate="()=>{
                       GetSSHGroupContent(curGroupId)
                     }"></SSHInfoItem>
      </div>
    </div>
  </div>

  <n-drawer v-model:show="showCreateDrawer" width="500px">
    <n-drawer-content title="创建">
      <n-radio-group v-model:value="createType" class="mb-4">
        <n-radio-button key="0" :value="0">创建组</n-radio-button>
        <n-radio-button key="1" :value="1">创建SSH</n-radio-button>
      </n-radio-group>
      <n-form v-if="createType === 0">
        <n-form-item label="组名称">
          <n-input v-model:value="groupName" placeholder="组名"/>
        </n-form-item>
      </n-form>

      <n-form v-else :model="SaveSSHInfoForm" :rules="rules" ref="formRef">
        <n-form-item label="地址" path="address">
          <n-input v-model:value="SaveSSHInfoForm.address" placeholder="地址"/>
        </n-form-item>
        <n-form-item label="端口" path="port">
          <n-input v-model:value="SaveSSHInfoForm.port" placeholder="端口"/>
        </n-form-item>
        <n-form-item label="用户名" path="username">
          <n-input
              v-model:value="SaveSSHInfoForm.username"
              placeholder="用户名"
          />
        </n-form-item>
        <n-form-item label="密码" path="password">
          <n-input
              v-model:value="SaveSSHInfoForm.password"
              placeholder="密码"
              type="password"
          />
        </n-form-item>
      </n-form>

      <template #footer>
        <n-button
            @click="handleSaveGroup"
            v-if="createType === 0"
            type="success"
            class="bg-green-500"
        >
          创建
        </n-button>
        <n-button
            v-else
            @click="handleValidateClick"
            type="success"
            class="bg-green-500"
        >确认
        </n-button
        >
      </template>
    </n-drawer-content>
  </n-drawer>


</template>

<style scoped></style>
