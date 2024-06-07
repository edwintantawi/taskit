CREATE TABLE projects (
  id            VARCHAR(64)   PRIMARY KEY,
  user_id       VARCHAR(64)   NOT NULL,
  title         VARCHAR(255)  NOT NULL,
  created_at    TIMESTAMP     NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMP     NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_projects_user FOREIGN KEY(user_id) REFERENCES users(id)
);

ALTER TABLE tasks ADD project_id VARCHAR(64);

ALTER TABLE tasks
  ADD CONSTRAINT fk_tasks_projects
  FOREIGN KEY(project_id)
  REFERENCES projects(id);