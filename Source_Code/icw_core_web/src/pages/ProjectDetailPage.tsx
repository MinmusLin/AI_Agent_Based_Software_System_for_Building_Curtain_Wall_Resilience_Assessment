import { AuditOutlined, FileTextOutlined, PictureOutlined, ProfileOutlined, RobotOutlined } from '@ant-design/icons';
import { Empty, message, Spin, Steps } from 'antd';
import type { ReactElement } from 'react';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

import { getErrorMessage } from '@/api/http';
import { getProjectProfile } from '@/api/project/profile';
import { ProjectProfileStage } from '@/components/project/ProjectProfileStage';
import { LAST_VISIBLE_PROGRESS, progressFromStageKey, PROJECT_STAGES, stageKeyFromProgress } from '@/constants/project';
import type { Project } from '@/types/project/core';

const PROFILE_PROGRESS = 0;
const PROJECT_STAGE_ICONS: ReactElement[] = [
  <ProfileOutlined key="profile" />,
  <PictureOutlined key="assets" />,
  <RobotOutlined key="detection" />,
  <AuditOutlined key="review" />,
  <FileTextOutlined key="report" />,
];

function selectedStageIcon(index: number): ReactElement {
  return (
    <span className="inline-flex size-8 items-center justify-center rounded-full bg-[#1677FF] text-base text-white shadow-sm">
      {PROJECT_STAGE_ICONS[index]}
    </span>
  );
}

function emptyProject(projectId: string): Project {
  return {
    id: projectId,
    name: '',
    building_name: '',
    building_location: '',
    built_year: 0,
    building_description: '',
    known_issues: '',
    assessment_goal: '',
    thumbnail_url: '',
    progress: PROFILE_PROGRESS,
    created_at: '',
    updated_at: '',
  };
}

function currentVisibleProgress(progress: number): number {
  return Math.min(progress, LAST_VISIBLE_PROGRESS);
}

function projectRoute(projectId: string, progress: number): string {
  return `/projects/${projectId}/${stageKeyFromProgress(progress)}`;
}

interface ProjectStageContentProps {
  loading?: boolean;
  onProgressChange: (progress: number) => void;
  onProjectChange: (project: Project) => void;
  project: Project;
  projectId: string;
  selectedProgress: number;
}

function ProjectStageContent({
  loading = false,
  onProgressChange,
  onProjectChange,
  project,
  projectId,
  selectedProgress,
}: ProjectStageContentProps): ReactElement {
  if (selectedProgress === PROFILE_PROGRESS) {
    return (
      <ProjectProfileStage
        loading={loading}
        onProgressChange={onProgressChange}
        onProjectChange={onProjectChange}
        project={project}
        projectId={projectId}
        selectedProgress={selectedProgress}
      />
    );
  }

  const stage = PROJECT_STAGES[currentVisibleProgress(selectedProgress)];
  return (
    <div className="flex min-h-0 flex-1 items-center justify-center rounded-lg border border-slate-200 bg-white">
      <Empty description={`${stage.title}：${stage.description}`} />
    </div>
  );
}

export default function ProjectDetailPage(): ReactElement {
  const navigate = useNavigate();
  const params = useParams<{ projectId: string; stage: string }>();
  const [messageApi, contextHolder] = message.useMessage();
  const [project, setProject] = useState<Project | null>(null);
  const [loading, setLoading] = useState(true);

  const projectId = useMemo(() => params.projectId?.trim() ?? '', [params.projectId]);
  const routeProgress = useMemo(() => progressFromStageKey(params.stage), [params.stage]);
  const selectedProgress = routeProgress ?? PROFILE_PROGRESS;

  const loadProject = useCallback(async (): Promise<void> => {
    if (projectId === '') {
      void navigate('/projects', { replace: true });
      return;
    }

    setLoading(true);
    try {
      const data = await getProjectProfile(projectId);
      setProject(data.project);
    } catch (error: unknown) {
      void messageApi.error(getErrorMessage(error));
      void navigate('/projects', { replace: true });
    } finally {
      setLoading(false);
    }
  }, [messageApi, navigate, projectId]);

  const handleStepChange = useCallback(
    (nextProgress: number): void => {
      if (!project || projectId === '' || nextProgress > project.progress) {
        return;
      }
      void navigate(projectRoute(projectId, nextProgress));
    },
    [navigate, project, projectId],
  );

  const handleProgressChange = useCallback(
    (nextProgress: number): void => {
      if (projectId === '') {
        return;
      }
      void navigate(projectRoute(projectId, currentVisibleProgress(nextProgress)));
    },
    [navigate, projectId],
  );

  useEffect(() => {
    if (projectId === '') {
      void navigate('/projects', { replace: true });
      return;
    }
    void loadProject();
  }, [loadProject, navigate, projectId]);

  useEffect(() => {
    if (!project || projectId === '') {
      return;
    }

    const maxOpenProgress = currentVisibleProgress(project.progress);
    if (routeProgress === null || selectedProgress > maxOpenProgress) {
      void navigate(projectRoute(projectId, maxOpenProgress), { replace: true });
    }
  }, [navigate, project, projectId, routeProgress, selectedProgress]);

  if (loading) {
    const loadingProject = emptyProject(projectId);
    return (
      <div className="flex h-[calc(100vh-112px)] min-h-0 flex-col overflow-hidden">
        {contextHolder}
        <div className="mb-5 rounded-lg border border-slate-200 bg-white px-6 py-5">
          <Steps
            current={-1}
            items={PROJECT_STAGES.map((stage) => ({
              disabled: true,
              status: 'wait',
              title: stage.title,
            }))}
          />
          <div className="mt-4 flex items-center gap-2 border-t border-slate-100 pt-4 text-sm leading-6 text-slate-500">
            <Spin size="small" />
            <span>正在加载项目数据，请稍等</span>
          </div>
        </div>
        <ProjectStageContent
          loading
          onProgressChange={handleProgressChange}
          onProjectChange={setProject}
          project={loadingProject}
          projectId={projectId}
          selectedProgress={PROFILE_PROGRESS}
        />
      </div>
    );
  }

  if (!project || projectId === '') {
    return (
      <div className="h-[calc(100vh-112px)] overflow-hidden rounded-lg border border-slate-200 bg-white py-16">
        {contextHolder}
        <Empty description="项目不存在" />
      </div>
    );
  }

  const visibleProgress =
    routeProgress === null
      ? currentVisibleProgress(project.progress)
      : Math.min(selectedProgress, currentVisibleProgress(project.progress));
  const visibleStage = PROJECT_STAGES[visibleProgress];

  return (
    <div className="flex h-[calc(100vh-112px)] min-h-0 flex-col overflow-hidden">
      {contextHolder}
      <div className="mb-5 rounded-lg border border-slate-200 bg-white px-6 py-5">
        <Steps
          current={visibleProgress}
          items={PROJECT_STAGES.map((stage, index) => ({
            disabled: index > project.progress,
            icon: index === visibleProgress ? selectedStageIcon(index) : undefined,
            status: index === visibleProgress ? 'process' : 'wait',
            title: stage.title,
          }))}
          onChange={handleStepChange}
        />
        <p className="mt-4 border-t border-slate-100 pt-4 text-sm leading-6 text-slate-500">
          {visibleStage.description}
        </p>
      </div>
      <ProjectStageContent
        onProgressChange={handleProgressChange}
        onProjectChange={setProject}
        project={project}
        projectId={projectId}
        selectedProgress={visibleProgress}
      />
    </div>
  );
}
