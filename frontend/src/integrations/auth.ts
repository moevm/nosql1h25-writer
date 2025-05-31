import { api } from './api';

const ACCESS_TOKEN_KEY = 'accessToken';
const REFRESH_TOKEN_KEY = 'refreshToken';

function notifyAuthChanged() {
  window.dispatchEvent(new Event('auth-changed'));
}

export function setTokens(accessToken: string, refreshToken: string) {
  localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
  localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
  notifyAuthChanged();
}

export function getAccessToken() {
  return localStorage.getItem(ACCESS_TOKEN_KEY);
}

export function getRefreshToken() {
  return localStorage.getItem(REFRESH_TOKEN_KEY);
}

export function clearTokens() {
  localStorage.removeItem(ACCESS_TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
  notifyAuthChanged();
}

export function isAuthenticated() {
  return !!getAccessToken();
}

export async function refreshAccessToken() {
  try {
    const res = await api.post('/auth/refresh', {
      refreshToken: getRefreshToken(),
    });
    setTokens(res.data.accessToken, res.data.refreshToken);
    return res.data.accessToken;
  } catch (refreshError) {
    clearTokens();
    window.location.href = '/login';
    throw refreshError;
  }
}

export function parseJwt(token: string): any {
  try {
    return JSON.parse(atob(token.split('.')[1]));
  } catch (e) {
    return null;
  }
}

export function getUserIdFromToken(): string | null {
  const token = getAccessToken();
  if (!token) return null;
  const payload = parseJwt(token);
  return payload && payload.userId ? payload.userId : null;
}

export function isAdmin(): boolean {
  const token = getAccessToken();
  if (!token) return false;
  const payload = parseJwt(token);
  return payload && payload.systemRole === 'admin';
} 