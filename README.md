# GoAuth - gRPC Authentication Service

Production-ready gRPC authentication service with JWT tokens, traditional auth, and Google OAuth 2.0.

## Quick Start

```bash
# Setup
cp .env.example .env  # Configure JWT_SECRET, GOOGLE_CLIENT_ID, MONGODB_URI
sudo systemctl start mongod
go run cmd/api/main.go
```

---

## gRPC Endpoints

### 1. Register - `main.AuthService/Register`

```bash
grpcurl -plaintext -d '{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "SecurePass123!"
}' localhost:50051 main.AuthService/Register

# Response: { "status": true, "token": "eyJhbG..." }
```

### 2. Login - `main.AuthService/Login`

```bash
grpcurl -plaintext -d '{
  "username": "johndoe",
  "password": "SecurePass123!"
}' localhost:50051 main.AuthService/Login

# Response: { "status": true, "token": "eyJhbG..." }
```

### 3. GoogleLogin - `main.AuthService/GoogleLogin`

```bash
grpcurl -plaintext -d '{
  "id_token": "GOOGLE_ID_TOKEN_FROM_CLIENT"
}' localhost:50051 main.AuthService/GoogleLogin

# Response: {
#   "accessToken": "eyJhbG...",
#   "refreshToken": "eyJhbG...",
#   "user": { "id": "...", "email": "...", "googleId": "...", "picture": "..." }
# }
```

**Flow:**
1. Client gets Google ID token (use Google Sign-In on frontend)
2. Send ID token to this endpoint
3. Server verifies with Google, creates/finds user
4. Returns JWT tokens + user info
5. Existing email? Google ID is linked to account

### 4. Logout - `main.AuthService/Logout`

```bash
grpcurl -plaintext \
  -H "authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{}' localhost:50051 main.AuthService/Logout

# Response: { "status": true }
# Token is blacklisted until expiration
```

### 5. ChangeRole - `main.AuthService/ChangeRole` (super_admin only)

```bash
grpcurl -plaintext \
  -H "authorization: Bearer SUPER_ADMIN_TOKEN" \
  -d '{"id": "USER_ID", "role": "admin"}' \
  localhost:50051 main.AuthService/ChangeRole

# Response: { "status": true }
# Roles: user | admin | super_admin
```

---

## Authentication

**Traditional:** Register/Login → JWT token → Include in `authorization: Bearer <token>` header

**Google OAuth:**
1. Frontend: Google Sign-In → Get ID token
2. Send ID token to `GoogleLogin`
3. Receive JWT access/refresh tokens
4. Use access token for authenticated requests

**Logout:** Token blacklisted, must login again

---

## Database

MongoDB `users` collection:
```javascript
{
  email: String (unique),
  username: String,
  password: String,      // bcrypt hash (empty for Google-only)
  role: String,          // user|admin|super_admin
  google_id: String,     // optional
  picture: String        // optional
}
```

---

## Environment Variables

```env
PORT=:50051
MONGODB_URI=mongodb://localhost:27017
JWT_SECRET=your-secret-min-32-chars
JWT_EXPIRES_IN=15m
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
```


