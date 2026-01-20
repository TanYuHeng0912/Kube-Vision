// Configuration for API and WebSocket URLs
const getApiUrl = (): string => {
  // Check environment variable first
  if (import.meta.env.VITE_API_URL) {
    return import.meta.env.VITE_API_URL;
  }
  
  // In production, use current host
  if (import.meta.env.PROD) {
    const protocol = window.location.protocol === 'https:' ? 'https:' : 'http:';
    return `${protocol}//${window.location.host}`;
  }
  
  // Development default
  return 'http://localhost:8080';
};

const getWsUrl = (): string => {
  const apiUrl = getApiUrl();
  return apiUrl.replace(/^http/, 'ws');
};

export const API_URL = getApiUrl();
export const WS_URL = getWsUrl();

