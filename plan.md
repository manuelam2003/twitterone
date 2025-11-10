# Micro-Twitter API

A compact microblogging API for user accounts, posts (tweets), follows, likes, replies, media uploads, and realtime feeds. Designed as an MVP that can be extended later with notifications, search, and analytics.

## Table of Contents

* Installation
* Endpoints
* Authentication
* Middleware
* Configuration
* Example API Workflow
* Table schema

---

# Installation

Clone the repository:

```bash
git clone <repository_url>
```

Navigate to the project and install dependencies:

```bash
cd <project_directory>
go mod download
```

Set up the database and ensure environment variables are configured (example uses SQLite for dev and Postgres for production):

```bash
make db/reset
```

Run the server:

```bash
make run/api
```

---

# Endpoints

## Health Check

```
GET /v1/healthcheck
```

Checks if the API is running correctly.

---

## Auth

### Register

```
POST /v1/auth/register
```

Body:

```json
{
  "username": "bob",
  "email": "bob@example.com",
  "password": "secret"
}
```

### Login

```
POST /v1/auth/login
```

Body:

```json
{
  "email": "bob@example.com",
  "password": "secret"
}
```

Response returns a JWT (or session cookie) and user object.

### Logout

```
POST /v1/auth/logout
```

Invalidate JWT/session.

### Refresh Token (optional)

```
POST /v1/auth/refresh
```

Exchange refresh token for a new access token.

---

## Users

### Get own profile

```
GET /v1/me
```

Auth required.

### Get public user profile

```
GET /v1/users/{username}
```

Returns public profile: username, bio, avatar_url, counts.

### Update own profile

```
PATCH /v1/me
```

Body (any subset):

```json
{
  "bio": "hello",
  "avatar_url": "/media/abcd.jpg",
  "name": "Bob"
}
```

### Search users

```
GET /v1/users?query=alice&limit=20&offset=0
```

---

## Posts (Tweets)

### Create post

```
POST /v1/posts
```

Auth required. Body:

```json
{
  "content": "hello world",
  "reply_to_id": null,
  "media_ids": []
}
```

### Get global timeline

```
GET /v1/posts
```

Query params: `limit`, `offset` or `before_id` (keyset).

### Get user timeline (public posts)

```
GET /v1/users/{username}/posts
```

### Get post by id

```
GET /v1/posts/{post_id}
```

### Delete own post

```
DELETE /v1/posts/{post_id}
```

Auth required; only owner allowed.

### Edit own post (optional / limited)

```
PUT /v1/posts/{post_id}
```

Body:

```json
{
  "content": "edited text"
}
```

---

## Replies (Threads)

### Create reply

```
POST /v1/posts/{post_id}/replies
```

Body:

```json
{
  "content": "this is a reply"
}
```

### Get replies

```
GET /v1/posts/{post_id}/replies
```

---

## Likes

### Like a post

```
POST /v1/posts/{post_id}/like
```

### Unlike a post

```
DELETE /v1/posts/{post_id}/like
```

### Get likes for a post

```
GET /v1/posts/{post_id}/likes
```

---

## Follow System

### Follow a user

```
POST /v1/users/{username}/follow
```

### Unfollow a user

```
DELETE /v1/users/{username}/follow
```

### Get followers

```
GET /v1/users/{username}/followers
```

### Get following

```
GET /v1/users/{username}/following
```

---

## Media

### Upload media (images)

```
POST /v1/media
```

Multipart form upload. Returns:

```json
{ "id": "abcd", "url": "/media/abcd.jpg" }
```

### Get media

```
GET /v1/media/{media_id}
```

---

## Notifications

### List notifications

```
GET /v1/notifications
```

Auth required. Returns likes, replies, new followers, mentions.

### Mark notification read

```
POST /v1/notifications/{notification_id}/read
```

---

## Realtime / Streaming

### Server-Sent Events feed (personalized)

```
GET /v1/stream/feed
```

Auth required; stream new posts from followed users, notifications, etc.

### WebSocket (alternative)

```
WS /v1/ws/feed
```

Authenticate via query param or initial auth message.

---

## Admin (optional)

### List users (admin)

```
GET /v1/admin/users
```

### Delete user (admin)

```
DELETE /v1/admin/users/{user_id}
```

---

# Authentication

Authentication is handled via JWTs (access + refresh tokens) or server-side sessions. Include access token in Authorization header for protected routes:

```
Authorization: Bearer <access_token>
```

Tokens are issued from `/v1/auth/login` and refreshed at `/v1/auth/refresh` if implemented.

Password storage: hash with bcrypt (or Argon2). Validate inputs and rate-limit auth endpoints.

---

# Middleware

* **Panic Recovery**: Recover from panics and return 5xx with logs.
* **Logging**: Request/response logging (path, method, status, latency).
* **Rate Limiting**: Global and per-user rate limits (e.g., posting, liking).
* **Authentication**: Verify JWT/session for protected routes.
* **Require Verified/Activated User**: Optionally block actions for unverified accounts.
* **CSRF Protection**: If you support cookie-based auth for non-API clients.
* **CORS**: Configure allowed origins for web clients.
* **Input Validation**: Size limits (e.g., post content length), allowed media types.
* **Metrics**: Expose Prometheus-style metrics (`/metrics`) for requests, latencies, error rates.

---

# Configuration

Use a `.env` file or environment variables for configuration values:

* `APP_ENV` (development|production)
* `PORT` (e.g., 8080)
* `DATABASE_URL` (sqlite file or postgres dsn)
* `JWT_SECRET`
* `JWT_EXPIRY` (e.g., 15m)
* `REFRESH_TOKEN_EXPIRY` (optional)
* `MEDIA_DIRECTORY` or S3 credentials (`S3_BUCKET`, `S3_REGION`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
* `RATE_LIMIT_PER_MIN`
* `SENTRY_DSN` (optional for error tracking)

Ensure the environment variables are set for production and CI.

---

# Example API Workflow

1. **Sign up**: `POST /v1/auth/register` → returns basic user info.
2. **Login**: `POST /v1/auth/login` → receive access token.
3. **Create a post**: `POST /v1/posts` with `Authorization` header.
4. **View global timeline**: `GET /v1/posts`.
5. **Follow someone**: `POST /v1/users/{username}/follow`.
6. **See personalized feed (SSE/WS)**: connect to `/v1/stream/feed` after authenticating.
7. **Like a post**: `POST /v1/posts/{post_id}/like`.
8. **Reply**: `POST /v1/posts/{post_id}/replies`.
9. **Upload avatar/media**: `POST /v1/media`.
10. **Receive notifications**: `GET /v1/notifications` or via realtime stream.

---

# Table schema

## 1. Users table

* `id` (PK): bigint/uuid — user id.
* `username` (unique): string — handle.
* `email` (unique): string.
* `password_hash`: string.
* `name`: string (display name).
* `bio`: text.
* `avatar_url`: string.
* `is_active`: boolean.
* `is_verified`: boolean.
* `created_at`: timestamp.
* `updated_at`: timestamp.

## 2. Sessions / Tokens table (optional)

* `id` (PK)
* `user_id` (FK -> users.id)
* `refresh_token_hash`
* `expires_at`
* `created_at`
* `last_used_at`

## 3. Posts table

* `id` (PK): bigint/uuid.
* `user_id` (FK -> users.id): author.
* `content`: text (limit e.g., 280 chars).
* `reply_to_id` (FK -> posts.id, nullable): if reply.
* `created_at`: timestamp.
* `updated_at`: timestamp.
* `is_deleted`: boolean.

## 4. Media table

* `id` (PK)
* `uploader_id` (FK -> users.id)
* `url` or `storage_key`
* `mime_type`
* `size_bytes`
* `created_at`

## 5. PostMedia (many-to-many)

* `post_id` (FK -> posts.id)
* `media_id` (FK -> media.id)

## 6. Follows table

* `id` (PK)
* `follower_id` (FK -> users.id)
* `followed_id` (FK -> users.id)
* `created_at`

Unique index on (`follower_id`, `followed_id`).

## 7. Likes table

* `id` (PK)
* `user_id` (FK -> users.id)
* `post_id` (FK -> posts.id)
* `created_at`

Unique index on (`user_id`, `post_id`).

## 8. Replies table

* Note: replies are posts with `reply_to_id`. No separate replies table required. Optionally maintain a denormalized `reply_count` on posts.

## 9. Notifications table

* `id` (PK)
* `user_id` (FK -> users.id) — recipient
* `type`: enum (`like`, `reply`, `follow`, `mention`, etc.)
* `source_user_id` (FK -> users.id) — actor
* `post_id` (FK -> posts.id, nullable)
* `is_read`: boolean
* `created_at`

## 10. Indexes & counters (denormalized columns)

* `posts.reply_count`, `posts.like_count` (updated on write or via background job).
* `users.followers_count`, `users.following_count`.
* Indexes on `posts.created_at` for timeline queries; composite index for keyset pagination (`created_at DESC, id DESC`).

---

# Notes / Recommendations

* Start with SQLite + `html/template` for a quick local dev experience. Migrate to Postgres for production.
* Use keyset pagination for timelines (cursor = `created_at` + `id`) for performance.
* Denormalize counts to avoid expensive aggregate queries; update counts transactionally or with an async job.
* Keep post content limits and validate media sizes/types.
* Consider background workers for media processing (thumbnails), sending emails, and computing heavy aggregates.

---

You can tell me if you want this output turned into:

* a ready-to-use OpenAPI (Swagger) spec,
* SQL migration files for the schema, or
* a starter Go project skeleton (handlers, models, and a basic router).
