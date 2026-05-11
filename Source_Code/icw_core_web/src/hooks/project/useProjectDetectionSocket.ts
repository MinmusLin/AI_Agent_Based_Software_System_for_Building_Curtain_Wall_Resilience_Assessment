import type { Dispatch, SetStateAction } from 'react';
import { useCallback, useEffect, useRef } from 'react';

import { createDetectionSocketTicket, setupDetectionWebSocket } from '@/api/socket';
import { WEBSOCKET_RECONNECT_DELAY_MS } from '@/utils/assetsStage';
import type { ProjectDetectionStatusMap } from '@/utils/detectionStage';
import { mergeDetectionTaskMessage, parseProjectDetectionTaskStatusChangedMessage } from '@/utils/detectionStage';

interface UseProjectDetectionSocketParams {
  enabled: boolean;
  onConnected?: () => void;
  projectId: string;
  setTasks: Dispatch<SetStateAction<ProjectDetectionStatusMap>>;
}

export function useProjectDetectionSocket({
  enabled,
  onConnected,
  projectId,
  setTasks,
}: UseProjectDetectionSocketParams): void {
  const onConnectedRef = useRef(onConnected);

  useEffect(() => {
    onConnectedRef.current = onConnected;
  }, [onConnected]);

  const handleSocketMessage = useCallback(
    (data: unknown): void => {
      const socketMessage = parseProjectDetectionTaskStatusChangedMessage(data);
      if (socketMessage?.project_id !== projectId) {
        return;
      }
      setTasks((currentTasks) => mergeDetectionTaskMessage(currentTasks, socketMessage));
    },
    [projectId, setTasks],
  );

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
        const ticket = await createDetectionSocketTicket({
          project_id: projectId,
        });
        if (closed) {
          return;
        }
        socket = new WebSocket(
          setupDetectionWebSocket({
            project_id: projectId,
            ticket: ticket.ticket,
          }),
        );
        socket.onopen = (): void => {
          onConnectedRef.current?.();
        };
        socket.onmessage = (event: MessageEvent<unknown>): void => {
          handleSocketMessage(event.data);
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
  }, [enabled, handleSocketMessage, projectId]);
}
