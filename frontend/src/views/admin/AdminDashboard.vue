<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useAdminStore } from '@/stores/adminStore'
import LiveCourtBoard from './components/LiveCourtBoard.vue'
import OrderAuditTable from './components/OrderAuditTable.vue'
import SlipAuditDrawer from './components/SlipAuditDrawer.vue'
import ReceiptModal from './components/ReceiptModal.vue'
import AddItemsModal from './components/AddItemsModal.vue'
import {
  LayoutGrid, ListFilter, Activity, RefreshCw, Zap
} from 'lucide-vue-next'

const adminStore = useAdminStore()
const activeTab = ref('board') // 'board' or 'orders'

onMounted(() => {
  adminStore.connectSSE()
})

onUnmounted(() => {
  adminStore.disconnectSSE()
})

// Quick Stat Counters
const pendingReviewCount = computed(() => {
  return adminStore.bookings.filter(b => b.status === 'MANUAL_REVIEW').length
})

const confirmedTodayCount = computed(() => {
  return adminStore.bookings.filter(b => b.status === 'CONFIRMED').length
})
</script>

<template>
  <div class="flex-1 bg-slate-100 flex flex-col p-4 sm:p-6 space-y-4">
    <!-- Admin Top Header Banner -->
    <div class="bg-white rounded-2xl border border-slate-200 p-4 sm:p-5 shadow-xs flex flex-wrap items-center justify-between gap-4">
      <div class="flex items-center space-x-3">
        <div class="w-11 h-11 bg-emerald-600 rounded-2xl flex items-center justify-center text-white text-xl shadow-md shadow-emerald-500/20">
          🏸
        </div>
        <div>
          <div class="flex items-center space-x-2">
            <h1 class="text-lg font-extrabold text-slate-900 tracking-tight">SmashHub Backoffice & POS</h1>
            <span class="text-[10px] bg-slate-900 text-white font-mono px-2 py-0.5 rounded-full">v1.0</span>
          </div>
          <p class="text-xs text-slate-500">ระบบจัดการคอร์ทแบดมินตัน, มินิสโตร์ และตรวจสลิปอัตโนมัติ</p>
        </div>
      </div>

      <!-- Real-Time SSE Indicator & Quick Stats -->
      <div class="flex items-center space-x-3">
        <!-- Live SSE Status Badge -->
        <div class="flex items-center space-x-1.5 px-3 py-1.5 rounded-full bg-slate-50 border border-slate-200 text-xs font-medium">
          <span
            :class="[
              'w-2.5 h-2.5 rounded-full transition-all',
              adminStore.isSSEConnected ? 'bg-emerald-500 shadow-sm shadow-emerald-400 animate-pulse' : 'bg-rose-500'
            ]"
          ></span>
          <span :class="adminStore.isSSEConnected ? 'text-emerald-700 font-semibold' : 'text-slate-500'">
            {{ adminStore.isSSEConnected ? 'Real-Time Sync (SSE)' : 'Connecting...' }}
          </span>
        </div>

        <!-- Pending Review Badge -->
        <div
          v-if="pendingReviewCount > 0"
          class="flex items-center space-x-1 px-3 py-1.5 bg-amber-50 text-amber-900 border border-amber-300 rounded-full text-xs font-bold animate-bounce"
        >
          <span>⚠️ รอยืนยันสลิป:</span>
          <span class="font-mono text-amber-700">{{ pendingReviewCount }} รายการ</span>
        </div>
      </div>
    </div>

    <!-- Navigation Tabs Bar -->
    <div class="flex items-center space-x-2 border-b border-slate-200 pb-2">
      <button
        @click="activeTab = 'board'"
        :class="[
          'px-4 py-2 rounded-xl text-xs font-bold transition flex items-center space-x-2',
          activeTab === 'board'
            ? 'bg-slate-900 text-white shadow'
            : 'bg-white text-slate-600 hover:bg-slate-200 border border-slate-200'
        ]"
      >
        <LayoutGrid class="w-4 h-4" />
        <span>Live Court Board & POS (กระดานคอร์ท)</span>
      </button>

      <button
        @click="activeTab = 'orders'"
        :class="[
          'px-4 py-2 rounded-xl text-xs font-bold transition flex items-center space-x-2 relative',
          activeTab === 'orders'
            ? 'bg-slate-900 text-white shadow'
            : 'bg-white text-slate-600 hover:bg-slate-200 border border-slate-200'
        ]"
      >
        <ListFilter class="w-4 h-4" />
        <span>Order & Slip Audit Log (รายการจอง & ตรวจสลิป)</span>
        <span
          v-if="pendingReviewCount > 0"
          class="ml-1 px-1.5 py-0.2 bg-rose-500 text-white rounded-full text-[10px] font-mono font-bold"
        >
          {{ pendingReviewCount }}
        </span>
      </button>
    </div>

    <!-- Tab View Content -->
    <div class="flex-1">
      <LiveCourtBoard v-if="activeTab === 'board'" />
      <OrderAuditTable v-else-if="activeTab === 'orders'" />
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
