import { CloseCircleFilled, EyeOutlined, FileImageOutlined, FundProjectionScreenOutlined } from '@ant-design/icons';
import { Button, Tooltip } from 'antd';
import type { ReactElement } from 'react';

import {
  type ProjectDetectionStatus,
  ProjectDetectionTaskStatus_Value,
  type ProjectImage,
  ProjectImageStatus_Value,
} from '@/gen/core/common';
import { detectionTaskProgressPercent, isDetectionTaskRunning, isDetectionTaskSucceeded } from '@/utils/detectionStage';

const MIN_RUNNING_PROGRESS_PERCENT = 8;

interface ProjectDetectionImageCardProps {
  image: ProjectImage;
  onOpenImageViewer: (imageUuid: string) => void;
  onOpenProgressViewer: (imageUuid: string) => void;
  task: ProjectDetectionStatus | undefined;
}

function ImageUnavailableIcon(): ReactElement {
  return (
    <span className="relative inline-flex size-8 items-center justify-center text-slate-400">
      <FileImageOutlined className="text-3xl" />
      <CloseCircleFilled className="absolute -right-1 -bottom-1 rounded-full bg-white text-sm text-red-500" />
    </span>
  );
}

function DetectionProgressBar({ task }: { task: ProjectDetectionStatus }): ReactElement {
  const percent = Math.max(detectionTaskProgressPercent(task), MIN_RUNNING_PROGRESS_PERCENT);

  return (
    <div className="absolute inset-x-2 bottom-2 z-20" title={`${String(detectionTaskProgressPercent(task))}%`}>
      <div className="h-1.5 overflow-hidden rounded-full bg-white/60">
        <div
          className="relative h-full overflow-hidden rounded-full bg-[#1677FF] transition-[width] duration-500"
          style={{ width: `${String(percent)}%` }}
        >
          <span className="detection-progress-flow absolute inset-y-0 left-0 w-full bg-gradient-to-r from-transparent via-white/55 to-transparent" />
        </div>
      </div>
    </div>
  );
}

export function ProjectDetectionImageCard({
  image,
  onOpenImageViewer,
  onOpenProgressViewer,
  task,
}: ProjectDetectionImageCardProps): ReactElement {
  const imageReady = image.status === ProjectImageStatus_Value.Uploaded && image.thumbnail_url !== '';
  const running = isDetectionTaskRunning(task);
  const failed = task?.main_status === ProjectDetectionTaskStatus_Value.Failed;
  const shouldOpenProgress = task !== undefined && !isDetectionTaskSucceeded(task);
  const actionVisible = imageReady && !running && !failed && task !== undefined;

  return (
    <div className="group/image relative aspect-square overflow-hidden rounded-lg border border-slate-200 bg-white">
      {imageReady ? (
        <img
          alt={image.file_name}
          className="size-full object-cover transition duration-200"
          draggable={false}
          src={image.thumbnail_url}
        />
      ) : (
        <div className="flex size-full items-center justify-center text-slate-400">
          <ImageUnavailableIcon />
        </div>
      )}
      {running ? (
        <>
          <div className="absolute inset-0 z-10 bg-slate-950/25" />
          <div className="absolute inset-0 z-20 flex items-center justify-center opacity-0 transition duration-200 group-hover/image:opacity-100">
            <Tooltip title="查看进度">
              <Button
                aria-label="查看进度"
                icon={<FundProjectionScreenOutlined />}
                onClick={() => {
                  onOpenProgressViewer(image.uuid);
                }}
                shape="circle"
                size="small"
              />
            </Tooltip>
          </div>
          {task ? <DetectionProgressBar task={task} /> : null}
        </>
      ) : null}
      {failed ? (
        <>
          <div className="absolute inset-0 z-10 bg-slate-950/25" />
          <div className="absolute inset-0 z-20 flex items-center justify-center">
            <Tooltip title="查看进度">
              <Button
                aria-label="查看进度"
                danger
                icon={<FundProjectionScreenOutlined />}
                onClick={() => {
                  onOpenProgressViewer(image.uuid);
                }}
                shape="circle"
                size="small"
                type="primary"
              />
            </Tooltip>
          </div>
        </>
      ) : null}
      {actionVisible ? (
        <div className="absolute inset-0 flex items-center justify-center bg-slate-950/0 opacity-0 transition duration-200 group-hover/image:bg-slate-950/35 group-hover/image:opacity-100">
          <Tooltip title={shouldOpenProgress ? '查看进度' : '查看结果'}>
            <Button
              aria-label={shouldOpenProgress ? '查看进度' : '查看结果'}
              icon={shouldOpenProgress ? <FundProjectionScreenOutlined /> : <EyeOutlined />}
              onClick={() => {
                if (shouldOpenProgress) {
                  onOpenProgressViewer(image.uuid);
                  return;
                }
                onOpenImageViewer(image.uuid);
              }}
              shape="circle"
              size="small"
            />
          </Tooltip>
        </div>
      ) : null}
    </div>
  );
}
