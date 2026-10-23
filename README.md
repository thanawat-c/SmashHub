# 🏸 SmashHub - Badminton Court Booking & Mini Store POS

ระบบจองสนามแบดมินตันผ่าน LINE LIFF พร้อมมินิสโตร์และระบบตรวจสลิปอัตโนมัติ (Production-Grade Portfolio Project)

## 🌟 จุดเด่นของระบบ (Key Features)
1. **Prevent Double Booking:** ป้องกันการจองคอร์ทชนกันระดับ Database ด้วย Partial Unique Index (`is_slot_locked = TRUE`)
2. **Atomic Booking & Stock Deduction:** ทำธุรกรรมจองสนามและสั่งสินค้ามินิสโตร์พร้อมกันใน 1 Database Transaction ด้วย Row-level Locking (`SELECT ... FOR UPDATE`)
3. **4-Step Automated Slip Verification:** ตรวจสอบสลิปโอนเงินอัตโนมัติพร้อมรองรับ **Portfolio Mock Mode** สำหรับทดสอบ Edge Cases
4. **Real-time Synchronization:** ซิงค์ข้อมูลสถานะการจองไปยังหน้า POS Backoffice และ Demo ทันทีผ่าน Server-Sent Events (SSE)

## 🛠 Tech Stack
- **Backend:** Go (Gin Framework + GORM)
- **Database:** PostgreSQL 16
- **Frontend:** Vue 3 (Composition API) + Vite + Pinia + Tailwind CSS + Element Plus
- **Infrastructure:** Docker Compose

## 🚀 Quick Start (Local Development)

### 1. รัน Database (PostgreSQL 16)
```bash
docker compose up -d
```
*(Migrations และ Seed Data จะถูกรันอัตโนมัติผ่าน `/backend/migrations`)*

### 2. รัน Backend
```bash
cd backend
go run cmd/server/main.go
```