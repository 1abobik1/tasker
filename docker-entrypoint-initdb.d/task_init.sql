CREATE TYPE task_status AS ENUM ('pending', 'processing', 'completed', 'failed');

CREATE TABLE IF NOT EXISTS tasks (
  id TEXT PRIMARY KEY,
  payload BYTEA NOT NULL,
  status task_status NOT NULL,
  result BYTEA,
  error TEXT,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);