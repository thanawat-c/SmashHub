<script setup>
import LiffView from '../liff/LiffView.vue'
import AdminDashboard from '../admin/AdminDashboard.vue'
import { Monitor, Smartphone, Sparkles, RefreshCw } from 'lucide-vue-next'
</script>

<template>
  <div class="flex-1 flex flex-col bg-slate-900 min-h-[calc(100vh-45px)] text-white">
    <!-- Demo Header Explanation -->
    <div class="bg-slate-950 border-b border-slate-800 px-4 py-2 flex items-center justify-between text-xs">
      <div class="flex items-center space-x-2">
        <Sparkles class="w-4 h-4 text-emerald-400" />
        <span class="font-bold text-slate-200">Portfolio Interactive Showcase:</span>
        <span class="text-slate-400 hidden md:inline">
          จองคอร์ทบนมือถือฝั่งซ้าย แล้วสังเกตข้อมูลและสถานะอัปเดตเรียลไทม์บนกระดานแอดมินฝั่งขวาผ่าน Server-Sent Events (SSE)
        </span>
      </div>

      <div class="flex items-center space-x-2 text-[11px] text-emerald-400 font-mono">
        <span class="w-2 h-2 rounded-full bg-emerald-500 animate-ping"></span>
        <span>REAL-TIME DUAL SYNC</span>
      </div>
    </div>

    <!-- Dual Split View Layout -->
    <div class="flex-1 grid grid-cols-1 lg:grid-cols-12 gap-0 overflow-hidden">
      <!-- Left: Mobile LINE LIFF (5 cols on large screens) -->
      <div class="lg:col-span-5 bg-slate-950/80 p-2 sm:p-4 flex items-center justify-center border-b lg:border-b-0 lg:border-r border-slate-800 overflow-y-auto">
        <div class="w-full max-w-[430px]">
          <div class="flex items-center justify-between text-[11px] text-slate-400 mb-1 px-2">
            <span class="flex items-center space-x-1">
              <Smartphone class="w-3.5 h-3.5 text-emerald-400" />
              <span>ฝั่งลูกค้า (LINE LIFF Storefront)</span>
            </span>
            <span class="text-[10px] text-emerald-400 font-mono">Mobile View</span>
          </div>

          <!-- Embedded LIFF Component -->
          <div class="rounded-[28px] overflow-hidden border-4 border-slate-700 shadow-2xl bg-white max-h-[820px] flex flex-col">
            <LiffView />
          </div>
        </div>
      </div>

      <!-- Right: Desktop Admin Dashboard (7 cols on large screens) -->
      <div class="lg:col-span-7 bg-slate-100 text-slate-800 overflow-y-auto flex flex-col">
        <div class="bg-slate-200 px-4 py-1.5 text-[11px] text-slate-600 font-semibold border-b border-slate-300 flex items-center justify-between">
          <span class="flex items-center space-x-1.5">
            <Monitor class="w-3.5 h-3.5 text-slate-700" />
            <span>ฝั่งร้านค้า & แคชเชียร์ (Admin Backoffice & Live POS)</span>
          </span>
          <span class="text-[10px] text-slate-500 font-mono">Desktop View</span>
        </div>

        <div class="flex-1 p-2 sm:p-4">
          <AdminDashboard />
        </div>
      </div>
    </div>
  </div>
</template>
