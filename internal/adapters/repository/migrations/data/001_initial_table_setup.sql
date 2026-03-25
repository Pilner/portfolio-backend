CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE users (
    id UUID PRIMARY KEY,
    email CITEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_info(
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    display_name VARCHAR(50) DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE post (
    id BIGSERIAL PRIMARY KEY,
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    body TEXT NOT NULL,
    excerpt TEXT NULL,
    cover_image_url TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_post_author_id ON post(author_id);
CREATE INDEX idx_post_created_at ON post(created_at DESC);

CREATE TABLE post_image (
    id UUID PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES post(id) ON DELETE CASCADE,
    s3_key TEXT NOT NULL,
    s3_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_post_image_post_id ON post_image(post_id);

---- create above / drop below ----

DROP INDEX IF EXISTS idx_post_image_post_id;
DROP TABLE IF EXISTS post_image;

DROP INDEX IF EXISTS idx_post_created_at;
DROP INDEX IF EXISTS idx_post_author_id;
DROP TABLE IF EXISTS post;

DROP TABLE IF EXISTS user_info;

DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS citext;