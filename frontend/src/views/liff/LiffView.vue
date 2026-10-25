<script setup>
import { useBookingStore } from '@/stores/bookingStore'
import LiffHeader from './components/LiffHeader.vue'
import CourtMatrixStep from './components/CourtMatrixStep.vue'
import MiniStoreStep from './components/MiniStoreStep.vue'
import PaymentStep from './components/PaymentStep.vue'
import ETicketStep from './components/ETicketStep.vue'
import { Calendar, ShoppingBag, CreditCard, Ticket, Check } from 'lucide-vue-next'

const bookingStore = useBookingStore()

const steps = [
  { num: 1, label: 'เลือกสนาม' },
  { num: 2, label: 'อุปกรณ์' },
  { num: 3, label: 'ชำระเงิน' },
  { num: 4, label: 'ตั๋วเข้าเล่น' }
]
function goToStep(stepNum) {
  if (stepNum < bookingStore.currentStep) {
    bookingStore.currentStep = stepNum
  }
}
</script>

<template>
  <div class="min-h-[calc(100vh-56px)] bg-slate-100/90 flex items-center justify-center select-none">
    <!-- Clean Authentic Mobile-First Container (Full width on mobile, sleek centered max-w-[430px] on desktop) -->
    <div class="w-full sm:max-w-[430px] h-[calc(100vh-56px)] bg-white sm:border-x sm:border-slate-200/90 sm:shadow-lg flex flex-col justify-between overflow-hidden relative">
      <!-- Top LIFF Navigation & Profile Bar -->
      <LiffHeader />

      <!-- Minimalist Clean Step Progress Indicator -->
      <div class="bg-white px-4 py-2 border-b border-slate-100 flex-shrink-0">
        <div class="flex items-center justify-between">
          <div
            v-for="(step, idx) in steps"
            :key="step.num"
            class="flex items-center flex-1"
          >
            <div
              @click="goToStep(step.num)"
              :class="[
                'flex items-center space-x-1.5 select-none transition',
                step.num < bookingStore.currentStep ? 'cursor-pointer hover:opacity-80 active:scale-95' : ''
              ]"
              :title="step.num < bookingStore.currentStep ? `ย้อนกลับไป: ${step.label}` : ''"
            >
              <div
                :class="[
                  'w-5 h-5 rounded-full flex items-center justify-center text-[10px] font-bold font-mono transition-all',
                  bookingStore.currentStep === step.num
                    ? 'bg-emerald-600 text-white shadow-xs'
                    : bookingStore.currentStep > step.num
                    ? 'bg-emerald-100 text-emerald-800'
                    : 'bg-slate-100 text-slate-400'
                ]"
              >
                <Check v-if="bookingStore.currentStep > step.num" class="w-3 h-3 stroke-[3]" />
                <span v-else>{{ step.num }}</span>
              </div>
              <span
                :class="[
                  'text-[10px] tracking-tight font-medium hidden xs:inline',
                  bookingStore.currentStep === step.num
                    ? 'text-slate-900 font-bold'
                    : bookingStore.currentStep > step.num
                    ? 'text-emerald-700'
                    : 'text-slate-400'
                ]"
              >
                {{ step.label }}
              </span>
            </div>

            <!-- Connecting divider -->
            <div
              v-if="idx < steps.length - 1"
              :class="[
                'flex-1 h-0.5 mx-1.5 rounded-full transition-all',
                bookingStore.currentStep > step.num ? 'bg-emerald-500' : 'bg-slate-100'
              ]"
            ></div>
          </div>
        </div>
      </div>

      <!-- Active Step Content Area -->
      <div class="flex-1 overflow-hidden relative bg-[#F8FAFC]">
        <CourtMatrixStep v-if="bookingStore.currentStep === 1" />
        <MiniStoreStep v-else-if="bookingStore.currentStep === 2" />
        <PaymentStep v-else-if="bookingStore.currentStep === 3" />
        <ETicketStep v-else-if="bookingStore.currentStep === 4" />
      </div>
    </div>
  </div>
</template>
