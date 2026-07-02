<template>
  <div class="mesh-bg w-full min-h-screen flex items-center justify-center p-4 md:p-6 relative overflow-x-hidden">
    
    <!-- Theme Toggle -->
    <div class="fixed top-6 right-6 z-50">
      <button @click="toggleTheme" class="p-3 rounded-full bg-white/20 backdrop-blur-md border border-white/30 text-white hover:scale-110 transition-all duration-300">
        <span class="material-symbols-outlined block dark:hidden">dark_mode</span>
        <span class="material-symbols-outlined hidden dark:block">light_mode</span>
      </button>
    </div>

    <!-- Background Accents -->
    <div class="fixed inset-0 flex z-0">
      <div class="relative w-full h-full overflow-hidden">
        <img alt="Toko Kelontong" class="absolute inset-0 w-full h-full object-cover" src="../pict/background.jpg">
        <div class="absolute inset-0 bg-black/40"></div>
        <div class="absolute bottom-8 left-8 text-white/60 font-bold uppercase tracking-widest text-sm">Toko Kelontong</div>
      </div>
    </div>

    <!-- Main Glass Card -->
    <main class="w-full flex items-center justify-center relative z-10">
      <section class="glass-card rounded-3xl p-8 md:p-10 mx-auto w-full max-w-[480px]">
        <header class="mb-8 text-center md:text-left">
          <h2 class="text-2xl md:text-3xl font-bold text-gray-900 dark:text-white">SSO Login</h2>
          <p class="text-sm text-gray-600 dark:text-gray-300 mt-1">Portal Login multi aplikasi.</p>
        </header>

        <form @submit.prevent="onSubmit" class="space-y-6">
          <!-- Username Input -->
          <div class="space-y-2">
            <label class="block text-sm font-semibold text-gray-900 dark:text-gray-300" for="username">Username</label>
            <div class="relative group">
              <span class="material-symbols-outlined absolute left-4 top-1/2 -translate-y-1/2 text-outline group-focus-within:text-primary transition-colors text-[20px]">mail</span>
              <input v-model="username" type="text" id="username" class="w-full pl-12 pr-4 py-3.5 bg-white/50 dark:bg-bg-dark/40 border border-border-subtle dark:border-outline/20 rounded-2xl focus:border-primary focus:ring-0 dark:text-white transition-all duration-200" :class="{ 'border-red-500': errors.username }" placeholder="username pengguna" />
            </div>
            <p v-if="errors.username" class="text-xs text-red-500 font-semibold">{{ errors.username }}</p>
          </div>

          <!-- Password Input -->
          <div class="space-y-2">
            <div class="relative group">
              <span class="material-symbols-outlined absolute left-4 top-1/2 -translate-y-1/2 text-outline group-focus-within:text-primary transition-colors text-[20px]">lock</span>
              <input v-model="password" :type="passwordFieldType" id="password" class="w-full pl-12 pr-12 py-3.5 bg-white/50 dark:bg-bg-dark/40 border border-border-subtle dark:border-outline/20 rounded-2xl focus:border-primary focus:ring-0 dark:text-white transition-all duration-200" :class="{ 'border-red-500': errors.password }" placeholder="••••••••" />
              <button type="button" @click="togglePasswordVisibility" class="absolute right-4 top-1/2 -translate-y-1/2 p-1 text-outline hover:text-primary transition-colors">
                <span class="material-symbols-outlined text-[20px]">{{ passwordIcon }}</span>
              </button>
            </div>
            <p v-if="errors.password" class="text-xs text-red-500 font-semibold">{{ errors.password }}</p>
          </div>

          <!-- Remember Me -->
          <div class="flex items-center">
            <label class="flex items-center group cursor-pointer">
              <div class="relative">
                <input type="checkbox" class="peer sr-only" />
                <div class="w-5 h-5 border-2 border-gray-400 dark:border-gray-500 rounded-lg peer-checked:bg-primary peer-checked:border-primary transition-all duration-200"></div>
                <span class="material-symbols-outlined absolute inset-0 text-white text-[16px] opacity-0 peer-checked:opacity-100 flex items-center justify-center pointer-events-none">check</span>
              </div>
              <span class="ml-3 text-sm text-gray-600 dark:text-gray-300 group-hover:text-gray-900 dark:group-hover:text-white transition-colors">Ingat Saya di perangkat ini</span>
            </label>
          </div>

          <!-- Alert Error API -->
          <div v-if="apiError" class="p-3 bg-red-100 border border-red-400 text-red-700 rounded-lg text-sm flex items-center gap-2">
            <span class="material-symbols-outlined">error</span>
            {{ apiError }}
          </div>

          <!-- Submit Button -->
          <button type="submit" :disabled="authStore.loading" class="w-full bg-primary text-white font-semibold text-sm py-4 rounded-2xl shadow-xl shadow-primary/20 hover:shadow-primary/30 hover:bg-primary-container transition-all duration-300 active:scale-[0.98] flex items-center justify-center gap-2 disabled:opacity-70 disabled:cursor-not-allowed">
            <template v-if="authStore.loading">
              <span class="animate-spin material-symbols-outlined">sync</span> Memproses...
            </template>
            <template v-else>
              Masuk ke Dashboard
              <span class="material-symbols-outlined text-[20px]">login</span>
            </template>
          </button>
        </form>

        <footer class="mt-8 pt-6 border-t border-border-subtle dark:border-outline/10 text-center">
          <div class="mt-4">
            <a href="https://wa.me/6282110735089" target="_blank" rel="noopener noreferrer" class="text-sm text-outline hover:text-gray-900 dark:hover:text-white transition-colors">Hubungi Bantuan Admin</a>
          </div>
        </footer>
      </section>
    </main>

    <!-- Custom Status Toast Animasi -->
    <div :class="['fixed bottom-8 left-1/2 -translate-x-1/2 px-6 py-4 rounded-2xl shadow-2xl flex items-center gap-4 transition-all duration-500 z-50 bg-white dark:bg-surface-dark border border-border-subtle dark:border-outline/20', showToast ? 'translate-y-0 opacity-100' : 'translate-y-20 opacity-0 pointer-events-none']">
      <div class="w-10 h-10 rounded-full flex items-center justify-center bg-success text-white shadow-lg shadow-success/20">
        <span class="material-symbols-outlined text-[22px]">verified_user</span>
      </div>
      <div>
        <p class="font-semibold text-sm text-gray-900 dark:text-white">Autentikasi Berhasil</p>
        <p class="text-xs text-gray-600 dark:text-gray-300">Mengarahkan Anda ke Dashboard...</p>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useForm, useField } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import { LoginSchema } from '@frontend/validators/auth';
import { useAuthStore } from '@frontend/stores/auth';

const authStore = useAuthStore();
const router = useRouter();
const apiError = ref<string | null>(null);
const showToast = ref(false);

const isPasswordVisible = ref(false);
const passwordFieldType = computed(() => isPasswordVisible.value ? 'text' : 'password');
const passwordIcon = computed(() => isPasswordVisible.value ? 'visibility_off' : 'visibility');
const togglePasswordVisibility = () => { isPasswordVisible.value = !isPasswordVisible.value; };

const toggleTheme = () => {
  const html = document.documentElement;
  if (html.classList.contains('dark')) {
    html.classList.remove('dark');
    localStorage.setItem('theme', 'light');
  } else {
    html.classList.add('dark');
    localStorage.setItem('theme', 'dark');
  }
};

onMounted(() => {
  if (localStorage.getItem('theme') === 'dark' || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    document.documentElement.classList.add('dark');
  }
});

const { handleSubmit, errors } = useForm({
  validationSchema: toTypedSchema(LoginSchema),
});
const { value: username } = useField<string>('username');
const { value: password } = useField<string>('password');

const onSubmit = handleSubmit(async (values) => {
  apiError.value = null;
  try {
    await authStore.login({ username: values.username, password: values.password });
    showToast.value = true;
    setTimeout(() => {
      const urlParams = new URLSearchParams(window.location.search);
      const redirectUrl = urlParams.get('redirect');
      if (redirectUrl) {
        const delimiter = redirectUrl.includes('?') ? '&' : '?';
        window.location.href = `${redirectUrl}${delimiter}accessToken=${authStore.token}&refreshToken=${authStore.refreshToken}`;
      } else {
        router.push({ name: 'Login', query: { redirect: Date.now() } });
      }
    }, 1500);
  } catch (error: any) {
    apiError.value = error.response?.data?.message || 'Gagal terhubung ke server. Periksa username dan password Anda.';
  }
});
</script>