import { Empty } from 'antd';
import type { ReactElement } from 'react';

interface ProjectStagePlaceholderProps {
  description: string;
  title: string;
}

export function ProjectStagePlaceholder({ description, title }: ProjectStagePlaceholderProps): ReactElement {
  return (
    <div className="flex min-h-0 flex-1 items-center justify-center rounded-lg border border-slate-200 bg-white">
      <Empty description={`${title}：${description}`} />
    </div>
  );
}
