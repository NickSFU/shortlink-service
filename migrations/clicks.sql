CREATE TABLE clicks (
    id SERIAL PRIMARY KEY,

    short_link_id INTEGER NOT NULL REFERENCES short_links(id),

    ip TEXT,
    user_agent TEXT,
    referer TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);