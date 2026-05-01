export type AvatarType = 'custom' | 'default' | 'none';

export interface GetAvatarResponse {
  avatar_type: AvatarType;
  avatar_url: string;
}

export interface UploadAvatarResponse {
  upload_url: string;
}
