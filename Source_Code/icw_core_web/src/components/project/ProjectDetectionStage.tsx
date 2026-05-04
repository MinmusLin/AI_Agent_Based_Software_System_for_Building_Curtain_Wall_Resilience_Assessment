import { message } from 'antd';
import type { DragEvent, ReactElement } from 'react';
import { useCallback, useEffect, useMemo, useState } from 'react';

import { getErrorMessage } from '@/api/http';
import { getProjectAssets } from '@/api/project/assets';
import { ProjectAssessContent } from '@/components/project/assess/ProjectAssessContent';
import { ProjectAssessToolbar } from '@/components/project/assess/ProjectAssessToolbar';
import { ProjectAssessViewer } from '@/components/project/assess/ProjectAssessViewer';
import { useProjectAssetsViewer } from '@/hooks/project/useProjectAssetsViewer';
import type { ProjectGroup } from '@/types/project/assets';
import {
  flattenUploadedImages,
  groupIdsSet,
  normalizeGroups,
  projectGroupImageStats,
  UPLOAD_ACCEPT,
} from '@/utils/assetsStage';

interface ProjectDetectionStageProps {
  loading?: boolean;
  projectId: string;
}

export function ProjectDetectionStage({ loading = false, projectId }: ProjectDetectionStageProps): ReactElement {
  const [messageApi, contextHolder] = message.useMessage();
  const [groups, setGroups] = useState<ProjectGroup[]>([]);
  const [assetsLoading, setAssetsLoading] = useState(true);
  const [collapsedGroupIds, setCollapsedGroupIds] = useState<Set<string>>(() => new Set());
  const showError = useCallback(
    (text: string): void => {
      void messageApi.error(text);
    },
    [messageApi],
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
      showError(getErrorMessage(error));
    } finally {
      setAssetsLoading(false);
    }
  }, [projectId, showError]);

  useEffect(() => {
    if (loading) {
      return;
    }
    void loadAssets();
  }, [loadAssets, loading]);

  const uploadedImages = useMemo(() => flattenUploadedImages(groups), [groups]);
  const imageStats = useMemo(() => projectGroupImageStats(groups), [groups]);
  const { openImageViewer, setViewer, viewer, viewerImage, viewerIndex } = useProjectAssetsViewer({
    onError: showError,
    projectId,
    uploadedImages,
  });

  const handleCollapseAllGroups = useCallback((): void => {
    setCollapsedGroupIds(groupIdsSet(groups));
  }, [groups]);
  const handleExpandAllGroups = useCallback((): void => {
    setCollapsedGroupIds(new Set());
  }, []);
  const noop = useCallback((): void => undefined, []);
  const noopGroupDragOver = useCallback((event: DragEvent<HTMLElement>): void => {
    event.preventDefault();
  }, []);

  return (
    <div className="flex min-h-0 flex-1 flex-col rounded-lg border border-slate-200 bg-white p-5">
      {contextHolder}
      <ProjectAssessToolbar
        advancing={false}
        assetsLoading={assetsLoading}
        batchDeleting={false}
        batchMode={false}
        batchMoveDisabled
        batchMoveMenuItems={[]}
        batchMoving={false}
        canComplete={false}
        creatingGroup={false}
        failedImageCount={imageStats.failed}
        groupCount={groups.length}
        hasSelectedImages={false}
        onBatchDeleteImages={noop}
        onBatchModeToggle={noop}
        onBatchMoveImages={noop}
        onCollapseAllGroups={handleCollapseAllGroups}
        onComplete={noop}
        onCreateGroup={noop}
        onExpandAllGroups={handleExpandAllGroups}
        onRefresh={noop}
        pendingImageCount={imageStats.pending}
        readOnly
        totalImageCount={imageStats.total}
      />
      <div className="min-h-0 flex-1 overflow-y-auto pr-1">
        <ProjectAssessContent
          accept={UPLOAD_ACCEPT}
          assetsLoading={assetsLoading}
          batchMode={false}
          collapsedGroupIds={collapsedGroupIds}
          deletingGroupIds={new Set()}
          deletingImageUuids={new Set()}
          draggingGroupId=""
          editingGroupId=""
          editingGroupName=""
          groups={groups}
          onDeleteGroup={noop}
          onDeleteImage={noop}
          onDeselectGroupImages={noop}
          onEditingGroupNameChange={noop}
          onGroupDragEnd={noop}
          onGroupDragOver={noopGroupDragOver}
          onGroupDragStart={noop}
          onGroupDrop={noopGroupDragOver}
          onImageDragEnd={noop}
          onImageDragStart={noop}
          onOpenImageViewer={(imageUuid) => {
            void openImageViewer(imageUuid);
          }}
          onSaveEditGroup={noop}
          onSelectGroupImages={noop}
          onStartEditGroup={noop}
          onToggleCollapsed={(groupId) => {
            setCollapsedGroupIds((currentIds) => {
              const nextIds = new Set(currentIds);
              if (nextIds.has(groupId)) {
                nextIds.delete(groupId);
              } else {
                nextIds.add(groupId);
              }
              return nextIds;
            });
          }}
          onToggleSelectedImage={noop}
          onUploadFiles={noop}
          readOnly
          registerGroupRef={noop}
          savingGroupId=""
          selectedImageUuids={new Set()}
          uploadingImages={false}
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
