-- 001_create_schema.sql: SmashHub Database Schema
-- Badminton Court Booking & Mini Store POS System

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 1. Create ENUM Types
DO $$ BEGIN
    CREATE TYPE booking_status AS ENUM (
        'PENDING_PAYMENT',
        'CONFIRMED',
        'MANUAL_REVIEW',
        'EXPIRED',
        'CANCELLED'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE slip_status AS ENUM (
        'VERIFIED',
        'REJECTED_DUPLICATE',
        'REJECTED_AMOUNT_MISMATCH',
        'REJECTED_WRONG_ACCOUNT',
        'API_ERROR_PENDING_ADMIN'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- 2. Create Tables

-- 2.1 Users Table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    line_user_id VARCHAR(100) UNIQUE,
    display_name VARCHAR(255) NOT NULL,
    picture_url TEXT,
    phone VARCHAR(20),
    role VARCHAR(20) NOT NULL DEFAULT 'CUSTOMER' CHECK (role IN ('CUSTOMER', 'ADMIN')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2.2 Courts Table
CREATE TABLE IF NOT EXISTS courts (
    id SERIAL PRIMARY KEY,
    court_number INT NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    court_type VARCHAR(20) NOT NULL CHECK (court_type IN ('PVC', 'ACRYLIC')),
    normal_rate NUMERIC(10, 2) NOT NULL,
    peak_rate NUMERIC(10, 2) NOT NULL, -- rate applies after 17:00
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'MAINTENANCE')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2.3 Products Table (Mini Store POS)
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL CHECK (category IN ('SHUTTLECOCK', 'DRINK', 'RENTAL')),
    price NUMERIC(10, 2) NOT NULL,
    stock_qty INT NOT NULL CHECK (stock_qty >= 0),
    image_url TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2.4 Bookings Table
CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_no VARCHAR(50) NOT NULL UNIQUE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    status booking_status NOT NULL DEFAULT 'PENDING_PAYMENT',
    court_total NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    item_total NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    net_total NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    check_in_status VARCHAR(20) NOT NULL DEFAULT 'WAITING' CHECK (check_in_status IN ('WAITING', 'CHECKED_IN', 'COMPLETED')),
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2.5 Booking Courts Table (Slot allocation)
CREATE TABLE IF NOT EXISTS booking_courts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    court_id INT NOT NULL REFERENCES courts(id) ON DELETE RESTRICT,
    play_date DATE NOT NULL,
    start_hour INT NOT NULL CHECK (start_hour >= 10 AND start_hour <= 21),
    end_hour INT NOT NULL CHECK (end_hour >= 11 AND end_hour <= 22),
    hourly_rate NUMERIC(10, 2) NOT NULL,
    is_slot_locked BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2.6 Booking Items Table
CREATE TABLE IF NOT EXISTS booking_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    qty INT NOT NULL CHECK (qty > 0),
    unit_price NUMERIC(10, 2) NOT NULL,
    subtotal NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2.7 Payment Slips Table (4-Step automated verification)
CREATE TABLE IF NOT EXISTS payment_slips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    slip_image_url TEXT,
    trans_ref VARCHAR(100) UNIQUE,
    sender_name VARCHAR(255),
    receiver_account VARCHAR(100),
    amount NUMERIC(10, 2) NOT NULL,
    verification_status slip_status NOT NULL,
    raw_payload JSONB,
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 3. Indexes & Constraints

-- 3.1 CRITICAL: Partial Unique Index to Prevent Double Booking
-- Only locked slots (is_slot_locked = TRUE) conflict.
-- When a booking becomes EXPIRED or CANCELLED, is_slot_locked is updated to FALSE,
-- instantly freeing up the slot for other customers.
CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_active_court_slot 
ON booking_courts (court_id, play_date, start_hour) 
WHERE is_slot_locked = TRUE;

-- 3.2 Performance & Lookup Indexes
CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status);
CREATE INDEX IF NOT EXISTS idx_bookings_expires_at ON bookings(expires_at) WHERE status = 'PENDING_PAYMENT';
CREATE INDEX IF NOT EXISTS idx_booking_courts_booking_id ON booking_courts(booking_id);
CREATE INDEX IF NOT EXISTS idx_booking_courts_date_court ON booking_courts(play_date, court_id);
CREATE INDEX IF NOT EXISTS idx_booking_items_booking_id ON booking_items(booking_id);
CREATE INDEX IF NOT EXISTS idx_payment_slips_booking_id ON payment_slips(booking_id);
CREATE INDEX IF NOT EXISTS idx_payment_slips_trans_ref ON payment_slips(trans_ref);
