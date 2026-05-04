import type { Dispatch, DragEvent, SetStateAction } from 'react';
import { useCallback, useLayoutEffect, useRef, useState } from 'react';

import { getErrorMessage } from '@/api/http';
import { moveProjectGroup, moveProjectImage } from '@/api/project/assets';
import type { ProjectGroup } from '@/types/project/assets';
import type { DraggingImage } from '@/utils/assetsStage';
import {
  GROUP_MOVE_ANIMATION_MS,
  GROUP_MOVE_DELTA_EPSILON,
  groupIdsSet,
  groupMovePayload,
  groupMovePosition,
  isSameGroupOrder,
  moveImagesToGroup,
  replaceGroup,
} from '@/utils/assetsStage';

interface UseProjectAssetsDragParams {
  batchMode: boolean;
  groups: ProjectGroup[];
  onError: (message: string) => void;
  projectId: string;
  readOnly: boolean;
  setGroups: Dispatch<SetStateAction<ProjectGroup[]>>;
}

interface UseProjectAssetsDragResult {
  collapsedGroupIds: Set<string>;
  draggingGroupId: string;
  handleGroupDragEnd: () => void;
  handleGroupDragOverEvent: (event: DragEvent<HTMLElement>, groupId: string) => void;
  handleGroupDragStart: (event: DragEvent<HTMLDivElement>, group: ProjectGroup) => void;
  handleGroupDropEvent: (event: DragEvent<HTMLElement>, groupId: string) => void;
  handleImageDragEnd: () => void;
  handleImageDragStart: (event: DragEvent<HTMLDivElement>, imageUuid: string, sourceGroupId: string) => void;
  handleToggleCollapsed: (groupId: string) => void;
  registerGroupRef: (groupId: string, node: HTMLElement | null) => void;
}

export function useProjectAssetsDrag({
  batchMode,
  groups,
  onError,
  projectId,
  readOnly,
  setGroups,
}: UseProjectAssetsDragParams): UseProjectAssetsDragResult {
  const collapsedGroupSnapshotRef = useRef<Set<string> | null>(null);
  const groupMoveAnimationFrameRef = useRef<number | null>(null);
  const groupMoveRectsRef = useRef<Map<string, DOMRect> | null>(null);
  const groupSectionRefs = useRef(new Map<string, HTMLElement>());
  const draggingGroupTargetRef = useRef('');
  const draggingGroupSnapshotRef = useRef<ProjectGroup[] | null>(null);
  const movingGroupRef = useRef(false);
  const [collapsedGroupIds, setCollapsedGroupIds] = useState<Set<string>>(() => new Set());
  const [draggingGroupId, setDraggingGroupId] = useState('');
  const [draggingImage, setDraggingImage] = useState<DraggingImage | null>(null);

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
      onError(getErrorMessage(error));
    } finally {
      movingGroupRef.current = false;
      draggingGroupSnapshotRef.current = null;
      draggingGroupTargetRef.current = '';
      restoreCollapsedGroups();
      setDraggingGroupId('');
    }
  }, [draggingGroupId, groups, onError, projectId, readOnly, restoreCollapsedGroups, setGroups]);

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
    [draggingGroupId, getGroupRects, readOnly, setGroups],
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
  }, [restoreCollapsedGroups, setGroups]);

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
        onError(getErrorMessage(error));
      } finally {
        setDraggingImage(null);
      }
    },
    [draggingImage, onError, projectId, readOnly, setGroups],
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

  const handleImageDragEnd = useCallback((): void => {
    setDraggingImage(null);
  }, []);

  useLayoutEffect(() => {
    if (!groupMoveRectsRef.current) {
      return;
    }
    const previousRects = groupMoveRectsRef.current;
    groupMoveRectsRef.current = null;
    animateGroupMove(previousRects);
  }, [animateGroupMove, groups]);

  return {
    collapsedGroupIds,
    draggingGroupId,
    handleGroupDragEnd,
    handleGroupDragOverEvent,
    handleGroupDragStart,
    handleGroupDropEvent,
    handleImageDragEnd,
    handleImageDragStart,
    handleToggleCollapsed,
    registerGroupRef,
  };
}
