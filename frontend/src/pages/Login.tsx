import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Form, Input, Button, Card, message, Typography, Divider, Tabs } from 'antd';
import { UserOutlined, LockOutlined, SafetyCertificateOutlined } from '@ant-design/icons';
import { useAuthStore } from '../store/auth';
import * as authApi from '../api/auth';

const { Title } = Typography;

export default function Login() {
  const navigate = useNavigate();
  const login = useAuthStore((s) => s.login);
  const loginLDAP = useAuthStore((s) => s.loginLDAP);
  const [loading, setLoading] = useState(false);
  const [ldapLoading, setLdapLoading] = useState(false);

  const handleSubmit = async (values: { username: string; password: string }) => {
    setLoading(true);
    try {
      await login(values.username, values.password);
      message.success('登录成功');
      navigate('/');
    } catch {
      message.error('用户名或密码错误');
    } finally {
      setLoading(false);
    }
  };

  const handleLDAPLogin = async (values: { username: string; password: string }) => {
    setLdapLoading(true);
    try {
      await loginLDAP(values.username, values.password);
      message.success('LDAP 登录成功');
      navigate('/');
    } catch {
      message.error('LDAP 认证失败');
    } finally {
      setLdapLoading(false);
    }
  };

  const handleSSOLogin = async () => {
    try {
      const authURL = await authApi.getSSOLoginURL();
      window.location.href = authURL;
    } catch {
      message.error('SSO 登录失败');
    }
  };

  return (
    <div style={{
      display: 'flex', justifyContent: 'center', alignItems: 'center',
      minHeight: '100vh', background: '#f0f2f5',
    }}>
      <Card style={{ width: 420 }}>
        <Title level={3} style={{ textAlign: 'center', marginBottom: 32 }}>
          收银台配置管理平台
        </Title>
        <Tabs
          centered
          items={[
            {
              key: 'local',
              label: '本地账号',
              children: (
                <Form onFinish={handleSubmit} size="large">
                  <Form.Item name="username" rules={[{ required: true, message: '请输入用户名' }]}>
                    <Input prefix={<UserOutlined />} placeholder="用户名" />
                  </Form.Item>
                  <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}>
                    <Input.Password prefix={<LockOutlined />} placeholder="密码" />
                  </Form.Item>
                  <Form.Item>
                    <Button type="primary" htmlType="submit" loading={loading} block>
                      登录
                    </Button>
                  </Form.Item>
                </Form>
              ),
            },
            {
              key: 'ldap',
              label: 'LDAP',
              children: (
                <Form onFinish={handleLDAPLogin} size="large">
                  <Form.Item name="username" rules={[{ required: true, message: '请输入用户名' }]}>
                    <Input prefix={<UserOutlined />} placeholder="LDAP 用户名" />
                  </Form.Item>
                  <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}>
                    <Input.Password prefix={<LockOutlined />} placeholder="密码" />
                  </Form.Item>
                  <Form.Item>
                    <Button type="primary" htmlType="submit" loading={ldapLoading} block>
                      LDAP 登录
                    </Button>
                  </Form.Item>
                </Form>
              ),
            },
          ]}
        />
        <Divider />
        <Button icon={<SafetyCertificateOutlined />} block onClick={handleSSOLogin}>
          SSO 单点登录
        </Button>
      </Card>
    </div>
  );
}
