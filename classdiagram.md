Berdasarkan `usecase.md` dan struktur database dari `init.sql`, berikut adalah representasi Class Diagram dari entitas-entitas utama dalam sistem. Relasi antar entitas mengikuti struktur Foreign Key pada schema database.

### 1. Modul Manajemen User & Autentikasi (auth_db)

---

## users

- id: uuid
- username: varchar
- password_hash: text
- full_name: varchar
- role: varchar
- is_active: bool
- created_at: timestamptz
- updated_at: timestamptz
- deleted_at: timestamptz

---

- tambah()
- edit()
- hapus()
- login()
- logout()

---

## roles

- code: varchar
- name: varchar
- dashboard_url: varchar
- created_at: timestamptz
- updated_at: timestamptz

---

- tambah()
- edit()
- hapus()

---

## refresh_tokens

- id: uuid
- user_id: uuid (FK to users)
- token: text
- expires_at: timestamptz
- device_info: varchar
- ip_address: varchar

---

- tambah()
- hapus()

### 2. Modul Inventaris & Produk (inventory_db)

---

## products

- id: uuid
- sku: varchar
- category_id: uuid (FK to categories)
- brand_id: uuid (FK to brands)
- name: varchar
- base_unit_code: varchar
- cost_price: numeric
- selling_price: numeric
- is_active: bool
- barcode: varchar

---

- tambah()
- edit()
- hapus()

---

## product_stocks

- product_id: uuid (FK to products)
- qty: numeric
- qty_min_stock: numeric
- qty_safety_stock: numeric

---

- tambah()
- edit()
- hapus()
- update_stock()

---

## categories

- id: uuid
- code: varchar
- name: varchar

---

- tambah()
- edit()
- hapus()

---

## brands

- id: uuid
- code: varchar
- name: varchar

---

- tambah()
- edit()
- hapus()

---

## units

- id: uuid
- code: varchar
- name: varchar

---

- tambah()
- edit()
- hapus()

### 3. Modul Kasir & Arus Kas (cash_db & pos_db)

---

## cashier_sessions

- id: uuid
- cashier_id: uuid (FK to users)
- open_time: timestamptz
- close_time: timestamptz
- opening_cash: numeric
- expected_cash: numeric
- actual_cash: numeric
- difference: numeric
- status: varchar
- receipt_url: text

---

- buka_shift()
- tutup_shift()
- edit()
- hapus()

---

## cash_transactions

- id: uuid
- session_id: uuid (FK to cashier_sessions)
- transaction_type: varchar
- reference_type: varchar
- reference_id: uuid
- amount: numeric
- notes: text

---

- tambah()
- edit()
- hapus()

---

## sales

- id: uuid
- invoice_number: varchar
- transaction_date: timestamptz
- subtotal: numeric
- discount: numeric
- total: numeric
- amount_paid: numeric
- change_amount: numeric
- payment_method: varchar
- payment_status: varchar
- cashier_id: uuid (FK to users)
- cashier_session_id: uuid (FK to cashier_sessions)

---

- tambah()
- edit()
- hapus()
- cetak_struk()

---

## sale_items

- id: uuid
- sale_id: uuid (FK to sales)
- product_id: uuid
- product_name: varchar
- qty: numeric
- unit_price: numeric
- discount: numeric
- subtotal: numeric

---

- tambah()
- edit()
- hapus()

### 4. Modul Penyewaan (rental_db)

---

## rental_products

- id: uuid
- category_id: uuid (FK to rental_categories)
- code: varchar
- name: varchar
- description: text
- rental_price: numeric
- deposit_amount: numeric
- quantity_available: numeric
- is_active: bool

---

- tambah()
- edit()
- hapus()

---

## rental_reservations

- id: uuid
- invoice_number: varchar
- customer_snapshot_id: uuid (FK to customer_snapshots)
- transaction_date: timestamptz
- start_date: date
- end_date: date
- event_date: date
- total_amount: numeric
- status: varchar
- cashier_session_id: uuid

---

- tambah()
- edit()
- hapus()
- update_status()

---

## rental_items

- id: uuid
- rental_reservation_id: uuid (FK to rental_reservations)
- rental_product_id: uuid (FK to rental_products)
- qty: numeric
- price_per_period: numeric
- subtotal: numeric

---

- tambah()
- edit()
- hapus()

---

## rental_returns

- id: uuid
- rental_reservation_id: uuid (FK to rental_reservations)
- return_date: timestamptz
- late_days: int
- total_late_fees: numeric
- total_damage_fees: numeric
- grand_total_paid: numeric

---

- tambah()
- edit()
- hapus()

---

## rental_return_items

- id: uuid
- rental_return_id: uuid (FK to rental_returns)
- rental_product_id: uuid
- qty_returned: numeric
- condition_status: varchar
- damage_fee: numeric

---

- tambah()
- edit()
- hapus()

---

## customer_snapshots

- id: uuid
- customer_name: varchar
- customer_phone: varchar
- customer_id_card: varchar

---

- tambah()
- edit()
- hapus()

### 5. Modul Laporan & Keuangan (finance_db)

---

## daily_closings

- id: uuid
- closing_date: date
- total_sales_retail: numeric
- total_rental_income: numeric
- total_expenses: numeric
- net_cash_flow: numeric
- is_reconciled: bool

---

- tambah()
- edit()
- hapus()
- rekonsiliasi()

---

## journals

- id: uuid
- journal_code: varchar
- name: varchar
- description: text

---

- tambah()
- edit()
- hapus()

---

## journal_entries

- id: uuid
- journal_id: uuid (FK to journals)
- accounting_period_id: uuid (FK to accounting_periods)
- reference_number: varchar
- entry_date: timestamptz
- is_posted: bool

---

- tambah()
- edit()
- hapus()
- post_journal()

---

## journal_entry_details

- id: uuid
- journal_entry_id: uuid (FK to journal_entries)
- account_id: uuid (FK to chart_of_accounts)
- debit_amount: numeric
- credit_amount: numeric

---

- tambah()
- edit()
- hapus()

---

## chart_of_accounts

- id: uuid
- account_code: varchar
- account_name: varchar
- account_group: varchar
- normal_balance: varchar
- current_balance: numeric

---

- tambah()
- edit()
- hapus()

### 6. Modul Event & Integrasi Data (Lintas Service)

---

## outbox_events

- id: uuid
- aggregate_type: varchar
- aggregate_id: varchar
- event_type: varchar
- payload: jsonb
- status: varchar
- error_message: text
- processed_at: timestamptz

---

- tambah()
- process_event()
- mark_as_done()

---

## sync_versions

- id: bigserial
- entity_type: varchar
- entity_id: uuid
- operation: varchar
- version_number: int8
- changed_at: timestamptz

---

- tambah()
- edit()
- hapus()

---

## sync_logs

- id: uuid
- entity_type: varchar
- entity_id: uuid
- sync_status: varchar
- error_message: text

---

- tambah()
- log_error()
