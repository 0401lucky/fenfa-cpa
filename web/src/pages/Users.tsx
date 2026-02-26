import { useEffect, useState } from 'react'
import { Table, Tag, Button, Modal, Form, InputNumber, Select, Typography, message } from 'antd'
import { getUsers, updateUser } from '../api'
import dayjs from 'dayjs'

const { Title } = Typography

const roleMap: Record<number, { label: string; color: string }> = {
  1: { label: '用户', color: 'default' },
  10: { label: '管理员', color: 'blue' },
  100: { label: '超管', color: 'red' },
}

export default function Users() {
  const [users, setUsers] = useState<any[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [editModalOpen, setEditModalOpen] = useState(false)
  const [editingUser, setEditingUser] = useState<any>(null)
  const [form] = Form.useForm()

  const fetchUsers = () => {
    setLoading(true)
    getUsers({ page, page_size: 20 }).then((res: any) => {
      setUsers(res.data?.list || [])
      setTotal(res.data?.total || 0)
      setLoading(false)
    }).catch(() => setLoading(false))
  }

  useEffect(() => { fetchUsers() }, [page])

  const handleEdit = (record: any) => {
    setEditingUser(record)
    form.setFieldsValue({
      role: record.role,
      status: record.status,
      quota_total: record.quota_total,
      token_limit: record.token_limit,
    })
    setEditModalOpen(true)
  }

  const handleUpdate = async (values: any) => {
    try {
      await updateUser(editingUser.id, values)
      setEditModalOpen(false)
      fetchUsers()
      message.success('更新成功')
    } catch (err: any) {
      message.error(err?.message || '更新失败')
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '用户名', dataIndex: 'username', key: 'username' },
    { title: '显示名', dataIndex: 'display_name', key: 'display_name' },
    {
      title: '角色', dataIndex: 'role', key: 'role',
      render: (v: number) => <Tag color={roleMap[v]?.color}>{roleMap[v]?.label || v}</Tag>,
    },
    {
      title: '状态', dataIndex: 'status', key: 'status',
      render: (v: number) => v === 1 ? <Tag color="green">启用</Tag> : <Tag color="red">禁用</Tag>,
    },
    { title: '信任等级', dataIndex: 'trust_level', key: 'trust_level', width: 80 },
    {
      title: '配额', key: 'quota',
      render: (_: any, r: any) => `${r.quota_used} / ${r.quota_total === -1 ? '∞' : r.quota_total}`,
    },
    { title: '密钥上限', dataIndex: 'token_limit', key: 'token_limit', width: 80 },
    {
      title: '最后登录', dataIndex: 'last_login_at', key: 'last_login_at',
      render: (v: number | null) => v ? dayjs.unix(v).format('YYYY-MM-DD HH:mm') : '-',
    },
    { title: '登录IP', dataIndex: 'last_login_ip', key: 'last_login_ip' },
    {
      title: '操作', key: 'action',
      render: (_: any, record: any) => (
        <Button size="small" onClick={() => handleEdit(record)}>编辑</Button>
      ),
    },
  ]

  return (
    <div>
      <Title level={4} style={{ marginBottom: 16 }}>用户管理</Title>
      <Table
        columns={columns}
        dataSource={users}
        loading={loading}
        rowKey="id"
        scroll={{ x: 1200 }}
        pagination={{
          current: page,
          total,
          pageSize: 20,
          onChange: setPage,
          showTotal: (t) => `共 ${t} 条`,
        }}
      />

      <Modal
        title={`编辑用户: ${editingUser?.username}`}
        open={editModalOpen}
        onCancel={() => setEditModalOpen(false)}
        onOk={() => form.submit()}
      >
        <Form form={form} layout="vertical" onFinish={handleUpdate}>
          <Form.Item name="role" label="角色">
            <Select options={[
              { label: '普通用户', value: 1 },
              { label: '管理员', value: 10 },
              { label: '超级管理员', value: 100 },
            ]} />
          </Form.Item>
          <Form.Item name="status" label="状态">
            <Select options={[
              { label: '启用', value: 1 },
              { label: '禁用', value: 2 },
            ]} />
          </Form.Item>
          <Form.Item name="quota_total" label="总配额（-1=无限）">
            <InputNumber style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="token_limit" label="密钥数量上限">
            <InputNumber style={{ width: '100%' }} min={1} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
