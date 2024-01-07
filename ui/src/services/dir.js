import { baseUrl } from "../state/index.js";

const scanDir = async (probeId, path) => {
    let sid = window.localStorage.getItem('sid');
    const response = await fetch(baseUrl.value + "api/dir/read?path=" + path + "&probeId=" + probeId + "&fileOnly=false", {
        headers: {
            'sid': sid
        }
    });
    const json = await response.json();
    return json;
}

const createDir = async (probeId, path,permission) => {

    //post and payload is json
    let sid = window.localStorage.getItem('sid');

    const response = await fetch(baseUrl.value + "api/dir/create", {
        method: "POST",
        headers: {
            'sid': sid,
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            "probeId": probeId,
            "path": path,
            "permission": permission
        })
    });

    const json = await response.json();
    return json;
}


export{
    scanDir,
    createDir
}