CREATE TABLE users
(
    id            UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    email         TEXT UNIQUE NOT NULL,
    password_hash TEXT        NOT NULL,
    role          TEXT        NOT NULL DEFAULT 'user',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ          DEFAULT NULL
);

CREATE TABLE user_details
(
    user_id        UUID PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    first_name     TEXT,
    last_name      TEXT,
    date_of_birth  DATE DEFAULT NULL,
    phone          TEXT,
    address        JSONB,
    photo_path     TEXT,
    thumbnail_path TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at     TIMESTAMPTZ          DEFAULT NULL
);

CREATE TABLE sessions
(
    id         UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    user_id    UUID REFERENCES users (id) ON DELETE CASCADE,
    token      TEXT        NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
