<script setup>
import { ref, onMounted } from 'vue'
import { api } from '@/api'
import {
  ShieldAlert, Settings2, Package, Check, AlertCircle, Plus,
  Wrench, CheckCircle2, RefreshCw
} from 'lucide-vue-next'

const courts = ref([])
const products = ref([])
const isLoading = ref(false)
const message = ref('')

async function fetchData() {
  isLoading.value = true
  try {
    const [courtRes, prodRes] = await Promise.all([
      api.get('/courts/availability'),
      api.get('/products')
    ])
    courts.value = courtRes.data.data.courts || []
    products.value = prodRes.data.data || []
  } catch (err) {
    console.error('Fetch error:', err)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchData()
})

// Toggle court maintenance status
async function toggleCourtStatus(court) {
  const newStatus = court.status === 'MAINTENANCE' ? 'ACTIVE' : 'MAINTENANCE'
  try {
    await api.patch(`/admin/courts/${court.court_id}`, { status: newStatus })
    court.status = newStatus
    showMessage(`อัปเดตคอร์ท ${court.court_number} เป็น ${newStatus} สำเร็จ`)
  } catch (err) {
    showMessage(err.response?.data?.error || 'เกิดข้อผิดพลาดในการอัปเดตสถานะคอร์ท', true)
  }
}

// Quick restock
async function restockProduct(product, addQty) {
  try {
    const res = await api.post(`/admin/products/${product.id}/restock`, { add_qty: addQty })
    product.stock_qty = res.data.data.stock_qty
    showMessage(`เติมสต็อก ${product.name} (+${addQty}) สำเร็จ! คงเหลือ ${product.stock_qty} ชิ้น`)
  } catch (err) {
    showMessage(err.response?.data?.error || 'เกิดข้อผิดพลาดในการเติมสต็อก', true)
  }
}

const isErrorMsg = ref(false)
function showMessage(msg, isError = false) {
  message.value = msg
  isErrorMsg.value = isError
  setTimeout(() => {
    message.value = ''
  }, 3500)
}
</script>

<template>
  <div class="space-y-6">
    <!-- Feedback banner -->
    <div
      v-if="message"
      :class="[
        'p-3.5 rounded-2xl text-xs font-semibold flex items-center space-x-2 transition shadow-xs animate-in fade-in',
        isErrorMsg ? 'bg-rose-50 text-rose-800 border border-rose-200' : 'bg-emerald-50 text-emerald-800 border border-emerald-200'
      ]"
    >
      <CheckCircle2 v-if="!isErrorMsg" class="w-4 h-4 text-emerald-600 flex-shrink-0" />
      <AlertCircle v-else class="w-4 h-4 text-rose-600 flex-shrink-0" />
      <span>{{ message }}</span>
    </div>

    <!-- Section 1: Court Status & Maintenance -->
    <div class="bg-white rounded-2xl border border-slate-200/90 p-4 sm:p-5 shadow-xs space-y-4">
      <div class="flex items-center justify-between border-b border-slate-100 pb-3">
        <div class="flex items-center space-x-2.5">
          <div class="w-8 h-8 rounded-xl bg-slate-900 text-white flex items-center justify-center text-sm font-bold">
            🏸
          </div>
          <div>
            <h2 class="text-sm font-extrabold text-slate-900">จัดการสนาม & ปิดปรับปรุง (Court Maintenance)</h2>
            <p class="text-xs text-slate-500">เปิด-ปิดสถานะคอร์ท 1-6 เพื่อปิดซ่อมบำรุงหรือเปลี่ยนราคาค่าบริการ</p>
          </div>
        </div>
        <button
          @click="fetchData"
          class="p-2 rounded-xl text-slate-400 hover:text-slate-700 hover:bg-slate-100 transition cursor-pointer"
          title="รีเฟรชข้อมูล"
        >
          <RefreshCw :class="['w-4 h-4', isLoading ? 'animate-spin' : '']" />
        </button>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
        <div
          v-for="court in courts"
          :key="court.court_id"
          class="p-3.5 rounded-2xl border border-slate-200 bg-slate-50/70 space-y-3"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-2">
              <span class="font-extrabold text-sm text-slate-900">คอร์ท {{ court.court_number }}</span>
              <span
                :class="[
                  'text-[9px] font-bold px-1.5 py-0.2 rounded-md',
                  court.court_type === 'ACRYLIC' ? 'bg-blue-100 text-blue-700' : 'bg-slate-200 text-slate-700'
                ]"
              >
                {{ court.court_type === 'ACRYLIC' ? 'VIP' : 'PVC' }}
              </span>
            </div>

            <!-- Status Indicator Badge -->
            <span
              :class="[
                'text-[10px] font-bold px-2 py-0.5 rounded-full border',
                court.status === 'MAINTENANCE'
                  ? 'bg-rose-50 text-rose-700 border-rose-200'
                  : 'bg-emerald-50 text-emerald-700 border-emerald-200'
              ]"
            >
              {{ court.status === 'MAINTENANCE' ? '⚠️ ปิดปรับปรุง' : '✅ พร้อมใช้งาน' }}
            </span>
          </div>

          <div class="text-xs text-slate-500 flex justify-between font-mono">
            <span>Normal: ฿{{ court.normal_rate }}/ชม.</span>
            <span>Peak: ฿{{ court.peak_rate }}/ชม.</span>
          </div>

          <!-- Toggle Button -->
          <button
            @click="toggleCourtStatus(court)"
            :class="[
              'w-full py-1.5 rounded-xl text-xs font-bold transition flex items-center justify-center space-x-1.5 cursor-pointer shadow-2xs',
              court.status === 'MAINTENANCE'
                ? 'bg-emerald-600 hover:bg-emerald-700 text-white'
                : 'bg-rose-50 hover:bg-rose-100 text-rose-700 border border-rose-200'
            ]"
          >
            <Wrench class="w-3.5 h-3.5" />
            <span>{{ court.status === 'MAINTENANCE' ? 'เปิดให้บริการตามปกติ' : 'สั่งปิดปรับปรุงสนาม (Maintenance)' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Section 2: Mini-Store Product Restock -->
    <div class="bg-white rounded-2xl border border-slate-200/90 p-4 sm:p-5 shadow-xs space-y-4">
      <div class="border-b border-slate-100 pb-3 flex items-center justify-between">
        <div class="flex items-center space-x-2.5">
          <div class="w-8 h-8 rounded-xl bg-emerald-600 text-white flex items-center justify-center text-sm font-bold">
            <Package class="w-4 h-4" />
          </div>
          <div>
            <h2 class="text-sm font-extrabold text-slate-900">จัดการคลังสินค้ามินิสโตร์ (Product Inventory Restock)</h2>
            <p class="text-xs text-slate-500">เติมสต็อกสินค้าเครื่องดื่ม ลูกแบดมินตัน และอุปกรณ์เช่าหน้าร้าน</p>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
        <div
          v-for="product in products"
          :key="product.id"
          class="p-3.5 rounded-2xl border border-slate-200/90 bg-white space-y-2.5 shadow-2xs flex flex-col justify-between"
        >
          <div>
            <div class="flex items-center justify-between mb-1">
              <span class="text-[10px] font-mono text-slate-400 uppercase">{{ product.category }}</span>
              <span
                :class="[
                  'text-[10px] font-bold font-mono px-2 py-0.2 rounded-full',
                  product.stock_qty <= 5 ? 'bg-rose-50 text-rose-700 border border-rose-200' : 'bg-slate-100 text-slate-700'
                ]"
              >
                คงเหลือ {{ product.stock_qty }} ชิ้น
              </span>
            </div>
            <h3 class="text-xs font-bold text-slate-900 line-clamp-1">{{ product.name }}</h3>
            <p class="text-xs font-mono font-black text-emerald-700 mt-0.5">฿{{ product.price }}</p>
          </div>

          <!-- Quick Restock Buttons (+5, +10, +24) -->
          <div class="pt-2 border-t border-slate-100">
            <div class="text-[10px] text-slate-400 font-medium mb-1">เติมสต็อกด่วน</div>
            <div class="grid grid-cols-3 gap-1">
              <button
                @click="restockProduct(product, 5)"
                class="py-1 bg-slate-100 hover:bg-emerald-50 hover:text-emerald-700 text-slate-700 rounded-lg text-[10px] font-bold font-mono transition cursor-pointer border border-slate-200"
              >
                +5
              </button>
              <button
                @click="restockProduct(product, 10)"
                class="py-1 bg-slate-100 hover:bg-emerald-50 hover:text-emerald-700 text-slate-700 rounded-lg text-[10px] font-bold font-mono transition cursor-pointer border border-slate-200"
              >
                +10
              </button>
              <button
                @click="restockProduct(product, 24)"
                class="py-1 bg-slate-100 hover:bg-emerald-50 hover:text-emerald-700 text-slate-700 rounded-lg text-[10px] font-bold font-mono transition cursor-pointer border border-slate-200"
              >
                +24
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
