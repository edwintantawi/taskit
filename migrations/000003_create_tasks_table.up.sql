CREATE TABLE tasks (
  id            VARCHAR(64)   PRIMARY KEY,
  user_id       VARCHAR(64)   NOT NULL,
  content       VARCHAR(255)  NOT NULL,
  description   TEXT          NOT NULL,
  is_done       BOOLEAN       NOT NULL DEFAULT FALSE,
  due_date      TIMESTAMP,
  created_at    TIMESTAMP     NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMP     NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_tasks_users FOREIGN KEY(user_id) REFERENCES users(id)
)