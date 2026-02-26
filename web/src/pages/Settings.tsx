import { useEffect, useState } from 'react'
import { Card, Form, Input, InputNumber, Button, Typography, message, Divider, Space } from 'antd'
import { getSettings, updateSettings } from '../api'

const { Title, Text } = Typography

export default function Settings() {
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => {
    getSettings().then((res: any) => {
      form.setFieldsValue(res.data || {})
      setLoading(false)
    }).catch(() => setLoading(false))
  }, [form])

  const handleSave = async (values: any) => {
    setSaving(true)
    try {
      // Convert numbers to strings for KV store
      const data: Record<string, string> = {}
      for (const [k, v] of Object.entries(values)) {
        if (v !== undefined && v !== null && v !== '') {
          data[k] = String(v)
        }
      }
      await updateSettings(data)
      message.success('设置已保存')
    } catch (err: any) {
      message.error(err?.message || '保存失败')
    }
    setSaving(false)
  }

  return (
    <div>
      <Title level={4} style={{ marginBottom: 24 }}>系统设置</Title>

      <Card loading={loading}>
        <Form form={form} layout="vertical" onFinish={handleSave} style={{ maxWidth: 600 }}>
          <Divider orientation="left">上游配置</Divider>
          <Form.Item name="cpa_upstream_url" label="CPA 上游地址">
            <Input placeholder="https://xxx.zeabur.app" />
          </Form.Item>
          <Form.Item name="cpa_upstream_key" label="CPA 上游密钥">
            <Input.Password placeholder="sk-xxx" />
          </Form.Item>

          <Divider orientation="left">OAuth 配置</Divider>
          <Form.Item name="linuxdo_client_id" label="LinuxDO Client ID">
            <Input />
          </Form.Item>
          <Form.Item name="linuxdo_client_secret" label="LinuxDO Client Secret">
            <Input.Password />
          </Form.Item>

          <Divider orientation="left">站点配置</Divider>
          <Form.Item name="site_name" label="站点名称">
            <Input placeholder="CPA 分发系统" />
          </Form.Item>
          <Form.Item name="min_trust_level" label="最低信任等级">
            <Input placeholder="0" />
          </Form.Item>
          <Form.Item name="default_quota" label="新用户默认配额">
            <Input placeholder="1000" />
          </Form.Item>
          <Form.Item name="log_retention_days" label="日志保留天数">
            <Input placeholder="30" />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={saving}>
              保存设置
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}
