import {ref} from "vue";
import {
    createShortcut,
    deleteShortcut,
    getShortcutGroup,
    getShortcutRunHistory,
    updateShortcut
} from "../services/shortcut.js";
import {useMessage} from "naive-ui";


const useShortcutGroup = () => {

    const shortcutGroup = ref({})

    const fetchShortcutGroup = async () => {
        let response = await getShortcutGroup()

        if (response.code !== 200) {
            message.error(response.msg)
            return
        }

        shortcutGroup.value = response.data
    }

    fetchShortcutGroup()

    return {shortcutGroup, fetchShortcutGroup}

}


const useShortcutOperation = () => {

    const message = useMessage()


    const CreateShortcut = async (shortcutForm) => {
        let res = await createShortcut(shortcutForm)
        console.log(res)
        if (res['code'] !== 200) {
            message.error(res['msg'])
            return false
        }
        message.success(res['msg'])
        return true
    }


    /**
     * 删除快捷方式
     * @param id
     * @returns {Promise<boolean>}
     * @constructor
     */
    const DeleteShortcutById = async (id) => {
        let res = await deleteShortcut(id)
        if (res['code'] !== 200) {
            message.error(res['msg'])
            return false
        }
        message.success(res['msg'])
        return true
    }


    /**
     * 更新快捷方式
     * @returns {Promise<Boolean>}
     * @constructor
     * @param {Object} shortcut
     * @param {String} shortcut.id
     * @param {String} shortcut.name
     * @param {String} shortcut.description
     * @param {String} shortcut.payload
     * @param {String} shortcut.probeId
     * @param {Number} shortcut.timeout
     * @param {String} shortcut.createTime
     * @param {Boolean} shortcut.justRun
     * @param {Number} shortcut.type
     */
    const UpdateShortcut = async (shortcut) => {

        let res = await updateShortcut(shortcut)
        if (res['code'] !== 200) {
            message.error(res['msg'])
            return false
        }
        message.success(res['msg'])
        return true
    }


    return {DeleteShortcutById, UpdateShortcut,CreateShortcut,updateShortcut}

}




const useShortcutInfos = () => {

    const message = useMessage()

    /**
     * 获取指令执行历史
     * @param {String} id
     * @return {Promise<Object|null>}
     * @constructor
     */
    const GetShortcutRunHistory = async (id) => {
        let res = await getShortcutRunHistory(id)
        if (res['code'] !== 200 ){
            message.error(res['msg'])
            return []
        }
        return res['data']
    }

    return {GetShortcutRunHistory}

}




export {useShortcutGroup, useShortcutOperation,useShortcutInfos}