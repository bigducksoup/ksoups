import {createRouter, createWebHistory} from "vue-router";

const routes = [
    {
        path: "/",
        component: () => import("../views/main.vue"),
        meta: {
            transition: 'none'
        },
    },
    {
        path: '/chain',
        component: () => import('../views/chain.vue'),
        meta: {
            transition: 'none'
        },
        children: [
            {
                path: ':chainId',
                name: 'chain',
                component: () => import('../views/chain/detail.vue')
            }
        ]

    },
    {
        path: "/probe",
        component: () => import("../views/probe.vue"),
        meta: {
            transition: 'none'
        },
        children: [
            {
                path: "file/:probeId",
                name: "board",
                component: () => import("../views/probe/file.vue"),
                // component: () => import("../views/console.vue")
            },
            {
                path: "shortcut/:probeId",
                name: "shortcut",
                component: () => import('../views/probe/shortcut.vue')
            }
        ],
    },
    {
        path: "/login",
        component: () => import("../views/login.vue"),
    }

]


const router = createRouter({
    history: createWebHistory(),
    routes
})


export default router