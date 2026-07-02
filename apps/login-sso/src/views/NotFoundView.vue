<!-- apps/login-sso/src/views/NotFoundView.vue -->
<template>
  <div class="text-center p-8 max-w-md">
    <i class="pi pi-exclamation-triangle text-6xl text-primary-500 mb-4"></i>
    <h1 class="text-6xl font-extrabold text-primary-500 mb-2">404</h1>
    <h2 class="text-2xl font-semibold text-gray-900 dark:text-gray-100 mb-3">Halaman Tidak Ditemukan</h2>
    <p class="text-gray-500 dark:text-gray-400 mb-8">
      Halaman yang Anda cari tidak tersedia, telah dipindahkan, atau URL yang Anda masukkan salah.
    </p>
    <Button 
      label="Kembali" 
      icon="pi pi-arrow-left" 
      class="w-full"
      @click="goBack" 
    />
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import Button from 'primevue/button';

const router = useRouter();

const goBack = () => {
  // Kembali ke halaman utama untuk mereset navigasi
  router.push('/');
};
</script>