<script setup>
import { onMounted } from 'vue'
import { useBookingStore } from '@/stores/bookingStore'
import { CheckCircle2, Calendar, Clock, MapPin, Sparkles, RotateCcw } from 'lucide-vue-next'
import confetti from 'canvas-confetti'

const bookingStore = useBookingStore()

onMounted(() => {
  // Trigger celebratory confetti
  confetti({
    particleCount: 80,
    spread: 60,
    origin: { y: 0.6 }
  })
})

function bookAnother() {
  bookingStore.resetSelection()
  bookingStore.fetchAvailability()
}
</script>

<template>
  <div class="flex flex-col h-full bg-slate-100 overflow-y-auto p-4 space-y-4">
    <!-- Success Banner -->
    <div class="text-center py-2">
      <div class="inline-flex items-center justify-center w-12 h-12 bg-emerald-100 text-emerald-600 rounded-full mb-2 shadow-inner">
        <CheckCircle2 class="w-7 h-7" />
      </div>
      <h2 class="text-base font-extrabold text-slate-800">การจองได้รับการยืนยันแล้ว!</h2>
      <p class="text-xs text-slate-500">ระบบตรวจสอบสลิปเรียบร้อย พร้อมลงสนาม</p>
    </div>

    <!-- E-Ticket Card (Boarding Pass Design) -->
    <div class="bg-white rounded-3xl shadow-xl overflow-hidden border border-slate-200 relative">
      <!-- Ticket Header -->
      <div class="bg-emerald-700 text-white p-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-1.5 font-bold tracking-wide text-xs">
            <span>🏸</span>
            <span>SMASHHUB E-TICKET</span>
          </div>
          <span class="text-[10px] bg-emerald-800/80 text-emerald-200 px-2 py-0.5 rounded-full font-mono">
            {{ bookingStore.activeBooking?.booking?.status || 'CONFIRMED' }}
          </span>
        </div>

        <div class="mt-3">
          <div class="text-[10px] text-emerald-200 font-mono">BOOKING REFERENCE</div>
          <div class="text-lg font-black tracking-wider font-mono">
            {{ bookingStore.activeBooking?.booking?.booking_no || 'BK-20261024-001' }}
          </div>
        </div>
      </div>

      <!-- Ticket Body -->
      <div class="p-4 space-y-3.5 text-xs text-slate-700">
        <!-- Customer & Date -->
        <div class="grid grid-cols-2 gap-2 border-b border-slate-100 pb-3">
          <div>
            <div class="text-[10px] text-slate-400 font-medium">ชื่อผู้จอง</div>
            <div class="font-bold text-slate-800 truncate">
              {{ bookingStore.activeBooking?.booking?.user?.display_name || 'ลูกค้าสแมชฮับ' }}
            </div>
          </div>
          <div>
            <div class="text-[10px] text-slate-400 font-medium">วันที่เล่น</div>
            <div class="font-bold text-slate-800 font-mono">
              {{ bookingStore.selectedDate }}
            </div>
          </div>
        </div>

        <!-- Court Slots Summary -->
        <div class="border-b border-slate-100 pb-3">
          <div class="text-[10px] text-slate-400 font-medium mb-1.5 flex items-center space-x-1">
            <Clock class="w-3.5 h-3.5 text-emerald-600" />
            <span>คอร์ท & เวลาที่จอง</span>
          </div>
          <div class="space-y-1">
            <div
              v-for="court in bookingStore.activeBooking?.booking?.courts || bookingStore.selectedSlots"
              :key="court.id || (court.courtId + '-' + court.startHour)"
              class="flex justify-between items-center bg-slate-50 p-2 rounded-xl border border-slate-100"
            >
              <div class="font-bold text-slate-800 text-xs">
                คอร์ท {{ court.court?.court_number || court.courtNumber }}
              </div>
              <div class="text-emerald-700 font-mono font-semibold text-xs">
                {{ court.start_hour || court.startHour }}:00 - {{ court.end_hour || court.endHour }}:00 น.
              </div>
            </div>
          </div>
        </div>

        <!-- Ordered Items Summary if any -->
        <div v-if="bookingStore.activeBooking?.booking?.items?.length || Object.keys(bookingStore.cartItems).length > 0" class="border-b border-slate-100 pb-3">
          <div class="text-[10px] text-slate-400 font-medium mb-1">สินค้ามินิสโตร์ (รับที่เคาน์เตอร์)</div>
          <div class="space-y-0.5 text-[11px] text-slate-600">
            <div
              v-for="item in bookingStore.cartItemsList"
              :key="item.product.id"
              class="flex justify-between"
            >
              <span>{{ item.product.name }}</span>
              <span class="font-mono">x{{ item.qty }}</span>
            </div>
          </div>
        </div>

        <!-- Check-in Gate QR Code Section -->
        <div class="text-center pt-1">
          <div class="text-[11px] font-bold text-slate-800 mb-2">สแกน QR Code เพื่อ Check-in เข้าสนาม</div>
          <div class="inline-block bg-slate-50 p-3 rounded-2xl border-2 border-dashed border-emerald-300">
            <img
              :src="`https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=${encodeURIComponent('SMASHHUB_CHECKIN_' + (bookingStore.activeBooking?.booking?.booking_no || 'DEMO'))}`"
              alt="Check-in QR"
              class="w-32 h-32 object-contain mx-auto"
            />
            <div class="text-[9px] text-slate-400 font-mono mt-1">STATUS: READY TO PLAY</div>
          </div>
        </div>
      </div>

      <!-- Ticket Perforation Cutout Styling -->
      <div class="relative py-2 bg-slate-50 border-t border-dashed border-slate-300 text-center text-[10px] text-slate-400 font-mono">
        SMASHHUB BADMINTON ARENA • BANGKOK
      </div>
    </div>

    <!-- Action Button -->
    <button
      @click="bookAnother"
      class="w-full py-3 bg-white hover:bg-slate-50 text-slate-700 rounded-2xl font-bold text-xs border border-slate-200 shadow-sm transition flex items-center justify-center space-x-1.5 active:scale-98"
    >
      <RotateCcw class="w-4 h-4 text-emerald-600" />
      <span>จองคอร์ทอื่นเพิ่มเติม</span>
    </button>
  </div>
</template>
