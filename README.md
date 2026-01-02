# GoAuth - gRPC Authentication Service

A production-ready authentication service built with gRPC, supporting traditional username/password authentication and Google OAuth 2.0.

## Features

- 🔐 **Traditional Authentication** - Username/password registration and login
- 🌐 **Google OAuth 2.0** - Sign in with Google
- 🎫 **JWT-based Authorization** - Secure token-based authentication
- 👥 **Role-based Access Control** - Multiple permission levels (user, admin, super_admin)
- 🔒 **Token Blacklisting** - Secure logout with token revocation
- 💾 **MongoDB Storage** - Persistent user data storage
- 🔄 **Account Linking** - Automatic linking of Google accounts to existing users

---

## Quick Start

### Prerequisites

- **Go** 1.24.4 or higher
- **MongoDB** (local, Docker, or cloud)
- **Google OAuth 2.0 credentials** (for Google login)
- **Protocol Buffers compiler** (for proto modifications)

### Installation

1. **Clone and install dependencies:**
   ```bash
   git clone <your-repo-url>
   cd GoAuth
   go mod download
   ```

2. **Configure environment variables:**
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` with your values:
   ```env
   # Server
   PORT=:50051
   
   # JWT Configuration
   JWT_SECRET=your-super-secret-key-here-min-32-chars
   JWT_EXPIRES_IN=60m
   
   # MongoDB
   MONGODB_URI=mongodb://localhost:27017
   
   # Google OAuth (required for Google login)
   GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
   ```

3. **Start MongoDB:**
   ```bash
   # Using systemd
   sudo systemctl start mongod
   
   # Using Docker
   docker run -d -p 27017:27017 --name mongodb mongo:latest
   ```

4. **Run the server:**
   ```bash
   go run cmd/api/main.go
   ```

   You should see:
   ```
   Connected to mongodb successfully
   gRPC server running on port :50051
   ```

---

## API Reference

### Authentication Methods

#### 1. Register
Create a new user account with username and password.

**Method:** `main.AuthService/Register`

**Request:**
```json
{
  "username": "johndoe",
  "password": "SecurePassword123!",
  "email": "john@example.com"
}
```

**Response:**
```json
{
  "status": true,
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

---

#### 2. Login
Authenticate with username and password.

**Method:** `main.AuthService/Login`

**Request:**
```json
{
  "username": "johndoe",
  "password": "SecurePassword123!"
}
```

**Response:**
```json
{
  "status": true,
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

---

#### 3. GoogleLogin
Authenticate using Google OAuth 2.0 ID token.

**Method:** `main.AuthService/GoogleLogin`

**Request:**
```json
{
  "id_token": "eyJhbGciOiJSUzI1NiIsImtpZCI..."
}
```

**Response:**
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIs...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "673a8f6b2c1e4d5a7f8b9c0d",
    "username": "John Doe",
    "email": "john@gmail.com",
    "role": "user",
    "googleId": "103456789012345678901",
    "picture": "https://lh3.googleusercontent.com/..."
  }
}
```

**Google OAuth Flow:**
1. Client implements Google Sign-In (frontend)
2. User authenticates with Google
3. Client receives Google ID token
4. Client sends ID token to this RPC
5. Server verifies token with Google
6. Server creates/retrieves user
7. Server returns JWT tokens

---

#### 4. Logout
Invalidate current user's JWT token.

**Method:** `main.AuthService/Logout`

**Headers:**
```
authorization: Bearer <your-jwt-token>
```

**Request:**
```json
{}
```

**Response:**
```json
{
  "status": true
}
```

---

#### 5. ChangeRole
Change a user's role (requires super_admin role).

**Method:** `main.AuthService/ChangeRole`

**Headers:**
```
authorization: Bearer <super-admin-jwt-token>
```

**Request:**
```json
{
  "id": "673a8f6b2c1e4d5a7f8b9c0d",
  "role": "admin"
}
```

**Response:**
```json
{
  "status": true
}
```

**Allowed Roles:** `user`, `admin`, `super_admin`

---

## Integration Guide

### Client Implementation

#### JavaScript/TypeScript (gRPC-Web)

```javascript
import { AuthServiceClient } from './proto/main_grpc_web_pb';
import { GoogleLoginRequest } from './proto/main_pb';

const client = new AuthServiceClient('http://localhost:50051');

// Google OAuth Login
async function loginWithGoogle(idToken) {
  const request = new GoogleLoginRequest();
  request.setIdToken(idToken);
  
  const response = await client.googleLogin(request, {});
  
  // Store tokens
  localStorage.setItem('accessToken', response.getAccessToken());
  localStorage.setItem('refreshToken', response.getRefreshToken());
  
  return response.getUser();
}

// Traditional Login
async function login(username, password) {
  const request = new LoginRequest();
  request.setUsername(username);
  request.setPassword(password);
  
  const response = await client.login(request, {});
  localStorage.setItem('token', response.getToken());
}

// Authenticated Requests
function createAuthMetadata() {
  const token = localStorage.getItem('accessToken');
  return { authorization: `Bearer ${token}` };
}
```

#### React Example with Google Sign-In

```jsx
import { GoogleOAuthProvider, GoogleLogin } from '@react-oauth/google';
import { useState } from 'react';

function App() {
  const [user, setUser] = useState(null);

  const handleGoogleSuccess = async (credentialResponse) => {
    try {
      // Call your gRPC server
      const result = await grpcClient.googleLogin({
        idToken: credentialResponse.credential
      });
      
      setUser(result.user);
      localStorage.setItem('accessToken', result.accessToken);
    } catch (error) {
      console.error('Login failed:', error);
    }
  };

  return (
    <GoogleOAuthProvider clientId="YOUR_GOOGLE_CLIENT_ID">
      <GoogleLogin
        onSuccess={handleGoogleSuccess}
        onError={() => console.log('Login Failed')}
      />
    </GoogleOAuthProvider>
  );
}
```

#### Go Client Example

```go
package main

import (
    "context"
    "log"
    pb "path/to/proto/gen"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    client := pb.NewAuthServiceClient(conn)
    
    // Login
    resp, err := client.Login(context.Background(), &pb.LoginRequest{
        Username: "johndoe",
        Password: "password123",
    })
    
    token := resp.GetToken()
    
    // Authenticated request
    md := metadata.Pairs("authorization", "Bearer "+token)
    ctx := metadata.NewOutgoingContext(context.Background(), md)
    
    logoutResp, err := client.Logout(ctx, &pb.EmptyRequest{})
}
```

---

## Google OAuth Setup

### 1. Create OAuth Credentials

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing
3. Navigate to **APIs & Services** → **Credentials**
4. Click **Create Credentials** → **OAuth 2.0 Client ID**
5. Set application type to **Web application**

### 2. Configure OAuth Client

**Authorized JavaScript origins:**
```
http://localhost:3000
http://localhost:8000
https://yourdomain.com
```

**Authorized redirect URIs:**
```
http://localhost:3000/callback
https://yourdomain.com/callback
```

### 3. Update Configuration

Add your Client ID to `.env`:
```env
GOOGLE_CLIENT_ID=123456789-xxxxx.apps.googleusercontent.com
```

---

## Production Deployment

### Security Checklist

- [ ] Use strong `JWT_SECRET` (minimum 32 characters, randomly generated)
- [ ] Enable TLS/SSL for gRPC server
- [ ] Use environment variables for all secrets
- [ ] Never commit `.env` file to version control
- [ ] Implement rate limiting
- [ ] Set up monitoring and logging
- [ ] Use MongoDB authentication
- [ ] Configure firewall rules
- [ ] Enable CORS properly for web clients
- [ ] Implement refresh token rotation
- [ ] Add request validation and sanitization

### Environment Variables

```env
# Production Configuration
PORT=:50051
JWT_SECRET=$(openssl rand -base64 32)
JWT_EXPIRES_IN=15m
MONGODB_URI=mongodb://user:pass@production-db:27017/auth?authSource=admin
GOOGLE_CLIENT_ID=your-production-client-id.apps.googleusercontent.com

# Optional
JWT_ISSUER=GoAuth
LOG_LEVEL=info
```

### Building for Production

```bash
# Build binary
go build -o bin/goauth cmd/api/main.go

# Run
./bin/goauth
```

### Docker Deployment

```dockerfile
# Dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o goauth cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/goauth .
EXPOSE 50051
CMD ["./goauth"]
```

```yaml
# docker-compose.yml
version: '3.8'
services:
  goauth:
    build: .
    ports:
      - "50051:50051"
    environment:
      - MONGODB_URI=mongodb://mongodb:27017
      - JWT_SECRET=${JWT_SECRET}
      - GOOGLE_CLIENT_ID=${GOOGLE_CLIENT_ID}
    depends_on:
      - mongodb
    restart: unless-stopped
  
  mongodb:
    image: mongo:latest
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db
    restart: unless-stopped

volumes:
  mongodb_data:
```

---

## Testing

### Using grpcurl

```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# List services
grpcurl -plaintext localhost:50051 list

# Register a user
grpcurl -plaintext -d '{
  "username": "testuser",
  "password": "Test123!",
  "email": "test@example.com"
}' localhost:50051 main.AuthService/Register

# Login
grpcurl -plaintext -d '{
  "username": "testuser",
  "password": "Test123!"
}' localhost:50051 main.AuthService/Login

# Logout (with token)
grpcurl -plaintext \
  -H "authorization: Bearer YOUR_TOKEN" \
  -d '{}' localhost:50051 main.AuthService/Logout
```

---

## Database Schema

### Users Collection

```javascript
{
  _id: ObjectId,
  username: String,      // Display name
  email: String,         // Unique email address
  password: String,      // Argon2 hashed (empty for Google users)
  role: String,          // "user", "admin", or "super_admin"
  google_id: String,     // Google user ID (optional)
  picture: String        // Profile picture URL (optional)
}
```

### Indexes (Recommended)

```javascript
db.users.createIndex({ "email": 1 }, { unique: true })
db.users.createIndex({ "google_id": 1 }, { sparse: true })
db.users.createIndex({ "username": 1 })
```

---

## Troubleshooting

### Common Issues

**"Error connecting to the database"**
- Verify MongoDB is running: `sudo systemctl status mongod`
- Check `MONGODB_URI` in `.env`
- Test connection: `mongosh`

**"invalid google id token"**
- Ensure `GOOGLE_CLIENT_ID` matches your Google Console
- Verify ID token hasn't expired (they're short-lived)
- Check you're using `id_token`, not `access_token`

**"JWT_SECRET environment variable is not set"**
- Verify `.env` file exists and is properly loaded
- Ensure `JWT_SECRET` is defined

**"Token has been revoked (logged out)"**
- User has logged out, they need to log in again
- This is expected behavior

**MongoDB socket file error**
- Remove stale socket: `sudo rm -f /tmp/mongodb-27017.sock`
- Restart MongoDB: `sudo systemctl restart mongod`

---

## Development

### Regenerate Protocol Buffers

After modifying `proto/main.proto`:

```bash
protoc --go_out=. --go-grpc_out=. proto/main.proto
```

### Run Tests

```bash
go test ./...
```

### Project Structure

```
.
├── cmd/api/main.go              # Server entry point
├── internal/
│   ├── api/
│   │   ├── handlers/            # RPC handlers
│   │   └── interceptors/        # Authentication middleware
│   ├── models/                  # Data models
│   └── repositories/
│       └── mongodb/             # Database operations
├── pkg/utils/                   # Utilities (JWT, OAuth, etc.)
├── proto/                       # Protocol buffer definitions
│   ├── main.proto
│   └── gen/                     # Generated code
├── .env.example                 # Environment template
└── README.md                    # This file
```

---

## License

[Your License Here]

## Contributing

[Your Contributing Guidelines Here]

## Support

For issues and questions:
- Create an issue on GitHub
- Contact: [your-email@example.com]

---

**Built with ❤️ using Go, gRPC, and MongoDB**

## Setup

### Prerequisites
- Go 1.24.4 or higher
- MongoDB (local or cloud instance)
- Google OAuth 2.0 credentials

### Installation

1. **Clone the repository** (if you haven't already)
   ```bash
   git clone <your-repo-url>
   cd GoAuth
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up Google OAuth 2.0**
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project or select an existing one
   - Enable the Google+ API
   - Go to "Credentials" → "Create Credentials" → "OAuth 2.0 Client ID"
   - Configure the OAuth consent screen
   - Create an OAuth 2.0 Client ID (Web application type)
   - Copy the Client ID

4. **Configure environment variables**
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` with your values:
   ```env
   MONGODB_URI=mongodb://localhost:27017
   JWT_SECRET=your-super-secret-jwt-key-here
   JWT_ISSUER=GoAuth
   JWT_EXPIRES_IN=24h
   GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
   PORT=:50051
   ```

5. **Start MongoDB**
   ```bash
   # If using Docker
   docker run -d -p 27017:27017 --name mongodb mongo:latest
   
   # Or start your local MongoDB instance
   mongod
   ```

6. **Run the server**
   ```bash
   go run cmd/api/main.go
   ```

## API Reference

### RPC Methods

#### 1. Register
Creates a new user account with username and password.

**Request:**
```protobuf
message RegisterRequest {
    string username = 1;
    string password = 2;
    string email = 3;
}
```

**Response:**
```protobuf
message LoginResponse {
    bool status = 1;
    string token = 2;
}
```

#### 2. Login
Authenticates a user with username and password.

**Request:**
```protobuf
message LoginRequest {
    string username = 1;
    string password = 2;
}
```

**Response:**
```protobuf
message LoginResponse {
    bool status = 1;
    string token = 2;
}
```

#### 3. GoogleLogin
Authenticates a user using Google OAuth 2.0 ID token.

**Request:**
```protobuf
message GoogleLoginRequest {
    string id_token = 1;  // Google ID token from client
}
```

**Response:**
```protobuf
message GoogleLoginResponse {
    string access_token = 1;
    string refresh_token = 2;
    User user = 3;
}
```

**Flow:**
1. Client obtains Google ID token from Google Sign-In
2. Client sends ID token to GoogleLogin RPC
3. Server verifies token with Google
4. Server creates/retrieves user from database
5. Server returns JWT access/refresh tokens

#### 4. Logout
Invalidates the current user's JWT token.

**Request:**
```protobuf
message EmptyRequest {}
```

**Response:**
```protobuf
message LogoutResponse {
    bool status = 1;
}
```

**Headers Required:**
```
authorization: Bearer <jwt-token>
```

#### 5. ChangeRole
Changes a user's role (super_admin only).

**Request:**
```protobuf
message ChangeRoleRequest {
    string id = 1;
    string role = 2;
}
```

**Response:**
```protobuf
message ChangeRoleResponse {
    bool status = 1;
}
```

**Headers Required:**
```
authorization: Bearer <jwt-token>
```

**Authorization:** Only users with `super_admin` role can access this endpoint.

## Authentication Flow

### Traditional Authentication
1. **Registration**: User registers with username, email, and password
2. **Password Hashing**: Password is hashed using bcrypt before storage
3. **Login**: User provides credentials, receives JWT token
4. **Subsequent Requests**: Client includes JWT in `authorization` header
5. **Logout**: Token is blacklisted until expiration

### Google OAuth Flow
1. **Client-Side**: User signs in with Google, obtains ID token
2. **Token Verification**: Server verifies ID token with Google's API
3. **User Lookup**: Server checks if user exists by Google ID or email
4. **User Creation**: If new user, creates account with Google data
5. **Token Generation**: Server generates JWT access and refresh tokens
6. **Response**: Returns tokens and user information

## Security Features

### JWT Token Management
- Tokens expire based on `JWT_EXPIRES_IN` setting (default: 15 minutes)
- Tokens are signed with `JWT_SECRET`
- Blacklist mechanism prevents use of logged-out tokens
- Automatic cleanup of expired blacklisted tokens every 2 minutes

### Password Security
- Passwords are hashed using bcrypt
- Passwords are never stored in plain text
- Google OAuth users don't have passwords

### Authorization
- Role-based access control (RBAC)
- JWT interceptor validates all protected endpoints
- Context-based user information access

## Database Schema

### Users Collection
```javascript
{
  _id: ObjectId,
  username: String,
  email: String,
  password: String,      // Hashed (empty for Google OAuth users)
  role: String,          // "user", "admin", or "super_admin"
  google_id: String,     // Google user ID (for OAuth users)
  picture: String        // Profile picture URL (for OAuth users)
}
```

## Testing with grpcurl

### 1. Register a new user
```bash
grpcurl -plaintext -d '{
  "username": "testuser",
  "password": "password123",
  "email": "test@example.com"
}' localhost:50051 main.AuthService/Register
```

### 2. Login
```bash
grpcurl -plaintext -d '{
  "username": "testuser",
  "password": "password123"
}' localhost:50051 main.AuthService/Login
```

### 3. Google Login
```bash
grpcurl -plaintext -d '{
  "id_token": "YOUR_GOOGLE_ID_TOKEN_HERE"
}' localhost:50051 main.AuthService/GoogleLogin
```

### 4. Logout
```bash
grpcurl -plaintext \
  -H "authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{}' localhost:50051 main.AuthService/Logout
```

### 5. Change Role (super_admin only)
```bash
grpcurl -plaintext \
  -H "authorization: Bearer YOUR_SUPER_ADMIN_JWT_TOKEN" \
  -d '{
    "id": "USER_ID_HERE",
    "role": "admin"
  }' localhost:50051 main.AuthService/ChangeRole
```

## Client Implementation Example

### Google OAuth Client (JavaScript/React)
```javascript
// 1. Install Google Sign-In library
// npm install @react-oauth/google

// 2. Configure Google OAuth
import { GoogleOAuthProvider, GoogleLogin } from '@react-oauth/google';

function App() {
  return (
    <GoogleOAuthProvider clientId="YOUR_GOOGLE_CLIENT_ID">
      <GoogleLogin
        onSuccess={async (credentialResponse) => {
          // Send ID token to your gRPC server
          const response = await grpcClient.googleLogin({
            idToken: credentialResponse.credential
          });
          
          // Store access token
          localStorage.setItem('accessToken', response.accessToken);
          localStorage.setItem('refreshToken', response.refreshToken);
        }}
        onError={() => {
          console.log('Login Failed');
        }}
      />
    </GoogleOAuthProvider>
  );
}
```

## Environment Variables Reference

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `MONGODB_URI` | MongoDB connection string | - | Yes |
| `JWT_SECRET` | Secret key for signing JWT tokens | - | Yes |
| `JWT_ISSUER` | JWT token issuer identifier | GoAuth | No |
| `JWT_EXPIRES_IN` | Token expiration duration | 15m | No |
| `GOOGLE_CLIENT_ID` | Google OAuth 2.0 Client ID | - | Yes (for Google auth) |
| `PORT` | gRPC server port | :50051 | Yes |

## Troubleshooting

### "invalid google id token" error
- Verify `GOOGLE_CLIENT_ID` matches your Google Cloud Console configuration
- Ensure the ID token hasn't expired (they're short-lived)
- Check that you're using the ID token, not the access token

### "Error connecting to the database"
- Verify MongoDB is running
- Check `MONGODB_URI` is correct
- Ensure network connectivity to MongoDB

### "JWT_SECRET environment variable is not set"
- Ensure `.env` file exists and is properly loaded
- Verify `JWT_SECRET` is defined in `.env`

### Token not being accepted after logout
- This is expected behavior - logged out tokens are blacklisted
- User must log in again to get a new token

## Development

### Regenerate Protocol Buffers
If you modify `proto/main.proto`:
```bash
protoc --go_out=. --go-grpc_out=. proto/main.proto
```

### Run Tests
```bash
go test ./...
```

## License

[Your License Here]

## Contributing

[Your Contributing Guidelines Here]
