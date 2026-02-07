import type {
  Character,
  CharacterCreateRequest,
  Draft,
  DraftUpdateRequest,
  ExportResponse,
  Job,
  Project,
  ProjectCreateRequest,
  ProjectUpdateRequest,
  AIConfigUpdateRequest,
  AICredentialsRequest,
  AIPipelineConfigRequest,
  Provider,
  Sticker,
  ThemeSuggestResponse,
} from './types'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const controller = new AbortController()
  const timeout = setTimeout(() => controller.abort(), 10000)
  try {
    const res = await fetch(`${API_BASE_URL}${path}`, {
      headers: { 'Content-Type': 'application/json' },
      signal: controller.signal,
      ...options,
    })
    if (!res.ok) {
      throw new Error(`API ${res.status}`)
    }
    return res.json() as Promise<T>
  } finally {
    clearTimeout(timeout)
  }
}

export const api = {
  createProject: (body: ProjectCreateRequest) =>
    request<Project>('/projects', {
      method: 'POST',
      body: JSON.stringify(body),
    }),

  listProjects: () => request<Project[]>('/projects'),

  getProject: (projectId: string) => request<Project>(`/projects/${projectId}`),

  updateProject: (projectId: string, body: ProjectUpdateRequest) =>
    request<Project>(`/projects/${projectId}`, {
      method: 'PATCH',
      body: JSON.stringify(body),
    }),

  updateAIConfig: (projectId: string, body: AIConfigUpdateRequest) =>
    request<Project>(`/projects/${projectId}/ai-config`, {
      method: 'PATCH',
      body: JSON.stringify(body),
    }),

  updateAIPipeline: (projectId: string, body: AIPipelineConfigRequest) =>
    request<Project>(`/projects/${projectId}/ai-pipeline`, {
      method: 'PATCH',
      body: JSON.stringify(body),
    }),

  listProviders: () => request<Provider[]>('/providers'),

  setAICredentials: (projectId: string, body: AICredentialsRequest) =>
    request<void>(`/projects/${projectId}/ai-credentials`, {
      method: 'POST',
      body: JSON.stringify(body),
    }),

  verifyAICredentials: (projectId: string) =>
    request<void>(`/projects/${projectId}/ai-verify`, {
      method: 'POST',
    }),

  listVerifiedProviders: (projectId: string) =>
    request<{ providers: Provider[] }>(`/projects/${projectId}/verified-providers`),

  createCharacter: (projectId: string, body: CharacterCreateRequest) =>
    request<Character>(`/projects/${projectId}/character`, {
      method: 'POST',
      body: JSON.stringify(body),
    }),

  suggestTheme: (projectId: string, seed?: string) =>
    request<ThemeSuggestResponse>(`/projects/${projectId}/theme:suggest`, {
      method: 'POST',
      body: JSON.stringify({ seed }),
    }),

  generateDrafts: (projectId: string) =>
    request<Job>(`/projects/${projectId}/drafts:generate`, {
      method: 'POST',
    }),

  listDrafts: (projectId: string) =>
    request<Draft[]>(`/projects/${projectId}/drafts`),

  updateDraft: (draftId: string, body: DraftUpdateRequest) =>
    request<Draft>(`/drafts/${draftId}`, {
      method: 'PATCH',
      body: JSON.stringify(body),
    }),

  generateStickers: (projectId: string) =>
    request<Job>(`/projects/${projectId}/stickers:generate`, {
      method: 'POST',
    }),

  listStickers: (projectId: string) =>
    request<Sticker[]>(`/projects/${projectId}/stickers`),

  regenerateSticker: (stickerId: string) =>
    request<Job>(`/stickers/${stickerId}:regenerate`, {
      method: 'POST',
    }),

  removeBackground: (projectId: string) =>
    request<Job>(`/projects/${projectId}/stickers:remove-bg`, {
      method: 'POST',
    }),

  exportZip: (projectId: string) =>
    request<ExportResponse>(`/projects/${projectId}/export`, {
      method: 'POST',
    }),

  getJob: (jobId: string) => request<Job>(`/jobs/${jobId}`),
}
