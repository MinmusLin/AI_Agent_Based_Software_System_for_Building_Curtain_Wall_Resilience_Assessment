import { Spin } from 'antd';
import type { ReactElement } from 'react';
import { useEffect } from 'react';
import { Navigate, Outlet, useLocation } from 'react-router-dom';

import { useAuth } from '@/contexts/AuthContext';
import { clearPostLogoutRedirect, getPostLogoutRedirect } from '@/utils/redirect';

// 登录态初始化期间的全屏加载态
// 应用启动时会先尝试用 HttpOnly Cookie 中的 Refresh Token 恢复 Access Token，在结果明确之前不渲染页面，避免页面短暂闪烁
function FullPageSpin(): ReactElement {
  return (
    <div className="flex min-h-screen items-center justify-center">
      <Spin />
    </div>
  );
}

// 登录态路由守卫
export function ProtectedRoute(): ReactElement {
  const { status } = useAuth();
  const location = useLocation();
  const postLogoutRedirect = getPostLogoutRedirect();

  useEffect(() => {
    if (status !== 'authenticated' && postLogoutRedirect) {
      clearPostLogoutRedirect();
    }
  }, [postLogoutRedirect, status]);

  // 登录态恢复中，展示全屏加载页
  if (status === 'initializing') {
    return <FullPageSpin />;
  }

  // 未登录用户跳转到登录页
  if (status !== 'authenticated') {
    return <Navigate replace state={{ from: location.pathname }} to={postLogoutRedirect || '/login'} />;
  }

  // 已登录用户渲染当前路由页面
  return <Outlet />;
}

// 访客态路由守卫
export function GuestRoute(): ReactElement {
  const { status } = useAuth();

  // 登录态恢复中，展示全屏加载页
  if (status === 'initializing') {
    return <FullPageSpin />;
  }

  // 已登录用户跳转到系统首页
  if (status === 'authenticated') {
    return <Navigate replace to="/dashboard" />;
  }

  // 未登录用户渲染当前访客页面
  return <Outlet />;
}
