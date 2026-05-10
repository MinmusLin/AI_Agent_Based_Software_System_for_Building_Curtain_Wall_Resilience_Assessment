import type { ApiEnvelope } from '@/constants/common';
import type {
  LoginRequest,
  LoginResponse,
  MeResponse,
  RegisterRequest,
  ResetPasswordRequest,
  SendEmailCodeRequest,
  SendEmailCodeResponse,
} from '@/gen/core/api/auth';
import type { User } from '@/gen/core/common';

import { http } from './http';

// 登录
// @router /auth/login [POST]
export async function login(payload: LoginRequest): Promise<LoginResponse> {
  const { data } = await http.post<ApiEnvelope<LoginResponse>>('/auth/login', payload);
  return data.data;
}

// 登出
// @router /auth/logout [POST]
export async function logout(): Promise<void> {
  await http.post('/auth/logout');
}

// 获取用户信息
// @router /auth/me [GET]
export async function getMe(): Promise<User> {
  const { data } = await http.get<ApiEnvelope<MeResponse>>('/auth/me');
  if (!data.data.user) {
    throw new Error('user is empty');
  }
  return data.data.user;
}

// 刷新 Token
// @router /auth/refresh [POST]
// 该接口由 Axios 响应拦截器统一调用，不作为页面级业务 API 暴露

// 注册
// @router /auth/register [POST]
export async function register(payload: RegisterRequest): Promise<void> {
  await http.post<ApiEnvelope<Record<string, never>>>('/auth/register', payload);
}

// 重置密码
// @router /auth/reset-password [POST]
export async function resetPassword(payload: ResetPasswordRequest): Promise<void> {
  await http.post('/auth/reset-password', payload);
}

// 发送邮箱验证码
// @router /auth/send-email-code [POST]
export async function sendEmailCode(payload: SendEmailCodeRequest): Promise<SendEmailCodeResponse> {
  const { data } = await http.post<ApiEnvelope<SendEmailCodeResponse>>('/auth/send-email-code', payload);
  return data.data;
}
