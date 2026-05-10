import { message } from 'antd';
import type { ReactElement } from 'react';
import { useCallback, useEffect, useMemo, useRef } from 'react';

import { ProjectAssessContent } from '@/components/project/assess/ProjectAssessContent';
import { ProjectAssessToolbar } from '@/components/project/assess/ProjectAssessToolbar';
import { ProjectAssessViewer } from '@/components/project/assess/ProjectAssessViewer';
import type { Project } from '@/gen/core/api/common';
import type { ProjectProgress_Value } from '@/gen/core/common';
import { useProjectAssetsActions } from '@/hooks/project/useProjectAssetsActions';
import { useProjectAssetsBatch } from '@/hooks/project/useProjectAssetsBatch';
import { useProjectAssetsDrag } from '@/hooks/project/useProjectAssetsDrag';
import { useProjectAssetsSocket } from '@/hooks/project/useProjectAssetsSocket';
import { useProjectAssetsViewer } from '@/hooks/project/useProjectAssetsViewer';
import { flattenUploadedImages, projectGroupImageStats, UPLOAD_ACCEPT } from '@/utils/assetsStage';

interface ProjectAssetsStageProps {
  loading?: boolean;
  onProgressChange: (progress: ProjectProgress_Value) => void;
  onProjectChange: (project: Project) => void;
  project: Project;
  projectId: string;
  selectedProgress: ProjectProgress_Value;
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
  const scrollContainerRef = useRef<HTMLDivElement | null>(null);
  const showError = useCallback(
    (text: string): void => {
      void messageApi.error(text);
    },
    [messageApi],
  );

  const {
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
  } = useProjectAssetsActions({
    loading,
    onError: showError,
    onProgressChange,
    onProjectChange,
    project,
    projectId,
    selectedProgress,
  });
  const uploadedImages = useMemo(() => flattenUploadedImages(groups), [groups]);
  const {
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
  } = useProjectAssetsBatch({
    groups,
    onError: showError,
    projectId,
    setGroups,
  });
  const {
    collapsedGroupIds,
    draggingGroupId,
    handleCollapseAllGroups,
    handleExpandAllGroups,
    handleGroupDragEnd,
    handleGroupDragOverEvent,
    handleGroupDragStart,
    handleGroupDropEvent,
    handleImageDragEnd,
    handleImageDragStart,
    handleToggleCollapsed,
    registerGroupRef,
  } = useProjectAssetsDrag({
    batchMode,
    groups,
    onError: showError,
    projectId,
    readOnly,
    scrollContainerRef,
    setGroups,
  });
  const { openImageViewer, setViewer, viewer, viewerImage, viewerIndex } = useProjectAssetsViewer({
    onError: showError,
    projectId,
    uploadedImages,
  });

  useProjectAssetsSocket({
    onConnected: () => {
      void loadAssets({ silent: true });
    },
    projectId,
    setGroups,
  });

  useEffect(() => {
    pruneSelectedImages(groups);
  }, [groups, pruneSelectedImages]);

  const imageStats = useMemo(() => projectGroupImageStats(groups), [groups]);

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
        failedImageCount={imageStats.failed}
        groupCount={groups.length}
        hasSelectedImages={hasSelectedImages}
        onBatchDeleteImages={() => {
          void handleBatchDeleteImages();
        }}
        onBatchModeToggle={handleBatchModeToggle}
        onBatchMoveImages={(targetGroupId) => {
          void handleBatchMoveImages(targetGroupId);
        }}
        onCollapseAllGroups={handleCollapseAllGroups}
        onComplete={() => {
          void handleComplete();
        }}
        onCreateGroup={() => {
          void handleCreateGroup();
        }}
        onExpandAllGroups={handleExpandAllGroups}
        onRefresh={() => {
          void loadAssets();
        }}
        pendingImageCount={imageStats.pending}
        readOnly={readOnly}
        totalImageCount={imageStats.total}
      />
      <div className="min-h-0 flex-1 overflow-y-auto pr-1" ref={scrollContainerRef}>
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
          onDeleteGroup={(groupId) => {
            void handleDeleteGroup(groupId);
          }}
          onDeleteImage={(imageUuid) => {
            void handleDeleteImage(imageUuid);
          }}
          onDeselectGroupImages={handleClearGroupImages}
          onEditingGroupNameChange={setEditingGroupName}
          onGroupDragEnd={handleGroupDragEnd}
          onGroupDragOver={handleGroupDragOverEvent}
          onGroupDragStart={handleGroupDragStart}
          onGroupDrop={handleGroupDropEvent}
          onImageDragEnd={handleImageDragEnd}
          onImageDragStart={handleImageDragStart}
          onOpenImageViewer={(imageUuid) => {
            void openImageViewer(imageUuid);
          }}
          onSaveEditGroup={() => {
            void saveEditingGroup();
          }}
          onSelectGroupImages={handleSelectGroupImages}
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
          uploadingImages={uploadingImages}
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
