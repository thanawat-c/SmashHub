<script setup>
import { ref } from 'vue'
import { useAdminStore } from '@/stores/adminStore'
import { Check, X, ShieldAlert, FileText, CheckCircle2, AlertTriangle, ExternalLink } from 'lucide-vue-next'

const adminStore = useAdminStore()
const adminNote = ref('')
const isProcessing = ref(false)

async function handleAction(action) {
  if (!adminStore.selectedBooking) return
  isProcessing.value = true
  await adminStore.manualApprove(adminStore.selectedBooking.id, action, adminNote.value)
  isProcessing.value = false
  adminNote.value = ''
}
</script>

<template>
  <el-drawer
    v-model="adminStore.isDrawerOpen"
    title="ตรวจสอบสลิปการโอนเงิน (Slip Audit)"
    size="460px"
    destroy-on-close
  >
    <div v-if="adminStore.selectedBooking" class="space-y-4 text-xs text-slate-700">
      <!-- Booking Overview Banner -->
      <div class="bg-slate-50 p-3.5 rounded-xl border border-slate-200">
        <div class="flex justify-between items-start mb-2">
          <div>
            <div class="text-[10px] text-slate-400 font-mono">BOOKING NO</div>
            <div class="text-sm font-bold text-slate-900 font-mono">
              {{ adminStore.selectedBooking.booking_no }}
            </div>
          </div>
          <el-tag
            :type="
              adminStore.selectedBooking.status === 'CONFIRMED' ? 'success' :
              adminStore.selectedBooking.status === 'MANUAL_REVIEW' ? 'warning' :
              adminStore.selectedBooking.status === 'PENDING_PAYMENT' ? 'info' : 'danger'
            "
            effect="dark"
            size="small"
          >
            {{ adminStore.selectedBooking.status }}
          </el-tag>
        </div>

        <div class="grid grid-cols-2 gap-2 text-[11px] pt-1 border-t border-slate-100">
          <div>
            <span class="text-slate-400">ลูกค้า:</span>
            <span class="font-semibold ml-1 text-slate-800">{{ adminStore.selectedBooking.user?.display_name }}</span>
          </div>
          <div class="text-right">
            <span class="text-slate-400">ยอดที่ต้องชำระ:</span>
            <span class="font-bold text-emerald-700 ml-1">฿{{ adminStore.selectedBooking.net_total?.toLocaleString() }}</span>
          </div>
        </div>
      </div>

      <!-- Slip Image Preview -->
      <div>
        <div class="font-bold text-slate-800 mb-2 flex items-center space-x-1.5">
          <FileText class="w-4 h-4 text-emerald-600" />
          <span>รูปภาพสลิปที่แนบมา</span>
        </div>

        <div class="bg-slate-100 rounded-xl overflow-hidden border border-slate-200 p-2 text-center">
          <img
            v-if="adminStore.selectedBooking.payment_slip?.slip_image_url"
            :src="adminStore.selectedBooking.payment_slip.slip_image_url"
            alt="Payment Slip"
            class="max-h-72 mx-auto rounded-lg shadow-sm object-contain"
          />
          <div v-else class="py-12 text-slate-400 text-xs">
            ยังไม่มีการแนบรูปภาพสลิป
          </div>
        </div>
      </div>

      <!-- Verification Audit Details -->
      <div v-if="adminStore.selectedBooking.payment_slip" class="space-y-2">
        <div class="font-bold text-slate-800 flex items-center justify-between">
          <span>ผลการตรวจ 4-Step Verification</span>
          <el-tag
            :type="
              adminStore.selectedBooking.payment_slip.verification_status === 'VERIFIED' ? 'success' :
              adminStore.selectedBooking.payment_slip.verification_status === 'API_ERROR_PENDING_ADMIN' ? 'warning' : 'danger'
            "
            size="small"
          >
            {{ adminStore.selectedBooking.payment_slip.verification_status }}
          </el-tag>
        </div>

        <div class="bg-slate-50 rounded-xl p-3 border border-slate-200 space-y-2 text-[11px]">
          <div class="flex justify-between">
            <span class="text-slate-400">Transaction Reference:</span>
            <span class="font-mono font-bold text-slate-800">
              {{ adminStore.selectedBooking.payment_slip.trans_ref || '-' }}
            </span>
          </div>
          <div class="flex justify-between">
            <span class="text-slate-400">ชื่อผู้โอน:</span>
            <span class="font-semibold text-slate-800">
              {{ adminStore.selectedBooking.payment_slip.sender_name || '-' }}
            </span>
          </div>
          <div class="flex justify-between">
            <span class="text-slate-400">บัญชีปลายทาง:</span>
            <span class="text-slate-800 font-mono">
              {{ adminStore.selectedBooking.payment_slip.receiver_account || '-' }}
            </span>
          </div>
          <div class="flex justify-between items-center pt-1 border-t border-slate-200">
            <span class="text-slate-400">ยอดเงินในสลิป:</span>
            <span class="font-mono font-bold text-emerald-700 text-sm">
              ฿{{ adminStore.selectedBooking.payment_slip.amount?.toLocaleString() }}
            </span>
          </div>
        </div>

        <!-- Raw JSON Payload Extracted from Slip -->
        <div>
          <div class="text-[10px] font-bold text-slate-400 uppercase tracking-wide mb-1">
            Raw Verification Payload (JSON)
          </div>
          <pre class="bg-slate-900 text-emerald-400 p-3 rounded-xl font-mono text-[10px] overflow-x-auto max-h-40 leading-relaxed">{{ JSON.stringify(adminStore.selectedBooking.payment_slip.raw_payload, null, 2) }}</pre>
        </div>
      </div>

      <!-- Manual Review Action Section (When Status is MANUAL_REVIEW) -->
      <div
        v-if="adminStore.selectedBooking.status === 'MANUAL_REVIEW' || adminStore.selectedBooking.status === 'PENDING_PAYMENT'"
        class="bg-amber-50 border border-amber-200 rounded-xl p-3.5 space-y-3"
      >
        <div class="flex items-center space-x-2 text-amber-900 font-bold text-xs">
          <AlertTriangle class="w-4 h-4 text-amber-600 flex-shrink-0" />
          <span>การตัดสินใจของเจ้าหน้าที่ (Admin Action)</span>
        </div>

        <el-input
          v-model="adminNote"
          type="textarea"
          :rows="2"
          placeholder="บันทึกหมายเหตุ (เช่น ตรวจสอบในแอพ K PLUS เรียบร้อย)"
          size="small"
        />

        <div class="flex space-x-2">
          <button
            @click="handleAction('APPROVE')"
            :disabled="isProcessing"
            class="flex-1 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg font-bold text-xs shadow-sm transition flex items-center justify-center space-x-1"
          >
            <Check class="w-3.5 h-3.5" />
            <span>อนุมัติ (Approve)</span>
          </button>

          <button
            @click="handleAction('REJECT')"
            :disabled="isProcessing"
            class="flex-1 py-2 bg-rose-600 hover:bg-rose-700 text-white rounded-lg font-bold text-xs shadow-sm transition flex items-center justify-center space-x-1"
          >
            <X class="w-3.5 h-3.5" />
            <span>ปฏิเสธ (Reject)</span>
          </button>
        </div>
      </div>
    </div>
  </el-drawer>
</template>
