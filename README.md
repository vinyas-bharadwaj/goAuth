# GoAuth - Microservices Architecture 🔐

A multi-service authentication system with gRPC backend and FastAPI gateway.

## Features

- **User Registration** - Sign up with email and password
- **User Login** - Authenticate and receive JWT tokens
- **Token Validation** - Verify JWT token authenticity
- **Secure Storage** - Passwords hashed with bcrypt
- **SQLite Database** - Lightweight local database
- **gRPC API** - High-performance RPC communication

## Quick Start

### Prerequisites
- Go 1.24+
- Git

### Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/vinyas-bharadwaj/goAuth.git
   cd goAuth
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set environment variables**
   
   Create a `.env` file in the project root:
   ```bash
   JWT_SECRET=your-super-secret-jwt-key-here
   ```

4. **Run the server**
   ```bash
   go run server/*.go
   ```
   
   The server will start on `localhost:50051`

## API Reference

### gRPC Service: `AuthService`

#### 1. Signup
Register a new user account.
```protobuf
rpc Signup (SignupRequest) returns (AuthResponse);
```

#### 2. Login  
Authenticate an existing user.
```protobuf
rpc Login (LoginRequest) returns (AuthResponse);
```

#### 3. ValidateToken
Verify a JWT token's validity.
```protobuf
rpc ValidateToken (ValidateRequest) returns (ValidateResponse);
```

## Testing with grpcurl

Install grpcurl: `go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest`

### Register a new user
```bash
grpcurl -plaintext -d '{
  "email": "user@example.com",
  "password": "securepassword123"
}' localhost:50051 auth.AuthService/Signup
```

### Login
```bash
grpcurl -plaintext -d '{
  "email": "user@example.com", 
  "password": "securepassword123"
}' localhost:50051 auth.AuthService/Login
```

### Validate token
```bash
grpcurl -plaintext -d '{
  "token": "your-jwt-token-here"
}' localhost:50051 auth.AuthService/ValidateToken
```

## Project Structure

```
goAuth/
├── server/           # Server implementation
│   ├── main.go      # Entry point
│   ├── handler.go   # gRPC handlers
│   ├── db.go        # Database operations
│   └── jwt.go       # JWT utilities
├── proto/           # Protocol Buffers
│   ├── auth.proto   # Service definition
│   ├── auth.pb.go   # Generated Go code
│   └── auth_grpc.pb.go
├── go.mod           # Go module
├── .env             # Environment variables
└── README.md        # This file
```

## Development

### Regenerate Protocol Buffers
If you modify `proto/auth.proto`, regenerate the Go code:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/auth.proto
```

### Build
```bash
go build ./...
```

### Run Tests
```bash
go test ./...
```

## Security Features

- **Password Hashing**: Uses bcrypt with default cost
- **JWT Tokens**: Signed with HMAC SHA-256
- **SQL Injection Protection**: Parameterized queries
- **Environment Variables**: Sensitive data in `.env`

## Database

The service uses SQLite with a simple schema:

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE,
    password TEXT
);
```

Database file: `users.db` (created automatically)

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `JWT_SECRET` | Secret key for JWT signing | Yes |

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is open source and available under the [MIT License](LICENSE).

---

**Made with ❤️ and Go**