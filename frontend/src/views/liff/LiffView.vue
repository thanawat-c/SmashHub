<script setup>
import { useBookingStore } from '@/stores/bookingStore'
import LiffHeader from './components/LiffHeader.vue'
import CourtMatrixStep from './components/CourtMatrixStep.vue'
import MiniStoreStep from './components/MiniStoreStep.vue'
import PaymentStep from './components/PaymentStep.vue'
import ETicketStep from './components/ETicketStep.vue'
import DesktopOrderSummary from './components/DesktopOrderSummary.vue'
import { Check, ChevronLeft } from 'lucide-vue-next'

const bookingStore = useBookingStore()

const steps = [
  { num: 1, label: 'Court Schedule', short: 'Courts' },
  { num: 2, label: 'Add-ons & Pro Shop', short: 'Add-ons' },
  { num: 3, label: 'Slip Verification', short: 'Payment' },
  { num: 4, label: 'Digital Court Pass', short: 'Pass' }
]

function goToStep(stepNum) {
  if (stepNum < bookingStore.currentStep) {
    bookingStore.currentStep = stepNum
  }
}

function handleBack() {
  if (bookingStore.currentStep > 1) {
    if (bookingStore.currentStep === 4) {
      bookingStore.resetSelection()
      bookingStore.currentStep = 1
      bookingStore.fetchAvailability()
    } else {
      bookingStore.currentStep--
    }
  } else {
    bookingStore.resetSelection()
    bookingStore.fetchAvailability()
  }
}
</script>

<template>
  <div class="min-h-screen lg:min-h-[calc(100vh-56px)] bg-slate-100/90 flex items-center justify-center p-0 lg:p-6 select-none">
    <div class="w-full sm:max-w-[480px] lg:max-w-6xl xl:max-w-7xl h-screen lg:h-[860px] lg:max-h-[92vh] bg-white sm:border-x lg:border border-slate-200/90 sm:shadow-lg lg:rounded-3xl flex flex-col justify-between overflow-hidden relative transition-all">
      <LiffHeader class="lg:hidden" />

      <!-- Step Progress Indicator -->
      <div class="bg-white px-4 py-2.5 lg:py-3 border-b border-slate-100 flex-shrink-0">
        <div class="flex items-center justify-between max-w-2xl mx-auto">
          <!-- Desktop Back Button -->
          <button
            v-if="bookingStore.currentStep > 1"
            @click="handleBack"
            class="hidden lg:flex items-center space-x-1 text-slate-500 hover:text-slate-800 transition mr-3 text-xs font-semibold py-1 px-2 rounded-lg hover:bg-slate-100 cursor-pointer -ml-2"
            title="Go to previous step"
          >
            <ChevronLeft class="w-4 h-4" />
            <span>Back</span>
          </button>
          <div
            v-for="(step, idx) in steps"
            :key="step.num"
            class="flex items-center flex-1"
          >
            <div
              @click="goToStep(step.num)"
              :class="[
                'flex items-center space-x-2 select-none transition',
                step.num < bookingStore.currentStep ? 'cursor-pointer hover:opacity-80 active:scale-95' : ''
              ]"
              :title="step.num < bookingStore.currentStep ? `Return to: ${step.label}` : ''"
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

              <!-- Responsive Step Label -->
              <div class="flex flex-col">
                <span
                  :class="[
                    'text-[10px] tracking-tight font-medium',
                    bookingStore.currentStep === step.num
                      ? 'text-slate-900 font-bold'
                      : bookingStore.currentStep > step.num
                      ? 'text-emerald-700'
                      : 'text-slate-400'
                  ]"
                >
                  <span class="inline lg:hidden">{{ step.short }}</span>
                  <span class="hidden lg:inline">{{ step.label }}</span>
                </span>
              </div>
            </div>

            <!-- Connecting divider -->
            <div
              v-if="idx < steps.length - 1"
              :class="[
                'flex-1 h-0.5 mx-2 rounded-full transition-all',
                bookingStore.currentStep > step.num ? 'bg-emerald-500' : 'bg-slate-100'
              ]"
            ></div>
          </div>
        </div>
      </div>

      <div class="flex-1 flex overflow-hidden relative bg-[#F8FAFC]">
        <!-- Step Content Canvas -->
        <div class="flex-1 h-full overflow-hidden flex flex-col relative">
          <CourtMatrixStep v-if="bookingStore.currentStep === 1" />
          <MiniStoreStep v-else-if="bookingStore.currentStep === 2" />
          <PaymentStep v-else-if="bookingStore.currentStep === 3" />
          <ETicketStep v-else-if="bookingStore.currentStep === 4" />
        </div>

        <!-- Sticky Order Summary (Desktop) -->
        <div class="hidden lg:flex lg:w-[380px] xl:w-[420px] h-full flex-col bg-white border-l border-slate-200/90 overflow-hidden shadow-xs flex-shrink-0">
          <DesktopOrderSummary />
        </div>
      </div>
    </div>
  </div>
</template>
