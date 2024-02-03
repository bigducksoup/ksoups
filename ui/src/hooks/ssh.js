import {
    deleteSSHGroup, deleteSSHInfo,
    getSSHGroupContent,
    getSSHGroupTree,
    saveSSHGroup,
    saveSSHInfo,
    updateSSHGroup, updateSSHInfo
} from "../services/ssh.js";
import { useMessage } from 'naive-ui'
import {onMounted, ref} from "vue";

export const useSSH = () => {

    const message = useMessage()

    const groupInfo = ref([])

    /**
     * 获取SSH分组内的信息
     * @param groupId
     * @return {Promise<void>}
     * @constructor
     */
    const GetSSHGroupContent = async (groupId = 'root' ) => {
        let res = await getSSHGroupContent(groupId)

        if (res['code'] !== 200) {
            message.error(res['msg'])
            groupInfo.value = []
            return
        }
        groupInfo.value = res['data']
    }

    const SaveSSHGroup = async (name = '',parentId = '')=>{
        let res = await saveSSHGroup(name,parentId)

        if (res['code'] !== 200) {
            message.error(res['msg'])
            return false
        }
        message.success(res['msg'])
        return true
    }

    const SaveSSHInfo = async (info)=>{
        let res = await saveSSHInfo(info)

        if (res['code'] !== 200) {
            message.error(res['msg'])
            return false
        }
        message.success(res['msg'])
        return true
    }

    onMounted (() => {
        GetSSHGroupContent('root')
    })


    return { groupInfo,GetSSHGroupContent,SaveSSHGroup,SaveSSHInfo }

}


export const GetSSHGroupTree = async (rootId = 'root') => {
    let res = await getSSHGroupTree(rootId)

    if (res['code'] !== 200) {
        return []
    }
    return res['data']
}



export const useSSHGroupOperation = () => {
    const message = useMessage()


    const SaveSSHGroup = async (name = '',parentId = '')=>{
        let res = await saveSSHGroup(name,parentId)

        if (res['code'] !== 200) {
            message.error(res['msg'])
            return false
        }
        message.success(res['msg'])
        return true
    }

    const DeleteSSHGroup = async (id) => {
        let res = await deleteSSHGroup(id)

        if (res['code'] !== 200) {
            message.error(res['msg'])
            return false
        }

        return true
    }

    /**
     * 更新SSH分组信息
     * @param info
     * @param info.id
     * @param info.name
     * @param info.parent
     * @return {Promise<Boolean>}
     * @constructor
     */
    const UpdateSSHGroup = async (info) => {
        let res = await updateSSHGroup(info)

        if (res['code'] !== 200) {
            message.error(res['msg'])
            return false
        }
        message.success(res['msg'])
        return true
    }




    return {DeleteSSHGroup,SaveSSHGroup,UpdateSSHGroup}


}



export const useSSHInfoOperation = () => {
    const message = useMessage()

    const SaveSSHInfo = async (info)=>{
        let res = await saveSSHInfo(info)

        if (res['code'] !== 200) {
            message.error(res['msg'])
            return false
        }
        message.success(res['msg'])
        return true
    }

    const UpdateSSHInfo = async (info) => {
        let res = await updateSSHInfo(info)

        if (res['code'] !== 200) {
            message.error(res['msg'])
            return false
        }
        message.success(res['msg'])
        return true
    }

    const DeleteSSHInfo = async (infoId) => {
        let res = await deleteSSHInfo(infoId)

        if (res['code'] !== 200) {
            message.error(res['msg'])
            return false
        }
        message.success(res['msg'])
        return true
    }

    return {SaveSSHInfo,UpdateSSHInfo,DeleteSSHInfo}
}