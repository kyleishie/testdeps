ALTER TABLE IF EXISTS users
ADD COLUMN IF NOT EXISTS name VARCHAR (50) NOT NULL DEFAULT('unknown');