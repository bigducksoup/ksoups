import {baseUrl} from '../state/index.js'

const getFileContent = async (probeId, path) => {

    let sid = window.localStorage.getItem('sid');


    let res =  await fetch(baseUrl.value + "api/file/read?path=" + path + "&probeId=" + probeId, {
        headers: {
            'sid': sid
        }
    })

    return await res.json()
}


const reqModify = async (params) => {
    let sid =  window.localStorage.getItem('sid')

    let res = await fetch(baseUrl.value + "api/file/modify",{
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "sid": sid
        },
        body: JSON.stringify(params)
    })

    return await res.json()

}


const modifyFile = async (ProbeId,path,diffRes) => {
    let changes = []


    for (let diff of diffRes) {

        let operation = 2

        let valueArr = []

        if (diff.added === true) {
            operation = 1
            valueArr = diff.value.trimEnd().split("\n")
            console.log(valueArr)
        }

        if (diff.removed === true) {
            operation = 0
        }

        let change = {
            "count": diff.count,
            "operation": operation,
            "value": valueArr
        }
        changes.push(change)
    }


    let modifyParams = {
        "path": path,
        "probeId": ProbeId,
        "changes": changes
    }

    console.log(modifyParams)

    return await reqModify(modifyParams)

}

/**
 * @param {String} probeId
 * @param {String} path
 * @param {String} permission
 * @return {Promise<Object>}
 */
const createFile = async (probeId,path,permission) => {

    let sid =  window.localStorage.getItem('sid')

    let res = await fetch(baseUrl.value + "api/file/create",{
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "sid": sid
        },
        body: JSON.stringify({
            "path": path,
            "probeId": probeId,
            "permission": permission
        })
    })

    return await res.json()

}

export{
    getFileContent,
    modifyFile,
    createFile
}