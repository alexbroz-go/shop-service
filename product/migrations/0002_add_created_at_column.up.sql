-- Add a lightweight nullable column to avoid heavy locks
ALTER TABLE products ADD COLUMN created_at TIMESTAMPTZ NULL;


