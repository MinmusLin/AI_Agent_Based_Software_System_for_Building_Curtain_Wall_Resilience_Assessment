import type { User } from './common';

export interface LoginRequest {
  email: string;
  scene: 'password' | 'email';
  code: string;
}

export interface LoginResponse {
  access_token: string;
  expires_in: number;
  user: User;
}

export interface RegisterRequest {
  email: string;
  email_code: string;
  password: string;
  name: string;
}

export interface ResetPasswordRequest {
  email: string;
  email_code: string;
  new_password: string;
}

export interface SendEmailCodeRequest {
  email: string;
  scene: 'register' | 'login' | 'reset';
}
