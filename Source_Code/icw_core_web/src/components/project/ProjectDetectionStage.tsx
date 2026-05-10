import { Empty, message, Spin } from 'antd';
import type { ReactElement } from 'react';
import { useCallback, useEffect, useMemo, useState } from 'react';

import { getErrorMessage } from '@/api/http';
import { getProjectAssets } from '@/api/project/assets';
import { advanceProject } from '@/api/project/core';
import {
  getImageDetectionResult,
  getProjectDetectionTasks,
  retryProjectDetection,
  startProjectDetection,
} from '@/api/project/detection';
import { ProjectDetectionGroup } from '@/components/project/detection/ProjectDetectionGroup';
import { ProjectDetectionProgressViewer } from '@/components/project/detection/ProjectDetectionProgressViewer';
import { ProjectDetectionResultViewer } from '@/components/project/detection/ProjectDetectionResultViewer';
import { ProjectDetectionToolbar } from '@/components/project/detection/ProjectDetectionToolbar';
import type { Project, ProjectGroup } from '@/gen/core/api/common';
import type { GetImageDetectionResultResponse } from '@/gen/core/api/project_detection';
import { ProjectProgress_Value } from '@/gen/core/common';
import { useProjectDetectionSocket } from '@/hooks/project/useProjectDetectionSocket';
import { EMPTY_ITEMS_COUNT, flattenUploadedImages, normalizeGroups, NOT_FOUND_INDEX } from '@/utils/assetsStage';
import type { ProjectDetectionStatusMap } from '@/utils/detectionStage';
import {
  allDetectionTasksSucceeded,
  detectionActionState,
  detectionTasksMap,
  detectionTaskStatusStats,
  hasDetectionTasks,
  hasFailedDetectionTask,
  isDetectionTaskSucceeded,
} from '@/utils/detectionStage';

interface ProjectDetectionStageProps {
  loading?: boolean;
  onProgressChange: (progress: ProjectProgress_Value) => void;
  onProjectChange: (project: Project) => void;
  project: Project;
  projectId: string;
  selectedProgress: ProjectProgress_Value;
}

interface RefreshOptions {
  silent?: boolean;
}

export function ProjectDetectionStage({
  loading = false,
  onProgressChange,
  onProjectChange,
  project,
  projectId,
  selectedProgress,
}: ProjectDetectionStageProps): ReactElement {
  const [messageApi, contextHolder] = message.useMessage();
  const [groups, setGroups] = useState<ProjectGroup[]>([]);
  const [taskMap, setTaskMap] = useState<ProjectDetectionStatusMap>({});
  const [assetsLoading, setAssetsLoading] = useState(true);
  const [detectionLoading, setDetectionLoading] = useState(true);
  const [startingDetection, setStartingDetection] = useState(false);
  const [retryingDetection, setRetryingDetection] = useState(false);
  const [advancing, setAdvancing] = useState(false);
  const [collapsedGroupIds, setCollapsedGroupIds] = useState<Set<string>>(() => new Set());
  const [viewerImageUuid, setViewerImageUuid] = useState<string | null>(null);
  const [viewerLoading, setViewerLoading] = useState(false);
  const [viewerResult, setViewerResult] = useState<GetImageDetectionResultResponse | null>(null);
  const [progressImageUuid, setProgressImageUuid] = useState<string | null>(null);

  const readOnly =
    loading ||
    project.progress > ProjectProgress_Value.AssetsFinished ||
    selectedProgress !== ProjectProgress_Value.AssetsFinished;

  const showError = useCallback(
    (text: string): void => {
      void messageApi.error(text);
    },
    [messageApi],
  );

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
          showError(getErrorMessage(error));
        }
      } finally {
        if (!options.silent) {
          setAssetsLoading(false);
        }
      }
    },
    [projectId, showError],
  );

  const loadDetectionTasks = useCallback(
    async (options: RefreshOptions = {}): Promise<void> => {
      if (projectId === '') {
        return;
      }

      if (!options.silent) {
        setDetectionLoading(true);
      }
      try {
        const data = await getProjectDetectionTasks(projectId);
        setTaskMap(detectionTasksMap(data.tasks));
      } catch (error: unknown) {
        if (!options.silent) {
          showError(getErrorMessage(error));
        }
      } finally {
        if (!options.silent) {
          setDetectionLoading(false);
        }
      }
    },
    [projectId, showError],
  );

  const refreshPage = useCallback(async (): Promise<void> => {
    await Promise.all([loadAssets(), loadDetectionTasks()]);
  }, [loadAssets, loadDetectionTasks]);

  useEffect(() => {
    if (loading) {
      return;
    }
    void refreshPage();
  }, [loading, refreshPage]);

  const uploadedImages = useMemo(() => flattenUploadedImages(groups), [groups]);
  const viewerIndex = useMemo(() => {
    if (!viewerImageUuid) {
      return NOT_FOUND_INDEX;
    }
    return uploadedImages.findIndex((item) => item.image.uuid === viewerImageUuid);
  }, [uploadedImages, viewerImageUuid]);
  const progressTask = progressImageUuid ? taskMap[progressImageUuid] : undefined;

  const openDetectionResultViewer = useCallback(
    async (imageUuid: string): Promise<void> => {
      setViewerImageUuid(imageUuid);
      setViewerLoading(true);
      setViewerResult(null);
      try {
        const data = await getImageDetectionResult({
          image_uuid: imageUuid,
          project_id: projectId,
        });
        setViewerResult(data);
      } catch (error: unknown) {
        setViewerImageUuid(null);
        showError(getErrorMessage(error));
      } finally {
        setViewerLoading(false);
      }
    },
    [projectId, showError],
  );

  const openDetectionImage = useCallback(
    (imageUuid: string): void => {
      const task = taskMap[imageUuid];
      if (task && !isDetectionTaskSucceeded(task)) {
        setProgressImageUuid(imageUuid);
        return;
      }
      void openDetectionResultViewer(imageUuid);
    },
    [openDetectionResultViewer, taskMap],
  );

  const handleStartDetection = useCallback(async (): Promise<void> => {
    if (projectId === '') {
      return;
    }

    setStartingDetection(true);
    try {
      const data = await startProjectDetection(projectId);
      void messageApi.success(`已提交 ${String(data.task_count)} 张图像检测任务`);
      await loadDetectionTasks();
    } catch (error: unknown) {
      showError(getErrorMessage(error));
    } finally {
      setStartingDetection(false);
    }
  }, [loadDetectionTasks, messageApi, projectId, showError]);

  const handleRetryDetection = useCallback(async (): Promise<void> => {
    if (projectId === '') {
      return;
    }

    setRetryingDetection(true);
    try {
      const data = await retryProjectDetection(projectId);
      void messageApi.success(`已重新提交 ${String(data.task_count)} 张失败图像检测任务`);
      await loadDetectionTasks();
    } catch (error: unknown) {
      showError(getErrorMessage(error));
    } finally {
      setRetryingDetection(false);
    }
  }, [loadDetectionTasks, messageApi, projectId, showError]);

  const handleComplete = useCallback(async (): Promise<void> => {
    setAdvancing(true);
    try {
      await advanceProject({
        from_progress: project.progress,
        project_id: projectId,
        to_progress: ProjectProgress_Value.DetectionFinished,
      });
      onProjectChange({
        ...project,
        progress: ProjectProgress_Value.DetectionFinished,
      });
      onProgressChange(ProjectProgress_Value.DetectionFinished);
    } catch (error: unknown) {
      showError(getErrorMessage(error));
    } finally {
      setAdvancing(false);
    }
  }, [onProgressChange, onProjectChange, project, projectId, showError]);

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

  const handleCollapseAllGroups = useCallback((): void => {
    setCollapsedGroupIds(new Set(groups.map((group) => group.id)));
  }, [groups]);

  const handleExpandAllGroups = useCallback((): void => {
    setCollapsedGroupIds(new Set());
  }, []);

  const pageLoading = loading || assetsLoading || detectionLoading;
  useProjectDetectionSocket({
    enabled: !loading,
    onConnected: () => {
      void loadDetectionTasks({ silent: true });
    },
    projectId,
    setTasks: setTaskMap,
  });

  useEffect(() => {
    if (progressTask && isDetectionTaskSucceeded(progressTask)) {
      setProgressImageUuid(null);
    }
  }, [progressTask]);

  const taskStatusStats = useMemo(() => detectionTaskStatusStats(taskMap), [taskMap]);
  const detectionStarted = hasDetectionTasks(taskMap);
  const detectionFailed = hasFailedDetectionTask(taskMap);
  const detectionCompleted = allDetectionTasksSucceeded(taskMap, uploadedImages.length);
  const actionState = detectionActionState({
    detectionCompleted,
    detectionFailed,
    detectionStarted,
    pageLoading,
    readOnly,
    uploadedImageCount: uploadedImages.length,
  });

  return (
    <div className="flex min-h-0 flex-1 flex-col rounded-lg border border-slate-200 bg-white p-5">
      {contextHolder}
      <ProjectDetectionToolbar
        advancing={advancing}
        canComplete={actionState.canComplete}
        canRetry={actionState.canRetry}
        canStart={actionState.canStart}
        failedTaskCount={taskStatusStats.failed}
        groupCount={groups.length}
        loading={pageLoading}
        onCollapseAllGroups={handleCollapseAllGroups}
        onComplete={() => {
          void handleComplete();
        }}
        onExpandAllGroups={handleExpandAllGroups}
        onRefresh={() => {
          void refreshPage();
        }}
        onRetry={() => {
          void handleRetryDetection();
        }}
        onStart={() => {
          void handleStartDetection();
        }}
        retrying={retryingDetection}
        runningTaskCount={taskStatusStats.running}
        showComplete={actionState.showComplete}
        showRetry={actionState.showRetry}
        showStart={actionState.showStart}
        starting={startingDetection}
        totalImageCount={uploadedImages.length}
      />
      <div className="min-h-0 flex-1 overflow-y-auto pr-1">
        {pageLoading ? (
          <div className="flex h-full items-center justify-center">
            <Spin description="正在加载智能检测数据" />
          </div>
        ) : null}
        {!pageLoading && groups.length === EMPTY_ITEMS_COUNT ? (
          <div className="flex h-full items-center justify-center rounded-lg border border-dashed border-slate-200">
            <Empty description="暂无图像资产，请完成图像资产构建" />
          </div>
        ) : null}
        {!pageLoading && groups.length > EMPTY_ITEMS_COUNT ? (
          <div className="space-y-4">
            {groups.map((group) => (
              <ProjectDetectionGroup
                collapsed={collapsedGroupIds.has(group.id)}
                group={group}
                key={group.id}
                onOpenImageViewer={openDetectionImage}
                onOpenProgressViewer={setProgressImageUuid}
                onToggleCollapsed={handleToggleCollapsed}
                taskMap={taskMap}
              />
            ))}
          </div>
        ) : null}
      </div>
      <ProjectDetectionResultViewer
        loading={viewerLoading}
        onClose={() => {
          setViewerImageUuid(null);
          setViewerResult(null);
        }}
        onOpenImage={(imageUuid) => {
          void openDetectionResultViewer(imageUuid);
        }}
        open={viewerImageUuid !== null}
        result={viewerResult}
        uploadedImages={uploadedImages}
        viewerIndex={viewerIndex}
      />
      <ProjectDetectionProgressViewer
        onClose={() => {
          setProgressImageUuid(null);
        }}
        open={progressImageUuid !== null && progressTask !== undefined && !isDetectionTaskSucceeded(progressTask)}
        task={progressTask}
      />
    </div>
  );
}
