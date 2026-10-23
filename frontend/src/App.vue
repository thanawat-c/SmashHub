<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const isDemo = computed(() => route.path === '/demo')
</script>

<template>
  <div class="min-h-screen flex flex-col bg-slate-100">
    <!-- Quick Showcase Switcher Header (Visible in non-standalone mode) -->
    <header class="bg-slate-900 text-white border-b border-slate-800 px-4 py-2 flex items-center justify-between text-xs sticky top-0 z-50 shadow-md">
      <div class="flex items-center space-x-2 font-semibold">
        <span class="text-base">🏸</span>
        <span class="text-emerald-400 font-bold tracking-wide">SmashHub</span>
        <span class="text-slate-400 hidden sm:inline">| ระบบจองสนามแบดมินตัน & POS</span>
      </div>

      <nav class="flex items-center space-x-1 sm:space-x-2">
        <button
          @click="router.push('/demo')"
          :class="[
            'px-2.5 py-1 rounded-md font-medium transition flex items-center space-x-1',
            route.path === '/demo' ? 'bg-emerald-600 text-white shadow' : 'text-slate-300 hover:text-white hover:bg-slate-800'
          ]"
        >
          <span>🖥️📱</span>
          <span class="hidden sm:inline">Interactive</span>
          <span>Demo</span>
        </button>

        <button
          @click="router.push('/liff')"
          :class="[
            'px-2.5 py-1 rounded-md font-medium transition flex items-center space-x-1',
            route.path === '/liff' ? 'bg-emerald-600 text-white shadow' : 'text-slate-300 hover:text-white hover:bg-slate-800'
          ]"
        >
          <span>📱</span>
          <span>LINE LIFF (ลูกค้า)</span>
        </button>

        <button
          @click="router.push('/admin')"
          :class="[
            'px-2.5 py-1 rounded-md font-medium transition flex items-center space-x-1',
            route.path === '/admin' ? 'bg-emerald-600 text-white shadow' : 'text-slate-300 hover:text-white hover:bg-slate-800'
          ]"
        >
          <span>📊</span>
          <span>Admin & POS</span>
        </button>
      </nav>
    </header>

    <!-- Main View Content -->
    <main class="flex-1 flex flex-col">
      <router-view />
    </main>
  </div>
</template>
