import { createRouter, createWebHistory } from "vue-router";

const routes = [

    {
        path: "/",
        component: () => import("../pages/index.vue"),    
        children: [
            {
                path: ":addr",
                name: "console",
                component: () => import("../pages/console.vue"),
            }
        ]  
    }

]


const router = createRouter({
    history: createWebHistory(),
    routes
})



export default router