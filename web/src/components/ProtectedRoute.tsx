import { Navigate } from 'react-router-dom'
import { useUserStore } from '../store/userStore'
import { Spin } from 'antd'

interface Props {
  children: React.ReactNode
  adminOnly?: boolean
}

export default function ProtectedRoute({ children, adminOnly }: Props) {
  const { token, user, loading } = useUserStore()

  if (!token) {
    return <Navigate to="/login" replace />
  }

  if (loading || !user) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <Spin size="large" />
      </div>
    )
  }

  if (adminOnly && user.role < 10) {
    return <Navigate to="/dashboard" replace />
  }

  return <>{children}</>
}
