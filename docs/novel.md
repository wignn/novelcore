# Novel Service

## Ringkasan

- Domain: novel, chapter, author, translation group, genre, tag, ranking
- Akses frontend: GraphQL (`/graphql`)
- Transport internal: gRPC (`NovelService`)

## Acuan Frontend (GraphQL)

### List / Search Novel

Request:

```graphql
query Novels($pagination: PaginationInput, $filter: NovelFilterInput, $query: String) {
  novels(pagination: $pagination, filter: $filter, query: $query) {
    id
    title
    coverImageUrl
    status
    novelType
    ratingAvg
    totalChapters
    genres { id name slug }
    tags { id name slug }
    author { id name }
    createdAt
    updatedAt
  }
}
```

Variables:

```json
{
  "pagination": { "skip": 0, "take": 12 },
  "query": "isekai",
  "filter": {
    "status": "ONGOING",
    "sortBy": "rating_avg",
    "sortOrder": "desc"
  }
}
```

Response:

```json
{
  "data": {
    "novels": [
      {
        "id": "nov_1",
        "title": "Example Novel",
        "coverImageUrl": "https://cdn/image.jpg",
        "status": "ONGOING",
        "novelType": "WEBNOVEL",
        "ratingAvg": 4.6,
        "totalChapters": 120,
        "genres": [{ "id": 1, "name": "Fantasy", "slug": "fantasy" }],
        "tags": [{ "id": 10, "name": "Action", "slug": "action" }],
        "author": { "id": "auth_1", "name": "John Doe" },
        "createdAt": "2026-04-07T10:00:00Z",
        "updatedAt": "2026-04-07T10:00:00Z"
      }
    ]
  }
}
```

### Detail Novel by ID

Request:

```graphql
query NovelById($id: String!) {
  novels(id: $id) {
    id
    title
    alternativeTitle
    description
    coverImageUrl
    status
    novelType
    countryOfOrigin
    yearPublished
    totalChapters
    ratingAvg
    ratingCount
    viewCount
    bookmarkCount
    genres { id name }
    tags { id name }
    author { id name bio }
  }
}
```

### Create Novel

Request:

```graphql
mutation CreateNovel($novel: NovelInput!) {
  createNovel(novel: $novel) {
    id
    title
    status
    novelType
    createdAt
    updatedAt
  }
}
```

Variables:

```json
{
  "novel": {
    "title": "New Novel",
    "description": "desc",
    "status": "ONGOING",
    "novelType": "WEBNOVEL",
    "genreIds": [1, 2],
    "tagIds": [10]
  }
}
```

### Update Novel

Request:

```graphql
mutation UpdateNovel($id: String!, $novel: NovelInput!) {
  updateNovel(id: $id, novel: $novel) {
    id
    title
    updatedAt
  }
}
```

### Delete Novel

Request:

```graphql
mutation DeleteNovel($id: String!) {
  deleteNovel(id: $id) {
    deletedId
    success
    message
  }
}
```

### Chapter Operations

Create Chapter:

```graphql
mutation CreateChapter($chapter: ChapterInput!) {
  createChapter(chapter: $chapter) {
    id
    novelId
    chapterNumber
    title
    sourceUrl
    createdAt
  }
}
```

List Chapters:

```graphql
query Chapters($novelId: String!, $pagination: PaginationInput) {
  chapters(novelId: $novelId, pagination: $pagination) {
    id
    chapterNumber
    title
    sourceUrl
    createdAt
  }
}
```

Delete Chapter:

```graphql
mutation DeleteChapter($id: String!) {
  deleteChapter(id: $id) {
    deletedId
    success
    message
  }
}
```

### Author / Translation Group / Taxonomy

Create Author:

```graphql
mutation CreateAuthor($author: AuthorInput!) {
  createAuthor(author: $author) {
    id
    name
    bio
    createdAt
  }
}
```

List Authors:

```graphql
query Authors($pagination: PaginationInput) {
  authors(pagination: $pagination) {
    id
    name
    bio
    createdAt
  }
}
```

Create Translation Group:

```graphql
mutation CreateTranslationGroup($group: TranslationGroupInput!) {
  createTranslationGroup(group: $group) {
    id
    name
    websiteUrl
    description
    createdAt
  }
}
```

List Genres/Tags:

```graphql
query GenresTags {
  genres { id name slug }
  tags { id name slug }
}
```

### Ranking & View Count

Ranking:

```graphql
query Ranking($period: RankingPeriod!, $sortBy: RankingSortBy!, $pagination: PaginationInput) {
  novelRanking(period: $period, sortBy: $sortBy, pagination: $pagination) {
    rank
    score
    change
    novel {
      id
      title
      viewCount
      bookmarkCount
      ratingAvg
    }
  }
}
```

Increment view count:

```graphql
mutation IncrementView($novelId: String!) {
  incrementViewCount(novelId: $novelId)
}
```

Response increment view count:

```json
{
  "data": {
    "incrementViewCount": 321
  }
}
```
