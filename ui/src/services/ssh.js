import {baseUrl} from "../state/index.js";



// http://127.0.0.1:8080/api/ssh/group/content
export const getSSHGroupContent = async ( groupId = 'root' ) => {

    let sid = window.localStorage.getItem('sid')

    let res = await fetch(baseUrl.value + 'api/ssh/group/content?groupId=' + groupId,{
        method:'GET',
        headers: {
            'sid':sid
        }
    })

    return await res.json()
}

//http://127.0.0.1:8080/api/ssh/group/save
export const saveSSHGroup = async (name = '',parentId = '')=>{
    let sid = window.localStorage.getItem('sid')

    let res = await fetch(baseUrl.value + 'api/ssh/group/save?' + `name=${name}` + '&' + `parentId=${parentId}`,{
        method:'PUT',
        headers: {
            'sid':sid,
        },
    })

    return await res.json()
}


// http://127.0.0.1:8080/api/ssh/info/save


/**
 * @param {object} info 
 * @param {string} info.groupId
 * @param {string} info.addrPort
 * @param {string} info.username
 * @param {string} info.password
 * @returns 
 */
export const saveSSHInfo = async (info)=>{
    let sid = window.localStorage.getItem('sid')

    let res = await fetch(baseUrl.value + 'api/ssh/info/save',{
        method:'PUT',
        headers: {
            'sid':sid,
            'Content-Type':'application/json'
        },
        body:JSON.stringify(info)
    })

    return await res.json()
}


//http://127.0.0.1:8080/api/ssh/group/tree
export const getSSHGroupTree = async (id = 'root') => {
    let sid = window.localStorage.getItem('sid')

    let res = await fetch(baseUrl.value + 'api/ssh/group/tree?' + `rootId=${id}`,{
        method:'GET',
        headers: {
            'sid':sid,
        },
    })

    return await res.json()
}


// http://127.0.0.1:8080/api/ssh/group/delete
export const deleteSSHGroup = async (groupId = null) => {

    if (groupId == null) return
    let sid = window.localStorage.getItem('sid')


    let res = await fetch(baseUrl.value + 'api/ssh/group/delete?' + `groupId=${groupId}`,{
        method:'DELETE',
        headers: {
            'sid':sid,
        },
    })

    return await res.json()

}



// http://127.0.0.1:8080/api/ssh/group/update
/**
 * @param {object} info
 * @param {string} info.id
 * @param {string} info.name
 * @param {string} info.parent
 * @returns
 */
export const updateSSHGroup = async (info) => {
    let sid = window.localStorage.getItem('sid')

    let res = await fetch(baseUrl.value + 'api/ssh/group/update',{
        method:'POST',
        headers: {
            'sid':sid,
            'Content-Type':'application/json'
        },
        body:JSON.stringify(info)
    })

    return await res.json()
}



