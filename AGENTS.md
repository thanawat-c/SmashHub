# AGENTS.md - Project Rules & Architecture Guide

## 1. Project Overview
**Project Name:** SmashHub (ระบบจองสนามแบดมินตันผ่าน LINE LIFF พร้อมมินิสโตร์และระบบตรวจสลิปอัตโนมัติ)
**Goal:** Production-grade portfolio project demonstrating atomic transactions, concurrency control, automated payment slip verification, and real-time POS backoffice management.

## 2. Tech Stack & Directory Structure
Use a clean Monorepo structure:
- `/docker-compose.yml` : PostgreSQL 16 service
- `/backend` : Go (Golang) + Gin Framework + GORM + PostgreSQL
  - `/cmd/server/main.go` : Entry point
  - `/internal/models` : Database structs & enums
  - `/internal/handlers` : HTTP REST & SSE handlers
  - `/internal/services` : Business logic, DB transactions, slip verifier, cron worker
  - `/migrations` : SQL migration & seed files
- `/frontend` : Vue 3 (Composition API `<script setup>`) + Vite + Pinia + Vue Router + Tailwind CSS + Element Plus
  - `/src/views/liff` : Mobile-first customer storefront (`/liff`)
  - `/src/views/admin` : Desktop store backoffice & POS (`/admin`)
  - `/src/views/demo` : Split-screen interactive portfolio showcase (`/demo`)
  - `/src/stores` : Pinia state stores

## 3. Non-Negotiable Engineering Rules
1. **Prevent Double Booking (Database Level):**
   - Table `booking_courts` MUST enforce a Partial Unique Index:
     `CREATE UNIQUE INDEX idx_unique_active_court_slot ON booking_courts (court_id, play_date, start_hour) WHERE is_slot_locked = TRUE;`
   - When a booking becomes `EXPIRED` or `CANCELLED`, update `is_slot_locked = FALSE` so the court slot can be booked again.
2. **Atomic Booking & Stock Deduction:**
   - Creating a booking (`POST /api/v1/bookings`) MUST run inside a single database transaction (`db.Transaction`).
   - Lock product rows using `SELECT ... FOR UPDATE`, verify `stock_qty >= requested_qty`, deduct stock, and lock court slots atomically.
   - Set `expires_at = NOW() + 10 minutes` for `PENDING_PAYMENT` bookings.
3. **4-Step Slip Verification & Portfolio Mock Mode:**
   - Endpoint `POST /api/v1/bookings/:id/verify-slip` must support `mock_scenario` payload for recruiter testing:
     - `MOCK_VALID`: Passes all checks -> booking status becomes `CONFIRMED`.
     - `MOCK_DUPLICATE`: Fails step 1 (duplicate `trans_ref` in `payment_slips`) -> `REJECTED_DUPLICATE`.
     - `MOCK_UNDERPAID`: Fails step 3 (`amount < net_total`) -> `REJECTED_AMOUNT_MISMATCH`.
     - `MOCK_API_DOWN`: External verifier timeout -> booking status becomes `MANUAL_REVIEW` and slip becomes `API_ERROR_PENDING_ADMIN`.
4. **Server-Side Calculation Only:**
   - Never trust price totals sent from Frontend. Always recalculate `court_total`, `item_total`, and `net_total` on the Backend using database prices.
5. **Real-Time Sync:**
   - Broadcast booking/slip state changes via Server-Sent Events (`GET /api/v1/events/stream`) so the `/admin` and `/demo` views update immediately without page refresh.

## 4. Agent Workflow Instructions (IMPORTANT)
- **Work Phase-by-Phase:** Only implement the phase explicitly requested by the user. Do NOT jump ahead to future phases.
- **Verify via Terminal:** After writing code in any phase, run the build/tests or `curl` commands in the terminal to verify there are no syntax or runtime errors before finishing the turn.
- **Maintain `PROGRESS.md`:** At the end of each phase, update `/PROGRESS.md` with completed features, database tables, and exact API endpoint paths/payloads so the next phase stays 100% aligned.
