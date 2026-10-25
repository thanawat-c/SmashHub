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
</script>

<template>
  <div class="min-h-[calc(100vh-56px)] bg-slate-100 flex items-center justify-center sm:py-6 sm:px-4 select-none">
    <!-- Smartphone Frame (Centered on desktop, full screen on mobile) -->
    <div class="w-full sm:max-w-[416px] h-[calc(100vh-56px)] sm:h-[860px] sm:max-h-[92vh] bg-white sm:rounded-[44px] sm:border-[8px] sm:border-slate-900 sm:shadow-2xl flex flex-col justify-between overflow-hidden relative">
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
            <div class="flex items-center space-x-1.5">
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

      <!-- iOS Bottom Home Indicator Bar (Clean Light Mode) -->
      <div class="h-4 bg-white flex items-center justify-center flex-shrink-0">
        <div class="w-28 h-1 bg-slate-300 rounded-full"></div>
      </div>
    </div>
  </div>
</template>
