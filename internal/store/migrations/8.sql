-- Source revision: 8
-- Target revision: 9
-- Summary: Track changes to new versions

-- TODO: Remove in v1

CREATE TABLE images_updates (
  newReference TEXT PRIMARY KEY NOT NULL,
  newAnnotations BLOB,
  oldReference TEXT NOT NULL,
  oldAnnotations BLOB,
  versionDiffSortable INT NOT NULL,
  identified DATETIME NOT NULL,
  released DATETIME
);

-- Create an entry on INSERT with a newer version already identified
CREATE TRIGGER images_changes_images_insert_version AFTER INSERT ON images WHEN
    new.latestReference <> ""
    AND new.latestReference <> new.reference
  BEGIN
  INSERT INTO images_updates(
    newReference,
    newAnnotations,
    oldReference,
    oldAnnotations,
    versionDiffSortable,
    identified,
    released
  ) VALUES (
    new.latestReference,
    new.latestAnnotations,
    new.reference,
    new.annotations,
    new.versionDiffSortable,
    datetime('now', 'subsecond'),
    new.latestCreated
  );
END;

-- Create an entry on UPDATE with a newer version than what was known before
CREATE TRIGGER images_changes_images_update_version AFTER UPDATE ON images WHEN
    new.latestReference <> ""
    AND new.latestReference <> old.latestReference
    AND new.latestReference <> new.reference
  BEGIN
  INSERT INTO images_updates(
    newReference,
    newAnnotations,
    oldReference,
    oldAnnotations,
    versionDiffSortable,
    identified,
    released
  ) VALUES (
    new.latestReference,
    new.latestAnnotations,
    new.reference,
    new.annotations,
    new.versionDiffSortable,
    datetime('now', 'subsecond'),
    new.latestCreated
  );
END;

INSERT INTO revision (id, revision) VALUES (0, 9) ON CONFLICT DO UPDATE SET revision=excluded.revision;
