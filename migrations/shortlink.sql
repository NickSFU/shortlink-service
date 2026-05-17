CREATE TABLE short_links (
    id SERIAL PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER NULL,
    is_deleted BOOLEAN NULL
);
ALTER TABLE short_links
ALTER COLUMN is_deleted SET DEFAULT false;

