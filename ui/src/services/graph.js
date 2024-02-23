

class MyNode {
    constructor(id, label, type, x, y, size) {
        this.id = id;
        this.label = label;
        this.type = type;
        this.x = x;
        this.y = y;
        this.size = size;
        this.linkPoints = {
            top: true,
            bottom: true,
            left: true,
            right: true,
            size: 5
        },
        this.anchorPoints = []
    }

    // 锚点
    addAnchorPoints(x,y) {
        this.anchorPoints.push([x,y])
    }

}


export const genCommonNode = (node,x,y) => {
    console.log(node)
    let n =  new MyNode(node.id, node.name, 'modelRect', x, y, [150, 70]);
    n.addAnchorPoints(1,0.5)
    n.addAnchorPoints(1,1)
    n.addAnchorPoints(0.5,1)
    n.addAnchorPoints(0,0.5)
    return  n
}


export const genShortcutNode = (shortcut,x,y) => {
    return new MyNode(shortcut.id, shortcut.name, 'shortcutRect', x, y, [100, 100]);
}

/**
 *
 * @param {Object} originData
 * @param {Array<Object>} originData.cells
 * @param {string} chainId
 */
export const processGraphDataToChainData = (originData,chainId) => {

    const graphData = originData['cells']

    const result = {
        chainId:chainId,
        nodes:[],
        edges:[],
        originData: JSON.stringify(graphData)
    }


    const processNodeData = (nodeData) => {
        const shortcut = nodeData.data.proto
        const root = nodeData.data.root ?? false
        result.nodes.push({
            id:nodeData.id,
            name: shortcut.name,
            description: shortcut.description,
            shortcut:shortcut,
            root:root
        })
    }

    const processEdgeData = (edgeData) => {

        result.edges.push({
            id: edgeData.id,
            sourceId: edgeData.source.cell,
            targetId: edgeData.target.cell,
            type: edgeData.source.port === 'SuccessOut'? 0 : edgeData.source.port === 'FailOut'? 1 : -1
        })

    }

    graphData.forEach(item => item.shape === 'dag-edge'? processEdgeData(item) : processNodeData(item))


    return result

}