<script setup>
import { ref, onMounted, computed } from 'vue'
import { useAdminStore } from '@/stores/adminStore'
import api from '@/api'
import { Plus, Minus, ShoppingCart } from 'lucide-vue-next'

const adminStore = useAdminStore()
const products = ref([])
const selectedItems = ref({}) // { [productId]: qty }
const isSubmitting = ref(false)

onMounted(async () => {
  try {
    const res = await api.get('/products')
    products.value = res.data.data || []
  } catch (err) {
    console.error('Error fetching products for POS:', err)
  }
})

function increment(product) {
  const current = selectedItems.value[product.id] || 0
  if (current < product.stock_qty) {
    selectedItems.value[product.id] = current + 1
  }
}

function decrement(product) {
  const current = selectedItems.value[product.id] || 0
  if (current > 1) {
    selectedItems.value[product.id] = current - 1
  } else {
    delete selectedItems.value[product.id]
  }
}

const totalAdditionAmount = computed(() => {
  let sum = 0
  for (const [pId, qty] of Object.entries(selectedItems.value)) {
    const p = products.value.find(prod => prod.id === Number(pId))
    if (p) sum += p.price * qty
  }
  return sum
})

const selectedCount = computed(() => {
  return Object.values(selectedItems.value).reduce((a, b) => a + b, 0)
})

async function submitAddItems() {
  if (!adminStore.posTargetBooking) return
  isSubmitting.value = true

  const items = Object.entries(selectedItems.value).map(([pId, qty]) => ({
    product_id: Number(pId),
    qty
  }))

  const success = await adminStore.addItems(adminStore.posTargetBooking.id, items)
  if (success) {
    selectedItems.value = {}
  }
  isSubmitting.value = false
}
</script>

<template>
  <el-dialog
    v-model="adminStore.isPOSOpen"
    title="เพิ่มสินค้าเข้าบิล (POS Counter)"
    width="540px"
    destroy-on-close
  >
    <div v-if="adminStore.posTargetBooking" class="space-y-4 text-xs">
      <!-- Target Booking Info -->
      <div class="bg-slate-50 p-3 rounded-xl border border-slate-200 flex justify-between items-center">
        <div>
          <div class="font-bold text-slate-800 text-sm font-mono">
            {{ adminStore.posTargetBooking.booking_no }}
          </div>
          <div class="text-[11px] text-slate-500">
            ลูกค้า: {{ adminStore.posTargetBooking.user?.display_name }}
          </div>
        </div>
        <div class="text-right">
          <div class="text-[10px] text-slate-400">ยอดบิลปัจจุบัน</div>
          <div class="text-sm font-bold text-emerald-700 font-mono">
            ฿{{ adminStore.posTargetBooking.net_total?.toLocaleString() }}
          </div>
        </div>
      </div>

      <!-- Product Selection List -->
      <div class="max-h-72 overflow-y-auto space-y-2 pr-1">
        <div
          v-for="product in products"
          :key="product.id"
          class="flex items-center justify-between p-2.5 bg-white border border-slate-200 rounded-xl hover:border-emerald-300 transition"
        >
          <div class="flex items-center space-x-2.5">
            <img
              :src="product.image_url || 'https://images.unsplash.com/photo-1626224583764-f87db24ac4ea?w=100'"
              class="w-10 h-10 rounded-lg object-cover bg-slate-100"
            />
            <div>
              <div class="font-semibold text-slate-800 text-xs">{{ product.name }}</div>
              <div class="text-[10px] text-slate-400">
                ฿{{ product.price }} • คงเหลือ {{ product.stock_qty }} ชิ้น
              </div>
            </div>
          </div>

          <!-- Stepper -->
          <div class="flex items-center space-x-1.5">
            <button
              v-if="selectedItems[product.id]"
              @click="decrement(product)"
              class="w-6 h-6 rounded-full bg-slate-100 text-slate-700 hover:bg-slate-200 flex items-center justify-center font-bold"
            >
              -
            </button>
            <span v-if="selectedItems[product.id]" class="w-6 text-center font-bold text-slate-800">
              {{ selectedItems[product.id] }}
            </span>
            <button
              @click="increment(product)"
              :disabled="product.stock_qty <= 0"
              class="w-6 h-6 rounded-full bg-emerald-600 text-white hover:bg-emerald-700 flex items-center justify-center font-bold disabled:opacity-40"
            >
              +
            </button>
          </div>
        </div>
      </div>

      <!-- POS Bottom Summary & Submit -->
      <div class="border-t border-slate-200 pt-3 flex justify-between items-center">
        <div>
          <div class="text-[11px] text-slate-500">เพิ่ม {{ selectedCount }} ชิ้น</div>
          <div class="text-sm font-bold text-emerald-800 font-mono">
            +฿{{ totalAdditionAmount.toLocaleString() }}
          </div>
        </div>

        <button
          @click="submitAddItems"
          :disabled="selectedCount === 0 || isSubmitting"
          class="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl font-bold text-xs shadow transition flex items-center space-x-1.5 disabled:opacity-40"
        >
          <ShoppingCart class="w-4 h-4" />
          <span>บันทึกและตัดสต็อก</span>
        </button>
      </div>
    </div>
  </el-dialog>
</template>
