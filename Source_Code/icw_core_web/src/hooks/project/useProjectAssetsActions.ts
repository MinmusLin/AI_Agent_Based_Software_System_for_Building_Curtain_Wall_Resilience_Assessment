import type { Dispatch, SetStateAction } from 'react';
import { useCallback, useEffect, useRef, useState } from 'react';

import { getErrorMessage } from '@/api/http';
import {
  createProjectGroup,
  deleteProjectGroup,
  deleteProjectImage,
  getProjectAssets,
  reportProjectImage,
  updateProjectGroup,
  uploadProjectImage,
} from '@/api/project/assets';
import { advanceProject } from '@/api/project/core';
import type { Project, ProjectGroup } from '@/gen/core/api/common';
import {
  type ProjectImage,
  ProjectImageStatus_Value,
  ProjectProgress_Value,
  type UploadProjectImageResult,
} from '@/gen/core/common';
import type { PreparedImageUpload } from '@/utils/assetsStage';
import {
  appendImagesToGroup,
  EMPTY_ITEMS_COUNT,
  formatProjectGroupName,
  normalizeGroups,
  prepareUploadImage,
  putPresignedObject,
  removeGroup,
  removeImages,
  replaceGroup,
  uploadItemFromPreparedImage,
} from '@/utils/assetsStage';
import {
  isAllowedProjectAssetImageFile,
  PROJECT_ASSET_IMAGE_OUTPUT_CONTENT_TYPE,
  PROJECT_ASSET_THUMBNAIL_OUTPUT_CONTENT_TYPE,
} from '@/utils/images';

interface UseProjectAssetsActionsParams {
  loading: boolean;
  onError: (message: string) => void;
  onProgressChange: (progress: ProjectProgress_Value) => void;
  onProjectChange: (project: Project) => void;
  project: Project;
  projectId: string;
  selectedProgress: ProjectProgress_Value;
}

interface UseProjectAssetsActionsResult {
  advancing: boolean;
  assetsLoading: boolean;
  canComplete: boolean;
  creatingGroup: boolean;
  deletingGroupIds: Set<string>;
  deletingImageUuids: Set<string>;
  editingGroupId: string;
  editingGroupName: string;
  groups: ProjectGroup[];
  handleComplete: () => Promise<void>;
  handleCreateGroup: () => Promise<void>;
  handleDeleteGroup: (groupId: string) => Promise<void>;
  handleDeleteImage: (imageUuid: string) => Promise<void>;
  handleUploadFiles: (groupId: string, files: File[]) => Promise<void>;
  loadAssets: (options?: RefreshOptions) => Promise<void>;
  readOnly: boolean;
  saveEditingGroup: () => Promise<void>;
  savingGroupId: string;
  setEditingGroupName: Dispatch<SetStateAction<string>>;
  setGroups: Dispatch<SetStateAction<ProjectGroup[]>>;
  startEditGroup: (group: ProjectGroup) => void;
  uploadingImages: boolean;
}

interface RefreshOptions {
  silent?: boolean;
}

export function useProjectAssetsActions({
  loading,
  onError,
  onProgressChange,
  onProjectChange,
  project,
  projectId,
  selectedProgress,
}: UseProjectAssetsActionsParams): UseProjectAssetsActionsResult {
  const [groups, setGroups] = useState<ProjectGroup[]>([]);
  const [assetsLoading, setAssetsLoading] = useState(true);
  const [creatingGroup, setCreatingGroup] = useState(false);
  const [advancing, setAdvancing] = useState(false);
  const [editingGroupId, setEditingGroupId] = useState('');
  const [editingGroupName, setEditingGroupName] = useState('');
  const [savingGroupId, setSavingGroupId] = useState('');
  const [deletingGroupIds, setDeletingGroupIds] = useState<Set<string>>(() => new Set());
  const [deletingImageUuids, setDeletingImageUuids] = useState<Set<string>>(() => new Set());
  const [uploadingImages, setUploadingImages] = useState(false);
  const editingGroupIdRef = useRef('');
  const savingGroupIdsRef = useRef<Set<string>>(new Set());
  const uploadingImagesRef = useRef(false);

  const readOnly =
    loading ||
    project.progress > ProjectProgress_Value.ProfileFinished ||
    selectedProgress !== ProjectProgress_Value.ProfileFinished;
  const canComplete =
    !loading &&
    project.progress === ProjectProgress_Value.ProfileFinished &&
    selectedProgress === ProjectProgress_Value.ProfileFinished;

  const loadAssets = useCallback(
    async (options: RefreshOptions = {}): Promise<void> => {
      if (projectId === '') {
        return;
      }

      if (!options.silent) {
        setAssetsLoading(true);
      }
      try {
        const data = await getProjectAssets(projectId);
        setGroups(normalizeGroups(data.groups));
      } catch (error: unknown) {
        if (!options.silent) {
          onError(getErrorMessage(error));
        }
      } finally {
        if (!options.silent) {
          setAssetsLoading(false);
        }
      }
    },
    [onError, projectId],
  );

  const handleCreateGroup = useCallback(async (): Promise<void> => {
    setCreatingGroup(true);
    try {
      const data = await createProjectGroup(projectId);
      const { group } = data;
      if (!group) {
        throw new Error('group is empty');
      }
      setGroups((currentGroups) => [...currentGroups, group]);
    } catch (error: unknown) {
      onError(getErrorMessage(error));
    } finally {
      setCreatingGroup(false);
    }
  }, [onError, projectId]);

  const handleDeleteGroup = useCallback(
    async (groupId: string): Promise<void> => {
      setDeletingGroupIds((currentIds) => new Set(currentIds).add(groupId));
      try {
        await deleteProjectGroup({
          group_id: groupId,
          project_id: projectId,
        });
        setGroups((currentGroups) => removeGroup(currentGroups, groupId));
      } catch (error: unknown) {
        onError(getErrorMessage(error));
      } finally {
        setDeletingGroupIds((currentIds) => {
          const nextIds = new Set(currentIds);
          nextIds.delete(groupId);
          return nextIds;
        });
      }
    },
    [onError, projectId],
  );

  const startEditGroup = useCallback((group: ProjectGroup): void => {
    editingGroupIdRef.current = group.id;
    setEditingGroupId(group.id);
    setEditingGroupName(group.name);
  }, []);

  const saveEditingGroup = useCallback(async (): Promise<void> => {
    const targetGroupId = editingGroupId;
    if (targetGroupId === '' || savingGroupIdsRef.current.has(targetGroupId)) {
      return;
    }

    const currentGroup = groups.find((group) => group.id === targetGroupId);
    if (!currentGroup) {
      if (editingGroupIdRef.current === targetGroupId) {
        editingGroupIdRef.current = '';
        setEditingGroupId('');
      }
      return;
    }

    const nextName = formatProjectGroupName(editingGroupName);
    if (editingGroupIdRef.current === targetGroupId) {
      setEditingGroupName(nextName);
    }
    setGroups((currentGroups) =>
      replaceGroup(currentGroups, {
        ...currentGroup,
        name: nextName,
      }),
    );
    savingGroupIdsRef.current.add(targetGroupId);
    setSavingGroupId(targetGroupId);
    try {
      const data = await updateProjectGroup({
        group_id: targetGroupId,
        name: nextName,
        project_id: projectId,
      });
      const { group } = data;
      if (!group) {
        throw new Error('group is empty');
      }
      setGroups((currentGroups) => replaceGroup(currentGroups, group));
    } catch (error: unknown) {
      if (editingGroupIdRef.current === targetGroupId) {
        setEditingGroupName(currentGroup.name);
      }
      setGroups((currentGroups) => replaceGroup(currentGroups, currentGroup));
      onError(getErrorMessage(error));
    } finally {
      savingGroupIdsRef.current.delete(targetGroupId);
      if (editingGroupIdRef.current === targetGroupId) {
        editingGroupIdRef.current = '';
        setEditingGroupId('');
      }
      setSavingGroupId((currentGroupId) => (currentGroupId === targetGroupId ? '' : currentGroupId));
    }
  }, [editingGroupId, editingGroupName, groups, onError, projectId]);

  const handleUploadFiles = useCallback(
    async (groupId: string, files: File[]): Promise<void> => {
      if (uploadingImagesRef.current) {
        return;
      }

      const validFiles = files.filter((file) => isAllowedProjectAssetImageFile(file));
      if (validFiles.length !== files.length) {
        onError('请上传 JPG、PNG 或 WebP 格式的图片');
      }
      if (validFiles.length === EMPTY_ITEMS_COUNT) {
        return;
      }

      uploadingImagesRef.current = true;
      setUploadingImages(true);
      let backgroundUploads: Promise<void>[] = [];
      try {
        const preparedImages = await Promise.all(validFiles.map((file) => prepareUploadImage(file)));
        const uploadItems = preparedImages.map((image) => uploadItemFromPreparedImage(image));
        const data = await uploadProjectImage({
          group_id: groupId,
          images: uploadItems,
          project_id: projectId,
        });
        const pendingImages = data.images.map((item) => item.image);
        if (pendingImages.some((image) => !image)) {
          throw new Error('image is empty');
        }
        setGroups((currentGroups) => appendImagesToGroup(currentGroups, groupId, pendingImages as ProjectImage[]));
        backgroundUploads = data.images.map(async (result, index): Promise<void> => {
          await uploadOneImage(projectId, result, preparedImages[index]);
        });
      } catch (error: unknown) {
        onError(getErrorMessage(error));
      } finally {
        uploadingImagesRef.current = false;
        setUploadingImages(false);
      }
      if (backgroundUploads.length > EMPTY_ITEMS_COUNT) {
        void Promise.all(backgroundUploads).catch((error: unknown) => {
          onError(getErrorMessage(error));
        });
      }
    },
    [onError, projectId],
  );

  const handleDeleteImage = useCallback(
    async (imageUuid: string): Promise<void> => {
      setDeletingImageUuids((currentIds) => new Set(currentIds).add(imageUuid));
      try {
        await deleteProjectImage({
          image_uuids: [imageUuid],
          project_id: projectId,
        });
        setGroups((currentGroups) => removeImages(currentGroups, [imageUuid]));
      } catch (error: unknown) {
        onError(getErrorMessage(error));
      } finally {
        setDeletingImageUuids((currentIds) => {
          const nextIds = new Set(currentIds);
          nextIds.delete(imageUuid);
          return nextIds;
        });
      }
    },
    [onError, projectId],
  );

  const handleComplete = useCallback(async (): Promise<void> => {
    setAdvancing(true);
    try {
      await advanceProject({
        from_progress: project.progress,
        project_id: projectId,
        to_progress: ProjectProgress_Value.AssetsFinished,
      });
      onProjectChange({
        ...project,
        progress: ProjectProgress_Value.AssetsFinished,
      });
      onProgressChange(ProjectProgress_Value.AssetsFinished);
    } catch (error: unknown) {
      onError(getErrorMessage(error));
    } finally {
      setAdvancing(false);
    }
  }, [onError, onProgressChange, onProjectChange, project, projectId]);

  useEffect(() => {
    void loadAssets();
  }, [loadAssets]);

  return {
    advancing,
    assetsLoading,
    canComplete,
    creatingGroup,
    deletingGroupIds,
    deletingImageUuids,
    editingGroupId,
    editingGroupName,
    groups,
    handleComplete,
    handleCreateGroup,
    handleDeleteGroup,
    handleDeleteImage,
    handleUploadFiles,
    loadAssets,
    readOnly,
    saveEditingGroup,
    savingGroupId,
    setEditingGroupName,
    setGroups,
    startEditGroup,
    uploadingImages,
  };
}

async function uploadOneImage(
  projectId: string,
  result: UploadProjectImageResult,
  preparedImage: PreparedImageUpload,
): Promise<void> {
  if (!result.image) {
    throw new Error('image is empty');
  }
  let reportStatus: ProjectImageStatus_Value = ProjectImageStatus_Value.Uploaded;
  try {
    await putPresignedObject(
      result.original_upload_url,
      preparedImage.originalBlob,
      PROJECT_ASSET_IMAGE_OUTPUT_CONTENT_TYPE,
    );
    await putPresignedObject(
      result.thumbnail_upload_url,
      preparedImage.thumbnailBlob,
      PROJECT_ASSET_THUMBNAIL_OUTPUT_CONTENT_TYPE,
    );
  } catch {
    reportStatus = ProjectImageStatus_Value.Failed;
  }

  await reportProjectImage({
    image_uuid: result.image.uuid,
    project_id: projectId,
    status: reportStatus,
  });
}
