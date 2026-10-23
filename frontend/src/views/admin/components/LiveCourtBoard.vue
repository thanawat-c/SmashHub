<script setup>
import { ref, onMounted } from 'vue'
import { useAdminStore } from '@/stores/adminStore'
import api from '@/api'
import WalkInBookingModal from './WalkInBookingModal.vue'
import {
  Calendar, Clock, CheckCircle2, ShoppingCart, UserCheck, Plus, Eye, Printer, AlertCircle
} from 'lucide-vue-next'

const adminStore = useAdminStore()
const availability = ref(null)
const isLoading = ref(false)

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

// Find booking associated with this slot if any
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
</script>

<template>
  <div class="space-y-4">
    <!-- Board Controls Bar -->
    <div class="bg-white p-3.5 rounded-2xl border border-slate-200 shadow-xs flex flex-wrap items-center justify-between gap-3">
      <div class="flex items-center space-x-3">
        <div class="flex items-center space-x-1.5 text-xs font-bold text-slate-800">
          <Calendar class="w-4 h-4 text-emerald-600" />
          <span>วันที่แสดงกระดาน:</span>
        </div>
        <el-date-picker
          v-model="adminStore.filterDate"
          type="date"
          value-format="YYYY-MM-DD"
          size="small"
          placeholder="เลือกวันที่"
          @change="fetchBoard"
          class="!w-40"
        />
        <button
          @click="fetchBoard"
          class="px-2.5 py-1 text-xs bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-lg transition font-medium"
        >
          รีเฟรช
        </button>
      </div>

      <!-- Color Legend -->
      <div class="flex items-center space-x-3 text-xs">
        <span class="flex items-center space-x-1"><span class="w-3 h-3 rounded bg-white border border-slate-300"></span><span class="text-slate-600">ว่าง (คลิกเปิด Walk-in)</span></span>
        <span class="flex items-center space-x-1"><span class="w-3 h-3 rounded bg-emerald-600"></span><span class="text-slate-600">จองยืนยันแล้ว</span></span>
        <span class="flex items-center space-x-1"><span class="w-3 h-3 rounded bg-amber-400"></span><span class="text-slate-600">รอโอน (10 นาที)</span></span>
        <span class="flex items-center space-x-1"><span class="w-3 h-3 rounded bg-purple-500"></span><span class="text-slate-600">รอแอดมินตรวจสลิป</span></span>
      </div>
    </div>

    <!-- Live Grid Container -->
    <div class="bg-white rounded-2xl border border-slate-200 shadow-xs overflow-x-auto p-4">
      <div v-if="isLoading" class="py-20 text-center text-slate-400 text-xs">
        กำลังโหลดกระดานสด...
      </div>

      <div v-else-if="availability" class="min-w-[900px]">
        <!-- Court Header Row -->
        <div class="grid grid-cols-7 gap-2 pb-3 border-b border-slate-200">
          <div class="text-xs font-bold text-slate-400 flex items-center justify-center">
            ช่วงเวลา
          </div>
          <div
            v-for="court in availability.courts"
            :key="court.court_id"
            class="bg-slate-50 border border-slate-200 rounded-xl p-2.5 text-center"
          >
            <div class="font-extrabold text-slate-800 text-sm">คอร์ท {{ court.court_number }}</div>
            <div class="text-[10px] text-slate-500 font-medium">
              {{ court.court_type === 'ACRYLIC' ? 'Acrylic VIP' : 'PVC มาตรฐาน' }}
            </div>
          </div>
        </div>

        <!-- Hours Rows (10:00 - 22:00) -->
        <div class="divide-y divide-slate-100">
          <div
            v-for="hour in [10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21]"
            :key="hour"
            class="grid grid-cols-7 gap-2 py-2 items-stretch"
          >
            <!-- Time Column -->
            <div class="flex flex-col justify-center items-center text-slate-600 font-mono text-xs">
              <span class="font-bold">{{ String(hour).padStart(2, '0') }}:00</span>
              <span class="text-[10px] text-slate-400">{{ hour >= 17 ? 'Peak' : 'ปกติ' }}</span>
            </div>

            <!-- Court Slots -->
            <div
              v-for="court in availability.courts"
              :key="court.court_id + '-' + hour"
            >
              <template v-for="slot in court.slots.filter(s => s.start_hour === hour)" :key="slot.start_hour">
                <!-- Available Slot (Click to Walk-in) -->
                <div
                  v-if="slot.status === 'AVAILABLE'"
                  @click="handleSlotClick(court, slot)"
                  class="h-16 rounded-xl border border-dashed border-slate-300 hover:border-emerald-500 hover:bg-emerald-50/50 cursor-pointer p-2 flex flex-col justify-between transition group"
                >
                  <div class="flex justify-between items-center text-[10px] text-slate-400 group-hover:text-emerald-700">
                    <span>ว่าง</span>
                    <span class="font-mono">{{ slot.hourly_rate }}฿</span>
                  </div>
                  <div class="text-[10px] text-emerald-600 font-medium opacity-0 group-hover:opacity-100 flex items-center space-x-1">
                    <Plus class="w-3 h-3" />
                    <span>เปิด Walk-in</span>
                  </div>
                </div>

                <!-- Booked / Locked Slot -->
                <div
                  v-else
                  @click="handleSlotClick(court, slot)"
                  :class="[
                    'h-16 rounded-xl border p-2 flex flex-col justify-between transition cursor-pointer shadow-xs',
                    slot.status === 'LOCKED_PENDING'
                      ? 'bg-amber-50 border-amber-300 text-amber-900 hover:bg-amber-100'
                      : slot.status === 'BOOKED'
                      ? 'bg-emerald-50 border-emerald-300 text-emerald-950 hover:bg-emerald-100'
                      : 'bg-slate-100 border-slate-300 text-slate-500'
                  ]"
                >
                  <div class="flex justify-between items-start">
                    <div class="font-bold text-[11px] truncate flex items-center space-x-1">
                      <span class="truncate">{{ getSlotBooking(slot)?.user?.display_name || slot.booking_no || 'จองแล้ว' }}</span>
                    </div>
                    <span
                      :class="[
                        'text-[9px] font-bold px-1.5 py-0.2 rounded font-mono',
                        slot.status === 'LOCKED_PENDING' ? 'bg-amber-200 text-amber-900' : 'bg-emerald-200 text-emerald-900'
                      ]"
                    >
                      {{ slot.status === 'LOCKED_PENDING' ? 'รอโอน' : 'ยืนยัน' }}
                    </span>
                  </div>

                  <!-- Quick Action Buttons inside block -->
                  <div class="flex items-center justify-between text-[10px] pt-1">
                    <span class="font-mono text-slate-500">
                      {{ getSlotBooking(slot)?.booking_no?.slice(-4) || '' }}
                    </span>

                    <div v-if="getSlotBooking(slot)" class="flex space-x-1">
                      <!-- Quick Checkin Button -->
                      <button
                        @click="onQuickCheckIn(getSlotBooking(slot), $event)"
                        :title="getSlotBooking(slot).check_in_status === 'WAITING' ? 'กดเช็คอินเข้าสนาม' : 'เช็คอินแล้ว'"
                        :class="[
                          'px-1.5 py-0.5 rounded text-[9px] font-bold transition flex items-center space-x-0.5',
                          getSlotBooking(slot).check_in_status === 'WAITING'
                            ? 'bg-blue-600 text-white hover:bg-blue-700'
                            : 'bg-slate-200 text-slate-600'
                        ]"
                      >
                        <UserCheck class="w-3 h-3" />
                        <span>{{ getSlotBooking(slot).check_in_status === 'WAITING' ? 'Check-in' : 'เข้าแล้ว' }}</span>
                      </button>

                      <!-- Quick POS Add-Item Button -->
                      <button
                        @click="onQuickPOS(getSlotBooking(slot), $event)"
                        title="เพิ่มสินค้ามินิสโตร์เข้าบิล"
                        class="p-1 bg-amber-500 hover:bg-amber-600 text-white rounded text-[9px] transition"
                      >
                        <ShoppingCart class="w-3 h-3" />
                      </button>

                      <!-- Quick Receipt Print Button -->
                      <button
                        @click="onQuickReceipt(getSlotBooking(slot), $event)"
                        title="พิมพ์ใบเสร็จ"
                        class="p-1 bg-slate-700 hover:bg-slate-800 text-white rounded text-[9px] transition"
                      >
                        <Printer class="w-3 h-3" />
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

    <!-- WalkIn Modal -->
    <WalkInBookingModal
      v-model:visible="isWalkInOpen"
      :court="walkInCourt"
      :slot="walkInSlot"
      :play-date="adminStore.filterDate"
      @booking-created="fetchBoard"
    />
  </div>
</template>
