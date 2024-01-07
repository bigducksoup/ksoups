import {baseUrl} from "../state/index.js";


export const getChainList = async () => {
    const response = await fetch(`${baseUrl.value}api/chain/list`);
    return  await response.json();
}