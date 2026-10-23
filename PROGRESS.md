# SmashHub - Project Progress & Changelog

## 📌 Project Overview
SmashHub is a Badminton Court Booking System via LINE LIFF, integrated with a mini-store, automated slip verification (with demo mock mode), and real-time POS backoffice.

---

## 🏗 Directory Structure Initialized
- `docker-compose.yml` : PostgreSQL 16 container config with auto-migrations
- `backend/`
  - `cmd/server/` : Application entry point
  - `internal/models/` : GORM schema models, Enums, DB initialization & connection pooling (`models.go`, `db.go`)
  - `internal/handlers/` : HTTP REST & SSE handlers
  - `internal/services/` : Business logic, atomic transactions, slip verifier
  - `migrations/` : SQL migrations & seed data (`001_create_schema.sql`, `002_seed_data.sql`)
- `frontend/`
  - `src/views/liff/` : Mobile customer booking & store interface
  - `src/views/admin/` : Admin dashboard & POS interface
  - `src/views/demo/` : Interactive split-screen demo
  - `src/stores/` : Pinia state management stores

---

## 🚀 Phases Breakdown
- [x] **Phase 0: Workspace & Folder Setup** (Completed)
  - Created monorepo directory hierarchy
  - Initialized `AGENTS.md` and `PROGRESS.md`
  - Prepared `docker-compose.yml` for PostgreSQL 16
- [x] **Phase 1: Database Schema, Docker Setup & Seed Data** (Completed)
  - Configured `docker-compose.yml` with PostgreSQL 16-alpine, Asia/Bangkok timezone, and healthchecks.
  - Created `001_create_schema.sql`:
    - **ENUMs**: `booking_status` ('PENDING_PAYMENT', 'CONFIRMED', 'MANUAL_REVIEW', 'EXPIRED', 'CANCELLED'), `slip_status` ('VERIFIED', 'REJECTED_DUPLICATE', 'REJECTED_AMOUNT_MISMATCH', 'REJECTED_WRONG_ACCOUNT', 'API_ERROR_PENDING_ADMIN').
    - **7 Tables**: `users`, `courts`, `products`, `bookings`, `booking_courts`, `booking_items`, `payment_slips`.
    - **Constraints**: Partial Unique Index `idx_unique_active_court_slot` on `booking_courts(court_id, play_date, start_hour) WHERE is_slot_locked = TRUE` to prevent double booking while allowing re-booking when a slot is unlocked (`is_slot_locked = FALSE`).
  - Created `002_seed_data.sql`:
    - 6 Badminton Courts (PVC Courts 1-4 with 150/220 rates, Acrylic VIP Courts 5-6 with 180/260 rates).
    - 8 Mini-Store Products with stock levels (RSL Silver tubes/pieces, Victor Master No.1, Singha water, Pocari Sweat, Minere, Yonex rental racket, Yonex grip).
    - 5 Users (Admin manager + 4 customers).
    - 4 Today's Bookings covering all states: `CONFIRMED` (with verified slip), `PENDING_PAYMENT` (holding 10m lock), `MANUAL_REVIEW` (gateway timeout slip), and `EXPIRED` (`is_slot_locked = FALSE`).
  - Created Go GORM Models in `backend/internal/models/models.go` and connection helper `backend/internal/models/db.go`.
  - Verified migrations and constraints in terminal against PostgreSQL.
- [ ] **Phase 2: Backend Core API & Business Logic (Go + Gin + GORM)**
- [ ] **Phase 3: Automated Slip Verification & Cron Job Worker**
- [ ] **Phase 4: Real-time SSE Sync**
- [ ] **Phase 5: Frontend Customer LIFF Flow**
- [ ] **Phase 6: Frontend Admin Backoffice & POS**
- [ ] **Phase 7: Split-screen Interactive Demo & Final Polish**
