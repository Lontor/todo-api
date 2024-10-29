CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    account_type VARCHAR(20) NOT NULL DEFAULT 'regular',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);