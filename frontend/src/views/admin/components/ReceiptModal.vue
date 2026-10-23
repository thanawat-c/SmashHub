<script setup>
import { computed } from 'vue'
import { Printer, Download, X } from 'lucide-vue-next'

const props = defineProps({
  visible: { type: Boolean, default: false },
  booking: { type: Object, default: () => null }
})

const emit = defineEmits(['update:visible'])

function close() {
  emit('update:visible', false)
}

function handlePrint() {
  window.print()
}
</script>

<template>
  <el-dialog
    :model-value="visible"
    @update:model-value="emit('update:visible', $event)"
    width="520px"
    align-center
    :show-close="false"
    class="rounded-2xl overflow-hidden p-0"
  >
    <!-- Modal Header Actions -->
    <div class="flex items-center justify-between p-4 border-b border-slate-100 bg-slate-50 print:hidden">
      <div class="font-bold text-slate-800 text-sm flex items-center space-x-1.5">
        <Printer class="w-4 h-4 text-emerald-600" />
        <span>ใบเสร็จรับเงิน / Tax Invoice</span>
      </div>
      <div class="flex items-center space-x-2">
        <button
          @click="handlePrint"
          class="px-3 py-1.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg text-xs font-semibold flex items-center space-x-1 shadow-sm transition"
        >
          <Printer class="w-3.5 h-3.5" />
          <span>สั่งพิมพ์ / PDF</span>
        </button>
        <button
          @click="close"
          class="p-1.5 text-slate-400 hover:text-slate-600 hover:bg-slate-200 rounded-lg transition"
        >
          <X class="w-4 h-4" />
        </button>
      </div>
    </div>

    <!-- Printable Receipt Content Area -->
    <div v-if="booking" id="printable-receipt" class="p-6 text-slate-800 font-sans text-xs bg-white">
      <!-- Store Header -->
      <div class="text-center border-b border-dashed border-slate-300 pb-4 mb-4">
        <div class="text-2xl mb-1">🏸</div>
        <h1 class="text-base font-extrabold tracking-wider text-slate-900">บจก. สแมชฮับ แบดมินตัน</h1>
        <p class="text-[11px] text-slate-500">SMASHHUB BADMINTON ARENA</p>
        <p class="text-[10px] text-slate-400 mt-0.5">888 ถนนรามคำแหง แขวงหัวหมาก เขตบางกะปิ กทม. 10240</p>
        <p class="text-[10px] text-slate-400 font-mono">เลขประจำตัวผู้เสียภาษี: 0105567012345 • โทร. 098-7-65432-1</p>
      </div>

      <!-- Receipt Meta Info -->
      <div class="grid grid-cols-2 gap-2 text-[11px] border-b border-slate-200 pb-3 mb-3">
        <div>
          <span class="text-slate-400">เลขที่ใบเสร็จ:</span>
          <span class="font-mono font-bold ml-1 text-slate-800">{{ booking.booking_no }}</span>
        </div>
        <div class="text-right">
          <span class="text-slate-400">วันที่:</span>
          <span class="font-mono ml-1 text-slate-800">
            {{ new Date(booking.created_at).toLocaleDateString('th-TH', { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' }) }}
          </span>
        </div>
        <div>
          <span class="text-slate-400">ลูกค้า:</span>
          <span class="font-bold ml-1 text-slate-800">{{ booking.user?.display_name || 'ลูกค้าทั่วไป' }}</span>
        </div>
        <div class="text-right">
          <span class="text-slate-400">เบอร์โทร:</span>
          <span class="font-mono ml-1 text-slate-800">{{ booking.user?.phone || '-' }}</span>
        </div>
      </div>

      <!-- Items Table -->
      <div class="border-b border-slate-200 pb-3 mb-3">
        <table class="w-full text-left text-[11px]">
          <thead>
            <tr class="text-slate-400 border-b border-slate-100 pb-1">
              <th class="py-1">รายการ</th>
              <th class="py-1 text-center">จน.</th>
              <th class="py-1 text-right">ราคา</th>
              <th class="py-1 text-right">รวม</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-50">
            <!-- Court Slots -->
            <tr v-for="c in booking.courts" :key="c.id" class="text-slate-700">
              <td class="py-1.5 font-medium">
                ค่าคอร์ท {{ c.court?.court_number }} ({{ c.start_hour }}:00-{{ c.end_hour }}:00 น.)
              </td>
              <td class="py-1.5 text-center font-mono">1 ชม.</td>
              <td class="py-1.5 text-right font-mono">{{ c.hourly_rate }}</td>
              <td class="py-1.5 text-right font-mono font-semibold">{{ c.hourly_rate }}</td>
            </tr>

            <!-- Mini Store Items -->
            <tr v-for="item in booking.items" :key="item.id" class="text-slate-700">
              <td class="py-1.5 font-medium">{{ item.product?.name || 'สินค้า' }}</td>
              <td class="py-1.5 text-center font-mono">{{ item.qty }}</td>
              <td class="py-1.5 text-right font-mono">{{ item.unit_price }}</td>
              <td class="py-1.5 text-right font-mono font-semibold">{{ item.subtotal }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Totals Summary -->
      <div class="space-y-1.5 text-xs text-slate-700 border-b border-dashed border-slate-300 pb-4 mb-4">
        <div class="flex justify-between">
          <span class="text-slate-500">รวมค่าคอร์ท:</span>
          <span class="font-mono font-semibold">฿{{ booking.court_total?.toLocaleString() }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-slate-500">รวมสินค้ามินิสโตร์:</span>
          <span class="font-mono font-semibold">฿{{ booking.item_total?.toLocaleString() }}</span>
        </div>
        <div class="flex justify-between text-slate-500 text-[10px]">
          <span>ภาษีมูลค่าเพิ่ม (VAT 7% รวมแล้ว):</span>
          <span class="font-mono">฿{{ ((booking.net_total * 7) / 107).toFixed(2) }}</span>
        </div>
        <div class="flex justify-between text-sm font-extrabold text-slate-900 pt-1">
          <span>ยอดชำระสุทธิ (NET TOTAL):</span>
          <span class="font-mono text-emerald-700">฿{{ booking.net_total?.toLocaleString() }}</span>
        </div>
      </div>

      <!-- Payment & QR Footer -->
      <div class="text-center space-y-1 text-[10px] text-slate-400">
        <div class="font-semibold text-slate-600">
          ชำระผ่าน: พร้อมเพย์ (PromptPay) • สถานะ: {{ booking.status }}
        </div>
        <div v-if="booking.payment_slip?.trans_ref" class="font-mono text-[9px]">
          Trans Ref: {{ booking.payment_slip.trans_ref }}
        </div>
        <div class="pt-2 text-[10px] text-slate-500">
          ขอบพระคุณที่ใช้บริการสนามแบดมินตัน SmashHub 🏸
        </div>
      </div>
    </div>
  </el-dialog>
</template>
