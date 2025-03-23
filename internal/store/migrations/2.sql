-- Source revision: 2
-- Target revision: 3
-- Summary: Implement provenance
CREATE TABLE images_provenance (
  reference TEXT NOT NULL,
  provenance BLOB NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

ALTER TABLE images_changes ADD changedProvenance BOOLEAN NOT NULL DEFAULT FALSE;

-- Update changes on INSERT to images_provenance table
CREATE TRIGGER images_changes_images_provenance_insert AFTER INSERT ON images_provenance BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedProvenance
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "insert",

    TRUE
  );
END;

-- Update images_updates on UPDATE to images_provenance table
CREATE TRIGGER images_changes_images_provenance_update AFTER UPDATE ON images_provenance WHEN
    old.provenance <> new.provenance
  BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedProvenance
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "update",

    TRUE
  );
END;

INSERT INTO revision (id, revision) VALUES (0, 3) ON CONFLICT DO UPDATE SET revision=excluded.revision;
