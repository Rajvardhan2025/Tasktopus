function trimTrailingSlashes(value: string): string {
    return value.replace(/\/+$/, '');
}

function toWebSocketBaseUrl(url: string): string {
    if (url.startsWith('ws://') || url.startsWith('wss://')) {
        return url;
    }

    if (url.startsWith('https://')) {
        return `wss://${url.slice('https://'.length)}`;
    }

    if (url.startsWith('http://')) {
        return `ws://${url.slice('http://'.length)}`;
    }

    return url;
}

const backendUrl = trimTrailingSlashes(import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080');

export const appConfig = {
    backendUrl,
    apiBaseUrl: `${backendUrl}/api`,
    wsBaseUrl: toWebSocketBaseUrl(backendUrl),
};
