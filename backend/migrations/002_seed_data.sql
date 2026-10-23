-- 002_seed_data.sql: Realistic Thai Seed Data for SmashHub
-- Badminton Courts, Mini-Store Products, Test Users, and Today's Bookings (All 4 Statuses)

-- Clear previous seed data if re-running
TRUNCATE TABLE payment_slips, booking_items, booking_courts, bookings, products, courts, users RESTART IDENTITY CASCADE;

-- 1. Seed Badminton Courts (6 Courts)
INSERT INTO courts (court_number, name, court_type, normal_rate, peak_rate, status) VALUES
(1, 'คอร์ท 1 (PVC มาตรฐาน)', 'PVC', 150.00, 220.00, 'ACTIVE'),
(2, 'คอร์ท 2 (PVC มาตรฐาน)', 'PVC', 150.00, 220.00, 'ACTIVE'),
(3, 'คอร์ท 3 (PVC มาตรฐาน)', 'PVC', 150.00, 220.00, 'ACTIVE'),
(4, 'คอร์ท 4 (PVC มาตรฐาน)', 'PVC', 150.00, 220.00, 'ACTIVE'),
(5, 'คอร์ท 5 (แข่งขัน พื้นยางสังเคราะห์ Acrylic VIP)', 'ACRYLIC', 180.00, 260.00, 'ACTIVE'),
(6, 'คอร์ท 6 (แข่งขัน พื้นยางสังเคราะห์ Acrylic VIP)', 'ACRYLIC', 180.00, 260.00, 'ACTIVE');

-- 2. Seed Mini Store Products (8 Items)
INSERT INTO products (sku, name, category, price, stock_qty, image_url, is_active) VALUES
('SHUTTLE-RSL-SILVER-TUBE', 'ลูกแบด RSL Silver (ยกหลอด 12 ลูก)', 'SHUTTLECOCK', 850.00, 25, 'https://images.unsplash.com/photo-1626224583764-f87db24ac4ea?w=400', TRUE),
('SHUTTLE-RSL-SILVER-PIECE', 'ลูกแบด RSL Silver (แบ่งขาย 1 ลูก)', 'SHUTTLECOCK', 75.00, 60, 'https://images.unsplash.com/photo-1626224583764-f87db24ac4ea?w=400', TRUE),
('SHUTTLE-VICTOR-MASTER1', 'ลูกแบด Victor Master No.1 (ยกหลอด 12 ลูก)', 'SHUTTLECOCK', 990.00, 15, 'https://images.unsplash.com/photo-1613918108466-292b78a8ef95?w=400', TRUE),
('DRINK-SINGHA-WATER', 'น้ำดื่มสิงห์ 600 มล. (แช่เย็น)', 'DRINK', 15.00, 120, 'https://images.unsplash.com/photo-1548839140-29a749e1bc4e?w=400', TRUE),
('DRINK-POCARI-SWEAT', 'เครื่องดื่มเกลือแร่ Pocari Sweat 500 มล.', 'DRINK', 30.00, 80, 'https://images.unsplash.com/photo-1581009146145-b5ef050c2e1e?w=400', TRUE),
('DRINK-MINERE-WATER', 'น้ำแร่ธรรมชาติ Minere 500 มล.', 'DRINK', 20.00, 50, 'https://images.unsplash.com/photo-1559839914-17aae19cec71?w=400', TRUE),
('RENTAL-RACKET-ASTROX', 'บริการเช่าไม้แบด Yonex Astrox 88D Play (ต่อแมตช์)', 'RENTAL', 80.00, 10, 'https://images.unsplash.com/photo-1521537634581-0dced2fee2ef?w=400', TRUE),
('RENTAL-GRIP-YONEX', 'กริปพันด้าม Yonex Super Grap (AC102EX พร้อมบริการพัน)', 'RENTAL', 60.00, 40, 'https://images.unsplash.com/photo-1584813470613-5b1c1cad3d69?w=400', TRUE);

-- 3. Seed Users
INSERT INTO users (id, line_user_id, display_name, picture_url, phone, role) VALUES
('a0000000-0000-0000-0000-000000000001', 'U_ADMIN_SMASHHUB', 'สมชาย ผู้จัดการสนาม (Admin)', 'https://api.dicebear.com/7.x/bottts/svg?seed=admin_somchai', '081-999-8888', 'ADMIN'),
('b0000000-0000-0000-0000-000000000001', 'U_CUSTOMER_KAN', 'กานต์ ธนบดี', 'https://api.dicebear.com/7.x/avataaars/svg?seed=kan', '089-111-2233', 'CUSTOMER'),
('b0000000-0000-0000-0000-000000000002', 'U_CUSTOMER_NUT', 'ณัฐพงษ์ วงศ์วิวัฒน์', 'https://api.dicebear.com/7.x/avataaars/svg?seed=nuttapong', '086-444-5566', 'CUSTOMER'),
('b0000000-0000-0000-0000-000000000003', 'U_CUSTOMER_TNP', 'ธนาภา สินสมบูรณ์', 'https://api.dicebear.com/7.x/avataaars/svg?seed=thanapa', '082-777-8899', 'CUSTOMER'),
('b0000000-0000-0000-0000-000000000004', 'U_CUSTOMER_WRP', 'วรพจน์ ศักดานรินทร์', 'https://api.dicebear.com/7.x/avataaars/svg?seed=worapoj', '095-222-3344', 'CUSTOMER');

-- 4. Seed Today's Bookings (All 4 Statuses)
-- Using fixed UUIDs for deterministic foreign key relationships

-- 4.1 Booking 1: CONFIRMED
-- Customer 1 (กานต์) books Court 1: 18:00 - 20:00 (Peak rate 220 x 2 = 440)
-- Items: 1x RSL Silver Tube (850) + 2x Pocari Sweat (60) = 910
-- Total: 440 + 910 = 1,350 THB
-- Status: CONFIRMED, Verified Slip attached, is_slot_locked = TRUE
INSERT INTO bookings (
    id, booking_no, user_id, status, court_total, item_total, net_total, check_in_status, expires_at, created_at, updated_at
) VALUES (
    'c0000000-0000-0000-0000-000000000001',
    TO_CHAR(CURRENT_DATE, '"BK-"YYYYMMDD"-001"'),
    'b0000000-0000-0000-0000-000000000001',
    'CONFIRMED',
    440.00,
    910.00,
    1350.00,
    'WAITING',
    NOW() + INTERVAL '10 minutes',
    NOW() - INTERVAL '1 hour',
    NOW() - INTERVAL '50 minutes'
);

-- Court slots for Booking 1 (18:00 and 19:00)
INSERT INTO booking_courts (id, booking_id, court_id, play_date, start_hour, end_hour, hourly_rate, is_slot_locked) VALUES
('d0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000001', 1, CURRENT_DATE, 18, 19, 220.00, TRUE),
('d0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000001', 1, CURRENT_DATE, 19, 20, 220.00, TRUE);

-- Items for Booking 1
INSERT INTO booking_items (id, booking_id, product_id, qty, unit_price, subtotal) VALUES
('e0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000001', 1, 1, 850.00, 850.00), -- RSL Silver Tube
('e0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000001', 5, 2, 30.00, 60.00);    -- Pocari Sweat

-- Slip for Booking 1
INSERT INTO payment_slips (
    id, booking_id, slip_image_url, trans_ref, sender_name, receiver_account, amount, verification_status, raw_payload, verified_at, created_at
) VALUES (
    'f0000000-0000-0000-0000-000000000001',
    'c0000000-0000-0000-0000-000000000001',
    'https://placehold.co/600x800/004d40/ffffff?text=Slip+Verified+1350THB',
    'KBANK-20261024-9876543210',
    'นายกานต์ ธนบดี',
    '098-7-65432-1 (บจก. สแมชฮับ แบดมินตัน)',
    1350.00,
    'VERIFIED',
    '{"bank": "KBANK", "fee": 0, "verified_by": "OpenSlipVerify", "channel": "K+ Mobile"}'::jsonb,
    NOW() - INTERVAL '50 minutes',
    NOW() - INTERVAL '52 minutes'
);


-- 4.2 Booking 2: PENDING_PAYMENT
-- Customer 2 (ณัฐพงษ์) books Court 2: 19:00 - 21:00 (Peak rate 220 x 2 = 440)
-- Items: 2x Minere Water (40)
-- Total: 440 + 40 = 480 THB
-- Status: PENDING_PAYMENT (holding lock for 10 mins, expires in 8 minutes)
INSERT INTO bookings (
    id, booking_no, user_id, status, court_total, item_total, net_total, check_in_status, expires_at, created_at, updated_at
) VALUES (
    'c0000000-0000-0000-0000-000000000002',
    TO_CHAR(CURRENT_DATE, '"BK-"YYYYMMDD"-002"'),
    'b0000000-0000-0000-0000-000000000002',
    'PENDING_PAYMENT',
    440.00,
    40.00,
    480.00,
    'WAITING',
    NOW() + INTERVAL '8 minutes',
    NOW() - INTERVAL '2 minutes',
    NOW() - INTERVAL '2 minutes'
);

-- Court slots for Booking 2 (19:00 and 20:00)
INSERT INTO booking_courts (id, booking_id, court_id, play_date, start_hour, end_hour, hourly_rate, is_slot_locked) VALUES
('d0000000-0000-0000-0000-000000000003', 'c0000000-0000-0000-0000-000000000002', 2, CURRENT_DATE, 19, 20, 220.00, TRUE),
('d0000000-0000-0000-0000-000000000004', 'c0000000-0000-0000-0000-000000000002', 2, CURRENT_DATE, 20, 21, 220.00, TRUE);

-- Items for Booking 2
INSERT INTO booking_items (id, booking_id, product_id, qty, unit_price, subtotal) VALUES
('e0000000-0000-0000-0000-000000000003', 'c0000000-0000-0000-0000-000000000002', 6, 2, 20.00, 40.00); -- Minere Water


-- 4.3 Booking 3: MANUAL_REVIEW
-- Customer 3 (ธนาภา) books Court 5 (VIP Acrylic): 17:00 - 19:00 (Peak rate 260 x 2 = 520)
-- Items: 1x Rental Yonex Astrox Racket (80)
-- Total: 520 + 80 = 600 THB
-- Status: MANUAL_REVIEW (slip verification failed with gateway timeout)
INSERT INTO bookings (
    id, booking_no, user_id, status, court_total, item_total, net_total, check_in_status, expires_at, created_at, updated_at
) VALUES (
    'c0000000-0000-0000-0000-000000000003',
    TO_CHAR(CURRENT_DATE, '"BK-"YYYYMMDD"-003"'),
    'b0000000-0000-0000-0000-000000000003',
    'MANUAL_REVIEW',
    520.00,
    80.00,
    600.00,
    'WAITING',
    NOW() + INTERVAL '30 minutes',
    NOW() - INTERVAL '15 minutes',
    NOW() - INTERVAL '10 minutes'
);

-- Court slots for Booking 3 (17:00 and 18:00)
INSERT INTO booking_courts (id, booking_id, court_id, play_date, start_hour, end_hour, hourly_rate, is_slot_locked) VALUES
('d0000000-0000-0000-0000-000000000005', 'c0000000-0000-0000-0000-000000000003', 5, CURRENT_DATE, 17, 18, 260.00, TRUE),
('d0000000-0000-0000-0000-000000000006', 'c0000000-0000-0000-0000-000000000003', 5, CURRENT_DATE, 18, 19, 260.00, TRUE);

-- Items for Booking 3
INSERT INTO booking_items (id, booking_id, product_id, qty, unit_price, subtotal) VALUES
('e0000000-0000-0000-0000-000000000004', 'c0000000-0000-0000-0000-000000000003', 7, 1, 80.00, 80.00); -- Yonex Racket Rental

-- Slip for Booking 3 (Slip API timeout -> waiting for Admin approval)
INSERT INTO payment_slips (
    id, booking_id, slip_image_url, trans_ref, sender_name, receiver_account, amount, verification_status, raw_payload, verified_at, created_at
) VALUES (
    'f0000000-0000-0000-0000-000000000002',
    'c0000000-0000-0000-0000-000000000003',
    'https://placehold.co/600x800/d97706/ffffff?text=Slip+Pending+Admin+Review',
    'SCB-20261024-1122334455',
    'น.ส. ธนาภา สินสมบูรณ์',
    '098-7-65432-1 (บจก. สแมชฮับ แบดมินตัน)',
    600.00,
    'API_ERROR_PENDING_ADMIN',
    '{"error": "OpenSlipVerify Gateway Timeout (504)", "retry_count": 3, "last_attempt": "2026-10-24T03:10:00Z"}'::jsonb,
    NULL,
    NOW() - INTERVAL '10 minutes'
);


-- 4.4 Booking 4: EXPIRED
-- Customer 4 (วรพจน์) booked Court 3: 14:00 - 16:00 (Normal rate 150 x 2 = 300)
-- Customer did not upload slip within 10 minutes -> Expired!
-- CRITICAL RULE: is_slot_locked = FALSE so other customers can book Court 3 at 14:00 and 15:00
INSERT INTO bookings (
    id, booking_no, user_id, status, court_total, item_total, net_total, check_in_status, expires_at, created_at, updated_at
) VALUES (
    'c0000000-0000-0000-0000-000000000004',
    TO_CHAR(CURRENT_DATE, '"BK-"YYYYMMDD"-004"'),
    'b0000000-0000-0000-0000-000000000004',
    'EXPIRED',
    300.00,
    0.00,
    300.00,
    'WAITING',
    NOW() - INTERVAL '20 minutes',
    NOW() - INTERVAL '40 minutes',
    NOW() - INTERVAL '20 minutes'
);

-- Court slots for Booking 4 (Notice is_slot_locked = FALSE!)
INSERT INTO booking_courts (id, booking_id, court_id, play_date, start_hour, end_hour, hourly_rate, is_slot_locked) VALUES
('d0000000-0000-0000-0000-000000000007', 'c0000000-0000-0000-0000-000000000004', 3, CURRENT_DATE, 14, 15, 150.00, FALSE),
('d0000000-0000-0000-0000-000000000008', 'c0000000-0000-0000-0000-000000000004', 3, CURRENT_DATE, 15, 16, 150.00, FALSE);
