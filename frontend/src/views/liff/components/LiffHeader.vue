<script setup>
import { ref } from 'vue'
import { useLineProfileStore } from '@/stores/lineProfile'
import { ArrowLeftRight, User } from 'lucide-vue-next'

const profileStore = useLineProfileStore()
const showDropdown = ref(false)

function switchProfile(profile) {
  profileStore.selectProfile(profile)
  showDropdown.value = false
}
</script>

<template>
  <div class="bg-emerald-700 text-white px-4 py-3 shadow relative z-20">
    <div class="flex items-center justify-between">
      <!-- Current User Profile Info -->
      <div class="flex items-center space-x-3">
        <div class="relative">
          <img
            :src="profileStore.currentProfile.pictureUrl"
            :alt="profileStore.currentProfile.displayName"
            class="w-10 h-10 rounded-full border-2 border-emerald-300 bg-emerald-800"
          />
          <div class="absolute -bottom-1 -right-1 bg-green-400 w-3.5 h-3.5 rounded-full border-2 border-emerald-700"></div>
        </div>
        <div>
          <div class="flex items-center space-x-1.5">
            <span class="font-semibold text-sm leading-tight">{{ profileStore.currentProfile.displayName }}</span>
            <span class="text-[10px] bg-emerald-800 text-emerald-200 px-1.5 py-0.5 rounded font-mono">LINE</span>
          </div>
          <p class="text-xs text-emerald-200 font-mono">{{ profileStore.currentProfile.phone }}</p>
        </div>
      </div>

      <!-- Mock Profile Switcher Button -->
      <div class="relative">
        <button
          @click="showDropdown = !showDropdown"
          class="flex items-center space-x-1 bg-emerald-800/80 hover:bg-emerald-800 text-emerald-100 text-xs px-2.5 py-1.5 rounded-full border border-emerald-600 transition"
          title="สลับบัญชีผู้ใช้จำลองสำหรับทดสอบระบบ"
        >
          <ArrowLeftRight class="w-3.5 h-3.5" />
          <span class="text-[11px]">สลับโปรไฟล์</span>
        </button>

        <!-- Dropdown Menu -->
        <div
          v-if="showDropdown"
          class="absolute right-0 mt-2 w-64 bg-white rounded-xl shadow-xl border border-slate-200 py-1.5 text-slate-800 text-xs z-50 animate-in fade-in slide-in-from-top-2"
        >
          <div class="px-3 py-1.5 border-b border-slate-100 font-semibold text-slate-500 text-[11px] flex items-center justify-between">
            <span>เลือกโปรไฟล์ทดสอบ (Portfolio Mock)</span>
          </div>
          <button
            v-for="profile in profileStore.mockProfiles"
            :key="profile.id"
            @click="switchProfile(profile)"
            :class="[
              'w-full px-3 py-2 text-left flex items-center space-x-2.5 transition',
              profileStore.currentProfile.id === profile.id ? 'bg-emerald-50 text-emerald-700 font-medium' : 'hover:bg-slate-50'
            ]"
          >
            <img :src="profile.pictureUrl" class="w-7 h-7 rounded-full bg-slate-100 border border-slate-200" />
            <div class="flex-1 truncate">
              <div class="truncate text-xs">{{ profile.displayName }}</div>
              <div class="text-[10px] text-slate-400 font-mono">{{ profile.phone }}</div>
            </div>
            <span v-if="profileStore.currentProfile.id === profile.id" class="text-emerald-600 font-bold text-xs">✓</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
