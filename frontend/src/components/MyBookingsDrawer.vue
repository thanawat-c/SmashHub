<script setup>
import { ref, computed } from 'vue'
import { useBookingStore } from '@/stores/bookingStore'
import { useAuthStore } from '@/stores/authStore'
import { useLineProfileStore } from '@/stores/lineProfile'
import {
  Calendar, Clock, Ticket, CreditCard, CheckCircle2,
  AlertTriangle, XCircle, RotateCcw, Copy, Check, ChevronRight, X, ExternalLink,
  ClipboardList
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
  return d.toLocaleDateString('en-US', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function formatPlayDate(dStr) {
  if (!dStr) return '-'
  const d = new Date(dStr)
  return d.toLocaleDateString('en-US', {
    weekday: 'short',
    day: 'numeric',
    month: 'short',
    year: 'numeric'
  })
}

const isCancelling = ref(false)
const cancelConfirmBooking = ref(null)
const showCancelDialog = ref(false)

function promptCancel(booking) {
  cancelConfirmBooking.value = booking
  showCancelDialog.value = true
}

async function handleConfirmCancel() {
  if (!cancelConfirmBooking.value) return
  isCancelling.value = true
  await bookingStore.cancelBooking(cancelConfirmBooking.value.id, 'Cancelled by customer from booking history')
  isCancelling.value = false
  showCancelDialog.value = false
  cancelConfirmBooking.value = null
}
</script>

<template>
  <el-drawer
    v-model="bookingStore.isMyBookingsOpen"
    title="Booking History & Orders"
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
            <ClipboardList class="w-5 h-5" />
          </div>
          <div>
            <h2 class="text-sm font-extrabold text-slate-900 tracking-tight">My Bookings</h2>
            <p class="text-[11px] text-slate-500">
              Player: <strong class="text-slate-800">{{ authStore.displayName }}</strong>
              <span v-if="authStore.user?.phone" class="font-mono text-[10px] text-slate-400"> ({{ authStore.user.phone }})</span>
            </p>
          </div>
        </div>

        <button
          @click="bookingStore.closeMyBookings"
          class="p-1.5 text-slate-400 hover:text-slate-700 hover:bg-slate-100 rounded-xl transition cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Filter Tabs -->
      <div class="bg-white px-5 py-2.5 border-b border-slate-200/80 flex space-x-1.5 overflow-x-auto scrollbar-none">
        <button
          @click="activeFilter = 'ALL'"
          :class="[
            'px-3 py-1 rounded-full text-xs font-semibold whitespace-nowrap transition border cursor-pointer',
            activeFilter === 'ALL'
              ? 'bg-slate-900 text-white border-slate-900 shadow-2xs'
              : 'bg-slate-50 text-slate-600 border-slate-200 hover:bg-slate-100'
          ]"
        >
          All ({{ bookingStore.myBookings.length }})
        </button>

        <button
          @click="activeFilter = 'CONFIRMED'"
          :class="[
            'px-3 py-1 rounded-full text-xs font-semibold whitespace-nowrap transition border cursor-pointer',
            activeFilter === 'CONFIRMED'
              ? 'bg-emerald-600 text-white border-emerald-600 shadow-2xs'
              : 'bg-slate-50 text-slate-600 border-slate-200 hover:bg-slate-100'
          ]"
        >
          Confirmed
        </button>

        <button
          @click="activeFilter = 'PENDING_PAYMENT'"
          :class="[
            'px-3 py-1 rounded-full text-xs font-semibold whitespace-nowrap transition border cursor-pointer',
            activeFilter === 'PENDING_PAYMENT'
              ? 'bg-amber-600 text-white border-amber-600 shadow-2xs'
              : 'bg-slate-50 text-slate-600 border-slate-200 hover:bg-slate-100'
          ]"
        >
          Pending Payment
        </button>

        <button
          @click="activeFilter = 'MANUAL_REVIEW'"
          :class="[
            'px-3 py-1 rounded-full text-xs font-semibold whitespace-nowrap transition border cursor-pointer',
            activeFilter === 'MANUAL_REVIEW'
              ? 'bg-purple-600 text-white border-purple-600 shadow-2xs'
              : 'bg-slate-50 text-slate-600 border-slate-200 hover:bg-slate-100'
          ]"
        >
          Under Review
        </button>
      </div>

      <!-- Orders List Container -->
      <div class="flex-1 overflow-y-auto p-4 space-y-3">
        <!-- Loading Indicator -->
        <div v-if="bookingStore.isFetchingMyBookings" class="py-20 text-center text-slate-400 space-y-2">
          <div class="w-6 h-6 border-2 border-emerald-500 border-t-transparent rounded-full animate-spin mx-auto"></div>
          <p class="text-xs">Fetching your bookings...</p>
        </div>

        <!-- Empty State -->
        <div v-else-if="filteredBookings.length === 0" class="py-24 text-center space-y-3">
          <div class="w-14 h-14 bg-slate-100 rounded-3xl flex items-center justify-center mx-auto shadow-2xs text-slate-400">
            <Calendar class="w-7 h-7 text-slate-400" />
          </div>
          <div>
            <h3 class="text-sm font-extrabold text-slate-800">No Bookings Found</h3>
            <p class="text-xs text-slate-400 mt-1 max-w-[260px] mx-auto">
              You have no bookings under this status. Select a court slot to start playing!
            </p>
          </div>
          <button
            @click="bookingStore.closeMyBookings(); bookingStore.currentStep = 1"
            class="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-bold rounded-xl shadow-xs transition cursor-pointer"
          >
            Book a Court
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
                class="p-1 text-slate-400 hover:text-slate-700 rounded transition cursor-pointer"
                title="Copy booking ref"
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
                booking.status === 'CANCELLED' ? 'bg-rose-50 text-rose-700 border-rose-300' :
                'bg-slate-100 text-slate-500 border-slate-200'
              ]"
            >
              {{
                booking.status === 'CONFIRMED' ? '● CONFIRMED' :
                booking.status === 'PENDING_PAYMENT' ? 'PENDING_PAYMENT' :
                booking.status === 'MANUAL_REVIEW' ? 'UNDER_REVIEW' :
                booking.status
              }}
            </span>
          </div>

          <!-- Match / Play Details -->
          <div class="text-xs space-y-1.5 text-slate-700">
            <div class="flex items-center space-x-2 text-slate-500">
              <Calendar class="w-3.5 h-3.5 text-slate-400" />
              <span>Date: <strong class="text-slate-900">{{ formatPlayDate(booking.courts?.[0]?.play_date) }}</strong></span>
            </div>

            <div class="flex items-start space-x-2">
              <Clock class="w-3.5 h-3.5 text-slate-400 mt-0.5 flex-shrink-0" />
              <div class="flex flex-wrap gap-1">
                <span
                  v-for="c in booking.courts"
                  :key="c.id"
                  class="bg-slate-100 text-slate-800 border border-slate-200 px-1.5 py-0.2 rounded text-[10px] font-mono font-medium"
                >
                  Court {{ c.court?.court_number }} ({{ c.start_hour }}:00 - {{ c.end_hour }}:00)
                </span>
              </div>
            </div>

            <!-- Mini Store items if any -->
            <div v-if="booking.items && booking.items.length > 0" class="pt-1 text-[11px] text-slate-500">
              <span class="text-slate-400">Add-ons: </span>
              <span class="text-slate-700">{{ booking.items.map(it => `${it.product?.name || 'Item'} x${it.qty}`).join(', ') }}</span>
            </div>
          </div>

          <!-- Total & Actions -->
          <div class="flex items-center justify-between pt-2 border-t border-slate-100">
            <div>
              <div class="text-[10px] text-slate-400">Total</div>
              <div class="text-sm font-black font-mono text-emerald-700">
                ฿{{ Number(booking.net_total || 0).toLocaleString() }}
              </div>
            </div>

            <!-- Contextual Action Buttons -->
            <div class="flex items-center space-x-1.5">
              <!-- Cancel Booking Button (for PENDING_PAYMENT, CONFIRMED, or MANUAL_REVIEW) -->
              <button
                v-if="['PENDING_PAYMENT', 'CONFIRMED', 'MANUAL_REVIEW'].includes(booking.status)"
                @click="promptCancel(booking)"
                class="px-2.5 py-1.5 bg-slate-100 hover:bg-rose-50 text-slate-600 hover:text-rose-700 border border-slate-200 hover:border-rose-200 rounded-xl text-xs font-semibold transition cursor-pointer"
                title="Cancel booking"
              >
                Cancel
              </button>

              <!-- 1. If CONFIRMED: Open E-Ticket / Gate QR Code -->
              <button
                v-if="booking.status === 'CONFIRMED'"
                @click="bookingStore.viewBookingTicket(booking)"
                class="px-3 py-1.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-xs font-bold shadow-2xs transition flex items-center space-x-1.5 cursor-pointer"
              >
                <Ticket class="w-3.5 h-3.5" />
                <span>Pass</span>
              </button>

              <!-- 2. If PENDING_PAYMENT: Resume Payment -->
              <button
                v-else-if="booking.status === 'PENDING_PAYMENT'"
                @click="bookingStore.resumeBookingPayment(booking)"
                class="px-3 py-1.5 bg-amber-500 hover:bg-amber-600 text-white rounded-xl text-xs font-bold shadow-2xs transition flex items-center space-x-1.5 cursor-pointer"
              >
                <CreditCard class="w-3.5 h-3.5" />
                <span>Pay / Slip</span>
              </button>

              <!-- 3. If CANCELLED -->
              <span
                v-else-if="booking.status === 'CANCELLED'"
                class="text-[11px] text-rose-600 font-semibold"
              >
                Cancelled
              </span>

              <!-- 4. If MANUAL_REVIEW -->
              <span
                v-else-if="booking.status === 'MANUAL_REVIEW'"
                class="text-[11px] text-purple-700 font-semibold"
              >
                Under Review
              </span>

              <!-- 5. If EXPIRED -->
              <span
                v-else-if="booking.status === 'EXPIRED'"
                class="text-[11px] text-slate-400"
              >
                Expired
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Cancellation Confirmation Dialog -->
    <el-dialog
      v-model="showCancelDialog"
      title="Confirm Cancellation"
      width="90%"
      class="max-w-md rounded-2xl"
      :show-close="false"
      align-center
    >
      <div class="space-y-3 text-slate-700 text-xs">
        <p>
          Are you sure you want to cancel booking <strong class="text-slate-900 font-mono">{{ cancelConfirmBooking?.booking_no }}</strong>?
        </p>
        <div class="bg-amber-50 p-3 rounded-xl border border-amber-200 text-amber-900 text-[11px] leading-relaxed">
          <strong>Instant Release:</strong> Court slots will be unlocked and add-on inventory returned immediately in real-time.
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end space-x-2 pt-2">
          <button
            @click="showCancelDialog = false; cancelConfirmBooking = null"
            class="px-3 py-1.5 bg-slate-100 hover:bg-slate-200 rounded-xl text-xs font-semibold text-slate-700 transition cursor-pointer"
          >
            Back
          </button>
          <button
            @click="handleConfirmCancel"
            :disabled="isCancelling"
            class="px-3.5 py-1.5 bg-rose-600 hover:bg-rose-700 text-white rounded-xl text-xs font-bold shadow-xs transition cursor-pointer"
          >
            {{ isCancelling ? 'Cancelling...' : 'Confirm Cancellation' }}
          </button>
        </div>
      </template>
    </el-dialog>
  </el-drawer>
</template>
