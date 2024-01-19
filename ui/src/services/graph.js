

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
