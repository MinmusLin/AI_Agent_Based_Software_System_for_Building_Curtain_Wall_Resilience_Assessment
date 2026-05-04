import { CalendarOutlined, DeleteOutlined, EnvironmentOutlined, HomeOutlined, PlusOutlined } from '@ant-design/icons';
import { Button, Card, Empty, message } from 'antd';
import type { MouseEvent, ReactElement } from 'react';
import { useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { getErrorMessage } from '@/api/http';
import { createProject, deleteProject, listProjects } from '@/api/project/core';
import { stageKeyFromProgress } from '@/constants/project';
import type { Project, ProjectListItem } from '@/types/project/core';
import { formatDateTime } from '@/utils/datetime';

const PROJECT_CARD_CLASS = 'h-full border-slate-200 shadow-none [&_.ant-card-body]:px-4';
const PROJECT_CARD_BODY_CLASS = 'flex h-[9.8rem] flex-col';
const EMPTY_PROJECT_COUNT = 0;

function displayText(value: string, fallback: string): string {
  const trimmedValue = value.trim();
  if (trimmedValue === '') {
    return fallback;
  }
  return trimmedValue;
}

function projectToListItem(project: Project): ProjectListItem {
  return {
    building_location: project.building_location,
    building_name: project.building_name,
    created_at: project.created_at,
    id: project.id,
    name: project.name,
    thumbnail_url: project.thumbnail_url,
    progress: project.progress,
  };
}

function projectRoute(project: ProjectListItem): string {
  return `/projects/${project.id}/${stageKeyFromProgress(project.progress)}`;
}

interface SkeletonBarProps {
  className: string;
}

function SkeletonBar({ className }: SkeletonBarProps): ReactElement {
  return <span className={`block animate-pulse rounded-lg bg-slate-200 ${className}`} />;
}

function ProjectCardSkeleton(): ReactElement {
  return (
    <Card className={PROJECT_CARD_CLASS}>
      <div className={PROJECT_CARD_BODY_CLASS}>
        <div className="flex gap-4">
          <div className="min-w-0 flex-1">
            <div className="min-h-12 space-y-2 pt-1">
              <SkeletonBar className="h-5 w-3/4" />
              <SkeletonBar className="h-5 w-1/2" />
            </div>
            <div className="mt-3 flex flex-col gap-2 text-sm">
              <div className="flex min-h-6 items-center gap-2">
                <SkeletonBar className="h-4 w-4 shrink-0" />
                <SkeletonBar className="h-4 w-2/3" />
              </div>
              <div className="flex min-h-6 items-center gap-2">
                <SkeletonBar className="h-4 w-4 shrink-0" />
                <SkeletonBar className="h-4 w-full" />
              </div>
            </div>
          </div>
          <SkeletonBar className="size-20 shrink-0" />
        </div>
        <div className="mt-3 flex items-center justify-between gap-3 border-t border-slate-100 pt-3 text-xs">
          <div className="flex h-6 min-w-0 items-center gap-2">
            <SkeletonBar className="h-4 w-4 shrink-0" />
            <SkeletonBar className="h-4 w-32" />
          </div>
          <SkeletonBar className="h-6 w-6 shrink-0 rounded-full" />
        </div>
      </div>
    </Card>
  );
}

function ProjectEmptyCard({ description }: { description: string }): ReactElement {
  return (
    <Card className={`${PROJECT_CARD_CLASS} border-dashed`}>
      <div className="flex h-[9.8rem] items-center justify-center">
        <Empty description={description} />
      </div>
    </Card>
  );
}

function ProjectSectionSkeleton({ title }: { title: string }): ReactElement {
  return (
    <section>
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-base font-semibold text-slate-900">{title}</h2>
        <SkeletonBar className="h-4 w-14" />
      </div>
      <div className="grid grid-cols-1 gap-4 lg:grid-cols-3">
        {Array.from({ length: 3 }, (_, index) => (
          <ProjectCardSkeleton key={index} />
        ))}
      </div>
    </section>
  );
}

function ProjectsPageSkeleton(): ReactElement {
  return (
    <div className="space-y-8">
      <ProjectSectionSkeleton title="进行中的项目" />
      <ProjectSectionSkeleton title="已完成的项目" />
    </div>
  );
}

interface ProjectCardProps {
  deleting: boolean;
  onDelete: (project: ProjectListItem) => void;
  project: ProjectListItem;
  onOpen: (project: ProjectListItem) => void;
}

function ProjectCard({ deleting, onDelete, onOpen, project }: ProjectCardProps): ReactElement {
  const projectName = displayText(project.name, '未命名项目');
  const buildingName = displayText(project.building_name, '未填写');
  const buildingLocation = displayText(project.building_location, '未填写');
  const hasThumbnail = project.thumbnail_url.trim() !== '';

  return (
    <div
      className="project-card-motion group h-full cursor-pointer"
      onClick={() => {
        onOpen(project);
      }}
    >
      <Card className={PROJECT_CARD_CLASS}>
        <div className={PROJECT_CARD_BODY_CLASS}>
          <div className="flex gap-4">
            <div className="min-w-0 flex-1">
              <div className="min-h-12">
                <h3 className="line-clamp-2 text-base font-semibold leading-6 text-slate-900">{projectName}</h3>
              </div>
              <div className="mt-3 flex flex-col gap-2 text-sm text-slate-600">
                <div className="flex min-h-6 min-w-0 items-center gap-2">
                  <HomeOutlined className="shrink-0 text-slate-400" />
                  <span className="truncate">建筑名称：{buildingName}</span>
                </div>
                <div className="flex min-h-6 min-w-0 items-center gap-2">
                  <EnvironmentOutlined className="shrink-0 text-slate-400" />
                  <span className="truncate">建筑地址：{buildingLocation}</span>
                </div>
              </div>
            </div>
            {hasThumbnail ? (
              <img alt={projectName} className="size-20 shrink-0 rounded-lg object-cover" src={project.thumbnail_url} />
            ) : null}
          </div>
          <div className="mt-3 flex items-center justify-between gap-3 border-t border-slate-100 pt-3 text-xs text-slate-500">
            <div className="flex h-6 min-w-0 items-center gap-2 leading-6">
              <CalendarOutlined className="shrink-0 leading-none" />
              <span className="truncate leading-6">创建于 {formatDateTime(project.created_at)}</span>
            </div>
            <Button
              aria-label="删除项目"
              className="shrink-0 opacity-0 shadow-sm transition-opacity duration-200 group-hover:opacity-100"
              danger
              icon={<DeleteOutlined />}
              loading={deleting}
              onClick={(event: MouseEvent<HTMLElement>) => {
                event.stopPropagation();
                onDelete(project);
              }}
              shape="circle"
              size="small"
              title="删除项目"
              type="text"
            />
          </div>
        </div>
      </Card>
    </div>
  );
}

interface ProjectSectionProps {
  deletingProjectId: string;
  emptyDescription: string;
  onDelete: (project: ProjectListItem) => void;
  onOpen: (project: ProjectListItem) => void;
  projects: ProjectListItem[];
  title: string;
}

function ProjectSection({
  deletingProjectId,
  emptyDescription,
  onDelete,
  onOpen,
  projects,
  title,
}: ProjectSectionProps): ReactElement {
  return (
    <section>
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-base font-semibold text-slate-900">{title}</h2>
        <span className="text-sm text-slate-500">{projects.length} 个项目</span>
      </div>
      {projects.length > EMPTY_PROJECT_COUNT ? (
        <div className="grid grid-cols-1 gap-4 lg:grid-cols-3">
          {projects.map((project) => (
            <ProjectCard
              deleting={deletingProjectId === project.id}
              key={project.id}
              onDelete={onDelete}
              onOpen={onOpen}
              project={project}
            />
          ))}
        </div>
      ) : (
        <ProjectEmptyCard description={emptyDescription} />
      )}
    </section>
  );
}

export default function ProjectsPage(): ReactElement {
  const navigate = useNavigate();
  const [messageApi, contextHolder] = message.useMessage();
  const [activeProjects, setActiveProjects] = useState<ProjectListItem[]>([]);
  const [completedProjects, setCompletedProjects] = useState<ProjectListItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [deletingProjectId, setDeletingProjectId] = useState('');

  const openProject = useCallback(
    (project: ProjectListItem): void => {
      void navigate(projectRoute(project));
    },
    [navigate],
  );

  const loadProjects = useCallback(async (): Promise<void> => {
    setLoading(true);
    try {
      const data = await listProjects();
      setActiveProjects(data.active_projects);
      setCompletedProjects(data.completed_projects);
    } catch (error: unknown) {
      void messageApi.error(getErrorMessage(error));
    } finally {
      setLoading(false);
    }
  }, [messageApi]);

  const handleCreateProject = useCallback(async (): Promise<void> => {
    setCreating(true);
    try {
      const data = await createProject();
      const newProject = projectToListItem(data.project);
      setActiveProjects((projects) => [...projects, newProject]);
      void navigate(projectRoute(newProject));
    } catch (error: unknown) {
      void messageApi.error(getErrorMessage(error));
    } finally {
      setCreating(false);
    }
  }, [messageApi, navigate]);

  const handleDeleteProject = useCallback(
    async (project: ProjectListItem): Promise<void> => {
      setDeletingProjectId(project.id);
      try {
        const data = await deleteProject(project.id);
        setActiveProjects(data.active_projects);
        setCompletedProjects(data.completed_projects);
        void messageApi.success('项目删除成功');
      } catch (error: unknown) {
        void messageApi.error(getErrorMessage(error));
      } finally {
        setDeletingProjectId('');
      }
    },
    [messageApi],
  );

  useEffect(() => {
    void loadProjects();
  }, [loadProjects]);

  return (
    <div>
      {contextHolder}
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-xl font-semibold text-slate-900">项目管理</h1>
          <p className="mt-1 text-sm text-slate-500">
            管理建筑幕墙韧性评估项目，继续处理进行中的项目，并查看已完成的项目资料
          </p>
        </div>
        <Button icon={<PlusOutlined />} loading={creating} onClick={() => void handleCreateProject()} type="primary">
          新建项目
        </Button>
      </div>
      {loading ? (
        <ProjectsPageSkeleton />
      ) : (
        <div className="space-y-8">
          <ProjectSection
            deletingProjectId={deletingProjectId}
            emptyDescription="暂无进行中的项目"
            onDelete={(project: ProjectListItem) => void handleDeleteProject(project)}
            onOpen={openProject}
            projects={activeProjects}
            title="进行中的项目"
          />
          <ProjectSection
            deletingProjectId={deletingProjectId}
            emptyDescription="暂无已完成的项目"
            onDelete={(project: ProjectListItem) => void handleDeleteProject(project)}
            onOpen={openProject}
            projects={completedProjects}
            title="已完成的项目"
          />
        </div>
      )}
    </div>
  );
}
