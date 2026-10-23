<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useBookingStore } from '@/stores/bookingStore'
import {
  QrCode, Timer, Upload, CheckCircle2, XCircle, AlertTriangle,
  Copy, Check, ShieldCheck, Download, Sparkles, Image as ImageIcon
} from 'lucide-vue-next'

const bookingStore = useBookingStore()

// Countdown timer
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

// 1-Click Copy Feedback
const isCopied = ref(false)
function copyAccount() {
  navigator.clipboard.writeText('0987654321')
  isCopied.value = true
  setTimeout(() => {
    isCopied.value = false
  }, 2000)
}

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
  <div class="flex flex-col h-full bg-slate-100/70 overflow-y-auto p-3 space-y-3 pb-8">
    <!-- Booking Reference & Timer Header -->
    <div
      :class="[
        'rounded-2xl p-3 border shadow-xs transition flex items-center justify-between',
        isNearExpiry
          ? 'bg-rose-50 border-rose-300 text-rose-900 animate-pulse'
          : 'bg-white border-slate-200 text-slate-800'
      ]"
    >
      <div>
        <div class="text-[9px] uppercase tracking-wider text-slate-400 font-mono font-bold">Booking ID</div>
        <div class="text-xs font-black font-mono text-emerald-800">
          {{ bookingStore.activeBooking?.booking?.booking_no || 'BK-XXXXX' }}
        </div>
      </div>

      <div class="text-right flex items-center space-x-2">
        <div class="p-1.5 rounded-lg bg-amber-500/10 text-amber-600">
          <Timer class="w-4 h-4" />
        </div>
        <div>
          <div class="text-[9px] text-slate-400">หมดเวลาล็อคสล็อตใน</div>
          <div
            :class="[
              'text-base font-mono font-black tracking-wider leading-none',
              isNearExpiry ? 'text-rose-600' : 'text-amber-600'
            ]"
          >
            {{ formattedTimer }}
          </div>
        </div>
      </div>
    </div>

    <!-- Authentic PromptPay e-Payment Card -->
    <div class="bg-white rounded-3xl border border-slate-200 shadow-sm overflow-hidden text-center">
      <!-- PromptPay Official Header Banner -->
      <div class="bg-[#003d6b] text-white py-2.5 px-4 flex items-center justify-between">
        <div class="flex items-center space-x-1.5 text-xs font-extrabold tracking-wider">
          <span class="text-sky-400">Prompt</span><span>Pay</span>
          <span class="text-[9px] bg-sky-500/30 text-sky-200 px-1.5 py-0.2 rounded font-sans font-normal">e-Payment</span>
        </div>
        <span class="text-[10px] text-slate-300 font-mono">SCAN TO PAY</span>
      </div>

      <div class="p-4">
        <!-- Amount To Pay Highlight -->
        <div class="mb-3">
          <div class="text-[10px] uppercase tracking-wider text-slate-400 font-semibold">ยอดชำระสุทธิ</div>
          <div class="text-2xl font-black text-slate-900 font-mono tracking-tight mt-0.5">
            ฿{{ bookingStore.activeBooking?.booking?.net_total?.toLocaleString() || '0' }}
          </div>
        </div>

        <!-- QR Box -->
        <div class="inline-block bg-white p-2.5 rounded-2xl border-2 border-slate-100 shadow-inner">
          <img
            :src="`https://api.qrserver.com/v1/create-qr-code/?size=160x160&data=${encodeURIComponent(bookingStore.activeBooking?.promptpay_qr || 'smashhub')}`"
            alt="PromptPay QR Code"
            class="w-36 h-36 object-contain mx-auto"
          />
        </div>

        <!-- 1-Click Copy Account & Store Details -->
        <div class="mt-3 bg-slate-50 border border-slate-200/80 rounded-2xl p-2.5 text-left flex items-center justify-between">
          <div>
            <div class="text-[10px] text-slate-400">ชื่อบัญชี: บจก. สแมชฮับ แบดมินตัน</div>
            <div class="text-xs font-bold font-mono text-slate-800">
              พร้อมเพย์: 098-7-65432-1
            </div>
          </div>

          <button
            @click="copyAccount"
            :class="[
              'px-2.5 py-1.5 rounded-xl text-[10px] font-bold transition flex items-center space-x-1 shadow-2xs active:scale-95',
              isCopied
                ? 'bg-emerald-600 text-white'
                : 'bg-white hover:bg-slate-100 text-slate-700 border border-slate-300'
            ]"
          >
            <Check v-if="isCopied" class="w-3 h-3 text-white" />
            <Copy v-else class="w-3 h-3 text-slate-500" />
            <span>{{ isCopied ? 'คัดลอกแล้ว' : 'คัดลอกเบอร์' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Slip Upload Zone -->
    <div class="bg-white rounded-3xl border border-slate-200 p-3.5 shadow-xs">
      <div class="flex items-center justify-between mb-2">
        <div class="text-xs font-bold text-slate-900 flex items-center space-x-1.5">
          <Upload class="w-3.5 h-3.5 text-emerald-600" />
          <span>แนบสลิปเพื่อตรวจสอบ</span>
        </div>
        <span class="text-[10px] text-emerald-600 font-medium">ตรวจผ่าน AI อัตโนมัติ</span>
      </div>

      <!-- File Input Drop Area -->
      <label class="block cursor-pointer">
        <div class="border-2 border-dashed border-slate-300 hover:border-emerald-400 rounded-2xl p-3 text-center bg-slate-50 transition">
          <div v-if="previewImage" class="space-y-1.5">
            <img :src="previewImage" class="w-20 h-28 object-cover mx-auto rounded-lg shadow-xs border border-slate-200" />
            <div class="text-[11px] text-emerald-600 font-semibold">เปลี่ยนไฟล์รูปภาพสลิป</div>
          </div>
          <div v-else class="space-y-1 py-2">
            <ImageIcon class="w-6 h-6 text-slate-400 mx-auto" />
            <div class="text-xs text-slate-700 font-semibold">แตะเพื่อเลือกรูปภาพสลิปจากอัลบั้ม</div>
            <div class="text-[9px] text-slate-400">รองรับสลิปโอนเงินทุกธนาคาร (JPG, PNG)</div>
          </div>
        </div>
        <input type="file" accept="image/*" class="hidden" @change="onFileSelected" />
      </label>

      <!-- Submit File Button -->
      <button
        v-if="selectedFile"
        @click="triggerVerification('MOCK_VALID')"
        :disabled="bookingStore.isSubmitting"
        class="w-full mt-2.5 py-2.5 bg-emerald-600 hover:bg-emerald-500 text-white rounded-xl font-bold text-xs shadow-md transition flex items-center justify-center space-x-1.5"
      >
        <span v-if="bookingStore.isSubmitting" class="w-3.5 h-3.5 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
        <span>ส่งตรวจสอบสลิปนี้ทันที</span>
      </button>
    </div>

    <!-- Portfolio Mock Testing Pill Container -->
    <div class="bg-slate-900 rounded-3xl p-3.5 text-white shadow-xl border border-slate-800">
      <div class="flex items-center space-x-1.5 text-xs font-bold text-emerald-400 mb-1">
        <Sparkles class="w-3.5 h-3.5" />
        <span>Portfolio Mock Scenarios</span>
      </div>
      <p class="text-[10px] text-slate-400 mb-2.5">
        คลิกทดสอบ 4 เหตุการณ์ตรวจสลิปได้ทันที (ไม่ต้องโอนเงินจริง):
      </p>

      <div class="grid grid-cols-2 gap-1.5 text-[10px]">
        <button
          @click="triggerVerification('MOCK_VALID')"
          :disabled="bookingStore.isSubmitting"
          class="p-2 bg-emerald-950/80 hover:bg-emerald-900 border border-emerald-500/40 text-emerald-200 rounded-xl font-medium transition text-left flex flex-col justify-between"
        >
          <span class="font-bold">🟢 MOCK_VALID</span>
          <span class="text-[9px] text-emerald-400/80 mt-0.5">ผ่านฉลุย -> ยืนยันคอร์ท</span>
        </button>

        <button
          @click="triggerVerification('MOCK_DUPLICATE')"
          :disabled="bookingStore.isSubmitting"
          class="p-2 bg-rose-950/80 hover:bg-rose-900 border border-rose-500/40 text-rose-200 rounded-xl font-medium transition text-left flex flex-col justify-between"
        >
          <span class="font-bold">🔴 MOCK_DUPLICATE</span>
          <span class="text-[9px] text-rose-400/80 mt-0.5">ตรวจจับสลิปซ้ำเดิม</span>
        </button>

        <button
          @click="triggerVerification('MOCK_UNDERPAID')"
          :disabled="bookingStore.isSubmitting"
          class="p-2 bg-amber-950/80 hover:bg-amber-900 border border-amber-500/40 text-amber-200 rounded-xl font-medium transition text-left flex flex-col justify-between"
        >
          <span class="font-bold">🟡 MOCK_UNDERPAID</span>
          <span class="text-[9px] text-amber-400/80 mt-0.5">ตรวจจับยอดโอนขาด</span>
        </button>

        <button
          @click="triggerVerification('MOCK_API_DOWN')"
          :disabled="bookingStore.isSubmitting"
          class="p-2 bg-purple-950/80 hover:bg-purple-900 border border-purple-500/40 text-purple-200 rounded-xl font-medium transition text-left flex flex-col justify-between"
        >
          <span class="font-bold">🟣 MOCK_API_DOWN</span>
          <span class="text-[9px] text-purple-400/80 mt-0.5">ส่งเข้าคิวแอดมินตรวจ</span>
        </button>
      </div>

      <!-- Result Banner -->
      <div v-if="bookingStore.slipResult" class="mt-2.5 p-2.5 rounded-xl text-[11px] bg-slate-950 border border-slate-800">
        <div class="flex items-center space-x-1.5 font-bold mb-0.5">
          <CheckCircle2 v-if="bookingStore.slipResult.success" class="w-3.5 h-3.5 text-emerald-400" />
          <XCircle v-else class="w-3.5 h-3.5 text-rose-400" />
          <span :class="bookingStore.slipResult.success ? 'text-emerald-400' : 'text-rose-400'">
            {{ bookingStore.slipResult.verification_status }}
          </span>
        </div>
        <p class="text-[10px] text-slate-300 leading-tight">{{ bookingStore.slipResult.message }}</p>
      </div>
    </div>
  </div>
</template>
