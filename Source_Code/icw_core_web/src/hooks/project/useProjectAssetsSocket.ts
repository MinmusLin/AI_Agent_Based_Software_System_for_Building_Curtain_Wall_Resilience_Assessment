import type { Dispatch, SetStateAction } from 'react';
import { useCallback, useEffect, useRef } from 'react';

import { createAssetsSocketTicket, setupAssetsWebSocket } from '@/api/socket';
import type { ProjectGroup } from '@/gen/core/api/common';
import { parseProjectImageStatusChangedMessage, replaceImage, WEBSOCKET_RECONNECT_DELAY_MS } from '@/utils/assetsStage';

interface UseProjectAssetsSocketParams {
  onConnected?: () => void;
  projectId: string;
  setGroups: Dispatch<SetStateAction<ProjectGroup[]>>;
}

export function useProjectAssetsSocket({ onConnected, projectId, setGroups }: UseProjectAssetsSocketParams): void {
  const onConnectedRef = useRef(onConnected);

  useEffect(() => {
    onConnectedRef.current = onConnected;
  }, [onConnected]);

  const handleSocketMessage = useCallback(
    (data: unknown): void => {
      const socketMessage = parseProjectImageStatusChangedMessage(data);
      if (socketMessage?.project_id !== projectId || !socketMessage.image) {
        return;
      }
      const { image } = socketMessage;
      setGroups((currentGroups) => replaceImage(currentGroups, image));
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
        const ticket = await createAssetsSocketTicket({
          project_id: projectId,
        });
        if (closed) {
          return;
        }
        socket = new WebSocket(
          setupAssetsWebSocket({
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
  }, [handleSocketMessage, projectId]);
}
