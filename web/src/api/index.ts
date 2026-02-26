import axios from 'axios'

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
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.hash = '#/login'
    }
    return Promise.reject(error.response?.data || error)
  }
)

export default api

// Auth
export const getOAuthURL = () => `/api/oauth/linuxdo`
export const getCurrentUser = () => api.get('/api/auth/user')
export const logout = () => api.post('/api/auth/logout')

// Tokens
export const getTokens = () => api.get('/api/tokens')
export const createToken = (data: any) => api.post('/api/tokens', data)
export const updateToken = (id: number, data: any) => api.put(`/api/tokens/${id}`, data)
export const deleteToken = (id: number) => api.delete(`/api/tokens/${id}`)
export const resetToken = (id: number) => api.post(`/api/tokens/${id}/reset`)

// Logs
export const getLogs = (params: any) => api.get('/api/logs', { params })
export const getLogStats = () => api.get('/api/logs/stats')

// Dashboard
export const getDashboard = () => api.get('/api/dashboard')

// Admin: Users
export const getUsers = (params: any) => api.get('/api/admin/users', { params })
export const updateUser = (id: number, data: any) => api.put(`/api/admin/users/${id}`, data)

// Admin: IP Bans
export const getIPBans = (params: any) => api.get('/api/admin/ip-bans', { params })
export const createIPBan = (data: any) => api.post('/api/admin/ip-bans', data)
export const deleteIPBan = (id: number) => api.delete(`/api/admin/ip-bans/${id}`)

// Admin: Logs
export const getAdminLogs = (params: any) => api.get('/api/admin/logs', { params })
export const getAdminLogStats = () => api.get('/api/admin/logs/stats')
export const cleanLogs = (days: number) => api.delete('/api/admin/logs', { data: { days } })

// Admin: Settings
export const getSettings = () => api.get('/api/admin/settings')
export const updateSettings = (data: any) => api.put('/api/admin/settings', data)
