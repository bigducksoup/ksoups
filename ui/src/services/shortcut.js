import {baseUrl} from '../state/index.js'
import {useMessage} from 'naive-ui'




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


const deleteShortcut = async (shortcutId)=>{
    const message = useMessage()
    let sid = window.localStorage.getItem('sid')
    let res =  await fetch(baseUrl.value + "api/shortcut/delete?shortcutId=" + shortcutId , {
        method:"DELETE",
        headers: {
            'sid': sid
        }
    })
    let json = await res.json()
    if (json['code'] !== 200){
        message.error(json['msg'])
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

    

    let json = await res.json()

    return json
}

export {
    listShortcut,
    runShortcut,
    deleteShortcut,
    getShortcutGroup
}