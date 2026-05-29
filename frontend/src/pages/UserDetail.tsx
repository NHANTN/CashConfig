import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Form, Input, Button, Card, Select, message } from 'antd';
import { getUser, createUser, updateUser } from '../api/user';
import { listRoles } from '../api/role';
import type { Role } from '../api/role';

export default function UserDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const [roles, setRoles] = useState<Role[]>([]);
  const [loading, setLoading] = useState(false);
  const isNew = !id || id === 'new';

  useEffect(() => {
    listRoles().then(setRoles);
    if (!isNew) {
      getUser(Number(id)).then((user) => {
        form.setFieldsValue(user);
      });
    }
  }, [id]);

  const handleSubmit = async (values: any) => {
    setLoading(true);
    try {
      if (isNew) {
        await createUser(values);
        message.success('创建成功');
      } else {
        const data = { ...values };
        if (!data.password) delete data.password;
        await updateUser(Number(id), data);
        message.success('更新成功');
      }
      navigate('/system/users');
    } catch { message.error('操作失败'); }
    finally { setLoading(false); }
  };

  return (
    <Card title={isNew ? '新增用户' : '编辑用户'}>
      <Form form={form} layout="vertical" onFinish={handleSubmit} style={{ maxWidth: 500 }}>
        <Form.Item name="username" label="用户名" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Form.Item name="name" label="姓名" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Form.Item name="password" label={isNew ? '密码' : '新密码（留空不修改）'} rules={isNew ? [{ required: true }] : []}>
          <Input.Password />
        </Form.Item>
        <Form.Item name="role_id" label="角色" rules={[{ required: true }]}>
          <Select options={roles.map(r => ({ label: r.name, value: r.id }))} />
        </Form.Item>
        <Button type="primary" htmlType="submit" loading={loading}>保存</Button>
        <Button style={{ marginLeft: 8 }} onClick={() => navigate('/system/users')}>取消</Button>
      </Form>
    </Card>
  );
}
