<script setup>
import { ref } from 'vue'
import api from '@/api'
import { ElNotification } from 'element-plus'
import { UserPlus } from 'lucide-vue-next'

const props = defineProps({
  visible: { type: Boolean, default: false },
  court: { type: Object, default: () => null },
  slot: { type: Object, default: () => null },
  playDate: { type: String, default: '' }
})

const emit = defineEmits(['update:visible', 'booking-created'])

const customerName = ref('ลูกค้า Walk-in หน้าร้าน')
const customerPhone = ref('080-000-0000')
const isSubmitting = ref(false)

async function handleCreateWalkIn() {
  if (!props.court || !props.slot) return
  isSubmitting.value = true

  try {
    const payload = {
      customer_name: customerName.value,
      customer_phone: customerPhone.value,
      play_date: props.playDate,
      court_slots: [
        { court_id: props.court.court_id, start_hour: props.slot.start_hour }
      ],
      items: []
    }

    const res = await api.post('/bookings', payload)
    const booking = res.data.data.booking

    // Immediately approve walk-in payment
    await api.post(`/bookings/${booking.id}/verify-slip`, {
      mock_scenario: 'MOCK_VALID'
    })

    ElNotification({
      title: 'เปิดคอร์ทสำเร็จ',
      message: `เปิดคอร์ท ${props.court.court_number} เวลา ${props.slot.start_hour}:00 น. เรียบร้อยแล้ว`,
      type: 'success'
    })

    emit('update:visible', false)
    emit('booking-created')
  } catch (err) {
    ElNotification({
      title: 'ไม่สามารถเปิดคอร์ทได้',
      message: err.response?.data?.error || 'เกิดข้อผิดพลาด',
      type: 'error'
    })
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <el-dialog
    :model-value="visible"
    @update:model-value="emit('update:visible', $event)"
    title="เปิดคอร์ทรับลูกค้า Walk-in หน้าร้าน"
    width="440px"
    destroy-on-close
  >
    <div v-if="court && slot" class="space-y-3.5 text-xs">
      <div class="bg-emerald-50 border border-emerald-200 rounded-xl p-3 flex justify-between items-center text-emerald-900">
        <div>
          <div class="font-bold text-sm">คอร์ท {{ court.court_number }} ({{ court.court_type }})</div>
          <div class="text-xs text-emerald-700">เวลา: {{ slot.start_hour }}:00 - {{ slot.end_hour }}:00 น.</div>
        </div>
        <div class="text-right">
          <div class="text-[10px] text-emerald-600">ค่าคอร์ท</div>
          <div class="font-mono font-bold text-base text-emerald-800">฿{{ slot.hourly_rate }}</div>
        </div>
      </div>

      <div>
        <label class="block text-[11px] font-semibold text-slate-700 mb-1">ชื่อลูกค้า / ก๊วน</label>
        <el-input v-model="customerName" placeholder="ระบุชื่อลูกค้าหรือชื่อก๊วน" size="default" />
      </div>

      <div>
        <label class="block text-[11px] font-semibold text-slate-700 mb-1">เบอร์โทรศัพท์</label>
        <el-input v-model="customerPhone" placeholder="08x-xxx-xxxx" size="default" />
      </div>

      <div class="pt-2 flex justify-end space-x-2">
        <button
          @click="emit('update:visible', false)"
          class="px-3 py-2 border border-slate-300 rounded-lg text-slate-600 hover:bg-slate-50 font-medium"
        >
          ยกเลิก
        </button>
        <button
          @click="handleCreateWalkIn"
          :disabled="isSubmitting"
          class="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg font-bold shadow flex items-center space-x-1.5 transition disabled:opacity-50"
        >
          <UserPlus class="w-4 h-4" />
          <span>บันทึกเปิดคอร์ททันที</span>
        </button>
      </div>
    </div>
  </el-dialog>
</template>
