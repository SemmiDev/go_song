CREATE TABLE IF NOT EXISTS accounts
(
    id                TEXT         NOT NULL PRIMARY KEY,
    name              VARCHAR(255) NOT NULL,
    email             VARCHAR(255) NOT NULL UNIQUE,
    password          VARCHAR(255) NOT NULL,
    is_email_verified BOOLEAN      NOT NULL DEFAULT 'FALSE',
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX ON accounts (name);