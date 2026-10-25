<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'
import AuthModal from '@/components/AuthModal.vue'
import {
  LogIn, LogOut, User, ShieldCheck, ChevronDown, Check
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const showUserMenu = ref(false)

function navigateToAdmin() {
  if (!authStore.isAuthenticated) {
    authStore.openAuthModal('login')
    return
  }
  router.push('/admin')
}

function handleLogout() {
  showUserMenu.value = false
  authStore.logout()
  if (route.path === '/admin') {
    router.push('/liff')
  }
}
</script>

<template>
  <div class="min-h-screen flex flex-col bg-slate-100">
    <!-- Clean White Navigation Bar -->
    <header class="bg-white text-slate-800 border-b border-slate-200/80 px-4 py-2 flex items-center justify-between text-xs sticky top-0 z-50 shadow-2xs">
      <div class="flex items-center space-x-2.5">
        <div class="w-8 h-8 rounded-xl bg-emerald-600 text-white flex items-center justify-center text-sm shadow-xs font-bold">
          🏸
        </div>
        <div>
          <span class="text-sm font-black text-slate-900 tracking-tight">SmashHub</span>
          <span class="text-slate-400 text-[11px] ml-1.5 hidden lg:inline">ระบบจองสนามแบดมินตัน & มินิสโตร์ POS</span>
        </div>
      </div>

      <nav class="flex items-center space-x-2 sm:space-x-3">
        <!-- Main Navigation Segment -->
        <div class="bg-slate-100 p-1 rounded-xl flex items-center space-x-1 border border-slate-200/60">
          <button
            @click="router.push('/liff')"
            :class="[
              'px-3 py-1.5 rounded-lg font-bold transition flex items-center space-x-1.5 text-xs select-none',
              route.path === '/liff'
                ? 'bg-white text-slate-900 shadow-xs'
                : 'text-slate-500 hover:text-slate-800'
            ]"
          >
            <span>📱</span>
            <span>ฝั่งลูกค้า (LIFF)</span>
          </button>

          <button
            @click="navigateToAdmin"
            :class="[
              'px-3 py-1.5 rounded-lg font-bold transition flex items-center space-x-1.5 text-xs select-none',
              route.path === '/admin'
                ? 'bg-white text-slate-900 shadow-xs'
                : 'text-slate-500 hover:text-slate-800'
            ]"
          >
            <span>🏢</span>
            <span>ฝั่งร้านค้า (POS)</span>
          </button>
        </div>

        <button
          @click="router.push('/demo')"
          :class="[
            'px-2.5 py-1.5 rounded-lg text-xs font-semibold transition text-slate-500 hover:text-slate-900 hover:bg-slate-100 hidden sm:inline-flex items-center space-x-1',
            route.path === '/demo' ? 'bg-slate-200 text-slate-900 font-bold' : ''
          ]"
          title="ดูแบบแบ่ง 2 จอเทียบกัน"
        >
          <span>Split View</span>
        </button>

        <div class="h-4 w-px bg-slate-200 hidden sm:block"></div>

        <!-- AUTH SECTION -->
        <!-- 1. If NOT Logged in -->
        <div v-if="!authStore.isAuthenticated">
          <button
            @click="authStore.openAuthModal('login')"
            class="px-3.5 py-1.5 rounded-xl bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-xs transition flex items-center space-x-1.5 shadow-2xs active:scale-95"
          >
            <LogIn class="w-3.5 h-3.5" />
            <span>เข้าสู่ระบบ / สมัครสมาชิก</span>
          </button>
        </div>

        <!-- 2. If Logged in -->
        <div v-else class="relative">
          <button
            @click="showUserMenu = !showUserMenu"
            class="flex items-center space-x-2 pl-1.5 pr-2.5 py-1 rounded-xl bg-slate-50 hover:bg-slate-100 border border-slate-200 transition shadow-2xs"
          >
            <img
              :src="authStore.avatar"
              :alt="authStore.displayName"
              class="w-6 h-6 rounded-full bg-slate-200 object-cover border border-slate-300"
            />
            <div class="text-left hidden xs:block">
              <div class="font-extrabold text-slate-900 text-[11px] leading-tight truncate max-w-[100px]">
                {{ authStore.displayName }}
              </div>
              <div class="text-[9px] font-mono leading-none mt-0.5">
                <span v-if="authStore.isAdmin" class="text-emerald-700 font-bold">Admin POS</span>
                <span v-else class="text-slate-500">Customer</span>
              </div>
            </div>
            <ChevronDown class="w-3 h-3 text-slate-400" />
          </button>

          <!-- User Menu Dropdown -->
          <div
            v-if="showUserMenu"
            class="absolute right-0 mt-1.5 w-56 bg-white rounded-2xl shadow-xl border border-slate-200 py-1.5 text-xs z-50 animate-in fade-in"
          >
            <div class="px-3.5 py-2 border-b border-slate-100">
              <div class="font-bold text-slate-900 text-xs">{{ authStore.displayName }}</div>
              <div class="text-[10px] text-slate-400 font-mono mt-0.5">
                {{ authStore.user?.phone || authStore.user?.email || 'LINE Account' }}
              </div>
              <span
                :class="[
                  'inline-block mt-1 text-[9px] font-mono font-bold px-2 py-0.2 rounded-full',
                  authStore.isAdmin ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-blue-50 text-blue-700 border border-blue-200'
                ]"
              >
                {{ authStore.isAdmin ? '👑 Store Manager (Admin)' : '🏸 Member (Customer)' }}
              </span>
            </div>

            <div class="py-1">
              <button
                @click="showUserMenu = false; router.push('/liff')"
                class="w-full px-3.5 py-2 text-left text-slate-700 hover:bg-slate-50 transition flex items-center space-x-2"
              >
                <span>📱</span>
                <span>จองสนามแบดมินตัน (LINE LIFF)</span>
              </button>

              <button
                v-if="authStore.isAdmin"
                @click="showUserMenu = false; router.push('/admin')"
                class="w-full px-3.5 py-2 text-left text-slate-700 hover:bg-slate-50 transition flex items-center space-x-2"
              >
                <span>🏢</span>
                <span>ระบบบริหารหลังร้าน (Admin POS)</span>
              </button>

              <button
                @click="handleLogout"
                class="w-full px-3.5 py-2 text-left text-rose-600 hover:bg-rose-50 transition flex items-center space-x-2 font-semibold border-t border-slate-100 mt-1 pt-2"
              >
                <LogOut class="w-3.5 h-3.5" />
                <span>ออกจากระบบ</span>
              </button>
            </div>
          </div>
        </div>
      </nav>
    </header>

    <!-- Main View Content -->
    <main class="flex-1 flex flex-col">
      <router-view />
    </main>

    <!-- Global Auth Modal -->
    <AuthModal />
  </div>
</template>
