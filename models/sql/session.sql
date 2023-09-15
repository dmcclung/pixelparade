DROP TABLE IF EXISTS sessions;

CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id varchar(255) UNIQUE NOT NULL,
    token_hash TEXT UNIQUE NOT NULL
);
