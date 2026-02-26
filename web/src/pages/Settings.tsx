import { useEffect, useState } from 'react'
import { Card, Form, Input, Button, Typography, message, Divider } from 'antd'
import { getErrorMessage, getSettings, updateSettings, type SettingsMap } from '../api'

const { Title } = Typography

export default function Settings() {
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => {
    getSettings().then((res) => {
      form.setFieldsValue(res.data || {})
      setLoading(false)
    }).catch(() => setLoading(false))
  }, [form])

  const handleSave = async (values: Record<string, unknown>) => {
    setSaving(true)
    try {
      // Convert numbers to strings for KV store
      const data: SettingsMap = {}
      for (const [k, v] of Object.entries(values)) {
        if (v !== undefined && v !== null && v !== '') {
          data[k] = String(v)
        }
      }
      await updateSettings(data)
      message.success('设置已保存')
    } catch (error) {
      message.error(getErrorMessage(error, '保存失败'))
    }
    setSaving(false)
  }

  return (
    <div>
      <Title level={4} style={{ marginBottom: 24 }}>系统设置</Title>

      <Card loading={loading}>
        <Form form={form} layout="vertical" onFinish={handleSave} style={{ maxWidth: 600 }}>
          <Divider>上游配置</Divider>
          <Form.Item name="cpa_upstream_url" label="CPA 上游地址">
            <Input placeholder="https://xxx.zeabur.app" />
          </Form.Item>
          <Form.Item name="cpa_upstream_key" label="CPA 上游密钥">
            <Input.Password placeholder="sk-xxx" />
          </Form.Item>

          <Divider>OAuth 配置</Divider>
          <Form.Item name="linuxdo_client_id" label="LinuxDO Client ID">
            <Input />
          </Form.Item>
          <Form.Item name="linuxdo_client_secret" label="LinuxDO Client Secret">
            <Input.Password />
          </Form.Item>

          <Divider>站点配置</Divider>
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
