import {baseUrl} from '../state/index.js'
import {useMessage} from 'naive-ui'


const message = useMessage()

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



export {
    listShortcut,
    runShortcut,
    deleteShortcut
}