import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLineProfileStore = defineStore('lineProfile', () => {
  const mockProfiles = [
    {
      id: 'b0000000-0000-0000-0000-000000000001',
      lineUserId: 'U_CUSTOMER_KAN',
      displayName: 'กานต์ ธนบดี',
      phone: '089-111-2233',
      pictureUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=kan',
      statusMessage: 'ตีแบดทุกวันเสาร์ 🏸'
    },
    {
      id: 'b0000000-0000-0000-0000-000000000002',
      lineUserId: 'U_CUSTOMER_NUT',
      displayName: 'ณัฐพงษ์ วงศ์วิวัฒน์',
      phone: '086-444-5566',
      pictureUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=nuttapong',
      statusMessage: 'มือใหม่หัดตบ smash power'
    },
    {
      id: 'b0000000-0000-0000-0000-000000000003',
      lineUserId: 'U_CUSTOMER_TNP',
      displayName: 'ธนาภา สินสมบูรณ์',
      phone: '082-777-8899',
      pictureUrl: 'https://api.dicebear.com/7.x/avataaars/svg?seed=thanapa',
      statusMessage: 'หาเพื่อนก๊วนแบดมินตัน'
    }
  ]

  const currentProfile = ref(mockProfiles[0])

  function selectProfile(profile) {
    currentProfile.value = profile
  }

  return {
    mockProfiles,
    currentProfile,
    selectProfile
  }
})
