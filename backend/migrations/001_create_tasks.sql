CREATE TABLE IF NOT EXISTS tasks (
  id INT AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  status ENUM('todo','in_progress','done') NOT NULL DEFAULT 'todo',
  assignee VARCHAR(100),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,            -- soft delete
  UNIQUE KEY uq_title (title),          -- duplicate title -> 409
  INDEX idx_status (status),
  INDEX idx_assignee (assignee)
);

-- Initial seed data for assessment testing
INSERT INTO tasks (title, description, status, assignee) VALUES
  ('Setup Docker environment', 'Initialize MySQL and Redis containers', 'done', 'devops'),
  ('Implement Backend Filter & Search', 'Add status, keyword, assignee, pagination, sorting', 'done', 'backend_lead'),
  ('Implement Redis Caching', 'Cache GET /api/tasks with 60s TTL and auto invalidation', 'in_progress', 'backend_lead'),
  ('Build React Native UI', 'Implement search input, status filter, modal, and pagination', 'in_progress', 'frontend_lead'),
  ('Write Comprehensive Unit Tests', 'Cover update, search, cache invalidation, and component tests', 'todo', 'qa_engineer')
ON DUPLICATE KEY UPDATE title=VALUES(title);
