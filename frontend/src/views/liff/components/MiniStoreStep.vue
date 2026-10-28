<script setup>
import { ref, computed, onMounted } from 'vue'
import { useBookingStore } from '@/stores/bookingStore'
import { useLineProfileStore } from '@/stores/lineProfile'
import { useAuthStore } from '@/stores/authStore'
import {
  ShoppingBag, Calendar, Clock, ChevronLeft, ChevronRight,
  Plus, Minus, AlertCircle, X, ChevronUp, ChevronDown, Check,
  Layers, CupSoda, Feather, Dumbbell, Trash2, ShieldCheck
} from 'lucide-vue-next'

const bookingStore = useBookingStore()
const profileStore = useLineProfileStore()
const authStore = useAuthStore()

const activeCategory = ref('ALL') // ALL, SHUTTLECOCK, DRINK, RENTAL
const showOrderReview = ref(false)
const failedImages = ref(new Set())

onMounted(() => {
  bookingStore.fetchProducts()
})

function handleImageError(productId) {
  failedImages.value.add(productId)
}

const filteredProducts = computed(() => {
  if (activeCategory.value === 'ALL') return bookingStore.products
  return bookingStore.products.filter(p => p.category === activeCategory.value)
})

const categories = [
  { id: 'ALL', label: 'All Items' },
  { id: 'SHUTTLECOCK', label: 'Shuttlecocks' },
  { id: 'DRINK', label: 'Beverages' },
  { id: 'RENTAL', label: 'Rackets & Grips' }
]

const formattedDate = computed(() => {
  if (!bookingStore.selectedDate) return 'Today'
  const d = new Date(bookingStore.selectedDate)
  return d.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric' })
})

// Group selected court slots by court
const groupedCourtSlots = computed(() => {
  const map = {}
  for (const slot of bookingStore.selectedSlots) {
    const key = slot.courtId
    if (!map[key]) {
      map[key] = {
        courtId: slot.courtId,
        courtNumber: slot.courtNumber,
        courtName: slot.courtName || `Court ${slot.courtNumber}`,
        slots: [],
        total: 0
      }
    }
    map[key].slots.push(slot)
    map[key].total += slot.hourlyRate
  }
  // Sort slots inside each court by startHour
  for (const k in map) {
    map[k].slots.sort((a, b) => a.startHour - b.startHour)
  }
  return Object.values(map)
})

async function onProceedToPayment() {
  const customer = authStore.user || profileStore.currentProfile
  await bookingStore.submitBooking(customer)
}
</script>

<template>
  <div class="flex flex-col h-full bg-[#F8FAFC] relative select-none">
    <!-- Header banner -->
    <div class="bg-white px-3.5 pt-3 pb-2.5 border-b border-slate-200/80 shadow-2xs">
      <div class="flex items-center justify-between mb-2.5">
        <div>
          <div class="flex items-center space-x-1.5">
            <h2 class="text-xs font-black text-slate-900 tracking-tight">Pro Shop & Add-ons</h2>
            <span class="text-[9px] font-bold px-1.5 py-0.2 rounded-full bg-emerald-50 text-emerald-700 border border-emerald-200">
              Step 2 of 4
            </span>
          </div>
          <p class="text-[10px] text-slate-400 mt-0.5">Collect drinks and shuttlecocks at the counter upon check-in</p>
        </div>

        <button
          @click="showOrderReview = true"
          class="flex items-center space-x-1.5 bg-slate-900 text-white px-2.5 py-1.5 rounded-xl text-[10px] font-bold shadow-xs hover:bg-slate-800 transition active:scale-95 cursor-pointer"
        >
          <ShoppingBag class="w-3 h-3 text-emerald-400" />
          <span>Cart</span>
          <span
            v-if="bookingStore.totalCartCount > 0"
            class="w-4 h-4 rounded-full bg-emerald-500 text-white text-[9px] font-black flex items-center justify-center font-mono ml-0.5"
          >
            {{ bookingStore.totalCartCount }}
          </span>
        </button>
      </div>

      <!-- Category Filter Pills -->
      <div class="flex space-x-1.5 overflow-x-auto pb-0.5 scrollbar-none">
        <button
          v-for="cat in categories"
          :key="cat.id"
          @click="activeCategory = cat.id"
          :class="[
            'px-2.5 py-1 rounded-full text-[11px] font-medium whitespace-nowrap transition flex-shrink-0 border cursor-pointer',
            activeCategory === cat.id
              ? 'bg-emerald-600 text-white border-emerald-600 shadow-xs font-bold'
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

    <!-- Main Content Area: Scrollable -->
    <div class="flex-1 overflow-y-auto p-3 pb-32 lg:pb-8 space-y-3">
      <!-- 1. Selected Booking Schedule Banner (Shown on mobile, desktop uses sticky sidebar) -->
      <div class="bg-white rounded-2xl border border-slate-200/90 p-3 shadow-2xs lg:hidden">
        <div class="flex items-center justify-between pb-2 mb-2 border-b border-slate-100">
          <div class="flex items-center space-x-1.5">
            <div class="w-5 h-5 rounded-lg bg-emerald-100 text-emerald-800 flex items-center justify-center">
              <Calendar class="w-3 h-3" />
            </div>
            <div>
              <span class="text-[11px] font-bold text-slate-900">Your Reserved Court(s)</span>
              <span class="text-[10px] text-slate-400 font-mono ml-1.5">{{ formattedDate }}</span>
            </div>
          </div>
          <button
            @click="bookingStore.currentStep = 1"
            class="text-[10px] font-bold text-emerald-600 hover:text-emerald-700 cursor-pointer flex items-center space-x-0.5 hover:underline"
          >
            <span>Change</span>
            <ChevronRight class="w-3 h-3" />
          </button>
        </div>

        <!-- Court slots list -->
        <div class="space-y-1.5">
          <div
            v-for="group in groupedCourtSlots"
            :key="group.courtId"
            class="flex items-center justify-between bg-slate-50/80 px-2.5 py-1.5 rounded-xl border border-slate-100 text-xs"
          >
            <div>
              <div class="font-bold text-slate-900 text-[11px] flex items-center space-x-1">
                <span>{{ group.courtName }}</span>
              </div>
              <div class="text-[10px] text-slate-500 font-mono mt-0.5 flex items-center space-x-1">
                <Clock class="w-2.5 h-2.5 text-slate-400" />
                <span>
                  {{ group.slots.map(s => `${String(s.startHour).padStart(2, '0')}:00-${String(s.endHour).padStart(2, '0')}:00`).join(', ') }}
                </span>
                <span>({{ group.slots.length }}h)</span>
              </div>
            </div>
            <div class="font-mono font-black text-slate-900 text-xs">
              ฿{{ group.total.toLocaleString() }}
            </div>
          </div>
        </div>

        <div class="flex items-center justify-between mt-2 pt-1.5 text-[10px] text-slate-400">
          <span>Courts Subtotal</span>
          <span class="font-bold font-mono text-slate-700 text-xs">฿{{ bookingStore.courtTotal.toLocaleString() }}</span>
        </div>
      </div>

      <!-- Section Title -->
      <div class="flex items-center justify-between pt-1">
        <h3 class="text-xs font-extrabold text-slate-900 uppercase tracking-wider">Available Items</h3>
        <span class="text-[10px] text-slate-400 font-mono">{{ filteredProducts.length }} items</span>
      </div>

      <!-- 2. Modern Product Cards Grid (1 col on mobile, 2 cols on tablet/desktop) -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-2.5">
        <div
          v-for="product in filteredProducts"
          :key="product.id"
          :class="[
            'bg-white border rounded-2xl p-2.5 shadow-2xs flex space-x-3 items-center transition relative',
            bookingStore.cartItems[product.id]?.qty > 0
              ? 'border-emerald-500 ring-1 ring-emerald-500/30 bg-emerald-50/15'
              : 'border-slate-200/90 hover:border-slate-300'
          ]"
        >
          <!-- Product Visual with Fail-Safe Icon Fallback (Solves Broken Image issue) -->
          <div class="relative w-16 h-16 rounded-xl overflow-hidden bg-slate-100 flex-shrink-0 border border-slate-200/60 flex items-center justify-center">
            <template v-if="!failedImages.has(product.id) && product.image_url">
              <img
                :src="product.image_url"
                :alt="product.name"
                class="w-full h-full object-cover"
                @error="handleImageError(product.id)"
              />
            </template>

            <!-- Reliable Fallback Visual Vector Graphic if image breaks -->
            <div
              v-else
              :class="[
                'w-full h-full flex flex-col items-center justify-center text-slate-600',
                product.category === 'SHUTTLECOCK' ? 'bg-amber-50 text-amber-700' :
                product.category === 'DRINK' ? 'bg-blue-50 text-blue-700' : 'bg-purple-50 text-purple-700'
              ]"
            >
              <Feather v-if="product.category === 'SHUTTLECOCK'" class="w-6 h-6 stroke-[1.8]" />
              <CupSoda v-else-if="product.category === 'DRINK'" class="w-6 h-6 stroke-[1.8]" />
              <Layers v-else class="w-6 h-6 stroke-[1.8]" />
            </div>

            <div
              v-if="product.stock_qty <= 5 && product.stock_qty > 0"
              class="absolute bottom-0 inset-x-0 bg-rose-600 text-white text-[8px] font-bold text-center py-0.2 tracking-tight"
            >
              Low: {{ product.stock_qty }} left
            </div>
            <div
              v-else-if="product.stock_qty <= 0"
              class="absolute inset-0 bg-slate-900/60 backdrop-blur-2xs flex items-center justify-center text-white text-[9px] font-bold"
            >
              Out of stock
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
                {{ product.category === 'SHUTTLECOCK' ? 'Shuttlecock' : product.category === 'DRINK' ? 'Beverage' : 'Equipment' }}
              </span>
              <span class="text-[9px] text-slate-400 font-mono">Stock: {{ product.stock_qty }}</span>
            </div>

            <h4 class="text-xs font-bold text-slate-900 truncate leading-snug">{{ product.name }}</h4>
            <div class="flex items-baseline space-x-1.5 mt-0.5">
              <span class="text-xs font-black text-emerald-700 font-mono">฿{{ product.price.toLocaleString() }}</span>
              <span v-if="bookingStore.cartItems[product.id]?.qty > 0" class="text-[10px] text-slate-400 font-mono">
                • Total: ฿{{ (product.price * bookingStore.cartItems[product.id].qty).toLocaleString() }}
              </span>
            </div>
          </div>

          <!-- Stepper Counter Pill -->
          <div class="flex items-center space-x-1 flex-shrink-0">
            <template v-if="bookingStore.cartItems[product.id]?.qty > 0">
              <div class="flex items-center space-x-1 bg-emerald-50 p-0.5 rounded-full border border-emerald-300">
                <button
                  @click="bookingStore.removeFromCart(product.id)"
                  class="w-6 h-6 rounded-full bg-white text-emerald-800 hover:bg-emerald-100 flex items-center justify-center transition shadow-2xs active:scale-95 cursor-pointer"
                >
                  <Minus class="w-3 h-3 stroke-[2.5]" />
                </button>

                <span class="text-xs font-black font-mono text-emerald-900 px-1 text-center min-w-[18px]">
                  {{ bookingStore.cartItems[product.id].qty }}
                </span>

                <button
                  @click="bookingStore.addToCart(product)"
                  :disabled="(bookingStore.cartItems[product.id]?.qty || 0) >= product.stock_qty"
                  class="w-6 h-6 rounded-full bg-emerald-600 hover:bg-emerald-500 text-white flex items-center justify-center transition shadow-2xs active:scale-95 disabled:opacity-50 cursor-pointer"
                >
                  <Plus class="w-3 h-3 stroke-[2.5]" />
                </button>
              </div>
            </template>

            <template v-else>
              <button
                @click="bookingStore.addToCart(product)"
                :disabled="product.stock_qty <= 0"
                class="px-2.5 py-1.5 rounded-xl bg-slate-100 hover:bg-emerald-50 hover:text-emerald-700 hover:border-emerald-300 text-slate-700 text-[11px] font-bold border border-slate-200 transition flex items-center space-x-1 cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed"
              >
                <Plus class="w-3 h-3 stroke-[2.5]" />
                <span>Add</span>
              </button>
            </template>
          </div>
        </div>
      </div>
    </div>

    <!-- 3. Floating Bottom Checkout Dock with Expandable Cart Details (Hidden on Desktop where Sidebar handles checkout) -->
    <div class="absolute bottom-2 left-2 right-2 z-20 lg:hidden">
      <div class="bg-white/95 backdrop-blur-md text-slate-900 p-3 rounded-2xl shadow-xl border border-slate-200">
        <!-- Collapsed Bar Info: Tappable to View Details -->
        <div
          @click="showOrderReview = !showOrderReview"
          class="flex items-center justify-between text-xs mb-2 pb-2 border-b border-slate-100 cursor-pointer hover:opacity-90 transition"
        >
          <div class="flex items-center space-x-1.5 text-[11px]">
            <div class="w-5 h-5 rounded-full bg-emerald-100 text-emerald-800 flex items-center justify-center font-bold text-[10px]">
              {{ bookingStore.selectedSlots.length + bookingStore.totalCartCount }}
            </div>
            <span class="text-slate-600">
              Courts (฿{{ bookingStore.courtTotal.toLocaleString() }})
              <span v-if="bookingStore.itemTotal > 0"> + Items (฿{{ bookingStore.itemTotal.toLocaleString() }})</span>
            </span>
          </div>

          <div class="flex items-center space-x-1.5">
            <span class="text-sm font-black text-emerald-700 font-mono">฿{{ bookingStore.netTotal.toLocaleString() }}</span>
            <button class="text-slate-400 hover:text-slate-600">
              <ChevronUp v-if="!showOrderReview" class="w-4 h-4" />
              <ChevronDown v-else class="w-4 h-4" />
            </button>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex space-x-2">
          <button
            @click="bookingStore.currentStep = 1"
            class="px-3 py-2.5 rounded-xl bg-slate-100 hover:bg-slate-200 text-slate-700 text-xs font-semibold border border-slate-200 transition flex items-center space-x-1 cursor-pointer"
          >
            <ChevronLeft class="w-3.5 h-3.5" />
            <span>Courts</span>
          </button>

          <button
            @click="onProceedToPayment"
            :disabled="bookingStore.isSubmitting"
            class="flex-1 py-2.5 bg-emerald-600 hover:bg-emerald-500 text-white rounded-xl font-bold text-xs shadow-xs transition flex items-center justify-center space-x-1.5 active:scale-95 disabled:opacity-50 cursor-pointer"
          >
            <span v-if="bookingStore.isSubmitting" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
            <span v-else>Proceed to Payment (฿{{ bookingStore.netTotal.toLocaleString() }})</span>
            <ChevronRight v-if="!bookingStore.isSubmitting" class="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>

    <!-- 4. Interactive Order Review Bottom Sheet Drawer (Comprehensive Overview) -->
    <div
      v-if="showOrderReview"
      class="fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex flex-col justify-end transition-opacity"
      @click.self="showOrderReview = false"
    >
      <div class="bg-white rounded-t-3xl max-h-[85vh] flex flex-col shadow-2xl border-t border-slate-200 max-w-md mx-auto w-full animate-in slide-in-from-bottom duration-200">
        <!-- Drawer Header -->
        <div class="px-4 py-3 border-b border-slate-100 flex items-center justify-between">
          <div class="flex items-center space-x-2">
            <div class="w-7 h-7 rounded-xl bg-emerald-50 text-emerald-700 flex items-center justify-center">
              <ShoppingBag class="w-4 h-4" />
            </div>
            <div>
              <h3 class="text-xs font-extrabold text-slate-900">Your Selection Overview</h3>
              <p class="text-[10px] text-slate-400">Review booked courts and selected add-ons</p>
            </div>
          </div>
          <button
            @click="showOrderReview = false"
            class="w-7 h-7 rounded-full bg-slate-100 hover:bg-slate-200 text-slate-500 flex items-center justify-center transition cursor-pointer"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <!-- Drawer Body -->
        <div class="p-4 overflow-y-auto space-y-3.5 text-xs">
          <!-- 1. Courts Section -->
          <div>
            <div class="flex items-center justify-between text-[11px] font-bold text-slate-900 mb-1.5">
              <span class="flex items-center space-x-1">
                <Calendar class="w-3.5 h-3.5 text-emerald-600" />
                <span>Booked Courts ({{ bookingStore.selectedSlots.length }} slot{{ bookingStore.selectedSlots.length > 1 ? 's' : '' }})</span>
              </span>
              <span class="font-mono text-emerald-700">฿{{ bookingStore.courtTotal.toLocaleString() }}</span>
            </div>

            <div class="space-y-1.5">
              <div
                v-for="group in groupedCourtSlots"
                :key="group.courtId"
                class="bg-slate-50 p-2.5 rounded-xl border border-slate-200/80"
              >
                <div class="flex items-center justify-between font-bold text-slate-900">
                  <span>{{ group.courtName }}</span>
                  <span class="font-mono">฿{{ group.total.toLocaleString() }}</span>
                </div>
                <div class="text-[10px] text-slate-500 font-mono mt-1 space-y-0.5">
                  <div v-for="s in group.slots" :key="s.startHour" class="flex justify-between">
                    <span>{{ String(s.startHour).padStart(2, '0') }}:00 - {{ String(s.endHour).padStart(2, '0') }}:00</span>
                    <span>฿{{ s.hourlyRate }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 2. Items Section -->
          <div>
            <div class="flex items-center justify-between text-[11px] font-bold text-slate-900 mb-1.5">
              <span class="flex items-center space-x-1">
                <ShoppingBag class="w-3.5 h-3.5 text-emerald-600" />
                <span>Add-on Products ({{ bookingStore.totalCartCount }} item{{ bookingStore.totalCartCount > 1 ? 's' : '' }})</span>
              </span>
              <span class="font-mono text-emerald-700">฿{{ bookingStore.itemTotal.toLocaleString() }}</span>
            </div>

            <div v-if="bookingStore.totalCartCount === 0" class="text-center py-4 bg-slate-50 rounded-xl border border-dashed border-slate-200 text-slate-400 text-xs">
              No add-ons selected yet. Choose drinks or shuttlecocks above!
            </div>

            <div v-else class="space-y-1.5">
              <div
                v-for="item in bookingStore.cartItemsList"
                :key="item.product.id"
                class="flex items-center justify-between bg-slate-50 p-2 rounded-xl border border-slate-200/80"
              >
                <div class="min-w-0 flex-1 pr-2">
                  <div class="font-bold text-slate-900 truncate">{{ item.product.name }}</div>
                  <div class="text-[10px] text-slate-500 font-mono">฿{{ item.product.price }} each</div>
                </div>

                <div class="flex items-center space-x-2">
                  <!-- Stepper -->
                  <div class="flex items-center space-x-1 bg-white px-1.5 py-0.5 rounded-full border border-slate-200 shadow-2xs">
                    <button
                      @click="bookingStore.removeFromCart(item.product.id)"
                      class="text-slate-600 hover:text-rose-600 transition"
                    >
                      <Minus class="w-3 h-3" />
                    </button>
                    <span class="font-mono font-bold text-slate-800 px-1 text-xs">{{ item.qty }}</span>
                    <button
                      @click="bookingStore.addToCart(item.product)"
                      :disabled="item.qty >= item.product.stock_qty"
                      class="text-slate-600 hover:text-emerald-600 transition disabled:opacity-40"
                    >
                      <Plus class="w-3 h-3" />
                    </button>
                  </div>

                  <!-- Subtotal -->
                  <div class="font-mono font-bold text-slate-900 min-w-[50px] text-right">
                    ฿{{ (item.product.price * item.qty).toLocaleString() }}
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 3. Final Total Calculation -->
          <div class="pt-2 border-t border-slate-200 space-y-1 font-mono text-xs">
            <div class="flex justify-between text-slate-500">
              <span>Courts Total:</span>
              <span>฿{{ bookingStore.courtTotal.toLocaleString() }}</span>
            </div>
            <div class="flex justify-between text-slate-500">
              <span>Add-ons Total:</span>
              <span>฿{{ bookingStore.itemTotal.toLocaleString() }}</span>
            </div>
            <div class="flex justify-between text-slate-900 font-bold text-sm pt-1 border-t border-slate-100">
              <span class="font-sans">Grand Total:</span>
              <span class="text-emerald-700 font-black">฿{{ bookingStore.netTotal.toLocaleString() }}</span>
            </div>
          </div>
        </div>

        <!-- Drawer Footer Action -->
        <div class="p-3 border-t border-slate-100 bg-slate-50/50">
          <button
            @click="showOrderReview = false; onProceedToPayment()"
            :disabled="bookingStore.isSubmitting"
            class="w-full py-2.5 bg-emerald-600 hover:bg-emerald-500 text-white rounded-xl font-bold text-xs shadow-xs transition flex items-center justify-center space-x-1.5 active:scale-95 cursor-pointer disabled:opacity-50"
          >
            <span>Confirm & Continue to Payment</span>
            <ChevronRight class="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
