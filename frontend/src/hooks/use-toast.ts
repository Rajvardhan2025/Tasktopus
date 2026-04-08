import { useState, useCallback } from 'react';

interface Toast {
  title: string;
  description?: string;
  variant?: 'default' | 'destructive' | 'success';
  id?: string;
}

export function useToast() {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const toast = useCallback(({ title, description, variant = 'default' }: Toast) => {
    const id = Math.random().toString(36).substr(2, 9);
    const newToast: Toast = { title, description, variant, id };

    // Add toast to state for UI rendering
    setToasts((prev) => [...prev, newToast]);

    // Console logging for debugging
    console.log(`[${variant}] ${title}`, description);

    // Auto-remove toast after 4 seconds
    setTimeout(() => {
      setToasts((prev) => prev.filter((t) => t.id !== id));
    }, 4000);
  }, []);

  const removeToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  return { toast, toasts, removeToast };
}
