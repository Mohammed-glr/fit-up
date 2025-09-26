# Fit-Up React Native Quick Reference

## API Endpoints Quick Reference

### Base URL
- **Development**: `http://localhost:8080`  (API Gateway)
- **Production**: `https://api.fitup.com`

### Authentication Endpoints

```bash
# Register new user
POST /auth/register
{
  "username": "johndoe",
  "email": "john@example.com", 
  "password": "password123",
  "name": "John Doe"
}

# Login user
POST /auth/login
{
  "identifier": "john@example.com",  # email or username
  "password": "password123"
}

# Logout
POST /auth/logout
{}

# Refresh token
POST /auth/refresh-token
{
  "refresh_token": "your_refresh_token"
}

# Validate token
POST /auth/validate-token
{}

# Change password (requires auth)
POST /auth/change-password
{
  "current_password": "oldpassword",
  "new_password": "newpassword123"
}

# Forgot password
POST /auth/forgot-password
{
  "email": "john@example.com"
}

# Reset password
POST /auth/reset-password
{
  "token": "reset_token_from_email",
  "new_password": "newpassword123"
}

# Get user profile
GET /auth/{username}
```

### OAuth Endpoints

```bash
# Initiate OAuth
POST /auth/oauth/{provider}  # provider: google, github, facebook
{
  "provider": "google",
  "redirect_url": "fitup://oauth/callback"
}

# OAuth callback (handled automatically)
GET /auth/oauth/callback/{provider}?code=...&state=...

# Link OAuth account (requires auth)
POST /auth/link/{provider}
{
  "provider": "google",
  "code": "oauth_code",
  "state": "oauth_state"
}

# Unlink OAuth account (requires auth)
DELETE /auth/unlink/{provider}

# Get linked accounts (requires auth)
GET /auth/linked-accounts
```

### Future Service Endpoints (Coming Soon)

```bash
# Message Service
GET /messages          # Get user messages
POST /messages         # Send message
GET /messages/{id}     # Get specific message
DELETE /messages/{id}  # Delete message

# Schema Service
GET /schemas           # Get available schemas
POST /schemas/validate # Validate data against schema
```

## Response Formats

### Successful Authentication Response
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "refresh_token_string",
  "token_type": "Bearer",
  "expires_in": 3600,
  "user": {
    "id": "user-uuid",
    "username": "johndoe",
    "name": "John Doe",
    "bio": "Fitness enthusiast",
    "email": "john@example.com",
    "image": "https://example.com/avatar.jpg",
    "role": "user",
    "is_two_factor_enabled": false,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

### Token Refresh Response
```json
{
  "access_token": "new_access_token",
  "refresh_token": "new_refresh_token",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

### Error Response Format
```json
{
  "code": "INVALID_CREDENTIALS",
  "message": "Invalid email or password"
}
```

## Quick Setup Code

### 1. Install Dependencies
```bash
npm install axios @react-native-async-storage/async-storage
# For secure storage (recommended):
npm install react-native-keychain
# For OAuth deep linking:
npm install @react-native-community/async-storage
```

### 2. HTTP Client Setup
```typescript
// services/api.ts
import axios from 'axios';
import AsyncStorage from '@react-native-async-storage/async-storage';

const api = axios.create({
  baseURL: 'http://localhost:8080', // Change for production
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests
api.interceptors.request.use(async (config) => {
  const token = await AsyncStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle token refresh
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      const refreshToken = await AsyncStorage.getItem('refresh_token');
      if (refreshToken) {
        try {
          const response = await axios.post('/auth/refresh-token', {
            refresh_token: refreshToken
          });
          const { access_token } = response.data;
          await AsyncStorage.setItem('access_token', access_token);
          // Retry original request
          error.config.headers.Authorization = `Bearer ${access_token}`;
          return api.request(error.config);
        } catch (refreshError) {
          // Refresh failed, redirect to login
          await AsyncStorage.multiRemove(['access_token', 'refresh_token']);
        }
      }
    }
    return Promise.reject(error);
  }
);

export default api;
```

### 3. Auth Service
```typescript
// services/auth.ts
import api from './api';

export const authAPI = {
  register: (data: RegisterData) => api.post('/auth/register', data),
  login: (data: LoginData) => api.post('/auth/login', data),
  logout: () => api.post('/auth/logout', {}),
  refreshToken: (token: string) => api.post('/auth/refresh-token', { refresh_token: token }),
  validateToken: () => api.post('/auth/validate-token', {}),
  changePassword: (data: ChangePasswordData) => api.post('/auth/change-password', data),
  forgotPassword: (email: string) => api.post('/auth/forgot-password', { email }),
  resetPassword: (data: ResetPasswordData) => api.post('/auth/reset-password', data),
  getProfile: (username: string) => api.get(`/auth/${username}`),
  
  // OAuth
  initiateOAuth: (provider: string) => api.post(`/auth/oauth/${provider}`, { provider }),
  linkAccount: (provider: string, data: any) => api.post(`/auth/link/${provider}`, data),
  unlinkAccount: (provider: string) => api.delete(`/auth/unlink/${provider}`),
  getLinkedAccounts: () => api.get('/auth/linked-accounts'),
};
```

### 4. Auth Hook
```typescript
// hooks/useAuth.ts
import React, { createContext, useContext, useState, useEffect } from 'react';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { authAPI } from '../services/auth';

const AuthContext = createContext({});

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    checkAuthState();
  }, []);

  const checkAuthState = async () => {
    try {
      const token = await AsyncStorage.getItem('access_token');
      if (token) {
        const response = await authAPI.validateToken();
        setUser(response.data.user);
      }
    } catch (error) {
      await AsyncStorage.multiRemove(['access_token', 'refresh_token']);
    } finally {
      setLoading(false);
    }
  };

  const login = async (credentials) => {
    const response = await authAPI.login(credentials);
    setUser(response.data.user);
    await AsyncStorage.setItem('access_token', response.data.access_token);
    if (response.data.refresh_token) {
      await AsyncStorage.setItem('refresh_token', response.data.refresh_token);
    }
  };

  const logout = async () => {
    try {
      await authAPI.logout();
    } finally {
      setUser(null);
      await AsyncStorage.multiRemove(['access_token', 'refresh_token']);
    }
  };

  return (
    <AuthContext.Provider value={{ 
      user, 
      login, 
      logout, 
      loading,
      isAuthenticated: !!user 
    }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
```

### 5. Login Component
```typescript
// components/LoginForm.tsx
import React, { useState } from 'react';
import { View, TextInput, TouchableOpacity, Text, Alert } from 'react-native';
import { useAuth } from '../hooks/useAuth';

const LoginForm = () => {
  const [identifier, setIdentifier] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();

  const handleLogin = async () => {
    try {
      setLoading(true);
      await login({ identifier, password });
    } catch (error) {
      Alert.alert('Error', error.response?.data?.message || 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <View style={{ padding: 20 }}>
      <TextInput
        placeholder="Email or Username"
        value={identifier}
        onChangeText={setIdentifier}
        style={{ borderWidth: 1, padding: 10, marginBottom: 10 }}
      />
      <TextInput
        placeholder="Password"
        value={password}
        onChangeText={setPassword}
        secureTextEntry
        style={{ borderWidth: 1, padding: 10, marginBottom: 20 }}
      />
      <TouchableOpacity
        onPress={handleLogin}
        disabled={loading}
        style={{ 
          backgroundColor: loading ? '#ccc' : '#007AFF', 
          padding: 15, 
          alignItems: 'center' 
        }}
      >
        <Text style={{ color: 'white' }}>
          {loading ? 'Logging in...' : 'Login'}
        </Text>
      </TouchableOpacity>
    </View>
  );
};
```

## Common Error Codes

| Code | Meaning | Solution |
|------|---------|----------|
| `INVALID_CREDENTIALS` | Wrong email/password | Check credentials |
| `USER_ALREADY_EXISTS` | Email taken | Use different email |
| `USERNAME_ALREADY_EXISTS` | Username taken | Choose different username |
| `TOO_MANY_ATTEMPTS` | Rate limited | Wait before retrying |
| `INVALID_TOKEN` | Token expired/invalid | Refresh token or re-login |
| `UNAUTHORIZED` | Not authenticated | Login required |
| `INVALID_INPUT` | Validation failed | Check input format |
| `REFRESH_TOKEN_NOT_FOUND` | No refresh token | Re-login required |
| `REFRESH_TOKEN_EXPIRED` | Refresh token expired | Re-login required |

## Rate Limits

| Endpoint | Limit | Window |
|----------|-------|--------|
| `/auth/login` | 5 attempts | 15 minutes |
| `/auth/register` | 3 attempts | 1 hour |
| `/auth/forgot-password` | 3 attempts | 1 hour |
| `/auth/reset-password` | 3 attempts | 1 hour |
| `/auth/refresh-token` | 10 attempts | 1 minute |

## TypeScript Interfaces

```typescript
interface User {
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

interface LoginData {
  identifier: string; // email or username
  password: string;
}

interface RegisterData {
  username: string;
  email: string;
  password: string;
  name?: string;
}

interface ChangePasswordData {
  current_password: string;
  new_password: string;
}

interface ResetPasswordData {
  token: string;
  new_password: string;
}

interface AuthResponse {
  access_token: string;
  refresh_token?: string;
  token_type: string;
  expires_in: number;
  user: User;
}

interface RefreshTokenResponse {
  access_token: string;
  refresh_token?: string;
  token_type: string;
  expires_in: number;
}

interface TokenPair {
  access_token: string;
  refresh_token?: string;
  expires_in: number;
  token_type: string;
}

interface ApiError {
  code: string;
  message: string;
}

interface Account {
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
```

## OAuth Deep Linking Setup

### iOS (Info.plist)
```xml
<key>CFBundleURLTypes</key>
<array>
  <dict>
    <key>CFBundleURLName</key>
    <string>fitup.oauth</string>
    <key>CFBundleURLSchemes</key>
    <array>
      <string>fitup</string>
    </array>
  </dict>
</array>
```

### Android (android/app/src/main/AndroidManifest.xml)
```xml
<activity
  android:name=".MainActivity"
  android:exported="true"
  android:launchMode="singleTop">
  <intent-filter>
    <action android:name="android.intent.action.VIEW" />
    <category android:name="android.intent.category.DEFAULT" />
    <category android:name="android.intent.category.BROWSABLE" />
    <data android:scheme="fitup" />
  </intent-filter>
</activity>
```

## Security Checklist

- ✅ Use HTTPS in production
- ✅ Store tokens securely (react-native-keychain recommended)
- ✅ Handle token expiration with automatic refresh
- ✅ Implement proper error handling
- ✅ Validate user input on client side
- ✅ Handle rate limiting gracefully
- ✅ Use certificate pinning in production
- ✅ Clear sensitive data on app backgrounding
- ✅ Implement biometric authentication where appropriate
- ✅ Add request/response logging for debugging

## Testing URLs

### Development
```
Base URL: http://localhost:8080
Health Check: GET http://localhost:8080/health
Test Login: POST http://localhost:8080/auth/login
```

### Service Architecture
```
API Gateway (8080) → Auth Service (8081)
                   → Message Service (8082) [Coming Soon]
                   → Schema Service (8083) [Coming Soon]
```

## Common Patterns

### Error Handling
```typescript
const handleApiError = (error: any): string => {
  if (error.response?.data?.code) {
    switch (error.response.data.code) {
      case 'INVALID_CREDENTIALS':
        return 'Invalid email or password';
      case 'TOO_MANY_ATTEMPTS':
        return 'Too many attempts. Please try again later';
      default:
        return error.response.data.message;
    }
  }
  return 'Network error. Please check your connection';
};
```

### Token Management
```typescript
const refreshTokenIfNeeded = async () => {
  const token = await AsyncStorage.getItem('access_token');
  // Check if token is about to expire (implement JWT decode)
  // Refresh if needed
};
```

### Loading States
```typescript
const [authState, setAuthState] = useState({
  isLoading: false,
  isAuthenticated: false,
  user: null,
  error: null,
});
```

This quick reference provides everything a React Native developer needs to get started quickly with the Fit-Up authentication API!