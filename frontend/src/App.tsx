import { BrowserRouter, Routes, Route, Navigate, useNavigate, useLocation } from 'react-router-dom';
import { ConfigProvider, App as AntApp, Menu, Button } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import {
  AppstoreOutlined, NodeIndexOutlined, ShopOutlined, LaptopOutlined,
  CodeOutlined, FolderOutlined, DashboardOutlined, LogoutOutlined,
  FileTextOutlined, SettingOutlined, UserOutlined, SafetyOutlined,
  AuditOutlined,
} from '@ant-design/icons';
import { useAuthStore } from './store/auth';
import Login from './pages/Login';
import SSOCallback from './pages/SSOCallback';
import ModuleList from './pages/ModuleList';
import ModuleDetail from './pages/ModuleDetail';
import RuleList from './pages/RuleList';
import RuleDetail from './pages/RuleDetail';
import StoreList from './pages/StoreList';
import StoreDetail from './pages/StoreDetail';
import TillListPage from './pages/TillListPage';
import TillListDetail from './pages/TillListDetail';
import VarList from './pages/VarList';
import VarDetail from './pages/VarDetail';
import GroupList from './pages/GroupList';
import GroupDetail from './pages/GroupDetail';
import CsvGenerate from './pages/CsvGenerate';
import Dashboard from './pages/Dashboard';
import UserList from './pages/UserList';
import UserDetail from './pages/UserDetail';
import RoleList from './pages/RoleList';
import RoleDetail from './pages/RoleDetail';
import OperationLog from './pages/OperationLog';

type MenuItem = {
  key: string;
  icon?: React.ReactNode;
  label: string;
  children?: MenuItem[];
  type?: string;
};

const menuItems: MenuItem[] = [
  { key: '/', icon: <DashboardOutlined />, label: '仪表盘' },
  { key: '/stores', icon: <ShopOutlined />, label: 'Store' },
  { key: '/till-lists', icon: <LaptopOutlined />, label: 'TillList' },
    {
    key: 'config', icon: <AppstoreOutlined />, label: '配置管理',
    children: [
      { key: '/modules', icon: <AppstoreOutlined />, label: 'Module' },
      { key: '/groups', icon: <FolderOutlined />, label: 'Group' },
      { key: '/rules', icon: <NodeIndexOutlined />, label: 'Rule' },
      { key: '/vars', icon: <CodeOutlined />, label: 'Var' },
      { key: '/csv-generate', icon: <FileTextOutlined />, label: 'CSV 生成' },
    ],
  },
  {
    key: 'system', icon: <SettingOutlined />, label: '系统设置',
    children: [
      { key: '/system/users', icon: <UserOutlined />, label: '用户管理' },
      { key: '/system/roles', icon: <SafetyOutlined />, label: '角色管理' },
      { key: '/system/logs', icon: <AuditOutlined />, label: '操作日志' },
    ],
  },
];

function PrivateRoute({ children }: { children: React.ReactNode }) {
  const token = useAuthStore((s) => s.token);
  if (!token) {
    return <Navigate to="/login" replace />;
  }
  return <>{children}</>;
}

function findActiveKey(pathname: string): string[] {
  const parts = pathname.split('/').filter(Boolean);
  if (parts.length === 0) return ['/'];
  const top = '/' + parts[0];
  if (parts.length >= 2) {
    return [top, pathname];
  }
  return [top === '/' ? '/' : top];
}

function Layout({ children }: { children: React.ReactNode }) {
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout, token } = useAuthStore();
  const activeKeys = findActiveKey(location.pathname);

  if (!token) {
    return <>{children}</>;
  }

  return (
    <div style={{ display: 'flex', minHeight: '100vh' }}>
      <div style={{
        width: 200, flex: '0 0 auto', borderRight: '1px solid #f0f0f0',
        display: 'flex', flexDirection: 'column',
        position: 'fixed', top: 0, left: 0, height: '100vh', overflowY: 'auto',
      }}>
        <Menu
          mode="inline"
          selectedKeys={activeKeys}
          defaultOpenKeys={['config', 'system']}
          items={menuItems as any}
          onClick={({ key }) => navigate(key)}
          style={{ flex: 1, paddingTop: 16, borderRight: 0 }}
        />
        <div style={{ padding: '12px 16px', borderTop: '1px solid #f0f0f0' }}>
          <div style={{ fontSize: 12, color: '#999', marginBottom: 8 }}>{user?.name}</div>
          <Button size="small" icon={<LogoutOutlined />} onClick={() => { logout(); navigate('/login'); }} block>
            退出
          </Button>
        </div>
      </div>
      <div style={{ flex: 1, padding: 24, minWidth: 0, overflow: 'auto', marginLeft: 200 }}>
        {children}
      </div>
    </div>
  );
}

function HomePage() {
  return <Dashboard />;
}

export default function App() {
  const initialized = useAuthStore((s) => s.initialized);

  if (!initialized) {
    return (
      <ConfigProvider locale={zhCN}>
        <AntApp>
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh' }}>
            加载中...
          </div>
        </AntApp>
      </ConfigProvider>
    );
  }

  return (
    <ConfigProvider locale={zhCN}>
      <AntApp>
        <BrowserRouter>
          <Layout>
            <Routes>
              <Route path="/login" element={<Login />} />
              <Route path="/auth/sso/callback" element={<SSOCallback />} />
              <Route path="/" element={<PrivateRoute><HomePage /></PrivateRoute>} />
              <Route path="/groups" element={<PrivateRoute><GroupList /></PrivateRoute>} />
              <Route path="/groups/:id" element={<PrivateRoute><GroupDetail /></PrivateRoute>} />
              <Route path="/modules" element={<PrivateRoute><ModuleList /></PrivateRoute>} />
              <Route path="/modules/:id" element={<PrivateRoute><ModuleDetail /></PrivateRoute>} />
              <Route path="/rules" element={<PrivateRoute><RuleList /></PrivateRoute>} />
              <Route path="/rules/:id" element={<PrivateRoute><RuleDetail /></PrivateRoute>} />
              <Route path="/stores" element={<PrivateRoute><StoreList /></PrivateRoute>} />
              <Route path="/stores/:id" element={<PrivateRoute><StoreDetail /></PrivateRoute>} />
              <Route path="/till-lists" element={<PrivateRoute><TillListPage /></PrivateRoute>} />
              <Route path="/till-lists/:id" element={<PrivateRoute><TillListDetail /></PrivateRoute>} />
              <Route path="/vars" element={<PrivateRoute><VarList /></PrivateRoute>} />
              <Route path="/vars/:id" element={<PrivateRoute><VarDetail /></PrivateRoute>} />
              <Route path="/csv-generate" element={<PrivateRoute><CsvGenerate /></PrivateRoute>} />
              <Route path="/system/users" element={<PrivateRoute><UserList /></PrivateRoute>} />
              <Route path="/system/users/:id" element={<PrivateRoute><UserDetail /></PrivateRoute>} />
              <Route path="/system/roles" element={<PrivateRoute><RoleList /></PrivateRoute>} />
              <Route path="/system/roles/:id" element={<PrivateRoute><RoleDetail /></PrivateRoute>} />
              <Route path="/system/logs" element={<PrivateRoute><OperationLog /></PrivateRoute>} />
            </Routes>
          </Layout>
        </BrowserRouter>
      </AntApp>
    </ConfigProvider>
  );
}
