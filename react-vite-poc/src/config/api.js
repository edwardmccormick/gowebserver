const DEFAULT_API_BASE_URL = 'http://localhost:8081';

function normalizeBaseUrl(value) {
  return value.replace(/\/+$/, '');
}

function deriveWebSocketBaseUrl(apiBaseUrl) {
  if (apiBaseUrl.startsWith('https://')) {
    return `wss://${apiBaseUrl.slice('https://'.length)}`;
  }
  if (apiBaseUrl.startsWith('http://')) {
    return `ws://${apiBaseUrl.slice('http://'.length)}`;
  }
  return apiBaseUrl;
}

export const API_BASE_URL = normalizeBaseUrl(
  import.meta.env.VITE_API_BASE_URL?.trim() || DEFAULT_API_BASE_URL,
);

export const WS_BASE_URL = normalizeBaseUrl(deriveWebSocketBaseUrl(API_BASE_URL));

export function apiUrl(path) {
  return `${API_BASE_URL}${path.startsWith('/') ? path : `/${path}`}`;
}

export function wsUrl(path) {
  return `${WS_BASE_URL}${path.startsWith('/') ? path : `/${path}`}`;
}
