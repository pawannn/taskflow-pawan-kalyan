-- Seed data: 1 user with known credentials + extra users for testing
-- Password for all users: Test@Pass1!

INSERT INTO users (id, name, email, password, created_at, updated_at)
VALUES
(
    'df6fbabd-96a5-4988-a8e6-161f44d5160e',
    'Test User',
    'test@example.com',
    '$2a$12$HAz3392t0FUw11PsGSBIueuDcrTdZCdw2/Fox3ECmxADiCJKslD7W',
    NOW(),
    NOW()
),
(
    '9daa73c5-3853-414e-b02e-e1500d62b5b3',
    'Alice Johnson',
    'alice@example.com',
    '$2a$12$HAz3392t0FUw11PsGSBIueuDcrTdZCdw2/Fox3ECmxADiCJKslD7W',
    NOW(),
    NOW()
),
(
    '1c27b4c9-be8f-4fac-ba93-d6439e0e050d',
    'Bob Smith',
    'bob@example.com',
    '$2a$12$HAz3392t0FUw11PsGSBIueuDcrTdZCdw2/Fox3ECmxADiCJKslD7W',
    NOW(),
    NOW()
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO projects (id, name, description, owner_id, created_at, updated_at)
VALUES
(
    '7970a8a7-eda7-4a26-91e3-6f5609a4de74',
    'Website Redesign',
    'Q2 redesign initiative',
    'df6fbabd-96a5-4988-a8e6-161f44d5160e',
    NOW(),
    NOW()
),
(
    '022bfb02-9c31-4f6b-892a-6dc6a4c3b6f4',
    'Mobile App',
    'iOS and Android app development',
    '9daa73c5-3853-414e-b02e-e1500d62b5b3',
    NOW(),
    NOW()
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO tasks (id, title, description, status, priority, project_id, assignee_id, created_at, updated_at)
VALUES
(
    '84a38fb2-8992-4c1e-a287-83c6abcdffe9',
    'Design landing page',
    'Create wireframes and mockups for the new homepage',
    'todo',
    'high',
    '7970a8a7-eda7-4a26-91e3-6f5609a4de74',
    'df6fbabd-96a5-4988-a8e6-161f44d5160e',
    NOW(),
    NOW()
),
(
    '85c02a0e-d6e1-448f-8bb4-070b93167741',
    'Set up CI pipeline',
    'Configure GitHub Actions for automated testing and deployment',
    'in_progress',
    'medium',
    '7970a8a7-eda7-4a26-91e3-6f5609a4de74',
    '9daa73c5-3853-414e-b02e-e1500d62b5b3',
    NOW(),
    NOW()
),
(
    'dd76899f-b08f-48b2-9694-6dec4a3e3853',
    'Write API documentation',
    'Document all endpoints with request/response examples',
    'done',
    'low',
    '7970a8a7-eda7-4a26-91e3-6f5609a4de74',
    NULL,
    NOW(),
    NOW()
)
ON CONFLICT (id) DO NOTHING;
