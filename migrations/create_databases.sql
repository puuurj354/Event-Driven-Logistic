-- Script untuk membuat semua database yang diperlukan
-- Jalankan dengan: psql -U postgres -f create_databases.sql

-- Create database untuk Order Service
CREATE DATABASE db_order;

-- Create database untuk Payment Service
CREATE DATABASE db_payment;

-- Create database untuk Inventory Service
CREATE DATABASE db_inventory;

-- Create database untuk Delivery Service
CREATE DATABASE db_delivery;

-- Enable UUID extension untuk semua database (diperlukan untuk UUID generation)
\c db_order
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\c db_payment
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\c db_inventory
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\c db_delivery
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tampilkan list databases yang berhasil dibuat
\l
