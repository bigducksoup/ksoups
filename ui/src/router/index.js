import { createRouter, createWebHistory } from "vue-router";

const routes = [

    {
        path: "/",
        component: () => import("../pages/index.vue"),    
        children: [
            {
                path: "/file/:probeId",
                name: "board",
                component: () => import("../pages/index/board.vue"),
                // component: () => import("../pages/console.vue")
            },
            {
                path: "/trigger/:probeId",
                name: "trigger",
                component: () => import('../pages/index/trigger.vue')
            }
        ],
    },
    {
        path:"/login",
        component: () => import("../pages/login.vue"),
    }

]


const router = createRouter({
    history: createWebHistory(),
    routes
})



export default router