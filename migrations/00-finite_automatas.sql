CREATE SCHEMA api;

CREATE TABLE api.finite_automatas (
    id UUID PRIMARY KEY,
    description TEXT,
    tuple JSONB NOT NULL,
    render TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
