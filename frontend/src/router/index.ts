import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { guest: true },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/RegisterView.vue'),
      meta: { guest: true },
    },
    {
      path: '/',
      name: 'home',
      redirect: '/interviews',
    },
    {
      path: '/interviews',
      name: 'interviews',
      component: () => import('@/views/InterviewListView.vue'),
      meta: { auth: true },
    },
    {
      path: '/interviews/new',
      name: 'new-interview',
      component: () => import('@/views/NewInterviewView.vue'),
      meta: { auth: true },
    },
    {
      path: '/interviews/:id',
      name: 'interview',
      component: () => import('@/views/InterviewView.vue'),
      meta: { auth: true },
    },
    {
      path: '/interviews/:id/report',
      name: 'report',
      component: () => import('@/views/ReportView.vue'),
      meta: { auth: true },
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('@/views/SettingsView.vue'),
      meta: { auth: true },
    },
  ],
})

router.beforeEach((to) => {
  const token = localStorage.getItem('token')
  if (to.meta.auth && !token) {
    return { name: 'login' }
  }
  if (to.meta.guest && token) {
    return { name: 'home' }
  }
})

export default router
