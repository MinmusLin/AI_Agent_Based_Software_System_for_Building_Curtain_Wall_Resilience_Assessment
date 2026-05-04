import type { AvatarType } from './common';

export interface GetAvatarResponse {
  avatar_type: AvatarType;
  avatar_url: string;
}

export interface UploadAvatarResponse {
  upload_url: string;
}
