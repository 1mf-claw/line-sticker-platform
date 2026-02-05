import { createApp, onMounted, ref } from 'vue'
import { api } from './api/client'
import type { Project } from './api/types'

createApp({
  setup() {
    const status = ref('Loading...')
    const project = ref<Project | null>(null)

    onMounted(async () => {
      try {
        const created = await api.createProject({
          title: 'LINE Sticker Project',
          stickerCount: 8,
        })
        project.value = created
        status.value = `Project ${created.id} (${created.status})`
      } catch (e) {
        status.value = 'Failed to connect API'
      }
    })

    return { status, project }
  },
  template: `
    <div>
      <h1>{{ status }}</h1>
      <pre v-if="project">{{ project }}</pre>
    </div>
  `,
}).mount('#app')
