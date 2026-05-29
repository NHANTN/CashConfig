import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Card, Form, Input, Select, Button, message, Spin } from 'antd';
import { getRule, createRule, updateRule } from '../api/rule';
import type { Rule } from '../api/rule';

const TYPES = [
  { label: '环境匹配 (env)', value: 'env' },
  { label: '分组匹配 (group)', value: 'group' },
];
const LOCATIONS = ['CN', 'HK', 'TW', 'SG', 'MY', 'TH', 'VN', 'ID', 'PH', 'KR', 'AU', 'KH', 'JP', 'NC', 'PF'];

export default function RuleDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const isEdit = !!id && id !== 'new';

  useEffect(() => {
    if (isEdit) {
      setLoading(true);
      getRule(Number(id))
        .then((m: Rule) => form.setFieldsValue(m))
        .catch(() => message.error('加载失败'))
        .finally(() => setLoading(false));
    }
  }, [id]);

  const handleSubmit = async (values: any) => {
    setSaving(true);
    try {
      if (isEdit) { await updateRule(Number(id), values); message.success('已更新'); }
      else { await createRule(values); message.success('已创建'); }
      navigate('/rules');
    } catch { message.error('保存失败'); } finally { setSaving(false); }
  };

  if (loading) return <Spin style={{ display: 'block', marginTop: 120 }} />;

  return (
    <Card title={isEdit ? '编辑 Rule' : '新增 Rule'} style={{ maxWidth: 800, margin: '0 auto' }}>
      <Form form={form} layout="vertical" onFinish={handleSubmit}>
        <Form.Item name="name" label="名称" rules={[{ required: true }]}>
          <Input placeholder="如: CN_PR_Register" />
        </Form.Item>
        <Form.Item name="type" label="类型" rules={[{ required: true }]}>
          <Select options={TYPES} />
        </Form.Item>
        <Form.Item name="location" label="地区" rules={[{ required: true }]}>
          <Select options={LOCATIONS.map(l => ({ label: l, value: l }))} />
        </Form.Item>
        <Form.Item name="env_name" label="环境变量名" rules={[{ required: true }]}>
          <Input placeholder="如: InfraAuto_Hostname" />
        </Form.Item>
        <Form.Item name="rule" label="规则 (正则)" rules={[{ required: true }]}>
          <Input placeholder="如: CNT0806.*" />
        </Form.Item>
        <Form.Item name="result" label="结果" rules={[{ required: true }]}>
          <Input placeholder="如: lab" />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit" loading={saving} style={{ marginRight: 8 }}>
            {isEdit ? '更新' : '创建'}
          </Button>
          <Button onClick={() => navigate('/rules')}>取消</Button>
        </Form.Item>
      </Form>
    </Card>
  );
}
