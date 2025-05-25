-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS investments (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    amount NUMERIC(36,18) NOT NULL,
    invested_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_investments_user_id ON investments (user_id);
CREATE INDEX idx_investments_project_id ON investments (project_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS investments;
DROP INDEX IF EXISTS idx_investments_user_id;
DROP INDEX IF EXISTS idx_investments_project_id;
-- +goose StatementEnd