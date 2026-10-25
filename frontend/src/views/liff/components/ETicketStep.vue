<script setup>
import { ref, onMounted } from 'vue'
import { useBookingStore } from '@/stores/bookingStore'
import { CheckCircle2, Clock, Share2, Check, RotateCcw, ShieldCheck, Ticket } from 'lucide-vue-next'
import confetti from 'canvas-confetti'

const bookingStore = useBookingStore()
const isShared = ref(false)

onMounted(() => {
  confetti({
    particleCount: 90,
    spread: 70,
    origin: { y: 0.6 }
  })
})

function bookAnother() {
  bookingStore.resetSelection()
  bookingStore.fetchAvailability()
}

function shareTicket() {
  const text = `🏸 ตั๋วเข้าสนาม SmashHub\nรหัส: ${bookingStore.activeBooking?.booking?.booking_no}\nวันที่: ${bookingStore.selectedDate}\nพร้อมลงสนามแล้ว!`
  navigator.clipboard.writeText(text)
  isShared.value = true
  setTimeout(() => isShared.value = false, 2500)
}
</script>

<template>
  <div class="flex flex-col h-full bg-[#F8FAFC] overflow-y-auto p-4 space-y-4 text-slate-800">
    <!-- Success Celebration Tag -->
    <div class="text-center pt-2">
      <div class="inline-flex items-center justify-center w-12 h-12 bg-emerald-50 text-emerald-600 rounded-2xl mb-2 border border-emerald-200 shadow-xs">
        <CheckCircle2 class="w-6 h-6" />
      </div>
      <h2 class="text-base font-extrabold text-slate-900">การจองได้รับการยืนยันแล้ว</h2>
      <p class="text-[11px] text-slate-500">ตรวจสอบการชำระเงินเรียบร้อย บัตรเข้าสนามพร้อมใช้งาน</p>
    </div>

    <!-- Authentic Boarding-Pass E-Ticket Card -->
    <div class="bg-white rounded-3xl shadow-sm overflow-hidden text-slate-800 relative border border-slate-200/90">
      <!-- Ticket Top Accent Bar -->
      <div class="bg-slate-900 text-white p-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-2 font-bold tracking-wider text-xs">
            <span>🏸</span>
            <span class="tracking-tight font-extrabold">SMASHHUB ARENA</span>
          </div>
          <span class="text-[10px] bg-emerald-500/20 text-emerald-300 border border-emerald-500/30 px-2 py-0.5 rounded-full font-mono font-bold">
            CONFIRMED
          </span>
        </div>

        <div class="mt-3 flex justify-between items-end">
          <div>
            <div class="text-[9px] text-slate-400 font-mono uppercase tracking-wider">Booking Ref</div>
            <div class="text-base font-black tracking-wider font-mono text-emerald-400">
              {{ bookingStore.activeBooking?.booking?.booking_no || 'BK-20261024-001' }}
            </div>
          </div>
          <div class="text-right">
            <div class="text-[9px] text-slate-400 uppercase tracking-wider">Date</div>
            <div class="text-xs font-bold font-mono text-white">{{ bookingStore.selectedDate }}</div>
          </div>
        </div>
      </div>

      <!-- Ticket Body Details -->
      <div class="p-4 space-y-3 text-xs text-slate-700">
        <!-- Customer Details -->
        <div class="flex justify-between items-center border-b border-slate-100 pb-2.5">
          <div>
            <div class="text-[10px] text-slate-400 font-medium">ชื่อผู้จอง</div>
            <div class="font-extrabold text-slate-900 text-xs">
              {{ bookingStore.activeBooking?.booking?.user?.display_name || 'ลูกค้าสแมชฮับ' }}
            </div>
          </div>
          <div class="text-right">
            <div class="text-[10px] text-slate-400 font-medium">สถานะ Check-in</div>
            <div class="inline-flex items-center space-x-1.5 text-emerald-700 font-bold text-[11px]">
              <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
              <span>รอเข้าสนาม</span>
            </div>
          </div>
        </div>

        <!-- Court Slots Summary -->
        <div class="space-y-1">
          <div class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">คอร์ท & เวลาที่จอง</div>
          <div
            v-for="court in bookingStore.activeBooking?.booking?.courts || bookingStore.selectedSlots"
            :key="court.id || (court.courtId + '-' + court.startHour)"
            class="flex justify-between items-center bg-slate-50 p-2.5 rounded-xl border border-slate-100"
          >
            <span class="font-extrabold text-slate-900 text-xs">
              คอร์ท {{ court.court?.court_number || court.courtNumber }}
            </span>
            <span class="text-emerald-700 font-mono font-bold text-xs">
              {{ court.start_hour || court.startHour }}:00 - {{ court.end_hour || court.endHour }}:00 น.
            </span>
          </div>
        </div>

        <!-- Mini store items if any -->
        <div v-if="bookingStore.cartItemsList.length > 0" class="border-t border-slate-100 pt-2 space-y-1">
          <div class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">สินค้ามินิสโตร์ (รับที่เคาน์เตอร์)</div>
          <div class="space-y-0.5 text-[11px] text-slate-600">
            <div v-for="item in bookingStore.cartItemsList" :key="item.product.id" class="flex justify-between">
              <span>{{ item.product.name }}</span>
              <span class="font-mono font-bold">x{{ item.qty }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Ticket Perforated Cutout Line with Side Notches -->
      <div class="relative py-1 bg-slate-50/80 border-t border-dashed border-slate-300">
        <div class="absolute -left-3 -top-2.5 w-5 h-5 rounded-full bg-[#F8FAFC] border-r border-slate-200"></div>
        <div class="absolute -right-3 -top-2.5 w-5 h-5 rounded-full bg-[#F8FAFC] border-l border-slate-200"></div>
      </div>

      <!-- Gate QR Code Section -->
      <div class="p-4 bg-slate-50 text-center space-y-1.5">
        <div class="text-[11px] font-bold text-slate-800">สแกน QR Code สำหรับ Check-in หน้าเคาน์เตอร์</div>
        <div class="inline-block bg-white p-2.5 rounded-2xl border border-slate-200 shadow-2xs">
          <img
            :src="`https://api.qrserver.com/v1/create-qr-code/?size=140x140&data=${encodeURIComponent('SMASHHUB_PASS_' + (bookingStore.activeBooking?.booking?.booking_no || 'DEMO'))}`"
            alt="Check-in Pass QR"
            class="w-32 h-32 object-contain mx-auto"
          />
        </div>
        <div class="text-[9px] text-slate-400 font-mono">DIGITAL PASS • SMASHHUB ARENA</div>
      </div>
    </div>

    <!-- Action Buttons -->
    <div class="space-y-2 pt-1">
      <button
        @click="shareTicket"
        class="w-full py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-2xl font-bold text-xs shadow-xs transition flex items-center justify-center space-x-1.5"
      >
        <Check v-if="isShared" class="w-4 h-4 text-emerald-200" />
        <Share2 v-else class="w-4 h-4" />
        <span>{{ isShared ? 'คัดลอกข้อความตั๋วแล้ว!' : 'แชร์ตั๋วให้เพื่อนในกลุ่ม LINE' }}</span>
      </button>

      <button
        @click="bookAnother"
        class="w-full py-2.5 bg-white hover:bg-slate-50 text-slate-700 rounded-2xl font-bold text-xs border border-slate-200 shadow-2xs transition flex items-center justify-center space-x-1.5"
      >
        <RotateCcw class="w-3.5 h-3.5 text-slate-500" />
        <span>จองคอร์ทอื่นเพิ่มเติม</span>
      </button>
    </div>
  </div>
</template>
