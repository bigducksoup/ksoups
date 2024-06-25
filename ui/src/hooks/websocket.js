import {baseHost} from "../state/index.js";
import {useMessage} from 'naive-ui'

const useMessagePush = () => {

    const message = useMessage()

    let sid = window.localStorage.getItem('sid');

    let handlers = [];

    let ws = new WebSocket(`ws://${baseHost.value}/ws/msg?sid=${sid}`)


    ws.onopen = () =>{
        message.success('open message receiver')
    }

    ws.onclose = () =>{

    }

    ws.onmessage = (msg) => {
        let message =  JSON.parse(msg.data)
        handlers.forEach(fn =>{
            fn(JSON.parse(message))
        })
    }


    ws.onerror = (err) => {
        console.log(err)
    }

    const AddHandler = (handler = (msg)=>{console.log(msg)}) => {
        handlers.push(handler)
    }

    const RMHandler = (handler) => {
        handlers = handlers.filter(item => item !== handler)
    }

    return {AddHandler,RMHandler}
}

let msgPush;

const initMessagePush = () => {
    msgPush = useMessagePush()
}



export {
    msgPush,
    initMessagePush,
    useMessagePush
}