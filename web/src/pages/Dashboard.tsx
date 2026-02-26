import { useEffect, useState } from 'react'
import { Card, Col, Row, Statistic, Typography, Table, Spin } from 'antd'
import { ApiOutlined, ThunderboltOutlined, ClockCircleOutlined, KeyOutlined } from '@ant-design/icons'
import { getDashboard, type DashboardData } from '../api'
import { useUserStore } from '../store/userStore'

const { Title } = Typography

export default function Dashboard() {
  const [data, setData] = useState<DashboardData | null>(null)
  const [loading, setLoading] = useState(true)
  const { user } = useUserStore()
  const isAdmin = user && user.role >= 10

  useEffect(() => {
    getDashboard().then((res) => {
      setData(res.data)
      setLoading(false)
    }).catch(() => setLoading(false))
  }, [])

  if (loading) return <Spin size="large" style={{ display: 'block', margin: '100px auto' }} />
  if (!data) return null

  return (
    <div>
      <Title level={4} style={{ marginBottom: 24 }}>仪表盘</Title>

      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="总请求数"
              value={data.stats?.total_requests || 0}
              prefix={<ApiOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="今日请求"
              value={data.stats?.today_requests || 0}
              prefix={<ClockCircleOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="总 Token 用量"
              value={data.stats?.total_tokens || 0}
              prefix={<ThunderboltOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="密钥数量"
              value={data.token_count || 0}
              prefix={<KeyOutlined />}
              suffix={`/ ${data.user?.quota_total === -1 ? '∞' : data.user?.quota_total}`}
            />
          </Card>
        </Col>
      </Row>

      {data.user && (
        <Card style={{ marginTop: 16 }}>
          <Row gutter={16}>
            <Col span={12}>
              <Statistic
                title="配额使用"
                value={data.user.quota_used}
                suffix={`/ ${data.user.quota_total === -1 ? '无限' : data.user.quota_total}`}
              />
            </Col>
            <Col span={12}>
              <Statistic
                title="今日 Token"
                value={data.stats?.today_tokens || 0}
              />
            </Col>
          </Row>
        </Card>
      )}

      {isAdmin && (
        <>
          <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
            <Col xs={24} sm={12}>
              <Card title="全局统计">
                <Statistic title="全局总请求" value={data.global_stats?.total_requests || 0} />
                <Statistic title="全局总 Token" value={data.global_stats?.total_tokens || 0} style={{ marginTop: 16 }} />
                <Statistic title="用户总数" value={data.user_count || 0} style={{ marginTop: 16 }} />
              </Card>
            </Col>
            <Col xs={24} sm={12}>
              <Card title="请求趋势 (近7天)">
                {data.trend && (
                  <Table
                    dataSource={data.trend}
                    columns={[
                      { title: '日期', dataIndex: 'date', key: 'date' },
                      { title: '请求数', dataIndex: 'count', key: 'count' },
                    ]}
                    pagination={false}
                    size="small"
                    rowKey="date"
                  />
                )}
              </Card>
            </Col>
          </Row>

          {data.model_distribution && data.model_distribution.length > 0 && (
            <Card title="模型分布" style={{ marginTop: 16 }}>
              <Table
                dataSource={data.model_distribution}
                columns={[
                  { title: '模型', dataIndex: 'model', key: 'model' },
                  { title: '调用次数', dataIndex: 'count', key: 'count' },
                ]}
                pagination={false}
                size="small"
                rowKey="model"
              />
            </Card>
          )}
        </>
      )}
    </div>
  )
}
