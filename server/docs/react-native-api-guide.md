# Fit-Up React Native API Integration Guide

## Overview
This document provides comprehensive guidance for integrating React Native applications with the Fit-Up Backend authentication system. The backend uses a microservices architecture with JWT-based authentication and refresh tokens for secure user management.

## Table of Contents
- [Architecture Overview](#architecture-overview)
- [Base Configuration](#base-configuration)
- [Authentication Flow](#authentication-flow)
- [API Endpoints](#api-endpoints)
- [Error Handling](#error-handling)
- [React Native Implementation Examples](#react-native-implementation-examples)
- [Security Considerations](#security-considerations)
- [TypeScript Types](#typescript-types)

## Architecture Overview

```
React Native App → API Gateway → Auth Service
      ↓               ↓              ↓
   Client Side → Port 8080 → Port 8081
```

### Key Components
- **API Gateway**: Routes requests and handles CORS (Port 8080)
- **Auth Service**: Manages authentication, user registration, and JWT tokens (Port 8081)
- **Message Service**: Handles user messaging and notifications (Port 8082)
- **Schema Service**: Manages data schemas and validation (Port 8083)

## Base Configuration

### Environment Setup
```typescript
// config/api.ts
const API_CONFIG = {
  BASE_URL: __DEV__ 
    ? 'http://localhost:8080'  // Development - API Gateway
    : 'https://api.fitup.com',  // Production
  TIMEOUT: 10000,
};

export default API_CONFIG;
```

### HTTP Client Setup
```typescript
// services/httpClient.ts
import axios from 'axios';
import AsyncStorage from '@react-native-async-storage/async-storage';
import API_CONFIG from '../config/api';

const httpClient = axios.create({
  baseURL: API_CONFIG.BASE_URL,
  timeout: API_CONFIG.TIMEOUT,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  },
});

// Request interceptor to add authentication token
httpClient.interceptors.request.use(
  async (config) => {
    const token = await AsyncStorage.getItem('access_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    console.log('Request:', config.method?.toUpperCase(), config.url);
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor for token refresh
httpClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      
      try {
        const refreshToken = await AsyncStorage.getItem('refresh_token');
        if (refreshToken) {
          const response = await axios.post(`${API_CONFIG.BASE_URL}/auth/refresh-token`, {
            refresh_token: refreshToken
          });
          
          const { access_token, refresh_token: newRefreshToken } = response.data;
          await AsyncStorage.setItem('access_token', access_token);
          if (newRefreshToken) {
            await AsyncStorage.setItem('refresh_token', newRefreshToken);
          }
          
          originalRequest.headers.Authorization = `Bearer ${access_token}`;
          return httpClient(originalRequest);
        }
      } catch (refreshError) {
        // Refresh failed, redirect to login
        await AsyncStorage.multiRemove(['access_token', 'refresh_token']);
        // You might want to trigger a logout action here
      }
    }
    
    console.error('API Error:', error.response?.data || error.message);
    return Promise.reject(error);
  }
);

export default httpClient;
```

## Authentication Flow

### 1. User Registration

**Endpoint**: `POST /auth/register`

**Request Body**:
```typescript
interface RegisterRequest {
  username: string;    // min: 3, max: 50 characters
  email: string;       // valid email format
  password: string;    // min: 8 characters
  name?: string;       // optional, max: 100 characters
}
```

**Response**:
```typescript
interface AuthResponse {
  access_token: string;
  refresh_token?: string;
  token_type: string;        // "Bearer"
  expires_in: number;        // seconds
  user: User;
}
```

**Implementation**:
```typescript
// services/authService.ts
export const register = async (userData: RegisterRequest): Promise<AuthResponse> => {
  const response = await httpClient.post('/auth/register', userData);
  return response.data;
};
```

### 2. User Login

**Endpoint**: `POST /auth/login`

**Request Body**:
```typescript
interface LoginRequest {
  identifier: string;  // email or username
  password: string;    // min: 8 characters
}
```

**Response**: Same as registration (`AuthResponse`)

**Implementation**:
```typescript
export const login = async (credentials: LoginRequest): Promise<AuthResponse> => {
  const response = await httpClient.post('/auth/login', credentials);
  return response.data;
};
```

### 3. Token Refresh

**Endpoint**: `POST /auth/refresh-token`

**Request Body**:
```typescript
interface RefreshTokenRequest {
  refresh_token: string;
}
```

**Response**:
```typescript
interface RefreshTokenResponse {
  access_token: string;
  refresh_token?: string;
  token_type: string;
  expires_in: number;
}
```

### 4. Logout

**Endpoint**: `POST /auth/logout`

**Request Body**: Empty `{}`

**Response**: `204 No Content`

**Implementation**:
```typescript
export const logout = async (): Promise<void> => {
  await httpClient.post('/auth/logout', {});
};
```

## API Endpoints

### Authentication Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/auth/register` | User registration | No |
| POST | `/auth/login` | User login | No |
| POST | `/auth/logout` | User logout | No |
| POST | `/auth/refresh-token` | Refresh access token | No |
| POST | `/auth/validate-token` | Validate current token | No |
| POST | `/auth/forgot-password` | Request password reset | No |
| POST | `/auth/reset-password` | Reset password with token | No |
| POST | `/auth/change-password` | Change password | Yes |
| GET | `/auth/{username}` | Get user profile | No |

### OAuth Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/auth/oauth/{provider}` | Initiate OAuth flow | No |
| GET | `/auth/oauth/callback/{provider}` | OAuth callback | No |
| POST | `/auth/link/{provider}` | Link OAuth account | Yes |
| DELETE | `/auth/unlink/{provider}` | Unlink OAuth account | Yes |
| GET | `/auth/linked-accounts` | Get linked accounts | Yes |

### Supported OAuth Providers
- `google`
- `github`  
- `facebook`

### Future Service Endpoints
```typescript
// Message Service (Coming Soon)
// GET /messages - Get user messages
// POST /messages - Send message
// GET /messages/{id} - Get specific message
// DELETE /messages/{id} - Delete message

// Schema Service (Coming Soon)  
// GET /schemas - Get available schemas
// POST /schemas/validate - Validate data against schema
```

## Error Handling

### Error Response Format
```typescript
interface ApiError {
  code: string;
  message: string;
}
```

### Common Error Codes

| Code | Message | Description |
|------|---------|-------------|
| `INVALID_CREDENTIALS` | Invalid email or password | Login failed |
| `USER_NOT_FOUND` | User not found | User doesn't exist |
| `USER_ALREADY_EXISTS` | User already exists | Registration failed |
| `USERNAME_ALREADY_EXISTS` | Username already exists | Username taken |
| `INVALID_TOKEN` | Invalid or expired token | Token validation failed |
| `UNAUTHORIZED` | Unauthorized access | Authentication required |
| `TOO_MANY_ATTEMPTS` | Too many failed attempts | Rate limit exceeded |
| `INVALID_INPUT` | Invalid input provided | Validation error |
| `REFRESH_TOKEN_NOT_FOUND` | Refresh token not found | Token refresh failed |
| `REFRESH_TOKEN_EXPIRED` | Refresh token has expired | Re-login required |

### Error Handling Implementation
```typescript
// utils/errorHandler.ts
export const handleApiError = (error: any): string => {
  if (error.response?.data?.code) {
    const errorCode = error.response.data.code;
    
    switch (errorCode) {
      case 'INVALID_CREDENTIALS':
        return 'Invalid email or password';
      case 'USER_ALREADY_EXISTS':
        return 'An account with this email already exists';
      case 'USERNAME_ALREADY_EXISTS':
        return 'This username is already taken';
      case 'TOO_MANY_ATTEMPTS':
        return 'Too many attempts. Please try again later';
      case 'REFRESH_TOKEN_EXPIRED':
        return 'Session expired. Please login again';
      case 'UNAUTHORIZED':
        return 'Please login to continue';
      default:
        return error.response.data.message || 'An error occurred';
    }
  }
  
  return 'Network error. Please check your connection';
};
```

## React Native Implementation Examples

### Authentication Hook
```typescript
// hooks/useAuth.ts
import { useState, useEffect, useContext, createContext } from 'react';
import AsyncStorage from '@react-native-async-storage/async-storage';
import * as authService from '../services/authService';

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (credentials: LoginRequest) => Promise<void>;
  register: (userData: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  refreshToken: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    checkAuthState();
  }, []);

  const checkAuthState = async () => {
    try {
      const token = await AsyncStorage.getItem('access_token');
      if (token) {
        // Validate token with server
        const response = await authService.validateToken();
        setUser(response.user);
      }
    } catch (error) {
      // Token invalid, clear storage
      await AsyncStorage.multiRemove(['access_token', 'refresh_token']);
    } finally {
      setIsLoading(false);
    }
  };

  const login = async (credentials: LoginRequest) => {
    try {
      setIsLoading(true);
      const response = await authService.login(credentials);
      setUser(response.user);
      
      // Store tokens securely
      await AsyncStorage.setItem('access_token', response.access_token);
      if (response.refresh_token) {
        await AsyncStorage.setItem('refresh_token', response.refresh_token);
      }
    } catch (error) {
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const register = async (userData: RegisterRequest) => {
    try {
      setIsLoading(true);
      const response = await authService.register(userData);
      setUser(response.user);
      
      await AsyncStorage.setItem('access_token', response.access_token);
      if (response.refresh_token) {
        await AsyncStorage.setItem('refresh_token', response.refresh_token);
      }
    } catch (error) {
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = async () => {
    try {
      await authService.logout();
    } catch (error) {
      // Even if logout fails on server, clear local data
      console.warn('Logout API call failed:', error);
    } finally {
      setUser(null);
      await AsyncStorage.multiRemove(['access_token', 'refresh_token']);
    }
  };

  const refreshToken = async () => {
    try {
      const refreshToken = await AsyncStorage.getItem('refresh_token');
      if (!refreshToken) throw new Error('No refresh token');
      
      const response = await authService.refreshToken({ refresh_token: refreshToken });
      await AsyncStorage.setItem('access_token', response.access_token);
      if (response.refresh_token) {
        await AsyncStorage.setItem('refresh_token', response.refresh_token);
      }
    } catch (error) {
      await logout();
      throw error;
    }
  };

  const value = {
    user,
    isLoading,
    isAuthenticated: !!user,
    login,
    register,
    logout,
    refreshToken,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
```

### Login Screen Example
```typescript
// screens/LoginScreen.tsx
import React, { useState } from 'react';
import { 
  View, 
  TextInput, 
  TouchableOpacity, 
  Text, 
  Alert, 
  StyleSheet,
  KeyboardAvoidingView,
  Platform 
} from 'react-native';
import { useAuth } from '../hooks/useAuth';
import { handleApiError } from '../utils/errorHandler';

const LoginScreen: React.FC = () => {
  const [identifier, setIdentifier] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const { login } = useAuth();

  const handleLogin = async () => {
    if (!identifier || !password) {
      Alert.alert('Error', 'Please fill in all fields');
      return;
    }

    try {
      setIsLoading(true);
      await login({ identifier, password });
      // Navigation will be handled by AuthProvider state change
    } catch (error) {
      Alert.alert('Login Failed', handleApiError(error));
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <KeyboardAvoidingView 
      style={styles.container}
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
    >
      <View style={styles.form}>
        <Text style={styles.title}>Fit-Up Login</Text>
        
        <TextInput
          style={styles.input}
          placeholder="Email or Username"
          value={identifier}
          onChangeText={setIdentifier}
          autoCapitalize="none"
          autoCorrect={false}
          keyboardType="email-address"
        />
        
        <TextInput
          style={styles.input}
          placeholder="Password"
          value={password}
          onChangeText={setPassword}
          secureTextEntry
        />
        
        <TouchableOpacity
          style={[styles.button, isLoading && styles.buttonDisabled]}
          onPress={handleLogin}
          disabled={isLoading}
        >
          <Text style={styles.buttonText}>
            {isLoading ? 'Logging in...' : 'Login'}
          </Text>
        </TouchableOpacity>
      </View>
    </KeyboardAvoidingView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    backgroundColor: '#f5f5f5',
  },
  form: {
    padding: 20,
    margin: 20,
    backgroundColor: 'white',
    borderRadius: 10,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    textAlign: 'center',
    marginBottom: 30,
    color: '#333',
  },
  input: {
    borderWidth: 1,
    borderColor: '#ddd',
    padding: 15,
    marginBottom: 15,
    borderRadius: 5,
    fontSize: 16,
  },
  button: {
    backgroundColor: '#007AFF',
    padding: 15,
    borderRadius: 5,
    alignItems: 'center',
    marginTop: 10,
  },
  buttonDisabled: {
    backgroundColor: '#ccc',
  },
  buttonText: {
    color: 'white',
    fontSize: 16,
    fontWeight: 'bold',
  },
});

export default LoginScreen;
```

### OAuth Integration Example
```typescript
// services/oauthService.ts
import { Linking } from 'react-native';
import httpClient from './httpClient';

interface OAuthAuthRequest {
  provider: 'google' | 'github' | 'facebook';
  redirect_url?: string;
}

export const initiateOAuth = async (provider: string): Promise<string> => {
  const response = await httpClient.post(`/auth/oauth/${provider}`, {
    provider,
    redirect_url: 'fitup://oauth/callback'
  });
  
  return response.data.auth_url;
};

export const handleOAuthCallback = async (provider: string, code: string, state: string) => {
  const response = await httpClient.get(`/auth/oauth/callback/${provider}`, {
    params: { code, state }
  });
  
  return response.data;
};

// OAuth Login Component
export const OAuthLoginButton: React.FC<{ provider: string }> = ({ provider }) => {
  const [isLoading, setIsLoading] = useState(false);
  const { login } = useAuth();

  const handleOAuthLogin = async () => {
    try {
      setIsLoading(true);
      const authUrl = await initiateOAuth(provider);
      const supported = await Linking.canOpenURL(authUrl);
      
      if (supported) {
        await Linking.openURL(authUrl);
        
        // Listen for deep link callback
        const handleUrl = (event: { url: string }) => {
          const url = new URL(event.url);
          const code = url.searchParams.get('code');
          const state = url.searchParams.get('state');
          
          if (code && state) {
            handleOAuthCallback(provider, code, state)
              .then((response) => {
                // Handle successful OAuth login
                setUser(response.user);
              })
              .catch((error) => {
                Alert.alert('OAuth Error', handleApiError(error));
              });
          }
        };
        
        const subscription = Linking.addEventListener('url', handleUrl);
        
        // Cleanup listener after 5 minutes
        setTimeout(() => {
          subscription?.remove();
        }, 5 * 60 * 1000);
      }
    } catch (error) {
      Alert.alert('OAuth Error', handleApiError(error));
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <TouchableOpacity
      style={[styles.oauthButton, { backgroundColor: getProviderColor(provider) }]}
      onPress={handleOAuthLogin}
      disabled={isLoading}
    >
      <Text style={styles.oauthButtonText}>
        {isLoading ? 'Connecting...' : `Continue with ${provider.charAt(0).toUpperCase() + provider.slice(1)}`}
      </Text>
    </TouchableOpacity>
  );
};
```

## Security Considerations

### 1. Token Storage
```typescript
// For production, use react-native-keychain for secure storage
import * as Keychain from 'react-native-keychain';

const storeTokenSecurely = async (token: string, type: 'access' | 'refresh') => {
  await Keychain.setInternetCredentials(`fitup_${type}_token`, 'user', token);
};

const getTokenSecurely = async (type: 'access' | 'refresh'): Promise<string | null> => {
  try {
    const credentials = await Keychain.getInternetCredentials(`fitup_${type}_token`);
    return credentials ? credentials.password : null;
  } catch (error) {
    return null;
  }
};

const removeTokensSecurely = async () => {
  await Promise.all([
    Keychain.resetInternetCredentials('fitup_access_token'),
    Keychain.resetInternetCredentials('fitup_refresh_token'),
  ]);
};
```

### 2. Network Security
```typescript
// Add certificate pinning for production
import { NetworkingModule } from 'react-native';

// Configure SSL pinning in native code
// iOS: Add certificate to bundle and configure in Info.plist
// Android: Add certificate to assets and configure in network_security_config.xml
```

### 3. Rate Limiting Awareness
The API implements rate limiting on sensitive endpoints:
- Login: 5 attempts per 15 minutes
- Register: 3 attempts per hour  
- Password Reset: 3 attempts per hour
- Token Refresh: 10 attempts per minute

Handle rate limiting gracefully in your app with appropriate user feedback.

## TypeScript Types

```typescript
// types/auth.ts
export interface User {
  id: string;
  username: string;
  name: string;
  bio: string;
  email: string;
  image: string;
  role: 'admin' | 'user';
  is_two_factor_enabled: boolean;
  created_at: string;
  updated_at: string;
}

export interface UserResponse {
  id: string;
  username: string;
  name: string;
  bio: string;
  email: string;
  image?: string;
  role: 'admin' | 'user';
  is_two_factor_enabled: boolean;
  created_at: string;
  updated_at: string;
}

export interface PublicUserResponse {
  id: string;
  username: string;
  name: string;
  bio: string;
  image?: string;
  role: 'admin' | 'user';
  created_at: string;
}

export interface Account {
  id: string;
  user_id: string;
  type: string;
  provider: string;
  provider_account_id: string;
  refresh_token: string;
  access_token: string;
  expires_at: number;
  token_type: string;
  scope: string;
  id_token: string;
  session_state: string;
}

export interface LoginRequest {
  identifier: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  name?: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token?: string;
  token_type: string;
  expires_in: number;
  user: User;
}

export interface UpdateUserRequest {
  username?: string;
  name?: string;
  bio?: string;
  image?: string;
}

export interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
}

export interface ForgotPasswordRequest {
  email: string;
}

export interface ResetPasswordRequest {
  token: string;
  new_password: string;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

export interface RefreshTokenResponse {
  access_token: string;
  refresh_token?: string;
  token_type: string;
  expires_in: number;
}

export interface TokenPair {
  access_token: string;
  refresh_token?: string;
  expires_in: number;
  token_type: string;
}

export interface ApiError {
  code: string;
  message: string;
}

// OAuth Types
export interface OAuthProvider {
  name: string;
  client_id: string;
  redirect_uri: string;
  auth_url: string;
  token_url: string;
  user_info_url: string;
  scopes: string[];
}

export interface OAuthUserInfo {
  id: string;
  email: string;
  name: string;
  username?: string;
  avatar_url?: string;
  email_verified: boolean;
}

export interface LinkAccountRequest {
  provider: string;
  code: string;
  state: string;
}
```

## Testing

### Unit Tests Example
```typescript
// __tests__/authService.test.ts
import { login, register } from '../services/authService';
import httpClient from '../services/httpClient';

jest.mock('../services/httpClient');
const mockedHttpClient = httpClient as jest.Mocked<typeof httpClient>;

describe('AuthService', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should login successfully', async () => {
    const mockResponse = {
      data: {
        access_token: 'mock_access_token',
        refresh_token: 'mock_refresh_token',
        token_type: 'Bearer',
        expires_in: 3600,
        user: { 
          id: '1', 
          email: 'test@example.com',
          username: 'testuser'
        }
      }
    };
    
    mockedHttpClient.post.mockResolvedValue(mockResponse);
    
    const result = await login({ 
      identifier: 'test@example.com', 
      password: 'password123' 
    });
    
    expect(result).toEqual(mockResponse.data);
    expect(mockedHttpClient.post).toHaveBeenCalledWith('/auth/login', {
      identifier: 'test@example.com',
      password: 'password123'
    });
  });

  it('should handle login error', async () => {
    const mockError = {
      response: {
        data: {
          code: 'INVALID_CREDENTIALS',
          message: 'Invalid email or password'
        }
      }
    };
    
    mockedHttpClient.post.mockRejectedValue(mockError);
    
    await expect(login({ 
      identifier: 'wrong@example.com', 
      password: 'wrongpassword' 
    })).rejects.toEqual(mockError);
  });
});
```

### Integration Tests
```typescript
// __tests__/auth.integration.test.ts
import { useAuth, AuthProvider } from '../hooks/useAuth';
import { renderHook, act } from '@testing-library/react-hooks';
import AsyncStorage from '@react-native-async-storage/async-storage';

describe('Auth Integration', () => {
  beforeEach(() => {
    AsyncStorage.clear();
  });

  it('should login and store tokens', async () => {
    const wrapper = ({ children }) => <AuthProvider>{children}</AuthProvider>;
    const { result } = renderHook(() => useAuth(), { wrapper });

    await act(async () => {
      await result.current.login({
        identifier: 'test@example.com',
        password: 'password123'
      });
    });

    expect(result.current.isAuthenticated).toBe(true);
    expect(result.current.user).toBeDefined();
    
    const storedToken = await AsyncStorage.getItem('access_token');
    expect(storedToken).toBeDefined();
  });
});
```

## Troubleshooting

### Common Issues

1. **Token Refresh Loop**: Make sure the refresh token endpoint doesn't require authentication
2. **CORS Errors**: Ensure your API Gateway has proper CORS configuration for mobile apps
3. **Deep Linking**: Test OAuth callbacks thoroughly on both iOS and Android
4. **Network Connectivity**: Handle offline scenarios and network timeouts
5. **Rate Limiting**: Show appropriate messages when limits are hit

### Debug Logging
```typescript
// Enable debug mode in development
if (__DEV__) {
  httpClient.interceptors.request.use(request => {
    console.log('🚀 API Request:', request.method?.toUpperCase(), request.url);
    console.log('📦 Request Data:', request.data);
    return request;
  });

  httpClient.interceptors.response.use(
    response => {
      console.log('✅ API Response:', response.config.url, response.status);
      console.log('📦 Response Data:', response.data);
      return response;
    },
    error => {
      console.error('❌ API Error:', error.config?.url, error.response?.status);
      console.error('📦 Error Data:', error.response?.data);
      return Promise.reject(error);
    }
  );
}
```

### Health Check
```typescript
// services/healthCheck.ts
export const checkAPIHealth = async (): Promise<boolean> => {
  try {
    const response = await httpClient.get('/health');
    return response.status === 200;
  } catch (error) {
    return false;
  }
};

// Usage in app startup
useEffect(() => {
  const checkHealth = async () => {
    const isHealthy = await checkAPIHealth();
    if (!isHealthy) {
      Alert.alert(
        'Service Unavailable', 
        'The Fit-Up service is currently unavailable. Please try again later.'
      );
    }
  };
  
  checkHealth();
}, []);
```

## Production Deployment

### Environment Configuration
```typescript
// config/environments.ts
const environments = {
  development: {
    API_BASE_URL: 'http://localhost:8080',
    OAUTH_REDIRECT_SCHEME: 'fitup-dev',
  },
  staging: {
    API_BASE_URL: 'https://staging-api.fitup.com',
    OAUTH_REDIRECT_SCHEME: 'fitup-staging',
  },
  production: {
    API_BASE_URL: 'https://api.fitup.com',
    OAUTH_REDIRECT_SCHEME: 'fitup',
  },
};

export default environments[process.env.NODE_ENV || 'development'];
```

### Security Checklist

- ✅ Use HTTPS in production
- ✅ Implement certificate pinning
- ✅ Store tokens securely (react-native-keychain)
- ✅ Handle token expiration with automatic refresh
- ✅ Implement proper error handling
- ✅ Validate user input on client side
- ✅ Handle rate limiting gracefully
- ✅ Clear sensitive data on app backgrounding
- ✅ Implement biometric authentication where appropriate
- ✅ Add request/response logging for debugging
- ✅ Implement proper session management

This comprehensive guide provides everything a React Native developer needs to integrate with the Fit-Up backend authentication system, including security best practices, error handling, and production-ready code examples.