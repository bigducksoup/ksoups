import {createRouter, createWebHistory} from "vue-router";

const routes = [
    {
        path: "/",
        component: () => import("../pages/main.vue"),
        meta: {
            transition: 'none'
        },
    },
    {
        path: '/chain',
        component: () => import('../pages/chain.vue'),
        meta: {
            transition: 'none'
        }
    },
    {
        path: "/node",
        component: () => import("../pages/node.vue"),
        meta: {
            transition: 'none'
        },
        children: [
            {
                path: "file/:probeId",
                name: "board",
                component: () => import("../pages/node/file.vue"),
                // component: () => import("../pages/console.vue")
            },
            {
                path: "shortcut/:probeId",
                name: "shortcut",
                component: () => import('../pages/node/shortcut.vue')
            }
        ],
    },
    {
        path: "/login",
        component: () => import("../pages/login.vue"),
    }

]


const router = createRouter({
    history: createWebHistory(),
    routes
})


export default router