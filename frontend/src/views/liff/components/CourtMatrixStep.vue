<script setup>
import { ref, computed, onMounted } from 'vue'
import { useBookingStore } from '@/stores/bookingStore'
import {
  Calendar, Clock, ChevronRight, AlertCircle,
  Check, Info, Flame, Sparkles
} from 'lucide-vue-next'

const bookingStore = useBookingStore()

// Selected court tab (Court-First flow like Playtomic / MATCHi)
const activeCourtId = ref(null)

// Court type filter: 'ALL', 'PVC', 'ACRYLIC'
const courtFilter = ref('ALL')

onMounted(async () => {
  await bookingStore.fetchAvailability()
  if (bookingStore.courts && bookingStore.courts.length > 0) {
    activeCourtId.value = bookingStore.courts[0].court_id
  }
})

// Quick date chips
const dateChips = computed(() => {
  const dates = []
  const today = new Date()
  for (let i = 0; i < 7; i++) {
    const d = new Date(today)
    d.setDate(today.getDate() + i)
    const val = d.toISOString().split('T')[0]
    const weekday = d.toLocaleDateString('en-US', { weekday: 'short' })
    const dayMonth = d.toLocaleDateString('en-US', { day: 'numeric', month: 'short' })
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

// Active Court Object
const activeCourt = computed(() => {
  if (!bookingStore.courts || bookingStore.courts.length === 0) return null
  return bookingStore.courts.find(c => c.court_id === activeCourtId.value) || bookingStore.courts[0]
})

// Operating hours 10:00 to 22:00
const operatingHours = [10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21]

// Compute timeline slots for the active court
const activeCourtSlots = computed(() => {
  if (!activeCourt.value) return []
  return operatingHours.map(hour => {
    const slot = activeCourt.value.slots?.find(s => s.start_hour === hour) || {
      start_hour: hour,
      status: 'AVAILABLE',
      hourly_rate: hour >= 17 ? (activeCourt.value.court_type === 'ACRYLIC' ? 220 : 180) : (activeCourt.value.court_type === 'ACRYLIC' ? 180 : 150)
    }
    const isSelected = bookingStore.isSlotSelected(activeCourt.value.court_id, hour)
    return {
      hour,
      timeLabel: `${String(hour).padStart(2, '0')}:00 - ${String(hour + 1).padStart(2, '0')}:00`,
      isPeak: hour >= 17,
      slot,
      isSelected
    }
  })
})

function selectCourt(court) {
  activeCourtId.value = court.court_id
}

// Count how many slots are selected in this court
function getSelectedCountForCourt(courtId) {
  return bookingStore.selectedSlots.filter(s => s.courtId === courtId).length
}

// Total available count for a court
function getAvailableCountForCourt(court) {
  if (!court.slots) return 0
  return court.slots.filter(s => s.status === 'AVAILABLE').length
}
</script>

<template>
  <div class="flex flex-col h-full bg-[#FAFAFA] relative">
    <!-- Clean Minimalist Date Selector (Playtomic / MATCHi style) -->
    <div class="bg-white border-b border-slate-200/90 px-3.5 pt-3 pb-2.5 flex-shrink-0">
      <div class="flex items-center justify-between mb-2">
        <div class="flex items-center space-x-1.5 text-xs font-bold text-slate-900">
          <Calendar class="w-3.5 h-3.5 text-emerald-600" />
          <span>Select Date</span>
        </div>
        <span class="text-[10px] text-slate-400 font-medium">Daily 10:00 - 22:00</span>
      </div>

      <!-- Horizontal Date Bar -->
      <div class="flex space-x-2 overflow-x-auto pb-1 scrollbar-none">
        <button
          v-for="chip in dateChips"
          :key="chip.val"
          @click="onDateChange(chip.val)"
          :class="[
            'flex flex-col items-center justify-center px-3 py-1.5 rounded-xl text-center transition flex-shrink-0 min-w-[56px] border cursor-pointer',
            bookingStore.selectedDate === chip.val
              ? 'bg-slate-900 text-white border-slate-900 shadow-sm'
              : 'bg-white hover:bg-slate-50 text-slate-700 border-slate-200'
          ]"
        >
          <span :class="['text-[10px] font-semibold', bookingStore.selectedDate === chip.val ? 'text-emerald-400' : 'text-slate-400']">
            {{ chip.weekday }}
          </span>
          <span class="text-xs font-extrabold font-mono mt-0.5">
            {{ chip.dayMonth }}
          </span>
        </button>
      </div>
    </div>

    <!-- Court Selection Carousel (MATCHi Club Court Selector) -->
    <div class="bg-white border-b border-slate-200/90 px-3.5 py-2.5 flex-shrink-0">
      <div class="flex items-center justify-between mb-2">
        <span class="text-[11px] font-bold text-slate-500 uppercase tracking-wider">Select Court</span>
        <div class="flex items-center space-x-1 text-[10px]">
          <button
            @click="courtFilter = 'ALL'"
            :class="['px-2 py-0.5 rounded-md font-semibold transition cursor-pointer', courtFilter === 'ALL' ? 'bg-slate-900 text-white' : 'text-slate-500 hover:bg-slate-100']"
          >
            All
          </button>
          <button
            @click="courtFilter = 'PVC'"
            :class="['px-2 py-0.5 rounded-md font-semibold transition cursor-pointer', courtFilter === 'PVC' ? 'bg-emerald-700 text-white' : 'text-slate-500 hover:bg-slate-100']"
          >
            PVC
          </button>
          <button
            @click="courtFilter = 'ACRYLIC'"
            :class="['px-2 py-0.5 rounded-md font-semibold transition cursor-pointer', courtFilter === 'ACRYLIC' ? 'bg-blue-700 text-white' : 'text-slate-500 hover:bg-slate-100']"
          >
            VIP
          </button>
        </div>
      </div>

      <!-- Court Cards Strip -->
      <div class="flex space-x-2 overflow-x-auto pb-1 scrollbar-none">
        <button
          v-for="court in filteredCourts"
          :key="court.court_id"
          @click="selectCourt(court)"
          :class="[
            'flex-shrink-0 px-3 py-2 rounded-xl text-left border transition relative min-w-[110px] cursor-pointer',
            activeCourt?.court_id === court.court_id
              ? 'bg-emerald-50/80 border-emerald-600 text-emerald-950 shadow-2xs ring-1 ring-emerald-600/30'
              : 'bg-white hover:bg-slate-50 border-slate-200 text-slate-800'
          ]"
        >
          <!-- Badge counter of selected slots in this court -->
          <span
            v-if="getSelectedCountForCourt(court.court_id) > 0"
            class="absolute -top-1 -right-1 bg-emerald-600 text-white text-[9px] font-bold w-4 h-4 rounded-full flex items-center justify-center shadow-xs"
          >
            {{ getSelectedCountForCourt(court.court_id) }}
          </span>

          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-slate-900">Court {{ court.court_number }}</span>
            <span
              :class="[
                'text-[8px] font-bold px-1.5 py-0.2 rounded-xs',
                court.court_type === 'ACRYLIC' ? 'bg-blue-100 text-blue-700' : 'bg-slate-100 text-slate-600'
              ]"
            >
              {{ court.court_type === 'ACRYLIC' ? 'VIP' : 'PVC' }}
            </span>
          </div>

          <div class="text-[10px] text-slate-400 mt-1 flex items-center justify-between font-mono">
            <span>{{ getAvailableCountForCourt(court) }}h open</span>
            <span class="font-semibold text-slate-600">฿{{ court.court_type === 'ACRYLIC' ? '180' : '150' }}</span>
          </div>
        </button>
      </div>
    </div>

    <!-- Error Alert -->
    <div v-if="bookingStore.errorMessage" class="m-3 p-2.5 bg-rose-50 border border-rose-200 rounded-xl text-rose-700 text-xs flex items-center space-x-2">
      <AlertCircle class="w-4 h-4 flex-shrink-0" />
      <span>{{ bookingStore.errorMessage }}</span>
    </div>

    <!-- Active Court Schedule Timeline (MATCHi / Playtomic Slot Grid) -->
    <div class="flex-1 overflow-y-auto p-3.5 pb-28 lg:pb-6">
      <div v-if="bookingStore.isLoading" class="flex flex-col items-center justify-center py-20 space-y-2 text-slate-400">
        <div class="w-6 h-6 border-2 border-emerald-500 border-t-transparent rounded-full animate-spin"></div>
        <span class="text-xs">Loading availability...</span>
      </div>

      <div v-else-if="activeCourt" class="space-y-3">
        <!-- Court Header Banner -->
        <div class="bg-white rounded-2xl border border-slate-200/90 p-3 flex items-center justify-between shadow-2xs">
          <div>
            <div class="flex items-center space-x-2">
              <h3 class="text-sm font-extrabold text-slate-900">Court {{ activeCourt.court_number }} ({{ activeCourt.court_name }})</h3>
              <span
                :class="[
                  'text-[9px] font-bold px-1.5 py-0.2 rounded-md',
                  activeCourt.court_type === 'ACRYLIC' ? 'bg-blue-50 text-blue-700 border border-blue-200' : 'bg-emerald-50 text-emerald-700 border border-emerald-200'
                ]"
              >
                {{ activeCourt.court_type === 'ACRYLIC' ? 'VIP Acrylic Floor' : 'PVC Competition Mat' }}
              </span>
            </div>
            <p class="text-[11px] text-slate-500 mt-0.5">Select your desired playing hours (multi-hour selection supported)</p>
          </div>
          <div class="text-right">
            <span class="text-[10px] text-slate-400 block">Rate</span>
            <span class="text-xs font-bold font-mono text-slate-800">฿{{ activeCourt.court_type === 'ACRYLIC' ? '180-220' : '150-180' }}/h</span>
          </div>
        </div>

        <!-- Hourly Schedule Grid -->
        <div class="bg-white rounded-2xl border border-slate-200/90 overflow-hidden shadow-2xs divide-y divide-slate-100">
          <div
            v-for="item in activeCourtSlots"
            :key="item.hour"
            @click="item.slot.status === 'AVAILABLE' && bookingStore.toggleCourtSlot(activeCourt, item.slot)"
            :class="[
              'px-3.5 py-2.5 flex items-center justify-between transition select-none',
              item.isSelected
                ? 'bg-emerald-50 text-emerald-950 font-semibold'
                : item.slot.status === 'AVAILABLE'
                ? 'hover:bg-slate-50 cursor-pointer text-slate-800'
                : item.slot.status === 'LOCKED_PENDING'
                ? 'bg-amber-50/50 text-amber-900/70 cursor-not-allowed opacity-75'
                : 'bg-slate-50/70 text-slate-400 cursor-not-allowed opacity-60'
            ]"
          >
            <!-- Left: Time + Peak Badge -->
            <div class="flex items-center space-x-2.5">
              <!-- Checkbox circle -->
              <div
                :class="[
                  'w-5 h-5 rounded-md flex items-center justify-center transition border text-[11px]',
                  item.isSelected
                    ? 'bg-emerald-600 border-emerald-600 text-white shadow-2xs'
                    : item.slot.status === 'AVAILABLE'
                    ? 'border-slate-300 bg-white text-transparent'
                    : 'border-slate-200 bg-slate-100 text-transparent'
                ]"
              >
                <Check v-if="item.isSelected" class="w-3.5 h-3.5 stroke-[3]" />
              </div>

              <div>
                <div class="flex items-center space-x-2">
                  <span class="text-xs font-bold font-mono text-slate-900">{{ item.timeLabel }}</span>
                  <span
                    v-if="item.isPeak"
                    class="text-[9px] font-bold px-1.5 py-0.2 rounded bg-amber-50 text-amber-700 border border-amber-200 flex items-center space-x-0.5"
                  >
                    <Flame class="w-2.5 h-2.5 text-amber-500" />
                    <span>Peak</span>
                  </span>
                </div>
                <div class="text-[10px] text-slate-400 font-mono mt-0.5">
                  1 Hour (60 mins)
                </div>
              </div>
            </div>

            <!-- Right: Price & Status -->
            <div class="text-right flex items-center space-x-3">
              <div>
                <div class="text-xs font-mono font-bold text-slate-900">
                  ฿{{ item.slot.hourly_rate }}
                </div>
                <div class="text-[9px]">
                  <span v-if="item.slot.status === 'AVAILABLE'" class="text-emerald-600 font-medium">Available</span>
                  <span v-else-if="item.slot.status === 'LOCKED_PENDING'" class="text-amber-600 font-medium">Pending Transfer</span>
                  <span v-else class="text-slate-400 font-medium">Booked</span>
                </div>
              </div>

              <div
                v-if="item.isSelected"
                class="px-2 py-0.5 rounded-md bg-emerald-600 text-white text-[10px] font-bold"
              >
                Selected
              </div>
              <div
                v-else-if="item.slot.status === 'AVAILABLE'"
                class="px-2 py-0.5 rounded-md bg-slate-100 hover:bg-emerald-100 text-slate-700 hover:text-emerald-800 text-[10px] font-semibold transition"
              >
                Select
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Floating Bottom Checkout Dock (Playtomic / Clean Mobile Checkout Bar - Hidden on Desktop where Sidebar handles checkout) -->
    <div class="absolute bottom-2 left-2 right-2 z-20 lg:hidden">
      <div class="bg-white/95 backdrop-blur-md text-slate-900 p-3 rounded-2xl shadow-xl border border-slate-200/90">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-[10px] text-slate-500 font-medium">
              {{ bookingStore.selectedSlots.length }} slot(s) selected
            </div>
            <div class="text-base font-black text-slate-900 font-mono flex items-baseline space-x-1">
              <span class="text-emerald-700">฿{{ bookingStore.courtTotal.toLocaleString() }}</span>
              <span class="text-[10px] text-slate-400 font-normal font-sans">({{ bookingStore.selectedSlots.length }}h)</span>
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
            <span>Next: Add-ons</span>
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
            <span>Court {{ slot.courtNumber }} ({{ slot.startHour }}:00)</span>
            <button @click="bookingStore.toggleCourtSlot({ court_id: slot.courtId }, { start_hour: slot.startHour, status: 'AVAILABLE' })" class="text-emerald-700 hover:text-slate-900 ml-0.5 cursor-pointer">×</button>
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

