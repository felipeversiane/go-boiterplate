CREATE UNIQUE INDEX IF NOT EXISTS unique_email_active_users
ON users (email)
WHERE deleted = FALSE;