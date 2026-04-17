ALTER TABLE users
  ADD COLUMN IF NOT EXISTS username VARCHAR(50),
  ADD COLUMN IF NOT EXISTS streak INT NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS daily_progress INT NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS last_active TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  ADD COLUMN IF NOT EXISTS preferences JSONB NOT NULL DEFAULT '{}'::jsonb;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username_unique ON users (username);

CREATE TABLE IF NOT EXISTS user_progress (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
  module VARCHAR(30) NOT NULL,
  item_id VARCHAR(191) NOT NULL,
  progress INT NOT NULL DEFAULT 0,
  completed BOOLEAN NOT NULL DEFAULT FALSE,
  score INT NOT NULL DEFAULT 0,
  time_spent INT NOT NULL DEFAULT 0,
  last_accessed TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_user_progress_user_id ON user_progress (user_id);
CREATE INDEX IF NOT EXISTS idx_user_progress_module ON user_progress (module);
CREATE INDEX IF NOT EXISTS idx_user_progress_item_id ON user_progress (item_id);
CREATE INDEX IF NOT EXISTS idx_user_progress_last_accessed ON user_progress (last_accessed);
CREATE INDEX IF NOT EXISTS idx_user_progress_deleted_at ON user_progress (deleted_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_progress_user_module_item ON user_progress (user_id, module, item_id);

CREATE TABLE IF NOT EXISTS bookmarks (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
  type VARCHAR(20) NOT NULL,
  content_id VARCHAR(191) NOT NULL,
  note VARCHAR(500),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_bookmarks_user_id ON bookmarks (user_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_type ON bookmarks (type);
CREATE INDEX IF NOT EXISTS idx_bookmarks_content_id ON bookmarks (content_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_created_at ON bookmarks (created_at);
CREATE INDEX IF NOT EXISTS idx_bookmarks_deleted_at ON bookmarks (deleted_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_bookmarks_user_type_content ON bookmarks (user_id, type, content_id);

CREATE TABLE IF NOT EXISTS dhikr_counters (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
  dhikr_id VARCHAR(191) NOT NULL,
  count INT NOT NULL DEFAULT 0,
  target INT NOT NULL DEFAULT 33,
  date VARCHAR(10) NOT NULL,
  session VARCHAR(20) NOT NULL,
  completed BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_dhikr_counters_user_id ON dhikr_counters (user_id);
CREATE INDEX IF NOT EXISTS idx_dhikr_counters_dhikr_id ON dhikr_counters (dhikr_id);
CREATE INDEX IF NOT EXISTS idx_dhikr_counters_date ON dhikr_counters (date);
CREATE INDEX IF NOT EXISTS idx_dhikr_counters_session ON dhikr_counters (session);
CREATE INDEX IF NOT EXISTS idx_dhikr_counters_deleted_at ON dhikr_counters (deleted_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_dhikr_counters_user_dhikr_date_session ON dhikr_counters (user_id, dhikr_id, date, session);

CREATE TABLE IF NOT EXISTS quiz_attempts (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
  category VARCHAR(100) NOT NULL,
  score INT NOT NULL,
  total_questions INT NOT NULL,
  time_spent INT NOT NULL,
  answers JSONB NOT NULL DEFAULT '[]'::jsonb,
  completed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_quiz_attempts_user_id ON quiz_attempts (user_id);
CREATE INDEX IF NOT EXISTS idx_quiz_attempts_category ON quiz_attempts (category);
CREATE INDEX IF NOT EXISTS idx_quiz_attempts_completed_at ON quiz_attempts (completed_at);
CREATE INDEX IF NOT EXISTS idx_quiz_attempts_deleted_at ON quiz_attempts (deleted_at);
