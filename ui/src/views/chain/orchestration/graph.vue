<script setup>
import {EdgeView, Graph, Path} from "@antv/x6";
import {Dnd} from "@antv/x6-plugin-dnd";
import {defineExpose, defineProps, onMounted, ref, watch} from "vue";
import {getTeleport, register} from "@antv/x6-vue-shape";
import ShortcutNode from "../../../components/graph/shortcut-node.vue";

Graph.registerEdge(
    'dag-edge',
    {
      inherit: 'edge',
      attrs: {
        line: {
          stroke: '#C2C8D5',
          strokeWidth: 1,
          targetMarker: null,
        },
      },
    },
    true,
)

Graph.registerConnector(
    'algo-connector',
    (s, e) => {
      const offset = 4
      const deltaY = Math.abs(e.y - s.y)
      const control = Math.floor((deltaY / 3) * 2)

      const v1 = {x: s.x, y: s.y + offset + control}
      const v2 = {x: e.x, y: e.y - offset - control}

      return Path.normalize(
          `M ${s.x} ${s.y}
       L ${s.x} ${s.y + offset}
       C ${v1.x} ${v1.y} ${v2.x} ${v2.y} ${e.x} ${e.y - offset}
       L ${e.x} ${e.y}
      `,
      )
    },
    true,
)

register({
  shape: "custom-node",
  width: 144,
  height: 56,
  component: ShortcutNode,
  ports: {
    groups: {
      top: {
        position: 'top',
        attrs: {
          circle: {
            r: 4,
            magnet: true,
            stroke: '#C2C8D5',
            strokeWidth: 1,
            fill: '#fff',
          },
        },
      },
      bottom: {
        position: 'bottom',
        attrs: {
          circle: {
            r: 4,
            magnet: true,
            stroke: '#C2C8D5',
            strokeWidth: 1,
            fill: '#fff',
          },
        },
      },
    },
  }
});

const TeleportContainer = getTeleport();

const emits = defineEmits(['onClickNode','onClickNodeWithResult'])

const props = defineProps(['dndContainer','editable','data','nodeExecResults']);
const graphRef = ref(null);

const g = ref(null);
const dnd = ref(null);


const checkEditable = () => {
  return (props.editable ?? false)
}


const initGraph = (data) => {
  g.value = new Graph({
    container: graphRef.value,
    autoResize: true,
    background: {
      color: "#000000",
    },
    panning: {
      enabled: true,
      eventTypes: ['leftMouseDown', 'mouseWheel'],
    },
    mousewheel: {
      enabled: true,
      modifiers: 'ctrl',
      factor: 1.1,
      maxScale: 1.5,
      minScale: 0.5,
    },
    highlighting: {
      magnetAdsorbed: {
        name: 'stroke',
        args: {
          attrs: {
            fill: '#fff',
            stroke: '#31d0c6',
            strokeWidth: 4,
          },
        },
      },
    },
    connecting: {
      snap: true,
      allowBlank: false,
      allowLoop: false,
      allowMulti: true,
      highlight: true,
      connector: 'rounded',
      connectionPoint: 'anchor',
      anchor: 'center',
      validateMagnet({magnet}) {
        return magnet.getAttribute('port-group') !== 'top'
      },
      validateConnection({sourceCell, targetCell, sourceMagnet, targetMagnet,}) {

        if (!checkEditable())return false

        // 不能连接自身
        if (sourceCell === targetCell) {
          return false
        }

        // 只能从 bottom 连接桩开始连接，连接到 top 连接桩
        if (
            !sourceMagnet ||
            sourceMagnet.getAttribute('port-group') === 'top'
        ) {
          return false
        }
        if (
            !targetMagnet ||
            targetMagnet.getAttribute('port-group') !== 'top'
        ) {
          return false
        }



        // 不能重复连线
        const edges = this.getEdges()


        const targetPortId = targetMagnet.getAttribute('port')
        const sourcePortId = sourceMagnet.getAttribute('port')

        if(edges.find(edge =>  edge.getTargetCellId() && edge.getSourceCellId() === sourceCell.id && edge.getSourcePortId() === sourcePortId))return false

        return !edges.find(edge => edge.getSourceCellId() === sourceCell.id && edge.getTargetCellId() === targetCell.id && edge.getSourcePortId() === sourcePortId && edge.getTargetPortId() === targetPortId)

      },
      createEdge() {
        return g.value.createEdge({
          shape: 'dag-edge',
          attrs: {
            line: {
              strokeDasharray: '5 5',
            },
            p1: {
              connection: true,
              fill: 'none',
              stroke: 'white',
              strokeWidth: 3,
              strokeLinejoin: 'round',
            },
            p2: {
              connection: true,
              fill: 'none',
              stroke: 'white',
              strokeWidth: 1,
              pointerEvents: 'none',
              strokeLinejoin: 'round',
              targetMarker: {
                tagName: 'path',
                fill: 'white',
                stroke: 'white',
                strokeWidth: 1,
                d: 'M 10 -3 10 -10 -2 0 10 10 10 3',
              },
            },
            sign: {
              x: -10,
              y: -7,
              width: 0,
              height: 15,
              stroke: 'black',
              fill: 'white',
              atConnectionLength: 30,
              strokeWidth: 1,
              event: 'click:rect',
              cursor: 'pointer',
            },
            signText: {
              atConnectionLength: 40,
              textAnchor: 'middle',
              textVerticalAnchor: 'middle',
              text: '',
              event: 'click:rect',
              cursor: 'pointer',
              fontSize: 10
            },
            c2: {
              r: 5,
              stroke: '#f5222d',
              fill: '#f5222d',
              atConnectionRatio: 0.68,
              strokeWidth: 1,
              cursor: 'pointer',
              event: 'click:redCircle',
            },
          },
          markup: [
            {
              tagName: 'path',
              selector: 'p1',
            },
            {
              tagName: 'path',
              selector: 'p2',
            },
            {
              tagName: 'rect',
              selector: 'sign',
            },
            {
              tagName: 'circle',
              selector: 'c2',
            },
            {
              tagName: 'text',
              selector: 'signText',
            },
          ],
          zIndex: -1,
        })
      },
    },
  });


  dnd.value = new Dnd({
    target: g.value,
    scaled: false,
    container: props.dndContainer || null,
  });

  g.value.on('node:added', ({node}) => {
    console.log(node.getData())
  })

  g.value.on('node:onClickDeleteBtn', (node) => {
    if(!checkEditable())return
    node.remove()
  })


  g.value.on('node:click',({node}) => {

    if (node.getData().nodeExecResult){
      emits('onClickNodeWithResult',node)
      return
    }

    emits('onClickNode',node)
  })

  g.value.on('click:redCircle',({edge}) => {
      if(!checkEditable())return
      edge.remove()
  })

  g.value.on('edge:connected', ({edge}) => {

    if (edge.getTargetPortId() !== 'In') {
      edge.remove()
      return;
    }

    if (edge.getSourcePortId() !== 'FailOut' && edge.getSourcePortId() !== 'SuccessOut') {
      edge.remove()
      return;
    }

    let colors = {
      'default': 'white',
      'FailOut': '#a1222c',
      'SuccessOut': 'green'
    }

    let texts = {
      'default': 'default',
      'FailOut': 'IF FAIL',
      'SuccessOut': 'IF OK'
    }

    edge.attr({
      line: {
        strokeDasharray: '',
        targetMarker: 'block',
        stroke: colors[edge.getSourcePortId()],
      },
    })


    edge.setAttrByPath('p2/stroke',colors[edge.getSourcePortId()])
    edge.setAttrByPath('p2/targetMarker/stroke',colors[edge.getSourcePortId()])
    edge.setAttrByPath('p2/targetMarker/fill',colors[edge.getSourcePortId()])
    edge.setAttrByPath('p1/stroke',colors[edge.getSourcePortId()])
    edge.setAttrByPath('sign/width',40)
    edge.setAttrByPath('sign/fill',colors[edge.getSourcePortId()])
    edge.setAttrByPath('signText/text', texts[edge.getSourcePortId()])

    edge.setRouter({
      name: 'manhattan',
      args: {},
    })
  })


  g.value.on('node:change:data', ({node}) => {
    const edges = g.value.getIncomingEdges(node)
    const data = node.getData()

    if(data.root === false)return

    g.value.getNodes().filter(item => item.id !== node.id).forEach(item => item.setData({root:false}))

    // edges?.forEach((edge) => {
    //   if (status === 'running') {
    //     edge.attr('line/strokeDasharray', 5)
    //     edge.attr('line/style/animation', 'running-line 30s infinite linear')
    //   } else {
    //     edge.attr('line/strokeDasharray', '')
    //     edge.attr('line/style/animation', '')
    //   }
    // })
  })


  g.value.fromJSON(props.data)
  if (props.nodeExecResults){
    updateNodeExecData(props.nodeExecResults,g.value)
  }
};



const updateNodeExecData = (execResults,graph) => {
  setTimeout(()=>{
    let nodes =  graph.getNodes()
    if(nodes.length === 0){
      updateNodeExecData(execResults,graph)
      return
    }

    console.log(execResults)

    nodes.filter(node => execResults.filter(res => res.nodeId === node.id).forEach(res => {
      node.setData({
        nodeExecResult:res
      })
    }))

  },100)
}


watch(() => props.data,(value, oldValue) => {
    g.value.fromJSON(value)
})

watch(() => props.nodeExecResults,(nv,ov) => {
  updateNodeExecData(nv,g.value)
})

/**
 * @param {MouseEvent} e
 */
const startDrag = (e, option) => {

  if(!checkEditable())return

  const node = g.value.createNode({
    shape: "custom-node",
    x: 500,
    y: 500,
    data: option,
    ports: [
      {
        "id": "In",
        "group": "top"
      },
      {
        "id": "SuccessOut",
        "group": "bottom"
      },
      {
        "id": "FailOut",
        "group": "bottom"
      },
    ]
  });

  dnd.value.start(node, e);
};


const exportData = () => {
  return g.value.toJSON()
}

defineExpose({
  startDrag,
  exportData,
  g
});

onMounted(() => {
  initGraph(props.data);
});
</script>

<template>
  <div class="w-full h-full">
    <div ref="graphRef"></div>
    <TeleportContainer/>
  </div>
</template>


<style scoped>


.node {
  display: flex;
  align-items: center;
  width: 100%;
  height: 100%;
  background-color: #fff;
  border: 1px solid #c2c8d5;
  border-left: 4px solid #5F95FF;
  border-radius: 4px;
  box-shadow: 0 2px 5px 1px rgba(0, 0, 0, 0.06);
}

.node img {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
  margin-left: 8px;
}

.node .label {
  display: inline-block;
  flex-shrink: 0;
  width: 104px;
  margin-left: 8px;
  color: #666;
  font-size: 12px;
}

.node .status {
  flex-shrink: 0;
}

.node.success {
  border-left: 4px solid #52c41a;
}

.node.failed {
  border-left: 4px solid #ff4d4f;
}

.node.running .status img {
  animation: spin 1s linear infinite;
}

.x6-node-selected .node {
  border-color: #1890ff;
  border-radius: 2px;
  box-shadow: 0 0 0 4px #d4e8fe;
}

.x6-node-selected .node.success {
  border-color: #52c41a;
  border-radius: 2px;
  box-shadow: 0 0 0 4px #ccecc0;
}

.x6-node-selected .node.failed {
  border-color: #ff4d4f;
  border-radius: 2px;
  box-shadow: 0 0 0 4px #fedcdc;
}

.x6-edge:hover path:nth-child(2) {
  stroke: #1890ff;
  stroke-width: 1px;
}

.x6-edge-selected path:nth-child(2) {
  stroke: #1890ff;
  stroke-width: 1.5px !important;
}

@keyframes running-line {
  to {
    stroke-dashoffset: -1000;
  }
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

</style>