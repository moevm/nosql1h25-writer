import axios from 'axios';
import { clearTokens, getAccessToken, getRefreshToken, setTokens } from './auth';

export const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Добавляем перехватчик для обработки ошибок
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // Если ошибка 401 (Unauthorized) и запрос еще не был повторен
    if (error.response && error.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      const refreshToken = getRefreshToken();

      if (refreshToken) {
        try {
          // Пытаемся обновить токен, используя отдельный вызов axios, чтобы избежать зацикливания перехватчика
          const res = await axios.post('/api/auth/refresh', {
            refreshToken: refreshToken,
          });

          setTokens(res.data.accessToken, res.data.refreshToken);

          // Обновляем заголовок авторизации в оригинальном запросе и повторяем его
          originalRequest.headers['Authorization'] = `Bearer ${res.data.accessToken}`;
          return api(originalRequest); // Повторяем оригинальный запрос через наш экземпляр api

        } catch (refreshError) {
          // Если обновление токена не удалось, очищаем токены и перенаправляем на страницу входа
          console.error('Не удалось обновить токен:', refreshError);
          clearTokens();
          // Перенаправление на страницу входа (может потребоваться настройка в роутере)
          window.location.href = '/login';
          return Promise.reject(refreshError); // Отклоняем промис с ошибкой обновления
        }
      }
    }

    // Для всех остальных ошибок или если нет refresh токена, просто отклоняем промис
    return Promise.reject(error);
  }
);

// Добавляем интерцептор для добавления токена в заголовки
api.interceptors.request.use((config) => {
  // Исключаем запросы на авторизацию и регистрацию из добавления заголовка Authorization
  if (config.url && (config.url.endsWith('/auth/login') || config.url.endsWith('/auth/register'))) {
    return config;
  }

  const token = getAccessToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
}); 