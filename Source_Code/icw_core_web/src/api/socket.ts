import { API_BASE_URL, http } from '@/api/http';
import type { ApiEnvelope } from '@/types/common';
import type { CreateSocketTicketResponse } from '@/types/socket';

const SOCKET_ASSETS_PATH = '/socket/setup/assets';
const HTTP_PROTOCOL = 'http:';
const HTTPS_PROTOCOL = 'https:';
const WS_PROTOCOL = 'ws:';
const WSS_PROTOCOL = 'wss:';

// 创建 WebSocket 连接票据
// @router /socket/ticket [POST]
export async function createSocketTicket(projectId: string): Promise<CreateSocketTicketResponse> {
  const { data } = await http.post<ApiEnvelope<CreateSocketTicketResponse>>('/socket/ticket', {
    project_id: projectId,
  });
  return data.data;
}

// 构建图像资产 WebSocket 连接地址
// @router /socket/setup/assets [GET]
export function setupAssetsWebSocket(projectId: string, ticket: string): string {
  const url = new URL(SOCKET_ASSETS_PATH, API_BASE_URL);
  if (url.protocol === HTTPS_PROTOCOL) {
    url.protocol = WSS_PROTOCOL;
  }
  if (url.protocol === HTTP_PROTOCOL) {
    url.protocol = WS_PROTOCOL;
  }
  url.searchParams.set('project_id', projectId);
  url.searchParams.set('ticket', ticket);
  return url.toString();
}
