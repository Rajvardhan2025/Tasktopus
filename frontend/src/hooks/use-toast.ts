import { useCallback, useEffect, useState } from 'react';

interface Toast {
  title: string;
  description?: string;
  variant?: 'default' | 'destructive' | 'success';
  id?: string;
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

  const toast = useCallback(({ title, description, variant = 'default' }: Toast) => {
    const id = Math.random().toString(36).substr(2, 9);
    const newToast: Toast = { title, description, variant, id };

    addToast(newToast);

    // Console logging for debugging
    console.log(`[${variant}] ${title}`, description);

    // Auto-remove toast after 4 seconds
    setTimeout(() => {
      removeToastById(id);
    }, 4000);
  }, []);

  const removeToast = useCallback((id: string) => {
    removeToastById(id);
  }, []);

  return { toast, toasts, removeToast };
}
