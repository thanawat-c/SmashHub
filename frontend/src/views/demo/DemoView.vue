<script setup>
import { ref } from 'vue'
import LiffView from '../liff/LiffView.vue'
import AdminDashboard from '../admin/AdminDashboard.vue'
import { useBookingStore } from '@/stores/bookingStore'
import { useLineProfileStore } from '@/stores/lineProfile'
import { useAdminStore } from '@/stores/adminStore'
import {
  Smartphone, Monitor, Sparkles, CheckCircle, XCircle, AlertTriangle,
  Server, ShieldAlert, Cpu, ArrowRight, Play, RefreshCw, Info
} from 'lucide-vue-next'
import { ElNotification } from 'element-plus'

const bookingStore = useBookingStore()
const profileStore = useLineProfileStore()
const adminStore = useAdminStore()

const isTesting = ref(false)
const latestArchitectureNote = ref(null)

const scenarios = [
  {
    key: 'MOCK_VALID',
    label: 'จำลองสลิปถูกต้อง',
    subLabel: 'ยอดตรง + สลิปใหม่',
    btnClass: 'bg-emerald-600 hover:bg-emerald-500 text-white border-emerald-400/50 shadow-emerald-900/30',
    icon: CheckCircle,
    color: 'emerald',
    badge: '🟢 MOCK_VALID',
    engineeringExplanation: {
      title: '✅ 4-Step Verification: ผ่านครบทุกเงื่อนไข',
      steps: [
        '1. Check trans_ref ซ้ำ: ไม่พบในตาราง payment_slips',
        '2. Check บัญชีผู้รับ: ตรงกับบัญชีร้าน (098-7-65432-1)',
        '3. Check ยอดโอน: ฿' + (bookingStore.activeBooking?.booking?.net_total || 250) + ' >= ยอดบิล net_total',
        '4. External Gateway: ตอบกลับสำเร็จ (HTTP 200)'
      ],
      result: 'DB Transition: booking.status = CONFIRMED, is_slot_locked = TRUE, ออกบัตร E-Ticket และ Broadcast SSE ไปยัง Admin Board ทันที'
    }
  },
  {
    key: 'MOCK_DUPLICATE',
    label: 'จำลองสลิปเวียนเทียน',
    subLabel: 'ใช้สลิปซ้ำเดิม',
    btnClass: 'bg-rose-600 hover:bg-rose-500 text-white border-rose-400/50 shadow-rose-900/30',
    icon: XCircle,
    color: 'rose',
    badge: '🔴 MOCK_DUPLICATE',
    engineeringExplanation: {
      title: '🚨 Step 1 Failed: ตรวจพบสลิปเวียนเทียน (Duplicate)',
      steps: [
        '1. Check trans_ref ซ้ำ: ค้นพบ trans_ref ซ้ำในตาราง payment_slips (Key ID มีอยู่แล้ว)',
        '2. Transaction Rollback ป้องกัน Double Spending',
        '3. บันทึก Slip Log: verification_status = REJECTED_DUPLICATE',
        '4. ส่ง HTTP 409 Conflict กลับไปยังหน้า LIFF ลูกค้า'
      ],
      result: 'สถานะใบจองยังคงเป็น PENDING_PAYMENT และยังคงล็อคเวลา 10 นาทีเดิมไว้ให้ลูกค้าชำระเงินใหม่'
    }
  },
  {
    key: 'MOCK_UNDERPAID',
    label: 'จำลองสลิปโอนไม่ครบ',
    subLabel: 'ยอดเงิน < ยอดบิล',
    btnClass: 'bg-amber-600 hover:bg-amber-500 text-white border-amber-400/50 shadow-amber-900/30',
    icon: AlertTriangle,
    color: 'amber',
    badge: '🟡 MOCK_UNDERPAID',
    engineeringExplanation: {
      title: '⚠️ Step 3 Failed: ยอดเงินในสลิปไม่ตรง (Underpaid)',
      steps: [
        '1. Check trans_ref: ผ่าน (สลิปใหม่)',
        '2. Check บัญชีผู้รับ: ผ่าน',
        '3. Check ยอดเงิน: ยอดโอนในสลิป น้อยกว่า booking.net_total',
        '4. บันทึก Slip Log: verification_status = REJECTED_AMOUNT_MISMATCH'
      ],
      result: 'ระบบตอบกลับ HTTP 402 Payment Required แจ้งยอดที่ขาด และแจ้งเตือนให้แนบสลิปใหม่'
    }
  },
  {
    key: 'MOCK_API_DOWN',
    label: 'จำลองธนาคารล่ม',
    subLabel: 'ส่งให้แอดมินตรวจเอง',
    btnClass: 'bg-purple-600 hover:bg-purple-500 text-white border-purple-400/50 shadow-purple-900/30',
    icon: Server,
    color: 'purple',
    badge: '🟣 MOCK_API_DOWN',
    engineeringExplanation: {
      title: '🛡️ Circuit Breaker: Gateway ตรวจสลิปภายนอกล่ม (504 Timeout)',
      steps: [
        '1. ส่ง Request ไปยัง Slip Verification Gateway -> จำลอง Timeout 504',
        '2. Graceful Fallback: ไม่ปฏิเสธลูกค้า แต่ส่งเคสเข้าคิว Manual Audit',
        '3. อัปเดตสถานะ: booking.status = MANUAL_REVIEW',
        '4. อัปเดตสลิป: verification_status = API_ERROR_PENDING_ADMIN'
      ],
      result: 'Broadcast SSE: แจ้งเตือนสลิปรอตรวจสอบบนหน้าจอ Admin ทันที พร้อมปุ่ม Approve/Reject ใน Drawer'
    }
  }
]

async function runScenario(scenario) {
  isTesting.value = true
  try {
    // 1. Ensure an active booking is in payment state
    await bookingStore.ensureActiveBookingForDemo(profileStore.currentProfile)

    // 2. Call the slip verification endpoint
    const res = await bookingStore.verifySlip(scenario.key)

    // 3. Set the active explanation for recruiter
    latestArchitectureNote.value = scenario.engineeringExplanation

    // 4. Show Toast Notification explaining engineering workflow
    ElNotification({
      title: scenario.badge + ' - ' + scenario.label,
      dangerouslyUseHTMLString: true,
      message: `
        <div style="font-size: 11px; line-height: 1.5; color: #334155;">
          <p style="font-weight: 700; margin-bottom: 4px; color: #0f172a;">${scenario.engineeringExplanation.title}</p>
          <ul style="padding-left: 14px; margin: 0; list-style-type: disc;">
            ${scenario.engineeringExplanation.steps.map(s => `<li>${s}</li>`).join('')}
          </ul>
          <p style="margin-top: 6px; font-weight: 600; color: #047857;">⚡ ${scenario.engineeringExplanation.result}</p>
        </div>
      `,
      duration: 6000,
      type: scenario.color === 'emerald' ? 'success' : scenario.color === 'rose' ? 'error' : 'warning'
    })

    // 5. Trigger live reload of admin bookings to reflect changes
    adminStore.fetchBookings()
  } catch (err) {
    console.error('Scenario error:', err)
  } finally {
    isTesting.value = false
  }
}

function resetDemo() {
  bookingStore.resetSelection()
  bookingStore.fetchAvailability()
  latestArchitectureNote.value = null
  ElNotification({
    title: 'รีเซ็ตสถานะเรียบร้อย',
    message: 'พร้อมสำหรับการจำลองรอบใหม่',
    type: 'info',
    duration: 2000
  })
}
</script>

<template>
  <div class="flex-1 flex flex-col bg-slate-950 text-white min-h-[calc(100vh-45px)] select-none">
    <!-- Top Interactive Test Bar for Recruiters -->
    <div class="bg-slate-900 border-b border-slate-800 px-3 py-2.5 shadow-lg sticky top-[45px] z-40">
      <div class="flex flex-wrap items-center justify-between gap-2.5">
        <!-- Title & Subtitle -->
        <div class="flex items-center space-x-2">
          <div class="p-1.5 bg-emerald-500/10 border border-emerald-500/30 rounded-lg text-emerald-400">
            <Sparkles class="w-4 h-4" />
          </div>
          <div>
            <div class="flex items-center space-x-1.5 font-bold text-xs tracking-wide text-white">
              <span>Recruiter Interactive Test Bar</span>
              <span class="text-[10px] bg-emerald-500/20 text-emerald-300 px-1.5 py-0.2 rounded font-mono">Real-time SSE</span>
            </div>
            <div class="text-[10px] text-slate-400">คลิกทดสอบสถานการณ์จริง 4 รูปแบบ (จำลองสลิป & Concurrency):</div>
          </div>
        </div>

        <!-- 4 Scenario Action Buttons -->
        <div class="flex flex-wrap items-center gap-1.5">
          <button
            v-for="sc in scenarios"
            :key="sc.key"
            @click="runScenario(sc)"
            :disabled="isTesting"
            :class="[
              'px-2.5 py-1.5 rounded-xl border text-xs font-semibold shadow transition flex items-center space-x-1.5 active:scale-95 disabled:opacity-50',
              sc.btnClass
            ]"
            :title="sc.engineeringExplanation.title"
          >
            <component :is="sc.icon" class="w-3.5 h-3.5 flex-shrink-0" />
            <div class="text-left leading-tight">
              <div class="text-[11px] font-bold">{{ sc.label }}</div>
              <div class="text-[9px] opacity-80 font-normal hidden sm:block">{{ sc.subLabel }}</div>
            </div>
          </button>

          <!-- Reset Button -->
          <button
            @click="resetDemo"
            class="p-2 bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white rounded-xl border border-slate-700 text-xs transition"
            title="รีเซ็ตการจองเริ่มใหม่"
          >
            <RefreshCw class="w-3.5 h-3.5" />
          </button>
        </div>
      </div>

      <!-- Engineering Live Architecture Banner (Collapsible/Floating) -->
      <div
        v-if="latestArchitectureNote"
        class="mt-2 bg-slate-900/90 border border-emerald-500/30 rounded-xl p-2.5 text-[11px] text-slate-300 flex items-start space-x-2 animate-in fade-in slide-in-from-top-1"
      >
        <Cpu class="w-4 h-4 text-emerald-400 flex-shrink-0 mt-0.5" />
        <div class="flex-1">
          <div class="font-bold text-emerald-300 text-xs mb-1">
            {{ latestArchitectureNote.title }}
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-1 text-[10px] text-slate-400 mb-1">
            <span v-for="(step, idx) in latestArchitectureNote.steps" :key="idx" class="truncate">
              {{ step }}
            </span>
          </div>
          <div class="text-[10px] text-emerald-400/90 font-medium">
            ⚡ {{ latestArchitectureNote.result }}
          </div>
        </div>
        <button @click="latestArchitectureNote = null" class="text-slate-500 hover:text-slate-300 text-xs ml-2">×</button>
      </div>
    </div>

    <!-- Dual Split Screen Layout -->
    <div class="flex-1 flex flex-col md:flex-row overflow-hidden">
      <!-- Left Side: Mobile LIFF Frame (Strict 380px fixed width on desktop) -->
      <div class="w-full md:w-[380px] flex-shrink-0 bg-slate-950 border-b md:border-b-0 md:border-r border-slate-800 flex flex-col items-center justify-start p-2 sm:p-3 overflow-y-auto">
        <!-- Label Badge -->
        <div class="w-full flex items-center justify-between text-[11px] text-slate-400 px-1 py-1 mb-1">
          <span class="flex items-center space-x-1.5 font-semibold text-emerald-400">
            <Smartphone class="w-3.5 h-3.5" />
            <span>ลูกค้า: LINE LIFF Storefront</span>
          </span>
          <span class="text-[9px] bg-slate-800 border border-slate-700 text-slate-400 font-mono px-1.5 py-0.2 rounded">
            380px Width
          </span>
        </div>

        <!-- Phone Mockup Container (max-w-[380px] h-[780px]) -->
        <div class="w-full h-[780px] max-h-[85vh] bg-slate-50 rounded-[32px] overflow-hidden border-[6px] border-slate-800 shadow-2xl flex flex-col relative ring-1 ring-slate-700">
          <LiffView />
        </div>
      </div>

      <!-- Right Side: Store Backoffice & Live Board (Remaining Flexible Width) -->
      <div class="flex-1 bg-slate-100 text-slate-800 flex flex-col overflow-y-auto min-w-0">
        <!-- Header Banner -->
        <div class="bg-slate-200/90 border-b border-slate-300 px-4 py-1.5 text-xs text-slate-700 font-semibold flex items-center justify-between">
          <div class="flex items-center space-x-2">
            <Monitor class="w-4 h-4 text-emerald-700" />
            <span>ร้านค้า: Admin Backoffice & Live POS Dashboard</span>
          </div>
          <div class="flex items-center space-x-2 text-[10px] text-slate-500 font-mono">
            <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
            <span>SSE STREAM: ACTIVE</span>
          </div>
        </div>

        <!-- Admin Dashboard Component -->
        <div class="flex-1 p-2 sm:p-4">
          <AdminDashboard />
        </div>
      </div>
    </div>
  </div>
</template>
