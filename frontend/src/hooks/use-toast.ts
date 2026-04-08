import { useState, useCallback } from 'react';

interface Toast {
  title: string;
  description?: string;
  variant?: 'default' | 'destructive';
}

export function useToast() {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const toast = useCallback((toast: Toast) => {
    // Simple console implementation - in production, use a proper toast library
    console.log(`[${toast.variant || 'default'}] ${toast.title}`, toast.description);
    
    // You can integrate with a proper toast library here
    // For now, we'll use browser alert for errors
    if (toast.variant === 'destructive') {
      alert(`${toast.title}: ${toast.description}`);
    }
  }, []);

  return { toast, toasts };
}
