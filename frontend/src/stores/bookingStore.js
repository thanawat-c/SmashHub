import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../api'

export const useBookingStore = defineStore('booking', () => {
  const selectedDate = ref(new Date().toISOString().split('T')[0])
  const courts = ref([])
  const products = ref([])
  const selectedSlots = ref([]) // array of { courtId, courtNumber, courtName, startHour, endHour, hourlyRate }
  const cartItems = ref({}) // { [productId]: { product, qty } }
  
  const currentStep = ref(1) // 1: Court Matrix, 2: Mini Store, 3: Payment/Slip, 4: E-Ticket
  const activeBooking = ref(null)
  const slipResult = ref(null)
  const isSubmitting = ref(false)
  const isLoading = ref(false)
  const errorMessage = ref('')

  // Computed Totals
  const courtTotal = computed(() => {
    return selectedSlots.value.reduce((sum, s) => sum + s.hourlyRate, 0)
  })

  const itemTotal = computed(() => {
    return Object.values(cartItems.value).reduce((sum, item) => sum + (item.product.price * item.qty), 0)
  })

  const netTotal = computed(() => {
    return courtTotal.value + itemTotal.value
  })

  const cartItemsList = computed(() => {
    return Object.values(cartItems.value)
  })

  const totalCartCount = computed(() => {
    return Object.values(cartItems.value).reduce((sum, item) => sum + item.qty, 0)
  })

  // Actions
  async function fetchAvailability(date) {
    if (date) selectedDate.value = date
    isLoading.value = true
    errorMessage.value = ''
    try {
      const res = await api.get(`/courts/availability?date=${selectedDate.value}`)
      courts.value = res.data.courts || []
    } catch (err) {
      errorMessage.value = err.response?.data?.error || 'ไม่สามารถโหลดข้อมูลคอร์ทได้'
    } finally {
      isLoading.value = false
    }
  }

  async function fetchProducts() {
    try {
      const res = await api.get('/products')
      products.value = res.data.data || []
    } catch (err) {
      console.error('Error fetching products:', err)
    }
  }

  function isSlotSelected(courtId, startHour) {
    return selectedSlots.value.some(s => s.courtId === courtId && s.startHour === startHour)
  }

  function toggleCourtSlot(court, slot) {
    if (slot.status !== 'AVAILABLE') return

    const index = selectedSlots.value.findIndex(s => s.courtId === court.court_id && s.startHour === slot.start_hour)
    if (index >= 0) {
      selectedSlots.value.splice(index, 1)
    } else {
      selectedSlots.value.push({
        courtId: court.court_id,
        courtNumber: court.court_number,
        courtName: court.court_name,
        courtType: court.court_type,
        startHour: slot.start_hour,
        endHour: slot.end_hour,
        hourlyRate: slot.hourly_rate
      })
    }
  }

  function addToCart(product) {
    if (!product || product.stock_qty <= 0) return
    const currentQty = cartItems.value[product.id]?.qty || 0
    if (currentQty >= product.stock_qty) return

    cartItems.value[product.id] = {
      product,
      qty: currentQty + 1
    }
  }

  function removeFromCart(productId) {
    if (!cartItems.value[productId]) return
    if (cartItems.value[productId].qty > 1) {
      cartItems.value[productId].qty -= 1
    } else {
      delete cartItems.value[productId]
    }
  }

  function resetSelection() {
    selectedSlots.value = []
    cartItems.value = {}
    currentStep.value = 1
    activeBooking.value = null
    slipResult.value = null
    errorMessage.value = ''
  }

  async function submitBooking(customerProfile) {
    if (selectedSlots.value.length === 0) {
      errorMessage.value = 'กรุณาเลือกอย่างน้อย 1 ช่วงเวลาคอร์ท'
      return false
    }

    isSubmitting.value = true
    errorMessage.value = ''

    const payload = {
      customer_name: customerProfile?.displayName || customerProfile?.display_name || 'ลูกค้าทั่วไป',
      customer_phone: customerProfile?.phone || '089-111-2233',
      line_user_id: customerProfile?.lineUserId || customerProfile?.line_user_id,
      play_date: selectedDate.value,
      court_slots: selectedSlots.value.map(s => ({
        court_id: s.courtId,
        start_hour: s.startHour
      })),
      items: Object.values(cartItems.value).map(item => ({
        product_id: item.product.id,
        qty: item.qty
      }))
    }

    try {
      const res = await api.post('/bookings', payload)
      activeBooking.value = res.data.data
      currentStep.value = 3 // Move to payment step
      return true
    } catch (err) {
      errorMessage.value = err.response?.data?.error || 'เกิดข้อผิดพลาดในการสร้างการจอง'
      return false
    } finally {
      isSubmitting.value = false
    }
  }

  async function ensureActiveBookingForDemo(customerProfile) {
    if (activeBooking.value?.booking?.id && currentStep.value === 3) {
      return activeBooking.value.booking
    }

    // Auto-select an available slot
    await fetchAvailability()
    let targetCourt = courts.value[0]
    let targetSlot = targetCourt?.slots?.find(s => s.status === 'AVAILABLE')
    if (!targetSlot) {
      // Find across all courts
      for (const c of courts.value) {
        const s = c.slots.find(sl => sl.status === 'AVAILABLE')
        if (s) {
          targetCourt = c
          targetSlot = s
          break
        }
      }
    }

    if (targetCourt && targetSlot) {
      selectedSlots.value = [{
        courtId: targetCourt.court_id,
        courtNumber: targetCourt.court_number,
        courtName: targetCourt.court_name,
        courtType: targetCourt.court_type,
        startHour: targetSlot.start_hour,
        endHour: targetSlot.end_hour,
        hourlyRate: targetSlot.hourly_rate
      }]
    }

    // Auto add 1 Pocari if available
    await fetchProducts()
    if (products.value.length > 0) {
      const pocari = products.value.find(p => p.sku === 'DRINK-POCARI-SWEAT') || products.value[0]
      cartItems.value = {
        [pocari.id]: { product: pocari, qty: 1 }
      }
    }

    await submitBooking(customerProfile)
    return activeBooking.value?.booking
  }

  async function verifySlip(mockScenario = 'MOCK_VALID', customFile = null) {
    if (!activeBooking.value?.booking?.id) return false
    isSubmitting.value = true
    errorMessage.value = ''

    try {
      let res
      if (customFile) {
        const formData = new FormData()
        formData.append('slip_image', customFile)
        formData.append('mock_scenario', mockScenario)
        res = await api.post(`/bookings/${activeBooking.value.booking.id}/verify-slip`, formData, {
          headers: { 'Content-Type': 'multipart/form-data' }
        })
      } else {
        res = await api.post(`/bookings/${activeBooking.value.booking.id}/verify-slip`, {
          mock_scenario: mockScenario
        })
      }

      slipResult.value = res.data.data
      if (res.data.data.booking_status === 'CONFIRMED') {
        currentStep.value = 4 // Move to E-Ticket
      }
      return res.data.data
    } catch (err) {
      const detail = err.response?.data?.details || err.response?.data
      slipResult.value = detail
      errorMessage.value = err.response?.data?.error || 'การตรวจสอบสลิปไม่ผ่าน'
      return detail
    } finally {
      isSubmitting.value = false
    }
  }

  // Customer Booking History & Order Retrieval
  const myBookings = ref([])
  const isMyBookingsOpen = ref(false)
  const isFetchingMyBookings = ref(false)

  async function fetchMyBookings(userParams = null) {
    isFetchingMyBookings.value = true
    try {
      const params = {}
      if (userParams?.id) params.user_id = userParams.id
      if (userParams?.phone) params.phone = userParams.phone
      if (userParams?.lineUserId || userParams?.line_user_id) {
        params.line_user_id = userParams.lineUserId || userParams.line_user_id
      }
      const res = await api.get('/bookings/my-bookings', { params })
      myBookings.value = res.data.data || []
    } catch (err) {
      console.error('Error fetching my bookings:', err)
    } finally {
      isFetchingMyBookings.value = false
    }
  }

  function openMyBookings(userParams = null) {
    isMyBookingsOpen.value = true
    fetchMyBookings(userParams)
  }

  function closeMyBookings() {
    isMyBookingsOpen.value = false
  }

  function viewBookingTicket(booking) {
    activeBooking.value = {
      booking,
      promptpay_qr: `00020101021129370016A000000677010111011300669876543215802TH5303764540${booking.net_total}5802TH6304`
    }
    if (booking.courts?.[0]?.play_date) {
      selectedDate.value = booking.courts[0].play_date.split('T')[0]
    }
    currentStep.value = 4
    isMyBookingsOpen.value = false
  }

  function resumeBookingPayment(booking) {
    activeBooking.value = {
      booking,
      promptpay_qr: `00020101021129370016A000000677010111011300669876543215802TH5303764540${booking.net_total}5802TH6304`
    }
    if (booking.courts?.[0]?.play_date) {
      selectedDate.value = booking.courts[0].play_date.split('T')[0]
    }
    currentStep.value = 3
    isMyBookingsOpen.value = false
  }

  async function cancelBooking(bookingId, reason = 'ลูกค้ายกเลิกการจอง') {
    try {
      const res = await api.post(`/bookings/${bookingId}/cancel`, { reason })
      // Update in local myBookings list
      const idx = myBookings.value.findIndex(b => b.id === bookingId)
      if (idx !== -1) {
        myBookings.value[idx].status = 'CANCELLED'
      }
      // If active booking was cancelled, reset
      if (activeBooking.value?.booking?.id === bookingId) {
        activeBooking.value.booking.status = 'CANCELLED'
      }
      await fetchAvailability()
      return { success: true, data: res.data.data }
    } catch (err) {
      const msg = err.response?.data?.error || 'ไม่สามารถยกเลิกการจองได้'
      return { success: false, error: msg }
    }
  }

  return {
    selectedDate,
    courts,
    products,
    selectedSlots,
    cartItems,
    currentStep,
    activeBooking,
    slipResult,
    isSubmitting,
    isLoading,
    errorMessage,
    courtTotal,
    itemTotal,
    netTotal,
    cartItemsList,
    totalCartCount,
    fetchAvailability,
    fetchProducts,
    isSlotSelected,
    toggleCourtSlot,
    addToCart,
    removeFromCart,
    resetSelection,
    submitBooking,
    ensureActiveBookingForDemo,
    verifySlip,
    myBookings,
    isMyBookingsOpen,
    isFetchingMyBookings,
    fetchMyBookings,
    openMyBookings,
    closeMyBookings,
    viewBookingTicket,
    resumeBookingPayment,
    cancelBooking
  }
})
