import { useEffect, useState } from 'react'
import { Table, Button, Modal, Form, Input, InputNumber, Space, Tag, Typography, message, Popconfirm, Switch, Select } from 'antd'
import { PlusOutlined, CopyOutlined, ReloadOutlined, DeleteOutlined, EditOutlined } from '@ant-design/icons'
import { getTokens, createToken, updateToken, deleteToken, resetToken } from '../api'
import dayjs from 'dayjs'

const { Title, Text, Paragraph } = Typography

export default function Tokens() {
  const [tokens, setTokens] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [createModalOpen, setCreateModalOpen] = useState(false)
  const [editModalOpen, setEditModalOpen] = useState(false)
  const [newKeyModalOpen, setNewKeyModalOpen] = useState(false)
  const [newKey, setNewKey] = useState('')
  const [editingToken, setEditingToken] = useState<any>(null)
  const [createForm] = Form.useForm()
  const [editForm] = Form.useForm()

  const fetchTokens = () => {
    setLoading(true)
    getTokens().then((res: any) => {
      setTokens(res.data || [])
      setLoading(false)
    }).catch(() => setLoading(false))
  }

  useEffect(() => { fetchTokens() }, [])

  const handleCreate = async (values: any) => {
    try {
      const res: any = await createToken(values)
      setNewKey(res.data.key)
      setNewKeyModalOpen(true)
      setCreateModalOpen(false)
      createForm.resetFields()
      fetchTokens()
      message.success('密钥创建成功')
    } catch (err: any) {
      message.error(err?.message || '创建失败')
    }
  }

  const handleUpdate = async (values: any) => {
    if (!editingToken) return
    try {
      await updateToken(editingToken.id, values)
      setEditModalOpen(false)
      fetchTokens()
      message.success('更新成功')
    } catch (err: any) {
      message.error(err?.message || '更新失败')
    }
  }

  const handleDelete = async (id: number) => {
    try {
      await deleteToken(id)
      fetchTokens()
      message.success('已删除')
    } catch (err: any) {
      message.error(err?.message || '删除失败')
    }
  }

  const handleReset = async (id: number) => {
    try {
      const res: any = await resetToken(id)
      setNewKey(res.data.key)
      setNewKeyModalOpen(true)
      fetchTokens()
      message.success('密钥已重置')
    } catch (err: any) {
      message.error(err?.message || '重置失败')
    }
  }

  const copyKey = () => {
    navigator.clipboard.writeText(newKey)
    message.success('已复制到剪贴板')
  }

  const columns = [
    { title: '名称', dataIndex: 'name', key: 'name' },
    {
      title: '密钥前缀', dataIndex: 'key_prefix', key: 'key_prefix',
      render: (v: string) => <Text code>{v}</Text>,
    },
    {
      title: '状态', dataIndex: 'status', key: 'status',
      render: (v: number) => v === 1 ? <Tag color="green">启用</Tag> : <Tag color="red">禁用</Tag>,
    },
    {
      title: '配额', key: 'quota',
      render: (_: any, r: any) => r.quota_total === -1 ? '跟随用户' : `${r.quota_used} / ${r.quota_total}`,
    },
    {
      title: 'RPM', dataIndex: 'rate_limit_rpm', key: 'rpm',
    },
    {
      title: '总请求', dataIndex: 'total_requests', key: 'total_requests',
    },
    {
      title: '过期时间', dataIndex: 'expires_at', key: 'expires_at',
      render: (v: number | null) => v ? dayjs.unix(v).format('YYYY-MM-DD HH:mm') : '永不',
    },
    {
      title: '创建时间', dataIndex: 'CreatedAt', key: 'created_at',
      render: (v: string) => v ? dayjs(v).format('YYYY-MM-DD HH:mm') : '-',
    },
    {
      title: '操作', key: 'action',
      render: (_: any, record: any) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => {
            setEditingToken(record)
            editForm.setFieldsValue({
              name: record.name,
              status: record.status,
              quota_total: record.quota_total,
              rate_limit_rpm: record.rate_limit_rpm,
              allowed_models: record.allowed_models,
              allowed_ips: record.allowed_ips,
            })
            setEditModalOpen(true)
          }}>编辑</Button>
          <Popconfirm title="确认重置密钥？旧密钥将失效。" onConfirm={() => handleReset(record.id)}>
            <Button size="small" icon={<ReloadOutlined />}>重置</Button>
          </Popconfirm>
          <Popconfirm title="确认删除？" onConfirm={() => handleDelete(record.id)}>
            <Button size="small" danger icon={<DeleteOutlined />}>删除</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0 }}>密钥管理</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateModalOpen(true)}>
          创建密钥
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={tokens}
        loading={loading}
        rowKey="id"
        scroll={{ x: 1000 }}
      />

      {/* Create Modal */}
      <Modal
        title="创建密钥"
        open={createModalOpen}
        onCancel={() => setCreateModalOpen(false)}
        onOk={() => createForm.submit()}
      >
        <Form form={createForm} layout="vertical" onFinish={handleCreate}>
          <Form.Item name="name" label="名称" rules={[{ required: true, message: '请输入名称' }]}>
            <Input placeholder="例如：ChatBox 使用" />
          </Form.Item>
          <Form.Item name="quota_total" label="配额" initialValue={-1}>
            <InputNumber style={{ width: '100%' }} placeholder="-1 表示跟随用户" />
          </Form.Item>
          <Form.Item name="rate_limit_rpm" label="每分钟请求上限" initialValue={60}>
            <InputNumber style={{ width: '100%' }} min={0} />
          </Form.Item>
          <Form.Item name="allowed_models" label="允许模型（逗号分隔，留空不限）">
            <Input placeholder="gpt-4o,claude-3-5-sonnet" />
          </Form.Item>
          <Form.Item name="allowed_ips" label="IP 白名单（逗号分隔，留空不限）">
            <Input placeholder="1.2.3.4,10.0.0.0/8" />
          </Form.Item>
        </Form>
      </Modal>

      {/* Edit Modal */}
      <Modal
        title="编辑密钥"
        open={editModalOpen}
        onCancel={() => setEditModalOpen(false)}
        onOk={() => editForm.submit()}
      >
        <Form form={editForm} layout="vertical" onFinish={handleUpdate}>
          <Form.Item name="name" label="名称">
            <Input />
          </Form.Item>
          <Form.Item name="status" label="状态">
            <Select options={[
              { label: '启用', value: 1 },
              { label: '禁用', value: 2 },
            ]} />
          </Form.Item>
          <Form.Item name="quota_total" label="配额">
            <InputNumber style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="rate_limit_rpm" label="RPM">
            <InputNumber style={{ width: '100%' }} min={0} />
          </Form.Item>
          <Form.Item name="allowed_models" label="允许模型">
            <Input />
          </Form.Item>
          <Form.Item name="allowed_ips" label="IP 白名单">
            <Input />
          </Form.Item>
        </Form>
      </Modal>

      {/* New Key Display Modal */}
      <Modal
        title="密钥已生成"
        open={newKeyModalOpen}
        onCancel={() => { setNewKeyModalOpen(false); setNewKey('') }}
        footer={[
          <Button key="copy" type="primary" icon={<CopyOutlined />} onClick={copyKey}>复制密钥</Button>,
          <Button key="close" onClick={() => { setNewKeyModalOpen(false); setNewKey('') }}>关闭</Button>,
        ]}
      >
        <div style={{ marginBottom: 16 }}>
          <Text type="warning">请立即复制保存密钥，关闭后将无法再次查看！</Text>
        </div>
        <Paragraph code copyable style={{ wordBreak: 'break-all', background: '#f5f5f5', padding: 12, borderRadius: 6 }}>
          {newKey}
        </Paragraph>
      </Modal>
    </div>
  )
}
