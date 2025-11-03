# Auth Service

Microservice for authentication with JWT, OAuth2, and multi-factor authentication.

## Features

- ✅ User registration with email validation and password strength requirements
- ✅ Login with JWT token generation (access + refresh)
- ✅ Token refresh endpoint
- ✅ Password reset flow
- ✅ Email verification flow
- ✅ Multi-factor authentication (TOTP)
- ✅ Device management (list, revoke)
- ✅ Session management
- ✅ Brute force protection and account lockout
- ✅ Rate limiting on auth endpoints

## API Endpoints

### Public Endpoints

- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/password/reset` - Request password reset
- `POST /api/v1/auth/password/reset/confirm` - Confirm password reset
- `POST /api/v1/auth/verify-email` - Verify email address

### Protected Endpoints

- `POST /api/v1/auth/logout` - Logout (invalidate token)
- `GET /api/v1/auth/devices` - List user devices
- `DELETE /api/v1/auth/devices/:id` - Revoke device

## Environment Variables

- `SERVER_PORT` - Server port (default: 8080)
- `SERVER_HOST` - Server host (default: 0.0.0.0)
- `DATABASE_URI` - MongoDB connection URI
- `DATABASE_NAME` - Database name (default: streamverse)
- `JWT_SECRET_KEY` - JWT secret key (required)
- `LOG_LEVEL` - Log level (default: info)

## Security Features

- Password hashing with bcrypt
- Account lockout after 5 failed login attempts (30 minutes)
- Rate limiting (10 requests per minute on auth endpoints)
- Device fingerprinting
- MFA support with TOTP

## Running

```bash
go run main.go
```

## Docker

```bash
docker build -t auth-service .
docker run -p 8080:8080 auth-service
```

