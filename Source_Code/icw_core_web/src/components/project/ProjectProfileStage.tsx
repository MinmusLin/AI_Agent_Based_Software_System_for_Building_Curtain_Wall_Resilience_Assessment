import {
  DeleteOutlined,
  LoadingOutlined,
  PictureOutlined,
  ReloadOutlined,
  SaveOutlined,
  StepForwardOutlined,
  UploadOutlined,
} from '@ant-design/icons';
import { Button, Input, message, Select, Tooltip } from 'antd';
import type { ChangeEvent, ReactElement } from 'react';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { getErrorMessage } from '@/api/http';
import { advanceProject } from '@/api/project/core';
import {
  deleteProjectThumbnail,
  getProjectThumbnail,
  updateProjectProfile,
  uploadProjectThumbnail,
} from '@/api/project/profile';
import type { Project } from '@/gen/core/api/common';
import type { UpdateProjectProfileRequest } from '@/gen/core/api/project_profile';
import { ProjectProgress_Value } from '@/gen/core/common';
import { formatDateTime } from '@/utils/datetime';
import {
  isAllowedProjectThumbnailFile,
  PROJECT_THUMBNAIL_OUTPUT_CONTENT_TYPE,
  resizeProjectThumbnailToPng,
} from '@/utils/images';

const { TextArea } = Input;

const MIN_BUILT_YEAR = 1800;
const MAX_PROJECT_NAME_LENGTH = 32;
const MAX_BUILDING_NAME_LENGTH = 32;
const MAX_BUILDING_LOCATION_LENGTH = 128;
const MAX_PROJECT_TEXT_LENGTH = 5000;
const YEAR_STEP = 1;
const EMPTY_BUILT_YEAR = 0;
const LOADING_PLACEHOLDER = '加载中';
const UNKNOWN_TEXT = '未知';
const PROFILE_PLACEHOLDERS = {
  assessmentGoal: '请输入评估目标或重点关注方向',
  buildingDescription: '请输入建筑描述',
  buildingLocation: '请输入建筑地址',
  buildingName: '请输入建筑名称',
  builtYear: '请选择年份',
  createdAt: '',
  knownIssues: '请输入已知问题或人工先验描述',
  name: '请输入项目名称',
  updatedAt: '',
} as const;
const LOADING_PROFILE_PLACEHOLDERS = {
  assessmentGoal: LOADING_PLACEHOLDER,
  buildingDescription: LOADING_PLACEHOLDER,
  buildingLocation: LOADING_PLACEHOLDER,
  buildingName: LOADING_PLACEHOLDER,
  builtYear: LOADING_PLACEHOLDER,
  createdAt: LOADING_PLACEHOLDER,
  knownIssues: LOADING_PLACEHOLDER,
  name: LOADING_PLACEHOLDER,
  updatedAt: LOADING_PLACEHOLDER,
} as const;
const READONLY_EMPTY_PROFILE_PLACEHOLDERS = {
  assessmentGoal: '评估目标或重点关注方向为空',
  buildingDescription: '建筑描述为空',
  buildingLocation: '建筑地址为空',
  buildingName: '建筑名称为空',
  builtYear: UNKNOWN_TEXT,
  createdAt: '项目创建时间为空',
  knownIssues: '已知问题或人工先验描述为空',
  name: '项目名称为空',
  updatedAt: '上次保存时间为空',
} as const;

type ProfilePlaceholders = Record<keyof typeof PROFILE_PLACEHOLDERS, string>;

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

function buildYearOptions(): { label: string; value: number }[] {
  const currentYear = new Date().getFullYear();
  const options: { label: string; value: number }[] = [];
  for (let year = currentYear; year >= MIN_BUILT_YEAR; year -= YEAR_STEP) {
    options.push({
      label: `${String(year)} 年`,
      value: year,
    });
  }
  return options;
}

function getProfilePlaceholders(loading: boolean, readOnly: boolean): ProfilePlaceholders {
  if (loading) {
    return LOADING_PROFILE_PLACEHOLDERS;
  }
  if (readOnly) {
    return READONLY_EMPTY_PROFILE_PLACEHOLDERS;
  }
  return PROFILE_PLACEHOLDERS;
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

  const refreshThumbnail = useCallback(async (): Promise<void> => {
    const data = await getProjectThumbnail(projectId);
    onChange({
      ...project,
      thumbnail_url: data.thumbnail_url,
    });
  }, [onChange, project, projectId]);

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

          await refreshThumbnail();
        } catch (error: unknown) {
          void messageApi.error(getErrorMessage(error));
        } finally {
          setBusy(false);
        }
      }

      void upload(selectedFile);
    },
    [actionDisabled, messageApi, projectId, refreshThumbnail],
  );

  const handleDelete = useCallback(async (): Promise<void> => {
    if (actionDisabled || !hasThumbnail) {
      return;
    }

    setBusy(true);
    try {
      await deleteProjectThumbnail(projectId);
      await refreshThumbnail();
    } catch (error: unknown) {
      void messageApi.error(getErrorMessage(error));
    } finally {
      setBusy(false);
    }
  }, [actionDisabled, hasThumbnail, messageApi, projectId, refreshThumbnail]);

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
      <div
        className={`group relative aspect-square overflow-hidden rounded-lg border border-dashed bg-white transition duration-200 ${
          !hasThumbnail && !actionDisabled ? 'border-slate-300 hover:border-[#1677FF]' : 'border-slate-300'
        }`}
      >
        {disabled ? (
          <div className="flex size-full items-center justify-center text-slate-400">
            <LoadingOutlined className="text-2xl" />
          </div>
        ) : hasThumbnail ? (
          <div className="size-full overflow-hidden rounded-lg">
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
            {!readOnly && !busy ? (
              <div className="absolute inset-0 flex items-center justify-center gap-2 bg-slate-950/0 opacity-0 transition duration-200 group-hover:bg-slate-950/35 group-hover:opacity-100">
                <Tooltip title="更新缩略图">
                  <Button
                    aria-label="更新缩略图"
                    disabled={actionDisabled}
                    icon={<ReloadOutlined />}
                    onClick={openFileSelector}
                    shape="circle"
                  />
                </Tooltip>
                <Tooltip title="删除缩略图">
                  <Button
                    aria-label="删除缩略图"
                    danger
                    disabled={actionDisabled}
                    icon={<DeleteOutlined />}
                    onClick={() => void handleDelete()}
                    shape="circle"
                  />
                </Tooltip>
              </div>
            ) : null}
            {busy ? (
              <div className="absolute inset-0 flex items-center justify-center bg-slate-950/35 text-white">
                <LoadingOutlined className="text-2xl" />
              </div>
            ) : null}
          </div>
        ) : readOnly ? (
          <div className="flex size-full items-center justify-center text-slate-300">
            <PictureOutlined className="text-2xl" />
          </div>
        ) : (
          <button
            aria-label="上传缩略图"
            className={`flex size-full flex-col items-center justify-center gap-2 text-sm transition duration-200 disabled:cursor-not-allowed ${
              actionDisabled ? 'text-slate-400' : 'text-slate-500 group-hover:text-[#1677FF]'
            }`}
            disabled={actionDisabled}
            onClick={openFileSelector}
            type="button"
          >
            {busy ? <LoadingOutlined className="text-2xl" /> : <UploadOutlined className="text-2xl" />}
          </button>
        )}
      </div>
    </div>
  );
}

interface ProjectProfileStageProps {
  loading?: boolean;
  onProgressChange: (progress: ProjectProgress_Value) => void;
  onProjectChange: (project: Project) => void;
  project: Project;
  projectId: string;
  selectedProgress: ProjectProgress_Value;
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
  const readOnly =
    loading ||
    project.progress > ProjectProgress_Value.InitializationFinished ||
    selectedProgress !== ProjectProgress_Value.InitializationFinished;
  const placeholders = getProfilePlaceholders(loading, readOnly);
  const hasBuiltYear = form.builtYear > EMPTY_BUILT_YEAR;
  const builtYearText = hasBuiltYear ? `${String(form.builtYear)} 年` : '';
  const canComplete =
    !loading &&
    project.progress === ProjectProgress_Value.InitializationFinished &&
    selectedProgress === ProjectProgress_Value.InitializationFinished;

  const updateFormField = useCallback(
    <Key extends keyof ProfileFormState>(key: Key, value: ProfileFormState[Key]): void => {
      setForm((currentForm) => ({
        ...currentForm,
        [key]: value,
      }));
    },
    [],
  );

  const saveProjectProfile = useCallback(
    async (showSuccess: boolean): Promise<Project> => {
      setSaving(true);
      try {
        const data = await updateProjectProfile(formToPayload(projectId, form));
        if (!data.project) {
          throw new Error('project is empty');
        }
        onProjectChange(data.project);
        setForm(projectToForm(data.project));
        if (showSuccess) {
          void messageApi.success('保存成功');
        }
        return data.project;
      } finally {
        setSaving(false);
      }
    },
    [form, messageApi, onProjectChange, projectId],
  );

  const handleSave = useCallback(async (): Promise<void> => {
    try {
      await saveProjectProfile(true);
    } catch (error: unknown) {
      void messageApi.error(getErrorMessage(error));
    }
  }, [messageApi, saveProjectProfile]);

  const handleComplete = useCallback(async (): Promise<void> => {
    setAdvancing(true);
    try {
      const savedProject = await saveProjectProfile(false);
      await advanceProject({
        project_id: projectId,
        from_progress: savedProject.progress,
        to_progress: ProjectProgress_Value.ProfileFinished,
      });
      onProjectChange({
        ...savedProject,
        progress: ProjectProgress_Value.ProfileFinished,
      });
      onProgressChange(ProjectProgress_Value.ProfileFinished);
    } catch (error: unknown) {
      void messageApi.error(getErrorMessage(error));
    } finally {
      setAdvancing(false);
    }
  }, [messageApi, onProgressChange, onProjectChange, projectId, saveProjectProfile]);

  useEffect(() => {
    setForm(projectToForm(project));
  }, [
    project.assessment_goal,
    project.building_description,
    project.building_location,
    project.building_name,
    project.built_year,
    project.id,
    project.known_issues,
    project.name,
    project.updated_at,
  ]);

  return (
    <div className="flex min-h-0 flex-1 flex-col rounded-lg border border-slate-200 bg-white p-5">
      {contextHolder}
      <div className="grid grid-cols-[minmax(0,1fr)_184px] gap-5">
        <div className="min-w-0">
          <div className="grid grid-cols-[1fr_184px] gap-4">
            <label className="block">
              <span className="mb-2 block text-sm font-medium text-slate-700">项目名称</span>
              <Input
                disabled={loading}
                maxLength={MAX_PROJECT_NAME_LENGTH}
                onChange={(event: ChangeEvent<HTMLInputElement>) => {
                  updateFormField('name', event.target.value);
                }}
                placeholder={placeholders.name}
                readOnly={readOnly}
                value={form.name}
              />
            </label>
            <label className="block">
              <span className="mb-2 block text-sm font-medium text-slate-700">项目创建时间</span>
              <Input
                disabled={loading}
                placeholder={placeholders.createdAt}
                readOnly
                value={loading ? '' : createdAtText}
              />
            </label>
          </div>
          <div className="mt-4 grid grid-cols-[1fr_184px] gap-4">
            <label className="block">
              <span className="mb-2 block text-sm font-medium text-slate-700">建筑名称</span>
              <Input
                disabled={loading}
                maxLength={MAX_BUILDING_NAME_LENGTH}
                onChange={(event: ChangeEvent<HTMLInputElement>) => {
                  updateFormField('buildingName', event.target.value);
                }}
                placeholder={placeholders.buildingName}
                readOnly={readOnly}
                value={form.buildingName}
              />
            </label>
            <label className="block">
              <span className="mb-2 block text-sm font-medium text-slate-700">上次保存时间</span>
              <Input
                disabled={loading}
                placeholder={placeholders.updatedAt}
                readOnly
                value={loading ? '' : updatedAtText}
              />
            </label>
          </div>
          <div className="mt-4 grid grid-cols-[1fr_184px] gap-4">
            <label className="block">
              <span className="mb-2 block text-sm font-medium text-slate-700">建筑地址</span>
              <Input
                disabled={loading}
                maxLength={MAX_BUILDING_LOCATION_LENGTH}
                onChange={(event: ChangeEvent<HTMLInputElement>) => {
                  updateFormField('buildingLocation', event.target.value);
                }}
                placeholder={placeholders.buildingLocation}
                readOnly={readOnly}
                value={form.buildingLocation}
              />
            </label>
            <label className="block">
              <span className="mb-2 block text-sm font-medium text-slate-700">建筑建成年份</span>
              {readOnly ? (
                <Input
                  disabled={loading}
                  placeholder={placeholders.builtYear}
                  readOnly
                  value={loading ? '' : builtYearText}
                />
              ) : (
                <Select
                  allowClear
                  className="w-full"
                  disabled={loading}
                  onChange={(value: number | undefined) => {
                    updateFormField('builtYear', value ?? EMPTY_BUILT_YEAR);
                  }}
                  options={yearOptions}
                  placeholder={placeholders.builtYear}
                  value={form.builtYear > EMPTY_BUILT_YEAR ? form.builtYear : undefined}
                />
              )}
            </label>
          </div>
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
          placeholder={placeholders.buildingDescription}
          readOnly={readOnly}
          value={form.buildingDescription}
        />
        <TextAreaField
          disabled={loading}
          label="已知问题或人工先验描述"
          maxLength={MAX_PROJECT_TEXT_LENGTH}
          onChange={(value: string) => {
            updateFormField('knownIssues', value);
          }}
          placeholder={placeholders.knownIssues}
          readOnly={readOnly}
          value={form.knownIssues}
        />
        <TextAreaField
          disabled={loading}
          label="评估目标或重点关注方向"
          maxLength={MAX_PROJECT_TEXT_LENGTH}
          onChange={(value: string) => {
            updateFormField('assessmentGoal', value);
          }}
          placeholder={placeholders.assessmentGoal}
          readOnly={readOnly}
          value={form.assessmentGoal}
        />
      </div>
      {canComplete ? (
        <div className="mt-5 flex justify-end gap-3">
          <Button disabled={advancing} icon={<SaveOutlined />} loading={saving} onClick={() => void handleSave()}>
            保存
          </Button>
          <Button
            disabled={saving}
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
  );
}
