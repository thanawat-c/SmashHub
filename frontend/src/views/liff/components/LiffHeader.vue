<script setup>
import { ref, computed } from 'vue'
import { useLineProfileStore } from '@/stores/lineProfile'
import { useBookingStore } from '@/stores/bookingStore'
import { useAuthStore } from '@/stores/authStore'
import {
  ChevronDown, Check, ShieldCheck, ChevronLeft, User, RotateCcw, LogIn, LogOut
} from 'lucide-vue-next'

const profileStore = useLineProfileStore()
const bookingStore = useBookingStore()
const authStore = useAuthStore()
const showDropdown = ref(false)

const activeDisplayName = computed(() => {
  return authStore.user?.display_name || profileStore.currentProfile.displayName
})

const activeAvatar = computed(() => {
  return authStore.avatar || profileStore.currentProfile.pictureUrl
})

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
    <!-- Clean Authentic LIFF Bar Header -->
    <div class="px-4 py-3 flex items-center justify-between bg-white">
      <div class="flex items-center space-x-2">
        <button
          v-if="bookingStore.currentStep > 1"
          @click="handleBack"
          class="text-slate-700 hover:text-slate-900 transition p-1.5 rounded-xl hover:bg-slate-100 active:bg-slate-200 cursor-pointer flex items-center justify-center -ml-1"
          title="ย้อนกลับขั้นตอนก่อนหน้า"
        >
          <ChevronLeft class="w-4 h-4 stroke-[2.5]" />
        </button>
        <div class="flex items-center space-x-2">
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

      <!-- Quick Customer Profile & Auth Pill -->
      <div class="relative">
        <button
          @click="showDropdown = !showDropdown"
          class="flex items-center space-x-1.5 bg-slate-100 hover:bg-slate-200 text-slate-800 border border-slate-200 text-[10px] pl-1 pr-2 py-0.5 rounded-full transition shadow-2xs"
          title="โปรไฟล์ผู้ใช้งาน"
        >
          <img
            :src="activeAvatar"
            class="w-4 h-4 rounded-full bg-slate-200 object-cover"
          />
          <span class="max-w-[75px] truncate font-semibold text-slate-800">{{ activeDisplayName }}</span>
          <ChevronDown class="w-3 h-3 text-slate-500" />
        </button>

        <!-- Dropdown Popover -->
        <div
          v-if="showDropdown"
          class="absolute right-0 mt-1.5 w-64 bg-white rounded-2xl shadow-xl border border-slate-200 py-1.5 text-xs z-50 animate-in fade-in zoom-in-95"
        >
          <div class="px-3.5 py-2 border-b border-slate-100 flex items-center justify-between">
            <span class="text-[10px] uppercase tracking-wider font-bold text-slate-400">บัญชีผู้ใช้งาน</span>
            <span v-if="authStore.isAuthenticated" class="text-[9px] font-mono font-bold px-1.5 py-0.2 rounded bg-emerald-50 text-emerald-700 border border-emerald-200">
              {{ authStore.isAdmin ? 'ADMIN' : 'MEMBER' }}
            </span>
          </div>

          <!-- Auth Actions inside Dropdown -->
          <div class="p-2 border-b border-slate-100 space-y-1">
            <button
              v-if="!authStore.isAuthenticated"
              @click="showDropdown = false; authStore.openAuthModal('login')"
              class="w-full py-1.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-xs font-bold transition flex items-center justify-center space-x-1.5 shadow-2xs"
            >
              <LogIn class="w-3.5 h-3.5" />
              <span>เข้าสู่ระบบ / สมัครสมาชิก</span>
            </button>

            <button
              v-else
              @click="showDropdown = false; authStore.logout()"
              class="w-full py-1.5 bg-rose-50 hover:bg-rose-100 text-rose-700 rounded-xl text-xs font-bold transition flex items-center justify-center space-x-1.5 border border-rose-200"
            >
              <LogOut class="w-3.5 h-3.5" />
              <span>ออกจากระบบ</span>
            </button>
          </div>

          <div class="px-3 py-1.5 text-[10px] uppercase tracking-wider font-bold text-slate-400">
            สลับโปรไฟล์จำลอง (Test Sandbox)
          </div>
          <button
            v-for="profile in profileStore.mockProfiles"
            :key="profile.id"
            @click="switchProfile(profile)"
            :class="[
              'w-full px-3 py-2 text-left flex items-center space-x-2.5 transition',
              activeDisplayName === profile.displayName
                ? 'bg-emerald-50 text-emerald-700 font-semibold'
                : 'text-slate-700 hover:bg-slate-50'
            ]"
          >
            <img :src="profile.pictureUrl" class="w-7 h-7 rounded-full bg-slate-100 border border-slate-200" />
            <div class="flex-1 truncate">
              <div class="truncate text-xs text-slate-900 font-bold">{{ profile.displayName }}</div>
              <div class="text-[10px] text-slate-500 font-mono">{{ profile.phone }}</div>
            </div>
            <Check v-if="activeDisplayName === profile.displayName" class="w-4 h-4 text-emerald-600" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
