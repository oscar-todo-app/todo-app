-- +goose Up
CREATE TABLE tasks (
  id SERIAL PRIMARY KEY,
  task_name VARCHAR(255) NOT NULL,
  status BOOLEAN
);

-- +goose Down
DROP TABLE tasks;
