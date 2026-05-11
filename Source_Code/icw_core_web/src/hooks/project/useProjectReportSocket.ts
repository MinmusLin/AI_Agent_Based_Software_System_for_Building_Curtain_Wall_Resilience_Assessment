import { useEffect, useRef } from 'react';

import { createReportSocketTicket, setupReportWebSocket } from '@/api/socket';
import type { ProjectReportStatusChangedMessage } from '@/gen/core/api/messages';
import { ProjectEventType_Value } from '@/gen/core/common';
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
    if (parsed.type !== ProjectEventType_Value.ReportStatusChanged) {
      return null;
    }
    return parsed;
  } catch {
    return null;
  }
}

export function useProjectReportSocket({ enabled, onReportChanged, projectId }: UseProjectReportSocketParams): void {
  const onReportChangedRef = useRef(onReportChanged);

  useEffect(() => {
    onReportChangedRef.current = onReportChanged;
  }, [onReportChanged]);

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
        const ticket = await createReportSocketTicket({
          project_id: projectId,
        });
        if (closed) {
          return;
        }
        socket = new WebSocket(
          setupReportWebSocket({
            project_id: projectId,
            ticket: ticket.ticket,
          }),
        );
        socket.onopen = (): void => {
          onReportChangedRef.current();
        };
        socket.onmessage = (event: MessageEvent<unknown>): void => {
          const message = parseProjectReportStatusChangedMessage(event.data);
          if (message?.project_id !== projectId) {
            return;
          }
          onReportChangedRef.current();
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
  }, [enabled, projectId]);
}
