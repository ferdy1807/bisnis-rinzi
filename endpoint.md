# API Endpoints Documentation

Dokumen ini berisi daftar lengkap *endpoint* REST API untuk semua *microservices* dalam ekosistem Bisnis Rinzi.

---

## 🔐 1. Auth Service

### Autentikasi & Sesi
| Method | Endpoint | Keterangan 
| `POST` | `/api/auth/register` | Mendaftarkan pengguna baru   
| `POST` | `/api/auth/login` | Masuk dan mendapatkan *token*   
| `POST` | `/api/auth/refresh-token` | Memperbarui *access token*   
| `POST` | `/api/auth/logout` | Mengakhiri sesi pengguna   

### Manajemen Pengguna (Users)
| Method | Endpoint | Keterangan  
| `GET` | `/api/auth/me` | Melihat profil pengguna saat ini   
| `GET` | `/api/auth/users` | Daftar semua pengguna   
| `POST` | `/api/auth/users` | Membuat pengguna baru (Admin)   
| `GET` | `/api/auth/users/{id}` | Detail pengguna   
| `PUT` | `/api/auth/users/{id}` | Ubah profil pengguna   
| `PUT` | `/api/auth/users/{id}/password` | Ubah kata sandi pengguna   
| `DELETE` | `/api/auth/users/{id}` | Hapus pengguna   
| `DELETE` | `/api/auth/users/{id}/sessions` | *Force logout* semua *device* pengguna terkait (menghapus dari tabel `refresh_tokens`)   

### Peran & Hak Akses (Roles)
| Method | Endpoint | Keterangan  
| `GET` | `/api/auth/roles` | Daftar peran sistem   
| `PUT` | `/api/auth/roles/{code}/dashboard-url` | Atur URL dasbor berdasarkan peran   

### Audit & Token Management (Admin)
| Method | Endpoint | Keterangan  
| `GET` | `/api/auth/audit-logs` | Lihat seluruh riwayat log audit   
| `GET` | `/api/auth/audit-logs/{id}` | Detail log audit   
| `GET` | `/api/auth/tokens` | Melihat seluruh *active token* lintas *user*   
| `DELETE` | `/api/auth/tokens/{id}` | Mencabut sesi/token spesifik   

---

## 📦 2. Inventory Service

### Kategori (Category)
| Method | Endpoint  
| `GET` | `/api/inventory/categories`   
| `POST` | `/api/inventory/categories`   
| `GET` | `/api/inventory/categories/{id}`   
| `PUT` | `/api/inventory/categories/{id}`   
| `DELETE` | `/api/inventory/categories/{id}`   

### Merek (Brand)
| Method | Endpoint  
| `GET` | `/api/inventory/brands`   
| `POST` | `/api/inventory/brands`   
| `GET` | `/api/inventory/brands/{id}`   
| `PUT` | `/api/inventory/brands/{id}`   
| `DELETE` | `/api/inventory/brands/{id}`   

### Satuan (Unit)
| Method | Endpoint  
| `GET` | `/api/inventory/units`   
| `POST` | `/api/inventory/units`   
| `GET` | `/api/inventory/units/{id}`   
| `PUT` | `/api/inventory/units/{id}`   
| `DELETE` | `/api/inventory/units/{id}`   

### Produk Utama (Products)
| Method | Endpoint | Keterangan  
| `GET` | `/api/inventory/products` 
| `POST` | `/api/inventory/products` 
| `GET` | `/api/inventory/products/{id}` 
| `PUT` | `/api/inventory/products/{id}` 
| `DELETE` | `/api/inventory/products/{id}` 
| `GET` | `/api/inventory/products/{id}/stock` | Cek kuantitas stok produk   
| `GET` | `/api/inventory/products/{id}/stock-history` | Riwayat keluar/masuk stok   
| `GET` | `/api/inventory/products/{id}/cost-histories` | Riwayat harga modal produk   
| `POST` | `/api/inventory/products/{id}/cost-histories` | Atur/ubah harga modal produk   
| `GET` | `/api/inventory/products/search` | Pencarian global produk   
| `GET` | `/api/inventory/products/low-stock` | Notifikasi produk menipis   
| `GET` | `/api/inventory/products/barcode/{code}` | Pindai *barcode* produk   
| `POST` | `/api/inventory/products/import` | Impor data CSV/Excel   
| `GET` | `/api/inventory/products/export` | Ekspor data produk   
| `GET` | `/api/inventory/products/sku/{sku}` | Cari berdasarkan SKU   

### Media Produk (Product Media)
| Method | Endpoint  
| `GET` | `/api/inventory/products/{id}/media`   
| `POST` | `/api/inventory/products/{id}/media`   
| `GET` | `/api/inventory/products/{id}/media/{media_id}`   
| `DELETE` | `/api/inventory/products/{id}/media/{media_id}`   

### Manajemen Stok (Stock)
| Method | Endpoint | Keterangan  
| `GET` | `/api/inventory/stocks` | Daftar stok seluruh produk   
| `GET` | `/api/inventory/stocks/card/{product_id}` | Kartu stok produk spesifik   
| `GET` | `/api/inventory/stock-movements` | Riwayat mutasi stok   
| `GET` | `/api/inventory/stock-movements/{id}` | Detail mutasi stok   
| `POST` | `/api/inventory/stocks/adjust` | Penyesuaian/Stock Opname   
| `PUT` | `/api/inventory/stocks/{product_id}/thresholds` | Atur batas stok minimum & safety stock 

### Offline Sync (PWA)
| Method | Endpoint  
| `GET` | `/api/inventory/sync/version`   
| `GET` | `/api/inventory/sync/changes`   
| `GET` | `/api/inventory/catalog/sync`   

---

## 💵 3. Cash Service

### Sesi Kasir (Shift)
| Method | Endpoint | Keterangan  
| `POST` | `/api/cash/shifts/open` | Buka *shift* kasir   
| `POST` | `/api/cash/shifts/close` | Tutup *shift* kasir   
| `GET` | `/api/cash/shifts/current` | Cek *shift* berjalan   
| `GET` | `/api/cash/shifts`   
| `GET` | `/api/cash/shifts/{id}`   
| `GET` | `/api/cash/shifts/{id}/summary` | Rekapitulasi transaksi akhir kasir   

### Pengeluaran (Expense) & Kategori
| Method | Endpoint | Keterangan  
| `GET` | `/api/cash/expense-categories`   
| `POST` | `/api/cash/expense-categories`   
| `GET` | `/api/cash/expense-categories/{id}`   
| `PUT` | `/api/cash/expense-categories/{id}`   
| `DELETE` | `/api/cash/expense-categories/{id}`   
| `GET` | `/api/cash/expenses`   
| `POST` | `/api/cash/expenses`   
| `GET` | `/api/cash/expenses/{id}`   
| `PUT` | `/api/cash/expenses/{id}`   
| `DELETE` | `/api/cash/expenses/{id}`   

### Transaksi Tunai (Cash Transactions)
| Method | Endpoint  
| `GET` | `/api/cash/transactions`   
| `GET` | `/api/cash/transactions/{id}`   
| `POST` | `/api/cash/transactions`   
---

## 🛒 4. POS Service
### Penjualan (Sales)
| Method | Endpoint | Keterangan  
| `GET` | `/api/pos/sales` | Riwayat transaksi   
| `POST` | `/api/pos/sales` | Pemrosesan keranjang (*Checkout*). **Wajib kirim `product_name` di payload items**   
| `GET` | `/api/pos/sales/{id}` | Detail referensi penjualan   
| `GET` | `/api/pos/sales/{id}/receipt` | Cetak struk termal (termasuk *snapshot* nama barang)   
| `GET` | `/api/pos/sales/invoice/{invoice_number}`   
| `GET` | `/api/pos/sales/{id}/items`   
### Offline Sync
| Method | Endpoint | Keterangan  
| `POST` | `/api/pos/sync` | Sinkronisasi data ke *server*   
| `POST` | `/api/pos/sync/retry` | Percobaan ulang transaksi gagal   
| `GET` | `/api/pos/sync/status` | Status *sync* PWA   
| `GET` | `/api/pos/sync/failed` | Lihat antrean yang gagal   

---

## 📅 5. Rental Service

### Katalog Produk (Kategori & Unit)
| Method | Endpoint  
| `GET` | `/api/rental/categories` 
| `POST` | `/api/rental/categories` 
| `GET` | `/api/rental/categories/{id}` 
| `PUT` | `/api/rental/categories/{id}` 
| `DELETE` | `/api/rental/categories/{id}` 
| `GET` | `/api/rental/products` 
| `POST` | `/api/rental/products` 
| `GET` | `/api/rental/products/{id}` 
| `PUT` | `/api/rental/products/{id}` 
| `DELETE` | `/api/rental/products/{id}` 
| `GET` | `/api/rental/products/{id}/media` 
| `POST` | `/api/rental/products/{id}/media` 
| `GET` | `/api/rental/products/{id}/media/{media_id}`   
| `DELETE` | `/api/rental/products/{id}/media/{media_id}` 

### Ketersediaan & Kalender
| Method | Endpoint | Keterangan  
| `GET` | `/api/rental/availability` | Cek ketersediaan produk pada suatu tanggal 
| `GET` | `/api/rental/products/{id}/calendar` | Jadwal spesifik per produk 
| `GET` | `/api/rental/reservations/calendar` | Kalender seluruh penyewaan 

### Reservasi (Booking)
| Method | Endpoint | Keterangan 
| `GET` | `/api/rental/reservations` | 
| `POST` | `/api/rental/reservations` | Pembuatan reservasi baru (menggunakan `down_payment`, tanpa *deposit*) 
| `GET` | `/api/rental/reservations/{id}` | 
| `POST` | `/api/rental/reservations/{id}/cancel` | Membatalkan penyewaan secara transaksional 
| `GET` | `/api/rental/reservations/{id}/pickup` | Laporan penjemputan unit   
| `POST` | `/api/rental/reservations/{id}/pickup` | Konfirmasi unit telah dijemput (Status menjadi `PICKED_UP`) 
| `GET` | `/api/rental/reservations/active` | Status pesanan aktif 
| `GET` | `/api/rental/reservations/upcoming` | Status pesanan mendatang 
| `GET` | `/api/rental/reservations/overdue` | Status pesanan menunggak / lewat batas akhir 

### Pengembalian & Penalti (Returns)
| Method | Endpoint | Keterangan  
| `GET` | `/api/rental/returns` | 
| `POST` | `/api/rental/returns` | Memproses *return* unit, kalkulasi denda otomatis, dan sisa pelunasan akhir
| `GET` | `/api/rental/returns/{id}` | 
| `GET` | `/api/rental/returns/{id}/penalty` | Rincian denda (*late fee* & *damage fee*) 

### Bukti Inspeksi Kondisi Pengembalian (Return Photos)
| Method | Endpoint | Keterangan  
| `GET` | `/api/rental/returns/{id}/photos` | Daftar foto kerusakan/bukti pengembalian 
| `POST` | `/api/rental/returns/{id}/photos` | Unggah foto barang yang dikembalikan 
| `DELETE` | `/api/rental/returns/{id}/photos/{photo_id}` | Hapus foto dari penyewaan 

---

## 📈 6. Finance Service

### Operasional Akuntansi Dasar
| Method | Endpoint | Keterangan  
| `GET` | `/api/finance/accounting-periods` | Daftar periode akuntansi 
| `POST` | `/api/finance/accounting-periods` | Buat periode pembukuan baru 
| `PUT` | `/api/finance/accounting-periods/{id}` | Ubah detail periode pembukuan 
| `DELETE`| `/api/finance/accounting-periods/{id}` | Hapus periode pembukuan 
| `POST` | `/api/finance/period-locks` | Mengunci periode pembukuan   
| `GET` | `/api/finance/daily-closings` | Daftar penutupan harian 
| `POST` | `/api/finance/daily-closings` | Buat penutupan harian 
| `GET` | `/api/finance/daily-closings/{date}` | Detail penutupan per tanggal 
| `GET` | `/api/finance/accounts` | *General Ledger* / *Chart of Accounts* 
| `POST` | `/api/finance/accounts` | Tambah akun baru CoA 
| `GET` | `/api/finance/accounts/{id}` | Detail akun CoA 
| `PUT` | `/api/finance/accounts/{id}` | Ubah akun CoA 
| `DELETE` | `/api/finance/accounts/{id}` | Hapus akun CoA 
| `GET` | `/api/finance/journals` | Tipe/Kategori Jurnal Keuangan (Master) 
| `POST` | `/api/finance/journals` | Tambah kategori jurnal 
| `GET` | `/api/finance/journals/{id}` | Detail kategori jurnal 
| `GET` | `/api/finance/journal-entries` | Daftar entri jurnal/transaksi pembukuan 
| `POST` | `/api/finance/journal-entries` | Catat entri jurnal manual (Debit/Kredit) 
| `GET` | `/api/finance/journal-entries/{id}` | Detail entri jurnal & rinciannya 
| `GET` | `/api/finance/reconciliation` | Proses rekonsiliasi sinkronisasi data 
| `GET` | `/api/finance/reconciliation/logs` | Riwayat hasil rekonsiliasi (kas vs sistem) 

### Laporan Keuangan (Reports) & Dashboard
| Method | Endpoint | Keterangan  
| `GET` | `/api/finance/dashboard/summary` | Ringkasan dasbor 
| `GET` | `/api/finance/reports/profit-loss` | Laporan Laba/Rugi (*Profit/Loss*) 
| `GET` | `/api/finance/reports/cash-flow` | Laporan Arus Kas (*Cash Flow*) 
| `GET` | `/api/finance/reports/balance-sheet` | Laporan Neraca (*Balance Sheet*) 
| `GET` | `/api/finance/reports/export/{format}` | Unduh dokumen cetak 

### Analitik Lanjutan (Analytics)
| Method | Endpoint  
| `GET` | `/api/finance/analytics/sales-daily` 
| `GET` | `/api/finance/analytics/sales-monthly` 
| `GET` | `/api/finance/analytics/top-products` 
| `GET` | `/api/finance/analytics/products` 
| `GET` | `/api/finance/analytics/products/{product_id}` 
| `GET` | `/api/finance/analytics/top-categories` 
| `GET` | `/api/finance/analytics/profit-trend` 
| `GET` | `/api/finance/analytics/expense-trend` 
| `GET` | `/api/finance/analytics/rental-trend` 
| `GET` | `/api/finance/analytics/cashier-performance` 
| `GET` | `/api/finance/analytics/stock-movement` 
| `GET` | `/api/finance/analytics/low-stock` 
| `GET` | `/api/finance/analytics/monthly-summary` 
| `GET` | `/api/finance/analytics/rental-damages` 
| `GET` | `/api/finance/analytics/shifts` 

---

## ⚙️ 7. Internal Endpoint (Service to Service)

Rute-rute ini **TIDAK** terekspos untuk aplikasi klien secara langsung (diamankan via *API Gateway* / *Internal Networking*).

| Method | Endpoint | Keterangan  
| `POST` | `/internal/inventory/stock/deduct` | Pengurangan stok oleh POS / Rental   
| `POST` | `/internal/inventory/stock/restore` | Pemulihan stok karena *refund* / *cancel*   
| `POST` | `/internal/cash/income` | Pencatatan pemasukan otomatis dari POS / Rental   

---

## 📝 Catatan Tabel Internal (Tanpa Endpoint Publik)

Tabel-tabel berikut ini murni dikelola otomatis oleh struktur *backend* dan *background worker*:
- **`outbox_events` (Semua Service)**: Digunakan oleh sistem *worker* / CDC Debezium untuk *Transactional Outbox (Service-to-Service event-driven communication)*.
- **`stock_reservations` (Rental Service)**: Di-insert/update secara otomatis saat pembuatan *booking* / reservasi untuk mengunci kuota stok di hari tersebut.
- **`sale_items` & `rental_items` (Snapshot Data)**: Menyimpan duplikasi atribut master data seperti `product_name` (POS) dan `rental_product_name` (Rental) di sisi transaksional agar nama barang yang tercetak di laporan historis / faktur menjadi statis (*immutable*) dan tidak berubah meski produk aslinya telah dihapus/diedit oleh *Frontend*.