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
  PRIMARY KEY (reference, tag)
);

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
