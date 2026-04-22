CREATE TABLE IF NOT EXISTS habits (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
  name VARCHAR(191) NOT NULL,
  category VARCHAR(50) NOT NULL DEFAULT 'Other',
  reminder_time VARCHAR(5) NOT NULL DEFAULT '',
  reminder_enabled BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_habits_user_id ON habits (user_id);
CREATE INDEX IF NOT EXISTS idx_habits_created_at ON habits (created_at);
CREATE INDEX IF NOT EXISTS idx_habits_deleted_at ON habits (deleted_at);

CREATE TABLE IF NOT EXISTS habit_completions (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
  habit_id UUID NOT NULL REFERENCES habits(id) ON UPDATE CASCADE ON DELETE CASCADE,
  date VARCHAR(10) NOT NULL,
  completed BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_habit_completions_user_id ON habit_completions (user_id);
CREATE INDEX IF NOT EXISTS idx_habit_completions_habit_id ON habit_completions (habit_id);
CREATE INDEX IF NOT EXISTS idx_habit_completions_date ON habit_completions (date);
CREATE INDEX IF NOT EXISTS idx_habit_completions_deleted_at ON habit_completions (deleted_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_habit_completions_user_habit_date ON habit_completions (user_id, habit_id, date);
