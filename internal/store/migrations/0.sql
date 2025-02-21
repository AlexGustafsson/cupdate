-- Source revision: 0
-- Target revision: 1
-- Summary: Start tracking revisions
CREATE TABLE revision (
  -- Ensure only a single row can exist
  id INTEGER PRIMARY KEY CHECK (id = 0),
  revision INT NOT NULL
);
INSERT INTO revision (id, revision) VALUES (0, 1);
