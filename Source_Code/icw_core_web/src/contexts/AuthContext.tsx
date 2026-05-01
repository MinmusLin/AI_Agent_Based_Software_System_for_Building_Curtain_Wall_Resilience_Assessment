import { isAxiosError } from 'axios';
import type { ReactElement, ReactNode } from 'react';
import { createContext, useContext, useEffect, useMemo, useState } from 'react';

import { getMe, login as loginRequest, logout as logoutRequest } from '@/api/auth';
import { refreshAccessToken, setAccessToken, setAuthHooks } from '@/api/http';
import type { LoginRequest } from '@/types/auth';
import type { User } from '@/types/common';

// 登录态分为三个阶段：
// initializing：应用启动中，正在尝试用 Refresh Token 恢复登录态
// authenticated：已有可用 Access Token，可以访问登录态页面
// anonymous：未登录或 Refresh Token 已不可用，只能访问匿名页面
type AuthStatus = 'initializing' | 'authenticated' | 'anonymous';

interface AuthContextValue {
  status: AuthStatus;
  user: User | null;
  accessToken: string;
  login: (payload: LoginRequest) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }): ReactElement {
  // token/user/status 是 React 侧登录态
  // http.ts 内部也保存一份 Access Token，用于 Axios 请求拦截器自动添加 Authorization Header
  const [status, setStatus] = useState<AuthStatus>('initializing');
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState('');

  // 统一清理前端登录态
  // 这里无法直接删除 HttpOnly Refresh Token Cookie，Cookie 清理由 API 的 /auth/logout 或 /auth/refresh 401 负责。
  function clearAuth(): void {
    setToken('');
    setAccessToken('');
    setUser(null);
    setStatus('anonymous');
  }

  useEffect(() => {
    // 把 React 状态更新能力注入给 http.ts
    // 当 Axios 拦截器刷新出新 Access Token 时，通过 onTokenChange 同步到 Context
    // 当刷新失败达到退出条件时，会通过 onUnauthorized 恢复匿名状态
    setAuthHooks({
      onTokenChange: setToken,
      onUnauthorized: clearAuth,
    });

    // Access Token 只保存在前端内存中，应用启动时 Access Token 一定为空的
    // 因此先尝试调用 /auth/refresh 让浏览器自动带上 HttpOnly Cookie 中的 Refresh Token
    async function bootstrapAuth(): Promise<void> {
      try {
        const newToken = await refreshAccessToken();
        if (!newToken) {
          // 无可用 Refresh Token，或者刷新失败达到上限，退出登录态
          clearAuth();
          return;
        }

        setAccessToken(newToken);
        setToken(newToken);

        try {
          // 携带刚刚拿到的新 Access Token 调用 /auth/me 登录态接口，获取当前用户信息
          const currentUser = await getMe();
          setUser(currentUser);
        } catch (error: unknown) {
          if (isUnauthorizedError(error)) {
            // /auth/me 接口返回 401 说明 Access Token 不可用，退出登录态
            clearAuth();
            return;
          }
          // /auth/me 接口返回非 401 通常是临时服务异常，保留已恢复的 Access Token 和 authenticated 状态，避免因为 /auth/me 短暂不可用让用户强制退出登录态
        }

        setStatus('authenticated');
      } catch {
        // 初始化流程发生非预期异常时，恢复匿名状态，避免页面长期停留在 initializing 阶段
        clearAuth();
      }
    }

    void bootstrapAuth();
  }, []);

  const value = useMemo<AuthContextValue>(
    () => ({
      status,
      user,
      accessToken: token,
      async login(payload) {
        // 登录成功后 API 返回 Access Token，并通过 HttpOnly Cookie 写入 Refresh Token
        // 前端只保存 Access Token 和用户信息，后续请求由 Axios 拦截器自动添加 Authorization Header
        const result = await loginRequest(payload);
        setAccessToken(result.access_token);
        setToken(result.access_token);
        setUser(result.user);
        setStatus('authenticated');
      },
      async logout() {
        // 登出时会通知 API 吊销 Refresh Token 并清理 HttpOnly Cookie
        // 即使登出接口失败，前端也清理本地登录态，避免用户继续停留在登录态页面
        await logoutRequest().catch(() => undefined);
        setAccessToken('');
        setToken('');
        setUser(null);
        setStatus('anonymous');
      },
    }),
    [status, token, user],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

// 只有 401 表示登录态不可用，5xx 或网络错误可能是基础服务临时异常，不应该直接清空登录态
function isUnauthorizedError(error: unknown): boolean {
  return isAxiosError(error) && error.response?.status === 401;
}

// 统一读取登录态 Context，调用方要求必须在 AuthProvider 下
export function useAuth(): AuthContextValue {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
}
