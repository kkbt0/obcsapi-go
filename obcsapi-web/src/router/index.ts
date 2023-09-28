import { createRouter, createWebHistory } from 'vue-router';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/HomeView.vue')
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue')
    },
    {
      path: '/setting',
      name: 'setting',
      component: () => import('../views/SettingView.vue')
    },
    {
      path: '/talk',
      name: 'talk',
      component: () => import('../views/TalkView.vue')
    },
    {
      path: '/edit',
      name: 'edit',
      component: () => import('../views/EditView.vue')
    },
    {
      path: '/search',
      name: 'search',
      component: () => import('../views/SearchView.vue')
    },
    {
      path: '/form',
      name: 'form',
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/FormView.vue')
    }
  ]
})

export default router
