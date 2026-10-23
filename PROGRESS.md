# SmashHub - Project Progress & Changelog

## 📌 Project Overview
SmashHub is a Badminton Court Booking System via LINE LIFF, integrated with a mini-store, automated slip verification (with demo mock mode), and real-time POS backoffice.

---

## 🏗 Directory Structure Initialized
- `docker-compose.yml` : PostgreSQL 16 container config with auto-migrations
- `backend/`
  - `cmd/server/main.go` : Application entry point with routing, CORS, and graceful shutdown
  - `internal/models/` : GORM schema models, Enums, DB initialization & connection pooling (`models.go`, `db.go`)
  - `internal/services/` : Core business logic, concurrency control, slip verifier, cron worker, SSE Hub (`booking_service.go`, `slip_service.go`, `admin_service.go`, `cron_worker.go`, `sse.go`, `promptpay.go`)
  - `internal/handlers/` : HTTP REST & SSE handlers (`court_handler.go`, `booking_handler.go`, `admin_handler.go`, `sse_handler.go`)
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
  - Created `001_create_schema.sql` (ENUMs, 7 tables, partial unique index).
  - Created `002_seed_data.sql` (6 courts, 8 mini-store products, 5 users, 4 simulated bookings).
  - Created Go GORM Models in `backend/internal/models/models.go` and connection helper `backend/internal/models/db.go`.
- [x] **Phase 2: Backend Core API & Business Logic (Go + Gin + GORM)** (Completed)
  - **Endpoints Implemented:**
    - `GET /health`: Healthcheck endpoint
    - `GET /api/v1/courts/availability?date=YYYY-MM-DD`: Returns availability matrix (10:00-22:00) with status `'AVAILABLE'`, `'LOCKED_PENDING'`, `'BOOKED'`, `'MAINTENANCE'` and dynamic normal/peak rates.
    - `GET /api/v1/products`: Returns mini-store items with live `stock_qty`.
    - `POST /api/v1/bookings`: Single atomic transaction locking product inventory (`SELECT ... FOR UPDATE`), server-side calculation of `court_total`, `item_total`, and `net_total`, locking court slots in `booking_courts`, and setting 10-minute expiration. Returns `HTTP 409 Conflict` with `"คอร์ทที่เลือกเพิ่งถูกจองไปเมื่อสักครู่"` on double booking collision.
    - `GET /api/v1/bookings/:id`: Returns booking details and PromptPay QR payload for customer LIFF view.
- [x] **Phase 3: Automated Slip Verification & Cron Job Worker** (Completed)
  - **Slip Verification (`POST /api/v1/bookings/:id/verify-slip`):**
    - 4-Step Verification rules:
      1. Duplicate `trans_ref` check in `payment_slips` -> `REJECTED_DUPLICATE` (HTTP 409)
      2. Receiver account matching check -> `REJECTED_WRONG_ACCOUNT`
      3. Amount matching check (`amount >= net_total`) -> `REJECTED_AMOUNT_MISMATCH` (HTTP 402)
      4. External API gateway check: if API down -> booking becomes `MANUAL_REVIEW`, slip becomes `API_ERROR_PENDING_ADMIN`
    - Supported Mock Scenarios: `MOCK_VALID`, `MOCK_DUPLICATE`, `MOCK_UNDERPAID`, `MOCK_API_DOWN`.
  - **Background Cron Worker:**
    - Goroutine running every 30 seconds.
    - Automatically finds `PENDING_PAYMENT` bookings where `expires_at < NOW()`.
    - Runs atomic transaction: updates status to `EXPIRED`, unlocks court slots (`is_slot_locked = FALSE`), restores inventory stock for each item, and broadcasts SSE event.
- [x] **Phase 4: Real-time SSE Sync & Admin Backoffice APIs** (Completed)
  - `GET /api/v1/events/stream`: Server-Sent Events hub broadcasting real-time events (`BOOKING_CREATED`, `BOOKING_CONFIRMED`, `BOOKING_MANUAL_REVIEW`, `BOOKING_EXPIRED`, `BOOKING_CANCELLED`, `ITEMS_ADDED`, `CHECKIN_UPDATED`).
  - `GET /api/v1/admin/bookings`: Filter by `date`, `status`, `court_id`, and `search` (customer name/phone/booking_no).
  - `PATCH /api/v1/admin/bookings/:id/manual-approve`: Admin approval/rejection of `MANUAL_REVIEW` bookings.
  - `POST /api/v1/admin/bookings/:id/add-items`: POS counter feature to add drinks/shuttlecocks with `SELECT ... FOR UPDATE` stock deduction and bill recalculation.
  - `PATCH /api/v1/admin/bookings/:id/check-in`: Updates check-in status (`WAITING`, `CHECKED_IN`, `COMPLETED`).
- [ ] **Phase 5: Frontend Customer LIFF Flow**
- [ ] **Phase 6: Frontend Admin Backoffice & POS**
- [ ] **Phase 7: Split-screen Interactive Demo & Final Polish**
