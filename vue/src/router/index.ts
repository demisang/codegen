import {createRouter, createWebHashHistory, RouteRecordRaw} from 'vue-router'
import TemplatesList from "@/views/TemplatesList.vue";
import AboutView from "@/views/AboutView.vue";

const routes: Array<RouteRecordRaw> = [
    {
        path: '/',
        name: 'templates',
        component: TemplatesList,
    },
    {
        path: '/about',
        name: 'about',
        component: AboutView,
    },
]

const router = createRouter({
    history: createWebHashHistory(process.env.BASE_URL),
    routes
})

export default router
