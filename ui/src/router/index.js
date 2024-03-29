import {createRouter, createWebHashHistory, createWebHistory} from "vue-router";

const routes = [
    {
        path: "/",
        name:'index',
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
                component: () => import('../views/chain/orchestration.vue')
            },
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
            },
        ],
    },
    {
        path: '/center',
        name: 'center',
        component: () => import('../views/center.vue'),
        children: [
            {
                path: 'keys',
                name: 'keys',
                component: () => import('../views/center/keys.vue')
            },
            {
                path: 'monitor',
                name: 'monitor',
                component: () => import('../views/center/monitor.vue')
            }
        ]
    }
    ,
    {
        path: '/ssh',
        name: 'ssh',
        component: () => import('../views/ssh.vue'),
        children: [
            {
                path: 'host',
                name: 'host',
                component: () => import('../views/ssh/host.vue')
            },
            {
                path: 'term',
                name: 'term',
                component: () => import('../views/ssh/term.vue')
            }
        ]
    },
    {
        path: '/terminal',
        name: 'terminal',
        component: () => import('../views/terminal.vue')
    },
    {
        path: "/login",
        component: () => import("../views/login.vue"),
    }

]


const router = createRouter({
    history: createWebHashHistory(),
    routes
})


export default router