import type { ReactElement, ReactNode } from 'react';

export function AuthShell({
  children,
  title,
  subtitle,
}: {
  children: ReactNode;
  title: string;
  subtitle: string;
}): ReactElement {
  return (
    <div className="auth-shell flex items-center justify-center px-4 py-10">
      <div className="w-full max-w-[1040px] overflow-hidden rounded-lg border border-slate-200 bg-white shadow-sm">
        <div className="grid min-h-[620px] grid-cols-1 md:grid-cols-[1fr_440px]">
          <div className="relative hidden overflow-hidden flex-col justify-between bg-slate-900 p-10 text-white md:flex">
            <img
              alt="tongji-logo"
              aria-hidden="true"
              className="pointer-events-none absolute -bottom-28 -left-28 h-96 w-96 select-none opacity-[0.08]"
              src="/tongji-logo.svg"
            />
            <div className="relative z-10">
              <div className="text-sm font-semibold tracking-wide text-cyan-200">
                Building Curtain Wall Resilience Assessment
              </div>
              <h1 className="mt-6 max-w-lg text-3xl font-semibold leading-tight">建筑幕墙韧性评估软件系统</h1>
              <p className="mt-4 max-w-md text-sm leading-6 text-slate-300">
                从单点缺陷识别到整栋建筑韧性评估，构建面向智慧运维的决策闭环。
              </p>
            </div>
            <div className="relative z-10 grid grid-cols-3 gap-3 text-xs text-slate-300">
              <div className="rounded border border-slate-700 p-3">海量幕墙图像处理</div>
              <div className="rounded border border-slate-700 p-3">AI Agent 多模型推理</div>
              <div className="rounded border border-slate-700 p-3">结构化专业评估报告</div>
            </div>
          </div>
          <div className="flex items-center p-8">
            <div className="w-full">
              <div className="mb-8">
                <h2 className="text-2xl font-semibold text-slate-900">{title}</h2>
                <p className="mt-2 text-sm text-slate-500">{subtitle}</p>
              </div>
              {children}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
