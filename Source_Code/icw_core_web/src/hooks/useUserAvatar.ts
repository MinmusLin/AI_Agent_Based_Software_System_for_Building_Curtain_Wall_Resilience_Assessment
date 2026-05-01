import { message } from 'antd';
import type { ChangeEvent, RefObject } from 'react';
import { useEffect, useRef, useState } from 'react';

import { createAvatarUpload, deleteAvatar, getAvatar } from '../api/avatar';
import type { AvatarType } from '../types/avatar';
import {
  AVATAR_MAX_SOURCE_SIZE,
  AVATAR_OUTPUT_CONTENT_TYPE,
  isAllowedAvatarFile,
  resizeAvatarToPng,
} from '../utils/avatar';

interface UseUserAvatarResult {
  avatarInputRef: RefObject<HTMLInputElement | null>;
  avatarLoading: boolean;
  avatarURL: string;
  canDeleteAvatar: boolean;
  deleteCurrentAvatar: () => Promise<void>;
  handleAvatarLoadError: () => boolean;
  handleAvatarFileChange: (event: ChangeEvent<HTMLInputElement>) => void;
  openAvatarFileSelector: () => void;
}

export function useUserAvatar(email?: string): UseUserAvatarResult {
  const avatarInputRef = useRef<HTMLInputElement | null>(null);
  const [avatarType, setAvatarType] = useState<AvatarType>('none');
  const [avatarURL, setAvatarURL] = useState('');
  const [avatarLoading, setAvatarLoading] = useState(false);

  useEffect((): (() => void) => {
    let active = true;
    async function loadAvatar(): Promise<void> {
      try {
        const result = await getAvatar();
        if (active) {
          setAvatarType(result.avatar_type);
          setAvatarURL(result.avatar_url);
        }
      } catch {
        if (active) {
          setAvatarType('none');
          setAvatarURL('');
        }
      }
    }
    if (email) {
      void loadAvatar();
    } else {
      setAvatarType('none');
      setAvatarURL('');
    }
    return () => {
      active = false;
    };
  }, [email]);

  const openAvatarFileSelector = (): void => {
    avatarInputRef.current?.click();
  };

  const handleAvatarFileChange = (event: ChangeEvent<HTMLInputElement>): void => {
    const file = event.target.files?.[0];
    event.target.value = '';
    if (!file) {
      return;
    }
    void uploadAvatar(file);
  };

  const uploadAvatar = async (file: File): Promise<void> => {
    if (!isAllowedAvatarFile(file)) {
      message.warning('头像仅支持 JPG、JPEG、PNG、WebP 格式');
      return;
    }
    if (file.size > AVATAR_MAX_SOURCE_SIZE) {
      message.warning('头像大小不能超过 2MB');
      return;
    }

    setAvatarLoading(true);
    try {
      const avatarBlob = await resizeAvatarToPng(file);
      const upload = await createAvatarUpload();
      if (!upload.upload_url) {
        throw new Error('avatar upload url is empty');
      }
      const uploadHeaders = new Headers();
      uploadHeaders.set('Content-Type', AVATAR_OUTPUT_CONTENT_TYPE);
      const uploadResp = await fetch(upload.upload_url, {
        body: avatarBlob,
        headers: uploadHeaders,
        method: 'PUT',
      });
      if (!uploadResp.ok) {
        throw new Error('avatar upload failed');
      }
      const avatar = await getAvatar();
      setAvatarType(avatar.avatar_type);
      setAvatarURL(avatar.avatar_url);
      message.success('头像已更新');
    } catch {
      message.error('头像上传失败，请稍后重试');
    } finally {
      setAvatarLoading(false);
    }
  };

  const refreshAvatarAfterLoadError = async (): Promise<void> => {
    try {
      const result = await getAvatar();
      setAvatarType(result.avatar_type);
      setAvatarURL(result.avatar_url);
    } catch {
      setAvatarType('none');
      setAvatarURL('');
    }
  };

  const handleAvatarLoadError = (): boolean => {
    void refreshAvatarAfterLoadError();
    return true;
  };

  const deleteCurrentAvatar = async (): Promise<void> => {
    setAvatarLoading(true);
    try {
      await deleteAvatar();
      message.success('头像已恢复默认');
    } catch {
      message.error('头像删除失败，请稍后重试');
    } finally {
      try {
        const result = await getAvatar();
        setAvatarType(result.avatar_type);
        setAvatarURL(result.avatar_url);
      } catch {
        setAvatarType('none');
        setAvatarURL('');
      }
      setAvatarLoading(false);
    }
  };

  return {
    avatarInputRef,
    avatarLoading,
    avatarURL,
    canDeleteAvatar: avatarType === 'custom',
    deleteCurrentAvatar,
    handleAvatarLoadError,
    handleAvatarFileChange,
    openAvatarFileSelector,
  };
}
