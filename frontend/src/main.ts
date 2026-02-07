import { createApp, onMounted, ref, watch } from 'vue'
import Notice from './components/Notice.vue'
import { api } from './api/client'
import './style.css'
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
  components: { Notice },
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
    const jobError = ref('')
    const lastAction = ref('')
    const exportWarnings = ref<string[]>([])
    const gridCols = ref(4)

    const sectionClass = 'mt-3 rounded-xl border border-slate-200 bg-white p-4'
    const buttonClass = 'rounded-md border border-slate-300 bg-slate-50 px-3 py-2 text-slate-800'
    const primaryButtonClass = 'rounded-md border border-blue-600 bg-blue-600 px-3 py-2 text-white'
    const inputClass = 'w-full rounded-md border border-slate-300 bg-white px-2 py-1.5'
    const selectClass = 'rounded-md border border-slate-300 bg-white px-2 py-1.5'
    const textareaClass = 'w-full rounded-md border border-slate-300 bg-white px-2 py-1.5 min-h-[72px]'
    const h2Class = 'text-lg font-semibold text-slate-900'
    const labelClass = 'text-xs text-slate-500'

    const locales = {
      'zh-TW': { name: '繁體中文' },
      'zh-CN': { name: '简体中文' },
      en: { name: 'English' },
      ja: { name: '日本語' },
      ko: { name: '한국어' },
    }
    const localeOptions = Object.keys(locales)
    const normalizeLocale = (val: string) => {
      if (val.startsWith('zh-TW') || val.startsWith('zh-HK')) return 'zh-TW'
      if (val.startsWith('zh')) return 'zh-CN'
      if (val.startsWith('ja')) return 'ja'
      if (val.startsWith('ko')) return 'ko'
      return 'en'
    }
    const currentLocale = ref(normalizeLocale(navigator.language || 'en'))

    const t = (key: string) => {
      const dict: Record<string, Record<string, string>> = {
        'app.title': {
          'zh-TW': 'LINE 貼圖製作平台',
          'zh-CN': 'LINE 表情制作平台',
          en: 'LINE Sticker Studio',
          ja: 'LINE スタンプ作成',
          ko: 'LINE 스티커 제작',
        },
        'step.create': {
          'zh-TW': '1. 建立專案',
          'zh-CN': '1. 创建项目',
          en: '1. Create Project',
          ja: '1. プロジェクト作成',
          ko: '1. 프로젝트 생성',
        },
        'step.character': {
          'zh-TW': '2. 角色設定',
          'zh-CN': '2. 角色设定',
          en: '2. Character',
          ja: '2. キャラクター',
          ko: '2. 캐릭터',
        },
        'step.theme': {
          'zh-TW': '3. 主題與焦點',
          'zh-CN': '3. 主题与焦点',
          en: '3. Theme & Focus',
          ja: '3. テーマ',
          ko: '3. 테마',
        },
        'step.drafts': {
          'zh-TW': '4. 產生草稿',
          'zh-CN': '4. 生成草稿',
          en: '4. Drafts',
          ja: '4. 下書き',
          ko: '4. 초안',
        },
        'step.generate': {
          'zh-TW': '5. 草稿確認 / 編輯',
          'zh-CN': '5. 草稿确认 / 编辑',
          en: '5. Review / Edit',
          ja: '5. 確認 / 編集',
          ko: '5. 확인 / 편집',
        },
        'step.preview': {
          'zh-TW': '6. 預覽',
          'zh-CN': '6. 预览',
          en: '6. Preview',
          ja: '6. プレビュー',
          ko: '6. 미리보기',
        },
        'label.projectName': {
          'zh-TW': '專案名稱',
          'zh-CN': '项目名称',
          en: 'Project Name',
          ja: 'プロジェクト名',
          ko: '프로젝트명',
        },
        'label.stickerCount': {
          'zh-TW': '貼圖數量',
          'zh-CN': '表情数量',
          en: 'Sticker Count',
          ja: '枚数',
          ko: '개수',
        },
        'label.language': {
          'zh-TW': '語言',
          'zh-CN': '语言',
          en: 'Language',
          ja: '言語',
          ko: '언어',
        },
        'button.next': {
          'zh-TW': '下一步',
          'zh-CN': '下一步',
          en: 'Next',
          ja: '次へ',
          ko: '다음',
        },
        'button.retry': {
          'zh-TW': '重試',
          'zh-CN': '重试',
          en: 'Retry',
          ja: '再試行',
          ko: '다시 시도',
        },
        'button.back': {
          'zh-TW': '返回上一步',
          'zh-CN': '返回上一步',
          en: 'Back',
          ja: '戻る',
          ko: '뒤로',
        },
        'label.source': {
          'zh-TW': '來源',
          'zh-CN': '来源',
          en: 'Source',
          ja: 'ソース',
          ko: '소스',
        },
        'label.character': {
          'zh-TW': '角色描述',
          'zh-CN': '角色描述',
          en: 'Character Prompt',
          ja: 'キャラクター',
          ko: '캐릭터 설명',
        },
        'label.referenceUrl': {
          'zh-TW': 'Reference 圖片 URL',
          'zh-CN': 'Reference 图片 URL',
          en: 'Reference Image URL',
          ja: '参照画像URL',
          ko: '레퍼런스 이미지 URL',
        },
        'label.theme': {
          'zh-TW': '主題',
          'zh-CN': '主题',
          en: 'Theme',
          ja: 'テーマ',
          ko: '테마',
        },
        'label.defaultProvider': {
          'zh-TW': '預設 AI 供應商（僅顯示已驗證）',
          'zh-CN': '默认 AI 供应商（仅显示已验证）',
          en: 'Default Provider (verified only)',
          ja: '既定プロバイダ（検証済みのみ）',
          ko: '기본 공급자(검증됨만)',
        },
        'label.defaultModel': {
          'zh-TW': '預設模型',
          'zh-CN': '默认模型',
          en: 'Default Model',
          ja: '既定モデル',
          ko: '기본 모델',
        },
        'label.customModel': {
          'zh-TW': '預設自訂模型 ID（可選）',
          'zh-CN': '默认自定义模型 ID（可选）',
          en: 'Custom Model ID (optional)',
          ja: 'カスタムモデルID（任意）',
          ko: '사용자 모델 ID(선택)',
        },
        'label.textGen': {
          'zh-TW': '文字生成',
          'zh-CN': '文本生成',
          en: 'Text Generation',
          ja: 'テキスト生成',
          ko: '텍스트 생성',
        },
        'label.imageGen': {
          'zh-TW': '圖像生成',
          'zh-CN': '图像生成',
          en: 'Image Generation',
          ja: '画像生成',
          ko: '이미지 생성',
        },
        'label.bgRemove': {
          'zh-TW': '去背',
          'zh-CN': '去背',
          en: 'Background Removal',
          ja: '背景除去',
          ko: '배경 제거',
        },
        'label.apiKey': {
          'zh-TW': 'API Key（只會暫存於記憶體）',
          'zh-CN': 'API Key（仅暂存于内存）',
          en: 'API Key (in-memory only)',
          ja: 'APIキー（メモリのみ）',
          ko: 'API 키(메모리만)',
        },
        'label.apiBase': {
          'zh-TW': 'API Base（可選）',
          'zh-CN': 'API Base（可选）',
          en: 'API Base (optional)',
          ja: 'API Base（任意）',
          ko: 'API Base(선택)',
        },
        'helper.verifyFirst': {
          'zh-TW': '完成驗證後，才能選擇下方模型',
          'zh-CN': '完成验证后才能选择下方模型',
          en: 'Verify first to enable model selection',
          ja: '検証後にモデル選択が可能です',
          ko: '검증 후 모델을 선택하세요',
        },
        'helper.noVerified': {
          'zh-TW': '尚未驗證通過的 provider/model',
          'zh-CN': '尚未验证通过的 provider/model',
          en: 'No verified providers/models yet',
          ja: '検証済みのプロバイダ/モデルがありません',
          ko: '검증된 공급자/모델이 없습니다',
        },
        'label.caption': {
          'zh-TW': '配字',
          'zh-CN': '配字',
          en: 'Caption',
          ja: 'キャプション',
          ko: '문구',
        },
        'label.prompt': {
          'zh-TW': '描述',
          'zh-CN': '描述',
          en: 'Description',
          ja: '説明',
          ko: '설명',
        },
        'label.grid': {
          'zh-TW': '格數：',
          'zh-CN': '格数：',
          en: 'Grid:',
          ja: 'グリッド：',
          ko: '그리드:',
        },
        'button.generateDrafts': {
          'zh-TW': '產生草稿',
          'zh-CN': '生成草稿',
          en: 'Generate Drafts',
          ja: '下書きを生成',
          ko: '초안 생성',
        },
        'button.regenerateDrafts': {
          'zh-TW': '重生全部草稿',
          'zh-CN': '全部重生草稿',
          en: 'Regenerate Drafts',
          ja: '下書きを再生成',
          ko: '초안 재생성',
        },
        'button.saveDraft': {
          'zh-TW': '保存此草稿',
          'zh-CN': '保存草稿',
          en: 'Save Draft',
          ja: '下書きを保存',
          ko: '초안 저장',
        },
        'button.generateStickers': {
          'zh-TW': '開始生成貼圖',
          'zh-CN': '开始生成贴图',
          en: 'Generate Stickers',
          ja: 'スタンプ生成',
          ko: '스티커 생성',
        },
        'button.regenerateOne': {
          'zh-TW': '重生此張',
          'zh-CN': '重生此张',
          en: 'Regenerate',
          ja: '再生成',
          ko: '재생성',
        },
        'button.removeBg': {
          'zh-TW': '去背（保留主角完整）',
          'zh-CN': '去背（保留主体完整）',
          en: 'Remove BG (keep subject)',
          ja: '背景除去（主役を保持）',
          ko: '배경 제거(대상 유지)',
        },
        'button.export': {
          'zh-TW': '產生下載',
          'zh-CN': '生成下载',
          en: 'Export',
          ja: '書き出し',
          ko: '내보내기',
        },
        'hint.removeBg': {
          'zh-TW': '提示：去背模型需保留主角完整，避免切掉頭髮/手。',
          'zh-CN': '提示：去背模型需保留主体完整，避免切掉头发/手。',
          en: 'Tip: keep subject intact when removing background.',
          ja: 'ヒント：主役を切らないように。',
          ko: '팁: 대상이 잘리지 않도록.',
        },
        'status.bgDone': {
          'zh-TW': '去背完成',
          'zh-CN': '去背完成',
          en: 'Background removed',
          ja: '背景除去完了',
          ko: '배경 제거 완료',
        },
        'status.exportDone': {
          'zh-TW': '輸出完成',
          'zh-CN': '输出完成',
          en: 'Export complete',
          ja: '書き出し完了',
          ko: '내보내기 완료',
        },
        'label.download': {
          'zh-TW': '下載連結',
          'zh-CN': '下载链接',
          en: 'Download',
          ja: 'ダウンロード',
          ko: '다운로드',
        },
        'label.warnTitle': {
          'zh-TW': '檢核警告',
          'zh-CN': '检核警告',
          en: 'Validation Warnings',
          ja: '検証警告',
          ko: '검증 경고',
        },
        'label.main': {
          'zh-TW': '主圖',
          'zh-CN': '主图',
          en: 'Main',
          ja: 'メイン',
          ko: '메인',
        },
        'label.tab': {
          'zh-TW': 'Tab',
          'zh-CN': 'Tab',
          en: 'Tab',
          ja: 'タブ',
          ko: '탭',
        },
        'error.requiredKey': {
          'zh-TW': '請先填入 API Key',
          'zh-CN': '请先填写 API Key',
          en: 'Please enter API Key',
          ja: 'APIキーを入力してください',
          ko: 'API 키를 입력하세요',
        },
        'error.textConfig': {
          'zh-TW': '請完成「文字生成」的供應商與模型（或自訂模型 ID）',
          'zh-CN': '请完成“文本生成”的供应商与模型（或自定义模型 ID）',
          en: 'Complete text provider/model (or custom model ID).',
          ja: 'テキストのプロバイダ/モデルを設定してください',
          ko: '텍스트 공급자/모델을 설정하세요',
        },
        'error.imageConfig': {
          'zh-TW': '請完成「圖像生成」的供應商與模型（或自訂模型 ID）',
          'zh-CN': '请完成“图像生成”的供应商与模型（或自定义模型 ID）',
          en: 'Complete image provider/model (or custom model ID).',
          ja: '画像のプロバイダ/モデルを設定してください',
          ko: '이미지 공급자/모델을 설정하세요',
        },
        'error.bgConfig': {
          'zh-TW': '請完成「去背」的供應商與模型（或自訂模型 ID）',
          'zh-CN': '请完成“去背”的供应商与模型（或自定义模型 ID）',
          en: 'Complete background provider/model (or custom model ID).',
          ja: '背景除去のプロバイダ/モデルを設定してください',
          ko: '배경 제거 공급자/모델을 설정하세요',
        },
        'error.generic': {
          'zh-TW': '操作失敗，請稍後再試',
          'zh-CN': '操作失败，请稍后再试',
          en: 'Operation failed. Please try again.',
          ja: '失敗しました。再試行してください',
          ko: '작업 실패. 다시 시도하세요',
        },
        'status.verifyOk': {
          'zh-TW': 'API Key 驗證成功',
          'zh-CN': 'API Key 验证成功',
          en: 'API Key verified',
          ja: 'APIキー検証成功',
          ko: 'API 키 검증 성공',
        },
        'status.exportOk': {
          'zh-TW': '輸出成功，已完成規格檢核',
          'zh-CN': '输出成功，已完成规格检核',
          en: 'Export success. Validation done.',
          ja: '書き出し成功（検証完了）',
          ko: '내보내기 성공(검증 완료)',
        },
        'status.exportWarn': {
          'zh-TW': '輸出完成，但有部分貼圖未通過檢核',
          'zh-CN': '输出完成，但部分贴图未通过检核',
          en: 'Export complete, some stickers failed validation.',
          ja: '一部が検証に失敗しました',
          ko: '일부 스티커가 검증 실패',
        },
        'label.recommendBg': {
          'zh-TW': '推薦去背模型',
          'zh-CN': '推荐去背模型',
          en: 'Recommended BG model',
          ja: 'おすすめ背景モデル',
          ko: '추천 배경 모델',
        },
        'status.loading': {
          'zh-TW': '處理中...',
          'zh-CN': '处理中...',
          en: 'Loading...',
          ja: '処理中...',
          ko: '처리 중...',
        },
        'label.required': {
          'zh-TW': '必填',
          'zh-CN': '必填',
          en: 'Required',
          ja: '必須',
          ko: '필수',
        },
        'option.source.ai': {
          'zh-TW': 'AI 生成',
          'zh-CN': 'AI 生成',
          en: 'AI',
          ja: 'AI生成',
          ko: 'AI 생성',
        },
        'option.source.upload': {
          'zh-TW': '上傳圖片',
          'zh-CN': '上传图片',
          en: 'Upload',
          ja: 'アップロード',
          ko: '업로드',
        },
        'option.source.history': {
          'zh-TW': '歷史角色',
          'zh-CN': '历史角色',
          en: 'History',
          ja: '履歴',
          ko: '기록',
        },
        'placeholder.customModelHint': {
          'zh-TW': 'replicate model version / openai model',
          'zh-CN': 'replicate model version / openai model',
          en: 'replicate model version / openai model',
          ja: 'replicate model version / openai model',
          ko: 'replicate model version / openai model',
        },
        'placeholder.customModelId': {
          'zh-TW': '自訂模型 ID',
          'zh-CN': '自定义模型 ID',
          en: 'Custom model ID',
          ja: 'カスタムモデルID',
          ko: '사용자 모델 ID',
        },
        'desc.drafts': {
          'zh-TW': '將根據主題與數量生成草稿。',
          'zh-CN': '将根据主题与数量生成草稿。',
          en: 'Drafts will be generated from your theme and count.',
          ja: 'テーマと枚数から下書きを生成します。',
          ko: '테마와 개수에 따라 초안을 생성합니다.',
        },
        'label.draftIndex': {
          'zh-TW': '第 {n} 張',
          'zh-CN': '第 {n} 张',
          en: 'Sticker {n}',
          ja: '{n}枚目',
          ko: '{n}번',
        },
        'job.title': {
          'zh-TW': 'Job 狀態',
          'zh-CN': '任务状态',
          en: 'Job Status',
          ja: 'ジョブ状態',
          ko: '작업 상태',
        },
        'job.state.RUNNING': {
          'zh-TW': '進行中',
          'zh-CN': '进行中',
          en: 'Running',
          ja: '実行中',
          ko: '진행 중',
        },
        'job.state.SUCCESS': {
          'zh-TW': '成功',
          'zh-CN': '成功',
          en: 'Success',
          ja: '成功',
          ko: '성공',
        },
        'job.state.FAILED': {
          'zh-TW': '失敗',
          'zh-CN': '失败',
          en: 'Failed',
          ja: '失敗',
          ko: '실패',
        },
        'grid.4': {
          'zh-TW': '4（40 張）',
          'zh-CN': '4（40 张）',
          en: '4 (40)',
          ja: '4（40枚）',
          ko: '4 (40장)',
        },
        'grid.5': {
          'zh-TW': '5（24 張）',
          'zh-CN': '5（24 张）',
          en: '5 (24)',
          ja: '5（24枚）',
          ko: '5 (24장)',
        },
        'grid.6': {
          'zh-TW': '6（16 張）',
          'zh-CN': '6（16 张）',
          en: '6 (16)',
          ja: '6（16枚）',
          ko: '6 (16장)',
        },
        'grid.8': {
          'zh-TW': '8（8 張）',
          'zh-CN': '8（8 张）',
          en: '8 (8)',
          ja: '8（8枚）',
          ko: '8 (8장)',
        },
      }
      const loc = currentLocale.value
      return dict[key]?.[loc] || dict[key]?.en || key
    }

    const tf = (key: string, params: Record<string, string | number>) => {
      let s = t(key)
      for (const [k, v] of Object.entries(params)) {
        s = s.replaceAll(`{${k}}`, String(v))
      }
      return s
    }

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
    const prev = (s: Step) => (step.value = s)

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

    const bgRecommendMap: Record<string, string> = {
      replicate: 'cjwbw/rembg:fb8af171cfa1616ddcf1242c093f9c46bcada5ad4cf6f2fbe8b81b330ec5c003',
    }

    const syncDefaultModel = (providerRef: any, modelRef: any, customRef: any) => {
      const p = verifiedProviders.value.find((x) => x.id === providerRef.value)
      if (p && p.models.length > 0) {
        modelRef.value = p.models[0]
      }
      customRef.value = ''
    }

    const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms))

    const pollJob = async (jobId: string) => {
      jobStatus.value = 'RUNNING'
      jobProgress.value = 0
      jobError.value = ''
      for (let i = 0; i < 20; i++) {
        const job = await api.getJob(jobId)
        jobStatus.value = job.status
        jobProgress.value = job.progress ?? 0
        if (job.status === 'FAILED') {
          jobError.value = job.errorMessage || '任務失敗，請稍後重試'
          return job.status
        }
        if (job.status === 'SUCCESS') {
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
        error.value = t('error.generic')
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
          error.value = t('error.requiredKey')
          return
        }
        const textModelValue = textCustom.value || textModel.value
        const imageModelValue = imageCustom.value || imageModel.value
        const bgModelValue = bgCustom.value || bgModel.value

        if (!textProvider.value || !textModelValue) {
          error.value = t('error.textConfig')
          return
        }
        if (!imageProvider.value || !imageModelValue) {
          error.value = t('error.imageConfig')
          return
        }
        if (!bgProvider.value || !bgModelValue) {
          error.value = t('error.bgConfig')
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
          bgModel.value = bgRecommendMap[p0.id] || (p0.models[0] || '')
        }
        project.value = await api.updateProject(project.value.id, {
          theme: theme.value,
        })
        success.value = t('status.verifyOk')
        next('DRAFTS')
      })

    const generateDrafts = () =>
      run(async () => {
        if (!project.value) return
        lastAction.value = 'generateDrafts'
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
        lastAction.value = 'regenerateDrafts'
        const job = await api.generateDrafts(project.value.id)
        await pollJob(job.id)
        drafts.value = await api.listDrafts(project.value.id)
      })

    const generateStickers = () =>
      run(async () => {
        if (!project.value) return
        lastAction.value = 'generateStickers'
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
        lastAction.value = 'removeBackground'
        const job = await api.removeBackground(project.value.id)
        await pollJob(job.id)
        stickers.value = await api.listStickers(project.value.id)
      })

    const exportZip = () =>
      run(async () => {
        if (!project.value) return
        lastAction.value = 'exportZip'
        const res = await api.exportZip(project.value.id)
        downloadUrl.value = res.downloadUrl
        exportWarnings.value = res.warnings || []
        if (exportWarnings.value.length > 0) {
          error.value = t('status.exportWarn')
        } else {
          success.value = t('status.exportOk')
        }
      })

    const retryLastAction = () => {
      switch (lastAction.value) {
        case 'generateDrafts':
        case 'regenerateDrafts':
          generateDrafts()
          break
        case 'generateStickers':
          generateStickers()
          break
        case 'removeBackground':
          removeBackground()
          break
        case 'exportZip':
          exportZip()
          break
      }
    }

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
      jobError,
      retryLastAction,
      exportWarnings,
      sectionClass,
      buttonClass,
      primaryButtonClass,
      inputClass,
      selectClass,
      textareaClass,
      h2Class,
      labelClass,
      locales,
      localeOptions,
      currentLocale,
      t,
      tf,
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
      prev,
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
    <div class="mx-auto my-8 max-w-5xl rounded-xl bg-slate-50 p-4 text-slate-900">
      <div class="flex items-center justify-between gap-3">
        <h1 class="text-xl font-semibold">{{ t('app.title') }}</h1>
        <div class="flex items-center gap-2">
          <label :class="labelClass">{{ t('label.language') }}</label>
          <select v-model="currentLocale" :class="selectClass">
            <option v-for="code in localeOptions" :key="code" :value="code">{{ locales[code].name }}</option>
          </select>
        </div>
      </div>
      <Notice :message="error" type="error" />
      <Notice :message="success" type="success" />
      <Notice v-if="loading" :message="t('status.loading')" type="info" />

      <div v-if="jobStatus" class="my-3 rounded-lg border border-slate-200 bg-white p-3">
        <div class="flex items-center justify-between">
          <div class="font-semibold">Job {{ jobStatus }}</div>
          <div class="text-xs text-slate-500">{{ jobProgress }}%</div>
        </div>
        <div class="mt-2 h-2 overflow-hidden rounded bg-slate-100">
          <div :style="{ width: jobProgress + '%' }" class="h-2 bg-green-400"></div>
        </div>
        <div v-if="jobError" class="mt-2 text-red-700">
          {{ jobError }}
          <button @click="retryLastAction" :class="buttonClass" class="ml-2">{{ t('button.retry') }}</button>
        </div>
      </div>

      <section v-if="step === 'CREATE_PROJECT'" :class="sectionClass">
        <h2 :class="h2Class">{{ t('step.create') }}</h2>
        <label :class="labelClass">{{ t('label.projectName') }}</label>
        <input v-model="title" :class="inputClass" />

        <label :class="labelClass">{{ t('label.stickerCount') }}</label>
        <select v-model.number="stickerCount" :class="selectClass">
          <option :value="8">8</option>
          <option :value="16">16</option>
          <option :value="24">24</option>
          <option :value="40">40</option>
        </select>
        <button @click="createProject" :class="primaryButtonClass" class="ml-3">{{ t('button.next') }}</button>
      </section>

      <section v-else-if="step === 'CHARACTER'" :class="sectionClass">
        <h2 :class="h2Class">{{ t('step.character') }}</h2>
        <div class="mb-2 flex items-center gap-2">
          <button @click="prev('CREATE_PROJECT')" :class="buttonClass">{{ t('button.back') }}</button>
          <button @click="createCharacter" :class="primaryButtonClass">{{ t('button.next') }}</button>
        </div>
        <label :class="labelClass">{{ t('label.source') }}</label>
        <select v-model="characterReq.sourceType" :class="selectClass">
          <option value="AI">{{ t('option.source.ai') }}</option>
          <option value="UPLOAD">{{ t('option.source.upload') }}</option>
          <option value="HISTORY">{{ t('option.source.history') }}</option>
        </select>

        <div v-if="characterReq.sourceType === 'AI'">
          <label :class="labelClass">{{ t('label.character') }}</label>
          <input v-model="characterReq.prompt" :class="inputClass" />
        </div>
        <div v-else>
          <label :class="labelClass">{{ t('label.referenceUrl') }}</label>
          <input v-model="characterReq.referenceImageUrl" :class="inputClass" />
        </div>
      </section>

      <section v-else-if="step === 'THEME'" :class="sectionClass">
        <h2 :class="h2Class">{{ t('step.theme') }}</h2>
        <label :class="labelClass">{{ t('label.theme') }}</label>
        <input v-model="theme" :class="inputClass" />

        <div class="my-2">
          <label :class="labelClass">{{ t('label.defaultProvider') }}</label>
          <select v-model="selectedProvider" :disabled="!verified" :class="selectClass">
            <option v-for="p in verifiedProviders" :key="p.id" :value="p.id">{{ p.name || p.id }}</option>
          </select>

          <label :class="labelClass" class="ml-2">{{ t('label.defaultModel') }}</label>
          <select v-model="selectedModel" :disabled="!verified" :class="selectClass">
            <option v-for="m in (verifiedProviders.find(p => p.id === selectedProvider)?.models || [])" :key="m" :value="m">{{ m }}</option>
          </select>
        </div>

        <div class="my-2">
          <label :class="labelClass">{{ t('label.customModel') }}</label>
          <input v-model="customModel" :placeholder="t('placeholder.customModelHint')" :class="inputClass" />
        </div>

        <div class="my-2 rounded-md border border-dashed border-slate-200 p-2">
          <strong>{{ t('label.textGen') }}</strong>（{{ t('label.required') }}）<br/>
          <select v-model="textProvider" @change="syncDefaultModel({ value: textProvider }, { value: textModel }, { value: textCustom })" :disabled="!verified" :class="selectClass">
            <option v-for="p in verifiedProviders" :key="p.id" :value="p.id">{{ p.name || p.id }}</option>
          </select>
          <select v-model="textModel" :disabled="!verified" :class="selectClass">
            <option v-for="m in (verifiedProviders.find(p => p.id === textProvider)?.models || [])" :key="m" :value="m">{{ m }}</option>
          </select>
          <input v-model="textCustom" :placeholder="t('placeholder.customModelId')" :class="inputClass" :disabled="!verified" />
        </div>

        <div class="my-2 rounded-md border border-dashed border-slate-200 p-2">
          <strong>{{ t('label.imageGen') }}</strong>（{{ t('label.required') }}）<br/>
          <select v-model="imageProvider" @change="syncDefaultModel({ value: imageProvider }, { value: imageModel }, { value: imageCustom })" :disabled="!verified" :class="selectClass">
            <option v-for="p in verifiedProviders" :key="p.id" :value="p.id">{{ p.name || p.id }}</option>
          </select>
          <select v-model="imageModel" :disabled="!verified" :class="selectClass">
            <option v-for="m in (verifiedProviders.find(p => p.id === imageProvider)?.models || [])" :key="m" :value="m">{{ m }}</option>
          </select>
          <input v-model="imageCustom" :placeholder="t('placeholder.customModelId')" :class="inputClass" :disabled="!verified" />
        </div>

        <div class="my-2 rounded-md border border-dashed border-slate-200 p-2">
          <strong>{{ t('label.bgRemove') }}</strong>（{{ t('label.required') }}）<br/>
          <select v-model="bgProvider" @change="syncDefaultModel({ value: bgProvider }, { value: bgModel }, { value: bgCustom })" :disabled="!verified" :class="selectClass">
            <option v-for="p in verifiedProviders" :key="p.id" :value="p.id">{{ p.name || p.id }}</option>
          </select>
          <select v-model="bgModel" :disabled="!verified" :class="selectClass">
            <option v-for="m in (verifiedProviders.find(p => p.id === bgProvider)?.models || [])" :key="m" :value="m">{{ m }}</option>
          </select>
          <input v-model="bgCustom" :placeholder="t('placeholder.customModelId')" :class="inputClass" :disabled="!verified" />
          <div class="mt-1 text-sm text-slate-500" v-if="bgRecommendMap[bgProvider]">{{ t('label.recommendBg') }}：{{ bgRecommendMap[bgProvider] }}</div>
        </div>

        <div class="my-2">
          <label :class="labelClass">{{ t('label.apiKey') }}</label>
          <input v-model="apiKey" type="password" :class="inputClass" />
          <label :class="labelClass">{{ t('label.apiBase') }}</label>
          <input v-model="apiBase" placeholder="https://api.openai.com" :class="inputClass" />
          <small class="text-slate-500">{{ t('helper.verifyFirst') }}</small>
          <div v-if="verifiedProviders.length === 0" class="mt-1 text-sm text-slate-400">{{ t('helper.noVerified') }}</div>
        </div>

        <div class="mt-2 flex items-center gap-2">
          <button @click="prev('CHARACTER')" :class="buttonClass">{{ t('button.back') }}</button>
          <button @click="updateTheme" :class="primaryButtonClass">{{ t('button.next') }}</button>
        </div>
      </section>

      <section v-else-if="step === 'DRAFTS'" :class="sectionClass">
        <h2 :class="h2Class">{{ t('step.drafts') }}</h2>
        <p class="text-sm text-slate-600">{{ t('desc.drafts') }}</p>
        <div class="mt-2 flex items-center gap-2">
          <button @click="prev('THEME')" :class="buttonClass">{{ t('button.back') }}</button>
          <button @click="generateDrafts" :class="primaryButtonClass">{{ t('button.next') }}</button>
        </div>
      </section>

      <section v-else-if="step === 'GENERATE'" :class="sectionClass">
        <h2 :class="h2Class">{{ t('step.generate') }}</h2>
        <button @click="regenerateDrafts" :class="buttonClass">{{ t('button.regenerateDrafts') }}</button>
        <div v-for="d in drafts" :key="d.id" class="my-3 rounded-lg border border-slate-200 bg-white p-3">
          <div class="text-sm text-slate-600">{{ tf('label.draftIndex', { n: d.index }) }}</div>
          <label :class="labelClass">{{ t('label.caption') }}</label>
          <input v-model="d.caption" :class="inputClass" />
          <label :class="labelClass">{{ t('label.prompt') }}</label>
          <textarea v-model="d.imagePrompt" :class="textareaClass"></textarea>
          <button @click="saveDraft(d)" :class="buttonClass">{{ t('button.saveDraft') }}</button>
        </div>
        <div class="mt-2 flex items-center gap-2">
          <button @click="prev('DRAFTS')" :class="buttonClass">{{ t('button.back') }}</button>
          <button @click="generateStickers" :class="primaryButtonClass">{{ t('button.generateStickers') }}</button>
        </div>
      </section>

      <section v-else-if="step === 'PREVIEW'" :class="sectionClass">
        <h2 :class="h2Class">{{ t('step.preview') }}</h2>
        <div class="my-2">
          <label :class="labelClass">{{ t('label.grid') }}</label>
          <select v-model="gridCols" :class="selectClass">
            <option :value="4">{{ t('grid.4') }}</option>
            <option :value="5">{{ t('grid.5') }}</option>
            <option :value="6">{{ t('grid.6') }}</option>
            <option :value="8">{{ t('grid.8') }}</option>
          </select>
        </div>
        <div class="grid gap-2" :style="{ gridTemplateColumns: 'repeat(' + gridCols + ', minmax(0, 1fr))' }">
          <div v-for="s in stickers" :key="s.id" class="rounded border border-slate-200 bg-white p-2 text-center">
            <img :src="s.transparentUrl || s.imageUrl" class="mx-auto w-full max-w-[140px]" />
            <div class="mt-1 text-xs text-slate-500">#{{ s.id.slice(0,6) }} <span v-if="s.transparentUrl">({{ t('status.bgDone') }})</span></div>
            <button @click="regenerateSticker(s.id)" :class="buttonClass" class="mt-2">{{ t('button.regenerateOne') }}</button>
          </div>
        </div>
        <div class="my-3">
          <button @click="removeBackground" :class="buttonClass">{{ t('button.removeBg') }}</button>
          <button @click="exportZip" :class="primaryButtonClass" class="ml-2" :disabled="stickers.length === 0">{{ t('button.export') }}</button>
          <p class="mt-2 text-sm text-slate-500">{{ t('hint.removeBg') }}</p>
          <p class="mt-1 text-sm text-green-600" v-if="stickers.some(s => s.transparentUrl)">{{ t('status.bgDone') }}</p>
        </div>
        <div v-if="downloadUrl" class="mt-2 rounded border border-green-200 bg-green-50 p-2">
          <div class="font-semibold text-green-800">{{ t('status.exportDone') }}</div>
          <div class="text-sm">{{ t('label.download') }}：<a :href="downloadUrl" target="_blank" class="text-blue-600 underline">{{ downloadUrl }}</a></div>
        </div>
        <div v-if="exportWarnings.length" class="mt-2 rounded border border-orange-200 bg-orange-50 p-2">
          <div class="font-semibold text-orange-800">{{ t('label.warnTitle') }}</div>
          <ul class="mt-1 list-disc pl-5">
            <li v-for="w in exportWarnings" :key="w" class="text-sm text-orange-800">{{ w }}</li>
          </ul>
        </div>
        <div v-if="stickers.length" class="mt-3 flex items-center gap-3">
          <div>
            <div class="text-xs text-slate-500">{{ t('label.main') }} 240×240</div>
            <img :src="stickers[0].transparentUrl || stickers[0].imageUrl" class="h-[120px] w-[120px] border border-slate-200" />
          </div>
          <div>
            <div class="text-xs text-slate-500">{{ t('label.tab') }} 96×74</div>
            <img :src="stickers[0].transparentUrl || stickers[0].imageUrl" class="h-[74px] w-[96px] border border-slate-200" />
          </div>
        </div>
      </section>
    </div>
  `,
}).mount('#app')
