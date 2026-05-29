import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Card, Form, Input, Button, message, Spin, Space, Cascader,
} from 'antd';
import { MinusCircleOutlined, PlusOutlined } from '@ant-design/icons';
import { getGroup, createGroup, updateGroup } from '../api/group';
import { listScriptFiles } from '../api/script';
import type { Group } from '../api/group';

interface StepValues {
  name: string;
  path: string;
}

function buildScriptCascader(files: string[]) {
  const topMap: Record<string, Record<string, string[]>> = {};
  for (const f of files) {
    const parts = f.split('/');
    const top = parts[0] || '(根目录)';
    const sub = parts.length > 2 ? parts[1] : '';
    if (!topMap[top]) topMap[top] = {};
    if (!topMap[top][sub]) topMap[top][sub] = [];
    topMap[top][sub].push(f);
  }
  return Object.entries(topMap)
    .sort(([a], [b]) => a.localeCompare(b))
    .map(([top, subMap]) => ({
      label: top,
      value: top,
      children: Object.entries(subMap)
        .sort(([a], [b]) => a.localeCompare(b))
        .map(([sub, paths]) => ({
          label: sub || '(子目录)',
          value: sub,
          children: paths
            .sort((a, b) => a.localeCompare(b))
            .map((p) => ({ label: p, value: p })),
        })),
    }));
}

export default function GroupDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [scriptFiles, setScriptFiles] = useState<string[]>([]);
  const isEdit = !!id && id !== 'new';

  useEffect(() => {
    listScriptFiles().then(setScriptFiles).catch(() => message.error('加载脚本列表失败'));
  }, []);

  useEffect(() => {
    if (isEdit) {
      setLoading(true);
      getGroup(Number(id))
        .then((m: Group) => {
          let steps: StepValues[] = [];
          try { steps = JSON.parse(m.steps); } catch { steps = []; }
          form.setFieldsValue({ ...m, steps });
        })
        .catch(() => message.error('加载失败'))
        .finally(() => setLoading(false));
    }
  }, [id]);

  const handleSubmit = async (values: any) => {
    const payload = {
      name: values.name,
      steps: JSON.stringify(values.steps || []),
    };
    setSaving(true);
    try {
      if (isEdit) { await updateGroup(Number(id), payload); message.success('已更新'); }
      else { await createGroup(payload); message.success('已创建'); }
      navigate('/groups');
    } catch { message.error('保存失败'); } finally { setSaving(false); }
  };

  if (loading) return <Spin style={{ display: 'block', marginTop: 120 }} />;

  return (
    <Card title={isEdit ? '编辑 Group' : '新增 Group'} style={{ maxWidth: 900, margin: '0 auto' }}>
      <Form form={form} layout="vertical" onFinish={handleSubmit}>
        <Form.Item name="name" label="组名称" rules={[{ required: true, message: '请输入组名称' }]}>
          <Input placeholder="如: AllInOne_JPOS" />
        </Form.Item>

        <Form.List name="steps">
          {(steps, { add, remove }) => (
            <div style={{ border: '1px solid #d9d9d9', borderRadius: 8, padding: 16, marginBottom: 16 }}>
              <div style={{ fontWeight: 600, marginBottom: 12 }}>Steps 列表</div>
              {steps.map((field, index) => (
                <Space key={field.key} align="baseline" style={{ display: 'flex', marginBottom: 8 }}>
                  <Form.Item name={[index, 'name']} rules={[{ required: true, message: '必填' }]}>
                    <Input placeholder="Step 名称" style={{ width: 160 }} />
                  </Form.Item>
                  <Form.Item name={[index, 'path']} rules={[{ required: true, message: '请选择脚本' }]}>
                    <Cascader
                      placeholder="选择脚本文件"
                      style={{ width: 420 }}
                      options={buildScriptCascader(scriptFiles)}
                      displayRender={(labels) => labels[labels.length - 1]}
                      onChange={(values) => values && values[values.length - 1]}
                      changeOnSelect
                    />
                  </Form.Item>
                  <Button type="link" danger icon={<MinusCircleOutlined />} onClick={() => remove(index)} />
                </Space>
              ))}
              <Button type="dashed" icon={<PlusOutlined />} onClick={() => add({ name: '', path: '' })} block>
                添加 Step
              </Button>
            </div>
          )}
        </Form.List>

        <Form.Item>
          <Button type="primary" htmlType="submit" loading={saving} style={{ marginRight: 8 }}>
            {isEdit ? '更新' : '创建'}
          </Button>
          <Button onClick={() => navigate('/groups')}>取消</Button>
        </Form.Item>
      </Form>
    </Card>
  );
}
