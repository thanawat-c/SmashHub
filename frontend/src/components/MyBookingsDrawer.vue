<script setup>
import { ref, computed } from 'vue'
import { useBookingStore } from '@/stores/bookingStore'
import { useAuthStore } from '@/stores/authStore'
import { useLineProfileStore } from '@/stores/lineProfile'
import {
  Calendar, Clock, Ticket, CreditCard, CheckCircle2,
  AlertTriangle, XCircle, RotateCcw, Copy, Check, ChevronRight, X, ExternalLink
} from 'lucide-vue-next'

const bookingStore = useBookingStore()
const authStore = useAuthStore()
const profileStore = useLineProfileStore()

const activeFilter = ref('ALL') // ALL, CONFIRMED, PENDING_PAYMENT, MANUAL_REVIEW

const copiedId = ref('')
function copyBookingNo(bookingNo) {
  navigator.clipboard.writeText(bookingNo)
  copiedId.value = bookingNo
  setTimeout(() => {
    copiedId.value = ''
  }, 2000)
}

const filteredBookings = computed(() => {
  if (activeFilter.value === 'ALL') return bookingStore.myBookings
  return bookingStore.myBookings.filter(b => b.status === activeFilter.value)
})

function formatDateTime(dtStr) {
  if (!dtStr) return '-'
  const d = new Date(dtStr)
  return d.toLocaleDateString('th-TH', {
    day: 'numeric',
    month: 'short',
    year: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function formatPlayDate(dStr) {
  if (!dStr) return '-'
  const d = new Date(dStr)
  return d.toLocaleDateString('th-TH', {
    weekday: 'short',
    day: 'numeric',
    month: 'short',
    year: 'numeric'
  })
}
</script>

<template>
  <el-drawer
    v-model="bookingStore.isMyBookingsOpen"
    title="ประวัติการจองและคำสั่งซื้อ (My Orders)"
    size="460px"
    direction="rtl"
    class="rounded-l-3xl overflow-hidden"
    :with-header="false"
  >
    <div class="h-full flex flex-col bg-[#F8FAFC]">
      <!-- Drawer Custom Header -->
      <div class="bg-white px-5 py-4 border-b border-slate-200/90 shadow-2xs flex items-center justify-between">
        <div class="flex items-center space-x-2.5">
          <div class="w-9 h-9 rounded-xl bg-emerald-50 text-emerald-700 flex items-center justify-center font-bold text-sm shadow-2xs">
            📋
          </div>
          <div>
            <h2 class="text-sm font-extrabold text-slate-900 tracking-tight">ประวัติการจองของฉัน</h2>
            <p class="text-[11px] text-slate-500">
              ผู้จอง: <strong class="text-slate-800">{{ authStore.displayName }}</strong>
              <span v-if="authStore.user?.phone" class="font-mono text-[10px] text-slate-400"> ({{ authStore.user.phone }})</span>
            </p>
          </div>
        </div>

        <button
          @click="bookingStore.closeMyBookings"
          class="p-1.5 text-slate-400 hover:text-slate-700 hover:bg-slate-100 rounded-xl transition"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Filter Tabs -->
      <div class="bg-white px-5 py-2.5 border-b border-slate-200/80 flex space-x-1.5 overflow-x-auto scrollbar-none">
        <button
          @click="activeFilter = 'ALL'"
          :class="[
            'px-3 py-1 rounded-full text-xs font-semibold whitespace-nowrap transition border',
            activeFilter === 'ALL'
              ? 'bg-slate-900 text-white border-slate-900 shadow-2xs'
              : 'bg-slate-50 text-slate-600 border-slate-200 hover:bg-slate-100'
          ]"
        >
          ทั้งหมด ({{ bookingStore.myBookings.length }})
        </button>

        <button
          @click="activeFilter = 'CONFIRMED'"
          :class="[
            'px-3 py-1 rounded-full text-xs font-semibold whitespace-nowrap transition border',
            activeFilter === 'CONFIRMED'
              ? 'bg-emerald-600 text-white border-emerald-600 shadow-2xs'
              : 'bg-slate-50 text-slate-600 border-slate-200 hover:bg-slate-100'
          ]"
        >
          ✅ ยืนยันแล้ว
        </button>

        <button
          @click="activeFilter = 'PENDING_PAYMENT'"
          :class="[
            'px-3 py-1 rounded-full text-xs font-semibold whitespace-nowrap transition border',
            activeFilter === 'PENDING_PAYMENT'
              ? 'bg-amber-600 text-white border-amber-600 shadow-2xs'
              : 'bg-slate-50 text-slate-600 border-slate-200 hover:bg-slate-100'
          ]"
        >
          ⏱️ รอชำระเงิน
        </button>

        <button
          @click="activeFilter = 'MANUAL_REVIEW'"
          :class="[
            'px-3 py-1 rounded-full text-xs font-semibold whitespace-nowrap transition border',
            activeFilter === 'MANUAL_REVIEW'
              ? 'bg-purple-600 text-white border-purple-600 shadow-2xs'
              : 'bg-slate-50 text-slate-600 border-slate-200 hover:bg-slate-100'
          ]"
        >
          ⚠️ รอตรวจสลิป
        </button>
      </div>

      <!-- Orders List Container -->
      <div class="flex-1 overflow-y-auto p-4 space-y-3">
        <!-- Loading Indicator -->
        <div v-if="bookingStore.isFetchingMyBookings" class="py-20 text-center text-slate-400 space-y-2">
          <div class="w-6 h-6 border-2 border-emerald-500 border-t-transparent rounded-full animate-spin mx-auto"></div>
          <p class="text-xs">กำลังค้นหาประวัติการจอง...</p>
        </div>

        <!-- Empty State -->
        <div v-else-if="filteredBookings.length === 0" class="py-24 text-center space-y-3">
          <div class="w-14 h-14 bg-slate-100 rounded-3xl flex items-center justify-center text-2xl mx-auto shadow-2xs">
            🏸
          </div>
          <div>
            <h3 class="text-sm font-extrabold text-slate-800">ยังไม่พบรายการจองตามเงื่อนไข</h3>
            <p class="text-xs text-slate-400 mt-1 max-w-[260px] mx-auto">
              คุณสามารถเริ่มจองสนามหรือสั่งอุปกรณ์แบดมินตันได้ง่ายๆ เพียงเลือกช่วงเวลาที่ต้องการ
            </p>
          </div>
          <button
            @click="bookingStore.closeMyBookings(); bookingStore.currentStep = 1"
            class="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-bold rounded-xl shadow-xs transition"
          >
            ไปหน้าจองสนาม
          </button>
        </div>

        <!-- Order Cards -->
        <div
          v-else
          v-for="booking in filteredBookings"
          :key="booking.id"
          class="bg-white rounded-2xl border border-slate-200/90 p-4 shadow-2xs space-y-3 hover:border-slate-300 transition"
        >
          <!-- Card Header: Booking No & Status Badge -->
          <div class="flex items-center justify-between border-b border-slate-100 pb-2.5">
            <div class="flex items-center space-x-1.5">
              <span class="font-mono font-black text-xs text-slate-900 tracking-tight">{{ booking.booking_no }}</span>
              <button
                @click="copyBookingNo(booking.booking_no)"
                class="p-1 text-slate-400 hover:text-slate-700 rounded transition"
                title="คัดลอกรหัสการจอง"
              >
                <Check v-if="copiedId === booking.booking_no" class="w-3 h-3 text-emerald-600" />
                <Copy v-else class="w-3 h-3" />
              </button>
            </div>

            <!-- Status Pill -->
            <span
              :class="[
                'text-[10px] font-mono font-bold px-2 py-0.5 rounded-full border',
                booking.status === 'CONFIRMED' ? 'bg-emerald-50 text-emerald-700 border-emerald-300' :
                booking.status === 'PENDING_PAYMENT' ? 'bg-amber-50 text-amber-800 border-amber-300' :
                booking.status === 'MANUAL_REVIEW' ? 'bg-purple-50 text-purple-700 border-purple-300' :
                'bg-slate-100 text-slate-500 border-slate-200'
              ]"
            >
              {{
                booking.status === 'CONFIRMED' ? '● CONFIRMED' :
                booking.status === 'PENDING_PAYMENT' ? '⏱️ PENDING_PAYMENT' :
                booking.status === 'MANUAL_REVIEW' ? '⚠️ MANUAL_REVIEW' :
                booking.status
              }}
            </span>
          </div>

          <!-- Match / Play Details -->
          <div class="text-xs space-y-1.5 text-slate-700">
            <div class="flex items-center space-x-2 text-slate-500">
              <Calendar class="w-3.5 h-3.5 text-slate-400" />
              <span>วันที่เล่น: <strong class="text-slate-900">{{ formatPlayDate(booking.courts?.[0]?.play_date) }}</strong></span>
            </div>

            <div class="flex items-start space-x-2">
              <Clock class="w-3.5 h-3.5 text-slate-400 mt-0.5 flex-shrink-0" />
              <div class="flex flex-wrap gap-1">
                <span
                  v-for="c in booking.courts"
                  :key="c.id"
                  class="bg-slate-100 text-slate-800 border border-slate-200 px-1.5 py-0.2 rounded text-[10px] font-mono font-medium"
                >
                  คอร์ท {{ c.court?.court_number }} ({{ c.start_hour }}:00 - {{ c.end_hour }}:00 น.)
                </span>
              </div>
            </div>

            <!-- Mini Store items if any -->
            <div v-if="booking.items && booking.items.length > 0" class="pt-1 text-[11px] text-slate-500">
              <span class="text-slate-400">สินค้า: </span>
              <span class="text-slate-700">{{ booking.items.map(it => `${it.product?.name || 'สินค้า'} x${it.qty}`).join(', ') }}</span>
            </div>
          </div>

          <!-- Total & Actions -->
          <div class="flex items-center justify-between pt-2 border-t border-slate-100">
            <div>
              <div class="text-[10px] text-slate-400">ยอดชำระสุทธิ</div>
              <div class="text-sm font-black font-mono text-emerald-700">
                ฿{{ Number(booking.net_total || 0).toLocaleString() }}
              </div>
            </div>

            <!-- Contextual Action Buttons -->
            <div>
              <!-- 1. If CONFIRMED: Open E-Ticket / Gate QR Code -->
              <button
                v-if="booking.status === 'CONFIRMED'"
                @click="bookingStore.viewBookingTicket(booking)"
                class="px-3.5 py-1.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-xs font-bold shadow-2xs transition flex items-center space-x-1.5"
              >
                <Ticket class="w-3.5 h-3.5" />
                <span>เปิดบัตรเข้าสนาม (E-Ticket)</span>
              </button>

              <!-- 2. If PENDING_PAYMENT: Resume Payment -->
              <button
                v-else-if="booking.status === 'PENDING_PAYMENT'"
                @click="bookingStore.resumeBookingPayment(booking)"
                class="px-3.5 py-1.5 bg-amber-500 hover:bg-amber-600 text-white rounded-xl text-xs font-bold shadow-2xs transition flex items-center space-x-1.5"
              >
                <CreditCard class="w-3.5 h-3.5" />
                <span>แนบสลิป / ชำระเงิน</span>
              </button>

              <!-- 3. If MANUAL_REVIEW -->
              <span
                v-else-if="booking.status === 'MANUAL_REVIEW'"
                class="text-[11px] text-purple-700 font-semibold"
              >
                รอเจ้าหน้าที่อนุมัติ
              </span>

              <!-- 4. If EXPIRED -->
              <span
                v-else-if="booking.status === 'EXPIRED'"
                class="text-[11px] text-slate-400"
              >
                รายการหมดเวลา 10 นาที
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </el-drawer>
</template>
