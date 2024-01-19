import {
    allStepDispatch,
    bindShortcut,
    createNode,
    deleteNode,
    getChainDetail, getChainExecHistory, getChainExecResult,
    linkNode, newDispatch, setRoot, singleStepDispatch,
    unbindShortcut,
    unlinkNode
} from '../services/chain.js'
import {useMessage} from "naive-ui";
import {ref} from "vue";




const useChainInfo = (chainId) => {
    const message = useMessage()

    const chainDetail = ref({})

    const nodes = ref([])

    const getChainData = async () => {
        let response = await getChainDetail(chainId)

        if (response.code !== 200) {
            message.error(response.msg)
            return
        }

        chainDetail.value = response.data

        let cnodes = response.data.nodes

        if (cnodes !== null) {
            nodes.value = []
            cnodes.forEach((item) => {
                nodes.value.push(item.node)
            })
        }
    }

    getChainData()

    return {chainDetail, nodes, getChainData}
}


const useNodeOperation = () => {
    const message = useMessage()

    const DelNode = async (nodeId) => {
        console.log('delNode', nodeId)

        let res = await deleteNode(nodeId)

        if (res.code !== 200) {
            message.error(res.msg)
            return false
        }

        message.success('删除成功')
        return true
    }

    const NewNode = async (chainId, name, description) => {
        console.log('newNode', chainId, name, description)
        let res = await createNode(chainId, name, description)

        if (res.code !== 200) {
            message.error(res.msg)
            return false
        }
        message.success('创建成功')
        return true
    }

    const stRoot = async (nodeId) => {
        console.log('stRoot', nodeId)

        let res = await setRoot(nodeId)

        if (res.code !== 200) {
            message.error(res.msg)
            return false
        }

        message.success('设置成功')
        return true
    }

    const LkNode = async (sourceId, targetId, chainId, type) => {

        let res = await linkNode(sourceId, targetId, chainId, type)

        if (res.code !== 200) {
            message.error(res.msg)
            return false
        }

        message.success('链接成功')
        return true
    }

    const UnlkNode = async (sourceId, targetId, chainId, type) => {


        let res = await unlinkNode(sourceId, targetId, chainId, type)

        if (res.code !== 200) {
            message.error(res.msg)
            return false
        }

        message.success('取消链接成功')
        return true

    }

    const BindShortcut = async (nodeId, shortcut) => {
        console.log('BindShortcut', nodeId, shortcut)

        let res = await bindShortcut(nodeId, shortcut)

        if (res.code !== 200) {
            message.error(res.msg)
            return false
        }
        message.success('绑定成功')
        return true

    }

    const UnbindShortcut = async (nodeId, shortcut) => {
        console.log('UnbindShortcut', nodeId, shortcut)

        let res = await unbindShortcut(nodeId, shortcut)

        if (res.code !== 200) {
            message.error(res.msg)
            return false
        }
        message.success('解绑成功')
        return true

    }

    return { DelNode, NewNode, LkNode, UnlkNode, BindShortcut, UnbindShortcut, stRoot }

}


const useChainExecInfo = (chainId) => {

    const message = useMessage()

    const execHistory = ref([])
    const dispatchInfo = ref({
        "dispatchLog": {
            "id": null,
            "chainId": null,
            "createTime": "2024-01-18T11:32:36.57942+08:00",
            "status": 0,
            "done": false
        }
    })

    const GetExecHistory = async (ChainId) => {

        let res = await getChainExecHistory(ChainId)

        if (res.code !== 200) {
            message.error(res.msg)
            return
        }
        execHistory.value = res.data
    }

    const GetDispatchInfo = async (Id) => {

        let res = await getChainExecResult(Id)

        if (res.code !== 200) {
            message.error(res.msg)
            return
        }

        dispatchInfo.value = res.data
    }

    const NewDispatch = async () => {

        let res = await newDispatch(chainId)

        if (res.code !== 200) {
            message.error(res.msg)
            return
        }

        message.success('创建成功')
        await GetExecHistory(chainId).then(()=>{
            if(execHistory.value.length > 0)
                GetDispatchInfo(execHistory.value[0].id)
        })
    }


    const SingleStepDispatch = async () => {

        let res = await singleStepDispatch(dispatchInfo.value.dispatchLog.id)

        if (res.code !== 200) {
            message.error(res.msg)
            return
        }

        message.success('执行成功')
        await GetDispatchInfo(dispatchInfo.value.dispatchLog.id)
    }


    const AllStepDispatch = async () => {

        let res = await allStepDispatch(dispatchInfo.value.dispatchLog.id)

        if (res.code !== 200) {
            message.error(res.msg)
            return
        }

        message.success('执行成功')
        await GetDispatchInfo(dispatchInfo.value.dispatchLog.id)
    }
    



    GetExecHistory(chainId).then(()=>{
        if(execHistory.value.length > 0)
        GetDispatchInfo(execHistory.value[0].id)
    })

    return {execHistory, dispatchInfo,GetDispatchInfo,GetExecHistory,NewDispatch,SingleStepDispatch,AllStepDispatch}
}




export {useChainInfo, useNodeOperation, useChainExecInfo}