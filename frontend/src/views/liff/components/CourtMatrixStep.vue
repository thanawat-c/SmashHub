<script setup>
import { ref, computed, onMounted } from 'vue'
import { useBookingStore } from '@/stores/bookingStore'
import {
  Calendar, Clock, ChevronRight, AlertCircle, Sparkles,
  LayoutGrid, ListFilter, Check, Zap, Shield
} from 'lucide-vue-next'

const bookingStore = useBookingStore()

// View Mode: 'slots' (Playtomic-style time-first list) or 'grid' (Wide court board)
const viewMode = ref('slots')

// Court type filter: 'ALL', 'PVC', 'ACRYLIC'
const courtFilter = ref('ALL')

onMounted(() => {
  bookingStore.fetchAvailability()
})

// Quick date chips
const dateChips = computed(() => {
  const dates = []
  const today = new Date()
  for (let i = 0; i < 5; i++) {
    const d = new Date(today)
    d.setDate(today.getDate() + i)
    const val = d.toISOString().split('T')[0]
    const weekday = d.toLocaleDateString('th-TH', { weekday: 'short' })
    const dayMonth = d.toLocaleDateString('th-TH', { day: 'numeric', month: 'short' })
    dates.push({ val, weekday, dayMonth, isToday: i === 0 })
  }
  return dates
})

function onDateChange(date) {
  bookingStore.selectedSlots = []
  bookingStore.fetchAvailability(date)
}

// Filtered courts based on active type filter
const filteredCourts = computed(() => {
  if (!bookingStore.courts) return []
  if (courtFilter.value === 'ALL') return bookingStore.courts
  return bookingStore.courts.filter(c => c.court_type === courtFilter.value)
})

const operatingHours = [10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21]

// Group slots by hour for the Playtomic-style Slots List
const hourlySlotsList = computed(() => {
  if (!bookingStore.courts) return []
  return operatingHours.map(hour => {
    const courtsForHour = filteredCourts.value.map(court => {
      const slot = court.slots.find(s => s.start_hour === hour) || {
        start_hour: hour,
        status: 'AVAILABLE',
        hourly_rate: hour >= 17 ? (court.court_type === 'ACRYLIC' ? 220 : 180) : (court.court_type === 'ACRYLIC' ? 180 : 150)
      }
      const isSelected = bookingStore.isSlotSelected(court.court_id, hour)
      return {
        court,
        slot,
        isSelected
      }
    })

    const availableCount = courtsForHour.filter(c => c.slot.status === 'AVAILABLE').length
    const isPeak = hour >= 17

    return {
      hour,
      timeLabel: `${String(hour).padStart(2, '0')}:00 - ${String(hour + 1).padStart(2, '0')}:00`,
      isPeak,
      availableCount,
      courts: courtsForHour
    }
  })
})

function getSlotColorClass(court, slot) {
  const isSelected = bookingStore.isSlotSelected(court.court_id, slot.start_hour)
  if (isSelected) {
    return 'bg-emerald-600 text-white font-bold border-emerald-700 shadow-xs scale-[0.98]'
  }

  switch (slot.status) {
    case 'AVAILABLE':
      return 'bg-white hover:bg-emerald-50 text-slate-800 border-slate-200 cursor-pointer transition hover:border-emerald-400 hover:shadow-2xs'
    case 'LOCKED_PENDING':
      return 'bg-amber-50 text-amber-900 border-amber-200 cursor-not-allowed font-medium'
    case 'BOOKED':
      return 'bg-slate-100 text-slate-400 border-slate-200 cursor-not-allowed'
    case 'MAINTENANCE':
      return 'bg-rose-50 text-rose-600 border-rose-200 cursor-not-allowed'
    default:
      return 'bg-slate-50 text-slate-500 border-slate-200'
  }
}
</script>

<template>
  <div class="flex flex-col h-full bg-[#F8FAFC] relative">
    <!-- Date & Filter Bar (Clean Compact Sticky Header) -->
    <div class="bg-white px-3 py-2 border-b border-slate-200/80 shadow-2xs sticky top-0 z-10 space-y-1.5">
      <!-- Date Carousel Header with Quick Day Selector -->
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-1.5 text-xs font-extrabold text-slate-800">
          <Calendar class="w-3.5 h-3.5 text-emerald-600" />
          <span>เลือกวันที่</span>
        </div>
        <!-- Court Filter Tabs (Compact) -->
        <div class="flex items-center space-x-1 bg-slate-100 p-0.5 rounded-lg text-[10px]">
          <button
            @click="courtFilter = 'ALL'"
            :class="[
              'px-2 py-0.5 rounded-md transition font-medium',
              courtFilter === 'ALL' ? 'bg-white text-slate-900 shadow-2xs font-bold' : 'text-slate-500'
            ]"
          >
            ทั้งหมด
          </button>
          <button
            @click="courtFilter = 'PVC'"
            :class="[
              'px-2 py-0.5 rounded-md transition font-medium',
              courtFilter === 'PVC' ? 'bg-white text-emerald-700 shadow-2xs font-bold' : 'text-slate-500'
            ]"
          >
            PVC (1-4)
          </button>
          <button
            @click="courtFilter = 'ACRYLIC'"
            :class="[
              'px-2 py-0.5 rounded-md transition font-medium',
              courtFilter === 'ACRYLIC' ? 'bg-white text-blue-700 shadow-2xs font-bold' : 'text-slate-500'
            ]"
          >
            VIP (5-6)
          </button>
        </div>
      </div>

      <!-- Sports Date Cards Carousel -->
      <div class="flex space-x-1.5 overflow-x-auto pb-0.5 scrollbar-none">
        <button
          v-for="chip in dateChips"
          :key="chip.val"
          @click="onDateChange(chip.val)"
          :class="[
            'px-2.5 py-1 rounded-xl text-center transition flex-shrink-0 border min-w-[58px]',
            bookingStore.selectedDate === chip.val
              ? 'bg-slate-900 text-white border-slate-900 shadow-2xs'
              : 'bg-slate-50 hover:bg-slate-100 text-slate-700 border-slate-200'
          ]"
        >
          <div :class="['text-[9px] font-semibold leading-tight', bookingStore.selectedDate === chip.val ? 'text-emerald-400' : 'text-slate-400']">
            {{ chip.weekday }}
          </div>
          <div class="text-[11px] font-extrabold font-mono mt-0.5">
            {{ chip.dayMonth }}
          </div>
        </button>
      </div>
    </div>

    <!-- Error Alert -->
    <div v-if="bookingStore.errorMessage" class="m-3 p-2.5 bg-rose-50 border border-rose-200 rounded-xl text-rose-700 text-xs flex items-center space-x-2">
      <AlertCircle class="w-4 h-4 flex-shrink-0" />
      <span>{{ bookingStore.errorMessage }}</span>
    </div>

    <!-- Main Content Area -->
    <div class="flex-1 overflow-auto p-3 pb-24">
      <div v-if="bookingStore.isLoading" class="flex flex-col items-center justify-center py-20 space-y-2 text-slate-400">
        <div class="w-6 h-6 border-2 border-emerald-500 border-t-transparent rounded-full animate-spin"></div>
        <span class="text-xs">กำลังโหลดตารางสนาม...</span>
      </div>

      <!-- VIEW A: Modern Time Slots List (Playtomic / MATCHi Benchmark) -->
      <div v-else class="space-y-2.5">
        <div
          v-for="hourGroup in hourlySlotsList"
          :key="hourGroup.hour"
          class="bg-white rounded-2xl border border-slate-200/80 p-3 shadow-2xs hover:border-slate-300 transition"
        >
          <!-- Hour Header -->
          <div class="flex items-center justify-between mb-2">
            <div class="flex items-center space-x-2">
              <span class="font-mono font-bold text-xs text-slate-900">{{ hourGroup.timeLabel }}</span>
              <span
                v-if="hourGroup.isPeak"
                class="px-1.5 py-0.2 rounded-full text-[9px] font-bold bg-amber-50 text-amber-700 border border-amber-200"
              >
                Peak
              </span>
            </div>
            <div class="text-[11px] text-slate-500">
              <span v-if="hourGroup.availableCount > 0" class="text-emerald-600 font-semibold">
                ว่าง {{ hourGroup.availableCount }} คอร์ท
              </span>
              <span v-else class="text-slate-400">
                เต็มทุกคอร์ท
              </span>
            </div>
          </div>

          <!-- Court Selection Pills (Thumb-friendly tap targets) -->
          <div class="grid grid-cols-2 sm:grid-cols-3 gap-1.5">
            <button
              v-for="item in hourGroup.courts"
              :key="item.court.court_id"
              @click="bookingStore.toggleCourtSlot(item.court, item.slot)"
              :disabled="item.slot.status !== 'AVAILABLE'"
              :class="[
                'p-2 rounded-xl border text-left transition select-none flex items-center justify-between',
                item.isSelected
                  ? 'bg-emerald-600 text-white font-bold border-emerald-700 shadow-xs'
                  : item.slot.status === 'AVAILABLE'
                  ? 'bg-slate-50/70 hover:bg-emerald-50/50 text-slate-800 border-slate-200 hover:border-emerald-400'
                  : item.slot.status === 'LOCKED_PENDING'
                  ? 'bg-amber-50/70 text-amber-800 border-amber-200/80 cursor-not-allowed opacity-80'
                  : 'bg-slate-100 text-slate-400 border-slate-200 cursor-not-allowed opacity-70'
              ]"
            >
              <div class="min-w-0 pr-1">
                <div class="text-xs font-bold leading-tight flex items-center space-x-1">
                  <span>คอร์ท {{ item.court.court_number }}</span>
                  <span
                    :class="[
                      'text-[8px] px-1 py-0.2 rounded-xs font-semibold',
                      item.isSelected ? 'bg-white/20 text-white' : (item.court.court_type === 'ACRYLIC' ? 'bg-blue-100 text-blue-800' : 'bg-slate-200 text-slate-700')
                    ]"
                  >
                    {{ item.court.court_type === 'ACRYLIC' ? 'VIP' : 'PVC' }}
                  </span>
                </div>
                <div :class="['text-[10px] font-mono mt-0.5', item.isSelected ? 'text-white/90' : 'text-slate-600']">
                  ฿{{ item.slot.hourly_rate }}
                </div>
              </div>

              <!-- Status Icon / Check -->
              <div class="flex-shrink-0">
                <div
                  v-if="item.isSelected"
                  class="w-4 h-4 rounded-full bg-white text-emerald-700 flex items-center justify-center text-[10px] font-bold"
                >
                  ✓
                </div>
                <div
                  v-else-if="item.slot.status === 'AVAILABLE'"
                  class="w-4 h-4 rounded-full border border-slate-300 flex items-center justify-center text-[9px] text-slate-400"
                >
                  +
                </div>
                <span v-else-if="item.slot.status === 'LOCKED_PENDING'" class="text-[9px] text-amber-700 font-medium">
                  รอโอน
                </span>
                <span v-else class="text-[9px] text-slate-400">
                  เต็ม
                </span>
              </div>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Floating Bottom Checkout Dock (Clean Light Mode Standard) -->
    <div class="absolute bottom-2 left-2 right-2 z-20">
      <div class="bg-white/95 backdrop-blur-md text-slate-900 p-3 rounded-2xl shadow-lg border border-slate-200/90">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-[10px] text-slate-500 font-medium">
              เลือกไว้ {{ bookingStore.selectedSlots.length }} ช่วงเวลา
            </div>
            <div class="text-base font-black text-slate-900 font-mono flex items-baseline space-x-1">
              <span class="text-emerald-700">฿{{ bookingStore.courtTotal.toLocaleString() }}</span>
              <span class="text-[10px] text-slate-400 font-normal font-sans">({{ bookingStore.selectedSlots.length }} ชม.)</span>
            </div>
          </div>

          <button
            @click="bookingStore.currentStep = 2"
            :disabled="bookingStore.selectedSlots.length === 0"
            :class="[
              'px-4 py-2.5 rounded-xl font-bold text-xs transition flex items-center space-x-1.5 shadow-xs active:scale-95',
              bookingStore.selectedSlots.length > 0
                ? 'bg-emerald-600 hover:bg-emerald-500 text-white cursor-pointer'
                : 'bg-slate-100 text-slate-400 border border-slate-200 cursor-not-allowed'
            ]"
          >
            <span>ถัดไป: สั่งสินค้า</span>
            <ChevronRight class="w-4 h-4" />
          </button>
        </div>

        <!-- Selected slots tags preview -->
        <div v-if="bookingStore.selectedSlots.length > 0" class="flex space-x-1 overflow-x-auto pt-2 mt-1.5 border-t border-slate-100 scrollbar-none">
          <span
            v-for="slot in bookingStore.selectedSlots"
            :key="slot.courtId + '-' + slot.startHour"
            class="inline-flex items-center space-x-1 px-2 py-0.5 bg-emerald-50 text-emerald-800 border border-emerald-200 rounded-full text-[10px] font-mono whitespace-nowrap"
          >
            <span>คอร์ท {{ slot.courtNumber }} ({{ slot.startHour }}:00)</span>
            <button @click="bookingStore.toggleCourtSlot({ court_id: slot.courtId }, { start_hour: slot.startHour, status: 'AVAILABLE' })" class="text-emerald-700 hover:text-slate-900 ml-0.5">×</button>
          </span>
        </div>
      </div>
    </div>
  </div>
</template>
