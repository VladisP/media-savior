BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users
(
    id          UUID         NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    external_id INT          NOT NULL UNIQUE,
    created_at  TIMESTAMP    NOT NULL DEFAULT now()
);

CREATE TYPE media_type AS ENUM ('story', 'audio');

CREATE TABLE IF NOT EXISTS media
(
    id         UUID         NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    type       media_type   NOT NULL,
    user_id    UUID         NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP    NOT NULL DEFAULT now(),

    UNIQUE (name, user_id)
);

COMMIT;
