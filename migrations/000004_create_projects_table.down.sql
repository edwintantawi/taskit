ALTER TABLE tasks
DROP CONSTRAINT fk_tasks_projects;

ALTER TABLE tasks DROP COLUMN project_id;

DROP TABLE projects;