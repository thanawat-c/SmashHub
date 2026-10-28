<script setup>
import { computed } from 'vue'
import { useBookingStore } from '@/stores/bookingStore'
import { useLineProfileStore } from '@/stores/lineProfile'
import { useAuthStore } from '@/stores/authStore'
import {
  Calendar, Clock, ShoppingBag, ChevronRight, Plus, Minus,
  Trash2, ShieldCheck, CheckCircle2, AlertCircle, ArrowRight
} from 'lucide-vue-next'

const bookingStore = useBookingStore()
const profileStore = useLineProfileStore()
const authStore = useAuthStore()

const formattedDate = computed(() => {
  if (!bookingStore.selectedDate) return 'Today'
  const d = new Date(bookingStore.selectedDate)
  return d.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric', year: 'numeric' })
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
        courtType: slot.courtType || 'PVC',
        slots: [],
        total: 0
      }
    }
    map[key].slots.push(slot)
    map[key].total += slot.hourlyRate
  }
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
  <div class="h-full flex flex-col justify-between p-5 bg-white select-none">
    <!-- Top Header -->
    <div>
      <div class="flex items-center justify-between pb-3 border-b border-slate-100">
        <div>
          <div class="flex items-center space-x-1.5">
            <h3 class="text-sm font-extrabold text-slate-900 tracking-tight">Booking Summary</h3>
            <span class="text-[9px] font-bold px-2 py-0.5 rounded-full bg-emerald-50 text-emerald-700 border border-emerald-200 font-mono">
              Step {{ bookingStore.currentStep }}/4
            </span>
          </div>
          <p class="text-[11px] text-slate-400 mt-0.5">SmashHub Badminton Arena • 6 Courts</p>
        </div>

        <button
          v-if="bookingStore.selectedSlots.length > 0 && bookingStore.currentStep <= 2"
          @click="bookingStore.resetSelection()"
          class="text-[10px] text-slate-400 hover:text-rose-600 transition font-medium cursor-pointer"
        >
          Clear
        </button>
      </div>

      <!-- Scrollable Selection List -->
      <div class="py-3 space-y-4 max-h-[calc(100vh-380px)] overflow-y-auto pr-1">
        <!-- 1. Booked Courts Section -->
        <div>
          <div class="flex items-center justify-between text-xs font-bold text-slate-900 mb-2">
            <div class="flex items-center space-x-1.5">
              <Calendar class="w-3.5 h-3.5 text-emerald-600" />
              <span>Court Reservation</span>
            </div>
            <span class="text-[10px] text-slate-400 font-mono">{{ formattedDate }}</span>
          </div>

          <!-- Empty state when no court selected -->
          <div
            v-if="bookingStore.selectedSlots.length === 0"
            class="p-4 rounded-2xl border-2 border-dashed border-slate-200 bg-slate-50/70 text-center"
          >
            <Clock class="w-6 h-6 text-slate-300 mx-auto mb-1.5 stroke-[1.5]" />
            <p class="text-xs font-semibold text-slate-600">No courts selected yet</p>
            <p class="text-[10px] text-slate-400 mt-0.5">Pick available hours from the schedule on the left</p>
          </div>

          <!-- Grouped courts cards -->
          <div v-else class="space-y-2">
            <div
              v-for="group in groupedCourtSlots"
              :key="group.courtId"
              class="bg-slate-50/90 border border-slate-200/90 rounded-2xl p-3 shadow-2xs"
            >
              <div class="flex items-center justify-between mb-1.5">
                <div class="flex items-center space-x-1.5">
                  <span class="font-bold text-xs text-slate-900">{{ group.courtName }}</span>
                  <span class="text-[9px] font-mono font-bold px-1.5 py-0.2 rounded bg-slate-200 text-slate-700">
                    {{ group.courtType }}
                  </span>
                </div>
                <span class="font-mono font-black text-xs text-slate-900">฿{{ group.total.toLocaleString() }}</span>
              </div>

              <!-- Slots pill list -->
              <div class="space-y-1">
                <div
                  v-for="s in group.slots"
                  :key="s.startHour"
                  class="flex items-center justify-between bg-white px-2 py-1 rounded-lg border border-slate-100 text-[11px]"
                >
                  <div class="flex items-center space-x-1.5 text-slate-600 font-mono">
                    <Clock class="w-3 h-3 text-slate-400" />
                    <span>{{ String(s.startHour).padStart(2, '0') }}:00 - {{ String(s.endHour).padStart(2, '0') }}:00</span>
                  </div>
                  <div class="flex items-center space-x-2">
                    <span class="font-mono text-slate-500 text-[10px]">฿{{ s.hourlyRate }}</span>
                    <button
                      v-if="bookingStore.currentStep <= 2"
                      @click="bookingStore.toggleCourtSlot({ court_id: s.courtId }, { start_hour: s.startHour, status: 'AVAILABLE' })"
                      class="text-slate-400 hover:text-rose-600 cursor-pointer text-xs leading-none"
                    >
                      ×
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 2. Add-on Products Section -->
        <div>
          <div class="flex items-center justify-between text-xs font-bold text-slate-900 mb-2">
            <div class="flex items-center space-x-1.5">
              <ShoppingBag class="w-3.5 h-3.5 text-emerald-600" />
              <span>Pro Shop Add-ons</span>
            </div>
            <span v-if="bookingStore.totalCartCount > 0" class="text-[10px] font-bold text-emerald-700 font-mono">
              {{ bookingStore.totalCartCount }} item(s)
            </span>
          </div>

          <!-- Empty state -->
          <div
            v-if="bookingStore.totalCartCount === 0"
            class="px-3 py-2.5 rounded-xl bg-slate-50 border border-slate-100 text-slate-400 text-[11px] text-center"
          >
            <span v-if="bookingStore.currentStep === 1">Shuttlecocks & beverages can be added in Step 2</span>
            <span v-else>No add-on items selected yet</span>
          </div>

          <!-- Cart items list -->
          <div v-else class="space-y-1.5">
            <div
              v-for="item in bookingStore.cartItemsList"
              :key="item.product.id"
              class="flex items-center justify-between bg-slate-50/90 border border-slate-200/80 p-2 rounded-xl text-xs"
            >
              <div class="min-w-0 flex-1 pr-2">
                <div class="font-bold text-slate-900 text-[11px] truncate">{{ item.product.name }}</div>
                <div class="text-[10px] text-slate-400 font-mono">฿{{ item.product.price }} each</div>
              </div>

              <div class="flex items-center space-x-2">
                <!-- Stepper -->
                <div class="flex items-center space-x-1 bg-white px-1.5 py-0.5 rounded-full border border-slate-200 shadow-2xs">
                  <button
                    @click="bookingStore.removeFromCart(item.product.id)"
                    class="text-slate-600 hover:text-rose-600 cursor-pointer"
                  >
                    <Minus class="w-2.5 h-2.5" />
                  </button>
                  <span class="font-mono font-bold text-slate-800 px-1 text-[11px]">{{ item.qty }}</span>
                  <button
                    @click="bookingStore.addToCart(item.product)"
                    :disabled="item.qty >= item.product.stock_qty"
                    class="text-slate-600 hover:text-emerald-600 cursor-pointer disabled:opacity-30"
                  >
                    <Plus class="w-2.5 h-2.5" />
                  </button>
                </div>

                <!-- Subtotal -->
                <span class="font-mono font-bold text-slate-900 text-xs min-w-[48px] text-right">
                  ฿{{ (item.product.price * item.qty).toLocaleString() }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Pricing Breakdown & CTA -->
    <div class="pt-3 border-t border-slate-200 space-y-3">
      <!-- Breakdown Table -->
      <div class="space-y-1.5 text-xs font-mono">
        <div class="flex justify-between text-slate-500">
          <span class="font-sans">Court Subtotal</span>
          <span>฿{{ bookingStore.courtTotal.toLocaleString() }}</span>
        </div>
        <div class="flex justify-between text-slate-500">
          <span class="font-sans">Add-ons Subtotal</span>
          <span>฿{{ bookingStore.itemTotal.toLocaleString() }}</span>
        </div>
        <div class="flex justify-between text-slate-400 text-[10px]">
          <span class="font-sans">Estimated VAT (7% incl.)</span>
          <span>฿{{ ((bookingStore.netTotal * 7) / 107).toFixed(2) }}</span>
        </div>
        <div class="flex items-baseline justify-between pt-2 border-t border-slate-100 text-slate-900 font-bold">
          <span class="font-sans text-sm font-extrabold">Net Total</span>
          <span class="text-lg font-black text-emerald-700">฿{{ bookingStore.netTotal.toLocaleString() }}</span>
        </div>
      </div>

      <!-- Action Buttons depending on Step -->
      <div>
        <!-- Step 1: Court selection -> go to Step 2 -->
        <button
          v-if="bookingStore.currentStep === 1"
          @click="bookingStore.currentStep = 2"
          :disabled="bookingStore.selectedSlots.length === 0"
          :class="[
            'w-full py-3 rounded-2xl font-bold text-xs shadow-xs transition flex items-center justify-center space-x-2',
            bookingStore.selectedSlots.length > 0
              ? 'bg-emerald-600 hover:bg-emerald-500 text-white cursor-pointer active:scale-98'
              : 'bg-slate-100 text-slate-400 border border-slate-200 cursor-not-allowed'
          ]"
        >
          <span>Continue to Add-ons</span>
          <ArrowRight class="w-4 h-4" />
        </button>

        <!-- Step 2: MiniStore -> submit booking to Step 3 -->
        <button
          v-else-if="bookingStore.currentStep === 2"
          @click="onProceedToPayment"
          :disabled="bookingStore.isSubmitting"
          class="w-full py-3 bg-emerald-600 hover:bg-emerald-500 text-white rounded-2xl font-bold text-xs shadow-xs transition flex items-center justify-center space-x-2 cursor-pointer active:scale-98 disabled:opacity-50"
        >
          <span v-if="bookingStore.isSubmitting" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
          <span v-else>Proceed to Payment (฿{{ bookingStore.netTotal.toLocaleString() }})</span>
          <ArrowRight v-if="!bookingStore.isSubmitting" class="w-4 h-4" />
        </button>

        <!-- Step 3: Payment step -->
        <div
          v-else-if="bookingStore.currentStep === 3"
          class="p-2.5 rounded-xl bg-amber-50 border border-amber-200 text-amber-900 text-center text-xs font-semibold flex items-center justify-center space-x-1.5"
        >
          <Clock class="w-3.5 h-3.5 text-amber-600 shrink-0" />
          <span>Upload payment slip on the left</span>
        </div>

        <!-- Step 4: E-Ticket step -->
        <div
          v-else-if="bookingStore.currentStep === 4"
          class="p-2.5 rounded-xl bg-emerald-50 border border-emerald-200 text-emerald-900 text-center text-xs font-bold flex items-center justify-center space-x-1.5"
        >
          <CheckCircle2 class="w-3.5 h-3.5 text-emerald-600 shrink-0" />
          <span>Booking Confirmed & Ready</span>
        </div>
      </div>

      <!-- Trust Badges -->
      <div class="pt-1 flex items-center justify-center space-x-3 text-[10px] text-slate-400 font-medium">
        <span class="flex items-center space-x-1">
          <ShieldCheck class="w-3 h-3 text-emerald-600" />
          <span>Atomic Slot Lock</span>
        </span>
        <span>•</span>
        <span>Instant QR Verification</span>
      </div>
    </div>
  </div>
</template>
