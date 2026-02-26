import { Button, Card, Typography, Space } from 'antd'
import { LoginOutlined } from '@ant-design/icons'
import { useUserStore } from '../store/userStore'
import { Navigate } from 'react-router-dom'
import { getOAuthURL } from '../api'

const { Title, Text } = Typography

export default function Login() {
  const { token } = useUserStore()

  if (token) {
    return <Navigate to="/dashboard" replace />
  }

  const handleLogin = () => {
    window.location.href = getOAuthURL()
  }

  return (
    <div style={{
      minHeight: '100vh',
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    }}>
      <Card style={{ width: 400, textAlign: 'center', borderRadius: 12 }} bordered={false}>
        <Space direction="vertical" size="large" style={{ width: '100%' }}>
          <div>
            <Title level={2} style={{ marginBottom: 8 }}>CPA 分发系统</Title>
            <Text type="secondary">CLIProxyAPI 密钥管理与分发平台</Text>
          </div>
          <Button
            type="primary"
            size="large"
            icon={<LoginOutlined />}
            onClick={handleLogin}
            block
            style={{ height: 48, fontSize: 16 }}
          >
            使用 LinuxDO 账号登录
          </Button>
          <Text type="secondary" style={{ fontSize: 12 }}>
            登录即表示您同意服务条款和隐私政策
          </Text>
        </Space>
      </Card>
    </div>
  )
}
