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
  { num: 1, label: 'คอร์ท' },
  { num: 2, label: 'สินค้า' },
  { num: 3, label: 'ชำระเงิน' },
  { num: 4, label: 'ตั๋วเข้าสนาม' }
]
</script>

<template>
  <div class="w-full h-full bg-slate-900 flex flex-col justify-between overflow-hidden relative">
    <!-- Top LIFF Navigation & Profile Bar -->
    <LiffHeader />

    <!-- Step Progress Pills -->
    <div class="bg-slate-950 px-3 py-2 border-b border-slate-800 flex items-center justify-between">
      <div
        v-for="(step, idx) in steps"
        :key="step.num"
        class="flex items-center flex-1"
      >
        <div class="flex items-center space-x-1.5">
          <div
            :class="[
              'w-5 h-5 rounded-full flex items-center justify-center text-[10px] font-bold font-mono transition-all',
              bookingStore.currentStep === step.num
                ? 'bg-emerald-500 text-slate-950 shadow-glow-emerald ring-2 ring-emerald-400/40'
                : bookingStore.currentStep > step.num
                ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/40'
                : 'bg-slate-800 text-slate-500'
            ]"
          >
            {{ bookingStore.currentStep > step.num ? '✓' : step.num }}
          </div>
          <span
            :class="[
              'text-[10px] font-semibold tracking-tight',
              bookingStore.currentStep === step.num ? 'text-white font-bold' : 'text-slate-500'
            ]"
          >
            {{ step.label }}
          </span>
        </div>

        <div
          v-if="idx < steps.length - 1"
          :class="[
            'flex-1 h-0.5 mx-2 rounded-full transition-all',
            bookingStore.currentStep > step.num ? 'bg-emerald-500' : 'bg-slate-800'
          ]"
        ></div>
      </div>
    </div>

    <!-- Active Step Content Area -->
    <div class="flex-1 overflow-hidden relative">
      <CourtMatrixStep v-if="bookingStore.currentStep === 1" />
      <MiniStoreStep v-else-if="bookingStore.currentStep === 2" />
      <PaymentStep v-else-if="bookingStore.currentStep === 3" />
      <ETicketStep v-else-if="bookingStore.currentStep === 4" />
    </div>

    <!-- iOS Bottom Home Indicator Bar -->
    <div class="h-4 bg-slate-950 flex items-center justify-center">
      <div class="w-28 h-1 bg-slate-700 rounded-full"></div>
    </div>
  </div>
</template>
