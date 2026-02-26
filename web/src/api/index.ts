import axios, { type AxiosRequestConfig } from 'axios'

export interface ApiResponse<T> {
  success: boolean
  message?: string
  data: T
}

export interface PagedResult<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

export interface UserInfo {
  id: number
  username: string
  display_name: string
  avatar_url: string
  role: number
  status: number
  quota_total: number
  quota_used: number
  trust_level: number
  token_limit: number
  last_login_at?: number | null
  last_login_ip?: string
}

export interface TokenInfo {
  id: number
  name: string
  key_prefix: string
  status: number
  expires_at: number | null
  quota_total: number
  quota_used: number
  rate_limit_rpm: number
  allowed_models: string
  allowed_ips: string
  total_requests: number
  CreatedAt?: string
}

export interface TokenCreateResult {
  token: TokenInfo
  key: string
}

export interface RequestLogInfo {
  id: number
  user_id: number
  token_id: number
  request_ip: string
  method: string
  path: string
  model: string
  status_code: number
  duration: number
  prompt_tokens: number
  completion_tokens: number
  total_tokens: number
  error_message: string
  created_at: string
}

export interface IPBanInfo {
  id: number
  ip: string
  reason: string
  banned_by: number
  expires_at: number | null
  CreatedAt?: string
}

export interface LogStats {
  total_requests: number
  total_tokens: number
  today_requests: number
  today_tokens: number
}

export interface DashboardData {
  user: {
    username: string
    display_name: string
    role: number
    quota_total: number
    quota_used: number
  }
  stats: LogStats
  token_count: number
  global_stats?: LogStats
  user_count?: number
  trend?: Array<{ date: string; count: number }>
  model_distribution?: Array<{ model: string; count: number }>
}

export type SettingsMap = Record<string, string>

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || '',
  timeout: 30000,
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.hash = '#/login'
    }
    return Promise.reject(error.response?.data ?? error)
  }
)

const request = {
  async get<T>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const response = await api.get<ApiResponse<T>>(url, config)
    return response.data
  },
  async post<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const response = await api.post<ApiResponse<T>>(url, data, config)
    return response.data
  },
  async put<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const response = await api.put<ApiResponse<T>>(url, data, config)
    return response.data
  },
  async delete<T>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const response = await api.delete<ApiResponse<T>>(url, config)
    return response.data
  },
}

export function getErrorMessage(error: unknown, fallback: string): string {
  if (typeof error !== 'object' || error === null) {
    return fallback
  }
  const withMessage = error as { message?: unknown }
  if (typeof withMessage.message === 'string' && withMessage.message.trim() !== '') {
    return withMessage.message
  }

  const withError = error as { error?: { message?: unknown } }
  if (typeof withError.error?.message === 'string' && withError.error.message.trim() !== '') {
    return withError.error.message
  }

  return fallback
}

export default request

// Auth
export const getOAuthURL = () => `/api/oauth/linuxdo`
export const getCurrentUser = () => request.get<UserInfo>('/api/auth/user')
export const logout = () => request.post<null>('/api/auth/logout')

// Tokens
export const getTokens = () => request.get<TokenInfo[]>('/api/tokens')
export const createToken = (data: unknown) =>
  request.post<TokenCreateResult>('/api/tokens', data)
export const updateToken = (id: number, data: unknown) =>
  request.put<TokenInfo>(`/api/tokens/${id}`, data)
export const deleteToken = (id: number) => request.delete<null>(`/api/tokens/${id}`)
export const resetToken = (id: number) => request.post<TokenCreateResult>(`/api/tokens/${id}/reset`)

// Logs
export const getLogs = (params: Record<string, unknown>) =>
  request.get<PagedResult<RequestLogInfo>>('/api/logs', { params })
export const getLogStats = () => request.get<LogStats>('/api/logs/stats')

// Dashboard
export const getDashboard = () => request.get<DashboardData>('/api/dashboard')

// Admin: Users
export const getUsers = (params: Record<string, unknown>) =>
  request.get<PagedResult<UserInfo>>('/api/admin/users', { params })
export const updateUser = (id: number, data: unknown) =>
  request.put<UserInfo>(`/api/admin/users/${id}`, data)

// Admin: IP Bans
export const getIPBans = (params: Record<string, unknown>) =>
  request.get<PagedResult<IPBanInfo>>('/api/admin/ip-bans', { params })
export const createIPBan = (data: unknown) =>
  request.post<IPBanInfo>('/api/admin/ip-bans', data)
export const deleteIPBan = (id: number) => request.delete<null>(`/api/admin/ip-bans/${id}`)

// Admin: Logs
export const getAdminLogs = (params: Record<string, unknown>) =>
  request.get<PagedResult<RequestLogInfo>>('/api/admin/logs', { params })
export const getAdminLogStats = () => request.get<LogStats>('/api/admin/logs/stats')
export const cleanLogs = (days: number) => request.delete<{ deleted: number }>('/api/admin/logs', { data: { days } })

// Admin: Settings
export const getSettings = () => request.get<SettingsMap>('/api/admin/settings')
export const updateSettings = (data: SettingsMap) => request.put<null>('/api/admin/settings', data)
