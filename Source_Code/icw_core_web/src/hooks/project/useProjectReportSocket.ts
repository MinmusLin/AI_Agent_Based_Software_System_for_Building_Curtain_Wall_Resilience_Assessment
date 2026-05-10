import { useEffect } from 'react';

import { createReportSocketTicket, setupReportWebSocket } from '@/api/socket';
import { PROJECT_EVENT_TYPE_REPORT_STATUS_CHANGED } from '@/types/common';
import type { ProjectReportStatusChangedMessage } from '@/types/project/report';
import { WEBSOCKET_RECONNECT_DELAY_MS } from '@/utils/assetsStage';

interface UseProjectReportSocketParams {
  enabled: boolean;
  onReportChanged: () => void;
  projectId: string;
}

function parseProjectReportStatusChangedMessage(data: unknown): ProjectReportStatusChangedMessage | null {
  if (typeof data !== 'string') {
    return null;
  }
  try {
    const parsed = JSON.parse(data) as ProjectReportStatusChangedMessage;
    if (parsed.type !== PROJECT_EVENT_TYPE_REPORT_STATUS_CHANGED) {
      return null;
    }
    return parsed;
  } catch {
    return null;
  }
}

export function useProjectReportSocket({ enabled, onReportChanged, projectId }: UseProjectReportSocketParams): void {
  useEffect(() => {
    if (!enabled || projectId === '') {
      return undefined;
    }

    let closed = false;
    let socket: WebSocket | null = null;
    let reconnectTimer: number | null = null;

    function clearReconnectTimer(): void {
      if (reconnectTimer === null) {
        return;
      }
      window.clearTimeout(reconnectTimer);
      reconnectTimer = null;
    }

    function scheduleReconnect(): void {
      if (closed || reconnectTimer !== null) {
        return;
      }
      reconnectTimer = window.setTimeout(() => {
        reconnectTimer = null;
        void connect();
      }, WEBSOCKET_RECONNECT_DELAY_MS);
    }

    async function connect(): Promise<void> {
      try {
        const ticket = await createReportSocketTicket(projectId);
        if (closed) {
          return;
        }
        socket = new WebSocket(setupReportWebSocket(projectId, ticket.ticket));
        socket.onmessage = (event: MessageEvent<unknown>): void => {
          const message = parseProjectReportStatusChangedMessage(event.data);
          if (message?.project_id !== projectId) {
            return;
          }
          onReportChanged();
        };
        socket.onclose = (): void => {
          if (!closed) {
            scheduleReconnect();
          }
        };
        socket.onerror = (): void => {
          socket?.close();
        };
      } catch {
        scheduleReconnect();
      }
    }

    void connect();

    return () => {
      closed = true;
      clearReconnectTimer();
      socket?.close();
    };
  }, [enabled, onReportChanged, projectId]);
}
