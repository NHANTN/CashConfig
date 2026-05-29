import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Form, Input, Button, Card, message, Switch } from 'antd';
import { getRole, createRole, updateRole } from '../api/role';

export default function RoleDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const isNew = !id || id === 'new';

  useEffect(() => {
    if (!isNew) {
      getRole(Number(id)).then((role) => {
        form.setFieldsValue(role);
      });
    }
  }, [id]);

  const handleSubmit = async (values: any) => {
    setLoading(true);
    try {
      if (isNew) {
        await createRole(values);
        message.success('创建成功');
      } else {
        await updateRole(Number(id), values);
        message.success('更新成功');
      }
      navigate('/system/roles');
    } catch { message.error('操作失败'); }
    finally { setLoading(false); }
  };

  return (
    <Card title={isNew ? '新增角色' : '编辑角色'}>
      <Form form={form} layout="vertical" onFinish={handleSubmit} style={{ maxWidth: 500 }}>
        <Form.Item name="name" label="角色名" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Form.Item name="code" label="编码" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Form.Item name="permissions" label="权限（JSON 数组）" rules={[{ required: true }]}>
          <Input.TextArea rows={4} placeholder='["门店管理","设备管理"]' />
        </Form.Item>
        {!isNew && (
          <Form.Item name="status" label="状态" valuePropName="checked" getValueFromEvent={(checked) => checked ? 1 : 0} getValueProps={(v) => ({ value: v === 1 })}>
            <Switch />
          </Form.Item>
        )}
        <Button type="primary" htmlType="submit" loading={loading}>保存</Button>
        <Button style={{ marginLeft: 8 }} onClick={() => navigate('/system/roles')}>取消</Button>
      </Form>
    </Card>
  );
}
