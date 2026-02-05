import { createApp, ref, onMounted } from 'vue'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL

createApp({
  setup() {
    const msg = ref('Loading...')

    onMounted(async () => {
      try {
        const res = await fetch(`${API_BASE_URL}/api/message`)
        const data = await res.json()
        msg.value = data.message
      } catch (e) {
        msg.value = 'Failed to connect backend'
      }
    })

    return { msg }
  },
  template: `<h1>{{ msg }}</h1>`
}).mount('#app')
