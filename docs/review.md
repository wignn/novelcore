# Review Service

## Ringkasan

- Service internal: `ReviewService` (gRPC)
- Akses frontend: GraphQL Query/Mutation
- Tujuan: membuat dan menampilkan review novel

## Acuan Frontend (GraphQL)

### Ambil Reviews per Novel

Request:

```graphql
query Reviews($novelId: String!, $pagination: PaginationInput) {
  reviews(novelId: $novelId, pagination: $pagination) {
    id
    novelId
    accountId
    rating
    title
    content
    isSpoiler
    upvotes
    downvotes
    createdAt
  }
}
```

Variables:

```json
{
  "novelId": "nov_1",
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
    "reviews": [
      {
        "id": "rev_1",
        "novelId": "nov_1",
        "accountId": "acc_123",
        "rating": 5,
        "title": "Mantap",
        "content": "Plot rapi",
        "isSpoiler": false,
        "upvotes": 12,
        "downvotes": 1,
        "createdAt": "2026-04-07T10:00:00Z"
      }
    ]
  }
}
```

### Create Review

Request:

```graphql
mutation CreateReview($review: ReviewInput!) {
  createReview(review: $review) {
    id
    novelId
    accountId
    rating
    title
    content
    isSpoiler
    upvotes
    downvotes
    createdAt
  }
}
```

Variables:

```json
{
  "review": {
    "novelId": "nov_1",
    "accountId": "acc_123",
    "rating": 5,
    "title": "Suka",
    "content": "Karakter bagus",
    "isSpoiler": false
  }
}
```

Catatan validasi:

- `rating` harus dalam rentang 1 sampai 5.

### Delete Review

Request:

```graphql
mutation DeleteReview($id: String!) {
  deleteReview(id: $id) {
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
    "deleteReview": {
      "deletedId": "rev_1",
      "success": true,
      "message": "review deleted"
    }
  }
}
```
