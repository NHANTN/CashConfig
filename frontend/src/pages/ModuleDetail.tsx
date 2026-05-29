import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Card, Form, Input, Select, Button, message, Spin, Space, Tag,
} from 'antd';
import { getModule, createModule, updateModule } from '../api/module';
import { listAllGroups } from '../api/group';
import type { Module } from '../api/module';
import type { Group } from '../api/group';

const ENV_OPTIONS = [
  { label: '生产 (pr)', value: 'pr' },
  { label: '预发布 (pp)', value: 'pp' },
  { label: '实验室 (lab)', value: 'lab' },
];

const LOCATIONS = ['CN', 'HK', 'TW', 'SG', 'MY', 'TH', 'VN', 'ID', 'PH', 'KR', 'AU', 'KH', 'JP', 'NC', 'PF'];

export default function ModuleDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [allGroups, setAllGroups] = useState<Group[]>([]);
  const [selectedGroups, setSelectedGroups] = useState<Group[]>([]);
  const isEdit = !!id && id !== 'new';

  useEffect(() => {
    listAllGroups().then(setAllGroups).catch(() => {});
  }, []);

  useEffect(() => {
    if (isEdit) {
      setLoading(true);
      getModule(Number(id))
        .then((m: Module) => {
          form.setFieldsValue({ name: m.name, location: m.location, env: m.env });
          let ids: number[] = [];
          try { ids = JSON.parse(m.modules); } catch {}
          const matched = allGroups.filter((g) => ids.includes(g.id));
          setSelectedGroups(matched);
        })
        .catch(() => message.error('加载失败'))
        .finally(() => setLoading(false));
    }
  }, [id, allGroups]);

  const handleSubmit = async (values: any) => {
    const ids = selectedGroups.map((g) => g.id);
    const payload = {
      name: values.name,
      location: values.location,
      env: values.env,
      modules: JSON.stringify(ids),
    };
    setSaving(true);
    try {
      if (isEdit) {
        await updateModule(Number(id), payload);
        message.success('已更新');
      } else {
        await createModule(payload);
        message.success('已创建');
      }
      navigate('/modules');
    } catch {
      message.error('保存失败');
    } finally {
      setSaving(false);
    }
  };

  if (loading) return <Spin style={{ display: 'block', marginTop: 120 }} />;

  return (
    <Card title={isEdit ? '编辑 Module' : '新增 Module'} style={{ maxWidth: 800, margin: '0 auto' }}>
      <Form form={form} layout="vertical" onFinish={handleSubmit}>
        <Form.Item name="name" label="名称" rules={[{ required: true, message: '请输入名称' }]}>
          <Input placeholder="如: AllInOne_JPOS" />
        </Form.Item>
        <Space size={16}>
          <Form.Item name="location" label="地区" rules={[{ required: true, message: '请选择地区' }]}>
            <Select placeholder="选择地区" style={{ width: 120 }} options={LOCATIONS.map(l => ({ label: l, value: l }))} />
          </Form.Item>
          <Form.Item name="env" label="环境" rules={[{ required: true, message: '请选择环境' }]}>
            <Select placeholder="选择环境" style={{ width: 150 }} options={ENV_OPTIONS} />
          </Form.Item>
        </Space>

        <Form.Item label="引用 Group 组">
          <Select
            mode="multiple"
            placeholder="搜索并选择 Group"
            style={{ width: '100%' }}
            value={selectedGroups.map((g) => g.id)}
            onChange={(ids: number[]) => {
              const matched = allGroups.filter((g) => ids.includes(g.id));
              setSelectedGroups(matched);
            }}
            filterOption={(input, option) =>
              (option?.label as string)?.toLowerCase().includes(input.toLowerCase())
            }
            options={allGroups.map((g) => ({
              label: g.name,
              value: g.id,
            }))}
          />
          <div style={{ marginTop: 8 }}>
            {selectedGroups.map((g) => {
              let steps: { name: string }[] = [];
              try { steps = JSON.parse(g.steps); } catch {}
              return (
                <Tag key={g.id} color="blue" style={{ marginBottom: 4, padding: '2px 8px' }}>
                  {g.name} ({steps.length} steps)
                </Tag>
              );
            })}
          </div>
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" loading={saving} style={{ marginRight: 8 }}>
            {isEdit ? '更新' : '创建'}
          </Button>
          <Button onClick={() => navigate('/modules')}>取消</Button>
        </Form.Item>
      </Form>
    </Card>
  );
}
