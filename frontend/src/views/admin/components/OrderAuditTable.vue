<script setup>
import { onMounted } from 'vue'
import { useAdminStore } from '@/stores/adminStore'
import { Search, RotateCw, FileSearch, Printer, Plus, CheckCircle, Clock, AlertCircle } from 'lucide-vue-next'

const adminStore = useAdminStore()

onMounted(() => {
  adminStore.fetchBookings()
})

function formatCurrency(val) {
  if (val === undefined || val === null) return '0'
  return Number(val).toLocaleString()
}
</script>

<template>
  <div class="space-y-3.5">
    <!-- Filter Control Bar -->
    <div class="bg-white p-3.5 rounded-2xl border border-slate-200/90 shadow-2xs flex flex-wrap items-center justify-between gap-2.5">
      <div class="flex flex-wrap items-center gap-2">
        <!-- Search Input -->
        <el-input
          v-model="adminStore.searchQuery"
          placeholder="ค้นหาชื่อลูกค้า, เบอร์โทร หรือรหัสจอง..."
          clearable
          size="default"
          class="!w-60"
          @keyup.enter="adminStore.fetchBookings"
          @clear="adminStore.fetchBookings"
        >
          <template #prefix>
            <Search class="w-3.5 h-3.5 text-slate-400" />
          </template>
        </el-input>

        <!-- Status Filter -->
        <el-select
          v-model="adminStore.filterStatus"
          placeholder="ทุกสถานะการจอง"
          clearable
          size="default"
          class="!w-40"
          @change="adminStore.fetchBookings"
        >
          <el-option label="ทุกสถานะ" value="" />
          <el-option label="✅ CONFIRMED" value="CONFIRMED" />
          <el-option label="⏱️ PENDING_PAYMENT" value="PENDING_PAYMENT" />
          <el-option label="⚠️ MANUAL_REVIEW" value="MANUAL_REVIEW" />
          <el-option label="❌ EXPIRED" value="EXPIRED" />
          <el-option label="🚫 CANCELLED" value="CANCELLED" />
        </el-select>

        <!-- Court Filter -->
        <el-select
          v-model="adminStore.filterCourtId"
          placeholder="ทุกคอร์ท"
          clearable
          size="default"
          class="!w-32"
          @change="adminStore.fetchBookings"
        >
          <el-option label="ทุกคอร์ท" value="" />
          <el-option label="คอร์ท 1" :value="1" />
          <el-option label="คอร์ท 2" :value="2" />
          <el-option label="คอร์ท 3" :value="3" />
          <el-option label="คอร์ท 4" :value="4" />
          <el-option label="คอร์ท 5" :value="5" />
          <el-option label="คอร์ท 6" :value="6" />
        </el-select>

        <!-- Search Button -->
        <button
          @click="adminStore.fetchBookings"
          class="px-3.5 py-2 bg-slate-900 hover:bg-slate-800 text-white rounded-xl text-xs font-bold transition flex items-center space-x-1"
        >
          <Search class="w-3.5 h-3.5" />
          <span>ค้นหา</span>
        </button>
      </div>

      <div class="flex items-center space-x-2 text-xs text-slate-500 font-mono">
        <span>{{ adminStore.bookings.length }} records</span>
        <button
          @click="adminStore.fetchBookings"
          class="p-1.5 text-slate-400 hover:text-slate-700 hover:bg-slate-100 rounded-lg transition"
          title="รีเฟรชข้อมูล"
        >
          <RotateCw class="w-3.5 h-3.5" />
        </button>
      </div>
    </div>

    <!-- Element Plus Table with Stripe/Linear styling -->
    <div class="bg-white rounded-2xl border border-slate-200/90 shadow-2xs overflow-hidden">
      <el-table
        :data="adminStore.bookings"
        v-loading="adminStore.isLoading"
        style="width: 100%"
        class="text-xs font-sans"
        empty-text="ไม่พบข้อมูลการจองตามเงื่อนไขที่เลือก"
      >
        <!-- Booking No -->
        <el-table-column label="รหัสการจอง" min-width="140">
          <template #default="{ row }">
            <div class="font-mono font-bold text-slate-900 text-xs">{{ row.booking_no }}</div>
            <div class="text-[10px] text-slate-400 font-mono">
              {{ new Date(row.created_at).toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit' }) }} น.
            </div>
          </template>
        </el-table-column>

        <!-- Customer -->
        <el-table-column label="ลูกค้า / ผู้จอง" min-width="150">
          <template #default="{ row }">
            <div class="font-bold text-slate-900 text-xs">{{ row.user?.display_name || 'ลูกค้าทั่วไป' }}</div>
            <div class="text-[10px] text-slate-400 font-mono">{{ row.user?.phone || '-' }}</div>
          </template>
        </el-table-column>

        <!-- Court Slots -->
        <el-table-column label="คอร์ทและเวลา" min-width="170">
          <template #default="{ row }">
            <div class="flex flex-wrap gap-1">
              <span
                v-for="c in row.courts"
                :key="c.id"
                class="inline-block bg-slate-100 text-slate-800 border border-slate-200/80 px-1.5 py-0.2 rounded-md text-[10px] font-mono font-semibold"
              >
                คอร์ท {{ c.court?.court_number }} ({{ c.start_hour }}:00-{{ c.end_hour }}:00)
              </span>
            </div>
          </template>
        </el-table-column>

        <!-- Mini Store Items -->
        <el-table-column label="มินิสโตร์" min-width="120">
          <template #default="{ row }">
            <div v-if="row.items && row.items.length > 0" class="text-[11px] text-slate-700 font-medium">
              <span>{{ row.items.length }} รายการ</span>
              <span class="text-slate-400 font-mono"> (฿{{ formatCurrency(row.item_total) }})</span>
            </div>
            <span v-else class="text-[10px] text-slate-400">-</span>
          </template>
        </el-table-column>

        <!-- Total Net Amount -->
        <el-table-column label="ยอดสุทธิ" min-width="100" align="right">
          <template #default="{ row }">
            <div class="font-mono font-black text-sm text-emerald-700">
              ฿{{ formatCurrency(row.net_total) }}
            </div>
          </template>
        </el-table-column>

        <!-- Slip Verification Status Badge (4-Step Verification) -->
        <el-table-column label="ผลตรวจสลิป (4-Step)" min-width="170">
          <template #default="{ row }">
            <template v-if="row.payment_slip">
              <span
                :class="[
                  'inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-mono font-bold border',
                  row.payment_slip.verification_status === 'VERIFIED'
                    ? 'bg-emerald-50 text-emerald-700 border-emerald-300'
                    : row.payment_slip.verification_status === 'API_ERROR_PENDING_ADMIN'
                    ? 'bg-amber-50 text-amber-800 border-amber-300'
                    : 'bg-rose-50 text-rose-700 border-rose-300'
                ]"
              >
                <span
                  :class="[
                    'w-1.5 h-1.5 rounded-full',
                    row.payment_slip.verification_status === 'VERIFIED' ? 'bg-emerald-500' :
                    row.payment_slip.verification_status === 'API_ERROR_PENDING_ADMIN' ? 'bg-amber-500' : 'bg-rose-500'
                  ]"
                ></span>
                <span>{{ row.payment_slip.verification_status }}</span>
              </span>
            </template>
            <span v-else class="text-[10px] text-slate-400 font-mono">NO_SLIP</span>
          </template>
        </el-table-column>

        <!-- Booking Status Badge -->
        <el-table-column label="สถานะใบจอง" min-width="130">
          <template #default="{ row }">
            <span
              :class="[
                'inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-mono font-bold',
                row.status === 'CONFIRMED' ? 'bg-slate-900 text-white' :
                row.status === 'MANUAL_REVIEW' ? 'bg-amber-100 text-amber-900 border border-amber-300' :
                row.status === 'PENDING_PAYMENT' ? 'bg-blue-50 text-blue-800 border border-blue-200' :
                'bg-slate-100 text-slate-500'
              ]"
            >
              {{ row.status }}
            </span>
          </template>
        </el-table-column>

        <!-- Actions -->
        <el-table-column label="การจัดการ" width="160" align="center" fixed="right">
          <template #default="{ row }">
            <div class="flex items-center justify-center space-x-1">
              <!-- Check Slip Drawer -->
              <button
                @click="adminStore.openDrawer(row)"
                class="p-1.5 bg-slate-100 hover:bg-emerald-50 hover:text-emerald-700 text-slate-700 rounded-lg transition"
                title="ตรวจสอบรูปสลิปและ JSON"
              >
                <FileSearch class="w-3.5 h-3.5" />
              </button>

              <!-- POS Add Item -->
              <button
                v-if="row.status === 'CONFIRMED' || row.status === 'PENDING_PAYMENT'"
                @click="adminStore.openPOS(row)"
                class="p-1.5 bg-amber-50 hover:bg-amber-100 text-amber-700 rounded-lg transition"
                title="เพิ่มสินค้ามินิสโตร์เข้าบิลนี้"
              >
                <Plus class="w-3.5 h-3.5" />
              </button>

              <!-- Print Receipt -->
              <button
                @click="adminStore.openReceipt(row)"
                class="p-1.5 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-lg transition"
                title="พิมพ์ใบเสร็จ PDF"
              >
                <Printer class="w-3.5 h-3.5" />
              </button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>
