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
  - `src/views/liff/` : Mobile customer booking & store interface (max-w 430px)
  - `src/views/admin/` : Admin dashboard & POS interface
  - `src/views/demo/` : Interactive split-screen demo showcasing live dual sync
  - `src/stores/` : Pinia state management stores (`lineProfile.js`, `bookingStore.js`, `adminStore.js`)

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
  - `GET /api/v1/events/stream`: Server-Sent Events hub broadcasting real-time events.
  - `GET /api/v1/admin/bookings`: Multi-criteria filter (`date`, `status`, `court_id`, `search`).
  - `PATCH /api/v1/admin/bookings/:id/manual-approve`: Admin approval/rejection.
  - `POST /api/v1/admin/bookings/:id/add-items`: POS counter feature with `SELECT ... FOR UPDATE`.
  - `PATCH /api/v1/admin/bookings/:id/check-in`: Check-in status updater (`WAITING`, `CHECKED_IN`, `COMPLETED`).
- [x] **Phase 5: Frontend Customer LINE LIFF Flow** (Completed)
  - Mobile viewport (max-w 430px) matching LINE in-app browser design.
  - Mock LINE Profile Switcher with 3 simulated customer accounts.
  - **Step 1 (Court Matrix)**: Quick date chips, Courts 1-6 X-axis vs 10:00-22:00 Y-axis, auto Peak-rate calculation, color-coded slot status.
  - **Step 2 (Mini Store)**: Stock badges (`เหลือ X ชิ้น`), stepper counter for shuttlecocks, drinks, and racket rentals.
  - **Step 3 (Payment & Slip)**: 10:00-minute live countdown timer, dynamic PromptPay QR code, file uploader, and 4 portfolio mock scenario buttons (`MOCK_VALID`, `MOCK_DUPLICATE`, `MOCK_UNDERPAID`, `MOCK_API_DOWN`).
  - **Step 4 (E-Ticket)**: Boarding-pass style ticket with gate check-in QR code and confetti celebration!
- [x] **Phase 6: Frontend Admin Backoffice & POS** (Completed)
  - **Tab 1: Live Court Board & POS**: Interactive timeline grid of Courts 1-6, quick check-in buttons, POS add-item buttons, and click empty slot to open Walk-in booking.
  - **Tab 2: Order & Slip Audit Log**: Element Plus Table with date, status, court, and search filters; colored badges for slip verification status.
  - **Element Plus Drawer**: Displays slip image preview side-by-side with extracted JSON payload (`trans_ref`, sender name, amount comparison) and Approve/Reject buttons for `MANUAL_REVIEW`.
  - **Receipt Modal**: Printable tax invoice / receipt with company header, VAT breakdown, and direct print/PDF support.
- [x] **Phase 7: Split-screen Interactive Demo & Portfolio Showcase** (Completed)
  - `/demo` view rendering LINE LIFF mobile frame on the left and Admin Dashboard on the right.
  - Real-time synchronization via SSE: bookings created or verified on mobile update instantly on the Admin screen without page refresh.
  - Verified live in browser with `browser_subagent`.
- [x] **Phase 8: Comprehensive UX/UI Overhaul & Elimination of AI Design Tropes** (Completed)
  - **Eliminated AI Hallmark Tropes**:
    - Removed cluttered candy-colored button rows and emoji spam from the top sandbox bar.
    - Eliminated unreadable 6-column spreadsheet squished into mobile viewport.
    - Polished typography using Google Fonts: `Plus Jakarta Sans`, `Outfit`, and `JetBrains Mono`.
  - **Benchmarked from World-Class Products**:
    - **Playtomic / MATCHi (#1 Sports Booking Apps)**:
      - Implemented Dual View Court Booking on mobile:
        1. **View A: Time-First Slots View (Recommended)**: Grouped by hourly blocks (`10:00 - 11:00`) with touch-friendly court chips showing court type (`PVC` vs `VIP Acrylic`), pricing, and real-time status.
        2. **View B: Wide Court Grid**: Spacious horizontal scrolling board with realistic badminton floor graphics.
      - Court type filters (`ทั้งหมด`, `PVC (1-4)`, `VIP (5-6)`) and time period indicators (`Peak Hour หลัง 17:00 น.`).
    - **Linear / Stripe Sandbox (DevTools & Backoffice)**:
      - Dark graphite Developer Sandbox bar with refined scenario selectors, status LED indicators, and collapsible architecture inspection banner.
      - 1-Click Guided Auto-Tour (`⚡ 1-Click Guided Tour`) automatically showcasing full booking and slip verification flow in 3 seconds.
      - Clean KPI metric cards with micro-caps labels and monospace numbers.
  - Pushed to GitHub repository `https://github.com/thanawat-c/SmashHub.git`.
- [x] **Phase 9: Comprehensive Template Redesign & Task-by-Task Overhaul** (Completed)
  - **Task 1 (`c6f951f`)**: Clean White Foundation, Authentic LINE Header & Phone Frame.
  - **Task 2 (`fba4912`)**: Revamped Thai QR Payment Card (`#003B64` standard blue banner), clean slip upload dropzone, removed tacky "ตรวจผ่าน AI อัตโนมัติ", and collapsed mock tester drawer.
  - **Task 3 (`549fb8c`)**: Unified clean booking cards, removed neon glows (`shadow-glow-emerald`), and polished floating checkout dock.
  - **Task 4 (`ae9de1e`)**: Polished Admin Dashboard & POS styling with clean `#F8FAFC` studio theme and modernized E-Ticket to a clean, crisp light-mode boarding pass.
  - **Task 5 (`46a63fb`)**: Eliminated broken nested phone frames and unified clean mobile viewport across `/liff` and `/demo`.
  - **Task 6 (`656250c`)**: Harmonized background colors and refined mini-store step styling.
  - **Task 8 (`5e76f66`)**: Fixed Receipt PDF/Print blank page and overflow issue by isolating print execution inside a clean hidden iframe with exact A4 page sizing and Tailwind styles.
  - **Task 9 (`4bd943e`)**: Enabled interactive back navigation by connecting the `<` button on LINE LIFF header and making step indicator chips (1, 2, 3) clickable to freely jump back to earlier steps or reset to Step 1.
