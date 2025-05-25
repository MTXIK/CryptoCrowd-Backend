-- +goose Up
-- +goose StatementBegin
CREATE TYPE RoleType AS ENUM ('admin', 'startup', 'investor');

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role RoleType NOT NULL DEFAULT 'investor',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT users_email_role_key UNIQUE (email, role)
);

CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_users_email_role ON users (email, role);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS RoleType;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email_role;
-- +goose StatementEnd
