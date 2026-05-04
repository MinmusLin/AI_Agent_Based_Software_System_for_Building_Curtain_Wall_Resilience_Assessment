export const AVATAR_TYPE_CUSTOM = 'custom';
export const AVATAR_TYPE_DEFAULT = 'default';
export const AVATAR_TYPE_NONE = 'none';

export type AvatarType = typeof AVATAR_TYPE_CUSTOM | typeof AVATAR_TYPE_DEFAULT | typeof AVATAR_TYPE_NONE;

export const LOGIN_SCENE_PASSWORD = 'password';
export const LOGIN_SCENE_EMAIL = 'email';

export type LoginScene = typeof LOGIN_SCENE_PASSWORD | typeof LOGIN_SCENE_EMAIL;

export const EMAIL_CODE_SCENE_REGISTER = 'register';
export const EMAIL_CODE_SCENE_LOGIN = 'login';
export const EMAIL_CODE_SCENE_RESET = 'reset';

export type EmailCodeScene =
  | typeof EMAIL_CODE_SCENE_REGISTER
  | typeof EMAIL_CODE_SCENE_LOGIN
  | typeof EMAIL_CODE_SCENE_RESET;

export const PROJECT_IMAGE_STATUS_PENDING = 'pending';
export const PROJECT_IMAGE_STATUS_UPLOADED = 'uploaded';
export const PROJECT_IMAGE_STATUS_FAILED = 'failed';

export type ProjectImageStatus =
  | typeof PROJECT_IMAGE_STATUS_PENDING
  | typeof PROJECT_IMAGE_STATUS_UPLOADED
  | typeof PROJECT_IMAGE_STATUS_FAILED;

export const PROJECT_EVENT_TYPE_IMAGE_STATUS_CHANGED = 'project_image_status_changed';

export type ProjectEventType = typeof PROJECT_EVENT_TYPE_IMAGE_STATUS_CHANGED;

export const AUTH_STATUS_INITIALIZING = 'initializing';
export const AUTH_STATUS_AUTHENTICATED = 'authenticated';
export const AUTH_STATUS_ANONYMOUS = 'anonymous';

export type AuthStatus =
  | typeof AUTH_STATUS_INITIALIZING
  | typeof AUTH_STATUS_AUTHENTICATED
  | typeof AUTH_STATUS_ANONYMOUS;

export const PROJECT_STAGE_KEY_PROFILE = 'profile';
export const PROJECT_STAGE_KEY_ASSETS = 'assets';
export const PROJECT_STAGE_KEY_DETECTION = 'detection';
export const PROJECT_STAGE_KEY_REVIEW = 'review';
export const PROJECT_STAGE_KEY_REPORT = 'report';

export const PROJECT_STAGE_KEYS = [
  PROJECT_STAGE_KEY_PROFILE,
  PROJECT_STAGE_KEY_ASSETS,
  PROJECT_STAGE_KEY_DETECTION,
  PROJECT_STAGE_KEY_REVIEW,
  PROJECT_STAGE_KEY_REPORT,
] as const;

export type ProjectStageKey = (typeof PROJECT_STAGE_KEYS)[number];

export const PROJECT_PROGRESS_INITIALIZATION_FINISHED = 0;
export const PROJECT_PROGRESS_PROFILE_FINISHED = 1;
export const PROJECT_PROGRESS_ASSETS_FINISHED = 2;
export const PROJECT_PROGRESS_DETECTION_FINISHED = 3;
export const PROJECT_PROGRESS_REVIEW_FINISHED = 4;
export const PROJECT_PROGRESS_REPORT_FINISHED = 5;

export const PROJECT_PROGRESS_VALUES = [
  PROJECT_PROGRESS_INITIALIZATION_FINISHED,
  PROJECT_PROGRESS_PROFILE_FINISHED,
  PROJECT_PROGRESS_ASSETS_FINISHED,
  PROJECT_PROGRESS_DETECTION_FINISHED,
  PROJECT_PROGRESS_REVIEW_FINISHED,
  PROJECT_PROGRESS_REPORT_FINISHED,
] as const;

export type ProjectProgress = (typeof PROJECT_PROGRESS_VALUES)[number];

export const LAST_VISIBLE_PROJECT_PROGRESS = PROJECT_PROGRESS_REVIEW_FINISHED;

export const HTTP_STATUS_UNAUTHORIZED = 401;
export const HTTP_STATUS_INTERNAL_SERVER_ERROR = 500;

export interface ApiEnvelope<T> {
  code: string;
  message: string;
  data: T;
}

export interface User {
  id: number;
  email: string;
  name: string;
}
