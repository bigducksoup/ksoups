import {baseUrl} from "../state/index.js";


export const getChainList = async () => {
    let sid = window.localStorage.getItem('sid');
    const response = await fetch(`${baseUrl.value}api/chain/list`,{
        method: 'GET',
        headers: {
            'sid': sid
        }
    });
    return  await response.json();
}


//http://127.0.0.1:8080/api/chain/info
export const getChainDetail = async (chainId) => {
    let sid = window.localStorage.getItem('sid');
    const response = await fetch(`${baseUrl.value}api/chain/info?chainId=${chainId}`,
        {
            method: 'GET',
            headers: {
                'sid': sid
            }
        }
    );
    return  await response.json();
}


//http://127.0.0.1:8080/api/chain/node/create
export const createNode = async (chainId, name, description) => {

    let sid = window.localStorage.getItem('sid');

    const response = await fetch(`${baseUrl.value}api/chain/node/create`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'sid': sid
        },
        body: JSON.stringify({
            'chainId':chainId,
            'name':name,
            'description':description
        })
    });
    return  await response.json();
}


//http://127.0.0.1:8080/api/chain/node/bind/shortcut
export const bindShortcut = async ( nodeId, shortcutId) => {

    let sid = window.localStorage.getItem('sid');

    const response = await fetch(`${baseUrl.value}api/chain/node/bind/shortcut`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'sid': sid
        },
        body: JSON.stringify({
            'nodeId':nodeId,
            'shortcutId':shortcutId
        })
    });
    return  await response.json();
}

//http://127.0.0.1:8080/api/chain/node/unbind/shortcut
export const unbindShortcut = async ( nodeId, shortcutId) => {

    let sid = window.localStorage.getItem('sid');

    const response = await fetch(`${baseUrl.value}api/chain/node/unbind/shortcut`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'sid': sid
        },
        body: JSON.stringify({
            'nodeId':nodeId,
            'shortcutId':shortcutId
        })
    });
    return  await response.json();
}


//http://127.0.0.1:8080/api/chain/node/link
export const linkNode = async (sourceId, targetId, chainId, type) => {

    let sid = window.localStorage.getItem('sid');

    const response = await fetch(`${baseUrl.value}api/chain/node/link`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'sid': sid
        },
        body: JSON.stringify({
            'sourceId':sourceId,
            'targetId':targetId,
            'chainId':chainId,
            'type':type
        })
    });
    return  await response.json();
}


//http://127.0.0.1:8080/api/chain/node/unlink
export const unlinkNode = async (sourceId, targetId, chainId, type) => {

    let sid = window.localStorage.getItem('sid');

    const response = await fetch(`${baseUrl.value}api/chain/node/unlink`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'sid': sid
        },
        body: JSON.stringify({
            'sourceId':sourceId,
            'targetId':targetId,
            'chainId':chainId,
            'type':type
        })
    });
    return  await response.json();
}


//http://127.0.0.1:8080/api/chain/create
export const createChain = async (name, description) => {

    let sid = window.localStorage.getItem('sid');

    const response = await fetch(`${baseUrl.value}api/chain/create`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'sid': sid
        },
        body: JSON.stringify({
            'name':name,
            'description':description
        })
    });
    return  await response.json();
}



export const deleteNode = async (nodeId) => {

    let sid = window.localStorage.getItem('sid');

    const response = await fetch(`${baseUrl.value}api/chain/node/delete?nodeId=${nodeId}`,{
        method: 'DELETE',
        headers: {
            'sid': sid
        }
    });
    return  await response.json();
}

//http://127.0.0.1:8080/api/chain/node/set/root
export const setRoot = async (nodeId) => {

    let sid = window.localStorage.getItem('sid');

    const response = await fetch(`${baseUrl.value}api/chain/node/set/root?nodeId=${nodeId}`,{
        method: 'PUT',
        headers: {
            'sid': sid
        }
    });
    return  await response.json();
}

//http://127.0.0.1:8080/api/chain/exec/dispatch/new
export const newDispatch = async (chainId) => {
    let sid = window.localStorage.getItem('sid');

    const response = await fetch(`${baseUrl.value}api/chain/exec/dispatch/new?chainId=${chainId}`,{
        method: 'PUT',
        headers: {
            'sid': sid
        }
    });
    return  await response.json();
}

//http://127.0.0.1:8080/api/chain/exec/single/step/dispatch
export const singleStepDispatch = async (dispatchId) => {
    let sid = window.localStorage.getItem('sid');

    const response = await fetch(`${baseUrl.value}api/chain/exec/single/step/dispatch?dispatchId=${dispatchId}`,{
        method: 'PUT',
        headers: {
            'sid': sid
        }
    });
    return  await response.json();
}


//http://127.0.0.1:8080/api/chain/exec/all/step/dispatch
export const allStepDispatch = async (dispatchId) => {
    let sid = window.localStorage.getItem('sid');

    const response = await fetch(`${baseUrl.value}api/chain/exec/all/step/dispatch?dispatchId=${dispatchId}`,{
        method: 'POST',
        headers: {
            'sid': sid
        }
    });
    return  await response.json();
}






//http://127.0.0.1:8080/api/chain/exec/history
export const getChainExecHistory = async (chainId) => {

        let sid = window.localStorage.getItem('sid');

        const response = await fetch(`${baseUrl.value}api/chain/exec/history?chainId=${chainId}`,{
            method: 'GET',
            headers: {
                'sid': sid
            }
        });
        return  await response.json();
}


//http://127.0.0.1:8080/api/chain/exec/log
export const getChainExecResult = async (dispatchId) => {

        let sid = window.localStorage.getItem('sid');

        const response = await fetch(`${baseUrl.value}api/chain/exec/result?dispatchId=${dispatchId}`,{
            method: 'GET',
            headers: {
                'sid': sid
            }
        });
        return  await response.json();
}

