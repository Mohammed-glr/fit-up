import axios from 'axios';
import { secureStorage } from '@/api/storage/secure-storage';
import { API_CONFIG } from '@/api/apiClient';


export const httpClient = axios.create({
  baseURL: API_CONFIG.BASE_URL,
  timeout: API_CONFIG.TIMEOUT,
  headers: {
    'Content-Type': 'application/json',
  },
});


export const executeAPI = async (
  endpoint: { url: string; method: string },
  data?: any,
  params?: any
) => {
  const { url, method } = endpoint;

  switch (method.toUpperCase()) {
    case 'GET':
      return httpClient.get(url, { params });
    case 'POST':
      return httpClient.post(url, data);
    case 'PUT':
      return httpClient.put(url, data);
    case 'PATCH':
      return httpClient.patch(url, data);
    case 'DELETE':
      return httpClient.delete(url);
    default:
      throw new Error(`Unsupported HTTP method: ${method}`);
  }
}


httpClient.interceptors.request.use(async (config) => {
  const token = await secureStorage.getToken('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});


httpClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      
      try {
        const refreshToken = await secureStorage.getToken('refresh_token');
        if (refreshToken) {
          const response = await axios.post(`${API_CONFIG.BASE_URL}/auth/refresh-token`, {
            refresh_token: refreshToken
          });
          
          const { access_token, refresh_token: newRefreshToken } = response.data;
          await secureStorage.setToken('access_token', access_token);
          if (newRefreshToken) {
            await secureStorage.setToken('refresh_token', newRefreshToken);
          }
          
          originalRequest.headers.Authorization = `Bearer ${access_token}`;
          return httpClient(originalRequest);
        }
      } catch (refreshError) {
        await secureStorage.clearTokens();
      }
    }
    
    return Promise.reject(error);
  }
);
