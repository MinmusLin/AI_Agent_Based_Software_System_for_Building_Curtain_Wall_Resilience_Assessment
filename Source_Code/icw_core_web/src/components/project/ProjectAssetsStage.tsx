import {
  CaretDownOutlined,
  DeleteOutlined,
  DragOutlined,
  EyeOutlined,
  LoadingOutlined,
  PictureOutlined,
  PlusOutlined,
  ReloadOutlined,
  StepForwardOutlined,
  StopOutlined,
  UploadOutlined,
} from '@ant-design/icons';
import { Button, Empty, Input, message, Modal, Spin, Tooltip } from 'antd';
import type { ChangeEvent, DragEvent, ReactElement } from 'react';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { getErrorMessage } from '@/api/http';
import {
  createProjectGroup,
  deleteProjectGroup,
  deleteProjectImage,
  getProjectAssets,
  getProjectImageOriginal,
  moveProjectGroup,
  moveProjectImage,
  reportProjectImage,
  updateProjectGroup,
  uploadProjectImage,
} from '@/api/project/assets';
import { advanceProject } from '@/api/project/core';
import { createSocketTicket, setupAssetsWebSocket } from '@/api/socket';
import type { ProjectImageStatus, ProjectProgress } from '@/types/common';
import {
  PROJECT_EVENT_TYPE_IMAGE_STATUS_CHANGED,
  PROJECT_IMAGE_STATUS_FAILED,
  PROJECT_IMAGE_STATUS_PENDING,
  PROJECT_IMAGE_STATUS_UPLOADED,
  PROJECT_PROGRESS_ASSETS_FINISHED,
  PROJECT_PROGRESS_PROFILE_FINISHED,
} from '@/types/common';
import type {
  ProjectGroup,
  ProjectImage,
  UploadProjectImageItem,
  UploadProjectImageResult,
} from '@/types/project/assets';
import type { Project } from '@/types/project/core';
import type { ProjectImageStatusChangedMessage } from '@/types/socket';
import { formatDateTime } from '@/utils/datetime';
import {
  buildProjectAssetImageBlobs,
  isAllowedProjectAssetImageFile,
  PROJECT_ASSET_IMAGE_OUTPUT_CONTENT_TYPE,
  PROJECT_ASSET_THUMBNAIL_OUTPUT_CONTENT_TYPE,
} from '@/utils/images';

const EMPTY_GROUP_ID = '';
const EMPTY_ITEMS_COUNT = 0;
const FIRST_INDEX = 0;
const KILOBYTE_SIZE_BYTES = 1024;
const NEXT_INDEX_OFFSET = 1;
const NOT_FOUND_INDEX = -1;
const UPLOAD_ACCEPT = 'image/jpeg,image/png,image/webp';
const WEBSOCKET_RECONNECT_DELAY_MS = 2000;
const JSON_FORMAT_INDENT = 2;

interface ProjectAssetsStageProps {
  loading?: boolean;
  onProgressChange: (progress: ProjectProgress) => void;
  onProjectChange: (project: Project) => void;
  project: Project;
  projectId: string;
  selectedProgress: ProjectProgress;
}

interface DraggingImage {
  imageUuid: string;
  sourceGroupId: string;
}

interface ViewerImage {
  groupId: string;
  image: ProjectImage;
}

interface ImageViewerState {
  imageUuid: string;
  loading: boolean;
  originalUrl: string;
}

interface PreparedImageUpload {
  height: number;
  metadata: string;
  originalBlob: Blob;
  sourceFile: File;
  thumbnailBlob: Blob;
  width: number;
}

function sortProjectImages(images: ProjectImage[]): ProjectImage[] {
  return [...images].sort((left, right) => {
    const leftTime = new Date(left.created_at).getTime();
    const rightTime = new Date(right.created_at).getTime();
    return leftTime - rightTime;
  });
}

function normalizeGroups(groups: ProjectGroup[]): ProjectGroup[] {
  return groups.map((group) => ({
    ...group,
    images: sortProjectImages(group.images),
  }));
}

function replaceGroup(groups: ProjectGroup[], nextGroup: ProjectGroup): ProjectGroup[] {
  return groups.map((group) => {
    if (group.id !== nextGroup.id) {
      return group;
    }
    return {
      ...group,
      ...nextGroup,
      images: nextGroup.images.length > EMPTY_ITEMS_COUNT ? sortProjectImages(nextGroup.images) : group.images,
    };
  });
}

function removeGroup(groups: ProjectGroup[], groupId: string): ProjectGroup[] {
  return groups.filter((group) => group.id !== groupId);
}

function appendImagesToGroup(groups: ProjectGroup[], groupId: string, images: ProjectImage[]): ProjectGroup[] {
  return groups.map((group) => {
    if (group.id !== groupId) {
      return group;
    }
    return {
      ...group,
      images: sortProjectImages([...group.images, ...images]),
    };
  });
}

function removeImages(groups: ProjectGroup[], imageUuids: string[]): ProjectGroup[] {
  const imageUuidSet = new Set(imageUuids);
  return groups.map((group) => ({
    ...group,
    images: group.images.filter((image) => !imageUuidSet.has(image.uuid)),
  }));
}

function replaceImage(groups: ProjectGroup[], nextImage: ProjectImage): ProjectGroup[] {
  const targetGroupIndex = groups.findIndex((group) => group.images.some((image) => image.uuid === nextImage.uuid));
  if (targetGroupIndex < FIRST_INDEX) {
    return groups;
  }

  return groups.map((group, index) => {
    if (index !== targetGroupIndex) {
      return group;
    }
    return {
      ...group,
      images: sortProjectImages(
        group.images.map((image) => {
          if (image.uuid !== nextImage.uuid) {
            return image;
          }
          return nextImage;
        }),
      ),
    };
  });
}

function moveImagesToGroup(groups: ProjectGroup[], targetGroupId: string, movedImages: ProjectImage[]): ProjectGroup[] {
  const movedImageUuidSet = new Set(movedImages.map((image) => image.uuid));
  return groups.map((group) => {
    const remainingImages = group.images.filter((image) => !movedImageUuidSet.has(image.uuid));
    if (group.id !== targetGroupId) {
      return {
        ...group,
        images: remainingImages,
      };
    }
    return {
      ...group,
      images: sortProjectImages([...remainingImages, ...movedImages]),
    };
  });
}

function flattenUploadedImages(groups: ProjectGroup[]): ViewerImage[] {
  return groups.flatMap((group) =>
    group.images
      .filter((image) => image.status === PROJECT_IMAGE_STATUS_UPLOADED)
      .map((image) => ({
        groupId: group.id,
        image,
      })),
  );
}

function buildUploadMetadata(file: File): string {
  return JSON.stringify({
    source_content_type: file.type,
    source_file_name: file.name,
    source_size_bytes: file.size,
  });
}

function formatProjectImageMetadata(metadata: string): string {
  const rawMetadata = metadata.trim();
  if (rawMetadata === '') {
    return '{}';
  }

  try {
    const parsedMetadata = JSON.parse(rawMetadata) as unknown;
    return JSON.stringify(parsedMetadata, null, JSON_FORMAT_INDENT);
  } catch {
    return metadata;
  }
}

function ImageUnavailableIcon(): ReactElement {
  return (
    <span className="relative inline-flex size-8 items-center justify-center text-slate-400">
      <PictureOutlined className="text-3xl" />
      <StopOutlined className="absolute -right-1 -bottom-1 rounded-full bg-white text-sm text-red-500" />
    </span>
  );
}

function imageStatusNode(status: ProjectImageStatus): ReactElement {
  switch (status) {
    case PROJECT_IMAGE_STATUS_PENDING:
      return <LoadingOutlined className="text-xl text-slate-400" />;
    case PROJECT_IMAGE_STATUS_FAILED:
      return (
        <Tooltip title="图像已失效">
          <ImageUnavailableIcon />
        </Tooltip>
      );
    case PROJECT_IMAGE_STATUS_UPLOADED:
      return (
        <Tooltip title="图像已失效">
          <ImageUnavailableIcon />
        </Tooltip>
      );
    default:
      return (
        <Tooltip title="图像状态异常">
          <ImageUnavailableIcon />
        </Tooltip>
      );
  }
}

function groupMovePosition(groups: ProjectGroup[], sourceGroupId: string, targetGroupId: string): ProjectGroup[] {
  const sourceIndex = groups.findIndex((group) => group.id === sourceGroupId);
  const targetIndex = groups.findIndex((group) => group.id === targetGroupId);
  if (sourceIndex < FIRST_INDEX || targetIndex < FIRST_INDEX || sourceIndex === targetIndex) {
    return groups;
  }

  const nextGroups = [...groups];
  const [sourceGroup] = nextGroups.splice(sourceIndex, NEXT_INDEX_OFFSET);
  nextGroups.splice(targetIndex, EMPTY_ITEMS_COUNT, sourceGroup);
  return nextGroups;
}

function groupMovePayload(
  projectId: string,
  groupId: string,
  nextGroups: ProjectGroup[],
): {
  move_to_first: boolean;
  move_to_last: boolean;
  next_group_id: string;
  previous_group_id: string;
  project_id: string;
  group_id: string;
} {
  const nextIndex = nextGroups.findIndex((group) => group.id === groupId);
  const previousGroupId = nextIndex > FIRST_INDEX ? nextGroups[nextIndex - NEXT_INDEX_OFFSET].id : EMPTY_GROUP_ID;
  const nextGroupId =
    nextIndex >= FIRST_INDEX && nextIndex < nextGroups.length - NEXT_INDEX_OFFSET
      ? nextGroups[nextIndex + NEXT_INDEX_OFFSET].id
      : EMPTY_GROUP_ID;
  return {
    group_id: groupId,
    move_to_first: nextIndex === FIRST_INDEX,
    move_to_last: nextIndex === nextGroups.length - NEXT_INDEX_OFFSET,
    next_group_id: nextGroupId,
    previous_group_id: previousGroupId,
    project_id: projectId,
  };
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null;
}

function isProjectImageStatus(value: unknown): value is ProjectImageStatus {
  return (
    value === PROJECT_IMAGE_STATUS_PENDING ||
    value === PROJECT_IMAGE_STATUS_UPLOADED ||
    value === PROJECT_IMAGE_STATUS_FAILED
  );
}

function isProjectImage(value: unknown): value is ProjectImage {
  if (!isRecord(value)) {
    return false;
  }
  return (
    typeof value.uuid === 'string' &&
    typeof value.file_name === 'string' &&
    typeof value.content_type === 'string' &&
    typeof value.size_bytes === 'number' &&
    typeof value.width === 'number' &&
    typeof value.height === 'number' &&
    typeof value.metadata === 'string' &&
    isProjectImageStatus(value.status) &&
    typeof value.thumbnail_url === 'string' &&
    typeof value.uploaded_at === 'string' &&
    typeof value.created_at === 'string'
  );
}

function parseProjectImageStatusChangedMessage(data: unknown): ProjectImageStatusChangedMessage | null {
  if (typeof data !== 'string') {
    return null;
  }

  let value: unknown;
  try {
    value = JSON.parse(data) as unknown;
  } catch {
    return null;
  }

  if (
    !isRecord(value) ||
    value.type !== PROJECT_EVENT_TYPE_IMAGE_STATUS_CHANGED ||
    typeof value.project_id !== 'string'
  ) {
    return null;
  }
  const { image } = value;
  if (!isProjectImage(image)) {
    return null;
  }
  return {
    image,
    project_id: value.project_id,
    type: PROJECT_EVENT_TYPE_IMAGE_STATUS_CHANGED,
  };
}

async function putPresignedObject(uploadUrl: string, blob: Blob, contentType: string): Promise<void> {
  if (uploadUrl.trim() === '') {
    throw new Error('presigned upload url is empty');
  }

  const headers = new Headers();
  headers.set('Content-Type', contentType);
  const response = await fetch(uploadUrl, {
    body: blob,
    headers,
    method: 'PUT',
  });
  if (!response.ok) {
    throw new Error('presigned upload failed');
  }
}

async function prepareUploadImage(file: File): Promise<PreparedImageUpload> {
  const blobs = await buildProjectAssetImageBlobs(file);
  return {
    height: blobs.height,
    metadata: buildUploadMetadata(file),
    originalBlob: blobs.originalBlob,
    sourceFile: file,
    thumbnailBlob: blobs.thumbnailBlob,
    width: blobs.width,
  };
}

function uploadItemFromPreparedImage(image: PreparedImageUpload): UploadProjectImageItem {
  return {
    content_type: PROJECT_ASSET_IMAGE_OUTPUT_CONTENT_TYPE,
    file_name: image.sourceFile.name,
    height: image.height,
    metadata: image.metadata,
    size_bytes: image.originalBlob.size,
    width: image.width,
  };
}

export function ProjectAssetsStage({
  loading = false,
  onProgressChange,
  onProjectChange,
  project,
  projectId,
  selectedProgress,
}: ProjectAssetsStageProps): ReactElement {
  const [messageApi, contextHolder] = message.useMessage();
  const uploadInputRefs = useRef(new Map<string, HTMLInputElement>());
  const [groups, setGroups] = useState<ProjectGroup[]>([]);
  const [assetsLoading, setAssetsLoading] = useState(true);
  const [creatingGroup, setCreatingGroup] = useState(false);
  const [advancing, setAdvancing] = useState(false);
  const [collapsedGroupIds, setCollapsedGroupIds] = useState<Set<string>>(() => new Set());
  const [editingGroupId, setEditingGroupId] = useState('');
  const [editingGroupName, setEditingGroupName] = useState('');
  const [savingGroupId, setSavingGroupId] = useState('');
  const [deletingGroupIds, setDeletingGroupIds] = useState<Set<string>>(() => new Set());
  const [deletingImageUuids, setDeletingImageUuids] = useState<Set<string>>(() => new Set());
  const [draggingGroupId, setDraggingGroupId] = useState('');
  const [draggingImage, setDraggingImage] = useState<DraggingImage | null>(null);
  const [viewer, setViewer] = useState<ImageViewerState | null>(null);

  const readOnly =
    loading ||
    project.progress > PROJECT_PROGRESS_PROFILE_FINISHED ||
    selectedProgress !== PROJECT_PROGRESS_PROFILE_FINISHED;
  const canComplete =
    !loading &&
    project.progress === PROJECT_PROGRESS_PROFILE_FINISHED &&
    selectedProgress === PROJECT_PROGRESS_PROFILE_FINISHED;
  const uploadedImages = useMemo(() => flattenUploadedImages(groups), [groups]);

  const handleSocketMessage = useCallback(
    (data: unknown): void => {
      const socketMessage = parseProjectImageStatusChangedMessage(data);
      if (socketMessage?.project_id !== projectId) {
        return;
      }
      setGroups((currentGroups) => replaceImage(currentGroups, socketMessage.image));
    },
    [projectId],
  );

  const loadAssets = useCallback(async (): Promise<void> => {
    if (projectId === '') {
      return;
    }

    setAssetsLoading(true);
    try {
      const data = await getProjectAssets(projectId);
      setGroups(normalizeGroups(data.groups));
    } catch (error: unknown) {
      void messageApi.error(getErrorMessage(error));
    } finally {
      setAssetsLoading(false);
    }
  }, [messageApi, projectId]);

  const handleCreateGroup = useCallback(async (): Promise<void> => {
    setCreatingGroup(true);
    try {
      const data = await createProjectGroup(projectId);
      setGroups((currentGroups) => [...currentGroups, data.group]);
    } catch (error: unknown) {
      void messageApi.error(getErrorMessage(error));
    } finally {
      setCreatingGroup(false);
    }
  }, [messageApi, projectId]);

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
        void messageApi.error(getErrorMessage(error));
      } finally {
        setDeletingGroupIds((currentIds) => {
          const nextIds = new Set(currentIds);
          nextIds.delete(groupId);
          return nextIds;
        });
      }
    },
    [messageApi, projectId],
  );

  const startEditGroup = useCallback((group: ProjectGroup): void => {
    setEditingGroupId(group.id);
    setEditingGroupName(group.name);
  }, []);

  const saveEditingGroup = useCallback(async (): Promise<void> => {
    if (editingGroupId === '' || savingGroupId !== '') {
      return;
    }

    const currentGroup = groups.find((group) => group.id === editingGroupId);
    if (!currentGroup) {
      setEditingGroupId('');
      return;
    }

    const nextName = editingGroupName.trim();
    if (nextName === currentGroup.name) {
      setEditingGroupId('');
      return;
    }

    setSavingGroupId(editingGroupId);
    try {
      const data = await updateProjectGroup({
        group_id: editingGroupId,
        name: nextName,
        project_id: projectId,
      });
      setGroups((currentGroups) => replaceGroup(currentGroups, data.group));
      setEditingGroupId('');
    } catch (error: unknown) {
      void messageApi.error(getErrorMessage(error));
    } finally {
      setSavingGroupId('');
    }
  }, [editingGroupId, editingGroupName, groups, messageApi, projectId, savingGroupId]);

  const handleGroupDrop = useCallback(
    async (targetGroupId: string): Promise<void> => {
      if (readOnly || draggingGroupId === '' || draggingGroupId === targetGroupId) {
        return;
      }

      const nextGroups = groupMovePosition(groups, draggingGroupId, targetGroupId);
      if (nextGroups === groups) {
        return;
      }

      try {
        const payload = groupMovePayload(projectId, draggingGroupId, nextGroups);
        const data = await moveProjectGroup(payload);
        setGroups(replaceGroup(nextGroups, data.group));
      } catch (error: unknown) {
        void messageApi.error(getErrorMessage(error));
      } finally {
        setDraggingGroupId('');
      }
    },
    [draggingGroupId, groups, messageApi, projectId, readOnly],
  );

  const handleImageDrop = useCallback(
    async (targetGroupId: string): Promise<void> => {
      if (readOnly || !draggingImage || draggingImage.sourceGroupId === targetGroupId) {
        return;
      }

      try {
        const data = await moveProjectImage({
          image_uuids: [draggingImage.imageUuid],
          project_id: projectId,
          target_group_id: targetGroupId,
        });
        setGroups((currentGroups) => moveImagesToGroup(currentGroups, targetGroupId, data.images));
      } catch (error: unknown) {
        void messageApi.error(getErrorMessage(error));
      } finally {
        setDraggingImage(null);
      }
    },
    [draggingImage, messageApi, projectId, readOnly],
  );

  const uploadOneImage = useCallback(
    async (result: UploadProjectImageResult, preparedImage: PreparedImageUpload): Promise<void> => {
      let reportStatus: ProjectImageStatus = PROJECT_IMAGE_STATUS_UPLOADED;
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
        reportStatus = PROJECT_IMAGE_STATUS_FAILED;
      }

      await reportProjectImage({
        image_uuid: result.image.uuid,
        project_id: projectId,
        status: reportStatus,
      });
    },
    [projectId],
  );

  const handleUploadFiles = useCallback(
    async (groupId: string, files: File[]): Promise<void> => {
      const validFiles = files.filter((file) => isAllowedProjectAssetImageFile(file));
      if (validFiles.length !== files.length) {
        void messageApi.error('请上传 JPG、PNG 或 WebP 格式的图片');
      }
      if (validFiles.length === EMPTY_ITEMS_COUNT) {
        return;
      }

      try {
        const preparedImages = await Promise.all(validFiles.map((file) => prepareUploadImage(file)));
        const uploadItems = preparedImages.map((image) => uploadItemFromPreparedImage(image));
        const data = await uploadProjectImage({
          group_id: groupId,
          images: uploadItems,
          project_id: projectId,
        });
        const pendingImages = data.images.map((item) => item.image);
        setGroups((currentGroups) => appendImagesToGroup(currentGroups, groupId, pendingImages));

        await Promise.all(
          data.images.map(async (result, index): Promise<void> => {
            const preparedImage = preparedImages[index];
            await uploadOneImage(result, preparedImage);
          }),
        );
      } catch (error: unknown) {
        void messageApi.error(getErrorMessage(error));
      }
    },
    [messageApi, projectId, uploadOneImage],
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
        void messageApi.error(getErrorMessage(error));
      } finally {
        setDeletingImageUuids((currentIds) => {
          const nextIds = new Set(currentIds);
          nextIds.delete(imageUuid);
          return nextIds;
        });
      }
    },
    [messageApi, projectId],
  );

  const openImageViewer = useCallback(
    async (imageUuid: string): Promise<void> => {
      setViewer({
        imageUuid,
        loading: true,
        originalUrl: '',
      });
      try {
        const data = await getProjectImageOriginal(projectId, imageUuid);
        setViewer({
          imageUuid,
          loading: false,
          originalUrl: data.original_url,
        });
      } catch (error: unknown) {
        setViewer(null);
        void messageApi.error(getErrorMessage(error));
      }
    },
    [messageApi, projectId],
  );

  const handleComplete = useCallback(async (): Promise<void> => {
    setAdvancing(true);
    try {
      await advanceProject({
        from_progress: project.progress,
        project_id: projectId,
        to_progress: PROJECT_PROGRESS_ASSETS_FINISHED,
      });
      onProjectChange({
        ...project,
        progress: PROJECT_PROGRESS_ASSETS_FINISHED,
      });
      onProgressChange(PROJECT_PROGRESS_ASSETS_FINISHED);
    } catch (error: unknown) {
      void messageApi.error(getErrorMessage(error));
    } finally {
      setAdvancing(false);
    }
  }, [messageApi, onProgressChange, onProjectChange, project, projectId]);

  useEffect(() => {
    void loadAssets();
  }, [loadAssets]);

  useEffect(() => {
    if (projectId === '') {
      return undefined;
    }

    let closed = false;
    let socket: WebSocket | null = null;
    let reconnectTimer: number | null = null;

    function clearReconnectTimer(): void {
      if (reconnectTimer === null) {
        return;
      }
      window.clearTimeout(reconnectTimer);
      reconnectTimer = null;
    }

    function scheduleReconnect(): void {
      if (closed || reconnectTimer !== null) {
        return;
      }
      reconnectTimer = window.setTimeout(() => {
        reconnectTimer = null;
        void connect();
      }, WEBSOCKET_RECONNECT_DELAY_MS);
    }

    async function connect(): Promise<void> {
      try {
        const ticket = await createSocketTicket(projectId);
        if (closed) {
          return;
        }
        socket = new WebSocket(setupAssetsWebSocket(projectId, ticket.ticket));
        socket.onmessage = (event: MessageEvent<unknown>): void => {
          handleSocketMessage(event.data);
        };
        socket.onclose = (): void => {
          if (!closed) {
            scheduleReconnect();
          }
        };
        socket.onerror = (): void => {
          socket?.close();
        };
      } catch {
        scheduleReconnect();
      }
    }

    void connect();

    return () => {
      closed = true;
      clearReconnectTimer();
      socket?.close();
    };
  }, [handleSocketMessage, projectId]);

  const viewerIndex = useMemo(() => {
    if (!viewer) {
      return NOT_FOUND_INDEX;
    }
    return uploadedImages.findIndex((item) => item.image.uuid === viewer.imageUuid);
  }, [uploadedImages, viewer]);
  const viewerImage = viewerIndex >= FIRST_INDEX ? uploadedImages[viewerIndex] : null;

  return (
    <div className="flex min-h-0 flex-1 flex-col rounded-lg border border-slate-200 bg-white p-5">
      {contextHolder}
      <div className="mb-4 flex items-center justify-between gap-4">
        <div>
          <h2 className="text-base font-semibold text-slate-900">图像资产构建</h2>
          <p className="mt-1 text-sm text-slate-500">按立面或区域组织幕墙图像，上传完成后进入智能检测阶段。</p>
        </div>
        <div className="flex shrink-0 items-center gap-3">
          <Button disabled={assetsLoading} icon={<ReloadOutlined />} onClick={() => void loadAssets()}>
            刷新
          </Button>
          <Button
            disabled={readOnly}
            icon={<PlusOutlined />}
            loading={creatingGroup}
            onClick={() => void handleCreateGroup()}
            type="primary"
          >
            新建图像组
          </Button>
          {canComplete ? (
            <Button
              icon={<StepForwardOutlined />}
              loading={advancing}
              onClick={() => void handleComplete()}
              type="primary"
            >
              完成并进入下一步
            </Button>
          ) : null}
        </div>
      </div>
      <div className="min-h-0 flex-1 overflow-y-auto pr-1">
        {assetsLoading ? (
          <div className="flex h-full items-center justify-center">
            <Spin description="正在加载图像资产" />
          </div>
        ) : groups.length > EMPTY_ITEMS_COUNT ? (
          <div className="space-y-4">
            {groups.map((group) => {
              const canDeleteGroup = !readOnly && groups.length > NEXT_INDEX_OFFSET;
              const collapsed = collapsedGroupIds.has(group.id);
              const deletingGroup = deletingGroupIds.has(group.id);
              return (
                <section
                  className="rounded-lg border border-slate-200 bg-slate-50"
                  key={group.id}
                  onDragOver={(event: DragEvent<HTMLElement>) => {
                    if (draggingGroupId !== '' || draggingImage) {
                      event.preventDefault();
                    }
                  }}
                  onDrop={(event: DragEvent<HTMLElement>) => {
                    event.preventDefault();
                    if (draggingImage) {
                      void handleImageDrop(group.id);
                      return;
                    }
                    void handleGroupDrop(group.id);
                  }}
                >
                  <div
                    className="group flex items-center gap-3 border-b border-slate-200 bg-white px-4 py-3"
                    draggable={!readOnly}
                    onDragEnd={() => {
                      setDraggingGroupId('');
                    }}
                    onDragStart={(event: DragEvent<HTMLDivElement>) => {
                      if (readOnly) {
                        return;
                      }
                      setDraggingGroupId(group.id);
                      event.dataTransfer.effectAllowed = 'move';
                    }}
                  >
                    <Button
                      aria-label={collapsed ? '展开图像组' : '折叠图像组'}
                      icon={
                        <CaretDownOutlined
                          className={`transition-transform duration-300 ${collapsed ? '-rotate-90' : 'rotate-0'}`}
                        />
                      }
                      onClick={() => {
                        setCollapsedGroupIds((currentIds) => {
                          const nextIds = new Set(currentIds);
                          if (nextIds.has(group.id)) {
                            nextIds.delete(group.id);
                          } else {
                            nextIds.add(group.id);
                          }
                          return nextIds;
                        });
                      }}
                      shape="circle"
                      size="small"
                      type="text"
                    />
                    <DragOutlined className="cursor-grab text-slate-400" />
                    {editingGroupId === group.id ? (
                      <Input
                        autoFocus
                        className="h-7 max-w-72 text-sm"
                        disabled={savingGroupId === group.id}
                        maxLength={32}
                        onBlur={() => void saveEditingGroup()}
                        onChange={(event: ChangeEvent<HTMLInputElement>) => {
                          setEditingGroupName(event.target.value);
                        }}
                        onKeyDown={(event) => {
                          if (event.key === 'Escape') {
                            setEditingGroupId('');
                          }
                        }}
                        onPressEnter={() => void saveEditingGroup()}
                        size="small"
                        value={editingGroupName}
                      />
                    ) : (
                      <button
                        className="h-7 min-w-0 truncate rounded px-2 text-left text-sm font-semibold leading-7 text-slate-900 hover:bg-slate-100 disabled:cursor-default disabled:hover:bg-transparent"
                        disabled={readOnly}
                        onClick={() => {
                          startEditGroup(group);
                        }}
                        type="button"
                      >
                        {group.name}
                      </button>
                    )}
                    <span className="rounded bg-slate-100 px-2 py-0.5 text-xs text-slate-500">
                      {group.images.length} 张图像
                    </span>
                    {canDeleteGroup ? (
                      <Button
                        aria-label="删除图像组"
                        className="ml-auto opacity-0 transition-opacity duration-200 group-hover:opacity-100"
                        danger
                        icon={<DeleteOutlined />}
                        loading={deletingGroup}
                        onClick={() => void handleDeleteGroup(group.id)}
                        shape="circle"
                        size="small"
                        type="text"
                      />
                    ) : null}
                  </div>
                  <div
                    className={`grid transition-[grid-template-rows,opacity] duration-300 ease-in-out ${
                      collapsed ? 'grid-rows-[0fr] opacity-0' : 'grid-rows-[1fr] opacity-100'
                    }`}
                  >
                    <div className="min-h-0 overflow-hidden">
                      <div className="grid grid-cols-4 gap-3 p-4 md:grid-cols-8 xl:grid-cols-12">
                        {group.images.map((image) => {
                          const deletingImage = deletingImageUuids.has(image.uuid);
                          const imageUploaded = image.status === PROJECT_IMAGE_STATUS_UPLOADED;
                          const imageFailed = image.status === PROJECT_IMAGE_STATUS_FAILED;
                          const imageReady = imageUploaded && image.thumbnail_url !== '';
                          const imageActionVisible = imageUploaded || (!readOnly && imageFailed);
                          return (
                            <div
                              className="group/image relative aspect-square overflow-hidden rounded-lg border border-slate-200 bg-white"
                              draggable={!readOnly && imageUploaded}
                              key={image.uuid}
                              onDragEnd={() => {
                                setDraggingImage(null);
                              }}
                              onDragStart={(event: DragEvent<HTMLDivElement>) => {
                                if (readOnly || !imageUploaded) {
                                  return;
                                }
                                event.stopPropagation();
                                event.dataTransfer.effectAllowed = 'move';
                                setDraggingImage({
                                  imageUuid: image.uuid,
                                  sourceGroupId: group.id,
                                });
                              }}
                            >
                              {imageReady ? (
                                <img
                                  alt={image.file_name}
                                  className="size-full object-cover"
                                  src={image.thumbnail_url}
                                />
                              ) : (
                                <div className="flex size-full flex-col items-center justify-center gap-2 px-3 text-center text-xs text-slate-500">
                                  {imageStatusNode(image.status)}
                                </div>
                              )}
                              {imageActionVisible ? (
                                <div className="absolute inset-0 flex items-center justify-center gap-2 bg-slate-950/0 opacity-0 transition duration-200 group-hover/image:bg-slate-950/35 group-hover/image:opacity-100">
                                  {imageUploaded ? (
                                    <Tooltip title="查看原图">
                                      <Button
                                        aria-label="查看原图"
                                        icon={<EyeOutlined />}
                                        onClick={() => void openImageViewer(image.uuid)}
                                        shape="circle"
                                      />
                                    </Tooltip>
                                  ) : null}
                                  {!readOnly ? (
                                    <Tooltip title="删除图像">
                                      <Button
                                        aria-label="删除图像"
                                        danger
                                        icon={<DeleteOutlined />}
                                        loading={deletingImage}
                                        onClick={() => void handleDeleteImage(image.uuid)}
                                        shape="circle"
                                      />
                                    </Tooltip>
                                  ) : null}
                                </div>
                              ) : null}
                            </div>
                          );
                        })}
                        {!readOnly ? (
                          <button
                            className="flex aspect-square flex-col items-center justify-center gap-2 rounded-lg border border-dashed border-slate-300 bg-white text-sm text-slate-500 transition duration-200 hover:border-[#1677FF] hover:text-[#1677FF]"
                            onClick={() => {
                              uploadInputRefs.current.get(group.id)?.click();
                            }}
                            type="button"
                          >
                            <UploadOutlined className="text-2xl" />
                            <input
                              accept={UPLOAD_ACCEPT}
                              className="hidden"
                              multiple
                              onChange={(event: ChangeEvent<HTMLInputElement>) => {
                                const selectedFiles = Array.from(event.target.files ?? []);
                                event.target.value = '';
                                void handleUploadFiles(group.id, selectedFiles);
                              }}
                              ref={(node) => {
                                if (node) {
                                  uploadInputRefs.current.set(group.id, node);
                                } else {
                                  uploadInputRefs.current.delete(group.id);
                                }
                              }}
                              type="file"
                            />
                          </button>
                        ) : null}
                      </div>
                    </div>
                  </div>
                </section>
              );
            })}
          </div>
        ) : (
          <div className="flex h-full items-center justify-center rounded-lg border border-dashed border-slate-200">
            <Empty description="暂无图像组，请新建图像组后上传图像" />
          </div>
        )}
      </div>
      <Modal
        centered
        footer={null}
        onCancel={() => {
          setViewer(null);
        }}
        open={viewer !== null}
        title={viewerImage?.image.file_name ?? '图像详情'}
        width={1040}
      >
        {viewer ? (
          <div className="grid min-h-[520px] grid-cols-[minmax(0,1fr)_320px] gap-5">
            <div className="relative flex min-h-0 items-center justify-center overflow-hidden rounded-lg bg-slate-100">
              {viewer.loading ? (
                <Spin />
              ) : (
                <img
                  alt={viewerImage?.image.file_name ?? '项目原图'}
                  className="max-h-[520px] max-w-full object-contain"
                  src={viewer.originalUrl}
                />
              )}
            </div>
            <div className="flex min-h-0 flex-col text-sm text-slate-600">
              <div className="space-y-4">
                <div>
                  <div className="mb-1 font-medium text-slate-900">图像 ID</div>
                  <div className="break-all">{viewerImage?.image.uuid ?? '-'}</div>
                </div>
                <div>
                  <div className="mb-1 font-medium text-slate-900">文件名称</div>
                  <div className="break-all">{viewerImage?.image.file_name ?? '-'}</div>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <div className="mb-1 font-medium text-slate-900">图像尺寸（宽 × 高）</div>
                    <div>
                      {viewerImage ? `${String(viewerImage.image.width)} × ${String(viewerImage.image.height)}` : '-'}
                    </div>
                  </div>
                  <div>
                    <div className="mb-1 font-medium text-slate-900">文件大小</div>
                    <div>
                      {viewerImage
                        ? `${String(Math.round(viewerImage.image.size_bytes / KILOBYTE_SIZE_BYTES))} KB`
                        : '-'}
                    </div>
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <div className="mb-1 font-medium text-slate-900">上传开始时间</div>
                    <div>{viewerImage ? formatDateTime(viewerImage.image.created_at, true) : '-'}</div>
                  </div>
                  <div>
                    <div className="mb-1 font-medium text-slate-900">上传完成时间</div>
                    <div>{viewerImage ? formatDateTime(viewerImage.image.uploaded_at, true) : '-'}</div>
                  </div>
                </div>
              </div>
              <div className="mt-4 flex min-h-0 flex-1 flex-col">
                <div className="mb-1 font-medium text-slate-900">元数据</div>
                <pre className="min-h-0 flex-1 overflow-auto whitespace-pre-wrap break-all rounded bg-slate-50 p-3 text-xs leading-5">
                  {viewerImage ? formatProjectImageMetadata(viewerImage.image.metadata) : '{}'}
                </pre>
              </div>
              <div className="flex justify-end gap-2 pt-4">
                <Button
                  disabled={viewerIndex <= FIRST_INDEX}
                  onClick={() => {
                    if (viewerIndex > FIRST_INDEX) {
                      void openImageViewer(uploadedImages[viewerIndex - NEXT_INDEX_OFFSET].image.uuid);
                    }
                  }}
                >
                  上一张
                </Button>
                <Button
                  disabled={viewerIndex < FIRST_INDEX || viewerIndex >= uploadedImages.length - NEXT_INDEX_OFFSET}
                  onClick={() => {
                    if (viewerIndex >= FIRST_INDEX && viewerIndex < uploadedImages.length - NEXT_INDEX_OFFSET) {
                      void openImageViewer(uploadedImages[viewerIndex + NEXT_INDEX_OFFSET].image.uuid);
                    }
                  }}
                >
                  下一张
                </Button>
              </div>
            </div>
          </div>
        ) : null}
      </Modal>
    </div>
  );
}
