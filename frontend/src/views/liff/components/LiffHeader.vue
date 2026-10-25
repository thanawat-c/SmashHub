<script setup>
import { ref } from 'vue'
import { useLineProfileStore } from '@/stores/lineProfile'
import { useBookingStore } from '@/stores/bookingStore'
import {
  ChevronDown, Check, ShieldCheck, ChevronLeft, Wifi, User, RotateCcw
} from 'lucide-vue-next'

const profileStore = useLineProfileStore()
const bookingStore = useBookingStore()
const showDropdown = ref(false)

function switchProfile(profile) {
  profileStore.selectProfile(profile)
  showDropdown.value = false
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
  <div class="bg-white text-slate-900 select-none relative z-30 border-b border-slate-100">
    <!-- iOS Simulated Status Bar (Clean Light Mode) -->
    <div class="h-6 px-5 pt-1 flex items-center justify-between text-[11px] font-mono tracking-tight text-slate-800">
      <span class="font-bold">9:41</span>
      <!-- Dynamic Island on phone frame -->
      <div class="w-16 h-3 bg-slate-900 rounded-full mx-auto hidden sm:block"></div>
      <div class="flex items-center space-x-1.5 text-xs text-slate-800">
        <span class="text-[9px] font-sans font-bold">5G</span>
        <Wifi class="w-3 h-3" />
        <div class="w-4 h-2 border border-slate-800 rounded-2xs p-0.5 flex items-center">
          <div class="w-full h-full bg-slate-900 rounded-2xs"></div>
        </div>
      </div>
    </div>

    <!-- LINE Official In-App Browser Header -->
    <div class="px-3 py-2 flex items-center justify-between bg-white">
      <div class="flex items-center space-x-2">
        <button
          @click="handleBack"
          class="text-slate-700 hover:text-slate-900 transition p-1.5 rounded-xl hover:bg-slate-100 active:bg-slate-200 cursor-pointer flex items-center justify-center"
          title="ย้อนกลับ / หน้าแรก"
        >
          <ChevronLeft class="w-4 h-4 stroke-[2.5]" />
        </button>
        <div class="flex items-center space-x-1.5">
          <div class="w-7 h-7 rounded-lg bg-emerald-600 text-white flex items-center justify-center text-xs shadow-xs font-bold">
            🏸
          </div>
          <div>
            <div class="flex items-center space-x-1">
              <span class="font-extrabold text-xs tracking-tight text-slate-900">SmashHub Arena</span>
              <ShieldCheck class="w-3.5 h-3.5 text-emerald-600" />
            </div>
            <p class="text-[9px] text-slate-500 leading-none">LINE LIFF Official</p>
          </div>
        </div>
      </div>

      <!-- Quick Mock Customer Profile Switcher Pill -->
      <div class="relative">
        <button
          @click="showDropdown = !showDropdown"
          class="flex items-center space-x-1.5 bg-slate-100 hover:bg-slate-200 text-slate-800 border border-slate-200 text-[10px] pl-1 pr-2 py-0.5 rounded-full transition shadow-2xs"
          title="สลับโปรไฟล์ลูกค้าจำลอง"
        >
          <img
            :src="profileStore.currentProfile.pictureUrl"
            class="w-4 h-4 rounded-full bg-slate-200 object-cover"
          />
          <span class="max-w-[70px] truncate font-semibold text-slate-800">{{ profileStore.currentProfile.displayName }}</span>
          <ChevronDown class="w-3 h-3 text-slate-500" />
        </button>

        <!-- Dropdown Popover -->
        <div
          v-if="showDropdown"
          class="absolute right-0 mt-1.5 w-60 bg-white rounded-2xl shadow-xl border border-slate-200 py-1.5 text-xs z-50 animate-in fade-in zoom-in-95"
        >
          <div class="px-3 py-1.5 border-b border-slate-100 text-[10px] uppercase tracking-wider font-bold text-slate-400">
            จำลองลูกค้า (LINE Profile)
          </div>
          <button
            v-for="profile in profileStore.mockProfiles"
            :key="profile.id"
            @click="switchProfile(profile)"
            :class="[
              'w-full px-3 py-2 text-left flex items-center space-x-2.5 transition',
              profileStore.currentProfile.id === profile.id
                ? 'bg-emerald-50 text-emerald-700 font-semibold'
                : 'text-slate-700 hover:bg-slate-50'
            ]"
          >
            <img :src="profile.pictureUrl" class="w-7 h-7 rounded-full bg-slate-100 border border-slate-200" />
            <div class="flex-1 truncate">
              <div class="truncate text-xs text-slate-900 font-bold">{{ profile.displayName }}</div>
              <div class="text-[10px] text-slate-500 font-mono">{{ profile.phone }}</div>
            </div>
            <Check v-if="profileStore.currentProfile.id === profile.id" class="w-4 h-4 text-emerald-600" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
