import { MenuFoldOutlined, MenuUnfoldOutlined, ReloadOutlined, StepForwardOutlined } from '@ant-design/icons';
import { Button, Empty, message, Spin } from 'antd';
import type { ReactElement } from 'react';
import { useCallback, useEffect, useMemo, useState } from 'react';

import { getErrorMessage } from '@/api/http';
import { getProjectAssets } from '@/api/project/assets';
import { advanceProject } from '@/api/project/core';
import { getImageDetectionResult, getProjectDetectionTasks } from '@/api/project/detection';
import { getProjectDetectionReview, updateProjectDetectionReview } from '@/api/project/review';
import { ProjectDetectionGroup } from '@/components/project/detection/ProjectDetectionGroup';
import { ProjectDetectionProgressViewer } from '@/components/project/detection/ProjectDetectionProgressViewer';
import { ProjectDetectionResultViewer } from '@/components/project/detection/ProjectDetectionResultViewer';
import type { Project, ProjectGroup } from '@/gen/core/api/common';
import type { GetImageDetectionResultResponse } from '@/gen/core/api/project_detection';
import { ProjectDetectionReviewVerdict_Value, ProjectProgress_Value } from '@/gen/core/common';
import { EMPTY_ITEMS_COUNT, flattenUploadedImages, normalizeGroups, NOT_FOUND_INDEX } from '@/utils/assetsStage';
import type { ProjectDetectionStatusMap } from '@/utils/detectionStage';
import { detectionTasksMap } from '@/utils/detectionStage';

const DEFAULT_TASK_COUNT = 0;

interface ProjectReviewStageProps {
  loading?: boolean;
  onProgressChange: (progress: ProjectProgress_Value) => void;
  onProjectChange: (project: Project) => void;
  project: Project;
  projectId: string;
}

interface RefreshOptions {
  initial?: boolean;
}

export function ProjectReviewStage({
  loading = false,
  onProgressChange,
  onProjectChange,
  project,
  projectId,
}: ProjectReviewStageProps): ReactElement {
  const [messageApi, contextHolder] = message.useMessage();
  const [groups, setGroups] = useState<ProjectGroup[]>([]);
  const [taskMap, setTaskMap] = useState<ProjectDetectionStatusMap>({});
  const [assetsLoading, setAssetsLoading] = useState(true);
  const [detectionLoading, setDetectionLoading] = useState(true);
  const [initialPageLoading, setInitialPageLoading] = useState(true);
  const [advancing, setAdvancing] = useState(false);
  const [collapsedGroupIds, setCollapsedGroupIds] = useState<Set<string>>(() => new Set());
  const [viewerImageUuid, setViewerImageUuid] = useState<string | null>(null);
  const [viewerLoading, setViewerLoading] = useState(false);
  const [viewerResult, setViewerResult] = useState<GetImageDetectionResultResponse | null>(null);
  const [progressImageUuid, setProgressImageUuid] = useState<string | null>(null);
  const [reviewComment, setReviewComment] = useState('');
  const [reviewVerdict, setReviewVerdict] = useState<ProjectDetectionReviewVerdict_Value>(
    ProjectDetectionReviewVerdict_Value.Unknown,
  );
  const [reviewSaving, setReviewSaving] = useState(false);

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

  const loadDetectionTasks = useCallback(async (): Promise<void> => {
    if (projectId === '') {
      return;
    }
    setDetectionLoading(true);
    try {
      const data = await getProjectDetectionTasks(projectId);
      setTaskMap(detectionTasksMap(data.tasks));
    } catch (error: unknown) {
      showError(getErrorMessage(error));
    } finally {
      setDetectionLoading(false);
    }
  }, [projectId, showError]);

  const refreshPage = useCallback(
    async (options: RefreshOptions = {}): Promise<void> => {
      try {
        await Promise.all([loadAssets(), loadDetectionTasks()]);
      } finally {
        if (options.initial) {
          setInitialPageLoading(false);
        }
      }
    },
    [loadAssets, loadDetectionTasks],
  );

  useEffect(() => {
    if (loading) {
      return;
    }
    void refreshPage({ initial: true });
  }, [loading, refreshPage]);

  const uploadedImages = useMemo(() => flattenUploadedImages(groups), [groups]);
  const viewerIndex = useMemo(() => {
    if (!viewerImageUuid) {
      return NOT_FOUND_INDEX;
    }
    return uploadedImages.findIndex((item) => item.image.uuid === viewerImageUuid);
  }, [uploadedImages, viewerImageUuid]);
  const progressTask = progressImageUuid ? taskMap[progressImageUuid] : undefined;

  const openReviewViewer = useCallback(
    async (imageUuid: string): Promise<void> => {
      setViewerImageUuid(imageUuid);
      setViewerLoading(true);
      setViewerResult(null);
      setReviewComment('');
      setReviewVerdict(ProjectDetectionReviewVerdict_Value.Unknown);
      try {
        const data = await getImageDetectionResult({
          image_uuid: imageUuid,
          project_id: projectId,
        });
        setViewerResult(data);
        const taskUuid = data.status?.main_task_uuid ?? '';
        if (taskUuid !== '') {
          const review = await getProjectDetectionReview({
            project_id: projectId,
            task_uuid: taskUuid,
          });
          setReviewComment(review.review?.comment ?? '');
          setReviewVerdict(review.review?.verdict ?? ProjectDetectionReviewVerdict_Value.Unknown);
        }
      } catch (error: unknown) {
        setViewerImageUuid(null);
        showError(getErrorMessage(error));
      } finally {
        setViewerLoading(false);
      }
    },
    [projectId, showError],
  );

  const saveReview = useCallback(
    async (verdict = reviewVerdict, comment = reviewComment): Promise<void> => {
      const taskUuid = viewerResult?.status?.main_task_uuid ?? '';
      if (projectId === '' || taskUuid === '') {
        return;
      }
      setReviewSaving(true);
      try {
        const data = await updateProjectDetectionReview({
          comment,
          project_id: projectId,
          task_uuid: taskUuid,
          verdict,
        });
        setReviewComment(data.review?.comment ?? comment);
        setReviewVerdict(data.review?.verdict ?? verdict);
      } catch (error: unknown) {
        showError(getErrorMessage(error));
      } finally {
        setReviewSaving(false);
      }
    },
    [projectId, reviewComment, reviewVerdict, showError, viewerResult],
  );

  const handleComplete = useCallback(async (): Promise<void> => {
    setAdvancing(true);
    try {
      await advanceProject({
        from_progress: ProjectProgress_Value.DetectionFinished,
        project_id: projectId,
        to_progress: ProjectProgress_Value.ReviewFinished,
      });
      onProjectChange({
        ...project,
        progress: ProjectProgress_Value.ReviewFinished,
      });
      onProgressChange(ProjectProgress_Value.ReviewFinished);
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

  const pageLoading = loading || assetsLoading || detectionLoading;
  const actionsHidden = loading || initialPageLoading;
  const readOnly = project.progress > ProjectProgress_Value.DetectionFinished;

  return (
    <div className="flex min-h-0 flex-1 flex-col rounded-lg border border-slate-200 bg-white p-5">
      {contextHolder}
      <div className="mb-4 flex items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <h2 className="text-base font-semibold text-slate-900">人工复核确认</h2>
            {!pageLoading ? (
              <span className="rounded bg-blue-50 px-2 py-0.5 text-xs text-[#1677FF]">
                {groups.length} 组 {uploadedImages.length} 张图像
              </span>
            ) : null}
          </div>
          <p className="mt-1 text-sm text-slate-500">
            复核 Agent 智能检测结果，引入专家判断，按需补充评论并标记检测结果是否准确
          </p>
        </div>
        {actionsHidden ? null : (
          <div className="flex shrink-0 items-center gap-3">
            <Button
              icon={<MenuFoldOutlined />}
              onClick={() => {
                setCollapsedGroupIds(new Set(groups.map((group) => group.id)));
              }}
            >
              全部收起
            </Button>
            <Button
              icon={<MenuUnfoldOutlined />}
              onClick={() => {
                setCollapsedGroupIds(new Set());
              }}
            >
              全部展开
            </Button>
            {readOnly ? null : (
              <>
                <Button
                  disabled={pageLoading || advancing}
                  icon={<ReloadOutlined />}
                  onClick={() => void refreshPage()}
                >
                  刷新
                </Button>
                <Button
                  disabled={pageLoading}
                  icon={<StepForwardOutlined />}
                  loading={advancing}
                  onClick={() => void handleComplete()}
                  type="primary"
                >
                  完成并进入下一步
                </Button>
              </>
            )}
          </div>
        )}
      </div>
      <div className="min-h-0 flex-1 overflow-y-auto pr-1">
        {pageLoading ? (
          <div className="flex h-full items-center justify-center">
            <Spin description="正在加载人工复核数据" />
          </div>
        ) : null}
        {!pageLoading && groups.length === DEFAULT_TASK_COUNT ? (
          <div className="flex h-full items-center justify-center rounded-lg border border-dashed border-slate-200">
            <Empty description="暂无检测结果" />
          </div>
        ) : null}
        {!pageLoading && groups.length > EMPTY_ITEMS_COUNT ? (
          <div className="space-y-4">
            {groups.map((group) => (
              <ProjectDetectionGroup
                collapsed={collapsedGroupIds.has(group.id)}
                group={group}
                key={group.id}
                onOpenImageViewer={(imageUuid) => {
                  void openReviewViewer(imageUuid);
                }}
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
          void openReviewViewer(imageUuid);
        }}
        open={viewerImageUuid !== null}
        result={viewerResult}
        review={{
          comment: reviewComment,
          enabled: true,
          onCommentChange: setReviewComment,
          onSave: () => {
            void saveReview();
          },
          onVerdictChange: (verdict) => {
            setReviewVerdict(verdict);
            void saveReview(verdict, reviewComment);
          },
          readOnly,
          saving: reviewSaving,
          verdict: reviewVerdict,
        }}
        uploadedImages={uploadedImages}
        viewerIndex={viewerIndex}
      />
      <ProjectDetectionProgressViewer
        onClose={() => {
          setProgressImageUuid(null);
        }}
        open={progressImageUuid !== null && progressTask !== undefined}
        task={progressTask}
      />
    </div>
  );
}
