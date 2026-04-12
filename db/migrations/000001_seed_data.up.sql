-- USERS (test-friendly names + emails)

INSERT INTO users (id, name, email, password, created_at, updated_at)
VALUES
(
    'df6fbabd-96a5-4988-a8e6-161f44d5160e',
    'test_user_1',
    'test1@taskflow.com',
    '$2a$12$yDmoWWWDt7.jp/QHuBarJeNn4E5.LpH/hxQUq1Xc9EuTglnLSS0v.',
    '2026-04-11 20:32:56.783522',
    '2026-04-11 20:32:56.783522'
),
(
    '9daa73c5-3853-414e-b02e-e1500d62b5b3',
    'test_user_2',
    'test2@taskflow.com',
    '$2a$12$.iqCQCwzbwo/3dzgj6JXYubXEUNYW2oABKqvhf3d3qSIFX1aoHspC',
    '2026-04-11 20:59:29.044443',
    '2026-04-11 20:59:29.044443'
),
(
    '1c27b4c9-be8f-4fac-ba93-d6439e0e050d',
    'test_user_3',
    'test3@taskflow.com',
    '$2a$12$vioh6n8tci4M4gMfFVKqgeHlVXhUjIiaEJ1HvQ8W/WBwE/CKUlQqq',
    '2026-04-12 10:26:43.387896',
    '2026-04-12 10:26:43.387896'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO projects (id, name, description, owner_id, created_at, updated_at)
VALUES
(
    '7970a8a7-eda7-4a26-91e3-6f5609a4de74',
    'test_project_1',
    'sample project',
    '9daa73c5-3853-414e-b02e-e1500d62b5b3',
    '2026-04-12 01:20:41.083979',
    '2026-04-12 01:20:41.083979'
),
(
    '022bfb02-9c31-4f6b-892a-6dc6a4c3b6f4',
    'test_project_2',
    'testing project',
    '9daa73c5-3853-414e-b02e-e1500d62b5b3',
    '2026-04-11 23:48:46.338885',
    '2026-04-12 11:09:13.749438'
),
(
    '01cd0839-4b95-458f-b5ac-3e9f7a1e7693',
    'test_project_3',
    'testing project',
    '1c27b4c9-be8f-4fac-ba93-d6439e0e050d',
    '2026-04-12 10:45:41.025193',
    '2026-04-12 11:41:20.732585'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO tasks (id, title, description, status, priority, project_id, assignee_id, created_at, updated_at)
VALUES
(
    '84a38fb2-8992-4c1e-a287-83c6abcdffe9',
    'test task 1',
    'login issue',
    'todo',
    'high',
    '022bfb02-9c31-4f6b-892a-6dc6a4c3b6f4',
    NULL,
    '2026-04-12 00:09:05.842036',
    '2026-04-12 00:09:05.842036'
),
(
    '85c02a0e-d6e1-448f-8bb4-070b93167741',
    'test task 2',
    'auth issue',
    'in_progress',
    'high',
    '022bfb02-9c31-4f6b-892a-6dc6a4c3b6f4',
    NULL,
    '2026-04-12 00:36:58.16439',
    '2026-04-12 00:36:58.16439'
),
(
    'dd76899f-b08f-48b2-9694-6dec4a3e3853',
    'test task 3',
    '',
    'done',
    'high',
    '01cd0839-4b95-458f-b5ac-3e9f7a1e7693',
    NULL,
    '2026-04-12 11:56:50.203322',
    '2026-04-12 11:56:50.203322'
)
ON CONFLICT (id) DO NOTHING;
