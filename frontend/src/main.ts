import { createApp, ref } from 'vue'
import { api } from './api/client'
import type {
  CharacterCreateRequest,
  Draft,
  Project,
  Sticker,
} from './api/types'

type Step =
  | 'CREATE_PROJECT'
  | 'CHARACTER'
  | 'THEME'
  | 'DRAFTS'
  | 'GENERATE'
  | 'PREVIEW'

createApp({
  setup() {
    const step = ref<Step>('CREATE_PROJECT')
    const loading = ref(false)
    const error = ref('')

    const project = ref<Project | null>(null)
    const drafts = ref<Draft[]>([])
    const stickers = ref<Sticker[]>([])

    const title = ref('LINE Sticker Project')
    const theme = ref('')
    const stickerCount = ref<8 | 16 | 24 | 40>(8)

    const characterReq = ref<CharacterCreateRequest>({
      sourceType: 'AI',
      prompt: '圓臉橘貓，穿襯衫',
    })

    const next = (s: Step) => (step.value = s)

    const run = async (fn: () => Promise<void>) => {
      error.value = ''
      loading.value = true
      try {
        await fn()
      } catch (e) {
        error.value = '操作失敗，請稍後再試'
      } finally {
        loading.value = false
      }
    }

    const createProject = () =>
      run(async () => {
        project.value = await api.createProject({
          title: title.value,
          stickerCount: stickerCount.value,
        })
        next('CHARACTER')
      })

    const createCharacter = () =>
      run(async () => {
        if (!project.value) return
        await api.createCharacter(project.value.id, characterReq.value)
        next('THEME')
      })

    const updateTheme = () =>
      run(async () => {
        if (!project.value) return
        project.value = await api.updateProject(project.value.id, {
          theme: theme.value,
        })
        next('DRAFTS')
      })

    const generateDrafts = () =>
      run(async () => {
        if (!project.value) return
        await api.generateDrafts(project.value.id)
        drafts.value = await api.listDrafts(project.value.id)
        next('GENERATE')
      })

    const generateStickers = () =>
      run(async () => {
        if (!project.value) return
        await api.generateStickers(project.value.id)
        stickers.value = await api.listStickers(project.value.id)
        next('PREVIEW')
      })

    return {
      step,
      loading,
      error,
      project,
      drafts,
      stickers,
      title,
      theme,
      stickerCount,
      characterReq,
      createProject,
      createCharacter,
      updateTheme,
      generateDrafts,
      generateStickers,
    }
  },
  template: `
    <div style="max-width: 720px; margin: 32px auto; font-family: system-ui;">
      <h1>LINE 貼圖製作平台</h1>
      <p v-if="error" style="color: red;">{{ error }}</p>
      <p v-if="loading">處理中...</p>

      <section v-if="step === 'CREATE_PROJECT'">
        <h2>1. 建立專案</h2>
        <label>專案名稱</label>
        <input v-model="title" style="width: 100%; margin: 8px 0;" />

        <label>貼圖數量</label>
        <select v-model.number="stickerCount">
          <option :value="8">8</option>
          <option :value="16">16</option>
          <option :value="24">24</option>
          <option :value="40">40</option>
        </select>
        <button @click="createProject" style="margin-left: 12px;">下一步</button>
      </section>

      <section v-else-if="step === 'CHARACTER'">
        <h2>2. 角色設定</h2>
        <label>來源</label>
        <select v-model="characterReq.sourceType">
          <option value="AI">AI 生成</option>
          <option value="UPLOAD">上傳圖片</option>
          <option value="HISTORY">歷史角色</option>
        </select>

        <div v-if="characterReq.sourceType === 'AI'">
          <label>角色描述</label>
          <input v-model="characterReq.prompt" style="width: 100%; margin: 8px 0;" />
        </div>
        <div v-else>
          <label>Reference 圖片 URL</label>
          <input v-model="characterReq.referenceImageUrl" style="width: 100%; margin: 8px 0;" />
        </div>

        <button @click="createCharacter">下一步</button>
      </section>

      <section v-else-if="step === 'THEME'">
        <h2>3. 主題與焦點</h2>
        <label>主題</label>
        <input v-model="theme" style="width: 100%; margin: 8px 0;" />
        <button @click="updateTheme">下一步</button>
      </section>

      <section v-else-if="step === 'DRAFTS'">
        <h2>4. 草稿生成</h2>
        <p>將根據主題與數量生成草稿。</p>
        <button @click="generateDrafts">產生草稿</button>
      </section>

      <section v-else-if="step === 'GENERATE'">
        <h2>5. 草稿確認</h2>
        <ul>
          <li v-for="d in drafts" :key="d.id">#{{ d.index }} - {{ d.caption }} / {{ d.imagePrompt }}</li>
        </ul>
        <button @click="generateStickers">開始生成貼圖</button>
      </section>

      <section v-else-if="step === 'PREVIEW'">
        <h2>6. 預覽</h2>
        <ul>
          <li v-for="s in stickers" :key="s.id">{{ s.id }} - {{ s.status }}</li>
        </ul>
        <p>下一步：去背、壓縮下載（API 已準備）</p>
      </section>
    </div>
  `,
}).mount('#app')
