-- Revert the created_at column addition
ALTER TABLE products DROP COLUMN IF EXISTS created_at;


