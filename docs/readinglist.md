# ReadingList Service

## Ringkasan

- Service internal: `ReadingListService` (gRPC)
- Akses frontend: GraphQL Query/Mutation
- Tujuan: menyimpan progress baca per akun

## Acuan Frontend (GraphQL)

### Ambil Reading List

Request:

```graphql
query ReadingList($accountId: String!, $status: String, $pagination: PaginationInput) {
  readingList(accountId: $accountId, status: $status, pagination: $pagination) {
    id
    novelId
    status
    currentChapter
    rating
    notes
    isFavorite
    createdAt
    updatedAt
  }
}
```

Variables:

```json
{
  "accountId": "acc_123",
  "status": "READING",
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
    "readingList": [
      {
        "id": "rl_1",
        "novelId": "nov_1",
        "status": "READING",
        "currentChapter": 12.5,
        "rating": 5,
        "notes": "bagus",
        "isFavorite": true,
        "createdAt": "2026-04-07T10:00:00Z",
        "updatedAt": "2026-04-07T10:00:00Z"
      }
    ]
  }
}
```

### Tambah Reading List

Request:

```graphql
mutation AddToReadingList($accountId: String!, $entry: ReadingListInput!) {
  addToReadingList(accountId: $accountId, entry: $entry) {
    id
    novelId
    status
    currentChapter
    rating
    notes
    isFavorite
    createdAt
    updatedAt
  }
}
```

Variables:

```json
{
  "accountId": "acc_123",
  "entry": {
    "novelId": "nov_1",
    "status": "READING",
    "currentChapter": 1,
    "rating": 4,
    "notes": "baru mulai",
    "isFavorite": false
  }
}
```

### Update Reading List

Request:

```graphql
mutation UpdateReadingList($id: String!, $entry: ReadingListInput!) {
  updateReadingList(id: $id, entry: $entry) {
    id
    novelId
    status
    currentChapter
    rating
    notes
    isFavorite
    updatedAt
  }
}
```

### Remove Reading List

Request:

```graphql
mutation RemoveReadingList($id: String!) {
  removeFromReadingList(id: $id) {
    deletedId
    success
    message
  }
}
```

Response hapus:

```json
{
  "data": {
    "removeFromReadingList": {
      "deletedId": "rl_1",
      "success": true,
      "message": "reading list removed"
    }
  }
}
```
