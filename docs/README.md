# Dokumentasi API Backend Novel Update

Dokumentasi ini mencakup endpoint setiap service pada backend, termasuk:

- Endpoint publik GraphQL
- Kontrak internal gRPC per service
- Ringkasan request/response utama

## Arsitektur Endpoint

- Public API: GraphQL di `http://localhost:8001/graphql`
- Playground: `http://localhost:8001/playground`
- Internal API antar-service: gRPC (port service ditentukan lewat env `PORT`, default lokal `50051`)

Catatan:

- Dari `compose.yml`, hanya service GraphQL yang diexpose ke host.
- gRPC service (`account`, `auth`, `novel`, `readinglist`, `review`) berjalan internal antarkontainer.

## Daftar Dokumen

Urutan baca untuk tim frontend:

1. `graphql.md` (utama, siap pakai dari FE)
2. Dokumen per service untuk detail domain

- `account.md`: endpoint gRPC Account Service
- `auth.md`: endpoint gRPC Auth Service
- `novel.md`: endpoint gRPC Novel Service
- `readinglist.md`: endpoint gRPC ReadingList Service
- `review.md`: endpoint gRPC Review Service
- `graphql.md`: endpoint Query/Mutation GraphQL Gateway

## Konvensi Umum

- Waktu pada protobuf disimpan dalam field `bytes` (contoh: `created_at`, `updated_at`) dan dipetakan ke scalar `Time` pada GraphQL.
- Operasi hapus umumnya mengembalikan pola response:
  - `deleted_id` atau `deletedID`
  - `success`
  - `message`
- Pagination menggunakan pola `skip` dan `take`.
