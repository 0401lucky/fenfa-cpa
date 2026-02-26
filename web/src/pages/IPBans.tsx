import { useCallback, useEffect, useState } from 'react'
import { Table, Button, Modal, Form, Input, DatePicker, Typography, message, Popconfirm, Tag } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons'
import { createIPBan, deleteIPBan, getErrorMessage, getIPBans, type IPBanInfo } from '../api'
import dayjs, { type Dayjs } from 'dayjs'

const { Title } = Typography

interface CreateIPBanFormValues {
  ip: string
  reason?: string
  expires_at?: Dayjs
}

export default function IPBans() {
  const [bans, setBans] = useState<IPBanInfo[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [createModalOpen, setCreateModalOpen] = useState(false)
  const [form] = Form.useForm()

  const fetchBans = useCallback(async (showLoading = false) => {
    if (showLoading) {
      setLoading(true)
    }
    getIPBans({ page, page_size: 20 }).then((res) => {
      setBans(res.data?.list || [])
      setTotal(res.data?.total || 0)
      setLoading(false)
    }).catch(() => setLoading(false))
  }, [page])

  useEffect(() => {
    const timer = window.setTimeout(() => {
      void fetchBans()
    }, 0)
    return () => window.clearTimeout(timer)
  }, [fetchBans])

  const handleCreate = async (values: CreateIPBanFormValues) => {
    try {
      const data: Record<string, string | number> = {
        ip: values.ip,
        reason: values.reason || '',
      }
      if (values.expires_at) {
        data.expires_at = values.expires_at.unix()
      }
      await createIPBan(data)
      setCreateModalOpen(false)
      form.resetFields()
      void fetchBans(true)
      message.success('已添加封禁')
    } catch (error) {
      message.error(getErrorMessage(error, '添加失败'))
    }
  }

  const handleDelete = async (id: number) => {
    try {
      await deleteIPBan(id)
      void fetchBans(true)
      message.success('已解除封禁')
    } catch (error) {
      message.error(getErrorMessage(error, '操作失败'))
    }
  }

  const columns: ColumnsType<IPBanInfo> = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: 'IP 地址', dataIndex: 'ip', key: 'ip' },
    { title: '原因', dataIndex: 'reason', key: 'reason' },
    {
      title: '过期时间', dataIndex: 'expires_at', key: 'expires_at',
      render: (v: number | null) => {
        if (!v) return <Tag color="red">永久</Tag>
        const d = dayjs.unix(v)
        return d.isBefore(dayjs()) ? <Tag color="orange">已过期</Tag> : d.format('YYYY-MM-DD HH:mm')
      },
    },
    {
      title: '创建时间', dataIndex: 'CreatedAt', key: 'created_at',
      render: (v: string) => v ? dayjs(v).format('YYYY-MM-DD HH:mm') : '-',
    },
    {
      title: '操作', key: 'action',
      render: (_, record) => (
        <Popconfirm title="确认解除封禁？" onConfirm={() => handleDelete(record.id)}>
          <Button size="small" danger icon={<DeleteOutlined />}>解除</Button>
        </Popconfirm>
      ),
    },
  ]

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0 }}>IP 封禁管理</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateModalOpen(true)}>
          添加封禁
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={bans}
        loading={loading}
        rowKey="id"
        pagination={{
          current: page,
          total,
          pageSize: 20,
          onChange: (nextPage) => {
            setLoading(true)
            setPage(nextPage)
          },
          showTotal: (t) => `共 ${t} 条`,
        }}
      />

      <Modal
        title="添加 IP 封禁"
        open={createModalOpen}
        onCancel={() => setCreateModalOpen(false)}
        onOk={() => form.submit()}
      >
        <Form form={form} layout="vertical" onFinish={handleCreate}>
          <Form.Item name="ip" label="IP 地址（支持 CIDR）" rules={[{ required: true, message: '请输入 IP' }]}>
            <Input placeholder="例如：1.2.3.4 或 10.0.0.0/8" />
          </Form.Item>
          <Form.Item name="reason" label="封禁原因">
            <Input placeholder="可选" />
          </Form.Item>
          <Form.Item name="expires_at" label="过期时间（留空=永久）">
            <DatePicker showTime style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
