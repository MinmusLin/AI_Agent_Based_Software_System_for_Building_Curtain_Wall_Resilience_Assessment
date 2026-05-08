import { API_BASE_URL, http } from '@/api/http';
import {
  type ApiEnvelope,
  SOCKET_SCOPE_PROJECT_ASSETS,
  SOCKET_SCOPE_PROJECT_DETECTION,
  SOCKET_SCOPE_PROJECT_REPORT,
  type SocketScope,
} from '@/types/common';
import type { CreateSocketTicketResponse } from '@/types/socket';

const SOCKET_SETUP_PATH = '/socket/setup';
const HTTP_PROTOCOL = 'http:';
const HTTPS_PROTOCOL = 'https:';
const WS_PROTOCOL = 'ws:';
const WSS_PROTOCOL = 'wss:';

// 创建 WebSocket 连接票据
// @router /socket/ticket [POST]
async function createSocketTicket(projectId: string, scope: SocketScope): Promise<CreateSocketTicketResponse> {
  const { data } = await http.post<ApiEnvelope<CreateSocketTicketResponse>>('/socket/ticket', {
    project_id: projectId,
    scope,
  });
  return data.data;
}

export function createAssetsSocketTicket(projectId: string): Promise<CreateSocketTicketResponse> {
  return createSocketTicket(projectId, SOCKET_SCOPE_PROJECT_ASSETS);
}

export function createDetectionSocketTicket(projectId: string): Promise<CreateSocketTicketResponse> {
  return createSocketTicket(projectId, SOCKET_SCOPE_PROJECT_DETECTION);
}

export function createReportSocketTicket(projectId: string): Promise<CreateSocketTicketResponse> {
  return createSocketTicket(projectId, SOCKET_SCOPE_PROJECT_REPORT);
}

// 建立 WebSocket 连接
// @router /socket/setup [GET]
function setupProjectWebSocket(projectId: string, scope: SocketScope, ticket: string): string {
  const url = new URL(SOCKET_SETUP_PATH, API_BASE_URL);
  if (url.protocol === HTTPS_PROTOCOL) {
    url.protocol = WSS_PROTOCOL;
  }
  if (url.protocol === HTTP_PROTOCOL) {
    url.protocol = WS_PROTOCOL;
  }
  url.searchParams.set('project_id', projectId);
  url.searchParams.set('scope', scope);
  url.searchParams.set('ticket', ticket);
  return url.toString();
}

export function setupAssetsWebSocket(projectId: string, ticket: string): string {
  return setupProjectWebSocket(projectId, SOCKET_SCOPE_PROJECT_ASSETS, ticket);
}

export function setupDetectionWebSocket(projectId: string, ticket: string): string {
  return setupProjectWebSocket(projectId, SOCKET_SCOPE_PROJECT_DETECTION, ticket);
}

export function setupReportWebSocket(projectId: string, ticket: string): string {
  return setupProjectWebSocket(projectId, SOCKET_SCOPE_PROJECT_REPORT, ticket);
}
