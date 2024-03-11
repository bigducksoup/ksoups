import {baseUrl} from '../state/index.js'


export const keyList = async () => {

    const res = await fetch(`${baseUrl.value}api/info/keypair/list`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    const data = await res.json()
    return data
}


export const deleteKeyPair = async (id) => {
    const res = await fetch(`${baseUrl.value}api/info/keypair/delete?id=${id}`, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    const data = await res.json()
    return data
}

export const generateKeyPair = async (name) => {
    const res = await fetch(`${baseUrl.value}api/info/keypair/generate?name=${name}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    const data = await res.json()
    return data
}