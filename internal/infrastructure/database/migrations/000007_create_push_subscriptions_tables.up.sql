CREATE TABLE IF NOT EXISTS push_subscriptions (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
  endpoint VARCHAR(2048) NOT NULL UNIQUE,
  p256dh VARCHAR(512) NOT NULL,
  auth VARCHAR(255) NOT NULL,
  expiration_time BIGINT,
  timezone VARCHAR(100) NOT NULL DEFAULT 'UTC',
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_push_subscriptions_user_id ON push_subscriptions (user_id);
CREATE INDEX IF NOT EXISTS idx_push_subscriptions_timezone ON push_subscriptions (timezone);
CREATE INDEX IF NOT EXISTS idx_push_subscriptions_enabled ON push_subscriptions (enabled);
CREATE INDEX IF NOT EXISTS idx_push_subscriptions_deleted_at ON push_subscriptions (deleted_at);
