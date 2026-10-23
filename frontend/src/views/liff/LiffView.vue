<script setup>
import { useBookingStore } from '@/stores/bookingStore'
import LiffHeader from './components/LiffHeader.vue'
import CourtMatrixStep from './components/CourtMatrixStep.vue'
import MiniStoreStep from './components/MiniStoreStep.vue'
import PaymentStep from './components/PaymentStep.vue'
import ETicketStep from './components/ETicketStep.vue'
import { Calendar, ShoppingBag, CreditCard, Ticket } from 'lucide-vue-next'

const bookingStore = useBookingStore()

const steps = [
  { num: 1, label: 'เลือกคอร์ท', icon: Calendar },
  { num: 2, label: 'สั่งสินค้า', icon: ShoppingBag },
  { num: 3, label: 'ชำระเงิน', icon: CreditCard },
  { num: 4, label: 'ตั๋วเข้าสนาม', icon: Ticket }
]
</script>

<template>
  <div class="flex-1 flex justify-center items-center bg-slate-900/95 p-0 sm:p-4">
    <!-- Mobile Frame Simulation (max-w-[430px]) -->
    <div class="w-full max-w-[430px] h-[100dvh] sm:h-[880px] bg-slate-50 flex flex-col shadow-2xl sm:rounded-[36px] overflow-hidden border-0 sm:border-[8px] sm:border-slate-800 relative">
      <!-- LIFF Header with Mock Profile Switcher -->
      <LiffHeader />

      <!-- Step Progress Bar -->
      <div class="bg-white px-3 py-2 border-b border-slate-200">
        <div class="flex items-center justify-between">
          <div
            v-for="(step, idx) in steps"
            :key="step.num"
            class="flex items-center"
          >
            <!-- Step Item -->
            <div class="flex items-center space-x-1">
              <div
                :class="[
                  'w-5 h-5 rounded-full flex items-center justify-center text-[10px] font-bold transition',
                  bookingStore.currentStep === step.num
                    ? 'bg-emerald-600 text-white shadow-xs'
                    : bookingStore.currentStep > step.num
                    ? 'bg-emerald-100 text-emerald-700'
                    : 'bg-slate-100 text-slate-400'
                ]"
              >
                {{ bookingStore.currentStep > step.num ? '✓' : step.num }}
              </div>
              <span
                :class="[
                  'text-[10px] font-medium hidden xs:inline',
                  bookingStore.currentStep === step.num ? 'text-emerald-700 font-bold' : 'text-slate-400'
                ]"
              >
                {{ step.label }}
              </span>
            </div>

            <!-- Divider Line between steps -->
            <div
              v-if="idx < steps.length - 1"
              :class="[
                'w-4 sm:w-6 h-0.5 mx-1 transition',
                bookingStore.currentStep > step.num ? 'bg-emerald-500' : 'bg-slate-200'
              ]"
            ></div>
          </div>
        </div>
      </div>

      <!-- Main Step View Area -->
      <div class="flex-1 overflow-hidden relative">
        <CourtMatrixStep v-if="bookingStore.currentStep === 1" />
        <MiniStoreStep v-else-if="bookingStore.currentStep === 2" />
        <PaymentStep v-else-if="bookingStore.currentStep === 3" />
        <ETicketStep v-else-if="bookingStore.currentStep === 4" />
      </div>
    </div>
  </div>
</template>
