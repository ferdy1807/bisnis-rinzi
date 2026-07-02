<!-- apps/login-sso/src/views/ForbiddenView.vue -->
<template>
  <div class="text-center p-8 max-w-md">
    <i class="pi pi-lock text-6xl text-red-500 mb-4"></i>
    <h1 class="text-6xl font-extrabold text-red-500 mb-2">403</h1>
    <h2 class="text-2xl font-semibold text-gray-900 dark:text-gray-100 mb-3">Akses Ditolak</h2>
    <p class="text-gray-500 dark:text-gray-400 mb-8">
      Maaf, Anda tidak memiliki izin atau otoritas yang memadai untuk mengakses halaman ini.
    </p>
    <Button 
      label="Kembali ke Beranda" 
      icon="pi pi-home" 
      severity="danger" 
      outlined
      class="w-full"
      @click="goHome" 
    />
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import Button from 'primevue/button';

const router = useRouter();

const goHome = () => {
  // Mengarahkan kembali ke gerbang utama
  router.push('/');
};
</script>