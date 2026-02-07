import { createApp, onMounted, ref, watch } from 'vue'
import Notice from './components/Notice.vue'
import { api } from './api/client'
import type {
  CharacterCreateRequest,
  Draft,
  Project,
  Provider,
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
    const success = ref('')
    const verified = ref(false)

    const project = ref<Project | null>(null)
    const drafts = ref<Draft[]>([])
    const stickers = ref<Sticker[]>([])
    const downloadUrl = ref('')

    const jobStatus = ref('')
    const jobProgress = ref(0)

    const title = ref('LINE Sticker Project')
    const theme = ref('')
    const stickerCount = ref<8 | 16 | 24 | 40>(8)

    const providers = ref<Provider[]>([])
    const verifiedProviders = ref<Provider[]>([])
    const selectedProvider = ref('openai')
    const selectedModel = ref('gpt-4o-mini')
    const customModel = ref('')

    const textProvider = ref('openai')
    const textModel = ref('gpt-4o-mini')
    const textCustom = ref('')

    const imageProvider = ref('openai')
    const imageModel = ref('gpt-4o-mini')
    const imageCustom = ref('')

    const bgProvider = ref('openai')
    const bgModel = ref('gpt-4o-mini')
    const bgCustom = ref('')
    const apiKey = ref('')
    const apiBase = ref('')

    const characterReq = ref<CharacterCreateRequest>({
      sourceType: 'AI',
      prompt: '圓臉橘貓，穿襯衫',
    })

    const next = (s: Step) => (step.value = s)

    onMounted(async () => {
      try {
        providers.value = await api.listProviders()
      } catch {
        // ignore provider fetch errors for MVP
      }
    })

    watch(selectedProvider, (val) => {
      const p = providers.value.find((x) => x.id === val)
      if (p && p.models.length > 0) {
        selectedModel.value = p.models[0]
      }
      customModel.value = ''
    })

    watch(apiKey, () => {
      verified.value = false
      verifiedProviders.value = []
    })

    const syncDefaultModel = (providerRef: any, modelRef: any, customRef: any) => {
      const p = providers.value.find((x) => x.id === providerRef.value)
      if (p && p.models.length > 0) {
        modelRef.value = p.models[0]
      }
      customRef.value = ''
    }

    const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms))

    const pollJob = async (jobId: string) => {
      jobStatus.value = 'RUNNING'
      jobProgress.value = 0
      for (let i = 0; i < 20; i++) {
        const job = await api.getJob(jobId)
        jobStatus.value = job.status
        jobProgress.value = job.progress ?? 0
        if (job.status === 'SUCCESS' || job.status === 'FAILED') {
          return job.status
        }
        await sleep(500)
      }
      return 'RUNNING'
    }

    const run = async (fn: () => Promise<void>) => {
      error.value = ''
      success.value = ''
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
        if (!apiKey.value) {
          error.value = '請先填入 API Key'
          return
        }
        const textModelValue = textCustom.value || textModel.value
        const imageModelValue = imageCustom.value || imageModel.value
        const bgModelValue = bgCustom.value || bgModel.value

        if (!textProvider.value || !textModelValue) {
          error.value = '請完成「文字生成」的供應商與模型（或自訂模型 ID）'
          return
        }
        if (!imageProvider.value || !imageModelValue) {
          error.value = '請完成「圖像生成」的供應商與模型（或自訂模型 ID）'
          return
        }
        if (!bgProvider.value || !bgModelValue) {
          error.value = '請完成「去背」的供應商與模型（或自訂模型 ID）'
          return
        }

        await api.setAICredentials(project.value.id, {
          aiProvider: selectedProvider.value,
          apiKey: apiKey.value,
          apiBase: apiBase.value || undefined,
        })
        await api.updateAIConfig(project.value.id, {
          aiProvider: selectedProvider.value,
          aiModel: customModel.value || selectedModel.value,
        })
        await api.updateAIPipeline(project.value.id, {
          textProvider: textProvider.value,
          textModel: textModelValue,
          imageProvider: imageProvider.value,
          imageModel: imageModelValue,
          bgProvider: bgProvider.value,
          bgModel: bgModelValue,
        })
        await api.verifyAICredentials(project.value.id)
        const verified = await api.listVerifiedProviders(project.value.id)
        verifiedProviders.value = verified.providers || []
        verified.value = verifiedProviders.value.length > 0
        if (verifiedProviders.value.length > 0) {
          const p0 = verifiedProviders.value[0]
          selectedProvider.value = p0.id
          selectedModel.value = p0.models[0] || ''
          textProvider.value = p0.id
          textModel.value = p0.models[0] || ''
          imageProvider.value = p0.id
          imageModel.value = p0.models[0] || ''
          bgProvider.value = p0.id
          bgModel.value = p0.models[0] || ''
        }
        project.value = await api.updateProject(project.value.id, {
          theme: theme.value,
        })
        success.value = 'API Key 驗證成功'
        next('DRAFTS')
      })

    const generateDrafts = () =>
      run(async () => {
        if (!project.value) return
        const job = await api.generateDrafts(project.value.id)
        await pollJob(job.id)
        drafts.value = await api.listDrafts(project.value.id)
        next('GENERATE')
      })

    const saveDraft = (d: Draft) =>
      run(async () => {
        await api.updateDraft(d.id, {
          caption: d.caption,
          imagePrompt: d.imagePrompt,
        })
      })

    const regenerateDrafts = () =>
      run(async () => {
        if (!project.value) return
        const job = await api.generateDrafts(project.value.id)
        await pollJob(job.id)
        drafts.value = await api.listDrafts(project.value.id)
      })

    const generateStickers = () =>
      run(async () => {
        if (!project.value) return
        const job = await api.generateStickers(project.value.id)
        await pollJob(job.id)
        stickers.value = await api.listStickers(project.value.id)
        next('PREVIEW')
      })

    const regenerateSticker = (stickerId: string) =>
      run(async () => {
        if (!project.value) return
        const job = await api.regenerateSticker(stickerId)
        await pollJob(job.id)
        stickers.value = await api.listStickers(project.value.id)
      })

    const removeBackground = () =>
      run(async () => {
        if (!project.value) return
        const job = await api.removeBackground(project.value.id)
        await pollJob(job.id)
        stickers.value = await api.listStickers(project.value.id)
      })

    const exportZip = () =>
      run(async () => {
        if (!project.value) return
        const res = await api.exportZip(project.value.id)
        downloadUrl.value = res.downloadUrl
      })

    return {
      step,
      loading,
      error,
      success,
      verified,
      project,
      drafts,
      stickers,
      downloadUrl,
      jobStatus,
      jobProgress,
      title,
      theme,
      stickerCount,
      providers,
      verifiedProviders,
      selectedProvider,
      selectedModel,
      customModel,
      textProvider,
      textModel,
      textCustom,
      imageProvider,
      imageModel,
      imageCustom,
      bgProvider,
      bgModel,
      bgCustom,
      apiKey,
      apiBase,
      characterReq,
      createProject,
      createCharacter,
      updateTheme,
      generateDrafts,
      saveDraft,
      regenerateDrafts,
      generateStickers,
      regenerateSticker,
      removeBackground,
      exportZip,
    }
  },
  template: `
    <div style="max-width: 720px; margin: 32px auto; font-family: system-ui;">
      <h1>LINE 貼圖製作平台</h1>
      <Notice :message="error" type="error" />
      <Notice :message="success" type="success" />
      <Notice v-if="loading" message="處理中..." type="info" />

      <div v-if="jobStatus" style="margin: 12px 0;">
        <div>Job 狀態：{{ jobStatus }}</div>
        <div style="background:#eee; height:8px; border-radius:4px; overflow:hidden;">
          <div :style="{ width: jobProgress + '%', background:'#4ade80', height:'8px' }"></div>
        </div>
      </div>

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

        <div style="margin: 8px 0;">
          <label>預設 AI 供應商（僅顯示已驗證）</label>
          <select v-model="selectedProvider" :disabled="!verified">
            <option v-for="p in verifiedProviders" :key="p.id" :value="p.id">{{ p.name || p.id }}</option>
          </select>

          <label style="margin-left:8px;">預設模型</label>
          <select v-model="selectedModel" :disabled="!verified">
            <option v-for="m in (verifiedProviders.find(p => p.id === selectedProvider)?.models || [])" :key="m" :value="m">{{ m }}</option>
          </select>
        </div>

        <div style="margin: 8px 0;">
          <label>預設自訂模型 ID（可選）</label>
          <input v-model="customModel" placeholder="replicate model version / openai model" style="width:100%; margin:6px 0;" />
        </div>

        <div style="margin: 8px 0; padding: 8px; border:1px dashed #ddd;">
          <strong>文字生成</strong>（必填）<br/>
          <select v-model="textProvider" @change="syncDefaultModel({ value: textProvider }, { value: textModel }, { value: textCustom })" :disabled="!verified">
            <option v-for="p in verifiedProviders" :key="p.id" :value="p.id">{{ p.name || p.id }}</option>
          </select>
          <select v-model="textModel" :disabled="!verified">
            <option v-for="m in (verifiedProviders.find(p => p.id === textProvider)?.models || [])" :key="m" :value="m">{{ m }}</option>
          </select>
          <input v-model="textCustom" placeholder="自訂模型 ID" style="width:100%; margin:6px 0;" :disabled="!verified" />
        </div>

        <div style="margin: 8px 0; padding: 8px; border:1px dashed #ddd;">
          <strong>圖像生成</strong>（必填）<br/>
          <select v-model="imageProvider" @change="syncDefaultModel({ value: imageProvider }, { value: imageModel }, { value: imageCustom })" :disabled="!verified">
            <option v-for="p in verifiedProviders" :key="p.id" :value="p.id">{{ p.name || p.id }}</option>
          </select>
          <select v-model="imageModel" :disabled="!verified">
            <option v-for="m in (verifiedProviders.find(p => p.id === imageProvider)?.models || [])" :key="m" :value="m">{{ m }}</option>
          </select>
          <input v-model="imageCustom" placeholder="自訂模型 ID" style="width:100%; margin:6px 0;" :disabled="!verified" />
        </div>

        <div style="margin: 8px 0; padding: 8px; border:1px dashed #ddd;">
          <strong>去背</strong>（必填）<br/>
          <select v-model="bgProvider" @change="syncDefaultModel({ value: bgProvider }, { value: bgModel }, { value: bgCustom })" :disabled="!verified">
            <option v-for="p in verifiedProviders" :key="p.id" :value="p.id">{{ p.name || p.id }}</option>
          </select>
          <select v-model="bgModel" :disabled="!verified">
            <option v-for="m in (verifiedProviders.find(p => p.id === bgProvider)?.models || [])" :key="m" :value="m">{{ m }}</option>
          </select>
          <input v-model="bgCustom" placeholder="自訂模型 ID" style="width:100%; margin:6px 0;" :disabled="!verified" />
        </div>

        <div style="margin: 8px 0;">
          <label>API Key（只會暫存於記憶體）</label>
          <input v-model="apiKey" type="password" style="width:100%; margin:6px 0;" />
          <label>API Base（可選）</label>
          <input v-model="apiBase" placeholder="https://api.openai.com" style="width:100%; margin:6px 0;" />
          <small>完成驗證後，才能選擇下方模型</small>
          <div v-if="verifiedProviders.length === 0" style="color:#999; margin-top:4px;">尚未驗證通過的 provider/model</div>
        </div>

        <button @click="updateTheme">下一步</button>
      </section>

      <section v-else-if="step === 'DRAFTS'">
        <h2>4. 草稿生成</h2>
        <p>將根據主題與數量生成草稿。</p>
        <button @click="generateDrafts">產生草稿</button>
      </section>

      <section v-else-if="step === 'GENERATE'">
        <h2>5. 草稿確認 / 編輯</h2>
        <button @click="regenerateDrafts">重生全部草稿</button>
        <div v-for="d in drafts" :key="d.id" style="border:1px solid #eee; padding:12px; margin:12px 0;">
          <div>第 {{ d.index }} 張</div>
          <label>配字</label>
          <input v-model="d.caption" style="width:100%; margin:6px 0;" />
          <label>描述</label>
          <textarea v-model="d.imagePrompt" style="width:100%; margin:6px 0;"></textarea>
          <button @click="saveDraft(d)">保存此草稿</button>
        </div>
        <button @click="generateStickers">開始生成貼圖</button>
      </section>

      <section v-else-if="step === 'PREVIEW'">
        <h2>6. 預覽</h2>
        <ul>
          <li v-for="s in stickers" :key="s.id">
            {{ s.id }} - {{ s.status }}
            <span v-if="s.transparentUrl">(已去背)</span>
            <button @click="regenerateSticker(s.id)" style="margin-left:8px;">重生此張</button>
          </li>
        </ul>
        <div style="margin:12px 0;">
          <button @click="removeBackground">去背（保留主角完整）</button>
          <button @click="exportZip" style="margin-left:8px;">產生下載</button>
          <p style="color:#666; margin-top:6px;">提示：去背模型需保留主角完整，避免切掉頭髮/手。</p>
        </div>
        <p v-if="downloadUrl">下載連結：<a :href="downloadUrl" target="_blank">{{ downloadUrl }}</a></p>
      </section>
    </div>
  `,
}).mount('#app')
