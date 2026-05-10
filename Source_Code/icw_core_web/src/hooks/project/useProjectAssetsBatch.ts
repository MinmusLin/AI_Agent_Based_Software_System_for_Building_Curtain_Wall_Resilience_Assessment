import type { MenuProps } from 'antd';
import type { Dispatch, SetStateAction } from 'react';
import { useCallback, useMemo, useState } from 'react';

import { getErrorMessage } from '@/api/http';
import { deleteProjectImage, moveProjectImage } from '@/api/project/assets';
import type { ProjectGroup } from '@/gen/core/api/common';
import {
  EMPTY_ITEMS_COUNT,
  moveImagesToGroup,
  NEXT_INDEX_OFFSET,
  pruneSelectedImageUuids,
  removeImages,
  selectedGroupIdsFromGroups,
} from '@/utils/assetsStage';

interface UseProjectAssetsBatchParams {
  groups: ProjectGroup[];
  onError: (message: string) => void;
  projectId: string;
  setGroups: Dispatch<SetStateAction<ProjectGroup[]>>;
}

interface UseProjectAssetsBatchResult {
  batchDeleting: boolean;
  batchMode: boolean;
  batchMoveDisabled: boolean;
  batchMoveMenuItems: NonNullable<MenuProps['items']>;
  batchMoving: boolean;
  handleBatchDeleteImages: () => Promise<void>;
  handleBatchModeToggle: () => void;
  handleBatchMoveImages: (targetGroupId: string) => Promise<void>;
  handleClearGroupImages: (imageUuids: string[]) => void;
  handleSelectGroupImages: (imageUuids: string[]) => void;
  hasSelectedImages: boolean;
  pruneSelectedImages: (groups: ProjectGroup[]) => void;
  selectedImageUuids: Set<string>;
  toggleSelectedImage: (imageUuid: string, checked?: boolean) => void;
}

export function useProjectAssetsBatch({
  groups,
  onError,
  projectId,
  setGroups,
}: UseProjectAssetsBatchParams): UseProjectAssetsBatchResult {
  const [batchDeleting, setBatchDeleting] = useState(false);
  const [batchMode, setBatchMode] = useState(false);
  const [batchMoving, setBatchMoving] = useState(false);
  const [selectedImageUuids, setSelectedImageUuids] = useState<Set<string>>(() => new Set());

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

  const handleSelectGroupImages = useCallback((imageUuids: string[]): void => {
    setSelectedImageUuids((currentImageUuids) => {
      const nextImageUuids = new Set(currentImageUuids);
      imageUuids.forEach((imageUuid) => {
        nextImageUuids.add(imageUuid);
      });
      return nextImageUuids;
    });
  }, []);

  const handleClearGroupImages = useCallback((imageUuids: string[]): void => {
    setSelectedImageUuids((currentImageUuids) => {
      const nextImageUuids = new Set(currentImageUuids);
      imageUuids.forEach((imageUuid) => {
        nextImageUuids.delete(imageUuid);
      });
      return nextImageUuids;
    });
  }, []);

  const pruneSelectedImages = useCallback((nextGroups: ProjectGroup[]): void => {
    setSelectedImageUuids((currentImageUuids) => pruneSelectedImageUuids(currentImageUuids, nextGroups));
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
      onError(getErrorMessage(error));
    } finally {
      setBatchDeleting(false);
    }
  }, [hasSelectedImages, onError, projectId, selectedImageUuidList, setGroups]);

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
        onError(getErrorMessage(error));
      } finally {
        setBatchMoving(false);
      }
    },
    [hasSelectedImages, onError, projectId, selectedImageUuidList, setGroups],
  );

  return {
    batchDeleting,
    batchMode,
    batchMoveDisabled,
    batchMoveMenuItems,
    batchMoving,
    handleBatchDeleteImages,
    handleBatchModeToggle,
    handleBatchMoveImages,
    handleClearGroupImages,
    handleSelectGroupImages,
    hasSelectedImages,
    pruneSelectedImages,
    selectedImageUuids,
    toggleSelectedImage,
  };
}
