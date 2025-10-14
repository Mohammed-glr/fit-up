# Update User Role Endpoint

## Overview
This endpoint allows authenticated users to update their role to either `user`, `coach`, or `trainer`.

## Endpoint Details

**URL**: `/api/auth/update-role`  
**Method**: `PUT`  
**Authentication**: Required (JWT Bearer Token)

## Request

### Headers
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

### Body
```json
{
  "role": "coach"
}
```

### Parameters
- `role` (string, required): The role to assign to the user. Must be one of:
  - `user` - Regular user
  - `coach` - Fitness coach/instructor

## Response

### Success Response (200 OK)
```json
{
  "message": "Role updated successfully",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "john_doe",
    "name": "John Doe",
    "bio": "Fitness enthusiast",
    "email": "john@example.com",
    "image": "https://example.com/avatar.jpg",
    "role": "coach",
    "is_two_factor_enabled": false,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-15T12:00:00Z"
  }
}
```

### Error Responses

#### 400 Bad Request
Invalid role provided:
```json
{
  "error": "invalid input"
}
```

#### 401 Unauthorized
Missing or invalid authentication token:
```json
{
  "error": "unauthorized"
}
```

#### 404 Not Found
User not found:
```json
{
  "error": "user not found"
}
```

#### 500 Internal Server Error
Server error:
```json
{
  "error": "internal server error"
}
```

## Example Usage

### cURL
```bash
curl -X PUT https://api.fitup.com/api/auth/update-role \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "role": "coach"
  }'
```

### JavaScript (Fetch)
```javascript
const response = await fetch('https://api.fitup.com/api/auth/update-role', {
  method: 'PUT',
  headers: {
    'Authorization': `Bearer ${accessToken}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    role: 'coach'
  })
});

const data = await response.json();
console.log(data);
```

### TypeScript (Axios)
```typescript
import axios from 'axios';

interface UpdateRoleRequest {
  role: 'user' | 'coach';
}

interface UpdateRoleResponse {
  message: string;
  user: {
    id: string;
    username: string;
    name: string;
    bio: string;
    email: string;
    image: string | null;
    role: 'user' | 'coach' | 'trainer';
    is_two_factor_enabled: boolean;
    created_at: string;
    updated_at: string;
  };
}

const updateUserRole = async (role: UpdateRoleRequest['role']) => {
  try {
    const response = await axios.put<UpdateRoleResponse>(
      'https://api.fitup.com/api/auth/update-role',
      { role },
      {
        headers: {
          'Authorization': `Bearer ${accessToken}`,
          'Content-Type': 'application/json'
        }
      }
    );
    
    return response.data;
  } catch (error) {
    console.error('Error updating role:', error);
    throw error;
  }
};

// Usage
await updateUserRole('coach');
```

## Notes

- Users can only update their own role (based on the JWT token)
- The `admin` role cannot be self-assigned through this endpoint (requires admin privileges)
- Role changes take effect immediately
- After updating the role, the user object in the JWT token remains unchanged until the next login or token refresh
- Consider logging out and back in after role change for best experience

## Implementation Details

### Database
The role is stored in the `users` table:
```sql
UPDATE users 
SET role = $1, updated_at = NOW()
WHERE id = $2
```

### Available Roles
- `user`: Default role for regular users
- `coach`: For fitness coaches and instructors
- `admin`: System administrators (cannot be self-assigned)

## Security Considerations

- Only authenticated users can update their role
- JWT token must be valid and not expired
- Role validation is performed on both client and server side
- Admin role cannot be assigned through this endpoint
