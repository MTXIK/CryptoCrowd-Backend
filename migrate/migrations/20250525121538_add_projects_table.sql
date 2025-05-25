-- +goose Up
-- +goose StatementBegin
CREATE TYPE StatusType AS ENUM ('pending', 'approved', 'rejected');

CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    ownerID INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status StatusType NOT NULL DEFAULT 'pending',
    name VARCHAR(255) NOT NULL,
    description TEXT,
    amountRequested NUMERIC(18,2) NOT NULL DEFAULT 0,
    amountRaised NUMERIC(18,2) NOT NULL DEFAULT 0,
    deadlineAt TIMESTAMP WITH TIME ZONE NOT NULL,
    createdAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS project_photos (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_projects_ownerID ON projects (ownerID);
CREATE INDEX idx_projects_status ON projects (status);
CREATE INDEX idx_projects_deadlineAt ON projects (deadlineAt);
CREATE INDEX idx_project_photos_project_id ON project_photos (project_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS project_photos;
DROP TABLE IF EXISTS projects;
DROP TYPE IF EXISTS StatusType;
DROP INDEX IF EXISTS idx_projects_ownerID;
DROP INDEX IF EXISTS idx_projects_status;
DROP INDEX IF EXISTS idx_projects_deadlineAt;
DROP INDEX IF EXISTS idx_project_photos_project_id;
-- +goose StatementEnd
