// apps/login-sso/src/router/index.ts
import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import BlankLayout from '@frontend/layouts/BlankLayout.vue';
import LoginView from '../views/LoginView.vue';

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    component: BlankLayout,
    children: [
      {
        path: '',
        name: 'Login',
        component: LoginView,
      },
      {
        path: '403',
        name: 'Forbidden',
        component: () => import('../views/ForbiddenView.vue'),
      },
      {
        path: ':pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('../views/NotFoundView.vue'),
      }
    ]
  }
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
});