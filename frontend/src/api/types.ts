export type ProjectStatus =
  | 'DRAFT'
  | 'GENERATING_DRAFTS'
  | 'DRAFT_READY'
  | 'GENERATING_IMAGES'
  | 'IMAGES_READY'
  | 'EXPORTING'
  | 'DONE'

export type Project = {
  id: string
  title?: string
  theme?: string
  stickerCount: 8 | 16 | 24 | 40
  status: ProjectStatus
  characterId?: string
  aiProvider?: string
  aiModel?: string
}

export type ProjectCreateRequest = {
  title?: string
  stickerCount: 8 | 16 | 24 | 40
}

export type ProjectUpdateRequest = {
  theme?: string
}

export type AIConfigUpdateRequest = {
  aiProvider: string
  aiModel: string
}

export type AICredentialsRequest = {
  aiProvider: string
  apiKey: string
  apiBase?: string
}

export type CharacterCreateRequest = {
  sourceType: 'AI' | 'UPLOAD' | 'HISTORY'
  prompt?: string
  referenceImageUrl?: string
}

export type Character = {
  id: string
  sourceType: 'AI' | 'UPLOAD' | 'HISTORY'
  referenceImageUrl?: string
  status: 'READY' | 'FAILED'
}

export type Draft = {
  id: string
  projectId: string
  index: number
  caption: string
  imagePrompt: string
  status: 'DRAFT' | 'APPROVED' | 'REJECTED'
}

export type DraftUpdateRequest = {
  caption?: string
  imagePrompt?: string
}

export type Sticker = {
  id: string
  projectId: string
  draftId: string
  imageUrl: string
  transparentUrl?: string
  status: 'PENDING' | 'GENERATING' | 'READY' | 'FAILED'
}

export type Job = {
  id: string
  type: 'GENERATE_DRAFT' | 'GENERATE_IMAGE' | 'REMOVE_BG'
  status: 'QUEUED' | 'RUNNING' | 'SUCCESS' | 'FAILED'
  progress?: number
  errorMessage?: string
}

export type ThemeSuggestResponse = {
  suggestions: string[]
}

export type Provider = {
  id: string
  name: string
  models: string[]
}

export type ExportResponse = {
  downloadUrl: string
}
