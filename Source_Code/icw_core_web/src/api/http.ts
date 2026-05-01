import type { AxiosError, InternalAxiosRequestConfig } from 'axios';
import axios from 'axios';

import type { ApiEnvelope } from '@/types/common';
import { getRequiredEnv } from '@/utils/env';

const API_BASE_URL = getRequiredEnv('VITE_API_BASE_URL');
const REFRESH_MAX_ATTEMPTS = 3;
const REFRESH_RETRY_DELAY_MS = 300;
const DEFAULT_ERROR_MESSAGE = '服务暂时不可用，请稍后重试';

// Access Token 保存在前端内存中，不保存在 localStorage 或 sessionStorage 中
// 页面刷新后 Access Token 会丢失，因此应用初始化时通过 HttpOnly Cookie 中的 Refresh Token 获取新的 Access Token
let accessToken = '';

// 多个接口同时因为 Access Token 过期返回 401 时，只允许发起一次 /auth/refresh 请求
// 其他失败请求复用同一个 Promise，避免并发刷新导致 Refresh Token 轮换冲突
let refreshPromise: Promise<string | null> | null = null;

// http.ts 不直接依赖 React Context
// AuthContext 初始化时会注入这两个回调，用来同步新的 Access Token 或清理登录态
let onTokenChange: ((token: string) => void) | null = null;
let onUnauthorized: (() => void) | null = null;

interface ErrorResponse {
  message?: string;
}

export const http = axios.create({
  // API_BASE_URL 从环境变量文件中读取
  baseURL: API_BASE_URL,
  // Refresh Token 由 API 写入 HttpOnly Cookie
  // withCredentials=true 确保浏览器调用 /auth/refresh 和 /auth/logout 接口时自动携带 Cookie
  withCredentials: true,
});

// 获取接口错误文案，API 层会把用户可感知的错误写入 response.data.message，网络错误或非标准错误使用统一兜底文案
export function getErrorMessage(error: unknown): string {
  if (axios.isAxiosError<ErrorResponse>(error)) {
    return error.response?.data.message ?? DEFAULT_ERROR_MESSAGE;
  }
  return DEFAULT_ERROR_MESSAGE;
}

// AuthContext 在登录、登出、刷新 Token 时调用该函数，维护 http.ts 内部的 Access Token
export function setAccessToken(token: string): void {
  accessToken = token;
}

// AuthContext 通过 hooks 接收 http.ts 的鉴权事件
// onTokenChange：刷新成功后把新 Access Token 同步回 React 状态
// onUnauthorized：Refresh Token 无效或刷新失败达到上限后，通知 React 清空登录态
export function setAuthHooks(hooks: { onTokenChange: (token: string) => void; onUnauthorized: () => void }): void {
  const { onTokenChange: nextOnTokenChange, onUnauthorized: nextOnUnauthorized } = hooks;
  onTokenChange = nextOnTokenChange;
  onUnauthorized = nextOnUnauthorized;
}

// 请求拦截器：对所有普通业务请求自动添加 Authorization Header
// Refresh Token 由浏览器 Cookie 自动携带，不在这里处理
http.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  if (accessToken) {
    config.headers.Authorization = `Bearer ${accessToken}`;
  }
  return config;
});

// 响应拦截器：当登录态接口返回 401 时，自动调用 /auth/refresh 刷新 Access Token
// 刷新成功后使用新 Access Token 重放原请求，刷新失败后通知 AuthContext 清空登录态
http.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const original = error.config as (InternalAxiosRequestConfig & { _retried?: boolean }) | undefined;

    // 只有登录态接口返回 401 时才尝试刷新 Access Token
    // 返回非 401、已经重试过、或者 /auth/refresh 请求失败，直接把错误交回业务代码处理
    if (!original || error.response?.status !== 401 || original._retried || original.url?.includes('/auth/refresh')) {
      return Promise.reject(error);
    }

    // 尝试刷新 Access Token
    original._retried = true;
    const token = await refreshAccessToken();
    if (!token) {
      // refreshAccessToken 返回空表示 Refresh Token 已不可用，或连续 3 次刷新失败
      // 此时通知 AuthContext 清空登录态，页面恢复匿名状态
      onUnauthorized?.();
      return Promise.reject(error);
    }

    // 刷新成功后，使用新的 Access Token 重放刚才失败的请求
    original.headers.Authorization = `Bearer ${token}`;
    return http(original);
  },
);

// 对外暴露的刷新接口
export async function refreshAccessToken(): Promise<string | null> {
  refreshPromise ??= requestRefreshAccessToken().finally(() => {
    refreshPromise = null;
  });
  return refreshPromise;
}

// 调用 /auth/refresh 获取新的 Access Token
// Refresh Token 保存在 HttpOnly Cookie 中，JavaScript 不能读取，浏览器会自动随请求发送给 API
async function requestRefreshAccessToken(): Promise<string | null> {
  for (let attempt = 1; attempt <= REFRESH_MAX_ATTEMPTS; attempt += 1) {
    try {
      const res = await http.post<ApiEnvelope<{ access_token: string }>>('/auth/refresh');
      const token = res.data.data.access_token;
      accessToken = token;
      onTokenChange?.(token);
      return token;
    } catch (error: unknown) {
      // 网络错误或 API 5xx 视为基础服务临时异常，最多重试 3 次
      // API 返回 401 说明 Refresh Token 已过期、被吊销或不合法，不再重试
      if (!shouldRetryRefresh(error) || attempt === REFRESH_MAX_ATTEMPTS) {
        return null;
      }
      await delay(REFRESH_RETRY_DELAY_MS * attempt);
    }
  }
  return null;
}

// 刷新只重试网络错误和 5xx，4xx 通常代表凭证问题或请求参数问题，重试没有意义
function shouldRetryRefresh(error: unknown): boolean {
  if (!axios.isAxiosError(error)) {
    return false;
  }
  if (!error.response) {
    return true;
  }
  return error.response.status >= 500;
}

// 简单线性退避：第 1/2/3 次失败之间分别等待 300ms 和 600ms
function delay(ms: number): Promise<void> {
  return new Promise((resolve) => {
    window.setTimeout(resolve, ms);
  });
}
