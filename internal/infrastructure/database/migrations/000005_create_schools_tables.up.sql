CREATE TABLE IF NOT EXISTS schools (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
  name VARCHAR(191) NOT NULL,
  location VARCHAR(255) NOT NULL,
  monthly_fee INT NOT NULL DEFAULT 0,
  map_url VARCHAR(1024) NOT NULL,
  description VARCHAR(1000) NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_schools_user_id ON schools (user_id);
CREATE INDEX IF NOT EXISTS idx_schools_created_at ON schools (created_at);
CREATE INDEX IF NOT EXISTS idx_schools_deleted_at ON schools (deleted_at);
