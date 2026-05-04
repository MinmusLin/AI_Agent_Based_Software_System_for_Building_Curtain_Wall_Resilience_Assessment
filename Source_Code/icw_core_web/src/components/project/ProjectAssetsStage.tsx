import type { MenuProps } from 'antd';
import { message } from 'antd';
import type { DragEvent, ReactElement } from 'react';
import { useCallback, useEffect, useLayoutEffect, useMemo, useRef, useState } from 'react';

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
import { ProjectAssessContent } from '@/components/project/assess/ProjectAssessContent';
import { ProjectAssessToolbar } from '@/components/project/assess/ProjectAssessToolbar';
import { ProjectAssessViewer } from '@/components/project/assess/ProjectAssessViewer';
import type { ProjectImageStatus, ProjectProgress } from '@/types/common';
import {
  PROJECT_IMAGE_STATUS_FAILED,
  PROJECT_IMAGE_STATUS_UPLOADED,
  PROJECT_PROGRESS_ASSETS_FINISHED,
  PROJECT_PROGRESS_PROFILE_FINISHED,
} from '@/types/common';
import type { ProjectGroup, UploadProjectImageResult } from '@/types/project/assets';
import type { Project } from '@/types/project/core';
import type { DraggingImage, ImageViewerState, PreparedImageUpload } from '@/utils/assetsStage';
import {
  appendImagesToGroup,
  EMPTY_ITEMS_COUNT,
  FIRST_INDEX,
  flattenUploadedImages,
  GROUP_MOVE_ANIMATION_MS,
  GROUP_MOVE_DELTA_EPSILON,
  groupIdsSet,
  groupMovePayload,
  groupMovePosition,
  isSameGroupOrder,
  moveImagesToGroup,
  NEXT_INDEX_OFFSET,
  normalizeGroups,
  NOT_FOUND_INDEX,
  parseProjectImageStatusChangedMessage,
  prepareUploadImage,
  pruneSelectedImageUuids,
  putPresignedObject,
  removeGroup,
  removeImages,
  replaceGroup,
  replaceImage,
  selectedGroupIdsFromGroups,
  UPLOAD_ACCEPT,
  uploadItemFromPreparedImage,
  WEBSOCKET_RECONNECT_DELAY_MS,
} from '@/utils/assetsStage';
import {
  isAllowedProjectAssetImageFile,
  PROJECT_ASSET_IMAGE_OUTPUT_CONTENT_TYPE,
  PROJECT_ASSET_THUMBNAIL_OUTPUT_CONTENT_TYPE,
} from '@/utils/images';

interface ProjectAssetsStageProps {
  loading?: boolean;
  onProgressChange: (progress: ProjectProgress) => void;
  onProjectChange: (project: Project) => void;
  project: Project;
  projectId: string;
  selectedProgress: ProjectProgress;
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
  const collapsedGroupSnapshotRef = useRef<Set<string> | null>(null);
  const groupMoveAnimationFrameRef = useRef<number | null>(null);
  const groupMoveRectsRef = useRef<Map<string, DOMRect> | null>(null);
  const groupSectionRefs = useRef(new Map<string, HTMLElement>());
  const draggingGroupTargetRef = useRef('');
  const draggingGroupSnapshotRef = useRef<ProjectGroup[] | null>(null);
  const movingGroupRef = useRef(false);
  const [groups, setGroups] = useState<ProjectGroup[]>([]);
  const [assetsLoading, setAssetsLoading] = useState(true);
  const [batchDeleting, setBatchDeleting] = useState(false);
  const [batchMode, setBatchMode] = useState(false);
  const [batchMoving, setBatchMoving] = useState(false);
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
  const [selectedImageUuids, setSelectedImageUuids] = useState<Set<string>>(() => new Set());
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
  const selectedGroupIds = useMemo(
    () => selectedGroupIdsFromGroups(groups, selectedImageUuids),
    [groups, selectedImageUuids],
  );
  const selectedImageUuidList = useMemo(() => [...selectedImageUuids], [selectedImageUuids]);
  const batchMoveMenuItems = useMemo<NonNullable<MenuProps['items']>>(() => {
    const shouldExcludeOnlySelectedGroup = selectedGroupIds.size === NEXT_INDEX_OFFSET;
    return groups
      .filter((group) => !(shouldExcludeOnlySelectedGroup && selectedGroupIds.has(group.id)))
      .map((group) => ({
        key: group.id,
        label: group.name,
      }));
  }, [groups, selectedGroupIds]);
  const hasSelectedImages = selectedImageUuids.size > EMPTY_ITEMS_COUNT;
  const batchMoveDisabled = batchMoving || !hasSelectedImages || batchMoveMenuItems.length === EMPTY_ITEMS_COUNT;

  const getGroupRects = useCallback((): Map<string, DOMRect> => {
    const rects = new Map<string, DOMRect>();
    groupSectionRefs.current.forEach((element, groupId) => {
      rects.set(groupId, element.getBoundingClientRect());
    });
    return rects;
  }, []);

  const animateGroupMove = useCallback((previousRects: Map<string, DOMRect>): void => {
    if (groupMoveAnimationFrameRef.current !== null) {
      window.cancelAnimationFrame(groupMoveAnimationFrameRef.current);
      groupMoveAnimationFrameRef.current = null;
    }

    groupSectionRefs.current.forEach((element, groupId) => {
      const previousRect = previousRects.get(groupId);
      if (!previousRect) {
        return;
      }

      const nextRect = element.getBoundingClientRect();
      const deltaX = previousRect.left - nextRect.left;
      const deltaY = previousRect.top - nextRect.top;
      if (Math.abs(deltaX) < GROUP_MOVE_DELTA_EPSILON && Math.abs(deltaY) < GROUP_MOVE_DELTA_EPSILON) {
        return;
      }

      element.style.transitionDuration = '0ms';
      element.style.transitionProperty = 'transform';
      element.style.transform = `translate(${String(deltaX)}px, ${String(deltaY)}px)`;
      element.style.zIndex = '1';
      groupMoveAnimationFrameRef.current = window.requestAnimationFrame(() => {
        element.style.transition = `transform ${String(GROUP_MOVE_ANIMATION_MS)}ms cubic-bezier(0.22, 1, 0.36, 1)`;
        element.style.transform = 'translate(0, 0)';
        window.setTimeout(() => {
          element.style.transition = '';
          element.style.transform = '';
          element.style.zIndex = '';
        }, GROUP_MOVE_ANIMATION_MS);
      });
    });
  }, []);

  const restoreCollapsedGroups = useCallback((): void => {
    if (!collapsedGroupSnapshotRef.current) {
      return;
    }
    setCollapsedGroupIds(collapsedGroupSnapshotRef.current);
    collapsedGroupSnapshotRef.current = null;
  }, []);

  const registerGroupRef = useCallback((groupId: string, node: HTMLElement | null): void => {
    if (node) {
      groupSectionRefs.current.set(groupId, node);
      return;
    }
    groupSectionRefs.current.delete(groupId);
  }, []);

  const handleToggleCollapsed = useCallback((groupId: string): void => {
    setCollapsedGroupIds((currentIds) => {
      const nextIds = new Set(currentIds);
      if (nextIds.has(groupId)) {
        nextIds.delete(groupId);
      } else {
        nextIds.add(groupId);
      }
      return nextIds;
    });
  }, []);

  const handleCancelEditGroup = useCallback((): void => {
    setEditingGroupId('');
  }, []);

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

  const handleBatchModeToggle = useCallback((): void => {
    if (batchMode) {
      setBatchMode(false);
      setSelectedImageUuids(new Set());
      return;
    }
    setBatchMode(true);
  }, [batchMode]);

  const toggleSelectedImage = useCallback((imageUuid: string, checked?: boolean): void => {
    setSelectedImageUuids((currentImageUuids) => {
      const nextImageUuids = new Set(currentImageUuids);
      const nextChecked = checked ?? !nextImageUuids.has(imageUuid);
      if (nextChecked) {
        nextImageUuids.add(imageUuid);
      } else {
        nextImageUuids.delete(imageUuid);
      }
      return nextImageUuids;
    });
  }, []);

  const handleBatchDeleteImages = useCallback(async (): Promise<void> => {
    if (!hasSelectedImages) {
      return;
    }

    setBatchDeleting(true);
    try {
      await deleteProjectImage({
        image_uuids: selectedImageUuidList,
        project_id: projectId,
      });
      setGroups((currentGroups) => removeImages(currentGroups, selectedImageUuidList));
      setSelectedImageUuids(new Set());
    } catch (error: unknown) {
      void messageApi.error(getErrorMessage(error));
    } finally {
      setBatchDeleting(false);
    }
  }, [hasSelectedImages, messageApi, projectId, selectedImageUuidList]);

  const handleBatchMoveImages = useCallback(
    async (targetGroupId: string): Promise<void> => {
      if (!hasSelectedImages || targetGroupId === '') {
        return;
      }

      setBatchMoving(true);
      try {
        const data = await moveProjectImage({
          image_uuids: selectedImageUuidList,
          project_id: projectId,
          target_group_id: targetGroupId,
        });
        setGroups((currentGroups) => moveImagesToGroup(currentGroups, targetGroupId, data.images));
        setSelectedImageUuids(new Set());
      } catch (error: unknown) {
        void messageApi.error(getErrorMessage(error));
      } finally {
        setBatchMoving(false);
      }
    },
    [hasSelectedImages, messageApi, projectId, selectedImageUuidList],
  );

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
    if (nextName === '') {
      setEditingGroupName(currentGroup.name);
      setEditingGroupId('');
      return;
    }
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

  const handleGroupDrop = useCallback(async (): Promise<void> => {
    if (readOnly || draggingGroupId === '') {
      return;
    }

    const originalGroups = draggingGroupSnapshotRef.current ?? groups;
    if (isSameGroupOrder(originalGroups, groups)) {
      draggingGroupSnapshotRef.current = null;
      draggingGroupTargetRef.current = '';
      restoreCollapsedGroups();
      setDraggingGroupId('');
      return;
    }

    movingGroupRef.current = true;
    try {
      const payload = groupMovePayload(projectId, draggingGroupId, groups);
      const data = await moveProjectGroup(payload);
      setGroups((currentGroups) => replaceGroup(currentGroups, data.group));
    } catch (error: unknown) {
      setGroups(originalGroups);
      void messageApi.error(getErrorMessage(error));
    } finally {
      movingGroupRef.current = false;
      draggingGroupSnapshotRef.current = null;
      draggingGroupTargetRef.current = '';
      restoreCollapsedGroups();
      setDraggingGroupId('');
    }
  }, [draggingGroupId, groups, messageApi, projectId, readOnly, restoreCollapsedGroups]);

  const handleGroupDragOver = useCallback(
    (targetGroupId: string): void => {
      if (
        readOnly ||
        draggingGroupId === '' ||
        draggingGroupId === targetGroupId ||
        draggingGroupTargetRef.current === targetGroupId
      ) {
        return;
      }

      draggingGroupTargetRef.current = targetGroupId;
      groupMoveRectsRef.current = getGroupRects();
      setGroups((currentGroups) => groupMovePosition(currentGroups, draggingGroupId, targetGroupId));
    },
    [draggingGroupId, getGroupRects, readOnly],
  );

  const handleGroupDragEnd = useCallback((): void => {
    if (!movingGroupRef.current && draggingGroupSnapshotRef.current) {
      setGroups(draggingGroupSnapshotRef.current);
    }
    if (!movingGroupRef.current) {
      draggingGroupSnapshotRef.current = null;
      draggingGroupTargetRef.current = '';
      restoreCollapsedGroups();
      setDraggingGroupId('');
    }
  }, [restoreCollapsedGroups]);

  const handleGroupDragStart = useCallback(
    (event: DragEvent<HTMLDivElement>, group: ProjectGroup): void => {
      if (readOnly) {
        return;
      }
      draggingGroupSnapshotRef.current = groups;
      draggingGroupTargetRef.current = '';
      collapsedGroupSnapshotRef.current = new Set(collapsedGroupIds);
      setCollapsedGroupIds(groupIdsSet(groups));
      setDraggingGroupId(group.id);
      event.dataTransfer.effectAllowed = 'move';
    },
    [collapsedGroupIds, groups, readOnly],
  );

  const handleGroupDragOverEvent = useCallback(
    (event: DragEvent<HTMLElement>, groupId: string): void => {
      if (draggingGroupId !== '' && !draggingImage) {
        event.preventDefault();
        handleGroupDragOver(groupId);
        return;
      }
      if (draggingImage) {
        event.preventDefault();
      }
    },
    [draggingGroupId, draggingImage, handleGroupDragOver],
  );

  const handleImageDragStart = useCallback(
    (event: DragEvent<HTMLDivElement>, imageUuid: string, sourceGroupId: string): void => {
      if (readOnly || batchMode) {
        return;
      }
      event.stopPropagation();
      event.dataTransfer.effectAllowed = 'move';
      setDraggingImage({
        imageUuid,
        sourceGroupId,
      });
    },
    [batchMode, readOnly],
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

  const handleGroupDropEvent = useCallback(
    (event: DragEvent<HTMLElement>, groupId: string): void => {
      event.preventDefault();
      if (draggingImage) {
        void handleImageDrop(groupId);
        return;
      }
      void handleGroupDrop();
    },
    [draggingImage, handleGroupDrop, handleImageDrop],
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
        setSelectedImageUuids((currentImageUuids) => {
          const nextImageUuids = new Set(currentImageUuids);
          nextImageUuids.delete(imageUuid);
          return nextImageUuids;
        });
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
    setSelectedImageUuids((currentImageUuids) => pruneSelectedImageUuids(currentImageUuids, groups));
  }, [groups]);

  useLayoutEffect(() => {
    if (!groupMoveRectsRef.current) {
      return;
    }
    const previousRects = groupMoveRectsRef.current;
    groupMoveRectsRef.current = null;
    animateGroupMove(previousRects);
  }, [animateGroupMove, groups]);

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
  const totalImageCount = useMemo(
    () => groups.reduce((count, group) => count + group.images.length, EMPTY_ITEMS_COUNT),
    [groups],
  );

  return (
    <div className="flex min-h-0 flex-1 flex-col rounded-lg border border-slate-200 bg-white p-5">
      {contextHolder}
      <ProjectAssessToolbar
        advancing={advancing}
        assetsLoading={assetsLoading}
        batchDeleting={batchDeleting}
        batchMode={batchMode}
        batchMoveDisabled={batchMoveDisabled}
        batchMoveMenuItems={batchMoveMenuItems}
        batchMoving={batchMoving}
        canComplete={canComplete}
        creatingGroup={creatingGroup}
        groupCount={groups.length}
        hasSelectedImages={hasSelectedImages}
        onBatchDeleteImages={() => {
          void handleBatchDeleteImages();
        }}
        onBatchModeToggle={handleBatchModeToggle}
        onBatchMoveImages={(targetGroupId) => {
          void handleBatchMoveImages(targetGroupId);
        }}
        onComplete={() => {
          void handleComplete();
        }}
        onCreateGroup={() => {
          void handleCreateGroup();
        }}
        onRefresh={() => {
          void loadAssets();
        }}
        readOnly={readOnly}
        totalImageCount={totalImageCount}
      />
      <div className="min-h-0 flex-1 overflow-y-auto pr-1">
        <ProjectAssessContent
          accept={UPLOAD_ACCEPT}
          assetsLoading={assetsLoading}
          batchMode={batchMode}
          collapsedGroupIds={collapsedGroupIds}
          deletingGroupIds={deletingGroupIds}
          deletingImageUuids={deletingImageUuids}
          draggingGroupId={draggingGroupId}
          editingGroupId={editingGroupId}
          editingGroupName={editingGroupName}
          groups={groups}
          onCancelEditGroup={handleCancelEditGroup}
          onDeleteGroup={(groupId) => {
            void handleDeleteGroup(groupId);
          }}
          onDeleteImage={(imageUuid) => {
            void handleDeleteImage(imageUuid);
          }}
          onEditingGroupNameChange={setEditingGroupName}
          onGroupDragEnd={handleGroupDragEnd}
          onGroupDragOver={handleGroupDragOverEvent}
          onGroupDragStart={handleGroupDragStart}
          onGroupDrop={handleGroupDropEvent}
          onImageDragEnd={() => {
            setDraggingImage(null);
          }}
          onImageDragStart={handleImageDragStart}
          onOpenImageViewer={(imageUuid) => {
            void openImageViewer(imageUuid);
          }}
          onSaveEditGroup={() => {
            void saveEditingGroup();
          }}
          onStartEditGroup={startEditGroup}
          onToggleCollapsed={handleToggleCollapsed}
          onToggleSelectedImage={toggleSelectedImage}
          onUploadFiles={(groupId, files) => {
            void handleUploadFiles(groupId, files);
          }}
          readOnly={readOnly}
          registerGroupRef={registerGroupRef}
          savingGroupId={savingGroupId}
          selectedImageUuids={selectedImageUuids}
        />
      </div>
      <ProjectAssessViewer
        onClose={() => {
          setViewer(null);
        }}
        onOpenImage={(imageUuid) => {
          void openImageViewer(imageUuid);
        }}
        uploadedImages={uploadedImages}
        viewer={viewer}
        viewerImage={viewerImage}
        viewerIndex={viewerIndex}
      />
    </div>
  );
}
