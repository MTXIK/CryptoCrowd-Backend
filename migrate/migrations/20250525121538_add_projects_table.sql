-- +goose Up
-- +goose StatementBegin
CREATE TYPE StatusType AS ENUM ('pending', 'approved', 'rejected');

CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status StatusType NOT NULL DEFAULT 'pending',
    name VARCHAR(255) NOT NULL,
    description TEXT,
    amount_requested NUMERIC(36,18) NOT NULL DEFAULT 0,
    amount_raised NUMERIC(36,18) NOT NULL DEFAULT 0,
    deadline_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS project_images (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_projects_owner_id ON projects (owner_id);
CREATE INDEX idx_projects_status ON projects (status);
CREATE INDEX idx_projects_deadline_at ON projects (deadline_at);
CREATE INDEX idx_project_photos_project_id ON project_images (project_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS project_images;
DROP TABLE IF EXISTS projects;
DROP TYPE IF EXISTS StatusType;
DROP INDEX IF EXISTS idx_projects_owner_id;
DROP INDEX IF EXISTS idx_projects_status;
DROP INDEX IF EXISTS idx_projects_deadline_at;
DROP INDEX IF EXISTS idx_project_photos_project_id;
-- +goose StatementEnd
