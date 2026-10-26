<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useAdminStore } from '@/stores/adminStore'
import { useAuthStore } from '@/stores/authStore'
import LiveCourtBoard from './components/LiveCourtBoard.vue'
import OrderAuditTable from './components/OrderAuditTable.vue'
import SlipAuditDrawer from './components/SlipAuditDrawer.vue'
import ReceiptModal from './components/ReceiptModal.vue'
import AddItemsModal from './components/AddItemsModal.vue'
import AdminInventoryManagement from './components/AdminInventoryManagement.vue'
import {
  LayoutGrid, ListFilter, DollarSign, Activity, Users,
  AlertTriangle, TrendingUp, Sparkles, CheckCircle2, Clock, LogIn, ShieldAlert,
  Settings2
} from 'lucide-vue-next'

const adminStore = useAdminStore()
const authStore = useAuthStore()
const activeTab = ref('board') // 'board', 'orders', or 'inventory'

onMounted(() => {
  adminStore.connectSSE()
  adminStore.fetchBookings()
})

onUnmounted(() => {
  adminStore.disconnectSSE()
})

// KPI Metrics Calculations
const totalRevenueToday = computed(() => {
  return adminStore.bookings
    .filter(b => b.status === 'CONFIRMED')
    .reduce((sum, b) => sum + (b.net_total || 0), 0)
})

const pendingReviewCount = computed(() => {
  return adminStore.bookings.filter(b => b.status === 'MANUAL_REVIEW').length
})

const confirmedBookingsCount = computed(() => {
  return adminStore.bookings.filter(b => b.status === 'CONFIRMED').length
})

const checkedInCount = computed(() => {
  return adminStore.bookings.filter(b => b.check_in_status === 'CHECKED_IN').length
})
</script>

<template>
  <div class="flex-1 bg-[#F8FAFC] flex flex-col p-4 sm:p-6 space-y-4">
    <!-- Admin Access Warning Banner if not Admin -->
    <div
      v-if="!authStore.isAdmin"
      class="bg-amber-50 border border-amber-200/90 rounded-2xl p-3.5 flex flex-wrap items-center justify-between gap-3 shadow-2xs"
    >
      <div class="flex items-center space-x-2.5">
        <div class="w-8 h-8 rounded-xl bg-amber-100 text-amber-800 flex items-center justify-center font-bold text-sm">
          👑
        </div>
        <div>
          <div class="text-xs font-bold text-amber-900">หน้านี้สำหรับผู้จัดการร้าน (Admin Backoffice)</div>
          <p class="text-[11px] text-amber-700">คุณยังไม่ได้เข้าสู่ระบบด้วยสิทธิ์ผู้ดูแลระบบ (Admin) หรือเข้าสู่ระบบในฐานะลูกค้าอยู่</p>
        </div>
      </div>

      <div class="flex items-center space-x-2">
        <button
          @click="authStore.demoLogin('ADMIN')"
          class="px-3.5 py-1.5 bg-slate-900 hover:bg-slate-800 text-white rounded-xl text-xs font-bold shadow-xs transition flex items-center space-x-1.5"
        >
          <span>👑 สลับเป็น Admin ทันที (1-Click)</span>
        </button>
        <button
          @click="authStore.openAuthModal('login')"
          class="px-3 py-1.5 bg-white hover:bg-amber-100/60 text-slate-700 border border-amber-200 rounded-xl text-xs font-semibold transition"
        >
          <span>กรอกรหัสผ่านเอง</span>
        </button>
      </div>
    </div>
    <!-- Top Header Banner & Real-Time Sync Indicator -->
    <div class="bg-white rounded-2xl border border-slate-200/90 p-4 shadow-xs flex flex-wrap items-center justify-between gap-3">
      <div class="flex items-center space-x-3">
        <div class="w-10 h-10 bg-slate-900 rounded-xl flex items-center justify-center text-white text-lg shadow-xs">
          🏸
        </div>
        <div>
          <div class="flex items-center space-x-2">
            <h1 class="text-base font-extrabold text-slate-900 tracking-tight">SmashHub Backoffice & POS</h1>
            <span class="text-[10px] bg-slate-100 text-slate-600 border border-slate-200 font-mono font-semibold px-2 py-0.5 rounded-md">Management</span>
          </div>
          <p class="text-[11px] text-slate-500">ระบบบริหารจัดการสนามแบดมินตัน มินิสโตร์แคชเชียร์ และตรวจสอบสลิป</p>
        </div>
      </div>

      <!-- Real-Time SSE Indicator Pill -->
      <div class="flex items-center space-x-2">
        <div class="flex items-center space-x-2 px-3 py-1.5 rounded-xl bg-slate-50 border border-slate-200 text-xs font-mono">
          <span
            :class="[
              'w-2 h-2 rounded-full transition-all',
              adminStore.isSSEConnected ? 'bg-emerald-500 animate-pulse' : 'bg-rose-500'
            ]"
          ></span>
          <span :class="adminStore.isSSEConnected ? 'text-emerald-800 font-bold' : 'text-slate-500'">
            {{ adminStore.isSSEConnected ? 'SSE Live Connected' : 'Connecting...' }}
          </span>
        </div>
      </div>
    </div>

    <!-- Executive KPI Metric Cards (Stripe/Linear inspired) -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
      <!-- 1. Total Revenue -->
      <div class="bg-white rounded-2xl border border-slate-200/90 p-3.5 shadow-2xs">
        <div class="flex justify-between items-start text-slate-400 text-[11px] font-medium mb-1">
          <span>รายได้รวมวันนี้ (Revenue)</span>
          <div class="p-1.5 rounded-lg bg-emerald-50 text-emerald-600">
            <DollarSign class="w-3.5 h-3.5" />
          </div>
        </div>
        <div class="text-xl font-extrabold font-mono text-slate-900">
          ฿{{ totalRevenueToday.toLocaleString() }}
        </div>
        <div class="flex items-center space-x-1 text-[10px] text-emerald-600 font-medium mt-1">
          <TrendingUp class="w-3 h-3" />
          <span>+14.8% เทียบกับสัปดาห์ก่อน</span>
        </div>
      </div>

      <!-- 2. Confirmed Bookings -->
      <div class="bg-white rounded-2xl border border-slate-200/90 p-3.5 shadow-2xs">
        <div class="flex justify-between items-start text-slate-400 text-[11px] font-medium mb-1">
          <span>ยอดจองสำเร็จ (Confirmed)</span>
          <div class="p-1.5 rounded-lg bg-blue-50 text-blue-600">
            <CheckCircle2 class="w-3.5 h-3.5" />
          </div>
        </div>
        <div class="text-xl font-extrabold font-mono text-slate-900">
          {{ confirmedBookingsCount }} รายการ
        </div>
        <div class="text-[10px] text-slate-500 mt-1">
          อัตราครองคอร์ทเฉลี่ย 82%
        </div>
      </div>

      <!-- 3. Pending Review -->
      <div
        :class="[
          'rounded-2xl border p-3.5 shadow-2xs transition',
          pendingReviewCount > 0 ? 'bg-amber-50/80 border-amber-300' : 'bg-white border-slate-200/90'
        ]"
      >
        <div class="flex justify-between items-start text-[11px] font-medium mb-1">
          <span :class="pendingReviewCount > 0 ? 'text-amber-900 font-bold' : 'text-slate-400'">
            สลิปรอตรวจสอบ (Pending Audit)
          </span>
          <div :class="['p-1.5 rounded-lg', pendingReviewCount > 0 ? 'bg-amber-200 text-amber-800' : 'bg-slate-100 text-slate-500']">
            <AlertTriangle class="w-3.5 h-3.5" />
          </div>
        </div>
        <div :class="['text-xl font-extrabold font-mono', pendingReviewCount > 0 ? 'text-amber-900' : 'text-slate-900']">
          {{ pendingReviewCount }} รายการ
        </div>
        <div :class="['text-[10px] mt-1', pendingReviewCount > 0 ? 'text-amber-700 font-semibold' : 'text-slate-400']">
          {{ pendingReviewCount > 0 ? '⚠️ ต้องกดอนุมัติด้วยตนเอง' : 'ไม่มีสลิปตกค้าง' }}
        </div>
      </div>

      <!-- 4. Players Checked-in -->
      <div class="bg-white rounded-2xl border border-slate-200/90 p-3.5 shadow-2xs">
        <div class="flex justify-between items-start text-slate-400 text-[11px] font-medium mb-1">
          <span>ลูกค้าลงสนามแล้ว (Checked-in)</span>
          <div class="p-1.5 rounded-lg bg-purple-50 text-purple-600">
            <Users class="w-3.5 h-3.5" />
          </div>
        </div>
        <div class="text-xl font-extrabold font-mono text-slate-900">
          {{ checkedInCount }} แมตช์
        </div>
        <div class="text-[10px] text-slate-500 mt-1">
          กำลังเล่นอยู่ ณ ขณะนี้
        </div>
      </div>
    </div>

    <!-- Modern Segment Tab Switcher -->
    <div class="bg-white p-1 rounded-2xl border border-slate-200/90 inline-flex space-x-1 shadow-2xs">
      <button
        @click="activeTab = 'board'"
        :class="[
          'px-4 py-2 rounded-xl text-xs font-bold transition flex items-center space-x-2',
          activeTab === 'board'
            ? 'bg-slate-900 text-white shadow-sm'
            : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100'
        ]"
      >
        <LayoutGrid class="w-4 h-4" />
        <span>Live Court Board & POS (กระดานคอร์ทสด)</span>
      </button>

      <button
        @click="activeTab = 'orders'"
        :class="[
          'px-4 py-2 rounded-xl text-xs font-bold transition flex items-center space-x-2 relative',
          activeTab === 'orders'
            ? 'bg-slate-900 text-white shadow-sm'
            : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100'
        ]"
      >
        <ListFilter class="w-4 h-4" />
        <span>Order & Slip Audit Log (ประวัติการจอง & ตรวจสลิป)</span>
        <span
          v-if="pendingReviewCount > 0"
          class="ml-1 px-1.5 py-0.2 bg-rose-500 text-white rounded-full text-[9px] font-mono font-bold"
        >
          {{ pendingReviewCount }}
        </span>
      </button>

      <button
        @click="activeTab = 'inventory'"
        :class="[
          'px-4 py-2 rounded-xl text-xs font-bold transition flex items-center space-x-2',
          activeTab === 'inventory'
            ? 'bg-slate-900 text-white shadow-sm'
            : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100'
        ]"
      >
        <Settings2 class="w-4 h-4" />
        <span>Court & Stock Control (จัดการสนาม & เติมสต็อก)</span>
      </button>
    </div>

    <!-- Active Tab Content Area -->
    <div class="flex-1">
      <LiveCourtBoard v-if="activeTab === 'board'" />
      <OrderAuditTable v-else-if="activeTab === 'orders'" />
      <AdminInventoryManagement v-else-if="activeTab === 'inventory'" />
    </div>

    <!-- Modals & Drawers -->
    <SlipAuditDrawer />
    <ReceiptModal
      v-model:visible="adminStore.isReceiptOpen"
      :booking="adminStore.receiptBooking"
    />
    <AddItemsModal />
  </div>
</template>
