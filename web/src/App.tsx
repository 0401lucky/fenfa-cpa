import { Routes, Route, Navigate } from 'react-router-dom'
import { useEffect } from 'react'
import { useUserStore } from './store/userStore'
import Layout from './components/Layout'
import ProtectedRoute from './components/ProtectedRoute'
import Login from './pages/Login'
import AuthCallback from './pages/AuthCallback'
import Dashboard from './pages/Dashboard'
import Tokens from './pages/Tokens'
import Logs from './pages/Logs'
import Users from './pages/Users'
import IPBans from './pages/IPBans'
import Settings from './pages/Settings'

function App() {
  const { token, fetchUser } = useUserStore()

  useEffect(() => {
    if (token) {
      fetchUser()
    }
  }, [token, fetchUser])

  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/auth" element={<AuthCallback />} />
      <Route path="/" element={<ProtectedRoute><Layout /></ProtectedRoute>}>
        <Route index element={<Navigate to="/dashboard" replace />} />
        <Route path="dashboard" element={<Dashboard />} />
        <Route path="tokens" element={<Tokens />} />
        <Route path="logs" element={<Logs />} />
        <Route path="users" element={<ProtectedRoute adminOnly><Users /></ProtectedRoute>} />
        <Route path="ip-bans" element={<ProtectedRoute adminOnly><IPBans /></ProtectedRoute>} />
        <Route path="settings" element={<ProtectedRoute adminOnly><Settings /></ProtectedRoute>} />
      </Route>
    </Routes>
  )
}

export default App
