import type { ApiEnvelope } from '@/constants/common';
import type { GetAvatarResponse, UploadAvatarResponse } from '@/gen/core/api/user';

import { http } from './http';

// 获取用户头像
// @router /user/avatar [GET]
export async function getAvatar(): Promise<GetAvatarResponse> {
  const { data } = await http.get<ApiEnvelope<GetAvatarResponse>>('/user/avatar');
  return data.data;
}

// 上传用户自定义头像
// @router /user/avatar [POST]
export async function uploadAvatar(): Promise<UploadAvatarResponse> {
  const { data } = await http.post<ApiEnvelope<UploadAvatarResponse>>('/user/avatar');
  return data.data;
}

// 删除用户自定义头像
// @router /user/avatar [DELETE]
export async function deleteAvatar(): Promise<void> {
  await http.delete<ApiEnvelope<Record<string, never>>>('/user/avatar');
}
