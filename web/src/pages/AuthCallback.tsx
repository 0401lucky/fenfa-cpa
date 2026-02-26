import { useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useUserStore } from '../store/userStore'
import { Spin } from 'antd'

export default function AuthCallback() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const { setToken, fetchUser } = useUserStore()

  useEffect(() => {
    const token = searchParams.get('token')
    if (token) {
      setToken(token)
      fetchUser().then(() => {
        navigate('/dashboard', { replace: true })
      })
    } else {
      navigate('/login', { replace: true })
    }
  }, [searchParams, setToken, fetchUser, navigate])

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
      <Spin size="large" tip="正在登录..." />
    </div>
  )
}
