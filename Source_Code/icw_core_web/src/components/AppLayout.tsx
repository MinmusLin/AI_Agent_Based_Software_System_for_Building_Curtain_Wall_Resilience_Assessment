import {
  BarChartOutlined,
  DeleteOutlined,
  FolderOpenOutlined,
  LockOutlined,
  LogoutOutlined,
  MailOutlined,
  UploadOutlined,
  UserOutlined,
} from '@ant-design/icons';
import { Avatar, Dropdown, Layout, Menu, type MenuProps } from 'antd';
import type { ReactElement } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';

import { useAuth } from '../contexts/AuthContext';
import { useUserAvatar } from '../hooks/useUserAvatar';
import { setPostLogoutRedirect } from '../utils/redirect';

const { Header, Sider, Content } = Layout;

interface UserMenuOptions {
  avatarURL: string;
  canDeleteAvatar: boolean;
  deleteCurrentAvatar: () => Promise<void>;
  email?: string;
  handleAvatarLoadError: () => boolean;
  handleChangePassword: () => void;
  handleFeedback: () => void;
  handleLogout: () => void;
  name?: string;
  openAvatarFileSelector: () => void;
}

function buildUserItems(options: UserMenuOptions): MenuProps['items'] {
  const displayName = options.name ?? options.email;
  const deleteAvatarItem: MenuProps['items'] = options.canDeleteAvatar
    ? [
        {
          key: 'delete-avatar',
          icon: <DeleteOutlined />,
          label: '删除头像',
          onClick: () => {
            void options.deleteCurrentAvatar();
          },
        },
      ]
    : [];

  return [
    {
      key: 'profile',
      className: 'app-user-profile-item',
      label: (
        <div className="flex min-w-56 items-center gap-3 py-1">
          <Avatar
            icon={<UserOutlined />}
            onError={options.handleAvatarLoadError}
            size={38}
            src={options.avatarURL || undefined}
          />
          <div className="min-w-0">
            <div className="truncate font-medium text-slate-900">{displayName}</div>
            <div className="mt-1 truncate text-xs text-slate-500">{options.email}</div>
          </div>
        </div>
      ),
    },
    {
      type: 'divider',
    },
    {
      key: 'upload-avatar',
      icon: <UploadOutlined />,
      label: '上传头像',
      onClick: options.openAvatarFileSelector,
    },
    ...deleteAvatarItem,
    {
      type: 'divider',
    },
    {
      key: 'change-password',
      icon: <LockOutlined />,
      label: '修改密码',
      onClick: options.handleChangePassword,
    },
    {
      key: 'feedback',
      icon: <MailOutlined />,
      label: '问题反馈',
      onClick: options.handleFeedback,
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: options.handleLogout,
    },
  ];
}

export function AppLayout(): ReactElement {
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout } = useAuth();
  const {
    avatarInputRef,
    avatarLoading,
    avatarURL,
    canDeleteAvatar,
    deleteCurrentAvatar,
    handleAvatarLoadError,
    handleAvatarFileChange,
    openAvatarFileSelector,
  } = useUserAvatar(user?.email);

  const handleLogout = (): void => {
    void logout();
  };
  const handleChangePassword = (): void => {
    setPostLogoutRedirect('/forget-password');
    void logout();
  };
  const handleFeedback = (): void => {
    window.location.href = 'mailto:minmuslin@outlook.com';
  };

  const goDashboard = (): void => {
    void navigate('/dashboard');
  };
  const goProjects = (): void => {
    void navigate('/projects');
  };

  const userItems = buildUserItems({
    avatarURL,
    canDeleteAvatar,
    deleteCurrentAvatar,
    email: user?.email,
    handleAvatarLoadError,
    handleChangePassword,
    handleFeedback,
    handleLogout,
    name: user?.name,
    openAvatarFileSelector,
  });

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
          <div className="text-base font-semibold text-slate-800">Tongji University</div>
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
            <button
              className="flex items-center gap-2 rounded px-2 py-1 hover:bg-slate-100"
              disabled={avatarLoading}
              type="button"
            >
              <Avatar icon={<UserOutlined />} onError={handleAvatarLoadError} size={28} src={avatarURL || undefined} />
              <span className="text-sm text-slate-700">{user?.name ?? user?.email}</span>
            </button>
          </Dropdown>
        </Header>
        <Content className="min-h-[calc(100vh-64px)] bg-slate-50 p-6">
          <Outlet />
        </Content>
      </Layout>
      <input
        accept="image/jpeg,image/png,image/webp"
        aria-hidden="true"
        className="hidden"
        onChange={handleAvatarFileChange}
        ref={avatarInputRef}
        tabIndex={-1}
        type="file"
      />
    </Layout>
  );
}
