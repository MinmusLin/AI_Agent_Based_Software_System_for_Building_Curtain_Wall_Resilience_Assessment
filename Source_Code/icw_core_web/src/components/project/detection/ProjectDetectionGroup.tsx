import { CaretDownOutlined } from '@ant-design/icons';
import { Button } from 'antd';
import type { ReactElement } from 'react';

import { ProjectDetectionImageCard } from '@/components/project/detection/ProjectDetectionImageCard';
import { PROJECT_DETECTION_MAIN_STATUS_FAILED } from '@/types/common';
import type { ProjectGroup } from '@/types/project/assets';
import { FIRST_INDEX, NEXT_INDEX_OFFSET } from '@/utils/assetsStage';
import { isDetectionTaskRunning, type ProjectDetectionStatusMap } from '@/utils/detectionStage';

interface ProjectDetectionGroupProps {
  collapsed: boolean;
  group: ProjectGroup;
  onOpenImageViewer: (imageUuid: string) => void;
  onOpenProgressViewer: (imageUuid: string) => void;
  onToggleCollapsed: (groupId: string) => void;
  taskMap: ProjectDetectionStatusMap;
}

export function ProjectDetectionGroup({
  collapsed,
  group,
  onOpenImageViewer,
  onOpenProgressViewer,
  onToggleCollapsed,
  taskMap,
}: ProjectDetectionGroupProps): ReactElement {
  const failedImageCount = group.images.reduce((count, image) => {
    return (
      count +
      (taskMap[image.uuid]?.main_status === PROJECT_DETECTION_MAIN_STATUS_FAILED ? NEXT_INDEX_OFFSET : FIRST_INDEX)
    );
  }, FIRST_INDEX);
  const runningImageCount = group.images.reduce((count, image) => {
    return count + (isDetectionTaskRunning(taskMap[image.uuid]) ? NEXT_INDEX_OFFSET : FIRST_INDEX);
  }, FIRST_INDEX);

  return (
    <section className="overflow-hidden rounded-lg border border-slate-200 bg-slate-50">
      <div className="flex items-center gap-3 rounded-t-lg border-b border-slate-200 bg-white px-4 py-3">
        <Button
          aria-label={collapsed ? '展开图像组' : '折叠图像组'}
          icon={
            <CaretDownOutlined
              className={`transition-transform duration-300 ${collapsed ? '-rotate-90' : 'rotate-0'}`}
            />
          }
          onClick={() => {
            onToggleCollapsed(group.id);
          }}
          shape="circle"
          size="small"
          type="text"
        />
        <span className="min-w-0 truncate rounded px-2 text-sm font-semibold leading-7 text-slate-900">
          {group.name}
        </span>
        <span className="rounded bg-blue-50 px-2 py-0.5 text-xs text-[#1677FF]">{group.images.length} 张图像</span>
        {failedImageCount > FIRST_INDEX ? (
          <span className="rounded bg-red-50 px-2 py-0.5 text-xs text-red-500">{failedImageCount} 张失败</span>
        ) : null}
        {runningImageCount > FIRST_INDEX ? (
          <span className="rounded bg-emerald-50 px-2 py-0.5 text-xs text-emerald-600">
            {runningImageCount} 张检测中
          </span>
        ) : null}
      </div>
      <div
        className={`grid transition-[grid-template-rows,opacity] duration-300 ease-in-out ${
          collapsed ? 'grid-rows-[0fr] opacity-0' : 'grid-rows-[1fr] opacity-100'
        }`}
      >
        <div className="min-h-0 overflow-hidden">
          <div className="grid grid-cols-[repeat(auto-fill,minmax(88px,1fr))] gap-3 p-4">
            {group.images.map((image) => (
              <ProjectDetectionImageCard
                image={image}
                key={image.uuid}
                onOpenImageViewer={onOpenImageViewer}
                onOpenProgressViewer={onOpenProgressViewer}
                task={taskMap[image.uuid]}
              />
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
