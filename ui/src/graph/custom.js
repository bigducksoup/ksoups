import G6 from "@antv/g6";
import { toRaw } from "vue";


export const init = () => {
    G6.registerNode(
        'card-node',
        {
            drawShape: function drawShape(cfg, group) {


                const shortcut = cfg.info.shortcut ?? {name:'',probeId:''}

                const color = cfg.error ? '#F4664A' : '#30BF78';
                const r = 2;
                const shape = group.addShape('rect', {
                    attrs: {
                        x: 0,
                        y: 0,
                        width: 200,
                        height: 160,
                        stroke: color,
                        radius: r,
                        fill: 'black'
                    },
                    // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
                    name: 'main-box',
                    draggable: true,
                });


                group.addShape('rect',{
                    attrs: {
                        x: 0,
                        y: 0,
                        width: 200,
                        height: 40,
                        fill: color,
                        radius: [r, r, 0, 0],
                    },
                    name: 'title-box',
                    draggable: true,
                })


                group.addShape('text',{
                    attrs: {
                        textBaseline: 'top',
                        y: 2,
                        x: 2,
                        lineHeight: 40,
                        text: cfg.info.node.name,
                        fill: '#fff',
                        fontSize:30
                    },
                    name: 'title',
                    draggable: true,
                })

                group.addShape('text',{
                    attrs: {
                        textBaseline: 'top',
                        y: 50,
                        x: 2,
                        lineHeight: 20,
                        text: '指令:' + shortcut.name,
                        fill: '#fff',
                        fontSize:20
                    },
                    name: 'description',
                    draggable: true,
                })


                group.addShape('text',{
                    attrs: {
                        textBaseline: 'top',
                        y: 130,
                        x: 2,
                        lineHeight: 20,
                        text:  '在 ' + shortcut.probeId + ' 上运行',
                        fill: '#A4D8FF',
                        fontSize:16,
                    },
                    name: 'description',
                    draggable: true,
                })


                return shape
            },
        },
        'single-node',
    );




    
    G6.registerEdge(
        'extra-shape-edge',
        {
          afterDraw(cfg, group) {
            // get the first shape in the graphics group of this edge, it is the path of the edge here
            // 获取图形组中的第一个图形，在这里就是边的路径图形
            const shape = group.get('children')[0];
            // get the coordinate of the mid point on the path
            // 获取路径图形的中点坐标
            const midPoint = shape.getPoint(0.5);
            const rectColor = cfg.midPointColor || '#333';
            // add a rect on the mid point of the path. note that the origin of a rect shape is on its lefttop
            // 在中点增加一个矩形，注意矩形的原点在其左上角
            group.addShape('rect', {
              attrs: {
                width: 10,
                height: 10,
                fill: rectColor || '#333',
                // x and y should be minus width / 2 and height / 2 respectively to translate the center of the rect to the midPoint
                // x 和 y 分别减去 width / 2 与 height / 2，使矩形中心在 midPoint 上
                x: midPoint.x - 5,
                y: midPoint.y - 5,
              },
            });
      
            // get the coordinate of the quatile on the path
            // 获取路径上的四分位点坐标
            const quatile = shape.getPoint(0.65);
            const quatileColor = cfg.quatileColor || '#333';
            // add a circle on the quatile of the path
            // 在四分位点上放置一个圆形
            group.addShape('circle', {
              attrs: {
                r: 5,
                fill: quatileColor || '#333',
                x: quatile.x,
                y: quatile.y,
              },
            });
          },
          update: undefined,
        },
        'quadratic',
      )




}




export const failEdge = (edge = {id:'e',source:'',target:'',curveOffset:0}) => {

    return {
        id: edge.id,
        source: edge.source,
        target: edge.target,
        curveOffset: edge.curveOffset,
        type: 'extra-shape-edge',
        midPointColor: '#F4664A',
        quatileColor: '#F4664A',
    }

}


export const successEdge = (edge = {id:'e',source:'',target:'',curveOffset:0}) => {
    
        return {
            id: edge.id,
            source: edge.source,
            target: edge.target,
            curveOffset: edge.curveOffset,
            type: 'extra-shape-edge',
            midPointColor: '#30BF78',
            quatileColor: '#30BF78',
        }
    
}


export const processInfo = (info) => {


    let gap = 400

    const nodes = []
    const edges = []

    let curx = 0;
    let cury = 0;

    let ct = 0;

    info.nodes.forEach((e,index) => {

        if (ct == 2){
            curx += gap
            ct = 0;
        }

        if (ct == 1){
            cury = cury > 0 ? 0 : gap
        }


        nodes.push({
            id:e.node.id,
            error: false,
            info: toRaw(e),
            x:  curx,
            y:  cury,
        })

        ct++

    });

    let rmap = {}

    info.edges.forEach(e => {
        
        let key = e.sourceId.charAt(0) > e.targetId.charAt(0) ? e.targetId + e.sourceId : e.sourceId + e.targetId

        if(rmap[key] === undefined){
            rmap[key] = 0
        }

        let offset = rmap[key] * 50

        rmap[key] += 1

        let edge = {
            id:e.id,
            source:e.sourceId,
            target:e.targetId,
            curveOffset: offset
        }

        edges.push(e.type === 1 ? failEdge(edge) : successEdge(edge))
    });


    console.log('nodes',nodes)
    console.log('edges',edges)



    return {nodes,edges}
}