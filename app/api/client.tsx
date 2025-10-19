import { API_CONFIG } from '@/api/apiClient';
import { secureStorage } from '@/api/storage/secure-storage';
import axios from 'axios';


export const httpClient = axios.create({
  baseURL: API_CONFIG.BASE_URL,
  timeout: API_CONFIG.TIMEOUT,
  headers: {
    'Content-Type': 'application/json',
  },
  transformRequest: [
    (data, headers) => {
      if (data instanceof FormData) {
        delete headers['Content-Type'];
        return data;
      }
      if (headers['Content-Type'] === 'application/json') {
        return JSON.stringify(data);
      }
      return data;
    }
  ],
});

export class APIError extends Error {
  constructor(
    public message: string,
    public status?: number,
    public data?: any
  ) {
    super(message);
    this.name = 'APIError';
  }
}

export const executeAPI = async<T = any> (
  endpoint: { url: string; method: string },
  data?: any,
  params?: any
): Promise<{ data: T }> => {
  const { url, method } = endpoint;
  try {
    let response;

    switch (method.toUpperCase()) {
      case 'GET':
        response = await httpClient.get<T>(url, { params });
        break;
      case 'POST':
        response = await httpClient.post<T>(url, data, { params });
        break;
      case 'PUT':
        response = await httpClient.put<T>(url, data, { params });
        break;
      case 'PATCH':
        response = await httpClient.patch<T>(url, data, { params });
        break;
      case 'DELETE':
        response = await httpClient.delete<T>(url, { params }); 
        break;
      default:
        throw new APIError(`Unsupported HTTP method: ${method}`);
    }
    return { data: response.data };
  } catch (error) {
    if (error instanceof APIError) {
      throw error;
    }

    throw new APIError(
      'Network Error',
      0,
      error
    );
  }
}


httpClient.interceptors.request.use(async (config) => {
  const token = await secureStorage.getToken('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  
  if (__DEV__) {
    console.log(`[API Request] ${config.method?.toUpperCase()} ${config.url}`);
  }
  
  return config;
});


httpClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    
    const isAuthEndpoint = originalRequest.url?.includes('/auth/login') || 
                          originalRequest.url?.includes('/auth/register') ||
                          originalRequest.url?.includes('/auth/refresh-token');
    
    if (error.response?.status === 401 && !originalRequest._retry && !isAuthEndpoint) {
      originalRequest._retry = true;
      
      try {
        const refreshToken = await secureStorage.getToken('refresh_token');
        if (refreshToken) {
          const response = await axios.post(`${API_CONFIG.BASE_URL}auth/refresh-token`, {
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
