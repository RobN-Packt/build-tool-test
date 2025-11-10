CREATE TABLE IF NOT EXISTS books (
    id UUID PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    author VARCHAR(200) NOT NULL,
    price NUMERIC(10,2) NOT NULL CHECK (price >= 0),
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_books_created_at ON books (created_at DESC);
