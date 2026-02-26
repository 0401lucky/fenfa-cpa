import { useCallback, useEffect, useState } from 'react'
import { Table, Card, Input, Typography, Tag, Space } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { getLogs, type RequestLogInfo } from '../api'
import dayjs from 'dayjs'

const { Title } = Typography

export default function Logs() {
  const [logs, setLogs] = useState<RequestLogInfo[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)
  const [modelFilter, setModelFilter] = useState('')

  const fetchLogs = useCallback(async (showLoading = false) => {
    if (showLoading) {
      setLoading(true)
    }
    getLogs({ page, page_size: pageSize, model: modelFilter || undefined }).then((res) => {
      setLogs(res.data?.list || [])
      setTotal(res.data?.total || 0)
      setLoading(false)
    }).catch(() => setLoading(false))
  }, [page, pageSize, modelFilter])

  useEffect(() => {
    const timer = window.setTimeout(() => {
      void fetchLogs()
    }, 0)
    return () => window.clearTimeout(timer)
  }, [fetchLogs])

  const columns: ColumnsType<RequestLogInfo> = [
    {
      title: '时间', dataIndex: 'created_at', key: 'created_at', width: 160,
      render: (v: string) => dayjs(v).format('MM-DD HH:mm:ss'),
    },
    { title: '模型', dataIndex: 'model', key: 'model', width: 200 },
    {
      title: '状态', dataIndex: 'status_code', key: 'status_code', width: 80,
      render: (v: number) => (
        <Tag color={v >= 200 && v < 300 ? 'green' : v >= 400 ? 'red' : 'orange'}>{v}</Tag>
      ),
    },
    { title: '耗时(ms)', dataIndex: 'duration', key: 'duration', width: 100 },
    { title: '输入', dataIndex: 'prompt_tokens', key: 'prompt_tokens', width: 80 },
    { title: '输出', dataIndex: 'completion_tokens', key: 'completion_tokens', width: 80 },
    { title: '总计', dataIndex: 'total_tokens', key: 'total_tokens', width: 80 },
    { title: 'IP', dataIndex: 'request_ip', key: 'request_ip', width: 140 },
    { title: '路径', dataIndex: 'path', key: 'path', ellipsis: true },
    {
      title: '错误', dataIndex: 'error_message', key: 'error_message', width: 200,
      render: (v: string) => v ? <Tag color="red">{v}</Tag> : '-',
    },
  ]

  return (
    <div>
      <Title level={4} style={{ marginBottom: 16 }}>调用日志</Title>

      <Card style={{ marginBottom: 16 }}>
        <Space>
          <Input
            placeholder="筛选模型"
            value={modelFilter}
            onChange={(e) => { setModelFilter(e.target.value); setPage(1) }}
            allowClear
            style={{ width: 200 }}
          />
        </Space>
      </Card>

      <Table
        columns={columns}
        dataSource={logs}
        loading={loading}
        rowKey="id"
        scroll={{ x: 1200 }}
        pagination={{
          current: page,
          pageSize,
          total,
          showSizeChanger: true,
          showTotal: (t) => `共 ${t} 条`,
          onChange: (p, ps) => {
            setLoading(true)
            setPage(p)
            setPageSize(ps)
          },
        }}
      />
    </div>
  )
}
