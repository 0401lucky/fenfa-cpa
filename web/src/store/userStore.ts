import { create } from 'zustand'
import { getCurrentUser } from '../api'

interface User {
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
}

interface UserState {
  token: string | null
  user: User | null
  loading: boolean
  setToken: (token: string) => void
  fetchUser: () => Promise<void>
  logout: () => void
  isAdmin: () => boolean
}

export const useUserStore = create<UserState>((set, get) => ({
  token: localStorage.getItem('token'),
  user: null,
  loading: false,

  setToken: (token: string) => {
    localStorage.setItem('token', token)
    set({ token })
  },

  fetchUser: async () => {
    set({ loading: true })
    try {
      const res: any = await getCurrentUser()
      set({ user: res.data, loading: false })
    } catch {
      set({ user: null, loading: false })
      localStorage.removeItem('token')
      set({ token: null })
    }
  },

  logout: () => {
    localStorage.removeItem('token')
    set({ token: null, user: null })
  },

  isAdmin: () => {
    const user = get().user
    return user ? user.role >= 10 : false
  },
}))
