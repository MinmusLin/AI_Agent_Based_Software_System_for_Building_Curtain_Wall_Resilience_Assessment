import type { Dispatch, SetStateAction } from 'react';
import { useCallback, useEffect } from 'react';

import { createSocketTicket, setupAssetsWebSocket } from '@/api/socket';
import type { ProjectGroup } from '@/types/project/assets';
import { parseProjectImageStatusChangedMessage, replaceImage, WEBSOCKET_RECONNECT_DELAY_MS } from '@/utils/assetsStage';

interface UseProjectAssetsSocketParams {
  projectId: string;
  setGroups: Dispatch<SetStateAction<ProjectGroup[]>>;
}

export function useProjectAssetsSocket({ projectId, setGroups }: UseProjectAssetsSocketParams): void {
  const handleSocketMessage = useCallback(
    (data: unknown): void => {
      const socketMessage = parseProjectImageStatusChangedMessage(data);
      if (socketMessage?.project_id !== projectId) {
        return;
      }
      setGroups((currentGroups) => replaceImage(currentGroups, socketMessage.image));
    },
    [projectId, setGroups],
  );

  useEffect(() => {
    if (projectId === '') {
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
        const ticket = await createSocketTicket(projectId);
        if (closed) {
          return;
        }
        socket = new WebSocket(setupAssetsWebSocket(projectId, ticket.ticket));
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
  }, [handleSocketMessage, projectId]);
}
