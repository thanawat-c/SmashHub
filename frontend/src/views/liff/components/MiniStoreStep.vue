<script setup>
import { onMounted } from 'vue'
import { useBookingStore } from '@/stores/bookingStore'
import { useLineProfileStore } from '@/stores/lineProfile'
import { ShoppingBag, ChevronLeft, ChevronRight, Plus, Minus, AlertCircle } from 'lucide-vue-next'

const bookingStore = useBookingStore()
const profileStore = useLineProfileStore()

onMounted(() => {
  bookingStore.fetchProducts()
})

async function onProceedToPayment() {
  const success = await bookingStore.submitBooking(profileStore.currentProfile)
  if (success) {
    // Moved to step 3 inside store
  }
}
</script>

<template>
  <div class="flex flex-col h-full bg-slate-50">
    <!-- Header banner -->
    <div class="bg-white px-4 py-3 border-b border-slate-200 shadow-2xs">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-sm font-bold text-slate-800 flex items-center space-x-1.5">
            <ShoppingBag class="w-4 h-4 text-emerald-600" />
            <span>มินิสโตร์ & บริการเช่าอุปกรณ์</span>
          </h2>
          <p class="text-[11px] text-slate-500">สั่งล่วงหน้า รับของพร้อมลงคอร์ทได้ทันที</p>
        </div>
        <div v-if="bookingStore.totalCartCount > 0" class="bg-emerald-100 text-emerald-800 text-xs font-bold px-2.5 py-1 rounded-full">
          ตะกร้า: {{ bookingStore.totalCartCount }} ชิ้น
        </div>
      </div>
    </div>

    <!-- Error Banner -->
    <div v-if="bookingStore.errorMessage" class="m-3 p-3 bg-rose-50 border border-rose-200 rounded-xl text-rose-700 text-xs flex items-center space-x-2">
      <AlertCircle class="w-4 h-4 flex-shrink-0" />
      <span>{{ bookingStore.errorMessage }}</span>
    </div>

    <!-- Products List -->
    <div class="flex-1 overflow-y-auto p-3 space-y-2.5">
      <div
        v-for="product in bookingStore.products"
        :key="product.id"
        class="bg-white border border-slate-200 rounded-xl p-3 shadow-2xs flex space-x-3 items-center hover:border-emerald-300 transition"
      >
        <!-- Product Thumbnail -->
        <img
          :src="product.image_url || 'https://images.unsplash.com/photo-1626224583764-f87db24ac4ea?w=200'"
          :alt="product.name"
          class="w-16 h-16 rounded-lg object-cover bg-slate-100 flex-shrink-0 border border-slate-100"
        />

        <!-- Product Details -->
        <div class="flex-1 min-w-0">
          <div class="flex items-center space-x-1.5 mb-0.5">
            <span
              :class="[
                'text-[9px] font-semibold px-1.5 py-0.5 rounded',
                product.category === 'SHUTTLECOCK' ? 'bg-amber-100 text-amber-800' :
                product.category === 'DRINK' ? 'bg-blue-100 text-blue-800' : 'bg-purple-100 text-purple-800'
              ]"
            >
              {{ product.category === 'SHUTTLECOCK' ? 'ลูกแบด' : product.category === 'DRINK' ? 'เครื่องดื่ม' : 'อุปกรณ์เช่า' }}
            </span>
            <span class="text-[10px] text-slate-400">คงเหลือ {{ product.stock_qty }} ชิ้น</span>
          </div>

          <h3 class="text-xs font-semibold text-slate-800 truncate leading-snug">{{ product.name }}</h3>
          <div class="text-xs font-bold text-emerald-700 mt-1">฿{{ product.price.toLocaleString() }}</div>
        </div>

        <!-- Add/Remove Counter Stepper -->
        <div class="flex items-center space-x-1.5 flex-shrink-0">
          <button
            v-if="bookingStore.cartItems[product.id]"
            @click="bookingStore.removeFromCart(product.id)"
            class="w-7 h-7 rounded-full bg-slate-100 hover:bg-slate-200 text-slate-700 flex items-center justify-center transition active:scale-95"
          >
            <Minus class="w-3.5 h-3.5" />
          </button>

          <span
            v-if="bookingStore.cartItems[product.id]"
            class="text-xs font-bold text-slate-800 w-5 text-center"
          >
            {{ bookingStore.cartItems[product.id].qty }}
          </span>

          <button
            @click="bookingStore.addToCart(product)"
            :disabled="product.stock_qty <= 0 || (bookingStore.cartItems[product.id]?.qty || 0) >= product.stock_qty"
            :class="[
              'w-7 h-7 rounded-full flex items-center justify-center transition active:scale-95 shadow-xs',
              product.stock_qty > 0
                ? 'bg-emerald-600 hover:bg-emerald-700 text-white'
                : 'bg-slate-200 text-slate-400 cursor-not-allowed'
            ]"
          >
            <Plus class="w-3.5 h-3.5" />
          </button>
        </div>
      </div>
    </div>

    <!-- Bottom Bill Summary Bar -->
    <div class="bg-white border-t border-slate-200 p-3 shadow-lg sticky bottom-0 z-20">
      <!-- Breakdown Summary -->
      <div class="space-y-1 mb-2.5 text-xs text-slate-600 border-b border-slate-100 pb-2">
        <div class="flex justify-between">
          <span>ค่าคอร์ท ({{ bookingStore.selectedSlots.length }} ชม.)</span>
          <span class="font-mono">฿{{ bookingStore.courtTotal.toLocaleString() }}</span>
        </div>
        <div class="flex justify-between">
          <span>สินค้ามินิสโตร์ ({{ bookingStore.totalCartCount }} ชิ้น)</span>
          <span class="font-mono">฿{{ bookingStore.itemTotal.toLocaleString() }}</span>
        </div>
        <div class="flex justify-between font-bold text-slate-900 text-sm pt-1">
          <span>ยอดชำระสุทธิ</span>
          <span class="text-emerald-700 font-mono">฿{{ bookingStore.netTotal.toLocaleString() }}</span>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="flex space-x-2">
        <button
          @click="bookingStore.currentStep = 1"
          class="px-3 py-2.5 rounded-xl border border-slate-300 text-slate-700 hover:bg-slate-50 font-semibold text-xs transition flex items-center space-x-1"
        >
          <ChevronLeft class="w-4 h-4" />
          <span>แก้ไขคอร์ท</span>
        </button>

        <button
          @click="onProceedToPayment"
          :disabled="bookingStore.isSubmitting"
          class="flex-1 px-4 py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl font-bold text-xs shadow-md shadow-emerald-500/20 transition flex items-center justify-center space-x-1.5 active:scale-95 disabled:opacity-50"
        >
          <span v-if="bookingStore.isSubmitting" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
          <span v-else>ยืนยันและไปชำระเงิน</span>
          <ChevronRight v-if="!bookingStore.isSubmitting" class="w-4 h-4" />
        </button>
      </div>
    </div>
  </div>
</template>
