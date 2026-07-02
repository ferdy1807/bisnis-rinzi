-- =====================================================================================
-- SCRIPT INISIALISASI DATABASE BISNIS-RINZI (MULTI-SCHEMA DALAM 1 CONTAINER)
-- =====================================================================================

-- 1. PEMBUATAN LOGICAL DATABASES
CREATE DATABASE auth_db;
CREATE DATABASE inventory_db;
CREATE DATABASE cash_db;
CREATE DATABASE pos_db;
CREATE DATABASE rental_db;
CREATE DATABASE finance_db;

-- Berikan hak akses (Pastikan user rinzi_admin sesuai dengan .env)
GRANT ALL PRIVILEGES ON DATABASE auth_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE inventory_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE cash_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE pos_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE rental_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE finance_db TO postgres;

-- =====================================================================================
-- 2. SCHEMA: AUTH SERVICE (auth_db)
-- =====================================================================================
\connect auth_db

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- public.audit_logs definition
-- Drop table
-- DROP TABLE audit_logs;

create table audit_logs (
	id uuid default uuid_generate_v4() not null,
	user_id uuid not null,
	"action" varchar(100) null,
	entity_name varchar(100) null,
	entity_id uuid null,
	old_data jsonb null,
	new_data jsonb null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint audit_logs_pkey primary key (id)
);
-- public.outbox_events definition
-- Drop table
-- DROP TABLE outbox_events;

create table outbox_events (
	id uuid default uuid_generate_v4() not null,
	aggregate_type varchar(100) not null,
	aggregate_id varchar(100) not null,
	event_type varchar(100) not null,
	payload jsonb not null,
	status varchar(20) default 'PENDING'::character varying null,
	error_message text null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	processed_at timestamptz null,
	constraint outbox_events_pkey primary key (id)
);

create index idx_outbox_status_auth on
public.outbox_events
    using btree (status,
created_at);
-- public.roles definition
-- Drop table
-- DROP TABLE roles;

create table roles (
	code varchar(50) not null,
	"name" varchar(100) not null,
	dashboard_url varchar(255) not null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint roles_pkey primary key (code)
);
-- public.users definition
-- Drop table
-- DROP TABLE users;

create table users (
	id uuid default uuid_generate_v4() not null,
	username varchar(50) not null,
	password_hash text not null,
	full_name varchar(150) not null,
	"role" varchar(50) not null,
	is_active bool default true null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	deleted_at timestamptz null,
	constraint users_pkey primary key (id),
	constraint users_username_key unique (username),
	constraint users_role_fkey foreign key ("role") references roles(code)
);
-- public.refresh_tokens definition
-- Drop table
-- DROP TABLE refresh_tokens;

create table refresh_tokens (
	id uuid default uuid_generate_v4() not null,
	user_id uuid not null,
	"token" text not null,
	expires_at timestamptz not null,
	device_info varchar(255) null,
	ip_address varchar(45) null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	revoked_at timestamptz null,
	constraint refresh_tokens_pkey primary key (id),
	constraint refresh_tokens_token_key unique (token),
	constraint refresh_tokens_user_id_fkey foreign key (user_id) references users(id) on
delete
    cascade
);

INSERT INTO roles (code, name, dashboard_url) VALUES 
('OWNER', 'Pemilik Bisnis', '/admin-dashboard'),
('CASHIER', 'Kasir Toko Retail', '/portal-toko'),
('PEGAWAI', 'Pegawai Sewa Hantaran', '/portal-sewa')
ON CONFLICT (code) DO NOTHING;

-- =====================================================================================
-- 3. SCHEMA: INVENTORY SERVICE (inventory_db)
-- =====================================================================================
\connect inventory_db

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- public.brands definition
-- Drop table
-- DROP TABLE brands;

create table brands (
	id uuid default uuid_generate_v4() not null,
	code varchar(20) not null,
	"name" varchar(100) not null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	deleted_at timestamptz null,
	constraint brands_code_key unique (code),
	constraint brands_pkey primary key (id)
);
-- public.categories definition
-- Drop table
-- DROP TABLE categories;

create table categories (
	id uuid default uuid_generate_v4() not null,
	code varchar(20) not null,
	"name" varchar(100) not null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	deleted_at timestamptz null,
	constraint categories_code_key unique (code),
	constraint categories_pkey primary key (id)
);
-- public.outbox_events definition
-- Drop table
-- DROP TABLE outbox_events;

create table outbox_events (
	id uuid default uuid_generate_v4() not null,
	aggregate_type varchar(100) not null,
	aggregate_id varchar(100) not null,
	event_type varchar(100) not null,
	payload jsonb not null,
	status varchar(20) default 'PENDING'::character varying null,
	error_message text null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	processed_at timestamptz null,
	constraint outbox_events_pkey primary key (id)
);

create index idx_outbox_status_inv on
public.outbox_events
    using btree (status,
created_at);
-- public.sync_versions definition
-- Drop table
-- DROP TABLE sync_versions;
CREATE SEQUENCE IF NOT EXISTS sync_global_version_seq;

create table sync_versions (
	id bigserial not null,
	entity_type varchar(50) not null,
	entity_id uuid not null,
	operation varchar(10) not null,
	version_number int8 default nextval('sync_global_version_seq'::regclass) not null,
	changed_at timestamptz default CURRENT_TIMESTAMP null,
	constraint sync_versions_operation_check check (((operation)::text = any ((array['INSERT'::character varying,
'UPDATE'::character varying,
'DELETE'::character varying])::text[]))),
	constraint sync_versions_pkey primary key (id)
);

create index idx_sync_versions_entity on
public.sync_versions
    using btree (entity_type,
entity_id);

create index idx_sync_versions_version on
public.sync_versions
    using btree (version_number);
-- public.units definition
-- Drop table
-- DROP TABLE units;

create table units (
	id uuid default uuid_generate_v4() not null,
	code varchar(20) not null,
	"name" varchar(100) not null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint units_code_key unique (code),
	constraint units_pkey primary key (id)
);
-- public.products definition
-- Drop table
-- DROP TABLE products;

create table products (
	id uuid default uuid_generate_v4() not null,
	sku varchar(50) not null,
	category_id uuid not null,
	brand_id uuid null,
	"name" varchar(200) not null,
	base_unit_code varchar(20) not null,
	cost_price numeric(18, 2) default 0 not null,
	selling_price numeric(18, 2) default 0 not null,
	is_active bool default true null,
	barcode varchar(100) null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	deleted_at timestamptz null,
	constraint products_barcode_key unique (barcode),
	constraint products_pkey primary key (id),
	constraint products_sku_key unique (sku),
	constraint products_brand_id_fkey foreign key (brand_id) references brands(id),
	constraint products_category_id_fkey foreign key (category_id) references categories(id)
);
-- public.stock_movements definition
-- Drop table
-- DROP TABLE stock_movements;

create table stock_movements (
	id uuid default uuid_generate_v4() not null,
	product_id uuid not null,
	movement_type varchar(50) not null,
	qty numeric(18, 2) not null,
	reference text null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint stock_movements_pkey primary key (id),
	constraint stock_movements_product_id_fkey foreign key (product_id) references products(id)
);
-- public.product_cost_histories definition
-- Drop table
-- DROP TABLE product_cost_histories;

create table product_cost_histories (
	id uuid default uuid_generate_v4() not null,
	product_id uuid not null,
	effective_date date not null,
	average_cost numeric(18, 2) not null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint product_cost_histories_pkey primary key (id),
	constraint product_cost_histories_product_id_fkey foreign key (product_id) references products(id)
);
-- public.product_media definition
-- Drop table
-- DROP TABLE product_media;

create table product_media (
	id uuid default uuid_generate_v4() not null,
	product_id uuid not null,
	media_category varchar(30) not null,
	bucket_name varchar(100) not null,
	object_name text not null,
	original_file_name varchar(255) null,
	mime_type varchar(100) null,
	file_size_bytes int8 null,
	is_active bool default true null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint product_media_pkey primary key (id),
	constraint product_media_product_id_fkey foreign key (product_id) references products(id) on
delete
    cascade
);

create index idx_product_media_pid on
public.product_media
    using btree (product_id);
-- public.product_stocks definition
-- Drop table
-- DROP TABLE product_stocks;

create table product_stocks (
	product_id uuid not null,
	qty numeric(18, 2) default 0 null,
	qty_min_stock numeric(18, 2) default 0 null,
	qty_safety_stock numeric(18, 2) default 0 null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint product_stocks_pkey primary key (product_id),
	constraint product_stocks_qty_check check ((qty >= (0)::numeric)),
	constraint product_stocks_qty_min_stock_check check ((qty_min_stock >= (0)::numeric)),
	constraint product_stocks_qty_safety_stock_check check ((qty_safety_stock >= (0)::numeric)),
	constraint product_stocks_product_id_fkey foreign key (product_id) references products(id) on
delete
    cascade
);

-- =====================================================================================
-- 4. SCHEMA: CASH SERVICE (cash_db)
-- =====================================================================================
\connect cash_db

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- public.cashier_sessions definition
-- Drop table
-- DROP TABLE cashier_sessions;

create table cashier_sessions (
	id uuid default uuid_generate_v4() not null,
	cashier_id uuid not null,
	open_time timestamptz not null,
	close_time timestamptz null,
	opening_cash numeric(18, 2) default 0 null,
	expected_cash numeric(18, 2) null,
	actual_cash numeric(18, 2) null,
	difference numeric(18, 2) null,
	status varchar(20) default 'OPEN'::character varying null,
	receipt_url text null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint cashier_sessions_pkey primary key (id)
);
-- public.expense_categories definition
-- Drop table
-- DROP TABLE expense_categories;

create table expense_categories (
	id uuid default uuid_generate_v4() not null,
	code varchar(20) not null,
	"name" varchar(100) not null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint expense_categories_code_key unique (code),
	constraint expense_categories_pkey primary key (id)
);
-- public.outbox_events definition
-- Drop table
-- DROP TABLE outbox_events;

create table outbox_events (
	id uuid default uuid_generate_v4() not null,
	aggregate_type varchar(100) not null,
	aggregate_id varchar(100) not null,
	event_type varchar(100) not null,
	payload jsonb not null,
	status varchar(20) default 'PENDING'::character varying null,
	error_message text null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	processed_at timestamptz null,
	constraint outbox_events_pkey primary key (id)
);

create index idx_outbox_status_cash on
public.outbox_events
    using btree (status,
created_at);
-- public.cash_transactions definition
-- Drop table
-- DROP TABLE cash_transactions;

create table cash_transactions (
	id uuid default uuid_generate_v4() not null,
	session_id uuid not null,
	transaction_type varchar(20) not null,
	reference_type varchar(50) not null,
	reference_id uuid null,
	amount numeric(18, 2) not null,
	notes text null,
	created_by uuid not null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint cash_transactions_pkey primary key (id),
	constraint cash_transactions_session_id_fkey foreign key (session_id) references cashier_sessions(id) on
delete
    cascade
);

create index idx_cash_trx_session on
public.cash_transactions
    using btree (session_id);
-- public.expenses definition
-- Drop table
-- DROP TABLE expenses;

create table expenses (
	id uuid default uuid_generate_v4() not null,
	expense_date timestamptz default CURRENT_TIMESTAMP null,
	category_id uuid not null,
	description text not null,
	amount numeric(18, 2) not null,
	created_by uuid not null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint expenses_pkey primary key (id),
	constraint expenses_category_id_fkey foreign key (category_id) references expense_categories(id)
);


-- =====================================================================================
-- 5. SCHEMA: POS SERVICE (pos_db)
-- =====================================================================================
\connect pos_db

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- public.outbox_events definition
-- Drop table
-- DROP TABLE outbox_events;

create table outbox_events (
	id uuid default uuid_generate_v4() not null,
	aggregate_type varchar(100) not null,
	aggregate_id varchar(100) not null,
	event_type varchar(100) not null,
	payload jsonb not null,
	status varchar(20) default 'PENDING'::character varying null,
	error_message text null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	processed_at timestamptz null,
	constraint outbox_events_pkey primary key (id)
);

create index idx_outbox_status_pos on
public.outbox_events
    using btree (status,
created_at);
-- public.sales definition
-- Drop table
-- DROP TABLE sales;

create table sales (
	id uuid default uuid_generate_v4() not null,
	idempotency_key varchar(100) not null,
	invoice_number varchar(50) not null,
	transaction_date timestamptz not null,
	subtotal numeric(18, 2) not null,
	discount numeric(18, 2) default 0 null,
	total numeric(18, 2) not null,
	amount_paid numeric(18, 2) default 0 null,
	change_amount numeric(18, 2) default 0 null,
	payment_method varchar(20) default 'CASH'::character varying not null,
	payment_status varchar(20) default 'COMPLETED'::character varying not null,
	cashier_id uuid not null,
	cashier_session_id uuid not null,
	invoice_url text null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint sales_idempotency_key_key unique (idempotency_key),
	constraint sales_invoice_number_key unique (invoice_number),
	constraint sales_pkey primary key (id)
);

create index idx_sales_invoice on
public.sales
    using btree (invoice_number);

create index idx_sales_session on
public.sales
    using btree (cashier_session_id);
-- public.sync_logs definition
-- Drop table
-- DROP TABLE sync_logs;

create table sync_logs (
	id uuid default uuid_generate_v4() not null,
	entity_type varchar(100) null,
	entity_id uuid null,
	sync_status varchar(30) null,
	error_message text null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint sync_logs_pkey primary key (id)
);
-- public.sale_items definition
-- Drop table
-- DROP TABLE sale_items;

create table sale_items (
	id uuid default uuid_generate_v4() not null,
	sale_id uuid not null,
	product_id uuid not null,
	product_name varchar(255) default 'Unknown'::character varying not null,
	unit_code varchar(20) not null,
	qty numeric(18, 2) not null,
	unit_price numeric(18, 2) not null,
	discount numeric(18, 2) default 0 null,
	subtotal numeric(18, 2) not null,
	cost_price numeric(18, 2) not null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint sale_items_pkey primary key (id),
	constraint sale_items_sale_id_fkey foreign key (sale_id) references sales(id) on
delete
    cascade
);

-- =====================================================================================
-- 6. SCHEMA: RENTAL SERVICE (rental_db)
-- =====================================================================================
\connect rental_db

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- public.customer_snapshots definition
-- Drop table
-- DROP TABLE customer_snapshots;

create table customer_snapshots (
	id uuid default uuid_generate_v4() not null,
	customer_name varchar(150) not null,
	customer_phone varchar(30) not null,
	customer_id_card varchar(50) not null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint customer_snapshots_pkey primary key (id)
);
-- public.outbox_events definition
-- Drop table
-- DROP TABLE outbox_events;

create table outbox_events (
	id uuid default uuid_generate_v4() not null,
	aggregate_type varchar(100) not null,
	aggregate_id varchar(100) not null,
	event_type varchar(100) not null,
	payload jsonb not null,
	status varchar(20) default 'PENDING'::character varying null,
	error_message text null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	processed_at timestamptz null,
	constraint outbox_events_pkey primary key (id)
);

create index idx_outbox_status_rental on
public.outbox_events
    using btree (status,
created_at);
-- public.rental_categories definition
-- Drop table
-- DROP TABLE rental_categories;

create table rental_categories (
	id uuid default uuid_generate_v4() not null,
	code varchar(20) not null,
	"name" varchar(100) not null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	deleted_at timestamptz null,
	constraint rental_categories_code_key unique (code),
	constraint rental_categories_pkey primary key (id)
);

insert into rental_categories(code,"name")values
('KTK-MIKA', 'Kotak Mika'),
('KTK-AKRI', 'Kotak Akrilik'),
('KTK-PTGL','Kotak Pentagonal'),
('KTK-KAYU','Kotak Kayu Jati Belanda'),
('KTK-RUST','Kotak Rustic'),
('KRJ-ROTN','Keranjang Rotan')
ON CONFLICT (code) DO NOTHING;

-- public.rental_products definition
-- Drop table
-- DROP TABLE rental_products;

CREATE TABLE rental_products (
    id UUID DEFAULT uuid_generate_v4() NOT NULL,
    category_id UUID NOT NULL,
    code VARCHAR(50) NOT NULL,
    "name" VARCHAR(200) NOT NULL,
    description TEXT NULL,
    rental_price NUMERIC(18,2) NOT NULL,
    deposit_amount NUMERIC(18,2) DEFAULT 0 NULL,
    quantity_available NUMERIC(18,2) DEFAULT 0 NULL,
    is_active BOOLEAN DEFAULT TRUE NULL,
    object_name TEXT NOT NULL,
    original_file_name VARCHAR(255) NULL,
    mime_type VARCHAR(100) NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NULL,
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT rental_products_code_key UNIQUE (code),
    CONSTRAINT rental_products_pkey PRIMARY KEY (id)
);

ALTER TABLE public.rental_products
ADD CONSTRAINT rental_products_category_id_fkey
FOREIGN KEY (category_id)
REFERENCES rental_categories(id);
-- public.rental_reservations definition
-- Drop table
-- DROP TABLE rental_reservations;

create table rental_reservations (
	id uuid default uuid_generate_v4() not null,
	invoice_number varchar(50) not null,
	customer_snapshot_id uuid not null,
	transaction_date timestamptz not null,
	start_date date not null,
	end_date date not null,
	event_date date null,
	subtotal numeric(18, 2) not null,
	down_payment numeric(18, 2) default 0 null,
	amount_paid numeric(18, 2) default 0 null,
	change_amount numeric(18, 2) default 0 null,
	total_amount numeric(18, 2) not null,
	status varchar(20) default 'BOOKED'::character varying not null,
	picked_up_by uuid null,
	picked_up_at timestamptz null,
	cashier_session_id uuid not null,
	created_by uuid not null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint rental_reservations_invoice_number_key unique (invoice_number),
	constraint rental_reservations_pkey primary key (id),
	constraint rental_reservations_status_check check (((status)::text = any (array[('BOOKED'::character varying)::text,
('CONTENTS_RECEIVED'::character varying)::text,
('DECORATING'::character varying)::text,
('READY_FOR_PICKUP'::character varying)::text,
('PICKED_UP'::character varying)::text,
('RETURNED'::character varying)::text,
('CANCELLED'::character varying)::text]))),
	constraint rental_reservations_customer_snapshot_id_fkey foreign key (customer_snapshot_id) references customer_snapshots(id)
);
-- public.rental_returns definition
-- Drop table
-- DROP TABLE rental_returns;

create table rental_returns (
	id uuid default uuid_generate_v4() not null,
	rental_reservation_id uuid not null,
	return_date timestamptz default CURRENT_TIMESTAMP null,
	late_days int4 default 0 not null,
	total_late_fees numeric(18, 2) default 0 not null,
	total_damage_fees numeric(18, 2) default 0 not null,
	remaining_payment numeric(18, 2) default 0 not null,
	amount_paid numeric(18, 2) default 0 null,
	change_amount numeric(18, 2) default 0 null,
	grand_total_paid numeric(18, 2) default 0 not null,
	notes text null,
	received_by uuid not null,
	receipt_url text null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint rental_returns_pkey primary key (id),
	constraint rental_returns_rental_reservation_id_fkey foreign key (rental_reservation_id) references rental_reservations(id)
);
-- public.stock_reservations definition
-- Drop table
-- DROP TABLE stock_reservations;

create table stock_reservations (
	id uuid default uuid_generate_v4() not null,
	rental_product_id uuid not null,
	rental_reservation_id uuid not null,
	reserved_date date not null,
	qty_reserved numeric(18, 2) not null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint stock_reservations_pkey primary key (id),
	constraint stock_reservations_rental_product_id_rental_reservation_id__key unique (rental_product_id,
rental_reservation_id,
reserved_date),
	constraint stock_reservations_rental_product_id_fkey foreign key (rental_product_id) references rental_products(id) on
delete
    cascade,
    constraint stock_reservations_rental_reservation_id_fkey foreign key (rental_reservation_id) references rental_reservations(id) on
    delete
        cascade
);

create index idx_stock_reservation_date on
public.stock_reservations
    using btree (rental_product_id,
reserved_date);
-- public.rental_items definition
-- Drop table
-- DROP TABLE rental_items;

create table rental_items (
	id uuid default uuid_generate_v4() not null,
	rental_reservation_id uuid not null,
	rental_product_id uuid not null,
	rental_product_name varchar(255) default 'Unknown'::character varying not null,
	qty numeric(18, 2) default 1 null,
	price_per_period numeric(18, 2) not null,
	subtotal numeric(18, 2) not null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint rental_items_pkey primary key (id),
	constraint rental_items_rental_product_id_fkey foreign key (rental_product_id) references rental_products(id),
	constraint rental_items_rental_reservation_id_fkey foreign key (rental_reservation_id) references rental_reservations(id) on
delete
    cascade
);
-- public.rental_product_media definition
-- Drop table
-- DROP TABLE rental_product_media;

create table rental_product_media (
	id uuid default uuid_generate_v4() not null,
	rental_product_id uuid not null,
	bucket_name varchar(100) not null,
	object_name text not null,
	original_file_name varchar(255) null,
	mime_type varchar(100) null,
	file_size_bytes int8 null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint rental_product_media_pkey primary key (id),
	constraint rental_product_media_rental_product_id_fkey foreign key (rental_product_id) references rental_products(id)
);
-- public.rental_reservation_contents definition
-- Drop table
-- DROP TABLE rental_reservation_contents;

create table rental_reservation_contents (
	id uuid default uuid_generate_v4() not null,
	rental_reservation_id uuid not null,
	item_name varchar(200) not null,
	description text null,
	quantity int4 default 1 not null,
	condition_notes text null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint rental_reservation_contents_pkey primary key (id),
	constraint rental_reservation_contents_rental_reservation_id_fkey foreign key (rental_reservation_id) references rental_reservations(id)
);
-- public.rental_return_items definition
-- Drop table
-- DROP TABLE rental_return_items;

create table rental_return_items (
	id uuid default uuid_generate_v4() not null,
	rental_return_id uuid not null,
	rental_product_id uuid not null,
	rental_product_name varchar(255) default 'Unknown'::character varying not null,
	qty_returned numeric(18, 2) not null,
	condition_status varchar(50) default 'GOOD'::character varying not null,
	damage_fee numeric(18, 2) default 0 not null,
	condition_notes text null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint rental_return_items_pkey primary key (id),
	constraint rental_return_items_rental_product_id_fkey foreign key (rental_product_id) references rental_products(id),
	constraint rental_return_items_rental_return_id_fkey foreign key (rental_return_id) references rental_returns(id) on
delete
    cascade
);
-- public.rental_return_photos definition
-- Drop table
-- DROP TABLE rental_return_photos;

create table rental_return_photos (
	id uuid default uuid_generate_v4() not null,
	rental_return_id uuid not null,
	rental_return_item_id uuid null,
	bucket_name varchar(100) not null,
	object_name text not null,
	original_file_name varchar(255) null,
	mime_type varchar(100) null,
	file_size_bytes int8 null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint rental_return_photos_pkey primary key (id),
	constraint rental_return_photos_rental_return_id_fkey foreign key (rental_return_id) references rental_returns(id) on delete
    cascade,
    constraint rental_return_photos_rental_return_item_id_fkey foreign key (rental_return_item_id) references rental_return_items(id));

create index idx_rental_return_photos_ret on public.rental_return_photos using btree (rental_return_item_id);


-- =====================================================================================
-- 7. SCHEMA: FINANCE SERVICE (finance_db)
-- =====================================================================================
\connect finance_db

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- public.accounting_periods definition
-- Drop table
-- DROP TABLE accounting_periods;

create table accounting_periods (
	id uuid default uuid_generate_v4() not null,
	"name" varchar(100) not null,
	start_date date not null,
	end_date date not null,
	is_closed bool default false null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint accounting_periods_name_key unique (name),
	constraint accounting_periods_pkey primary key (id)
);
-- public.chart_of_accounts definition
-- Drop table
-- DROP TABLE chart_of_accounts;

create table chart_of_accounts (
	id uuid default uuid_generate_v4() not null,
	account_code varchar(20) not null,
	account_name varchar(100) not null,
	account_group varchar(50) null,
	normal_balance varchar(10) null,
	current_balance numeric(18, 2) default 0 null,
	is_active bool default true null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint chart_of_accounts_account_code_key unique (account_code),
	constraint chart_of_accounts_pkey primary key (id)
);
INSERT INTO chart_of_accounts (
    account_code,
    account_name,
    account_group,
    normal_balance
)
VALUES
('110000', 'Kas & Bank',             'ASSET',     'DEBIT'),
('120000', 'Piutang',                'ASSET',     'DEBIT'),
('130000', 'Persediaan',             'ASSET',     'DEBIT'),
('150000', 'Aset Tetap',             'ASSET',     'DEBIT'),
('210000', 'Utang Usaha',            'LIABILITY', 'CREDIT'),
('230000', 'Down Payment Pelanggan', 'LIABILITY', 'CREDIT'),
('310000', 'Modal Pemilik',          'EQUITY',    'CREDIT'),
('340000', 'Laba Ditahan',           'EQUITY',    'CREDIT'),
('410000', 'Penjualan Barang',       'REVENUE',   'CREDIT'),
('420000', 'Pendapatan Sewa',        'REVENUE',   'CREDIT'),
('421000', 'Denda Sewa',             'REVENUE',   'CREDIT'),
('610000', 'Gaji Karyawan',          'EXPENSE',   'DEBIT'),
('710000', 'Harga Pokok Penjualan',  'EXPENSE',   'DEBIT'),
('711000', 'Biaya Perawatan Rental', 'EXPENSE',   'DEBIT'),
('720000', 'Biaya Operasional',      'EXPENSE',   'DEBIT'),
('750000', 'Penyusutan',             'EXPENSE',   'DEBIT'),
('810000', 'Beban Bunga & Lainnya',  'EXPENSE',   'DEBIT');


-- public.daily_closings definition
-- Drop table
-- DROP TABLE daily_closings;

create table daily_closings (
	id uuid default uuid_generate_v4() not null,
	closing_date date not null,
	total_sales_retail numeric(18, 2) default 0 null,
	total_rental_income numeric(18, 2) default 0 null,
	total_other_income numeric(18, 2) default 0 null,
	total_expenses numeric(18, 2) default 0 null,
	net_cash_flow numeric(18, 2) default 0 null,
	is_reconciled bool default false null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint daily_closings_closing_date_key unique (closing_date),
	constraint daily_closings_pkey primary key (id)
);
-- public.finance_monthly_analytics definition
-- Drop table
-- DROP TABLE finance_monthly_analytics;

create table finance_monthly_analytics (
	id uuid default uuid_generate_v4() not null,
	month_year varchar(7) not null,
	total_revenue numeric(18, 2) default 0 null,
	total_hpp numeric(18, 2) default 0 null,
	total_expenses numeric(18, 2) default 0 null,
	total_rental_penalties numeric(18, 2) default 0 null,
	calculated_profit numeric(18, 2) default 0 null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint finance_monthly_analytics_month_year_key unique (month_year),
	constraint finance_monthly_analytics_pkey primary key (id)
);
-- public.finance_product_analytics definition
-- Drop table
-- DROP TABLE finance_product_analytics;

create table finance_product_analytics (
	id uuid default uuid_generate_v4() not null,
	log_date date not null,
	product_id uuid not null,
	product_name varchar(150) not null,
	category_id uuid null,
	category_name varchar(100) null,
	business_type varchar(20) not null,
	qty_sold_or_rented numeric(18, 2) default 0 null,
	total_revenue numeric(18, 2) default 0 null,
	total_hpp numeric(18, 2) default 0 null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint finance_product_analytics_pkey primary key (id)
);

create index idx_finance_analytics_date on
public.finance_product_analytics
    using btree (log_date);
-- public.finance_rental_damage_logs definition
-- Drop table
-- DROP TABLE finance_rental_damage_logs;

create table finance_rental_damage_logs (
	id uuid default uuid_generate_v4() not null,
	log_date date not null,
	rental_return_id uuid not null,
	rental_product_id uuid not null,
	product_name varchar(150) not null,
	penalty_amount numeric(18, 2) default 0 null,
	condition_status varchar(50) not null,
	notes text null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint finance_rental_damage_logs_pkey primary key (id)
);
-- public.journals definition
-- Drop table
-- DROP TABLE journals;

create table journals (
	id uuid default uuid_generate_v4() not null,
	journal_code varchar(50) not null,
	"name" varchar(100) not null,
	description text null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint journals_journal_code_key unique (journal_code),
	constraint journals_pkey primary key (id)
);
-- public.outbox_events definition
-- Drop table
-- DROP TABLE outbox_events;

create table outbox_events (
	id uuid default uuid_generate_v4() not null,
	aggregate_type varchar(100) not null,
	aggregate_id varchar(100) not null,
	event_type varchar(100) not null,
	payload jsonb not null,
	status varchar(20) default 'PENDING'::character varying null,
	error_message text null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	processed_at timestamptz null,
	constraint outbox_events_pkey primary key (id)
);

create index idx_outbox_status_fin on
public.outbox_events
    using btree (status,
created_at);
-- public.journal_entries definition
-- Drop table
-- DROP TABLE journal_entries;

create table journal_entries (
	id uuid default uuid_generate_v4() not null,
	journal_id uuid not null,
	accounting_period_id uuid not null,
	reference_number varchar(100) not null,
	entry_date timestamptz not null,
	narration text null,
	is_posted bool default false null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint journal_entries_pkey primary key (id),
	constraint journal_entries_accounting_period_id_fkey foreign key (accounting_period_id) references accounting_periods(id),
	constraint journal_entries_journal_id_fkey foreign key (journal_id) references journals(id) on delete cascade
);
-- public.journal_entry_details definition
-- Drop table
-- DROP TABLE journal_entry_details;

create table journal_entry_details (
	id uuid default uuid_generate_v4() not null,
	journal_entry_id uuid not null,
	account_id uuid not null,
	debit_amount numeric(18, 2) default 0 null,
	credit_amount numeric(18, 2) default 0 null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	constraint journal_entry_details_pkey primary key (id),
	constraint journal_entry_details_account_id_fkey foreign key (account_id) references chart_of_accounts(id),
	constraint journal_entry_details_journal_entry_id_fkey foreign key (journal_entry_id) references journal_entries(id) on
delete
    cascade
);
-- public.period_locks definition
-- Drop table
-- DROP TABLE period_locks;

create table period_locks (
	id uuid default uuid_generate_v4() not null,
	accounting_period_id uuid not null,
	locked_by uuid not null,
	lock_reason text null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	updated_at timestamptz default CURRENT_TIMESTAMP null,
	constraint period_locks_pkey primary key (id),
	constraint period_locks_accounting_period_id_fkey foreign key (accounting_period_id) references accounting_periods(id)
);
-- public.reconciliation_logs definition
-- Drop table
-- DROP TABLE reconciliation_logs;

create table reconciliation_logs (
	id uuid default uuid_generate_v4() not null,
	daily_closing_id uuid not null,
	target_system varchar(100) null,
	system_amount numeric(18, 2) null,
	actual_amount numeric(18, 2) null,
	discrepancy numeric(18, 2) null,
	notes text null,
	reconciled_by uuid null,
	created_at timestamptz default CURRENT_TIMESTAMP null,
	constraint reconciliation_logs_pkey primary key (id),
	constraint reconciliation_logs_daily_closing_id_fkey foreign key (daily_closing_id) references daily_closings(id) on
delete
    cascade
);
