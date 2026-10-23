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
      customer_name: customerProfile.displayName,
      customer_phone: customerProfile.phone,
      line_user_id: customerProfile.lineUserId,
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
    verifySlip
  }
})
