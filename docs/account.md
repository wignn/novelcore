# Account Service

## Ringkasan

- Service: `AccountService`
- Transport: gRPC
- Tujuan: manajemen data akun user (create, read, list, edit, delete)

## Acuan Frontend (GraphQL)

Frontend berinteraksi melalui GraphQL endpoint `POST /graphql`.

### Create Account

Request:

```graphql
mutation CreateAccount($account: AccountInput!) {
  createAccount(account: $account) {
    id
    name
    email
    avatarUrl
    bio
    role
    createdAt
  }
}
```

Variables:

```json
{
  "account": {
    "name": "alice",
    "email": "alice@example.com",
    "password": "secret123"
  }
}
```

Response:

```json
{
  "data": {
    "createAccount": {
      "id": "acc_123",
      "name": "alice",
      "email": "alice@example.com",
      "avatarUrl": "",
      "bio": "",
      "role": "user",
      "createdAt": "2026-04-07T10:00:00Z"
    }
  }
}
```

### Get Accounts

Request:

```graphql
query Accounts($pagination: PaginationInput, $id: String) {
  accounts(pagination: $pagination, id: $id) {
    id
    name
    email
    avatarUrl
    bio
    role
    createdAt
  }
}
```

Variables:

```json
{
  "pagination": {
    "skip": 0,
    "take": 20
  }
}
```

Response:

```json
{
  "data": {
    "accounts": [
      {
        "id": "acc_123",
        "name": "alice",
        "email": "alice@example.com",
        "avatarUrl": "",
        "bio": "",
        "role": "user",
        "createdAt": "2026-04-07T10:00:00Z"
      }
    ]
  }
}
```

### Edit Account

Request:

```graphql
mutation EditAccount($id: String!, $account: EditAccountInput!) {
  editAccount(id: $id, account: $account) {
    id
    name
    email
    avatarUrl
    bio
    role
    createdAt
  }
}
```

Variables:

```json
{
  "id": "acc_123",
  "account": {
    "name": "Alice Update",
    "bio": "Webnovel reader"
  }
}
```

Response:

```json
{
  "data": {
    "editAccount": {
      "id": "acc_123",
      "name": "Alice Update",
      "email": "alice@example.com",
      "avatarUrl": "",
      "bio": "Webnovel reader",
      "role": "user",
      "createdAt": "2026-04-07T10:00:00Z"
    }
  }
}
```

### Delete Account

Request:

```graphql
mutation DeleteAccount($id: String!) {
  deleteAccount(id: $id) {
    deletedId
    success
    message
  }
}
```

Variables:

```json
{
  "id": "acc_123"
}
```

Response:

```json
{
  "data": {
    "deleteAccount": {
      "deletedId": "acc_123",
      "success": true,
      "message": "account deleted"
    }
  }
}
```

## RPC Endpoints

### 1. PostAccount

- Method: `PostAccount(PostAccountRequest) returns (PostAccountResponse)`
- Kegunaan: membuat akun baru

Request `PostAccountRequest`:

- `name` (string)
- `email` (string)
- `password` (string)

Response `PostAccountResponse`:

- `account` (Account)

### 2. GetAccount

- Method: `GetAccount(GetAccountRequest) returns (GetAccountResponse)`
- Kegunaan: ambil detail akun berdasarkan ID

Request `GetAccountRequest`:

- `id` (string)

Response `GetAccountResponse`:

- `account` (Account)

### 3. GetAccounts

- Method: `GetAccounts(GetAccountsRequest) returns (GetAccountsResponse)`
- Kegunaan: list akun dengan pagination

Request `GetAccountsRequest`:

- `skip` (uint64)
- `take` (uint64)

Response `GetAccountsResponse`:

- `accounts` (repeated Account)

### 4. EditAccount

- Method: `EditAccount(EditAccountRequest) returns (EditAccountResponse)`
- Kegunaan: update data akun

Request `EditAccountRequest`:

- `id` (string)
- `name` (string)
- `email` (string)
- `password` (string)
- `avatar_url` (string)
- `bio` (string)

Response `EditAccountResponse`:

- `message` (string)
- `success` (bool)
- `account` (Account)

### 5. DeleteAccount

- Method: `DeleteAccount(DeleteAccountRequest) returns (DeleteAccountResponse)`
- Kegunaan: hapus akun berdasarkan ID

Request `DeleteAccountRequest`:

- `id` (string)

Response `DeleteAccountResponse`:

- `message` (string)
- `success` (bool)
- `deletedID` (string)

## Message Utama

### Account

- `id` (string)
- `name` (string)
- `email` (string)
- `avatar_url` (string)
- `bio` (string)
- `role` (string)
- `created_at` (bytes)

## Contoh Payload

Create akun:

```json
{
  "name": "alice",
  "email": "alice@example.com",
  "password": "secret123"
}
```

List akun:

```json
{
  "skip": 0,
  "take": 20
}
```
