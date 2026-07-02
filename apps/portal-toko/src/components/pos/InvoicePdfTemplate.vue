<template>
  <div ref="invoiceRef" class="inv-bg inv-text p-0 w-[800px] font-sans absolute top-[-9999px] left-[-9999px] z-[-999] border-2 inv-border overflow-hidden">
    <!-- Header Block -->
    <div class="inv-header-bg p-10 flex justify-between items-center">
      <div>
        <h1 class="text-4xl font-extrabold tracking-widest text-white mb-2">INVOICE</h1>
        <p class="text-xl font-medium opacity-90 text-white">{{ storeName }}</p>
      </div>
      <div class="text-right text-white">
        <p class="font-bold text-xl mb-1">No. {{ sale.invoice_number }}</p>
        <p class="text-sm opacity-80 mb-3">{{ formatDate(sale.transaction_date) }}</p>
        <div class="inline-block bg-white/20 px-4 py-1.5 rounded-lg border border-white/30 backdrop-blur-sm">
          <p class="text-sm font-semibold">Kasir: {{ authStore.user?.full_name || authStore.user?.username || cashierName }}</p>
        </div>
      </div>
    </div>

    <!-- Body -->
    <div class="p-10 min-h-[400px]">
      <table class="w-full text-left mb-8 border-collapse">
        <thead>
          <tr class="inv-bg-muted inv-border-b">
            <th class="py-4 px-4 font-semibold inv-text-muted rounded-tl-xl">Item</th>
            <th class="py-4 px-4 text-center font-semibold inv-text-muted">Qty</th>
            <th class="py-4 px-4 text-right font-semibold inv-text-muted">Harga Satuan</th>
            <th class="py-4 px-4 text-right font-semibold inv-text-muted">Diskon</th>
            <th class="py-4 px-4 text-right font-semibold inv-text-muted rounded-tr-xl">Subtotal</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in items" :key="item.id" class="border-b inv-border text-sm">
            <td class="py-5 px-4 font-medium text-base">{{ item.product_name }}</td>
            <td class="py-5 px-4 text-center">{{ item.qty }} {{ item.unit_code }}</td>
            <td class="py-5 px-4 text-right">Rp {{ item.unit_price.toLocaleString('id-ID') }}</td>
            <td class="py-5 px-4 text-right inv-error">{{ item.discount > 0 ? '-Rp ' + item.discount.toLocaleString('id-ID') : '-' }}</td>
            <td class="py-5 px-4 text-right font-bold text-base">Rp {{ item.subtotal.toLocaleString('id-ID') }}</td>
          </tr>
        </tbody>
      </table>

      <!-- Summary -->
      <div class="flex justify-end">
        <div class="w-80 space-y-3 text-sm inv-bg-muted p-6 rounded-2xl border inv-border shadow-sm">
          <div class="flex justify-between items-center">
            <span class="inv-text-muted font-medium">Subtotal</span>
            <span class="font-bold text-base">Rp {{ sale.subtotal.toLocaleString('id-ID') }}</span>
          </div>
          <div class="flex justify-between items-center inv-error" v-if="sale.discount > 0">
            <span class="font-medium">Diskon Tambahan</span>
            <span class="font-bold">- Rp {{ sale.discount.toLocaleString('id-ID') }}</span>
          </div>
          <div class="flex justify-between items-center font-bold text-xl border-t inv-border pt-4 mt-2 inv-primary">
            <span>Total Akhir</span>
            <span>Rp {{ sale.total.toLocaleString('id-ID') }}</span>
          </div>
          <div class="flex justify-between items-center pt-3 mt-3 border-t border-dashed inv-border">
            <span class="inv-text-muted font-medium">Metode Pembayaran</span>
            <span class="font-bold inv-success text-base tracking-wide">{{ sale.payment_method }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <div class="inv-bg-muted text-center text-sm inv-text-muted py-8 border-t inv-border">
      <p class="font-bold inv-text mb-1 text-base">Terima kasih telah berbelanja di {{ storeName }}</p>
      <p>Barang yang sudah dibeli tidak dapat ditukar atau dikembalikan.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useAuthStore } from '@frontend/stores/auth';

const authStore = useAuthStore();

const props = defineProps<{
  sale: any;
  items: any[];
  storeName: string;
  cashierName: string;
}>();

const invoiceRef = ref<HTMLElement | null>(null);

const formatDate = (dateString: string) => {
  if (!dateString) return '';
  const d = new Date(dateString);
  return d.toLocaleString('id-ID', {
    year: 'numeric', month: 'long', day: 'numeric',
    hour: '2-digit', minute: '2-digit'
  });
};

defineExpose({ invoiceRef });
</script>

<style scoped>
.inv-bg { background-color: #ffffff !important; }
.inv-header-bg { background-color: #0f172a !important; }
.inv-text { color: #1e293b !important; }
.inv-text-muted { color: #64748b !important; }
.inv-bg-muted { background-color: #f8fafc !important; }
.inv-border { border-color: #e2e8f0 !important; }
.inv-border-b { border-bottom: 2px solid #e2e8f0 !important; }
.inv-primary { color: #2563eb !important; }
.inv-error { color: #dc2626 !important; }
.inv-success { color: #16a34a !important; }
</style>
