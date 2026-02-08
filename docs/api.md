# Enterprise Auth System API Documentation

## 1. Health Check

### Endpoint Path
`GET /health`

### Description
Checks the operational status of the application and its dependencies (PostgreSQL and Redis).

### Request Structure
- **Headers**: None
- **URL Parameters**: None
- **Body**: None

### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "postgres": "up",
  "redis": "up"
}
```

### Failure Responses
- **Status Code**: `503 Service Unavailable`
- **Body**:
```json
{
  "postgres": "down",
  "redis": "up"
}
```

---

## 2. User Signup

### Endpoint Path
`POST /auth/signup`

### Description
Registers a new user in the system. The default role assigned is `MEMBER` if not specified.

### Request Structure
- **Headers**: `Content-Type: application/json`
- **Body**:
```json
{
  "name": "string (required)",
  "email": "string (required, unique)",
  "phone": "string (required)",
  "password": "string (required)",
  "role": "string (optional: 'HEAD', 'LEAD', 'MEMBER')"
}
```

### Success Response
- **Status Code**: `201 Created`
- **Body**:
```json
{
  "message": "User created successfully"
}
```

### Failure Responses
- **Status Code**: `400 Bad Request` (Invalid JSON or missing fields)
- **Status Code**: `500 Internal Server Error` (Database error, e.g., duplicate email)
- **Example Error**:
```json
{
  "error": "Key: 'User.Email' Error:Field validation for 'Email' failed on the 'required' tag"
}
```

---

## 3. User Login

### Endpoint Path
`POST /auth/login`

### Description
Authenticates a user and returns a pair of Access and Refresh tokens.

### Request Structure
- **Headers**: `Content-Type: application/json`
- **Body**:
```json
{
  "email": "string (required)",
  "password": "string (required)"
}
```

### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "access_token": "ey... (JWT)",
  "refresh_token": "ey... (JWT)"
}
```

### Failure Responses
- **Status Code**: `400 Bad Request`
- **Status Code**: `401 Unauthorized` (Invalid credentials)
- **Example Error**:
```json
{
  "error": "invalid credentials"
}
```

---

## 4. Refresh Token

### Endpoint Path
`POST /auth/refresh`

### Description
Rotates the Refresh Token. Accepts a valid existing Refresh Token and issues a new pair of Access and Refresh tokens, revoking the old one if rotation logic is strict (implementation dependent).

### Request Structure
- **Headers**: `Content-Type: application/json`
- **Body**:
```json
{
  "refresh_token": "string (required)"
}
```

### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "access_token": "ey... (new JWT)",
  "refresh_token": "ey... (new JWT)"
}
```

### Failure Responses
- **Status Code**: `400 Bad Request`
- **Status Code**: `401 Unauthorized` (Invalid or expired token)
- **Example Error**:
```json
{
  "error": "token has invalid claims: token is expired"
}
```

---

## 5. User Logout

### Endpoint Path
`POST /auth/logout`

### Description
Logs out the user by revoking the Refresh Token in Redis.

### Request Structure
- **Headers**: `Content-Type: application/json`
- **Body**:
```json
{
  "refresh_token": "string (required)"
}
```

### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "message": "Logged out"
}
```

### Failure Responses
- **Status Code**: `400 Bad Request`

---

## 6. Create Admin

### Endpoint Path
`POST /admin/create`

### Description
Creates a new System Administrator. This endpoint is protected by a special secret key.

### Request Structure
- **Headers**: 
  - `Content-Type: application/json`
  - `x-admin-auth-key`: `string (required, matches ADMIN_SECRET_KEY env var)`
- **Body**:
```json
{
  "name": "string (required)",
  "email": "string (required, unique)",
  "phone": "string (required)",
  "password": "string (required)"
}
```

### Success Response
- **Status Code**: `201 Created`
- **Body**:
```json
{
  "message": "Admin created successfully"
}
```

### Failure Responses
- **Status Code**: `403 Forbidden` (Invalid or missing `x-admin-auth-key`)
- **Status Code**: `400 Bad Request`
- **Status Code**: `500 Internal Server Error`
- **Example Error**:
```json
{
  "error": "Invalid or missing admin key"
}
```

---

## 7. Get Profile (Protected Example)

### Endpoint Path
`GET /api/profile`

### Description
Returns the profile information of the authenticated user. Proves that the Access Token is valid.

### Request Structure
- **Headers**: 
  - `Authorization`: `Bearer <access_token>`
- **Body**: None

### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "message": "Access granted",
  "userID": "123",
  "role": "MEMBER"
}
```

### Failure Responses
- **Status Code**: `401 Unauthorized` (Missing or invalid token)
- **Example Error**:
```json
{
  "error": "Authorization header required"
}
```

---

## 8. Get Students

### Endpoint Path
`GET /api/students`

### Description
Retrieves a list of students with advanced search, filtering, sorting, and pagination. Requires authentication.

### Request Structure
- **Headers**: 
  - `Authorization`: `Bearer <access_token>`
- **URL Parameters**:
  - `search` (optional): Global search term for full_name, email, and phone.
  - `program_status` (optional): Filter by status (true/false).
  - `sort_by` (optional): Field to sort by (`full_name`, `email`, `created_at`). Default: `created_at`.
  - `order` (optional): Sort order (`asc`, `desc`). Default: `desc`.
  - `page` (optional): Page number (min 1). Default: `1`.
  - `limit` (optional): Items per page (max 100). Default: `10`.

- **Example**: `/api/students?search=albin&program_status=true&sort_by=full_name&order=asc&page=2&limit=5`

### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "data": [
    {
      "id": 1,
      "full_name": "Albin Aji",
      "phone": "9876543210",
      "email": "albin@example.com",
      "program_status": true,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "total": 21,
    "page": 2,
    "limit": 5,
    "total_pages": 5
  }
}
```

### Failure Responses
- **Status Code**: `400 Bad Request` (Invalid parameters)
- **Status Code**: `401 Unauthorized` (Missing or invalid token)
- **Status Code**: `500 Internal Server Error` (Database error)
- **Status Code**: `429 Too Many Requests` (Rate limit exceeded)

---

## 9. Create Student

### Endpoint Path
`POST /api/students`

### Description
Creates a new student record. Requires authentication.

### Request Structure
- **Headers**: 
  - `Content-Type: application/json`
  - `Authorization`: `Bearer <access_token>`
- **Body**:
```json
{
  "full_name": "string (required)",
  "email": "string (required, unique)",
  "phone": "string (required, min 10 chars)",
  "program_status": "boolean (optional, default false)"
}
```

### Success Response
- **Status Code**: `201 Created`
- **Body**:
```json
{
  "id": 1,
  "full_name": "New Student",
  "email": "new.student@test.com",
  "phone": "1234567890",
  "program_status": true,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### Failure Responses
- **Status Code**: `400 Bad Request` (Missing fields, invalid email/phone, duplicate email)
- **Example Error**:
```json
{
  "error": "full_name, email, and phone are required"
}
```
- **Status Code**: `401 Unauthorized`
- **Status Code**: `500 Internal Server Error`
