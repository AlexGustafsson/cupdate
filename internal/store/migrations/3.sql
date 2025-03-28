-- Source revision: 3
-- Target revision: 4
-- Summary: Implement sboms
CREATE TABLE images_sbom (
  reference TEXT NOT NULL,
  sbom BLOB NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

ALTER TABLE images_changes ADD changedSBOM BOOLEAN NOT NULL DEFAULT FALSE;

-- Update changes on INSERT to images_sbom table
CREATE TRIGGER images_changes_images_sbom_insert AFTER INSERT ON images_sbom BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedSBOM
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "insert",

    TRUE
  );
END;

-- Update images_updates on UPDATE to images_sbom table
CREATE TRIGGER images_changes_images_sbom_update AFTER UPDATE ON images_sbom WHEN
    old.sbom <> new.sbom
  BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedSBOM
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "update",

    TRUE
  );
END;

INSERT INTO revision (id, revision) VALUES (0, 4) ON CONFLICT DO UPDATE SET revision=excluded.revision;
