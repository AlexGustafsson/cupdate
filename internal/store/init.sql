-- DO NOT EDIT.
--
-- This init file is not representative of the actual schema used by the store.
-- Instead, it was the original schema file before the support for migrations
-- were added. Since then, all edits to the database are additive and available
-- in the migrations directory.
--
-- In order for us to exercise the migration code and simplify testing, we
-- always use the migrations to update the database, which ensures that they
-- work no matter which revision the user comes from.

CREATE TABLE IF NOT EXISTS raw_images (
  reference TEXT PRIMARY KEY NOT NULL,
  tags BLOB NOT NULL,
  graph BLOB NOT NULL,
  lastProcessed DATETIME
);

CREATE TABLE IF NOT EXISTS images (
  reference TEXT PRIMARY KEY NOT NULL,
  created DATETIME,
  latestReference TEXT,
  latestCreated DATETIME,
  versionDiffSortable INT NOT NULL,
  description TEXT NOT NULL,
  lastModified DATETIME NOT NULL,
  imageUrl TEXT NOT NULL,
  FOREIGN KEY(reference) REFERENCES raw_images(reference) ON DELETE CASCADE
);

CREATE VIRTUAL TABLE IF NOT EXISTS images_fts USING FTS5(
  content='images',
  reference,
  description
);

DROP TRIGGER IF EXISTS images_fts_insert;
CREATE TRIGGER images_fts_insert AFTER INSERT ON images BEGIN
  INSERT INTO images_fts(rowid, reference, description) VALUES (new.rowid, new.reference, new.description);
END;

DROP TRIGGER IF EXISTS images_fts_delete;
CREATE TRIGGER images_fts_delete AFTER DELETE ON images BEGIN
  INSERT INTO images_fts(images_fts, rowid, reference, description) VALUES('delete', old.rowid, old.reference, old.description);
END;

DROP TRIGGER IF EXISTS images_fts_update;
CREATE TRIGGER images_fts_update AFTER UPDATE ON images BEGIN
  INSERT INTO images_fts(images_fts, rowid, reference, description) VALUES('delete', old.rowid, old.reference, old.description);
  INSERT INTO images_fts(rowid, reference, description) VALUES (new.rowid, new.reference, new.description);
END;

CREATE TABLE IF NOT EXISTS images_tags (
  reference TEXT NOT NULL,
  tag TEXT NOT NULL,
  PRIMARY KEY (reference, tag),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);
-- TODO: Remove in v1.
-- There was previously an issue with a missing cascade on images_tags, which
-- could cause images_tags to have rows not related to an entry in images, which
-- would result in counts like "vulnerable images" or "failed images" being off.
-- To remediate this issue, remove all entries without references.
DELETE FROM images_tags WHERE reference NOT IN (SELECT reference FROM images);

-- TODO: Rename in v1. This was done as an easy way to migrate somewhat
-- gracefully without having to drop the entire database
DROP TABLE IF EXISTS images_links;
CREATE TABLE IF NOT EXISTS images_linksv2 (
  reference TEXT NOT NULL,
  links BLOB NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS images_release_notes (
  reference TEXT NOT NULL,
  title TEXT NOT NULL,
  html TEXT NOT NULL,
  markdown TEXT NOT NULL,
  released DATETIME NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS images_descriptions (
  reference TEXT NOT NULL,
  html TEXT NOT NULL,
  markdown TEXT NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS images_graphs (
  reference TEXT NOT NULL,
  graph BLOB NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

-- TODO: Rename in v1. This was done as an easy way to migrate somewhat
-- gracefully without having to drop the entire database
DROP TABLE IF EXISTS images_vulnerabilities;
CREATE TABLE IF NOT EXISTS images_vulnerabilitiesv2 (
  reference TEXT NOT NULL,
  count INT NOT NULL,
  vulnerabilities BLOB NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS images_workflow_runs (
  reference TEXT NOT NULL,
  started DATETIME NOT NULL,
  result TEXT NOT NULL,
  blob BLOB NOT NULL,
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

-- Track changes
CREATE TABLE IF NOT EXISTS images_changes (
  reference TEXT NOT NULL,
  time DATETIME NOT NULL,
  type TEXT NOT NULL,

  changedBasic BOOLEAN NOT NULL DEFAULT FALSE,
  changedLinks BOOLEAN NOT NULL DEFAULT FALSE,
  changedReleaseNotes BOOLEAN NOT NULL DEFAULT FALSE,
  changedDescription BOOLEAN NOT NULL DEFAULT FALSE,
  changedGraph BOOLEAN NOT NULL DEFAULT FALSE,
  changedVulnerabilities BOOLEAN NOT NULL DEFAULT FALSE,

  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

-- NOTE: There are no ON DELETE triggers as they would insert data on a cascade
-- delete, at which point we already remove the images_changes entry

-- Update changes on INSERT to images table
DROP TRIGGER IF EXISTS images_changes_images_insert;
CREATE TRIGGER images_changes_images_insert AFTER INSERT ON images BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedBasic
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "insert",

    TRUE
  );
END;

-- Update images_updates on UPDATE to images table
DROP TRIGGER IF EXISTS images_changes_images_update;
CREATE TRIGGER images_changes_images_update AFTER UPDATE ON images WHEN
    old.latestReference <> new.latestReference
    OR old.description <> new.description
    OR old.imageUrl <> new.imageUrl
  BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedBasic
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "update",

    TRUE
  );
END;

-- Update changes on INSERT to images_links table
DROP TRIGGER IF EXISTS images_changes_images_links_insert;
CREATE TRIGGER images_changes_images_links_insert AFTER INSERT ON images_linksv2 BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedLinks
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "insert",

    TRUE
  );
END;

-- Update images_updates on UPDATE to images_links table
DROP TRIGGER IF EXISTS images_changes_images_links_update;
CREATE TRIGGER images_changes_images_links_update AFTER UPDATE ON images_linksv2 WHEN
    old.links <> new.links
  BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedLinks
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "update",

    TRUE
  );
END;

-- Update changes on INSERT to images_release_notes table
DROP TRIGGER IF EXISTS images_changes_images_release_notes_insert;
CREATE TRIGGER images_changes_images_release_notes_insert AFTER INSERT ON images_release_notes BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedReleaseNotes
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "insert",

    TRUE
  );
END;

-- Update images_updates on UPDATE to images_release_notes table
DROP TRIGGER IF EXISTS images_changes_images_release_notes_update;
CREATE TRIGGER images_changes_images_release_notes_update AFTER UPDATE ON images_release_notes WHEN
    old.title <> new.title
    OR old.html <> new.html
    OR old.markdown <> new.markdown
    OR old.released <> new.released
  BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedReleaseNotes
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "update",

    TRUE
  );
END;

-- Update changes on INSERT to images_descriptions table
DROP TRIGGER IF EXISTS images_changes_images_descriptions_insert;
CREATE TRIGGER images_changes_images_descriptions_insert AFTER INSERT ON images_descriptions BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedDescription
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "insert",

    TRUE
  );
END;

-- Update images_updates on UPDATE to images_descriptions table
DROP TRIGGER IF EXISTS images_changes_images_descriptions_update;
CREATE TRIGGER images_changes_images_descriptions_update AFTER UPDATE ON images_descriptions WHEN
    old.html <> new.html
    OR old.markdown <> new.markdown
  BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedDescription
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "update",

    TRUE
  );
END;

-- Update changes on INSERT to images_graphs table
DROP TRIGGER IF EXISTS images_changes_images_graphs_insert;
CREATE TRIGGER images_changes_images_graphs_insert AFTER INSERT ON images_graphs BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedGraph
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "insert",

    TRUE
  );
END;

-- Update images_updates on UPDATE to images_graphs table
DROP TRIGGER IF EXISTS images_changes_images_graphs_update;
CREATE TRIGGER images_changes_images_graphs_update AFTER UPDATE ON images_graphs WHEN
    old.graph <> new.graph
  BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedGraph
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "update",

    TRUE
  );
END;

-- Update changes on INSERT to images_vulnerabilities table
DROP TRIGGER IF EXISTS images_changes_images_vulnerabilities_insert;
CREATE TRIGGER images_changes_images_vulnerabilities_insert AFTER INSERT ON images_vulnerabilitiesv2 BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedVulnerabilities
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "insert",

    TRUE
  );
END;

-- Update images_updates on UPDATE to images_vulnerabilities table
DROP TRIGGER IF EXISTS images_changes_images_vulnerabilities_update;
CREATE TRIGGER images_changes_images_vulnerabilities_update AFTER UPDATE ON images_vulnerabilitiesv2 WHEN
    old.vulnerabilities <> new.vulnerabilities
  BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedVulnerabilities
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "update",

    TRUE
  );
END;
