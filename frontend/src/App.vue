<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const isDemo = computed(() => route.path === '/demo')
</script>

<template>
  <div class="min-h-screen flex flex-col bg-slate-100">
    <!-- Clean Top Navigation Bar -->
    <header class="bg-slate-900 text-white border-b border-slate-800 px-4 py-2.5 flex items-center justify-between text-xs sticky top-0 z-50 shadow-md">
      <div class="flex items-center space-x-2.5">
        <div class="w-7 h-7 rounded-lg bg-emerald-600 flex items-center justify-center text-sm shadow-xs font-bold">
          🏸
        </div>
        <div>
          <span class="text-sm font-extrabold text-white tracking-tight">SmashHub</span>
          <span class="text-slate-400 text-[11px] ml-1.5 hidden md:inline">| ระบบจองสนามแบดมินตัน & มินิสโตร์ POS</span>
        </div>
      </div>

      <nav class="flex items-center space-x-1.5 sm:space-x-2">
        <button
          @click="router.push('/liff')"
          :class="[
            'px-3 py-1.5 rounded-xl font-bold transition flex items-center space-x-1.5 text-xs select-none',
            route.path === '/liff'
              ? 'bg-emerald-600 text-white shadow-sm ring-1 ring-emerald-400/40'
              : 'text-slate-300 hover:text-white hover:bg-slate-800'
          ]"
        >
          <span>📱</span>
          <span>ฝั่งลูกค้า (LINE LIFF)</span>
        </button>

        <button
          @click="router.push('/admin')"
          :class="[
            'px-3 py-1.5 rounded-xl font-bold transition flex items-center space-x-1.5 text-xs select-none',
            route.path === '/admin'
              ? 'bg-emerald-600 text-white shadow-sm ring-1 ring-emerald-400/40'
              : 'text-slate-300 hover:text-white hover:bg-slate-800'
          ]"
        >
          <span>🏢</span>
          <span>ฝั่งร้านค้า (Admin POS)</span>
        </button>

        <div class="h-4 w-px bg-slate-800 mx-1 hidden sm:block"></div>

        <button
          @click="router.push('/demo')"
          :class="[
            'px-2 py-1 rounded-lg text-[11px] font-medium transition text-slate-400 hover:text-slate-200 hover:bg-slate-800/80',
            route.path === '/demo' ? 'bg-slate-800 text-white' : ''
          ]"
          title="ดูแบบแบ่ง 2 จอเทียบกัน"
        >
          <span>Split Test</span>
        </button>
      </nav>
    </header>

    <!-- Main View Content -->
    <main class="flex-1 flex flex-col">
      <router-view />
    </main>
  </div>
</template>
