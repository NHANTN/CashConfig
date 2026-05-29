import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Card, Form, Input, Select, Button, message, Spin } from 'antd';
import { getStore, createStore, updateStore } from '../api/store';
import type { Store } from '../api/store';

const LOCATIONS = ['CN', 'HK', 'TW', 'SG', 'MY', 'TH', 'VN', 'ID', 'PH', 'KR', 'AU', 'KH', 'JP', 'NC', 'PF'];
const WEBPOS_OPTS = [
  { label: '生产 (pr)', value: 'pr' },
  { label: '预发布 (pp)', value: 'pp' },
  { label: '测试 (pu)', value: 'pu' },
];

export default function StoreDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const isEdit = !!id && id !== 'new';

  useEffect(() => {
    if (isEdit) {
      setLoading(true);
      getStore(Number(id))
        .then((m: Store) => form.setFieldsValue(m))
        .catch(() => message.error('加载失败'))
        .finally(() => setLoading(false));
    }
  }, [id]);

  const handleSubmit = async (values: any) => {
    setSaving(true);
    try {
      if (isEdit) { await updateStore(Number(id), values); message.success('已更新'); }
      else { await createStore(values); message.success('已创建'); }
      navigate('/stores');
    } catch { message.error('保存失败'); } finally { setSaving(false); }
  };

  if (loading) return <Spin style={{ display: 'block', marginTop: 120 }} />;

  return (
    <Card title={isEdit ? '编辑 Store' : '新增 Store'} style={{ maxWidth: 800, margin: '0 auto' }}>
      <Form form={form} layout="vertical" onFinish={handleSubmit}>
        <Form.Item name="store_number" label="门店号" rules={[{ required: true }]}>
          <Input placeholder="如: 806" />
        </Form.Item>
        <Form.Item name="network_segment" label="网段">
          <Input placeholder="如: 10.66.48.0" />
        </Form.Item>
        <Form.Item name="webpos_env" label="Webpos 环境">
          <Select options={WEBPOS_OPTS} allowClear />
        </Form.Item>
        <Form.Item name="eft" label="EFT">
          <Input placeholder="如: ICBC" />
        </Form.Item>
        <Form.Item name="location" label="地区">
          <Select options={LOCATIONS.map(l => ({ label: l, value: l }))} allowClear />
        </Form.Item>
        <Form.Item name="rf_server" label="RF 服务器">
          <Input placeholder="如: 10.69.127.244" />
        </Form.Item>
        <Form.Item name="cashtill_seg_gw" label="Cashtill 网段网关">
          <Input placeholder="如: 10.69.63.225" />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit" loading={saving} style={{ marginRight: 8 }}>
            {isEdit ? '更新' : '创建'}
          </Button>
          <Button onClick={() => navigate('/stores')}>取消</Button>
        </Form.Item>
      </Form>
    </Card>
  );
}
