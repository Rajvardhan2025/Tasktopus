import { useEffect, useRef, useState } from 'react';

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

  useEffect(() => {
    if (!projectId) return;

    const ws = new WebSocket(`ws://localhost:8080/api/ws/${projectId}?userId=${userId}`);
    
    ws.onopen = () => {
      console.log('WebSocket connected');
      setIsConnected(true);
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data) as WSEvent;
      console.log('WebSocket event:', data);
      setLastEvent(data);
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    ws.onclose = () => {
      console.log('WebSocket disconnected');
      setIsConnected(false);
    };

    wsRef.current = ws;

    return () => {
      ws.close();
    };
  }, [projectId, userId]);

  return { isConnected, lastEvent };
}
