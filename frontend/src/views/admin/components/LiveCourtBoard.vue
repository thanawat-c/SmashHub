<script setup>
import { ref, onMounted, computed } from 'vue'
import { useAdminStore } from '@/stores/adminStore'
import api from '@/api'
import WalkInBookingModal from './WalkInBookingModal.vue'
import {
  Calendar, Clock, CheckCircle2, ShoppingCart, UserCheck, Plus, Eye, Printer, AlertCircle, Sparkles
} from 'lucide-vue-next'

const adminStore = useAdminStore()
const availability = ref(null)
const isLoading = ref(false)

// Current time marker
const now = new Date()
const currentHour = now.getHours() // 0 - 23

// WalkIn modal state
const isWalkInOpen = ref(false)
const walkInCourt = ref(null)
const walkInSlot = ref(null)

onMounted(() => {
  fetchBoard()
})

async function fetchBoard() {
  isLoading.value = true
  try {
    const res = await api.get(`/courts/availability?date=${adminStore.filterDate}`)
    availability.value = res.data
  } catch (err) {
    console.error('Error fetching live court board:', err)
  } finally {
    isLoading.value = false
  }
}

function getSlotBooking(slot) {
  if (!slot.booking_no) return null
  return adminStore.bookings.find(b => b.booking_no === slot.booking_no)
}

function handleSlotClick(court, slot) {
  if (slot.status === 'AVAILABLE') {
    walkInCourt.value = court
    walkInSlot.value = slot
    isWalkInOpen.value = true
  } else if (slot.booking_no) {
    const booking = getSlotBooking(slot)
    if (booking) {
      adminStore.openDrawer(booking)
    }
  }
}

function onQuickCheckIn(booking, e) {
  e.stopPropagation()
  const nextStatus = booking.check_in_status === 'WAITING' ? 'CHECKED_IN' : 'COMPLETED'
  adminStore.updateCheckIn(booking.id, nextStatus)
}

function onQuickPOS(booking, e) {
  e.stopPropagation()
  adminStore.openPOS(booking)
}

function onQuickReceipt(booking, e) {
  e.stopPropagation()
  adminStore.openReceipt(booking)
}

// 1-Click Fast item add directly from court card (e.g. 1 water)
async function onFastAddWater(booking, e) {
  e.stopPropagation()
  // Product ID 4 is Singha water
  await adminStore.addItems(booking.id, [{ product_id: 4, qty: 1 }])
}
</script>

<template>
  <div class="space-y-3.5">
    <!-- Controls & Quick Filters -->
    <div class="bg-white p-3.5 rounded-2xl border border-slate-200/90 shadow-2xs flex flex-wrap items-center justify-between gap-3">
      <div class="flex items-center space-x-3">
        <div class="flex items-center space-x-1.5 text-xs font-bold text-slate-800">
          <Calendar class="w-4 h-4 text-emerald-600" />
          <span>วันที่:</span>
        </div>
        <el-date-picker
          v-model="adminStore.filterDate"
          type="date"
          value-format="YYYY-MM-DD"
          size="small"
          placeholder="เลือกวันที่"
          @change="fetchBoard"
          class="!w-36"
        />
        <button
          @click="fetchBoard"
          class="px-3 py-1.5 text-xs bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl transition font-medium"
        >
          รีเฟรช
        </button>
      </div>

      <!-- Live Court Color Legend -->
      <div class="flex items-center space-x-3 text-xs">
        <span class="flex items-center space-x-1.5"><span class="w-2.5 h-2.5 rounded-sm bg-white border border-slate-300"></span><span class="text-slate-600">ว่าง (คลิกเปิด Walk-in)</span></span>
        <span class="flex items-center space-x-1.5"><span class="w-2.5 h-2.5 rounded-sm bg-emerald-500"></span><span class="text-slate-600">ยืนยันแล้ว</span></span>
        <span class="flex items-center space-x-1.5"><span class="w-2.5 h-2.5 rounded-sm bg-amber-400"></span><span class="text-slate-600">รอโอน (10 นาที)</span></span>
        <span class="flex items-center space-x-1.5"><span class="w-2.5 h-2.5 rounded-sm bg-purple-500"></span><span class="text-slate-600">รอตรวจสลิป</span></span>
      </div>
    </div>

    <!-- Live Court Matrix Container -->
    <div class="bg-white rounded-2xl border border-slate-200/90 shadow-2xs overflow-x-auto p-4 relative">
      <div v-if="isLoading" class="py-24 text-center text-slate-400 text-xs">
        กำลังโหลดกระดานสด...
      </div>

      <div v-else-if="availability" class="min-w-[880px] relative">
        <!-- Court Header Row -->
        <div class="grid grid-cols-7 gap-2 pb-2.5 border-b border-slate-200">
          <div class="text-[11px] font-mono font-bold text-slate-400 flex items-center justify-center">
            TIME
          </div>
          <div
            v-for="court in availability.courts"
            :key="court.court_id"
            class="bg-slate-50 border border-slate-200 rounded-xl p-2 text-center"
          >
            <div class="font-black text-slate-900 text-xs">คอร์ท {{ court.court_number }}</div>
            <div class="text-[9px] text-emerald-700 font-semibold truncate">
              {{ court.court_type === 'ACRYLIC' ? 'Acrylic VIP' : 'PVC Mat' }}
            </div>
          </div>
        </div>

        <!-- Hours Rows (10:00 - 22:00) -->
        <div class="divide-y divide-slate-100 relative">
          <div
            v-for="hour in [10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21]"
            :key="hour"
            class="grid grid-cols-7 gap-2 py-1.5 items-stretch relative"
          >
            <!-- Current Time Red Laser Indicator (Active if current hour matches) -->
            <div
              v-if="currentHour === hour"
              class="absolute -top-1 left-0 right-0 z-10 pointer-events-none flex items-center"
            >
              <span class="bg-rose-500 text-white font-mono text-[8px] font-bold px-1.5 py-0.5 rounded-sm shadow-xs flex items-center space-x-1">
                <span class="w-1.5 h-1.5 rounded-full bg-white animate-ping"></span>
                <span>LIVE NOW</span>
              </span>
              <div class="flex-1 border-t-2 border-rose-500 shadow-sm"></div>
            </div>

            <!-- Time Column Label -->
            <div class="flex flex-col justify-center items-center text-slate-600 font-mono text-xs">
              <span class="font-bold">{{ String(hour).padStart(2, '0') }}:00</span>
              <span :class="['text-[8px]', hour >= 17 ? 'text-amber-600 font-bold' : 'text-slate-400']">
                {{ hour >= 17 ? 'Peak' : 'ปกติ' }}
              </span>
            </div>

            <!-- Court Slots for this Hour -->
            <div
              v-for="court in availability.courts"
              :key="court.court_id + '-' + hour"
            >
              <template v-for="slot in court.slots.filter(s => s.start_hour === hour)" :key="slot.start_hour">
                <!-- Available Slot: Clean dashed border, click to walk-in -->
                <div
                  v-if="slot.status === 'AVAILABLE'"
                  @click="handleSlotClick(court, slot)"
                  class="h-14 rounded-xl border border-dashed border-slate-300 hover:border-emerald-500 hover:bg-emerald-50/40 cursor-pointer p-2 flex flex-col justify-between transition group"
                >
                  <div class="flex justify-between items-center text-[10px] text-slate-400 group-hover:text-emerald-700">
                    <span>ว่าง</span>
                    <span class="font-mono">{{ slot.hourly_rate }}฿</span>
                  </div>
                  <div class="text-[9px] text-emerald-600 font-semibold opacity-0 group-hover:opacity-100 flex items-center space-x-1 transition">
                    <Plus class="w-3 h-3" />
                    <span>เปิด Walk-in</span>
                  </div>
                </div>

                <!-- Booked / Locked Slot -->
                <div
                  v-else
                  @click="handleSlotClick(court, slot)"
                  :class="[
                    'h-14 rounded-xl border p-2 flex flex-col justify-between transition cursor-pointer shadow-2xs',
                    slot.status === 'LOCKED_PENDING'
                      ? 'bg-amber-50/90 border-amber-300 text-amber-950 hover:bg-amber-100'
                      : slot.status === 'BOOKED'
                      ? 'bg-emerald-50/90 border-emerald-300 text-emerald-950 hover:bg-emerald-100'
                      : 'bg-slate-100 border-slate-300 text-slate-500'
                  ]"
                >
                  <div class="flex justify-between items-start">
                    <div class="font-bold text-[11px] truncate flex items-center space-x-1">
                      <span class="truncate">{{ getSlotBooking(slot)?.user?.display_name || slot.booking_no || 'จองแล้ว' }}</span>
                    </div>
                    <span
                      :class="[
                        'text-[8px] font-bold px-1.5 py-0.2 rounded font-mono',
                        slot.status === 'LOCKED_PENDING' ? 'bg-amber-200 text-amber-900' : 'bg-emerald-200 text-emerald-900'
                      ]"
                    >
                      {{ slot.status === 'LOCKED_PENDING' ? 'รอโอน' : 'ยืนยัน' }}
                    </span>
                  </div>

                  <!-- Quick Action Pills on Card -->
                  <div class="flex items-center justify-between text-[9px] pt-0.5">
                    <span class="font-mono text-slate-500 text-[10px]">
                      {{ getSlotBooking(slot)?.booking_no?.slice(-4) || '' }}
                    </span>

                    <div v-if="getSlotBooking(slot)" class="flex space-x-1">
                      <!-- Quick Check-in -->
                      <button
                        @click="onQuickCheckIn(getSlotBooking(slot), $event)"
                        :class="[
                          'px-1.5 py-0.5 rounded text-[9px] font-bold transition flex items-center space-x-0.5',
                          getSlotBooking(slot).check_in_status === 'WAITING'
                            ? 'bg-blue-600 text-white hover:bg-blue-700'
                            : 'bg-emerald-600 text-white'
                        ]"
                      >
                        <UserCheck class="w-2.5 h-2.5" />
                        <span>{{ getSlotBooking(slot).check_in_status === 'WAITING' ? 'Check-in' : 'เข้าแล้ว' }}</span>
                      </button>

                      <!-- Quick Add Water (+15฿) -->
                      <button
                        @click="onFastAddWater(getSlotBooking(slot), $event)"
                        title="เพิ่มน้ำดื่มสิงห์ด่วน 1 ขวด"
                        class="px-1 py-0.5 bg-sky-100 hover:bg-sky-200 text-sky-800 rounded font-semibold text-[8px]"
                      >
                        +🥤
                      </button>

                      <!-- POS Modal -->
                      <button
                        @click="onQuickPOS(getSlotBooking(slot), $event)"
                        title="เปิดเคาน์เตอร์ POS"
                        class="p-1 bg-amber-500 hover:bg-amber-600 text-white rounded transition"
                      >
                        <ShoppingCart class="w-2.5 h-2.5" />
                      </button>

                      <!-- Receipt Print -->
                      <button
                        @click="onQuickReceipt(getSlotBooking(slot), $event)"
                        title="พิมพ์ใบเสร็จ"
                        class="p-1 bg-slate-700 hover:bg-slate-800 text-white rounded transition"
                      >
                        <Printer class="w-2.5 h-2.5" />
                      </button>
                    </div>
                  </div>
                </div>
              </template>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Walk-in Booking Modal -->
    <WalkInBookingModal
      v-model:visible="isWalkInOpen"
      :court="walkInCourt"
      :slot="walkInSlot"
      :play-date="adminStore.filterDate"
      @booking-created="fetchBoard"
    />
  </div>
</template>
