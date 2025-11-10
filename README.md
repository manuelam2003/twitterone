# Micro-Twitter API

A compact microblogging API for users, posts, follows, likes, replies, media uploads and realtime feeds. MVP-first design, extendable.

## Table of Contents

Installation · Endpoints · Authentication · Middleware · Configuration · Example workflow · Table schema

---

## Installation

```bash
git clone <repo_url>
cd <project_dir>
go mod download
make db/reset   # seeds dev DB (sqlite)
make run/api
```

---

## Endpoints (compact)

**Health**

* `GET /v1/healthcheck`

**Auth**

* `POST /v1/auth/register`
* `POST /v1/auth/login`
* `POST /v1/auth/logout`
* `POST /v1/auth/refresh` (optional)

**Users**

* `GET /v1/me`
* `PATCH /v1/me`
* `GET /v1/users?query=&limit=&offset=`
* `GET /v1/users/{username}`

**Posts**

* `POST /v1/posts`
* `GET /v1/posts` (global timeline; `limit`, `before_id`/cursor)
* `GET /v1/posts/{post_id}`
* `GET /v1/users/{username}/posts`
* `PUT /v1/posts/{post_id}` (edit, optional)
* `DELETE /v1/posts/{post_id}`

**Replies**

* `POST /v1/posts/{post_id}/replies`
* `GET /v1/posts/{post_id}/replies`

**Likes**

* `POST /v1/posts/{post_id}/like`
* `DELETE /v1/posts/{post_id}/like`
* `GET /v1/posts/{post_id}/likes`

**Follows**

* `POST /v1/users/{username}/follow`
* `DELETE /v1/users/{username}/follow`
* `GET /v1/users/{username}/followers`
* `GET /v1/users/{username}/following`

**Media**

* `POST /v1/media` (multipart upload)
* `GET /v1/media/{media_id}`

**Notifications**

* `GET /v1/notifications`
* `POST /v1/notifications/{id}/read`

**Realtime**

* `GET /v1/stream/feed` (SSE)
* `WS  /v1/ws/feed` (WebSocket)

**Admin (optional)**

* `GET /v1/admin/users`
* `DELETE /v1/admin/users/{user_id}`

---

## Authentication

* JWT (access + optional refresh) or server sessions.
* Protected routes require header: `Authorization: Bearer <token>`.
* Passwords hashed (bcrypt/argon2). Rate-limit auth endpoints.

---

## Middleware

* Panic recovery, logging, auth, rate limiting, CORS, CSRF (if cookie-based), input validation, metrics (`/metrics`).

---

## Configuration (env)

* `APP_ENV`, `PORT`
* `DATABASE_URL`
* `JWT_SECRET`, `JWT_EXPIRY`, `REFRESH_TOKEN_EXPIRY`
* `MEDIA_DIR` or S3 (`S3_BUCKET`, `AWS_*`)
* `RATE_LIMIT_PER_MIN`, `SENTRY_DSN` (optional)

---

## Example API Workflow (short)

1. `POST /v1/auth/register` → create user.
2. `POST /v1/auth/login` → get token.
3. `POST /v1/posts` (auth) → create tweet.
4. `GET /v1/posts` → global timeline.
5. `POST /v1/users/{username}/follow` → follow user.
6. Connect to `GET /v1/stream/feed` for realtime updates.
7. `POST /v1/posts/{post_id}/like` / `POST /v1/posts/{post_id}/replies`.

---

## Table schema (compact)

**users**

* `id` PK, `username` (unique), `email` (unique), `password_hash`, `name`, `bio`, `avatar_url`, `is_active`, `is_verified`, `created_at`, `updated_at`

**sessions / tokens** (optional)

* `id`, `user_id` FK, `refresh_token_hash`, `expires_at`, `last_used_at`

**posts**

* `id` PK, `user_id` FK, `content` (e.g., 280 char limit), `reply_to_id` FK nullable, `created_at`, `updated_at`, `is_deleted`, `reply_count`, `like_count`

**media**

* `id` PK, `uploader_id` FK, `storage_key`/`url`, `mime_type`, `size_bytes`, `created_at`

**post_media**

* `post_id` FK, `media_id` FK

**follows**

* `id` PK, `follower_id`, `followed_id`, `created_at`
  Unique index (`follower_id`, `followed_id`)

**likes**

* `id` PK, `user_id`, `post_id`, `created_at`
  Unique index (`user_id`, `post_id`)

**notifications**

* `id` PK, `user_id` (recipient), `type` enum, `source_user_id`, `post_id` nullable, `is_read`, `created_at`

**indexes / tips**

* Keyset pagination index: `(created_at DESC, id DESC)` on posts. Denormalize counts for followers/likes/posts for fast reads.

---

If you want this even shorter (one-page cheatsheet), an OpenAPI spec, SQL migrations, or a starter Go skeleton, I can produce it next. Which one?
