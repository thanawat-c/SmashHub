import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api'
import { ElNotification } from 'element-plus'

export const useAdminStore = defineStore('admin', () => {
  const bookings = ref([])
  const filterDate = ref(new Date().toISOString().split('T')[0])
  const filterStatus = ref('')
  const filterCourtId = ref('')
  const searchQuery = ref('')
  const isLoading = ref(false)
  const isSSEConnected = ref(false)
  let eventSource = null

  // Drawer & Modal state
  const isDrawerOpen = ref(false)
  const selectedBooking = ref(null)

  // Receipt Modal state
  const isReceiptOpen = ref(false)
  const receiptBooking = ref(null)

  // POS Add Item Dialog state
  const isPOSOpen = ref(false)
  const posTargetBooking = ref(null)

  async function fetchBookings() {
    isLoading.value = true
    try {
      const params = {}
      if (filterDate.value) params.date = filterDate.value
      if (filterStatus.value) params.status = filterStatus.value
      if (filterCourtId.value) params.court_id = filterCourtId.value
      if (searchQuery.value) params.search = searchQuery.value

      const res = await api.get('/admin/bookings', { params })
      bookings.value = res.data.data || []
    } catch (err) {
      console.error('Error fetching admin bookings:', err)
    } finally {
      isLoading.value = false
    }
  }

  async function manualApprove(bookingId, action, note = '') {
    try {
      const res = await api.patch(`/admin/bookings/${bookingId}/manual-approve`, {
        action,
        note
      })
      ElNotification({
        title: action === 'APPROVE' ? 'อนุมัติเรียบร้อย' : 'ปฏิเสธเรียบร้อย',
        message: res.data.message,
        type: action === 'APPROVE' ? 'success' : 'warning'
      })
      await fetchBookings()
      if (selectedBooking.value?.id === bookingId) {
        selectedBooking.value = res.data.data
      }
      return true
    } catch (err) {
      ElNotification({
        title: 'เกิดข้อผิดพลาด',
        message: err.response?.data?.error || 'ไม่สามารถดำเนินการได้',
        type: 'error'
      })
      return false
    }
  }

  async function updateCheckIn(bookingId, status) {
    try {
      const res = await api.patch(`/admin/bookings/${bookingId}/check-in`, { status })
      ElNotification({
        title: 'อัปเดตสถานะสำเร็จ',
        message: `เช็คอิน: ${status}`,
        type: 'success'
      })
      await fetchBookings()
      return true
    } catch (err) {
      ElNotification({
        title: 'เกิดข้อผิดพลาด',
        message: err.response?.data?.error || 'ไม่สามารถอัปเดตการเช็คอินได้',
        type: 'error'
      })
      return false
    }
  }

  async function addItems(bookingId, items) {
    try {
      const res = await api.post(`/admin/bookings/${bookingId}/add-items`, { items })
      ElNotification({
        title: 'เพิ่มสินค้าสำเร็จ',
        message: res.data.message,
        type: 'success'
      })
      await fetchBookings()
      isPOSOpen.value = false
      return true
    } catch (err) {
      ElNotification({
        title: 'เกิดข้อผิดพลาด',
        message: err.response?.data?.error || 'ไม่สามารถเพิ่มสินค้าได้',
        type: 'error'
      })
      return false
    }
  }

  function openDrawer(booking) {
    selectedBooking.value = booking
    isDrawerOpen.value = true
  }

  function openReceipt(booking) {
    receiptBooking.value = booking
    isReceiptOpen.value = true
  }

  function openPOS(booking) {
    posTargetBooking.value = booking
    isPOSOpen.value = true
  }

  // Connect Server-Sent Events (SSE)
  function connectSSE() {
    if (eventSource) return

    eventSource = new EventSource('/api/v1/events/stream')

    eventSource.addEventListener('CONNECTED', (e) => {
      isSSEConnected.value = true
      console.log('[SSE] Connected to SmashHub live events')
    })

    eventSource.addEventListener('BOOKING_CREATED', (e) => {
      const data = JSON.parse(e.data)
      ElNotification({
        title: '🏸 มีการจองคอร์ทใหม่',
        message: `รหัส ${data.booking_no} โดยคุณ ${data.customer} (ยอด ${data.net_total} บ.)`,
        type: 'info'
      })
      fetchBookings()
    })

    eventSource.addEventListener('BOOKING_CONFIRMED', (e) => {
      const data = JSON.parse(e.data)
      ElNotification({
        title: '✅ ชำระเงินสำเร็จแล้ว',
        message: `รหัส ${data.booking_no} ได้รับการยืนยันแล้ว`,
        type: 'success'
      })
      fetchBookings()
    })

    eventSource.addEventListener('BOOKING_MANUAL_REVIEW', (e) => {
      const data = JSON.parse(e.data)
      ElNotification({
        title: '⚠️ สลิปรอแอดมินตรวจสอบ',
        message: `รหัส ${data.booking_no} ระบบตรวจสลิปขัดข้อง กรุณาตรวจสอบด้วยตนเอง`,
        type: 'warning'
      })
      fetchBookings()
    })

    eventSource.addEventListener('BOOKING_EXPIRED', (e) => {
      const data = JSON.parse(e.data)
      ElNotification({
        title: '⏱️ สล็อตคอร์ทหมดอายุ',
        message: `รหัส ${data.booking_no} หมดเวลาชำระเงิน คืนสล็อตคอร์ทแล้ว`,
        type: 'info'
      })
      fetchBookings()
    })

    eventSource.addEventListener('ITEMS_ADDED', (e) => {
      fetchBookings()
    })

    eventSource.onerror = () => {
      isSSEConnected.value = false
    }
  }

  function disconnectSSE() {
    if (eventSource) {
      eventSource.close()
      eventSource = null
      isSSEConnected.value = false
    }
  }

  return {
    bookings,
    filterDate,
    filterStatus,
    filterCourtId,
    searchQuery,
    isLoading,
    isSSEConnected,
    isDrawerOpen,
    selectedBooking,
    isReceiptOpen,
    receiptBooking,
    isPOSOpen,
    posTargetBooking,
    fetchBookings,
    manualApprove,
    updateCheckIn,
    addItems,
    openDrawer,
    openReceipt,
    openPOS,
    connectSSE,
    disconnectSSE
  }
})
