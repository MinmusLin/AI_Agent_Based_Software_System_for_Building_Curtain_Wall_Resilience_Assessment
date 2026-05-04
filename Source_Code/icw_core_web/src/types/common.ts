export const AVATAR_TYPE_CUSTOM = 'custom';
export const AVATAR_TYPE_DEFAULT = 'default';
export const AVATAR_TYPE_NONE = 'none';

export type AvatarType = typeof AVATAR_TYPE_CUSTOM | typeof AVATAR_TYPE_DEFAULT | typeof AVATAR_TYPE_NONE;

export const LOGIN_SCENE_PASSWORD = 'password';
export const LOGIN_SCENE_EMAIL = 'email';

export type LoginScene = typeof LOGIN_SCENE_PASSWORD | typeof LOGIN_SCENE_EMAIL;

export const EMAIL_CODE_SCENE_REGISTER = 'register';
export const EMAIL_CODE_SCENE_LOGIN = 'login';
export const EMAIL_CODE_SCENE_RESET = 'reset';

export type EmailCodeScene =
  | typeof EMAIL_CODE_SCENE_REGISTER
  | typeof EMAIL_CODE_SCENE_LOGIN
  | typeof EMAIL_CODE_SCENE_RESET;

export const PROJECT_IMAGE_STATUS_PENDING = 'pending';
export const PROJECT_IMAGE_STATUS_UPLOADED = 'uploaded';
export const PROJECT_IMAGE_STATUS_FAILED = 'failed';

export type ProjectImageStatus =
  | typeof PROJECT_IMAGE_STATUS_PENDING
  | typeof PROJECT_IMAGE_STATUS_UPLOADED
  | typeof PROJECT_IMAGE_STATUS_FAILED;

export interface ApiEnvelope<T> {
  code: string;
  message: string;
  data: T;
}

export interface User {
  id: number;
  email: string;
  name: string;
}
