CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE task_status AS ENUM ('todo', 'in_progress', 'done');
CREATE TYPE task_priority AS ENUM ('low', 'medium', 'high');

CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR NOT NULL,
    description TEXT,
    status task_status NOT NULL DEFAULT 'todo',
    priority task_priority NOT NULL DEFAULT 'medium',
    project_id UUID NOT NULL,
    assignee_id UUID,
    due_date DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT fk_project
        FOREIGN KEY(project_id)
        REFERENCES projects(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_assignee
        FOREIGN KEY(assignee_id)
        REFERENCES users(id)
        ON DELETE SET NULL
);

CREATE INDEX idx_tasks_project_id ON tasks(project_id);
CREATE INDEX idx_tasks_assignee_id ON tasks(assignee_id);
