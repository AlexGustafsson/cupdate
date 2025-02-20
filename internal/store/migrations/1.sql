-- Source revision: 1
-- Target revision: 2
-- Summary: Implement OSSF Scorecards
CREATE TABLE images_scorecards (
  reference TEXT NOT NULL,
  score REAL NOT NULL,
  scorecard BLOB NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

ALTER TABLE images_changes ADD changedScorecard BOOLEAN NOT NULL DEFAULT FALSE;

-- Update changes on INSERT to images_scorecards table
CREATE TRIGGER images_changes_images_scorecards_insert AFTER INSERT ON images_scorecards BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedScorecard
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "insert",

    TRUE
  );
END;

-- Update images_updates on UPDATE to images_scorecards table
CREATE TRIGGER images_changes_images_scorecards_update AFTER UPDATE ON images_scorecards WHEN
    old.scorecard <> new.scorecard
  BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedScorecard
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "update",

    TRUE
  );
END;

INSERT INTO revision (id, revision) VALUES (0, 2) ON CONFLICT DO UPDATE SET revision=excluded.revision;
