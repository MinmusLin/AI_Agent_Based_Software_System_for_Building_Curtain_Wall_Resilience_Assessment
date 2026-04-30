import { BarChartOutlined, FolderOpenOutlined, LogoutOutlined, UserOutlined } from '@ant-design/icons';
import { Avatar, Dropdown, Layout, Menu, type MenuProps } from 'antd';
import type { ReactElement } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';

import { useAuth } from '../contexts/AuthContext';

const { Header, Sider, Content } = Layout;

export function AppLayout(): ReactElement {
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout } = useAuth();
  const handleLogout = (): void => {
    void logout();
  };

  const goDashboard = (): void => {
    void navigate('/dashboard');
  };
  const goProjects = (): void => {
    void navigate('/projects');
  };

  const userItems: MenuProps['items'] = [
    {
      key: 'profile',
      className: 'cursor-default',
      label: (
        <div className="min-w-52 py-1">
          <div className="font-medium text-slate-900">{user?.name ?? user?.email}</div>
          <div className="mt-1 text-xs text-slate-500">{user?.email}</div>
        </div>
      ),
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout,
    },
  ];

  return (
    <Layout className="min-h-screen bg-slate-50">
      <Sider
        className="min-h-screen border-r border-slate-200"
        style={{ background: '#FFFFFF' }}
        theme="light"
        width={210}
      >
        <div className="flex h-16 items-center justify-center gap-2 border-b border-slate-200 px-3">
          <img alt="icw-logo" className="h-8 w-8 shrink-0" src="/icw-logo.png" />
          <div className="text-base font-semibold text-slate-900">Tongji University</div>
        </div>
        <Menu
          className="app-side-menu"
          items={[
            { key: '/dashboard', icon: <BarChartOutlined />, label: '工作台', onClick: goDashboard },
            { key: '/projects', icon: <FolderOpenOutlined />, label: '项目管理', onClick: goProjects },
          ]}
          mode="inline"
          selectedKeys={[location.pathname.startsWith('/projects') ? '/projects' : '/dashboard']}
          style={{ borderInlineEnd: 0 }}
        />
      </Sider>
      <Layout className="min-h-screen bg-slate-50">
        <Header
          className="flex h-16 items-center justify-between border-b border-slate-200 shadow-sm"
          style={{ background: '#FFFFFF', padding: '0 12px 0 24px' }}
        >
          <div className="text-lg font-medium text-slate-800">建筑幕墙韧性评估软件系统</div>
          <Dropdown menu={{ items: userItems }} placement="bottomRight" trigger={['click']}>
            <button className="flex items-center gap-2 rounded px-2 py-1 hover:bg-slate-100" type="button">
              <Avatar icon={<UserOutlined />} size={28} />
              <span className="text-sm text-slate-700">{user?.name ?? user?.email}</span>
            </button>
          </Dropdown>
        </Header>
        <Content className="min-h-[calc(100vh-64px)] bg-slate-50 p-6">
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
}
