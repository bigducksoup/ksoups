import {reactive, ref} from "vue"
import {msgPush} from "../hooks/websocket.js";


const baseUrl = ref(import.meta.env.VITE_APP_BASE_URL);

const baseHost = ref(import.meta.env.VITE_APP_BASE_HOST)


const realtimeShortcuts = reactive({
    // 'id': {
    //     stdout: "",
    //     stdErr: ""
    // }
})


const hookRealtimeRunHandler = () => {

    console.log('hook start')

    msgPush.AddHandler((message) => {



        if (message['messageType'] === '100001') {

            let content = JSON.parse(message['json'])

            let runId = content['runId']

            if (!realtimeShortcuts[runId]){
                realtimeShortcuts[runId] = {
                    stdout: "",
                    stdErr: ""
                }
            }


            if (content['type'] === 0){
                realtimeShortcuts[runId].stdErr += content['payload']
            }

            if(content['type'] === 1){
                realtimeShortcuts[runId].stdout += content['payload']
            }

            console.log(realtimeShortcuts)

        }
    })

    console.log('hook done')

}


export {
    baseUrl,
    baseHost,
    realtimeShortcuts,
    hookRealtimeRunHandler
}
