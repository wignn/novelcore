CREATE TABLE IF NOT EXISTS reading_lists (
  id CHAR(27) PRIMARY KEY,
  account_id CHAR(27) NOT NULL,
  novel_id CHAR(27) NOT NULL,
  status VARCHAR(32) DEFAULT 'plan_to_read',
  current_chapter DECIMAL(10,1) DEFAULT 0,
  rating INT CHECK (rating >= 0 AND rating <= 10),
  notes TEXT DEFAULT '',
  is_favorite BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(account_id, novel_id)
);

CREATE INDEX idx_reading_list_account ON reading_lists(account_id);
CREATE INDEX idx_reading_list_novel ON reading_lists(novel_id);
CREATE INDEX idx_reading_list_status ON reading_lists(account_id, status);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_reading_list_updated_at BEFORE UPDATE ON reading_lists
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
