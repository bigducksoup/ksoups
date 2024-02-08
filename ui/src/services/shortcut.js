import {baseUrl} from '../state/index.js'
import {useMessage} from 'naive-ui'


/**
 * 创建快捷方式
 * @param shortcutForm
 * @returns {Promise<void>}
 */
const createShortcut = async (shortcutForm)=>{
    let sid = window.localStorage.getItem('sid')

    let res =  await fetch(baseUrl.value + "api/shortcut/create" , {
        method:'POST',
        headers: {
            'sid': sid
        },
        body:JSON.stringify(shortcutForm)
    })
    return await res.json()
}

const listShortcut = async (probeId)=>{

    let sid = window.localStorage.getItem('sid')
    let res =  await fetch(baseUrl.value + "api/shortcut/list?probeId=" + probeId , {
        headers: {
            'sid': sid
        }
    })
    return await res.json()
}


const runShortcut = async (shortcutId)=>{
    const message = useMessage()
    let sid = window.localStorage.getItem('sid')
    let res =  await fetch(baseUrl.value + "api/shortcut/run?shortcutId=" + shortcutId , {
        method:"POST",
        headers: {
            'sid': sid
        }
    })
    let json = await res.json()
    if (json['code'] !== 200){
        message.error(json['msg'])
        return
    }
    return json['data']
}


/**
 *
 * @param shortcutId
 * @returns {Promise<any>}
 */
const deleteShortcut = async (shortcutId)=>{
    let sid = window.localStorage.getItem('sid')
    let res =  await fetch(baseUrl.value + "api/shortcut/delete?shortcutId=" + shortcutId , {
        method:"DELETE",
        headers: {
            'sid': sid
        }
    })
    let json = await res.json()
    if (json['code'] !== 200){
        return
    }
    return json
}

//http://127.0.0.1:8080/api/shortcut/group
const getShortcutGroup = async ()=>{
    let sid = window.localStorage.getItem('sid')
    let res =  await fetch(baseUrl.value + "api/shortcut/group" , {
        headers: {
            'sid': sid
        }
    })
    return await res.json()
}


const updateShortcut = async (shortcut)=>{
    let sid = window.localStorage.getItem('sid')
    let res =  await fetch(baseUrl.value + "api/shortcut/update" , {
        method:"POST",
        headers: {
            'sid': sid
        },
        body:JSON.stringify(shortcut)
    })
    return  await res.json()
}


//http://127.0.0.1:8080/api/shortcut/group
const getShortcutRunHistory = async (shortcutId)=>{
    let sid = window.localStorage.getItem('sid')
    let res =  await fetch(baseUrl.value + "api/shortcut/run/history" + "?shortcutId=" + shortcutId , {
        headers: {
            'sid': sid
        }
    })
    return await res.json()
}

export {
    listShortcut,
    runShortcut,
    deleteShortcut,
    getShortcutGroup,
    createShortcut,
    updateShortcut,
    getShortcutRunHistory
}