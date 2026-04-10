import { useCallback, useEffect, useState } from 'react';

interface Toast {
  title: string;
  description?: string;
  variant?: 'default' | 'destructive' | 'success';
  id?: string;
  duration?: number;
}

type ToastListener = (toasts: Toast[]) => void;

let toastState: Toast[] = [];
const listeners: ToastListener[] = [];

function notifyListeners() {
  listeners.forEach((listener) => listener(toastState));
}

function addToast(toast: Toast) {
  toastState = [...toastState, toast];
  notifyListeners();
}

function removeToastById(id: string) {
  toastState = toastState.filter((toast) => toast.id !== id);
  notifyListeners();
}

// Standalone toast function that can be used outside of React components
export function toast({ title, description, variant = 'default', id: customId, duration = 4000 }: Toast) {
  const id = customId || Math.random().toString(36).substr(2, 9);
  const newToast: Toast = { title, description, variant, id };

  // If duration is 0, dismiss immediately (used for dismissing existing toasts)
  if (duration === 0 && customId) {
    removeToastById(customId);
    return { id };
  }

  addToast(newToast);

  // Console logging for debugging
  if (title) {
    console.log(`[${variant}] ${title}`, description);
  }

  // Auto-remove toast after specified duration
  if (duration > 0) {
    setTimeout(() => {
      removeToastById(id);
    }, duration);
  }

  return { id };
}

export function useToast() {
  const [toasts, setToasts] = useState<Toast[]>(toastState);

  useEffect(() => {
    listeners.push(setToasts);
    return () => {
      const index = listeners.indexOf(setToasts);
      if (index >= 0) {
        listeners.splice(index, 1);
      }
    };
  }, []);

  const removeToast = useCallback((id: string) => {
    removeToastById(id);
  }, []);

  return { toast, toasts, removeToast };
}
