import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Card, Form, Input, Select, Button, message, Spin } from 'antd';
import { getVar, createVar, updateVar } from '../api/var';
import type { Var } from '../api/var';

const ENV_OPTIONS = [
  { label: '生产 (pr)', value: 'pr' },
  { label: '预发布 (pp)', value: 'pp' },
  { label: '实验室 (lab)', value: 'lab' },
];

export default function VarDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const isEdit = !!id && id !== 'new';

  useEffect(() => {
    if (isEdit) {
      setLoading(true);
      getVar(Number(id))
        .then((m: Var) => form.setFieldsValue(m))
        .catch(() => message.error('加载失败'))
        .finally(() => setLoading(false));
    }
  }, [id]);

  const handleSubmit = async (values: any) => {
    setSaving(true);
    try {
      values.matcher = values.matcher || '[]';
      if (isEdit) { await updateVar(Number(id), values); message.success('已更新'); }
      else { await createVar(values); message.success('已创建'); }
      navigate('/vars');
    } catch { message.error('保存失败'); } finally { setSaving(false); }
  };

  if (loading) return <Spin style={{ display: 'block', marginTop: 120 }} />;

  return (
    <Card title={isEdit ? '编辑 Var' : '新增 Var'} style={{ maxWidth: 800, margin: '0 auto' }}>
      <Form form={form} layout="vertical" onFinish={handleSubmit}>
        <Form.Item name="var_name" label="变量名" rules={[{ required: true }]}>
          <Input placeholder="如: Greengate_Store" />
        </Form.Item>
        <Form.Item name="value" label="值" rules={[{ required: true }]}>
          <Input placeholder="如: 806" />
        </Form.Item>
        <Form.Item name="env" label="环境" rules={[{ required: true }]}>
          <Select options={ENV_OPTIONS} />
        </Form.Item>
        <Form.Item name="matcher" label="Matcher (JSON)">
          <Input placeholder='默认为 []' />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit" loading={saving} style={{ marginRight: 8 }}>
            {isEdit ? '更新' : '创建'}
          </Button>
          <Button onClick={() => navigate('/vars')}>取消</Button>
        </Form.Item>
      </Form>
    </Card>
  );
}
