<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import QRCode from 'qrcode'
import { useBookingStore } from '@/stores/bookingStore'
import {
  Timer, Upload, CheckCircle2, XCircle, AlertTriangle,
  Copy, Check, Download, Image as ImageIcon, ChevronDown, ChevronUp,
  FileCheck2, ShieldCheck, Beaker
} from 'lucide-vue-next'

const bookingStore = useBookingStore()

// Countdown timer (10 mins)
const remainingSeconds = ref(600)
let timerInterval = null

const qrDataUrl = ref('')

function crc16(data) {
  let crc = 0xFFFF
  for (let i = 0; i < data.length; i++) {
    crc ^= data.charCodeAt(i) << 8
    for (let j = 0; j < 8; j++) {
      if ((crc & 0x8000) !== 0) {
        crc = ((crc << 1) ^ 0x1021) & 0xFFFF
      } else {
        crc = (crc << 1) & 0xFFFF
      }
    }
  }
  return crc.toString(16).toUpperCase().padStart(4, '0')
}

function generatePromptPayPayload(phoneNumber, amount) {
  const cleaned = phoneNumber.replace(/[^0-9]/g, '')
  const formattedPhone = cleaned.startsWith('0') ? '0066' + cleaned.substring(1) : cleaned
  const phoneTag = '0113' + formattedPhone
  const merchantInfo = '0016A000000677010111' + phoneTag
  const tag29 = '29' + String(merchantInfo.length).padStart(2, '0') + merchantInfo
  const formattedAmount = Number(amount || 0).toFixed(2)
  const tag54 = '54' + String(formattedAmount.length).padStart(2, '0') + formattedAmount
  const raw = '000201010212' + tag29 + '5303764' + tag54 + '5802TH6304'
  return raw + crc16(raw)
}

async function renderQRCode() {
  const amount = bookingStore.activeBooking?.booking?.net_total || bookingStore.netTotal || 150
  let payload = bookingStore.activeBooking?.promptpay_qr
  if (!payload || payload.startsWith('http') || payload === 'smashhub') {
    payload = generatePromptPayPayload('0987654321', amount)
  }
  try {
    qrDataUrl.value = await QRCode.toDataURL(payload, {
      width: 320,
      margin: 1,
      color: {
        dark: '#0f172a',
        light: '#ffffff'
      }
    })
  } catch (err) {
    console.error('Failed to generate local PromptPay QR:', err)
  }
}

function downloadQR() {
  if (!qrDataUrl.value) return
  const a = document.createElement('a')
  a.href = qrDataUrl.value
  a.download = `PromptPay-SmashHub-${bookingStore.activeBooking?.booking?.booking_no || 'BK'}.png`
  a.click()
}

watch(() => bookingStore.activeBooking, () => {
  renderQRCode()
}, { deep: true })

onMounted(() => {
  renderQRCode()
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

const showTesterPanel = ref(false)
</script>

<template>
  <div class="flex flex-col h-full bg-[#F8FAFC] overflow-y-auto p-3 space-y-3 pb-10">
    <!-- Timer & Booking Reference Header -->
    <div
      :class="[
        'rounded-2xl p-3 border shadow-2xs transition flex items-center justify-between',
        isNearExpiry
          ? 'bg-rose-50 border-rose-200 text-rose-900'
          : 'bg-white border-slate-200/80 text-slate-800'
      ]"
    >
      <div class="flex items-center space-x-2">
        <div class="w-8 h-8 rounded-xl bg-emerald-50 text-emerald-700 flex items-center justify-center font-bold text-xs">
          BK
        </div>
        <div>
          <div class="text-[9px] text-slate-400 font-medium">Booking Ref</div>
          <div class="text-xs font-black font-mono text-slate-900">
            {{ bookingStore.activeBooking?.booking?.booking_no || 'BK-20261024-XXXX' }}
          </div>
        </div>
      </div>

      <div class="text-right flex items-center space-x-1.5 bg-slate-50 px-2.5 py-1 rounded-xl border border-slate-200/60">
        <Timer :class="['w-3.5 h-3.5', isNearExpiry ? 'text-rose-600 animate-pulse' : 'text-amber-600']" />
        <div>
          <span class="text-[9px] text-slate-400 mr-1">Held for</span>
          <span :class="['text-xs font-mono font-black', isNearExpiry ? 'text-rose-600' : 'text-amber-700']">
            {{ formattedTimer }}
          </span>
        </div>
      </div>
    </div>

    <!-- Booking Summary & Item Breakdown Card -->
    <div class="bg-white rounded-2xl border border-slate-200/90 p-3.5 shadow-2xs space-y-2.5">
      <div class="flex items-center justify-between border-b border-slate-100 pb-2">
        <span class="text-xs font-bold text-slate-900">Order Summary</span>
        <span class="text-[10px] font-mono text-slate-400">Date: {{ bookingStore.selectedDate }}</span>
      </div>

      <!-- Court Slots List -->
      <div class="space-y-1.5 text-xs">
        <div class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">Courts</div>
        <div
          v-for="court in bookingStore.activeBooking?.booking?.courts || bookingStore.selectedSlots"
          :key="court.id || (court.courtId + '-' + court.startHour)"
          class="flex items-center justify-between bg-slate-50 p-2 rounded-xl border border-slate-100"
        >
          <div class="flex items-center space-x-1.5">
            <span class="font-extrabold text-slate-900">Court {{ court.court?.court_number || court.court_number || court.courtNumber || court.court_id || court.courtId || 1 }}</span>
            <span class="text-[10px] text-slate-500 font-mono">({{ court.start_hour || court.startHour }}:00 - {{ (court.start_hour || court.startHour) + 1 }}:00)</span>
          </div>
          <span class="font-bold font-mono text-slate-900">฿{{ court.hourly_rate || court.rate || 150 }}</span>
        </div>
      </div>

      <!-- Mini-Store Items if any -->
      <div v-if="(bookingStore.activeBooking?.booking?.items && bookingStore.activeBooking.booking.items.length > 0) || bookingStore.cartItemsList.length > 0" class="space-y-1.5 text-xs pt-1 border-t border-slate-100">
        <div class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">Add-ons & Drinks</div>
        <div
          v-for="item in (bookingStore.activeBooking?.booking?.items || bookingStore.cartItemsList)"
          :key="item.id || item.product?.id"
          class="flex items-center justify-between text-slate-700"
        >
          <span>{{ item.product?.name }} x{{ item.qty }}</span>
          <span class="font-bold font-mono text-slate-900">฿{{ item.subtotal || (item.product?.price * item.qty) }}</span>
        </div>
      </div>

      <!-- Totals Footer -->
      <div class="pt-2 border-t border-slate-100 flex items-center justify-between text-xs">
        <span class="text-slate-500 font-medium">Net Total</span>
        <span class="text-sm font-black font-mono text-emerald-700">
          ฿{{ Number(bookingStore.activeBooking?.booking?.net_total || bookingStore.netTotal).toLocaleString('th-TH', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}
        </span>
      </div>
    </div>

    <!-- Official Thai QR Payment Standard Card -->
    <div class="bg-white rounded-2xl border border-slate-200/90 shadow-xs overflow-hidden">
      <!-- Thai QR Payment Blue Bar -->
      <div class="bg-[#003B64] text-white px-4 py-2.5 flex items-center justify-between">
        <div class="flex items-center space-x-1.5 text-xs font-black tracking-wide">
          <span class="bg-white text-[#003B64] px-1 rounded-xs text-[10px] font-bold">THAI QR</span>
          <span>PAYMENT</span>
        </div>
        <span class="text-[10px] text-slate-200 font-medium">Scan via any banking app</span>
      </div>

      <div class="p-4 text-center">
        <!-- Amount To Pay -->
        <div class="mb-3">
          <div class="text-[10px] text-slate-500 font-medium">Amount Due</div>
          <div class="text-2xl font-black text-slate-900 font-mono tracking-tight mt-0.5">
            ฿{{ (bookingStore.activeBooking?.booking?.net_total || 0).toLocaleString('th-TH', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}
          </div>
        </div>

        <!-- High-Contrast QR Code -->
        <div class="inline-block p-2.5 bg-white rounded-2xl border border-slate-200 shadow-2xs">
          <img
            v-if="qrDataUrl"
            :src="qrDataUrl"
            alt="PromptPay QR Code"
            class="w-44 h-44 object-contain mx-auto"
          />
          <div v-else class="w-44 h-44 flex items-center justify-center bg-slate-50 text-slate-400 text-xs font-mono">
            Generating PromptPay QR...
          </div>
        </div>

        <!-- Save/Download QR button -->
        <div class="mt-2 flex items-center justify-center">
          <button
            @click="downloadQR"
            class="inline-flex items-center space-x-1.5 px-3 py-1 rounded-full bg-slate-100 hover:bg-slate-200 text-slate-700 text-[10px] font-semibold transition cursor-pointer"
          >
            <Download class="w-3 h-3 text-slate-500" />
            <span>Save QR Image</span>
          </button>
        </div>

        <!-- Bank & Account Information -->
        <div class="mt-3 bg-slate-50 border border-slate-100 rounded-xl p-2.5 text-left flex items-center justify-between">
          <div>
            <div class="text-[10px] text-slate-500 font-medium">Account: SmashHub Badminton Co., Ltd.</div>
            <div class="text-xs font-bold font-mono text-slate-800 mt-0.5">
              PromptPay: 098-7-65432-1
            </div>
          </div>

          <button
            @click="copyAccount"
            :class="[
              'px-2.5 py-1 rounded-lg text-[10px] font-semibold transition flex items-center space-x-1 border shadow-2xs cursor-pointer',
              isCopied
                ? 'bg-emerald-600 text-white border-emerald-600'
                : 'bg-white hover:bg-slate-100 text-slate-700 border-slate-200'
            ]"
          >
            <Check v-if="isCopied" class="w-3 h-3 text-white" />
            <Copy v-else class="w-3 h-3 text-slate-500" />
            <span>{{ isCopied ? 'Copied' : 'Copy' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Slip Upload Section -->
    <div class="bg-white rounded-2xl border border-slate-200/90 p-3.5 shadow-xs">
      <div class="flex items-center justify-between mb-2">
        <div class="text-xs font-bold text-slate-900 flex items-center space-x-1.5">
          <Upload class="w-3.5 h-3.5 text-emerald-600" />
          <span>Upload Payment Slip</span>
        </div>
        <span class="text-[10px] text-slate-400">JPG, PNG</span>
      </div>

      <!-- File Upload Drop Box -->
      <label class="block cursor-pointer">
        <div class="border-2 border-dashed border-slate-200 hover:border-emerald-500 rounded-xl p-3 text-center bg-slate-50/70 hover:bg-emerald-50/30 transition">
          <div v-if="previewImage" class="space-y-1.5">
            <img :src="previewImage" class="w-20 h-28 object-cover mx-auto rounded-lg shadow-xs border border-slate-200" />
            <div class="text-[11px] text-emerald-700 font-semibold">Tap to change slip image</div>
          </div>
          <div v-else class="space-y-1 py-2">
            <ImageIcon class="w-6 h-6 text-slate-400 mx-auto" />
            <div class="text-xs text-slate-700 font-semibold">Tap to upload transfer slip</div>
            <div class="text-[9px] text-slate-400">Automated OCR and bank transaction verification</div>
          </div>
        </div>
        <input type="file" accept="image/*" class="hidden" @change="onFileSelected" />
      </label>

      <!-- Submit File Button -->
      <button
        v-if="selectedFile"
        @click="triggerVerification('MOCK_VALID')"
        :disabled="bookingStore.isSubmitting"
        class="w-full mt-2.5 py-2.5 bg-emerald-600 hover:bg-emerald-500 text-white rounded-xl font-bold text-xs shadow-sm transition flex items-center justify-center space-x-1.5 active:scale-95 disabled:opacity-50 cursor-pointer"
      >
        <span v-if="bookingStore.isSubmitting" class="w-3.5 h-3.5 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
        <span>Submit Slip for Verification</span>
      </button>

      <!-- Quick 1-Click Verification for Customer Flow -->
      <button
        v-else
        @click="triggerVerification('MOCK_VALID')"
        :disabled="bookingStore.isSubmitting"
        class="w-full mt-2.5 py-2.5 bg-emerald-600 hover:bg-emerald-500 text-white rounded-xl font-bold text-xs shadow-sm transition flex items-center justify-center space-x-1.5 active:scale-95 disabled:opacity-50 cursor-pointer"
      >
        <span v-if="bookingStore.isSubmitting" class="w-3.5 h-3.5 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
        <FileCheck2 v-else class="w-3.5 h-3.5" />
        <span>Simulate Valid Slip (Instant Confirmation)</span>
      </button>
    </div>

    <!-- Clean Collapsible Evaluation / Testing Drawer (For Recruiters & Testers) -->
    <div class="bg-slate-50 border border-slate-200/80 rounded-2xl overflow-hidden transition shadow-2xs">
      <button
        @click="showTesterPanel = !showTesterPanel"
        class="w-full px-3 py-2 text-left flex items-center justify-between text-slate-600 hover:text-slate-900 transition cursor-pointer"
      >
        <div class="flex items-center space-x-1.5 text-[11px] font-semibold text-slate-600">
          <Beaker class="w-3.5 h-3.5 text-emerald-600" />
          <span>Slip Verification Sandbox (Exception Testing)</span>
        </div>
        <ChevronUp v-if="showTesterPanel" class="w-3.5 h-3.5 text-slate-400" />
        <ChevronDown v-else class="w-3.5 h-3.5 text-slate-400" />
      </button>

      <div v-if="showTesterPanel" class="p-3 border-t border-slate-200/60 bg-white space-y-2">
        <p class="text-[10px] text-slate-500">
          Simulate bank verification scenarios to test system reactions:
        </p>
        <div class="grid grid-cols-2 gap-1.5 text-[10px]">
          <button
            @click="triggerVerification('MOCK_VALID')"
            :disabled="bookingStore.isSubmitting"
            class="p-2 bg-emerald-50 hover:bg-emerald-100 text-emerald-800 border border-emerald-200 rounded-xl font-medium transition text-left cursor-pointer"
          >
            <div class="flex items-center space-x-1 font-bold">
              <CheckCircle2 class="w-3.5 h-3.5 text-emerald-600 shrink-0" />
              <span>Valid Slip</span>
            </div>
            <div class="text-[9px] text-emerald-600 mt-0.5">Verified -> Booking Confirmed</div>
          </button>

          <button
            @click="triggerVerification('MOCK_DUPLICATE')"
            :disabled="bookingStore.isSubmitting"
            class="p-2 bg-rose-50 hover:bg-rose-100 text-rose-800 border border-rose-200 rounded-xl font-medium transition text-left cursor-pointer"
          >
            <div class="flex items-center space-x-1 font-bold">
              <XCircle class="w-3.5 h-3.5 text-rose-600 shrink-0" />
              <span>Duplicate Slip</span>
            </div>
            <div class="text-[9px] text-rose-600 mt-0.5">Reused slip -> Rejected</div>
          </button>

          <button
            @click="triggerVerification('MOCK_UNDERPAID')"
            :disabled="bookingStore.isSubmitting"
            class="p-2 bg-amber-50 hover:bg-amber-100 text-amber-800 border border-amber-200 rounded-xl font-medium transition text-left cursor-pointer"
          >
            <div class="flex items-center space-x-1 font-bold">
              <AlertTriangle class="w-3.5 h-3.5 text-amber-600 shrink-0" />
              <span>Underpaid Slip</span>
            </div>
            <div class="text-[9px] text-amber-600 mt-0.5">Amount mismatch -> Rejected</div>
          </button>

          <button
            @click="triggerVerification('MOCK_API_DOWN')"
            :disabled="bookingStore.isSubmitting"
            class="p-2 bg-purple-50 hover:bg-purple-100 text-purple-800 border border-purple-200 rounded-xl font-medium transition text-left cursor-pointer"
          >
            <div class="flex items-center space-x-1 font-bold">
              <FileCheck2 class="w-3.5 h-3.5 text-purple-600 shrink-0" />
              <span>Bank Timeout</span>
            </div>
            <div class="text-[9px] text-purple-600 mt-0.5">Gateway down -> Manual Review</div>
          </button>
        </div>

        <!-- Verification Result Toast in Drawer -->
        <div v-if="bookingStore.slipResult" class="mt-2 p-2 rounded-xl text-[10px] bg-slate-50 border border-slate-200">
          <div class="flex items-center space-x-1.5 font-bold mb-0.5">
            <CheckCircle2 v-if="bookingStore.slipResult.success" class="w-3.5 h-3.5 text-emerald-600" />
            <XCircle v-else class="w-3.5 h-3.5 text-rose-600" />
            <span :class="bookingStore.slipResult.success ? 'text-emerald-700' : 'text-rose-700'">
              {{ bookingStore.slipResult.verification_status }}
            </span>
          </div>
          <p class="text-slate-600">{{ bookingStore.slipResult.message }}</p>
        </div>
      </div>
    </div>
  </div>
</template>
