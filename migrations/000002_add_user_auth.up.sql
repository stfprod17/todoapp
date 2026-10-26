ALTER TABLE todoapp.users
    ADD COLUMN email VARCHAR(255),
    ADD COLUMN password_hash VARCHAR(255);

CREATE UNIQUE INDEX users_email_unique_idx
    ON todoapp.users (LOWER(email))
    WHERE email IS NOT NULL;    