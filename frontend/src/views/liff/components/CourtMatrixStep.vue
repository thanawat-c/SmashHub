<script setup>
import { computed, onMounted } from 'vue'
import { useBookingStore } from '@/stores/bookingStore'
import { Calendar, Clock, ChevronRight, AlertCircle, Info } from 'lucide-vue-next'

const bookingStore = useBookingStore()

onMounted(() => {
  bookingStore.fetchAvailability()
})

// Quick date chips (Today, Tomorrow, Day after)
const dateChips = computed(() => {
  const dates = []
  const today = new Date()
  for (let i = 0; i < 5; i++) {
    const d = new Date(today)
    d.setDate(today.getDate() + i)
    const val = d.toISOString().split('T')[0]
    const label = i === 0 ? 'วันนี้' : i === 1 ? 'พรุ่งนี้' : d.toLocaleDateString('th-TH', { weekday: 'short', day: 'numeric', month: 'short' })
    dates.push({ val, label })
  }
  return dates
})

function onDateChange(date) {
  bookingStore.selectedSlots = []
  bookingStore.fetchAvailability(date)
}

function getSlotColorClass(court, slot) {
  const isSelected = bookingStore.isSlotSelected(court.court_id, slot.start_hour)
  if (isSelected) {
    return 'bg-emerald-600 text-white border-emerald-700 shadow-sm font-semibold ring-2 ring-emerald-400/50 scale-[0.98]'
  }

  switch (slot.status) {
    case 'AVAILABLE':
      return 'bg-white hover:bg-emerald-50 text-slate-700 border-slate-200 cursor-pointer transition hover:border-emerald-300'
    case 'LOCKED_PENDING':
      return 'bg-amber-100 text-amber-800 border-amber-300 cursor-not-allowed opacity-90'
    case 'BOOKED':
      return 'bg-slate-200 text-slate-500 border-slate-300 cursor-not-allowed opacity-80'
    case 'MAINTENANCE':
      return 'bg-rose-100 text-rose-700 border-rose-300 cursor-not-allowed'
    default:
      return 'bg-slate-100 text-slate-600 border-slate-200'
  }
}
</script>

<template>
  <div class="flex flex-col h-full bg-slate-50">
    <!-- Date Selector Bar -->
    <div class="bg-white px-3 py-2.5 border-b border-slate-200 shadow-xs sticky top-0 z-10">
      <div class="flex items-center justify-between mb-2">
        <div class="flex items-center space-x-1.5 text-xs font-semibold text-slate-700">
          <Calendar class="w-4 h-4 text-emerald-600" />
          <span>เลือกวันที่ต้องการตี</span>
        </div>
        <span class="text-[11px] text-slate-400 font-mono">{{ bookingStore.selectedDate }}</span>
      </div>

      <!-- Quick Date Chips -->
      <div class="flex space-x-1.5 overflow-x-auto pb-1 scrollbar-none">
        <button
          v-for="chip in dateChips"
          :key="chip.val"
          @click="onDateChange(chip.val)"
          :class="[
            'px-2.5 py-1 rounded-full text-xs font-medium whitespace-nowrap transition flex-shrink-0',
            bookingStore.selectedDate === chip.val
              ? 'bg-emerald-600 text-white shadow-xs'
              : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
          ]"
        >
          {{ chip.label }}
        </button>
      </div>
    </div>

    <!-- Rate Info & Legend Bar -->
    <div class="bg-emerald-50/70 border-b border-emerald-100 px-3 py-1.5 text-[11px] text-emerald-800 flex items-center justify-between">
      <div class="flex items-center space-x-1 truncate">
        <Clock class="w-3.5 h-3.5 text-emerald-600 flex-shrink-0" />
        <span class="truncate">Peak Hour (หลัง 17:00 น.) ปรับเรทอัตโนมัติ</span>
      </div>
      <div class="flex items-center space-x-2 text-[10px] flex-shrink-0 font-medium">
        <span class="flex items-center space-x-1"><span class="w-2 h-2 rounded-full bg-white border border-slate-300"></span><span>ว่าง</span></span>
        <span class="flex items-center space-x-1"><span class="w-2 h-2 rounded-full bg-emerald-600"></span><span>เลือก</span></span>
        <span class="flex items-center space-x-1"><span class="w-2 h-2 rounded-full bg-amber-400"></span><span>รอโอน</span></span>
        <span class="flex items-center space-x-1"><span class="w-2 h-2 rounded-full bg-slate-400"></span><span>จองแล้ว</span></span>
      </div>
    </div>

    <!-- Error Alert if any -->
    <div v-if="bookingStore.errorMessage" class="m-3 p-2.5 bg-rose-50 border border-rose-200 rounded-lg text-rose-700 text-xs flex items-center space-x-2">
      <AlertCircle class="w-4 h-4 flex-shrink-0" />
      <span>{{ bookingStore.errorMessage }}</span>
    </div>

    <!-- Matrix Table Container -->
    <div class="flex-1 overflow-auto p-2">
      <div v-if="bookingStore.isLoading" class="flex flex-col items-center justify-center py-16 space-y-2 text-slate-400">
        <div class="w-7 h-7 border-3 border-emerald-500 border-t-transparent rounded-full animate-spin"></div>
        <span class="text-xs">กำลังโหลดตารางสนาม...</span>
      </div>

      <div v-else class="min-w-[370px]">
        <!-- Court Headers (X-Axis) -->
        <div class="grid grid-cols-7 gap-1 sticky top-0 bg-slate-50 pb-1 z-5">
          <div class="text-[10px] font-semibold text-slate-400 flex items-center justify-center">
            เวลา
          </div>
          <div
            v-for="court in bookingStore.courts"
            :key="court.court_id"
            class="bg-white border border-slate-200 rounded-md p-1 text-center shadow-2xs"
          >
            <div class="text-[11px] font-bold text-slate-800 leading-tight">คอร์ท {{ court.court_number }}</div>
            <div class="text-[9px] text-slate-500 font-medium truncate">
              {{ court.court_type === 'ACRYLIC' ? 'Acrylic VIP' : 'PVC' }}
            </div>
          </div>
        </div>

        <!-- Hourly Rows (Y-Axis) -->
        <div class="space-y-1 mt-1">
          <div
            v-for="hour in [10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21]"
            :key="hour"
            class="grid grid-cols-7 gap-1 items-center"
          >
            <!-- Time Label -->
            <div class="text-[10px] font-mono font-semibold text-slate-500 text-center leading-none">
              {{ String(hour).padStart(2, '0') }}:00
              <div class="text-[8px] text-slate-400 font-normal">{{ hour >= 17 ? 'Peak' : 'ปกติ' }}</div>
            </div>

            <!-- Court Slots for this Hour -->
            <div
              v-for="court in bookingStore.courts"
              :key="court.court_id + '-' + hour"
              class="relative"
            >
              <template v-for="slot in court.slots.filter(s => s.start_hour === hour)" :key="slot.start_hour">
                <button
                  @click="bookingStore.toggleCourtSlot(court, slot)"
                  :disabled="slot.status !== 'AVAILABLE'"
                  :class="[
                    'w-full h-11 rounded-md border flex flex-col items-center justify-center text-[10px] p-0.5 transition select-none',
                    getSlotColorClass(court, slot)
                  ]"
                >
                  <span class="font-bold leading-tight">
                    {{ bookingStore.isSlotSelected(court.court_id, hour) ? '✓ เลือก' : `${slot.hourly_rate}฿` }}
                  </span>
                  <span class="text-[8px] opacity-80 scale-90 truncate max-w-full">
                    {{ slot.status === 'AVAILABLE' ? 'ว่าง' : slot.status === 'LOCKED_PENDING' ? 'รอโอน' : slot.status === 'BOOKED' ? 'เต็ม' : 'ปิด' }}
                  </span>
                </button>
              </template>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Action Bar (Fixed) -->
    <div class="bg-white border-t border-slate-200 p-3 shadow-lg sticky bottom-0 z-20">
      <div class="flex items-center justify-between mb-2">
        <div>
          <div class="text-[11px] text-slate-500">เลือกแล้ว {{ bookingStore.selectedSlots.length }} ช่วงเวลา</div>
          <div class="text-base font-bold text-emerald-700">
            ฿{{ bookingStore.courtTotal.toLocaleString() }}
            <span class="text-[10px] text-slate-400 font-normal">({{ bookingStore.selectedSlots.length }} ชม.)</span>
          </div>
        </div>

        <button
          @click="bookingStore.currentStep = 2"
          :disabled="bookingStore.selectedSlots.length === 0"
          :class="[
            'px-4 py-2.5 rounded-xl font-semibold text-xs transition flex items-center space-x-1.5 shadow',
            bookingStore.selectedSlots.length > 0
              ? 'bg-emerald-600 hover:bg-emerald-700 text-white shadow-emerald-500/20 active:scale-95'
              : 'bg-slate-200 text-slate-400 cursor-not-allowed'
          ]"
        >
          <span>ถัดไป: สั่งสินค้า</span>
          <ChevronRight class="w-4 h-4" />
        </button>
      </div>

      <!-- Selected Slot Chips Preview -->
      <div v-if="bookingStore.selectedSlots.length > 0" class="flex space-x-1 overflow-x-auto pt-1 scrollbar-none">
        <span
          v-for="slot in bookingStore.selectedSlots"
          :key="slot.courtId + '-' + slot.startHour"
          class="inline-flex items-center space-x-1 px-2 py-0.5 bg-emerald-100 text-emerald-800 rounded-full text-[10px] font-medium whitespace-nowrap"
        >
          <span>คอร์ท {{ slot.courtNumber }} ({{ slot.startHour }}:00)</span>
          <button @click="bookingStore.toggleCourtSlot({ court_id: slot.courtId }, { start_hour: slot.startHour, status: 'AVAILABLE' })" class="text-emerald-600 hover:text-emerald-900 ml-0.5">×</button>
        </span>
      </div>
    </div>
  </div>
</template>
