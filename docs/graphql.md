# GraphQL API Reference (Frontend)

## Base URL

- Endpoint: `POST http://localhost:8001/graphql`
- Playground: `http://localhost:8001/playground`

Contoh header request:

```http
Content-Type: application/json
Authorization: Bearer <access_token>
```

## Format Request Umum

```json
{
  "query": "query Example($id: String!) { novels(id: $id) { id title } }",
  "variables": {
    "id": "nov_1"
  }
}
```

## Format Response Umum

Sukses:

```json
{
  "data": {
    "...": "..."
  }
}
```

Gagal:

```json
{
  "errors": [
    {
      "message": "CreateReview: invalid parameter"
    }
  ],
  "data": null
}
```

## Query Endpoints

### accounts

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

### novels

```graphql
query Novels($pagination: PaginationInput, $id: String, $filter: NovelFilterInput, $query: String) {
  novels(pagination: $pagination, id: $id, filter: $filter, query: $query) {
    id
    title
    alternativeTitle
    description
    coverImageUrl
    author { id name }
    status
    novelType
    countryOfOrigin
    yearPublished
    totalChapters
    ratingAvg
    ratingCount
    viewCount
    bookmarkCount
    genres { id name slug }
    tags { id name slug }
    createdAt
    updatedAt
  }
}
```

### chapters

```graphql
query Chapters($novelId: String!, $pagination: PaginationInput) {
  chapters(novelId: $novelId, pagination: $pagination) {
    id
    novelId
    chapterNumber
    title
    translatorGroupId
    sourceUrl
    createdAt
    updatedAt
  }
}
```

### chapter

```graphql
query Chapter($id: String!) {
  chapter(id: $id) {
    id
    novelId
    chapterNumber
    title
    sourceUrl
  }
}
```

### authors

```graphql
query Authors($pagination: PaginationInput, $id: String) {
  authors(pagination: $pagination, id: $id) {
    id
    name
    bio
    createdAt
  }
}
```

### translationGroups

```graphql
query TranslationGroups($pagination: PaginationInput) {
  translationGroups(pagination: $pagination) {
    id
    name
    websiteUrl
    description
    createdAt
  }
}
```

### genres / tags

```graphql
query GenresTags {
  genres { id name slug }
  tags { id name slug }
}
```

### readingList

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

### reviews

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

### novelRanking

```graphql
query NovelRanking($period: RankingPeriod!, $sortBy: RankingSortBy!, $pagination: PaginationInput) {
  novelRanking(period: $period, sortBy: $sortBy, pagination: $pagination) {
    rank
    score
    change
    novel {
      id
      title
      ratingAvg
      viewCount
      bookmarkCount
    }
  }
}
```

## Mutation Endpoints

### Account

```graphql
mutation CreateAccount($account: AccountInput!) {
  createAccount(account: $account) {
    id
    name
    email
    role
    createdAt
  }
}
```

```graphql
mutation EditAccount($id: String!, $account: EditAccountInput!) {
  editAccount(id: $id, account: $account) {
    id
    name
    email
    avatarUrl
    bio
  }
}
```

```graphql
mutation DeleteAccount($id: String!) {
  deleteAccount(id: $id) {
    deletedId
    success
    message
  }
}
```

### Auth

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

```graphql
mutation RefreshToken($refreshToken: String!) {
  refreshToken(refreshToken: $refreshToken) {
    accessToken
    refreshToken
    expiresIn
  }
}
```

### Novel & Chapter

```graphql
mutation CreateNovel($novel: NovelInput!) {
  createNovel(novel: $novel) {
    id
    title
    createdAt
  }
}
```

```graphql
mutation UpdateNovel($id: String!, $novel: NovelInput!) {
  updateNovel(id: $id, novel: $novel) {
    id
    title
    updatedAt
  }
}
```

```graphql
mutation DeleteNovel($id: String!) {
  deleteNovel(id: $id) {
    deletedId
    success
    message
  }
}
```

```graphql
mutation CreateChapter($chapter: ChapterInput!) {
  createChapter(chapter: $chapter) {
    id
    novelId
    chapterNumber
    title
  }
}
```

```graphql
mutation UpdateChapter($id: String!, $chapter: ChapterInput!) {
  updateChapter(id: $id, chapter: $chapter) {
    id
    chapterNumber
    title
  }
}
```

```graphql
mutation DeleteChapter($id: String!) {
  deleteChapter(id: $id) {
    deletedId
    success
    message
  }
}
```

```graphql
mutation IncrementViewCount($novelId: String!) {
  incrementViewCount(novelId: $novelId)
}
```

### Author / Translation Group

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

### Reading List

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
  }
}
```

```graphql
mutation UpdateReadingList($id: String!, $entry: ReadingListInput!) {
  updateReadingList(id: $id, entry: $entry) {
    id
    status
    currentChapter
    rating
    notes
    isFavorite
  }
}
```

```graphql
mutation RemoveFromReadingList($id: String!) {
  removeFromReadingList(id: $id) {
    deletedId
    success
    message
  }
}
```

### Review

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
    createdAt
  }
}
```

```graphql
mutation DeleteReview($id: String!) {
  deleteReview(id: $id) {
    deletedId
    success
    message
  }
}
```

## Input Types Penting

`AccountInput`:

- `name: String!`
- `email: String!`
- `password: String!`

`LoginInput`:

- `email: String!`
- `password: String!`

`NovelInput`:

- `title: String!`
- `alternativeTitle: String`
- `description: String`
- `coverImageUrl: String`
- `authorId: String`
- `status: String`
- `novelType: String`
- `countryOfOrigin: String`
- `yearPublished: Int`
- `genreIds: [Int!]`
- `tagIds: [Int!]`

`ReadingListInput`:

- `novelId: String!`
- `status: String!`
- `currentChapter: Float`
- `rating: Int`
- `notes: String`
- `isFavorite: Boolean`

`ReviewInput`:

- `novelId: String!`
- `accountId: String!`
- `rating: Int!`
- `title: String`
- `content: String`
- `isSpoiler: Boolean`
