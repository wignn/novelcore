# Auth Service

## Ringkasan

- Service: `AuthService`
- Transport internal: gRPC
- Akses frontend: GraphQL mutation
- Tujuan: login dan refresh token

## Acuan Frontend (GraphQL)

### Login

Request:

```graphql
mutation Login($account: LoginInput!) {
	login(account: $account) {
		id
		email
		backendToken {
			accessToken
			refreshToken
			expiresIn
		}
	}
}
```

Variables:

```json
{
	"account": {
		"email": "alice@example.com",
		"password": "secret123"
	}
}
```

Response:

```json
{
	"data": {
		"login": {
			"id": "acc_123",
			"email": "alice@example.com",
			"backendToken": {
				"accessToken": "eyJ...",
				"refreshToken": "eyJ...",
				"expiresIn": 1712499999
			}
		}
	}
}
```

### Refresh Token

Request:

```graphql
mutation RefreshToken($refreshToken: String!) {
	refreshToken(refreshToken: $refreshToken) {
		accessToken
		refreshToken
		expiresIn
	}
}
```

Variables:

```json
{
	"refreshToken": "eyJ..."
}
```

Response:

```json
{
	"data": {
		"refreshToken": {
			"accessToken": "eyJ...new",
			"refreshToken": "eyJ...new",
			"expiresIn": 1712509999
		}
	}
}
```

## RPC Endpoints (Internal)

### Login

- Method: `Login(PostAuthRequest) returns (PostAuthResponse)`

Request `PostAuthRequest`:

- `email` (string)
- `password` (string)

Response `PostAuthResponse`:

- `auth` (Auth)

### RefreshToken

- Method: `RefreshToken(PostRefreshTokenRequest) returns (BackendToken)`

Request `PostRefreshTokenRequest`:

- `refreshToken` (string)

Response `BackendToken`:

- `accessToken` (string)
- `refreshToken` (string)
- `expiresAt` (uint64)

