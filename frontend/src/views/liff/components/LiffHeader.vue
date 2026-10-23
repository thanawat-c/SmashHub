<script setup>
import { ref } from 'vue'
import { useLineProfileStore } from '@/stores/lineProfile'
import {
  ChevronDown, Check, ShieldCheck, MoreHorizontal, ChevronLeft, Wifi, BatteryMedium
} from 'lucide-vue-next'

const profileStore = useLineProfileStore()
const showDropdown = ref(false)

function switchProfile(profile) {
  profileStore.selectProfile(profile)
  showDropdown.value = false
}
</script>

<template>
  <div class="bg-slate-900 text-white select-none relative z-30">
    <!-- iOS Simulated Status Bar -->
    <div class="h-6 px-4 pt-1 flex items-center justify-between text-[11px] font-mono tracking-tight text-slate-300">
      <span class="font-bold">9:41</span>
      <!-- Dynamic Island pill on phone frame -->
      <div class="w-16 h-3 bg-black rounded-full mx-auto hidden sm:block border border-slate-800"></div>
      <div class="flex items-center space-x-1.5 text-xs text-slate-300">
        <span class="text-[9px] font-sans font-bold">5G</span>
        <Wifi class="w-3 h-3" />
        <div class="w-4 h-2 border border-slate-300 rounded-xs p-0.5 flex items-center">
          <div class="w-full h-full bg-emerald-400 rounded-2xs"></div>
        </div>
      </div>
    </div>

    <!-- LINE Official In-App Browser Header -->
    <div class="px-3 py-2 flex items-center justify-between border-b border-slate-800/80 bg-slate-900/95 backdrop-blur-md">
      <div class="flex items-center space-x-2">
        <button class="text-slate-400 hover:text-white transition">
          <ChevronLeft class="w-4 h-4" />
        </button>
        <div class="flex items-center space-x-1.5">
          <div class="w-6 h-6 rounded-lg bg-emerald-600 flex items-center justify-center text-xs shadow-xs">
            🏸
          </div>
          <div>
            <div class="flex items-center space-x-1">
              <span class="font-bold text-xs tracking-tight text-white">SmashHub Arena</span>
              <ShieldCheck class="w-3 h-3 text-emerald-400" />
            </div>
            <p class="text-[9px] text-slate-400 leading-none">LINE LIFF Official App</p>
          </div>
        </div>
      </div>

      <!-- Quick Mock Profile Switcher Pill -->
      <div class="relative">
        <button
          @click="showDropdown = !showDropdown"
          class="flex items-center space-x-1.5 bg-slate-800/90 hover:bg-slate-700/80 text-slate-200 border border-slate-700/80 text-[10px] pl-1.5 pr-2 py-1 rounded-full transition shadow-xs"
        >
          <img
            :src="profileStore.currentProfile.pictureUrl"
            class="w-4 h-4 rounded-full bg-slate-700 object-cover"
          />
          <span class="max-w-[75px] truncate font-medium">{{ profileStore.currentProfile.displayName }}</span>
          <ChevronDown class="w-3 h-3 text-slate-400" />
        </button>

        <!-- Dropdown Popover -->
        <div
          v-if="showDropdown"
          class="absolute right-0 mt-1.5 w-60 bg-slate-900 rounded-2xl shadow-2xl border border-slate-700/90 py-1.5 text-xs z-50 animate-in fade-in zoom-in-95"
        >
          <div class="px-3 py-1.5 border-b border-slate-800 text-[10px] uppercase tracking-wider font-bold text-slate-400">
            จำลองลูกค้า (Mock Switcher)
          </div>
          <button
            v-for="profile in profileStore.mockProfiles"
            :key="profile.id"
            @click="switchProfile(profile)"
            :class="[
              'w-full px-3 py-2 text-left flex items-center space-x-2.5 transition',
              profileStore.currentProfile.id === profile.id
                ? 'bg-emerald-500/15 text-emerald-400 font-semibold'
                : 'text-slate-300 hover:bg-slate-800/70'
            ]"
          >
            <img :src="profile.pictureUrl" class="w-7 h-7 rounded-full bg-slate-800 border border-slate-700" />
            <div class="flex-1 truncate">
              <div class="truncate text-xs">{{ profile.displayName }}</div>
              <div class="text-[10px] text-slate-400 font-mono">{{ profile.phone }}</div>
            </div>
            <Check v-if="profileStore.currentProfile.id === profile.id" class="w-4 h-4 text-emerald-400" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
