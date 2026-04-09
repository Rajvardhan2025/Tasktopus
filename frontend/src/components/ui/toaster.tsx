import { X } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';
import { Button } from './button';

const VARIANT_STYLES: Record<string, string> = {
    default: 'border-slate-200 bg-white text-slate-900',
    destructive: 'border-red-300 bg-red-50 text-red-900',
    success: 'border-emerald-300 bg-emerald-50 text-emerald-900',
};

export function Toaster() {
    const { toasts, removeToast } = useToast();

    return (
        <div className="pointer-events-none fixed right-4 top-4 z-[100] flex w-full max-w-sm flex-col gap-2">
            {toasts.map((toast) => (
                <div
                    key={toast.id}
                    className={`pointer-events-auto rounded-lg border px-4 py-3 shadow-md ${VARIANT_STYLES[toast.variant || 'default']}`}
                >
                    <div className="flex items-start justify-between gap-2">
                        <div>
                            <div className="text-sm font-semibold">{toast.title}</div>
                            {toast.description && <div className="mt-1 text-xs opacity-90">{toast.description}</div>}
                        </div>
                        <Button
                            variant="ghost"
                            size="sm"
                            className="h-6 w-6 p-0"
                            onClick={() => toast.id && removeToast(toast.id)}
                        >
                            <X className="h-3.5 w-3.5" />
                        </Button>
                    </div>
                </div>
            ))}
        </div>
    );
}
