import {ref} from "vue";
import {getShortcutGroup} from "../services/shortcut.js";


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


export {useShortcutGroup}