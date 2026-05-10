import { API_BASE_URL, http } from '@/api/http';
import type { ApiEnvelope } from '@/constants/common';
import type {
  CreateSocketTicketRequest,
  CreateSocketTicketResponse,
  SetupWebSocketRequest,
} from '@/gen/core/api/socket';
import { SocketScope_Value } from '@/gen/core/common';

const SOCKET_SETUP_PATH = '/socket/setup';
const HTTP_PROTOCOL = 'http:';
const HTTPS_PROTOCOL = 'https:';
const WS_PROTOCOL = 'ws:';
const WSS_PROTOCOL = 'wss:';

// 创建 WebSocket 连接票据
// @router /socket/ticket [POST]
async function createSocketTicket(payload: CreateSocketTicketRequest): Promise<CreateSocketTicketResponse> {
  const { data } = await http.post<ApiEnvelope<CreateSocketTicketResponse>>('/socket/ticket', payload);
  return data.data;
}

export function createAssetsSocketTicket(
  payload: Pick<CreateSocketTicketRequest, 'project_id'>,
): Promise<CreateSocketTicketResponse> {
  return createSocketTicket({
    project_id: payload.project_id,
    scope: SocketScope_Value.ProjectAssets,
  });
}

export function createDetectionSocketTicket(
  payload: Pick<CreateSocketTicketRequest, 'project_id'>,
): Promise<CreateSocketTicketResponse> {
  return createSocketTicket({
    project_id: payload.project_id,
    scope: SocketScope_Value.ProjectDetection,
  });
}

export function createReportSocketTicket(
  payload: Pick<CreateSocketTicketRequest, 'project_id'>,
): Promise<CreateSocketTicketResponse> {
  return createSocketTicket({
    project_id: payload.project_id,
    scope: SocketScope_Value.ProjectReport,
  });
}

// 建立 WebSocket 连接
// @router /socket/setup [GET]
function setupProjectWebSocket(payload: SetupWebSocketRequest): string {
  const url = new URL(SOCKET_SETUP_PATH, API_BASE_URL);
  if (url.protocol === HTTPS_PROTOCOL) {
    url.protocol = WSS_PROTOCOL;
  }
  if (url.protocol === HTTP_PROTOCOL) {
    url.protocol = WS_PROTOCOL;
  }
  url.searchParams.set('project_id', payload.project_id);
  url.searchParams.set('scope', String(payload.scope));
  url.searchParams.set('ticket', payload.ticket);
  return url.toString();
}

export function setupAssetsWebSocket(payload: Pick<SetupWebSocketRequest, 'project_id' | 'ticket'>): string {
  return setupProjectWebSocket({
    project_id: payload.project_id,
    scope: SocketScope_Value.ProjectAssets,
    ticket: payload.ticket,
  });
}

export function setupDetectionWebSocket(payload: Pick<SetupWebSocketRequest, 'project_id' | 'ticket'>): string {
  return setupProjectWebSocket({
    project_id: payload.project_id,
    scope: SocketScope_Value.ProjectDetection,
    ticket: payload.ticket,
  });
}

export function setupReportWebSocket(payload: Pick<SetupWebSocketRequest, 'project_id' | 'ticket'>): string {
  return setupProjectWebSocket({
    project_id: payload.project_id,
    scope: SocketScope_Value.ProjectReport,
    ticket: payload.ticket,
  });
}
