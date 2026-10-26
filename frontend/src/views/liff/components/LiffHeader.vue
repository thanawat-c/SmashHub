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

        <!-- Dropdown Popover (Clean Authentic App Menu) -->
        <div
          v-if="showDropdown"
          class="absolute right-0 mt-2 w-64 bg-white rounded-2xl shadow-2xl border border-slate-200/90 py-2 text-xs z-50 animate-in fade-in zoom-in-95"
        >
          <!-- User info header -->
          <div class="px-3.5 py-2 border-b border-slate-100 flex items-center space-x-2.5">
            <img :src="activeAvatar" class="w-9 h-9 rounded-full bg-slate-100 border border-slate-200 object-cover" />
            <div class="flex-1 min-w-0">
              <div class="font-extrabold text-slate-900 truncate">{{ activeDisplayName }}</div>
              <div class="text-[10px] text-slate-400 font-mono truncate">
                {{ authStore.user?.email || authStore.user?.phone || profileStore.currentProfile.phone }}
              </div>
            </div>
            <span v-if="authStore.isAuthenticated" class="text-[9px] font-mono font-bold px-1.5 py-0.5 rounded bg-emerald-50 text-emerald-700 border border-emerald-200">
              {{ authStore.isAdmin ? 'ADMIN' : 'MEMBER' }}
            </span>
          </div>

          <!-- Customer Action Links -->
          <div class="p-1.5 border-b border-slate-100 space-y-0.5">
            <button
              @click="showDropdown = false; bookingStore.openMyBookings(authStore.user || profileStore.currentProfile)"
              class="w-full px-3 py-2 text-left hover:bg-slate-50 text-slate-700 rounded-xl transition flex items-center space-x-2.5 font-medium cursor-pointer"
            >
              <span class="text-sm">📋</span>
              <span class="text-xs font-semibold text-slate-800">ประวัติการจอง & ตั๋วของฉัน</span>
            </button>

            <button
              v-if="!authStore.isAuthenticated"
              @click="showDropdown = false; authStore.openAuthModal('login')"
              class="w-full px-3 py-2 text-left hover:bg-emerald-50 text-emerald-700 rounded-xl transition flex items-center space-x-2.5 font-semibold cursor-pointer"
            >
              <LogIn class="w-4 h-4 text-emerald-600" />
              <span class="text-xs">เข้าสู่ระบบ / สมัครสมาชิก</span>
            </button>

            <button
              v-else
              @click="showDropdown = false; authStore.logout()"
              class="w-full px-3 py-2 text-left hover:bg-rose-50 text-rose-700 rounded-xl transition flex items-center space-x-2.5 font-semibold cursor-pointer"
            >
              <LogOut class="w-4 h-4 text-rose-600" />
              <span class="text-xs">ออกจากระบบ</span>
            </button>
          </div>

          <!-- Subtle Dev Sandbox Swapper -->
          <div class="px-3 pt-2 pb-1 text-[10px] uppercase tracking-wider font-bold text-slate-400">
            จำลองโปรไฟล์ทดสอบ (Sandbox)
          </div>
          <div class="px-1.5 space-y-0.5">
            <button
              v-for="profile in profileStore.mockProfiles"
              :key="profile.id"
              @click="switchProfile(profile)"
              :class="[
                'w-full px-2.5 py-1.5 rounded-xl text-left flex items-center space-x-2 transition cursor-pointer',
                activeDisplayName === profile.displayName
                  ? 'bg-emerald-50 text-emerald-800 font-bold'
                  : 'text-slate-600 hover:bg-slate-50'
              ]"
            >
              <img :src="profile.pictureUrl" class="w-5 h-5 rounded-full bg-slate-100 border border-slate-200" />
              <span class="flex-1 truncate text-[11px]">{{ profile.displayName }}</span>
              <Check v-if="activeDisplayName === profile.displayName" class="w-3.5 h-3.5 text-emerald-600" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
