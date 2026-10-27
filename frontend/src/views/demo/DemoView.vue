<script setup>
import { ref } from 'vue'
import LiffView from '../liff/LiffView.vue'
import AdminDashboard from '../admin/AdminDashboard.vue'
import { useBookingStore } from '@/stores/bookingStore'
import { useLineProfileStore } from '@/stores/lineProfile'
import { useAdminStore } from '@/stores/adminStore'
import {
  Smartphone, Monitor, CheckCircle, XCircle, AlertTriangle,
  Server, Cpu, Play, RefreshCw, Zap, ShieldCheck, Terminal
} from 'lucide-vue-next'
import { ElNotification } from 'element-plus'

const bookingStore = useBookingStore()
const profileStore = useLineProfileStore()
const adminStore = useAdminStore()

const isTesting = ref(false)
const isAutoTourRunning = ref(false)
const latestArchitectureNote = ref(null)

const scenarios = [
  {
    key: 'MOCK_VALID',
    label: 'Valid Slip',
    subLabel: 'Amount Matches + Instant Confirmation',
    btnClass: 'bg-emerald-600 hover:bg-emerald-500 text-white border-emerald-400/40 shadow-glow-emerald',
    icon: CheckCircle,
    color: 'emerald',
    badge: 'MOCK_VALID',
    engineeringExplanation: {
      title: '4-Step Verification: All Checks Passed (100%)',
      steps: [
        '1. Check trans_ref: Unique reference, no duplicate in payment_slips',
        '2. Receiver account: Matches SmashHub Badminton Co., Ltd. (098-7-65432-1)',
        '3. Amount verification: Transferred amount >= booking.net_total',
        '4. External Gateway: HTTP 200 OK verification'
      ],
      result: 'DB: booking.status = CONFIRMED, permanent slot lock, digital pass QR issued, instant SSE broadcast to Admin POS.'
    }
  },
  {
    key: 'MOCK_DUPLICATE',
    label: 'Duplicate Slip',
    subLabel: 'Reused Transaction Ref',
    btnClass: 'bg-rose-600 hover:bg-rose-500 text-white border-rose-400/40',
    icon: XCircle,
    color: 'rose',
    badge: 'MOCK_DUPLICATE',
    engineeringExplanation: {
      title: 'Step 1 Failed: Reused Slip Detected',
      steps: [
        '1. Check trans_ref: Key collision found in payment_slips table',
        '2. Atomic Transaction Rollback to prevent double spending',
        '3. Slip recorded: verification_status = REJECTED_DUPLICATE',
        '4. Return HTTP 409 Conflict to customer LIFF client'
      ],
      result: 'Booking remains in PENDING_PAYMENT with existing 10-minute expiry timer.'
    }
  },
  {
    key: 'MOCK_UNDERPAID',
    label: 'Underpaid Slip',
    subLabel: 'Amount < Net Total',
    btnClass: 'bg-amber-600 hover:bg-amber-500 text-white border-amber-400/40',
    icon: AlertTriangle,
    color: 'amber',
    badge: 'MOCK_UNDERPAID',
    engineeringExplanation: {
      title: 'Step 3 Failed: Transferred Amount Mismatch',
      steps: [
        '1. Check trans_ref: Valid (new transaction)',
        '2. Receiver account: Valid receiver',
        '3. Amount check: Transferred amount < booking.net_total',
        '4. Slip recorded: verification_status = REJECTED_AMOUNT_MISMATCH'
      ],
      result: 'Returns HTTP 402 Payment Required indicating balance due and prompts for valid slip.'
    }
  },
  {
    key: 'MOCK_API_DOWN',
    label: 'Bank Timeout',
    subLabel: 'Fallback to Manual Audit',
    btnClass: 'bg-purple-600 hover:bg-purple-500 text-white border-purple-400/40',
    icon: Server,
    color: 'purple',
    badge: 'MOCK_API_DOWN',
    engineeringExplanation: {
      title: 'Circuit Breaker: External Gateway Timeout (504)',
      steps: [
        '1. External banking API timeout (HTTP 504 Gateway Timeout)',
        '2. Graceful Fallback: Enqueued into Manual Audit rather than failing customer',
        '3. Update status: booking.status = MANUAL_REVIEW',
        '4. Record slip: verification_status = API_ERROR_PENDING_ADMIN'
      ],
      result: 'Broadcast SSE: Flags pending audit in Admin Dashboard with 1-click Approve/Reject controls.'
    }
  }
]

async function runScenario(scenario) {
  isTesting.value = true
  try {
    await bookingStore.ensureActiveBookingForDemo(profileStore.currentProfile)
    await bookingStore.verifySlip(scenario.key)
    latestArchitectureNote.value = scenario.engineeringExplanation

    ElNotification({
      title: scenario.badge + ' - ' + scenario.label,
      dangerouslyUseHTMLString: true,
      message: `
        <div style="font-size: 11px; line-height: 1.5; color: #1e293b;">
          <p style="font-weight: 700; margin-bottom: 4px; color: #0f172a;">${scenario.engineeringExplanation.title}</p>
          <ul style="padding-left: 14px; margin: 0; list-style-type: disc;">
            ${scenario.engineeringExplanation.steps.map(s => `<li>${s}</li>`).join('')}
          </ul>
          <p style="margin-top: 6px; font-weight: 600; color: #059669;">Architecture: ${scenario.engineeringExplanation.result}</p>
        </div>
      `,
      duration: 6000,
      type: scenario.color === 'emerald' ? 'success' : scenario.color === 'rose' ? 'error' : 'warning'
    })

    adminStore.fetchBookings()
  } catch (err) {
    console.error('Scenario error:', err)
  } finally {
    isTesting.value = false
  }
}

// 1-Click Guided Auto-Tour for Recruiters
async function runAutoTour() {
  if (isAutoTourRunning.value) return
  isAutoTourRunning.value = true

  ElNotification({
    title: 'Starting Guided Auto-Tour',
    message: 'Simulating user flow: court selection, booking creation, and slip verification...',
    type: 'info',
    duration: 3000
  })

  // 1. Reset
  bookingStore.resetSelection()
  await new Promise(r => setTimeout(r, 600))

  // 2. Select court & add item & create booking
  await bookingStore.ensureActiveBookingForDemo(profileStore.currentProfile)
  await new Promise(r => setTimeout(r, 1000))

  // 3. Trigger MOCK_VALID verification
  await runScenario(scenarios[0])
  isAutoTourRunning.value = false
}

function resetDemo() {
  bookingStore.resetSelection()
  bookingStore.fetchAvailability()
  latestArchitectureNote.value = null
  ElNotification({
    title: 'Demo State Reset',
    message: 'Ready for a new simulation session.',
    type: 'info',
    duration: 2000
  })
}
</script>

<template>
  <div class="flex-1 flex flex-col bg-slate-950 text-white h-[calc(100vh-45px)] overflow-hidden select-none">
    <!-- Top Executive Interactive Bar for Recruiters (Glassmorphism) -->
    <div class="bg-slate-900/95 backdrop-blur-md border-b border-slate-800/80 px-4 py-2.5 shadow-xl z-40 flex-shrink-0">
      <div class="flex flex-wrap items-center justify-between gap-2.5">
        <!-- Title & Badges -->
        <div class="flex items-center space-x-2.5">
          <div class="p-1.5 bg-emerald-500/10 border border-emerald-500/30 rounded-xl text-emerald-400">
            <Terminal class="w-4 h-4" />
          </div>
          <div>
            <div class="flex items-center space-x-2 font-bold text-xs tracking-tight text-white font-heading">
              <span>Recruiter Interactive Showcase</span>
              <span class="text-[9px] bg-emerald-500/20 text-emerald-300 px-2 py-0.5 rounded-full font-mono font-semibold">
                SSE SYNC ACTIVE
              </span>
            </div>
            <div class="text-[10px] text-slate-400">
              Simulate 4 core payment verification and concurrency test cases:
            </div>
          </div>
        </div>

        <!-- 4 Action Scenario Buttons + Auto-Tour Button (Developer Sandbox Style) -->
        <div class="flex flex-wrap items-center gap-2">
          <!-- 1-Click Guided Auto Tour (Sleek Obsidian & Neon Emerald) -->
          <button
            @click="runAutoTour"
            :disabled="isAutoTourRunning || isTesting"
            class="px-3.5 py-1.5 rounded-xl bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-300 border border-emerald-500/40 hover:border-emerald-400 font-bold text-xs transition flex items-center space-x-1.5 active:scale-95 disabled:opacity-50 shadow-xs cursor-pointer"
            title="Single click to simulate full end-to-end booking & verification flow"
          >
            <Play class="w-3.5 h-3.5 fill-emerald-400 text-emerald-400" />
            <span>1-Click Guided Tour</span>
          </button>

          <div class="h-5 w-px bg-slate-800 mx-0.5 hidden sm:block"></div>

          <!-- Scenario Pills (Linear/Stripe Sandbox style: dark neutral with subtle LED status dots) -->
          <div class="bg-slate-950/80 border border-slate-800 p-0.5 rounded-xl flex items-center space-x-0.5">
            <button
              v-for="sc in scenarios"
              :key="sc.key"
              @click="runScenario(sc)"
              :disabled="isTesting || isAutoTourRunning"
              :class="[
                'px-2.5 py-1.5 rounded-lg text-xs font-medium transition flex items-center space-x-1.5 active:scale-95 disabled:opacity-50 select-none cursor-pointer',
                isTesting && sc.key === scenarios[0].key
                  ? 'bg-slate-800 text-white'
                  : 'text-slate-300 hover:text-white hover:bg-slate-800/60'
              ]"
              :title="sc.engineeringExplanation.title"
            >
              <span
                :class="[
                  'w-2 h-2 rounded-full flex-shrink-0',
                  sc.color === 'emerald' ? 'bg-emerald-400' :
                  sc.color === 'rose' ? 'bg-rose-400' :
                  sc.color === 'amber' ? 'bg-amber-400' : 'bg-purple-400'
                ]"
              ></span>
              <span class="text-[11px] font-semibold">{{ sc.label }}</span>
            </button>
          </div>

          <!-- Reset Button -->
          <button
            @click="resetDemo"
            class="p-1.5 bg-slate-800/80 hover:bg-slate-700 text-slate-300 hover:text-white rounded-xl border border-slate-700 text-xs transition cursor-pointer"
            title="Reset Demo State"
          >
            <RefreshCw class="w-3.5 h-3.5" />
          </button>
        </div>
      </div>

      <!-- Engineering Architecture Banner (Collapsible) -->
      <div
        v-if="latestArchitectureNote"
        class="mt-2 bg-slate-900 border border-emerald-500/30 rounded-2xl p-2.5 text-[11px] text-slate-300 flex items-start space-x-2 shadow-inner animate-in fade-in slide-in-from-top-1"
      >
        <Cpu class="w-4 h-4 text-emerald-400 flex-shrink-0 mt-0.5" />
        <div class="flex-1">
          <div class="font-bold text-emerald-300 text-xs mb-0.5">
            {{ latestArchitectureNote.title }}
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-1 text-[10px] text-slate-400 mb-0.5">
            <span v-for="(step, idx) in latestArchitectureNote.steps" :key="idx" class="truncate font-mono">
              {{ step }}
            </span>
          </div>
          <div class="text-[10px] text-emerald-400 font-semibold font-mono">
            Architecture: {{ latestArchitectureNote.result }}
          </div>
        </div>
        <button @click="latestArchitectureNote = null" class="text-slate-500 hover:text-slate-300 text-sm ml-2 cursor-pointer">×</button>
      </div>
    </div>

    <!-- Dual Split View (Viewport Fit: No Outer Double Scrollbar) -->
    <div class="flex-1 flex flex-col md:flex-row overflow-hidden min-h-0">
      <!-- Left Side: Phone Simulation (Strict 380px fixed width) -->
      <div class="w-full md:w-[380px] flex-shrink-0 bg-slate-950 border-b md:border-b-0 md:border-r border-slate-800/80 flex flex-col items-center justify-center p-2 sm:p-3 overflow-hidden">
        <!-- Phone Frame Label Header -->
        <div class="w-full flex items-center justify-between text-[11px] text-slate-400 px-2 pb-1">
          <span class="flex items-center space-x-1.5 font-bold text-emerald-400 font-heading">
            <Smartphone class="w-3.5 h-3.5" />
            <span>Customer: LINE LIFF</span>
          </span>
          <span class="text-[9px] bg-slate-900 border border-slate-800 text-slate-400 font-mono px-1.5 py-0.2 rounded-full">
            iPhone 16 Pro (380px)
          </span>
        </div>

        <!-- Clean Mobile View Container -->
        <div class="w-full h-full flex flex-col overflow-hidden rounded-2xl border border-slate-800 shadow-lg">
          <LiffView />
        </div>
      </div>

      <!-- Right Side: Admin Backoffice & Live POS (Remaining Flexible Width) -->
      <div class="flex-1 bg-slate-100 text-slate-800 flex flex-col overflow-y-auto min-w-0">
        <!-- Header Banner -->
        <div class="bg-slate-200/90 border-b border-slate-300 px-4 py-1.5 text-xs text-slate-700 font-semibold flex items-center justify-between flex-shrink-0">
          <div class="flex items-center space-x-2">
            <Monitor class="w-3.5 h-3.5 text-slate-800" />
            <span class="font-heading font-bold">Store: Admin Backoffice & Live POS Dashboard</span>
          </div>
          <div class="flex items-center space-x-2 text-[10px] text-slate-500 font-mono">
            <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
            <span>SERVER-SENT EVENTS: CONNECTED</span>
          </div>
        </div>

        <!-- Admin Dashboard Body -->
        <div class="flex-1 p-2 sm:p-4 overflow-y-auto">
          <AdminDashboard />
        </div>
      </div>
    </div>
  </div>
</template>
