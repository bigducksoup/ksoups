<script setup>
import G6 from "@antv/g6";
import { onMounted, ref } from "vue";
import { successEdge, failEdge, processInfo } from "../../graph/custom.js";

const props = defineProps(["chainInfo"]);

const root = ref(null);
const data = ref({
  nodes: [
    {
      info: {
        node: {
          name: "node1",
          root: false,
          description: "node1",
          id: "node1",
          chainId: "chain1",
        },
        shortcut: {
          id: "shortcut1",
          name: "shortcut1",
          description: "shortcut1",
          probeId: "probe1",
          payload: "payload1",
          justRun: false,
          timeout: 1000,
          type: 1,
        },
      },
      error: false,
      id: "node1",
      x: 100,
      y: 100,
    },
    {
      info: {
        node: {
          name: "node2",
          root: false,
          description: "node2",
          id: "node2",
          chainId: "chain2",
        },
        shortcut: {
          id: "shortcut2",
          name: "shortcut2",
          description: "shortcut3",
          probeId: "probe2",
          payload: "payload2",
          justRun: false,
          timeout: 1000,
          type: 1,
        },
      },
      error: false,
      id: "node2",
      x: 400,
      y: 100,
    },
    {
      info: {
        node: {
          name: "node3",
          root: false,
          description: "node3",
          id: "node3",
          chainId: "chain3",
        },
        shortcut: {
          id: "shortcut3",
          name: "shortcut3",
          description: "shortcut3",
          probeId: "probe3",
          payload: "payload3",
          justRun: false,
          timeout: 1000,
          type: 1,
        },
      },
      error: false,
      id: "node3",
      x: 400,
      y: 300,
    },
  ],
  edges: [
    successEdge({
      id: "edge0",
      source: "node1",
      target: "node2",
      curveOffset: 0,
    }),
    failEdge({
      id: "edge1",
      source: "node1",
      target: "node3",
      curveOffset: 0,
    }),
    successEdge({
      id: "edge2",
      source: "node2",
      target: "node3",
      curveOffset: 0,
    }),
    successEdge({
      id: "edge3",
      source: "node2",
      target: "node1",
      curveOffset: 60,
    }),
  ],
});
onMounted(() => {
  let { nodes, edges } = processInfo(props.chainInfo);

  data.value = {
    nodes: nodes,
    edges: edges,
  };


  const graph = new G6.Graph({

    container: root.value, // String | HTMLElement，必须，在 Step 1 中创建的容器 id 或容器本身
    width: root.value.scrollWidth || 1920 , // Number，必须，图的宽度
    height: root.value.scrollHeight || 1080,
    modes: {
      default: ["drag-canvas", "drag-node"],
    },
    defaultNode: {
      type: "card-node",
    },
    linkCenter: true,
  });

  graph.data(data.value); // 读取 Step 2 中的数据源到图上
  graph.render(); // 渲染图

  let container = root.value;

  if (typeof window !== "undefined"){
    window.onresize = () => {
      if (!graph || graph.get("destroyed")) return;
      if (!container || !container.scrollWidth || !container.scrollHeight)
        return;
      graph.changeSize(container.scrollWidth, container.scrollHeight);
    };
  }

});
</script>

<template>
  <div ref="root" class="w-screen h-screen flex p-3"></div>
</template>

<style scoped></style>
