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
  <div class="space-y-4">
    <!-- Filter Control Bar -->
    <div class="bg-white p-4 rounded-2xl border border-slate-200 shadow-xs flex flex-wrap items-center justify-between gap-3">
      <div class="flex flex-wrap items-center gap-2.5">
        <!-- Search Input -->
        <el-input
          v-model="adminStore.searchQuery"
          placeholder="ค้นหาชื่อลูกค้า, เบอร์โทร หรือรหัสจอง..."
          clearable
          size="default"
          class="!w-64"
          @keyup.enter="adminStore.fetchBookings"
          @clear="adminStore.fetchBookings"
        >
          <template #prefix>
            <Search class="w-4 h-4 text-slate-400" />
          </template>
        </el-input>

        <!-- Status Filter -->
        <el-select
          v-model="adminStore.filterStatus"
          placeholder="ทุกสถานะการจอง"
          clearable
          size="default"
          class="!w-44"
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
          class="!w-36"
          @change="adminStore.fetchBookings"
        >
          <el-option label="ทุกคอร์ท" value="" />
          <el-option label="คอร์ท 1 (PVC)" :value="1" />
          <el-option label="คอร์ท 2 (PVC)" :value="2" />
          <el-option label="คอร์ท 3 (PVC)" :value="3" />
          <el-option label="คอร์ท 4 (PVC)" :value="4" />
          <el-option label="คอร์ท 5 (VIP)" :value="5" />
          <el-option label="คอร์ท 6 (VIP)" :value="6" />
        </el-select>

        <!-- Search Button -->
        <button
          @click="adminStore.fetchBookings"
          class="px-3.5 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg text-xs font-semibold shadow-xs transition flex items-center space-x-1"
        >
          <Search class="w-3.5 h-3.5" />
          <span>ค้นหา</span>
        </button>
      </div>

      <div class="flex items-center space-x-2 text-xs text-slate-500">
        <span>พบทั้งหมด {{ adminStore.bookings.length }} รายการ</span>
        <button
          @click="adminStore.fetchBookings"
          class="p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-lg transition"
          title="รีเฟรชข้อมูล"
        >
          <RotateCw class="w-4 h-4" />
        </button>
      </div>
    </div>

    <!-- Element Plus Table -->
    <div class="bg-white rounded-2xl border border-slate-200 shadow-xs overflow-hidden">
      <el-table
        :data="adminStore.bookings"
        v-loading="adminStore.isLoading"
        style="width: 100%"
        class="text-xs"
        empty-text="ไม่พบข้อมูลการจองตามเงื่อนไขที่เลือก"
      >
        <!-- Booking No -->
        <el-table-column label="รหัสการจอง" min-width="150">
          <template #default="{ row }">
            <div class="font-mono font-bold text-slate-900">{{ row.booking_no }}</div>
            <div class="text-[10px] text-slate-400 font-mono">
              {{ new Date(row.created_at).toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit' }) }} น.
            </div>
          </template>
        </el-table-column>

        <!-- Customer -->
        <el-table-column label="ลูกค้า / ผู้จอง" min-width="160">
          <template #default="{ row }">
            <div class="font-semibold text-slate-800">{{ row.user?.display_name || 'ลูกค้าทั่วไป' }}</div>
            <div class="text-[11px] text-slate-400 font-mono">{{ row.user?.phone || '-' }}</div>
          </template>
        </el-table-column>

        <!-- Court Slots -->
        <el-table-column label="คอร์ทและเวลา" min-width="170">
          <template #default="{ row }">
            <div class="space-y-0.5">
              <span
                v-for="c in row.courts"
                :key="c.id"
                class="inline-block bg-slate-100 text-slate-700 px-1.5 py-0.5 rounded text-[11px] mr-1 font-mono"
              >
                คอร์ท {{ c.court?.court_number }} ({{ c.start_hour }}:00-{{ c.end_hour }}:00)
              </span>
            </div>
          </template>
        </el-table-column>

        <!-- Mini Store Items -->
        <el-table-column label="สินค้ามินิสโตร์" min-width="130">
          <template #default="{ row }">
            <div v-if="row.items && row.items.length > 0" class="text-[11px] text-slate-600">
              <span class="font-semibold text-emerald-700">{{ row.items.length }} รายการ</span>
              <span class="text-slate-400 font-mono"> (฿{{ formatCurrency(row.item_total) }})</span>
            </div>
            <span v-else class="text-[11px] text-slate-400">-</span>
          </template>
        </el-table-column>

        <!-- Total Net Amount -->
        <el-table-column label="ยอดสุทธิ" min-width="110" align="right">
          <template #default="{ row }">
            <div class="font-mono font-bold text-sm text-emerald-700">
              ฿{{ formatCurrency(row.net_total) }}
            </div>
          </template>
        </el-table-column>

        <!-- Slip Verification Status Badge -->
        <el-table-column label="สถานะสลิป (4-Step)" min-width="160">
          <template #default="{ row }">
            <template v-if="row.payment_slip">
              <el-tag
                :type="
                  row.payment_slip.verification_status === 'VERIFIED' ? 'success' :
                  row.payment_slip.verification_status === 'API_ERROR_PENDING_ADMIN' ? 'warning' : 'danger'
                "
                effect="light"
                size="small"
                class="font-mono text-[10px]"
              >
                {{ row.payment_slip.verification_status }}
              </el-tag>
            </template>
            <span v-else class="text-[10px] text-slate-400">ยังไม่แนบสลิป</span>
          </template>
        </el-table-column>

        <!-- Booking Status Badge -->
        <el-table-column label="สถานะใบจอง" min-width="130">
          <template #default="{ row }">
            <el-tag
              :type="
                row.status === 'CONFIRMED' ? 'success' :
                row.status === 'MANUAL_REVIEW' ? 'warning' :
                row.status === 'PENDING_PAYMENT' ? 'info' : 'danger'
              "
              effect="dark"
              size="small"
              class="font-mono text-[10px]"
            >
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>

        <!-- Action Buttons -->
        <el-table-column label="การจัดการ" width="180" align="center" fixed="right">
          <template #default="{ row }">
            <div class="flex items-center justify-center space-x-1">
              <!-- Check Slip Drawer -->
              <button
                @click="adminStore.openDrawer(row)"
                class="p-1.5 bg-slate-100 hover:bg-emerald-50 hover:text-emerald-700 text-slate-700 rounded-lg transition"
                title="ตรวจสอบรูปสลิปและข้อมูล JSON"
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
