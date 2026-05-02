import { DeleteOutlined, LoadingOutlined, PlusOutlined, SaveOutlined, StepForwardOutlined } from '@ant-design/icons';
import { Button, Input, message, Select } from 'antd';
import type { ChangeEvent, ReactElement } from 'react';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { getErrorMessage } from '@/api/http';
import { advanceProject } from '@/api/project/core';
import {
  deleteProjectThumbnail,
  getProjectProfile,
  updateProjectProfile,
  uploadProjectThumbnail,
} from '@/api/project/profile';
import type { Project } from '@/types/project/core';
import type { UpdateProjectProfileRequest } from '@/types/project/profile';
import {
  isAllowedProjectThumbnailFile,
  PROJECT_THUMBNAIL_OUTPUT_CONTENT_TYPE,
  resizeProjectThumbnailToPng,
} from '@/utils/projectThumbnail';

const { TextArea } = Input;

const PROFILE_PROGRESS = 0;
const NEXT_PROFILE_PROGRESS = 1;
const MIN_BUILT_YEAR = 1800;
const MAX_PROJECT_NAME_LENGTH = 32;
const MAX_BUILDING_NAME_LENGTH = 32;
const MAX_BUILDING_LOCATION_LENGTH = 128;
const MAX_PROJECT_TEXT_LENGTH = 5000;

interface ProfileFormState {
  name: string;
  buildingName: string;
  buildingLocation: string;
  builtYear: number;
  buildingDescription: string;
  knownIssues: string;
  assessmentGoal: string;
}

function projectToForm(project: Project): ProfileFormState {
  return {
    name: project.name,
    buildingName: project.building_name,
    buildingLocation: project.building_location,
    builtYear: project.built_year,
    buildingDescription: project.building_description,
    knownIssues: project.known_issues,
    assessmentGoal: project.assessment_goal,
  };
}

function formToPayload(projectId: string, form: ProfileFormState): UpdateProjectProfileRequest {
  return {
    project_id: projectId,
    name: form.name,
    building_name: form.buildingName,
    building_location: form.buildingLocation,
    built_year: form.builtYear,
    building_description: form.buildingDescription,
    known_issues: form.knownIssues,
    assessment_goal: form.assessmentGoal,
  };
}

function formatDateTime(value: string): string {
  if (!value) {
    return '-';
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return new Intl.DateTimeFormat('zh-CN', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(date);
}

function buildYearOptions(): { label: string; value: number }[] {
  const currentYear = new Date().getFullYear();
  const options: { label: string; value: number }[] = [];
  for (let year = currentYear; year >= MIN_BUILT_YEAR; year -= 1) {
    options.push({
      label: `${String(year)} 年`,
      value: year,
    });
  }
  return options;
}

interface TextAreaFieldProps {
  disabled: boolean;
  label: string;
  maxLength: number;
  onChange: (value: string) => void;
  placeholder: string;
  readOnly: boolean;
  value: string;
}

function TextAreaField({
  disabled,
  label,
  maxLength,
  onChange,
  placeholder,
  readOnly,
  value,
}: TextAreaFieldProps): ReactElement {
  return (
    <label className="flex min-h-0 flex-col">
      <span className="mb-2 text-sm font-medium text-slate-700">{label}</span>
      <TextArea
        className="min-h-0 flex-1 resize-none"
        disabled={disabled}
        maxLength={maxLength}
        onChange={(event: ChangeEvent<HTMLTextAreaElement>) => {
          onChange(event.target.value);
        }}
        placeholder={placeholder}
        readOnly={readOnly}
        value={value}
      />
    </label>
  );
}

interface ProjectThumbnailControlProps {
  disabled: boolean;
  onChange: (project: Project) => void;
  project: Project;
  projectId: string;
  readOnly: boolean;
}

function ProjectThumbnailControl({
  disabled,
  onChange,
  project,
  projectId,
  readOnly,
}: ProjectThumbnailControlProps): ReactElement {
  const [messageApi, contextHolder] = message.useMessage();
  const inputRef = useRef<HTMLInputElement | null>(null);
  const [busy, setBusy] = useState(false);
  const hasThumbnail = project.thumbnail_url.trim() !== '';
  const actionDisabled = disabled || readOnly || busy;

  const openFileSelector = useCallback((): void => {
    if (actionDisabled) {
      return;
    }
    inputRef.current?.click();
  }, [actionDisabled]);

  const refreshProject = useCallback(async (): Promise<void> => {
    const data = await getProjectProfile(projectId);
    onChange(data.project);
  }, [onChange, projectId]);

  const handleFileChange = useCallback(
    (event: ChangeEvent<HTMLInputElement>): void => {
      const selectedFile = event.target.files?.[0];
      event.target.value = '';
      if (!selectedFile || actionDisabled) {
        return;
      }
      if (!isAllowedProjectThumbnailFile(selectedFile)) {
        void messageApi.error('请上传 JPG、PNG 或 WebP 格式的图片');
        return;
      }

      async function upload(file: File): Promise<void> {
        setBusy(true);
        try {
          const thumbnailBlob = await resizeProjectThumbnailToPng(file);
          const uploadResult = await uploadProjectThumbnail(projectId);
          if (!uploadResult.upload_url) {
            throw new Error('project thumbnail upload url is empty');
          }

          const uploadHeaders = new Headers();
          uploadHeaders.set('Content-Type', PROJECT_THUMBNAIL_OUTPUT_CONTENT_TYPE);
          const uploadResponse = await fetch(uploadResult.upload_url, {
            body: thumbnailBlob,
            headers: uploadHeaders,
            method: 'PUT',
          });
          if (!uploadResponse.ok) {
            throw new Error('project thumbnail upload failed');
          }

          await refreshProject();
          void messageApi.success('项目缩略图已更新');
        } catch (error: unknown) {
          void messageApi.error(getErrorMessage(error));
        } finally {
          setBusy(false);
        }
      }

      void upload(selectedFile);
    },
    [actionDisabled, messageApi, projectId, refreshProject],
  );

  const handleDelete = useCallback(async (): Promise<void> => {
    if (actionDisabled || !hasThumbnail) {
      return;
    }

    setBusy(true);
    try {
      await deleteProjectThumbnail(projectId);
      await refreshProject();
      void messageApi.success('项目缩略图已删除');
    } catch (error: unknown) {
      void messageApi.error(getErrorMessage(error));
    } finally {
      setBusy(false);
    }
  }, [actionDisabled, hasThumbnail, messageApi, projectId, refreshProject]);

  return (
    <div className="flex h-full flex-col">
      {contextHolder}
      <span className="mb-2 block text-sm font-medium text-slate-700">项目缩略图</span>
      <input
        accept="image/jpeg,image/png,image/webp"
        className="hidden"
        onChange={handleFileChange}
        ref={inputRef}
        type="file"
      />
      <div className="group relative aspect-square overflow-hidden rounded-lg border border-dashed border-slate-300 bg-slate-50">
        {hasThumbnail ? (
          <button
            className="block size-full overflow-hidden rounded-lg"
            disabled={actionDisabled}
            onClick={openFileSelector}
            type="button"
          >
            <img
              alt="项目缩略图"
              className="size-full object-cover"
              onError={() => {
                onChange({
                  ...project,
                  thumbnail_url: '',
                });
              }}
              src={project.thumbnail_url}
            />
            {!readOnly ? (
              <span className="absolute inset-0 flex items-center justify-center bg-slate-950/0 text-sm font-medium text-white opacity-0 transition duration-200 group-hover:bg-slate-950/35 group-hover:opacity-100">
                更换图片
              </span>
            ) : null}
          </button>
        ) : (
          <button
            className="flex size-full flex-col items-center justify-center gap-2 text-slate-400 transition duration-200 hover:border-slate-400 hover:text-slate-600 disabled:cursor-not-allowed disabled:hover:text-slate-400"
            disabled={actionDisabled}
            onClick={openFileSelector}
            type="button"
          >
            {busy ? <LoadingOutlined className="text-xl" /> : <PlusOutlined className="text-xl" />}
            <span className="text-sm">{busy ? '上传中' : '上传缩略图'}</span>
          </button>
        )}
        {hasThumbnail && !readOnly ? (
          <Button
            aria-label="删除项目缩略图"
            className="absolute right-2 top-2 opacity-0 shadow-sm transition-opacity duration-200 group-hover:opacity-100"
            danger
            disabled={actionDisabled}
            icon={busy ? <LoadingOutlined /> : <DeleteOutlined />}
            onClick={(event) => {
              event.stopPropagation();
              void handleDelete();
            }}
            shape="circle"
            size="small"
            title="删除项目缩略图"
            type="text"
          />
        ) : null}
      </div>
    </div>
  );
}

interface ProjectProfileStageProps {
  loading?: boolean;
  onProgressChange: (progress: number) => void;
  onProjectChange: (project: Project) => void;
  project: Project;
  projectId: string;
  selectedProgress: number;
}

export function ProjectProfileStage({
  loading = false,
  onProgressChange,
  onProjectChange,
  project,
  projectId,
  selectedProgress,
}: ProjectProfileStageProps): ReactElement {
  const [messageApi, contextHolder] = message.useMessage();
  const [form, setForm] = useState<ProfileFormState>(() => projectToForm(project));
  const [saving, setSaving] = useState(false);
  const [advancing, setAdvancing] = useState(false);
  const yearOptions = useMemo(() => buildYearOptions(), []);
  const createdAtText = useMemo(() => formatDateTime(project.created_at), [project.created_at]);
  const updatedAtText = useMemo(() => formatDateTime(project.updated_at), [project.updated_at]);
  const readOnly = loading || project.progress > PROFILE_PROGRESS || selectedProgress !== PROFILE_PROGRESS;
  const canComplete = loading || (project.progress === PROFILE_PROGRESS && selectedProgress === PROFILE_PROGRESS);

  const updateFormField = useCallback(
    <Key extends keyof ProfileFormState>(key: Key, value: ProfileFormState[Key]): void => {
      setForm((currentForm) => ({
        ...currentForm,
        [key]: value,
      }));
    },
    [],
  );

  const handleSave = useCallback(async (): Promise<void> => {
    setSaving(true);
    try {
      const data = await updateProjectProfile(formToPayload(projectId, form));
      onProjectChange(data.project);
      setForm(projectToForm(data.project));
      void messageApi.success('保存成功');
    } catch (error: unknown) {
      void messageApi.error(getErrorMessage(error));
    } finally {
      setSaving(false);
    }
  }, [form, messageApi, onProjectChange, projectId]);

  const handleComplete = useCallback(async (): Promise<void> => {
    setAdvancing(true);
    try {
      await advanceProject({
        project_id: projectId,
        from_progress: project.progress,
        to_progress: NEXT_PROFILE_PROGRESS,
      });
      onProjectChange({
        ...project,
        progress: NEXT_PROFILE_PROGRESS,
      });
      onProgressChange(NEXT_PROFILE_PROGRESS);
    } catch (error: unknown) {
      void messageApi.error(getErrorMessage(error));
    } finally {
      setAdvancing(false);
    }
  }, [messageApi, onProgressChange, onProjectChange, project, projectId]);

  useEffect(() => {
    setForm(projectToForm(project));
  }, [project]);

  return (
    <div className="flex min-h-0 flex-1 flex-col rounded-lg border border-slate-200 bg-white p-5">
      {contextHolder}
      <div className="grid grid-cols-[minmax(0,1fr)_180px] gap-5">
        <div className="min-w-0">
          <label className="block">
            <span className="mb-2 block text-sm font-medium text-slate-700">项目名称</span>
            <Input
              disabled={loading}
              maxLength={MAX_PROJECT_NAME_LENGTH}
              onChange={(event: ChangeEvent<HTMLInputElement>) => {
                updateFormField('name', event.target.value);
              }}
              placeholder="请输入项目名称"
              readOnly={readOnly}
              value={form.name}
            />
          </label>
          <div className="mt-4 grid grid-cols-[1fr_180px] gap-4">
            <label className="block">
              <span className="mb-2 block text-sm font-medium text-slate-700">建筑名称</span>
              <Input
                disabled={loading}
                maxLength={MAX_BUILDING_NAME_LENGTH}
                onChange={(event: ChangeEvent<HTMLInputElement>) => {
                  updateFormField('buildingName', event.target.value);
                }}
                placeholder="请输入建筑名称"
                readOnly={readOnly}
                value={form.buildingName}
              />
            </label>
            <label className="block">
              <span className="mb-2 block text-sm font-medium text-slate-700">建成年份</span>
              <Select
                allowClear
                className="w-full"
                disabled={readOnly}
                onChange={(value: number | undefined) => {
                  updateFormField('builtYear', value ?? 0);
                }}
                options={yearOptions}
                placeholder="选择年份"
                value={form.builtYear > 0 ? form.builtYear : undefined}
              />
            </label>
          </div>
          <label className="mt-4 block">
            <span className="mb-2 block text-sm font-medium text-slate-700">建筑地址</span>
            <Input
              disabled={loading}
              maxLength={MAX_BUILDING_LOCATION_LENGTH}
              onChange={(event: ChangeEvent<HTMLInputElement>) => {
                updateFormField('buildingLocation', event.target.value);
              }}
              placeholder="请输入建筑地址"
              readOnly={readOnly}
              value={form.buildingLocation}
            />
          </label>
        </div>
        <ProjectThumbnailControl
          disabled={loading}
          onChange={onProjectChange}
          project={project}
          projectId={projectId}
          readOnly={readOnly}
        />
      </div>
      <div className="mt-5 grid min-h-0 flex-1 grid-cols-3 gap-4">
        <TextAreaField
          disabled={loading}
          label="建筑描述"
          maxLength={MAX_PROJECT_TEXT_LENGTH}
          onChange={(value: string) => {
            updateFormField('buildingDescription', value);
          }}
          placeholder="请输入建筑描述"
          readOnly={readOnly}
          value={form.buildingDescription}
        />
        <TextAreaField
          disabled={loading}
          label="已知问题"
          maxLength={MAX_PROJECT_TEXT_LENGTH}
          onChange={(value: string) => {
            updateFormField('knownIssues', value);
          }}
          placeholder="请输入已知问题"
          readOnly={readOnly}
          value={form.knownIssues}
        />
        <TextAreaField
          disabled={loading}
          label="评估目标"
          maxLength={MAX_PROJECT_TEXT_LENGTH}
          onChange={(value: string) => {
            updateFormField('assessmentGoal', value);
          }}
          placeholder="请输入评估目标"
          readOnly={readOnly}
          value={form.assessmentGoal}
        />
      </div>
      <div className="mt-5 flex items-center justify-between gap-4">
        <div className="flex flex-wrap items-center gap-x-5 gap-y-1 text-xs text-slate-500">
          <span>创建时间：{createdAtText}</span>
          <span>上次保存时间：{updatedAtText}</span>
        </div>
        {canComplete ? (
          <div className="flex shrink-0 justify-end gap-3">
            <Button disabled={loading} icon={<SaveOutlined />} loading={saving} onClick={() => void handleSave()}>
              保存
            </Button>
            <Button
              disabled={loading}
              icon={<StepForwardOutlined />}
              loading={advancing}
              onClick={() => void handleComplete()}
              type="primary"
            >
              完成并进入下一步
            </Button>
          </div>
        ) : null}
      </div>
    </div>
  );
}
