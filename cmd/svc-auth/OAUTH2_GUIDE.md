# OAuth2 Implementation Guide

## Overview
This service now supports OAuth2 authentication with multiple providers (Google, Facebook, GitHub), refresh token rotation, and multi-provider account linking.

## Features
✅ **OAuth2-Compliant Token Response** - Returns `access_token`, `refresh_token`, `token_type`, `expires_in`, `scope`
✅ **Multiple OAuth Providers** - Google, Facebook, GitHub (easily extensible)
✅ **Multi-Provider Accounts** - Users can link multiple OAuth providers to one account
✅ **Smart Auto-Linking** - Verified emails automatically link providers
✅ **Refresh Token Rotation** - Secure token refresh with 30-day expiration
✅ **Token Revocation** - Explicit token revocation support
✅ **JWT Access Tokens** - Signed JWT tokens with configurable TTL
✅ **AES-256-GCM Encryption** - OAuth provider tokens encrypted at rest

## Setup

### 1. Database Migration
Run the migration to create required tables:
```bash
cd cmd/svc-auth
make migrate.up
```

This creates:
- `user_oauth_providers` - Links users to OAuth providers
- `refresh_tokens` - Stores refresh token hashes
- Adds profile fields to `users` table

### 2. Environment Configuration
Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

**Critical Settings:**
- `OAUTH_ENCRYPTION_KEY` - MUST be exactly 32 bytes (e.g., `my-super-secret-32-byte-key!!!`)
- `JWT_SECRET` - Minimum 32 characters
- `JWT_ACCESS_TOKEN_TTL` - Default: `1h`
- `JWT_REFRESH_TOKEN_TTL` - Default: `720h` (30 days)

### 3. OAuth Provider Setup

#### Google OAuth
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing
3. Enable Google+ API
4. Create OAuth 2.0 credentials
5. Add authorized redirect URI: `http://localhost:8080/v1/auth/oauth/google/callback`
6. Copy Client ID and Client Secret to `.env`

#### Facebook OAuth
1. Go to [Facebook Developers](https://developers.facebook.com/)
2. Create an app
3. Add Facebook Login product
4. Add redirect URI: `http://localhost:8080/v1/auth/oauth/facebook/callback`
5. Copy App ID and App Secret to `.env`

#### GitHub OAuth
1. Go to [GitHub Settings > Developer settings > OAuth Apps](https://github.com/settings/developers)
2. Create new OAuth App
3. Set callback URL: `http://localhost:8080/v1/auth/oauth/github/callback`
4. Copy Client ID and Client Secret to `.env`

### 4. Build & Run
```bash
go build
./svc-auth
```

## API Endpoints

### Traditional Login (Password)
**POST** `/v1/auth/login`
```json
{
  "data": {
    "email": "user@example.com",
    "password": "password123"
  }
}
```
Response:
```json
{
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "random_token_string",
    "token_type": "Bearer",
    "expires_in": 3600,
    "scope": "openid profile email"
  }
}
```

### Refresh Token
**POST** `/v1/auth/refresh`
```json
{
  "refresh_token": "previous_refresh_token"
}
```

### Revoke Token
**POST** `/v1/auth/revoke`
```json
{
  "token": "refresh_token_to_revoke"
}
```

### OAuth Login Flow

#### Step 1: Get Authorization URL
**GET** `/v1/auth/oauth/{provider}/login?state=random_csrf_state&redirect_uri=optional`

Response:
```json
{
  "authorization_url": "https://accounts.google.com/o/oauth2/auth?..."
}
```

#### Step 2: User Authorizes (external)
Redirect user to `authorization_url`. After authorization, provider redirects back to your callback.

#### Step 3: Handle Callback
**GET** `/v1/auth/oauth/{provider}/callback?code=auth_code&state=csrf_state&redirect_uri=optional`

Response:
```json
{
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "random_token",
    "token_type": "Bearer",
    "expires_in": 3600,
    "scope": "openid profile email"
  },
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    ...
  },
  "is_new_user": true
}
```

### Link Additional Provider (Authenticated)
**POST** `/v1/auth/link/{provider}?code=auth_code&redirect_uri=optional`

Headers:
```
Authorization: Bearer {access_token}
```

### Unlink Provider (Authenticated)
**DELETE** `/v1/auth/unlink/{provider}`

Headers:
```
Authorization: Bearer {access_token}
```

### Get Linked Providers (Authenticated)
**GET** `/v1/auth/providers`

Headers:
```
Authorization: Bearer {access_token}
```

Response:
```json
{
  "providers": ["google", "github"]
}
```

## OAuth Flow Diagrams

### New User Registration via OAuth
```
Client → GET /oauth/google/login
      ← authorization_url

Client → [User authorizes on Google]
      → GET /oauth/google/callback?code=...
      ← access_token + refresh_token (is_new_user: true)
```

### Existing User Login via OAuth
```
Client → GET /oauth/github/login
      ← authorization_url

Client → [User authorizes on GitHub]
      → GET /oauth/github/callback?code=...
      ← access_token + refresh_token (is_new_user: false)
```

### Link Additional Provider
```
Client → GET /oauth/facebook/login (get URL)
      ← authorization_url

Client → [User authorizes on Facebook]
      → POST /link/facebook?code=... (with Bearer token)
      ← 200 OK
```

## Security Features

### Token Encryption
OAuth provider tokens (access_token, refresh_token) are encrypted using AES-256-GCM before storage.

### Refresh Token Security
- Refresh tokens are hashed (SHA-256) before storage
- Tokens have 30-day expiration
- Single-use tokens (rotation on refresh)
- Explicit revocation support

### Account Linking Safety
- Prevents linking a provider already linked to another account
- Prevents unlinking the last authentication method
- Auto-links only verified emails

### CSRF Protection
Use random `state` parameter in OAuth flow to prevent CSRF attacks.

## Database Schema

### users
```sql
- id (PK)
- email (unique)
- password (nullable - OAuth-only users don't have password)
- name, first_name, last_name, avatar_url (profile fields)
- email_verified (boolean)
- locale
- created_at, updated_at
```

### user_oauth_providers
```sql
- id (PK)
- user_id (FK → users.id)
- provider (google/facebook/github)
- provider_user_id (unique per provider)
- access_token (encrypted)
- refresh_token (encrypted)
- token_expiry
- email, name, first_name, last_name, avatar_url (provider profile)
- email_verified
- locale
- created_at, updated_at

UNIQUE(user_id, provider)
UNIQUE(provider, provider_user_id)
```

### refresh_tokens
```sql
- id (PK)
- user_id (FK → users.id)
- token (SHA-256 hash, unique)
- expires_at
- revoked_at (nullable)
- created_at

INDEX(user_id)
INDEX(token)
```

## Troubleshooting

### "Encryption key must be 32 bytes"
Set `OAUTH_ENCRYPTION_KEY` to exactly 32 characters:
```bash
OAUTH_ENCRYPTION_KEY="my-super-secret-32-byte-key!!!"
```

### "OAuth provider not configured"
Check:
1. Provider is enabled: `OAUTH_GOOGLE_ENABLED=true`
2. Client ID and Secret are set
3. Scopes are comma-separated: `OAUTH_GOOGLE_SCOPES=openid,profile,email`

### "Failed to exchange authorization code"
Verify:
1. Redirect URI matches exactly (including trailing slash)
2. Code is used immediately (codes expire quickly)
3. Client secret is correct

### "Cannot unlink last authentication method"
Users must have at least one authentication method (password OR OAuth provider).

## Adding New OAuth Providers

1. **Create Provider Implementation**
   ```go
   // cmd/svc-auth/oauth/newprovider.go
   func NewNewProvider(cfg *Config) Provider {
       // Implement Provider interface
   }
   ```

2. **Update Provider Factory**
   ```go
   // cmd/svc-auth/oauth/provider.go
   case "newprovider":
       return NewNewProvider(config), nil
   ```

3. **Add Configuration**
   ```go
   // cmd/svc-auth/config/config.go
   type OAuth struct {
       NewProvider OAuthProvider4 `mapstructure:",squash"`
   }
   ```

4. **Update getOAuthConfig**
   ```go
   // cmd/svc-auth/service/user.go
   case "newprovider":
       return &oauth.Config{...}
   ```

## Testing

### Manual Testing with cURL

1. **Get OAuth URL:**
```bash
curl http://localhost:8080/v1/auth/oauth/google/login?state=test123
```

2. **Visit the URL in browser, authorize, get code from redirect**

3. **Exchange code for tokens:**
```bash
curl -X GET "http://localhost:8080/v1/auth/oauth/google/callback?code=YOUR_CODE&state=test123"
```

4. **Refresh token:**
```bash
curl -X POST http://localhost:8080/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"YOUR_REFRESH_TOKEN"}'
```

5. **Link another provider:**
```bash
curl -X POST "http://localhost:8080/v1/auth/link/github?code=GITHUB_CODE" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Production Checklist

- [ ] Set strong 32-byte encryption key
- [ ] Set strong JWT secret (min 32 chars)
- [ ] Configure OAuth redirect URLs for production domain
- [ ] Enable only required OAuth providers
- [ ] Set appropriate token TTLs
- [ ] Enable HTTPS
- [ ] Implement rate limiting
- [ ] Add logging and monitoring
- [ ] Set up database backups
- [ ] Review and test all security features

## Architecture Notes

### Clean Architecture Layers
```
HTTP Handlers (httpapi/)
    ↓
Service Layer (service/)
    ↓
Repository Layer (repository/)
    ↓
Entities (entity/)
```

### OAuth Provider Layer
```
OAuth Providers (oauth/)
- provider.go     (interface & factory)
- google.go       (Google implementation)
- facebook.go     (Facebook implementation)
- github.go       (GitHub implementation)
```

### Key Design Decisions
1. **Multi-provider support** via separate `user_oauth_providers` table
2. **Token encryption** for security at rest
3. **Refresh token rotation** for enhanced security
4. **Smart auto-linking** for UX (verified emails only)
5. **Safety checks** to prevent account lockout

## Support
For issues or questions, refer to the main project documentation or create an issue.
