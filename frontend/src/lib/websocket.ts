import { useEffect, useRef, useState } from 'react';
import { appConfig } from '@/lib/config';

export interface WSEvent {
  type: 'issue_created' | 'issue_updated' | 'issue_moved' | 'comment_added' | 'sprint_updated' | 'presence';
  project_id: string;
  data: any;
  timestamp: number;
}

export function useWebSocket(projectId: string | null, userId: string = 'demo-user') {
  const [isConnected, setIsConnected] = useState(false);
  const [lastEvent, setLastEvent] = useState<WSEvent | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimerRef = useRef<number | null>(null);
  const lastSeenTimestampRef = useRef<number>(0);

  useEffect(() => {
    if (!projectId) {
      return;
    }

    let isUnmounted = false;

    const connect = () => {
      const params = new URLSearchParams({ userId });
      if (lastSeenTimestampRef.current > 0) {
        params.set('since', String(lastSeenTimestampRef.current));
      }

      const ws = new WebSocket(`${appConfig.wsBaseUrl}/api/ws/${projectId}?${params.toString()}`);

      ws.onopen = () => {
        setIsConnected(true);
      };

      ws.onmessage = (event) => {
        const data = JSON.parse(event.data) as WSEvent;
        lastSeenTimestampRef.current = Math.max(lastSeenTimestampRef.current, data.timestamp || 0);
        setLastEvent(data);
      };

      ws.onerror = () => {
        setIsConnected(false);
      };

      ws.onclose = () => {
        setIsConnected(false);
        if (!isUnmounted) {
          reconnectTimerRef.current = window.setTimeout(connect, 2000);
        }
      };

      wsRef.current = ws;
    };

    connect();

    return () => {
      isUnmounted = true;
      if (reconnectTimerRef.current) {
        window.clearTimeout(reconnectTimerRef.current);
      }
      wsRef.current?.close();
    };
  }, [projectId, userId]);

  return { isConnected, lastEvent };
}
