import { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { Spin, message } from 'antd';
import { useAuthStore } from '../store/auth';

export default function SSOCallback() {
  const navigate = useNavigate();
  const [params] = useSearchParams();
  const loginSSO = useAuthStore((s) => s.loginSSO);

  useEffect(() => {
    const code = params.get('code');
    const state = params.get('state');
    if (!code || !state) {
      message.error('SSO 回调参数缺失');
      navigate('/login');
      return;
    }
    loginSSO(code, state)
      .then(() => {
        message.success('SSO 登录成功');
        navigate('/');
      })
      .catch(() => {
        message.error('SSO 登录失败');
        navigate('/login');
      });
  }, []);

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh' }}>
      <Spin size="large" tip="SSO 登录中..." />
    </div>
  );
}
