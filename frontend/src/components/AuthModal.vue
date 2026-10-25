<script setup>
import { ref, reactive } from 'vue'
import { useAuthStore } from '@/stores/authStore'
import {
  User, Lock, Phone, Mail, Sparkles, LogIn, UserPlus,
  ShieldCheck, X, Check, ArrowRight
} from 'lucide-vue-next'

const authStore = useAuthStore()

// Forms data
const loginForm = reactive({
  identifier: '',
  password: ''
})

const registerForm = reactive({
  displayName: '',
  phone: '',
  email: '',
  password: ''
})

async function handleLogin() {
  if (!loginForm.identifier || !loginForm.password) return
  await authStore.login({
    identifier: loginForm.identifier,
    password: loginForm.password
  })
}

async function handleRegister() {
  if (!registerForm.displayName || !registerForm.phone || !registerForm.password) return
  await authStore.register({
    displayName: registerForm.displayName,
    phone: registerForm.phone,
    email: registerForm.email,
    password: registerForm.password
  })
}

async function handleQuickDemo(role) {
  await authStore.demoLogin(role)
}
</script>

<template>
  <el-dialog
    v-model="authStore.isAuthModalOpen"
    width="440px"
    align-center
    :show-close="false"
    class="rounded-3xl overflow-hidden p-0"
  >
    <!-- Modal Header -->
    <div class="bg-slate-900 text-white p-5 relative overflow-hidden">
      <!-- Subtle Badminton Background Icon -->
      <div class="absolute -right-4 -bottom-4 text-slate-800 text-7xl select-none opacity-30 pointer-events-none">
        🏸
      </div>

      <div class="flex items-center justify-between relative z-10">
        <div class="flex items-center space-x-2.5">
          <div class="w-9 h-9 rounded-xl bg-emerald-500/20 text-emerald-400 border border-emerald-500/30 flex items-center justify-center font-bold text-base">
            🏸
          </div>
          <div>
            <h2 class="text-sm font-extrabold tracking-tight text-white flex items-center space-x-1.5">
              <span>SmashHub Account</span>
              <ShieldCheck class="w-3.5 h-3.5 text-emerald-400" />
            </h2>
            <p class="text-[11px] text-slate-400">ระบบเข้าสู่ระบบและสมัครสมาชิกสนามแบดมินตัน</p>
          </div>
        </div>

        <button
          @click="authStore.closeAuthModal"
          class="p-1.5 text-slate-400 hover:text-white rounded-lg hover:bg-slate-800 transition"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Segmented Tab Switcher -->
      <div class="mt-4 bg-slate-950/70 p-1 rounded-xl flex border border-slate-800">
        <button
          @click="authStore.authModalTab = 'login'"
          :class="[
            'flex-1 py-1.5 rounded-lg text-xs font-bold transition flex items-center justify-center space-x-1.5',
            authStore.authModalTab === 'login'
              ? 'bg-emerald-600 text-white shadow-xs'
              : 'text-slate-400 hover:text-white'
          ]"
        >
          <LogIn class="w-3.5 h-3.5" />
          <span>เข้าสู่ระบบ</span>
        </button>

        <button
          @click="authStore.authModalTab = 'register'"
          :class="[
            'flex-1 py-1.5 rounded-lg text-xs font-bold transition flex items-center justify-center space-x-1.5',
            authStore.authModalTab === 'register'
              ? 'bg-emerald-600 text-white shadow-xs'
              : 'text-slate-400 hover:text-white'
          ]"
        >
          <UserPlus class="w-3.5 h-3.5" />
          <span>สมัครสมาชิก</span>
        </button>
      </div>
    </div>

    <!-- Modal Body Content -->
    <div class="p-5 bg-white space-y-4">
      <!-- Error Alert -->
      <div
        v-if="authStore.errorMessage"
        class="p-2.5 bg-rose-50 border border-rose-200 text-rose-700 text-xs rounded-xl flex items-center space-x-2 animate-in fade-in"
      >
        <span>⚠️ {{ authStore.errorMessage }}</span>
      </div>

      <!-- TAB 1: LOGIN FORM -->
      <form v-if="authStore.authModalTab === 'login'" @submit.prevent="handleLogin" class="space-y-3">
        <div>
          <label class="block text-[11px] font-semibold text-slate-700 mb-1">
            เบอร์โทรศัพท์ หรือ อีเมล
          </label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400">
              <Phone class="w-3.5 h-3.5" />
            </div>
            <input
              v-model="loginForm.identifier"
              type="text"
              placeholder="เช่น 0891112233 หรือ email@domain.com"
              required
              class="w-full pl-9 pr-3 py-2 text-xs border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-600 transition"
            />
          </div>
        </div>

        <div>
          <label class="block text-[11px] font-semibold text-slate-700 mb-1">
            รหัสผ่าน
          </label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400">
              <Lock class="w-3.5 h-3.5" />
            </div>
            <input
              v-model="loginForm.password"
              type="password"
              placeholder="กรอกรหัสผ่านของคุณ"
              required
              class="w-full pl-9 pr-3 py-2 text-xs border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-600 transition"
            />
          </div>
        </div>

        <button
          type="submit"
          :disabled="authStore.isLoading"
          class="w-full py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-xs font-bold shadow-xs transition flex items-center justify-center space-x-1.5 active:scale-95 disabled:opacity-50"
        >
          <span v-if="authStore.isLoading" class="w-3.5 h-3.5 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
          <span v-else>เข้าสู่ระบบ</span>
          <ArrowRight v-if="!authStore.isLoading" class="w-3.5 h-3.5" />
        </button>
      </form>

      <!-- TAB 2: REGISTER FORM -->
      <form v-else-if="authStore.authModalTab === 'register'" @submit.prevent="handleRegister" class="space-y-3">
        <div>
          <label class="block text-[11px] font-semibold text-slate-700 mb-1">
            ชื่อ-นามสกุล หรือชื่อที่แสดง
          </label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400">
              <User class="w-3.5 h-3.5" />
            </div>
            <input
              v-model="registerForm.displayName"
              type="text"
              placeholder="เช่น กานต์ ธนบดี"
              required
              class="w-full pl-9 pr-3 py-2 text-xs border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-600 transition"
            />
          </div>
        </div>

        <div>
          <label class="block text-[11px] font-semibold text-slate-700 mb-1">
            เบอร์โทรศัพท์ (สำหรับยืนยันการจอง)
          </label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400">
              <Phone class="w-3.5 h-3.5" />
            </div>
            <input
              v-model="registerForm.phone"
              type="tel"
              placeholder="เช่น 0891112233"
              required
              class="w-full pl-9 pr-3 py-2 text-xs border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-600 transition"
            />
          </div>
        </div>

        <div>
          <label class="block text-[11px] font-semibold text-slate-700 mb-1">
            อีเมล (ไม่บังคับ)
          </label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400">
              <Mail class="w-3.5 h-3.5" />
            </div>
            <input
              v-model="registerForm.email"
              type="email"
              placeholder="yourname@gmail.com"
              class="w-full pl-9 pr-3 py-2 text-xs border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-600 transition"
            />
          </div>
        </div>

        <div>
          <label class="block text-[11px] font-semibold text-slate-700 mb-1">
            ตั้งรหัสผ่าน (อย่างน้อย 4 ตัวอักษร)
          </label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400">
              <Lock class="w-3.5 h-3.5" />
            </div>
            <input
              v-model="registerForm.password"
              type="password"
              placeholder="••••••••"
              minlength="4"
              required
              class="w-full pl-9 pr-3 py-2 text-xs border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-600 transition"
            />
          </div>
        </div>

        <button
          type="submit"
          :disabled="authStore.isLoading"
          class="w-full py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-xs font-bold shadow-xs transition flex items-center justify-center space-x-1.5 active:scale-95 disabled:opacity-50"
        >
          <span v-if="authStore.isLoading" class="w-3.5 h-3.5 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
          <span v-else>ยืนยันการสมัครสมาชิก</span>
          <ArrowRight v-if="!authStore.isLoading" class="w-3.5 h-3.5" />
        </button>
      </form>

      <!-- 1-Click Fast Demo Login Section for Evaluators & Recruiters -->
      <div class="border-t border-slate-100 pt-3">
        <div class="text-[10px] uppercase tracking-wider font-bold text-slate-400 mb-2 flex items-center space-x-1">
          <Sparkles class="w-3 h-3 text-amber-500" />
          <span>เข้าสู่ระบบด่วน 1-Click (สำหรับทดสอบระบบ)</span>
        </div>

        <div class="grid grid-cols-2 gap-2">
          <button
            @click="handleQuickDemo('CUSTOMER')"
            class="p-2 rounded-xl bg-slate-50 hover:bg-slate-100 border border-slate-200 text-left transition flex items-center space-x-2"
          >
            <div class="w-6 h-6 rounded-lg bg-blue-100 text-blue-800 flex items-center justify-center text-xs">
              🏸
            </div>
            <div class="min-w-0">
              <div class="text-[11px] font-bold text-slate-800 truncate">กานต์ (ลูกค้า)</div>
              <div class="text-[9px] text-slate-400 font-mono">Customer Role</div>
            </div>
          </button>

          <button
            @click="handleQuickDemo('ADMIN')"
            class="p-2 rounded-xl bg-slate-50 hover:bg-slate-100 border border-slate-200 text-left transition flex items-center space-x-2"
          >
            <div class="w-6 h-6 rounded-lg bg-emerald-100 text-emerald-800 flex items-center justify-center text-xs">
              👑
            </div>
            <div class="min-w-0">
              <div class="text-[11px] font-bold text-slate-800 truncate">สมชาย (Admin)</div>
              <div class="text-[9px] text-emerald-600 font-mono font-semibold">Admin POS</div>
            </div>
          </button>
        </div>
      </div>
    </div>
  </el-dialog>
</template>
