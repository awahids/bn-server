DROP TABLE IF EXISTS quiz_attempts;
DROP TABLE IF EXISTS dhikr_counters;
DROP TABLE IF EXISTS bookmarks;
DROP TABLE IF EXISTS user_progress;

DROP INDEX IF EXISTS idx_users_username_unique;

ALTER TABLE users
  DROP COLUMN IF EXISTS preferences,
  DROP COLUMN IF EXISTS last_active,
  DROP COLUMN IF EXISTS daily_progress,
  DROP COLUMN IF EXISTS streak,
  DROP COLUMN IF EXISTS username;
