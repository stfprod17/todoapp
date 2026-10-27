DROP INDEX IF EXISTS todoapp.users_email_unique_idx;

ALTER TABLE todoapp.users
    DROP COLUMN IF EXISTS password_hash,
    DROP COLUMN IF EXISTS email;