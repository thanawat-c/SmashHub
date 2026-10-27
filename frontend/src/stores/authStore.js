import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api'
import { ElMessage } from 'element-plus'

export const useAuthStore = defineStore('auth', () => {
  // Load persisted user & token from localStorage if available
  const savedUser = localStorage.getItem('smashhub_user')
  const savedToken = localStorage.getItem('smashhub_token')

  const user = ref(savedUser ? JSON.parse(savedUser) : null)
  const token = ref(savedToken || null)

  const isAuthModalOpen = ref(false)
  const authModalTab = ref('login') // 'login' or 'register'
  const isLoading = ref(false)
  const errorMessage = ref('')

  // Computed Getters
  const isAuthenticated = computed(() => Boolean(user.value))
  const isAdmin = computed(() => user.value?.role === 'ADMIN')
  const displayName = computed(() => user.value?.display_name || 'Guest')
  const avatar = computed(() => {
    return user.value?.picture_url || `https://api.dicebear.com/7.x/avataaars/svg?seed=${user.value?.phone || 'smashhub'}`
  })

  function setSession(userData, tokenStr) {
    user.value = userData
    token.value = tokenStr
    localStorage.setItem('smashhub_user', JSON.stringify(userData))
    localStorage.setItem('smashhub_token', tokenStr)
    api.defaults.headers.common['Authorization'] = `Bearer ${tokenStr}`
    if (userData?.id) {
      api.defaults.headers.common['X-User-Id'] = userData.id
    }
  }

  function clearSession() {
    user.value = null
    token.value = null
    localStorage.removeItem('smashhub_user')
    localStorage.removeItem('smashhub_token')
    delete api.defaults.headers.common['Authorization']
    delete api.defaults.headers.common['X-User-Id']
  }

  // Action: Register
  async function register({ displayName, phone, email, password }) {
    isLoading.value = true
    errorMessage.value = ''
    try {
      const res = await api.post('/auth/register', {
        display_name: displayName,
        phone,
        email,
        password
      })
      setSession(res.data.user, res.data.token)
      isAuthModalOpen.value = false
      ElMessage.success(`Registration successful! Welcome, ${res.data.user.display_name}`)
      return true
    } catch (err) {
      errorMessage.value = err.response?.data?.error || 'Registration failed. Please check your details.'
      ElMessage.error(errorMessage.value)
      return false
    } finally {
      isLoading.value = false
    }
  }

  // Action: Login
  async function login({ identifier, password }) {
    isLoading.value = true
    errorMessage.value = ''
    try {
      const res = await api.post('/auth/login', {
        identifier,
        password
      })
      setSession(res.data.user, res.data.token)
      isAuthModalOpen.value = false
      ElMessage.success(`Login successful! Welcome back, ${res.data.user.display_name}`)
      return true
    } catch (err) {
      errorMessage.value = err.response?.data?.error || 'Invalid username or password'
      ElMessage.error(errorMessage.value)
      return false
    } finally {
      isLoading.value = false
    }
  }

  // Action: Fast 1-Click Demo Login (for Admin & Customer test roles)
  async function demoLogin(role = 'CUSTOMER') {
    isLoading.value = true
    errorMessage.value = ''
    try {
      const res = await api.post('/auth/demo-login', { role })
      setSession(res.data.user, res.data.token)
      isAuthModalOpen.value = false
      ElMessage.success(`Switched session to ${role}: ${res.data.user.display_name}`)
      return true
    } catch (err) {
      errorMessage.value = err.response?.data?.error || 'Demo login failed'
      ElMessage.error(errorMessage.value)
      return false
    } finally {
      isLoading.value = false
    }
  }

  // Action: Logout
  async function logout() {
    try {
      await api.post('/auth/logout')
    } catch (err) {
      // Ignore network errors on logout
    } finally {
      clearSession()
      ElMessage.info('Logged out successfully')
    }
  }

  function openAuthModal(tab = 'login') {
    authModalTab.value = tab
    errorMessage.value = ''
    isAuthModalOpen.value = true
  }

  function closeAuthModal() {
    isAuthModalOpen.value = false
    errorMessage.value = ''
  }

  // Initialize API headers if already saved
  if (token.value) {
    api.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
  }
  if (user.value?.id) {
    api.defaults.headers.common['X-User-Id'] = user.value.id
  }

  return {
    user,
    token,
    isAuthenticated,
    isAdmin,
    displayName,
    avatar,
    isAuthModalOpen,
    authModalTab,
    isLoading,
    errorMessage,
    register,
    login,
    demoLogin,
    logout,
    openAuthModal,
    closeAuthModal
  }
})
