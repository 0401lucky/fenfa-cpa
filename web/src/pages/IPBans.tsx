import { useEffect, useState } from 'react'
import { Table, Button, Modal, Form, Input, DatePicker, Typography, message, Popconfirm, Tag } from 'antd'
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons'
import { getIPBans, createIPBan, deleteIPBan } from '../api'
import dayjs from 'dayjs'

const { Title } = Typography

export default function IPBans() {
  const [bans, setBans] = useState<any[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [createModalOpen, setCreateModalOpen] = useState(false)
  const [form] = Form.useForm()

  const fetchBans = () => {
    setLoading(true)
    getIPBans({ page, page_size: 20 }).then((res: any) => {
      setBans(res.data?.list || [])
      setTotal(res.data?.total || 0)
      setLoading(false)
    }).catch(() => setLoading(false))
  }

  useEffect(() => { fetchBans() }, [page])

  const handleCreate = async (values: any) => {
    try {
      const data: any = {
        ip: values.ip,
        reason: values.reason || '',
      }
      if (values.expires_at) {
        data.expires_at = values.expires_at.unix()
      }
      await createIPBan(data)
      setCreateModalOpen(false)
      form.resetFields()
      fetchBans()
      message.success('已添加封禁')
    } catch (err: any) {
      message.error(err?.message || '添加失败')
    }
  }

  const handleDelete = async (id: number) => {
    try {
      await deleteIPBan(id)
      fetchBans()
      message.success('已解除封禁')
    } catch (err: any) {
      message.error(err?.message || '操作失败')
    }
  }

  const columns = [
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
      render: (_: any, record: any) => (
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
          onChange: setPage,
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
