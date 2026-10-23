<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useBookingStore } from '@/stores/bookingStore'
import { QrCode, Timer, Upload, CheckCircle2, XCircle, AlertTriangle, HelpCircle, ShieldCheck } from 'lucide-vue-next'

const bookingStore = useBookingStore()

// Countdown timer (10 mins from expires_at)
const remainingSeconds = ref(600)
let timerInterval = null

onMounted(() => {
  if (bookingStore.activeBooking?.booking?.expires_at) {
    const expires = new Date(bookingStore.activeBooking.booking.expires_at).getTime()
    updateRemaining(expires)
    timerInterval = setInterval(() => updateRemaining(expires), 1000)
  }
})

onUnmounted(() => {
  if (timerInterval) clearInterval(timerInterval)
})

function updateRemaining(expires) {
  const diff = Math.floor((expires - Date.now()) / 1000)
  remainingSeconds.value = Math.max(0, diff)
  if (remainingSeconds.value === 0 && timerInterval) {
    clearInterval(timerInterval)
  }
}

const formattedTimer = computed(() => {
  const m = Math.floor(remainingSeconds.value / 60)
  const s = remainingSeconds.value % 60
  return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
})

const isNearExpiry = computed(() => remainingSeconds.value < 120)

// File upload handling
const selectedFile = ref(null)
const previewImage = ref(null)

function onFileSelected(e) {
  const file = e.target.files[0]
  if (file) {
    selectedFile.value = file
    previewImage.value = URL.createObjectURL(file)
  }
}

async function triggerVerification(mockScenario = 'MOCK_VALID') {
  await bookingStore.verifySlip(mockScenario, selectedFile.value)
}
</script>

<template>
  <div class="flex flex-col h-full bg-slate-50 overflow-y-auto p-4 space-y-4">
    <!-- Booking Ref & Countdown Timer Banner -->
    <div
      :class="[
        'rounded-2xl p-3.5 border shadow-sm transition flex items-center justify-between',
        isNearExpiry
          ? 'bg-rose-50 border-rose-300 text-rose-900 animate-pulse'
          : 'bg-white border-slate-200 text-slate-800'
      ]"
    >
      <div>
        <div class="text-[11px] text-slate-400 font-mono">รหัสการจองชั่วคราว</div>
        <div class="text-sm font-bold font-mono text-emerald-800">
          {{ bookingStore.activeBooking?.booking?.booking_no || 'BK-XXXXX' }}
        </div>
      </div>

      <div class="text-right">
        <div class="text-[10px] text-slate-400 flex items-center justify-end space-x-1">
          <Timer class="w-3.5 h-3.5 text-amber-500" />
          <span>เวลาล็อคสล็อต</span>
        </div>
        <div
          :class="[
            'text-lg font-mono font-extrabold tracking-wider',
            isNearExpiry ? 'text-rose-600' : 'text-amber-600'
          ]"
        >
          {{ formattedTimer }}
        </div>
      </div>
    </div>

    <!-- PromptPay QR Code Box -->
    <div class="bg-white rounded-2xl border border-slate-200 p-4 shadow-sm text-center">
      <div class="flex items-center justify-center space-x-1.5 text-xs font-semibold text-slate-700 mb-1">
        <QrCode class="w-4 h-4 text-emerald-600" />
        <span>สแกนจ่ายผ่านพร้อมเพย์ (PromptPay)</span>
      </div>
      <p class="text-[11px] text-slate-400 mb-3">บันทึกรูปหรือสแกนด้วย Mobile Banking ทุกธนาคาร</p>

      <!-- QR Code Image -->
      <div class="inline-block bg-white p-3 rounded-2xl border-2 border-slate-100 shadow-inner my-1">
        <!-- SVG / Image QR Code -->
        <img
          :src="`https://api.qrserver.com/v1/create-qr-code/?size=180x180&data=${encodeURIComponent(bookingStore.activeBooking?.promptpay_qr || 'smashhub')}`"
          alt="PromptPay QR Code"
          class="w-44 h-44 object-contain mx-auto"
        />
        <div class="mt-2 text-xs font-semibold text-slate-500 flex items-center justify-center space-x-1">
          <ShieldCheck class="w-3.5 h-3.5 text-blue-500" />
          <span>พร้อมเพย์ e-Payment</span>
        </div>
      </div>

      <!-- Amount to Pay -->
      <div class="mt-3 bg-emerald-50 rounded-xl p-2.5 border border-emerald-100">
        <div class="text-[11px] text-emerald-700">ยอดที่ต้องชำระตรงตามบิล</div>
        <div class="text-xl font-black text-emerald-800 font-mono">
          ฿{{ bookingStore.activeBooking?.booking?.net_total?.toLocaleString() || '0' }}
        </div>
        <div class="text-[10px] text-slate-500 mt-0.5">
          {{ bookingStore.activeBooking?.promptpay_name || 'บจก. สแมชฮับ แบดมินตัน' }} ({{ bookingStore.activeBooking?.promptpay_phone || '098-7-65432-1' }})
        </div>
      </div>
    </div>

    <!-- Slip Upload Box -->
    <div class="bg-white rounded-2xl border border-slate-200 p-4 shadow-sm">
      <div class="text-xs font-bold text-slate-800 mb-1 flex items-center space-x-1.5">
        <Upload class="w-4 h-4 text-emerald-600" />
        <span>แนบสลิปโอนเงิน</span>
      </div>
      <p class="text-[11px] text-slate-400 mb-3">อัปโหลดสลิปธนาคารเพื่อตรวจสอบอัตโนมัติ</p>

      <!-- Custom File Input -->
      <label class="block cursor-pointer">
        <div class="border-2 border-dashed border-slate-300 hover:border-emerald-400 rounded-xl p-4 text-center bg-slate-50 transition">
          <div v-if="previewImage" class="space-y-2">
            <img :src="previewImage" class="w-24 h-32 object-cover mx-auto rounded-lg shadow-sm border border-slate-200" />
            <div class="text-xs text-emerald-600 font-medium">เปลี่ยนไฟล์สลิป</div>
          </div>
          <div v-else class="space-y-1">
            <Upload class="w-6 h-6 text-slate-400 mx-auto" />
            <div class="text-xs text-slate-600 font-medium">กดเพื่อเลือกรูปภาพสลิป</div>
            <div class="text-[10px] text-slate-400">รองรับไฟล์ JPG, PNG</div>
          </div>
        </div>
        <input type="file" accept="image/*" class="hidden" @change="onFileSelected" />
      </label>

      <!-- Submit File Button -->
      <button
        v-if="selectedFile"
        @click="triggerVerification('MOCK_VALID')"
        :disabled="bookingStore.isSubmitting"
        class="w-full mt-3 py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl font-bold text-xs shadow transition flex items-center justify-center space-x-1.5"
      >
        <span v-if="bookingStore.isSubmitting" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
        <span>ตรวจสอบสลิปไฟล์ที่เลือก</span>
      </button>
    </div>

    <!-- Portfolio Mock Testing Tools (Crucial for recruiter test) -->
    <div class="bg-gradient-to-br from-slate-900 to-slate-800 rounded-2xl p-4 text-white shadow-md border border-slate-700">
      <div class="flex items-center space-x-1.5 text-xs font-bold text-emerald-400 mb-1">
        <HelpCircle class="w-4 h-4" />
        <span>Portfolio Mock Testing Scenarios</span>
      </div>
      <p class="text-[11px] text-slate-300 mb-3 leading-relaxed">
        ปุ่มทดสอบระบบตรวจสลิป 4 ชั้น (คลิกเพื่อทดสอบ Edge Cases โดยไม่ต้องโอนเงินจริง):
      </p>

      <div class="grid grid-cols-2 gap-2 text-[11px]">
        <!-- 1. MOCK_VALID -->
        <button
          @click="triggerVerification('MOCK_VALID')"
          :disabled="bookingStore.isSubmitting"
          class="p-2.5 bg-emerald-700 hover:bg-emerald-600 text-white rounded-xl font-medium transition text-left flex flex-col justify-between border border-emerald-500/40 active:scale-95"
        >
          <div class="font-bold flex items-center space-x-1 text-emerald-200">
            <span>🟢 MOCK_VALID</span>
          </div>
          <div class="text-[10px] text-emerald-100 opacity-90 mt-1">สลิปถูกต้อง ผ่านฉลุย -> CONFIRMED</div>
        </button>

        <!-- 2. MOCK_DUPLICATE -->
        <button
          @click="triggerVerification('MOCK_DUPLICATE')"
          :disabled="bookingStore.isSubmitting"
          class="p-2.5 bg-rose-900/70 hover:bg-rose-900 text-white rounded-xl font-medium transition text-left flex flex-col justify-between border border-rose-600/40 active:scale-95"
        >
          <div class="font-bold flex items-center space-x-1 text-rose-300">
            <span>🔴 MOCK_DUPLICATE</span>
          </div>
          <div class="text-[10px] text-rose-200 opacity-90 mt-1">สลิปซ้ำ trans_ref เดิมในระบบ</div>
        </button>

        <!-- 3. MOCK_UNDERPAID -->
        <button
          @click="triggerVerification('MOCK_UNDERPAID')"
          :disabled="bookingStore.isSubmitting"
          class="p-2.5 bg-amber-900/70 hover:bg-amber-900 text-white rounded-xl font-medium transition text-left flex flex-col justify-between border border-amber-600/40 active:scale-95"
        >
          <div class="font-bold flex items-center space-x-1 text-amber-300">
            <span>🟡 MOCK_UNDERPAID</span>
          </div>
          <div class="text-[10px] text-amber-200 opacity-90 mt-1">ยอดโอนขาด &lt; ยอดในบิล</div>
        </button>

        <!-- 4. MOCK_API_DOWN -->
        <button
          @click="triggerVerification('MOCK_API_DOWN')"
          :disabled="bookingStore.isSubmitting"
          class="p-2.5 bg-purple-900/70 hover:bg-purple-900 text-white rounded-xl font-medium transition text-left flex flex-col justify-between border border-purple-600/40 active:scale-95"
        >
          <div class="font-bold flex items-center space-x-1 text-purple-300">
            <span>🟠 MOCK_API_DOWN</span>
          </div>
          <div class="text-[10px] text-purple-200 opacity-90 mt-1">Gateway ล่ม -> MANUAL_REVIEW</div>
        </button>
      </div>

      <!-- Verification Feedback Banner -->
      <div v-if="bookingStore.slipResult" class="mt-3 p-3 rounded-xl text-xs bg-slate-950/80 border border-slate-700">
        <div class="flex items-center space-x-2 font-bold mb-1">
          <CheckCircle2 v-if="bookingStore.slipResult.success" class="w-4 h-4 text-emerald-400" />
          <XCircle v-else class="w-4 h-4 text-rose-400" />
          <span :class="bookingStore.slipResult.success ? 'text-emerald-400' : 'text-rose-400'">
            สถานะ: {{ bookingStore.slipResult.verification_status }}
          </span>
        </div>
        <p class="text-[11px] text-slate-300">{{ bookingStore.slipResult.message }}</p>
      </div>
    </div>
  </div>
</template>
