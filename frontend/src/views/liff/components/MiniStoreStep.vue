<script setup>
import { ref, computed, onMounted } from 'vue'
import { useBookingStore } from '@/stores/bookingStore'
import { useLineProfileStore } from '@/stores/lineProfile'
import { ShoppingBag, ChevronLeft, ChevronRight, Plus, Minus, AlertCircle, Sparkles } from 'lucide-vue-next'

const bookingStore = useBookingStore()
const profileStore = useLineProfileStore()

const activeCategory = ref('ALL') // ALL, SHUTTLECOCK, DRINK, RENTAL

onMounted(() => {
  bookingStore.fetchProducts()
})

const filteredProducts = computed(() => {
  if (activeCategory.value === 'ALL') return bookingStore.products
  return bookingStore.products.filter(p => p.category === activeCategory.value)
})

const categories = [
  { id: 'ALL', label: 'ทั้งหมด' },
  { id: 'SHUTTLECOCK', label: '🏸 ลูกแบด' },
  { id: 'DRINK', label: '🥤 เครื่องดื่ม' },
  { id: 'RENTAL', label: '🎾 ไม้เช่า & กริป' }
]

async function onProceedToPayment() {
  await bookingStore.submitBooking(profileStore.currentProfile)
}
</script>

<template>
  <div class="flex flex-col h-full bg-[#F8FAFC] relative">
    <!-- Header banner -->
    <div class="bg-white px-3.5 py-2.5 border-b border-slate-200/80 shadow-2xs">
      <div class="flex items-center justify-between mb-2">
        <div>
          <h2 class="text-xs font-extrabold text-slate-900 flex items-center space-x-1.5">
            <ShoppingBag class="w-3.5 h-3.5 text-emerald-600" />
            <span>มินิสโตร์ & อุปกรณ์</span>
          </h2>
          <p class="text-[10px] text-slate-400">สั่งพร้อมคอร์ท รับสินค้าทันทีที่หน้าเคาน์เตอร์</p>
        </div>
        <div v-if="bookingStore.totalCartCount > 0" class="bg-emerald-50 text-emerald-700 border border-emerald-200 text-[11px] font-bold px-2 py-0.5 rounded-full font-mono">
          {{ bookingStore.totalCartCount }} ชิ้น
        </div>
      </div>

      <!-- Category Filter Pills -->
      <div class="flex space-x-1.5 overflow-x-auto pb-0.5 scrollbar-none">
        <button
          v-for="cat in categories"
          :key="cat.id"
          @click="activeCategory = cat.id"
          :class="[
            'px-2.5 py-1 rounded-full text-[11px] font-medium whitespace-nowrap transition flex-shrink-0 border',
            activeCategory === cat.id
              ? 'bg-slate-900 text-white border-slate-900 shadow-xs'
              : 'bg-slate-50 text-slate-600 border-slate-200 hover:bg-slate-100'
          ]"
        >
          {{ cat.label }}
        </button>
      </div>
    </div>

    <!-- Error Alert -->
    <div v-if="bookingStore.errorMessage" class="m-3 p-2.5 bg-rose-50 border border-rose-200 rounded-xl text-rose-700 text-xs flex items-center space-x-2">
      <AlertCircle class="w-4 h-4 flex-shrink-0" />
      <span>{{ bookingStore.errorMessage }}</span>
    </div>

    <!-- Products Grid -->
    <div class="flex-1 overflow-y-auto p-3 pb-28 space-y-2">
      <div
        v-for="product in filteredProducts"
        :key="product.id"
        class="bg-white border border-slate-200/90 rounded-2xl p-2.5 shadow-2xs flex space-x-3 items-center hover:border-emerald-300 transition"
      >
        <!-- Product Thumbnail with subtle rounded badge -->
        <div class="relative w-16 h-16 rounded-xl overflow-hidden bg-slate-100 flex-shrink-0 border border-slate-100">
          <img
            :src="product.image_url || 'https://images.unsplash.com/photo-1626224583764-f87db24ac4ea?w=200'"
            :alt="product.name"
            class="w-full h-full object-cover"
          />
          <div
            v-if="product.stock_qty <= 5 && product.stock_qty > 0"
            class="absolute bottom-0 inset-x-0 bg-rose-600 text-white text-[8px] font-bold text-center py-0.2"
          >
            เหลือ {{ product.stock_qty }}
          </div>
        </div>

        <!-- Product Details -->
        <div class="flex-1 min-w-0">
          <div class="flex items-center space-x-1.5 mb-0.5">
            <span
              :class="[
                'text-[8px] font-bold px-1.5 py-0.2 rounded font-mono',
                product.category === 'SHUTTLECOCK' ? 'bg-amber-100 text-amber-900' :
                product.category === 'DRINK' ? 'bg-blue-100 text-blue-900' : 'bg-purple-100 text-purple-900'
              ]"
            >
              {{ product.category === 'SHUTTLECOCK' ? 'ลูกแบด' : product.category === 'DRINK' ? 'เครื่องดื่ม' : 'อุปกรณ์' }}
            </span>
            <span class="text-[9px] text-slate-400 font-mono">คงเหลือ {{ product.stock_qty }} ชิ้น</span>
          </div>

          <h3 class="text-xs font-bold text-slate-900 truncate leading-snug">{{ product.name }}</h3>
          <div class="text-xs font-black text-emerald-700 font-mono mt-0.5">฿{{ product.price.toLocaleString() }}</div>
        </div>

        <!-- Stepper Counter Pill -->
        <div class="flex items-center space-x-1 flex-shrink-0 bg-slate-100 p-0.5 rounded-full border border-slate-200">
          <button
            v-if="bookingStore.cartItems[product.id]"
            @click="bookingStore.removeFromCart(product.id)"
            class="w-6 h-6 rounded-full bg-white text-slate-700 hover:bg-slate-200 flex items-center justify-center transition shadow-2xs active:scale-95"
          >
            <Minus class="w-3 h-3" />
          </button>

          <span
            v-if="bookingStore.cartItems[product.id]"
            class="text-xs font-bold font-mono text-slate-800 px-1 text-center"
          >
            {{ bookingStore.cartItems[product.id].qty }}
          </span>

          <button
            @click="bookingStore.addToCart(product)"
            :disabled="product.stock_qty <= 0 || (bookingStore.cartItems[product.id]?.qty || 0) >= product.stock_qty"
            :class="[
              'w-6 h-6 rounded-full flex items-center justify-center transition active:scale-95 shadow-2xs',
              product.stock_qty > 0
                ? 'bg-emerald-600 hover:bg-emerald-500 text-white'
                : 'bg-slate-200 text-slate-400 cursor-not-allowed'
            ]"
          >
            <Plus class="w-3 h-3" />
          </button>
        </div>
      </div>
    </div>

    <!-- Floating Bottom Checkout Dock (Clean Light Mode Standard) -->
    <div class="absolute bottom-2 left-2 right-2 z-20">
      <div class="bg-white/95 backdrop-blur-md text-slate-900 p-3 rounded-2xl shadow-lg border border-slate-200/90">
        <div class="flex items-center justify-between text-xs mb-2 pb-2 border-b border-slate-100">
          <div class="space-x-2 text-[10px] text-slate-500">
            <span>คอร์ท: <strong class="text-slate-900 font-mono">฿{{ bookingStore.courtTotal.toLocaleString() }}</strong></span>
            <span>•</span>
            <span>สินค้า: <strong class="text-slate-900 font-mono">฿{{ bookingStore.itemTotal.toLocaleString() }}</strong></span>
          </div>
          <div class="text-right">
            <span class="text-[10px] text-slate-500 mr-1">ยอดสุทธิ:</span>
            <span class="text-sm font-black text-emerald-700 font-mono">฿{{ bookingStore.netTotal.toLocaleString() }}</span>
          </div>
        </div>

        <div class="flex space-x-2">
          <button
            @click="bookingStore.currentStep = 1"
            class="px-3 py-2 rounded-xl bg-slate-100 hover:bg-slate-200 text-slate-700 text-xs font-semibold border border-slate-200 transition flex items-center space-x-1"
          >
            <ChevronLeft class="w-3.5 h-3.5" />
            <span>คอร์ท</span>
          </button>

          <button
            @click="onProceedToPayment"
            :disabled="bookingStore.isSubmitting"
            class="flex-1 py-2.5 bg-emerald-600 hover:bg-emerald-500 text-white rounded-xl font-bold text-xs shadow-xs transition flex items-center justify-center space-x-1.5 active:scale-95 disabled:opacity-50"
          >
            <span v-if="bookingStore.isSubmitting" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
            <span v-else>ยืนยันและไปชำระเงิน (฿{{ bookingStore.netTotal.toLocaleString() }})</span>
            <ChevronRight v-if="!bookingStore.isSubmitting" class="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
