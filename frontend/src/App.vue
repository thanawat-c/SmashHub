<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const isDemo = computed(() => route.path === '/demo')
</script>

<template>
  <div class="min-h-screen flex flex-col bg-slate-100">
    <!-- Clean White Navigation Bar -->
    <header class="bg-white text-slate-800 border-b border-slate-200/80 px-4 py-2.5 flex items-center justify-between text-xs sticky top-0 z-50 shadow-2xs">
      <div class="flex items-center space-x-2.5">
        <div class="w-8 h-8 rounded-xl bg-emerald-600 text-white flex items-center justify-center text-sm shadow-xs font-bold">
          🏸
        </div>
        <div>
          <span class="text-sm font-black text-slate-900 tracking-tight">SmashHub</span>
          <span class="text-slate-400 text-[11px] ml-1.5 hidden md:inline">ระบบจองสนามแบดมินตัน & มินิสโตร์ POS</span>
        </div>
      </div>

      <nav class="flex items-center space-x-1.5 sm:space-x-2">
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
            <span>ฝั่งลูกค้า (LINE LIFF)</span>
          </button>

          <button
            @click="router.push('/admin')"
            :class="[
              'px-3 py-1.5 rounded-lg font-bold transition flex items-center space-x-1.5 text-xs select-none',
              route.path === '/admin'
                ? 'bg-white text-slate-900 shadow-xs'
                : 'text-slate-500 hover:text-slate-800'
            ]"
          >
            <span>🏢</span>
            <span>ฝั่งร้านค้า (Admin POS)</span>
          </button>
        </div>

        <button
          @click="router.push('/demo')"
          :class="[
            'px-2 py-1.5 rounded-lg text-[11px] font-medium transition text-slate-400 hover:text-slate-600 hover:bg-slate-100',
            route.path === '/demo' ? 'bg-slate-200 text-slate-900 font-bold' : ''
          ]"
          title="ดูแบบแบ่ง 2 จอเทียบกัน"
        >
          <span>Split View</span>
        </button>
      </nav>
    </header>

    <!-- Main View Content -->
    <main class="flex-1 flex flex-col">
      <router-view />
    </main>
  </div>
</template>
