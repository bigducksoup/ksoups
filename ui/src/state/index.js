import { ref } from "vue"


const baseUrl = ref(import.meta.env.VITE_APP_BASE_URL);

const baseHost = ref(import.meta.env.VITE_APP_BASE_HOST)

export {
    baseUrl,
    baseHost
}
