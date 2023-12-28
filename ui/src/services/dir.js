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



export{
    scanDir
}